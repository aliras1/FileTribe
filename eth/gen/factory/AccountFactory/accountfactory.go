// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package AccountFactory

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

// AccountFactoryABI is the input ABI used to generate the binding from.
const AccountFactoryABI = "[{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"ipfsId\",\"type\":\"string\"},{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"}],\"name\":\"create\",\"outputs\":[{\"internalType\":\"contractIAccount\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"contractIFileTribeDApp\",\"name\":\"parent\",\"type\":\"address\"}],\"name\":\"setParent\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// AccountFactoryBin is the compiled bytecode used for deploying new contracts.
var AccountFactoryBin = "0x608060405234801561001057600080fd5b50600080546001600160a01b0319163390811780835560405191926001600160a01b0391909116917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908290a35061127f8061006d6000396000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c80631499c59214610067578063715018a61461008f5780638da5cb5b146100975780638f32d59b146100bb578063d0bf552b146100d7578063f2fde38b146101a9575b600080fd5b61008d6004803603602081101561007d57600080fd5b50356001600160a01b03166101cf565b005b61008d610202565b61009f61025d565b604080516001600160a01b039092168252519081900360200190f35b6100c361026c565b604080519115158252519081900360200190f35b61009f600480360360808110156100ed57600080fd5b6001600160a01b03823516919081019060408101602082013564010000000081111561011857600080fd5b82018360208201111561012a57600080fd5b8035906020019184600183028401116401000000008311171561014c57600080fd5b91939092909160208101903564010000000081111561016a57600080fd5b82018360208201111561017c57600080fd5b8035906020019184600183028401116401000000008311171561019e57600080fd5b91935091503561027d565b61008d600480360360208110156101bf57600080fd5b50356001600160a01b031661034a565b6101d761026c565b6101e057600080fd5b600180546001600160a01b0319166001600160a01b0392909216919091179055565b61020a61026c565b61021357600080fd5b600080546040516001600160a01b03909116907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a3600080546001600160a01b0319169055565b6000546001600160a01b031690565b6000546001600160a01b0316331490565b6000600160009054906101000a90046001600160a01b03168787878787876040516102a7906103d5565b6001600160a01b038089168252871660208201526080810182905260a0604082018181529082018690526060820160c08301888880828437600083820152601f01601f1916909101848103835286815260200190508686808284376000838201819052604051601f909201601f19169093018190039c509a509098505050505050505050f08015801561033e573d6000803e3d6000fd5b50979650505050505050565b61035261026c565b61035b57600080fd5b61036481610367565b50565b6001600160a01b03811661037a57600080fd5b600080546040516001600160a01b03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a3600080546001600160a01b0319166001600160a01b0392909216919091179055565b610e68806103e38339019056fe60806040523480156200001157600080fd5b5060405162000e6838038062000e68833981810160405260a08110156200003757600080fd5b815160208301516040808501805191519395929483019291846401000000008211156200006357600080fd5b9083019060208201858111156200007957600080fd5b82516401000000008111828201881017156200009457600080fd5b82525081516020918201929091019080838360005b83811015620000c3578181015183820152602001620000a9565b50505050905090810190601f168015620000f15780820380516001836020036101000a031916815260200191505b50604052602001805160405193929190846401000000008211156200011557600080fd5b9083019060208201858111156200012b57600080fd5b82516401000000008111828201881017156200014657600080fd5b82525081516020918201929091019080838360005b83811015620001755781810151838201526020016200015b565b50505050905090810190601f168015620001a35780820380516001836020036101000a031916815260200191505b50604081905260209190910151600080546001600160a01b0319166001600160a01b038981169190911780835592955088945091909116917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908290a350600180546001600160a01b0319166001600160a01b03871617905582516200023190600290602086019062000256565b5081516200024790600390602085019062000256565b5060045550620002fb92505050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106200029957805160ff1916838001178555620002c9565b82800160010185558215620002c9579182015b82811115620002c9578251825591602001919060010190620002ac565b50620002d7929150620002db565b5090565b620002f891905b80821115620002d75760008155600101620002e2565b90565b610b5d806200030b6000396000f3fe608060405234801561001057600080fd5b50600436106100b45760003560e01c8063aff10f2811610071578063aff10f28146101e0578063b061d9a914610206578063dc2ddcae1461020e578063e9896cae146102b4578063eec30bfd146102bc578063f2fde38b146102c4576100b4565b806306fdde03146100b95780635bf89d9e14610136578063715018a61461018e57806383e86440146101985780638da5cb5b146101a05780638f32d59b146101c4575b600080fd5b6100c16102ea565b6040805160208082528351818301528351919283929083019185019080838360005b838110156100fb5781810151838201526020016100e3565b50505050905090810190601f1680156101285780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b61013e61037e565b60408051602080825283518183015283519192839290830191858101910280838360005b8381101561017a578181015183820152602001610162565b505050509050019250505060405180910390f35b6101966103df565b005b61019661043a565b6101a861058d565b604080516001600160a01b039092168252519081900360200190f35b6101cc61059c565b604080519115158252519081900360200190f35b610196600480360360208110156101f657600080fd5b50356001600160a01b03166105ad565b6101966106b5565b6101966004803603602081101561022457600080fd5b81019060208101813564010000000081111561023f57600080fd5b82018360208201111561025157600080fd5b8035906020019184600183028401116401000000008311171561027357600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610804945050505050565b6100c161097a565b6101966109db565b610196600480360360208110156102da57600080fd5b50356001600160a01b0316610a5e565b60028054604080516020601f60001961010060018716150201909416859004938401819004810282018101909252828152606093909290918301828280156103735780601f1061034857610100808354040283529160200191610373565b820191906000526020600020905b81548152906001019060200180831161035657829003601f168201915b505050505090505b90565b6060600580548060200260200160405190810160405280929190818152602001828054801561037357602002820191906000526020600020905b81546001600160a01b031681526001909101906020018083116103b8575050505050905090565b6103e761059c565b6103f057600080fd5b600080546040516001600160a01b03909116907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a3600080546001600160a01b0319169055565b6000335b60065482101561058957806001600160a01b03166006838154811061045f57fe5b6000918252602090912001546001600160a01b0316141561057e57600580546001810182556000919091527f036b6384b5eca791c62761152d0c79bb0604c104a5fb6f4eb0703f3154bb3db00180546001600160a01b0319166001600160a01b0383161790556006805460001981019081106104d757fe5b600091825260209091200154600680546001600160a01b0390921691849081106104fd57fe5b600091825260209091200180546001600160a01b0319166001600160a01b03929092169190911790556006805490610539906000198301610ae1565b50604080513081526001600160a01b038316602082015281517f4d7c243e154e530692e62f8539db65779f5cb85d58831956361697addede5adb929181900390910190a15b60019091019061043e565b5050565b6000546001600160a01b031690565b6000546001600160a01b0316331490565b60005b60055481101561058957816001600160a01b0316600582815481106105d157fe5b6000918252602090912001546001600160a01b031614156106aa576005805460001981019081106105fe57fe5b600091825260209091200154600580546001600160a01b03909216918390811061062457fe5b600091825260209091200180546001600160a01b0319166001600160a01b03929092169190911790556005805490610660906000198301610ae1565b50604080513081526001600160a01b038416602082015281517f0bbec966f892c752826c6f21835e9473e7aa1467ed30f73cb70360894022eed7929181900390910190a1506106b2565b6001016105b0565b50565b6000335b60065482101561058957806001600160a01b0316600683815481106106da57fe5b6000918252602090912001546001600160a01b031614156107f957600580546001810182556000919091527f036b6384b5eca791c62761152d0c79bb0604c104a5fb6f4eb0703f3154bb3db00180546001600160a01b0319166001600160a01b03831617905560068054600019810190811061075257fe5b600091825260209091200154600680546001600160a01b03909216918490811061077857fe5b600091825260209091200180546001600160a01b0319166001600160a01b039290921691909117905560068054906107b4906000198301610ae1565b50604080513081526001600160a01b038316602082015281517f7ea317cacaf634632d88c1fd20864b1f54f609bb063c678280aac6312715edef929181900390910190a15b6001909101906106b9565b61080c61059c565b61081557600080fd5b600154604051636e16ee5760e11b81526020600482018181528451602484015284516000946001600160a01b03169363dc2ddcae938793928392604401918501908083838b5b8381101561087357818101518382015260200161085b565b50505050905090810190601f1680156108a05780820380516001836020036101000a031916815260200191505b5092505050602060405180830381600087803b1580156108bf57600080fd5b505af11580156108d3573d6000803e3d6000fd5b505050506040513d60208110156108e957600080fd5b5051600580546001810182556000919091527f036b6384b5eca791c62761152d0c79bb0604c104a5fb6f4eb0703f3154bb3db00180546001600160a01b0319166001600160a01b03831690811790915560408051308152602081019290925280519293507fe39aa5facc1ab0c6f9442d5cac663064983c41007e4ff54800c9be350c4d2b1092918290030190a15050565b60038054604080516020601f60026000196101006001881615020190951694909404938401819004810282018101909252828152606093909290918301828280156103735780601f1061034857610100808354040283529160200191610373565b600680546001810182556000919091527ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f0180546001600160a01b03191633908117909155604080513081526020810183905281517f0527b775a0c0bf1079aefbf956c30a63bede4d66e9ef3b4a21b7ad4021b5fdaa929181900390910190a150565b610a6661059c565b610a6f57600080fd5b6106b2816001600160a01b038116610a8657600080fd5b600080546040516001600160a01b03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a3600080546001600160a01b0319166001600160a01b0392909216919091179055565b815481835581811115610b0557600083815260209020610b05918101908301610b0a565b505050565b61037b91905b80821115610b245760008155600101610b10565b509056fea265627a7a723158209633fe07705e010b9f291d2c5f31f308baa3fd80c2f61595f2220a63fddcc81d64736f6c634300050f0032a265627a7a72315820dee6524689e997ce6fba53a61a9d086060ae83289ed470eb059c8a3d2edd6da564736f6c634300050f0032"

