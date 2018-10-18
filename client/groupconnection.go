package client

import (
	"github.com/golang/glog"
	ipfsapi "ipfs-share/ipfs"
	"bytes"
	"encoding/base64"
	"github.com/pkg/errors"
)

const s_GroupChannelGetName  = "/getkey"

type GroupConnection struct {
	groupCtx *GroupContext

	channelState  chan []byte
	channelStop   chan int

	groupSubscription ipfsapi.IPubSubSubscription
}

func NewGroupConnection(groupCtx *GroupContext) *GroupConnection {
	glog.Infof("Creating group connection...")

	conn := GroupConnection{
		groupCtx: groupCtx,
		channelStop: make(chan int),
	}

	id := groupCtx.Group.Id.Data().([32]byte)
	groupIdString := base64.URLEncoding.EncodeToString(id[:])

	glog.Infof("%s: subscribing to ipfs pubsub topic %s", groupCtx.User.Name, groupIdString)

	sub, err := groupCtx.Ipfs.PubSubSubscribe(groupIdString)
	if err != nil {
		glog.Errorf("%s: could not ipfs subscribe to topic %s", groupCtx.User.Name, groupIdString)
		return nil
	}

	conn.groupSubscription = sub

	go conn.connectionListener()
	//go conn.StateListener()
	//go conn.heartBeat()

	return &conn
}


func (conn *GroupConnection) SendAll(msg []byte) error {
	id := conn.groupCtx.Group.Id.Data().([32]byte)
	topic := base64.URLEncoding.EncodeToString(id[:])

	encMsg := conn.groupCtx.Group.Boxer.BoxSeal(msg)
	msgString := base64.URLEncoding.EncodeToString(encMsg)

	if err := conn.groupCtx.Ipfs.PubSubPublish(topic, msgString); err != nil {
		return errors.Wrap(err, "could not send group message to all members")
	}

	return nil
}

func (conn *GroupConnection) StateListener() {
	// groupName := synch.groupCtx.Group.Name
	// for true {
	// 	select {
	// 	case _ = <-synch.channelStop:
	// 		glog.Infof("User '%s's group '%s' received a 'STOP' signal\n", synch.groupCtx.User.Name, synch.groupCtx.Group.Name)
	// 		close(synch.channelStop)
	// 		return
	// 	default:
	// 		time.Sleep(1 * time.Second)
	// 		state, err := synch.groupCtx.Network.GetGroupState(groupName)
	// 		if err != nil {
	// 			glog.Warningf("could not get group state: GroupConnection.StateListener: %s\n", err)
	// 			continue
	// 		}
	// 		if !bytes.Equal(state, synch.groupCtx.CalculateState(synch.groupCtx.Members, synch.groupCtx.Repo)) {
	// 			glog.Infof("group state changed")
	// 			go func() {
	// 				operationBytes, err := synch.groupCtx.Network.GetGroupOperation(groupName, state)
	// 				if err != nil {
	// 					glog.Errorf("could not get operation: GroupConnection.StateListener: %s", err)
	// 					return
	// 				}
	// 				var operation RawOperation
	// 				if err := json.Unmarshal(operationBytes, &operation); err != nil {
	// 					glog.Errorf("could not unmarshal operation: GroupConnection.StateListener: %s", err)
	// 					return
	// 				}
	// 				cmd, err := NewOperation(&operation)
	// 				if err != nil {
	// 					glog.Errorf("could not create command from operation: GroupConnection.StateListener: %s", err)
	// 					return
	// 				}
	// 				if err := cmd.Execute(synch.groupCtx); err != nil {
	// 					glog.Errorf("error while executing cmd: GroupConnection.StateListener: %s", err)
	// 					return
	// 				}
	// 			}()
	// 		}
	// 	}
	// }
}

