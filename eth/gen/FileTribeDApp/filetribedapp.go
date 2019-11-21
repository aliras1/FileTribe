// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package FileTribeDApp

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// FileTribeDAppABI is the input ABI used to generate the binding from.
const FileTribeDAppABI = "[{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractIAccount\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"AccountCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"msg\",\"type\":\"int256\"}],\"name\":\"Debug\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"msg\",\"type\":\"bytes\"}],\"name\":\"DebugBytes\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"GroupRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"group\",\"type\":\"address\"}],\"name\":\"onInvitationAccepted\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"onInvitationDeclined\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"contractIConsensusFactory\",\"name\":\"factory\",\"type\":\"address\"}],\"name\":\"setConsensusFactory\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"contractIAccountFactory\",\"name\":\"factory\",\"type\":\"address\"}],\"name\":\"setAccountFactory\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"contractIGroupFactory\",\"name\":\"factory\",\"type\":\"address\"}],\"name\":\"setGroupFactory\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"ipfsPeerId\",\"type\":\"string\"},{\"internalType\":\"bytes32\",\"name\":\"boxingKey\",\"type\":\"bytes32\"}],\"name\":\"createAccount\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"removeAccount\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"createGroup\",\"outputs\":[{\"internalType\":\"contractIGroup\",\"name\":\"group\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"contractIAccount\",\"name\":\"proposer\",\"type\":\"address\"}],\"name\":\"createConsensus\",\"outputs\":[{\"internalType\":\"contractIConsensus\",\"name\":\"consensus\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"getAccountOf\",\"outputs\":[{\"internalType\":\"contractIAccount\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// FileTribeDApp is an auto generated Go binding around an Ethereum contract.
type FileTribeDApp struct {
	FileTribeDAppCaller     // Read-only binding to the contract
	FileTribeDAppTransactor // Write-only binding to the contract
	FileTribeDAppFilterer   // Log filterer for contract events
}

