package client

import (
	"crypto/rand"
	"fmt"
	"path"
	"strings"
	"time"

	fs "ipfs-share/client/filestorage"
	"ipfs-share/ipfs"
	nw "ipfs-share/network"
	"github.com/golang/glog"
)

type UserContext struct {
	User     *User // TODO lock boxer
	Groups   []*GroupContext
	Repo     []*fs.FilePTP
	IPNSAddr string
	Network  *nw.Network
	IPFS     *ipfs.IPFS
	Storage  *fs.Storage // TODO lock

	channelMsg  chan nw.Message
	channelStop chan int
}

func NewUserContextFromSignUp(username, password, dataPath string, network *nw.Network, ipfs *ipfs.IPFS) (*UserContext, error) {
	fmt.Printf("[*] User '%s' signing up...\n", username)

	ipfsID, err := ipfs.ID()
	if err != nil {
		return nil, fmt.Errorf("could not get ipfs id: NewUserContextFromSignUp: %s", err)
	}
	user, err := SignUp(username, password, ipfsID.ID, network)
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

func NewUserContextFromSignIn(username, password, dataPath string, network *nw.Network, ipfs *ipfs.IPFS) (*UserContext, error) {
	fmt.Printf("[*] User '%s' signing in...\n", username)

	user, err := SignIn(username, password, network)
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

func NewUserContext(dataPath string, user *User, network *nw.Network, ipfs *ipfs.IPFS) (*UserContext, error) {
	var err error
	var uc UserContext
	uc.User = user
	uc.Network = network
	uc.IPFS = ipfs
	uc.Storage = fs.NewStorage(dataPath)
	uc.Storage.Init()
	uc.Groups = []*GroupContext{}
	uc.Repo, err = uc.Storage.BuildRepo(user.Name, &user.Boxer, network, ipfs)
	if err != nil {
		return nil, fmt.Errorf("could not build file repo: NewUserContext: %s", err)
	}
	ipfsID, err := ipfs.ID()
	if err != nil {
		return nil, fmt.Errorf("could not get ipfs id: NewUserContect: %s", err)
	}
	uc.IPNSAddr = ipfsID.ID

	uc.channelMsg = make(chan nw.Message)
	uc.channelStop = make(chan int)
	go MessageGetter(uc.User.Name, network, uc.channelMsg, uc.channelStop)
	go MessageProcessor(uc.channelMsg, uc.User.Name, &uc)

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

func MessageGetter(username string, network *nw.Network, channelMsg chan nw.Message, channelStop chan int) {
	for true {
		select {
		case <-channelStop:
			glog.Infof("User '%s' received 'STOP' signal: MessageGetter", username)
			close(channelMsg)
			close(channelStop)
			return
		default:
			msgs, err := network.GetMessages(username)
			if err != nil {
				glog.Errorf("could not get messages: MessageGetter: %s", err)
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
	glog.Infof("User '%s' is processing messages...\n", username)

	for msg := range channelMsg {
		glog.Infof("--> user '%s' got message", username)

		switch msg.Type {
		case "PTP READ CAP":
			fmt.Println("[*] Downloading PTP file")
			cap, err := fs.DownloadReadCAP(msg.From, username, msg.Message, &ctx.User.Boxer, ctx.Storage, ctx.Network, ctx.IPFS)
			if err != nil {
				glog.Errorf("could not download read cap '%s' while 'PTP READ CAP' message: MessageProcessor: %s\n", msg.Message, err)
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
			fmt.Printf("[*] Joining group '%s'...\n", strings.Split(msg.Message, ".")[0])
			glog.Infof("\t--> Donwloading group access cap '%s'...", msg.Message)

			cap, err := ctx.Storage.DownloadGroupAccessCAP(msg.From, username, msg.Message, &ctx.User.Boxer, ctx.Network, ctx.IPFS)
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
	group := NewGroup(groupname)
	if err := group.Register(uc.User.Name, uc.Network); err != nil {
		return fmt.Errorf("could not register group: UserContext.CreateGroup: %s", err)
	}
	if err := group.Save(uc.Storage); err != nil {
		return fmt.Errorf("could not save group: UserContext.CreateGroup: %s", err)
	}
	groupCtx, err := NewGroupContext(group, uc.User, uc.Network, uc.IPFS, uc.Storage)
	if err != nil {
		return fmt.Errorf("could not create new group context: UserContext.CreateGroup: %s", err)
	}
	uc.Groups = append(uc.Groups, groupCtx)

	fmt.Printf("[*] Group '%s' created!\n", groupname)

	return nil
}

func (uc *UserContext) AddAndShareFile(filePath string, shareWith []string) error {
	fmt.Printf("[*] PTP sharing file '%s' with user '%s'...\n", filePath, shareWith[0])

	if uc.isFileInRepo("/ipns/" + uc.IPNSAddr + "/files/" + path.Base(filePath)) {
		return fmt.Errorf("file is already added to the repo: UserContext.AddAndShareFile")
	}
	file, err := fs.NewSharedFilePTP(filePath, uc.User.Name, uc.Storage, uc.IPFS)
	if err != nil {
		return fmt.Errorf("could not create new shared file '%s': UserContext.AddAndShareFile: %s", file.Name, err)
	}
	if err := file.Share(shareWith, &uc.User.Boxer, uc.Storage, uc.Network, uc.IPFS); err != nil {
		return fmt.Errorf("could not share file '%s': UserContext.AddAndShareFile: %s", file.Name, err)
	}
	uc.addFileToRepo(file)

	fmt.Println("[*] PTP sharing ended")

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
