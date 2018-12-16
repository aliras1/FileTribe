package client

import (
	"bytes"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/golang/glog"
	comcommon "ipfs-share/client/communication/common"
	"ipfs-share/client/fs/caps"
	ethacc "ipfs-share/eth/gen/Account"
	ethapp "ipfs-share/eth/gen/Dipfshare"
	ethgroup "ipfs-share/eth/gen/Group"
	ethinv "ipfs-share/eth/gen/Invitation"
)

func (ctx *UserContext) HandleAccountCreatedEvents(app *ethapp.Dipfshare) {
	glog.Info("HandleAccountCreatedEvents...")
	ch := make(chan *ethapp.DipfshareAccountCreated)

	sub, err := app.WatchAccountCreated(&bind.WatchOpts{Context:ctx.eth.Auth.TxOpts.Context}, ch)
	if err != nil {
		glog.Errorf("could not subscribe to AccountCreated events")
		return
	}

	ctx.subs.Add(sub)

	for e := range ch {
		go ctx.onAccountCreated(e)
	}
}

func (ctx *UserContext) onAccountCreated(e *ethapp.DipfshareAccountCreated) {
	if !bytes.Equal(e.Owner.Bytes(), ctx.eth.Auth.Address.Bytes()) {
		return
	}

	acc, err := NewAccountFromStorage(ctx.storage, ctx.eth.Backend)
	if err != nil {
		glog.Errorf("could not create account on eth: error while NewAccountFromStorage: %s", err)
		return
	}

	contract, err := ethacc.NewAccount(e.Account, ctx.eth.Backend)
	if err != nil {
		glog.Errorf("could not create account on eth: could not instanciate new account")
		return
	}

	acc.SetContractAddress(e.Account)
	acc.SetContract(contract)

	if err := ctx.Init(acc); err != nil {
		glog.Errorf("could not initialize user context: %s", err)
		return
	}

	glog.Infof("Account created: %s", e.Account.String())
}

func (ctx *UserContext) HandleGroupInvitationEvents(acc *ethacc.Account) {
	glog.Info("groupInvitation handling...")
	ch := make(chan *ethacc.AccountNewInvitation)

	sub, err := acc.WatchNewInvitation(&bind.WatchOpts{Context:ctx.eth.Auth.TxOpts.Context}, ch)
	if err != nil {
		glog.Errorf("could not subscribe to AccountNewInvitation events")
		return
	}

	ctx.subs.Add(sub)

	for e := range ch {
		go ctx.onGroupInvitation(e)
	}
}

func (ctx *UserContext) onGroupInvitation(e *ethacc.AccountNewInvitation) {
	glog.Info("New INVITATION")

	inv, err := ethinv.NewInvitation(e.Inv, ctx.eth.Backend)
	if err != nil {
		glog.Errorf("could not create group invitation instance: %s", err)
		return
	}

	groupAddr, err := inv.Group(&bind.CallOpts{Pending:true})
	if err != nil {
		glog.Errorf("could not get group Address from invitation: %s", err)
		return
	}

	glog.Infof("%s: got a group invitation into %s", ctx.account.Name(), groupAddr.String())

	ctx.invitations.Add(inv)
}

func (ctx *UserContext) HandleInvitationAcceptedEvents(acc *ethacc.Account) {
	glog.Info("HandleInvitationAcceptedEvents...")
	ch := make(chan *ethacc.AccountInvitationAccepted)

	sub, err := acc.WatchInvitationAccepted(&bind.WatchOpts{Context:ctx.eth.Auth.TxOpts.Context}, ch)
	if err != nil {
		glog.Errorf("could not subscribe to InvitationAccepted events")
		return
	}

	ctx.subs.Add(sub)

	for e := range ch {
		go ctx.onInvitationAccepted(e)
	}
}

func (ctx *UserContext) onInvitationAccepted(e *ethacc.AccountInvitationAccepted) {
	if !bytes.Equal(e.Account.Bytes(), ctx.account.ContractAddress().Bytes()) {
		return
	}

	glog.Info("Invitation accepted")

	group, err := ethgroup.NewGroup(e.Group, ctx.eth.Backend)
	if err != nil {
		glog.Errorf("could not create new eth group instance: %s", err)
		return
	}

	members, err := group.Members(&bind.CallOpts{Pending:true})
	if err != nil {
		glog.Errorf("could not get group members from eth")
		return
	}

	// Get key
	for _, member := range members {
		if bytes.Equal(member.Bytes(), ctx.account.ContractAddress().Bytes()) {
			continue
		}

		contact, err := ctx.addressBook.Get(member)
		if err != nil {
			glog.Warningf("could not get contact for member: %s", member.String())
			continue
		}

		if err := ctx.p2p.StartGetGroupKeySession(
			comcommon.GetGroupKey,
			e.Group,
			contact,
			e.Account,
			ctx.onGetKeySuccess,
		);	err != nil {
			glog.Errorf("could not start get group key session: %s", err)
		}
	}
}

func (ctx *UserContext) onGetKeySuccess(cap *caps.GroupAccessCap) {
	exists := ctx.groups.Get(cap.Address)
	if exists != nil {
		return
	}

	if err := ctx.storage.SaveGroupAccessCap(cap); err != nil {
		glog.Errorf("could not save group access cap: %s", err)
		return
	}

	contract, err := ethgroup.NewGroup(cap.Address, ctx.eth.Backend)
	if err != nil {
		glog.Errorf("could not create new eth group instance: %s", err)
		return
	}

	config := &GroupContextConfig{
		Group: NewGroupFromCap(cap, ctx.storage),
		Account: ctx.account,
		P2P: ctx.p2p,
		AddressBook: ctx.addressBook,
		Ipfs: ctx.ipfs,
		Storage: ctx.storage,
		Transactions: ctx.transactions,
		Eth:&GroupEth{
			Group:contract,
			Eth:ctx.eth,
		},
	}

	groupCtx, err := NewGroupContext(config)
	if err != nil {
		glog.Errorf("could not create new group ctx")
		return
	}

	ctx.groups.Put(groupCtx.Address(), groupCtx)

	glog.Info("group ctx created")
}

