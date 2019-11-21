// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package IConsensus

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

// IConsensusABI is the input ABI used to generate the binding from.
const IConsensusABI = "[{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"counter\",\"type\":\"uint256\"}],\"name\":\"propose\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getProposer\",\"outputs\":[{\"internalType\":\"contractIAccount\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// IConsensus is an auto generated Go binding around an Ethereum contract.
type IConsensus struct {
	IConsensusCaller     // Read-only binding to the contract
	IConsensusTransactor // Write-only binding to the contract
	IConsensusFilterer   // Log filterer for contract events
}

// IConsensusCaller is an auto generated read-only Go binding around an Ethereum contract.
type IConsensusCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IConsensusTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IConsensusTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IConsensusFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IConsensusFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IConsensusSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IConsensusSession struct {
	Contract     *IConsensus       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IConsensusCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IConsensusCallerSession struct {
	Contract *IConsensusCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// IConsensusTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IConsensusTransactorSession struct {
	Contract     *IConsensusTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// IConsensusRaw is an auto generated low-level Go binding around an Ethereum contract.
type IConsensusRaw struct {
	Contract *IConsensus // Generic contract binding to access the raw methods on
}

// IConsensusCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IConsensusCallerRaw struct {
	Contract *IConsensusCaller // Generic read-only contract binding to access the raw methods on
}

// IConsensusTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IConsensusTransactorRaw struct {
	Contract *IConsensusTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIConsensus creates a new instance of IConsensus, bound to a specific deployed contract.
func NewIConsensus(address common.Address, backend bind.ContractBackend) (*IConsensus, error) {
	contract, err := bindIConsensus(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IConsensus{IConsensusCaller: IConsensusCaller{contract: contract}, IConsensusTransactor: IConsensusTransactor{contract: contract}, IConsensusFilterer: IConsensusFilterer{contract: contract}}, nil
}

// NewIConsensusCaller creates a new read-only instance of IConsensus, bound to a specific deployed contract.
func NewIConsensusCaller(address common.Address, caller bind.ContractCaller) (*IConsensusCaller, error) {
	contract, err := bindIConsensus(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IConsensusCaller{contract: contract}, nil
}

// NewIConsensusTransactor creates a new write-only instance of IConsensus, bound to a specific deployed contract.
func NewIConsensusTransactor(address common.Address, transactor bind.ContractTransactor) (*IConsensusTransactor, error) {
	contract, err := bindIConsensus(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IConsensusTransactor{contract: contract}, nil
}

// NewIConsensusFilterer creates a new log filterer instance of IConsensus, bound to a specific deployed contract.
func NewIConsensusFilterer(address common.Address, filterer bind.ContractFilterer) (*IConsensusFilterer, error) {
	contract, err := bindIConsensus(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IConsensusFilterer{contract: contract}, nil
}

// bindIConsensus binds a generic wrapper to an already deployed contract.
func bindIConsensus(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IConsensusABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IConsensus *IConsensusRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _IConsensus.Contract.IConsensusCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IConsensus *IConsensusRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IConsensus.Contract.IConsensusTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IConsensus *IConsensusRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IConsensus.Contract.IConsensusTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IConsensus *IConsensusCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _IConsensus.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IConsensus *IConsensusTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IConsensus.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IConsensus *IConsensusTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IConsensus.Contract.contract.Transact(opts, method, params...)
}

// GetId is a free data retrieval call binding the contract method 0x5d1ca631.
//
// Solidity: function getId() constant returns(uint256)
func (_IConsensus *IConsensusCaller) GetId(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _IConsensus.contract.Call(opts, out, "getId")
	return *ret0, err
}

// GetId is a free data retrieval call binding the contract method 0x5d1ca631.
//
// Solidity: function getId() constant returns(uint256)
func (_IConsensus *IConsensusSession) GetId() (*big.Int, error) {
	return _IConsensus.Contract.GetId(&_IConsensus.CallOpts)
}

// GetId is a free data retrieval call binding the contract method 0x5d1ca631.
//
// Solidity: function getId() constant returns(uint256)
func (_IConsensus *IConsensusCallerSession) GetId() (*big.Int, error) {
	return _IConsensus.Contract.GetId(&_IConsensus.CallOpts)
}

// GetProposer is a free data retrieval call binding the contract method 0xe9790d02.
//
// Solidity: function getProposer() constant returns(address)
func (_IConsensus *IConsensusCaller) GetProposer(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _IConsensus.contract.Call(opts, out, "getProposer")
	return *ret0, err
}

// GetProposer is a free data retrieval call binding the contract method 0xe9790d02.
//
// Solidity: function getProposer() constant returns(address)
func (_IConsensus *IConsensusSession) GetProposer() (common.Address, error) {
	return _IConsensus.Contract.GetProposer(&_IConsensus.CallOpts)
}

// GetProposer is a free data retrieval call binding the contract method 0xe9790d02.
//
// Solidity: function getProposer() constant returns(address)
func (_IConsensus *IConsensusCallerSession) GetProposer() (common.Address, error) {
	return _IConsensus.Contract.GetProposer(&_IConsensus.CallOpts)
}

// Propose is a paid mutator transaction binding the contract method 0xbade6033.
//
// Solidity: function propose(bytes payload, uint256 counter) returns()
func (_IConsensus *IConsensusTransactor) Propose(opts *bind.TransactOpts, payload []byte, counter *big.Int) (*types.Transaction, error) {
	return _IConsensus.contract.Transact(opts, "propose", payload, counter)
}

// Propose is a paid mutator transaction binding the contract method 0xbade6033.
//
// Solidity: function propose(bytes payload, uint256 counter) returns()
func (_IConsensus *IConsensusSession) Propose(payload []byte, counter *big.Int) (*types.Transaction, error) {
	return _IConsensus.Contract.Propose(&_IConsensus.TransactOpts, payload, counter)
}

// Propose is a paid mutator transaction binding the contract method 0xbade6033.
//
// Solidity: function propose(bytes payload, uint256 counter) returns()
func (_IConsensus *IConsensusTransactorSession) Propose(payload []byte, counter *big.Int) (*types.Transaction, error) {
	return _IConsensus.Contract.Propose(&_IConsensus.TransactOpts, payload, counter)
}
