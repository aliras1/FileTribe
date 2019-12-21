// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package FileTribeDApp

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

// FileTribeDAppABI is the input ABI used to generate the binding from.
const FileTribeDAppABI = "[{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractIAccount\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"AccountCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"msg\",\"type\":\"int256\"}],\"name\":\"Debug\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"msg\",\"type\":\"bytes\"}],\"name\":\"DebugBytes\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"GroupRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"contractIConsensusFactory\",\"name\":\"factory\",\"type\":\"address\"}],\"name\":\"setConsensusFactory\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"contractIAccountFactory\",\"name\":\"factory\",\"type\":\"address\"}],\"name\":\"setAccountFactory\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"contractIGroupFactory\",\"name\":\"factory\",\"type\":\"address\"}],\"name\":\"setGroupFactory\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"contractDkgFactory\",\"name\":\"factory\",\"type\":\"address\"}],\"name\":\"setDkgFactory\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"ipfsPeerId\",\"type\":\"string\"},{\"internalType\":\"bytes32\",\"name\":\"boxingKey\",\"type\":\"bytes32\"}],\"name\":\"createAccount\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"removeAccount\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"createGroup\",\"outputs\":[{\"internalType\":\"contractIGroup\",\"name\":\"group\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"contractIAccount\",\"name\":\"proposer\",\"type\":\"address\"}],\"name\":\"createConsensus\",\"outputs\":[{\"internalType\":\"contractIConsensus\",\"name\":\"consensus\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"createDkg\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"getAccountOf\",\"outputs\":[{\"internalType\":\"contractIAccount\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// FileTribeDAppBin is the compiled bytecode used for deploying new contracts.
var FileTribeDAppBin = "0x608060405234801561001057600080fd5b50600080546001600160a01b0319163390811780835560405191926001600160a01b0391909116917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908290a350610a188061006d6000396000f3fe608060405234801561001057600080fd5b50600436106100ea5760003560e01c8063837b3b931161008c57806395184d3b1161006657806395184d3b14610310578063addc1a7614610336578063dc2ddcae1461035c578063f2fde38b146103cc576100ea565b8063837b3b93146102c65780638da5cb5b146102ec5780638f32d59b146102f4576100ea565b80633c2aea2e116100c85780633c2aea2e1461017f57806345aa926f146102ae578063715018a6146102b657806375beca08146102be576100ea565b806315db2069146100ef578063186cd722146101315780631c38ea5614610159575b600080fd5b6101156004803603602081101561010557600080fd5b50356001600160a01b03166103f2565b604080516001600160a01b039092168252519081900360200190f35b6101576004803603602081101561014757600080fd5b50356001600160a01b031661047c565b005b6101576004803603602081101561016f57600080fd5b50356001600160a01b03166104af565b6101576004803603606081101561019557600080fd5b8101906020810181356401000000008111156101b057600080fd5b8201836020820111156101c257600080fd5b803590602001918460018302840111640100000000831117156101e457600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929594936020810193503591505064010000000081111561023757600080fd5b82018360208201111561024957600080fd5b8035906020019184600183028401116401000000008311171561026b57600080fd5b91908080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525092955050913592506104e2915050565b610115610706565b610157610787565b6101576107e2565b610157600480360360208110156102dc57600080fd5b50356001600160a01b0316610801565b610115610834565b6102fc610843565b604080519115158252519081900360200190f35b6101156004803603602081101561032657600080fd5b50356001600160a01b0316610854565b6101576004803603602081101561034c57600080fd5b50356001600160a01b0316610872565b6101156004803603602081101561037257600080fd5b81019060208101813564010000000081111561038d57600080fd5b82018360208201111561039f57600080fd5b803590602001918460018302840111640100000000831117156103c157600080fd5b5090925090506108a5565b610157600480360360208110156103e257600080fd5b50356001600160a01b0316610958565b60035460408051631f34340560e11b81526001600160a01b03848116600483015233602483015291516000939290921691633e68680a9160448082019260209290919082900301818787803b15801561044a57600080fd5b505af115801561045e573d6000803e3d6000fd5b505050506040513d602081101561047457600080fd5b505192915050565b610484610843565b61048d57600080fd5b600480546001600160a01b0319166001600160a01b0392909216919091179055565b6104b7610843565b6104c057600080fd5b600380546001600160a01b0319166001600160a01b0392909216919091179055565b336000908152600560205260409020546001600160a01b031615610546576040805162461bcd60e51b81526020600482015260166024820152754163636f756e7420616c72656164792065786973747360501b604482015290519081900360640190fd5b60015460405163d0bf552b60e01b81523360048201818152606483018590526080602484019081528751608485015287516000956001600160a01b03169463d0bf552b94938a938a938a93604481019160a49091019060208801908083838f5b838110156105be5781810151838201526020016105a6565b50505050905090810190601f1680156105eb5780820380516001836020036101000a031916815260200191505b50838103825285518152855160209182019187019080838360005b8381101561061e578181015183820152602001610606565b50505050905090810190601f16801561064b5780820380516001836020036101000a031916815260200191505b509650505050505050602060405180830381600087803b15801561066e57600080fd5b505af1158015610682573d6000803e3d6000fd5b505050506040513d602081101561069857600080fd5b50513360008181526005602090815260409182902080546001600160a01b0319166001600160a01b03861690811790915582519384529083015280519293507fac631f3001b55ea1509cf3d7e74898f85392a61a76e8149181ae1259622dabc892918290030190a150505050565b60048054604080516313db266360e31b81523393810193909352516000926001600160a01b0390921691639ed9331891602480830192602092919082900301818787803b15801561075657600080fd5b505af115801561076a573d6000803e3d6000fd5b505050506040513d602081101561078057600080fd5b5051905090565b61078f610843565b61079857600080fd5b600080546040516001600160a01b03909116907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a3600080546001600160a01b0319169055565b33600090815260056020526040902080546001600160a01b0319169055565b610809610843565b61081257600080fd5b600280546001600160a01b0319166001600160a01b0392909216919091179055565b6000546001600160a01b031690565b6000546001600160a01b0316331490565b6001600160a01b039081166000908152600560205260409020541690565b61087a610843565b61088357600080fd5b600180546001600160a01b0319166001600160a01b0392909216919091179055565b6002546040805163a15ab08d60e01b8152336004820181815260248301938452604483018690526000946001600160a01b03169363a15ab08d93889288929091606401848480828437600081840152601f19601f820116905080830192505050945050505050602060405180830381600087803b15801561092557600080fd5b505af1158015610939573d6000803e3d6000fd5b505050506040513d602081101561094f57600080fd5b50519392505050565b610960610843565b61096957600080fd5b61097281610975565b50565b6001600160a01b03811661098857600080fd5b600080546040516001600160a01b03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a3600080546001600160a01b0319166001600160a01b039290921691909117905556fea265627a7a7231582067b7b7e710ffed01e06d01cc12ba8ca49e3fd16a2e6a6ab4a2eb83fe0bf6d71a64736f6c634300050f0032"

