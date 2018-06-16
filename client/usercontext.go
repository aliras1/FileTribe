package client

import (
	"bytes"
	"crypto/rand"
	"time"
	// "io/ioutil"
	"encoding/base64"
	// "encoding/json"
	"fmt"
	"path"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/glog"
	ipfsapi "github.com/ipfs/go-ipfs-api"
	"golang.org/x/crypto/sha3"

	fs "ipfs-share/client/filestorage"
	"ipfs-share/crypto"
	"ipfs-share/eth"

	nw "ipfs-share/networketh"
)

type Friend struct {
	FriendshipID [32]byte
	Contact      *nw.Contact
	MyDirectory  string
	HisDirectory string
}

type WaitingFriendRequest struct {
	Friend *Friend
	Digest [32]byte
	Signer crypto.Signer
}

func (wfr *WaitingFriendRequest) Confirm(ctx *UserContext) error {
	sig, err := wfr.Signer.Sign(wfr.Digest[:])
	if err != nil {
		return fmt.Errorf("could not sign digest: WaitingFriendRequest:Confirm: %s", err)
	}

	var dir [32]byte
	_, err = rand.Read(dir[:])
	if err != nil {
		return fmt.Errorf("could not read random dir: UserContext.AddFriend: %s", err)
	}
	dirEnc, err := crypto.AuthSeal(
		[]byte(base64.URLEncoding.EncodeToString(dir[:])),
		&wfr.Friend.Contact.Boxer,
		&ctx.User.Boxer.SecretKey,
	)
	if err != nil {
		return fmt.Errorf("could not encrypt directory name: UserContext.AddFriend: %s", err)
	}

	tx, err := ctx.Network.Session.ConfirmFriendship(wfr.Friend.FriendshipID, sig, dirEnc)
	if err != nil {
		return fmt.Errorf("could not confirm friendship: WaitingFriendRequest:Confirm: %s", err)
	}

	glog.Info("confirm cost: ", tx.Cost().String())

	return nil
}

type UserContext struct {
	User           *User // TODO lock boxer
	Groups         []*GroupContext
	Friends        []*Friend
	PendingFriends []*Friend
	WaitingFriends []*WaitingFriendRequest
	Repo           []*fs.FilePTP
	IPNSAddr       string
	Network        *nw.Network
	IPFS           *ipfsapi.Shell
	Storage        *fs.Storage // TODO lock

	channelStop chan int
}

func NewUserContextFromSignUp(username, password, ethKeyPath, dataPath string, network *nw.Network, ipfs *ipfsapi.Shell) (*UserContext, error) {
	fmt.Printf("[*] User '%s' signing up...\n", username)

	t := time.Now()
	ipfsID, err := ipfs.ID()
	if err != nil {
		return nil, fmt.Errorf("could not get ipfs id: NewUserContextFromSignUp: %s", err)
	}
	glog.Info("ipfs get id: ", time.Since(t))

	glog.Info("IPFS ID of ", username, " : ", ipfsID.ID)
	user, err := SignUp(username, password, ethKeyPath, ipfsID.ID, network)
	if err != nil {
		return nil, fmt.Errorf("could not sign up: NewUserContextFromSignUp: %s", err)
	}
	uc, err := NewUserContext(dataPath, user, network, ipfs)
	if err != nil {
		return nil, fmt.Errorf("could not create new user context: NewUserContextFromSignUp: %s", err)
	}

	fmt.Printf("[*] Signed in\n")

	return uc, nil
}

func NewUserContextFromSignIn(username, password, ethKeyPath, dataPath string, network *nw.Network, ipfs *ipfsapi.Shell) (*UserContext, error) {
	fmt.Printf("[*] User '%s' signing in...\n", username)

	user, err := SignIn(username, password, ethKeyPath, network)
	if err != nil {
		return nil, fmt.Errorf("could not sign in: NewUserContextFromSignIn: %s", err)
	}
	uc, err := NewUserContext(dataPath, user, network, ipfs)
	if err != nil {
		return nil, fmt.Errorf("could not create new user context: NewUserContextFromSignIn: %s", err)
	}

	fmt.Println("[*] Signed in")

	return uc, nil
}