func (conn *GroupConnection) connectionListener() {
	glog.Infof("GroupConnection for user '%s' group '%s' is running...\n", conn.groupCtx.User.Name, conn.groupCtx.Group.Name)
	for {
		select {
		case <- conn.channelStop:
			{
				conn.groupSubscription.Cancel()
				close(conn.channelStop)
				return
			}
		default:
			{
				pubsubRecord, err := conn.groupSubscription.Next()
				if err != nil {
					glog.Warning("could not get next pubsub record")
					continue
				}

				glog.Infof("%s got a group message", conn.groupCtx.User.Name)

				encMsg, err := base64.URLEncoding.DecodeString((string)(pubsubRecord.Data()))
				if err != nil {
					glog.Warningf("could not url decode group message: %s", err)
					continue
				}

				msgData, ok := conn.groupCtx.Group.Boxer.BoxOpen(encMsg)
				if !ok {
					glog.Warningf("could not decrypt pubsub message")
					continue
				}

				msg, err := DecodeGroupMessage(msgData)
				if err != nil {
					glog.Warning("could not decode pubsub record message")
					continue
				}

				senderIsInGroup := false
				for _, member := range conn.groupCtx.Group.Members {
					if bytes.Equal(member.Bytes(), msg.From.Bytes()) {
						senderIsInGroup = true
						break
					}
				}

				if !senderIsInGroup {
					glog.Warningf("non group member %v has written to the group channel", msg.From.Bytes())
					continue
				}

				contact, err := msg.Validate(conn.groupCtx.Network, conn.groupCtx.Ipfs)
				if err != nil {
					glog.Warningf("invalid pubsub message to group %v from user %v", conn.groupCtx.Group.Id.Data(), msg.From.Bytes())
					continue
				}

				// TODO: check this with Ipfs address at the beginning
				if bytes.Equal(contact.Address.Bytes(), conn.groupCtx.User.Address.Bytes()) {
					continue
				}

				// append new contact to address book. if one already exists, therefore
				// it's P2P connection is not null, we will not try to create a new one
				// later
				if err := conn.groupCtx.AddressBook.Append(contact); err != nil {
					glog.Warningf("could not append elem: %s", err)
				}
				contact = conn.groupCtx.AddressBook.Get(contact.Id()).(*Contact)

				session := NewServerGroupSession(msg, contact, conn.groupCtx)
				conn.groupCtx.P2P.AddSession(session)
				go session.Run()
			}
		}
	}
}

func (conn *GroupConnection) processHeartBeat(heartBeat HeartBeat) error {
	// var heartBeat HeartBeat
	// err := json.Unmarshal(message.Data, &heartBeat)
	// if err != nil {
	// 	return err
	// }
	// member := s.groupCtx.Members.Get(heartBeat.From)
	// if member == nil {
	// 	return fmt.Errorf("heart beat from a non member user: GroupConnection.processHeartBeat")
	// }
	// _, ok := member.VerifyKey.Open(nil, heartBeat.Rand)
	// if !ok {
	// 	return fmt.Errorf("invalid heart beat: GroupConnection.processHeartBeat")
	// }
	return nil
}

func (conn *GroupConnection) processProposal(proposal Proposal) error {
	glog.Info("processing proposal")
	//transaction, err := s.authenticateProposal(message.From, message.Data)
	//if err != nil {
	//	return fmt.Errorf("could not authenticate proposal: GroupConnection.processProposal: %s", err)
	//}
	//if err := s.validateTransaction(transaction); err != nil {
	//	return fmt.Errorf("error while validating transaction: GroupConnection.processProposal: %s", err)
	//}
	//if err := s.approveTransaction(message.From, transaction); err != nil {
	//	return fmt.Errorf("could not approve proposal: GroupConnection.processProposal: %s", err)
	//}
	return nil
}

func (conn *GroupConnection) heartBeat() {
	// for {
	// 	time.Sleep(1 * time.Second)
	// 	var randomBytes [32]byte
	// 	rand.Read(randomBytes[:])
	// 	signedRand := s.groupCtx.User.Signer.VerifyKey.Sign(nil, randomBytes[:])
	// 	heartBeat := HeartBeat{
	// 		From: s.groupCtx.User.Name,
	// 		Rand: signedRand,
	// 	}
	// 	hbBytes, err := json.Marshal(heartBeat)
	// 	if err != nil {
	// 		glog.Error("could not marshal heart beat: GroupConnection.heartBeat: %s", err)
	// 		continue
	// 	}
	// 	groupMessage := GroupMessage{
	// 		From: "HB",
	// 		Type: s.groupCtx.User.Name,
	// 		Data: hbBytes,
	// 	}
	// 	msgBytes, err := json.Marshal(groupMessage)
	// 	if err != nil {
	// 		glog.Error("could not marshal groupMessage: GroupConnection.heartBeat: %s", err)
	// 		continue
	// 	}
	// 	if err := s.groupCtx.SendToAll(msgBytes); err != nil {
	// 		glog.Error("could not send group message: GroupConnection.heartBeat: %s", err)
	// 		continue
	// 	}
	// }
}

func (conn *GroupConnection) approveTransaction(proposer string, transaction *Transaction) error {
	//channelName := s.groupCtx.Group.Name + proposer
	//signature := s.groupCtx.User.SignTransaction(transaction)
	//approval := Approval{
	//	From:      s.groupCtx.User.Name,
	//	Signature: signature,
	//}
	//approvalBytes, err := json.Marshal(&approval)
	//if err != nil {
	//	return fmt.Errorf("could not marshal approval: GroupConnection.approveTransaction: %s", err)
	//}
	//groupEncApproval := s.groupCtx.Group.Boxer.BoxSeal(approvalBytes)
	//if err := s.groupCtx.Ipfs.PubSubPublish(channelName, base64.StdEncoding.EncodeToString(groupEncApproval)); err != nil {
	//	return fmt.Errorf("could not send approval: GroupConnection.approveTransaction: %s", err)
	//}
	return nil
}

