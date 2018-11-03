package client

import (
	"fmt"
	"github.com/golang/glog"
	//ipfsapi "github.com/ipfs/go-ipfs-api"

	ipfsapi "ipfs-share/ipfs"
	nw "ipfs-share/network"
	"github.com/pkg/errors"
	. "ipfs-share/collections"
	"ipfs-share/client/fs"
	"sync"
	"github.com/ethereum/go-ethereum/core/types"
)

type IUserFacade interface {
	CreateGroup(groupname string) error
	User() IUser
	Groups() []IGroupFacade
	SignOut()
	Transactions() ([]*types.Receipt, error)
}

type UserContext struct {
	user           IUser
	groups         *ConcurrentCollection
	addressBook    *ConcurrentCollection
	network        nw.INetwork
	ipfs           ipfsapi.IIpfs
	storage        *fs.Storage
	p2p            *P2PServer
	transactions   *ConcurrentCollection

	channelStop chan int
	lock sync.RWMutex
}

func NewUserContextFromSignUp(username, password, ethKeyPath, dataPath string, network nw.INetwork, ipfs ipfsapi.IIpfs, p2pPort string) (IUserFacade, error) {
	glog.Infof("[*] User '%s' signing up...", username)

	ipfsPeerId, err := ipfs.ID()
	if err != nil {
		return nil, errors.Wrap(err, "could not get ipfs peer id")
	}

	user, err := NewUser(username, password, ethKeyPath)
	if err != nil {
		return nil, errors.Wrap(err, "could not create new user")
	}

	// check if user is registered
	registered, err := network.IsUserRegistered(user.Address())
	if err != nil {
		return nil, errors.Wrap(err, "could not check if user is registered")
	}
	if registered {
		return nil, errors.New("user is already registered")
	}

	uc, err := NewUserContext(dataPath, user, network, ipfs, p2pPort)
	if err != nil {
		return nil, errors.Wrap(err, "could not create new user context")
	}

	boxer := user.Boxer()
	tx, err := network.RegisterUser(username, ipfsPeerId.ID, boxer.PublicKey.Value)
	if err != nil {
		return nil, errors.Wrap(err, "could not start transaction")
	}

	uc.transactions.Append(tx)

	return uc, nil
}

func NewUserContextFromSignIn(username, password, ethKeyPath, dataPath string, network nw.INetwork, ipfs ipfsapi.IIpfs, p2pPort string) (IUserFacade, error) {
	glog.Infof("[*] User '%s' signing in...", username)

	user, err := NewUser(username, password, ethKeyPath)
	if err != nil {
		return nil, errors.Wrap(err, "could not create new user")
	}

	registered, err := network.IsUserRegistered(user.Address())
	if err != nil {
		return nil, errors.Wrap(err, "could not check if user is registered")
	}
	if !registered {
		return nil, errors.New("user is not registered")
	}

	uc, err := NewUserContext(dataPath, user, network, ipfs, p2pPort)
	if err != nil {
		return nil, errors.Wrap(err, "could not create new user context")
	}

	glog.Info("[*] Signed in")

	return uc, nil
}

func NewUserContext(dataPath string, user IUser, network nw.INetwork, ipfs ipfsapi.IIpfs, p2pPort string) (*UserContext, error) {
	var err error
	var uc UserContext
	uc.user = user
	uc.network = network
	uc.ipfs = ipfs
	uc.storage = fs.NewStorage(dataPath)
	uc.groups = NewConcurrentCollection()
	uc.addressBook = NewConcurrentCollection()
	uc.transactions = NewConcurrentCollection()
	p2p, err := NewP2PConnection(p2pPort, &uc)
	if err != nil {
		return nil, errors.Wrap(err, "could not create P2P connection")
	}
	uc.p2p = p2p

	uc.channelStop = make(chan int)

	go HandleDebugEvents(network.GetDebugChannel())
	go HandleGroupInvitationEvents(network.GetGroupInvitationChannel(), uc.onGroupInvitationCallback)
	go HandleGroupUpdateIpfsEvents(network.GetGroupUpdateIpfsChannel(), uc.onGroupUpdateIpfsCallback)
	go HandleGroupRegisteredEvents(network.GetGroupRegisteredChannel(), uc.onGroupRegisteredCallback)

	if err := uc.BuildGroups(); err != nil {
		return nil, fmt.Errorf("could not build Groups: NewUserContext: %s", err)
	}

	return &uc, nil
}

