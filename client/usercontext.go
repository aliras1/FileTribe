package client

import (
	"bytes"
	"fmt"
	"strings"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/glog"
	//ipfsapi "github.com/ipfs/go-ipfs-api"

	ipfsapi "ipfs-share/ipfs"
	"ipfs-share/crypto"
	"ipfs-share/eth"
	nw "ipfs-share/networketh"
	"github.com/pkg/errors"
	. "ipfs-share/collections"
)

type Friend struct {
	Contact      *nw.Contact
	MyDirectory  string
	HisDirectory string
	From         bool
}

type WaitingFriendRequest struct {
	Friend *Friend
	Digest [32]byte
	Signer crypto.Signer
}

func (wfr *WaitingFriendRequest) Confirm(ctx *UserContext) error {
	//ipfsHash, err := groupCtx.Storage.MakeForDirectory(wfr.Friend.Contact.Address.String(), groupCtx.Ipfs)
	//if err != nil {
	//	return fmt.Errorf("could not create for directory: WaitingFriendRequest:Confirm: %s", err)
	//}
	//
	//message, err := NewFriendConfirmation(groupCtx.User.Address, wfr.Friend.Contact.Address, ipfsHash, groupCtx.User.Signer)
	//if err != nil {
	//	return fmt.Errorf("could not create new friend confirmation: WaitingFriendRequest:Confirm: %s", err)
	//}
	//
	//if err := message.DialP2PConn(&wfr.Friend.Contact.Boxer, &groupCtx.User.Boxer.PublicKey, groupCtx.User.Address, groupCtx.Network); err != nil {
	//	return fmt.Errorf("could not send message confirmation: WaitingFriendRequest:Confirm: %s", err)
	//}

	return nil
}

type UserContext struct {
	User           *User // TODO lock boxer
	Groups         *ConcurrentCollection
	Friends        []*Friend
	PendingFriends []*Friend
	WaitingFriends []*WaitingFriendRequest
	Repo           map[ethcommon.Address]map[[32]byte]*PTPFile
	IPNSAddr       string
	AddressBook    *ConcurrentCollection
	Network        nw.INetwork
	Ipfs           ipfsapi.IIpfs
	Storage        *Storage // TODO lock
	LatestBlock    uint64
	P2P            *P2PServer

	channelStop chan int
}

func NewUserContextFromSignUp(username, password, ethKeyPath, dataPath string, network nw.INetwork, ipfs ipfsapi.IIpfs, p2pPort string) (*UserContext, error) {
	glog.Infof("[*] User '%s' signing up...", username)

	ipfsPeerId, err := ipfs.ID()
	if err != nil {
		return nil, errors.Wrap(err, "could not get ipfs peer id")
	}

	user, err := SignUp(username, password, ipfsPeerId.ID, ethKeyPath, network)
	if err != nil {
		return nil, fmt.Errorf("could not sign up: NewUserContextFromSignUp: %s", err)
	}
	uc, err := NewUserContext(dataPath, user, network, ipfs, p2pPort)
	if err != nil {
		return nil, fmt.Errorf("could not create new user context: NewUserContextFromSignUp: %s", err)
	}

	glog.Info("[*] Signed in")

	return uc, nil
}

func NewUserContextFromSignIn(username, password, ethKeyPath, dataPath string, network nw.INetwork, ipfs ipfsapi.IIpfs, p2pPort string) (*UserContext, error) {
	glog.Infof("[*] User '%s' signing in...", username)

	user, err := SignIn(username, password, ethKeyPath, network)
	if err != nil {
		return nil, fmt.Errorf("could not sign in: NewUserContextFromSignIn: %s", err)
	}
	uc, err := NewUserContext(dataPath, user, network, ipfs, p2pPort)
	if err != nil {
		return nil, fmt.Errorf("could not create new user context: NewUserContextFromSignIn: %s", err)
	}

	glog.Info("[*] Signed in")

	return uc, nil
}

