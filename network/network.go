package network

import (
	"github.com/golang/glog"
	"crypto/ecdsa"
	// "encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"

	"ipfs-share/crypto"
	"ipfs-share/eth"
	"github.com/pkg/errors"
)

// const key = `{"address":"c4f45f1822b614116ea5b68d4020f3ae1a0179e5","crypto":{"cipher":"aes-128-ctr","ciphertext":"c47565906c488c5122c805a31a3e241d0839cda984903ec28aa07c8892deb5b0","cipherparams":{"iv":"d7814d0dc15a383630c0439c6ad2fea8"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"78d74296f7796969b5764bcfda6cf1cd2cd5bfc423fc0897313b9d23e7e0f219"},"mac":"d852362f275a61fd32acdf040a136a08dc0dc25ab69ddc3d54404b17e9f85826"},"id":"ce2a2147-38d2-4d99-95c1-4968ff6b7a0e","version":3}`
// const contractAddress = "0x41cf9ed28c99cc5ebd531bd1929a7e99c122fed8"

type Message struct {
	From    common.Address `json:"from"`
	Type    string         `json:"type"`
	Payload string         `json:"payload"`
}

type Contact struct {
	Address   common.Address
	Name      string
	IpfsPeerId string
	Boxer     crypto.AnonymPublicKey
}

type Approval struct {
	From common.Address
	Signature []byte
}


type INetwork interface {
	GetUser(address common.Address) (*Contact, error)
	RegisterUser(username, ipfsPeerId string, boxingKey [32]byte) error
	IsUserRegistered(common.Address) (bool, error)
	CreateGroup(id [32]byte, name string, ipfsPath string) error
	InviteUser(groupId [32]byte, newMember common.Address) error
	GetGroup(groupId [32]byte) (string, []common.Address, string, error)
	UpdateGroupIpfsPath(groupId [32]byte, newIpfsPath string, approvals []*Approval) error

	// contract events
	GetGroupInvitationSub() *event.Subscription
	GetGroupInvitationChannel() chan *eth.EthGroupInvitation

	GetGroupUpdateIpfsSub() *event.Subscription
	GetGroupUpdateIpfsChannel() chan *eth.EthGroupUpdateIpfsPath

	GetDebugSub() *event.Subscription
	GetDebugChannel() chan *eth.EthDebug

	Close()
}

type Network struct {
	Client *ethclient.Client

	Session *eth.EthSession
	Auth    *bind.TransactOpts

	// contract events
	groupInvitationSub     event.Subscription
	groupInvitationChannel chan *eth.EthGroupInvitation

	groupUpdateIpfsSub     event.Subscription
	groupUpdateIpfsChannel chan *eth.EthGroupUpdateIpfsPath

	debugSub     event.Subscription
	debugChannel chan *eth.EthDebug
}

func NewAccount(ks *keystore.KeyStore, ethKeyPath, password string) (*ecdsa.PrivateKey, error) {
	json, err := ioutil.ReadFile(ethKeyPath)
	if err != nil {
		return nil, err
	}
	key, err := keystore.DecryptKey(json, password)
	if err != nil {
		return nil, err
	}
	return key.PrivateKey, nil
}

func NewNetwork(wsAddress, ethKeyPath, contractAddress, password string) (INetwork, error) {

	conn, err := ethclient.Dial(wsAddress)
	if err != nil {
		return nil, fmt.Errorf("could not connect to ethereum node: NewNetwork(): %s", err)
	}

	dipfshare, err := eth.NewEth(common.HexToAddress(contractAddress), conn)
	if err != nil {
		return nil, fmt.Errorf("could not instantiate contract: NewNetwork: %s", err)
	}

	key, err := os.Open(ethKeyPath)
	if err != nil {
		return nil, fmt.Errorf("could not open key file '%s': NewNetwork: %s", ethKeyPath, err)
	}

	auth, err := bind.NewTransactor(key, password)
	if err != nil {
		return nil, fmt.Errorf("could not load account key data: NewNetwork: %s", err)
	}


	if conn == nil {
		glog.Info("conn is nil")
	}
	if auth == nil {
		glog.Info("auth is nil")
	}
	

	session := &eth.EthSession{
		Contract: dipfshare,
		CallOpts: bind.CallOpts{
			Pending: true,
		},
		TransactOpts: bind.TransactOpts{
			From:     auth.From,
			Signer:   auth.Signer,
			GasLimit: 2141592,
		},
	}

	if session.TransactOpts.Context == nil {
		glog.Info("context is nil")
	}

	channelDebug := make(chan *eth.EthDebug)
	channelGroupInvitation := make(chan *eth.EthGroupInvitation)
	channelGroupUpdateIpfs := make(chan *eth.EthGroupUpdateIpfsPath)

	network := Network{
		Client: conn,

		Session: session,
		Auth:    auth,

		debugChannel:       channelDebug,

		groupInvitationChannel: channelGroupInvitation,
		groupUpdateIpfsChannel: channelGroupUpdateIpfs,
	}

	if err := network.SubscribeToEvents(0); err != nil {
		return nil, err
	}
	
	return &network, nil
}



