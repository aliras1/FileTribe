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
	"context"
	"os"
	"sync"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/chequebook"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/glog"
	"github.com/pkg/errors"

	com "github.com/aliras1/FileTribe/client/communication"
	"github.com/aliras1/FileTribe/client/communication/common"
	"github.com/aliras1/FileTribe/client/fs"
	"github.com/aliras1/FileTribe/client/interfaces"
	. "github.com/aliras1/FileTribe/collections"
	ethApp "github.com/aliras1/FileTribe/eth/gen/FileTribeDApp"
	ethgroup "github.com/aliras1/FileTribe/eth/gen/Group"
	ipfsapi "github.com/aliras1/FileTribe/ipfs"
	"github.com/aliras1/FileTribe/tribecrypto"
)

type IUserFacade interface {
	SignUp(username string) error
	CreateGroup(groupname string) error
	AcceptInvitation(groupAddress ethcommon.Address) error
	User() interfaces.IAccount
	Groups() []IGroupFacade
	SignOut()
	Transactions() ([]*types.Receipt, error)
}

type UserContext struct {
	account      	interfaces.IAccount
	eth             *Eth
	groups       	*Map
	addressBook  	*common.AddressBook
	ipfs         	ipfsapi.IIpfs
	storage      	*fs.Storage
	p2p          	*com.P2PManager
	p2pPort 		string

	transactions 	*List
	invitations     *List
	subs 			*List

	channelStop 	chan int
	lock 			sync.RWMutex
}


func (ctx *UserContext) SignUp(username string) error {
	glog.Infof("[*] Account '%s' signing in...", username)

	acc, err := NewAccount(username, ctx.storage)
	if err != nil {
		return errors.Wrap(err, "could not create new account")
	}

	if err := acc.Save(); err != nil {
		return errors.Wrap(err, "could not save account")
	}

	ipfsId, err := ctx.ipfs.ID()
	if err != nil {
		return errors.Wrap(err, "could not get ipfs id")
	}

	tx, err := ctx.eth.App.CreateAccount(ctx.eth.Auth.TxOpts, username, ipfsId.ID, acc.Boxer().PublicKey.Value)
	if err != nil {
		return errors.Wrap(err, "could not send create account tx")
	}

	ctx.transactions.Add(tx)

	return nil
}

func (ctx *UserContext) SignIn(username string) error {
	//auth, err := NewAuth(ethKeyPath, password)
	//if err != nil {
	//	return errors.Wrap(err, "could not create Auth")
	//}
	//
	//acc, err := NewAccount()

	return errors.New("not implemented")
}

func (ctx *UserContext) IsMember(group ethcommon.Address, account ethcommon.Address) error {
	groupInt := ctx.groups.Get(group)
	if groupInt == nil {
		return errors.New("no group found")
	}

	if !groupInt.(*GroupContext).Group.IsMember(account) {
		return errors.New("account is not member of group")
	}

	return nil
}

func (ctx *UserContext) Boxer(group ethcommon.Address) (tribecrypto.SymmetricKey, error) {
	groupInt := ctx.groups.Get(group)
	if groupInt == nil {
		return tribecrypto.SymmetricKey{}, errors.New("no group found")
	}

	return groupInt.(*GroupContext).Group.Boxer(), nil
}

func (ctx *UserContext) ProposedBoxer(group ethcommon.Address, proposer ethcommon.Address) (tribecrypto.SymmetricKey, error) {
	groupInt := ctx.groups.Get(group)
	if groupInt == nil {
		return tribecrypto.SymmetricKey{}, errors.New("no group found")
	}

	boxerInt := groupInt.(*GroupContext).proposedKeys.Get(proposer)
	if boxerInt == nil {
		return tribecrypto.SymmetricKey{}, errors.New("no proposed key found")
	}

	return boxerInt.(tribecrypto.SymmetricKey), nil
}

func (ctx *UserContext) Init(acc interfaces.IAccount) error {
	p2p, err := com.NewP2PManager(
		ctx.p2pPort,
		acc,
		ctx.eth.Auth.Signer,
		ctx.addressBook,
		ctx,
		ctx.ipfs)
	if err != nil {
		return errors.Wrap(err, "could not create P2P connection")
	}

	ctx.account = acc
	ctx.p2p = p2p

	// Account events
	//go ctx.HandleDebugEvents(network.GetDebugChannel())
	go ctx.HandleGroupInvitationEvents(acc.Contract())
	go ctx.HandleGroupCreatedEvents(acc.Contract())
	go ctx.HandleInvitationAcceptedEvents(acc.Contract())

	if err := ctx.BuildGroups(); err != nil {
		return errors.Wrap(err, "could not build groups")
	}

	return nil
}

func NewUserContext(auth *Auth, backend chequebook.Backend, appContractAddress ethcommon.Address, ipfs ipfsapi.IIpfs, p2pPort string) (*UserContext, error) {
	var err error
	var ctx UserContext

	appContract, err := ethApp.NewFileTribeDApp(appContractAddress, backend)
	if err != nil {
		return nil, errors.Wrap(err, "could not create account contract instance")
	}

	ctx.eth = &Eth{
		Backend:	backend,
		App:		appContract,
		Auth:		auth,
	}
	ctx.p2pPort = p2pPort
	ctx.ipfs = ipfs
	ctx.groups = NewConcurrentMap()
	ctx.addressBook = common.NewAddressBook(backend, appContract, ipfs)
	ctx.transactions = NewConcurrentList()
	ctx.invitations = NewConcurrentList()
	ctx.subs = NewConcurrentList()
	ctx.channelStop = make(chan int)
	ctx.storage = fs.NewStorage(os.Getenv("HOME") + "/.filetribe/" + auth.Address.String())

	// app events
	go ctx.HandleAccountCreatedEvents(ctx.eth.App)

	return &ctx, nil
}

