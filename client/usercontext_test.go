package client

import (
	"fmt"
	"testing"
	"time"

	i "ipfs-share/ipfs"
	nw "ipfs-share/network"
)

func TestMessageGetter(t *testing.T) {
	ipfs, err := i.NewIPFS("http://127.0.0.1", 5001)
	network := nw.Network{"http://0.0.0.0:6000"}
	if err != nil {
		t.Fatal("could not connect to ipfs daemon")
	}
	_, err = NewUserContextFromSignUp("test_user", "pw", "./", &network, ipfs)
	if err != nil {
		t.Fatal(err)
	}

	network.SendMessage("from", "test_user", "hello friend!")
	network.SendMessage("from", "test_user", "hello friend, again!")
	fmt.Println("Sleeping...")
	time.Sleep(3 * time.Second)
	fmt.Println("End of test")
}

func TestSharingFromUserContext(t *testing.T) {
	ipfs, err := i.NewIPFS("http://127.0.0.1", 5001)
	network := nw.Network{"http://0.0.0.0:6000"}
	if err != nil {
		t.Fatal("could not connect to ipfs daemon")
	}
	uc1, err := NewUserContextFromSignUp("test_user", "pw", "./t1/", &network, ipfs)
	if err != nil {
		t.Fatal(err)
	}
	uc2, err := NewUserContextFromSignUp("test_user2", "pw", "./t2/", &network, ipfs)
	if err != nil {
		t.Fatal(err)
	}
	uc1.UserStorage.AddAndShareFile("usercontext.go", uc1.User.Username, []string{uc2.User.Username})
	uc1.UserStorage.AddAndShareFile("usercontext_test.go", uc1.User.Username, []string{uc2.User.Username})
	time.Sleep(3 * time.Second)
}

func TestNewUserContextFromSignIn(t *testing.T) {
	ipfs, err := i.NewIPFS("http://127.0.0.1", 5001)
	network := nw.Network{"http://0.0.0.0:6000"}
	if err != nil {
		t.Fatal("could not connect to ipfs daemon")
	}
	uc1, err := NewUserContextFromSignIn("test_user", "pw", "./t1/", &network, ipfs)
	if err != nil {
		t.Fatal(err)
	}
	uc2, err := NewUserContextFromSignIn("test_user2", "pw", "./t2/", &network, ipfs)
	if err != nil {
		t.Fatal(err)
	}
	uc1.UserStorage.List()
	uc2.UserStorage.List()
	time.Sleep(3 * time.Second)
}
