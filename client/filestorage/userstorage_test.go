package filestorage

import (
	"fmt"
	"testing"

	"ipfs-share/ipfs"
	"ipfs-share/network"
)

func TestUserStorage_AddAndShareFile(t *testing.T) {
	i, err := ipfs.NewIPFS("http://127.0.0.1", 5001)
	n := network.Network{"http://0.0.0.0:6000"}
	if err != nil {
		t.Fatal("could not connect to ipfs daemon")
	}
	us := UserStorage{[]*File{}, "test_user", "./data", i, &n}
	err = us.AddAndShareFile("./userstorage.go", []string{"lujza", "blanka"})
	if err != nil {
		t.Fatal()
	}
	fmt.Println(*us.RootDir[0])
}

func TestUserStorage_Init(t *testing.T) {

}
