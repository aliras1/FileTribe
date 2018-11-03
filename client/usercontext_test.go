package client

import (
	//"flag"
	"fmt"
	"net"
	"os"
	//"strings"
	"testing"
	"github.com/pkg/errors"
)


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

func C() error {
	return errors.New("shit happend")
}

func B() error {
	return errors.Wrap(C(), "B shit")
}

func A() error {
	return errors.Wrap(B(), "A shit")
}

func TestCucc(t *testing.T) {
	err := A()

	fmt.Println(err)
	fmt.Println(errors.Cause(err))
	fmt.Printf("%v", err)
}

func TestBigScenario(t *testing.T) {
	//flag.Set("alsologtostderr", fmt.Sprintf("%t", true))
	//var logLevel string
	//flag.StringVar(&logLevel, "-stderrthreshold", "INFO", "test")
	//
	//CleanUp()
	//
	//password := "pwd"
	//dir := "../test/keystore"
	//ks := keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)
	//
	//ti := time.Now()
	//keyAlice, ethKeyAlicePath, err := nw.NewAccount(ks, dir, password)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//keyBob, ethKeyBobPath, err := nw.NewAccount(ks, dir, password)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//keyCharlie, ethKeyCharliePath, err := nw.NewAccount(ks, dir, password)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//glog.Info("0: ", time.Since(ti))
	//ti = time.Now()
	//testNetwork, err := nw.NewTestNetwork(keyAlice, keyBob, keyCharlie)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//glog.Info("1: ", time.Since(ti))
	//
	//ti = time.Now()
	//testNetwork.SetAuthAlice()
	//alice, err := NewTestUser("alice", true, ethKeyAlicePath, testNetwork)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//glog.Info("2: ", time.Since(ti))
	//
	//ti = time.Now()
	//testNetwork.SetAuthBob()
	//bob, err := NewTestUser("bob", true, ethKeyBobPath, testNetwork)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//glog.Info("3: ", time.Since(ti))
	//
	//testNetwork.SetAuthCharlie()
	//charlie, err := NewTestUser("charlie", true, ethKeyCharliePath, testNetwork)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//glog.Info("----- fun begins -----")
	//
	//ti = time.Now()
	//if err := alice.AddFriend(bob.User.Address); err != nil {
	//	t.Fatal(err)
	//}
	//testNetwork.Simulator.Commit()
	//glog.Info("4: ", time.Since(ti))
	//
	//time.Sleep(2 * time.Second)
	//
	//if len(bob.WaitingFriends) < 1 {
	//	t.Fatal("no friend request")
	//}
	//
	//ti = time.Now()
	//bob.WaitingFriends[0].Confirm(bob)
	//
	//glog.Info("5: ", time.Since(ti))
	//
	//time.Sleep(3 * time.Second)
	//
	//if len(alice.Friends) < 1 {
	//	t.Fatal("alice has no friend")
	//}
	//if len(bob.Friends) < 1 {
	//	t.Fatal("bob has no friend")
	//}
	//
	//if strings.Compare(alice.Friends[0].MyDirectory, bob.Friends[0].HisDirectory) != 0 {
	//	t.Fatal("alices dir mismatch")
	//}
	//if strings.Compare(alice.Friends[0].HisDirectory, bob.Friends[0].MyDirectory) != 0 {
	//	t.Fatal("bobs dir mismatch")
	//}
	//
	//fileID1, err := alice.Commit("./user.go")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fileID2, err := alice.Commit("./user_test.go")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fileID3, err := bob.Commit("./storage.go")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//if len(alice.Repo[alice.User.Address]) < 2 {
	//	t.Fatal("no files by Alice")
	//}
	//if len(bob.Repo[bob.User.Address]) < 1 {
	//	t.Fatal("no files by Bob")
	//}
	//
	//if err := alice.Repo[alice.User.Address][fileID1].Share(alice.Friends[0], alice); err != nil {
	//	t.Fatal(err)
	//}
	//if err := alice.Repo[alice.User.Address][fileID2].Share(alice.Friends[0], alice); err != nil {
	//	t.Fatal(err)
	//}
	//if err := bob.Repo[bob.User.Address][fileID3].Share(bob.Friends[0], bob); err != nil {
	//	t.Fatal(err)
	//}
	//
	//testNetwork.Simulator.Commit()
	//// networkBob.Simulator.Commit()
	//
	//time.Sleep(10 * time.Second)
	//
	//alice.Files()
	//bob.Files()
	//
	//if len(bob.Repo[alice.User.Address]) < 2 {
	//	t.Fatal("no files by Bob2")
	//}
	//if len(alice.Repo[bob.User.Address]) < 1 {
	//	t.Fatal("no files by Alice2")
	//}
	//
	//alice.SignOut()
	//
	//fileID4, err := bob.Commit("./cap.go")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//if err := bob.Repo[bob.User.Address][fileID4].Share(bob.Friends[0], bob); err != nil {
	//	t.Fatal(err)
	//}
	//
	//networkBob.Simulator.Commit()
	//
	//alice, err = Alice(false, ethKeyAlicePath, testNetwork)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//time.Sleep(10 * time.Second)
	//
	//alice.Files()
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
	// 	t.Fatal("did not built any Groups")
	// }
	// fmt.Println("----- members -----")
	// fmt.Println(uc1.Groups[0].Members)
	// fmt.Println(uc2.Groups[0].Members)
	// fmt.Println("----- active members -----")
	// fmt.Println(uc1.Groups[0].Members)
	// fmt.Println(uc2.Groups[0].Members)

	// for i := 0; i < uc1.Groups[0].Members.Length(); i++ {
	// 	str1 := uc1.Groups[0].Members.Files[i].Name
	// 	str2 := uc2.Groups[0].Members.Files[i].Name
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
	// 	t.Fatal("#Groups do not match")
	// }
	// if len(uc1.Groups[0].Members.Files) != len(uc2.Groups[0].Members.Files) {
	// 	t.Fatal("members do not match")
	// }
	// if len(uc1.Groups[0].Members.Files) != len(uc3.Groups[0].Members.Files) {
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

	// // if err := network.DialP2PConn(
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
	// if err := uc1.Commit("usercontext.go", []string{uc2.User.Name}); err != nil {
	// 	t.Fatal(err)
	// }
	// if err := uc1.Commit("usercontext_test.go", []string{uc2.User.Name}); err != nil {
	// 	t.Fatal(err)
	// }
	// time.Sleep(3 * time.Second)
	// uc1.Files()
	// uc2.Files()
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
	// uc1.Files()
	// uc2.Files()
	// time.Sleep(3 * time.Second)
}
