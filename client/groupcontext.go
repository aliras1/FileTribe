package client

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

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

func (ml *MemberList) Save(groupName string, boxer crypto.SymmetricKey, storage *fs.Storage, ipfs *ipfs.IPFS) error {
	// store only in public, as on sign in groups are
	// built up from there
	memberBytes, err := json.Marshal(ml)
	if err != nil {
		return err
	}
	return storage.SaveGroupData(groupName, "members.json", boxer, memberBytes)
}

func NewMemberListFromFile(filePath string, key crypto.SymmetricKey, network *nw.Network) (*MemberList, error) {
	box, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	memberBytes, ok := key.BoxOpen(box)
	if !ok {
		return nil, errors.New("could not decrypt file " + filePath)
	}
	var memberList MemberList
	err = json.Unmarshal(memberBytes, &memberList)
	if err != nil {
		return nil, err
	}
	for i := 0; i < memberList.Length(); i++ {
		user := memberList.List[i].Name
		verifyKey, err := network.GetUserSigningKey(user)
		if err != nil {
			return nil, err
		}
		memberList.List[i].VerifyKey = verifyKey
	}
	return &memberList, nil
}

type ActiveMember struct {
	Member
	Time time.Time
}

type ActiveMemberList struct {
	List []ActiveMember
}

func (aml *ActiveMemberList) Length() int {
	return len(aml.List)
}

func (aml *ActiveMemberList) Get(user string) (ActiveMember, bool) {
	for i := 0; i < aml.Length(); i++ {
		if strings.Compare(aml.List[i].Name, user) == 0 {
			return aml.List[i], true
		}
	}
	return ActiveMember{}, false
}

func (aml *ActiveMemberList) Set(member Member) {
	for i := 0; i < aml.Length(); i++ {
		if strings.Compare(aml.List[i].Name, member.Name) == 0 {
			aml.List[i].Time = time.Now()
			return
		}
	}
	newActiveMember := ActiveMember{member, time.Now()}
	aml.List = append(aml.List, newActiveMember)
}

func (aml *ActiveMemberList) Refresh() {
	for {
		currentTime := time.Now()
		var newList []ActiveMember
		for i := 0; i < aml.Length(); i++ {
			if currentTime.Sub(aml.List[i].Time) < 2*time.Second {
				newList = append(newList, aml.List[i])
			}
		}
		aml.List = newList
		fmt.Println(aml.List)
		time.Sleep(2 * time.Second)
	}
}

type GroupContext struct {
	User          *User
	Group         *Group
	Repo          []*fs.File
	Members       *MemberList
	ActiveMembers *ActiveMemberList
	Network       *nw.Network
	IPFS          *ipfs.IPFS
	Storage       *fs.Storage
}

func (gc *GroupContext) Invite(username, newMember string, boxer *crypto.BoxingKeyPair, signer *crypto.SecretSigningKey) error {
	if gc.Members.Length() == 1 {
		cmd := InviteCMD{username, newMember}
		return cmd.Execute(gc)
	}

	prevHash := gc.Members.Hash()
	newMembers := gc.Members.Append(newMember, gc.Network)
	newHash := newMembers.Hash()

	proposalMsg := Proposal{username, "invite", []string{newMember}, prevHash, newHash}
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
	gc.sendToAll(groupMsgBytes)
	return nil
}

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

func (gc *GroupContext) Save() error {
	if err := gc.Group.Save(gc.Storage); err != nil {
		return err
	}
	if err := gc.Members.Save(gc.Group.GroupName, gc.Group.Boxer, gc.Storage, gc.IPFS); err != nil {
		return err
	}
	// ... //
	return gc.Storage.PublishPublicDir(gc.IPFS)
}

func (gc *GroupContext) sendToAll(data []byte) error {
	encGroupMsg := gc.Group.Boxer.BoxSeal(data)
	return gc.IPFS.PubsubPublish(gc.Group.GroupName, encGroupMsg)
}

func (gc *GroupContext) PullGroupData(from string) error {
	// TODO some hash agreement
	groupName := gc.Group.GroupName
	filePath, err := gc.Storage.DownloadGroupData(groupName, "members.json", from, gc.IPFS, gc.Network)
	if err != nil {
		return err
	}
	pml, err := NewMemberListFromFile(filePath, gc.Group.Boxer, gc.Network)
	if err != nil {
		return err
	}
	gc.Members = pml
	return nil
}
