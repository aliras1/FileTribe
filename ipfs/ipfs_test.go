package ipfs

import (
	"fmt"
	"testing"
)

func TestIPFS_AddDir(t *testing.T) {
	ipfs, err := NewIPFS("http://127.0.0.1", 5001)
	fmt.Println(err)
	ipfs.AddDir("./")
}

func TestIPFS_AddFile(t *testing.T) {
	ipfs, err := NewIPFS("http://127.0.0.1", 5001)
	if err != nil {
		t.Fatal("could not connect to ipfs daemon: " + err.Error())
	}
	mNode, err := ipfs.AddFile("./ipfs.go")
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(*mNode)
}

func TestIPFS_Get(t *testing.T) {
	ipfs, err := NewIPFS("http://127.0.0.1", 5001)
	if err != nil {
		t.Fatal(err)
	}
	ipfs.Get("Qmf6Dea2XP5GqmBdGmfpJNKpqaDCfDbSi1CnLHEz8B7aP9")
}

func TestIPFS_ID(t *testing.T) {
	ipfs, err := NewIPFS("http://127.0.0.1", 5001)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(ipfs.ID())
}