// DeployAccountFactory deploys a new Ethereum contract, binding an instance of AccountFactory to it.
func DeployAccountFactory(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *AccountFactory, error) {
	parsed, err := abi.JSON(strings.NewReader(AccountFactoryABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(AccountFactoryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AccountFactory{AccountFactoryCaller: AccountFactoryCaller{contract: contract}, AccountFactoryTransactor: AccountFactoryTransactor{contract: contract}, AccountFactoryFilterer: AccountFactoryFilterer{contract: contract}}, nil
}

// AccountFactory is an auto generated Go binding around an Ethereum contract.
type AccountFactory struct {
	AccountFactoryCaller     // Read-only binding to the contract
	AccountFactoryTransactor // Write-only binding to the contract
	AccountFactoryFilterer   // Log filterer for contract events
}

// AccountFactoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type AccountFactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccountFactoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AccountFactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccountFactoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AccountFactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccountFactorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AccountFactorySession struct {
	Contract     *AccountFactory   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AccountFactoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AccountFactoryCallerSession struct {
	Contract *AccountFactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// AccountFactoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AccountFactoryTransactorSession struct {
	Contract     *AccountFactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// AccountFactoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type AccountFactoryRaw struct {
	Contract *AccountFactory // Generic contract binding to access the raw methods on
}

// AccountFactoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AccountFactoryCallerRaw struct {
	Contract *AccountFactoryCaller // Generic read-only contract binding to access the raw methods on
}

// AccountFactoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AccountFactoryTransactorRaw struct {
	Contract *AccountFactoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAccountFactory creates a new instance of AccountFactory, bound to a specific deployed contract.
func NewAccountFactory(address common.Address, backend bind.ContractBackend) (*AccountFactory, error) {
	contract, err := bindAccountFactory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AccountFactory{AccountFactoryCaller: AccountFactoryCaller{contract: contract}, AccountFactoryTransactor: AccountFactoryTransactor{contract: contract}, AccountFactoryFilterer: AccountFactoryFilterer{contract: contract}}, nil
}

// NewAccountFactoryCaller creates a new read-only instance of AccountFactory, bound to a specific deployed contract.
func NewAccountFactoryCaller(address common.Address, caller bind.ContractCaller) (*AccountFactoryCaller, error) {
	contract, err := bindAccountFactory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AccountFactoryCaller{contract: contract}, nil
}

// NewAccountFactoryTransactor creates a new write-only instance of AccountFactory, bound to a specific deployed contract.
func NewAccountFactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*AccountFactoryTransactor, error) {
	contract, err := bindAccountFactory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AccountFactoryTransactor{contract: contract}, nil
}

// NewAccountFactoryFilterer creates a new log filterer instance of AccountFactory, bound to a specific deployed contract.
func NewAccountFactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*AccountFactoryFilterer, error) {
	contract, err := bindAccountFactory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AccountFactoryFilterer{contract: contract}, nil
}

