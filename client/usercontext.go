package client

import (
	"fmt"
	"github.com/golang/glog"
	//ipfsapi "github.com/ipfs/go-ipfs-api"

	ipfsapi "ipfs-share/ipfs"
	nw "ipfs-share/network"
	"github.com/pkg/errors"
	. "ipfs-share/collections"
)


type UserContext struct {
	User           IUser
	Groups         *ConcurrentCollection
	IPNSAddr       string
	AddressBook    *ConcurrentCollection
	Network        nw.INetwork
	Ipfs           ipfsapi.IIpfs
	Storage        *Storage
	LatestBlock    uint64
	P2P            *P2PServer
	Transactions   *ConcurrentCollection

	channelStop chan int
}

func NewUserContextFromSignUp(username, password, ethKeyPath, dataPath string, network nw.INetwork, ipfs ipfsapi.IIpfs, p2pPort string) (*UserContext, error) {
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

	uc.Transactions.Append(tx)

	return uc, nil
}

func NewUserContextFromSignIn(username, password, ethKeyPath, dataPath string, network nw.INetwork, ipfs ipfsapi.IIpfs, p2pPort string) (*UserContext, error) {
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
	uc.User = user
	uc.Network = network
	uc.Ipfs = ipfs
	uc.Storage = NewStorage(dataPath)
	uc.Groups = NewConcurrentCollection()
	uc.AddressBook = NewConcurrentCollection()
	uc.Transactions = NewConcurrentCollection()
	p2p, err := NewP2PConnection(p2pPort, &uc)
	if err != nil {
		return nil, errors.Wrap(err, "could not create P2P connection")
	}
	uc.P2P = p2p

	//if err := uc.Storage.LoadContextData(&uc); err != nil {
	//	glog.Warningf("could not load context data: %s", err)
	//}

	ipfsID, err := ipfs.ID()
	if err != nil {
		return nil, fmt.Errorf("could not get ipfs sessionId: NewUserContect: %s", err)
	}
	uc.IPNSAddr = ipfsID.ID

	uc.channelStop = make(chan int)

	go HandleDebugEvents(&uc)
	go HandleGroupInvitationEvents(&uc)
	go HandleGroupUpdateIpfsEvents(&uc)
	go HandleGroupRegisteredEvents(&uc)

	if err := uc.BuildGroups(); err != nil {
		return nil, fmt.Errorf("could not build Groups: NewUserContext: %s", err)
	}

	return &uc, nil
}


func (ctx *UserContext) Save() error {
	//if err := ctx.Storage.SaveContextData(ctx); err != nil {
	//	return fmt.Errorf("could not save context data: %s", err)
	//}

	return nil
}

func (ctx *UserContext) SignOut() {
	fmt.Printf("[*] User '%s' signing out...\n", ctx.User.Name)
	for groupCtx := range ctx.Groups.Iterator() {
		groupCtx.(*GroupContext).Stop()
	}

	ctx.Network.Close()

	if err := ctx.Save(); err != nil {
		glog.Errorf("could not save context state: UserContext.SignOut: %s", err)
	}
}

func (ctx *UserContext) BuildGroups() error {
	glog.Infof("Building Groups for user '%s'...", ctx.User.Name)
	caps, err := ctx.Storage.GetGroupCaps()
	if err != nil {
		return fmt.Errorf("[ERROR]: could not get group caps: UserContext.BuildGroups: %s", err)
	}
	for _, cap := range caps {
		groupCtx, err := NewGroupContextFromCAP(
			&cap,
			ctx.User,
			ctx.P2P,
			ctx.AddressBook,
			ctx.Network,
			ctx.Ipfs,
			ctx.Storage,
			ctx.Transactions,
		)
		if err != nil {
			return fmt.Errorf("could not create new group context: UserContext.BuildGroups: %s", err)
		}
		if err := ctx.Groups.Append(groupCtx); err != nil {
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
		ctx.User,
		ctx.P2P,
		ctx.AddressBook,
		ctx.Network,
		ctx.Ipfs,
		ctx.Storage,
		ctx.Transactions,
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

	if err := group.Save(ctx.Storage); err != nil {
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

	ctx.Transactions.Append(tx)

	if err := ctx.Groups.Append(groupCtx); err != nil {
		glog.Warningf("could not append elem: %s", err)
	}

	return nil
}

// List lists the content of the user's repository
func (ctx *UserContext) List() map[string][]string {
	list := make(map[string][]string)
	return list
}
