// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ConsensusFactory

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

// ConsensusFactoryABI is the input ABI used to generate the binding from.
const ConsensusFactoryABI = "[{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"contractIAccount\",\"name\":\"proposer\",\"type\":\"address\"},{\"internalType\":\"contractIGroup\",\"name\":\"group\",\"type\":\"address\"}],\"name\":\"create\",\"outputs\":[{\"internalType\":\"contractIConsensus\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"parent\",\"type\":\"address\"}],\"name\":\"setParent\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// ConsensusFactoryBin is the compiled bytecode used for deploying new contracts.
var ConsensusFactoryBin = "0x608060405234801561001057600080fd5b50600080546001600160a01b0319163390811780835560405191926001600160a01b0391909116917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908290a350610a548061006d6000396000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c80631499c592146100675780633e68680a1461008f578063715018a6146100d95780638da5cb5b146100e15780638f32d59b146100e9578063f2fde38b14610105575b600080fd5b61008d6004803603602081101561007d57600080fd5b50356001600160a01b031661012b565b005b6100bd600480360360408110156100a557600080fd5b506001600160a01b038135811691602001351661015e565b604080516001600160a01b039092168252519081900360200190f35b61008d6101a9565b6100bd610204565b6100f1610213565b604080519115158252519081900360200190f35b61008d6004803603602081101561011b57600080fd5b50356001600160a01b0316610224565b610133610213565b61013c57600080fd5b600180546001600160a01b0319166001600160a01b0392909216919091179055565b6000828260405161016e906102af565b6001600160a01b03928316815291166020820152604080519182900301906000f0801580156101a1573d6000803e3d6000fd5b509392505050565b6101b1610213565b6101ba57600080fd5b600080546040516001600160a01b03909116907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a3600080546001600160a01b0319169055565b6000546001600160a01b031690565b6000546001600160a01b0316331490565b61022c610213565b61023557600080fd5b61023e81610241565b50565b6001600160a01b03811661025457600080fd5b600080546040516001600160a01b03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a3600080546001600160a01b0319166001600160a01b0392909216919091179055565b610763806102bd8339019056fe608060405234801561001057600080fd5b506040516107633803806107638339818101604052604081101561003357600080fd5b508051602090910151600080546001600160a01b039384166001600160a01b031991821617909155600180549390921692169190911790556106e98061007a6000396000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c806312424e3f146100675780635d1ca63114610071578063a878f8581461008b578063a8e4fb9014610108578063bade60331461012c578063e9790d0214610108575b600080fd5b61006f61019c565b005b61007961042a565b60408051918252519081900360200190f35b610093610431565b6040805160208082528351818301528351919283929083019185019080838360005b838110156100cd5781810151838201526020016100b5565b50505050905090810190601f1680156100fa5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6101106104c4565b604080516001600160a01b039092168252519081900360200190f35b61006f6004803603604081101561014257600080fd5b81019060208101813564010000000081111561015d57600080fd5b82018360208201111561016f57600080fd5b8035906020019184600183028401116401000000008311171561019157600080fd5b9193509150356104d3565b6001546040805163288c314960e21b815233600482015290516001600160a01b039092169163a230c52491602480820192602092909190829003018186803b1580156101e757600080fd5b505afa1580156101fb573d6000803e3d6000fd5b505050506040513d602081101561021157600080fd5b5051610264576040805162461bcd60e51b815260206004820152601b60248201527f75736572206973206e6f74206d656d626572206f662067726f75700000000000604482015290519081900360640190fd5b600554604080516020808201939093523360601b8183015281516034818303018152605490910182528051908301206000818152600390935291205460ff16156102ae5750610428565b600081815260036020908152604091829020805460ff19166001908117909155548251630859bc9d60e31b815292516001600160a01b03909116926342cde4e8926004808301939192829003018186803b15801561030b57600080fd5b505afa15801561031f573d6000803e3d6000fd5b505050506040513d602081101561033557600080fd5b50516004805460010190819055111561042657600180546040516341bc844760e11b815260206004820190815260028054600019958116156101000295909501909416849004602483018190526001600160a01b0390931693638379088e9390928291604490910190849080156103ed5780601f106103c2576101008083540402835291602001916103ed565b820191906000526020600020905b8154815290600101906020018083116103d057829003601f168201915b505092505050600060405180830381600087803b15801561040d57600080fd5b505af1158015610421573d6000803e3d6000fd5b505050505b505b565b6005545b90565b60028054604080516020601f60001961010060018716150201909416859004938401819004810282018101909252828152606093909290918301828280156104ba5780601f1061048f576101008083540402835291602001916104ba565b820191906000526020600020905b81548152906001019060200180831161049d57829003601f168201915b5050505050905090565b6000546001600160a01b031690565b6001546001600160a01b03163314610532576040805162461bcd60e51b815260206004820152601d60248201527f6d73672e73656e646572206973206e6f742067726f7570206f776e6572000000604482015290519081900360640190fd5b61053e6002848461061c565b506005819055600160049081556000805460408051638da5cb5b60e01b81529051929385936001600160a01b0390931692638da5cb5b928083019260209291829003018186803b15801561059157600080fd5b505afa1580156105a5573d6000803e3d6000fd5b505050506040513d60208110156105bb57600080fd5b50516040805160208082019490945260609290921b6bffffffffffffffffffffffff19168282015280518083036034018152605490920181528151918301919091206000908152600390925290208054600160ff1990911617905550505050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061065d5782800160ff1982351617855561068a565b8280016001018555821561068a579182015b8281111561068a57823582559160200191906001019061066f565b5061069692915061069a565b5090565b61042e91905b8082111561069657600081556001016106a056fea265627a7a72315820fdd7c244e09d1e94f93027bb465fb232c76a5d4075162a3b24e698adcf5d415864736f6c634300050c0032a265627a7a723158204157d8db4ce399afe534e9076116b0504c97dab374b204ae8d3491d872b07d0464736f6c634300050c0032"