func (network *Network) SubscribeToEvents(latestBlock uint64) error {
	opts := &bind.WatchOpts{
		Context: network.Auth.Context,
		Start: &latestBlock,
	}

	subDebug, err := network.Session.Contract.WatchDebug(opts, network.debugChannel)
	if err != nil {
		return err
	}
	subGroupInvitation, err := network.Session.Contract.WatchGroupInvitation(opts, network.groupInvitationChannel)
	if err != nil {
		return err
	}
	subGroupUpdateIpfs, err := network.Session.Contract.WatchGroupUpdateIpfsPath(opts, network.groupUpdateIpfsChannel)
	if err != nil {
		return err
	}

	network.debugSub = subDebug
	network.groupInvitationSub = subGroupInvitation
	network.groupUpdateIpfsSub = subGroupUpdateIpfs

	return nil
}

func (network *Network) GetGroupInvitationSub() *event.Subscription {
	return &network.groupInvitationSub
}

func (network *Network) GetGroupInvitationChannel() chan *eth.EthGroupInvitation {
	return network.groupInvitationChannel
}

func (network *Network) GetGroupUpdateIpfsSub() *event.Subscription {
	return &network.groupUpdateIpfsSub
}

func (network *Network) GetGroupUpdateIpfsChannel() chan *eth.EthGroupUpdateIpfsPath {
	return network.groupUpdateIpfsChannel
}

func (network *Network) GetDebugSub() *event.Subscription {
	return &network.debugSub
}

func (network *Network) GetDebugChannel() chan *eth.EthDebug {
	return network.debugChannel
}

func (network *Network) UpdateGroupIpfsPath(groupId [32]byte, newIpfsPath string, approvals []*Approval) error {
	var members []common.Address
	var rs [][32]byte
	var ss [][32]byte
	var vs []uint8

	for _, approval := range approvals {
		if len(approval.Signature) != 65 {
			return errors.New("signature length must be 65")
		}

		members = append(members, approval.From)

		var r [32]byte
		copy(r[:], approval.Signature[:32])
		rs = append(rs, r)

		var s [32]byte
		copy(s[:], approval.Signature[32:64])
		ss = append(ss, s)

		v := approval.Signature[64]
		if v < 27 {
			v = v + 27;
		}
		vs = append(vs, v)
	}

	_, err := network.Session.UpdateGroupIpfsPath(groupId, newIpfsPath, members, rs, ss, vs)
	if err != nil {
		return errors.Wrapf(err, "could not send updateGroupIpfsPath transaction")
	}

	return nil
}

func (network *Network) RegisterUser(username, ipfsPeerId string, boxingKey [32]byte) error {
	_, err := network.Session.RegisterUser(username, ipfsPeerId, boxingKey)
	if err != nil {
		return fmt.Errorf("error while Network.RegisterUser(): %s", err)
	}
	return nil
}

func (network *Network) IsUserRegistered(id common.Address) (bool, error) {
	registered, err := network.Session.IsUserRegistered(id)
	if err != nil {
		return true, fmt.Errorf("error while Network.IsUserRegistered(): %s", err)
	}
	return registered, nil
}

func (network *Network) GetUser(address common.Address) (*Contact, error) {
	username, ipfsPeerId, boxingKey, err := network.Session.GetUser(address)
	if err != nil {
		return &Contact{}, fmt.Errorf("error while Network-GetUser(): %s", err)
	}
	return &Contact{
		Address:   address,
		Name:      username,
		IpfsPeerId: ipfsPeerId,
		Boxer:     crypto.AnonymPublicKey{&boxingKey},
	}, nil
}

func (network *Network) CreateGroup(id [32]byte, name string, ipfsPath string) error {
	_, err := network.Session.CreateGroup(id, name, ipfsPath)
	return err
}

func (network *Network) InviteUser(groupId [32]byte, newMember common.Address) error {
	_, err := network.Session.InviteUser(groupId, newMember)
	return err
}

func (network *Network) GetGroup(groupId [32]byte) (string, []common.Address, string, error) {
	return network.Session.GetGroup(groupId)
}


func (network *Network) Close() {
	network.debugSub.Unsubscribe()
	network.groupUpdateIpfsSub.Unsubscribe()
	network.groupInvitationSub.Unsubscribe()
	network.Client.Close()
}
