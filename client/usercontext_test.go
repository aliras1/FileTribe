package client

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang/glog"

	"github.com/ethereum/go-ethereum/accounts/keystore"

	ipfsapi "ipfs-share/ipfs"
	nw "ipfs-share/networketh"
)

func Alice(signup bool, ethKeyPath string, network *nw.Network) (*UserContext, error) {
	ipfs, err := ipfsapi.NewIPFS("http://127.0.0.1", 5001)
	if err != nil {
		return nil, fmt.Errorf("could not create new ipfs api conn: Alice: %s", err)
	}

	username := "alice"
	password := "pwd"
	homeDir := "./alice/"
	var alice *UserContext

	if signup {
		alice, err = NewUserContextFromSignUp(username, password, ethKeyPath, homeDir, network, ipfs)
		if err != nil {
			return nil, fmt.Errorf("could not sign up: Alice: %s", err)
		}
	} else {
		alice, err = NewUserContextFromSignIn(username, password, ethKeyPath, homeDir, network, ipfs)
		if err != nil {
			return nil, fmt.Errorf("could not sign in: Alice: %s", err)
		}
	}

	network.Simulator.Commit()

	reg, err := network.IsUserRegistered(alice.User.Address)
	if err != nil {
		return nil, err
	}
	if !reg {
		return nil, fmt.Errorf("bob not regged")
	}

	return alice, nil
}

func Bob(signup bool, ethKeyPath string, network *nw.Network) (*UserContext, error) {
	ipfs, err := ipfsapi.NewIPFS("http://127.0.0.1", 5001)
	if err != nil {
		return nil, fmt.Errorf("could not create new ipfs api conn: Bob: %s", err)
	}

	username := "bob"
	password := "pwd"
	homeDir := "./bob/"
	var bob *UserContext

	if signup {
		bob, err = NewUserContextFromSignUp(username, password, ethKeyPath, homeDir, network, ipfs)
		if err != nil {
			return nil, fmt.Errorf("could not sign up: Bob: %s", err)
		}
	} else {
		bob, err = NewUserContextFromSignIn(username, password, ethKeyPath, homeDir, network, ipfs)
		if err != nil {
			return nil, fmt.Errorf("could not sign in: Bob: %s", err)
		}
	}

	network.Simulator.Commit()

	reg, err := network.IsUserRegistered(bob.User.Address)
	if err != nil {
		return nil, err
	}
	if !reg {
		return nil, fmt.Errorf("bob not regged")
	}

	return bob, nil
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
	flag.Set("alsologtostderr", fmt.Sprintf("%t", true))
	var logLevel string
	flag.StringVar(&logLevel, "-stderrthreshold", "INFO", "test")

	CleanUp()

	password := "pwd"
	dir := "../test/keystore"
	ks := keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)
	keyAlice, ethKeyAlicePath, err := nw.NewAccount(ks, dir, password)
	if err != nil {
		t.Fatal(err)
	}
	keyBob, ethKeyBobPath, err := nw.NewAccount(ks, dir, password)
	if err != nil {
		t.Fatal(err)
	}

	networkAlice, networkBob, err := nw.NewTestNetwork(keyAlice, keyBob)
	if err != nil {
		t.Fatal(err)
	}

	alice, err := Alice(true, ethKeyAlicePath, networkAlice)
	if err != nil {
		t.Fatal(err)
	}
	bob, err := Bob(true, ethKeyBobPath, networkBob)
	if err != nil {
		t.Fatal(err)
	}

	glog.Info("----- fun begins -----")

	if err := alice.AddFriend(bob.User.Address); err != nil {
		t.Fatal(err)
	}
	networkAlice.Simulator.Commit()

	time.Sleep(2 * time.Second)

	if len(bob.WaitingFriends) < 1 {
		t.Fatal("no friend request")
	}
	bob.WaitingFriends[0].Confirm(bob)
	networkBob.Simulator.Commit()

	time.Sleep(2 * time.Second)

	if len(alice.Friends) < 1 {
		t.Fatal("alice has no friend")
	}
	if len(bob.Friends) < 1 {
		t.Fatal("bob has no friend")
	}

	if strings.Compare(alice.Friends[0].MyDirectory, bob.Friends[0].HisDirectory) != 0 {
		t.Fatal("alices dir mismatch")
	}
	if strings.Compare(alice.Friends[0].HisDirectory, bob.Friends[0].MyDirectory) != 0 {
		t.Fatal("bobs dir mismatch")
	}
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
