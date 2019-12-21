// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package IConsensusCallback

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

// IConsensusCallbackABI is the input ABI used to generate the binding from.
const IConsensusCallbackABI = "[{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"onConsensusSuccess\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"onConsensusFailure\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"threshold\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"isAuthorized\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// IConsensusCallback is an auto generated Go binding around an Ethereum contract.
type IConsensusCallback struct {
	IConsensusCallbackCaller     // Read-only binding to the contract
	IConsensusCallbackTransactor // Write-only binding to the contract
	IConsensusCallbackFilterer   // Log filterer for contract events
}

// IConsensusCallbackCaller is an auto generated read-only Go binding around an Ethereum contract.
type IConsensusCallbackCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IConsensusCallbackTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IConsensusCallbackTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IConsensusCallbackFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IConsensusCallbackFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IConsensusCallbackSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IConsensusCallbackSession struct {
	Contract     *IConsensusCallback // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// IConsensusCallbackCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IConsensusCallbackCallerSession struct {
	Contract *IConsensusCallbackCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// IConsensusCallbackTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IConsensusCallbackTransactorSession struct {
	Contract     *IConsensusCallbackTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// IConsensusCallbackRaw is an auto generated low-level Go binding around an Ethereum contract.
type IConsensusCallbackRaw struct {
	Contract *IConsensusCallback // Generic contract binding to access the raw methods on
}

// IConsensusCallbackCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IConsensusCallbackCallerRaw struct {
	Contract *IConsensusCallbackCaller // Generic read-only contract binding to access the raw methods on
}

// IConsensusCallbackTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IConsensusCallbackTransactorRaw struct {
	Contract *IConsensusCallbackTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIConsensusCallback creates a new instance of IConsensusCallback, bound to a specific deployed contract.
