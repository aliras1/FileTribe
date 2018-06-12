package client

import (
	"fmt"
	"net"
	"os"
	"testing"
	// "time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"ipfs-share/eth"
	ipfsapi "ipfs-share/ipfs"
	nw "ipfs-share/networketh"
)

func Alice(signup bool, network *nw.Network) (*UserContext, *nw.Network, error) {
	ipfs, err := ipfsapi.NewIPFS("http://127.0.0.1", 5001)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create new ipfs api conn: Alice: %s", err)
	}

	channel := make(chan *eth.EthMessageSent)

	start := uint64(0)
	watchOpts := &bind.WatchOpts{
		Start:   &start,
		Context: network.Auth.Context,
	}

	sub, err := network.Session.Contract.WatchMessageSent(watchOpts, channel)
	if err != nil {
		return nil, nil, err
	}
	aliceNet := &nw.Network{
		Session: network.Session,
		Auth:    network.Auth,
		MessageSentSubscription: sub,
		MessageSentChannel:      channel,
		Simulator:               network.Simulator,
	}
	username := "alice"
	password := "pwd"
	homeDir := "./alice/"
	var alice *UserContext

	if signup {
		alice, err = NewUserContextFromSignUp(username, password, homeDir, aliceNet, ipfs)
		if err != nil {
			return nil, nil, fmt.Errorf("could not sign up: Alice: %s", err)
		}
	} else {
		alice, err = NewUserContextFromSignIn(username, password, homeDir, aliceNet, ipfs)
		if err != nil {
			return nil, nil, fmt.Errorf("could not sign in: Alice: %s", err)
		}
	}

	aliceNet.Simulator.Commit()

	return alice, aliceNet, nil
}

func Bob(signup bool, network *nw.Network) (*UserContext, *nw.Network, error) {
	ipfs, err := ipfsapi.NewIPFS("http://127.0.0.1", 5001)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create new ipfs api conn: Bob: %s", err)
	}

	channel := make(chan *eth.EthMessageSent)

	start := uint64(0)
	watchOpts := &bind.WatchOpts{
		Start:   &start,
		Context: network.Auth.Context,
	}

	sub, err := network.Session.Contract.WatchMessageSent(watchOpts, channel)
	if err != nil {
		return nil, nil, err
	}
	bobNet := &nw.Network{
		Session: network.Session,
		Auth:    network.Auth,
		MessageSentSubscription: sub,
		MessageSentChannel:      channel,
		Simulator:               network.Simulator,
	}
	username := "bob"
	password := "pwd"
	homeDir := "./bob/"
	var bob *UserContext

	if signup {
		bob, err = NewUserContextFromSignUp(username, password, homeDir, bobNet, ipfs)
		if err != nil {
			return nil, nil, fmt.Errorf("could not sign up: Bob: %s", err)
		}
	} else {
		bob, err = NewUserContextFromSignIn(username, password, homeDir, bobNet, ipfs)
		if err != nil {
			return nil, nil, fmt.Errorf("could not sign in: Bob: %s", err)
		}
	}

	bobNet.Simulator.Commit()

	reg, err := bobNet.IsUserRegistered(bob.User.Address)
	if err != nil {
		return nil, nil, err
	}
	if !reg {
		return nil, nil, fmt.Errorf("bob not regged")
	}

	return bob, bobNet, nil
}

func Charlie(signup bool, network *nw.Network) (*UserContext, error) {
	ipfs, err := ipfsapi.NewIPFS("http://127.0.0.1", 5003)
	if err != nil {
		return nil, fmt.Errorf("could not create new ipfs api conn: Charlie: %s", err)
	}
	username := "charlie"
	password := "pwd"
	homeDir := "./charlie/"
	var charlie *UserContext

	if signup {
		charlie, err = NewUserContextFromSignUp(username, password, homeDir, network, ipfs)
		if err != nil {
			return nil, fmt.Errorf("could not sign up: Charlie: %s", err)
		}
	} else {
		charlie, err = NewUserContextFromSignIn(username, password, homeDir, network, ipfs)
		if err != nil {
			return nil, fmt.Errorf("could not sign in: Charlie: %s", err)
		}
	}

	return charlie, nil
}

func CleanUp() {
	os.RemoveAll("./alice")
	os.RemoveAll("./bob")
	os.RemoveAll("./charlie")
}

func GetIPAddress() (string, error) {
	name, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("could not get hostname: GetIPAddress: %s", err)
	}
	addrs, err := net.LookupHost(name)
	if err != nil {
		return "", fmt.Errorf("could not get addresses: GetIPAddress: %s", err)
	}
	if len(addrs) < 1 {
		return "", fmt.Errorf("no ip addresses found: GetIPAddress")
	}
	return addrs[0], nil
}