// FileTribeDAppCaller is an auto generated read-only Go binding around an Ethereum contract.
type FileTribeDAppCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FileTribeDAppTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FileTribeDAppTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FileTribeDAppFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FileTribeDAppFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FileTribeDAppSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FileTribeDAppSession struct {
	Contract     *FileTribeDApp    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FileTribeDAppCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FileTribeDAppCallerSession struct {
	Contract *FileTribeDAppCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// FileTribeDAppTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FileTribeDAppTransactorSession struct {
	Contract     *FileTribeDAppTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// FileTribeDAppRaw is an auto generated low-level Go binding around an Ethereum contract.
type FileTribeDAppRaw struct {
	Contract *FileTribeDApp // Generic contract binding to access the raw methods on
}

// FileTribeDAppCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FileTribeDAppCallerRaw struct {
	Contract *FileTribeDAppCaller // Generic read-only contract binding to access the raw methods on
}

// FileTribeDAppTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FileTribeDAppTransactorRaw struct {
	Contract *FileTribeDAppTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFileTribeDApp creates a new instance of FileTribeDApp, bound to a specific deployed contract.
func NewFileTribeDApp(address common.Address, backend bind.ContractBackend) (*FileTribeDApp, error) {
	contract, err := bindFileTribeDApp(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FileTribeDApp{FileTribeDAppCaller: FileTribeDAppCaller{contract: contract}, FileTribeDAppTransactor: FileTribeDAppTransactor{contract: contract}, FileTribeDAppFilterer: FileTribeDAppFilterer{contract: contract}}, nil
}

// NewFileTribeDAppCaller creates a new read-only instance of FileTribeDApp, bound to a specific deployed contract.
func NewFileTribeDAppCaller(address common.Address, caller bind.ContractCaller) (*FileTribeDAppCaller, error) {
	contract, err := bindFileTribeDApp(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FileTribeDAppCaller{contract: contract}, nil
}

// NewFileTribeDAppTransactor creates a new write-only instance of FileTribeDApp, bound to a specific deployed contract.
func NewFileTribeDAppTransactor(address common.Address, transactor bind.ContractTransactor) (*FileTribeDAppTransactor, error) {
	contract, err := bindFileTribeDApp(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FileTribeDAppTransactor{contract: contract}, nil
}

// NewFileTribeDAppFilterer creates a new log filterer instance of FileTribeDApp, bound to a specific deployed contract.
func NewFileTribeDAppFilterer(address common.Address, filterer bind.ContractFilterer) (*FileTribeDAppFilterer, error) {
	contract, err := bindFileTribeDApp(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FileTribeDAppFilterer{contract: contract}, nil
}

// bindFileTribeDApp binds a generic wrapper to an already deployed contract.
func bindFileTribeDApp(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(FileTribeDAppABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FileTribeDApp *FileTribeDAppRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _FileTribeDApp.Contract.FileTribeDAppCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FileTribeDApp *FileTribeDAppRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.FileTribeDAppTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FileTribeDApp *FileTribeDAppRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.FileTribeDAppTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FileTribeDApp *FileTribeDAppCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _FileTribeDApp.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FileTribeDApp *FileTribeDAppTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FileTribeDApp *FileTribeDAppTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.contract.Transact(opts, method, params...)
}

// GetAccountOf is a free data retrieval call binding the contract method 0x95184d3b.
//
// Solidity: function getAccountOf(address owner) constant returns(address)
func (_FileTribeDApp *FileTribeDAppCaller) GetAccountOf(opts *bind.CallOpts, owner common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _FileTribeDApp.contract.Call(opts, out, "getAccountOf", owner)
	return *ret0, err
}

// GetAccountOf is a free data retrieval call binding the contract method 0x95184d3b.
//
// Solidity: function getAccountOf(address owner) constant returns(address)
func (_FileTribeDApp *FileTribeDAppSession) GetAccountOf(owner common.Address) (common.Address, error) {
	return _FileTribeDApp.Contract.GetAccountOf(&_FileTribeDApp.CallOpts, owner)
}

// GetAccountOf is a free data retrieval call binding the contract method 0x95184d3b.
//
// Solidity: function getAccountOf(address owner) constant returns(address)
func (_FileTribeDApp *FileTribeDAppCallerSession) GetAccountOf(owner common.Address) (common.Address, error) {
	return _FileTribeDApp.Contract.GetAccountOf(&_FileTribeDApp.CallOpts, owner)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_FileTribeDApp *FileTribeDAppCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _FileTribeDApp.contract.Call(opts, out, "isOwner")
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_FileTribeDApp *FileTribeDAppSession) IsOwner() (bool, error) {
	return _FileTribeDApp.Contract.IsOwner(&_FileTribeDApp.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_FileTribeDApp *FileTribeDAppCallerSession) IsOwner() (bool, error) {
	return _FileTribeDApp.Contract.IsOwner(&_FileTribeDApp.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_FileTribeDApp *FileTribeDAppCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _FileTribeDApp.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_FileTribeDApp *FileTribeDAppSession) Owner() (common.Address, error) {
	return _FileTribeDApp.Contract.Owner(&_FileTribeDApp.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_FileTribeDApp *FileTribeDAppCallerSession) Owner() (common.Address, error) {
	return _FileTribeDApp.Contract.Owner(&_FileTribeDApp.CallOpts)
}

// CreateAccount is a paid mutator transaction binding the contract method 0x3c2aea2e.
//
// Solidity: function createAccount(string name, string ipfsPeerId, bytes32 boxingKey) returns()
func (_FileTribeDApp *FileTribeDAppTransactor) CreateAccount(opts *bind.TransactOpts, name string, ipfsPeerId string, boxingKey [32]byte) (*types.Transaction, error) {
	return _FileTribeDApp.contract.Transact(opts, "createAccount", name, ipfsPeerId, boxingKey)
}

// CreateAccount is a paid mutator transaction binding the contract method 0x3c2aea2e.
//
// Solidity: function createAccount(string name, string ipfsPeerId, bytes32 boxingKey) returns()
func (_FileTribeDApp *FileTribeDAppSession) CreateAccount(name string, ipfsPeerId string, boxingKey [32]byte) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.CreateAccount(&_FileTribeDApp.TransactOpts, name, ipfsPeerId, boxingKey)
}

// CreateAccount is a paid mutator transaction binding the contract method 0x3c2aea2e.
//
// Solidity: function createAccount(string name, string ipfsPeerId, bytes32 boxingKey) returns()
func (_FileTribeDApp *FileTribeDAppTransactorSession) CreateAccount(name string, ipfsPeerId string, boxingKey [32]byte) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.CreateAccount(&_FileTribeDApp.TransactOpts, name, ipfsPeerId, boxingKey)
}

// CreateConsensus is a paid mutator transaction binding the contract method 0x15db2069.
//
// Solidity: function createConsensus(address proposer) returns(address consensus)
func (_FileTribeDApp *FileTribeDAppTransactor) CreateConsensus(opts *bind.TransactOpts, proposer common.Address) (*types.Transaction, error) {
	return _FileTribeDApp.contract.Transact(opts, "createConsensus", proposer)
}

// CreateConsensus is a paid mutator transaction binding the contract method 0x15db2069.
//
// Solidity: function createConsensus(address proposer) returns(address consensus)
func (_FileTribeDApp *FileTribeDAppSession) CreateConsensus(proposer common.Address) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.CreateConsensus(&_FileTribeDApp.TransactOpts, proposer)
}

// CreateConsensus is a paid mutator transaction binding the contract method 0x15db2069.
//
// Solidity: function createConsensus(address proposer) returns(address consensus)
func (_FileTribeDApp *FileTribeDAppTransactorSession) CreateConsensus(proposer common.Address) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.CreateConsensus(&_FileTribeDApp.TransactOpts, proposer)
}

// CreateGroup is a paid mutator transaction binding the contract method 0xdc2ddcae.
//
// Solidity: function createGroup(string name) returns(address group)
func (_FileTribeDApp *FileTribeDAppTransactor) CreateGroup(opts *bind.TransactOpts, name string) (*types.Transaction, error) {
	return _FileTribeDApp.contract.Transact(opts, "createGroup", name)
}

// CreateGroup is a paid mutator transaction binding the contract method 0xdc2ddcae.
//
// Solidity: function createGroup(string name) returns(address group)
func (_FileTribeDApp *FileTribeDAppSession) CreateGroup(name string) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.CreateGroup(&_FileTribeDApp.TransactOpts, name)
}

// CreateGroup is a paid mutator transaction binding the contract method 0xdc2ddcae.
//
// Solidity: function createGroup(string name) returns(address group)
func (_FileTribeDApp *FileTribeDAppTransactorSession) CreateGroup(name string) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.CreateGroup(&_FileTribeDApp.TransactOpts, name)
}

// OnInvitationAccepted is a paid mutator transaction binding the contract method 0x0a561621.
//
// Solidity: function onInvitationAccepted(address group) returns()
func (_FileTribeDApp *FileTribeDAppTransactor) OnInvitationAccepted(opts *bind.TransactOpts, group common.Address) (*types.Transaction, error) {
	return _FileTribeDApp.contract.Transact(opts, "onInvitationAccepted", group)
}

// OnInvitationAccepted is a paid mutator transaction binding the contract method 0x0a561621.
//
// Solidity: function onInvitationAccepted(address group) returns()
func (_FileTribeDApp *FileTribeDAppSession) OnInvitationAccepted(group common.Address) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.OnInvitationAccepted(&_FileTribeDApp.TransactOpts, group)
}

// OnInvitationAccepted is a paid mutator transaction binding the contract method 0x0a561621.
//
// Solidity: function onInvitationAccepted(address group) returns()
func (_FileTribeDApp *FileTribeDAppTransactorSession) OnInvitationAccepted(group common.Address) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.OnInvitationAccepted(&_FileTribeDApp.TransactOpts, group)
}

// OnInvitationDeclined is a paid mutator transaction binding the contract method 0xb061d9a9.
//
// Solidity: function onInvitationDeclined() returns()
func (_FileTribeDApp *FileTribeDAppTransactor) OnInvitationDeclined(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FileTribeDApp.contract.Transact(opts, "onInvitationDeclined")
}

// OnInvitationDeclined is a paid mutator transaction binding the contract method 0xb061d9a9.
//
// Solidity: function onInvitationDeclined() returns()
func (_FileTribeDApp *FileTribeDAppSession) OnInvitationDeclined() (*types.Transaction, error) {
	return _FileTribeDApp.Contract.OnInvitationDeclined(&_FileTribeDApp.TransactOpts)
}

// OnInvitationDeclined is a paid mutator transaction binding the contract method 0xb061d9a9.
//
// Solidity: function onInvitationDeclined() returns()
func (_FileTribeDApp *FileTribeDAppTransactorSession) OnInvitationDeclined() (*types.Transaction, error) {
	return _FileTribeDApp.Contract.OnInvitationDeclined(&_FileTribeDApp.TransactOpts)
}

// RemoveAccount is a paid mutator transaction binding the contract method 0x75beca08.
//
// Solidity: function removeAccount() returns()
func (_FileTribeDApp *FileTribeDAppTransactor) RemoveAccount(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FileTribeDApp.contract.Transact(opts, "removeAccount")
}

// RemoveAccount is a paid mutator transaction binding the contract method 0x75beca08.
//
// Solidity: function removeAccount() returns()
func (_FileTribeDApp *FileTribeDAppSession) RemoveAccount() (*types.Transaction, error) {
	return _FileTribeDApp.Contract.RemoveAccount(&_FileTribeDApp.TransactOpts)
}

// RemoveAccount is a paid mutator transaction binding the contract method 0x75beca08.
//
// Solidity: function removeAccount() returns()
func (_FileTribeDApp *FileTribeDAppTransactorSession) RemoveAccount() (*types.Transaction, error) {
	return _FileTribeDApp.Contract.RemoveAccount(&_FileTribeDApp.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FileTribeDApp *FileTribeDAppTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FileTribeDApp.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FileTribeDApp *FileTribeDAppSession) RenounceOwnership() (*types.Transaction, error) {
	return _FileTribeDApp.Contract.RenounceOwnership(&_FileTribeDApp.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FileTribeDApp *FileTribeDAppTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _FileTribeDApp.Contract.RenounceOwnership(&_FileTribeDApp.TransactOpts)
}

// SetAccountFactory is a paid mutator transaction binding the contract method 0xaddc1a76.
//
// Solidity: function setAccountFactory(address factory) returns()
func (_FileTribeDApp *FileTribeDAppTransactor) SetAccountFactory(opts *bind.TransactOpts, factory common.Address) (*types.Transaction, error) {
	return _FileTribeDApp.contract.Transact(opts, "setAccountFactory", factory)
}

// SetAccountFactory is a paid mutator transaction binding the contract method 0xaddc1a76.
//
// Solidity: function setAccountFactory(address factory) returns()
func (_FileTribeDApp *FileTribeDAppSession) SetAccountFactory(factory common.Address) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.SetAccountFactory(&_FileTribeDApp.TransactOpts, factory)
}

// SetAccountFactory is a paid mutator transaction binding the contract method 0xaddc1a76.
//
// Solidity: function setAccountFactory(address factory) returns()
func (_FileTribeDApp *FileTribeDAppTransactorSession) SetAccountFactory(factory common.Address) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.SetAccountFactory(&_FileTribeDApp.TransactOpts, factory)
}

// SetConsensusFactory is a paid mutator transaction binding the contract method 0x1c38ea56.
//
// Solidity: function setConsensusFactory(address factory) returns()
func (_FileTribeDApp *FileTribeDAppTransactor) SetConsensusFactory(opts *bind.TransactOpts, factory common.Address) (*types.Transaction, error) {
	return _FileTribeDApp.contract.Transact(opts, "setConsensusFactory", factory)
}

// SetConsensusFactory is a paid mutator transaction binding the contract method 0x1c38ea56.
//
// Solidity: function setConsensusFactory(address factory) returns()
func (_FileTribeDApp *FileTribeDAppSession) SetConsensusFactory(factory common.Address) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.SetConsensusFactory(&_FileTribeDApp.TransactOpts, factory)
}