func NewUserContext(dataPath string, user *User, network *nw.Network, ipfs *ipfsapi.Shell) (*UserContext, error) {
	var err error
	var uc UserContext
	uc.User = user
	uc.Network = network
	uc.IPFS = ipfs
	uc.Storage = fs.NewStorage(dataPath)
	uc.Storage.Init()
	uc.Groups = []*GroupContext{}
	uc.Repo, err = uc.Storage.BuildRepo(user.Name, network, ipfs)
	if err != nil {
		return nil, fmt.Errorf("could not build file repo: NewUserContext: %s", err)
	}
	uc.Friends = []*Friend{}
	uc.PendingFriends = []*Friend{}
	uc.WaitingFriends = []*WaitingFriendRequest{}

	ipfsID, err := ipfs.ID()
	if err != nil {
		return nil, fmt.Errorf("could not get ipfs id: NewUserContect: %s", err)
	}
	uc.IPNSAddr = ipfsID.ID

	uc.channelStop = make(chan int)
	go MessageProcessor(&uc)
	go NewFriendRequestHandler(&uc)
	go FriendshipConfirmedHandler(&uc)
	go DebugHandler(&uc)

	if err := uc.BuildGroups(); err != nil {
		return nil, fmt.Errorf("could not build groups: NewUserContext: %s", err)
	}
	if err := uc.Storage.PublishPublicDir(uc.IPFS); err != nil {
		return nil, fmt.Errorf("could not publish public dir: NewUserContext: %s", err)
	}
	return &uc, nil
}

func (uc *UserContext) SignOut() {
	fmt.Printf("[*] User '%s' signing out...\n", uc.User.Name)
	for _, groupCtx := range uc.Groups {
		groupCtx.Stop()
	}
	uc.channelStop <- 1
}

func (uc *UserContext) AddFriend(address common.Address) error {
	var r [32]byte
	_, err := rand.Read(r[:])
	if err != nil {
		return fmt.Errorf("could not read random: UserContext.AddFriend: %s", err)
	}
	addresses := append(r[:], uc.User.Address.Bytes()...)
	addresses = append(addresses, address.Bytes()...)

	id := sha3.Sum256(addresses)

	toUser, err := uc.Network.GetUser(address)
	if err != nil {
		return fmt.Errorf("could not retrieve user: UserContext.AddFriend: %s", err)
	}

	// sign then encrypt `from`
	fromDigest := sha3.Sum256(uc.User.Address.Bytes())
	fromSig, err := uc.User.Signer.Sign(fromDigest[:])
	if err != nil {
		return fmt.Errorf("could not sign from: UserContext.AddFriend: %s", err)
	}
	fromEnc, err := toUser.Boxer.Seal(append(fromSig[:64], uc.User.Address.Bytes()...))
	if err != nil {
		return fmt.Errorf("could not encrpyt from: UserContext.AddFriend: %s", err)
	}

	// sign then encrypt `to`
	toDigest := sha3.Sum256(address.Bytes())
	toSig, err := uc.User.Signer.Sign(toDigest[:])
	if err != nil {
		return fmt.Errorf("could not sign to: UserContext.AddFriend: %s", err)
	}
	toEnc, err := uc.User.Boxer.Seal(append(toSig[:64], address.Bytes()...))
	if err != nil {
		return fmt.Errorf("could not encrpyt to: UserContext.AddFriend: %s", err)
	}

	// generate signer
	_, err = rand.Read(r[:])
	if err != nil {
		return fmt.Errorf("could not read random digest: UserContext.AddFriend: %s", err)
	}

	key, err := ethcrypto.GenerateKey()
	if err != nil {
		return fmt.Errorf("could not generate key: UserContext:AddFriend: %s", err)
	}

	rawKey := ethcrypto.FromECDSA(key)
	keyEnc, err := toUser.Boxer.Seal(rawKey)
	if err != nil {
		return fmt.Errorf("could not enc key: UserContext.AddFriend: %s", err)
	}

	// share /ipns/id/for/to directory
	_, err = rand.Read(r[:])
	if err != nil {
		return fmt.Errorf("could not read random dir: UserContext.AddFriend: %s", err)
	}
	dirEnc, err := crypto.AuthSeal(
		[]byte(base64.URLEncoding.EncodeToString(r[:])),
		&toUser.Boxer,
		&uc.User.Boxer.SecretKey,
	)
	if err != nil {
		return fmt.Errorf("could not encrypt directory name: UserContext.AddFriend: %s", err)
	}

	_, ok := crypto.AuthOpen(dirEnc, &toUser.Boxer, &uc.User.Boxer.SecretKey)
	if !ok {
		return fmt.Errorf("bad orig enc")
	}
	glog.Info("direnc: ", len(dirEnc), " -> ", dirEnc)

	verifyAddress := ethcrypto.PubkeyToAddress(key.PublicKey)

	tx, err := uc.Network.Session.AddFriend(id, fromEnc, toEnc, keyEnc, dirEnc, r, verifyAddress)
	if err != nil {
		return fmt.Errorf("could not send transaction: UserContext.AddFriend: %s", err)
	}
	
	glog.Info("addfriend cost: ", tx.Cost().String())
	
	return nil
}

