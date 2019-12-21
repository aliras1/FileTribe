// Copyright (c) 2019 Laszlo Sari
//
// FileTribe is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// FileTribe is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//

package client

import (
	"crypto/rand"
	"encoding/base64"
	"path"
	"sync"
	"time"

	"github.com/aliras1/FileTribe/asynctask"
	"github.com/aliras1/FileTribe/client/tasks"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/glog"
	"github.com/pkg/errors"

	com "github.com/aliras1/FileTribe/client/communication"
	"github.com/aliras1/FileTribe/client/communication/common"
	sesscommon "github.com/aliras1/FileTribe/client/communication/sessions/common"
	"github.com/aliras1/FileTribe/client/fs"
	"github.com/aliras1/FileTribe/client/interfaces"
	. "github.com/aliras1/FileTribe/collections"
	ethacc "github.com/aliras1/FileTribe/eth/gen/Account"
	ethcons "github.com/aliras1/FileTribe/eth/gen/Consensus"
	ipfsapi "github.com/aliras1/FileTribe/ipfs"
	"github.com/aliras1/FileTribe/tribecrypto"
)

// IGroupFacade is an interface to main.go through which it can communicate
// with a GroupContext
type IGroupFacade interface {
	Address() ethcommon.Address
	Name() string
	GrantWriteAccess(filePath string, user ethcommon.Address) error
	RevokeWriteAccess(filePath string, user ethcommon.Address) error
	CommitChanges() error
	Invite(user ethcommon.Address, hasInviteRigth bool) error
	Leave() error
	ListFiles() []FileView
	ListMembers() []MemberView
}

// MemberView is a view of a group member. These objects are sent back
// to main.go when it lists group members
type MemberView struct {
	Name    string
	Address string
}

// FileView is a view of a file objects. These objects are sent back
// to main.go when it lists the group repository
type FileView struct {
	Name        string
	WriteAccess []MemberView
}

type GroupState int

const (
	KeyInvalid GroupState = 0
)

type StateTransition struct {
	State GroupState
	Value bool
}

// GroupContext represents a groups current state and is responsible for
// all the communication, storage, encryption work
type GroupContext struct {
	account                  interfaces.Account
	Group                    interfaces.Group
	P2P                      *com.P2PManager
	Repo                     *fs.GroupRepo
	GroupConnection          *com.GroupConnection
	AddressBook              *common.AddressBook
	state                    map[GroupState]bool
	eth                      *GroupEth
	Ipfs                     ipfsapi.IIpfs
	Storage                  *fs.Storage
	Transactions             *List
	broadcastChannel         *ipfsapi.PubSubSubscription
	proposals                *Map
	proposedKeys             *Map
	proposedPayloads         *Map
	subs                     *List
	keyInvalidValueChangedCh chan bool
	stopCh                   chan struct{}
	keyEventCh               chan *asynctask.Event
	proposedKeyEventCh       chan *asynctask.Event
	lock                     sync.Mutex
}

// GroupContextConfig is a configuration struct for creating GroupContext
type GroupContextConfig struct {
	Group        interfaces.Group
	Account      interfaces.Account
	P2P          *com.P2PManager
	AddressBook  *common.AddressBook
	Eth          *GroupEth
	Ipfs         ipfsapi.IIpfs
	Storage      *fs.Storage
	Transactions *List
}

// NewGroupContext creates a GroupContext with data described in the
// provided configuration object
func NewGroupContext(config *GroupContextConfig) (*GroupContext, error) {

	groupContext := &GroupContext{
		account:                  config.Account,
		Group:                    config.Group,
		P2P:                      config.P2P,
		GroupConnection:          nil,
		AddressBook:              config.AddressBook,
		eth:                      config.Eth,
		Ipfs:                     config.Ipfs,
		Storage:                  config.Storage,
		Transactions:             config.Transactions,
		state:                    make(map[GroupState]bool),
		subs:                     NewConcurrentList(),
		proposals:                NewConcurrentMap(),
		proposedKeys:             NewConcurrentMap(),
		proposedPayloads:         NewConcurrentMap(),
		keyInvalidValueChangedCh: make(chan bool),
		stopCh:                   make(chan struct{}),
		keyEventCh:               make(chan *asynctask.Event, 1000),
		proposedKeyEventCh:       make(chan *asynctask.Event, 1000),
	}

	repo, err := fs.NewGroupRepo(config.Group, config.Account.ContractAddress(), config.Storage, config.Ipfs)
	if err != nil {
		return nil, errors.Wrap(err, "could not create group repo")
	}

	groupContext.Repo = repo
	groupContext.GroupConnection = com.NewGroupConnection(
		config.Group,
		repo,
		config.Account,
		config.AddressBook,
		onSessionClosed,
		config.P2P,
		config.Ipfs,
	)

	go groupContext.HandleGroupInvitationSentEvents(config.Eth.Group)
	go groupContext.HandleGroupInvitationAcceptedEvents(config.Eth.Group)
	go groupContext.HandleNewConsensusEvents(config.Eth.Group)
	go groupContext.HandleIpfsHashChangedEvents(config.Eth.Group)
	go groupContext.HandleDebugEvents(config.Eth.Group)

	go groupContext.handleStateKeyInvalid()
	go groupContext.HandleKeyEvents()
	go groupContext.HandleProposedKeyEvents()

	go groupContext.StartHeartbeat()

	return groupContext, nil
}