// DeployConsensusFactory deploys a new Ethereum contract, binding an instance of ConsensusFactory to it.
func DeployConsensusFactory(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ConsensusFactory, error) {
	parsed, err := abi.JSON(strings.NewReader(ConsensusFactoryABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ConsensusFactoryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ConsensusFactory{ConsensusFactoryCaller: ConsensusFactoryCaller{contract: contract}, ConsensusFactoryTransactor: ConsensusFactoryTransactor{contract: contract}, ConsensusFactoryFilterer: ConsensusFactoryFilterer{contract: contract}}, nil
}

// ConsensusFactory is an auto generated Go binding around an Ethereum contract.
type ConsensusFactory struct {
	ConsensusFactoryCaller     // Read-only binding to the contract
	ConsensusFactoryTransactor // Write-only binding to the contract
	ConsensusFactoryFilterer   // Log filterer for contract events
}

// ConsensusFactoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type ConsensusFactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConsensusFactoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ConsensusFactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConsensusFactoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ConsensusFactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConsensusFactorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ConsensusFactorySession struct {
	Contract     *ConsensusFactory // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ConsensusFactoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ConsensusFactoryCallerSession struct {
	Contract *ConsensusFactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// ConsensusFactoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ConsensusFactoryTransactorSession struct {
	Contract     *ConsensusFactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// ConsensusFactoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type ConsensusFactoryRaw struct {
	Contract *ConsensusFactory // Generic contract binding to access the raw methods on
}

// ConsensusFactoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ConsensusFactoryCallerRaw struct {
	Contract *ConsensusFactoryCaller // Generic read-only contract binding to access the raw methods on
}

// ConsensusFactoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ConsensusFactoryTransactorRaw struct {
	Contract *ConsensusFactoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewConsensusFactory creates a new instance of ConsensusFactory, bound to a specific deployed contract.
func NewConsensusFactory(address common.Address, backend bind.ContractBackend) (*ConsensusFactory, error) {
	contract, err := bindConsensusFactory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ConsensusFactory{ConsensusFactoryCaller: ConsensusFactoryCaller{contract: contract}, ConsensusFactoryTransactor: ConsensusFactoryTransactor{contract: contract}, ConsensusFactoryFilterer: ConsensusFactoryFilterer{contract: contract}}, nil
}

// NewConsensusFactoryCaller creates a new read-only instance of ConsensusFactory, bound to a specific deployed contract.
func NewConsensusFactoryCaller(address common.Address, caller bind.ContractCaller) (*ConsensusFactoryCaller, error) {
	contract, err := bindConsensusFactory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ConsensusFactoryCaller{contract: contract}, nil
}

// NewConsensusFactoryTransactor creates a new write-only instance of ConsensusFactory, bound to a specific deployed contract.
func NewConsensusFactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*ConsensusFactoryTransactor, error) {
	contract, err := bindConsensusFactory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ConsensusFactoryTransactor{contract: contract}, nil
}

// NewConsensusFactoryFilterer creates a new log filterer instance of ConsensusFactory, bound to a specific deployed contract.
func NewConsensusFactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*ConsensusFactoryFilterer, error) {
	contract, err := bindConsensusFactory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ConsensusFactoryFilterer{contract: contract}, nil
}

