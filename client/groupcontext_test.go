package client

import (
	"testing"

	nw "ipfs-share/network"

	"github.com/golang/glog"
	"fmt"
	"flag"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/golang/mock/gomock"
	"time"

	ipfsapi "ipfs-share/ipfs"
	"github.com/libp2p/go-libp2p-peer"
	"crypto/ecdsa"
	"ipfs-share/utils"
	"io/ioutil"
)

type FakePubSubRecord struct {
	from string
	data string
}

func (r *FakePubSubRecord) From() peer.ID {
	return peer.ID(r.from)
}

func (r *FakePubSubRecord) Data() []byte {
	return []byte(r.data)
}

func (r *FakePubSubRecord) SeqNo() int64 {
	return 0
}

func (r *FakePubSubRecord) TopicIDs() []string {
	return []string{}
}

const (
	ALICE = 0
	BOB = 1
	CHARLIE = 2
)

var (
	shells []*ipfsapi.Ipfs
	subs []*ipfsapi.PubSubSubscription
	pubsubRecords = []FakePubSubRecord{}
	controller *gomock.Controller
)

func NewTestUser(username string, signup bool, ethKeyPath string, shellIdx int, network nw.INetwork, p2pPort string) (*UserContext, error) {
	t := time.Now()
	glog.Info("ipfs inst: ", time.Since(t))
	password := "pwd"
	homeDir := "./" + username + "/"
	var testUser *UserContext
	var err error

	var ipfs ipfsapi.IIpfs

	switch username {
	case "alice":
		{
			ipfs = ipfsapi.NewIpfs("http://127.0.0.1:5001")
		}
	case "bob":
		{
			ipfs = ipfsapi.NewIpfs("http://127.0.0.1:5002")
		}
	case "charlie":
		{
			ipfs = ipfsapi.NewIpfs("http://127.0.0.1:5003")
		}
	default:
		{
			ipfs = nil
		}
	}

	if signup {
		testUser, err = NewUserContextFromSignUp(username, password, ethKeyPath, homeDir, network, ipfs, p2pPort)
		if err != nil {
			return nil, fmt.Errorf("could not sign up: %s: %s", username, err)
		}
	} else {
		testUser, err = NewUserContextFromSignIn(username, password, ethKeyPath, homeDir, network, ipfs, p2pPort)
		if err != nil {
			return nil, fmt.Errorf("could not sign in: %s: %s", username, err)
		}
	}

	reg, err := network.IsUserRegistered(testUser.User.Address())
	if err != nil {
		return nil, err
	}
	if !reg {
		return nil, fmt.Errorf("%s not regged", username)
	}

	return testUser, nil
}


