// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Consensus

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

// ConsensusABI is the input ABI used to generate the binding from.
const ConsensusABI = "[{\"inputs\":[{\"internalType\":\"contractIAccount\",\"name\":\"proposer\",\"type\":\"address\"},{\"internalType\":\"contractIConsensusCallback\",\"name\":\"callback\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"state\",\"type\":\"uint256\"}],\"name\":\"Debug\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"msg\",\"type\":\"address\"}],\"name\":\"DebugCons\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"propose\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"approve\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"decline\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getProposer\",\"outputs\":[{\"internalType\":\"contractIAccount\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"payload\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"proposer\",\"outputs\":[{\"internalType\":\"contractIAccount\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// ConsensusBin is the compiled bytecode used for deploying new contracts.
var ConsensusBin = "0x608060405234801561001057600080fd5b506040516109ac3803806109ac8339818101604052604081101561003357600080fd5b508051602090910151600080546001600160a01b039384166001600160a01b031991821617909155600180549390921692169190911790556109328061007a6000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c8063a8e4fb901161005b578063a8e4fb9014610123578063ab04010714610147578063bade60331461014f578063e9790d02146101235761007d565b806312424e3f146100825780635d1ca6311461008c578063a878f858146100a6575b600080fd5b61008a6101bf565b005b61009461044d565b60408051918252519081900360200190f35b6100ae610454565b6040805160208082528351818301528351919283929083019185019080838360005b838110156100e85781810151838201526020016100d0565b50505050905090810190601f1680156101155780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b61012b6104e7565b604080516001600160a01b039092168252519081900360200190f35b61008a6104f6565b61008a6004803603604081101561016557600080fd5b81019060208101813564010000000081111561018057600080fd5b82018360208201111561019257600080fd5b803590602001918460018302840111640100000000831117156101b457600080fd5b91935091503561071c565b600154604080516301fd3f7760e71b815233600482015290516001600160a01b039092169163fe9fbb8091602480820192602092909190829003018186803b15801561020a57600080fd5b505afa15801561021e573d6000803e3d6000fd5b505050506040513d602081101561023457600080fd5b5051610287576040805162461bcd60e51b815260206004820152601b60248201527f75736572206973206e6f74206d656d626572206f662067726f75700000000000604482015290519081900360640190fd5b600654604080516020808201939093523360601b8183015281516034818303018152605490910182528051908301206000818152600390935291205460ff16156102d1575061044b565b600081815260036020908152604091829020805460ff19166001908117909155548251630859bc9d60e31b815292516001600160a01b03909116926342cde4e8926004808301939192829003018186803b15801561032e57600080fd5b505afa158015610342573d6000803e3d6000fd5b505050506040513d602081101561035857600080fd5b5051600480546001019081905511156104495760018054604051637e2e5ddf60e01b815260206004820190815260028054600019958116156101000295909501909416849004602483018190526001600160a01b0390931693637e2e5ddf9390928291604490910190849080156104105780601f106103e557610100808354040283529160200191610410565b820191906000526020600020905b8154815290600101906020018083116103f357829003601f168201915b505092505050600060405180830381600087803b15801561043057600080fd5b505af1158015610444573d6000803e3d6000fd5b505050505b505b565b6006545b90565b60028054604080516020601f60001961010060018716150201909416859004938401819004810282018101909252828152606093909290918301828280156104dd5780601f106104b2576101008083540402835291602001916104dd565b820191906000526020600020905b8154815290600101906020018083116104c057829003601f168201915b5050505050905090565b6000546001600160a01b031690565b600154604080516301fd3f7760e71b815233600482015290516001600160a01b039092169163fe9fbb8091602480820192602092909190829003018186803b15801561054157600080fd5b505afa158015610555573d6000803e3d6000fd5b505050506040513d602081101561056b57600080fd5b50516105be576040805162461bcd60e51b815260206004820152601b60248201527f75736572206973206e6f74206d656d626572206f662067726f75700000000000604482015290519081900360640190fd5b600654604080516020808201939093523360601b8183015281516034818303018152605490910182528051908301206000818152600390935291205460ff1615610608575061044b565b600081815260036020908152604091829020805460ff19166001908117909155548251630859bc9d60e31b815292516001600160a01b03909116926342cde4e8926004808301939192829003018186803b15801561066557600080fd5b505afa158015610679573d6000803e3d6000fd5b505050506040513d602081101561068f57600080fd5b5051600580546001019081905511156104495760018054604051636f04df4160e01b815260206004820190815260028054600019958116156101000295909501909416849004602483018190526001600160a01b0390931693636f04df419390928291604490910190849080156104105780601f106103e557610100808354040283529160200191610410565b6001546001600160a01b0316331461077b576040805162461bcd60e51b815260206004820152601d60248201527f6d73672e73656e646572206973206e6f742067726f7570206f776e6572000000604482015290519081900360640190fd5b61078760028484610865565b506006819055600160049081556000805460408051638da5cb5b60e01b81529051929385936001600160a01b0390931692638da5cb5b928083019260209291829003018186803b1580156107da57600080fd5b505afa1580156107ee573d6000803e3d6000fd5b505050506040513d602081101561080457600080fd5b50516040805160208082019490945260609290921b6bffffffffffffffffffffffff19168282015280518083036034018152605490920181528151918301919091206000908152600390925290208054600160ff1990911617905550505050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106108a65782800160ff198235161785556108d3565b828001600101855582156108d3579182015b828111156108d35782358255916020019190600101906108b8565b506108df9291506108e3565b5090565b61045191905b808211156108df57600081556001016108e956fea265627a7a7231582010dd144f086299db7b2ab7b8c83e7061e66018a94575eb98bda458c435a3bb9964736f6c634300050f0032"

