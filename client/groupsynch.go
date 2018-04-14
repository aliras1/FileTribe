package client

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"ipfs-share/crypto"
	"ipfs-share/ipfs"
)

type Synchronizer struct {
	userSigner *crypto.SigningKeyPair
	userBoxer  *crypto.BoxingKeyPair
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
	go synch.heartBeat()

	return &synch
}

func (s *Synchronizer) MessageProcessor() {
	fmt.Println("synch forking...")
	for psMsg := range s.channelPubSub {
		groupMsgBytes, ok := psMsg.Decrypt(s.groupCtx.Group.Boxer)
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
		case "HB":
			if err := s.processHeartBeat(groupMsg); err != nil {
				log.Println(err)
				continue
			}
		case "proposal":
			if err := s.processProposal(groupMsg); err != nil {
				log.Println(err)
				continue
			}
		case "commit":
			if err := s.processCommitMsg(groupMsg); err != nil {
				log.Println(err)
				continue
			}
		}
	}
}

func (s *Synchronizer) processHeartBeat(message GroupMessage) error {
	var heartBeat HeartBeat
	err := json.Unmarshal(message.Data, &heartBeat)
	if err != nil {
		return err
	}
	member, in := s.groupCtx.Members.Get(heartBeat.From)
	if !in {
		return errors.New("heart beat from a non member user")
	}
	_, ok := member.VerifyKey.Open(nil, heartBeat.Rand)
	if !ok {
		return errors.New("invalid heart beat message")
	}
	s.groupCtx.ActiveMembers.Set(member)
	return nil
}

func (s *Synchronizer) processProposal(message GroupMessage) error {
	proposal, err := s.verifyProposal(message.Data)
	if err != nil {
		return err
	}
	// TODO dont need signature
	err = s.validateProposal(proposal)
	if err != nil {
		return err
	}
	return s.approveProposal(s.username, proposal)
}

func (s *Synchronizer) processCommitMsg(message GroupMessage) error {
	var commitMsg CommitMsg
	err := json.Unmarshal(message.Data, &commitMsg)
	if err != nil {
		return err
	}
	fmt.Println("commit msg from : " + commitMsg.Proposal.From)
	if len(commitMsg.SignedBy) < s.groupCtx.Members.Length()/2 {
		return errors.New("not enough approvals")
	}
	proposalBytes, err := json.Marshal(commitMsg.Proposal)
	if err != nil {
		return err
	}
	hash := sha256.Sum256(proposalBytes)
	numValidApprovals := 0
	for _, sign := range commitMsg.SignedBy {
		approval := Approval{sign.User, hash}
		approvalBytes, err := json.Marshal(approval)
		if err != nil {
			log.Println(err)
			continue
		}
		signedApproval := sign.Signature[:]
		signedApproval = append(signedApproval, approvalBytes...)
		verifyKey, err := s.groupCtx.Network.GetUserSigningKey(sign.User)
		if err != nil {
			log.Println(err)
			continue
		}
		_, ok := verifyKey.Open(nil, signedApproval)
		if !ok {
			log.Println("invalid approval in commit")
			continue
		}
		numValidApprovals += 1
	}
	if numValidApprovals < s.groupCtx.Members.Length()/2 {
		return errors.New("not enough valid approvals")
	}
	fmt.Println(numValidApprovals)
	cmd := CMDFromProposal(commitMsg.Proposal)
	return cmd.Execute(s.groupCtx)
}

func (s *Synchronizer) heartBeat() {
	for {
		var randomBytes [32]byte
		rand.Read(randomBytes[:])
		signedRand := s.userSigner.SecretKey.Sign(nil, randomBytes[:])
		heartBeat := HeartBeat{s.username, signedRand}
		hbBytes, err := json.Marshal(heartBeat)
		if err != nil {
			log.Println(err)
			continue
		}
		msg := GroupMessage{"HB", hbBytes}
		msgBytes, err := json.Marshal(msg)
		if err != nil {
			log.Println(err)
			continue
		}
		s.groupCtx.sendToAll(msgBytes)
		time.Sleep(1 * time.Second)
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
	groupEncApproval := s.groupCtx.Group.Boxer.BoxSeal(userSignedApproval)
	return s.groupCtx.IPFS.PubsubPublish(channelName, groupEncApproval)
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
		prevHash := s.groupCtx.Members.Hash()
		otherPrevHash := proposal.PrevHash
		if !bytes.Equal(prevHash[:], otherPrevHash[:]) {
			return errors.New("prev hashes do not match")
		}
		newMembers := s.groupCtx.Members.Append(newMember, s.groupCtx.Network)
		newHash := newMembers.Hash()
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
