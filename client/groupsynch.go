package client

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang/glog"

	"ipfs-share/ipfs"
)

type Synchronizer struct {
	groupCtx *GroupContext

	channelPubSub chan ipfs.PubsubMessage
	channelState  chan []byte
	channelStop   chan int
}

func NewSynchronizer(groupCtx *GroupContext) *Synchronizer {
	fmt.Printf("[*] Creating synchronizer...\n")
	var synch Synchronizer
	synch.groupCtx = groupCtx
	synch.channelStop = make(chan int)
	synch.channelPubSub = make(chan ipfs.PubsubMessage)

	go synch.groupCtx.IPFS.PubsubSubscribe(synch.groupCtx.User.Name, synch.groupCtx.Group.Name, synch.channelPubSub)
	go synch.MessageProcessor()
	go synch.StateListener()
	//go synch.heartBeat()

	return &synch
}

func (synch *Synchronizer) StateListener() {
	groupName := synch.groupCtx.Group.Name
	for true {
		select {
		case _ = <- synch.channelStop:
			glog.Infof("User '%s's group '%s' received a 'STOP' signal\n", synch.groupCtx.User.Name, synch.groupCtx.Group.Name)
			close(synch.channelStop)
			return
		default:
			time.Sleep(1 * time.Second)
			state, err := synch.groupCtx.Network.GetGroupState(groupName)
			if err != nil {
				glog.Warningf("could not get group state: Synchronizer.StateListener: %s\n", err)
				continue
			}
			if !bytes.Equal(state, synch.groupCtx.CalculateState(synch.groupCtx.Members, synch.groupCtx.Repo)) {
				glog.Infof("group state changed")
				go func() {
					operationBytes, err := synch.groupCtx.Network.GetGroupOperation(groupName, state)
					if err != nil {
						glog.Errorf("could not get operation: Synchronizer.StateListener: %s", err)
						return
					}
					var operation RawOperation
					if err := json.Unmarshal(operationBytes, &operation); err != nil {
						glog.Errorf("could not unmarshal operation: Synchronizer.StateListener: %s", err)
						return
					}
					cmd, err := NewOperation(&operation)
					if err != nil {
						glog.Errorf("could not create command from operation: Synchronizer.StateListener: %s", err)
						return
					}
					if err := cmd.Execute(synch.groupCtx); err != nil {
						glog.Errorf("error while executing cmd: Synchronizer.StateListener: %s", err)
						return
					}
				}()
			}
		}
	}
}

