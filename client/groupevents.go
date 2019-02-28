package client

import (
	"bytes"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	comcommon "ipfs-share/client/communication/common"
	"ipfs-share/client/fs/caps"
	ethcons "ipfs-share/eth/gen/Consensus"
	ethgroup "ipfs-share/eth/gen/Group"
	"ipfs-share/utils"
)

const (
	IPFS_HASH = iota
	KEY = 		iota
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

	ctype, err := cons.Ctype(&bind.CallOpts{Pending:true})
	if err != nil {
		glog.Errorf("could not get consensus type: %s", err)
		return
	}

	payload, err := cons.Payload(&bind.CallOpts{Pending:true})
	if err != nil {
		glog.Errorf("could not get consensus payload: %s", err)
		return
	}

	proposer, err := cons.Proposer(&bind.CallOpts{Pending:true})
	if err != nil {
		glog.Errorf("could not get proposer from consensus: %s", err)
		return
	}

	switch ctype {
	case IPFS_HASH:
		key := groupCtx.Group.Boxer()
		ipfsHash, ok := key.BoxOpen(payload)
		if !ok {
			glog.Errorf("could not decrypt consensus payload")
			return
		}

		if err := groupCtx.Repo.IsValidChangeSet(string(ipfsHash), proposer); err != nil {
			glog.Errorf("invalid changeset: %s", err)
			return
		}

		if err := groupCtx.approveConsensus(cons); err != nil {
			glog.Errorf("could not approve consensus")
			return
		}

	case KEY:
		// Get proposed key
		for _, member := range groupCtx.Group.Members() {
			if bytes.Equal(member.Bytes(), groupCtx.account.ContractAddress().Bytes()) {
				continue
			}

			contact, err := groupCtx.AddressBook.Get(member)
			if err != nil {
				glog.Warningf("could not get contact for member: %s", member.String())
				continue
			}

			if err := groupCtx.P2P.StartGetGroupKeySession(
				comcommon.GetProposedGroupKey,
				e.Group,
				contact,
				groupCtx.account.ContractAddress(),
				groupCtx.onGetProposedKeySuccess,
			);	err != nil {
				glog.Errorf("could not start get group key session: %s", err)
			}
		}
	}
}

func (groupCtx *GroupContext) onGetProposedKeySuccess(cap *caps.GroupAccessCap) {
	glog.Errorf("not implemented: onGetProposedKeySuccess")
}

func (groupCtx *GroupContext) approveConsensus(cons *ethcons.Consensus) error {
	digest, err := cons.Digest(&bind.CallOpts{Pending:true})
	if err != nil {
		return errors.Wrap(err, "could not get digest from consensus")
	}

	sig, err := groupCtx.eth.Auth.Signer.Sign(digest[:])
	if err != nil {
		return errors.WithMessage(err, "could not sign digest")
	}

	r, s, v, err := utils.SigToRSV(sig)
	if err != nil {
		return errors.Wrap(err, "could not convert sig to r,s,v")
	}

	tx, err := cons.Approve(groupCtx.eth.Auth.TxOpts, r, s, v)
	if err != nil {
		return errors.Wrap(err, "could not send consensus approve tx")
	}

	groupCtx.Transactions.Add(tx)

	return nil
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

	if err := groupCtx.Group.SetIpfsHash(e.IpfsHash); err != nil {
		glog.Errorf("could not set ipfs hash of group")
		return
	}

	if err := groupCtx.Group.Save(); err != nil {
		glog.Errorf("could not save group: %s", err)
	}

	if err := groupCtx.Repo.Update(groupCtx.Group.IpfsHash()); err != nil {
		glog.Errorf("could not update group repo: %s", err)
	}
}

func (groupCtx *GroupContext) HandleKeyDirtyEvents(group *ethgroup.Group) {
	glog.Info("HandleKeyDirtyEvents...")
	ch := make(chan *ethgroup.GroupKeyDirty)

	sub, err := group.WatchKeyDirty(&bind.WatchOpts{Context:groupCtx.eth.Auth.TxOpts.Context}, ch)
	if err != nil {
		glog.Errorf("could not subscribe to GroupKeyDirty events: %s", err)
		return
	}

	groupCtx.subs.Add(sub)

	for e := range ch {
		if !bytes.Equal(e.Group.Bytes(), groupCtx.Group.Address().Bytes()) {
			continue
		}

		glog.Info("new KeyDirty event")

		if err := groupCtx.onKeyDirty(); err != nil {
			glog.Errorf("error while KeyDirty: %s", err)
		}
	}
}

