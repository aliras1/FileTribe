package client

import (
	"bytes"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/glog"
	"ipfs-share/crypto"
	ethcons "ipfs-share/eth/gen/Consensus"
	ethgroup "ipfs-share/eth/gen/Group"
)


func (groupCtx *GroupContext) HandleGroupInvitationSentEvents(group *ethgroup.Group) {
	glog.Info("HandleGroupInvitationSentEvents...")
	ch := make(chan *ethgroup.GroupInvitationSent)

	sub, err := group.WatchInvitationSent(&bind.WatchOpts{Context:groupCtx.eth.Auth.TxOpts.Context}, ch)
	if err != nil {
		glog.Errorf("could not subscribe to GroupInvitationSent events")
		return
	}

	groupCtx.subs.Add(sub)

	for e := range ch {
		glog.Infof("Group Invitation sent to: %s", e.Account.String())
	}
}

func (groupCtx *GroupContext) HandleGroupInvitationAcceptedEvents(group *ethgroup.Group) {
	glog.Info("HandleGroupInvitationAcceptedEvents...")
	ch := make(chan *ethgroup.GroupInvitationAccepted)

	sub, err := group.WatchInvitationAccepted(&bind.WatchOpts{Context:groupCtx.eth.Auth.TxOpts.Context}, ch)
	if err != nil {
		glog.Errorf("could not subscribe to InvitationAccepted events")
		return
	}

	groupCtx.subs.Add(sub)

	for e := range ch {
		glog.Infof("Group Invitation accepted by: %s", e.Account.String())
		groupCtx.Group.AddMember(e.Account)
	}
}

func (groupCtx *GroupContext) HandleNewConsensusEvents(group *ethgroup.Group) {
	glog.Info("HandleNewConsensusEvents...")
	ch := make(chan *ethgroup.GroupNewConsensus)

	sub, err := group.WatchNewConsensus(&bind.WatchOpts{Context:groupCtx.eth.Auth.TxOpts.Context}, ch)
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
	glog.Info("new CONSENSUS")

	cons, err := ethcons.NewConsensus(e.Consensus, groupCtx.eth.Backend)
	if err != nil {
		glog.Errorf("could not create new consensus instance from eth: %s", err)
		return
	}

	proposer, err := cons.Proposer(&bind.CallOpts{Pending:true})
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

	payload, err := cons.Payload(&bind.CallOpts{Pending:true})
	if err != nil {
		glog.Errorf("could not get consensus payload: %s", err)
		return
	}

	glog.Infof("stored payload from: %s", proposer.String())
	groupCtx.proposedPayloads.Put(proposer, payload)


	// TODO: only ask k' from those that have approved
	// 1. get those that voted
	// 2. foreach voter: start a get proposed group key session
	glog.Infof("my account addr: %s", (groupCtx.account.ContractAddress()).String())
	glog.Info(groupCtx.Group.Members())

	// Get proposed key
	for _, member := range groupCtx.Group.Members() {
		if bytes.Equal(member.Bytes(), groupCtx.account.ContractAddress().Bytes()) {
			continue
		}

		glog.Infof("speaking to: %s", member.String())

		contact, err := groupCtx.AddressBook.Get(member)
		if err != nil {
			glog.Warningf("could not get contact for member: %s", member.String())
			continue
		}

		if err := groupCtx.P2P.StartGetProposedGroupKeySession(
			e.Group,
			member,
			contact,
			groupCtx.account.ContractAddress(),
			groupCtx.onGetProposedKeySuccess,
		);	err != nil {
			glog.Errorf("could not start get group key session: %s", err)
		}
	}
}

func (groupCtx *GroupContext) onGetProposedKeySuccess(proposer ethcommon.Address, boxer tribecrypto.SymmetricKey) {
	// TODO: check if the received key is correct, i.e. the payload can be decrypted

	glog.Infof("GOT proposed key: %v with proposer: %s", boxer.Key, proposer.String())

	groupCtx.proposedKeys.Put(proposer, boxer)

	payloadInt := groupCtx.proposedPayloads.Get(proposer)
	if payloadInt == nil {
		glog.Error("payload is nil")
		return
	}

	consensusAddress, err := groupCtx.eth.Group.GetConsensus(&bind.CallOpts{Pending:true}, proposer)
	if err != nil {
		glog.Errorf("could not get member's consensus: %s", err)
		return
	}

	consensus, err := ethcons.NewConsensus(consensusAddress, groupCtx.eth.Backend)
	if err != nil {
		glog.Errorf("could not create new consensus instance from eth: %s", err)
		return
	}

	ipfsHash, ok := boxer.BoxOpen(payloadInt.([]byte))
	if !ok {
		glog.Errorf("could not decrypt consensus payload")
		return
	}

	if err := groupCtx.Repo.IsValidChangeSet(string(ipfsHash), boxer, proposer); err != nil {
		glog.Errorf("invalid changeset: %s", err)
		return
	}

	if err := groupCtx.approveConsensus(consensus); err != nil {
		glog.Errorf("could not approve consensus: %s", err)
		return
	}

	groupCtx.proposedPayloads.Put(proposer, nil)
}

func (groupCtx *GroupContext) HandleIpfsHashChangedEvents(group *ethgroup.Group) {
	glog.Info("HandleIpfsHashChangedEvents...")
	ch := make(chan *ethgroup.GroupIpfsHashChanged)

	sub, err := group.WatchIpfsHashChanged(&bind.WatchOpts{Context:groupCtx.eth.Auth.TxOpts.Context}, ch)
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

	newBoxerInt := groupCtx.proposedKeys.Get(e.Proposer)
	if newBoxerInt == nil {
		for _, member := range groupCtx.Group.Members() {
			if bytes.Equal(member.Bytes(), groupCtx.account.ContractAddress().Bytes()) {
				continue
			}

			contact, err := groupCtx.AddressBook.Get(member)
			if err != nil {
				glog.Warningf("could not get contact for member: %s", member.String())
				continue
			}

			onGetKeySuccess := func(_ ethcommon.Address, newBoxer tribecrypto.SymmetricKey) {
				// if already got k' --> return
				currentBoxer := groupCtx.Group.Boxer()
				if bytes.Equal(currentBoxer.Key[:], newBoxer.Key[:]) {
					return
				}
				groupCtx.Group.SetBoxer(newBoxer)
				if err := groupCtx.Update(); err != nil {
					glog.Errorf("could not update group context: %s", err)
				}
			}

			if err := groupCtx.P2P.StartGetGroupKeySession(
				e.Group,
				contact,
				groupCtx.account.ContractAddress(),
				onGetKeySuccess,
			);	err != nil {
				glog.Errorf("could not start get group key session: %s", err)
			}
		}
	} else {
		groupCtx.Group.SetBoxer(newBoxerInt.(tribecrypto.SymmetricKey))
		if err := groupCtx.Update(); err != nil {
			glog.Errorf("could not update group context: %s", err)
		}
	}
}
