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
	"crypto/rand"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/glog"

	"github.com/aliras1/FileTribe/client/fs/meta"
	ethacc "github.com/aliras1/FileTribe/eth/gen/Account"
	ethapp "github.com/aliras1/FileTribe/eth/gen/FileTribeDApp"
	ethgroup "github.com/aliras1/FileTribe/eth/gen/Group"
	"github.com/aliras1/FileTribe/tribecrypto"
)

// HandleAccountCreatedEvents listens to 'AccountCreated' blockchain events
// and if one belongs to the current user, creates its appropriate UserContext
func (ctx *UserContext) HandleAccountCreatedEvents(app *ethapp.FileTribeDApp) {
	glog.Info("HandleAccountCreatedEvents...")
	ch := make(chan *ethapp.FileTribeDAppAccountCreated)

	sub, err := app.WatchAccountCreated(&bind.WatchOpts{Context: ctx.eth.Auth.TxOpts().Context}, ch)
	if err != nil {
		glog.Errorf("could not subscribe to AccountCreated events: %s", err)
		return
	}

	ctx.subs.Add(sub)

	for e := range ch {
		go ctx.onAccountCreated(e)
	}
}

func (ctx *UserContext) onAccountCreated(e *ethapp.FileTribeDAppAccountCreated) {
	if !bytes.Equal(e.Owner.Bytes(), ctx.eth.Auth.Address().Bytes()) {
		return
	}

	acc, err := NewAccountFromStorage(ctx.storage, ctx.eth.Backend)
	if err != nil {
		glog.Errorf("could not create account on eth: error while NewAccountFromStorage: %s", err)
		return
	}

	if err := acc.SetContract(e.Account, ctx.eth.Backend); err != nil {
		glog.Errorf("could not set account contract: %s, err")
		return
	}

	if err := acc.Save(); err != nil {
		glog.Errorf("could not save IAccount: %s", err)
		return
	}

	if err := ctx.Init(acc); err != nil {
		glog.Errorf("could not initialize user context: %s", err)
		return
	}

	glog.Infof("account created: %s --> %s (%s)", ctx.account.Name(), e.Account.String(), e.Owner.String())
}

// HandleGroupInvitationEvents listens to GroupCreated blockchain events
// and upon receiving one, it stores the invitation
func (ctx *UserContext) HandleGroupInvitationEvents(acc *ethacc.Account) {
	glog.Info("groupInvitation handling...")
	ch := make(chan *ethacc.AccountNewInvitation)

	sub, err := acc.WatchNewInvitation(&bind.WatchOpts{Context: ctx.eth.Auth.TxOpts().Context}, ch)
	if err != nil {
		glog.Errorf("could not subscribe to AccountNewInvitation events: %s", err)
		return
	}

	ctx.subs.Add(sub)

	for e := range ch {
		go ctx.onGroupInvitation(e)
	}
}

func (ctx *UserContext) onGroupInvitation(e *ethacc.AccountNewInvitation) {
	glog.Info("New INVITATION")

	glog.Infof("%s: got a group invitation into %s", ctx.account.Name(), e.Group.String())

	ctx.invitations.Add(e.Group)
}