func TestBigScenario(t *testing.T) {
	// flag.Set("alsologtostderr", fmt.Sprintf("%t", true))
	// var logLevel string
	// flag.StringVar(&logLevel, "logLevel", "4", "test")

	// CleanUp()
	// ip, err := GetIPAddress()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// network := &nw.Network{Address: "http://" + ip + ":6000"}
	// alice, err := Alice(true, network)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// // +----------------+
	// // | GROUP CREATION |
	// // +----------------+
	// // create some groups
	// if err := alice.CreateGroup("alice_group"); err != nil {
	// 	t.Fatal(err)
	// }
	// if err := alice.CreateGroup("alice_group2"); err != nil {
	// 	t.Fatal(err)
	// }

	// // check if can not create duplicates
	// err = alice.CreateGroup("alice_group")
	// if !strings.Contains(err.Error(), "group name already exists") {
	// 	t.Fatal(err)
	// }
	// err = alice.CreateGroup("alice_group2")
	// if !strings.Contains(err.Error(), "group name already exists") {
	// 	t.Fatal(err)
	// }

	// // check groups
	// if len(alice.Groups) < 2 {
	// 	t.Fatalf("%d number of groups found instead of 2", len(alice.Groups))
	// }

	// // sign in and check if groups are built up
	// alice.SignOut()
	// time.Sleep(2 * time.Second)
	// alice, err = Alice(false, network)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// if len(alice.Groups) < 2 {
	// 	t.Fatalf("%d number of groups found instead of 2", len(alice.Groups))
	// }

	// bob, err := Bob(true, network)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// charlie, err := Charlie(true, network)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// // +------------+
	// // | INVITATION |
	// // +------------+

	// // invite bob
	// if err := alice.Groups[0].Invite(bob.User.Name); err != nil {
	// 	t.Fatal(err)
	// }
	// // concurrently invite charlie
	// // just the first invitations should be successful
	// if err := alice.Groups[0].Invite(charlie.User.Name); err != nil {
	// 	t.Fatal(err)
	// }

	// time.Sleep(210 * time.Second)
	// if len(bob.Groups) < 1 {
	// 	t.Fatal("bob has no groups")
	// }
	// if len(alice.Groups[0].Members.List) != len(bob.Groups[0].Members.List) {
	// 	t.Fatal("members do not match")
	// }

	// // sign out and in with bob and check if he builds up the group
	// bob.SignOut()
	// bob, err = Bob(false, network)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// if len(bob.Groups) < 1 {
	// 	t.Fatal("bob has no groups")
	// }
	// if len(alice.Groups[0].Members.List) != len(bob.Groups[0].Members.List) {
	// 	t.Fatal("members do not match")
	// }

	// // invite charlie, consensus needed now
	// if err := alice.Groups[0].Invite(charlie.User.Name); err != nil {
	// 	t.Fatal(err)
	// }
	// time.Sleep(130 * time.Second)
	// if len(charlie.Groups) < 1 {
	// 	t.Fatal("charlie has no groups")
	// }
	// if len(alice.Groups[0].Members.List) != len(bob.Groups[0].Members.List) &&
	// 	len(alice.Groups[0].Members.List) != len(charlie.Groups[0].Members.List) {

	// 	t.Fatal("members do not match")
	// }

	// // sign out and in with charlie and check if he builds up the group
	// charlie.SignOut()
	// charlie, err = Charlie(false, network)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// if len(charlie.Groups) < 1 {
	// 	t.Fatal("charlie has no groups")
	// }
	// if len(alice.Groups[0].Members.List) != len(charlie.Groups[0].Members.List) {
	// 	t.Fatal("members do not match")
	// }

	// CleanUp()
}

func TestSignInAndBuildUpAfterInviteTest(t *testing.T) {
	// ipfs, err := ipfsapi.NewIPFS("http://127.0.0.1", 5001)
	// network := nw.Network{"http://0.0.0.0:6000"}
	// if err != nil {
	// 	t.Fatal("could not connect to ipfs daemon")
	// }
	// username1 := "test_user"
	// username2 := "test_user2"
	// uc1, err := NewUserContextFromSignIn(username1, "pw", "./test1/", &network, ipfs)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// uc2, err := NewUserContextFromSignIn(username2, "pw", "./test2/", &network, ipfs)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// time.Sleep(3 * time.Second)
	// fmt.Println(uc1)
	// fmt.Println(uc2)
	// if len(uc1.Groups) < 1 || len(uc2.Groups) < 1 {
	// 	t.Fatal("did not built any groups")
	// }
	// fmt.Println("----- members -----")
	// fmt.Println(uc1.Groups[0].Members)
	// fmt.Println(uc2.Groups[0].Members)
	// fmt.Println("----- active members -----")
	// fmt.Println(uc1.Groups[0].Members)
	// fmt.Println(uc2.Groups[0].Members)

	// for i := 0; i < uc1.Groups[0].Members.Length(); i++ {
	// 	str1 := uc1.Groups[0].Members.List[i].Name
	// 	str2 := uc2.Groups[0].Members.List[i].Name
	// 	if strings.Compare(str1, str2) != 0 {
	// 		t.Fatal("group members do not match")
	// 	}
	// }
}

