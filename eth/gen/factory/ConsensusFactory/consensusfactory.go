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
const ConsensusFactoryABI = "[{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"contractIAccount\",\"name\":\"proposer\",\"type\":\"address\"},{\"internalType\":\"contractIConsensusCallback\",\"name\":\"callback\",\"type\":\"address\"}],\"name\":\"create\",\"outputs\":[{\"internalType\":\"contractIConsensus\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"parent\",\"type\":\"address\"}],\"name\":\"setParent\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// ConsensusFactoryBin is the compiled bytecode used for deploying new contracts.
var ConsensusFactoryBin = "0x608060405234801561001057600080fd5b50600080546001600160a01b0319163390811780835560405191926001600160a01b0391909116917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908290a350610c9d8061006d6000396000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c80631499c592146100675780633e68680a1461008f578063715018a6146100d95780638da5cb5b146100e15780638f32d59b146100e9578063f2fde38b14610105575b600080fd5b61008d6004803603602081101561007d57600080fd5b50356001600160a01b031661012b565b005b6100bd600480360360408110156100a557600080fd5b506001600160a01b038135811691602001351661015e565b604080516001600160a01b039092168252519081900360200190f35b61008d6101a9565b6100bd610204565b6100f1610213565b604080519115158252519081900360200190f35b61008d6004803603602081101561011b57600080fd5b50356001600160a01b0316610224565b610133610213565b61013c57600080fd5b600180546001600160a01b0319166001600160a01b0392909216919091179055565b6000828260405161016e906102af565b6001600160a01b03928316815291166020820152604080519182900301906000f0801580156101a1573d6000803e3d6000fd5b509392505050565b6101b1610213565b6101ba57600080fd5b600080546040516001600160a01b03909116907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a3600080546001600160a01b0319169055565b6000546001600160a01b031690565b6000546001600160a01b0316331490565b61022c610213565b61023557600080fd5b61023e81610241565b50565b6001600160a01b03811661025457600080fd5b600080546040516001600160a01b03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a3600080546001600160a01b0319166001600160a01b0392909216919091179055565b6109ac806102bd8339019056fe608060405234801561001057600080fd5b506040516109ac3803806109ac8339818101604052604081101561003357600080fd5b508051602090910151600080546001600160a01b039384166001600160a01b031991821617909155600180549390921692169190911790556109328061007a6000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c8063a8e4fb901161005b578063a8e4fb9014610123578063ab04010714610147578063bade60331461014f578063e9790d02146101235761007d565b806312424e3f146100825780635d1ca6311461008c578063a878f858146100a6575b600080fd5b61008a6101bf565b005b61009461044d565b60408051918252519081900360200190f35b6100ae610454565b6040805160208082528351818301528351919283929083019185019080838360005b838110156100e85781810151838201526020016100d0565b50505050905090810190601f1680156101155780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b61012b6104e7565b604080516001600160a01b039092168252519081900360200190f35b61008a6104f6565b61008a6004803603604081101561016557600080fd5b81019060208101813564010000000081111561018057600080fd5b82018360208201111561019257600080fd5b803590602001918460018302840111640100000000831117156101b457600080fd5b91935091503561071c565b600154604080516301fd3f7760e71b815233600482015290516001600160a01b039092169163fe9fbb8091602480820192602092909190829003018186803b15801561020a57600080fd5b505afa15801561021e573d6000803e3d6000fd5b505050506040513d602081101561023457600080fd5b5051610287576040805162461bcd60e51b815260206004820152601b60248201527f75736572206973206e6f74206d656d626572206f662067726f75700000000000604482015290519081900360640190fd5b600654604080516020808201939093523360601b8183015281516034818303018152605490910182528051908301206000818152600390935291205460ff16156102d1575061044b565b600081815260036020908152604091829020805460ff19166001908117909155548251630859bc9d60e31b815292516001600160a01b03909116926342cde4e8926004808301939192829003018186803b15801561032e57600080fd5b505afa158015610342573d6000803e3d6000fd5b505050506040513d602081101561035857600080fd5b5051600480546001019081905511156104495760018054604051637e2e5ddf60e01b815260206004820190815260028054600019958116156101000295909501909416849004602483018190526001600160a01b0390931693637e2e5ddf9390928291604490910190849080156104105780601f106103e557610100808354040283529160200191610410565b820191906000526020600020905b8154815290600101906020018083116103f357829003601f168201915b505092505050600060405180830381600087803b15801561043057600080fd5b505af1158015610444573d6000803e3d6000fd5b505050505b505b565b6006545b90565b60028054604080516020601f60001961010060018716150201909416859004938401819004810282018101909252828152606093909290918301828280156104dd5780601f106104b2576101008083540402835291602001916104dd565b820191906000526020600020905b8154815290600101906020018083116104c057829003601f168201915b5050505050905090565b6000546001600160a01b031690565b600154604080516301fd3f7760e71b815233600482015290516001600160a01b039092169163fe9fbb8091602480820192602092909190829003018186803b15801561054157600080fd5b505afa158015610555573d6000803e3d6000fd5b505050506040513d602081101561056b57600080fd5b50516105be576040805162461bcd60e51b815260206004820152601b60248201527f75736572206973206e6f74206d656d626572206f662067726f75700000000000604482015290519081900360640190fd5b600654604080516020808201939093523360601b8183015281516034818303018152605490910182528051908301206000818152600390935291205460ff1615610608575061044b565b600081815260036020908152604091829020805460ff19166001908117909155548251630859bc9d60e31b815292516001600160a01b03909116926342cde4e8926004808301939192829003018186803b15801561066557600080fd5b505afa158015610679573d6000803e3d6000fd5b505050506040513d602081101561068f57600080fd5b5051600580546001019081905511156104495760018054604051636f04df4160e01b815260206004820190815260028054600019958116156101000295909501909416849004602483018190526001600160a01b0390931693636f04df419390928291604490910190849080156104105780601f106103e557610100808354040283529160200191610410565b6001546001600160a01b0316331461077b576040805162461bcd60e51b815260206004820152601d60248201527f6d73672e73656e646572206973206e6f742067726f7570206f776e6572000000604482015290519081900360640190fd5b61078760028484610865565b506006819055600160049081556000805460408051638da5cb5b60e01b81529051929385936001600160a01b0390931692638da5cb5b928083019260209291829003018186803b1580156107da57600080fd5b505afa1580156107ee573d6000803e3d6000fd5b505050506040513d602081101561080457600080fd5b50516040805160208082019490945260609290921b6bffffffffffffffffffffffff19168282015280518083036034018152605490920181528151918301919091206000908152600390925290208054600160ff1990911617905550505050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106108a65782800160ff198235161785556108d3565b828001600101855582156108d3579182015b828111156108d35782358255916020019190600101906108b8565b506108df9291506108e3565b5090565b61045191905b808211156108df57600081556001016108e956fea265627a7a7231582010dd144f086299db7b2ab7b8c83e7061e66018a94575eb98bda458c435a3bb9964736f6c634300050f0032a265627a7a72315820034057784801bf36040e0d1422d1cc42bada2e408f1179278f5ff3050f50efc464736f6c634300050f0032"

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
// Solidity: function create(address proposer, address callback) returns(address)
func (_ConsensusFactory *ConsensusFactoryTransactor) Create(opts *bind.TransactOpts, proposer common.Address, callback common.Address) (*types.Transaction, error) {
	return _ConsensusFactory.contract.Transact(opts, "create", proposer, callback)
}

// Create is a paid mutator transaction binding the contract method 0x3e68680a.
//
// Solidity: function create(address proposer, address callback) returns(address)
func (_ConsensusFactory *ConsensusFactorySession) Create(proposer common.Address, callback common.Address) (*types.Transaction, error) {
	return _ConsensusFactory.Contract.Create(&_ConsensusFactory.TransactOpts, proposer, callback)
}

// Create is a paid mutator transaction binding the contract method 0x3e68680a.
//
// Solidity: function create(address proposer, address callback) returns(address)
func (_ConsensusFactory *ConsensusFactoryTransactorSession) Create(proposer common.Address, callback common.Address) (*types.Transaction, error) {
	return _ConsensusFactory.Contract.Create(&_ConsensusFactory.TransactOpts, proposer, callback)
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
