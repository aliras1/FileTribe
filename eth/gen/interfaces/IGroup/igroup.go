// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package IGroup

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

// IGroupABI is the input ABI used to generate the binding from.
const IGroupABI = "[{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"isMember\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"onChangeIpfsHashConsensus\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"threshold\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// IGroup is an auto generated Go binding around an Ethereum contract.
type IGroup struct {
	IGroupCaller     // Read-only binding to the contract
	IGroupTransactor // Write-only binding to the contract
	IGroupFilterer   // Log filterer for contract events
}

// IGroupCaller is an auto generated read-only Go binding around an Ethereum contract.
type IGroupCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IGroupTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IGroupTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IGroupFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IGroupFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IGroupSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IGroupSession struct {
	Contract     *IGroup           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IGroupCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IGroupCallerSession struct {
	Contract *IGroupCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// IGroupTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IGroupTransactorSession struct {
	Contract     *IGroupTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IGroupRaw is an auto generated low-level Go binding around an Ethereum contract.
type IGroupRaw struct {
	Contract *IGroup // Generic contract binding to access the raw methods on
}

// IGroupCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IGroupCallerRaw struct {
	Contract *IGroupCaller // Generic read-only contract binding to access the raw methods on
}

// IGroupTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IGroupTransactorRaw struct {
	Contract *IGroupTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIGroup creates a new instance of IGroup, bound to a specific deployed contract.
func NewIGroup(address common.Address, backend bind.ContractBackend) (*IGroup, error) {
	contract, err := bindIGroup(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IGroup{IGroupCaller: IGroupCaller{contract: contract}, IGroupTransactor: IGroupTransactor{contract: contract}, IGroupFilterer: IGroupFilterer{contract: contract}}, nil
}

// NewIGroupCaller creates a new read-only instance of IGroup, bound to a specific deployed contract.
func NewIGroupCaller(address common.Address, caller bind.ContractCaller) (*IGroupCaller, error) {
	contract, err := bindIGroup(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IGroupCaller{contract: contract}, nil
}

// NewIGroupTransactor creates a new write-only instance of IGroup, bound to a specific deployed contract.
func NewIGroupTransactor(address common.Address, transactor bind.ContractTransactor) (*IGroupTransactor, error) {
	contract, err := bindIGroup(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IGroupTransactor{contract: contract}, nil
}

// NewIGroupFilterer creates a new log filterer instance of IGroup, bound to a specific deployed contract.
func NewIGroupFilterer(address common.Address, filterer bind.ContractFilterer) (*IGroupFilterer, error) {
	contract, err := bindIGroup(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IGroupFilterer{contract: contract}, nil
}

// bindIGroup binds a generic wrapper to an already deployed contract.
func bindIGroup(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IGroupABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IGroup *IGroupRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _IGroup.Contract.IGroupCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IGroup *IGroupRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IGroup.Contract.IGroupTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IGroup *IGroupRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IGroup.Contract.IGroupTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IGroup *IGroupCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _IGroup.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IGroup *IGroupTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IGroup.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IGroup *IGroupTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IGroup.Contract.contract.Transact(opts, method, params...)
}

// IsMember is a free data retrieval call binding the contract method 0xa230c524.
//
// Solidity: function isMember(address owner) constant returns(bool)
func (_IGroup *IGroupCaller) IsMember(opts *bind.CallOpts, owner common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _IGroup.contract.Call(opts, out, "isMember", owner)
	return *ret0, err
}

// IsMember is a free data retrieval call binding the contract method 0xa230c524.
//
// Solidity: function isMember(address owner) constant returns(bool)
func (_IGroup *IGroupSession) IsMember(owner common.Address) (bool, error) {
	return _IGroup.Contract.IsMember(&_IGroup.CallOpts, owner)
}

// IsMember is a free data retrieval call binding the contract method 0xa230c524.
//
// Solidity: function isMember(address owner) constant returns(bool)
func (_IGroup *IGroupCallerSession) IsMember(owner common.Address) (bool, error) {
	return _IGroup.Contract.IsMember(&_IGroup.CallOpts, owner)
}

// Threshold is a free data retrieval call binding the contract method 0x42cde4e8.
//
// Solidity: function threshold() constant returns(uint256)
func (_IGroup *IGroupCaller) Threshold(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _IGroup.contract.Call(opts, out, "threshold")
	return *ret0, err
}

// Threshold is a free data retrieval call binding the contract method 0x42cde4e8.
//
// Solidity: function threshold() constant returns(uint256)
func (_IGroup *IGroupSession) Threshold() (*big.Int, error) {
	return _IGroup.Contract.Threshold(&_IGroup.CallOpts)
}

// Threshold is a free data retrieval call binding the contract method 0x42cde4e8.
//
// Solidity: function threshold() constant returns(uint256)
func (_IGroup *IGroupCallerSession) Threshold() (*big.Int, error) {
	return _IGroup.Contract.Threshold(&_IGroup.CallOpts)
}

// OnChangeIpfsHashConsensus is a paid mutator transaction binding the contract method 0x8379088e.
//
// Solidity: function onChangeIpfsHashConsensus(bytes payload) returns()
func (_IGroup *IGroupTransactor) OnChangeIpfsHashConsensus(opts *bind.TransactOpts, payload []byte) (*types.Transaction, error) {
	return _IGroup.contract.Transact(opts, "onChangeIpfsHashConsensus", payload)
}

// OnChangeIpfsHashConsensus is a paid mutator transaction binding the contract method 0x8379088e.
//
// Solidity: function onChangeIpfsHashConsensus(bytes payload) returns()
func (_IGroup *IGroupSession) OnChangeIpfsHashConsensus(payload []byte) (*types.Transaction, error) {
	return _IGroup.Contract.OnChangeIpfsHashConsensus(&_IGroup.TransactOpts, payload)
}

// OnChangeIpfsHashConsensus is a paid mutator transaction binding the contract method 0x8379088e.
//
// Solidity: function onChangeIpfsHashConsensus(bytes payload) returns()
func (_IGroup *IGroupTransactorSession) OnChangeIpfsHashConsensus(payload []byte) (*types.Transaction, error) {
	return _IGroup.Contract.OnChangeIpfsHashConsensus(&_IGroup.TransactOpts, payload)
}