func TestGroupInviteWithMoreMembers(t *testing.T) {
	// ipfs, err := ipfsapi.NewIPFS("http://127.0.0.1", 5001)
	// network := nw.Network{"http://0.0.0.0:6000"}
	// if err != nil {
	// 	t.Fatal("could not connect to ipfs daemon")
	// }
	// username1 := "test_user"
	// username2 := "test_user2"
	// username3 := "test_user3"
	// uc1, err := NewUserContextFromSignIn(username1, "pw", "./test1/", &network, ipfs)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// uc2, err := NewUserContextFromSignIn(username2, "pw", "./test2/", &network, ipfs)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// uc3, err := NewUserContextFromSignUp(username3, "pw", "./test3/", &network, ipfs)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// fmt.Println(uc1.Groups)
	// fmt.Println(uc2.Groups)
	// fmt.Println(uc3.Groups)
	// if err := uc1.Groups[0].Invite(username3); err != nil {
	// 	t.Fatal(err)
	// }
	// time.Sleep(130 * time.Second)
	// if len(uc1.Groups) != len(uc3.Groups) && len(uc2.Groups) != len(uc3.Groups) {
	// 	t.Fatal("#groups do not match")
	// }
	// if len(uc1.Groups[0].Members.List) != len(uc2.Groups[0].Members.List) {
	// 	t.Fatal("members do not match")
	// }
	// if len(uc1.Groups[0].Members.List) != len(uc3.Groups[0].Members.List) {
	// 	t.Fatal("members do not match")
	// }
}

func TestMessages(t *testing.T) {
	// network, err := nw.NewTestNetwork()
	// if err != nil {
	// 	t.Fatal("could not connect to eth node")
	// }
	// alice, _, err := Alice(true, network)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// bob, _, err := Bob(true, network)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// // if err := network.SendMessage(
	// // 	&bob.User.Boxer.PublicKey,
	// // 	&alice.User.Signer,
	// // 	alice.User.EthKey.Address,
	// // 	"test_type",
	// // 	"hello friend!",
	// // ); err != nil {
	// // 	t.Fatal(err)
	// // }

	// network.Simulator.Commit()

	// fmt.Println("Sleeping...")
	// time.Sleep(3 * time.Second)
	// fmt.Println("End of test")

	// CleanUp()
}

func TestSharingFromUserContext(t *testing.T) {
	// ipfs, err := ipfsapi.NewIPFS("http://127.0.0.1", 5001)
	// network := nw.Network{"http://0.0.0.0:6000"}
	// if err != nil {
	// 	t.Fatal("could not connect to ipfs daemon")
	// }
	// uc1, err := NewUserContextFromSignUp("test_user", "pw", "./t1/", &network, ipfs)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// uc2, err := NewUserContextFromSignUp("test_user2", "pw", "./t2/", &network, ipfs)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// if err := uc1.AddAndShareFile("usercontext.go", []string{uc2.User.Name}); err != nil {
	// 	t.Fatal(err)
	// }
	// if err := uc1.AddAndShareFile("usercontext_test.go", []string{uc2.User.Name}); err != nil {
	// 	t.Fatal(err)
	// }
	// time.Sleep(3 * time.Second)
	// uc1.List()
	// uc2.List()
}

func TestNewUserContextFromSignIn(t *testing.T) {
	// ipfs, err := ipfsapi.NewIPFS("http://127.0.0.1", 5001)
	// network := nw.Network{"http://0.0.0.0:6000"}
	// if err != nil {
	// 	t.Fatal("could not connect to ipfs daemon")
	// }
	// uc1, err := NewUserContextFromSignIn("test_user", "pw", "./t1/", &network, ipfs)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// uc2, err := NewUserContextFromSignIn("test_user2", "pw", "./t2/", &network, ipfs)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// uc1.List()
	// uc2.List()
	// time.Sleep(3 * time.Second)
}