// SetConsensusFactory is a paid mutator transaction binding the contract method 0x1c38ea56.
//
// Solidity: function setConsensusFactory(address factory) returns()
func (_FileTribeDApp *FileTribeDAppTransactorSession) SetConsensusFactory(factory common.Address) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.SetConsensusFactory(&_FileTribeDApp.TransactOpts, factory)
}

// SetGroupFactory is a paid mutator transaction binding the contract method 0x837b3b93.
//
// Solidity: function setGroupFactory(address factory) returns()
func (_FileTribeDApp *FileTribeDAppTransactor) SetGroupFactory(opts *bind.TransactOpts, factory common.Address) (*types.Transaction, error) {
	return _FileTribeDApp.contract.Transact(opts, "setGroupFactory", factory)
}

// SetGroupFactory is a paid mutator transaction binding the contract method 0x837b3b93.
//
// Solidity: function setGroupFactory(address factory) returns()
func (_FileTribeDApp *FileTribeDAppSession) SetGroupFactory(factory common.Address) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.SetGroupFactory(&_FileTribeDApp.TransactOpts, factory)
}

// SetGroupFactory is a paid mutator transaction binding the contract method 0x837b3b93.
//
// Solidity: function setGroupFactory(address factory) returns()
func (_FileTribeDApp *FileTribeDAppTransactorSession) SetGroupFactory(factory common.Address) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.SetGroupFactory(&_FileTribeDApp.TransactOpts, factory)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FileTribeDApp *FileTribeDAppTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _FileTribeDApp.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FileTribeDApp *FileTribeDAppSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.TransferOwnership(&_FileTribeDApp.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FileTribeDApp *FileTribeDAppTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.TransferOwnership(&_FileTribeDApp.TransactOpts, newOwner)
}

