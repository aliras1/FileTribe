package client

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"ipfs-share/crypto"
	"ipfs-share/ipfs"
	"log"
)

type Synchronizer struct {
	userSigner *crypto.SigningKeyPair
	groupCtx   *GroupContext
	username   string

	channelPubSub chan ipfs.PubsubMessage
}

func NewSynchronizer(username string, userSigner *crypto.SigningKeyPair, groupCtx *GroupContext) *Synchronizer {
	var synch Synchronizer
	synch.userSigner = userSigner
	synch.groupCtx = groupCtx
	synch.username = username
	synch.channelPubSub = make(chan ipfs.PubsubMessage)

	go synch.groupCtx.IPFS.PubsubSubscribe(synch.groupCtx.Group.GroupName, synch.channelPubSub)
	go synch.MessageProcessor()

	return &synch
}

func (s *Synchronizer) MessageProcessor() {
	fmt.Println("forking...")
	for psMsg := range s.channelPubSub {
		groupMsgBytes, ok := psMsg.Verify(s.groupCtx.Group.Signer.PublicKey)
		if !ok {
			log.Println("invalid group message")
			continue
		}
		var groupMsg GroupMessage
		if err := json.Unmarshal(groupMsgBytes, &groupMsg); err != nil {
			log.Println(err)
			continue
		}
		switch groupMsg.Type {
		case "proposal":
			proposal, err := s.verifyProposal(groupMsg.Data)
			if err != nil {
				log.Println(err)
				continue
			}
			// TODO dont need signature
			err = s.validateProposal(proposal)
			if err != nil {
				log.Println(err)
				continue
			}
			s.approveProposal(s.username, proposal)
		}
	}
}

func (s *Synchronizer) approveProposal(username string, proposal *Proposal) error {
	channelName := s.groupCtx.Group.GroupName + proposal.From
	proposalBytes, err := json.Marshal(proposal)
	if err != nil {
		return err
	}
	hash := sha256.Sum256(proposalBytes)
	approval := Approval{username, hash}
	approvalBytes, err := json.Marshal(&approval)
	if err != nil {
		return err
	}
	userSignedApproval := s.userSigner.SecretKey.Sign(nil, approvalBytes)
	groupSignedApproval := s.groupCtx.Group.Signer.SecretKey.Sign(nil, userSignedApproval)
	return s.groupCtx.IPFS.PubsubPublish(channelName, groupSignedApproval)
}

// verify if the proposal really comes from the given user
func (s *Synchronizer) verifyProposal(data []byte) (*Proposal, error) {
	var proposal Proposal
	err := json.Unmarshal(data[64:], &proposal) // first 64 bytes are the signature
	if err != nil {
		return nil, err
	}
	verifyKey, err := s.groupCtx.Network.GetUserSigningKey(proposal.From)
	if err != nil {
		return nil, err
	}
	_, ok := verifyKey.Open(nil, data)
	if !ok {
		return nil, errors.New("invalid proposal message")
	}
	return &proposal, nil
}

// validate the content of the proposal
func (s *Synchronizer) validateProposal(proposal *Proposal) error {
	switch proposal.CMD {
	case "invite":
		if len(proposal.Args) < 1 {
			return errors.New("invalid #Args in invite proposal")
		}
		newMember := proposal.Args[0]
		prevHash := hashOfMembers(s.groupCtx.Members)
		otherPrevHash := proposal.PrevHash
		if !bytes.Equal(prevHash[:], otherPrevHash[:]) {
			return errors.New("prev hashes do not match")
		}
		newMembers := append(s.groupCtx.Members, newMember)
		newHash := hashOfMembers(newMembers)
		otherNewHash := proposal.NewHash
		if !bytes.Equal(newHash[:], otherNewHash[:]) {
			return errors.New("new hashes do not match")
		}
		// TODO check if user has appropriate rights
		return nil

	default:
		return errors.New("invalid cmd")
	}
}
