package client

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"ipfs-share/client/communication/common"
	"ipfs-share/client/fs/caps"
	"path"
	"sync"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/glog"
	"github.com/pkg/errors"

	com "ipfs-share/client/communication"
	sesscommon "ipfs-share/client/communication/sessions/common"
	"ipfs-share/client/fs"
	"ipfs-share/client/interfaces"
	. "ipfs-share/collections"
	"ipfs-share/crypto"
	ipfsapi "ipfs-share/ipfs"
	nw "ipfs-share/network"
)

type IGroupFacade interface {
	Id() IIdentifier
	GrantWriteAccess(filePath string, user ethcommon.Address) error
	RevokeWriteAccess(filePath string, user ethcommon.Address) error
	CommitChanges() error
	Invite(user ethcommon.Address, hasInviteRigth bool) error
	Leave() error
	ListFiles() []string
	ListMembers() []ethcommon.Address
}

type GroupContext struct {
	User             interfaces.IUser
	Group            interfaces.IGroup
	P2P 			 *com.P2PManager
	Repo             *fs.GroupRepo
	GroupConnection  *com.GroupConnection
	AddressBook		 *ConcurrentCollection
	Network          nw.INetwork
	Ipfs             ipfsapi.IIpfs
	Storage          *fs.Storage
	Transactions     *ConcurrentCollection
	broadcastChannel *ipfsapi.PubSubSubscription
	proposedKeys 	 map[string]crypto.SymmetricKey
	lock 		     sync.Mutex
}

func (groupCtx *GroupContext) Id() IIdentifier {
	return groupCtx.Group.Id()
}

func NewGroupContext(
	group interfaces.IGroup,
	user interfaces.IUser,
	p2p *com.P2PManager,
	addressBook *ConcurrentCollection,
	network nw.INetwork,
	ipfs ipfsapi.IIpfs,
	storage *fs.Storage,
	transactions *ConcurrentCollection,
) (*GroupContext, error) {

	groupContext := &GroupContext{
		User:            user,
		Group:           group,
		P2P: p2p,
		GroupConnection: nil,
		AddressBook:     addressBook,
		Network:         network,
		Ipfs:            ipfs,
		Storage:         storage,
		Transactions:    transactions,
		proposedKeys:    make(map[string]crypto.SymmetricKey),
	}

	repo, err := fs.NewGroupRepo(group, user.Address(), storage, ipfs)
	if err != nil {
		return nil, errors.Wrap(err, "could not create group repo")
	}

	groupContext.Repo = repo
	groupContext.GroupConnection = com.NewGroupConnection(
		group,
		repo,
		user,
		addressBook,
		onSessionClosed,
		p2p,
		ipfs,
		network)

	return groupContext, nil
}

func onSessionClosed(session sesscommon.ISession) {
	glog.Infof("session %d closed with error: %s", session.Id().Data().(uint32), session.Error())
}

func NewGroupContextFromCAP(
	cap *caps.GroupAccessCap,
	user interfaces.IUser,
	p2p *com.P2PManager,
	addressBook *ConcurrentCollection,
	network nw.INetwork,
	ipfs ipfsapi.IIpfs,
	storage *fs.Storage,
	transactions *ConcurrentCollection,
) (*GroupContext, error) {
	group := NewGroupFromCap(cap)
	gc, err := NewGroupContext(group, user, p2p, addressBook, network, ipfs, storage, transactions)
	if err != nil {
		return nil, fmt.Errorf("could not create group context: NewGroupContextFromCAP: %s", err)
	}

	return gc, nil
}

func (groupCtx *GroupContext) Update() error {
	name, members, encIpfsHash, leader, err := groupCtx.Network.GetGroup(groupCtx.Group.Id().Data().([32]byte))
	if err != nil {
		return errors.Wrapf(err, "could not get group %v", groupCtx.Group.Id().Data())
	}

	if err := groupCtx.Group.Update(name, members, encIpfsHash, leader); err != nil {
		return errors.Wrap(err, "could not update group")
	}

	if err := groupCtx.Save(); err != nil {
		return errors.Wrap(err, "could not save group")
	}

	if err := groupCtx.Repo.Update(groupCtx.Group.IpfsHash()); err != nil {
		return errors.Wrap(err, "could not update group repo")
	}

	return nil
}

func (groupCtx *GroupContext) Leave() error {
	tx, err := groupCtx.Network.LeaveGroup(groupCtx.Group.Id().Data().([32]byte))
	if err != nil {
		return errors.Wrap(err, "could not send leave group tx")
	}

	groupCtx.Transactions.Append(tx)

	return nil
}

func (groupCtx *GroupContext) Stop() {
	groupCtx.GroupConnection.Kill()
}