func NewIConsensusCallback(address common.Address, backend bind.ContractBackend) (*IConsensusCallback, error) {
	contract, err := bindIConsensusCallback(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IConsensusCallback{IConsensusCallbackCaller: IConsensusCallbackCaller{contract: contract}, IConsensusCallbackTransactor: IConsensusCallbackTransactor{contract: contract}, IConsensusCallbackFilterer: IConsensusCallbackFilterer{contract: contract}}, nil
}

// NewIConsensusCallbackCaller creates a new read-only instance of IConsensusCallback, bound to a specific deployed contract.
func NewIConsensusCallbackCaller(address common.Address, caller bind.ContractCaller) (*IConsensusCallbackCaller, error) {
	contract, err := bindIConsensusCallback(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IConsensusCallbackCaller{contract: contract}, nil
}

// NewIConsensusCallbackTransactor creates a new write-only instance of IConsensusCallback, bound to a specific deployed contract.
func NewIConsensusCallbackTransactor(address common.Address, transactor bind.ContractTransactor) (*IConsensusCallbackTransactor, error) {
	contract, err := bindIConsensusCallback(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IConsensusCallbackTransactor{contract: contract}, nil
}

// NewIConsensusCallbackFilterer creates a new log filterer instance of IConsensusCallback, bound to a specific deployed contract.
func NewIConsensusCallbackFilterer(address common.Address, filterer bind.ContractFilterer) (*IConsensusCallbackFilterer, error) {
	contract, err := bindIConsensusCallback(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IConsensusCallbackFilterer{contract: contract}, nil
}

// bindIConsensusCallback binds a generic wrapper to an already deployed contract.
func bindIConsensusCallback(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IConsensusCallbackABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IConsensusCallback *IConsensusCallbackRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _IConsensusCallback.Contract.IConsensusCallbackCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IConsensusCallback *IConsensusCallbackRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IConsensusCallback.Contract.IConsensusCallbackTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IConsensusCallback *IConsensusCallbackRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IConsensusCallback.Contract.IConsensusCallbackTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IConsensusCallback *IConsensusCallbackCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _IConsensusCallback.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IConsensusCallback *IConsensusCallbackTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IConsensusCallback.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IConsensusCallback *IConsensusCallbackTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IConsensusCallback.Contract.contract.Transact(opts, method, params...)
}

// IsAuthorized is a free data retrieval call binding the contract method 0xfe9fbb80.
//
// Solidity: function isAuthorized(address sender) constant returns(bool)
func (_IConsensusCallback *IConsensusCallbackCaller) IsAuthorized(opts *bind.CallOpts, sender common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _IConsensusCallback.contract.Call(opts, out, "isAuthorized", sender)
	return *ret0, err
}

// IsAuthorized is a free data retrieval call binding the contract method 0xfe9fbb80.
//
// Solidity: function isAuthorized(address sender) constant returns(bool)
func (_IConsensusCallback *IConsensusCallbackSession) IsAuthorized(sender common.Address) (bool, error) {
	return _IConsensusCallback.Contract.IsAuthorized(&_IConsensusCallback.CallOpts, sender)
}

// IsAuthorized is a free data retrieval call binding the contract method 0xfe9fbb80.
//
// Solidity: function isAuthorized(address sender) constant returns(bool)
func (_IConsensusCallback *IConsensusCallbackCallerSession) IsAuthorized(sender common.Address) (bool, error) {
	return _IConsensusCallback.Contract.IsAuthorized(&_IConsensusCallback.CallOpts, sender)
}

// Threshold is a free data retrieval call binding the contract method 0x42cde4e8.
//
// Solidity: function threshold() constant returns(uint256)
func (_IConsensusCallback *IConsensusCallbackCaller) Threshold(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _IConsensusCallback.contract.Call(opts, out, "threshold")
	return *ret0, err
}

// Threshold is a free data retrieval call binding the contract method 0x42cde4e8.
//
// Solidity: function threshold() constant returns(uint256)
func (_IConsensusCallback *IConsensusCallbackSession) Threshold() (*big.Int, error) {
	return _IConsensusCallback.Contract.Threshold(&_IConsensusCallback.CallOpts)
}

// Threshold is a free data retrieval call binding the contract method 0x42cde4e8.
//
// Solidity: function threshold() constant returns(uint256)
func (_IConsensusCallback *IConsensusCallbackCallerSession) Threshold() (*big.Int, error) {
	return _IConsensusCallback.Contract.Threshold(&_IConsensusCallback.CallOpts)
}

// OnConsensusFailure is a paid mutator transaction binding the contract method 0x6f04df41.
//
// Solidity: function onConsensusFailure(bytes payload) returns()
func (_IConsensusCallback *IConsensusCallbackTransactor) OnConsensusFailure(opts *bind.TransactOpts, payload []byte) (*types.Transaction, error) {
	return _IConsensusCallback.contract.Transact(opts, "onConsensusFailure", payload)
}

// OnConsensusFailure is a paid mutator transaction binding the contract method 0x6f04df41.
//
// Solidity: function onConsensusFailure(bytes payload) returns()
func (_IConsensusCallback *IConsensusCallbackSession) OnConsensusFailure(payload []byte) (*types.Transaction, error) {
	return _IConsensusCallback.Contract.OnConsensusFailure(&_IConsensusCallback.TransactOpts, payload)
}

// OnConsensusFailure is a paid mutator transaction binding the contract method 0x6f04df41.
//
// Solidity: function onConsensusFailure(bytes payload) returns()
func (_IConsensusCallback *IConsensusCallbackTransactorSession) OnConsensusFailure(payload []byte) (*types.Transaction, error) {
	return _IConsensusCallback.Contract.OnConsensusFailure(&_IConsensusCallback.TransactOpts, payload)
}

// OnConsensusSuccess is a paid mutator transaction binding the contract method 0x7e2e5ddf.
//
// Solidity: function onConsensusSuccess(bytes payload) returns()
func (_IConsensusCallback *IConsensusCallbackTransactor) OnConsensusSuccess(opts *bind.TransactOpts, payload []byte) (*types.Transaction, error) {
	return _IConsensusCallback.contract.Transact(opts, "onConsensusSuccess", payload)
}

// OnConsensusSuccess is a paid mutator transaction binding the contract method 0x7e2e5ddf.
//
// Solidity: function onConsensusSuccess(bytes payload) returns()
func (_IConsensusCallback *IConsensusCallbackSession) OnConsensusSuccess(payload []byte) (*types.Transaction, error) {
	return _IConsensusCallback.Contract.OnConsensusSuccess(&_IConsensusCallback.TransactOpts, payload)
}

// OnConsensusSuccess is a paid mutator transaction binding the contract method 0x7e2e5ddf.
//
// Solidity: function onConsensusSuccess(bytes payload) returns()
func (_IConsensusCallback *IConsensusCallbackTransactorSession) OnConsensusSuccess(payload []byte) (*types.Transaction, error) {
	return _IConsensusCallback.Contract.OnConsensusSuccess(&_IConsensusCallback.TransactOpts, payload)
}
