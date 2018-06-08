// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package eth

import (
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// EthABI is the input ABI used to generate the binding from.
const EthABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"UserRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"message\",\"type\":\"bytes1[]\"}],\"name\":\"MessageSent\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\"},{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"boxingKey\",\"type\":\"bytes32\"},{\"name\":\"verifyKey\",\"type\":\"bytes32\"},{\"name\":\"ipfsAddr\",\"type\":\"string\"}],\"name\":\"registerUser\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"isUserRegistered\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"getUser\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"string\"},{\"name\":\"\",\"type\":\"bytes32\"},{\"name\":\"\",\"type\":\"bytes32\"},{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"message\",\"type\":\"bytes1[]\"}],\"name\":\"sendMessage\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// EthBin is the compiled bytecode used for deploying new contracts.
const EthBin = `0x608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550610a7a806100606000396000f30060806040526004361061006d576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063311607db1461007257806351b6d61f146100bb5780636517579c146101945780638da5cb5b146102fb578063ad9eb55714610352575b600080fd5b34801561007e57600080fd5b506100a160048036038101908080356000191690602001909291905050506103b8565b604051808215151515815260200191505060405180910390f35b3480156100c757600080fd5b506101926004803603810190808035600019169060200190929190803590602001908201803590602001908080601f016020809104026020016040519081016040528093929190818152602001838380828437820191505050505050919291929080356000191690602001909291908035600019169060200190929190803590602001908201803590602001908080601f01602080910402602001604051908101604052809392919081815260200183838082843782019150505050505091929192905050506103fe565b005b3480156101a057600080fd5b506101c36004803603810190808035600019169060200190929190505050610637565b604051808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001806020018560001916600019168152602001846000191660001916815260200180602001838103835287818151815260200191508051906020019080838360005b8381101561025557808201518184015260208101905061023a565b50505050905090810190601f1680156102825780820380516001836020036101000a031916815260200191505b50838103825284818151815260200191508051906020019080838360005b838110156102bb5780820151818401526020810190506102a0565b50505050905090810190601f1680156102e85780820380516001836020036101000a031916815260200191505b5097505050505050505060405180910390f35b34801561030757600080fd5b506103106108b4565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561035e57600080fd5b506103b6600480360381019080803590602001908201803590602001908080602002602001604051908101604052809392919081815260200183836020028082843782019150505050505091929192905050506108d9565b005b600060016000836000191660001916815260200190815260200160002060050160009054906101000a900460ff16156103f457600190506103f9565b600090505b919050565b60016000836000191660001916815260200190815260200160002060050160009054906101000a900460ff1615151561049f576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260178152602001807f557365726e616d6520616c72656164792065786973747300000000000000000081525060200191505060405180910390fd5b60c0604051908101604052803373ffffffffffffffffffffffffffffffffffffffff16815260200185815260200184600019168152602001836000191681526020018281526020016001151581525060016000876000191660001916815260200190815260200160002060008201518160000160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550602082015181600101908051906020019061056c929190610954565b50604082015181600201906000191690556060820151816003019060001916905560808201518160040190805190602001906105a9929190610954565b5060a08201518160050160006101000a81548160ff0219169083151502179055509050507f54db7a5cb4735e1aac1f53db512d3390390bb6637bd30ad4bf9fc98667d9b9b933604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390a15050505050565b6000606060008060606106486109d4565b60016000886000191660001916815260200190815260200160002060050160009054906101000a900460ff16151561067f57600080fd5b60016000886000191660001916815260200190815260200160002060c060405190810160405290816000820160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001600182018054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156107955780601f1061076a57610100808354040283529160200191610795565b820191906000526020600020905b81548152906001019060200180831161077857829003601f168201915b50505050508152602001600282015460001916600019168152602001600382015460001916600019168152602001600482018054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561085b5780601f106108305761010080835404028352916020019161085b565b820191906000526020600020905b81548152906001019060200180831161083e57829003601f168201915b505050505081526020016005820160009054906101000a900460ff161515151581525050905080600001518160200151826040015183606001518460800151839350809050955095509550955095505091939590929450565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b7ff8bf0ac14a97aa75283adb748197d00e04cddea87df0f19cdb29aae4ca69bf83816040518080602001828103825283818151815260200191508051906020019060200280838360005b8381101561093e578082015181840152602081019050610923565b505050509050019250505060405180910390a150565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061099557805160ff19168380011785556109c3565b828001600101855582156109c3579182015b828111156109c25782518255916020019190600101906109a7565b5b5090506109d09190610a29565b5090565b60c060405190810160405280600073ffffffffffffffffffffffffffffffffffffffff168152602001606081526020016000801916815260200160008019168152602001606081526020016000151581525090565b610a4b91905b80821115610a47576000816000905550600101610a2f565b5090565b905600a165627a7a72305820acb020bf1bd6045e844fa481ef79958ace44dc76abf32a9ce9aed74d839ff0910029`

// DeployEth deploys a new Ethereum contract, binding an instance of Eth to it.
func DeployEth(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Eth, error) {
	parsed, err := abi.JSON(strings.NewReader(EthABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(EthBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Eth{EthCaller: EthCaller{contract: contract}, EthTransactor: EthTransactor{contract: contract}, EthFilterer: EthFilterer{contract: contract}}, nil
}

// Eth is an auto generated Go binding around an Ethereum contract.
type Eth struct {
	EthCaller     // Read-only binding to the contract
	EthTransactor // Write-only binding to the contract
	EthFilterer   // Log filterer for contract events
}

// EthCaller is an auto generated read-only Go binding around an Ethereum contract.
type EthCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EthTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EthFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EthSession struct {
	Contract     *Eth              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EthCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EthCallerSession struct {
	Contract *EthCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// EthTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EthTransactorSession struct {
	Contract     *EthTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EthRaw is an auto generated low-level Go binding around an Ethereum contract.
type EthRaw struct {
	Contract *Eth // Generic contract binding to access the raw methods on
}

// EthCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EthCallerRaw struct {
	Contract *EthCaller // Generic read-only contract binding to access the raw methods on
}

// EthTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EthTransactorRaw struct {
	Contract *EthTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEth creates a new instance of Eth, bound to a specific deployed contract.
func NewEth(address common.Address, backend bind.ContractBackend) (*Eth, error) {
	contract, err := bindEth(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Eth{EthCaller: EthCaller{contract: contract}, EthTransactor: EthTransactor{contract: contract}, EthFilterer: EthFilterer{contract: contract}}, nil
}

// NewEthCaller creates a new read-only instance of Eth, bound to a specific deployed contract.
func NewEthCaller(address common.Address, caller bind.ContractCaller) (*EthCaller, error) {
	contract, err := bindEth(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EthCaller{contract: contract}, nil
}

// NewEthTransactor creates a new write-only instance of Eth, bound to a specific deployed contract.
func NewEthTransactor(address common.Address, transactor bind.ContractTransactor) (*EthTransactor, error) {
	contract, err := bindEth(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EthTransactor{contract: contract}, nil
}

// NewEthFilterer creates a new log filterer instance of Eth, bound to a specific deployed contract.
func NewEthFilterer(address common.Address, filterer bind.ContractFilterer) (*EthFilterer, error) {
	contract, err := bindEth(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EthFilterer{contract: contract}, nil
}

// bindEth binds a generic wrapper to an already deployed contract.
func bindEth(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(EthABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Eth *EthRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Eth.Contract.EthCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Eth *EthRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Eth.Contract.EthTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Eth *EthRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Eth.Contract.EthTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Eth *EthCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Eth.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Eth *EthTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Eth.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Eth *EthTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Eth.Contract.contract.Transact(opts, method, params...)
}

// GetUser is a free data retrieval call binding the contract method 0x6517579c.
//
// Solidity: function getUser(id bytes32) constant returns(address, string, bytes32, bytes32, string)
func (_Eth *EthCaller) GetUser(opts *bind.CallOpts, id [32]byte) (common.Address, string, [32]byte, [32]byte, string, error) {
	var (
		ret0 = new(common.Address)
		ret1 = new(string)
		ret2 = new([32]byte)
		ret3 = new([32]byte)
		ret4 = new(string)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
		ret3,
		ret4,
	}
	err := _Eth.contract.Call(opts, out, "getUser", id)
	return *ret0, *ret1, *ret2, *ret3, *ret4, err
}

// GetUser is a free data retrieval call binding the contract method 0x6517579c.
//
// Solidity: function getUser(id bytes32) constant returns(address, string, bytes32, bytes32, string)
func (_Eth *EthSession) GetUser(id [32]byte) (common.Address, string, [32]byte, [32]byte, string, error) {
	return _Eth.Contract.GetUser(&_Eth.CallOpts, id)
}

// GetUser is a free data retrieval call binding the contract method 0x6517579c.
//
// Solidity: function getUser(id bytes32) constant returns(address, string, bytes32, bytes32, string)
func (_Eth *EthCallerSession) GetUser(id [32]byte) (common.Address, string, [32]byte, [32]byte, string, error) {
	return _Eth.Contract.GetUser(&_Eth.CallOpts, id)
}

// IsUserRegistered is a free data retrieval call binding the contract method 0x311607db.
//
// Solidity: function isUserRegistered(id bytes32) constant returns(bool)
func (_Eth *EthCaller) IsUserRegistered(opts *bind.CallOpts, id [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Eth.contract.Call(opts, out, "isUserRegistered", id)
	return *ret0, err
}

// IsUserRegistered is a free data retrieval call binding the contract method 0x311607db.
//
// Solidity: function isUserRegistered(id bytes32) constant returns(bool)
func (_Eth *EthSession) IsUserRegistered(id [32]byte) (bool, error) {
	return _Eth.Contract.IsUserRegistered(&_Eth.CallOpts, id)
}

// IsUserRegistered is a free data retrieval call binding the contract method 0x311607db.
//
// Solidity: function isUserRegistered(id bytes32) constant returns(bool)
func (_Eth *EthCallerSession) IsUserRegistered(id [32]byte) (bool, error) {
	return _Eth.Contract.IsUserRegistered(&_Eth.CallOpts, id)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Eth *EthCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Eth.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Eth *EthSession) Owner() (common.Address, error) {
	return _Eth.Contract.Owner(&_Eth.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Eth *EthCallerSession) Owner() (common.Address, error) {
	return _Eth.Contract.Owner(&_Eth.CallOpts)
}

// RegisterUser is a paid mutator transaction binding the contract method 0x51b6d61f.
//
// Solidity: function registerUser(id bytes32, name string, boxingKey bytes32, verifyKey bytes32, ipfsAddr string) returns()
func (_Eth *EthTransactor) RegisterUser(opts *bind.TransactOpts, id [32]byte, name string, boxingKey [32]byte, verifyKey [32]byte, ipfsAddr string) (*types.Transaction, error) {
	return _Eth.contract.Transact(opts, "registerUser", id, name, boxingKey, verifyKey, ipfsAddr)
}

// RegisterUser is a paid mutator transaction binding the contract method 0x51b6d61f.
//
// Solidity: function registerUser(id bytes32, name string, boxingKey bytes32, verifyKey bytes32, ipfsAddr string) returns()
func (_Eth *EthSession) RegisterUser(id [32]byte, name string, boxingKey [32]byte, verifyKey [32]byte, ipfsAddr string) (*types.Transaction, error) {
	return _Eth.Contract.RegisterUser(&_Eth.TransactOpts, id, name, boxingKey, verifyKey, ipfsAddr)
}

// RegisterUser is a paid mutator transaction binding the contract method 0x51b6d61f.
//
// Solidity: function registerUser(id bytes32, name string, boxingKey bytes32, verifyKey bytes32, ipfsAddr string) returns()
func (_Eth *EthTransactorSession) RegisterUser(id [32]byte, name string, boxingKey [32]byte, verifyKey [32]byte, ipfsAddr string) (*types.Transaction, error) {
	return _Eth.Contract.RegisterUser(&_Eth.TransactOpts, id, name, boxingKey, verifyKey, ipfsAddr)
}

// SendMessage is a paid mutator transaction binding the contract method 0xad9eb557.
//
// Solidity: function sendMessage(message bytes1[]) returns()
func (_Eth *EthTransactor) SendMessage(opts *bind.TransactOpts, message [][1]byte) (*types.Transaction, error) {
	return _Eth.contract.Transact(opts, "sendMessage", message)
}

// SendMessage is a paid mutator transaction binding the contract method 0xad9eb557.
//
// Solidity: function sendMessage(message bytes1[]) returns()
func (_Eth *EthSession) SendMessage(message [][1]byte) (*types.Transaction, error) {
	return _Eth.Contract.SendMessage(&_Eth.TransactOpts, message)
}

// SendMessage is a paid mutator transaction binding the contract method 0xad9eb557.
//
// Solidity: function sendMessage(message bytes1[]) returns()
func (_Eth *EthTransactorSession) SendMessage(message [][1]byte) (*types.Transaction, error) {
	return _Eth.Contract.SendMessage(&_Eth.TransactOpts, message)
}

// EthMessageSentIterator is returned from FilterMessageSent and is used to iterate over the raw logs and unpacked data for MessageSent events raised by the Eth contract.
type EthMessageSentIterator struct {
	Event *EthMessageSent // Event containing the contract specifics and raw log

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
func (it *EthMessageSentIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthMessageSent)
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
		it.Event = new(EthMessageSent)
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
func (it *EthMessageSentIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthMessageSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthMessageSent represents a MessageSent event raised by the Eth contract.
type EthMessageSent struct {
	Message [][1]byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMessageSent is a free log retrieval operation binding the contract event 0xf8bf0ac14a97aa75283adb748197d00e04cddea87df0f19cdb29aae4ca69bf83.
//
// Solidity: e MessageSent(message bytes1[])
func (_Eth *EthFilterer) FilterMessageSent(opts *bind.FilterOpts) (*EthMessageSentIterator, error) {

	logs, sub, err := _Eth.contract.FilterLogs(opts, "MessageSent")
	if err != nil {
		return nil, err
	}
	return &EthMessageSentIterator{contract: _Eth.contract, event: "MessageSent", logs: logs, sub: sub}, nil
}

// WatchMessageSent is a free log subscription operation binding the contract event 0xf8bf0ac14a97aa75283adb748197d00e04cddea87df0f19cdb29aae4ca69bf83.
//
// Solidity: e MessageSent(message bytes1[])
func (_Eth *EthFilterer) WatchMessageSent(opts *bind.WatchOpts, sink chan<- *EthMessageSent) (event.Subscription, error) {

	logs, sub, err := _Eth.contract.WatchLogs(opts, "MessageSent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthMessageSent)
				if err := _Eth.contract.UnpackLog(event, "MessageSent", log); err != nil {
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

// EthUserRegisteredIterator is returned from FilterUserRegistered and is used to iterate over the raw logs and unpacked data for UserRegistered events raised by the Eth contract.
type EthUserRegisteredIterator struct {
	Event *EthUserRegistered // Event containing the contract specifics and raw log

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
func (it *EthUserRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthUserRegistered)
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
		it.Event = new(EthUserRegistered)
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
func (it *EthUserRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthUserRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthUserRegistered represents a UserRegistered event raised by the Eth contract.
type EthUserRegistered struct {
	Addr common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterUserRegistered is a free log retrieval operation binding the contract event 0x54db7a5cb4735e1aac1f53db512d3390390bb6637bd30ad4bf9fc98667d9b9b9.
//
// Solidity: e UserRegistered(addr address)
func (_Eth *EthFilterer) FilterUserRegistered(opts *bind.FilterOpts) (*EthUserRegisteredIterator, error) {

	logs, sub, err := _Eth.contract.FilterLogs(opts, "UserRegistered")
	if err != nil {
		return nil, err
	}
	return &EthUserRegisteredIterator{contract: _Eth.contract, event: "UserRegistered", logs: logs, sub: sub}, nil
}

// WatchUserRegistered is a free log subscription operation binding the contract event 0x54db7a5cb4735e1aac1f53db512d3390390bb6637bd30ad4bf9fc98667d9b9b9.
//
// Solidity: e UserRegistered(addr address)
func (_Eth *EthFilterer) WatchUserRegistered(opts *bind.WatchOpts, sink chan<- *EthUserRegistered) (event.Subscription, error) {

	logs, sub, err := _Eth.contract.WatchLogs(opts, "UserRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthUserRegistered)
				if err := _Eth.contract.UnpackLog(event, "UserRegistered", log); err != nil {
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
