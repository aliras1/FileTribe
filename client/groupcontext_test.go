package client

import (
	"testing"
	"time"

	"ipfs-share/ipfs"
	nw "ipfs-share/network"
)

func TestGroupContext_Invite(t *testing.T) {
	network := nw.Network{"http://0.0.0.0:6000"}
	ipfs, err := ipfs.NewIPFS("http://127.0.0.1", 5001)
	if err != nil {
		t.Fatal(err)
	}
	username1 := "user1"
	username2 := "user2"
	uc1, err := NewUserContextFromSignUp(username1, "pw", "./t1/", &network, ipfs)
	if err != nil {
		t.Fatal(err)
	}
	_, err = NewUserContextFromSignUp(username2, "pw", "./t2/", &network, ipfs)
	if err != nil {
		t.Fatal(err)
	}
	group := NewGroup("cucc")

	memberList := &MemberList{[]Member{}}
	memberList = memberList.Append(username1, &network)
	memberList = memberList.Append(username2, &network)

	gc1, err := NewGroupContext(group, uc1.User, &network, ipfs, uc1.Storage)
	if err != nil {
		t.Fatal(err)
	}
	//gc2 := GroupContext{uc2.User, group, nil, memberList, &ActiveMemberList{}, &network, ipfs, uc2.Storage}
	//synch2 := NewSynchronizer(user2.Username, &user2.Signer, &gc2)
	//fmt.Println(synch2)

	err = gc1.Invite(uc1.User.Username, "goldmember")
	if err != nil {
		t.Fatal(err.Error())
	}
	time.Sleep(3 * time.Second)
}
