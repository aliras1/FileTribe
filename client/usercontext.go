package client

import (
	"crypto/rand"
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
	Repo    []*fs.FilePTP
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
	uc.Repo, err = uc.Storage.BuildRepo(user.Username, &user.Boxer, network, ipfs)
	if err != nil {
		log.Println(err)
	}

	uc.channelMsg = make(chan nw.Message)
	uc.channelSig = make(chan os.Signal)
	go MessageGetter(uc.User.Username, network, uc.channelMsg, uc.channelSig)
	go MessageProcessor(uc.channelMsg, uc.User.Username, &uc)

	if err := uc.BuildGroups(); err != nil {
		log.Println(err)
	}
	uc.Storage.PublishPublicDir(uc.IPFS)
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
			cap, err := fs.DownloadReadCAP(msg.From, username, msg.Message, &ctx.User.Boxer, ctx.Storage, ctx.Network, ctx.IPFS)
			if err != nil {
				log.Println(err)
				continue
			}
			cap.Store(ctx.Storage)
			fmt.Println(cap)
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
				log.Println("====== E ======")
				log.Println(err)
				continue
			}
			if err := cap.Store(ctx.Storage); err != nil {
				log.Println("====== E2 ======")
				log.Println(err)
				continue
			}
			cap.Boxer.RNG = rand.Reader
			err = ctx.CreateGroupFromCAP(cap, msg.From)
			if err != nil {
				log.Println("====== A ======")
				log.Println(err)
				continue
			}
			//ctx.Storage.PublishPublicDir(ctx.IPFS)
		}
	}
}

func (uc *UserContext) BuildGroups() error {
	caps, err := uc.Storage.GetGroupCAPs()
	if err != nil {
		return err
	}
	for _, cap := range caps {
		if err := uc.CreateGroupFromCAP(&cap, ""); err != nil {
			return err
		}
	}
	return nil
}

func (uc *UserContext) CreateGroup(groupName string) error {
	group := NewGroup(groupName)
	err := group.Register(uc.Network)
	if err != nil {
		return err
	}
	uc.Storage.CreateGroupStorage(groupName)
	groupCtx := GroupContext{uc.User, group, nil,
		&MemberList{[]Member{{uc.User.Username, uc.User.Signer.PublicKey}}},
		&ActiveMemberList{}, uc.Network, uc.IPFS, uc.Storage}

	NewSynchronizer(uc.User.Username, &uc.User.Signer, &groupCtx)
	uc.Groups = append(uc.Groups, &groupCtx)
	return groupCtx.Save()
}

func (uc *UserContext) CreateGroupFromCAP(cap *fs.GroupAccessCAP, from string) error {
	group := &Group{cap.GroupName, cap.Boxer}
	uc.Storage.CreateGroupStorage(group.GroupName)

	memberList := &MemberList{[]Member{}}
	if strings.Compare(from, "") != 0 {
		memberList = memberList.Append(from, uc.Network)
	}
	memberList = memberList.Append(uc.User.Username, uc.Network)
	groupCtx := GroupContext{uc.User, group, nil, memberList,
		&ActiveMemberList{}, uc.Network, uc.IPFS, uc.Storage}

	// load local members to have an idea, who are the good guys.
	// note, that by an invite request 'from' is not a null string
	// therefore we have a valid member we can contact
	if err := groupCtx.LoadGroupData("members.json"); err != nil {
		return err
	}
	NewSynchronizer(uc.User.Username, &uc.User.Signer, &groupCtx)
	time.Sleep(1 * time.Second) // wait for heartbeats of active members
	uc.Groups = append(uc.Groups, &groupCtx)
	if err := groupCtx.PullGroupData("members.json"); err != nil {
		return err
	}
	return groupCtx.Save()
}

func (uc *UserContext) AddAndShareFile(filePath string, shareWith []string) error {
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
		if strings.Compare(path.Base(i.Name), path.Base(filePath)) == 0 {
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

func (uc *UserContext) List() {
	fmt.Println(uc.User.Username)
	for _, f := range uc.Repo {
		fmt.Print("\t--> ")
		fmt.Println(*f)
	}
}
