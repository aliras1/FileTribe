package client

import (
	"github.com/ethereum/go-ethereum/common"
	"ipfs-share/utils"
	"testing"

	"flag"
	"fmt"
	"github.com/golang/glog"
	"time"

	"crypto/ecdsa"
	"io/ioutil"
)


func TestGroupContext_Invite(t *testing.T) {
	flag.Set("alsologtostderr", fmt.Sprintf("%t", true))
	var logLevel string
	flag.StringVar(&logLevel, "-stderrthreshold", "INFO", "test")

	password := "pwd"
	dir := "../test/keystore/"

	ethKeyAlicePath := dir + "UTC--2018-10-10T08-19-58.398032114Z--ab083e63cfc7525634642075d49a0de31374bc0f"
	keyAlice, err := NewEthAccount(ethKeyAlicePath, password)
	if err != nil {
		t.Fatal(err)
	}

	ethKeyBobPath := dir + "UTC--2018-10-10T08-20-04.769949175Z--be9678b9882dac288093b9d38ea7382f21479c77"
	keyBob, err := NewEthAccount(ethKeyBobPath, password)
	if err != nil {
		t.Fatal(err)
	}

	ethKeyCharliePath := dir + "UTC--2018-10-10T08-20-10.903818650Z--d7ad6058180005d6639653f1d0216e481a43af79"
	keyCharlie, err := NewEthAccount(ethKeyCharliePath, password)
	if err != nil {
		t.Fatal(err)
	}

	sim, appAddr, err := createApp([]*ecdsa.PrivateKey{keyAlice, keyBob, keyCharlie})
	if err != nil {
		t.Fatal(err)
	}

	alice, err := NewTestCtx("alice", true, ethKeyAlicePath, sim, appAddr, "2000")
	if err != nil {
		t.Fatal(err)
	}

	bob, err := NewTestCtx("bob", true, ethKeyBobPath, sim, appAddr, "2001")
	if err != nil {
		t.Fatal(err)
	}

	charlie, err := NewTestCtx("charlie", true, ethKeyCharliePath, sim, appAddr, "2002")
	if err != nil {
		t.Fatal(err)
	}

	glog.Info("----- fun begins -----")

	if err := alice.CreateGroup("GRUPPE"); err != nil {
		t.Fatal(err)
	}
	sim.Commit()

	time.Sleep(1 * time.Second)

	if len(alice.Groups()) != 1 {
		t.Fatal("no groupAtAlice found by alice")
	}

	//bobUser := bob.account()
	//charlieUser := charlie.account()

	aliceGroups := alice.groups.ToList()
	groupAtAlice := aliceGroups[0].(*GroupContext)
	if err := groupAtAlice.Invite(bob.account.ContractAddress(), true); err != nil {
		t.Fatal(err)
	}
	if err := groupAtAlice.Invite(charlie.account.ContractAddress(), true); err != nil {
		t.Fatal(err)
	}
	sim.Commit()

	time.Sleep(1 * time.Second)

	if bob.invitations.Count() == 0 {
		t.Fatal("no invitations at bob")
	}
	if charlie.invitations.Count() == 0 {
		t.Fatal("no invitations at charlie")
	}

	if err := bob.AcceptInvitation(bob.invitations.Get(0).(common.Address)); err != nil {
		t.Fatal(err)
	}
	sim.Commit()
	if err := charlie.AcceptInvitation(bob.invitations.Get(0).(common.Address)); err != nil {
		t.Fatal(err)
	}
	sim.Commit()

	time.Sleep(5 * time.Second)

	if bob.groups.Count() < 1 {
		t.Fatal("no group by bob")
	}
	if charlie.groups.Count() < 1 {
		t.Fatal("no group by charlie")
	}

	time.Sleep(1 * time.Second)

	fmt.Println("----------- Alice init commit ------------")

	fileAlice := "./0xab083E63Cfc7525634642075d49A0DE31374bc0f/data/userdata/root/" + groupAtAlice.Group.Address().String() + "/rrrepo.go"
	if err := utils.CopyFile("./account.go", fileAlice); err != nil {
		t.Fatal(err)
	}
	if err := groupAtAlice.CommitChanges(); err != nil {
		t.Fatal(err)
	}
	sim.Commit()

	time.Sleep(2 * time.Second)

	sim.Commit()

	time.Sleep(5 * time.Second)


	fmt.Println("----------- Bob tries to change the file ------------")
	bobGroups := bob.Groups()
	groupAtBob := bobGroups[0].(*GroupContext)
	fileBob := "./0xbE9678b9882dAc288093b9D38ea7382f21479c77/data/userdata/root/" + groupAtBob.Group.Address().String() + "/rrrepo.go"
	if err := AppendToFile(fileBob, "Bob's modification (should fail)\n"); err != nil {
		t.Fatal(err)
	}

	if err := groupAtBob.CommitChanges(); err != nil {
		t.Fatal(err)
	}
	sim.Commit()
	time.Sleep(2 * time.Second)
	sim.Commit()
	time.Sleep(3 * time.Second)

	fmt.Println("----------- Grant W access to only alice  ------------")
	if err := groupAtAlice.GrantWriteAccess(fileAlice, bob.account.ContractAddress()); err != nil {
		t.Fatal(err)
	}
	if err := groupAtAlice.CommitChanges(); err != nil {
		t.Fatal(err)
	}
	sim.Commit()
	time.Sleep(2 * time.Second)
	sim.Commit()
	time.Sleep(3 * time.Second)


	fmt.Println("----------- Bob modif  ------------")
	if err := AppendToFile(fileBob, "Bob's modification\n"); err != nil {
		t.Fatal(err)
	}

	if err := groupAtBob.CommitChanges(); err != nil {
		t.Fatal(err)
	}
	sim.Commit()

	time.Sleep(5 * time.Second)
	//fmt.Println("----------- Alice modif  ------------")
	//fakeNetwork.SetAuth(ALICE)
	//if err := AppendToFile(fileAlice, "Alice's modification\n"); err != nil {
	//	t.Fatal(err)
	//}
	//
	//if err := groupAtAlice.CommitChanges(); err != nil {
	//	t.Fatal(err)
	//}
	//
	//time.Sleep(500 * time.Second)
}

func AppendToFile(path string, data string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, append(file, []byte(data)...), 644)
}