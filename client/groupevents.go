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
	"bytes"
	"encoding/base64"
	"github.com/aliras1/FileTribe/client/interfaces"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/golang/glog"

	ethAccount "github.com/aliras1/FileTribe/eth/gen/Account"
	ethcons "github.com/aliras1/FileTribe/eth/gen/Consensus"
	ethgroup "github.com/aliras1/FileTribe/eth/gen/Group"
	"github.com/aliras1/FileTribe/tribecrypto"
)

// HandleGroupInvitationSentEvents listens to GroupInvitationSent events on the blockchain
func (groupCtx *GroupContext) HandleGroupInvitationSentEvents(group *ethgroup.Group) {
	glog.Info("HandleGroupInvitationSentEvents...")
	ch := make(chan *ethgroup.GroupInvitationSent)

	sub, err := group.WatchInvitationSent(&bind.WatchOpts{Context: groupCtx.eth.Auth.TxOpts().Context}, ch)
	if err != nil {
		glog.Errorf("could not subscribe to GroupInvitationSent events")
		return
	}

	groupCtx.subs.Add(sub)

	for e := range ch {
		glog.Infof("Group Invitation sent to: %s", e.Account.String())
	}
}

// HandleGroupInvitationAcceptedEvents listens to GroupInvitationAccepted events on the
// blockchain and adds the new member to the group, if it receives one
func (groupCtx *GroupContext) HandleGroupInvitationAcceptedEvents(group *ethgroup.Group) {
	glog.Info("HandleGroupInvitationAcceptedEvents...")
	ch := make(chan *ethgroup.GroupInvitationAccepted)

	sub, err := group.WatchInvitationAccepted(&bind.WatchOpts{Context: groupCtx.eth.Auth.TxOpts().Context}, ch)
	if err != nil {
		glog.Errorf("could not subscribe to InvitationAccepted events")
		return
	}

	groupCtx.subs.Add(sub)

	for e := range ch {
		glog.Infof("Group Invitation accepted by: %s", e.Account.String())

		account, err := ethAccount.NewAccount(e.Account, groupCtx.eth.Backend)
		if err != nil {
			glog.Errorf("could not get Account object from address: %s", e.Account.String())
			return
		}

		accountOwner, err := account.Owner(&bind.CallOpts{Pending: true})
		if err != nil {
			glog.Errorf("could not get owner of account: %s", err)
			return
		}

		groupCtx.Group.AddMember(accountOwner)
	}
}

// HandleNewConsensusEvents listens to NewConsensus events on the blockchain
// and checks if the target of the consensus is correct. If so it approves it
func (groupCtx *GroupContext) HandleNewConsensusEvents(group *ethgroup.Group) {
	glog.Info("HandleNewConsensusEvents...")
	ch := make(chan *ethgroup.GroupNewConsensus)

	sub, err := group.WatchNewConsensus(&bind.WatchOpts{Context: groupCtx.eth.Auth.TxOpts().Context}, ch)
	if err != nil {
		glog.Errorf("could not subscribe to NewConsensus events: %s", err)
		return
	}

	groupCtx.subs.Add(sub)

	for e := range ch {
		groupCtx.onNewConsensus(e)
	}
}

func (groupCtx *GroupContext) onNewConsensus(e *ethgroup.GroupNewConsensus) {
	glog.Infof("new CONSENSUS: %s", e.Consensus.String())

	cons, err := ethcons.NewConsensus(e.Consensus, groupCtx.eth.Backend)
	if err != nil {
		glog.Errorf("could not create new consensus instance from eth: %s", err)
		return
	}

	proposer, err := cons.Proposer(&bind.CallOpts{Pending: true})
	if err != nil {
		glog.Errorf("could not get the proposer of consensus: %s", err)
		return
	}
	glog.Infof("proposer of cons: %s", proposer.String())
	glog.Infof("account addr: %s", groupCtx.account.ContractAddress().String())

	if bytes.Equal(proposer.Bytes(), groupCtx.account.ContractAddress().Bytes()) {
		glog.Info("own consensus")
		return
	}

	payload, err := cons.Payload(&bind.CallOpts{Pending: true})
	if err != nil {
		glog.Errorf("could not get consensus payload: %s", err)
		return
	}

	glog.Infof("stored payload '%v' from: %s", payload, proposer.String())
	proposalKey := base64.StdEncoding.EncodeToString(payload)
	groupCtx.proposals.Put(proposalKey, &interfaces.Proposal{EncIpfsHash:payload, Proposer:proposer})

	groupCtx.getProposedKey(proposalKey)
}


// HandleIpfsHashChangedEvents listens to IpfsHashChanged events on the blockchain
// and if it receives one, it updates the group IPFS hash and fetches its contents
func (groupCtx *GroupContext) HandleIpfsHashChangedEvents(group *ethgroup.Group) {
	glog.Info("HandleIpfsHashChangedEvents...")
	ch := make(chan *ethgroup.GroupIpfsHashChanged)

	sub, err := group.WatchIpfsHashChanged(&bind.WatchOpts{Context: groupCtx.eth.Auth.TxOpts().Context}, ch)
	if err != nil {
		glog.Errorf("could not subscribe to IpfsHashChanged events: %s", err)
		return
	}

	groupCtx.subs.Add(sub)

	for e := range ch {
		groupCtx.onIpfsHashChanged(e)
	}
}

