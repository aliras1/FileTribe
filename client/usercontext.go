package client

import (
	"fmt"
	"os"
	"time"

	fs "ipfs-share/client/filestorage"
	"ipfs-share/ipfs"
	nw "ipfs-share/network"
)

type UserContext struct {
	User        *User
	Network     *nw.Network
	UserStorage *fs.UserStorage
	Messages    []*nw.Message
}

func NewUserContextFromUnamePassw(username, password string, network *nw.Network, ipfs *ipfs.IPFS) *UserContext {
	var uc UserContext
	uc.User = NewUser(username, password)
	uc.Network = network
	uc.UserStorage = &fs.UserStorage{[]*fs.File{}, "test_user", "./data", ipfs, network}
	uc.Messages = []*nw.Message{}

	channelMsg := make(chan nw.Message)
	channelSig := make(chan os.Signal)
	go MessageGetter(uc.User.Username, network, channelMsg, channelSig)
	go MessageProcessor(channelMsg)
	fmt.Println("forked")

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

func MessageProcessor(channelMsg chan nw.Message) {
	for msg := range channelMsg {
		fmt.Print("msgproc: ")
		fmt.Println(msg)
	}
}
