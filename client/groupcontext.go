package client

import (
	"fmt"
	// "github.com/golang/glog"
	ipfsapi "ipfs-share/ipfs"
	ethcommon "github.com/ethereum/go-ethereum/common"

	nw "ipfs-share/network"
	"github.com/pkg/errors"
	. "ipfs-share/collections"
	"sync"
	"path"
	"github.com/golang/glog"
)


type GroupContext struct {
	User             IUser
	Group            IGroup
	P2P *P2PServer
	Repo             *GroupRepo
	GroupConnection  *GroupConnection
	AddressBook *ConcurrentCollection
	Network          nw.INetwork
	Ipfs             ipfsapi.IIpfs
	Storage          *Storage
	Transactions     *ConcurrentCollection
	broadcastChannel *ipfsapi.PubSubSubscription
	lock sync.Mutex
}

func (groupCtx *GroupContext) Id() IIdentifier {
	return groupCtx.Group.Id()
}

func NewGroupContext(
	group IGroup,
	user IUser,
	p2p *P2PServer,
	addressBook *ConcurrentCollection,
	network nw.INetwork,
	ipfs ipfsapi.IIpfs,
	storage *Storage,
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

	repo, err := NewGroupRepo(group, user.Address(), storage, ipfs)
	if err != nil {
		return nil, errors.Wrap(err, "could not create group repo")
	}

	groupContext.Repo = repo
	groupContext.GroupConnection = NewGroupConnection(groupContext)

	return groupContext, nil
}

func NewGroupContextFromCAP(
	cap *GroupAccessCap,
	user IUser,
	p2p *P2PServer,
	addressBook *ConcurrentCollection,
	network nw.INetwork,
	ipfs ipfsapi.IIpfs,
	storage *Storage,
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
	name, members, encIpfsHash, err := groupCtx.Network.GetGroup(groupCtx.Group.Id().Data().([32]byte))
	if err != nil {
		return errors.Wrapf(err, "could not get group %v", groupCtx.Group.Id().Data())
	}

	if err := groupCtx.Group.Update(name, members, encIpfsHash); err != nil {
		return errors.Wrap(err, "could not update group")
	}

	if err := groupCtx.Group.Save(groupCtx.Storage); err != nil {
		return errors.Wrapf(err, "could not save group")
	}

	if err := groupCtx.Repo.Update(groupCtx.Group.IpfsHash()); err != nil {
		return errors.Wrap(err, "could not update group repo")
	}

	return nil
}

func (groupCtx *GroupContext) Stop() {
	groupCtx.GroupConnection.Kill()
}

func (groupCtx *GroupContext) CommitChanges() error {
	hash, err := groupCtx.Repo.CommitChanges()
	if err != nil {
		return errors.Wrap(err, "could commit group repo's changes")
	}

	session := NewCommitChangesGroupSessionClient(
		hash,
		groupCtx.User,
		groupCtx.Group,
		groupCtx.GroupConnection,
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
	return fmt.Errorf("not implemented: GroupContext.SaveMetadata")
}


// Loads the locally available group meta data
func (groupCtx *GroupContext) LoadGroupData(data string) error {
	return fmt.Errorf("not implemented GroupContext.LoadGroupData")
}

func (groupCtx *GroupContext) GrantWriteAccess(filePath string, user ethcommon.Address) error {
	if !groupCtx.Group.IsMember(user) {
		return errors.New("can not grant write access to non group members")
	}

	var file *File
	fileInt := groupCtx.Repo.files.Get(NewStringId(path.Base(filePath)))
	if fileInt == nil {
		tmpFile, err := NewGroupFile(
			filePath,
			[]ethcommon.Address{groupCtx.User.Address()},
			groupCtx.Group,
			groupCtx.Storage,
			groupCtx.Ipfs)
		if err != nil {
			return errors.Wrap(err, "could not create new group file")
		}
		file = tmpFile
	} else {
		file = fileInt.(*File)
	}

	if err := file.GrantWriteAccess(groupCtx.User.Address(), user); err != nil {
		return errors.Wrap(err, "could not grant write access to user")
	}

	return nil
}

