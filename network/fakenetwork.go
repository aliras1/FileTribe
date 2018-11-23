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
	"github.com/ethereum/go-ethereum/core/types"
)

type FakeNetwork struct {
	accCount 				int
	currentAcc 				int

	Client 					*eth.Eth
	Simulator 				*backends.SimulatedBackend

	auths    				[]*bind.TransactOpts

	// GroupInvitationEvent for users
	groupInvitationSubs     []event.Subscription
	groupInvitationChannels []chan *eth.EthGroupInvitation

	// GroupUpdateIpfs events for users
	groupUpdateIpfsSubs     []event.Subscription
	groupUpdateIpfsChannels []chan *eth.EthGroupUpdateIpfsHash

	groupRegisteredSubs     []event.Subscription
	groupRegisteredChannels []chan *eth.EthGroupRegistered

	keyDirtySubs     		[]event.Subscription
	keyDirtyChannels 		[]chan *eth.EthKeyDirty

	groupKeyChangedSubs     []event.Subscription
	groupKeyChangedChannels []chan *eth.EthGroupKeyChanged

	groupLeftSubs     		[]event.Subscription
	groupLeftChannels	 	[]chan *eth.EthGroupLeft

	// Debug event
	debugSubs     			[]event.Subscription
	debugChannels 			[]chan *eth.EthDebug
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
		alloc[auth.From] = core.GenesisAccount{Balance: big.NewInt(1000000000000000)}
	}

	simulator := backends.NewSimulatedBackend(core.GenesisAlloc(alloc), 9000000)

	_, _, ethclient, err := eth.DeployEth(auths[0], simulator)
	if err != nil {
		return nil, errors.Wrap(err, "could not deploy contract on simulated chain")
	}

	opts := &bind.WatchOpts{
		Context: auths[0].Context,
	}

	var channelGroupInvitations []chan *eth.EthGroupInvitation
	var subGroupInvitations []event.Subscription

	var channelGroupUpdateIpfss []chan *eth.EthGroupUpdateIpfsHash
	var subGroupUpdateIpfss []event.Subscription

	var channelGroupRegistered []chan *eth.EthGroupRegistered
	var subGroupGroupRegistered []event.Subscription

	var channelDebugs []chan *eth.EthDebug
	var subDebugs []event.Subscription

	var channelKeyDirtys []chan *eth.EthKeyDirty
	var subKeyDirtys []event.Subscription

	var groupKeyChangedChannels []chan *eth.EthGroupKeyChanged
	var groupKeyChangedSubs []event.Subscription

	var groupLeftChannels []chan *eth.EthGroupLeft
	var groupLeftSubs []event.Subscription

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
		chUpdtIpfs := make(chan *eth.EthGroupUpdateIpfsHash)
		subUpdtIpfs, err := ethclient.WatchGroupUpdateIpfsHash(opts, chUpdtIpfs)
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

		// KeyDirty event
		chKeyDirty := make(chan *eth.EthKeyDirty)
		subKeyDirty, err := ethclient.WatchKeyDirty(opts, chKeyDirty)
		if err != nil {
			return nil, err
		}
		channelKeyDirtys = append(channelKeyDirtys, chKeyDirty)
		subKeyDirtys = append(subKeyDirtys, subKeyDirty)

		// GroupKeyChanged event
		chGroupKeyChanged := make(chan *eth.EthGroupKeyChanged)
		subGroupKeyChanged, err := ethclient.WatchGroupKeyChanged(opts, chGroupKeyChanged)
		if err != nil {
			return nil, err
		}
		groupKeyChangedChannels = append(groupKeyChangedChannels, chGroupKeyChanged)
		groupKeyChangedSubs = append(groupKeyChangedSubs, subGroupKeyChanged)

		// GroupLeft event
		chGroupLeft := make(chan *eth.EthGroupLeft)
		subGroupLeft, err := ethclient.WatchGroupLeft(opts, chGroupLeft)
		if err != nil {
			return nil, err
		}
		groupLeftChannels = append(groupLeftChannels, chGroupLeft)
		groupLeftSubs = append(groupLeftSubs, subGroupLeft)
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

		groupKeyChangedChannels: groupKeyChangedChannels,
		groupKeyChangedSubs:	 groupKeyChangedSubs,

		debugSubs:               subDebugs,
		debugChannels:           channelDebugs,

		keyDirtyChannels:		 channelKeyDirtys,
		keyDirtySubs:			 subKeyDirtys,

		groupLeftChannels:		 groupLeftChannels,
		groupLeftSubs:			 groupLeftSubs,
	}

	return testNetwork, nil
}

