package networketh

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"path"
	// "strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	// c "github.com/ethereum/go-ethereum/crypto"
	// "github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"

	"ipfs-share/crypto"
	"ipfs-share/eth"
)

const key = `{"address":"c4f45f1822b614116ea5b68d4020f3ae1a0179e5","crypto":{"cipher":"aes-128-ctr","ciphertext":"c47565906c488c5122c805a31a3e241d0839cda984903ec28aa07c8892deb5b0","cipherparams":{"iv":"d7814d0dc15a383630c0439c6ad2fea8"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"78d74296f7796969b5764bcfda6cf1cd2cd5bfc423fc0897313b9d23e7e0f219"},"mac":"d852362f275a61fd32acdf040a136a08dc0dc25ab69ddc3d54404b17e9f85826"},"id":"ce2a2147-38d2-4d99-95c1-4968ff6b7a0e","version":3}`
const contractAddress = "0x41cf9ed28c99cc5ebd531bd1929a7e99c122fed8"

type Message struct {
	From    common.Address `json:"from"`
	Type    string         `json:"type"`
	Payload string         `json:"payload"`
}

type Contact struct {
	Address     common.Address
	Name        string
	Boxer       crypto.AnonymPublicKey
	VerifyKey   crypto.VerifyKey
	IPFSAddress string
}

type Network struct {
	Session *eth.EthSession
	Auth    *bind.TransactOpts

	MessageSentSub     event.Subscription
	MessageSentChannel chan *eth.EthMessageSent

	NewFriendRequestSub     event.Subscription
	NewFriendRequestChannel chan *eth.EthNewFriendRequest

	FriendshipConfirmedSub     event.Subscription
	FriendshipConfirmedChannel chan *eth.EthFriendshipConfirmed

	DebugSub     event.Subscription
	DebugChannel chan *eth.EthDebug

	Simulator *backends.SimulatedBackend
}

func NewAccount(ks *keystore.KeyStore, dir, password string) (*ecdsa.PrivateKey, string, error) {
	acc, err := ks.NewAccount(password)
	if err != nil {
		return nil, "", err
	}
	ethKeyPath := dir + "/" + path.Base(acc.URL.String())
	json, err := ioutil.ReadFile(ethKeyPath)
	if err != nil {
		return nil, "", err
	}
	key, err := keystore.DecryptKey(json, password)
	if err != nil {
		return nil, "", err
	}
	return key.PrivateKey, ethKeyPath, nil
}

func newTestNet(ethclient *eth.Eth, auth *bind.TransactOpts, backend *backends.SimulatedBackend) (*Network, error) {
	session := &eth.EthSession{
		Contract: ethclient,
		CallOpts: bind.CallOpts{
			Pending: true,
		},
		TransactOpts: bind.TransactOpts{
			From:     auth.From,
			Signer:   auth.Signer,
			GasLimit: 2141592, // orig 3141592
		},
	}
	channelMessageSent := make(chan *eth.EthMessageSent)
	channelNewFriendRequest := make(chan *eth.EthNewFriendRequest)
	channelFriendshipConfirmed := make(chan *eth.EthFriendshipConfirmed)
	channelDebug := make(chan *eth.EthDebug)

	start := uint64(0)
	watchOpts := &bind.WatchOpts{
		Start:   &start,
		Context: auth.Context,
	}

	subMessageSent, err := session.Contract.WatchMessageSent(watchOpts, channelMessageSent)
	if err != nil {
		return nil, err
	}
	subNewFriendRequest, err := session.Contract.WatchNewFriendRequest(watchOpts, channelNewFriendRequest)
	if err != nil {
		return nil, err
	}
	subFriendshipConfirmed, err := session.Contract.WatchFriendshipConfirmed(watchOpts, channelFriendshipConfirmed)
	if err != nil {
		return nil, err
	}
	subDebug, err := session.Contract.WatchDebug(watchOpts, channelDebug)
	if err != nil {
		return nil, err
	}

	network := &Network{
		Session: session,
		Auth:    auth,

		MessageSentSub:     subMessageSent,
		MessageSentChannel: channelMessageSent,

		NewFriendRequestSub:     subNewFriendRequest,
		NewFriendRequestChannel: channelNewFriendRequest,

		FriendshipConfirmedSub:     subFriendshipConfirmed,
		FriendshipConfirmedChannel: channelFriendshipConfirmed,

		DebugSub:     subDebug,
		DebugChannel: channelDebug,

		Simulator: backend,
	}
	return network, nil
}

