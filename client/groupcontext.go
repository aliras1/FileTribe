package client

import (
	"crypto/sha256"
	"encoding/json"

	"fmt"
	fs "ipfs-share/client/filestorage"
	"ipfs-share/crypto"
	"ipfs-share/ipfs"
	nw "ipfs-share/network"
	"log"
	"time"
)

type GroupContext struct {
	Group   *Group
	Repo    []*fs.File
	Members []string
	Network *nw.Network
	IPFS    *ipfs.IPFS

	groupSynch *Synchronizer
}

func hashOfMembers(members []string) [32]byte {
	var data []byte
	for _, member := range members {
		data = append(data, []byte(member)...)
	}
	return sha256.Sum256(data)
}

func (gc *GroupContext) Invite(username, newMember string, userSignKey, groupSignKey *crypto.SecretSigningKey) error {
	channelName := gc.Group.GroupName + username
	fmt.Println("proposal channel name: " + channelName)
	channel := make(chan ipfs.PubsubMessage)
	go gc.IPFS.PubsubSubscribe(channelName, channel)

	prevHash := hashOfMembers(gc.Members)
	newMembers := append(gc.Members, newMember)
	newHash := hashOfMembers(newMembers)

	proposalMsg := Proposal{username, "invite", []string{newMember}, prevHash, newHash}
	proposalMsgBytes, err := json.Marshal(proposalMsg)
	if err != nil {
		return err
	}
	signedProposalMsg := userSignKey.Sign(nil, proposalMsgBytes)
	groupMsg := GroupMessage{"proposal", signedProposalMsg}
	groupMsgBytes, err := json.Marshal(groupMsg)
	if err != nil {
		return err
	}
	signedGroupMsg := groupSignKey.Sign(nil, groupMsgBytes)
	gc.sendToAll(signedGroupMsg)

	// ------------ verifying --------------
	commitMsg := CommitMsg{proposalMsg, []SignedBy{}}
	hash := sha256.Sum256(proposalMsgBytes)
	loop := true
	for loop {
		select {
		case psm := <-channel:
			signedBy, err := ValidateApproval(&psm, hash, gc.Group.Signer.PublicKey, gc.Network)
			if err != nil {
				log.Println(err)
				continue
			}
			commitMsg.SignedBy = append(commitMsg.SignedBy, signedBy)
			if len(commitMsg.SignedBy) >= len(gc.Members)/2 {
				fmt.Println(commitMsg)
				//commit(commitMsg)
				loop = false
			}

		case <-time.After(5 * time.Second):
			fmt.Println("timeout")
			loop = false
		}
	}

	return nil
}

func (gc *GroupContext) sendToAll(data []byte) error {
	return gc.IPFS.PubsubPublish(gc.Group.GroupName, data)
}
