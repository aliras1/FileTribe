// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Account

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

// AccountABI is the input ABI used to generate the binding from.
const AccountABI = "[{\"inputs\":[{\"internalType\":\"contractIFileTribeDApp\",\"name\":\"fileTribe\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"ipfsPeerId\",\"type\":\"string\"},{\"internalType\":\"bytes32\",\"name\":\"boxingKey\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"msg\",\"type\":\"string\"}],\"name\":\"Debug\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"msg\",\"type\":\"bytes\"}],\"name\":\"DebugBytesAcc\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractIAccount\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractIGroup\",\"name\":\"group\",\"type\":\"address\"}],\"name\":\"GroupCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractIAccount\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractIGroup\",\"name\":\"group\",\"type\":\"address\"}],\"name\":\"GroupLeft\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractIAccount\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractIGroup\",\"name\":\"group\",\"type\":\"address\"}],\"name\":\"InvitationAccepted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractIAccount\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractIGroup\",\"name\":\"group\",\"type\":\"address\"}],\"name\":\"InvitationDeclined\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractIAccount\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractIGroup\",\"name\":\"group\",\"type\":\"address\"}],\"name\":\"NewInvitation\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"createGroup\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"invite\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"onInvitationAccepted\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"onInvitationDeclined\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"contractIGroup\",\"name\":\"group\",\"type\":\"address\"}],\"name\":\"onGroupLeft\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"groups\",\"outputs\":[{\"internalType\":\"contractIGroup[]\",\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"ipfsId\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// AccountBin is the compiled bytecode used for deploying new contracts.
var AccountBin = "0x60806040523480156200001157600080fd5b5060405162000e6838038062000e68833981810160405260a08110156200003757600080fd5b815160208301516040808501805191519395929483019291846401000000008211156200006357600080fd5b9083019060208201858111156200007957600080fd5b82516401000000008111828201881017156200009457600080fd5b82525081516020918201929091019080838360005b83811015620000c3578181015183820152602001620000a9565b50505050905090810190601f168015620000f15780820380516001836020036101000a031916815260200191505b50604052602001805160405193929190846401000000008211156200011557600080fd5b9083019060208201858111156200012b57600080fd5b82516401000000008111828201881017156200014657600080fd5b82525081516020918201929091019080838360005b83811015620001755781810151838201526020016200015b565b50505050905090810190601f168015620001a35780820380516001836020036101000a031916815260200191505b50604081905260209190910151600080546001600160a01b0319166001600160a01b038981169190911780835592955088945091909116917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908290a350600180546001600160a01b0319166001600160a01b03871617905582516200023190600290602086019062000256565b5081516200024790600390602085019062000256565b5060045550620002fb92505050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106200029957805160ff1916838001178555620002c9565b82800160010185558215620002c9579182015b82811115620002c9578251825591602001919060010190620002ac565b50620002d7929150620002db565b5090565b620002f891905b80821115620002d75760008155600101620002e2565b90565b610b5d806200030b6000396000f3fe608060405234801561001057600080fd5b50600436106100b45760003560e01c8063aff10f2811610071578063aff10f28146101e0578063b061d9a914610206578063dc2ddcae1461020e578063e9896cae146102b4578063eec30bfd146102bc578063f2fde38b146102c4576100b4565b806306fdde03146100b95780635bf89d9e14610136578063715018a61461018e57806383e86440146101985780638da5cb5b146101a05780638f32d59b146101c4575b600080fd5b6100c16102ea565b6040805160208082528351818301528351919283929083019185019080838360005b838110156100fb5781810151838201526020016100e3565b50505050905090810190601f1680156101285780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b61013e61037e565b60408051602080825283518183015283519192839290830191858101910280838360005b8381101561017a578181015183820152602001610162565b505050509050019250505060405180910390f35b6101966103df565b005b61019661043a565b6101a861058d565b604080516001600160a01b039092168252519081900360200190f35b6101cc61059c565b604080519115158252519081900360200190f35b610196600480360360208110156101f657600080fd5b50356001600160a01b03166105ad565b6101966106b5565b6101966004803603602081101561022457600080fd5b81019060208101813564010000000081111561023f57600080fd5b82018360208201111561025157600080fd5b8035906020019184600183028401116401000000008311171561027357600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610804945050505050565b6100c161097a565b6101966109db565b610196600480360360208110156102da57600080fd5b50356001600160a01b0316610a5e565b60028054604080516020601f60001961010060018716150201909416859004938401819004810282018101909252828152606093909290918301828280156103735780601f1061034857610100808354040283529160200191610373565b820191906000526020600020905b81548152906001019060200180831161035657829003601f168201915b505050505090505b90565b6060600580548060200260200160405190810160405280929190818152602001828054801561037357602002820191906000526020600020905b81546001600160a01b031681526001909101906020018083116103b8575050505050905090565b6103e761059c565b6103f057600080fd5b600080546040516001600160a01b03909116907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a3600080546001600160a01b0319169055565b6000335b60065482101561058957806001600160a01b03166006838154811061045f57fe5b6000918252602090912001546001600160a01b0316141561057e57600580546001810182556000919091527f036b6384b5eca791c62761152d0c79bb0604c104a5fb6f4eb0703f3154bb3db00180546001600160a01b0319166001600160a01b0383161790556006805460001981019081106104d757fe5b600091825260209091200154600680546001600160a01b0390921691849081106104fd57fe5b600091825260209091200180546001600160a01b0319166001600160a01b03929092169190911790556006805490610539906000198301610ae1565b50604080513081526001600160a01b038316602082015281517f4d7c243e154e530692e62f8539db65779f5cb85d58831956361697addede5adb929181900390910190a15b60019091019061043e565b5050565b6000546001600160a01b031690565b6000546001600160a01b0316331490565b60005b60055481101561058957816001600160a01b0316600582815481106105d157fe5b6000918252602090912001546001600160a01b031614156106aa576005805460001981019081106105fe57fe5b600091825260209091200154600580546001600160a01b03909216918390811061062457fe5b600091825260209091200180546001600160a01b0319166001600160a01b03929092169190911790556005805490610660906000198301610ae1565b50604080513081526001600160a01b038416602082015281517f0bbec966f892c752826c6f21835e9473e7aa1467ed30f73cb70360894022eed7929181900390910190a1506106b2565b6001016105b0565b50565b6000335b60065482101561058957806001600160a01b0316600683815481106106da57fe5b6000918252602090912001546001600160a01b031614156107f957600580546001810182556000919091527f036b6384b5eca791c62761152d0c79bb0604c104a5fb6f4eb0703f3154bb3db00180546001600160a01b0319166001600160a01b03831617905560068054600019810190811061075257fe5b600091825260209091200154600680546001600160a01b03909216918490811061077857fe5b600091825260209091200180546001600160a01b0319166001600160a01b039290921691909117905560068054906107b4906000198301610ae1565b50604080513081526001600160a01b038316602082015281517f7ea317cacaf634632d88c1fd20864b1f54f609bb063c678280aac6312715edef929181900390910190a15b6001909101906106b9565b61080c61059c565b61081557600080fd5b600154604051636e16ee5760e11b81526020600482018181528451602484015284516000946001600160a01b03169363dc2ddcae938793928392604401918501908083838b5b8381101561087357818101518382015260200161085b565b50505050905090810190601f1680156108a05780820380516001836020036101000a031916815260200191505b5092505050602060405180830381600087803b1580156108bf57600080fd5b505af11580156108d3573d6000803e3d6000fd5b505050506040513d60208110156108e957600080fd5b5051600580546001810182556000919091527f036b6384b5eca791c62761152d0c79bb0604c104a5fb6f4eb0703f3154bb3db00180546001600160a01b0319166001600160a01b03831690811790915560408051308152602081019290925280519293507fe39aa5facc1ab0c6f9442d5cac663064983c41007e4ff54800c9be350c4d2b1092918290030190a15050565b60038054604080516020601f60026000196101006001881615020190951694909404938401819004810282018101909252828152606093909290918301828280156103735780601f1061034857610100808354040283529160200191610373565b600680546001810182556000919091527ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f0180546001600160a01b03191633908117909155604080513081526020810183905281517f0527b775a0c0bf1079aefbf956c30a63bede4d66e9ef3b4a21b7ad4021b5fdaa929181900390910190a150565b610a6661059c565b610a6f57600080fd5b6106b2816001600160a01b038116610a8657600080fd5b600080546040516001600160a01b03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a3600080546001600160a01b0319166001600160a01b0392909216919091179055565b815481835581811115610b0557600083815260209020610b05918101908301610b0a565b505050565b61037b91905b80821115610b245760008155600101610b10565b509056fea265627a7a72315820edfed3d65c45874d3ec493b7308512b0c00f2596606856da64994a89b2f417f264736f6c634300050c0032"

