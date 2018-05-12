package client

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/golang/glog"

	fs "ipfs-share/client/filestorage"
	"ipfs-share/crypto"
	"ipfs-share/ipfs"
	nw "ipfs-share/network"
)


type Member struct {
	Name      string                  `json:"name"`
	VerifyKey crypto.PublicSigningKey `json:"-"`
}

type MemberList struct {
	List []Member
}

func NewMemberList() *MemberList {
	return &MemberList{ List:[]Member{}}
}

func (ml *MemberList) Length() int {
	return len(ml.List)
}

func (ml *MemberList) Bytes() []byte {
	var data []byte
	for _, member := range ml.List {
		data = append(data, []byte(member.Name)...)
	}
	return data
}

func (ml *MemberList) Append(user string, network *nw.Network) *MemberList {
	verifyKey, err := network.GetUserVerifyKey(user)
	if err != nil {
		glog.Errorf("could not get user verify key: MemberList.Append: %s", err)
		return ml
	}
	newList := make([]Member, len(ml.List))
	copy(newList, ml.List)
	newList = append(newList, Member{user, verifyKey})
	return &MemberList{List: newList}
}

func (ml *MemberList) Get(user string) *Member {
	for i := 0; i < ml.Length(); i++ {
		if strings.Compare(ml.List[i].Name, user) == 0 {
			return &ml.List[i]
		}
	}
	return nil
}

type GroupContext struct {
	User         *User
	Group        *Group
	Repo         *fs.GroupRepo
	Members      *MemberList
	Synchronizer *Synchronizer
	Network      *nw.Network
	IPFS         *ipfs.IPFS
	Storage      *fs.Storage
}

func NewGroupContext(group *Group, user *User, network *nw.Network, ipfs *ipfs.IPFS, storage *fs.Storage) (*GroupContext, error) {
	members := NewMemberList()
	memberStrings, err := network.GetGroupMembers(group.Name)
	if err != nil {
		return nil, fmt.Errorf("could not get group members of '%s': NewGroupContext: %s", group.Name, err)
	}
	for _, member := range memberStrings {
		members = members.Append(member, network)
	}
	repo := &fs.GroupRepo{
		Files: []*fs.FileGroup{},
	}
	groupContext := &GroupContext{
		User:         user,
		Group:        group,
		Repo:         repo,
		Members:      members,
		Synchronizer: nil,
		Network:      network,
		IPFS:         ipfs,
		Storage:      storage,
	}
	groupContext.Synchronizer = NewSynchronizer(groupContext)
	return groupContext, nil
}

func NewGroupContextFromCAP(cap *fs.GroupAccessCAP, user *User, network *nw.Network, ipfs *ipfs.IPFS, storage *fs.Storage) (*GroupContext, error) {
	group := &Group{
		Name:  cap.GroupName,
		Boxer: cap.Boxer,
	}
	gc, err := NewGroupContext(group, user, network, ipfs, storage)
	if err != nil {
		return nil, fmt.Errorf("could not create group context: NewGroupContextFromCAP: %s", err)
	}
	return gc, nil
}

func (gc *GroupContext) Stop() {
	gc.Synchronizer.Kill()
}

func (gc *GroupContext) CalculateState(members *MemberList, repo *fs.GroupRepo) []byte {
	digest := append(members.Bytes(), repo.Bytes()...)
	hash := sha256.Sum256(digest)
	return hash[:]
}

func (gc *GroupContext) GetState() ([]byte, error) {
	state, err := gc.Network.GetGroupState(gc.Group.Name)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve state from network: GroupContext.GetState: %s", err)
	}
	return state, nil
}

func (gc *GroupContext) AddAndShareFile(filePath string) error {
	file, err := fs.NewSharedFileGroup(filePath, gc.Group.Name, gc.Group.Boxer, gc.Storage, gc.IPFS)
	if err != nil {
		return fmt.Errorf("could not create new shared file group: GroupContext.AddAndShareFile: %s", err)
	}

	newRepo := &fs.GroupRepo{
		Files: append(gc.Repo.Files, file),
	}
	newState := gc.CalculateState(gc.Members, newRepo)
	operation := NewShareFileOperation(gc.User.Name, file.Name, file.IPFSHash)
	transaction := &Transaction{
		PrevState: gc.CalculateState(gc.Members, gc.Repo),
		State: newState,
		Operation: operation.RawOperation(),
		SignedBy: []SignedBy{},
	}
	signature := gc.User.SignTransaction(transaction)
	transaction.SignedBy = []SignedBy{
		{
			Username:  gc.User.Name,
			Signature: signature,
		},
	}
	transactionJSON, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("could not marshal transaction: GroupContext.AddAndShareFile: %s", err)
	}
	if err := gc.Network.GroupShare(gc.Group.Name, transactionJSON); err != nil {
		return fmt.Errorf("error while network call: GroupContext.AddANdShareFile: %s", err)
	}

	fmt.Printf("[*] file '%s' shared with group '%s'\n", filePath, gc.Group.Name)

	return nil
}

func (gc *GroupContext) Invite(newMember string) error {
	fmt.Printf("[*] Inviting user '%s' into group '%s'...\n", newMember, gc.Group.Name)

	prevHash := gc.CalculateState(gc.Members, gc.Repo)
	newMembers := gc.Members.Append(newMember, gc.Network)
	newHash := gc.CalculateState(newMembers, gc.Repo)

	operation := NewInviteOperation(gc.User.Name, newMember)
	transaction := Transaction{
		PrevState: prevHash[:],
		State:     newHash[:],
		Operation: operation.RawOperation(),
		SignedBy:  []SignedBy{},
	}
	// fork down the collection of signatures for the operation
	go gc.Synchronizer.CollectApprovals(&transaction)

	// send out the proposed transaction to be signed
	transactionBytes, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("could not marshal transaction: GroupContext.Invite: %s", err)
	}
	signedTransactionBytes := gc.User.Signer.SecretKey.Sign(nil, transactionBytes)
	groupMsg := GroupMessage{
		Type: "PROPOSAL",
		From: gc.User.Name,
		Data: signedTransactionBytes,
	}
	groupMsgBytes, err := json.Marshal(groupMsg)
	if err != nil {
		return fmt.Errorf("could not marshal group message: GroupContext.Invite: %s", err)
	}
	if err := gc.sendToAll(groupMsgBytes); err != nil {
		return fmt.Errorf("could not send group message: GroupContext.Invite: %s", err)
	}
	return nil
}

func (gc *GroupContext) Save() error {
	return fmt.Errorf("not implemented: GroupContext.Save")
}

// Sends pubsub messages to all members of the group
func (gc *GroupContext) sendToAll(data []byte) error {
	encGroupMsg := gc.Group.Boxer.BoxSeal(data)
	if err := gc.IPFS.PubsubPublish(gc.Group.Name, encGroupMsg); err != nil {
		return fmt.Errorf("could not pubsub publish: GroupContext.sendToAll: %s", err)
	}
	return nil
}

// Pulls from others the given group meta data
func (gc *GroupContext) PullGroupData(data string) error {
	return fmt.Errorf("not implemented: GroupContext.PullGroupData")
}

// Loads the locally available group meta data (stored in the
// data/public/for/GROUP/ directory).
func (gc *GroupContext) LoadGroupData(data string) error {
	return fmt.Errorf("not implemented GroupContext.LoadGroupData")
}
