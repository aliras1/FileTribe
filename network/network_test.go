package network

import (
	"testing"
	"github.com/ethereum/go-ethereum/common"
)

const (
	ethWSAddress   = "ws://172.18.0.2:8001"
	//ethKeyPath = "../misc/eth/ethkeystore/UTC--2018-05-19T19-06-19.498239404Z--c4f45f1822b614116ea5b68d4020f3ae1a0179e5"
	ethKeyPath = "../misc/eth/ethkeystore/UTC--2018-05-19T19-09-28.926310760Z--2938fb7144b41ec6d7a34c1380eb6ec92535c1c3"
	contractAddress = "0x41cf9ed28c99cc5ebd531bd1929a7e99c122fed8"


	password = "pwd"
)

func TestRegisterUser(t *testing.T) {
	network, err := NewNetwork(ethWSAddress, ethKeyPath, contractAddress, password)
	if err != nil {
		t.Fatal(err)
	}
	//address := common.HexToAddress("0xc4f45f1822b614116ea5b68d4020f3ae1a0179e5")
	if err := network.RegisterUser("bob", "asd", [32]byte{2}); err != nil {
		t.Fatal(err)
	}
}

func TestIsUserRegistered(t *testing.T) {
	network, err := NewNetwork(ethWSAddress, ethKeyPath, contractAddress, password)
	if err != nil {
		t.Fatal(err)
	}

	address := common.HexToAddress("2938fb7144b41ec6d7a34c1380eb6ec92535c1c3")
	ok, err := network.IsUserRegistered(address)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("not registered")
	}
}