func (groupCtx *GroupContext) CommitChanges() error {
	hash, err := groupCtx.Repo.CommitChanges(groupCtx.Group.Boxer())
	if err != nil {
		return errors.Wrap(err, "could commit group repo's changes")
	}

	if err := groupCtx.P2P.StartCommitSession(
		hash,
		groupCtx.User,
		groupCtx.Group,
		groupCtx.broadcast,
		groupCtx.OnCommitClientSuccess,
	); err != nil {
		return errors.Wrap(err, "could not start session")
	}

	return nil
}


func (groupCtx *GroupContext) Invite(newMember ethcommon.Address, hasInviteRight bool) error {
	fmt.Printf("[*] Inviting user '%s' into group '%s'...\n", newMember, groupCtx.Group.Name)

	tx, err := groupCtx.Network.InviteUser(groupCtx.Group.Id().Data().([32]byte), newMember, hasInviteRight)
	if err != nil {
		return errors.Wrap(err, "could not send invite user tx")
	}

	groupCtx.Transactions.Append(tx)

	return nil
}


func (groupCtx *GroupContext) Save() error {
	cap := &caps.GroupAccessCap{
		Boxer: groupCtx.Group.Boxer(),
		GroupId: groupCtx.Group.Id().Data().([32]byte),
	}

	if err := groupCtx.Storage.SaveGroupAccessCap(cap); err != nil {
		return errors.Wrap(err, "could not save group cap")
	}

	return nil
}


// Loads the locally available group meta data
func (groupCtx *GroupContext) LoadGroupData(data string) error {
	return fmt.Errorf("not implemented GroupContext.LoadGroupData")
}

func (groupCtx *GroupContext) GrantWriteAccess(filePath string, user ethcommon.Address) error {
	if !groupCtx.Group.IsMember(user) {
		return errors.New("can not grant write access to non group members")
	}

	file := groupCtx.Repo.Get(NewStringId(path.Base(filePath)))
	if file == nil {
		tmpFile, err := fs.NewGroupFile(
			filePath,
			[]ethcommon.Address{groupCtx.User.Address()},
			groupCtx.Group.Id().ToString(),
			groupCtx.Storage,)
		if err != nil {
			return errors.Wrap(err, "could not create new group file")
		}
		file = tmpFile
	}

	if err := file.GrantWriteAccess(groupCtx.User.Address(), user); err != nil {
		return errors.Wrap(err, "could not grant write access to user")
	}

	return nil
}

func (groupCtx *GroupContext) RevokeWriteAccess(filePath string, user ethcommon.Address) error {
	if !groupCtx.Group.IsMember(user) {
		return errors.New("can not revoke write access from non group members")
	}

	file := groupCtx.Repo.Get(NewStringId(path.Base(filePath)))
	if file == nil {
		tmpFile, err := fs.NewGroupFile(
			filePath,
			[]ethcommon.Address{groupCtx.User.Address()},
			groupCtx.Group.Id().ToString(),
			groupCtx.Storage,)
		if err != nil {
			return errors.Wrap(err, "could not create new group file")
		}
		file = tmpFile
	}

	if err := file.RevokeWriteAccess(groupCtx.User.Address(), user); err != nil {
		return errors.Wrap(err, "could not revoke write access to user")
	}

	return nil
}

func (groupCtx *GroupContext) OnKeyDirty() error {
	leader, err := groupCtx.Network.GetGroupLeader(groupCtx.Group.Id().Data().([32]byte))
	if err != nil {
		return errors.Wrap(err, "could not get group leader")
	}

	// if user is not the leader --> return
	if !bytes.Equal(leader.Bytes(), groupCtx.User.Address().Bytes()) {
		return nil
	}

	newBoxer, err := crypto.NewSymmetricKey()
	if err != nil {
		return errors.Wrap(err, "could not create new group key")
	}

	newIpfsHash, err := groupCtx.Repo.ReEncrypt(newBoxer)
	if err != nil {
		return errors.Wrap(err, "could not re-encrypt group repo")
	}

	encNewIpfsHash := newBoxer.BoxSeal([]byte(newIpfsHash))
	encNewIpfsHashBase64 := base64.StdEncoding.EncodeToString(encNewIpfsHash)

	// TODO: use a cache class instead and save state
	groupCtx.proposedKeys[encNewIpfsHashBase64] = newBoxer

	err = groupCtx.P2P.StartChangeGroupKeySession(
		newBoxer,
		encNewIpfsHash,
		groupCtx.User,
		groupCtx.Group,
		groupCtx.p2pBroadcast,
		groupCtx.OnChangeGroupKeyClientSuccess)

	if err != nil {
		return errors.Wrap(err, "could not start new session")
	}

	return nil
}

