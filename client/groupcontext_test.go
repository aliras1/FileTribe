package client

import (
	"fmt"
	"ipfs-share/ipfs"
	nw "ipfs-share/network"
	"testing"
	"time"
)

func TestGroupContext_Invite(t *testing.T) {
	network := nw.Network{"http://0.0.0.0:6000"}
	ipfs, err := ipfs.NewIPFS("http://127.0.0.1", 5001)
	if err != nil {
		t.Fatal(err)
	}
	ipfsID, err := ipfs.ID()
	if err != nil {
		t.Fatal(err)
	}
	username1 := "user1"
	username2 := "user2"
	user1, err := SignUp(username1, "pw1", ipfsID.ID, &network)
	if err != nil {
		t.Fatal(err)
	}
	user2, err := SignUp(username2, "pw2", ipfsID.ID, &network)
	if err != nil {
		t.Fatal(err)
	}
	group := NewGroup("cucc")

	gc1 := GroupContext{user1, group, nil, []string{username1, username2}, &network, ipfs, nil}

	gc2 := GroupContext{user2, group, nil, []string{username1, username2}, &network, ipfs, nil}
	synch2 := NewSynchronizer(user2.Username, &user2.Signer, &gc2)
	fmt.Println(synch2)

	err = gc1.Invite(user1.Username, "goldmember", &user1.Boxer, &user1.Signer.SecretKey)
	if err != nil {
		t.Fatal(err.Error())
	}
	time.Sleep(3 * time.Second)
}
