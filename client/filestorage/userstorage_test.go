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
	us := NewUserStorage("./", ipfs, &network)
	err = us.AddAndShareFile("./userstorage.go", "test_user", []string{"lujza", "blanka"})
	if err != nil {
		t.Fatal()
	}
	fmt.Println(*us.RootDir[0])
}

func TestUserStorage_Init(t *testing.T) {

}