func TestGroupContext_Invite(t *testing.T) {
	flag.Set("alsologtostderr", fmt.Sprintf("%t", true))
	var logLevel string
	flag.StringVar(&logLevel, "-stderrthreshold", "INFO", "test")

	controller = gomock.NewController(t)
	defer controller.Finish()

	CleanUp()

	password := "pwd"
	dir := "../test/keystore/"
	ks := keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)

	ethKeyAlicePath := dir + "UTC--2018-10-10T08-19-58.398032114Z--ab083e63cfc7525634642075d49a0de31374bc0f"
	keyAlice, err := nw.NewAccount(ks, ethKeyAlicePath, password)
	if err != nil {
		t.Fatal(err)
	}

	ethKeyBobPath := dir + "UTC--2018-10-10T08-20-04.769949175Z--be9678b9882dac288093b9d38ea7382f21479c77"
	keyBob, err := nw.NewAccount(ks, ethKeyBobPath, password)
	if err != nil {
		t.Fatal(err)
	}

	ethKeyCharliePath := dir + "UTC--2018-10-10T08-20-10.903818650Z--d7ad6058180005d6639653f1d0216e481a43af79"
	keyCharlie, err := nw.NewAccount(ks, ethKeyCharliePath, password)
	if err != nil {
		t.Fatal(err)
	}

	fakeNetwork, err := nw.NewTestNetwork([]*ecdsa.PrivateKey{keyAlice, keyBob, keyCharlie})
	if err != nil {
		t.Fatal(err)
	}

	fakeNetwork.SetAuth(ALICE)
	alice, err := NewTestUser("alice", true, ethKeyAlicePath, 0, fakeNetwork, "2000")
	if err != nil {
		t.Fatal(err)
	}

	fakeNetwork.SetAuth(BOB)
	bob, err := NewTestUser("bob", true, ethKeyBobPath, 1, fakeNetwork, "2001")
	if err != nil {
		t.Fatal(err)
	}

	fakeNetwork.SetAuth(CHARLIE)
	charlie, err := NewTestUser("charlie", true, ethKeyCharliePath, 2, fakeNetwork, "2002")
	if err != nil {
		t.Fatal(err)
	}

	glog.Info("----- fun begins -----")

	fakeNetwork.SetAuth(ALICE)

	if err := alice.CreateGroup("GRUPPE"); err != nil {
		t.Fatal(err)
	}

	if alice.Groups.Count() != 1 {
		t.Fatal("no groupAtAlice found by alice")
	}

	groupAtAlice := alice.Groups.FirstOrDefault(nil).(*GroupContext)
	groupAtAlice.Invite(bob.User.Address(), true)
	groupAtAlice.Invite(charlie.User.Address(), true)

	time.Sleep(5 * time.Second)

	if bob.Groups.Count() != 1 {
		t.Fatal("no group found by bob")
	}
	if charlie.Groups.Count() != 1 {
		t.Fatal("no group found by charlie")
	}

	fmt.Println(alice.Groups.FirstOrDefault(nil).(*GroupContext).Group.IpfsHash)
	fmt.Println(bob.Groups.FirstOrDefault(nil).(*GroupContext).Group.IpfsHash)
	fmt.Println(charlie.Groups.FirstOrDefault(nil).(*GroupContext).Group.IpfsHash)

	if alice.Groups.FirstOrDefault(nil).(*GroupContext).Group.CountMembers() != 3 {
		t.Fatal("alice's groupAtAlice has not got enough members")
	}
	if bob.Groups.FirstOrDefault(nil).(*GroupContext).Group.CountMembers() != 3 {
		t.Fatal("bob's groupAtAlice has not got enough members")
	}
	if charlie.Groups.FirstOrDefault(nil).(*GroupContext).Group.CountMembers() != 3 {
		t.Fatal("charlie's groupAtAlice has not got enough members")
	}

	fileAlice := "./alice/data/userdata/root/" + groupAtAlice.Group.Id().ToString() + "/rrrepo.go"
	if err := utils.CopyFile("./grouprepo.go", fileAlice); err != nil {
		t.Fatal(err)
	}
	if err := groupAtAlice.CommitChanges(); err != nil {
		t.Fatal(err)
	}

	time.Sleep(5 * time.Second)
	fmt.Println("----------- change file ------------")

	groupAtBob := bob.Groups.FirstOrDefault(nil).(*GroupContext)
	fileBob := "./bob/data/userdata/root/" + groupAtBob.Group.Id().ToString() + "/rrrepo.go"
	if err := AppendToFile(fileBob, "Bob's modification (should fail)\n"); err != nil {
		t.Fatal(err)
	}

	if err := groupAtBob.CommitChanges(); err != nil {
		t.Fatal(err)
	}

	time.Sleep(5 * time.Second)
	fmt.Println("----------- Grant W access to only alice  ------------")

	if err := groupAtAlice.GrantWriteAccess(fileAlice, bob.User.Address()); err != nil {
		t.Fatal(err)
	}
	if err := groupAtAlice.CommitChanges(); err != nil {
		t.Fatal(err)
	}

	time.Sleep(5 * time.Second)
	if err := AppendToFile(fileBob, "Bob's modification\n"); err != nil {
		t.Fatal(err)
	}

	if err := groupAtBob.CommitChanges(); err != nil {
		t.Fatal(err)
	}

	time.Sleep(5 * time.Second)
	if err := AppendToFile(fileAlice, "Alice's modification\n"); err != nil {
		t.Fatal(err)
	}

	if err := groupAtAlice.CommitChanges(); err != nil {
		t.Fatal(err)
	}



	time.Sleep(500 * time.Second)

	fmt.Println(alice)
	for c := range alice.AddressBook.Iterator() {
		cc := c.(*Contact)
		cc.Send([]byte{2})
	}

	fmt.Println(bob)
	for c := range bob.AddressBook.Iterator() {
		cc := c.(*Contact)
		cc.Send([]byte{2})
	}
	fmt.Println(charlie)
	for c := range charlie.AddressBook.Iterator() {
		cc := c.(*Contact)
		cc.Send([]byte{2})
	}
}

func AppendToFile(path string, data string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, append(file, []byte(data)...), 644)
}