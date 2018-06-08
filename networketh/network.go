package networketh

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	c "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"

	"ipfs-share/crypto"
	"ipfs-share/eth"
)

const key = `{"address":"c4f45f1822b614116ea5b68d4020f3ae1a0179e5","crypto":{"cipher":"aes-128-ctr","ciphertext":"c47565906c488c5122c805a31a3e241d0839cda984903ec28aa07c8892deb5b0","cipherparams":{"iv":"d7814d0dc15a383630c0439c6ad2fea8"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"78d74296f7796969b5764bcfda6cf1cd2cd5bfc423fc0897313b9d23e7e0f219"},"mac":"d852362f275a61fd32acdf040a136a08dc0dc25ab69ddc3d54404b17e9f85826"},"id":"ce2a2147-38d2-4d99-95c1-4968ff6b7a0e","version":3}`
const contractAddress = "0x41cf9ed28c99cc5ebd531bd1929a7e99c122fed8"

type Message struct {
	From    [32]byte `json:"from"`
	Type    string   `json:"type"`
	Payload string   `json:"payload"`
}

type Network struct {
	Session *eth.EthSession
	Auth    *bind.TransactOpts

	MessageSentSubscription event.Subscription
	MessageSentChannel      chan *eth.EthMessageSent

	Simulator *backends.SimulatedBackend
}

func NewTestNetwork() (*Network, error) {
	key, _ := c.GenerateKey()
	auth := bind.NewKeyedTransactor(key)

	simulator := backends.NewSimulatedBackend(core.GenesisAlloc{
		auth.From: core.GenesisAccount{Balance: big.NewInt(10000000000)},
	})

	_, _, ethclient, err := eth.DeployEth(auth, simulator)
	if err != nil {
		return nil, fmt.Errorf("could not deploy contract on cimulated chain")
	}

	session := &eth.EthSession{
		Contract: ethclient,
		CallOpts: bind.CallOpts{
			Pending: true,
		},
		TransactOpts: bind.TransactOpts{
			From:     auth.From,
			Signer:   auth.Signer,
			GasLimit: 3141592,
		},
	}
	channel := make(chan *eth.EthMessageSent)

	start := uint64(0)
	watchOpts := &bind.WatchOpts{
		Start:   &start,
		Context: auth.Context,
	}

	sub, err := session.Contract.WatchMessageSent(watchOpts, channel)
	if err != nil {
		return nil, err
	}

	network := &Network{
		Session: session,
		Auth:    auth,
		MessageSentSubscription: sub,
		MessageSentChannel:      channel,
		Simulator:               simulator,
	}
	return network, nil
}

func NewNetwork() (*Network, error) {

	conn, err := ethclient.Dial("ws://127.0.0.1:8001")
	if err != nil {
		return nil, fmt.Errorf("could not connect to ethereum node: NewNetwork(): %s", err)
	}

	dipfshare, err := eth.NewEth(common.HexToAddress(contractAddress), conn)
	if err != nil {
		return nil, fmt.Errorf("could not instantiate contract: NewNetwork: %s", err)
	}

	auth, err := bind.NewTransactor(strings.NewReader(key), "pwd")
	if err != nil {
		return nil, fmt.Errorf("could not load account key data: NewNetwork: %s", err)
	}

	session := &eth.EthSession{
		Contract: dipfshare,
		CallOpts: bind.CallOpts{
			Pending: true,
		},
		TransactOpts: bind.TransactOpts{
			From:     auth.From,
			Signer:   auth.Signer,
			GasLimit: 3141592, // should be fine tuned...
		},
	}

	channel := make(chan *eth.EthMessageSent)

	start := uint64(0)
	watchOpts := &bind.WatchOpts{
		Start:   &start,
		Context: auth.Context,
	}
	messageSentSubscription, err := session.Contract.WatchMessageSent(watchOpts, channel)
	if err != nil {
		return nil, fmt.Errorf("could not subscript to 'MessageSent' event: NewNetwork: %s", err)
	}

	network := Network{
		Session: session,
		Auth:    auth,
		MessageSentSubscription: messageSentSubscription,
		MessageSentChannel:      channel,
	}

	return &network, nil
}

func (network *Network) RegisterUser(userID [32]byte, username string, boxingKey, verifyKey [32]byte, ipfsAddress string) error {
	_, err := network.Session.RegisterUser(userID, username, boxingKey, verifyKey, ipfsAddress)
	if err != nil {
		return fmt.Errorf("error while Network.RegisterUser(): %s", err)
	}
	return nil
}

func (network *Network) IsUserRegistered(id [32]byte) (bool, error) {
	registered, err := network.Session.IsUserRegistered(id)
	if err != nil {
		return true, fmt.Errorf("error while Network.IsUserRegistered(): %s", err)
	}
	return registered, nil
}

func (network *Network) GetUser(id [32]byte) (common.Address, string, [32]byte, [32]byte, string, error) {
	address, username, boxingKey, verifyKey, ipfsAddress, err := network.Session.GetUser(id)
	if err != nil {
		return common.Address{}, "", [32]byte{}, [32]byte{}, "", fmt.Errorf("error while Network-GetUser(): %s", err)
	}
	return address, username, boxingKey, verifyKey, ipfsAddress, nil
}

func (network *Network) SendMessage(boxer *crypto.AnonymPublicKey, signer *crypto.SigningKeyPair, from [32]byte, msgType, payload string) error {
	message := Message{
		From:    from,
		Type:    msgType,
		Payload: payload,
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("could not marshal message")
	}

	messageSigned := signer.SigningKey.Sign(messageJSON)
	messageEnc, err := boxer.Seal(messageSigned)
	if err != nil {
		return fmt.Errorf("could not encrypt data: Network.SendMessage")
	}

	fmt.Println("net enc:")
	fmt.Println(messageEnc)

	var msg [][1]byte
	for _, b := range messageEnc {
		msg = append(msg, [1]byte{b})
	}
	if _, err := network.Session.SendMessage(msg); err != nil {
		return fmt.Errorf("error while Network.SendMessage(): %s", err)
	}
	return nil
}