// authenticate if the proposal really comes from the given user
func (conn *GroupConnection) authenticateProposal(author string, data []byte) (*Transaction, error) {
	// verifyKey, err := s.groupCtx.Network.GetUserVerifyKey(author)
	// if err != nil {
	// 	return nil, fmt.Errorf("could not get user verify key: GroupConnection.authenticateProposal: %s", err)
	// }
	// transactionBytes, ok := verifyKey.Open(nil, data)
	// if !ok {
	// 	return nil, fmt.Errorf("invalid proposal message from user '%s': GroupConnection.authenticateProposal", author)
	// }
	// var proposal Transaction
	// if err := json.Unmarshal(transactionBytes, &proposal); err != nil {
	// 	return nil, fmt.Errorf("could not unmarshal proposal: GroupConnection.authenticateProposal: %s", err)
	// }
	// return &proposal, nil
	return nil, nil
}

// validate the content of a transaction
func (conn *GroupConnection) validateTransaction(transaction *Transaction) error {
	//state, err := s.groupCtx.GetState()
	//if err != nil {
	//	return fmt.Errorf("could not get state of group '%s': GroupConnection.validateTransaction: %s", s.groupCtx.Group.Name, err)
	//}
	//if !bytes.Equal(state, transaction.PrevState) {
	//	return fmt.Errorf("invlaid prev state in transaction proposal: GroupConnection.validateTransaction")
	//}
	//
	//operation, err := NewOperation(&transaction.Operation)
	//if err != nil {
	//	return fmt.Errorf("could not unmarshal operation: GroupConnection.validateTransaction: %s", err)
	//}
	//if err := operation.Validate(transaction.State, s.groupCtx); err != nil {
	//	return fmt.Errorf("error while validating transaction: GroupConnection.validateTransaction: %s", err)
	//}
	return nil
}

// By operations (e.g. Invite()) a given number of valid approvals
// is needed to be able to commit the current operation. This func
// collects these approvals and upon receiving enough approvals it
// commits the operation
func (conn *GroupConnection) CollectApprovals(transaction *Transaction) {
	// channelName := synch.groupCtx.Group.Name + synch.groupCtx.User.Name
	// channel := make(chan ipfs.PubsubMessage)
	// go synch.groupCtx.Ipfs.PubsubSubscribe(synch.groupCtx.User.Name, channelName, channel)

	// for {
	// 	if len(transaction.SignedBy) > synch.groupCtx.Members.Length()/2 {
	// 		transactionBytes, err := json.Marshal(transaction)
	// 		if err != nil {
	// 			glog.Error("could not marshal transaction: GroupConnection.CollectApprovals: %s", err)
	// 			synch.groupCtx.Ipfs.Kill(channelName)
	// 			return
	// 		}
	// 		if err := synch.groupCtx.Network.GroupInvite(synch.groupCtx.Group.Name, transactionBytes); err != nil {
	// 			glog.Error("could not call invite transaction: GroupConnection.CollectApprovals: %s", err)
	// 			synch.groupCtx.Ipfs.Kill(channelName)
	// 			return
	// 		}
	// 		synch.groupCtx.Ipfs.Kill(channelName)
	// 		return
	// 	}

	// 	select {
	// 	case pubsubMessage := <-channel:
	// 		approvalBytes, ok := pubsubMessage.Decrypt(synch.groupCtx.Group.Boxer)
	// 		if !ok {
	// 			glog.Error("invalid group pubsub msg: GroupConnection.CollectApprovals")
	// 			continue
	// 		}
	// 		var approval Approval
	// 		if err := json.Unmarshal(approvalBytes, &approval); err != nil {
	// 			glog.Error("could not unmarshal approval: GroupConnection.CollectApprovals: %s", err)
	// 			continue
	// 		}
	// 		glog.Infof("got an approval from user '%s'", approval.From)
	// 		if err := approval.Validate(transaction.Bytes(), synch.groupCtx.Group.Boxer, synch.groupCtx.Network); err != nil {
	// 			glog.Error("could not validate approval: GroupConnection.CollectApprovals: %s", err)
	// 			continue
	// 		}
	// 		signedBy := SignedBy{
	// 			Username:  approval.From,
	// 			Signature: approval.Signature,
	// 		}
	// 		transaction.SignedBy = append(transaction.SignedBy, signedBy)
	// 	case <-time.After(5 * time.Second):
	// 		glog.Warningf("timeout: GroupConnection.CollectApprovals")
	// 		synch.groupCtx.Ipfs.Kill(channelName)
	// 		return
	// 	}
	// }
}

func (conn *GroupConnection) Kill() {
	// collectingChannel := synch.groupCtx.Group.Name + synch.groupCtx.User.Name
	// synch.groupCtx.Ipfs.Kill(synch.groupCtx.Group.Name)
	// synch.groupCtx.Ipfs.Kill(collectingChannel)
	conn.channelStop <- 1
}
