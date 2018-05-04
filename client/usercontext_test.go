package client

import (
	"fmt"
	"strings"
	"testing"
	"time"

	ipfsapi "ipfs-share/ipfs"
	nw "ipfs-share/network"
)

func TestUserContext_CreateGroup(t *testing.T) {
	network := &nw.Network{"http://0.0.0.0:6000"}
	ipfs, err := ipfsapi.NewIPFS("http://127.0.0.1", 5001)
	// sign up and create group
	fmt.Println("signing up and creating group...")
	uc, err := NewUserContextFromSignUp("testuser", "pw", "./testuser/", network, ipfs)
	if err != nil {
		t.Fatal(err)
	}
	if err := uc.CreateGroup("cucc_group"); err != nil {
		t.Fatal(err)
	}
	if err := uc.CreateGroup("cucc_group2"); err != nil {
		t.Fatal(err)
	}
	// sign in and check if group is correct
	fmt.Println("signing in and checking group...")
	uc, err = NewUserContextFromSignIn("testuser", "pw", "./testuser/", network, ipfs)
	if err != nil {
		t.Fatal(err)
	}
	if len(uc.Groups) < 2 {
		t.Fatal("no groups found")
	}
	for _, group := range uc.Groups {
		fmt.Print(group.Group.GroupName + ": ")
		fmt.Println(group.Members)
	}
}

func TestGroupInvite(t *testing.T) {
	ipfs, err := ipfsapi.NewIPFS("http://127.0.0.1", 5001)
	network := nw.Network{"http://0.0.0.0:6000"}
	if err != nil {
		t.Fatal("could not connect to ipfs daemon")
	}
	username1 := "test_user"
	username2 := "test_user2"
	uc1, err := NewUserContextFromSignUp(username1, "pw", "./t1/", &network, ipfs)
	if err != nil {
		t.Fatal(err)
	}
	uc2, err := NewUserContextFromSignUp(username2, "pw", "./t2/", &network, ipfs)
	if err != nil {
		t.Fatal(err)
	}
	groupname := "test_group"
	if err := uc1.CreateGroup(groupname); err != nil {
		t.Fatal(err)
	}
	if err := uc1.Groups[0].Invite(uc1.User.Username, uc2.User.Username); err != nil {
		t.Fatal(err)
	}
}

func TestSignInAndBuildUpAfterInviteTest(t *testing.T) {
	ipfs, err := ipfsapi.NewIPFS("http://127.0.0.1", 5001)
	network := nw.Network{"http://0.0.0.0:6000"}
	if err != nil {
		t.Fatal("could not connect to ipfs daemon")
	}
	username1 := "test_user"
	username2 := "test_user2"
	uc1, err := NewUserContextFromSignIn(username1, "pw", "./t1/", &network, ipfs)
	if err != nil {
		t.Fatal(err)
	}
	uc2, err := NewUserContextFromSignIn(username2, "pw", "./t2/", &network, ipfs)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(3 * time.Second)
	fmt.Println(uc1)
	fmt.Println(uc2)
	if len(uc1.Groups) < 1 || len(uc2.Groups) < 1 {
		t.Fatal("did not built any groups")
	}
	fmt.Println("----- members -----")
	fmt.Println(uc1.Groups[0].Members)
	fmt.Println(uc2.Groups[0].Members)
	fmt.Println("----- active members -----")
	fmt.Println(uc1.Groups[0].Members)
	fmt.Println(uc2.Groups[0].Members)

	for i := 0; i < uc1.Groups[0].Members.Length(); i++ {
		str1 := uc1.Groups[0].Members.List[i].Name
		str2 := uc2.Groups[0].Members.List[i].Name
		if strings.Compare(str1, str2) != 0 {
			t.Fatal("group members do not match")
		}
	}
}

func TestMessageGetter(t *testing.T) {
	ipfs, err := ipfsapi.NewIPFS("http://127.0.0.1", 5001)
	network := nw.Network{"http://0.0.0.0:6000"}
	if err != nil {
		t.Fatal("could not connect to ipfs daemon")
	}
	_, err = NewUserContextFromSignUp("test_user", "pw", "./t1", &network, ipfs)
	if err != nil {
		t.Fatal(err)
	}

	network.SendMessage("from", "test_user", "test", "hello friend!")
	network.SendMessage("from", "test_user", "test", "hello friend, again!")
	fmt.Println("Sleeping...")
	time.Sleep(3 * time.Second)
	fmt.Println("End of test")
}

func TestSharingFromUserContext(t *testing.T) {
	ipfs, err := ipfsapi.NewIPFS("http://127.0.0.1", 5001)
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
	if err := uc1.AddAndShareFile("usercontext.go", []string{uc2.User.Username}); err != nil {
		t.Fatal(err)
	}
	if err := uc1.AddAndShareFile("usercontext_test.go", []string{uc2.User.Username}); err != nil {
		t.Fatal(err)
	}
	time.Sleep(3 * time.Second)
	uc1.List()
	uc2.List()
}

func TestNewUserContextFromSignIn(t *testing.T) {
	ipfs, err := ipfsapi.NewIPFS("http://127.0.0.1", 5001)
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
	uc1.List()
	uc2.List()
	time.Sleep(3 * time.Second)
}
