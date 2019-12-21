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
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/glog"
	"github.com/pkg/errors"

	com "github.com/aliras1/FileTribe/client/communication"
	"github.com/aliras1/FileTribe/client/communication/common"
	"github.com/aliras1/FileTribe/client/fs"
	"github.com/aliras1/FileTribe/client/interfaces"
	. "github.com/aliras1/FileTribe/collections"
	ethaccount "github.com/aliras1/FileTribe/eth/gen/Account"
	ethapp "github.com/aliras1/FileTribe/eth/gen/FileTribeDApp"
	ethgroup "github.com/aliras1/FileTribe/eth/gen/Group"
	ipfsapi "github.com/aliras1/FileTribe/ipfs"
	"github.com/aliras1/FileTribe/tribecrypto"
)

// IUserFacade is an interface through which main.go can communicate
// with its UserContext
type IUserFacade interface {
	SignUp(username string) error
	CreateGroup(groupname string) error
	AcceptInvitation(groupAddress ethcommon.Address) error
	User() interfaces.Account
	Groups() []IGroupFacade
	SignOut()
	Transactions() ([]*types.Transaction, error)
}

// UserContext stores all the user data and it is responsible
// for handling communication, events, encryption, etc.
type UserContext struct {
	account     interfaces.Account
	eth         *Eth
	groups      *Map
	addressBook *common.AddressBook
	ipfs        ipfsapi.IIpfs
	storage     *fs.Storage
	p2p         *com.P2PManager
	p2pPort     string

	transactions *List
	invitations  *List
	subs         *List

	channelStop chan int
	lock        sync.RWMutex
}

// NewUserContext creates a new UserContext with the data provided
func NewUserContext(auth interfaces.Auth, backend bind.ContractBackend, appContractAddress ethcommon.Address, ipfs ipfsapi.IIpfs, p2pPort string) (*UserContext, error) {
	var err error
	var ctx UserContext

	appContract, err := ethapp.NewFileTribeDApp(appContractAddress, backend)
	if err != nil {
		return nil, errors.Wrap(err, "could not create account contract instance")
	}

	ctx.eth = &Eth{
		Backend: backend,
		App:     appContract,
		Auth:    auth,
	}
	ctx.p2pPort = p2pPort
	ctx.ipfs = ipfs
	ctx.groups = NewConcurrentMap()
	ctx.addressBook = common.NewAddressBook(backend, appContract, ipfs)
	ctx.transactions = NewConcurrentList()
	ctx.invitations = NewConcurrentList()
	ctx.subs = NewConcurrentList()
	ctx.channelStop = make(chan int)
	ctx.storage = fs.NewStorage(os.Getenv("HOME"))

	accountAddress, err := appContract.GetAccountOf(&bind.CallOpts{}, auth.Address())
	if err != nil {
		return nil, errors.Wrap(err, "could not get account address")
	}

	fmt.Print(accountAddress.Hex())

	if strings.EqualFold(accountAddress.String(), "0x"+strings.Repeat("0", 40)) {
		fmt.Println("No FileTribe account found associated to the current ethereum account. Use 'filetribe signup <username>' to sign up")

		go ctx.HandleAccountCreatedEvents(ctx.eth.App)
	} else {
		accountContract, err := ethaccount.NewAccount(accountAddress, ctx.eth.Backend)
		if err != nil {
			return nil, errors.Wrap(err, "could not create account contract object from address")
		}

		accountName, err := accountContract.Name(&bind.CallOpts{})
		if err != nil {
			return nil, errors.Wrap(err, "could not get account name from contract")
		}

		ctx.storage.Init(accountName)

		account, err := NewAccountFromStorage(ctx.storage, ctx.eth.Backend)
		if err != nil {
			return nil, errors.Wrap(err, "could not create account object")
		}

		if !bytes.Equal(account.ContractAddress().Bytes(), accountAddress.Bytes()) {
			if err := account.SetContract(accountAddress, ctx.eth.Backend); err != nil {
				return nil, errors.Wrap(err, "could not set contract")
			}

			if err := account.Save(); err != nil {
				return nil, errors.Wrap(err, "could not save account object")
			}
		}

		if err := ctx.Init(account); err != nil {
			return nil, errors.Wrap(err, "could not initialize user context")
		}
	}

	return &ctx, nil
}

// SignUp creates a new account, saves it and registers it on the blockchain
func (ctx *UserContext) SignUp(username string) error {
	glog.Infof("[*] account '%s' signing in...", username)

	acc, err := NewAccount(username, ctx.eth.Auth.Address(), ctx.storage)
	if err != nil {
		return errors.Wrap(err, "could not create new account")
	}

	ctx.storage.Init(username)

	if err := acc.Save(); err != nil {
		return errors.Wrap(err, "could not save account")
	}

	ipfsID, err := ctx.ipfs.ID()
	if err != nil {
		return errors.Wrap(err, "could not get ipfs id")
	}

	tx, err := ctx.eth.App.CreateAccount(ctx.eth.Auth.TxOpts(), username, ipfsID.ID, acc.Boxer().PublicKey.Value)
	if err != nil {
		return errors.Wrap(err, "could not send create account tx")
	}

	ctx.transactions.Add(tx)

	return nil
}

// IsMember returns whether a given user is a member of a given group or not.
// It is used by communication sessions that have no direct access to GroupContexts.
func (ctx *UserContext) IsMember(group ethcommon.Address, accountOwner ethcommon.Address) error {
	groupInt := ctx.groups.Get(group)
	if groupInt == nil {
		return errors.New("no group found")
	}

	if !groupInt.(*GroupContext).Group.IsMember(accountOwner) {
		return errors.New("account is not member of group")
	}

	return nil
}