// HandleInvitationAcceptedEvents listens to InvitationAccapted blockchain events
// and upon receiving one, it tries to get the group key and upon its success it
// creates the group's appropriate GroupContext
func (ctx *UserContext) HandleInvitationAcceptedEvents(acc *ethacc.Account) {
	glog.Info("HandleInvitationAcceptedEvents...")
	ch := make(chan *ethacc.AccountInvitationAccepted)

	sub, err := acc.WatchInvitationAccepted(&bind.WatchOpts{Context: ctx.eth.Auth.TxOpts().Context}, ch)
	if err != nil {
		glog.Errorf("could not subscribe to InvitationAccepted events: %s", err)
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

	memberOwners, err := group.MemberOwners(&bind.CallOpts{Pending: true})
	if err != nil {
		glog.Errorf("could not get group members from eth: %s", err)
		return
	}

	// Get key
	for _, memberOwner := range memberOwners {
		if bytes.Equal(memberOwner.Bytes(), ctx.eth.Auth.Address().Bytes()) {
			continue
		}

		member, err := ctx.eth.App.GetAccount(&bind.CallOpts{Pending: true}, memberOwner)
		if err != nil {
			glog.Errorf("could not get owner's account: %s", err)
			continue
		}

		contact, err := ctx.addressBook.GetFromAccountAddress(member)
		if err != nil {
			glog.Warningf("could not get contact for member: %s", member.String())
			continue
		}

		if err := ctx.p2p.StartGetGroupKeySession(
			e.Group,
			contact,
			e.Account,
			ctx.onGetKeySuccess,
		); err != nil {
			glog.Errorf("could not start get group key session: %s", err)
		}
	}
}

func (ctx *UserContext) onGetKeySuccess(groupAddressBytes []byte, boxer tribecrypto.SymmetricKey) {
	groupAddress := ethcommon.BytesToAddress(groupAddressBytes)
	exists := ctx.groups.Get(groupAddress)
	if exists != nil {
		return
	}

	groupMeta := &meta.GroupMeta{Address: groupAddress, Boxer: boxer}
	groupContract, err := ethgroup.NewGroup(groupMeta.Address, ctx.eth.Backend)
	if err != nil {
		glog.Errorf("could not create new eth group instance: %s", err)
		return
	}

	config := &GroupContextConfig{
		Group:        NewGroupFromMeta(groupMeta, groupContract, ctx.storage),
		Account:      ctx.account,
		P2P:          ctx.p2p,
		AddressBook:  ctx.addressBook,
		Ipfs:         ctx.ipfs,
		Storage:      ctx.storage,
		Transactions: ctx.transactions,
		Eth: &GroupEth{
			Group: groupContract,
			Eth:   ctx.eth,
		},
	}

	groupCtx, err := NewGroupContext(config)
	if err != nil {
		glog.Errorf("could not create new group ctx: %s", err)
		return
	}

	if err := groupCtx.Update(); err != nil {
		glog.Errorf("could not update group ctx: %s", err)
		return
	}

	ctx.groups.Put(groupCtx.Address(), groupCtx)

	glog.Info("group ctx created")
}

// HandleGroupCreatedEvents listens to GroupCreated blockchain events
// and upon receiving one, it creates the group's appropriate GroupContext
func (ctx *UserContext) HandleGroupCreatedEvents(acc *ethacc.Account) {
	glog.Info("GroupCreatedEvents...")
	ch := make(chan *ethacc.AccountGroupCreated)

	sub, err := acc.WatchGroupCreated(&bind.WatchOpts{Context: ctx.eth.Auth.TxOpts().Context}, ch)
	if err != nil {
		glog.Errorf("could not subscribe to GroupCreated events")
		return
	}

	ctx.subs.Add(sub)

	for e := range ch {
		go ctx.onGroupCreated(e)
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

	var secretKeyBytes [32]byte
	if _, err := rand.Read(secretKeyBytes[:]); err != nil {
		glog.Errorf("could not read rand: %s", err)
		return
	}

	boxer := tribecrypto.SymmetricKey{
		Key: secretKeyBytes,
		RNG: rand.Reader,
	}

	groupMeta := &meta.GroupMeta{
		Address: e.Group,
		Boxer:boxer,
	}

	config := &GroupContextConfig{
		Group:        NewGroupFromMeta(groupMeta, groupContract, ctx.storage),
		Account:      ctx.account,
		P2P:          ctx.p2p,
		AddressBook:  ctx.addressBook,
		Ipfs:         ctx.ipfs,
		Storage:      ctx.storage,
		Transactions: ctx.transactions,
		Eth: &GroupEth{
			Group: groupContract,
			Eth:   ctx.eth,
		},
	}

	groupCtx, err := NewGroupContext(config)
	if err != nil {
		glog.Errorf("could not create new group context: %s", err)
		return
	}

	if err := groupCtx.Save(); err != nil {
		glog.Errorf("could not save group: %s", err)
		return
	}

	ctx.groups.Put(e.Group, groupCtx)

	glog.Infof("Group created: %s", groupMeta.Address.String())
}