func (ctx *UserContext) User() IUser {
	return ctx.user
}

func (ctx *UserContext) Save() error {
	//if err := ctx.Storage.SaveContextData(ctx); err != nil {
	//	return fmt.Errorf("could not save context data: %s", err)
	//}

	return nil
}

func (ctx *UserContext) SignOut() {
	fmt.Printf("[*] User '%s' signing out...\n", ctx.user.Name())
	for groupCtx := range ctx.groups.Iterator() {
		groupCtx.(*GroupContext).Stop()
	}

	ctx.network.Close()

	if err := ctx.Save(); err != nil {
		glog.Errorf("could not save context state: UserContext.SignOut: %s", err)
	}
}

func (ctx *UserContext) BuildGroups() error {
	glog.Infof("Building Groups for user '%s'...", ctx.user.Name())
	caps, err := ctx.storage.GetGroupCaps()
	if err != nil {
		return fmt.Errorf("[ERROR]: could not get group caps: UserContext.BuildGroups: %s", err)
	}
	for _, cap := range caps {
		groupCtx, err := NewGroupContextFromCAP(
			&cap,
			ctx.user,
			ctx.p2p,
			ctx.addressBook,
			ctx.network,
			ctx.ipfs,
			ctx.storage,
			ctx.transactions,
		)
		if err != nil {
			return fmt.Errorf("could not create new group context: UserContext.BuildGroups: %s", err)
		}
		if err := ctx.groups.Append(groupCtx); err != nil {
			glog.Warningf("could not append elem: %s", err)
		}
	}
	glog.Infof("Building Groups ended")
	return nil
}

func (ctx *UserContext) CreateGroup(groupname string) error {
	group := NewGroup(groupname)
	groupCtx, err := NewGroupContext(
		group,
		ctx.user,
		ctx.p2p,
		ctx.addressBook,
		ctx.network,
		ctx.ipfs,
		ctx.storage,
		ctx.transactions,
	)
	if err != nil {
		return fmt.Errorf("could not create new group context: UserContext.CreateGroup: %s", err)
	}

	boxer := groupCtx.Group.Boxer()
	ipfsHash := groupCtx.Repo.ipfsHash
	encIpfsHash := boxer.BoxSeal([]byte(ipfsHash))

	if err := group.SetIpfsHash(ipfsHash, encIpfsHash); err != nil {
		return errors.Wrap(err, "could not set ipfs hash of group")
	}

	if err := group.Save(ctx.storage); err != nil {
		return fmt.Errorf("could not save group: UserContext.CreateGroup: %s", err)
	}

	tx, err := groupCtx.Network.CreateGroup(
		group.Id().Data().([32]byte),
		group.Name(),
		group.EncryptedIpfsHash(),
	)
	if err != nil {
		return fmt.Errorf("could not register group: UserContext.CreateGroup: %s", err)
	}

	ctx.transactions.Append(tx)

	if err := ctx.groups.Append(groupCtx); err != nil {
		glog.Warningf("could not append elem: %s", err)
	}

	return nil
}

func (ctx *UserContext) Groups() []IGroupFacade {
	var groups []IGroupFacade

	for groupCtxInt := range ctx.groups.Iterator() {
		groups = append(groups, groupCtxInt.(IGroupFacade))
	}

	return groups
}

// Files lists the content of the user's repository
func (ctx *UserContext) List() map[string][]string {
	list := make(map[string][]string)
	return list
}

func (ctx *UserContext) Transactions() ([]*types.Receipt, error) {
	var list []*types.Receipt

	for txInt := range ctx.transactions.Iterator() {
		r, err := txInt.(*nw.Transaction).Receipt(ctx.network)
		if err != nil {
			return nil, errors.Wrap(err, "could not get tx receipt")
		}

		list = append(list, r)
	}

	return list, nil
}