func (groupCtx *GroupContext) onIpfsHashChanged(e *ethgroup.GroupIpfsHashChanged) {
	glog.Info("IPFS HASH changed")

	if !bytes.Equal(e.Group.Bytes(), groupCtx.Group.Address().Bytes()) {
		return
	}

	for pk := range groupCtx.proposals.KIterator() {
		glog.Infof("%s: [[ PK ]]: %s", groupCtx.account.Name(), pk)
	}

	proposalKey := base64.StdEncoding.EncodeToString(e.IpfsHash)

	glog.Infof("{{{ proposal key: }}}  %s", proposalKey)

	proposalInt := groupCtx.proposals.Get(proposalKey)
	if proposalInt == nil || tribecrypto.BoxerIsNull(proposalInt.(*interfaces.Proposal).Boxer) {
		groupCtx.getKey()
	} else {
		if err := groupCtx.Group.SetBoxer(proposalInt.(*interfaces.Proposal).Boxer); err != nil {
			glog.Errorf("could not set new boxer: %s", err)
			return
		}

		if err := groupCtx.Update(); err != nil {
			glog.Errorf("could not update group context: %s", err)
		}
	}
}

func (groupCtx *GroupContext) HandleKeyEvents() {
	for e := range groupCtx.keyEventCh {
		glog.Infof("======= GOT KEY: %v =======", e.Args)
		key := e.Args.(tribecrypto.SymmetricKey)
		if tribecrypto.BoxerIsNull(key) {
			continue
		}

		currentKey := groupCtx.Group.Boxer()
		if bytes.Equal(key.Key[:], currentKey.Key[:]) {
			continue
		}

		if err := groupCtx.Group.SetBoxer(key); err != nil {
			glog.Errorf("could not set boxer: %s", err)
			continue
		}

		if err := groupCtx.Update(); err != nil {
			glog.Errorf("could not update group context: %s", err)
			continue
		}

		groupCtx.SetStateKeyInvalid(false)
		e.Sender.Cancel()
	}
}

func (groupCtx *GroupContext) HandleProposedKeyEvents() {
	for e := range groupCtx.proposedKeyEventCh {
		glog.Infof("%s: ======= GOT PROPOSED KEY: %v =======", groupCtx.account.Name(), e.Args)
		eventData := e.Args.(interfaces.Proposal)

		// TODO: check if the received key is correct, i.e. the payload can be decrypted

		proposal := groupCtx.proposals.Get(string(eventData.EncIpfsHash)).(*interfaces.Proposal)
		if proposal == nil {
			glog.Errorf("no proposal found to: %v", eventData.EncIpfsHash)
			continue
		}

		if proposal.EncIpfsHash == nil {
			glog.Error("payload is nil")
			continue
		}

		proposerAccount, err := ethAccount.NewAccount(proposal.Proposer, groupCtx.eth.Backend)
		if err != nil {
			glog.Errorf("could not get the account of the proposer: %s", err)
			continue
		}

		proposerOwner, err := proposerAccount.Owner(&bind.CallOpts{Pending: true})
		if err != nil {
			glog.Errorf("could not get owner of account: %s", err)
			continue
		}

		consensusAddress, err := groupCtx.eth.Group.GetConsensus(&bind.CallOpts{Pending: true}, proposerOwner)
		if err != nil {
			glog.Errorf("could not get member's consensus: %s", err)
			continue
		}

		consensus, err := ethcons.NewConsensus(consensusAddress, groupCtx.eth.Backend)
		if err != nil {
			glog.Errorf("could not create new consensus instance from eth: %s", err)
			continue
		}

		ipfsHash, ok := eventData.Boxer.BoxOpen(proposal.EncIpfsHash)
		if !ok {
			glog.Errorf("could not decrypt consensus payload")
			continue
		}

		if err := groupCtx.Repo.IsValidChangeSet(string(ipfsHash), eventData.Boxer, proposal.Proposer); err != nil {
			glog.Errorf("invalid changeset: %s", err)
			continue
		}

		proposal.Boxer = eventData.Boxer

		if err := groupCtx.approveConsensus(consensus); err != nil {
			glog.Errorf("could not approve consensus: %s", err)
			continue
		}

		glog.Infof("%s: consensus %s approved", groupCtx.account.Name(), consensusAddress.String())
	}
}

func (groupCtx *GroupContext) HandleDebugEvents(group *ethgroup.Group) {
	glog.Info("HandleDebugEvents...")
	ch := make(chan *ethgroup.GroupDebug)

	sub, err := group.WatchDebug(&bind.WatchOpts{Context: groupCtx.eth.Auth.TxOpts().Context}, ch)
	if err != nil {
		glog.Errorf("could not subscribe to IpfsHashChanged events: %s", err)
		return
	}

	groupCtx.subs.Add(sub)

	for e := range ch {
		glog.Infof("DEBUG: %s", e.Msg.String())
	}
}

func (groupCtx *GroupContext) HandleDebugConsEvents(cons *ethcons.Consensus) {
	glog.Info("HandleDebugConsEvents...")
	ch := make(chan *ethcons.ConsensusDebugCons)

	sub, err := cons.WatchDebugCons(&bind.WatchOpts{Context: groupCtx.eth.Auth.TxOpts().Context}, ch)
	if err != nil {
		glog.Errorf("could not subscribe to IpfsHashChanged events: %s", err)
		return
	}

	groupCtx.subs.Add(sub)

	for e := range ch {
		glog.Infof("DEBUG CONS: %s", e.Msg.String())
	}
}