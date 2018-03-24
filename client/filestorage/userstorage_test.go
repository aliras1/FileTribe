package filestorage

import (
	"fmt"
	"testing"

	i "ipfs-share/ipfs"
	nw "ipfs-share/network"
)

func TestUserStorage_AddAndShareFile(t *testing.T) {
	ipfs, err := i.NewIPFS("http://127.0.0.1", 5001)
	network := nw.Network{"http://0.0.0.0:6000"}
	if err != nil {
		t.Fatal("could not connect to ipfs daemon")
	}
	us := NewUserStorage("./", "test_user", &network, ipfs)
	err = us.AddAndShareFile("./userstorage.go", "test_user", []string{"lujza", "blanka"})
	if err != nil {
		t.Fatal(err)
	}
	us.List()
}

func TestUserStorage_build_shared(t *testing.T) {
	ipfs, err := i.NewIPFS("http://127.0.0.1", 5001)
	network := nw.Network{"http://0.0.0.0:6000"}
	if err != nil {
		t.Fatal("could not connect to ipfs daemon")
	}
	us := NewUserStorage("../t1", "test_user", &network, ipfs)
	if us == nil {
		t.Fatal("could not instantiate UserStorage")
	}
	for _, file := range us.RootDir {
		fmt.Println(file)
	}
}

func TestUserStorage_build_caps(t *testing.T) {
	ipfs, err := i.NewIPFS("http://127.0.0.1", 5001)
	network := nw.Network{"http://0.0.0.0:6000"}
	if err != nil {
		t.Fatal("could not connect to ipfs daemon")
	}
	us := NewUserStorage("../t2", "test_user2", &network, ipfs)
	if us == nil {
		t.Fatal("could not instantiate UserStorage")
	}
	for _, file := range us.RootDir {
		fmt.Println(file)
	}
}
