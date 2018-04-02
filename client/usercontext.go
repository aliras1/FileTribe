package client

import (
	"crypto/rand"
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
	User    *User // TODO lock boxer
	Groups  []*GroupContext
	Repo    []*fs.File
	Network *nw.Network
	IPFS    *ipfs.IPFS
	Storage *fs.Storage // TODO lock

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
	uc.Storage = fs.NewStorage(dataPath)
	uc.Groups = []*GroupContext{}
	uc.Repo, err = uc.Storage.BuildRepo(ipfs)
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
				// TODO validate message signature
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
		switch msg.Type {
		case "PTP READ CAP":
			cap, err := ctx.Storage.DownloadReadCAP(msg.From, username, msg.Message, &ctx.User.Boxer, ctx.Network, ctx.IPFS)
			if err != nil {
				log.Println(err)
				continue
			}
			file, err := fs.NewFileFromCAP(cap, ctx.Storage, ctx.IPFS)
			if err != nil {
				log.Println(err)
				continue
			}
			ctx.addFileToRepo(file)

			fmt.Println("content of root directory: ")
			ctx.List()
		case "GROUP INVITE":
			fmt.Println("group invite....")
			cap, err := ctx.Storage.DownloadGroupAccessCAP(msg.From, username, msg.Message, &ctx.User.Boxer, ctx.Network, ctx.IPFS)
			if err != nil {
				log.Println(err)
				continue
			}
			cap.Boxer.RNG = rand.Reader
			err = ctx.NewGroupFromCAP(msg.From, cap)
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}
}

func (uc *UserContext) CreateGroup(groupName string) error {
	group := NewGroup(groupName)
	err := group.Register(uc.Network)
	if err != nil {
		return err
	}
	uc.Storage.CreateGroupStorage(groupName)
	groupCtx := GroupContext{uc.User, group, nil,
		&MemberList{[]Member{{uc.User.Username, nil}}},
		&ActiveMemberList{}, uc.Network, uc.IPFS, uc.Storage}

	NewSynchronizer(uc.User.Username, &uc.User.Signer, &groupCtx)
	uc.Groups = append(uc.Groups, &groupCtx)
	return groupCtx.Save()
}

func (uc *UserContext) NewGroupFromCAP(from string, cap *fs.GroupAccessCAP) error {
	group := &Group{cap.GroupName, cap.Boxer}
	uc.Storage.CreateGroupStorage(group.GroupName)
	groupCtx := GroupContext{uc.User, group, nil,
		&MemberList{[]Member{{uc.User.Username, nil}}},
		&ActiveMemberList{}, uc.Network, uc.IPFS, uc.Storage}

	NewSynchronizer(uc.User.Username, &uc.User.Signer, &groupCtx)
	uc.Groups = append(uc.Groups, &groupCtx)
	fmt.Println("pulling g data")
	err := groupCtx.PullGroupData(from)
	if err != nil {
		return err
	}
	fmt.Print("members: ")
	fmt.Println(groupCtx.Members)
	return groupCtx.Save()
}

func (uc *UserContext) AddAndShareFile(filePath string, shareWith []string) error {
	if uc.isFileInRepo(filePath) {
		return errors.New("file already in root dir")
	}
	file, err := fs.NewSharedFile(filePath, uc.User.Username, uc.Storage, uc.IPFS)
	if err != nil {
		return err
	}
	err = file.Share(shareWith, &uc.User.Boxer, uc.Storage, uc.Network, uc.IPFS)
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
