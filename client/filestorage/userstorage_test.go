package filestorage

import (
	"testing"

	"ipfs-share/ipfs"
)

func TestDirectory_AppendEntry(t *testing.T) {
	us := UserStorage{}
	ipfsApi, err := ipfs.NewIPFS("http://127.0.0.1", 5001)
	if err != nil {
		t.Fatal()
	}
	_ = us.Init(".", ipfsApi)
	err = us.AddAndShareDirectory("/home/aliras/tmp", []string{}, []string{})
	if err != nil {
		t.Fatal()
	}
}

func TestUserStorage_Init(t *testing.T) {

}