func DebugHandler(ctx *UserContext) {
	glog.Info("debug handleing...")
	for debug := range ctx.Network.DebugChannel {
		glog.Info("incoming address: ", debug.Addr.Bytes())
	}
}

func NewFriendRequestHandler(ctx *UserContext) {
	glog.Infof("User '%s' is processing new friend requests...\n", ctx.User.Name)
	fmt.Printf("new friend requests proc of '%s'\n", ctx.User.Name)

	for newFriendRequest := range ctx.Network.NewFriendRequestChannel {
		// sender wants to be friends with me
		from, err := ctx.User.Boxer.Open(newFriendRequest.From)
		if err == nil {
			fromSig := from[:64]
			fromBytes := from[64:]
			fromDigest := sha3.Sum256(fromBytes)
			address := common.BytesToAddress(fromBytes)
			contact, err := ctx.Network.GetUser(address)
			if err != nil {
				glog.Warningf("could not get user from: NewFriendRequestHandler: %s", err)
				continue
			}
			if ok := contact.VerifyKey.Verify(fromDigest[:], fromSig); !ok {
				glog.Warning("could not verify friend request from: NewFriendRequestHandler")
				continue
			}
			signerBytes, err := ctx.User.Boxer.Open(newFriendRequest.SigningKey)
			if err != nil {
				glog.Warningf("could not decrypt signer: %s", err)
				continue
			}
			signer, err := ethcrypto.ToECDSA(signerBytes)
			if err != nil {
				glog.Warningf("could not load signer: %s", err)
			}

			friend := &Friend{
				FriendshipID: newFriendRequest.Id,
				Contact:      contact,
			}
			waitingRequest := &WaitingFriendRequest{
				Friend: friend,
				Digest: newFriendRequest.Digest,
				Signer: crypto.Signer{
					PrivateKey: signer,
				},
			}

			ctx.WaitingFriends = append(ctx.WaitingFriends, waitingRequest)

			continue
		}

		to, err := ctx.User.Boxer.Open(newFriendRequest.To)
		if err == nil {
			toSig := to[:64]
			toBytes := to[64:]
			toDigest := sha3.Sum256(toBytes)
			address := common.BytesToAddress(toBytes)
			contact, err := ctx.Network.GetUser(address)
			if err != nil {
				glog.Warning("could not get user to: NewFriendRequestHandler")
				continue
			}

			vk := crypto.VerifyKey(ethcrypto.CompressPubkey(&ctx.User.Signer.PrivateKey.PublicKey))
			if ok := vk.Verify(toDigest[:], toSig); !ok {
				glog.Warning("could not verify friend request to: NewFriendRequestHandler")
			}

			friend := &Friend{
				FriendshipID: newFriendRequest.Id,
				Contact:      contact,
			}
			ctx.PendingFriends = append(ctx.PendingFriends, friend)
			continue
		}
	}
}

