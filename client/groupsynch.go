package client

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"ipfs-share/ipfs"
)

type Synchronizer struct {
	groupCtx *GroupContext

	channelPubSub chan ipfs.PubsubMessage
}

func NewSynchronizer(groupCtx *GroupContext) *Synchronizer {
	var synch Synchronizer
	synch.groupCtx = groupCtx
	synch.channelPubSub = make(chan ipfs.PubsubMessage)

	go synch.groupCtx.IPFS.PubsubSubscribe(synch.groupCtx.Group.GroupName, synch.channelPubSub)
	go synch.MessageProcessor()
	go synch.heartBeat()

	return &synch
}

func (s *Synchronizer) MessageProcessor() {
	fmt.Println("synch forking...")
	for pubsubMessage := range s.channelPubSub {
		groupMessageBytes, ok := pubsubMessage.Decrypt(s.groupCtx.Group.Boxer)
		if !ok {
			log.Printf("could not decrypt group message: Synchronizer: MessageProcessor")
			continue
		}

		var groupMessage GroupMessage
		if err := json.Unmarshal(groupMessageBytes, &groupMessage); err != nil {
			log.Printf("could not unmarshal group message: Synchronizer, MessageProzessor: %s", err)
			continue
		}

		switch groupMessage.Type {
		case "HB":
			if err := s.processHeartBeat(groupMessage); err != nil {
				log.Println(err)
				continue
			}
		case "PROPOSAL":
			if err := s.processProposal(groupMessage); err != nil {
				log.Printf("error while processing PROPOSAL: Synchronizer.MessageProcessor: %s", err)
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
	return nil
}

func (s *Synchronizer) processProposal(message GroupMessage) error {
	transaction, err := s.authenticateProposal(message.From, message.Data)
	if err != nil {
		return fmt.Errorf("could not authenticate proposal: Synchronizer.processProposal: %s", err)
	}
	if err := s.validateTransaction(transaction); err != nil {
		return fmt.Errorf("error while validating transaction: Synchronizer.processProposal: %s", err)
	}
	if err := s.approveTransaction(message.From, transaction); err != nil {
		return fmt.Errorf("could not approve proposal: Synchronizer.processProposal: %s", err)
	}
	return nil
}

func (s *Synchronizer) heartBeat() {
	for {
		var randomBytes [32]byte
		rand.Read(randomBytes[:])
		signedRand := s.groupCtx.User.Signer.SecretKey.Sign(nil, randomBytes[:])
		heartBeat := HeartBeat{s.groupCtx.User.Username, signedRand}
		hbBytes, err := json.Marshal(heartBeat)
		if err != nil {
			log.Println(err)
			continue
		}
		msg := GroupMessage{"HB", s.groupCtx.User.Username, hbBytes}
		msgBytes, err := json.Marshal(msg)
		if err != nil {
			log.Println(err)
			continue
		}
		s.groupCtx.sendToAll(msgBytes)
		time.Sleep(1 * time.Second)
	}
}

func (s *Synchronizer) approveTransaction(username string, transaction *Transaction) error {
	channelName := s.groupCtx.Group.GroupName + username
	signature := s.groupCtx.User.SignTransaction(transaction)
	approval := Approval{
		From:      username,
		Signature: signature,
	}
	approvalBytes, err := json.Marshal(&approval)
	if err != nil {
		return fmt.Errorf("could not marshal approval: Synchronizer.approveTransaction: %s", err)
	}
	groupEncApproval := s.groupCtx.Group.Boxer.BoxSeal(approvalBytes)
	if err := s.groupCtx.IPFS.PubsubPublish(channelName, groupEncApproval); err != nil {
		return fmt.Errorf("could not send approval: Synchronizer.approveTransaction: %s", err)
	}
	return nil
}

// authenticate if the proposal really comes from the given user
func (s *Synchronizer) authenticateProposal(author string, data []byte) (*Transaction, error) {
	verifyKey, err := s.groupCtx.Network.GetUserSigningKey(author)
	if err != nil {
		return nil, fmt.Errorf("could not get user verify key: Synchronizer.authenticateProposal: %s", err)
	}
	transactionBytes, ok := verifyKey.Open(nil, data)
	if !ok {
		return nil, fmt.Errorf("invalid proposal message from user '%s': Synchronizer.authenticateProposal", author)
	}
	var proposal Transaction
	if err := json.Unmarshal(transactionBytes, &proposal); err != nil {
		return nil, fmt.Errorf("could not unmarshal proposal: Synchronizer.authenticateProposal: %s", err)
	}
	return &proposal, nil
}

// validate the content of a transaction
func (s *Synchronizer) validateTransaction(transaction *Transaction) error {
	state, err := s.groupCtx.GetState()
	if err != nil {
		return fmt.Errorf("could not get state of group '%s': Synchronizer.validateTransaction: %s", s.groupCtx.Group.GroupName, err)
	}
	if !bytes.Equal(state, transaction.PrevState) {
		return fmt.Errorf("invlaid prev state in transaction proposal: Synchronizer.validateTransaction")
	}

	args := strings.Split(transaction.Operation.Data, " ")

	switch transaction.Operation.Type {
	case "INVITE":
		if len(args) < 2 {
			return fmt.Errorf("invalid #Args in invite transaction: Synchronizer.validateTransaction")
		}
		// inviter is args[0]. it should be checked, if he has rights to invite
		newMember := args[1]

		newMembers := s.groupCtx.Members.Append(newMember, s.groupCtx.Network)
		newState := newMembers.Hash()
		if !bytes.Equal(newState[:], transaction.State) {
			return fmt.Errorf("invalid new state in transaction proposal: Synchronizer.validateTransaction")
		}
		return nil

	default:
		return fmt.Errorf("invalid operation type: Synchronizer.vlaidateTransaction")
	}
}

// By operations (e.g. Invite()) a given number of valid approvals
// is needed to be able to commit the current operation. This func
// collects these approvals and upon receiving enough approvals it
// commits the operation
func (synch *Synchronizer) CollectApprovals(transaction *Transaction) {
	channelName := synch.groupCtx.Group.GroupName + synch.groupCtx.User.Username
	channel := make(chan ipfs.PubsubMessage)
	go synch.groupCtx.IPFS.PubsubSubscribe(channelName, channel)

	for {
		if len(transaction.SignedBy) > synch.groupCtx.Members.Length()/2 {
			new_member := strings.Split(transaction.Operation.Data, " ")[1]
			if err := synch.groupCtx.Storage.CreateGroupAccessCAPForUser(
				new_member,
				synch.groupCtx.Group.GroupName,
				synch.groupCtx.Group.Boxer,
				&synch.groupCtx.User.Boxer,
				synch.groupCtx.Network,
			); err != nil {
				log.Printf("could not create ga cap for new member: Synchronizer.CollectApprovals: %s", err)
				return
			}
			if err := synch.groupCtx.Storage.PublishPublicDir(synch.groupCtx.IPFS); err != nil {
				log.Printf("could not publish public dir: Synchronizer.CollectApprovals: %s", err)
				return
			}

			transactionBytes, err := json.Marshal(transaction)
			if err != nil {
				log.Printf("could not marshal transaction: Synchronizer.CollectApprovals: %s", err)
				return
			}
			if err := synch.groupCtx.Network.GroupInvite(synch.groupCtx.Group.GroupName, transactionBytes); err != nil {
				log.Printf("could not call invite transaction: Synchronizer.CollectApprovals: %s", err)
				return
			}
			return
		}

		select {
		case pubsubMessage := <-channel:
			signedBy, err := ValidateApproval(&pubsubMessage, synch.groupCtx.Group.Boxer, synch.groupCtx.Network)
			if err != nil {
				log.Printf("could not validate approval: Synchronizer.CollectApprovals: %s", err)
				continue
			}
			transaction.SignedBy = append(transaction.SignedBy, signedBy)
		case <-time.After(5 * time.Second):
			log.Printf("timeout: Synchronizer.CollectApprovals")
			return
		}
	}
}
