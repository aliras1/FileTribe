package network

import (
	"ipfs-share/eth"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/core"
	"math/big"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"ipfs-share/crypto"
	"github.com/ethereum/go-ethereum/event"
	"github.com/pkg/errors"

	"github.com/golang/glog"
)

type FakeNetwork struct {
	Client *eth.Eth
	Auth    *bind.TransactOpts
	Simulator *backends.SimulatedBackend

	aliceAuth    *bind.TransactOpts
	bobAuth    *bind.TransactOpts
	charlieAuth    *bind.TransactOpts

	groupInvitationSub     event.Subscription
	groupInvitationChannel chan *eth.EthGroupInvitation
	groupUpdateIpfsSub     event.Subscription
	groupUpdateIpfsChannel chan *eth.EthGroupUpdateIpfsPath
	debugSub event.Subscription
	debugChannel chan *eth.EthDebug

	// GroupInvitationEvent for users
	groupInvitationSubAlice     event.Subscription
	groupInvitationChannelAlice chan *eth.EthGroupInvitation
	groupInvitationSubBob     event.Subscription
	groupInvitationChannelBob chan *eth.EthGroupInvitation
	groupInvitationSubCharlie     event.Subscription
	groupInvitationChannelCharlie chan *eth.EthGroupInvitation

	// GroupUpdateIpfs events for users
	groupUpdateIpfsSubAlice     event.Subscription
	groupUpdateIpfsChannelAlice chan *eth.EthGroupUpdateIpfsPath
	groupUpdateIpfsSubBob     event.Subscription
	groupUpdateIpfsChannelBob chan *eth.EthGroupUpdateIpfsPath
	groupUpdateIpfsSubCharlie     event.Subscription
	groupUpdateIpfsChannelCharlie chan *eth.EthGroupUpdateIpfsPath

	// Debug event
	debugSubAlice     event.Subscription
	debugChannelAlice chan *eth.EthDebug
	debugSubBob     event.Subscription
	debugChannelBob chan *eth.EthDebug
	debugSubCharlie     event.Subscription
	debugChannelCharlie chan *eth.EthDebug
}

func (network *FakeNetwork) SetAuthAlice() {
	network.Auth = network.aliceAuth

	network.groupInvitationSub = network.groupInvitationSubAlice
	network.groupInvitationChannel = network.groupInvitationChannelAlice

	network.groupUpdateIpfsSub = network.groupUpdateIpfsSubAlice
	network.groupUpdateIpfsChannel = network.groupUpdateIpfsChannelAlice

	network.debugSub = network.debugSubAlice
	network.debugChannel = network.debugChannelAlice
}

func (network *FakeNetwork) SetAuthBob() {
	network.Auth = network.bobAuth

	network.groupInvitationSub = network.groupInvitationSubBob
	network.groupInvitationChannel = network.groupInvitationChannelBob

	network.groupUpdateIpfsSub = network.groupUpdateIpfsSubBob
	network.groupUpdateIpfsChannel = network.groupUpdateIpfsChannelBob

	network.debugSub = network.debugSubBob
	network.debugChannel = network.debugChannelBob
}

func (network *FakeNetwork) SetAuthCharlie() {
	network.Auth = network.charlieAuth

	network.groupInvitationSub = network.groupInvitationSubCharlie
	network.groupInvitationChannel = network.groupInvitationChannelCharlie

	network.groupUpdateIpfsSub = network.groupUpdateIpfsSubCharlie
	network.groupUpdateIpfsChannel = network.groupUpdateIpfsChannelCharlie

	network.debugSub = network.debugSubCharlie
	network.debugChannel = network.debugChannelCharlie
}

