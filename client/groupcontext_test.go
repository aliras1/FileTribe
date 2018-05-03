package client

import (
	"fmt"
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

	gc1 := GroupContext{uc1.User, group, nil, memberList, &ActiveMemberList{}, &network, ipfs, uc1.Storage}
	//gc2 := GroupContext{uc2.User, group, nil, memberList, &ActiveMemberList{}, &network, ipfs, uc2.Storage}
	//synch2 := NewSynchronizer(user2.Username, &user2.Signer, &gc2)
	//fmt.Println(synch2)

	err = gc1.Invite(uc1.User.Username, "goldmember", &uc1.User.Boxer, &uc1.User.Signer.SecretKey)
	if err != nil {
		t.Fatal(err.Error())
	}
	time.Sleep(3 * time.Second)
}

func TestActiveMemberList(t *testing.T) {
	aml := ActiveMemberList{[]ActiveMember{}}
	go aml.Refresh()

	m1 := Member{"m1", nil}
	m2 := Member{"m2", nil}
	m3 := Member{"m3", nil}
	aml.Set(m1)
	aml.Set(m2)
	aml.Set(m3)

	time.Sleep(1 * time.Second)
	aml.Set(m1)
	time.Sleep(3 * time.Second)
	aml.Set(m3)
	time.Sleep(2 * time.Second)
}

func TestHeartBeat(t *testing.T) {
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
	username3 := "user3"
	user1, err := SignUp(username1, "pw1", ipfsID.ID, &network)
	if err != nil {
		t.Fatal(err)
	}
	user2, err := SignUp(username2, "pw2", ipfsID.ID, &network)
	if err != nil {
		t.Fatal(err)
	}
	user3, err := SignUp(username3, "pw3", ipfsID.ID, &network)
	if err != nil {
		t.Fatal(err)
	}
	group := NewGroup("cucc")

	memberList := &MemberList{[]Member{}}
	memberList = memberList.Append(username1, &network)
	memberList = memberList.Append(username2, &network)
	memberList = memberList.Append(username3, &network)
	fmt.Println(memberList)

	gc1 := GroupContext{user1, group, nil, memberList, &ActiveMemberList{}, &network, ipfs, nil}
	gc2 := GroupContext{user2, group, nil, memberList, &ActiveMemberList{}, &network, ipfs, nil}
	synch1 := NewSynchronizer(user1.Username, &user1.Signer, &gc1)
	synch2 := NewSynchronizer(user2.Username, &user2.Signer, &gc2)
	time.Sleep(4 * time.Second)
	if _, in := synch1.groupCtx.ActiveMembers.Get(username1); !in {
		t.Fatal(username1 + " not in synch1")
	}
	if _, in := synch1.groupCtx.ActiveMembers.Get(username2); !in {
		t.Fatal(username2 + " not in synch1")
	}
	if _, in := synch1.groupCtx.ActiveMembers.Get(username3); in {
		t.Fatal(username3 + " in synch1")
	}

	if _, in := synch2.groupCtx.ActiveMembers.Get(username1); !in {
		t.Fatal(username1 + " not in synch2")
	}
	if _, in := synch2.groupCtx.ActiveMembers.Get(username2); !in {
		t.Fatal(username2 + " not in synch2")
	}
	if _, in := synch2.groupCtx.ActiveMembers.Get(username3); in {
		t.Fatal(username3 + " in synch2")
	}
	gc3 := GroupContext{user3, group, nil, memberList, &ActiveMemberList{}, &network, ipfs, nil}
	synch3 := NewSynchronizer(user1.Username, &user1.Signer, &gc3)
	time.Sleep(2 * time.Second)
	if _, in := synch3.groupCtx.ActiveMembers.Get(username1); !in {
		t.Fatal(username1 + " not in synch3")
	}
	if _, in := synch3.groupCtx.ActiveMembers.Get(username2); !in {
		t.Fatal(username2 + " not in synch3")
	}
	if _, in := synch3.groupCtx.ActiveMembers.Get(username3); in {
		t.Fatal(username3 + " not in synch3")
	}
}