func (s *Synchronizer) MessageProcessor() {
	fmt.Printf("[*] Synchronizer for user '%s' group '%s' is running...\n", s.groupCtx.User.Name, s.groupCtx.Group.Name)
	for pubsubMessage := range s.channelPubSub {
		glog.Infof("--> user '%s' in group '%s' recieved a message", s.groupCtx.User.Name, s.groupCtx.Group.Name)
		groupMessageBytes, ok := pubsubMessage.Decrypt(s.groupCtx.Group.Boxer)
		if !ok {
			glog.Error("could not decrypt group message: Synchronizer: MessageProcessor")
			continue
		}

		var groupMessage GroupMessage
		if err := json.Unmarshal(groupMessageBytes, &groupMessage); err != nil {
			glog.Error("could not unmarshal group message: Synchronizer, MessageProzessor: %s", err)
			continue
		}

		switch groupMessage.Type {
		case "HB":
			if err := s.processHeartBeat(groupMessage); err != nil {
				glog.Error("could not process heart beat: Synchronizer.MessageProcessor: %s", err)
				continue
			}
		case "PROPOSAL":
			if err := s.processProposal(groupMessage); err != nil {
				glog.Error("error while processing PROPOSAL: Synchronizer.MessageProcessor: %s", err)
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
	member := s.groupCtx.Members.Get(heartBeat.From)
	if member == nil {
		return fmt.Errorf("heart beat from a non member user: Synchronizer.processHeartBeat")
	}
	_, ok := member.VerifyKey.Open(nil, heartBeat.Rand)
	if !ok {
		return fmt.Errorf("invalid heart beat: Synchronizer.processHeartBeat")
	}
	return nil
}

func (s *Synchronizer) processProposal(message GroupMessage) error {
	glog.Info("processing proposal")
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
		time.Sleep(1 * time.Second)
		var randomBytes [32]byte
		rand.Read(randomBytes[:])
		signedRand := s.groupCtx.User.Signer.SecretKey.Sign(nil, randomBytes[:])
		heartBeat := HeartBeat{
			From: s.groupCtx.User.Name,
			Rand: signedRand,
		}
		hbBytes, err := json.Marshal(heartBeat)
		if err != nil {
			glog.Error("could not marshal heart beat: Synchronizer.heartBeat: %s", err)
			continue
		}
		groupMessage := GroupMessage{
			 From:"HB",
			 Type: s.groupCtx.User.Name,
			 Data: hbBytes,
		}
		msgBytes, err := json.Marshal(groupMessage)
		if err != nil {
			glog.Error("could not marshal groupMessage: Synchronizer.heartBeat: %s", err)
			continue
		}
		if err := s.groupCtx.sendToAll(msgBytes); err != nil {
			glog.Error("could not send group message: Synchronizer.heartBeat: %s", err)
			continue
		}
	}
}

func (s *Synchronizer) approveTransaction(proposer string, transaction *Transaction) error {
	channelName := s.groupCtx.Group.Name + proposer
	signature := s.groupCtx.User.SignTransaction(transaction)
	approval := Approval{
		From:      s.groupCtx.User.Name,
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
	verifyKey, err := s.groupCtx.Network.GetUserVerifyKey(author)
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
		return fmt.Errorf("could not get state of group '%s': Synchronizer.validateTransaction: %s", s.groupCtx.Group.Name, err)
	}
	if !bytes.Equal(state, transaction.PrevState) {
		return fmt.Errorf("invlaid prev state in transaction proposal: Synchronizer.validateTransaction")
	}

	operation, err := NewOperation(&transaction.Operation)
	if err != nil {
		return fmt.Errorf("could not unmarshal operation: Synchronizer.validateTransaction: %s", err)
	}
	if err := operation.Validate(transaction.State, s.groupCtx); err != nil {
		return fmt.Errorf("error while validating transaction: Synchronizer.validateTransaction: %s", err)
	}
	return nil
}

// By operations (e.g. Invite()) a given number of valid approvals
// is needed to be able to commit the current operation. This func
// collects these approvals and upon receiving enough approvals it
// commits the operation
func (synch *Synchronizer) CollectApprovals(transaction *Transaction) {
	channelName := synch.groupCtx.Group.Name + synch.groupCtx.User.Name
	channel := make(chan ipfs.PubsubMessage)
	go synch.groupCtx.IPFS.PubsubSubscribe(synch.groupCtx.User.Name, channelName, channel)

	for {
		if len(transaction.SignedBy) > synch.groupCtx.Members.Length()/2 {
			transactionBytes, err := json.Marshal(transaction)
			if err != nil {
				glog.Error("could not marshal transaction: Synchronizer.CollectApprovals: %s", err)
				synch.groupCtx.IPFS.Kill(channelName)
				return
			}
			if err := synch.groupCtx.Network.GroupInvite(synch.groupCtx.Group.Name, transactionBytes); err != nil {
				glog.Error("could not call invite transaction: Synchronizer.CollectApprovals: %s", err)
				synch.groupCtx.IPFS.Kill(channelName)
				return
			}
			synch.groupCtx.IPFS.Kill(channelName)
			return
		}

		select {
		case pubsubMessage := <-channel:
			approvalBytes, ok := pubsubMessage.Decrypt(synch.groupCtx.Group.Boxer)
			if !ok {
				glog.Error("invalid group pubsub msg: Synchronizer.CollectApprovals")
				continue
			}
			var approval Approval
			if err := json.Unmarshal(approvalBytes, &approval); err != nil {
				glog.Error("could not unmarshal approval: Synchronizer.CollectApprovals: %s", err)
				continue
			}
			glog.Infof("got an approval from user '%s'", approval.From)
			if err := approval.Validate(transaction.Bytes(), synch.groupCtx.Group.Boxer, synch.groupCtx.Network); err != nil {
				glog.Error("could not validate approval: Synchronizer.CollectApprovals: %s", err)
				continue
			}
			signedBy := SignedBy{
				Username: approval.From,
				Signature: approval.Signature,
			}
			transaction.SignedBy = append(transaction.SignedBy, signedBy)
		case <-time.After(5 * time.Second):
			glog.Warningf("timeout: Synchronizer.CollectApprovals")
			synch.groupCtx.IPFS.Kill(channelName)
			return
		}
	}
}

func (synch *Synchronizer) Kill() {
	collectingChannel := synch.groupCtx.Group.Name + synch.groupCtx.User.Name
	synch.groupCtx.IPFS.Kill(synch.groupCtx.Group.Name)
	synch.groupCtx.IPFS.Kill(collectingChannel)
	synch.channelStop <- 1
}