func NewTestNetwork(keyAlice, keyBob, keyCharlie *ecdsa.PrivateKey) (*FakeNetwork, error) {
	authAlice := bind.NewKeyedTransactor(keyAlice)
	//authAlice.GasLimit = 2141592

	authBob := bind.NewKeyedTransactor(keyBob)
	//authBob.GasLimit = 2141592

	authCharlie := bind.NewKeyedTransactor(keyCharlie)
	//authCharlie.GasLimit = 2141592

	simulator := backends.NewSimulatedBackend(core.GenesisAlloc{
		authAlice.From: core.GenesisAccount{Balance: big.NewInt(10000000000)},
		authBob.From:   core.GenesisAccount{Balance: big.NewInt(10000000000)},
		authCharlie.From:   core.GenesisAccount{Balance: big.NewInt(10000000000)},

	}, )

	_, _, ethclient, err := eth.DeployEth(authAlice, simulator)
	if err != nil {
		return nil, fmt.Errorf("could not deploy contract on simulated chain")
	}

	opts := &bind.WatchOpts{
		Context: authAlice.Context,
	}

	// Subscribing to events for each user...

	// GroupInvitation event
	channelGroupInvitationAlice := make(chan *eth.EthGroupInvitation)
	subGroupInvitationAlice, err := ethclient.WatchGroupInvitation(opts, channelGroupInvitationAlice)
	if err != nil {
		return nil, err
	}

	channelGroupInvitationBob := make(chan *eth.EthGroupInvitation)
	subGroupInvitationBob, err := ethclient.WatchGroupInvitation(opts, channelGroupInvitationBob)
	if err != nil {
		return nil, err
	}

	channelGroupInvitationCharlie := make(chan *eth.EthGroupInvitation)
	subGroupInvitationCharlie, err := ethclient.WatchGroupInvitation(opts, channelGroupInvitationCharlie)
	if err != nil {
		return nil, err
	}

	// GroupUpdateIpfs event
	channelGroupUpdateIpfsAlice := make(chan *eth.EthGroupUpdateIpfsPath)
	subGroupUpdateIpfsAlice, err := ethclient.WatchGroupUpdateIpfsPath(opts, channelGroupUpdateIpfsAlice)
	if err != nil {
		return nil, err
	}

	channelGroupUpdateIpfsBob := make(chan *eth.EthGroupUpdateIpfsPath)
	subGroupUpdateIpfsBob, err := ethclient.WatchGroupUpdateIpfsPath(opts, channelGroupUpdateIpfsBob)
	if err != nil {
		return nil, err
	}
	channelGroupUpdateIpfsCharlie := make(chan *eth.EthGroupUpdateIpfsPath)
	subGroupUpdateIpfsCharlie, err := ethclient.WatchGroupUpdateIpfsPath(opts, channelGroupUpdateIpfsCharlie)
	if err != nil {
		return nil, err
	}

	// Debug event
	channelDebugAlice := make(chan *eth.EthDebug)
	subGDebugAlice, err := ethclient.WatchDebug(opts, channelDebugAlice)
	if err != nil {
		return nil, err
	}

	channelDebugBob := make(chan *eth.EthDebug)
	subDebugBob, err := ethclient.WatchDebug(opts, channelDebugBob)
	if err != nil {
		return nil, err
	}
	channelDebugCharlie := make(chan *eth.EthDebug)
	subDebugCharlie, err := ethclient.WatchDebug(opts, channelDebugCharlie)
	if err != nil {
		return nil, err
	}


	testNetwork := &FakeNetwork{
		Client: ethclient,
		Auth: nil,
		Simulator: simulator,

		aliceAuth: authAlice,
		bobAuth: authBob,
		charlieAuth: authCharlie,

		groupInvitationChannelAlice: channelGroupInvitationAlice,
		groupInvitationSubAlice:     subGroupInvitationAlice,
		groupUpdateIpfsChannelAlice: channelGroupUpdateIpfsAlice,
		groupUpdateIpfsSubAlice: subGroupUpdateIpfsAlice,
		debugSubAlice: subGDebugAlice,
		debugChannelAlice: channelDebugAlice,

		groupInvitationChannelBob: channelGroupInvitationBob,
		groupInvitationSubBob:     subGroupInvitationBob,
		groupUpdateIpfsChannelBob: channelGroupUpdateIpfsBob,
		groupUpdateIpfsSubBob: subGroupUpdateIpfsBob,
		debugSubBob: subDebugBob,
		debugChannelBob: channelDebugBob,

		groupInvitationChannelCharlie: channelGroupInvitationCharlie,
		groupInvitationSubCharlie: subGroupInvitationCharlie,
		groupUpdateIpfsChannelCharlie: channelGroupUpdateIpfsCharlie,
		groupUpdateIpfsSubCharlie: subGroupUpdateIpfsCharlie,
		debugSubCharlie: subDebugCharlie,
		debugChannelCharlie: channelDebugCharlie,
	}

	return testNetwork, nil
}