func NewTestNetwork(keyAlice, keyBob *ecdsa.PrivateKey) (*Network, *Network, error) {
	authAlice := bind.NewKeyedTransactor(keyAlice)
	authBob := bind.NewKeyedTransactor(keyBob)

	simulator := backends.NewSimulatedBackend(core.GenesisAlloc{
		authAlice.From: core.GenesisAccount{Balance: big.NewInt(10000000000)},
		authBob.From:   core.GenesisAccount{Balance: big.NewInt(10000000000)},
	})

	_, _, ethclient, err := eth.DeployEth(authAlice, simulator)
	if err != nil {
		return nil, nil, fmt.Errorf("could not deploy contract on simulated chain")
	}

	networkAlice, err := newTestNet(ethclient, authAlice, simulator)
	if err != nil {
		return nil, nil, err
	}

	networkBob, err := newTestNet(ethclient, authBob, simulator)
	if err != nil {
		return nil, nil, err
	}

	return networkAlice, networkBob, nil
}

func NewNetwork() (*Network, error) {

	// conn, err := ethclient.Dial("ws://127.0.0.1:8001")
	// if err != nil {
	// 	return nil, fmt.Errorf("could not connect to ethereum node: NewNetwork(): %s", err)
	// }

	// dipfshare, err := eth.NewEth(common.HexToAddress(contractAddress), conn)
	// if err != nil {
	// 	return nil, fmt.Errorf("could not instantiate contract: NewNetwork: %s", err)
	// }

	// auth, err := bind.NewTransactor(strings.NewReader(key), "pwd")
	// if err != nil {
	// 	return nil, fmt.Errorf("could not load account key data: NewNetwork: %s", err)
	// }

	// session := &eth.EthSession{
	// 	Contract: dipfshare,
	// 	CallOpts: bind.CallOpts{
	// 		Pending: true,
	// 	},
	// 	TransactOpts: bind.TransactOpts{
	// 		From:     auth.From,
	// 		Signer:   auth.Signer,
	// 		GasLimit: 3141592, // should be fine tuned...
	// 	},
	// }

	// channel := make(chan *eth.EthMessageSent)

	// start := uint64(0)
	// watchOpts := &bind.WatchOpts{
	// 	Start:   &start,
	// 	Context: auth.Context,
	// }
	// messageSentSubscription, err := session.Contract.WatchMessageSent(watchOpts, channel)
	// if err != nil {
	// 	return nil, fmt.Errorf("could not subscript to 'MessageSent' event: NewNetwork: %s", err)
	// }

	// network := Network{
	// 	Session: session,
	// 	Auth:    auth,
	// 	MessageSentSubscription: messageSentSubscription,
	// 	MessageSentChannel:      channel,
	// }

	// return &network, nil
	return nil, nil
}

func goToEthByteArray(array []byte) [][1]byte {
	var newArray [][1]byte
	for _, b := range array {
		newArray = append(newArray, [1]byte{b})
	}
	return newArray
}

func ethToGoByteArray(array [][1]byte) []byte {
	var newArray []byte
	for _, b := range array {
		newArray = append(newArray, b[0])
	}
	return newArray
}

func (network *Network) RegisterUser(username string, boxingKey [32]byte, verifyKey []byte, ipfsAddress string) error {
	_, err := network.Session.RegisterUser(username, boxingKey, verifyKey, ipfsAddress)
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
	username, boxingKey, verifyKey, ipfsAddress, err := network.Session.GetUser(address)
	if err != nil {
		return &Contact{}, fmt.Errorf("error while Network-GetUser(): %s", err)
	}
	return &Contact{
		Address:     address,
		Name:        username,
		Boxer:       crypto.AnonymPublicKey{&boxingKey},
		VerifyKey:   crypto.VerifyKey(verifyKey),
		IPFSAddress: ipfsAddress,
	}, nil
}

func (network *Network) SendMessage(boxer *crypto.AnonymPublicKey, signer *crypto.Signer, from common.Address, msgType, payload string) error {
	message := Message{
		From:    from,
		Type:    msgType,
		Payload: payload,
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("could not marshal message")
	}

	messageSigned, err := signer.Sign(messageJSON)
	if err != nil {
		return fmt.Errorf("could not sign message: Network.SendMessage")
	}
	messageEnc, err := boxer.Seal(messageSigned)
	if err != nil {
		return fmt.Errorf("could not encrypt data: Network.SendMessage")
	}

	fmt.Println("net enc:")
	fmt.Println(messageEnc)

	if _, err := network.Session.SendMessage(messageEnc); err != nil {
		return fmt.Errorf("error while Network.SendMessage(): %s", err)
	}
	return nil
}
