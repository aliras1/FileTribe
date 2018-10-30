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
	"encoding/base64"
)


type GroupContext struct {
	User             *User
	Group            *Group
	P2P *P2PServer
	Repo             *GroupRepo
	GroupConnection  *GroupConnection
	AddressBook *ConcurrentCollection
	Network          nw.INetwork
	Ipfs             ipfsapi.IIpfs
	Storage          *Storage
	broadcastChannel *ipfsapi.PubSubSubscription
	lock sync.Mutex
}

func (groupCtx *GroupContext) Id() IIdentifier {
	return groupCtx.Group.Id
}

func NewGroupContext(
	group *Group,
	user *User,
	p2p *P2PServer,
	addressBook *ConcurrentCollection,
	network nw.INetwork,
	ipfs ipfsapi.IIpfs,
	storage *Storage,
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
	}

	repo, err := NewGroupRepo(groupContext)
	if err != nil {
		return nil, errors.Wrap(err, "could not create group repo")
	}

	groupContext.Repo = repo
	groupContext.GroupConnection = NewGroupConnection(groupContext)

	return groupContext, nil
}

func NewGroupContextFromCAP(
	cap *GroupAccessCap,
	user *User,
	p2p *P2PServer,
	addressBook *ConcurrentCollection,
	network nw.INetwork,
	ipfs ipfsapi.IIpfs,
	storage *Storage,
) (*GroupContext, error) {
	group := &Group{
		Id:  NewBytesId(cap.GroupId),
		Boxer: cap.Boxer,
	}

	gc, err := NewGroupContext(group, user, p2p, addressBook, network, ipfs, storage)
	if err != nil {
		return nil, fmt.Errorf("could not create group context: NewGroupContextFromCAP: %s", err)
	}

	return gc, nil
}

func (groupCtx *GroupContext) Update() error {
	groupCtx.lock.Lock()
	defer groupCtx.lock.Unlock()

	name, members, encIpfsPathBase64, err := groupCtx.Network.GetGroup(groupCtx.Group.Id.Data().([32]byte))
	if err != nil {
		return errors.Wrapf(err, "could not get group %v", groupCtx.Group.Id.Data())
	}

	encIpfsPath, err := base64.URLEncoding.DecodeString(encIpfsPathBase64)
	if err != nil {
		return errors.Wrap(err, "could not base64 decode encrypted ipfs path")
	}
	ipfsPathBytes, ok := groupCtx.Group.Boxer.BoxOpen(encIpfsPath)
	if !ok {
		return errors.New("could not decrypt ipfs path")
	}
	ipfsPath := string(ipfsPathBytes)

	groupCtx.Group.Name = name
	groupCtx.Group.Members = members
	groupCtx.Group.IpfsHash = ipfsPath
	groupCtx.Group.EncryptedIpfsHash = encIpfsPathBase64

	if err := groupCtx.Group.Save(groupCtx.Storage); err != nil {
		return errors.Wrapf(err, "could not save group")
	}

	if err := groupCtx.Repo.Update(ipfsPath); err != nil {
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

	session := NewCommitChangesGroupSessionClient(hash, groupCtx)
	groupCtx.P2P.AddSession(session)
	go session.Run()

	return nil
}

func (groupCtx *GroupContext) Invite(newMember ethcommon.Address, hasInviteRight bool) error {
	fmt.Printf("[*] Inviting user '%s' into group '%s'...\n", newMember, groupCtx.Group.Name)

	if err := groupCtx.Network.InviteUser(groupCtx.Group.Id.Data().([32]byte), newMember, hasInviteRight); err != nil {
		return fmt.Errorf("could not invite user: GroupContext::Invite(): %s", err)
	}

	return nil
}


func (groupCtx *GroupContext) Save() error {
	return fmt.Errorf("not implemented: GroupContext.SaveMetadata")
}


// Loads the locally available group meta data
func (groupCtx *GroupContext) LoadGroupData(data string) error {
	return fmt.Errorf("not implemented GroupContext.LoadGroupData")
}