// FileTribeDAppAccountCreatedIterator is returned from FilterAccountCreated and is used to iterate over the raw logs and unpacked data for AccountCreated events raised by the FileTribeDApp contract.
type FileTribeDAppAccountCreatedIterator struct {
	Event *FileTribeDAppAccountCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FileTribeDAppAccountCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FileTribeDAppAccountCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FileTribeDAppAccountCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FileTribeDAppAccountCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FileTribeDAppAccountCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FileTribeDAppAccountCreated represents a AccountCreated event raised by the FileTribeDApp contract.
type FileTribeDAppAccountCreated struct {
	Owner   common.Address
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAccountCreated is a free log retrieval operation binding the contract event 0xac631f3001b55ea1509cf3d7e74898f85392a61a76e8149181ae1259622dabc8.
//
// Solidity: event AccountCreated(address owner, address account)
func (_FileTribeDApp *FileTribeDAppFilterer) FilterAccountCreated(opts *bind.FilterOpts) (*FileTribeDAppAccountCreatedIterator, error) {

	logs, sub, err := _FileTribeDApp.contract.FilterLogs(opts, "AccountCreated")
	if err != nil {
		return nil, err
	}
	return &FileTribeDAppAccountCreatedIterator{contract: _FileTribeDApp.contract, event: "AccountCreated", logs: logs, sub: sub}, nil
}

// WatchAccountCreated is a free log subscription operation binding the contract event 0xac631f3001b55ea1509cf3d7e74898f85392a61a76e8149181ae1259622dabc8.
//
// Solidity: event AccountCreated(address owner, address account)
func (_FileTribeDApp *FileTribeDAppFilterer) WatchAccountCreated(opts *bind.WatchOpts, sink chan<- *FileTribeDAppAccountCreated) (event.Subscription, error) {

	logs, sub, err := _FileTribeDApp.contract.WatchLogs(opts, "AccountCreated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FileTribeDAppAccountCreated)
				if err := _FileTribeDApp.contract.UnpackLog(event, "AccountCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAccountCreated is a log parse operation binding the contract event 0xac631f3001b55ea1509cf3d7e74898f85392a61a76e8149181ae1259622dabc8.
//
// Solidity: event AccountCreated(address owner, address account)
func (_FileTribeDApp *FileTribeDAppFilterer) ParseAccountCreated(log types.Log) (*FileTribeDAppAccountCreated, error) {
	event := new(FileTribeDAppAccountCreated)
	if err := _FileTribeDApp.contract.UnpackLog(event, "AccountCreated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// FileTribeDAppDebugIterator is returned from FilterDebug and is used to iterate over the raw logs and unpacked data for Debug events raised by the FileTribeDApp contract.
type FileTribeDAppDebugIterator struct {
	Event *FileTribeDAppDebug // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FileTribeDAppDebugIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FileTribeDAppDebug)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FileTribeDAppDebug)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FileTribeDAppDebugIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FileTribeDAppDebugIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FileTribeDAppDebug represents a Debug event raised by the FileTribeDApp contract.
type FileTribeDAppDebug struct {
	Msg *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterDebug is a free log retrieval operation binding the contract event 0xaad112e1b1cb380614d4d2c495b213d9ed63560382320d6a49cc64793149f93c.
//
// Solidity: event Debug(int256 msg)
func (_FileTribeDApp *FileTribeDAppFilterer) FilterDebug(opts *bind.FilterOpts) (*FileTribeDAppDebugIterator, error) {

	logs, sub, err := _FileTribeDApp.contract.FilterLogs(opts, "Debug")
	if err != nil {
		return nil, err
	}
	return &FileTribeDAppDebugIterator{contract: _FileTribeDApp.contract, event: "Debug", logs: logs, sub: sub}, nil
}

// WatchDebug is a free log subscription operation binding the contract event 0xaad112e1b1cb380614d4d2c495b213d9ed63560382320d6a49cc64793149f93c.
//
// Solidity: event Debug(int256 msg)
func (_FileTribeDApp *FileTribeDAppFilterer) WatchDebug(opts *bind.WatchOpts, sink chan<- *FileTribeDAppDebug) (event.Subscription, error) {

	logs, sub, err := _FileTribeDApp.contract.WatchLogs(opts, "Debug")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FileTribeDAppDebug)
				if err := _FileTribeDApp.contract.UnpackLog(event, "Debug", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDebug is a log parse operation binding the contract event 0xaad112e1b1cb380614d4d2c495b213d9ed63560382320d6a49cc64793149f93c.
//
// Solidity: event Debug(int256 msg)
func (_FileTribeDApp *FileTribeDAppFilterer) ParseDebug(log types.Log) (*FileTribeDAppDebug, error) {
	event := new(FileTribeDAppDebug)
	if err := _FileTribeDApp.contract.UnpackLog(event, "Debug", log); err != nil {
		return nil, err
	}
	return event, nil
}

// FileTribeDAppDebugBytesIterator is returned from FilterDebugBytes and is used to iterate over the raw logs and unpacked data for DebugBytes events raised by the FileTribeDApp contract.
type FileTribeDAppDebugBytesIterator struct {
	Event *FileTribeDAppDebugBytes // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FileTribeDAppDebugBytesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FileTribeDAppDebugBytes)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FileTribeDAppDebugBytes)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FileTribeDAppDebugBytesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FileTribeDAppDebugBytesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FileTribeDAppDebugBytes represents a DebugBytes event raised by the FileTribeDApp contract.
type FileTribeDAppDebugBytes struct {
	Msg []byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterDebugBytes is a free log retrieval operation binding the contract event 0xaf5a5af90a78ece430d7df503b54fc4070844db69884a9a4afb00710a4816e53.
//
// Solidity: event DebugBytes(bytes msg)
func (_FileTribeDApp *FileTribeDAppFilterer) FilterDebugBytes(opts *bind.FilterOpts) (*FileTribeDAppDebugBytesIterator, error) {

	logs, sub, err := _FileTribeDApp.contract.FilterLogs(opts, "DebugBytes")
	if err != nil {
		return nil, err
	}
	return &FileTribeDAppDebugBytesIterator{contract: _FileTribeDApp.contract, event: "DebugBytes", logs: logs, sub: sub}, nil
}

// WatchDebugBytes is a free log subscription operation binding the contract event 0xaf5a5af90a78ece430d7df503b54fc4070844db69884a9a4afb00710a4816e53.
//
// Solidity: event DebugBytes(bytes msg)
func (_FileTribeDApp *FileTribeDAppFilterer) WatchDebugBytes(opts *bind.WatchOpts, sink chan<- *FileTribeDAppDebugBytes) (event.Subscription, error) {

	logs, sub, err := _FileTribeDApp.contract.WatchLogs(opts, "DebugBytes")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FileTribeDAppDebugBytes)
				if err := _FileTribeDApp.contract.UnpackLog(event, "DebugBytes", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDebugBytes is a log parse operation binding the contract event 0xaf5a5af90a78ece430d7df503b54fc4070844db69884a9a4afb00710a4816e53.
//
// Solidity: event DebugBytes(bytes msg)
func (_FileTribeDApp *FileTribeDAppFilterer) ParseDebugBytes(log types.Log) (*FileTribeDAppDebugBytes, error) {
	event := new(FileTribeDAppDebugBytes)
	if err := _FileTribeDApp.contract.UnpackLog(event, "DebugBytes", log); err != nil {
		return nil, err
	}
	return event, nil
}

// FileTribeDAppGroupRegisteredIterator is returned from FilterGroupRegistered and is used to iterate over the raw logs and unpacked data for GroupRegistered events raised by the FileTribeDApp contract.
type FileTribeDAppGroupRegisteredIterator struct {
	Event *FileTribeDAppGroupRegistered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FileTribeDAppGroupRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FileTribeDAppGroupRegistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FileTribeDAppGroupRegistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FileTribeDAppGroupRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FileTribeDAppGroupRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FileTribeDAppGroupRegistered represents a GroupRegistered event raised by the FileTribeDApp contract.
type FileTribeDAppGroupRegistered struct {
	Id  [32]byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterGroupRegistered is a free log retrieval operation binding the contract event 0xb78b34f3219f25d6305837697a0e5d110975d6be50317c9a2e815823306c3576.
//
// Solidity: event GroupRegistered(bytes32 id)
func (_FileTribeDApp *FileTribeDAppFilterer) FilterGroupRegistered(opts *bind.FilterOpts) (*FileTribeDAppGroupRegisteredIterator, error) {

	logs, sub, err := _FileTribeDApp.contract.FilterLogs(opts, "GroupRegistered")
	if err != nil {
		return nil, err
	}
	return &FileTribeDAppGroupRegisteredIterator{contract: _FileTribeDApp.contract, event: "GroupRegistered", logs: logs, sub: sub}, nil
}

// WatchGroupRegistered is a free log subscription operation binding the contract event 0xb78b34f3219f25d6305837697a0e5d110975d6be50317c9a2e815823306c3576.
//
// Solidity: event GroupRegistered(bytes32 id)
func (_FileTribeDApp *FileTribeDAppFilterer) WatchGroupRegistered(opts *bind.WatchOpts, sink chan<- *FileTribeDAppGroupRegistered) (event.Subscription, error) {

	logs, sub, err := _FileTribeDApp.contract.WatchLogs(opts, "GroupRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FileTribeDAppGroupRegistered)
				if err := _FileTribeDApp.contract.UnpackLog(event, "GroupRegistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseGroupRegistered is a log parse operation binding the contract event 0xb78b34f3219f25d6305837697a0e5d110975d6be50317c9a2e815823306c3576.
//
// Solidity: event GroupRegistered(bytes32 id)
func (_FileTribeDApp *FileTribeDAppFilterer) ParseGroupRegistered(log types.Log) (*FileTribeDAppGroupRegistered, error) {
	event := new(FileTribeDAppGroupRegistered)
	if err := _FileTribeDApp.contract.UnpackLog(event, "GroupRegistered", log); err != nil {
		return nil, err
	}
	return event, nil
}

// FileTribeDAppOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the FileTribeDApp contract.
type FileTribeDAppOwnershipTransferredIterator struct {
	Event *FileTribeDAppOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FileTribeDAppOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FileTribeDAppOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FileTribeDAppOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FileTribeDAppOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FileTribeDAppOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FileTribeDAppOwnershipTransferred represents a OwnershipTransferred event raised by the FileTribeDApp contract.
type FileTribeDAppOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_FileTribeDApp *FileTribeDAppFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*FileTribeDAppOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _FileTribeDApp.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &FileTribeDAppOwnershipTransferredIterator{contract: _FileTribeDApp.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_FileTribeDApp *FileTribeDAppFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FileTribeDAppOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _FileTribeDApp.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FileTribeDAppOwnershipTransferred)
				if err := _FileTribeDApp.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_FileTribeDApp *FileTribeDAppFilterer) ParseOwnershipTransferred(log types.Log) (*FileTribeDAppOwnershipTransferred, error) {
	event := new(FileTribeDAppOwnershipTransferred)
	if err := _FileTribeDApp.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}
