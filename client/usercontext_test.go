package client

import (
	"fmt"
	"strings"
	"testing"
	"time"

	ipfsapi "ipfs-share/ipfs"
	nw "ipfs-share/network"
)

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
	if err := uc1.CreateGroup("test_group"); err != nil {
		t.Fatal(err)
	}
	if err := uc1.Groups[0].Invite(username1, username2, &uc1.User.Boxer, &uc1.User.Signer.SecretKey); err != nil {
		t.Fatal(err)
	}
	fmt.Println(uc2)
	time.Sleep(45 * time.Second)
	if len(uc1.Groups) != len(uc2.Groups) {
		t.Fatal("#groups does not match")
	}
	if uc1.Groups[0].Members.Length() != uc2.Groups[0].Members.Length() {
		t.Fatal("#(group members) does not match")
	}
	for i := 0; i < uc1.Groups[0].Members.Length(); i++ {
		str1 := uc1.Groups[0].Members.List[i].Name
		str2 := uc2.Groups[0].Members.List[i].Name
		if strings.Compare(str1, str2) != 0 {
			t.Fatal("group members do not match")
		}
	}
	fmt.Println("invite 1st member succeeded")

	username3 := "test_user3"
	uc3, err := NewUserContextFromSignUp(username3, "pw", "./t3/", &network, ipfs)
	if err != nil {
		t.Fatal(err)
	}
	if err := uc1.Groups[0].Invite(username1, username3, &uc1.User.Boxer, &uc1.User.Signer.SecretKey); err != nil {
		t.Fatal(err)
	}
	time.Sleep(45 * time.Second)
	fmt.Println(uc1.Groups[0].Members)
	fmt.Println(uc2.Groups[0].Members)
	if len(uc1.Groups) != len(uc3.Groups) {
		t.Fatal("#groups does not match")
	}
	if uc1.Groups[0].Members.Length() != uc3.Groups[0].Members.Length() {
		t.Fatal("#(group members) does not match")
	}
	for i := 0; i < uc1.Groups[0].Members.Length(); i++ {
		str1 := uc1.Groups[0].Members.List[i].Name
		str2 := uc3.Groups[0].Members.List[i].Name
		if strings.Compare(str1, str2) != 0 {
			t.Fatal("group members do not match")
		}
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
	fmt.Println(uc1.Groups[0].ActiveMembers)
	fmt.Println(uc2.Groups[0].ActiveMembers)

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
	uc1, err := NewUserContextFromSignUp("test_user", "pw", "./t1/", &network, ipfs)
	if err != nil {
		t.Fatal(err)
	}
	uc2, err := NewUserContextFromSignUp("test_user2", "pw", "./t2/", &network, ipfs)
	if err != nil {
		t.Fatal(err)
	}
	uc1.AddAndShareFile("usercontext.go", []string{uc2.User.Username})
	//uc1.Storage.AddAndShareFile("usercontext_test.go", uc1.User.Username, []string{uc2.User.Username}, &uc1.User.Boxer)
	time.Sleep(3 * time.Second)
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
