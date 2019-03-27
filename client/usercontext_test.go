package client

import (
	"crypto/ecdsa"
	"flag"
	"fmt"
	"io/ioutil"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/golang/glog"
	"github.com/pkg/errors"

	ethapp "github.com/aliras1/FileTribe/eth/gen/FileTribeDApp"
	"github.com/aliras1/FileTribe/eth/gen/factory/AccountFactory"
	"github.com/aliras1/FileTribe/eth/gen/factory/ConsensusFactory"
	"github.com/aliras1/FileTribe/eth/gen/factory/GroupFactory"
	ipfsapi "github.com/aliras1/FileTribe/ipfs"
)


func NewEthAccount(ethKeyPath, password string) (*ecdsa.PrivateKey, error) {
	json, err := ioutil.ReadFile(ethKeyPath)
	if err != nil {
		return nil, err
	}
	key, err := keystore.DecryptKey(json, password)
	if err != nil {
		return nil, err
	}
	return key.PrivateKey, nil
}

func createApp(keys []*ecdsa.PrivateKey) (*backends.SimulatedBackend, common.Address, error) {
	var auths []*bind.TransactOpts
	for _, key := range keys {
		auth := bind.NewKeyedTransactor(key)
		auth.GasLimit = 470000000

		auths = append(auths, auth)
	}

	alloc := make(map[common.Address]core.GenesisAccount)
	for _, auth := range auths {
		alloc[auth.From] = core.GenesisAccount{Balance: big.NewInt(10000000000000000)}
	}

	simulator := backends.NewSimulatedBackend(core.GenesisAlloc(alloc), 480000000)

	appAdrr, _, app, err := ethapp.DeployFileTribeDApp(auths[0], simulator)
	if err != nil {
		return nil, common.Address{}, errors.Wrap(err, "could not deploy app contract on simulated chain")
	}



	simulator.Commit()

	// AccountFactory
	accFactoryAddr, _, accFactory, err := AccountFactory.DeployAccountFactory(auths[0], simulator)
	if err != nil {
		return nil, common.Address{}, errors.Wrap(err, "could not deploy AccountFactory contract on simulated chain")
	}

	simulator.Commit()

	if _, err := accFactory.SetParent(auths[0], appAdrr); err != nil {
		return nil, common.Address{}, errors.Wrap(err, "could not set acc factory's parent")
	}

	simulator.Commit()

	if _, err := app.SetAccountFactory(auths[0], accFactoryAddr); err != nil {
		return nil, common.Address{}, errors.Wrap(err, "could not set app's acc factory")
	}

	simulator.Commit()

	// GroupFactory
	groupFactoryAddr, _, groupFactory, err := GroupFactory.DeployGroupFactory(auths[0], simulator)
	if err != nil {
		return nil, common.Address{}, errors.Wrap(err, "could not deploy GroupFactory contract on simulated chain")
	}

	simulator.Commit()

	if _, err := groupFactory.SetParent(auths[0], appAdrr); err != nil {
		return nil, common.Address{}, errors.Wrap(err, "could not set group factory's parent")
	}

	simulator.Commit()

	if _, err := app.SetGroupFactory(auths[0], groupFactoryAddr); err != nil {
		return nil, common.Address{}, errors.Wrap(err, "could not set app's group factory")
	}

	simulator.Commit()

	// ConsensusFactory
	consFactoryAddr, _, consFactory, err := ConsensusFactory.DeployConsensusFactory(auths[0], simulator)
	if err != nil {
		return nil, common.Address{}, errors.Wrap(err, "could not deploy ConsensusFactory contract on simulated chain")
	}

	simulator.Commit()

	if _, err := consFactory.SetParent(auths[0], appAdrr); err != nil {
		return nil, common.Address{}, errors.Wrap(err, "could not set cons factory's parent")
	}

	simulator.Commit()

	if _, err := app.SetConsensusFactory(auths[0], consFactoryAddr); err != nil {
		return nil, common.Address{}, errors.Wrap(err, "could not set app's cons factory")
	}

	simulator.Commit()

	return simulator, appAdrr, nil
}


func NewTestCtx(username string, signup bool, ethKeyPath string, sim *backends.SimulatedBackend, appAddr common.Address, p2pPort string) (*UserContext, error) {
	t := time.Now()
	glog.Info("ipfs inst: ", time.Since(t))
	password := "pwd"
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

	auth, err := NewAuth(ethKeyPath, password)
	if err != nil {
		panic(fmt.Sprintf("could not load account key data: NewNetwork: %s", err))
	}

	ctx, err := NewUserContext(auth, sim, appAddr, ipfs, p2pPort)

	if signup {
		err = ctx.SignUp(username)
		if err != nil {
			return nil, fmt.Errorf("could not sign up: %s: %s", username, err)
		}
		time.Sleep(2 * time.Second)
		sim.Commit()
	} else {
		//testUser, err = NewUserContextFromSignIn(username, password, ethKeyPath, homeDir, network, ipfs, p2pPort)
		//if err != nil {
		//	return nil, fmt.Errorf("could not sign in: %s: %s", username, err)
		//}
		return nil, errors.New("no signin")
	}

	return ctx, nil
}


func TestUserContext_SignUp(t *testing.T) {
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
	//bob, err := NewTestCtx("bob", true, ethKeyBobPath, sim, appAddr, "2001")
	//charles, err := NewTestCtx("charles", true, ethKeyCharliePath, sim, appAddr, "2002")
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	if alice.account == nil {
		t.Fatal("no account found by alice")
	}

	if err := alice.CreateGroup("gruppe"); err != nil {
		t.Fatal(err)
	}
	sim.Commit()

	time.Sleep(1 * time.Second)

	if alice.groups.Count() < 1 {
		t.Fatal("no groups found at alice")
	}
}