// bindAccountFactory binds a generic wrapper to an already deployed contract.
func bindAccountFactory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AccountFactoryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AccountFactory *AccountFactoryRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AccountFactory.Contract.AccountFactoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AccountFactory *AccountFactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccountFactory.Contract.AccountFactoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AccountFactory *AccountFactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccountFactory.Contract.AccountFactoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AccountFactory *AccountFactoryCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AccountFactory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AccountFactory *AccountFactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccountFactory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AccountFactory *AccountFactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccountFactory.Contract.contract.Transact(opts, method, params...)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_AccountFactory *AccountFactoryCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _AccountFactory.contract.Call(opts, out, "isOwner")
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_AccountFactory *AccountFactorySession) IsOwner() (bool, error) {
	return _AccountFactory.Contract.IsOwner(&_AccountFactory.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_AccountFactory *AccountFactoryCallerSession) IsOwner() (bool, error) {
	return _AccountFactory.Contract.IsOwner(&_AccountFactory.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_AccountFactory *AccountFactoryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _AccountFactory.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_AccountFactory *AccountFactorySession) Owner() (common.Address, error) {
	return _AccountFactory.Contract.Owner(&_AccountFactory.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_AccountFactory *AccountFactoryCallerSession) Owner() (common.Address, error) {
	return _AccountFactory.Contract.Owner(&_AccountFactory.CallOpts)
}

// Create is a paid mutator transaction binding the contract method 0xd0bf552b.
//
// Solidity: function create(address owner, string name, string ipfsId, bytes32 key) returns(address)
func (_AccountFactory *AccountFactoryTransactor) Create(opts *bind.TransactOpts, owner common.Address, name string, ipfsId string, key [32]byte) (*types.Transaction, error) {
	return _AccountFactory.contract.Transact(opts, "create", owner, name, ipfsId, key)
}

// Create is a paid mutator transaction binding the contract method 0xd0bf552b.
//
// Solidity: function create(address owner, string name, string ipfsId, bytes32 key) returns(address)
func (_AccountFactory *AccountFactorySession) Create(owner common.Address, name string, ipfsId string, key [32]byte) (*types.Transaction, error) {
	return _AccountFactory.Contract.Create(&_AccountFactory.TransactOpts, owner, name, ipfsId, key)
}

// Create is a paid mutator transaction binding the contract method 0xd0bf552b.
//
// Solidity: function create(address owner, string name, string ipfsId, bytes32 key) returns(address)
func (_AccountFactory *AccountFactoryTransactorSession) Create(owner common.Address, name string, ipfsId string, key [32]byte) (*types.Transaction, error) {
	return _AccountFactory.Contract.Create(&_AccountFactory.TransactOpts, owner, name, ipfsId, key)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_AccountFactory *AccountFactoryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccountFactory.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_AccountFactory *AccountFactorySession) RenounceOwnership() (*types.Transaction, error) {
	return _AccountFactory.Contract.RenounceOwnership(&_AccountFactory.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_AccountFactory *AccountFactoryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _AccountFactory.Contract.RenounceOwnership(&_AccountFactory.TransactOpts)
}

// SetParent is a paid mutator transaction binding the contract method 0x1499c592.
//
// Solidity: function setParent(address parent) returns()
func (_AccountFactory *AccountFactoryTransactor) SetParent(opts *bind.TransactOpts, parent common.Address) (*types.Transaction, error) {
	return _AccountFactory.contract.Transact(opts, "setParent", parent)
}

// SetParent is a paid mutator transaction binding the contract method 0x1499c592.
//
// Solidity: function setParent(address parent) returns()
func (_AccountFactory *AccountFactorySession) SetParent(parent common.Address) (*types.Transaction, error) {
	return _AccountFactory.Contract.SetParent(&_AccountFactory.TransactOpts, parent)
}

// SetParent is a paid mutator transaction binding the contract method 0x1499c592.
//
// Solidity: function setParent(address parent) returns()
func (_AccountFactory *AccountFactoryTransactorSession) SetParent(parent common.Address) (*types.Transaction, error) {
	return _AccountFactory.Contract.SetParent(&_AccountFactory.TransactOpts, parent)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_AccountFactory *AccountFactoryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _AccountFactory.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_AccountFactory *AccountFactorySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _AccountFactory.Contract.TransferOwnership(&_AccountFactory.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_AccountFactory *AccountFactoryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _AccountFactory.Contract.TransferOwnership(&_AccountFactory.TransactOpts, newOwner)
}

// AccountFactoryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the AccountFactory contract.
type AccountFactoryOwnershipTransferredIterator struct {
	Event *AccountFactoryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *AccountFactoryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccountFactoryOwnershipTransferred)
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
		it.Event = new(AccountFactoryOwnershipTransferred)
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
func (it *AccountFactoryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccountFactoryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccountFactoryOwnershipTransferred represents a OwnershipTransferred event raised by the AccountFactory contract.
type AccountFactoryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_AccountFactory *AccountFactoryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*AccountFactoryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _AccountFactory.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &AccountFactoryOwnershipTransferredIterator{contract: _AccountFactory.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_AccountFactory *AccountFactoryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AccountFactoryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _AccountFactory.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccountFactoryOwnershipTransferred)
				if err := _AccountFactory.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_AccountFactory *AccountFactoryFilterer) ParseOwnershipTransferred(log types.Log) (*AccountFactoryOwnershipTransferred, error) {
	event := new(AccountFactoryOwnershipTransferred)
	if err := _AccountFactory.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}