func (ctx *UserContext) HandleGroupConsensusSuccessfulEvents(ch chan *ethgroup.GroupConsensusReached) {
	glog.Info("groupUpdateIpfs handling...")

	for e := range ch {
		ctx.onGroupConsensus(e)
	}
}

func (ctx *UserContext) HandleGroupCreatedEvents(acc *ethacc.Account) {
	glog.Info("GroupCreatedEvents...")
	ch := make(chan *ethacc.AccountGroupCreated)

	sub, err := acc.WatchGroupCreated(&bind.WatchOpts{Context:ctx.eth.Auth.TxOpts.Context}, ch)
	if err != nil {
		glog.Errorf("could not subscribe to GroupCreated events")
		return
	}

	ctx.subs.Add(sub)

	for e := range ch {
		go ctx.onGroupCreated(e)
	}
}

func (ctx *UserContext) onGroupConsensus(e *ethgroup.GroupConsensusReached) {
	ctx.lock.Lock()
	defer ctx.lock.Unlock()

	glog.Info("got update ipfs event message")

	groupCtxInterface := ctx.groups.Get(e.Group)
	if groupCtxInterface == nil {
		return
	}

	groupCtx := groupCtxInterface.(*GroupContext)

	boxer := groupCtx.Group.Boxer()
	newIpfsHash, ok := boxer.BoxOpen(e.IpfsHash)
	if !ok {
		glog.Errorf("could not decrpyt new ipfs hash")
		return
	}

	if err := groupCtx.Update(); err != nil {
		glog.Errorf("could not update group context: %s", err)
	}

	if err := groupCtx.Repo.Update(string(newIpfsHash)); err != nil {
		glog.Errorf("could not update group %s's repo with ipfs hash %s: %s", groupCtx.Group.Address().String(), e.IpfsHash, err)
	}
}

func (ctx *UserContext) onGroupCreated(e *ethacc.AccountGroupCreated) {
	glog.Info("got a group created event")

	if !bytes.Equal(e.Account.Bytes(), ctx.account.ContractAddress().Bytes()) {
		return
	}

	groupContract, err := ethgroup.NewGroup(e.Group, ctx.eth.Backend)
	if err != nil {
		glog.Errorf("could not create new group contract instance: %s", err)
		return
	}

	groupName, err := groupContract.Name(&bind.CallOpts{Pending:true})
	if err != nil {
		glog.Errorf("could not get group name: %s", err)
		return
	}

	group := NewGroup(e.Group, groupName, ctx.storage)

	config := &GroupContextConfig{
		Group: group,
		Account: ctx.account,
		P2P: ctx.p2p,
		AddressBook: ctx.addressBook,
		Ipfs: ctx.ipfs,
		Storage: ctx.storage,
		Transactions: ctx.transactions,
		Eth:&GroupEth{
			Group:groupContract,
			Eth:ctx.eth,
		},
	}

	groupCtx, err := NewGroupContext(config)
	if err != nil {
		glog.Errorf("could not create new group context: %s", err)
		return
	}

	boxer := groupCtx.Group.Boxer()
	ipfsHash := groupCtx.Repo.IpfsHash()
	encIpfsHash := boxer.BoxSeal([]byte(ipfsHash))

	if err := group.SetIpfsHash(encIpfsHash); err != nil {
		glog.Errorf("could not set ipfs hash of group: %s", err)
		return
	}

	if err := groupCtx.Save(); err != nil {
		glog.Errorf("could not save group: %s", err)
		return
	}

	ctx.groups.Put(e.Group, groupCtx)

	glog.Infof("Group created: %s", group.Address().String())
}


func (ctx *UserContext) HandleKeyDirtyEvents(ch chan *ethgroup.GroupKeyDirty) {
	glog.Info("keyDirty handling...")

	for keyDirty := range ch {
		glog.Info("got Dirty Key event")
		go ctx.onKeyDirty(keyDirty)
	}
}

func (ctx *UserContext) onKeyDirty(keyDirty *ethgroup.GroupKeyDirty) {
	//id := NewBytesId(keyDirty.GroupId)
	//
	//groupCtxInt := ctx.groups.Get(id)
	//if groupCtxInt == nil {
	//	return
	//}
	//
	//groupCtx := groupCtxInt.(*GroupContext)
	//if err := groupCtx.Update(); err != nil {
	//	glog.Errorf("could not update group context: %s", err)
	//	return
	//}
	//
	//if err := groupCtx.onKeyDirty(); err != nil {
	//	glog.Errorf( "error while changing group key: %s", err)
	//}
}

//func (ctx *UserContext) HandleGroupLeftEvents(ch chan *eth.EthGroupLeft) {
//	glog.Info("GroupLeft handling...")
//
//	for groupLeft := range ch {
//		glog.Info("got GroupLeft event")
//		go ctx.onGroupLeft(groupLeft)
//	}
//}
//
//func (ctx *UserContext) onGroupLeft(event *eth.EthGroupLeft) {
//	if !bytes.Equal(event.account.Bytes(), ctx.account.Address().Bytes()) {
//		return
//	}
//
//	groupId := NewBytesId(event.GroupId)
//
//	if err := ctx.DeleteGroup(groupId); err != nil {
//		glog.Errorf("could not delete group: %s", err)
//	}
//}