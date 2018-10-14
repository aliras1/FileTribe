package networketh

import (
	// "bytes"
	"fmt"
	"math/big"
	// "strings"
	"testing"
	"time"

	// "ipfs-share/client"
	"ipfs-share/eth"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
)

func NewBlockchain() (*backends.SimulatedBackend, common.Address, *bind.TransactOpts, error) {
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)

	simulator := backends.NewSimulatedBackend(core.GenesisAlloc{
		auth.From: core.GenesisAccount{Balance: big.NewInt(10000000000)},
	})

	address, _, _, err := eth.DeployEth(auth, simulator)
	if err != nil {
		return nil, common.Address{}, nil, fmt.Errorf("could not deploy contract on cimulated chain")
	}

	return simulator, address, auth, nil
}

type UserData struct {
	username    string
	password    string
	boxingKey   [32]byte
	verifyKey   [32]byte
	ipfsAddress string
}

func Alice() *UserData {
	return &UserData{
		username:    "Alice",
		password:    "pwd",
		boxingKey:   [32]byte{2},
		verifyKey:   [32]byte{1}, // = id
		ipfsAddress: "as2356lksdjfaf723jgajsdf",
	}
}

func Bob() *UserData {
	return &UserData{
		username:    "Bob",
		password:    "pwd",
		boxingKey:   [32]byte{2},
		verifyKey:   [32]byte{1}, // = id
		ipfsAddress: "1231hkhhashdahdahas12",
	}
}

func TestNewNetwork(t *testing.T) {
	if _, err := NewNetwork(); err != nil {
		t.Fatal(err)
	}
}

func registerUser(network *Network, username string, boxingKey, verifyKey [32]byte, ipfsAddress string) error {
	registered, err := network.IsUserRegistered(network.Auth.From)
	if err != nil {
		return err
	}
	if registered {
		return fmt.Errorf("how is he registered?")
	}

	if err := network.RegisterUser(username, boxingKey, verifyKey, ipfsAddress); err != nil {
		return err
	}

	network.Simulator.Commit()

	return nil
}



func TestRegisterAndRetrieveUser(t *testing.T) {
	dir := "../test/keystore"
	ks := keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)
	keyAlice, _, err := newAccount(ks, dir, "pwd")
	if err != nil {
		t.Fatal(err)
	}
	keyBob, _, err := newAccount(ks, dir, "pwd")
	if err != nil {
		t.Fatal(err)
	}

	networkAlice, networkBob, err := NewTestNetwork(keyAlice, keyBob)
	if err != nil {
		t.Fatal(err)
	}

	aliceData := Alice()
	bobData := Bob()

	if err := registerUser(networkAlice, aliceData.username, aliceData.boxingKey, aliceData.verifyKey, aliceData.ipfsAddress); err != nil {
		t.Fatal(err)
	}
	if err := registerUser(networkBob, bobData.username, bobData.boxingKey, bobData.verifyKey, bobData.ipfsAddress); err != nil {
		t.Fatal(err)
	}

	// Test AliceNet
	registered, err := networkAlice.IsUserRegistered(networkAlice.Auth.From)
	if err != nil {
		t.Fatal(err)
	}
	if !registered {
		t.Fatal("Alice should be registered on alice net")
	}
	registered, err = networkAlice.IsUserRegistered(networkBob.Auth.From)
	if err != nil {
		t.Fatal(err)
	}
	if !registered {
		t.Fatal("Bob should be registered on alice net")
	}
	// Test Bob net
	registered, err = networkBob.IsUserRegistered(networkAlice.Auth.From)
	if err != nil {
		t.Fatal(err)
	}
	if !registered {
		t.Fatal("Alice should be registered on bob net")
	}
	registered, err = networkBob.IsUserRegistered(networkBob.Auth.From)
	if err != nil {
		t.Fatal(err)
	}
	if !registered {
		t.Fatal("Bob should be registered on bob net")
	}

	// _, uName, bKey, vKey, ipfs, err := network.GetUser(id)
	// if strings.Compare(uName, aliceData.username) != 0 {
	// 	t.Fatal("usernames do not match")
	// }
	// if !bytes.Equal(bKey[:], aliceData.boxingKey[:]) {
	// 	t.Fatal("boxing keys do not match")
	// }
	// if !bytes.Equal(vKey[:], aliceData.verifyKey[:]) {
	// 	t.Fatal("verify keys do not match")
	// }
	// if strings.Compare(ipfs, aliceData.ipfsAddress) != 0 {
	// 	t.Fatal("ipfs addresses do not match")
	// }
}

