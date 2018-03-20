package ipfs

import (
	"fmt"
	"strings"
	"testing"
)

func TestIPFS_AddDir(t *testing.T) {
	ipfs, err := NewIPFS("http://127.0.0.1", 5001)
	if err != nil {
		t.Fatal(err)
	}
	ipfs.AddDir("./data/filestorage/data/public")
}

func TestIPFS_List(t *testing.T) {
	ipfs, err := NewIPFS("http://127.0.0.1", 5001)
	if err != nil {
		t.Fatal(err)
	}
	merkleNodes, err := ipfs.AddDir("../client/filestorage/data/public")
	if err != nil {
		t.Fatal(err)
	}
	for _, mn := range merkleNodes {
		if strings.Compare(mn.Name, "public") == 0 {
			listObjects, err := ipfs.List("/ipfs/" + mn.Hash)
			if err != nil {
				t.Fatal(err)
			}
			for _, lo := range listObjects.Objects {
				fmt.Println(lo)
			}
			break
		}
	}
}

func TestIPFS_NamePublish(t *testing.T) {
	ipfs, err := NewIPFS("http://127.0.0.1", 5001)
	if err != nil {
		t.Fatal(err)
	}
	merkleNodes, err := ipfs.AddDir("../client/filestorage/data/public")
	if err != nil {
		t.Fatal(err)
	}
	for _, mn := range merkleNodes {
		if strings.Compare(mn.Name, "public") == 0 {
			err = ipfs.NamePublish(mn.Hash)
			if err != nil {
				t.Fatal(err)
			}
			break
		}
	}
	ipfsID, err := ipfs.ID()
	if err != nil {
		t.Fatal(err)
	}
	listObjects, err := ipfs.List("/ipns/" + ipfsID.ID + "/for")
	if err != nil {
		t.Fatal(err)
	}
	for _, lo := range listObjects.Objects {
		fmt.Println(lo)
	}
}

func TestIPFS_AddFile(t *testing.T) {
	ipfs, err := NewIPFS("http://127.0.0.1", 5001)
	if err != nil {
		t.Fatal("could not connect to ipfs daemon: " + err.Error())
	}
	mNode, err := ipfs.AddFile("./ipfs.go")
	if err != nil {
		t.Fatal(err)
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
