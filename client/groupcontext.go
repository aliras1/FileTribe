package client

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	fs "ipfs-share/client/filestorage"
	"ipfs-share/crypto"
	"ipfs-share/ipfs"
	nw "ipfs-share/network"
)

type Signedby struct {
	Username  string `json:"username"`
	Signature []byte `json:"signature"`
}

type Transaction struct {
	PrevState []byte     `json:"prev_state"`
	State     []byte     `json:"state"`
	Operation string     `json:"operation"`
	SignedBy  []Signedby `json:"signed_by"`
}

type Member struct {
	Name      string                  `json:"name"`
	VerifyKey crypto.PublicSigningKey `json:"-"`
}

type MemberList struct {
	List []Member
}

func NewMemberList() *MemberList {
	return &MemberList{[]Member{}}
}

func (ml *MemberList) Length() int {
	return len(ml.List)
}

func (ml *MemberList) Hash() [32]byte {
	var data []byte
	for _, member := range ml.List {
		data = append(data, []byte(member.Name)...)
	}
	return sha256.Sum256(data)
}

func (ml *MemberList) Append(user string, network *nw.Network) *MemberList {
	verifyKey, err := network.GetUserSigningKey(user)
	if err != nil {
		log.Println(err)
		return ml
	}
	newList := make([]Member, len(ml.List))
	copy(newList, ml.List)
	newList = append(newList, Member{user, verifyKey})
	return &MemberList{newList}
}

func (ml *MemberList) Get(user string) (Member, bool) {
	for i := 0; i < ml.Length(); i++ {
		if strings.Compare(ml.List[i].Name, user) == 0 {
			return ml.List[i], true
		}
	}
	return Member{}, false
}

type GroupContext struct {
	User    *User
	Group   *Group
	Repo    []*fs.FilePTP
	Members *MemberList
	Network *nw.Network
	IPFS    *ipfs.IPFS
	Storage *fs.Storage
}

func NewGroupContext(group *Group, user *User, network *nw.Network, ipfs *ipfs.IPFS, storage *fs.Storage) (*GroupContext, error) {
	members := NewMemberList()
	memberStrings, err := network.GetGroupMembers(group.GroupName)
	if err != nil {
		return nil, fmt.Errorf("could not get group members of '%s': NewGroupContext: %s", group.GroupName, err)
	}
	for _, member := range memberStrings {
		members = members.Append(member, network)
	}
	return &GroupContext{
		User:    user,
		Group:   group,
		Repo:    nil,
		Members: members,
		Network: network,
		IPFS:    ipfs,
		Storage: storage,
	}, nil
}

func NewGroupContextFromCAP(cap *fs.GroupAccessCAP, user *User, network *nw.Network, ipfs *ipfs.IPFS, storage *fs.Storage) (*GroupContext, error) {
	return nil, fmt.Errorf("not implemented: NewGroupContextFromCAP")
}

func (gc *GroupContext) State() []byte {
	state := gc.Members.Hash()
	return state[:]
}

// By operations (e.g. Invite()) a given number of valid approvals
// is needed to be able to commit the current operation. This func
// collects these approvals and upon receiving enough approvals it
// commits the operation
func (gc *GroupContext) collectApprovals(username string, proposal Proposal) {
	channelName := gc.Group.GroupName + username
	channel := make(chan ipfs.PubsubMessage)
	go gc.IPFS.PubsubSubscribe(channelName, channel)

	commitMsg := CommitMsg{proposal, []SignedBy{}}
	proposalMsgBytes, err := json.Marshal(proposal)
	if err != nil {
		log.Println(err)
		return
	}
	hash := sha256.Sum256(proposalMsgBytes)
	for {
		select {
		case psm := <-channel:
			signedBy, err := ValidateApproval(&psm, hash, gc.Group.Boxer, gc.Network)
			if err != nil {
				log.Println(err)
				continue
			}
			commitMsg.SignedBy = append(commitMsg.SignedBy, signedBy)
			if len(commitMsg.SignedBy) >= gc.Members.Length()/2 {
				commitMsgBytes, err := json.Marshal(commitMsg)
				if err != nil {
					log.Println(err)
					return
				}
				groupMsg := GroupMessage{"commit", commitMsgBytes}
				groupMsgBytes, err := json.Marshal(groupMsg)
				if err != nil {
					log.Println(err)
					return
				}
				if err := gc.sendToAll(groupMsgBytes); err != nil {
					log.Println(err)
					return
				}
				// current user also gets the commit message
				// which will be executed in the Synchronizer
				return
			}

		case <-time.After(5 * time.Second):
			fmt.Println("timeout")
			return
		}
	}
}

func (gc *GroupContext) Invite(username, newMember string) error {
	/*if gc.Members.Length() == 1 {
		cmd := InviteCMD{username, newMember}
		return cmd.Execute(gc)
	}*/

	prevHash := gc.Members.Hash()
	newMembers := gc.Members.Append(newMember, gc.Network)
	newHash := newMembers.Hash()

	/*proposalMsg := Proposal{username, "invite", []string{newMember}, prevHash, newHash}
	go gc.collectApprovals(username, proposalMsg)

	proposalMsgBytes, err := json.Marshal(proposalMsg)
	if err != nil {
		return err
	}
	signedProposalMsg := signer.Sign(nil, proposalMsgBytes)
	groupMsg := GroupMessage{"proposal", signedProposalMsg}
	groupMsgBytes, err := json.Marshal(groupMsg)
	if err != nil {
		return err
	}
	gc.sendToAll(groupMsgBytes)*/
	var rawTransaction []byte
	rawTransaction = append(rawTransaction, prevHash[:]...)
	rawTransaction = append(rawTransaction, newHash[:]...)
	rawTransaction = append(rawTransaction, []byte(newMember)...)
	signature := gc.User.Signer.SecretKey.Sign(nil, rawTransaction)[:64]

	transaction := Transaction{
		PrevState: prevHash[:],
		State:     newHash[:],
		Operation: newMember,
		SignedBy: []Signedby{
			Signedby{
				Username:  username,
				Signature: signature,
			},
		},
	}

	transactionBytes, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("could not marshal transaction data: GroupContext.Invite: %s", err)
	}

	if err := gc.Network.GroupInvite(gc.Group.GroupName, transactionBytes); err != nil {
		return fmt.Errorf("could not execute invite through th enetwork: %s", err)
	}

	return nil
}

func (gc *GroupContext) Save() error {
	if err := gc.Group.Save(gc.Storage); err != nil {
		return err
	}
	// should take out publish public dir from here, because it
	// publishes too often by signing in
	return gc.Storage.PublishPublicDir(gc.IPFS)
	//return nil
}

// Sends pubsub messages to all members of the group
func (gc *GroupContext) sendToAll(data []byte) error {
	encGroupMsg := gc.Group.Boxer.BoxSeal(data)
	return gc.IPFS.PubsubPublish(gc.Group.GroupName, encGroupMsg)
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