func TestSendMessage(t *testing.T) {
	network, err := NewNetwork()
	if err != nil {
		t.Fatal(err)
	}

	// aliceData := Alice()
	// alice := client.NewUser(aliceData.username, aliceData.password)

	// if err := network.RegisterUser(aliceData.username, aliceData.boxingKey, aliceData.verifyKey, aliceData.ipfsAddress); err != nil {
	// 	t.Fatal(err)
	// }

	// fmt.Println("1st sleep")
	// time.Sleep(30 * time.Second)

	go func() {
		fmt.Println("itten")
		for {
			select {
			case err := <-network.MessageSentSubscription.Err():
				t.Fatal(err)
			case log := <-network.MessageSentChannel:
				fmt.Println(log)
				network.MessageSentSubscription.Unsubscribe()
				return
			}
		}
	}()

	// message := []byte("Hello friend!")
	// if err := network.DialP2PConn(message); err != nil {
	// 	t.Fatal(err)
	// }

	fmt.Println("2nd sleep")
	time.Sleep(30 * time.Second)
}


func TestTheTest(t *testing.T) {
	// key, _ := crypto.GenerateKey()
	// auth := bind.NewKeyedTransactor(key)

	// simulator := backends.NewSimulatedBackend(core.GenesisAlloc{
	// 	auth.From: core.GenesisAccount{Balance: big.NewInt(10000000000)},
	// })

	// _, _, ethclient, err := eth.DeployEth(auth, simulator)
	// if err != nil {
	// 	t.Fatal("could not deploy contract on cimulated chain")
	// }

	// session := &eth.EthSession{
	// 	Contract: ethclient,
	// 	CallOpts: bind.CallOpts{
	// 		Pending: true,
	// 	},
	// 	TransactOpts: bind.TransactOpts{
	// 		From:     auth.From,
	// 		Signer:   auth.Signer,
	// 		GasLimit: 3141592,
	// 	},
	// }
	// fmt.Println("itten5")
	// channel := make(chan *eth.EthMessageSent)

	// start := uint64(0)
	// watchOpts := &bind.WatchOpts{
	// 	Start:   &start,
	// 	Context: auth.Context,
	// }

	// sub, err := session.Contract.WatchMessageSent(watchOpts, channel)
	// if err != nil {
	// 	t.Fatal("could not subscript to 'MessageSent' event: NewNetwork: %s " + err.Error())
	// }

	// go func() {
	// 	fmt.Println("itten")
	// 	for {
	// 		select {
	// 		case err := <-sub.Err():
	// 			t.Fatal(err)
	// 		case log := <-channel:
	// 			fmt.Println(log)
	// 			sub.Unsubscribe()
	// 			return
	// 		}
	// 	}
	// }()

	// fmt.Println("itten4")
	// aliceData := Alice()
	// id := aliceData.verifyKey

	// registered, err := session.IsUserRegistered(id)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// if registered {
	// 	t.Fatal("how is he registered?")
	// }

	// if _, err := session.RegisterUser(aliceData.username, aliceData.boxingKey, aliceData.verifyKey, aliceData.ipfsAddress); err != nil {
	// 	t.Fatal(err)
	// }

	// simulator.Commit()

	// registered, err = session.IsUserRegistered(id)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// if !registered {
	// 	t.Fatal("user should be registered")
	// }

	// _, uName, bKey, vKey, ipfs, err := session.GetUser(id)
	// if strings.Compare(uName, aliceData.username) != 0 {
	// 	t.Fatal("usernames do not match")
	// }
	// if !bytes.Equal(bKey[:], aliceData.boxingKey[:]) {
	// 	t.Fatal("boxing keys do not match")
	// }
	// if !bytes.Equal(vKey[:], aliceData.verifyKey[:]) {
	// 	t.Fatal("verify keys do not match")
	// }
	// if strings.Compare(ipfs, aliceData.ipfsAddress) != 0 {
	// 	t.Fatal("ipfs addresses do not match")
	// }

	// // message := []byte("hello friend")
	// // var msg [][1]byte
	// // for _, b := range message {
	// // 	msg = append(msg, [1]byte{b})
	// // }
	// // if _, err := session.DialP2PConn(msg); err != nil {
	// // 	t.Fatal(err)
	// // }

	// simulator.Commit()

	// time.Sleep(2 * time.Second)
}