func (groupCtx *GroupContext) ReEncrpyt() error {
	hash, err := groupCtx.Repo.ReEncrypt(groupCtx.Group.Boxer())
	if err != nil {
		return errors.Wrap(err, "could not re-encrypt group repo")
	}

	if err := groupCtx.P2P.StartCommitSession(
		hash,
		groupCtx.User,
		groupCtx.Group,
		groupCtx.broadcast,
		groupCtx.OnCommitClientSuccess,
	); err != nil {
		return errors.Wrap(err, "could not start new session")
	}

	return nil
}

func (groupCtx *GroupContext) GetKey(encNewIpfsHash []byte) error {
	encNewIpfsHashBase64 := base64.StdEncoding.EncodeToString(encNewIpfsHash)
	newBoxer, ok := groupCtx.proposedKeys[encNewIpfsHashBase64]
	if ok {
		groupCtx.onGetKeySuccess(newBoxer)
	} else {
		for _, member := range groupCtx.Group.Members() {
			if bytes.Equal(member.Bytes(), groupCtx.User.Address().Bytes()) {
				continue
			}

			var contact *common.Contact
			contactInt := groupCtx.AddressBook.Get(NewAddressId(&member))
			if contactInt == nil {
				c, err := groupCtx.Network.GetUser(member)
				if err != nil {
					glog.Warningf("could not get user in Group.GetKey(): %s", err)
					continue
				}

				contact = common.NewContact(c, groupCtx.Ipfs)
				groupCtx.AddressBook.Append(contact)
			} else {
				contact = contactInt.(*common.Contact)
			}

			err := groupCtx.P2P.StartGetGroupKeySession(
				groupCtx.Group.Id().Data().([32]byte),
				contact,
				groupCtx.User,
				groupCtx.Storage,
				func(cap *caps.GroupAccessCap) {
					groupCtx.onGetKeySuccess(cap.Boxer)
				})
			if err != nil {
				glog.Errorf("could not start get group key session: %s", err)
			}
		}
	}

	return nil
}

func (groupCtx *GroupContext) onGetKeySuccess(boxer crypto.SymmetricKey) {
	groupCtx.Group.SetBoxer(boxer)

	if err := groupCtx.Save(); err != nil {
		glog.Errorf("could not save new key: %s", err)
		return
	}

	if err := groupCtx.Update(); err != nil {
		glog.Errorf("could not update group: %s", err)
	}
}

func (groupCtx *GroupContext) ListFiles() []string {
	var fileNames []string
	files := groupCtx.Repo.Files()

	for _, file := range files {
		fileNames = append(fileNames, file.Cap.FileName)
	}

	return fileNames
}

func (groupCtx *GroupContext) ListMembers() []ethcommon.Address {
	return groupCtx.Group.Members()
}

func (groupCtx *GroupContext) OnChangeGroupKeyClientSuccess(args []interface{}, approvals []*nw.Approval) {
	if len(args) < 1 {
		glog.Error("args should be min. of length 1")
	}

	encNewIpfsHash := args[0].([]byte)

	groupId := groupCtx.Group.Id().Data().([32]byte)
	tx, err := groupCtx.Network.ChangeGroupKey(groupId, encNewIpfsHash, approvals)
	if err != nil {
		glog.Errorf("could not send change group key transaction: %s", err)
		return
	}

	groupCtx.Transactions.Append(tx)
}

func (groupCtx *GroupContext) OnCommitClientSuccess(args []interface{}, approvals []*nw.Approval) {
	if len(args) < 1 {
		glog.Error("args should be of length 1")
	}

	encNewIpfsHash := args[0].([]byte)

	groupId := groupCtx.Group.Id().Data().([32]byte)
	tx, err := groupCtx.Network.UpdateGroupIpfsHash(groupId, encNewIpfsHash, approvals)
	if err != nil {
		glog.Errorf("could not send update group ipfs hash transaction: %s", err)
		return
	}

	groupCtx.Transactions.Append(tx)
}

func (groupCtx *GroupContext) broadcast(msg []byte) error {
	return groupCtx.GroupConnection.Broadcast(msg)
}

func (groupCtx *GroupContext) p2pBroadcast(msg []byte) error {
	for _, member := range groupCtx.Group.Members() {
		var contact *common.Contact
		contactInt := groupCtx.AddressBook.Get(NewAddressId(&member))
		if contactInt == nil {
			netContact, err := groupCtx.Network.GetUser(member)
			if err != nil {
				glog.Errorf("could not get user from network: %s", err)
				continue
			}

			contact = common.NewContact(netContact, groupCtx.Ipfs)
			groupCtx.AddressBook.Append(contact)
		} else {
			contact = contactInt.(*common.Contact)
		}

		go func() {
			if err := contact.Send(msg); err != nil {
				glog.Errorf("error while sending p2p message: %s", err)
			}
		}()
	}

	return nil
}