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
const EthABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"UserRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"message\",\"type\":\"bytes1[]\"}],\"name\":\"MessageSent\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"boxingKey\",\"type\":\"bytes32\"},{\"name\":\"verifyKey\",\"type\":\"bytes1[]\"},{\"name\":\"ipfsAddr\",\"type\":\"string\"}],\"name\":\"registerUser\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"address\"}],\"name\":\"isUserRegistered\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"address\"}],\"name\":\"getUser\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"},{\"name\":\"\",\"type\":\"bytes32\"},{\"name\":\"\",\"type\":\"bytes1[]\"},{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"message\",\"type\":\"bytes1[]\"}],\"name\":\"sendMessage\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// EthBin is the compiled bytecode used for deploying new contracts.
const EthBin = `0x608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550610c47806100606000396000f30060806040526004361061006d576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063163f7522146100725780636f77926b146100cd5780638178c0261461024c5780638da5cb5b1461034c578063ad9eb557146103a3575b600080fd5b34801561007e57600080fd5b506100b3600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610409565b604051808215151515815260200191505060405180910390f35b3480156100d957600080fd5b5061010e600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610473565b604051808060200185600019166000191681526020018060200180602001848103845288818151815260200191508051906020019080838360005b83811015610164578082015181840152602081019050610149565b50505050905090810190601f1680156101915780820380516001836020036101000a031916815260200191505b50848103835286818151815260200191508051906020019060200280838360005b838110156101cd5780820151818401526020810190506101b2565b50505050905001848103825285818151815260200191508051906020019080838360005b8381101561020c5780820151818401526020810190506101f1565b50505050905090810190601f1680156102395780820380516001836020036101000a031916815260200191505b5097505050505050505060405180910390f35b34801561025857600080fd5b5061034a600480360381019080803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290803560001916906020019092919080359060200190820180359060200190808060200260200160405190810160405280939291908181526020018383602002808284378201915050505050509192919290803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290505050610784565b005b34801561035857600080fd5b506103616109aa565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156103af57600080fd5b50610407600480360381019080803590602001908201803590602001908080602002602001604051908101604052809392919081815260200183836020028082843782019150505050505091929192905050506109cf565b005b6000600160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060040160009054906101000a900460ff1615610469576001905061046e565b600090505b919050565b60606000606080610482610a4a565b600160008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060040160009054906101000a900460ff1615156104dd57600080fd5b600160008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060a06040519081016040529081600082018054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156105c15780601f10610596576101008083540402835291602001916105c1565b820191906000526020600020905b8154815290600101906020018083116105a457829003601f168201915b505050505081526020016001820154600019166000191681526020016002820180548060200260200160405190810160405280929190818152602001828054801561068f57602002820191906000526020600020906000905b82829054906101000a90047f0100000000000000000000000000000000000000000000000000000000000000027effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168152602001906001019060208260000104928301926001038202915080841161061a5790505b50505050508152602001600382018054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156107315780601f1061070657610100808354040283529160200191610731565b820191906000526020600020905b81548152906001019060200180831161071457829003601f168201915b505050505081526020016004820160009054906101000a900460ff161515151581525050905080600001518160200151826040015183606001518393508191508090509450945094509450509193509193565b600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060040160009054906101000a900460ff16151515610849576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260178152602001807f557365726e616d6520616c72656164792065786973747300000000000000000081525060200191505060405180910390fd5b60a0604051908101604052808581526020018460001916815260200183815260200182815260200160011515815250600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008201518160000190805190602001906108d3929190610a7f565b50602082015181600101906000191690556040820151816002019080519060200190610900929190610aff565b50606082015181600301908051906020019061091d929190610a7f565b5060808201518160040160006101000a81548160ff0219169083151502179055509050507f54db7a5cb4735e1aac1f53db512d3390390bb6637bd30ad4bf9fc98667d9b9b933604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390a150505050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b7ff8bf0ac14a97aa75283adb748197d00e04cddea87df0f19cdb29aae4ca69bf83816040518080602001828103825283818151815260200191508051906020019060200280838360005b83811015610a34578082015181840152602081019050610a19565b505050509050019250505060405180910390a150565b60a060405190810160405280606081526020016000801916815260200160608152602001606081526020016000151581525090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10610ac057805160ff1916838001178555610aee565b82800160010185558215610aee579182015b82811115610aed578251825591602001919060010190610ad2565b5b509050610afb9190610bc6565b5090565b82805482825590600052602060002090601f01602090048101928215610bb55791602002820160005b83821115610b8657835183826101000a81548160ff02191690837f0100000000000000000000000000000000000000000000000000000000000000900402179055509260200192600101602081600001049283019260010302610b28565b8015610bb35782816101000a81549060ff0219169055600101602081600001049283019260010302610b86565b505b509050610bc29190610beb565b5090565b610be891905b80821115610be4576000816000905550600101610bcc565b5090565b90565b610c1891905b80821115610c1457600081816101000a81549060ff021916905550600101610bf1565b5090565b905600a165627a7a72305820a84f36d3251172953ab2093a58ed24f89e8e605147d1168b709c31a278699e5a0029`

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

// GetUser is a free data retrieval call binding the contract method 0x6f77926b.
//
// Solidity: function getUser(id address) constant returns(string, bytes32, bytes1[], string)
func (_Eth *EthCaller) GetUser(opts *bind.CallOpts, id common.Address) (string, [32]byte, [][1]byte, string, error) {
	var (
		ret0 = new(string)
		ret1 = new([32]byte)
		ret2 = new([][1]byte)
		ret3 = new(string)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
		ret3,
	}
	err := _Eth.contract.Call(opts, out, "getUser", id)
	return *ret0, *ret1, *ret2, *ret3, err
}

// GetUser is a free data retrieval call binding the contract method 0x6f77926b.
//
// Solidity: function getUser(id address) constant returns(string, bytes32, bytes1[], string)
func (_Eth *EthSession) GetUser(id common.Address) (string, [32]byte, [][1]byte, string, error) {
	return _Eth.Contract.GetUser(&_Eth.CallOpts, id)
}

// GetUser is a free data retrieval call binding the contract method 0x6f77926b.
//
// Solidity: function getUser(id address) constant returns(string, bytes32, bytes1[], string)
func (_Eth *EthCallerSession) GetUser(id common.Address) (string, [32]byte, [][1]byte, string, error) {
	return _Eth.Contract.GetUser(&_Eth.CallOpts, id)
}

// IsUserRegistered is a free data retrieval call binding the contract method 0x163f7522.
//
// Solidity: function isUserRegistered(id address) constant returns(bool)
func (_Eth *EthCaller) IsUserRegistered(opts *bind.CallOpts, id common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Eth.contract.Call(opts, out, "isUserRegistered", id)
	return *ret0, err
}

// IsUserRegistered is a free data retrieval call binding the contract method 0x163f7522.
//
// Solidity: function isUserRegistered(id address) constant returns(bool)
func (_Eth *EthSession) IsUserRegistered(id common.Address) (bool, error) {
	return _Eth.Contract.IsUserRegistered(&_Eth.CallOpts, id)
}

// IsUserRegistered is a free data retrieval call binding the contract method 0x163f7522.
//
// Solidity: function isUserRegistered(id address) constant returns(bool)
func (_Eth *EthCallerSession) IsUserRegistered(id common.Address) (bool, error) {
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

// RegisterUser is a paid mutator transaction binding the contract method 0x8178c026.
//
// Solidity: function registerUser(name string, boxingKey bytes32, verifyKey bytes1[], ipfsAddr string) returns()
func (_Eth *EthTransactor) RegisterUser(opts *bind.TransactOpts, name string, boxingKey [32]byte, verifyKey [][1]byte, ipfsAddr string) (*types.Transaction, error) {
	return _Eth.contract.Transact(opts, "registerUser", name, boxingKey, verifyKey, ipfsAddr)
}

// RegisterUser is a paid mutator transaction binding the contract method 0x8178c026.
//
// Solidity: function registerUser(name string, boxingKey bytes32, verifyKey bytes1[], ipfsAddr string) returns()
func (_Eth *EthSession) RegisterUser(name string, boxingKey [32]byte, verifyKey [][1]byte, ipfsAddr string) (*types.Transaction, error) {
	return _Eth.Contract.RegisterUser(&_Eth.TransactOpts, name, boxingKey, verifyKey, ipfsAddr)
}

// RegisterUser is a paid mutator transaction binding the contract method 0x8178c026.
//
// Solidity: function registerUser(name string, boxingKey bytes32, verifyKey bytes1[], ipfsAddr string) returns()
func (_Eth *EthTransactorSession) RegisterUser(name string, boxingKey [32]byte, verifyKey [][1]byte, ipfsAddr string) (*types.Transaction, error) {
	return _Eth.Contract.RegisterUser(&_Eth.TransactOpts, name, boxingKey, verifyKey, ipfsAddr)
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