func (network *FakeNetwork) TransactionReceipt(tx *types.Transaction) (*types.Receipt, error) {
	return network.Simulator.TransactionReceipt(network.auths[network.currentAcc].Context, tx.Hash())
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

func (network *FakeNetwork) GetGroupUpdateIpfsChannel() chan *eth.EthGroupUpdateIpfsHash {
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

func (network *FakeNetwork) GetKeyDirtySub() *event.Subscription {
	return &network.keyDirtySubs[network.currentAcc]
}
func (network *FakeNetwork) GetKeyDirtyChannel() chan *eth.EthKeyDirty {
	return network.keyDirtyChannels[network.currentAcc]
}

func (network *FakeNetwork) GetGroupLeftSub() *event.Subscription {
	return &network.groupLeftSubs[network.currentAcc]
}

func (network *FakeNetwork) GetGroupLeftChannel() chan *eth.EthGroupLeft {
	return network.groupLeftChannels[network.currentAcc]
}

func (network *FakeNetwork) GetGroupKeyChangedSub() *event.Subscription {
	return &network.groupKeyChangedSubs[network.currentAcc]
}

func (network *FakeNetwork) GetGroupKeyChangedChannel() chan *eth.EthGroupKeyChanged {
	return network.groupKeyChangedChannels[network.currentAcc]
}

func (network *FakeNetwork) IsUserRegistered(id common.Address) (bool, error) {
	registered, err := network.Client.IsUserRegistered(&bind.CallOpts{Pending: true}, id)
	if err != nil {
		return true, fmt.Errorf("error while FakeNetwork.IsUserRegistered(): %s", err)
	}
	return registered, nil
}


func (network *FakeNetwork) LeaveGroup(groupId [32]byte) (*Transaction, error) {

	tx, err := network.Client.LeaveGroup(network.auths[network.currentAcc], groupId)
	if err != nil {
		return nil, errors.Wrap(err, "could not send leaveGroup transaction")
	}

	network.Simulator.Commit()

	return &Transaction{tx: tx}, err
}

func (network *FakeNetwork) GetGroupLeader(groupId [32]byte) (common.Address, error) {
	leader, err := network.Client.GetLeader(&bind.CallOpts{Pending: true}, groupId)
	if err != nil {
		return common.Address{}, err
	}

	return leader, nil
}

func (network *FakeNetwork) UpdateGroupIpfsHash(groupId [32]byte, newIpfsHash []byte, approvals []*Approval) (*Transaction, error) {
	members, rs, ss, vs, err := prepareApprovals(approvals)
	if err != nil {
		return nil, errors.Wrap(err, "could not prepare approvals")
	}

	auth := network.auths[network.currentAcc]
	glog.Error("auth: " + auth.From.String())
	tx, err := network.Client.UpdateGroupIpfsHash(auth, groupId, newIpfsHash, members, rs, ss, vs)
	if err != nil {
		return nil, errors.Wrapf(err, "could not send updateGroupIpfsPath transaction")
	}
	glog.Error(tx.Nonce())

	network.Simulator.Commit()

	glog.Info("FakeNetwork.UpdateGroupIpfsHash ended")

	return &Transaction{tx: tx}, nil
}

func (network *FakeNetwork) ChangeGroupKey(groupId [32]byte, newIpfsHash []byte, approvals []*Approval) (*Transaction, error) {
	members, rs, ss, vs, err := prepareApprovals(approvals)
	if err != nil {
		return nil, errors.Wrap(err, "could not prepare approvals")
	}

	auth := network.auths[network.currentAcc]
	glog.Error("auth: " + auth.From.String())
	tx, err := network.Client.ChangeGroupKey(auth, groupId, newIpfsHash, members, rs, ss, vs)
	if err != nil {
		return nil, errors.Wrapf(err, "could not send changeGroupKey transaction")
	}
	glog.Error(tx.Nonce())

	network.Simulator.Commit()

	glog.Info("FakeNetwork.ChangeGroupKey ended")

	return &Transaction{tx: tx}, nil
}

func (network *FakeNetwork) RegisterUser(username, ipfsPeerId string, boxingKey [32]byte) (*Transaction, error) {
	auth := network.auths[network.currentAcc]
	auth.GasLimit = 2000000

	tx, err := network.Client.RegisterUser(auth, username, ipfsPeerId, boxingKey)
	if err != nil {
		return nil, fmt.Errorf("error while FakeNetwork.RegisterUser(): %s", err)
	}

	network.Simulator.Commit()

	return &Transaction{tx: tx}, nil
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
		Boxer:     crypto.AnonymPublicKey{boxingKey},
	}, nil
}

func (network *FakeNetwork) CreateGroup(id [32]byte, name string, ipfsHash []byte) (*Transaction, error) {
	tx, err := network.Client.CreateGroup(network.auths[network.currentAcc], id, name, ipfsHash)
	if err != nil {
		return nil, fmt.Errorf("error while FakeNetwork.CreateGroup(): %s", err)
	}

	network.Simulator.Commit()

	return &Transaction{tx: tx}, nil
}

func (network *FakeNetwork) GetGroup(groupId [32]byte) (string, []common.Address, []byte, common.Address, error) {
	groupData, err := network.Client.GetGroup(&bind.CallOpts{Pending: true}, groupId)
	return groupData.Name, groupData.Members, groupData.IpfsHash, groupData.Leader, err
}

func (network *FakeNetwork) InviteUser(groupId [32]byte, newMember common.Address, canInvite bool) (*Transaction, error) {
	tx, err := network.Client.InviteUser(network.auths[network.currentAcc], groupId, newMember, canInvite)
	if err != nil {
		return nil, fmt.Errorf("error while FakeNetwork.InviteUser(): %s", err)
	}

	network.Simulator.Commit()

	return &Transaction{tx: tx}, nil
}

func (network *FakeNetwork) Close() {

}