func (groupCtx *GroupContext) StartHeartbeat() {
	for {
		<- time.After(2 * time.Minute)
		timeEncoded, err := time.Now().MarshalJSON()
		if err != nil {
			glog.Errorf("could not encode the timestamp")
			continue
		}
		msg, err := common.NewMessage(groupCtx.account.Owner(), common.Heartbeat, 0, timeEncoded, groupCtx.eth.Auth.Sign)
		if err != nil {
			glog.Errorf("could not create new heartbeat message: %s", err)
			continue
		}
		msgBytes, err := msg.Encode()
		if err != nil {
			glog.Errorf("could not encode heartbeat message: %s", err)
			continue
		}
		
		if err := groupCtx.GroupConnection.Broadcast(msgBytes); err != nil {
			glog.Errorf("could not broadcast hearbeat message: %s", err)
		}
	}
}

func onSessionClosed(session sesscommon.Session) {
	glog.Infof("session %d closed with error: %s", session.ID(), session.Error())
}

func (groupCtx *GroupContext) SetStateKeyInvalid(value bool) {
	groupCtx.keyInvalidValueChangedCh <- value
}

func (groupCtx *GroupContext) handleStateKeyInvalid() {
	for {
		select {
		case value := <-groupCtx.keyInvalidValueChangedCh:
			glog.Infof("/////// changed to: %v", value)
			if groupCtx.state[KeyInvalid] != value {
				groupCtx.state[KeyInvalid] = value
			}
			if groupCtx.state[KeyInvalid] {
				groupCtx.getKey()
			}

		case <-time.After(5 * time.Second):
			if groupCtx.state[KeyInvalid] {
				groupCtx.getKey()
			}

		case <-groupCtx.stopCh:
			// todo: handle stop
		}
	}
}

func (groupCtx *GroupContext) getKey() {
	tasks.NewGetGroupKeyAsyncTask(
		groupCtx.account,
		groupCtx.Group,
		groupCtx.P2P,
		groupCtx.keyEventCh,
	).Execute()
}

func (groupCtx *GroupContext) getProposedKey(proposalKey string) {
	tasks.NewGetProposalKeyAsyncTask(
		groupCtx.account,
		groupCtx.Group,
		proposalKey,
		groupCtx.P2P,
		groupCtx.proposedKeyEventCh,
	).Execute()
}

// Update fetches all the current group information from the blockchain
// and refreshes the GroupContext with its contents
func (groupCtx *GroupContext) Update() error {
	if err := groupCtx.Group.Update(); err != nil {
		return errors.Wrap(err, "could not update group")
	}

	if err := groupCtx.Group.Save(); err != nil {
		return errors.Wrap(err, "could not save group")
	}

	if err := groupCtx.Repo.Update(groupCtx.Group.IpfsHash()); err != nil {
		return errors.Wrap(err, "could not update group repo")
	}

	return nil
}

func (groupCtx *GroupContext) UpdateGroupKey() {
	if len(groupCtx.Group.EncryptedIpfsHash()) == 0 {
		groupCtx.SetStateKeyInvalid(false)
		return
	}

	boxer := groupCtx.Group.Boxer()
	if _, ok := boxer.BoxOpen(groupCtx.Group.EncryptedIpfsHash()); ok {
		return
	}

	groupCtx.SetStateKeyInvalid(true)
}

// Leave invokes the 'Leave' operation of the group on the blockchain
func (groupCtx *GroupContext) Leave() error {
	tx, err := groupCtx.eth.Group.Leave(groupCtx.eth.Auth.TxOpts())
	if err != nil {
		return errors.Wrap(err, "could not send leave group tx")
	}

	groupCtx.Transactions.Add(tx)

	return nil
}

// Cancel kills all IPFS pubsub group connection - NOT USED
func (groupCtx *GroupContext) Stop() {
	groupCtx.GroupConnection.Kill()
}