// DeployAccount deploys a new Ethereum contract, binding an instance of Account to it.
func DeployAccount(auth *bind.TransactOpts, backend bind.ContractBackend, fileTribe common.Address, owner common.Address, name string, ipfsPeerId string, boxingKey [32]byte) (common.Address, *types.Transaction, *Account, error) {
	parsed, err := abi.JSON(strings.NewReader(AccountABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(AccountBin), backend, fileTribe, owner, name, ipfsPeerId, boxingKey)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Account{AccountCaller: AccountCaller{contract: contract}, AccountTransactor: AccountTransactor{contract: contract}, AccountFilterer: AccountFilterer{contract: contract}}, nil
}

// Account is an auto generated Go binding around an Ethereum contract.
type Account struct {
	AccountCaller     // Read-only binding to the contract
	AccountTransactor // Write-only binding to the contract
	AccountFilterer   // Log filterer for contract events
}

// AccountCaller is an auto generated read-only Go binding around an Ethereum contract.
type AccountCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccountTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AccountTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccountFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AccountFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccountSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AccountSession struct {
	Contract     *Account          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AccountCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AccountCallerSession struct {
	Contract *AccountCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// AccountTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AccountTransactorSession struct {
	Contract     *AccountTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// AccountRaw is an auto generated low-level Go binding around an Ethereum contract.
type AccountRaw struct {
	Contract *Account // Generic contract binding to access the raw methods on
}

// AccountCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AccountCallerRaw struct {
	Contract *AccountCaller // Generic read-only contract binding to access the raw methods on
}

// AccountTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AccountTransactorRaw struct {
	Contract *AccountTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAccount creates a new instance of Account, bound to a specific deployed contract.
func NewAccount(address common.Address, backend bind.ContractBackend) (*Account, error) {
	contract, err := bindAccount(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Account{AccountCaller: AccountCaller{contract: contract}, AccountTransactor: AccountTransactor{contract: contract}, AccountFilterer: AccountFilterer{contract: contract}}, nil
}

// NewAccountCaller creates a new read-only instance of Account, bound to a specific deployed contract.
func NewAccountCaller(address common.Address, caller bind.ContractCaller) (*AccountCaller, error) {
	contract, err := bindAccount(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AccountCaller{contract: contract}, nil
}

// NewAccountTransactor creates a new write-only instance of Account, bound to a specific deployed contract.
func NewAccountTransactor(address common.Address, transactor bind.ContractTransactor) (*AccountTransactor, error) {
	contract, err := bindAccount(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AccountTransactor{contract: contract}, nil
}

// NewAccountFilterer creates a new log filterer instance of Account, bound to a specific deployed contract.
func NewAccountFilterer(address common.Address, filterer bind.ContractFilterer) (*AccountFilterer, error) {
	contract, err := bindAccount(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AccountFilterer{contract: contract}, nil
}

// bindAccount binds a generic wrapper to an already deployed contract.
func bindAccount(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AccountABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Account *AccountRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Account.Contract.AccountCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Account *AccountRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Account.Contract.AccountTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Account *AccountRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Account.Contract.AccountTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Account *AccountCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Account.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Account *AccountTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Account.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Account *AccountTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Account.Contract.contract.Transact(opts, method, params...)
}

// Groups is a free data retrieval call binding the contract method 0x5bf89d9e.
//
// Solidity: function groups() constant returns(address[])
func (_Account *AccountCaller) Groups(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _Account.contract.Call(opts, out, "groups")
	return *ret0, err
}

// Groups is a free data retrieval call binding the contract method 0x5bf89d9e.
//
// Solidity: function groups() constant returns(address[])
func (_Account *AccountSession) Groups() ([]common.Address, error) {
	return _Account.Contract.Groups(&_Account.CallOpts)
}

// Groups is a free data retrieval call binding the contract method 0x5bf89d9e.
//
// Solidity: function groups() constant returns(address[])
func (_Account *AccountCallerSession) Groups() ([]common.Address, error) {
	return _Account.Contract.Groups(&_Account.CallOpts)
}

// IpfsId is a free data retrieval call binding the contract method 0xe9896cae.
//
// Solidity: function ipfsId() constant returns(string)
func (_Account *AccountCaller) IpfsId(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Account.contract.Call(opts, out, "ipfsId")
	return *ret0, err
}

// IpfsId is a free data retrieval call binding the contract method 0xe9896cae.
//
// Solidity: function ipfsId() constant returns(string)
func (_Account *AccountSession) IpfsId() (string, error) {
	return _Account.Contract.IpfsId(&_Account.CallOpts)
}

// IpfsId is a free data retrieval call binding the contract method 0xe9896cae.
//
// Solidity: function ipfsId() constant returns(string)
func (_Account *AccountCallerSession) IpfsId() (string, error) {
	return _Account.Contract.IpfsId(&_Account.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Account *AccountCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Account.contract.Call(opts, out, "isOwner")
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Account *AccountSession) IsOwner() (bool, error) {
	return _Account.Contract.IsOwner(&_Account.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Account *AccountCallerSession) IsOwner() (bool, error) {
	return _Account.Contract.IsOwner(&_Account.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_Account *AccountCaller) Name(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Account.contract.Call(opts, out, "name")
	return *ret0, err
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_Account *AccountSession) Name() (string, error) {
	return _Account.Contract.Name(&_Account.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_Account *AccountCallerSession) Name() (string, error) {
	return _Account.Contract.Name(&_Account.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Account *AccountCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Account.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Account *AccountSession) Owner() (common.Address, error) {
	return _Account.Contract.Owner(&_Account.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Account *AccountCallerSession) Owner() (common.Address, error) {
	return _Account.Contract.Owner(&_Account.CallOpts)
}

// CreateGroup is a paid mutator transaction binding the contract method 0xdc2ddcae.
//
// Solidity: function createGroup(string name) returns()
func (_Account *AccountTransactor) CreateGroup(opts *bind.TransactOpts, name string) (*types.Transaction, error) {
	return _Account.contract.Transact(opts, "createGroup", name)
}

// CreateGroup is a paid mutator transaction binding the contract method 0xdc2ddcae.
//
// Solidity: function createGroup(string name) returns()
func (_Account *AccountSession) CreateGroup(name string) (*types.Transaction, error) {
	return _Account.Contract.CreateGroup(&_Account.TransactOpts, name)
}

// CreateGroup is a paid mutator transaction binding the contract method 0xdc2ddcae.
//
// Solidity: function createGroup(string name) returns()
func (_Account *AccountTransactorSession) CreateGroup(name string) (*types.Transaction, error) {
	return _Account.Contract.CreateGroup(&_Account.TransactOpts, name)
}

// Invite is a paid mutator transaction binding the contract method 0xeec30bfd.
//
// Solidity: function invite() returns()
func (_Account *AccountTransactor) Invite(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Account.contract.Transact(opts, "invite")
}

// Invite is a paid mutator transaction binding the contract method 0xeec30bfd.
//
// Solidity: function invite() returns()
func (_Account *AccountSession) Invite() (*types.Transaction, error) {
	return _Account.Contract.Invite(&_Account.TransactOpts)
}

// Invite is a paid mutator transaction binding the contract method 0xeec30bfd.
//
// Solidity: function invite() returns()
func (_Account *AccountTransactorSession) Invite() (*types.Transaction, error) {
	return _Account.Contract.Invite(&_Account.TransactOpts)
}

// OnGroupLeft is a paid mutator transaction binding the contract method 0xaff10f28.
//
// Solidity: function onGroupLeft(address group) returns()
func (_Account *AccountTransactor) OnGroupLeft(opts *bind.TransactOpts, group common.Address) (*types.Transaction, error) {
	return _Account.contract.Transact(opts, "onGroupLeft", group)
}

// OnGroupLeft is a paid mutator transaction binding the contract method 0xaff10f28.
//
// Solidity: function onGroupLeft(address group) returns()
func (_Account *AccountSession) OnGroupLeft(group common.Address) (*types.Transaction, error) {
	return _Account.Contract.OnGroupLeft(&_Account.TransactOpts, group)
}

// OnGroupLeft is a paid mutator transaction binding the contract method 0xaff10f28.
//
// Solidity: function onGroupLeft(address group) returns()
func (_Account *AccountTransactorSession) OnGroupLeft(group common.Address) (*types.Transaction, error) {
	return _Account.Contract.OnGroupLeft(&_Account.TransactOpts, group)
}

// OnInvitationAccepted is a paid mutator transaction binding the contract method 0x83e86440.
//
// Solidity: function onInvitationAccepted() returns()
func (_Account *AccountTransactor) OnInvitationAccepted(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Account.contract.Transact(opts, "onInvitationAccepted")
}

// OnInvitationAccepted is a paid mutator transaction binding the contract method 0x83e86440.
//
// Solidity: function onInvitationAccepted() returns()
func (_Account *AccountSession) OnInvitationAccepted() (*types.Transaction, error) {
	return _Account.Contract.OnInvitationAccepted(&_Account.TransactOpts)
}

// OnInvitationAccepted is a paid mutator transaction binding the contract method 0x83e86440.
//
// Solidity: function onInvitationAccepted() returns()
func (_Account *AccountTransactorSession) OnInvitationAccepted() (*types.Transaction, error) {
	return _Account.Contract.OnInvitationAccepted(&_Account.TransactOpts)
}

// OnInvitationDeclined is a paid mutator transaction binding the contract method 0xb061d9a9.
//
// Solidity: function onInvitationDeclined() returns()
func (_Account *AccountTransactor) OnInvitationDeclined(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Account.contract.Transact(opts, "onInvitationDeclined")
}

// OnInvitationDeclined is a paid mutator transaction binding the contract method 0xb061d9a9.
//
// Solidity: function onInvitationDeclined() returns()
func (_Account *AccountSession) OnInvitationDeclined() (*types.Transaction, error) {
	return _Account.Contract.OnInvitationDeclined(&_Account.TransactOpts)
}

// OnInvitationDeclined is a paid mutator transaction binding the contract method 0xb061d9a9.
//
// Solidity: function onInvitationDeclined() returns()
func (_Account *AccountTransactorSession) OnInvitationDeclined() (*types.Transaction, error) {
	return _Account.Contract.OnInvitationDeclined(&_Account.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Account *AccountTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Account.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Account *AccountSession) RenounceOwnership() (*types.Transaction, error) {
	return _Account.Contract.RenounceOwnership(&_Account.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Account *AccountTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Account.Contract.RenounceOwnership(&_Account.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Account *AccountTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Account.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Account *AccountSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Account.Contract.TransferOwnership(&_Account.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Account *AccountTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Account.Contract.TransferOwnership(&_Account.TransactOpts, newOwner)
}

// AccountDebugIterator is returned from FilterDebug and is used to iterate over the raw logs and unpacked data for Debug events raised by the Account contract.
type AccountDebugIterator struct {
	Event *AccountDebug // Event containing the contract specifics and raw log

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
func (it *AccountDebugIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccountDebug)
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
		it.Event = new(AccountDebug)
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
func (it *AccountDebugIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccountDebugIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccountDebug represents a Debug event raised by the Account contract.
type AccountDebug struct {
	Msg string
	Raw types.Log // Blockchain specific contextual infos
}

// FilterDebug is a free log retrieval operation binding the contract event 0x7cdb51e9dbbc205231228146c3246e7f914aa6d4a33170e43ecc8e3593481d1a.
//
// Solidity: event Debug(string msg)
func (_Account *AccountFilterer) FilterDebug(opts *bind.FilterOpts) (*AccountDebugIterator, error) {

	logs, sub, err := _Account.contract.FilterLogs(opts, "Debug")
	if err != nil {
		return nil, err
	}
	return &AccountDebugIterator{contract: _Account.contract, event: "Debug", logs: logs, sub: sub}, nil
}

// WatchDebug is a free log subscription operation binding the contract event 0x7cdb51e9dbbc205231228146c3246e7f914aa6d4a33170e43ecc8e3593481d1a.
//
// Solidity: event Debug(string msg)
func (_Account *AccountFilterer) WatchDebug(opts *bind.WatchOpts, sink chan<- *AccountDebug) (event.Subscription, error) {

	logs, sub, err := _Account.contract.WatchLogs(opts, "Debug")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccountDebug)
				if err := _Account.contract.UnpackLog(event, "Debug", log); err != nil {
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

// ParseDebug is a log parse operation binding the contract event 0x7cdb51e9dbbc205231228146c3246e7f914aa6d4a33170e43ecc8e3593481d1a.
//
// Solidity: event Debug(string msg)
func (_Account *AccountFilterer) ParseDebug(log types.Log) (*AccountDebug, error) {
	event := new(AccountDebug)
	if err := _Account.contract.UnpackLog(event, "Debug", log); err != nil {
		return nil, err
	}
	return event, nil
}

// AccountDebugBytesAccIterator is returned from FilterDebugBytesAcc and is used to iterate over the raw logs and unpacked data for DebugBytesAcc events raised by the Account contract.
type AccountDebugBytesAccIterator struct {
	Event *AccountDebugBytesAcc // Event containing the contract specifics and raw log

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
func (it *AccountDebugBytesAccIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccountDebugBytesAcc)
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
		it.Event = new(AccountDebugBytesAcc)
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
func (it *AccountDebugBytesAccIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccountDebugBytesAccIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccountDebugBytesAcc represents a DebugBytesAcc event raised by the Account contract.
type AccountDebugBytesAcc struct {
	Msg []byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterDebugBytesAcc is a free log retrieval operation binding the contract event 0x77e5005705a2864b0490bd315067ca7cf9a4637a493dde671083d6c87851ba66.
//
// Solidity: event DebugBytesAcc(bytes msg)
func (_Account *AccountFilterer) FilterDebugBytesAcc(opts *bind.FilterOpts) (*AccountDebugBytesAccIterator, error) {

	logs, sub, err := _Account.contract.FilterLogs(opts, "DebugBytesAcc")
	if err != nil {
		return nil, err
	}
	return &AccountDebugBytesAccIterator{contract: _Account.contract, event: "DebugBytesAcc", logs: logs, sub: sub}, nil
}

// WatchDebugBytesAcc is a free log subscription operation binding the contract event 0x77e5005705a2864b0490bd315067ca7cf9a4637a493dde671083d6c87851ba66.
//
// Solidity: event DebugBytesAcc(bytes msg)
func (_Account *AccountFilterer) WatchDebugBytesAcc(opts *bind.WatchOpts, sink chan<- *AccountDebugBytesAcc) (event.Subscription, error) {

	logs, sub, err := _Account.contract.WatchLogs(opts, "DebugBytesAcc")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccountDebugBytesAcc)
				if err := _Account.contract.UnpackLog(event, "DebugBytesAcc", log); err != nil {
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

// ParseDebugBytesAcc is a log parse operation binding the contract event 0x77e5005705a2864b0490bd315067ca7cf9a4637a493dde671083d6c87851ba66.
//
// Solidity: event DebugBytesAcc(bytes msg)
func (_Account *AccountFilterer) ParseDebugBytesAcc(log types.Log) (*AccountDebugBytesAcc, error) {
	event := new(AccountDebugBytesAcc)
	if err := _Account.contract.UnpackLog(event, "DebugBytesAcc", log); err != nil {
		return nil, err
	}
	return event, nil
}

// AccountGroupCreatedIterator is returned from FilterGroupCreated and is used to iterate over the raw logs and unpacked data for GroupCreated events raised by the Account contract.
type AccountGroupCreatedIterator struct {
	Event *AccountGroupCreated // Event containing the contract specifics and raw log

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
func (it *AccountGroupCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccountGroupCreated)
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
		it.Event = new(AccountGroupCreated)
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
func (it *AccountGroupCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccountGroupCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccountGroupCreated represents a GroupCreated event raised by the Account contract.
type AccountGroupCreated struct {
	Account common.Address
	Group   common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterGroupCreated is a free log retrieval operation binding the contract event 0xe39aa5facc1ab0c6f9442d5cac663064983c41007e4ff54800c9be350c4d2b10.
//
// Solidity: event GroupCreated(address account, address group)
func (_Account *AccountFilterer) FilterGroupCreated(opts *bind.FilterOpts) (*AccountGroupCreatedIterator, error) {

	logs, sub, err := _Account.contract.FilterLogs(opts, "GroupCreated")
	if err != nil {
		return nil, err
	}
	return &AccountGroupCreatedIterator{contract: _Account.contract, event: "GroupCreated", logs: logs, sub: sub}, nil
}

// WatchGroupCreated is a free log subscription operation binding the contract event 0xe39aa5facc1ab0c6f9442d5cac663064983c41007e4ff54800c9be350c4d2b10.
//
// Solidity: event GroupCreated(address account, address group)
func (_Account *AccountFilterer) WatchGroupCreated(opts *bind.WatchOpts, sink chan<- *AccountGroupCreated) (event.Subscription, error) {

	logs, sub, err := _Account.contract.WatchLogs(opts, "GroupCreated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccountGroupCreated)
				if err := _Account.contract.UnpackLog(event, "GroupCreated", log); err != nil {
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

// ParseGroupCreated is a log parse operation binding the contract event 0xe39aa5facc1ab0c6f9442d5cac663064983c41007e4ff54800c9be350c4d2b10.
//
// Solidity: event GroupCreated(address account, address group)
func (_Account *AccountFilterer) ParseGroupCreated(log types.Log) (*AccountGroupCreated, error) {
	event := new(AccountGroupCreated)
	if err := _Account.contract.UnpackLog(event, "GroupCreated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// AccountGroupLeftIterator is returned from FilterGroupLeft and is used to iterate over the raw logs and unpacked data for GroupLeft events raised by the Account contract.
type AccountGroupLeftIterator struct {
	Event *AccountGroupLeft // Event containing the contract specifics and raw log

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
func (it *AccountGroupLeftIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccountGroupLeft)
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
		it.Event = new(AccountGroupLeft)
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
func (it *AccountGroupLeftIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccountGroupLeftIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccountGroupLeft represents a GroupLeft event raised by the Account contract.
type AccountGroupLeft struct {
	Account common.Address
	Group   common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterGroupLeft is a free log retrieval operation binding the contract event 0x0bbec966f892c752826c6f21835e9473e7aa1467ed30f73cb70360894022eed7.
//
// Solidity: event GroupLeft(address account, address group)
func (_Account *AccountFilterer) FilterGroupLeft(opts *bind.FilterOpts) (*AccountGroupLeftIterator, error) {

	logs, sub, err := _Account.contract.FilterLogs(opts, "GroupLeft")
	if err != nil {
		return nil, err
	}
	return &AccountGroupLeftIterator{contract: _Account.contract, event: "GroupLeft", logs: logs, sub: sub}, nil
}

// WatchGroupLeft is a free log subscription operation binding the contract event 0x0bbec966f892c752826c6f21835e9473e7aa1467ed30f73cb70360894022eed7.
//
// Solidity: event GroupLeft(address account, address group)
func (_Account *AccountFilterer) WatchGroupLeft(opts *bind.WatchOpts, sink chan<- *AccountGroupLeft) (event.Subscription, error) {

	logs, sub, err := _Account.contract.WatchLogs(opts, "GroupLeft")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccountGroupLeft)
				if err := _Account.contract.UnpackLog(event, "GroupLeft", log); err != nil {
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

// ParseGroupLeft is a log parse operation binding the contract event 0x0bbec966f892c752826c6f21835e9473e7aa1467ed30f73cb70360894022eed7.
//
// Solidity: event GroupLeft(address account, address group)
func (_Account *AccountFilterer) ParseGroupLeft(log types.Log) (*AccountGroupLeft, error) {
	event := new(AccountGroupLeft)
	if err := _Account.contract.UnpackLog(event, "GroupLeft", log); err != nil {
		return nil, err
	}
	return event, nil
}

// AccountInvitationAcceptedIterator is returned from FilterInvitationAccepted and is used to iterate over the raw logs and unpacked data for InvitationAccepted events raised by the Account contract.
type AccountInvitationAcceptedIterator struct {
	Event *AccountInvitationAccepted // Event containing the contract specifics and raw log

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
func (it *AccountInvitationAcceptedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccountInvitationAccepted)
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
		it.Event = new(AccountInvitationAccepted)
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
func (it *AccountInvitationAcceptedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccountInvitationAcceptedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccountInvitationAccepted represents a InvitationAccepted event raised by the Account contract.
type AccountInvitationAccepted struct {
	Account common.Address
	Group   common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInvitationAccepted is a free log retrieval operation binding the contract event 0x4d7c243e154e530692e62f8539db65779f5cb85d58831956361697addede5adb.
//
// Solidity: event InvitationAccepted(address account, address group)
func (_Account *AccountFilterer) FilterInvitationAccepted(opts *bind.FilterOpts) (*AccountInvitationAcceptedIterator, error) {

	logs, sub, err := _Account.contract.FilterLogs(opts, "InvitationAccepted")
	if err != nil {
		return nil, err
	}
	return &AccountInvitationAcceptedIterator{contract: _Account.contract, event: "InvitationAccepted", logs: logs, sub: sub}, nil
}

// WatchInvitationAccepted is a free log subscription operation binding the contract event 0x4d7c243e154e530692e62f8539db65779f5cb85d58831956361697addede5adb.
//
// Solidity: event InvitationAccepted(address account, address group)
func (_Account *AccountFilterer) WatchInvitationAccepted(opts *bind.WatchOpts, sink chan<- *AccountInvitationAccepted) (event.Subscription, error) {

	logs, sub, err := _Account.contract.WatchLogs(opts, "InvitationAccepted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccountInvitationAccepted)
				if err := _Account.contract.UnpackLog(event, "InvitationAccepted", log); err != nil {
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

// ParseInvitationAccepted is a log parse operation binding the contract event 0x4d7c243e154e530692e62f8539db65779f5cb85d58831956361697addede5adb.
//
// Solidity: event InvitationAccepted(address account, address group)
func (_Account *AccountFilterer) ParseInvitationAccepted(log types.Log) (*AccountInvitationAccepted, error) {
	event := new(AccountInvitationAccepted)
	if err := _Account.contract.UnpackLog(event, "InvitationAccepted", log); err != nil {
		return nil, err
	}
	return event, nil
}

// AccountInvitationDeclinedIterator is returned from FilterInvitationDeclined and is used to iterate over the raw logs and unpacked data for InvitationDeclined events raised by the Account contract.
type AccountInvitationDeclinedIterator struct {
	Event *AccountInvitationDeclined // Event containing the contract specifics and raw log

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
func (it *AccountInvitationDeclinedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccountInvitationDeclined)
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
		it.Event = new(AccountInvitationDeclined)
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
func (it *AccountInvitationDeclinedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccountInvitationDeclinedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccountInvitationDeclined represents a InvitationDeclined event raised by the Account contract.
type AccountInvitationDeclined struct {
	Account common.Address
	Group   common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInvitationDeclined is a free log retrieval operation binding the contract event 0x7ea317cacaf634632d88c1fd20864b1f54f609bb063c678280aac6312715edef.
//
// Solidity: event InvitationDeclined(address account, address group)
func (_Account *AccountFilterer) FilterInvitationDeclined(opts *bind.FilterOpts) (*AccountInvitationDeclinedIterator, error) {

	logs, sub, err := _Account.contract.FilterLogs(opts, "InvitationDeclined")
	if err != nil {
		return nil, err
	}
	return &AccountInvitationDeclinedIterator{contract: _Account.contract, event: "InvitationDeclined", logs: logs, sub: sub}, nil
}

// WatchInvitationDeclined is a free log subscription operation binding the contract event 0x7ea317cacaf634632d88c1fd20864b1f54f609bb063c678280aac6312715edef.
//
// Solidity: event InvitationDeclined(address account, address group)
func (_Account *AccountFilterer) WatchInvitationDeclined(opts *bind.WatchOpts, sink chan<- *AccountInvitationDeclined) (event.Subscription, error) {

	logs, sub, err := _Account.contract.WatchLogs(opts, "InvitationDeclined")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccountInvitationDeclined)
				if err := _Account.contract.UnpackLog(event, "InvitationDeclined", log); err != nil {
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

// ParseInvitationDeclined is a log parse operation binding the contract event 0x7ea317cacaf634632d88c1fd20864b1f54f609bb063c678280aac6312715edef.
//
// Solidity: event InvitationDeclined(address account, address group)
func (_Account *AccountFilterer) ParseInvitationDeclined(log types.Log) (*AccountInvitationDeclined, error) {
	event := new(AccountInvitationDeclined)
	if err := _Account.contract.UnpackLog(event, "InvitationDeclined", log); err != nil {
		return nil, err
	}
	return event, nil
}

// AccountNewInvitationIterator is returned from FilterNewInvitation and is used to iterate over the raw logs and unpacked data for NewInvitation events raised by the Account contract.
type AccountNewInvitationIterator struct {
	Event *AccountNewInvitation // Event containing the contract specifics and raw log

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
func (it *AccountNewInvitationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccountNewInvitation)
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
		it.Event = new(AccountNewInvitation)
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
func (it *AccountNewInvitationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccountNewInvitationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccountNewInvitation represents a NewInvitation event raised by the Account contract.
type AccountNewInvitation struct {
	Account common.Address
	Group   common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterNewInvitation is a free log retrieval operation binding the contract event 0x0527b775a0c0bf1079aefbf956c30a63bede4d66e9ef3b4a21b7ad4021b5fdaa.
//
// Solidity: event NewInvitation(address account, address group)
func (_Account *AccountFilterer) FilterNewInvitation(opts *bind.FilterOpts) (*AccountNewInvitationIterator, error) {

	logs, sub, err := _Account.contract.FilterLogs(opts, "NewInvitation")
	if err != nil {
		return nil, err
	}
	return &AccountNewInvitationIterator{contract: _Account.contract, event: "NewInvitation", logs: logs, sub: sub}, nil
}

// WatchNewInvitation is a free log subscription operation binding the contract event 0x0527b775a0c0bf1079aefbf956c30a63bede4d66e9ef3b4a21b7ad4021b5fdaa.
//
// Solidity: event NewInvitation(address account, address group)
func (_Account *AccountFilterer) WatchNewInvitation(opts *bind.WatchOpts, sink chan<- *AccountNewInvitation) (event.Subscription, error) {

	logs, sub, err := _Account.contract.WatchLogs(opts, "NewInvitation")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccountNewInvitation)
				if err := _Account.contract.UnpackLog(event, "NewInvitation", log); err != nil {
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

// ParseNewInvitation is a log parse operation binding the contract event 0x0527b775a0c0bf1079aefbf956c30a63bede4d66e9ef3b4a21b7ad4021b5fdaa.
//
// Solidity: event NewInvitation(address account, address group)
func (_Account *AccountFilterer) ParseNewInvitation(log types.Log) (*AccountNewInvitation, error) {
	event := new(AccountNewInvitation)
	if err := _Account.contract.UnpackLog(event, "NewInvitation", log); err != nil {
		return nil, err
	}
	return event, nil
}

// AccountOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Account contract.
type AccountOwnershipTransferredIterator struct {
	Event *AccountOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *AccountOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccountOwnershipTransferred)
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
		it.Event = new(AccountOwnershipTransferred)
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
func (it *AccountOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccountOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccountOwnershipTransferred represents a OwnershipTransferred event raised by the Account contract.
type AccountOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Account *AccountFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*AccountOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Account.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &AccountOwnershipTransferredIterator{contract: _Account.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Account *AccountFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AccountOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Account.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccountOwnershipTransferred)
				if err := _Account.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Account *AccountFilterer) ParseOwnershipTransferred(log types.Log) (*AccountOwnershipTransferred, error) {
	event := new(AccountOwnershipTransferred)
	if err := _Account.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}
