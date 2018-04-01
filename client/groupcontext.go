package client

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	fs "ipfs-share/client/filestorage"
	"ipfs-share/crypto"
	"ipfs-share/ipfs"
	nw "ipfs-share/network"
)

type MemberList []string

func (ml *MemberList) Save(groupName string, boxer crypto.SymmetricKey, storage *fs.Storage, ipfs *ipfs.IPFS) error {
	// store only in public, as on sign in groups are
	// built up from there
	memberBytes, err := json.Marshal(ml)
	if err != nil {
		return err
	}
	return storage.SaveGroupData(groupName, "members.json", boxer, memberBytes)
}

func NewMemberListFromFile(filePath string, key crypto.SymmetricKey) (*MemberList, error) {
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
	return &memberList, nil
}

type GroupContext struct {
	Group   *Group
	Repo    []*fs.File
	Members MemberList
	Network *nw.Network
	IPFS    *ipfs.IPFS
	Storage *fs.Storage
}

func (gc *GroupContext) Invite(username, newMember string, boxer *crypto.BoxingKeyPair, signer *crypto.SecretSigningKey) error {
	if len(gc.Members) == 1 {
		cmd := InviteCMD{username, boxer, newMember, gc}
		cmd.Execute(gc.Network)
		return nil
	}

	prevHash := hashOfMembers(gc.Members)
	newMembers := append(gc.Members, newMember)
	newHash := hashOfMembers(newMembers)

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
			if len(commitMsg.SignedBy) >= len(gc.Members)/2 {
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
				gc.sendToAll(groupMsgBytes)
				// TODO execute cmd

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
	pml, err := NewMemberListFromFile(filePath, gc.Group.Boxer)
	if err != nil {
		return err
	}
	gc.Members = *pml
	return nil
}

func hashOfMembers(members []string) [32]byte {
	var data []byte
	for _, member := range members {
		data = append(data, []byte(member)...)
	}
	return sha256.Sum256(data)
}
