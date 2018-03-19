package filestorage

import (
	"fmt"
	"ipfs-share/client"
	"ipfs-share/ipfs"
	"testing"
)

func TestUserStorage_AddAndShareFile(t *testing.T) {
	i, err := ipfs.NewIPFS("http://127.0.0.1", 5001)
	n := client.Network{"http://0.0.0.0:6000"}
	if err != nil {
		t.Fatal("could not connect to ipfs daemon")
	}
	us := UserStorage{"test_user", "./data", []*File{}, i, &n}
	err = us.AddAndShareFile("./userstorage.go", []string{"lujza", "blanka"})
	if err != nil {
		t.Fatal()
	}
	fmt.Println(*us.RootDir[0])
}

func TestUserStorage_Init(t *testing.T) {

}