// GetBoxerOfGroup returns the secret key of the given group. It is used
// by communication sessions that have no direct access to GroupContexts
func (ctx *UserContext) GetBoxerOfGroup(group ethcommon.Address) (tribecrypto.SymmetricKey, error) {
	groupInt := ctx.groups.Get(group)
	if groupInt == nil {
		return tribecrypto.SymmetricKey{}, errors.New("no group found")
	}

	return groupInt.(*GroupContext).Group.Boxer(), nil
}

// GetProposedBoxerOfGroup returns the secret key of the given group. It is
// used by communication sessions that have no direct access to GroupContexts
func (ctx *UserContext) GetProposedBoxerOfGroup(group ethcommon.Address, proposalKey []byte) (tribecrypto.SymmetricKey, error) {
	groupInt := ctx.groups.Get(group)
	if groupInt == nil {
		return tribecrypto.SymmetricKey{}, errors.New("no group found")
	}

	proposalInt := groupInt.(*GroupContext).proposals.Get(string(proposalKey))
	if proposalInt == nil {
		return tribecrypto.SymmetricKey{}, errors.New("no proposed key found")
	}

	return proposalInt.(*interfaces.Proposal).Boxer, nil
}

// Init initializes a UserContext: it starts the P2P manager and the event handlers
func (ctx *UserContext) Init(acc interfaces.Account) error {
	p2p, err := com.NewP2PManager(
		ctx.p2pPort,
		acc,
		ctx.eth.Auth.Sign,
		ctx.addressBook,
		ctx,
		ctx.ipfs)
	if err != nil {
		return errors.Wrap(err, "could not create P2P connection")
	}

	ctx.account = acc
	ctx.p2p = p2p

	// account events
	//go ctx.HandleDebugEvents(network.GetDebugChannel())
	go ctx.HandleGroupInvitationEvents(acc.Contract())
	go ctx.HandleGroupCreatedEvents(acc.Contract())
	go ctx.HandleInvitationAcceptedEvents(acc.Contract())

	if err := ctx.BuildGroups(); err != nil {
		return errors.Wrap(err, "could not build groups")
	}

	return nil
}

// User returns the account interface
func (ctx *UserContext) User() interfaces.Account {
	return ctx.account
}

// Save saves all UserContext data
func (ctx *UserContext) Save() error {
	//if err := ctx.Storage.SaveContextData(ctx); err != nil {
	//	return fmt.Errorf("could not save context data: %s", err)
	//}

	return nil
}

// SignOut tries to gracefully stop all started threads and processes
func (ctx *UserContext) SignOut() {
	glog.Infof("[*] account '%s' signing out...\n", ctx.account.Name())
	for groupCtx := range ctx.groups.VIterator() {
		groupCtx.(*GroupContext).Stop()
	}

	if err := ctx.Save(); err != nil {
		glog.Errorf("could not save context state: UserContext.SignOut: %s", err)
	}
}

// BuildGroups builds up all groups found on disk
func (ctx *UserContext) BuildGroups() error {
	glog.Infof("Building Groups for account '%s'...", ctx.account.Name())

	groupDatas, err := ctx.storage.GetGroupDatas()
	if err != nil {
		return errors.Wrap(err, "could not get group groupDatas")
	}

	for _, groupData := range groupDatas {
		groupContract, err := ethgroup.NewGroup(groupData.Address, ctx.eth.Backend)
		if err != nil {
			return errors.Wrap(err, "could not create new eth group instance")
		}

		config := &GroupContextConfig{
			Group:        NewGroupFromGroupData(groupData, groupContract, ctx.storage),
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
			return errors.Wrap(err, "could not create new group context")
		}

		if err := groupCtx.Update(); err != nil {
			return errors.Wrap(err, "could not update group ctx")
		}

		ctx.groups.Put(groupCtx.Address(), groupCtx)
	}

	glog.Infof("Building Groups ended")

	return nil
}

// CreateGroup creates a group through a blockchain method invoke
func (ctx *UserContext) CreateGroup(groupname string) error {
	tx, err := ctx.account.Contract().CreateGroup(ctx.eth.Auth.TxOpts(), groupname)
	if err != nil {
		return errors.Wrap(err, "could not send create group tx")
	}

	ctx.transactions.Add(tx)

	return nil
}

// AcceptInvitation accepts a group invitation
func (ctx *UserContext) AcceptInvitation(groupAddress ethcommon.Address) error {
	for otherAddressInt := range ctx.invitations.Iterator() {
		otherAddress := otherAddressInt.(ethcommon.Address)

		if bytes.Equal(groupAddress.Bytes(), otherAddress.Bytes()) {
			group, err := ethgroup.NewGroup(groupAddress, ctx.eth.Backend)
			if err != nil {
				return errors.Wrap(err, "could not get group contract instance")
			}

			tx, err := group.Join(ctx.eth.Auth.TxOpts())
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

// Groups returns a list of group facades
func (ctx *UserContext) Groups() []IGroupFacade {
	var groups []IGroupFacade

	for groupCtxInt := range ctx.groups.VIterator() {
		groups = append(groups, groupCtxInt.(IGroupFacade))
	}

	glog.Info(groups)

	return groups
}

// ListFiles lists the files names in the account's repository
func (ctx *UserContext) ListFiles() map[string][]string {
	list := make(map[string][]string)
	return list
}

// Transactions returns a list of transactions initiated by the user
func (ctx *UserContext) Transactions() ([]*types.Transaction, error) {
	var list []*types.Transaction

	for txInt := range ctx.transactions.Iterator() {
		list = append(list, txInt.(*types.Transaction))
	}

	return list, nil
}