func (network *FakeNetwork) GetGroupInvitationSub() *event.Subscription {
	return &network.groupInvitationSub
}

func (network *FakeNetwork) GetGroupInvitationChannel() chan *eth.EthGroupInvitation {
	return network.groupInvitationChannel
}

func (network *FakeNetwork) GetGroupUpdateIpfsSub() *event.Subscription {
	return &network.groupUpdateIpfsSub
}

func (network *FakeNetwork) GetGroupUpdateIpfsChannel() chan *eth.EthGroupUpdateIpfsPath {
	return network.groupUpdateIpfsChannel
}

func (network *FakeNetwork) GetDebugSub() *event.Subscription {
	return &network.debugSub
}

func (network *FakeNetwork) GetDebugChannel() chan *eth.EthDebug {
	return network.debugChannel
}

func (network *FakeNetwork) IsUserRegistered(id common.Address) (bool, error) {
	registered, err := network.Client.IsUserRegistered(&bind.CallOpts{Pending: true}, id)
	if err != nil {
		return true, fmt.Errorf("error while FakeNetwork.IsUserRegistered(): %s", err)
	}
	return registered, nil
}



func (network *FakeNetwork) UpdateGroupIpfsPath(groupId [32]byte, newIpfsPath string, approvals []*Approval) error {
	//network.Simulator.EstimateGas(network.Auth.Context, )

	var members []common.Address
	var rs [][32]byte
	var ss [][32]byte
	var vs []uint8

	for _, approval := range approvals {
		if len(approval.Signature) != 65 {
			return errors.New("signature length must be 65")
		}

		members = append(members, approval.From)

		var r [32]byte
		copy(r[:], approval.Signature[:32])
		rs = append(rs, r)

		var s [32]byte
		copy(s[:], approval.Signature[32:64])
		ss = append(ss, s)

		v := approval.Signature[64]
		if v < 27 {
			v = v + 27;
		}
		vs = append(vs, v)
	}

	auth := network.Auth
	auth.GasLimit = 2541592
	_, err := network.Client.UpdateGroupIpfsPath(auth, groupId, newIpfsPath, members, rs, ss, vs)
	if err != nil {
		return errors.Wrapf(err, "could not send updateGroupIpfsPath transaction")
	}

	network.Simulator.Commit()

	glog.Info("FakeNetwork.UpdateGroupIpfsPath ended")

	return nil
}

func (network *FakeNetwork) RegisterUser(username, ipfsPeerId string, boxingKey [32]byte) error {
	//auth := network.Auth
	//auth.GasLimit = 2141592

	_, err := network.Client.RegisterUser(network.Auth, username, ipfsPeerId, boxingKey)
	if err != nil {
		return fmt.Errorf("error while FakeNetwork.RegisterUser(): %s", err)
	}

	network.Simulator.Commit()

	return nil
}

func (network *FakeNetwork) GetUser(address common.Address) (*Contact, error) {
	username, ipfsPeerId, boxingKey, err := network.Client.GetUser(&bind.CallOpts{Pending: true}, address)
	if err != nil {
		return &Contact{}, fmt.Errorf("error while FakeNetwork.GetUser(): %s", err)
	}

	return &Contact{
		Address:   address,
		Name:      username,
		IpfsPeerId: ipfsPeerId,
		Boxer:     crypto.AnonymPublicKey{&boxingKey},
	}, nil
}

func (network *FakeNetwork) CreateGroup(id [32]byte, name string, ipfsPath string) error {
	_, err := network.Client.CreateGroup(network.Auth, id, name, ipfsPath)
	if err != nil {
		return fmt.Errorf("error while FakeNetwork.CreateGroup(): %s", err)
	}

	network.Simulator.Commit()

	return nil
}

func (network *FakeNetwork) GetGroup(groupId [32]byte) (string, []common.Address, string, error) {
	return network.Client.GetGroup(&bind.CallOpts{Pending: true}, groupId)
}

func (network *FakeNetwork) InviteUser(groupId [32]byte, newMember common.Address) error {
	_, err := network.Client.InviteUser(network.Auth, groupId, newMember)
	if err != nil {
		return fmt.Errorf("error while FakeNetwork.InviteUser(): %s", err)
	}

	network.Simulator.Commit()

	return nil
}

func (network *FakeNetwork) Close() {

}