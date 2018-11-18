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
	"sync"
	"math/big"
	"crypto/rand"
	"context"
	"github.com/ethereum/go-ethereum/core/types"
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
	RegisterUser(username, ipfsPeerId string, boxingKey [32]byte) (*Transaction, error)
	IsUserRegistered(id common.Address) (bool, error)
	CreateGroup(id [32]byte, name string, ipfsHash []byte) (*Transaction, error)
	InviteUser(groupId [32]byte, newMember common.Address, canInvite bool) (*Transaction, error)
	GetGroup(groupId [32]byte) (name string, members []common.Address, ipfsHash []byte, keySalt [32]byte, keyHash [32]byte, err error)
	UpdateGroupIpfsHash(groupId [32]byte, newIpfsHash []byte, approvals []*Approval) (*Transaction, error)
	LeaveGroup(groupId [32]byte) (*Transaction, error)
	GetGroupLeader(groupId [32]byte) (common.Address, error)
	TransactionReceipt(tx *types.Transaction) (*types.Receipt, error)

	// contract events
	GetGroupInvitationSub() *event.Subscription
	GetGroupInvitationChannel() chan *eth.EthGroupInvitation

	GetGroupUpdateIpfsSub() *event.Subscription
	GetGroupUpdateIpfsChannel() chan *eth.EthGroupUpdateIpfsHash

	GetDebugSub() *event.Subscription
	GetDebugChannel() chan *eth.EthDebug

	GetGroupRegisteredSub() *event.Subscription
	GetGroupRegisteredChannel() chan *eth.EthGroupRegistered

	GetKeyDirtySub() *event.Subscription
	GetKeyDirtyChannel() chan *eth.EthKeyDirty

	Close()
}