// bindConsensusFactory binds a generic wrapper to an already deployed contract.
func bindConsensusFactory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ConsensusFactoryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConsensusFactory *ConsensusFactoryRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ConsensusFactory.Contract.ConsensusFactoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConsensusFactory *ConsensusFactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConsensusFactory.Contract.ConsensusFactoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConsensusFactory *ConsensusFactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConsensusFactory.Contract.ConsensusFactoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConsensusFactory *ConsensusFactoryCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ConsensusFactory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConsensusFactory *ConsensusFactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConsensusFactory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConsensusFactory *ConsensusFactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConsensusFactory.Contract.contract.Transact(opts, method, params...)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_ConsensusFactory *ConsensusFactoryCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ConsensusFactory.contract.Call(opts, out, "isOwner")
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_ConsensusFactory *ConsensusFactorySession) IsOwner() (bool, error) {
	return _ConsensusFactory.Contract.IsOwner(&_ConsensusFactory.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_ConsensusFactory *ConsensusFactoryCallerSession) IsOwner() (bool, error) {
	return _ConsensusFactory.Contract.IsOwner(&_ConsensusFactory.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_ConsensusFactory *ConsensusFactoryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ConsensusFactory.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_ConsensusFactory *ConsensusFactorySession) Owner() (common.Address, error) {
	return _ConsensusFactory.Contract.Owner(&_ConsensusFactory.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_ConsensusFactory *ConsensusFactoryCallerSession) Owner() (common.Address, error) {
	return _ConsensusFactory.Contract.Owner(&_ConsensusFactory.CallOpts)
}

// Create is a paid mutator transaction binding the contract method 0x3e68680a.
//
// Solidity: function create(address proposer, address group) returns(address)
func (_ConsensusFactory *ConsensusFactoryTransactor) Create(opts *bind.TransactOpts, proposer common.Address, group common.Address) (*types.Transaction, error) {
	return _ConsensusFactory.contract.Transact(opts, "create", proposer, group)
}

// Create is a paid mutator transaction binding the contract method 0x3e68680a.
//
// Solidity: function create(address proposer, address group) returns(address)
func (_ConsensusFactory *ConsensusFactorySession) Create(proposer common.Address, group common.Address) (*types.Transaction, error) {
	return _ConsensusFactory.Contract.Create(&_ConsensusFactory.TransactOpts, proposer, group)
}

// Create is a paid mutator transaction binding the contract method 0x3e68680a.
//
// Solidity: function create(address proposer, address group) returns(address)
func (_ConsensusFactory *ConsensusFactoryTransactorSession) Create(proposer common.Address, group common.Address) (*types.Transaction, error) {
	return _ConsensusFactory.Contract.Create(&_ConsensusFactory.TransactOpts, proposer, group)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ConsensusFactory *ConsensusFactoryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConsensusFactory.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ConsensusFactory *ConsensusFactorySession) RenounceOwnership() (*types.Transaction, error) {
	return _ConsensusFactory.Contract.RenounceOwnership(&_ConsensusFactory.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ConsensusFactory *ConsensusFactoryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ConsensusFactory.Contract.RenounceOwnership(&_ConsensusFactory.TransactOpts)
}

// SetParent is a paid mutator transaction binding the contract method 0x1499c592.
//
// Solidity: function setParent(address parent) returns()
func (_ConsensusFactory *ConsensusFactoryTransactor) SetParent(opts *bind.TransactOpts, parent common.Address) (*types.Transaction, error) {
	return _ConsensusFactory.contract.Transact(opts, "setParent", parent)
}

// SetParent is a paid mutator transaction binding the contract method 0x1499c592.
//
// Solidity: function setParent(address parent) returns()
func (_ConsensusFactory *ConsensusFactorySession) SetParent(parent common.Address) (*types.Transaction, error) {
	return _ConsensusFactory.Contract.SetParent(&_ConsensusFactory.TransactOpts, parent)
}

// SetParent is a paid mutator transaction binding the contract method 0x1499c592.
//
// Solidity: function setParent(address parent) returns()
func (_ConsensusFactory *ConsensusFactoryTransactorSession) SetParent(parent common.Address) (*types.Transaction, error) {
	return _ConsensusFactory.Contract.SetParent(&_ConsensusFactory.TransactOpts, parent)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ConsensusFactory *ConsensusFactoryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ConsensusFactory.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ConsensusFactory *ConsensusFactorySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ConsensusFactory.Contract.TransferOwnership(&_ConsensusFactory.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ConsensusFactory *ConsensusFactoryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ConsensusFactory.Contract.TransferOwnership(&_ConsensusFactory.TransactOpts, newOwner)
}

// ConsensusFactoryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ConsensusFactory contract.
type ConsensusFactoryOwnershipTransferredIterator struct {
	Event *ConsensusFactoryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ConsensusFactoryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConsensusFactoryOwnershipTransferred)
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
		it.Event = new(ConsensusFactoryOwnershipTransferred)
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
func (it *ConsensusFactoryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConsensusFactoryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConsensusFactoryOwnershipTransferred represents a OwnershipTransferred event raised by the ConsensusFactory contract.
type ConsensusFactoryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ConsensusFactory *ConsensusFactoryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ConsensusFactoryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ConsensusFactory.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ConsensusFactoryOwnershipTransferredIterator{contract: _ConsensusFactory.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ConsensusFactory *ConsensusFactoryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ConsensusFactoryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ConsensusFactory.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConsensusFactoryOwnershipTransferred)
				if err := _ConsensusFactory.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ConsensusFactory *ConsensusFactoryFilterer) ParseOwnershipTransferred(log types.Log) (*ConsensusFactoryOwnershipTransferred, error) {
	event := new(ConsensusFactoryOwnershipTransferred)
	if err := _ConsensusFactory.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}
