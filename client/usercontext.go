package client

import (
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/glog"
	"github.com/pkg/errors"

	com "ipfs-share/client/communication"
	"ipfs-share/client/fs"
	"ipfs-share/client/interfaces"
	. "ipfs-share/collections"
	ipfsapi "ipfs-share/ipfs"
	nw "ipfs-share/network"
)

type IUserFacade interface {
	CreateGroup(groupname string) error
	User() interfaces.IUser
	Groups() []IGroupFacade
	SignOut()
	Transactions() ([]*types.Receipt, error)
}

type UserContext struct {
	user           interfaces.IUser
	groups         *ConcurrentCollection
	addressBook    *ConcurrentCollection
	network        nw.INetwork
	ipfs           ipfsapi.IIpfs
	storage        *fs.Storage
	p2p            *com.P2PServer
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

func NewUserContext(dataPath string, user interfaces.IUser, network nw.INetwork, ipfs ipfsapi.IIpfs, p2pPort string) (*UserContext, error) {
	var err error
	var ctx UserContext
	ctx.user = user
	ctx.network = network
	ctx.ipfs = ipfs
	ctx.storage = fs.NewStorage(dataPath)
	ctx.groups = NewConcurrentCollection()
	ctx.addressBook = NewConcurrentCollection()
	ctx.transactions = NewConcurrentCollection()
	p2p, err := com.NewP2PConnection(
		p2pPort,
		user,
		ctx.addressBook,
		ctx.GetGroup,
		ipfs,
		network)
	if err != nil {
		return nil, errors.Wrap(err, "could not create P2P connection")
	}
	ctx.p2p = p2p

	ctx.channelStop = make(chan int)

	go ctx.HandleDebugEvents(network.GetDebugChannel())
	go ctx.HandleGroupInvitationEvents(network.GetGroupInvitationChannel())
	go ctx.HandleGroupUpdateIpfsEvents(network.GetGroupUpdateIpfsChannel())
	go ctx.HandleGroupRegisteredEvents(network.GetGroupRegisteredChannel())

	if err := ctx.BuildGroups(); err != nil {
		return nil, errors.Wrap(err, "could not build groups")
	}

	return &ctx, nil
}

func (ctx *UserContext) GetGroup(id [32]byte) interfaces.IGroup {
	groupCtxInt := ctx.groups.Get(NewBytesId(id))
	if (groupCtxInt == nil) {
		return nil
	}

	groupCtx := groupCtxInt.(*GroupContext)

	return groupCtx.Group
}

func (ctx *UserContext) User() interfaces.IUser {
	return ctx.user
}

func (ctx *UserContext) Save() error {
	//if err := ctx.Storage.SaveContextData(ctx); err != nil {
	//	return fmt.Errorf("could not save context data: %s", err)
	//}

	return nil
}

func (ctx *UserContext) SignOut() {
	glog.Infof("[*] User '%s' signing out...\n", ctx.user.Name())
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
		return errors.Wrap(err, "could not get group caps")
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
			return errors.Wrap(err, "could not create new group context")
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
		return errors.Wrap(err, "could not create new group context")
	}

	boxer := groupCtx.Group.Boxer()
	ipfsHash := groupCtx.Repo.IpfsHash()
	encIpfsHash := boxer.BoxSeal([]byte(ipfsHash))

	if err := group.SetIpfsHash(ipfsHash, encIpfsHash); err != nil {
		return errors.Wrap(err, "could not set ipfs hash of group")
	}

	if err := groupCtx.Save(); err != nil {
		return errors.Wrap(err, "could not save group")
	}

	tx, err := groupCtx.Network.CreateGroup(
		group.Id().Data().([32]byte),
		group.Name(),
		group.EncryptedIpfsHash(),
	)
	if err != nil {
		return errors.Wrap(err, "could not register group")
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