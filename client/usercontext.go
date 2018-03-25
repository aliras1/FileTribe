package client

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	fs "ipfs-share/client/filestorage"
	"ipfs-share/ipfs"
	nw "ipfs-share/network"
)

type UserContext struct {
	User        *User // TODO lock boxer
	Repo        []*fs.File
	Network     *nw.Network
	IPFS        *ipfs.IPFS
	UserStorage *fs.UserStorage // TODO lock

	channelMsg chan nw.Message
	channelSig chan os.Signal
}

func NewUserContextFromSignUp(username, password, dataPath string, network *nw.Network, ipfs *ipfs.IPFS) (*UserContext, error) {
	ipfsID, err := ipfs.ID()
	if err != nil {
		return nil, err
	}
	user, err := SignUp(username, password, ipfsID.ID, network)
	if err != nil {
		return nil, err
	}
	return NewUserContext(dataPath, user, network, ipfs), nil
}

func NewUserContextFromSignIn(username, password, dataPath string, network *nw.Network, ipfs *ipfs.IPFS) (*UserContext, error) {
	user, err := SignIn(username, password, network)
	if err != nil {
		return nil, err
	}
	return NewUserContext(dataPath, user, network, ipfs), nil
}

func NewUserContext(dataPath string, user *User, network *nw.Network, ipfs *ipfs.IPFS) *UserContext {
	var err error
	var uc UserContext
	uc.User = user
	uc.Network = network
	uc.IPFS = ipfs
	uc.UserStorage = fs.NewUserStorage(dataPath)
	uc.Repo, err = uc.UserStorage.BuildRepo(ipfs)
	if err != nil {
		log.Println(err)
	}

	uc.channelMsg = make(chan nw.Message)
	uc.channelSig = make(chan os.Signal)
	go MessageGetter(uc.User.Username, network, uc.channelMsg, uc.channelSig)
	go MessageProcessor(uc.channelMsg, uc.User.Username, &uc)

	return &uc
}

func MessageGetter(username string, network *nw.Network, channelMsg chan nw.Message, channelSig chan os.Signal) {
	for true {
		select {
		case sig := <-channelSig:
			fmt.Println("received signal", sig)
			close(channelMsg)
			close(channelSig)
			return
		default:
			msgs, err := network.GetMessages(username)
			if err != nil {
				fmt.Println(err.Error())
				break
			}
			for _, msg := range msgs {
				channelMsg <- *msg
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func MessageProcessor(channelMsg chan nw.Message, username string, ctx *UserContext) {
	fmt.Println("msg processing...")
	for msg := range channelMsg {
		fmt.Print("msgproc: ")
		fmt.Println(msg)

		cap, err := fs.DownloadCAP(msg.From, username, msg.Message, &ctx.User.Boxer, ctx.UserStorage, ctx.Network, ctx.IPFS)
		if err != nil {
			log.Println(err)
			continue
		}
		file, err := fs.NewFileFromCAP(cap, ctx.UserStorage, ctx.IPFS)
		if err != nil {
			log.Println(err)
			continue
		}
		ctx.addFileToRepo(file)

		fmt.Println("content of root directory: ")
		ctx.List()
	}
}

func SignUp(username, password, ipfsAddr string, network *nw.Network) (*User, error) {
	exists, err := network.IsUsernameRegistered(username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("user already exists")
	}
	user := NewUser(username, password)
	if user == nil {
		return nil, errors.New("could not generate user")
	}
	err = network.RegisterUsername(username, user.PublicKeyHash)
	if err != nil {
		return nil, err
	}
	network.PutSigningKey(user.PublicKeyHash, user.Signer.PublicKey)
	network.PutBoxingKey(user.PublicKeyHash, user.Boxer.PublicKey)
	network.PutIPFSAddr(user.PublicKeyHash, ipfsAddr)
	return user, nil
}

func SignIn(username, password string, network *nw.Network) (*User, error) {
	exists, err := network.IsUsernameRegistered(username)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("username does not exists")
	}
	user := NewUser(username, password)
	if user == nil {
		return nil, errors.New("could not generate user")
	}
	publicKeyHash, err := network.GetUserPublicKeyHash(username)
	if err != nil {
		return nil, err
	}
	if !publicKeyHash.Equals(&user.PublicKeyHash) {
		return nil, errors.New("incorrect password")
	}
	return user, nil
}

func (uc *UserContext) AddAndShareFile(filePath string, shareWith []string) error {
	if uc.isFileInRepo(filePath) {
		return errors.New("file already in root dir")
	}
	file, err := fs.NewSharedFile(filePath, uc.User.Username, uc.UserStorage, uc.IPFS)
	if err != nil {
		return err
	}
	err = file.Share(shareWith, &uc.User.Boxer, uc.UserStorage, uc.Network, uc.IPFS)
	if err != nil {
		return err
	}
	uc.addFileToRepo(file)
	return nil
}

func (uc *UserContext) isFileInRepo(filePath string) bool {
	for _, i := range uc.Repo {
		if strings.Compare(path.Base(i.Path), path.Base(filePath)) == 0 {
			return true
		}
	}
	return false
}

func (uc *UserContext) addFileToRepo(file *fs.File) {
	uc.Repo = append(uc.Repo, file)
}

func (uc *UserContext) getFileFromRepo(name string) *fs.File {
	for _, file := range uc.Repo {
		if strings.Compare(path.Base(file.Path), name) == 0 {
			return file
		}
	}
	return nil
}

func (uc *UserContext) List() {
	fmt.Println(uc.User.Username)
	for _, f := range uc.Repo {
		fmt.Print("\t--> ")
		fmt.Println(*f)
	}
}