// DeployConsensus deploys a new Ethereum contract, binding an instance of Consensus to it.
func DeployConsensus(auth *bind.TransactOpts, backend bind.ContractBackend, proposer common.Address, callback common.Address) (common.Address, *types.Transaction, *Consensus, error) {
	parsed, err := abi.JSON(strings.NewReader(ConsensusABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ConsensusBin), backend, proposer, callback)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Consensus{ConsensusCaller: ConsensusCaller{contract: contract}, ConsensusTransactor: ConsensusTransactor{contract: contract}, ConsensusFilterer: ConsensusFilterer{contract: contract}}, nil
}

// Consensus is an auto generated Go binding around an Ethereum contract.
type Consensus struct {
	ConsensusCaller     // Read-only binding to the contract
	ConsensusTransactor // Write-only binding to the contract
	ConsensusFilterer   // Log filterer for contract events
}

// ConsensusCaller is an auto generated read-only Go binding around an Ethereum contract.
type ConsensusCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConsensusTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ConsensusTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConsensusFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ConsensusFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConsensusSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ConsensusSession struct {
	Contract     *Consensus        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ConsensusCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ConsensusCallerSession struct {
	Contract *ConsensusCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// ConsensusTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ConsensusTransactorSession struct {
	Contract     *ConsensusTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// ConsensusRaw is an auto generated low-level Go binding around an Ethereum contract.
type ConsensusRaw struct {
	Contract *Consensus // Generic contract binding to access the raw methods on
}

// ConsensusCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ConsensusCallerRaw struct {
	Contract *ConsensusCaller // Generic read-only contract binding to access the raw methods on
}

// ConsensusTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ConsensusTransactorRaw struct {
	Contract *ConsensusTransactor // Generic write-only contract binding to access the raw methods on
}

