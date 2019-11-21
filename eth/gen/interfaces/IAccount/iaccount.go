// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package IAccount

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

// IAccountABI is the input ABI used to generate the binding from.
const IAccountABI = "[{\"constant\":false,\"inputs\":[],\"name\":\"invite\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"onInvitationAccepted\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"onInvitationDeclined\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"contractIGroup\",\"name\":\"group\",\"type\":\"address\"}],\"name\":\"onGroupLeft\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// IAccount is an auto generated Go binding around an Ethereum contract.
type IAccount struct {
	IAccountCaller     // Read-only binding to the contract
	IAccountTransactor // Write-only binding to the contract
	IAccountFilterer   // Log filterer for contract events
}

// IAccountCaller is an auto generated read-only Go binding around an Ethereum contract.
type IAccountCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAccountTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IAccountTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAccountFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IAccountFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAccountSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IAccountSession struct {
	Contract     *IAccount         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IAccountCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IAccountCallerSession struct {
	Contract *IAccountCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// IAccountTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IAccountTransactorSession struct {
	Contract     *IAccountTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// IAccountRaw is an auto generated low-level Go binding around an Ethereum contract.
type IAccountRaw struct {
	Contract *IAccount // Generic contract binding to access the raw methods on
}

// IAccountCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IAccountCallerRaw struct {
	Contract *IAccountCaller // Generic read-only contract binding to access the raw methods on
}

// IAccountTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IAccountTransactorRaw struct {
	Contract *IAccountTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIAccount creates a new instance of IAccount, bound to a specific deployed contract.
func NewIAccount(address common.Address, backend bind.ContractBackend) (*IAccount, error) {
	contract, err := bindIAccount(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IAccount{IAccountCaller: IAccountCaller{contract: contract}, IAccountTransactor: IAccountTransactor{contract: contract}, IAccountFilterer: IAccountFilterer{contract: contract}}, nil
}

// NewIAccountCaller creates a new read-only instance of IAccount, bound to a specific deployed contract.
func NewIAccountCaller(address common.Address, caller bind.ContractCaller) (*IAccountCaller, error) {
	contract, err := bindIAccount(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IAccountCaller{contract: contract}, nil
}

// NewIAccountTransactor creates a new write-only instance of IAccount, bound to a specific deployed contract.
func NewIAccountTransactor(address common.Address, transactor bind.ContractTransactor) (*IAccountTransactor, error) {
	contract, err := bindIAccount(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IAccountTransactor{contract: contract}, nil
}

// NewIAccountFilterer creates a new log filterer instance of IAccount, bound to a specific deployed contract.
func NewIAccountFilterer(address common.Address, filterer bind.ContractFilterer) (*IAccountFilterer, error) {
	contract, err := bindIAccount(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IAccountFilterer{contract: contract}, nil
}

// bindIAccount binds a generic wrapper to an already deployed contract.
func bindIAccount(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IAccountABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IAccount *IAccountRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _IAccount.Contract.IAccountCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IAccount *IAccountRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAccount.Contract.IAccountTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IAccount *IAccountRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAccount.Contract.IAccountTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IAccount *IAccountCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _IAccount.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IAccount *IAccountTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAccount.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IAccount *IAccountTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAccount.Contract.contract.Transact(opts, method, params...)
}

// Invite is a paid mutator transaction binding the contract method 0xeec30bfd.
//
// Solidity: function invite() returns()
func (_IAccount *IAccountTransactor) Invite(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAccount.contract.Transact(opts, "invite")
}

// Invite is a paid mutator transaction binding the contract method 0xeec30bfd.
//
// Solidity: function invite() returns()
func (_IAccount *IAccountSession) Invite() (*types.Transaction, error) {
	return _IAccount.Contract.Invite(&_IAccount.TransactOpts)
}

// Invite is a paid mutator transaction binding the contract method 0xeec30bfd.
//
// Solidity: function invite() returns()
func (_IAccount *IAccountTransactorSession) Invite() (*types.Transaction, error) {
	return _IAccount.Contract.Invite(&_IAccount.TransactOpts)
}

// OnGroupLeft is a paid mutator transaction binding the contract method 0xaff10f28.
//
// Solidity: function onGroupLeft(address group) returns()
func (_IAccount *IAccountTransactor) OnGroupLeft(opts *bind.TransactOpts, group common.Address) (*types.Transaction, error) {
	return _IAccount.contract.Transact(opts, "onGroupLeft", group)
}

// OnGroupLeft is a paid mutator transaction binding the contract method 0xaff10f28.
//
// Solidity: function onGroupLeft(address group) returns()
func (_IAccount *IAccountSession) OnGroupLeft(group common.Address) (*types.Transaction, error) {
	return _IAccount.Contract.OnGroupLeft(&_IAccount.TransactOpts, group)
}

// OnGroupLeft is a paid mutator transaction binding the contract method 0xaff10f28.
//
// Solidity: function onGroupLeft(address group) returns()
func (_IAccount *IAccountTransactorSession) OnGroupLeft(group common.Address) (*types.Transaction, error) {
	return _IAccount.Contract.OnGroupLeft(&_IAccount.TransactOpts, group)
}

// OnInvitationAccepted is a paid mutator transaction binding the contract method 0x83e86440.
//
// Solidity: function onInvitationAccepted() returns()
func (_IAccount *IAccountTransactor) OnInvitationAccepted(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAccount.contract.Transact(opts, "onInvitationAccepted")
}

// OnInvitationAccepted is a paid mutator transaction binding the contract method 0x83e86440.
//
// Solidity: function onInvitationAccepted() returns()
func (_IAccount *IAccountSession) OnInvitationAccepted() (*types.Transaction, error) {
	return _IAccount.Contract.OnInvitationAccepted(&_IAccount.TransactOpts)
}

// OnInvitationAccepted is a paid mutator transaction binding the contract method 0x83e86440.
//
// Solidity: function onInvitationAccepted() returns()
func (_IAccount *IAccountTransactorSession) OnInvitationAccepted() (*types.Transaction, error) {
	return _IAccount.Contract.OnInvitationAccepted(&_IAccount.TransactOpts)
}

// OnInvitationDeclined is a paid mutator transaction binding the contract method 0xb061d9a9.
//
// Solidity: function onInvitationDeclined() returns()
func (_IAccount *IAccountTransactor) OnInvitationDeclined(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAccount.contract.Transact(opts, "onInvitationDeclined")
}

// OnInvitationDeclined is a paid mutator transaction binding the contract method 0xb061d9a9.
//
// Solidity: function onInvitationDeclined() returns()
func (_IAccount *IAccountSession) OnInvitationDeclined() (*types.Transaction, error) {
	return _IAccount.Contract.OnInvitationDeclined(&_IAccount.TransactOpts)
}

// OnInvitationDeclined is a paid mutator transaction binding the contract method 0xb061d9a9.
//
// Solidity: function onInvitationDeclined() returns()
func (_IAccount *IAccountTransactorSession) OnInvitationDeclined() (*types.Transaction, error) {
	return _IAccount.Contract.OnInvitationDeclined(&_IAccount.TransactOpts)
}