func NewUserContext(dataPath string, user *User, network nw.INetwork, ipfs ipfsapi.IIpfs, p2pPort string) (*UserContext, error) {
	var err error
	var uc UserContext
	uc.User = user
	uc.Network = network
	uc.Ipfs = ipfs
	uc.Storage = NewStorage(dataPath)
	uc.Groups = NewConcurrentCollection()
	uc.AddressBook = NewConcurrentCollection()
	p2p, err := NewP2PConnection(p2pPort, &uc)
	if err != nil {
		return nil, errors.Wrap(err, "could not create P2P connection")
	}
	uc.P2P = p2p

	if err := uc.Storage.LoadContextData(&uc); err != nil {
		glog.Warningf("could not load context data: %s", err)

		uc.Friends = []*Friend{}
		uc.PendingFriends = []*Friend{}
		uc.WaitingFriends = []*WaitingFriendRequest{}
		uc.LatestBlock = 0
	}

	uc.Repo, err = uc.Storage.BuildRepo(&uc)
	if err != nil {
		return nil, fmt.Errorf("could not build file repo: NewUserContext: %s", err)
	}

	ipfsID, err := ipfs.ID()
	if err != nil {
		return nil, fmt.Errorf("could not get ipfs sessionId: NewUserContect: %s", err)
	}
	uc.IPNSAddr = ipfsID.ID

	uc.channelStop = make(chan int)

	//uc.SynchWithChain()
	//if err := uc.Network.SubscribeToEvents(uc.LatestBlock); err != nil {
	//	return nil, fmt.Errorf("could not subscrite to events: NewUserContext: %s", err)
	//}

	//go HandleMessageEventChannel(&uc)
	go HandleDebugEvents(&uc)
	go HandleGroupInvitationEvents(&uc)
	go HandleGroupUpdateIpfsEvents(&uc)

	if err := uc.BuildGroups(); err != nil {
		return nil, fmt.Errorf("could not build Groups: NewUserContext: %s", err)
	}

	return &uc, nil
}

func (ctx *UserContext) SynchWithChain() error {
	//if uc.LatestBlock != 0 {
	//	iterator, err := uc.Network.FilterMessageEvents(uc.LatestBlock)
	//	if err != nil {
	//		return fmt.Errorf("could not get message event iterator: NewUserContext: %s", err)
	//	}
	//
	//	event := iterator.Event
	//	for iterator.Next() {
	//		if event != nil {
	//			uc.LatestBlock = event.Raw.BlockNumber
	//			MessageEventProcessor(event, uc)
	//		}
	//		event = iterator.Event
	//	}
	//}

	return nil
}

func (ctx *UserContext) SaveState() error {
	if err := ctx.Storage.SaveContextData(ctx); err != nil {
		return fmt.Errorf("could not save context data: %s", err)
	}

	return nil
}

func (ctx *UserContext) SignOut() {
	fmt.Printf("[*] User '%s' signing out...\n", ctx.User.Name)
	for groupCtx := range ctx.Groups.Iterator() {
		groupCtx.(*GroupContext).Stop()
	}

	//uc.Network.debugSub.Unsubscribe()
	//uc.Network.MessageSentSub.Unsubscribe()

	if err := ctx.SaveState(); err != nil {
		glog.Errorf("could not save context state: UserContext.SignOut: %s", err)
	}
}

func (ctx *UserContext) AddFriend(address ethcommon.Address) error {
	// get user
	//toUser, err := uc.Network.GetUser(address)
	//if err != nil {
	//	return fmt.Errorf("could not retrieve user: UserContext.AddFriend: %s", err)
	//}
	//
	//// create /ipfs/for_to directory name
	//ipfsHash, err := uc.Storage.MakeForDirectory(toUser.Address.String(), uc.Ipfs)
	//if err != nil {
	//	return fmt.Errorf("could not create for directory: UserContext.AddFriend: %s", err)
	//}
	//
	//// create and send message
	//message, err := NewFriendRequest(uc.User.Address, address, ipfsHash, uc.User.Signer)
	//if err != nil {
	//	return fmt.Errorf("could not create message: UserContext.AddFriend: %s", err)
	//}
	//
	//if err := message.DialP2PConn(&toUser.Boxer, &uc.User.Boxer.PublicKey, uc.User.Address, uc.Network); err != nil {
	//	return fmt.Errorf("could not send message: UserContext.AddFriend: %s", err)
	//}

	return nil
}