// CommitChanges collects all changes in the group's root directory,
// creates a path from it and commits the changes on the blockchain
func (groupCtx *GroupContext) CommitChanges() error {
	var secretKeyBytes [32]byte
	if _, err := rand.Read(secretKeyBytes[:]); err != nil {
		return errors.Wrap(err, "could not read crypto/rand")
	}

	newKey := tribecrypto.SymmetricKey{
		Key: secretKeyBytes,
		RNG: rand.Reader,
	}

	glog.Infof("NEW KEY: %v", newKey.Key)

	hash, err := groupCtx.Repo.CommitChanges(newKey)
	if err != nil {
		return errors.Wrap(err, "could not commit group repo's changes")
	}

	encIpfsHash := newKey.BoxSeal([]byte(hash))
	proposalKey := base64.StdEncoding.EncodeToString(encIpfsHash)

	proposer := groupCtx.account.ContractAddress()
	groupCtx.proposals.Put(proposalKey, &interfaces.Proposal{Proposer: proposer, Boxer: newKey, EncIpfsHash: encIpfsHash})

	numOnlineMembers := 0
	for _, owner := range groupCtx.Group.MemberOwners() {
		contact, err := groupCtx.AddressBook.GetFromOwnerAddress(owner)
		if err != nil {
			return errors.Wrap(err, "could not get contact from owner")
		}

		if contact.IsOnline {
			numOnlineMembers += 1
		}		
	}

	if numOnlineMembers > groupCtx.Group.CountMembers() / 2 {
		newKeyBytes, err := newKey.Encode()
		if err != nil {
			return errors.Wrap(err, "could not encode new key")
		}
		commitData := common.CommitData{
			NewIpfsHash: encIpfsHash,
			NewKey: newKeyBytes,
		}
		commitDataBytes, err := commitData.Encode()
		if err != nil {
			return errors.Wrap(err, "could not encode commit data")
		}
		msg, err := common.NewMessage(groupCtx.account.Owner(), common.Commit, 0, commitDataBytes, groupCtx.eth.Auth.Sign)
		if err != nil {
			return errors.Wrap(err, "could not create new commit message")
		}
		msgBytes, err := msg.Encode()
		if err != nil {
			return errors.Wrap(err, "could not encode commit data message")
		}
		if err := groupCtx.GroupConnection.Broadcast(msgBytes); err != nil {
			return errors.Wrap(err, "could not broadcast commit message")
		}
	} else {
		tx, err := groupCtx.eth.Group.Commit(groupCtx.eth.Auth.TxOpts(), encIpfsHash)
		if err != nil {
			return errors.Wrap(err, "could not send change ipfs hash tx")
		}
	
		groupCtx.Transactions.Add(tx)
	}	

	return nil
}

// Invite invokes the 'Invite' method of the group on the blockchain
func (groupCtx *GroupContext) Invite(newMember ethcommon.Address, hasInviteRight bool) error {
	glog.Infof("[*] Inviting account '%s' into group '%s'...\n", newMember.String(), groupCtx.Group.Name())

	tx, err := groupCtx.eth.Group.Invite(groupCtx.eth.Auth.TxOpts(), newMember)
	if err != nil {
		return errors.Wrap(err, "could not send invite account tx")
	}

	groupCtx.Transactions.Add(tx)

	return nil
}

// Save stores group data on disk
func (groupCtx *GroupContext) Save() error {
	if err := groupCtx.Group.Save(); err != nil {
		return errors.Wrap(err, "could not save group")
	}

	return nil
}

// GrantWriteAccess adds the defined user to the write ACL in the file meta
func (groupCtx *GroupContext) GrantWriteAccess(filePath string, accountAddress ethcommon.Address) error {
	account, err := ethacc.NewAccount(accountAddress, groupCtx.eth.Backend)
	if err != nil {
		return errors.Wrap(err, "could not get Account object from address")
	}

	accountOwner, err := account.Owner(&bind.CallOpts{Pending: true})
	if err != nil {
		return errors.Wrap(err, "could not get owner of account")
	}

	if !groupCtx.Group.IsMember(accountOwner) {
		return errors.New("can not grant write access to non group members")
	}

	file := groupCtx.Repo.Get(path.Base(filePath))
	if file == nil {
		tmpFile, err := fs.NewGroupFile(
			filePath,
			[]ethcommon.Address{groupCtx.account.ContractAddress()},
			groupCtx.Group.Address().String(),
			groupCtx.Storage)
		if err != nil {
			return errors.Wrap(err, "could not create new group file")
		}
		file = tmpFile
	}

	if err := file.GrantWriteAccess(groupCtx.account.ContractAddress(), accountAddress); err != nil {
		return errors.Wrap(err, "could not grant write access to account")
	}

	return nil
}

