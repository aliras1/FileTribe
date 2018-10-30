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
	accCount int
	currentAcc int

	Client *eth.Eth
	Simulator *backends.SimulatedBackend

	auths    []*bind.TransactOpts

	// GroupInvitationEvent for users
	groupInvitationSubs     []event.Subscription
	groupInvitationChannels []chan *eth.EthGroupInvitation

	// GroupUpdateIpfs events for users
	groupUpdateIpfsSubs     []event.Subscription
	groupUpdateIpfsChannels []chan *eth.EthGroupUpdateIpfsPath

	groupRegisteredSubs     []event.Subscription
	groupRegisteredChannels []chan *eth.EthGroupRegistered

	// Debug event
	debugSubs     []event.Subscription
	debugChannels []chan *eth.EthDebug
}

func (network *FakeNetwork) SetAuth(accNum int) {
	if accNum >= len (network.auths) {
		glog.Error("invalid accNum: out of range")
	}

	network.currentAcc = accNum
}

func NewTestNetwork(keys []*ecdsa.PrivateKey) (*FakeNetwork, error) {
	accCount := len(keys)

	var auths []*bind.TransactOpts
	for _, key := range keys {
		auth := bind.NewKeyedTransactor(key)
		auth.GasLimit = 5000000

		auths = append(auths, auth)
	}

	alloc := make(map[common.Address]core.GenesisAccount)
	for _, auth := range auths {
		alloc[auth.From] = core.GenesisAccount{Balance: big.NewInt(10000000000)}
	}

	simulator := backends.NewSimulatedBackend(core.GenesisAlloc(alloc), 7000001)

	_, _, ethclient, err := eth.DeployEth(auths[0], simulator)
	if err != nil {
		return nil, errors.Wrap(err, "could not deploy contract on simulated chain")
	}

	opts := &bind.WatchOpts{
		Context: auths[0].Context,
	}

	var channelGroupInvitations []chan *eth.EthGroupInvitation
	var subGroupInvitations []event.Subscription
	var channelGroupUpdateIpfss []chan *eth.EthGroupUpdateIpfsPath
	var subGroupUpdateIpfss []event.Subscription
	var channelGroupRegistered []chan *eth.EthGroupRegistered
	var subGroupGroupRegistered []event.Subscription
	var channelDebugs []chan *eth.EthDebug
	var subDebugs []event.Subscription

	for i := 0; i < accCount; i++ {
		// Subscribing to events for each user...
		chGroupInv := make(chan *eth.EthGroupInvitation)
		subGroupInv, err := ethclient.WatchGroupInvitation(opts, chGroupInv)
		if err != nil {
			return nil, err
		}
		channelGroupInvitations = append(channelGroupInvitations, chGroupInv)
		subGroupInvitations = append(subGroupInvitations, subGroupInv)

		// GroupUpdateIpfs event
		chUpdtIpfs := make(chan *eth.EthGroupUpdateIpfsPath)
		subUpdtIpfs, err := ethclient.WatchGroupUpdateIpfsPath(opts, chUpdtIpfs)
		if err != nil {
			return nil, err
		}
		channelGroupUpdateIpfss = append(channelGroupUpdateIpfss, chUpdtIpfs)
		subGroupUpdateIpfss = append(subGroupUpdateIpfss, subUpdtIpfs)

		// GroupRegistered event
		chGrpRegistered := make(chan *eth.EthGroupRegistered)
		subGrpRegistered, err := ethclient.WatchGroupRegistered(opts, chGrpRegistered)
		if err != nil {
			return nil, err
		}
		channelGroupRegistered = append(channelGroupRegistered, chGrpRegistered)
		subGroupGroupRegistered = append(subGroupGroupRegistered, subGrpRegistered)

		// Debug event
		chDebug := make(chan *eth.EthDebug)
		subDebug, err := ethclient.WatchDebug(opts, chDebug)
		if err != nil {
			return nil, err
		}
		channelDebugs = append(channelDebugs, chDebug)
		subDebugs = append(subDebugs, subDebug)
	}


	testNetwork := &FakeNetwork{
		Client: ethclient,
		Simulator: simulator,

		auths: auths,

		groupInvitationChannels: channelGroupInvitations,
		groupInvitationSubs:     subGroupInvitations,
		groupUpdateIpfsChannels: channelGroupUpdateIpfss,
		groupUpdateIpfsSubs:     subGroupUpdateIpfss,
		groupRegisteredChannels: channelGroupRegistered,
		groupRegisteredSubs:     subGroupGroupRegistered,
		debugSubs:               subDebugs,
		debugChannels:           channelDebugs,
	}

	return testNetwork, nil
}

func (network *FakeNetwork) GetGroupInvitationSub() *event.Subscription {
	return &network.groupInvitationSubs[network.currentAcc]
}

func (network *FakeNetwork) GetGroupInvitationChannel() chan *eth.EthGroupInvitation {
	return network.groupInvitationChannels[network.currentAcc]
}

func (network *FakeNetwork) GetGroupUpdateIpfsSub() *event.Subscription {
	return &network.groupUpdateIpfsSubs[network.currentAcc]
}

func (network *FakeNetwork) GetGroupUpdateIpfsChannel() chan *eth.EthGroupUpdateIpfsPath {
	return network.groupUpdateIpfsChannels[network.currentAcc]
}

func (network *FakeNetwork) GetGroupRegisteredSub() *event.Subscription {
	return &network.groupRegisteredSubs[network.currentAcc]
}

func (network *FakeNetwork) GetGroupRegisteredChannel() chan *eth.EthGroupRegistered {
	return network.groupRegisteredChannels[network.currentAcc]
}

func (network *FakeNetwork) GetDebugSub() *event.Subscription {
	return &network.debugSubs[network.currentAcc]
}

func (network *FakeNetwork) GetDebugChannel() chan *eth.EthDebug {
	return network.debugChannels[network.currentAcc]
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

	auth := network.auths[network.currentAcc]
	auth.GasLimit = 3000000
	_, err := network.Client.UpdateGroupIpfsPath(auth, groupId, newIpfsPath, members, rs, ss, vs)
	if err != nil {
		return errors.Wrapf(err, "could not send updateGroupIpfsPath transaction")
	}

	network.Simulator.Commit()

	glog.Info("FakeNetwork.UpdateGroupIpfsPath ended")

	return nil
}

func (network *FakeNetwork) RegisterUser(username, ipfsPeerId string, boxingKey [32]byte) error {
	auth := network.auths[network.currentAcc]
	auth.GasLimit = 3000000

	_, err := network.Client.RegisterUser(auth, username, ipfsPeerId, boxingKey)
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
	_, err := network.Client.CreateGroup(network.auths[network.currentAcc], id, name, ipfsPath)
	if err != nil {
		return fmt.Errorf("error while FakeNetwork.CreateGroup(): %s", err)
	}

	network.Simulator.Commit()

	return nil
}

func (network *FakeNetwork) GetGroup(groupId [32]byte) (string, []common.Address, string, error) {
	return network.Client.GetGroup(&bind.CallOpts{Pending: true}, groupId)
}

func (network *FakeNetwork) InviteUser(groupId [32]byte, newMember common.Address, canInvite bool) error {
	_, err := network.Client.InviteUser(network.auths[network.currentAcc], groupId, newMember, canInvite)
	if err != nil {
		return fmt.Errorf("error while FakeNetwork.InviteUser(): %s", err)
	}

	network.Simulator.Commit()

	return nil
}

func (network *FakeNetwork) Close() {

}