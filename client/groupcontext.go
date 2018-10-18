package client

import (
	"github.com/ethereum/go-ethereum/common"
	"bytes"
	"fmt"
	// "github.com/golang/glog"
	ipfsapi "ipfs-share/ipfs"
	ethcommon "github.com/ethereum/go-ethereum/common"

	"ipfs-share/crypto"
	nw "ipfs-share/networketh"
	"github.com/pkg/errors"
	"github.com/golang/glog"

	. "ipfs-share/collections"
	"sync"
)

type Member struct {
	ID        common.Address         `json:"sessionId"`
	VerifyKey crypto.VerifyKey `json:"-"`
}

type MemberList struct {
	List []Member
}

func NewMemberList() *MemberList {
	return &MemberList{List: []Member{}}
}

func (ml *MemberList) Length() int {
	return len(ml.List)
}

func (ml *MemberList) Bytes() []byte {
	var data []byte
	for _, member := range ml.List {
		data = append(data, member.ID[:]...)
	}
	return data
}

func (ml *MemberList) Append(userID common.Address, network *nw.Network) *MemberList {
	// _, _, verifyKeyBytes, _, err := network.GetUser(userID)
	// if err != nil {
	// 	glog.Errorf("could not get user verify key: MemberList.Append: %s", err)
	// 	return ml
	// }
	// verifyKey := crypto.VerifyKey(verifyKeyBytes[:])
	// newList := make([]Member, len(ml.List))
	// copy(newList, ml.List)
	// newList = append(newList, Member{userID, verifyKey})
	// return &MemberList{List: newList}
	return nil
}

func (ml *MemberList) Get(userID [32]byte) *Member {
	for i := 0; i < ml.Length(); i++ {
		if bytes.Equal(ml.List[i].ID[:], userID[:]) {
			return &ml.List[i]
		}
	}
	return nil
}

type GroupContext struct {
	User             *User
	Group            *Group
	P2P *P2PServer
	Repo             *GroupRepo
	Members          *MemberList
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
	repo := &GroupRepo{
		Files: []*FileGroup{},
	}
	groupContext := &GroupContext{
		User:            user,
		Group:           group,
		P2P: p2p,
		Repo:            repo,
		GroupConnection: nil,
		AddressBook:     addressBook,
		Network:         network,
		Ipfs:            ipfs,
		Storage:         storage,
	}

	if err := groupContext.Update(); err != nil {
		glog.Errorf("could not update group %v", groupContext.Group.Id.Data())
	}

	groupContext.GroupConnection = NewGroupConnection(groupContext)

	return groupContext, nil
}

func NewGroupContextFromCAP(
	cap *GroupAccessCAP,
	user *User,
	p2p *P2PServer,
	addressBook *ConcurrentCollection,
	network nw.INetwork,
	ipfs ipfsapi.IIpfs,
	storage *Storage,
) (*GroupContext, error) {
	group := &Group{
		Id:  NewBytesId(cap.GroupID),
		Boxer: cap.Boxer,
	}

	gc, err := NewGroupContext(group, user, p2p, addressBook, network, ipfs, storage)
	if err != nil {
		return nil, fmt.Errorf("could not create group context: NewGroupContextFromCAP: %s", err)
	}

	return gc, nil

	return nil, fmt.Errorf("not implemented NewGroupContextFromCAP")
}

func (groupCtx *GroupContext) Update() error {
	groupCtx.lock.Lock()
	defer groupCtx.lock.Unlock()

	name, members, ipfsPath, err := groupCtx.Network.GetGroup(groupCtx.Group.Id.Data().([32]byte))
	if err != nil {
		return errors.Wrapf(err, "could not get group %v", groupCtx.Group.Id.Data())
	}

	// TODO: send updated event
	groupCtx.Group.Name = name
	groupCtx.Group.Members = members
	groupCtx.Group.IPFSPath = ipfsPath

	if err := groupCtx.Group.Save(groupCtx.Storage); err != nil {
		return errors.Wrapf(err, "could not save group")
	}

	return nil
}

func (groupCtx *GroupContext) Stop() {
	groupCtx.GroupConnection.Kill()
}

func (groupCtx *GroupContext) AddFile(filePath string) error {
	session := NewAddFileClientGroupSession("newPath", groupCtx)
	groupCtx.P2P.AddSession(session)
	go session.Run()

	return nil
}

func (groupCtx *GroupContext) Invite(newMember ethcommon.Address) error {
	fmt.Printf("[*] Inviting user '%s' into group '%s'...\n", newMember, groupCtx.Group.Name)

	if err := groupCtx.Network.InviteUser(groupCtx.Group.Id.Data().([32]byte), newMember); err != nil {
		return fmt.Errorf("could not invite user: GroupContext::Invite(): %s", err)
	}

	// prevHash := gc.CalculateState(gc.Members, gc.Repo)
	// newMembers := gc.Members.Append(newMember, gc.Network)
	// newHash := gc.CalculateState(newMembers, gc.Repo)

	// operation := NewInviteOperation(gc.User.Name, newMember)
	// transaction := Transaction{
	// 	PrevState: prevHash[:],
	// 	State:     newHash[:],
	// 	Operation: operation.RawOperation(),
	// 	SignedBy:  []SignedBy{},
	// }
	// // fork down the collection of signatures for the operation
	// go gc.GroupConnection.CollectApprovals(&transaction)

	// // send out the proposed transaction to be signed
	// transactionBytes, err := json.Marshal(transaction)
	// if err != nil {
	// 	return fmt.Errorf("could not marshal transaction: GroupContext.Invite: %s", err)
	// }
	// signedTransactionBytes := gc.User.Signer.VerifyKey.Sign(nil, transactionBytes)
	// groupMsg := GroupMessage{
	// 	Type: "PROPOSAL",
	// 	From: gc.User.Name,
	// 	Data: signedTransactionBytes,
	// }
	// groupMsgBytes, err := json.Marshal(groupMsg)
	// if err != nil {
	// 	return fmt.Errorf("could not marshal group message: GroupContext.Invite: %s", err)
	// }
	// if err := gc.SendToAll(groupMsgBytes); err != nil {
	// 	return fmt.Errorf("could not send group message: GroupContext.Invite: %s", err)
	// }
	return nil
}

func (groupCtx *GroupContext) Save() error {
	return fmt.Errorf("not implemented: GroupContext.Save")
}

// Sends pubsub messages to all members of the group
//func (groupCtx *GroupContext) SendToAll(data []byte) error {
//	encGroupMsg := groupCtx.Group.Boxer.BoxSeal(data)
//
//	if err := groupCtx.Ipfs.PubSubPublish(groupCtx.Group.Name, base64.StdEncoding.EncodeToString(encGroupMsg)); err != nil {
//		return fmt.Errorf("could not pubsub publish: GroupContext.SendToAll: %s", err)
//	}
//	return nil
//}

// Pulls from others the given group meta data
func (groupCtx *GroupContext) PullGroupData(data string) error {
	return fmt.Errorf("not implemented: GroupContext.PullGroupData")
}

// Loads the locally available group meta data (stored in the
// data/public/for/GROUP/ directory).
func (groupCtx *GroupContext) LoadGroupData(data string) error {
	return fmt.Errorf("not implemented GroupContext.LoadGroupData")
}