// NewConsensus creates a new instance of Consensus, bound to a specific deployed contract.
func NewConsensus(address common.Address, backend bind.ContractBackend) (*Consensus, error) {
	contract, err := bindConsensus(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Consensus{ConsensusCaller: ConsensusCaller{contract: contract}, ConsensusTransactor: ConsensusTransactor{contract: contract}, ConsensusFilterer: ConsensusFilterer{contract: contract}}, nil
}

// NewConsensusCaller creates a new read-only instance of Consensus, bound to a specific deployed contract.
func NewConsensusCaller(address common.Address, caller bind.ContractCaller) (*ConsensusCaller, error) {
	contract, err := bindConsensus(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ConsensusCaller{contract: contract}, nil
}

// NewConsensusTransactor creates a new write-only instance of Consensus, bound to a specific deployed contract.
func NewConsensusTransactor(address common.Address, transactor bind.ContractTransactor) (*ConsensusTransactor, error) {
	contract, err := bindConsensus(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ConsensusTransactor{contract: contract}, nil
}

// NewConsensusFilterer creates a new log filterer instance of Consensus, bound to a specific deployed contract.
func NewConsensusFilterer(address common.Address, filterer bind.ContractFilterer) (*ConsensusFilterer, error) {
	contract, err := bindConsensus(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ConsensusFilterer{contract: contract}, nil
}

// bindConsensus binds a generic wrapper to an already deployed contract.
func bindConsensus(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ConsensusABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Consensus *ConsensusRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Consensus.Contract.ConsensusCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Consensus *ConsensusRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Consensus.Contract.ConsensusTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Consensus *ConsensusRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Consensus.Contract.ConsensusTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Consensus *ConsensusCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Consensus.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Consensus *ConsensusTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Consensus.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Consensus *ConsensusTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Consensus.Contract.contract.Transact(opts, method, params...)
}

// GetId is a free data retrieval call binding the contract method 0x5d1ca631.
//
// Solidity: function getId() constant returns(uint256)
func (_Consensus *ConsensusCaller) GetId(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Consensus.contract.Call(opts, out, "getId")
	return *ret0, err
}

// GetId is a free data retrieval call binding the contract method 0x5d1ca631.
//
// Solidity: function getId() constant returns(uint256)
func (_Consensus *ConsensusSession) GetId() (*big.Int, error) {
	return _Consensus.Contract.GetId(&_Consensus.CallOpts)
}

// GetId is a free data retrieval call binding the contract method 0x5d1ca631.
//
// Solidity: function getId() constant returns(uint256)
func (_Consensus *ConsensusCallerSession) GetId() (*big.Int, error) {
	return _Consensus.Contract.GetId(&_Consensus.CallOpts)
}

// GetProposer is a free data retrieval call binding the contract method 0xe9790d02.
//
// Solidity: function getProposer() constant returns(address)
func (_Consensus *ConsensusCaller) GetProposer(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Consensus.contract.Call(opts, out, "getProposer")
	return *ret0, err
}

// GetProposer is a free data retrieval call binding the contract method 0xe9790d02.
//
// Solidity: function getProposer() constant returns(address)
func (_Consensus *ConsensusSession) GetProposer() (common.Address, error) {
	return _Consensus.Contract.GetProposer(&_Consensus.CallOpts)
}

// GetProposer is a free data retrieval call binding the contract method 0xe9790d02.
//
// Solidity: function getProposer() constant returns(address)
func (_Consensus *ConsensusCallerSession) GetProposer() (common.Address, error) {
	return _Consensus.Contract.GetProposer(&_Consensus.CallOpts)
}

// Payload is a free data retrieval call binding the contract method 0xa878f858.
//
// Solidity: function payload() constant returns(bytes)
func (_Consensus *ConsensusCaller) Payload(opts *bind.CallOpts) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _Consensus.contract.Call(opts, out, "payload")
	return *ret0, err
}

// Payload is a free data retrieval call binding the contract method 0xa878f858.
//
// Solidity: function payload() constant returns(bytes)
func (_Consensus *ConsensusSession) Payload() ([]byte, error) {
	return _Consensus.Contract.Payload(&_Consensus.CallOpts)
}

// Payload is a free data retrieval call binding the contract method 0xa878f858.
//
// Solidity: function payload() constant returns(bytes)
func (_Consensus *ConsensusCallerSession) Payload() ([]byte, error) {
	return _Consensus.Contract.Payload(&_Consensus.CallOpts)
}

// Proposer is a free data retrieval call binding the contract method 0xa8e4fb90.
//
// Solidity: function proposer() constant returns(address)
func (_Consensus *ConsensusCaller) Proposer(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Consensus.contract.Call(opts, out, "proposer")
	return *ret0, err
}

// Proposer is a free data retrieval call binding the contract method 0xa8e4fb90.
//
// Solidity: function proposer() constant returns(address)
func (_Consensus *ConsensusSession) Proposer() (common.Address, error) {
	return _Consensus.Contract.Proposer(&_Consensus.CallOpts)
}

// Proposer is a free data retrieval call binding the contract method 0xa8e4fb90.
//
// Solidity: function proposer() constant returns(address)
func (_Consensus *ConsensusCallerSession) Proposer() (common.Address, error) {
	return _Consensus.Contract.Proposer(&_Consensus.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x12424e3f.
//
// Solidity: function approve() returns()
func (_Consensus *ConsensusTransactor) Approve(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Consensus.contract.Transact(opts, "approve")
}

// Approve is a paid mutator transaction binding the contract method 0x12424e3f.
//
// Solidity: function approve() returns()
func (_Consensus *ConsensusSession) Approve() (*types.Transaction, error) {
	return _Consensus.Contract.Approve(&_Consensus.TransactOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x12424e3f.
//
// Solidity: function approve() returns()
func (_Consensus *ConsensusTransactorSession) Approve() (*types.Transaction, error) {
	return _Consensus.Contract.Approve(&_Consensus.TransactOpts)
}

// Decline is a paid mutator transaction binding the contract method 0xab040107.
//
// Solidity: function decline() returns()
func (_Consensus *ConsensusTransactor) Decline(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Consensus.contract.Transact(opts, "decline")
}

// Decline is a paid mutator transaction binding the contract method 0xab040107.
//
// Solidity: function decline() returns()
func (_Consensus *ConsensusSession) Decline() (*types.Transaction, error) {
	return _Consensus.Contract.Decline(&_Consensus.TransactOpts)
}

// Decline is a paid mutator transaction binding the contract method 0xab040107.
//
// Solidity: function decline() returns()
func (_Consensus *ConsensusTransactorSession) Decline() (*types.Transaction, error) {
	return _Consensus.Contract.Decline(&_Consensus.TransactOpts)
}

// Propose is a paid mutator transaction binding the contract method 0xbade6033.
//
// Solidity: function propose(bytes payload, uint256 id) returns()
func (_Consensus *ConsensusTransactor) Propose(opts *bind.TransactOpts, payload []byte, id *big.Int) (*types.Transaction, error) {
	return _Consensus.contract.Transact(opts, "propose", payload, id)
}

// Propose is a paid mutator transaction binding the contract method 0xbade6033.
//
// Solidity: function propose(bytes payload, uint256 id) returns()
func (_Consensus *ConsensusSession) Propose(payload []byte, id *big.Int) (*types.Transaction, error) {
	return _Consensus.Contract.Propose(&_Consensus.TransactOpts, payload, id)
}

// Propose is a paid mutator transaction binding the contract method 0xbade6033.
//
// Solidity: function propose(bytes payload, uint256 id) returns()
func (_Consensus *ConsensusTransactorSession) Propose(payload []byte, id *big.Int) (*types.Transaction, error) {
	return _Consensus.Contract.Propose(&_Consensus.TransactOpts, payload, id)
}

// ConsensusDebugIterator is returned from FilterDebug and is used to iterate over the raw logs and unpacked data for Debug events raised by the Consensus contract.
type ConsensusDebugIterator struct {
	Event *ConsensusDebug // Event containing the contract specifics and raw log

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
func (it *ConsensusDebugIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConsensusDebug)
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
		it.Event = new(ConsensusDebug)
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
func (it *ConsensusDebugIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConsensusDebugIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConsensusDebug represents a Debug event raised by the Consensus contract.
type ConsensusDebug struct {
	State *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterDebug is a free log retrieval operation binding the contract event 0x8a36f5a234186d446e36a7df36ace663a05a580d9bea2dd899c6dd76a075d5fa.
//
// Solidity: event Debug(uint256 state)
func (_Consensus *ConsensusFilterer) FilterDebug(opts *bind.FilterOpts) (*ConsensusDebugIterator, error) {

	logs, sub, err := _Consensus.contract.FilterLogs(opts, "Debug")
	if err != nil {
		return nil, err
	}
	return &ConsensusDebugIterator{contract: _Consensus.contract, event: "Debug", logs: logs, sub: sub}, nil
}

// WatchDebug is a free log subscription operation binding the contract event 0x8a36f5a234186d446e36a7df36ace663a05a580d9bea2dd899c6dd76a075d5fa.
//
// Solidity: event Debug(uint256 state)
func (_Consensus *ConsensusFilterer) WatchDebug(opts *bind.WatchOpts, sink chan<- *ConsensusDebug) (event.Subscription, error) {

	logs, sub, err := _Consensus.contract.WatchLogs(opts, "Debug")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConsensusDebug)
				if err := _Consensus.contract.UnpackLog(event, "Debug", log); err != nil {
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

// ParseDebug is a log parse operation binding the contract event 0x8a36f5a234186d446e36a7df36ace663a05a580d9bea2dd899c6dd76a075d5fa.
//
// Solidity: event Debug(uint256 state)
func (_Consensus *ConsensusFilterer) ParseDebug(log types.Log) (*ConsensusDebug, error) {
	event := new(ConsensusDebug)
	if err := _Consensus.contract.UnpackLog(event, "Debug", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ConsensusDebugConsIterator is returned from FilterDebugCons and is used to iterate over the raw logs and unpacked data for DebugCons events raised by the Consensus contract.
type ConsensusDebugConsIterator struct {
	Event *ConsensusDebugCons // Event containing the contract specifics and raw log

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
func (it *ConsensusDebugConsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConsensusDebugCons)
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
		it.Event = new(ConsensusDebugCons)
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
func (it *ConsensusDebugConsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConsensusDebugConsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConsensusDebugCons represents a DebugCons event raised by the Consensus contract.
type ConsensusDebugCons struct {
	Msg common.Address
	Raw types.Log // Blockchain specific contextual infos
}

// FilterDebugCons is a free log retrieval operation binding the contract event 0xdca97acc47326d0d02d07dbb1517e60b681494c362f8de4ac2e3398fbd22f130.
//
// Solidity: event DebugCons(address msg)
func (_Consensus *ConsensusFilterer) FilterDebugCons(opts *bind.FilterOpts) (*ConsensusDebugConsIterator, error) {

	logs, sub, err := _Consensus.contract.FilterLogs(opts, "DebugCons")
	if err != nil {
		return nil, err
	}
	return &ConsensusDebugConsIterator{contract: _Consensus.contract, event: "DebugCons", logs: logs, sub: sub}, nil
}

// WatchDebugCons is a free log subscription operation binding the contract event 0xdca97acc47326d0d02d07dbb1517e60b681494c362f8de4ac2e3398fbd22f130.
//
// Solidity: event DebugCons(address msg)
func (_Consensus *ConsensusFilterer) WatchDebugCons(opts *bind.WatchOpts, sink chan<- *ConsensusDebugCons) (event.Subscription, error) {

	logs, sub, err := _Consensus.contract.WatchLogs(opts, "DebugCons")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConsensusDebugCons)
				if err := _Consensus.contract.UnpackLog(event, "DebugCons", log); err != nil {
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

// ParseDebugCons is a log parse operation binding the contract event 0xdca97acc47326d0d02d07dbb1517e60b681494c362f8de4ac2e3398fbd22f130.
//
// Solidity: event DebugCons(address msg)
func (_Consensus *ConsensusFilterer) ParseDebugCons(log types.Log) (*ConsensusDebugCons, error) {
	event := new(ConsensusDebugCons)
	if err := _Consensus.contract.UnpackLog(event, "DebugCons", log); err != nil {
		return nil, err
	}
	return event, nil
}
