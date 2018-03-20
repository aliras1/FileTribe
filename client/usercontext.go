package client

import (
	"errors"
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
	channelMsg  chan nw.Message
	channelSig  chan os.Signal
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

func NewUserContext(dataPath string, user *User, network *nw.Network, ipfs *ipfs.IPFS) *UserContext {
	var uc UserContext
	uc.User = user
	uc.Network = network
	uc.UserStorage = fs.NewUserStorage(dataPath, ipfs, network)

	uc.channelMsg = make(chan nw.Message)
	uc.channelSig = make(chan os.Signal)
	go MessageGetter(uc.User.Username, network, uc.channelMsg, uc.channelSig)
	go MessageProcessor(uc.channelMsg, uc.User.Username, network, ipfs)
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

func MessageProcessor(channelMsg chan nw.Message, username string, network *nw.Network, ipfs *ipfs.IPFS) {
	for msg := range channelMsg {
		fmt.Print("msgproc: ")
		fmt.Println(msg)
		ipfsAddr, err := network.GetUserIPFSAddr(msg.From)
		if err != nil {
			fmt.Println(err)
		}
		listObjects, err := ipfs.List("/ipns/" + ipfsAddr + "/for/" + username)
		if err != nil {
			fmt.Println(err)
		}
		for _, lo := range listObjects.Objects {
			fmt.Println(lo)
		}
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