func FriendshipConfirmedHandler(ctx *UserContext) {
	glog.Infof("User '%s' is processing friendship confirmations...\n", ctx.User.Name)

	for confirmation := range ctx.Network.FriendshipConfirmedChannel {
		for _, pending := range ctx.PendingFriends {
			if bytes.Equal(confirmation.Id[:], pending.FriendshipID[:]) {
				myDir, ok := crypto.AuthOpen(confirmation.DirOfFromByTo, &pending.Contact.Boxer, &ctx.User.Boxer.SecretKey)
				if !ok {
					glog.Error("could not decrypt p my dir")
				}
				hisDir, ok := crypto.AuthOpen(confirmation.DirOfToByFrom, &pending.Contact.Boxer, &ctx.User.Boxer.SecretKey)
				if !ok {
					glog.Error("could not decrypt p his dir")
				}
				pending.MyDirectory = string(myDir)
				pending.HisDirectory = string(hisDir)

				ctx.Friends = append(ctx.Friends, pending)
				if err := ctx.Storage.MakeForDirectory(pending.HisDirectory, ctx.IPFS); err != nil {
					glog.Error("could not create for directory")
				}
				break
			}
		}

		for _, waiting := range ctx.WaitingFriends {
			if bytes.Equal(confirmation.Id[:], waiting.Friend.FriendshipID[:]) {
				myDir, ok := crypto.AuthOpen(
					confirmation.DirOfToByFrom,
					&waiting.Friend.Contact.Boxer,
					&ctx.User.Boxer.SecretKey,
				)
				if !ok {
					glog.Error("could not decrypt w my dir")
				}
				hisDir, ok := crypto.AuthOpen(
					confirmation.DirOfFromByTo,
					&waiting.Friend.Contact.Boxer,
					&ctx.User.Boxer.SecretKey,
				)
				if !ok {
					glog.Error("could not decrypt w his dir")
				}
				waiting.Friend.MyDirectory = string(myDir)
				waiting.Friend.HisDirectory = string(hisDir)

				ctx.Friends = append(ctx.Friends, waiting.Friend)
				break
			}
		}
	}
}

func unpackMessageSentEvent(event *eth.EthMessageSent, ctx *UserContext) (*nw.Message, error) {
	// var messageEnc []byte
	// for _, b := range event.Message {
	// 	messageEnc = append(messageEnc, b[0])
	// }

	// messageSigned, err := ctx.User.Boxer.Open(messageEnc)
	// if err != nil {
	// 	return nil, fmt.Errorf("could not decrypt event: MessageProcessor: %s", err)
	// }

	// messageJSON := messageSigned[64:]

	// var message nw.Message
	// if err := json.Unmarshal(messageJSON, &message); err != nil {
	// 	return nil, fmt.Errorf("could not unmarshal message: MessageProcessor: %s", err)
	// }

	// _, _, verifyKeyBytes, _, err := ctx.Network.GetUser(message.From)
	// if err != nil {
	// 	return nil, fmt.Errorf(
	// 		"could not get user '%s': MessageProcessor: %s",
	// 		base64.StdEncoding.EncodeToString(message.From[:]),
	// 		err,
	// 	)
	// }
	// verifyKey := crypto.VerifyKey(verifyKeyBytes[:])

	// return nil, fmt.Errorf("FIX IT")

	// ok := verifyKey.Verify(messageSigned, messageSigned)
	// if !ok {
	// 	return nil, fmt.Errorf("could not verify message: MessageProcessor")
	// }
	// return &message, nil
	return nil, nil
}

func MessageProcessor(ctx *UserContext) {
	glog.Infof("User '%s' is processing messages...\n", ctx.User.Name)
	fmt.Printf("msg proc of '%s'\n", ctx.User.Name)

	for messageSentEvent := range ctx.Network.MessageSentChannel {
		message, err := unpackMessageSentEvent(messageSentEvent, ctx)
		if err != nil {
			glog.Warning("error while processing 'MessageSentEvent': MessageProcessor: %s", err)
			continue
		}

		glog.Infof("--> user '%s' got message", ctx.User.Name)

		switch message.Type {
		case "PTP READ CAP":
			fmt.Println("[*] Downloading PTP file")
			cap, err := fs.DownloadReadCAP(message.From, ctx.User.Address, message.Payload, &ctx.User.Boxer, ctx.Storage, ctx.Network, ctx.IPFS)
			if err != nil {
				glog.Errorf("could not download read cap '%s' while 'PTP READ CAP' message: MessageProcessor: %s\n", message.Payload, err)
				continue
			}
			if ctx.isFileInRepo(cap.IPNSPath) {
				glog.Errorf("file '%s' is already in the repo", cap.IPNSPath)
				continue
			}
			file, err := fs.NewFileFromCAP(cap, ctx.Storage, ctx.IPFS)
			if err != nil {
				glog.Errorf("could not instantiate a new file from cap '%s' while 'PTP READ CAP' message: MessageProcessor: %s\n", cap.FileName, err)
				continue
			}
			ctx.addFileToRepo(file)

			fmt.Println("[*] PTP file downloaded")
		case "GROUP INVITE":
			fmt.Printf("[*] Joining group '%s'...\n", strings.Split(message.Payload, ".")[0])
			glog.Infof("\t--> Donwloading group access cap '%s'...", message.Payload)

			cap, err := ctx.Storage.DownloadGroupAccessCAP(message.From, ctx.User.Address, message.Payload, &ctx.User.Boxer, ctx.Network, ctx.IPFS)
			if err != nil {
				glog.Errorf("could not download ga cap: MessageProcessor: %s", err)
				continue
			}
			glog.Infof("\t<-- Donwloaded")
			if err := cap.Store(ctx.Storage); err != nil {
				glog.Errorf("could not store ga cap: MessageProcessor: %s", err)
				continue
			}
			cap.Boxer.RNG = rand.Reader
			groupCtx, err := NewGroupContextFromCAP(cap, ctx.User, ctx.Network, ctx.IPFS, ctx.Storage)
			if err != nil {
				glog.Errorf("could not create group context from cap: MessageProcessor: %s", err)
				continue
			}
			ctx.Groups = append(ctx.Groups, groupCtx)
			fmt.Printf("[*] Joined group '%s'\n", groupCtx.Group.Name)
		}
	}
}

