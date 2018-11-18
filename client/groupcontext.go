package client

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"ipfs-share/client/fs/caps"
	"path"
	"sync"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/glog"
	"github.com/pkg/errors"

	com "ipfs-share/client/communication"
	sessclients "ipfs-share/client/communication/sessions/clients"
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
	P2P 			 *com.P2PServer
	Repo             *fs.GroupRepo
	GroupConnection  *com.GroupConnection
	AddressBook		 *ConcurrentCollection
	Network          nw.INetwork
	Ipfs             ipfsapi.IIpfs
	Storage          *fs.Storage
	Transactions     *ConcurrentCollection
	broadcastChannel *ipfsapi.PubSubSubscription
	proposedKey 	 crypto.SymmetricKey
	lock 		     sync.Mutex
}

func (groupCtx *GroupContext) Id() IIdentifier {
	return groupCtx.Group.Id()
}

func NewGroupContext(
	group interfaces.IGroup,
	user interfaces.IUser,
	p2p *com.P2PServer,
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
	p2p *com.P2PServer,
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
	name, members, encIpfsHash, _, _, err := groupCtx.Network.GetGroup(groupCtx.Group.Id().Data().([32]byte))
	if err != nil {
		return errors.Wrapf(err, "could not get group %v", groupCtx.Group.Id().Data())
	}

	if err := groupCtx.Group.Update(name, members, encIpfsHash); err != nil {
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

	session := sessclients.NewCommitChangesGroupSessionClient(
		hash,
		groupCtx.User,
		groupCtx.Group,
		groupCtx.broadcast,
		groupCtx.P2P.SessionClosedChan,
		func(encIpfsHash []byte, approvals []*nw.Approval) {
		groupId := groupCtx.Group.Id().Data().([32]byte)
		tx, err := groupCtx.Network.UpdateGroupIpfsHash(groupId, encIpfsHash, approvals)
		if err != nil {
			glog.Errorf("could not send update group ipfs hash transaction: %s", err)
			return
		}

		groupCtx.Transactions.Append(tx)
	})

	groupCtx.P2P.AddSession(session)
	go session.Run()

	return nil
}

func (groupCtx *GroupContext) broadcast(msg []byte) error {
	return groupCtx.GroupConnection.Broadcast(msg)
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

	var newKey [32]byte
	if _, err := rand.Read(newKey[:]); err != nil {
		return errors.Wrap(err, "could not read from crypto.rand")
	}

	groupCtx.proposedKey = crypto.SymmetricKey{Key: newKey}

	hash, err := groupCtx.Repo.ReEncrypt(groupCtx.proposedKey)
	if err != nil {
		return errors.Wrap(err, "could not re-encrypt group repo")
	}

	session := sessclients.NewCommitChangesGroupSessionClient(
		hash,
		groupCtx.User,
		groupCtx.Group,
		groupCtx.broadcast,
		groupCtx.P2P.SessionClosedChan,
		func(encIpfsHash []byte, approvals []*nw.Approval) {
			groupId := groupCtx.Group.Id().Data().([32]byte)
			tx, err := groupCtx.Network.UpdateGroupIpfsHash(groupId, encIpfsHash, approvals)
			if err != nil {
				glog.Errorf("could not send update group ipfs hash transaction: %s", err)
				return
			}

			groupCtx.Transactions.Append(tx)
		})

	groupCtx.P2P.AddSession(session)
	go session.Run()

	return nil
}

func (groupCtx *GroupContext) ReEncrpyt() error {
	hash, err := groupCtx.Repo.ReEncrypt(groupCtx.Group.Boxer())
	if err != nil {
		return errors.Wrap(err, "could not re-encrypt group repo")
	}

	session := sessclients.NewCommitChangesGroupSessionClient(
		hash,
		groupCtx.User,
		groupCtx.Group,
		groupCtx.broadcast,
		groupCtx.P2P.SessionClosedChan,
		func(encIpfsHash []byte, approvals []*nw.Approval) {
			groupId := groupCtx.Group.Id().Data().([32]byte)
			tx, err := groupCtx.Network.UpdateGroupIpfsHash(groupId, encIpfsHash, approvals)
			if err != nil {
				glog.Errorf("could not send update group ipfs hash transaction: %s", err)
				return
			}

			groupCtx.Transactions.Append(tx)
		})

	groupCtx.P2P.AddSession(session)
	go session.Run()

	return nil
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