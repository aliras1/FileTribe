// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package IFileTribeDApp

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

// IFileTribeDAppABI is the input ABI used to generate the binding from.
const IFileTribeDAppABI = "[{\"constant\":false,\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"createGroup\",\"outputs\":[{\"internalType\":\"contractIGroup\",\"name\":\"group\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"contractIAccount\",\"name\":\"proposer\",\"type\":\"address\"}],\"name\":\"createConsensus\",\"outputs\":[{\"internalType\":\"contractIConsensus\",\"name\":\"consensus\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"getAccountOf\",\"outputs\":[{\"internalType\":\"contractIAccount\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"createDkg\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// IFileTribeDApp is an auto generated Go binding around an Ethereum contract.
type IFileTribeDApp struct {
	IFileTribeDAppCaller     // Read-only binding to the contract
	IFileTribeDAppTransactor // Write-only binding to the contract
	IFileTribeDAppFilterer   // Log filterer for contract events
}

// IFileTribeDAppCaller is an auto generated read-only Go binding around an Ethereum contract.
type IFileTribeDAppCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IFileTribeDAppTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IFileTribeDAppTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IFileTribeDAppFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IFileTribeDAppFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IFileTribeDAppSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IFileTribeDAppSession struct {
	Contract     *IFileTribeDApp   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IFileTribeDAppCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IFileTribeDAppCallerSession struct {
	Contract *IFileTribeDAppCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// IFileTribeDAppTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IFileTribeDAppTransactorSession struct {
	Contract     *IFileTribeDAppTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// IFileTribeDAppRaw is an auto generated low-level Go binding around an Ethereum contract.
type IFileTribeDAppRaw struct {
	Contract *IFileTribeDApp // Generic contract binding to access the raw methods on
}

// IFileTribeDAppCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IFileTribeDAppCallerRaw struct {
	Contract *IFileTribeDAppCaller // Generic read-only contract binding to access the raw methods on
}

// IFileTribeDAppTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IFileTribeDAppTransactorRaw struct {
	Contract *IFileTribeDAppTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIFileTribeDApp creates a new instance of IFileTribeDApp, bound to a specific deployed contract.
func NewIFileTribeDApp(address common.Address, backend bind.ContractBackend) (*IFileTribeDApp, error) {
	contract, err := bindIFileTribeDApp(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IFileTribeDApp{IFileTribeDAppCaller: IFileTribeDAppCaller{contract: contract}, IFileTribeDAppTransactor: IFileTribeDAppTransactor{contract: contract}, IFileTribeDAppFilterer: IFileTribeDAppFilterer{contract: contract}}, nil
}

// NewIFileTribeDAppCaller creates a new read-only instance of IFileTribeDApp, bound to a specific deployed contract.
func NewIFileTribeDAppCaller(address common.Address, caller bind.ContractCaller) (*IFileTribeDAppCaller, error) {
	contract, err := bindIFileTribeDApp(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IFileTribeDAppCaller{contract: contract}, nil
}

// NewIFileTribeDAppTransactor creates a new write-only instance of IFileTribeDApp, bound to a specific deployed contract.
func NewIFileTribeDAppTransactor(address common.Address, transactor bind.ContractTransactor) (*IFileTribeDAppTransactor, error) {
	contract, err := bindIFileTribeDApp(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IFileTribeDAppTransactor{contract: contract}, nil
}

// NewIFileTribeDAppFilterer creates a new log filterer instance of IFileTribeDApp, bound to a specific deployed contract.
func NewIFileTribeDAppFilterer(address common.Address, filterer bind.ContractFilterer) (*IFileTribeDAppFilterer, error) {
	contract, err := bindIFileTribeDApp(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IFileTribeDAppFilterer{contract: contract}, nil
}

// bindIFileTribeDApp binds a generic wrapper to an already deployed contract.
func bindIFileTribeDApp(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IFileTribeDAppABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IFileTribeDApp *IFileTribeDAppRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _IFileTribeDApp.Contract.IFileTribeDAppCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IFileTribeDApp *IFileTribeDAppRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IFileTribeDApp.Contract.IFileTribeDAppTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IFileTribeDApp *IFileTribeDAppRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IFileTribeDApp.Contract.IFileTribeDAppTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IFileTribeDApp *IFileTribeDAppCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _IFileTribeDApp.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IFileTribeDApp *IFileTribeDAppTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IFileTribeDApp.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IFileTribeDApp *IFileTribeDAppTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IFileTribeDApp.Contract.contract.Transact(opts, method, params...)
}

// GetAccountOf is a free data retrieval call binding the contract method 0x95184d3b.
//
// Solidity: function getAccountOf(address owner) constant returns(address)
func (_IFileTribeDApp *IFileTribeDAppCaller) GetAccountOf(opts *bind.CallOpts, owner common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _IFileTribeDApp.contract.Call(opts, out, "getAccountOf", owner)
	return *ret0, err
}

// GetAccountOf is a free data retrieval call binding the contract method 0x95184d3b.
//
// Solidity: function getAccountOf(address owner) constant returns(address)
func (_IFileTribeDApp *IFileTribeDAppSession) GetAccountOf(owner common.Address) (common.Address, error) {
	return _IFileTribeDApp.Contract.GetAccountOf(&_IFileTribeDApp.CallOpts, owner)
}

// GetAccountOf is a free data retrieval call binding the contract method 0x95184d3b.
//
// Solidity: function getAccountOf(address owner) constant returns(address)
func (_IFileTribeDApp *IFileTribeDAppCallerSession) GetAccountOf(owner common.Address) (common.Address, error) {
	return _IFileTribeDApp.Contract.GetAccountOf(&_IFileTribeDApp.CallOpts, owner)
}

// CreateConsensus is a paid mutator transaction binding the contract method 0x15db2069.
//
// Solidity: function createConsensus(address proposer) returns(address consensus)
func (_IFileTribeDApp *IFileTribeDAppTransactor) CreateConsensus(opts *bind.TransactOpts, proposer common.Address) (*types.Transaction, error) {
	return _IFileTribeDApp.contract.Transact(opts, "createConsensus", proposer)
}

// CreateConsensus is a paid mutator transaction binding the contract method 0x15db2069.
//
// Solidity: function createConsensus(address proposer) returns(address consensus)
func (_IFileTribeDApp *IFileTribeDAppSession) CreateConsensus(proposer common.Address) (*types.Transaction, error) {
	return _IFileTribeDApp.Contract.CreateConsensus(&_IFileTribeDApp.TransactOpts, proposer)
}

// CreateConsensus is a paid mutator transaction binding the contract method 0x15db2069.
//
// Solidity: function createConsensus(address proposer) returns(address consensus)
func (_IFileTribeDApp *IFileTribeDAppTransactorSession) CreateConsensus(proposer common.Address) (*types.Transaction, error) {
	return _IFileTribeDApp.Contract.CreateConsensus(&_IFileTribeDApp.TransactOpts, proposer)
}

// CreateDkg is a paid mutator transaction binding the contract method 0x45aa926f.
//
// Solidity: function createDkg() returns(address)
func (_IFileTribeDApp *IFileTribeDAppTransactor) CreateDkg(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IFileTribeDApp.contract.Transact(opts, "createDkg")
}

// CreateDkg is a paid mutator transaction binding the contract method 0x45aa926f.
//
// Solidity: function createDkg() returns(address)
func (_IFileTribeDApp *IFileTribeDAppSession) CreateDkg() (*types.Transaction, error) {
	return _IFileTribeDApp.Contract.CreateDkg(&_IFileTribeDApp.TransactOpts)
}

// CreateDkg is a paid mutator transaction binding the contract method 0x45aa926f.
//
// Solidity: function createDkg() returns(address)
func (_IFileTribeDApp *IFileTribeDAppTransactorSession) CreateDkg() (*types.Transaction, error) {
	return _IFileTribeDApp.Contract.CreateDkg(&_IFileTribeDApp.TransactOpts)
}

// CreateGroup is a paid mutator transaction binding the contract method 0xdc2ddcae.
//
// Solidity: function createGroup(string name) returns(address group)
func (_IFileTribeDApp *IFileTribeDAppTransactor) CreateGroup(opts *bind.TransactOpts, name string) (*types.Transaction, error) {
	return _IFileTribeDApp.contract.Transact(opts, "createGroup", name)
}

// CreateGroup is a paid mutator transaction binding the contract method 0xdc2ddcae.
//
// Solidity: function createGroup(string name) returns(address group)
func (_IFileTribeDApp *IFileTribeDAppSession) CreateGroup(name string) (*types.Transaction, error) {
	return _IFileTribeDApp.Contract.CreateGroup(&_IFileTribeDApp.TransactOpts, name)
}

// CreateGroup is a paid mutator transaction binding the contract method 0xdc2ddcae.
//
// Solidity: function createGroup(string name) returns(address group)
func (_IFileTribeDApp *IFileTribeDAppTransactorSession) CreateGroup(name string) (*types.Transaction, error) {
	return _IFileTribeDApp.Contract.CreateGroup(&_IFileTribeDApp.TransactOpts, name)
}