func (uc *UserContext) BuildGroups() error {
	glog.Infof("Building groups for user '%s'...", uc.User.Name)
	caps, err := uc.Storage.GetGroupCAPs()
	if err != nil {
		return fmt.Errorf("[ERROR]: could not get group caps: UserContext.BuildGroups: %s", err)
	}
	for _, cap := range caps {
		groupCtx, err := NewGroupContextFromCAP(&cap, uc.User, uc.Network, uc.IPFS, uc.Storage)
		if err != nil {
			return fmt.Errorf("could not create new group context: UserContext.BuildGroups: %s", err)
		}
		uc.Groups = append(uc.Groups, groupCtx)
	}
	glog.Infof("Building groups ended")
	return nil
}

func (uc *UserContext) CreateGroup(groupname string) error {
	// group := NewGroup(groupname)
	// if err := group.Register(uc.User.Name, uc.Network); err != nil {
	// 	return fmt.Errorf("could not register group: UserContext.CreateGroup: %s", err)
	// }
	// if err := group.Save(uc.Storage); err != nil {
	// 	return fmt.Errorf("could not save group: UserContext.CreateGroup: %s", err)
	// }
	// groupCtx, err := NewGroupContext(group, uc.User, uc.Network, uc.IPFS, uc.Storage)
	// if err != nil {
	// 	return fmt.Errorf("could not create new group context: UserContext.CreateGroup: %s", err)
	// }
	// uc.Groups = append(uc.Groups, groupCtx)

	// fmt.Printf("[*] Group '%s' created!\n", groupname)

	return nil
}

func (uc *UserContext) AddAndShareFile(filePath string, shareWith []string) error {
	// fmt.Printf("[*] PTP sharing file '%s' with user '%s'...\n", filePath, shareWith[0])

	// if uc.isFileInRepo("/ipns/" + uc.IPNSAddr + "/files/" + path.Base(filePath)) {
	// 	return fmt.Errorf("file is already added to the repo: UserContext.AddAndShareFile")
	// }
	// file, err := fs.NewSharedFilePTP(filePath, uc.User.Name, uc.Storage, uc.IPFS)
	// if err != nil {
	// 	return fmt.Errorf("could not create new shared file '%s': UserContext.AddAndShareFile: %s", file.Name, err)
	// }
	// if err := file.Share(shareWith, &uc.User.Boxer, uc.Storage, uc.Network, uc.IPFS); err != nil {
	// 	return fmt.Errorf("could not share file '%s': UserContext.AddAndShareFile: %s", file.Name, err)
	// }
	// uc.addFileToRepo(file)

	// fmt.Println("[*] PTP sharing ended")

	return nil
}

func (uc *UserContext) isFileInRepo(ipnsPath string) bool {
	for _, file := range uc.Repo {
		if strings.Compare(file.IPNSPath, ipnsPath) == 0 {
			return true
		}
	}
	return false
}

func (uc *UserContext) addFileToRepo(file *fs.FilePTP) {
	uc.Repo = append(uc.Repo, file)
}

func (uc *UserContext) getFileFromRepo(name string) *fs.FilePTP {
	for _, file := range uc.Repo {
		if strings.Compare(path.Base(file.Name), name) == 0 {
			return file
		}
	}
	return nil
}

func (uc *UserContext) List() string {
	str := "MyFilen"
	for _, file := range uc.Repo {
		str += "\t--> " + file.Name + "\n"
	}
	for _, groupCtx := range uc.Groups {
		str += groupCtx.Group.Name
		str += groupCtx.Repo.List()
	}
	return str
}
