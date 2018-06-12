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
const EthABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"UserRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"MessageSent\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"boxingKey\",\"type\":\"bytes32\"},{\"name\":\"verifyKey\",\"type\":\"bytes\"},{\"name\":\"ipfsAddr\",\"type\":\"string\"}],\"name\":\"registerUser\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"address\"}],\"name\":\"isUserRegistered\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"address\"}],\"name\":\"getUser\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"},{\"name\":\"\",\"type\":\"bytes32\"},{\"name\":\"\",\"type\":\"bytes\"},{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"sendMessage\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\"},{\"name\":\"from\",\"type\":\"bytes\"},{\"name\":\"to\",\"type\":\"bytes\"},{\"name\":\"signingKey\",\"type\":\"bytes\"},{\"name\":\"check\",\"type\":\"bytes32\"},{\"name\":\"verifyAddress\",\"type\":\"address\"}],\"name\":\"addFriend\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\"},{\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"confirmFriendship\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// EthBin is the compiled bytecode used for deploying new contracts.
const EthBin = `0x608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506110d9806100606000396000f300608060405260043610610083576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680630386c56114610088578063163f7522146101b95780635b981e2b146102145780636f77926b1461031757806382646a58146104ba5780638da5cb5b14610523578063a1dfb5ab1461057a575b600080fd5b34801561009457600080fd5b506101b76004803603810190808035600019169060200190929190803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290803590602001908201803590602001908080601f01602080910402602001604051908101604052809392919081815260200183838082843782019150505050505091929192908035600019169060200190929190803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506105f1565b005b3480156101c557600080fd5b506101fa600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610730565b604051808215151515815260200191505060405180910390f35b34801561022057600080fd5b50610315600480360381019080803590602001908201803590602001908080601f01602080910402602001604051908101604052809392919081815260200183838082843782019150505050505091929192908035600019169060200190929190803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290803590602001908201803590602001908080601f016020809104026020016040519081016040528093929190818152602001838380828437820191505050505050919291929050505061079a565b005b34801561032357600080fd5b50610358600480360381019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506109c0565b604051808060200185600019166000191681526020018060200180602001848103845288818151815260200191508051906020019080838360005b838110156103ae578082015181840152602081019050610393565b50505050905090810190601f1680156103db5780820380516001836020036101000a031916815260200191505b50848103835286818151815260200191508051906020019080838360005b838110156104145780820151818401526020810190506103f9565b50505050905090810190601f1680156104415780820380516001836020036101000a031916815260200191505b50848103825285818151815260200191508051906020019080838360005b8381101561047a57808201518184015260208101905061045f565b50505050905090810190601f1680156104a75780820380516001836020036101000a031916815260200191505b5097505050505050505060405180910390f35b3480156104c657600080fd5b50610521600480360381019080803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290505050610cb7565b005b34801561052f57600080fd5b50610538610d56565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561058657600080fd5b506105ef6004803603810190808035600019169060200190929190803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290505050610d7b565b005b60c060405190810160405280868152602001858152602001848152602001836000191681526020018273ffffffffffffffffffffffffffffffffffffffff168152602001600015158152506002600088600019166000191681526020019081526020016000206000820151816000019080519060200190610673929190610f53565b506020820151816001019080519060200190610690929190610f53565b5060408201518160020190805190602001906106ad929190610f53565b506060820151816003019060001916905560808201518160040160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060a08201518160040160146101000a81548160ff021916908315150217905550905050505050505050565b6000600160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060040160009054906101000a900460ff16156107905760019050610795565b600090505b919050565b600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060040160009054906101000a900460ff1615151561085f576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260178152602001807f557365726e616d6520616c72656164792065786973747300000000000000000081525060200191505060405180910390fd5b60a0604051908101604052808581526020018460001916815260200183815260200182815260200160011515815250600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008201518160000190805190602001906108e9929190610fd3565b50602082015181600101906000191690556040820151816002019080519060200190610916929190610f53565b506060820151816003019080519060200190610933929190610fd3565b5060808201518160040160006101000a81548160ff0219169083151502179055509050507f54db7a5cb4735e1aac1f53db512d3390390bb6637bd30ad4bf9fc98667d9b9b933604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390a150505050565b606060006060806109cf611053565b600160008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060040160009054906101000a900460ff161515610a2a57600080fd5b600160008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060a06040519081016040529081600082018054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610b0e5780601f10610ae357610100808354040283529160200191610b0e565b820191906000526020600020905b815481529060010190602001808311610af157829003601f168201915b50505050508152602001600182015460001916600019168152602001600282018054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610bc25780601f10610b9757610100808354040283529160200191610bc2565b820191906000526020600020905b815481529060010190602001808311610ba557829003601f168201915b50505050508152602001600382018054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610c645780601f10610c3957610100808354040283529160200191610c64565b820191906000526020600020905b815481529060010190602001808311610c4757829003601f168201915b505050505081526020016004820160009054906101000a900460ff161515151581525050905080600001518160200151826040015183606001518393508191508090509450945094509450509193509193565b7f8c5261668696ce22758910d05bab8f186d6eb247ceac2af2e82c7dc17669b036816040518080602001828103825283818151815260200191508051906020019080838360005b83811015610d19578082015181840152602081019050610cfe565b50505050905090810190601f168015610d465780820380516001836020036101000a031916815260200191505b509250505060405180910390a150565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000610da56002600085600019166000191681526020019081526020016000206003015483610e5b565b905060026000846000191660001916815260200190815260200160002060040160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16141515610e1f57600080fd5b600160026000856000191660001916815260200190815260200160002060040160146101000a81548160ff021916908315150217905550505050565b60008060008060418551141515610e755760009350610f4a565b6020850151925060408501519150606085015160001a9050601b8160ff161015610ea057601b810190505b601b8160ff1614158015610eb85750601c8160ff1614155b15610ec65760009350610f4a565b600186828585604051600081526020016040526040518085600019166000191681526020018460ff1660ff1681526020018360001916600019168152602001826000191660001916815260200194505050505060206040516020810390808403906000865af1158015610f3d573d6000803e3d6000fd5b5050506020604051035193505b50505092915050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10610f9457805160ff1916838001178555610fc2565b82800160010185558215610fc2579182015b82811115610fc1578251825591602001919060010190610fa6565b5b509050610fcf9190611088565b5090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061101457805160ff1916838001178555611042565b82800160010185558215611042579182015b82811115611041578251825591602001919060010190611026565b5b50905061104f9190611088565b5090565b60a060405190810160405280606081526020016000801916815260200160608152602001606081526020016000151581525090565b6110aa91905b808211156110a657600081600090555060010161108e565b5090565b905600a165627a7a723058201e0876eca741b4480ac67a5d155a0ee74b046c8b659477f7c713e72f59e2a9360029`

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
// Solidity: function getUser(id address) constant returns(string, bytes32, bytes, string)
func (_Eth *EthCaller) GetUser(opts *bind.CallOpts, id common.Address) (string, [32]byte, []byte, string, error) {
	var (
		ret0 = new(string)
		ret1 = new([32]byte)
		ret2 = new([]byte)
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
// Solidity: function getUser(id address) constant returns(string, bytes32, bytes, string)
func (_Eth *EthSession) GetUser(id common.Address) (string, [32]byte, []byte, string, error) {
	return _Eth.Contract.GetUser(&_Eth.CallOpts, id)
}

// GetUser is a free data retrieval call binding the contract method 0x6f77926b.
//
// Solidity: function getUser(id address) constant returns(string, bytes32, bytes, string)
func (_Eth *EthCallerSession) GetUser(id common.Address) (string, [32]byte, []byte, string, error) {
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

// AddFriend is a paid mutator transaction binding the contract method 0x0386c561.
//
// Solidity: function addFriend(id bytes32, from bytes, to bytes, signingKey bytes, check bytes32, verifyAddress address) returns()
func (_Eth *EthTransactor) AddFriend(opts *bind.TransactOpts, id [32]byte, from []byte, to []byte, signingKey []byte, check [32]byte, verifyAddress common.Address) (*types.Transaction, error) {
	return _Eth.contract.Transact(opts, "addFriend", id, from, to, signingKey, check, verifyAddress)
}

// AddFriend is a paid mutator transaction binding the contract method 0x0386c561.
//
// Solidity: function addFriend(id bytes32, from bytes, to bytes, signingKey bytes, check bytes32, verifyAddress address) returns()
func (_Eth *EthSession) AddFriend(id [32]byte, from []byte, to []byte, signingKey []byte, check [32]byte, verifyAddress common.Address) (*types.Transaction, error) {
	return _Eth.Contract.AddFriend(&_Eth.TransactOpts, id, from, to, signingKey, check, verifyAddress)
}

// AddFriend is a paid mutator transaction binding the contract method 0x0386c561.
//
// Solidity: function addFriend(id bytes32, from bytes, to bytes, signingKey bytes, check bytes32, verifyAddress address) returns()
func (_Eth *EthTransactorSession) AddFriend(id [32]byte, from []byte, to []byte, signingKey []byte, check [32]byte, verifyAddress common.Address) (*types.Transaction, error) {
	return _Eth.Contract.AddFriend(&_Eth.TransactOpts, id, from, to, signingKey, check, verifyAddress)
}

// ConfirmFriendship is a paid mutator transaction binding the contract method 0xa1dfb5ab.
//
// Solidity: function confirmFriendship(id bytes32, signature bytes) returns()
func (_Eth *EthTransactor) ConfirmFriendship(opts *bind.TransactOpts, id [32]byte, signature []byte) (*types.Transaction, error) {
	return _Eth.contract.Transact(opts, "confirmFriendship", id, signature)
}

// ConfirmFriendship is a paid mutator transaction binding the contract method 0xa1dfb5ab.
//
// Solidity: function confirmFriendship(id bytes32, signature bytes) returns()
func (_Eth *EthSession) ConfirmFriendship(id [32]byte, signature []byte) (*types.Transaction, error) {
	return _Eth.Contract.ConfirmFriendship(&_Eth.TransactOpts, id, signature)
}

// ConfirmFriendship is a paid mutator transaction binding the contract method 0xa1dfb5ab.
//
// Solidity: function confirmFriendship(id bytes32, signature bytes) returns()
func (_Eth *EthTransactorSession) ConfirmFriendship(id [32]byte, signature []byte) (*types.Transaction, error) {
	return _Eth.Contract.ConfirmFriendship(&_Eth.TransactOpts, id, signature)
}

// RegisterUser is a paid mutator transaction binding the contract method 0x5b981e2b.
//
// Solidity: function registerUser(name string, boxingKey bytes32, verifyKey bytes, ipfsAddr string) returns()
func (_Eth *EthTransactor) RegisterUser(opts *bind.TransactOpts, name string, boxingKey [32]byte, verifyKey []byte, ipfsAddr string) (*types.Transaction, error) {
	return _Eth.contract.Transact(opts, "registerUser", name, boxingKey, verifyKey, ipfsAddr)
}

// RegisterUser is a paid mutator transaction binding the contract method 0x5b981e2b.
//
// Solidity: function registerUser(name string, boxingKey bytes32, verifyKey bytes, ipfsAddr string) returns()
func (_Eth *EthSession) RegisterUser(name string, boxingKey [32]byte, verifyKey []byte, ipfsAddr string) (*types.Transaction, error) {
	return _Eth.Contract.RegisterUser(&_Eth.TransactOpts, name, boxingKey, verifyKey, ipfsAddr)
}

// RegisterUser is a paid mutator transaction binding the contract method 0x5b981e2b.
//
// Solidity: function registerUser(name string, boxingKey bytes32, verifyKey bytes, ipfsAddr string) returns()
func (_Eth *EthTransactorSession) RegisterUser(name string, boxingKey [32]byte, verifyKey []byte, ipfsAddr string) (*types.Transaction, error) {
	return _Eth.Contract.RegisterUser(&_Eth.TransactOpts, name, boxingKey, verifyKey, ipfsAddr)
}

// SendMessage is a paid mutator transaction binding the contract method 0x82646a58.
//
// Solidity: function sendMessage(message bytes) returns()
func (_Eth *EthTransactor) SendMessage(opts *bind.TransactOpts, message []byte) (*types.Transaction, error) {
	return _Eth.contract.Transact(opts, "sendMessage", message)
}

// SendMessage is a paid mutator transaction binding the contract method 0x82646a58.
//
// Solidity: function sendMessage(message bytes) returns()
func (_Eth *EthSession) SendMessage(message []byte) (*types.Transaction, error) {
	return _Eth.Contract.SendMessage(&_Eth.TransactOpts, message)
}

// SendMessage is a paid mutator transaction binding the contract method 0x82646a58.
//
// Solidity: function sendMessage(message bytes) returns()
func (_Eth *EthTransactorSession) SendMessage(message []byte) (*types.Transaction, error) {
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
	Message []byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMessageSent is a free log retrieval operation binding the contract event 0x8c5261668696ce22758910d05bab8f186d6eb247ceac2af2e82c7dc17669b036.
//
// Solidity: e MessageSent(message bytes)
func (_Eth *EthFilterer) FilterMessageSent(opts *bind.FilterOpts) (*EthMessageSentIterator, error) {

	logs, sub, err := _Eth.contract.FilterLogs(opts, "MessageSent")
	if err != nil {
		return nil, err
	}
	return &EthMessageSentIterator{contract: _Eth.contract, event: "MessageSent", logs: logs, sub: sub}, nil
}

// WatchMessageSent is a free log subscription operation binding the contract event 0x8c5261668696ce22758910d05bab8f186d6eb247ceac2af2e82c7dc17669b036.
//
// Solidity: e MessageSent(message bytes)
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