// DeployFileTribeDApp deploys a new Ethereum contract, binding an instance of FileTribeDApp to it.
func DeployFileTribeDApp(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *FileTribeDApp, error) {
	parsed, err := abi.JSON(strings.NewReader(FileTribeDAppABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(FileTribeDAppBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &FileTribeDApp{FileTribeDAppCaller: FileTribeDAppCaller{contract: contract}, FileTribeDAppTransactor: FileTribeDAppTransactor{contract: contract}, FileTribeDAppFilterer: FileTribeDAppFilterer{contract: contract}}, nil
}

// FileTribeDApp is an auto generated Go binding around an Ethereum contract.
type FileTribeDApp struct {
	FileTribeDAppCaller     // Read-only binding to the contract
	FileTribeDAppTransactor // Write-only binding to the contract
	FileTribeDAppFilterer   // Log filterer for contract events
}

// FileTribeDAppCaller is an auto generated read-only Go binding around an Ethereum contract.
type FileTribeDAppCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FileTribeDAppTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FileTribeDAppTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FileTribeDAppFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FileTribeDAppFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FileTribeDAppSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FileTribeDAppSession struct {
	Contract     *FileTribeDApp    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FileTribeDAppCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FileTribeDAppCallerSession struct {
	Contract *FileTribeDAppCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// FileTribeDAppTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FileTribeDAppTransactorSession struct {
	Contract     *FileTribeDAppTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// FileTribeDAppRaw is an auto generated low-level Go binding around an Ethereum contract.
type FileTribeDAppRaw struct {
	Contract *FileTribeDApp // Generic contract binding to access the raw methods on
}

// FileTribeDAppCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FileTribeDAppCallerRaw struct {
	Contract *FileTribeDAppCaller // Generic read-only contract binding to access the raw methods on
}

// FileTribeDAppTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FileTribeDAppTransactorRaw struct {
	Contract *FileTribeDAppTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFileTribeDApp creates a new instance of FileTribeDApp, bound to a specific deployed contract.
func NewFileTribeDApp(address common.Address, backend bind.ContractBackend) (*FileTribeDApp, error) {
	contract, err := bindFileTribeDApp(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FileTribeDApp{FileTribeDAppCaller: FileTribeDAppCaller{contract: contract}, FileTribeDAppTransactor: FileTribeDAppTransactor{contract: contract}, FileTribeDAppFilterer: FileTribeDAppFilterer{contract: contract}}, nil
}

// NewFileTribeDAppCaller creates a new read-only instance of FileTribeDApp, bound to a specific deployed contract.
func NewFileTribeDAppCaller(address common.Address, caller bind.ContractCaller) (*FileTribeDAppCaller, error) {
	contract, err := bindFileTribeDApp(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FileTribeDAppCaller{contract: contract}, nil
}

// NewFileTribeDAppTransactor creates a new write-only instance of FileTribeDApp, bound to a specific deployed contract.
func NewFileTribeDAppTransactor(address common.Address, transactor bind.ContractTransactor) (*FileTribeDAppTransactor, error) {
	contract, err := bindFileTribeDApp(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FileTribeDAppTransactor{contract: contract}, nil
}

// NewFileTribeDAppFilterer creates a new log filterer instance of FileTribeDApp, bound to a specific deployed contract.
func NewFileTribeDAppFilterer(address common.Address, filterer bind.ContractFilterer) (*FileTribeDAppFilterer, error) {
	contract, err := bindFileTribeDApp(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FileTribeDAppFilterer{contract: contract}, nil
}

// bindFileTribeDApp binds a generic wrapper to an already deployed contract.
func bindFileTribeDApp(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(FileTribeDAppABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FileTribeDApp *FileTribeDAppRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _FileTribeDApp.Contract.FileTribeDAppCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FileTribeDApp *FileTribeDAppRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.FileTribeDAppTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FileTribeDApp *FileTribeDAppRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.FileTribeDAppTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FileTribeDApp *FileTribeDAppCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _FileTribeDApp.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FileTribeDApp *FileTribeDAppTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FileTribeDApp *FileTribeDAppTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.contract.Transact(opts, method, params...)
}

// GetAccountOf is a free data retrieval call binding the contract method 0x95184d3b.
//
// Solidity: function getAccountOf(address owner) constant returns(address)
func (_FileTribeDApp *FileTribeDAppCaller) GetAccountOf(opts *bind.CallOpts, owner common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _FileTribeDApp.contract.Call(opts, out, "getAccountOf", owner)
	return *ret0, err
}

// GetAccountOf is a free data retrieval call binding the contract method 0x95184d3b.
//
// Solidity: function getAccountOf(address owner) constant returns(address)
func (_FileTribeDApp *FileTribeDAppSession) GetAccountOf(owner common.Address) (common.Address, error) {
	return _FileTribeDApp.Contract.GetAccountOf(&_FileTribeDApp.CallOpts, owner)
}

// GetAccountOf is a free data retrieval call binding the contract method 0x95184d3b.
//
// Solidity: function getAccountOf(address owner) constant returns(address)
func (_FileTribeDApp *FileTribeDAppCallerSession) GetAccountOf(owner common.Address) (common.Address, error) {
	return _FileTribeDApp.Contract.GetAccountOf(&_FileTribeDApp.CallOpts, owner)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_FileTribeDApp *FileTribeDAppCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _FileTribeDApp.contract.Call(opts, out, "isOwner")
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_FileTribeDApp *FileTribeDAppSession) IsOwner() (bool, error) {
	return _FileTribeDApp.Contract.IsOwner(&_FileTribeDApp.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_FileTribeDApp *FileTribeDAppCallerSession) IsOwner() (bool, error) {
	return _FileTribeDApp.Contract.IsOwner(&_FileTribeDApp.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_FileTribeDApp *FileTribeDAppCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _FileTribeDApp.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_FileTribeDApp *FileTribeDAppSession) Owner() (common.Address, error) {
	return _FileTribeDApp.Contract.Owner(&_FileTribeDApp.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_FileTribeDApp *FileTribeDAppCallerSession) Owner() (common.Address, error) {
	return _FileTribeDApp.Contract.Owner(&_FileTribeDApp.CallOpts)
}

// CreateAccount is a paid mutator transaction binding the contract method 0x3c2aea2e.
//
// Solidity: function createAccount(string name, string ipfsPeerId, bytes32 boxingKey) returns()
func (_FileTribeDApp *FileTribeDAppTransactor) CreateAccount(opts *bind.TransactOpts, name string, ipfsPeerId string, boxingKey [32]byte) (*types.Transaction, error) {
	return _FileTribeDApp.contract.Transact(opts, "createAccount", name, ipfsPeerId, boxingKey)
}

// CreateAccount is a paid mutator transaction binding the contract method 0x3c2aea2e.
//
// Solidity: function createAccount(string name, string ipfsPeerId, bytes32 boxingKey) returns()
func (_FileTribeDApp *FileTribeDAppSession) CreateAccount(name string, ipfsPeerId string, boxingKey [32]byte) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.CreateAccount(&_FileTribeDApp.TransactOpts, name, ipfsPeerId, boxingKey)
}

// CreateAccount is a paid mutator transaction binding the contract method 0x3c2aea2e.
//
// Solidity: function createAccount(string name, string ipfsPeerId, bytes32 boxingKey) returns()
func (_FileTribeDApp *FileTribeDAppTransactorSession) CreateAccount(name string, ipfsPeerId string, boxingKey [32]byte) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.CreateAccount(&_FileTribeDApp.TransactOpts, name, ipfsPeerId, boxingKey)
}

// CreateConsensus is a paid mutator transaction binding the contract method 0x15db2069.
//
// Solidity: function createConsensus(address proposer) returns(address consensus)
func (_FileTribeDApp *FileTribeDAppTransactor) CreateConsensus(opts *bind.TransactOpts, proposer common.Address) (*types.Transaction, error) {
	return _FileTribeDApp.contract.Transact(opts, "createConsensus", proposer)
}

// CreateConsensus is a paid mutator transaction binding the contract method 0x15db2069.
//
// Solidity: function createConsensus(address proposer) returns(address consensus)
func (_FileTribeDApp *FileTribeDAppSession) CreateConsensus(proposer common.Address) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.CreateConsensus(&_FileTribeDApp.TransactOpts, proposer)
}

// CreateConsensus is a paid mutator transaction binding the contract method 0x15db2069.
//
// Solidity: function createConsensus(address proposer) returns(address consensus)
func (_FileTribeDApp *FileTribeDAppTransactorSession) CreateConsensus(proposer common.Address) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.CreateConsensus(&_FileTribeDApp.TransactOpts, proposer)
}

// CreateDkg is a paid mutator transaction binding the contract method 0x45aa926f.
//
// Solidity: function createDkg() returns(address)
func (_FileTribeDApp *FileTribeDAppTransactor) CreateDkg(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FileTribeDApp.contract.Transact(opts, "createDkg")
}

// CreateDkg is a paid mutator transaction binding the contract method 0x45aa926f.
//
// Solidity: function createDkg() returns(address)
func (_FileTribeDApp *FileTribeDAppSession) CreateDkg() (*types.Transaction, error) {
	return _FileTribeDApp.Contract.CreateDkg(&_FileTribeDApp.TransactOpts)
}

// CreateDkg is a paid mutator transaction binding the contract method 0x45aa926f.
//
// Solidity: function createDkg() returns(address)
func (_FileTribeDApp *FileTribeDAppTransactorSession) CreateDkg() (*types.Transaction, error) {
	return _FileTribeDApp.Contract.CreateDkg(&_FileTribeDApp.TransactOpts)
}

// CreateGroup is a paid mutator transaction binding the contract method 0xdc2ddcae.
//
// Solidity: function createGroup(string name) returns(address group)
func (_FileTribeDApp *FileTribeDAppTransactor) CreateGroup(opts *bind.TransactOpts, name string) (*types.Transaction, error) {
	return _FileTribeDApp.contract.Transact(opts, "createGroup", name)
}

// CreateGroup is a paid mutator transaction binding the contract method 0xdc2ddcae.
//
// Solidity: function createGroup(string name) returns(address group)
func (_FileTribeDApp *FileTribeDAppSession) CreateGroup(name string) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.CreateGroup(&_FileTribeDApp.TransactOpts, name)
}

// CreateGroup is a paid mutator transaction binding the contract method 0xdc2ddcae.
//
// Solidity: function createGroup(string name) returns(address group)
func (_FileTribeDApp *FileTribeDAppTransactorSession) CreateGroup(name string) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.CreateGroup(&_FileTribeDApp.TransactOpts, name)
}

// RemoveAccount is a paid mutator transaction binding the contract method 0x75beca08.
//
// Solidity: function removeAccount() returns()
func (_FileTribeDApp *FileTribeDAppTransactor) RemoveAccount(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FileTribeDApp.contract.Transact(opts, "removeAccount")
}

// RemoveAccount is a paid mutator transaction binding the contract method 0x75beca08.
//
// Solidity: function removeAccount() returns()
func (_FileTribeDApp *FileTribeDAppSession) RemoveAccount() (*types.Transaction, error) {
	return _FileTribeDApp.Contract.RemoveAccount(&_FileTribeDApp.TransactOpts)
}

// RemoveAccount is a paid mutator transaction binding the contract method 0x75beca08.
//
// Solidity: function removeAccount() returns()
func (_FileTribeDApp *FileTribeDAppTransactorSession) RemoveAccount() (*types.Transaction, error) {
	return _FileTribeDApp.Contract.RemoveAccount(&_FileTribeDApp.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FileTribeDApp *FileTribeDAppTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FileTribeDApp.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FileTribeDApp *FileTribeDAppSession) RenounceOwnership() (*types.Transaction, error) {
	return _FileTribeDApp.Contract.RenounceOwnership(&_FileTribeDApp.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FileTribeDApp *FileTribeDAppTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _FileTribeDApp.Contract.RenounceOwnership(&_FileTribeDApp.TransactOpts)
}

// SetAccountFactory is a paid mutator transaction binding the contract method 0xaddc1a76.
//
// Solidity: function setAccountFactory(address factory) returns()
func (_FileTribeDApp *FileTribeDAppTransactor) SetAccountFactory(opts *bind.TransactOpts, factory common.Address) (*types.Transaction, error) {
	return _FileTribeDApp.contract.Transact(opts, "setAccountFactory", factory)
}

// SetAccountFactory is a paid mutator transaction binding the contract method 0xaddc1a76.
//
// Solidity: function setAccountFactory(address factory) returns()
func (_FileTribeDApp *FileTribeDAppSession) SetAccountFactory(factory common.Address) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.SetAccountFactory(&_FileTribeDApp.TransactOpts, factory)
}

// SetAccountFactory is a paid mutator transaction binding the contract method 0xaddc1a76.
//
// Solidity: function setAccountFactory(address factory) returns()
func (_FileTribeDApp *FileTribeDAppTransactorSession) SetAccountFactory(factory common.Address) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.SetAccountFactory(&_FileTribeDApp.TransactOpts, factory)
}

// SetConsensusFactory is a paid mutator transaction binding the contract method 0x1c38ea56.
//
// Solidity: function setConsensusFactory(address factory) returns()
func (_FileTribeDApp *FileTribeDAppTransactor) SetConsensusFactory(opts *bind.TransactOpts, factory common.Address) (*types.Transaction, error) {
	return _FileTribeDApp.contract.Transact(opts, "setConsensusFactory", factory)
}

// SetConsensusFactory is a paid mutator transaction binding the contract method 0x1c38ea56.
//
// Solidity: function setConsensusFactory(address factory) returns()
func (_FileTribeDApp *FileTribeDAppSession) SetConsensusFactory(factory common.Address) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.SetConsensusFactory(&_FileTribeDApp.TransactOpts, factory)
}

// SetConsensusFactory is a paid mutator transaction binding the contract method 0x1c38ea56.
//
// Solidity: function setConsensusFactory(address factory) returns()
func (_FileTribeDApp *FileTribeDAppTransactorSession) SetConsensusFactory(factory common.Address) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.SetConsensusFactory(&_FileTribeDApp.TransactOpts, factory)
}

// SetDkgFactory is a paid mutator transaction binding the contract method 0x186cd722.
//
// Solidity: function setDkgFactory(address factory) returns()
func (_FileTribeDApp *FileTribeDAppTransactor) SetDkgFactory(opts *bind.TransactOpts, factory common.Address) (*types.Transaction, error) {
	return _FileTribeDApp.contract.Transact(opts, "setDkgFactory", factory)
}

// SetDkgFactory is a paid mutator transaction binding the contract method 0x186cd722.
//
// Solidity: function setDkgFactory(address factory) returns()
func (_FileTribeDApp *FileTribeDAppSession) SetDkgFactory(factory common.Address) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.SetDkgFactory(&_FileTribeDApp.TransactOpts, factory)
}

// SetDkgFactory is a paid mutator transaction binding the contract method 0x186cd722.
//
// Solidity: function setDkgFactory(address factory) returns()
func (_FileTribeDApp *FileTribeDAppTransactorSession) SetDkgFactory(factory common.Address) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.SetDkgFactory(&_FileTribeDApp.TransactOpts, factory)
}

// SetGroupFactory is a paid mutator transaction binding the contract method 0x837b3b93.
//
// Solidity: function setGroupFactory(address factory) returns()
func (_FileTribeDApp *FileTribeDAppTransactor) SetGroupFactory(opts *bind.TransactOpts, factory common.Address) (*types.Transaction, error) {
	return _FileTribeDApp.contract.Transact(opts, "setGroupFactory", factory)
}

// SetGroupFactory is a paid mutator transaction binding the contract method 0x837b3b93.
//
// Solidity: function setGroupFactory(address factory) returns()
func (_FileTribeDApp *FileTribeDAppSession) SetGroupFactory(factory common.Address) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.SetGroupFactory(&_FileTribeDApp.TransactOpts, factory)
}

// SetGroupFactory is a paid mutator transaction binding the contract method 0x837b3b93.
//
// Solidity: function setGroupFactory(address factory) returns()
func (_FileTribeDApp *FileTribeDAppTransactorSession) SetGroupFactory(factory common.Address) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.SetGroupFactory(&_FileTribeDApp.TransactOpts, factory)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FileTribeDApp *FileTribeDAppTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _FileTribeDApp.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FileTribeDApp *FileTribeDAppSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.TransferOwnership(&_FileTribeDApp.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FileTribeDApp *FileTribeDAppTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _FileTribeDApp.Contract.TransferOwnership(&_FileTribeDApp.TransactOpts, newOwner)
}

// FileTribeDAppAccountCreatedIterator is returned from FilterAccountCreated and is used to iterate over the raw logs and unpacked data for AccountCreated events raised by the FileTribeDApp contract.
type FileTribeDAppAccountCreatedIterator struct {
	Event *FileTribeDAppAccountCreated // Event containing the contract specifics and raw log

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
func (it *FileTribeDAppAccountCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FileTribeDAppAccountCreated)
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
		it.Event = new(FileTribeDAppAccountCreated)
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
func (it *FileTribeDAppAccountCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FileTribeDAppAccountCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FileTribeDAppAccountCreated represents a AccountCreated event raised by the FileTribeDApp contract.
type FileTribeDAppAccountCreated struct {
	Owner   common.Address
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAccountCreated is a free log retrieval operation binding the contract event 0xac631f3001b55ea1509cf3d7e74898f85392a61a76e8149181ae1259622dabc8.
//
// Solidity: event AccountCreated(address owner, address account)
func (_FileTribeDApp *FileTribeDAppFilterer) FilterAccountCreated(opts *bind.FilterOpts) (*FileTribeDAppAccountCreatedIterator, error) {

	logs, sub, err := _FileTribeDApp.contract.FilterLogs(opts, "AccountCreated")
	if err != nil {
		return nil, err
	}
	return &FileTribeDAppAccountCreatedIterator{contract: _FileTribeDApp.contract, event: "AccountCreated", logs: logs, sub: sub}, nil
}

// WatchAccountCreated is a free log subscription operation binding the contract event 0xac631f3001b55ea1509cf3d7e74898f85392a61a76e8149181ae1259622dabc8.
//
// Solidity: event AccountCreated(address owner, address account)
func (_FileTribeDApp *FileTribeDAppFilterer) WatchAccountCreated(opts *bind.WatchOpts, sink chan<- *FileTribeDAppAccountCreated) (event.Subscription, error) {

	logs, sub, err := _FileTribeDApp.contract.WatchLogs(opts, "AccountCreated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FileTribeDAppAccountCreated)
				if err := _FileTribeDApp.contract.UnpackLog(event, "AccountCreated", log); err != nil {
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

// ParseAccountCreated is a log parse operation binding the contract event 0xac631f3001b55ea1509cf3d7e74898f85392a61a76e8149181ae1259622dabc8.
//
// Solidity: event AccountCreated(address owner, address account)
func (_FileTribeDApp *FileTribeDAppFilterer) ParseAccountCreated(log types.Log) (*FileTribeDAppAccountCreated, error) {
	event := new(FileTribeDAppAccountCreated)
	if err := _FileTribeDApp.contract.UnpackLog(event, "AccountCreated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// FileTribeDAppDebugIterator is returned from FilterDebug and is used to iterate over the raw logs and unpacked data for Debug events raised by the FileTribeDApp contract.
type FileTribeDAppDebugIterator struct {
	Event *FileTribeDAppDebug // Event containing the contract specifics and raw log

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
func (it *FileTribeDAppDebugIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FileTribeDAppDebug)
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
		it.Event = new(FileTribeDAppDebug)
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
func (it *FileTribeDAppDebugIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FileTribeDAppDebugIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FileTribeDAppDebug represents a Debug event raised by the FileTribeDApp contract.
type FileTribeDAppDebug struct {
	Msg *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterDebug is a free log retrieval operation binding the contract event 0xaad112e1b1cb380614d4d2c495b213d9ed63560382320d6a49cc64793149f93c.
//
// Solidity: event Debug(int256 msg)
func (_FileTribeDApp *FileTribeDAppFilterer) FilterDebug(opts *bind.FilterOpts) (*FileTribeDAppDebugIterator, error) {

	logs, sub, err := _FileTribeDApp.contract.FilterLogs(opts, "Debug")
	if err != nil {
		return nil, err
	}
	return &FileTribeDAppDebugIterator{contract: _FileTribeDApp.contract, event: "Debug", logs: logs, sub: sub}, nil
}

// WatchDebug is a free log subscription operation binding the contract event 0xaad112e1b1cb380614d4d2c495b213d9ed63560382320d6a49cc64793149f93c.
//
// Solidity: event Debug(int256 msg)
func (_FileTribeDApp *FileTribeDAppFilterer) WatchDebug(opts *bind.WatchOpts, sink chan<- *FileTribeDAppDebug) (event.Subscription, error) {

	logs, sub, err := _FileTribeDApp.contract.WatchLogs(opts, "Debug")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FileTribeDAppDebug)
				if err := _FileTribeDApp.contract.UnpackLog(event, "Debug", log); err != nil {
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

// ParseDebug is a log parse operation binding the contract event 0xaad112e1b1cb380614d4d2c495b213d9ed63560382320d6a49cc64793149f93c.
//
// Solidity: event Debug(int256 msg)
func (_FileTribeDApp *FileTribeDAppFilterer) ParseDebug(log types.Log) (*FileTribeDAppDebug, error) {
	event := new(FileTribeDAppDebug)
	if err := _FileTribeDApp.contract.UnpackLog(event, "Debug", log); err != nil {
		return nil, err
	}
	return event, nil
}

// FileTribeDAppDebugBytesIterator is returned from FilterDebugBytes and is used to iterate over the raw logs and unpacked data for DebugBytes events raised by the FileTribeDApp contract.
type FileTribeDAppDebugBytesIterator struct {
	Event *FileTribeDAppDebugBytes // Event containing the contract specifics and raw log

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
func (it *FileTribeDAppDebugBytesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FileTribeDAppDebugBytes)
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
		it.Event = new(FileTribeDAppDebugBytes)
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
func (it *FileTribeDAppDebugBytesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FileTribeDAppDebugBytesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FileTribeDAppDebugBytes represents a DebugBytes event raised by the FileTribeDApp contract.
type FileTribeDAppDebugBytes struct {
	Msg []byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterDebugBytes is a free log retrieval operation binding the contract event 0xaf5a5af90a78ece430d7df503b54fc4070844db69884a9a4afb00710a4816e53.
//
// Solidity: event DebugBytes(bytes msg)
func (_FileTribeDApp *FileTribeDAppFilterer) FilterDebugBytes(opts *bind.FilterOpts) (*FileTribeDAppDebugBytesIterator, error) {

	logs, sub, err := _FileTribeDApp.contract.FilterLogs(opts, "DebugBytes")
	if err != nil {
		return nil, err
	}
	return &FileTribeDAppDebugBytesIterator{contract: _FileTribeDApp.contract, event: "DebugBytes", logs: logs, sub: sub}, nil
}

// WatchDebugBytes is a free log subscription operation binding the contract event 0xaf5a5af90a78ece430d7df503b54fc4070844db69884a9a4afb00710a4816e53.
//
// Solidity: event DebugBytes(bytes msg)
func (_FileTribeDApp *FileTribeDAppFilterer) WatchDebugBytes(opts *bind.WatchOpts, sink chan<- *FileTribeDAppDebugBytes) (event.Subscription, error) {

	logs, sub, err := _FileTribeDApp.contract.WatchLogs(opts, "DebugBytes")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FileTribeDAppDebugBytes)
				if err := _FileTribeDApp.contract.UnpackLog(event, "DebugBytes", log); err != nil {
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

// ParseDebugBytes is a log parse operation binding the contract event 0xaf5a5af90a78ece430d7df503b54fc4070844db69884a9a4afb00710a4816e53.
//
// Solidity: event DebugBytes(bytes msg)
func (_FileTribeDApp *FileTribeDAppFilterer) ParseDebugBytes(log types.Log) (*FileTribeDAppDebugBytes, error) {
	event := new(FileTribeDAppDebugBytes)
	if err := _FileTribeDApp.contract.UnpackLog(event, "DebugBytes", log); err != nil {
		return nil, err
	}
	return event, nil
}

// FileTribeDAppGroupRegisteredIterator is returned from FilterGroupRegistered and is used to iterate over the raw logs and unpacked data for GroupRegistered events raised by the FileTribeDApp contract.
type FileTribeDAppGroupRegisteredIterator struct {
	Event *FileTribeDAppGroupRegistered // Event containing the contract specifics and raw log

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
func (it *FileTribeDAppGroupRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FileTribeDAppGroupRegistered)
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
		it.Event = new(FileTribeDAppGroupRegistered)
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
func (it *FileTribeDAppGroupRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FileTribeDAppGroupRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FileTribeDAppGroupRegistered represents a GroupRegistered event raised by the FileTribeDApp contract.
type FileTribeDAppGroupRegistered struct {
	Id  [32]byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterGroupRegistered is a free log retrieval operation binding the contract event 0xb78b34f3219f25d6305837697a0e5d110975d6be50317c9a2e815823306c3576.
//
// Solidity: event GroupRegistered(bytes32 id)
func (_FileTribeDApp *FileTribeDAppFilterer) FilterGroupRegistered(opts *bind.FilterOpts) (*FileTribeDAppGroupRegisteredIterator, error) {

	logs, sub, err := _FileTribeDApp.contract.FilterLogs(opts, "GroupRegistered")
	if err != nil {
		return nil, err
	}
	return &FileTribeDAppGroupRegisteredIterator{contract: _FileTribeDApp.contract, event: "GroupRegistered", logs: logs, sub: sub}, nil
}

// WatchGroupRegistered is a free log subscription operation binding the contract event 0xb78b34f3219f25d6305837697a0e5d110975d6be50317c9a2e815823306c3576.
//
// Solidity: event GroupRegistered(bytes32 id)
func (_FileTribeDApp *FileTribeDAppFilterer) WatchGroupRegistered(opts *bind.WatchOpts, sink chan<- *FileTribeDAppGroupRegistered) (event.Subscription, error) {

	logs, sub, err := _FileTribeDApp.contract.WatchLogs(opts, "GroupRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FileTribeDAppGroupRegistered)
				if err := _FileTribeDApp.contract.UnpackLog(event, "GroupRegistered", log); err != nil {
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

// ParseGroupRegistered is a log parse operation binding the contract event 0xb78b34f3219f25d6305837697a0e5d110975d6be50317c9a2e815823306c3576.
//
// Solidity: event GroupRegistered(bytes32 id)
func (_FileTribeDApp *FileTribeDAppFilterer) ParseGroupRegistered(log types.Log) (*FileTribeDAppGroupRegistered, error) {
	event := new(FileTribeDAppGroupRegistered)
	if err := _FileTribeDApp.contract.UnpackLog(event, "GroupRegistered", log); err != nil {
		return nil, err
	}
	return event, nil
}

// FileTribeDAppOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the FileTribeDApp contract.
type FileTribeDAppOwnershipTransferredIterator struct {
	Event *FileTribeDAppOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *FileTribeDAppOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FileTribeDAppOwnershipTransferred)
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
		it.Event = new(FileTribeDAppOwnershipTransferred)
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
func (it *FileTribeDAppOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FileTribeDAppOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FileTribeDAppOwnershipTransferred represents a OwnershipTransferred event raised by the FileTribeDApp contract.
type FileTribeDAppOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_FileTribeDApp *FileTribeDAppFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*FileTribeDAppOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _FileTribeDApp.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &FileTribeDAppOwnershipTransferredIterator{contract: _FileTribeDApp.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_FileTribeDApp *FileTribeDAppFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FileTribeDAppOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _FileTribeDApp.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FileTribeDAppOwnershipTransferred)
				if err := _FileTribeDApp.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_FileTribeDApp *FileTribeDAppFilterer) ParseOwnershipTransferred(log types.Log) (*FileTribeDAppOwnershipTransferred, error) {
	event := new(FileTribeDAppOwnershipTransferred)
	if err := _FileTribeDApp.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}