func (ctx *UserContext) GetGroupData(addr ethcommon.Address) (interfaces.IGroup, *fs.GroupRepo) {
	groupCtxInt := ctx.groups.Get(addr)
	if groupCtxInt == nil {
		return nil, nil
	}

	groupCtx := groupCtxInt.(*GroupContext)

	return groupCtx.Group, groupCtx.Repo
}

func (ctx *UserContext) User() interfaces.IAccount {
	return ctx.account
}

func (ctx *UserContext) Save() error {
	//if err := ctx.Storage.SaveContextData(ctx); err != nil {
	//	return fmt.Errorf("could not save context data: %s", err)
	//}

	return nil
}

func (ctx *UserContext) SignOut() {
	glog.Infof("[*] Account '%s' signing out...\n", ctx.account.Name())
	for groupCtx := range ctx.groups.VIterator() {
		groupCtx.(*GroupContext).Stop()
	}

	if err := ctx.Save(); err != nil {
		glog.Errorf("could not save context state: UserContext.SignOut: %s", err)
	}
}

func (ctx *UserContext) BuildGroups() error {
	//glog.Infof("Building Groups for account '%s'...", ctx.account.Name())
	//caps, err := ctx.storage.GetGroupCaps()
	//if err != nil {
	//	return errors.Wrap(err, "could not get group caps")
	//}
	//
	//for _, cap := range caps {
	//	groupContract, err := ethGroup.NewGroup(cap.Address, ctx.backend)
	//
	//	groupCtx, err := NewGroupContextFromCAP(
	//		&cap,
	//		ctx.account,
	//		ctx.p2p,
	//		ctx.addressBook,
	//		groupContract,
	//		ctx.ipfs,
	//		ctx.storage,
	//		ctx.transactions,
	//	)
	//	if err != nil {
	//		return errors.Wrap(err, "could not create new group context")
	//	}
	//	if err := ctx.groups.Put(groupCtx); err != nil {
	//		glog.Warningf("could not append elem: %s", err)
	//	}
	//}
	//glog.Infof("Building Groups ended")
	return nil
}

func (ctx *UserContext) CreateGroup(groupname string) error {
	tx, err := ctx.account.Contract().CreateGroup(ctx.eth.Auth.TxOpts, groupname)
	if err != nil {
		return errors.Wrap(err, "could not send create group tx")
	}

	ctx.transactions.Add(tx)

	return nil
}

func (ctx *UserContext) AcceptInvitation(groupAddress ethcommon.Address) error {
	for otherAddressInt := range ctx.invitations.Iterator() {
		otherAddress := otherAddressInt.(ethcommon.Address)

		if bytes.Equal(groupAddress.Bytes(), otherAddress.Bytes()) {
			group, err := ethgroup.NewGroup(groupAddress, ctx.eth.Backend)
			if err != nil {
				return errors.Wrap(err, "could not get group contract instance")
			}

			tx, err := group.Join(ctx.eth.Auth.TxOpts)
			if err != nil {
				return errors.Wrap(err, "could not send accept invitation tx")
			}

			ctx.transactions.Add(tx)

			return nil
		}
	}

	return errors.New("Group not found in invitations")
}

func (ctx *UserContext) disposeGroup(groupAddr ethcommon.Address) error {
	groupCtxInt := ctx.groups.Delete(groupAddr)
	if groupCtxInt == nil {
		return errors.New("no group found")
	}

	groupCtx := groupCtxInt.(*GroupContext)
	groupCtx.Stop()

	//if err := ctx.storage.RemoveGroupDir(groupId.Data().([32]byte)); err != nil {
	//	return errors.Wrap(err, "could not remove group dir")
	//}

	return nil
}

func (ctx *UserContext) Groups() []IGroupFacade {
	var groups []IGroupFacade

	for groupCtxInt := range ctx.groups.VIterator() {
		groups = append(groups, groupCtxInt.(IGroupFacade))
	}

	glog.Info(groups)

	return groups
}

// Files lists the content of the account's repository
func (ctx *UserContext) List() map[string][]string {
	list := make(map[string][]string)
	return list
}

func (ctx *UserContext) Transactions() ([]*types.Receipt, error) {
	var list []*types.Receipt

	for txInt := range ctx.transactions.Iterator() {
		r, err := ctx.eth.Backend.TransactionReceipt(context.Background(), txInt.(*types.Transaction).Hash())
		if err != nil {
			return nil, errors.Wrap(err, "could not get tx receipt")
		}

		list = append(list, r)
	}

	return list, nil
}

func (ctx *UserContext) OnChangeGroupKeyServerSessionSuccess(args []interface{}, groupId IIdentifier) {
	//if len(args) < 2 {
	//	glog.Error("error while OnServerSessionSuccess: invalid number of args")
	//	return
	//}
	//
	//boxer := args[1].(crypto.SymmetricKey)
	//encNewIpfsHash := args[0].([]byte)
	//encNewIpfsHashBase64 := base64.StdEncoding.EncodeToString(encNewIpfsHash)
	//
	//groupCtxInt := ctx.groups.Get(groupId)
	//if groupCtxInt == nil {
	//	glog.Error("no group found")
	//}
	//
	//groupCtx := groupCtxInt.(*GroupContext)
	//groupCtx.proposedKeys[encNewIpfsHashBase64] = boxer
}