type Network struct {
	Client *ethclient.Client

	Session *eth.EthSession
	Auth    *bind.TransactOpts

	nonce *big.Int

	// contract events
	groupInvitationSub     event.Subscription
	groupInvitationChannel chan *eth.EthGroupInvitation

	groupUpdateIpfsSub     event.Subscription
	groupUpdateIpfsChannel chan *eth.EthGroupUpdateIpfsHash

	groupRegisteredSub     event.Subscription
	groupRegisteredChannel chan *eth.EthGroupRegistered

	keyDirtySub     event.Subscription
	keyDirtyChannel chan *eth.EthKeyDirty

	debugSub     event.Subscription
	debugChannel chan *eth.EthDebug

	lock sync.Mutex
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
	glog.Infof("ethereum address: %s", wsAddress)

	conn, err := ethclient.Dial(wsAddress)
	if err != nil {
		return nil, fmt.Errorf("could not connect to ethereum node: NewNetwork(): %s", err)
	}

	dipfshare, err := eth.NewEth(common.HexToAddress(contractAddress), conn)
	if err != nil {
		return nil, fmt.Errorf("could not instantiate contract: NewNetwork: %s", err)
	}

	keyFile, err := os.Open(ethKeyPath)
	if err != nil {
		return nil, fmt.Errorf("could not open key file '%s': NewNetwork: %s", ethKeyPath, err)
	}
	defer keyFile.Close()

	auth, err := bind.NewTransactor(keyFile, password)
	if err != nil {
		return nil, fmt.Errorf("could not load account key data: NewNetwork: %s", err)
	}


	if conn == nil {
		glog.Info("conn is nil")
	}
	if auth == nil {
		glog.Info("auth is nil")
	}

	//Max random value, a 130-bits integer, i.e 2^130 - 1
	max := new(big.Int)
	max.Exp(big.NewInt(2), big.NewInt(130), nil).Sub(max, big.NewInt(1))

	//Generate cryptographically strong pseudo-random between 0 - max
	nonce, err := rand.Int(rand.Reader, max)
	if err != nil {
		return nil, errors.Wrap(err, "could not generate tx nonce")
	}

	session := &eth.EthSession{
		Contract: dipfshare,
		CallOpts: bind.CallOpts{
			Pending: true,
		},
		TransactOpts: bind.TransactOpts{
			From:     auth.From,
			Signer:   auth.Signer,
			GasLimit: 3000000,
			Context: context.Background(),
		},
	}

	if session.TransactOpts.Context == nil {
		glog.Info("context is nil")
	}

	channelDebug := make(chan *eth.EthDebug)
	channelGroupInvitation := make(chan *eth.EthGroupInvitation)
	channelGroupUpdateIpfs := make(chan *eth.EthGroupUpdateIpfsHash)

	network := Network{
		Client: conn,

		Session: session,
		Auth:    auth,

		debugChannel:       channelDebug,

		groupInvitationChannel: channelGroupInvitation,
		groupUpdateIpfsChannel: channelGroupUpdateIpfs,
		groupRegisteredChannel: make(chan *eth.EthGroupRegistered),
		keyDirtyChannel: make(chan *eth.EthKeyDirty),

		nonce: nonce,
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
	subGroupUpdateIpfs, err := network.Session.Contract.WatchGroupUpdateIpfsHash(opts, network.groupUpdateIpfsChannel)
	if err != nil {
		return err
	}
	subGroupRegistered, err := network.Session.Contract.WatchGroupRegistered(opts, network.groupRegisteredChannel)
	if err != nil {
		return err
	}
	subKeyDirty, err := network.Session.Contract.WatchKeyDirty(opts, network.keyDirtyChannel)
	if err != nil {
		return err
	}

	network.debugSub = subDebug
	network.groupInvitationSub = subGroupInvitation
	network.groupUpdateIpfsSub = subGroupUpdateIpfs
	network.groupRegisteredSub = subGroupRegistered
	network.keyDirtySub = subKeyDirty

	return nil
}

func (network *Network) TransactionReceipt(tx *types.Transaction) (*types.Receipt, error) {
	return network.Client.TransactionReceipt(context.Background(), tx.Hash())
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

func (network *Network) GetGroupUpdateIpfsChannel() chan *eth.EthGroupUpdateIpfsHash {
	return network.groupUpdateIpfsChannel
}

func (network *Network) GetGroupRegisteredSub() *event.Subscription {
	return &network.groupRegisteredSub
}

func (network *Network) GetGroupRegisteredChannel() chan *eth.EthGroupRegistered {
	return network.groupRegisteredChannel
}

func (network *Network) GetDebugSub() *event.Subscription {
	return &network.debugSub
}

func (network *Network) GetDebugChannel() chan *eth.EthDebug {
	return network.debugChannel
}

func (network *Network) GetKeyDirtySub() *event.Subscription {
	return &network.keyDirtySub
}
func (network *Network) GetKeyDirtyChannel() chan *eth.EthKeyDirty {
	return network.keyDirtyChannel
}

func prepareApprovals(approvals []*Approval) ([]common.Address, [][32]byte, [][32]byte, []uint8, error) {
	var members []common.Address
	var rs [][32]byte
	var ss [][32]byte
	var vs []uint8

	for _, approval := range approvals {
		if len(approval.Signature) != 65 {
			return nil, nil, nil, nil, errors.New("signature length must be 65")
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

	return members, rs, ss, vs, nil
}

func (network *Network) LeaveGroup(groupId [32]byte) (*Transaction, error) {

	tx, err := network.Session.LeaveGroup(groupId)
	if err != nil {
		return nil, errors.Wrap(err, "could not send leaveGroup transaction")
	}

	return &Transaction{tx: tx}, err
}

func (network *Network) UpdateGroupIpfsHash(groupId [32]byte, newIpfsHash []byte, approvals []*Approval) (*Transaction, error) {
	members, rs, ss, vs, err := prepareApprovals(approvals)
	if err != nil {
		return nil, errors.Wrap(err, "could not prepare approvals")
	}

	tx, err := network.Session.UpdateGroupIpfsHash(groupId, newIpfsHash, members, rs, ss, vs)
	if err != nil {
		return nil, errors.Wrap(err, "could not send updateGroupIpfsPath transaction")
	}

	return &Transaction{tx: tx}, err
}

func (network *Network) GetGroupLeader(groupId [32]byte) (common.Address, error) {
	leader, err := network.Session.GetLeader(groupId)
	if err != nil {
		return common.Address{}, err
	}

	return leader, nil
}

func (network *Network) RegisterUser(username, ipfsPeerId string, boxingKey [32]byte) (*Transaction, error) {
	tx, err := network.Session.RegisterUser(username, ipfsPeerId, boxingKey)
	if err != nil {
		return nil, err
	}

	glog.Infof("tx reg user nonce: %d", tx.Nonce())

	return &Transaction{tx: tx}, err
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
		return &Contact{}, err
	}

	return &Contact{
		Address:   address,
		Name:      username,
		IpfsPeerId: ipfsPeerId,
		Boxer:     crypto.AnonymPublicKey{boxingKey},
	}, nil
}

func (network *Network) CreateGroup(id [32]byte, name string, ipfsHash []byte) (*Transaction, error) {
	tx, err := network.Session.CreateGroup(id, name, ipfsHash)
	if err != nil {
		return nil, err
	}

	glog.Infof("tx create group nonce: %d", tx.Nonce())

	return &Transaction{tx: tx}, err
}

func (network *Network) InviteUser(groupId [32]byte, newMember common.Address, canInvite bool) (*Transaction, error) {
	tx, err := network.Session.InviteUser(groupId, newMember, canInvite)
	if err != nil {
		return nil, err
	}

	glog.Infof("tx invite new member nonce: %d", tx.Nonce())

	return &Transaction{tx: tx}, err
}

func (network *Network) GetGroup(groupId [32]byte) (string, []common.Address, []byte, [32]byte, [32]byte, error) {
	return network.Session.GetGroup(groupId)
}


func (network *Network) Close() {
	network.debugSub.Unsubscribe()
	network.groupUpdateIpfsSub.Unsubscribe()
	network.groupInvitationSub.Unsubscribe()
	network.Client.Close()
}