func unpackMessageSentEvent(event *eth.EthMessageSent, ctx *UserContext) (*nw.Message, error) {
	// var messageEnc []byte
	// for _, b := range event.Message {
	// 	messageEnc = append(messageEnc, b[0])
	// }

	// messageSigned, err := groupCtx.User.Boxer.Open(messageEnc)
	// if err != nil {
	// 	return nil, fmt.Errorf("could not decrypt event: connectionListener: %s", err)
	// }

	// messageJSON := messageSigned[64:]

	// var message nw.Message
	// if err := json.Unmarshal(messageJSON, &message); err != nil {
	// 	return nil, fmt.Errorf("could not unmarshal message: connectionListener: %s", err)
	// }

	// _, _, verifyKeyBytes, _, err := groupCtx.Network.GetUser(message.From)
	// if err != nil {
	// 	return nil, fmt.Errorf(
	// 		"could not get user '%s': connectionListener: %s",
	// 		base64.StdEncoding.EncodeToString(message.From[:]),
	// 		err,
	// 	)
	// }
	// verifyKey := crypto.VerifyKey(verifyKeyBytes[:])

	// return nil, fmt.Errorf("FIX IT")

	// ok := verifyKey.Verify(messageSigned, messageSigned)
	// if !ok {
	// 	return nil, fmt.Errorf("could not verify message: connectionListener")
	// }
	// return &message, nil
	return nil, nil
}

func (ctx *UserContext) BuildGroups() error {
	glog.Infof("Building Groups for user '%s'...", ctx.User.Name)
	caps, err := ctx.Storage.GetGroupCAPs()
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
	if err := group.CreateOnNetwork(ctx.User.Name, ctx.Network); err != nil {
		return fmt.Errorf("could not register group: UserContext.CreateGroup: %s", err)
	}
	if err := group.Save(ctx.Storage); err != nil {
		return fmt.Errorf("could not save group: UserContext.CreateGroup: %s", err)
	}
	groupCtx, err := NewGroupContext(
		group,
		ctx.User,
		ctx.P2P,
		ctx.AddressBook,
		ctx.Network,
		ctx.Ipfs,
		ctx.Storage,
	)
	if err != nil {
		return fmt.Errorf("could not create new group context: UserContext.CreateGroup: %s", err)
	}
	if err := ctx.Groups.Append(groupCtx); err != nil {
		glog.Warningf("could not append elem: %s", err)
	}

	return nil
}

func (ctx *UserContext) AddFile(filePath string) ([32]byte, error) {
	file, err := NewPTPFile(filePath, ctx)
	if err != nil {
		return [32]byte{}, fmt.Errorf("could not create new file ptp'%s': UserContext.AddFile: %s", filePath, err)
	}

	// file is user's, add to under his address
	ctx.addFileToRepo(ctx.User.Address, file)

	return file.CAP.ID, nil
}

func (ctx *UserContext) isFileInRepo(ipfsHash string) bool {
	for _, friendFiles := range ctx.Repo {
		for _, file := range friendFiles {
			if strings.Compare(file.CAP.IPFSHash, ipfsHash) == 0 {
				return true
			}
		}
	}
	return false
}

func (ctx *UserContext) addFileToRepo(address ethcommon.Address, file *PTPFile) {
	ctx.Repo[address][file.CAP.ID] = file
}

func (ctx *UserContext) getFileFromRepo(id [32]byte) *PTPFile {
	for _, friendFiles := range ctx.Repo {
		for _, file := range friendFiles {
			if bytes.Equal(file.CAP.ID[:], id[:]) {
				return file
			}
		}
	}
	return nil
}

// List lists the content of the user's repository
func (ctx *UserContext) List() map[string][]string {
	list := make(map[string][]string)

	for address, files := range ctx.Repo {
		fileList := []string{}

		glog.Info(address.String())

		for _, file := range files {
			glog.Info("\t--> " + file.CAP.FileName)
			fileList = append(fileList, file.CAP.FileName)
		}

		list[address.String()] = fileList
	}

	return list
}