// RevokeWriteAccess removes the defined user from the write ACL in the file meta
func (groupCtx *GroupContext) RevokeWriteAccess(filePath string, user ethcommon.Address) error {
	if !groupCtx.Group.IsMember(user) {
		return errors.New("can not revoke write access from non group members")
	}

	file := groupCtx.Repo.Get(path.Base(filePath))
	if file == nil {
		tmpFile, err := fs.NewGroupFile(
			filePath,
			[]ethcommon.Address{groupCtx.account.ContractAddress()},
			groupCtx.Group.Address().String(),
			groupCtx.Storage)
		if err != nil {
			return errors.Wrap(err, "could not create new group file")
		}
		file = tmpFile
	}

	if err := file.RevokeWriteAccess(groupCtx.account.ContractAddress(), user); err != nil {
		return errors.Wrap(err, "could not revoke write access to account")
	}

	return nil
}

func (groupCtx *GroupContext) startGetKey(encNewIpfsHash []byte) error {
	//newBoxer, ok := groupCtx.proposedKeys[encNewIpfsHashBase64]
	//
	//if ok {
	//	groupCtx.onGetKeySuccess(newBoxer)
	//} else {
	//	for _, member := range groupCtx.Group.MemberOwners() {
	//		if bytes.Equal(member.Bytes(), groupCtx.account.ContractAddress().Bytes()) {
	//			continue
	//		}
	//
	//		c, err := groupCtx.AddressBook.Get(member)
	//		if err != nil {
	//			glog.Warningf("could not get contact for member: %s", member.String())
	//			continue
	//		}
	//
	//		if err := groupCtx.P2P.StartGetGroupKeySession(
	//			groupCtx.Group.EthAccountAddress(),
	//			c,
	//			groupCtx.account.ContractAddress(),
	//			func(cap *caps.GroupAccessCap) {
	//				groupCtx.onGetKeySuccess(cap.GetBoxerOfGroup)
	//			},
	//		);	err != nil {
	//			glog.Errorf("could not start get group key session: %s", err)
	//		}
	//	}
	//}

	return nil
}

// Address returns the smart contract address of the group
func (groupCtx *GroupContext) Address() ethcommon.Address {
	return groupCtx.Group.Address()
}

// Name returns the group name
func (groupCtx *GroupContext) Name() string {
	return groupCtx.Group.Name()
}

// ListFiles returns a list of type string with the group files
func (groupCtx *GroupContext) ListFiles() []FileView {
	var list []FileView

	for _, file := range groupCtx.Repo.Files() {
		var acl []MemberView
		for _, address := range file.Meta.WriteAccessList {
			member := MemberView{Address: address.String()}

			contact, err := groupCtx.AddressBook.GetFromOwnerAddress(address)
			if err != nil {
				glog.Errorf("could not get contact for address: '%s': %s", address, err)
				member.Name = "<error>"
			} else {
				member.Name = contact.Name
			}

			acl = append(acl, member)
		}

		list = append(list, FileView{Name: file.Meta.FileName, WriteAccess: acl})
	}

	return list
}

// ListMembers returns a list of the members addresses
func (groupCtx *GroupContext) ListMembers() []MemberView {
	var list []MemberView
	addresses := groupCtx.Group.MemberOwners()

	for _, address := range addresses {
		member := MemberView{Address: address.String()}

		contact, err := groupCtx.AddressBook.GetFromOwnerAddress(address)
		if err != nil {
			glog.Errorf("could not get contact for address '%s': %s", address.String(), err)
			member.Name = "<error>"
		} else {
			member.Name = contact.Name
		}

		list = append(list, member)
	}

	return list
}

func (groupCtx *GroupContext) broadcast(msg []byte) error {
	return groupCtx.GroupConnection.Broadcast(msg)
}

func (groupCtx *GroupContext) p2pBroadcast(msg []byte) error {
	for _, member := range groupCtx.Group.MemberOwners() {

		c, err := groupCtx.AddressBook.GetFromOwnerAddress(member)
		if err != nil {
			glog.Warningf("could not get contact for member: %s", member)
		}

		go func() {
			if err := c.Send(msg); err != nil {
				glog.Errorf("error while sending p2p message: %s", err)
			}
		}()
	}

	return nil
}

func (groupCtx *GroupContext) approveConsensus(cons *ethcons.Consensus) error {
	tx, err := cons.Approve(groupCtx.eth.Auth.TxOpts())
	if err != nil {
		return errors.Wrapf(err, "could not send consensus approve tx with arguments: %v", groupCtx.eth.Auth.TxOpts)
	}

	groupCtx.Transactions.Add(tx)

	return nil
}
