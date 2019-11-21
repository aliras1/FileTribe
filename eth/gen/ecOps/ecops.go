// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ecOps

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

// EcOpsABI is the input ABI used to generate the binding from.
const EcOpsABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"b\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"p\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"q\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"p0\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"scalar\",\"type\":\"uint256\"}],\"name\":\"ecmul\",\"outputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"p1\",\"type\":\"uint256[2]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"p0\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"p1\",\"type\":\"uint256[2]\"}],\"name\":\"ecadd\",\"outputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"p2\",\"type\":\"uint256[2]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"base\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"e\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"m\",\"type\":\"uint256\"}],\"name\":\"modExp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"r\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"g1XToYSquared\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"y\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"ySqr\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"q\",\"type\":\"uint256\"}],\"name\":\"calcQuadRes\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"result\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"hashToG1\",\"outputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"point\",\"type\":\"uint256[2]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"x\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[4]\",\"name\":\"w\",\"type\":\"uint256[4]\"},{\"internalType\":\"uint256[2]\",\"name\":\"y\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[4]\",\"name\":\"z\",\"type\":\"uint256[4]\"}],\"name\":\"pairingCheck\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"p1\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"p2\",\"type\":\"uint256[2]\"}],\"name\":\"isEqualPoints\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isEqual\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"p1\",\"type\":\"uint256[2]\"}],\"name\":\"isInG1\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256[4]\",\"name\":\"p1\",\"type\":\"uint256[4]\"}],\"name\":\"isInG2\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// EcOpsBin is the compiled bytecode used for deploying new contracts.
var EcOpsBin = "0x610b83610026600b82828239805160001a60731461001957fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600436106100d95760003560e01c80639ae8886a11610096578063d8522ba511610070578063d8522ba514610391578063e9d1c41f146103e3578063f171866214610457578063fd3ab282146104fd576100d9565b80639ae8886a146102ca578063a03980bb146102d2578063a68889061461031d576100d9565b80633148f14f146100de57806349ac97d1146101195780634df7e3d0146101fa5780636594e60b1461020257806372dcadb314610225578063776710d714610242575b600080fd5b610107600480360360608110156100f457600080fd5b5080359060208101359060400135610505565b60408051918252519081900360200190f35b6101e6600480360361018081101561013057600080fd5b6040805180820182529183019291818301918390600290839083908082843760009201919091525050604080516080818101909252929594938181019392509060049083908390808284376000920191909152505060408051808201825292959493818101939250906002908390839080828437600092019190915250506040805160808181019092529295949381810193925090600490839083908082843760009201919091525091945061054d9350505050565b604080519115158252519081900360200190f35b610107610603565b6101076004803603604081101561021857600080fd5b5080359060200135610608565b6101076004803603602081101561023b57600080fd5b503561063a565b61028f6004803603606081101561025857600080fd5b60408051808201825291830192918183019183906002908390839080828437600092019190915250919450509035915061065f9050565b6040518082600260200280838360005b838110156102b757818101518382015260200161029f565b5050505090500191505060405180910390f35b610107610698565b6101e6600480360360408110156102e857600080fd5b604080518082018252918301929181830191839060029083908390808284376000920191909152509194506106aa9350505050565b6101e66004803603608081101561033357600080fd5b6040805180820182529183019291818301918390600290839083908082843760009201919091525050604080518082018252929594938181019392509060029083908390808284376000920191909152509194506107389350505050565b6101e6600480360360808110156103a757600080fd5b810190808060800190600480602002604051908101604052809291908260046020028082843760009201919091525091945061075a9350505050565b61028f600480360360808110156103f957600080fd5b6040805180820182529183019291818301918390600290839083908082843760009201919091525050604080518082018252929594938181019392509060029083908390808284376000920191909152509194506108469350505050565b61028f6004803603602081101561046d57600080fd5b81019060208101813564010000000081111561048857600080fd5b82018360208201111561049a57600080fd5b803590602001918460018302840111640100000000831117156104bc57600080fd5b91908080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525092955061088b945050505050565b610107610a73565b600060405160208152602080820152602060408201528460608201528360808201528260a082015260208160c08360056107d05a03fa61054457600080fd5b51949350505050565b6000610557610a97565b5060408051610180810182528651815260208088015181830152865182840152868101516060808401919091528784015160808401528088015160a0840152865160c084015286820151600080516020610b2f8339815191520360e084015285516101008401529085015161012083015291840151610140820152908301516101608201526105e4610ab6565b6000602082610180856008600019fa5050516001149695505050505050565b600381565b60006003808316908114156106335760046002198401046001810161062e868287610505565b935050505b5092915050565b6000610656826003600080516020610b2f833981519152610505565b60030192915050565b610667610ad4565b61066f610af2565b83518152602080850151908201526040808201849052826060836007600019fa61063357600080fd5b600080516020610b2f83398151915281565b80516000901580156106be57506020820151155b156106cb57506001610733565b8151600090600080516020610b2f8339815191529080098351909150600080516020610b2f8339815191529082099050600080516020610b2f833981519152600382086020840151909150600090600080516020610b2f833981519152908009919091149150505b919050565b80518251600091148015610753575060208083015190840151145b9392505050565b6000610764610a97565b6040805161018081018252600181526002602080830191909152855182840152850151606082015290840151608082015260a0810184600360200201518152602001600181526020016002600080516020610b2f833981519152038152602001846000600481106107d157fe5b60200201518152602001846001600481106107e857fe5b60200201518152602001846002600481106107ff57fe5b602002015181526020018460036004811061081657fe5b602002015190529050610827610ab6565b6001602082610180856008600019fa61083e575060005b949350505050565b61084e610ad4565b610856610b10565b83518152602080850151828201528351604080840191909152908401516060830152826080836006600019fa61063357600080fd5b610893610ad4565b60008080805b60008487604051602001808360ff1660ff1660f81b815260010182805190602001908083835b602083106108de5780518252601f1990920191602091820191016108bf565b6001836020036101000a0380198251168184511680821785525050505050509050019250505060405160208183030381529060405280519060200120905084806001019550508060001c9350600080516020610b2f833981519152848161094157fe5b069350600061094f8561063a565b9050600061096b82600080516020610b2f833981519152610608565b90506000610989826002600080516020610b2f833981519152610505565b905082811415610a5b5781955060ff97506002888b604051602001808360ff1660ff1660f81b815260010182805190602001908083835b602083106109df5780518252601f1990920191602091820191016109c0565b6001836020036101000a0380198251168184511680821785525050505050509050019250505060405160208183030381529060405280519060200120601f60208110610a2757fe5b1a81610a2f57fe5b0660ff1694508460011415610a525785600080516020610b2f8339815191520395505b50505050610a64565b50505050610899565b50908352602083015250919050565b7f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f000000181565b604051806101800160405280600c906020820280388339509192915050565b60405180602001604052806001906020820280388339509192915050565b60405180604001604052806002906020820280388339509192915050565b60405180606001604052806003906020820280388339509192915050565b6040518060800160405280600490602082028038833950919291505056fe30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47a265627a7a72315820f767bccc5a8ddb760965b1fead4562deb4f025af215fe4360dc083a51405b28164736f6c634300050c0032"

// DeployEcOps deploys a new Ethereum contract, binding an instance of EcOps to it.
func DeployEcOps(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *EcOps, error) {
	parsed, err := abi.JSON(strings.NewReader(EcOpsABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(EcOpsBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &EcOps{EcOpsCaller: EcOpsCaller{contract: contract}, EcOpsTransactor: EcOpsTransactor{contract: contract}, EcOpsFilterer: EcOpsFilterer{contract: contract}}, nil
}

// EcOps is an auto generated Go binding around an Ethereum contract.
type EcOps struct {
	EcOpsCaller     // Read-only binding to the contract
	EcOpsTransactor // Write-only binding to the contract
	EcOpsFilterer   // Log filterer for contract events
}

// EcOpsCaller is an auto generated read-only Go binding around an Ethereum contract.
type EcOpsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EcOpsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EcOpsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EcOpsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EcOpsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EcOpsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EcOpsSession struct {
	Contract     *EcOps            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EcOpsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EcOpsCallerSession struct {
	Contract *EcOpsCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// EcOpsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EcOpsTransactorSession struct {
	Contract     *EcOpsTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EcOpsRaw is an auto generated low-level Go binding around an Ethereum contract.
type EcOpsRaw struct {
	Contract *EcOps // Generic contract binding to access the raw methods on
}

// EcOpsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EcOpsCallerRaw struct {
	Contract *EcOpsCaller // Generic read-only contract binding to access the raw methods on
}

// EcOpsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EcOpsTransactorRaw struct {
	Contract *EcOpsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEcOps creates a new instance of EcOps, bound to a specific deployed contract.
func NewEcOps(address common.Address, backend bind.ContractBackend) (*EcOps, error) {
	contract, err := bindEcOps(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EcOps{EcOpsCaller: EcOpsCaller{contract: contract}, EcOpsTransactor: EcOpsTransactor{contract: contract}, EcOpsFilterer: EcOpsFilterer{contract: contract}}, nil
}

// NewEcOpsCaller creates a new read-only instance of EcOps, bound to a specific deployed contract.
func NewEcOpsCaller(address common.Address, caller bind.ContractCaller) (*EcOpsCaller, error) {
	contract, err := bindEcOps(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EcOpsCaller{contract: contract}, nil
}

// NewEcOpsTransactor creates a new write-only instance of EcOps, bound to a specific deployed contract.
func NewEcOpsTransactor(address common.Address, transactor bind.ContractTransactor) (*EcOpsTransactor, error) {
	contract, err := bindEcOps(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EcOpsTransactor{contract: contract}, nil
}

// NewEcOpsFilterer creates a new log filterer instance of EcOps, bound to a specific deployed contract.
func NewEcOpsFilterer(address common.Address, filterer bind.ContractFilterer) (*EcOpsFilterer, error) {
	contract, err := bindEcOps(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EcOpsFilterer{contract: contract}, nil
}

// bindEcOps binds a generic wrapper to an already deployed contract.
func bindEcOps(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(EcOpsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EcOps *EcOpsRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _EcOps.Contract.EcOpsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EcOps *EcOpsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EcOps.Contract.EcOpsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EcOps *EcOpsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EcOps.Contract.EcOpsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EcOps *EcOpsCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _EcOps.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EcOps *EcOpsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EcOps.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EcOps *EcOpsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EcOps.Contract.contract.Transact(opts, method, params...)
}

// B is a free data retrieval call binding the contract method 0x4df7e3d0.
//
// Solidity: function b() constant returns(uint256)
func (_EcOps *EcOpsCaller) B(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _EcOps.contract.Call(opts, out, "b")
	return *ret0, err
}

// B is a free data retrieval call binding the contract method 0x4df7e3d0.
//
// Solidity: function b() constant returns(uint256)
func (_EcOps *EcOpsSession) B() (*big.Int, error) {
	return _EcOps.Contract.B(&_EcOps.CallOpts)
}

// B is a free data retrieval call binding the contract method 0x4df7e3d0.
//
// Solidity: function b() constant returns(uint256)
func (_EcOps *EcOpsCallerSession) B() (*big.Int, error) {
	return _EcOps.Contract.B(&_EcOps.CallOpts)
}

// CalcQuadRes is a free data retrieval call binding the contract method 0x6594e60b.
//
// Solidity: function calcQuadRes(uint256 ySqr, uint256 q) constant returns(uint256 result)
func (_EcOps *EcOpsCaller) CalcQuadRes(opts *bind.CallOpts, ySqr *big.Int, q *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _EcOps.contract.Call(opts, out, "calcQuadRes", ySqr, q)
	return *ret0, err
}

// CalcQuadRes is a free data retrieval call binding the contract method 0x6594e60b.
//
// Solidity: function calcQuadRes(uint256 ySqr, uint256 q) constant returns(uint256 result)
func (_EcOps *EcOpsSession) CalcQuadRes(ySqr *big.Int, q *big.Int) (*big.Int, error) {
	return _EcOps.Contract.CalcQuadRes(&_EcOps.CallOpts, ySqr, q)
}

// CalcQuadRes is a free data retrieval call binding the contract method 0x6594e60b.
//
// Solidity: function calcQuadRes(uint256 ySqr, uint256 q) constant returns(uint256 result)
func (_EcOps *EcOpsCallerSession) CalcQuadRes(ySqr *big.Int, q *big.Int) (*big.Int, error) {
	return _EcOps.Contract.CalcQuadRes(&_EcOps.CallOpts, ySqr, q)
}

// Ecadd is a free data retrieval call binding the contract method 0xe9d1c41f.
//
// Solidity: function ecadd(uint256[2] p0, uint256[2] p1) constant returns(uint256[2] p2)
func (_EcOps *EcOpsCaller) Ecadd(opts *bind.CallOpts, p0 [2]*big.Int, p1 [2]*big.Int) ([2]*big.Int, error) {
	var (
		ret0 = new([2]*big.Int)
	)
	out := ret0
	err := _EcOps.contract.Call(opts, out, "ecadd", p0, p1)
	return *ret0, err
}

// Ecadd is a free data retrieval call binding the contract method 0xe9d1c41f.
//
// Solidity: function ecadd(uint256[2] p0, uint256[2] p1) constant returns(uint256[2] p2)
func (_EcOps *EcOpsSession) Ecadd(p0 [2]*big.Int, p1 [2]*big.Int) ([2]*big.Int, error) {
	return _EcOps.Contract.Ecadd(&_EcOps.CallOpts, p0, p1)
}

// Ecadd is a free data retrieval call binding the contract method 0xe9d1c41f.
//
// Solidity: function ecadd(uint256[2] p0, uint256[2] p1) constant returns(uint256[2] p2)
func (_EcOps *EcOpsCallerSession) Ecadd(p0 [2]*big.Int, p1 [2]*big.Int) ([2]*big.Int, error) {
	return _EcOps.Contract.Ecadd(&_EcOps.CallOpts, p0, p1)
}

// Ecmul is a free data retrieval call binding the contract method 0x776710d7.
//
// Solidity: function ecmul(uint256[2] p0, uint256 scalar) constant returns(uint256[2] p1)
func (_EcOps *EcOpsCaller) Ecmul(opts *bind.CallOpts, p0 [2]*big.Int, scalar *big.Int) ([2]*big.Int, error) {
	var (
		ret0 = new([2]*big.Int)
	)
	out := ret0
	err := _EcOps.contract.Call(opts, out, "ecmul", p0, scalar)
	return *ret0, err
}

// Ecmul is a free data retrieval call binding the contract method 0x776710d7.
//
// Solidity: function ecmul(uint256[2] p0, uint256 scalar) constant returns(uint256[2] p1)
func (_EcOps *EcOpsSession) Ecmul(p0 [2]*big.Int, scalar *big.Int) ([2]*big.Int, error) {
	return _EcOps.Contract.Ecmul(&_EcOps.CallOpts, p0, scalar)
}

// Ecmul is a free data retrieval call binding the contract method 0x776710d7.
//
// Solidity: function ecmul(uint256[2] p0, uint256 scalar) constant returns(uint256[2] p1)
func (_EcOps *EcOpsCallerSession) Ecmul(p0 [2]*big.Int, scalar *big.Int) ([2]*big.Int, error) {
	return _EcOps.Contract.Ecmul(&_EcOps.CallOpts, p0, scalar)
}

// G1XToYSquared is a free data retrieval call binding the contract method 0x72dcadb3.
//
// Solidity: function g1XToYSquared(uint256 x) constant returns(uint256 y)
func (_EcOps *EcOpsCaller) G1XToYSquared(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _EcOps.contract.Call(opts, out, "g1XToYSquared", x)
	return *ret0, err
}

// G1XToYSquared is a free data retrieval call binding the contract method 0x72dcadb3.
//
// Solidity: function g1XToYSquared(uint256 x) constant returns(uint256 y)
func (_EcOps *EcOpsSession) G1XToYSquared(x *big.Int) (*big.Int, error) {
	return _EcOps.Contract.G1XToYSquared(&_EcOps.CallOpts, x)
}

// G1XToYSquared is a free data retrieval call binding the contract method 0x72dcadb3.
//
// Solidity: function g1XToYSquared(uint256 x) constant returns(uint256 y)
func (_EcOps *EcOpsCallerSession) G1XToYSquared(x *big.Int) (*big.Int, error) {
	return _EcOps.Contract.G1XToYSquared(&_EcOps.CallOpts, x)
}

// HashToG1 is a free data retrieval call binding the contract method 0xf1718662.
//
// Solidity: function hashToG1(bytes message) constant returns(uint256[2] point)
func (_EcOps *EcOpsCaller) HashToG1(opts *bind.CallOpts, message []byte) ([2]*big.Int, error) {
	var (
		ret0 = new([2]*big.Int)
	)
	out := ret0
	err := _EcOps.contract.Call(opts, out, "hashToG1", message)
	return *ret0, err
}

// HashToG1 is a free data retrieval call binding the contract method 0xf1718662.
//
// Solidity: function hashToG1(bytes message) constant returns(uint256[2] point)
func (_EcOps *EcOpsSession) HashToG1(message []byte) ([2]*big.Int, error) {
	return _EcOps.Contract.HashToG1(&_EcOps.CallOpts, message)
}

// HashToG1 is a free data retrieval call binding the contract method 0xf1718662.
//
// Solidity: function hashToG1(bytes message) constant returns(uint256[2] point)
func (_EcOps *EcOpsCallerSession) HashToG1(message []byte) ([2]*big.Int, error) {
	return _EcOps.Contract.HashToG1(&_EcOps.CallOpts, message)
}

// IsEqualPoints is a free data retrieval call binding the contract method 0xa6888906.
//
// Solidity: function isEqualPoints(uint256[2] p1, uint256[2] p2) constant returns(bool isEqual)
func (_EcOps *EcOpsCaller) IsEqualPoints(opts *bind.CallOpts, p1 [2]*big.Int, p2 [2]*big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _EcOps.contract.Call(opts, out, "isEqualPoints", p1, p2)
	return *ret0, err
}

// IsEqualPoints is a free data retrieval call binding the contract method 0xa6888906.
//
// Solidity: function isEqualPoints(uint256[2] p1, uint256[2] p2) constant returns(bool isEqual)
func (_EcOps *EcOpsSession) IsEqualPoints(p1 [2]*big.Int, p2 [2]*big.Int) (bool, error) {
	return _EcOps.Contract.IsEqualPoints(&_EcOps.CallOpts, p1, p2)
}

// IsEqualPoints is a free data retrieval call binding the contract method 0xa6888906.
//
// Solidity: function isEqualPoints(uint256[2] p1, uint256[2] p2) constant returns(bool isEqual)
func (_EcOps *EcOpsCallerSession) IsEqualPoints(p1 [2]*big.Int, p2 [2]*big.Int) (bool, error) {
	return _EcOps.Contract.IsEqualPoints(&_EcOps.CallOpts, p1, p2)
}

// IsInG1 is a free data retrieval call binding the contract method 0xa03980bb.
//
// Solidity: function isInG1(uint256[2] p1) constant returns(bool)
func (_EcOps *EcOpsCaller) IsInG1(opts *bind.CallOpts, p1 [2]*big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _EcOps.contract.Call(opts, out, "isInG1", p1)
	return *ret0, err
}

// IsInG1 is a free data retrieval call binding the contract method 0xa03980bb.
//
// Solidity: function isInG1(uint256[2] p1) constant returns(bool)
func (_EcOps *EcOpsSession) IsInG1(p1 [2]*big.Int) (bool, error) {
	return _EcOps.Contract.IsInG1(&_EcOps.CallOpts, p1)
}

// IsInG1 is a free data retrieval call binding the contract method 0xa03980bb.
//
// Solidity: function isInG1(uint256[2] p1) constant returns(bool)
func (_EcOps *EcOpsCallerSession) IsInG1(p1 [2]*big.Int) (bool, error) {
	return _EcOps.Contract.IsInG1(&_EcOps.CallOpts, p1)
}

// IsInG2 is a free data retrieval call binding the contract method 0xd8522ba5.
//
// Solidity: function isInG2(uint256[4] p1) constant returns(bool)
func (_EcOps *EcOpsCaller) IsInG2(opts *bind.CallOpts, p1 [4]*big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _EcOps.contract.Call(opts, out, "isInG2", p1)
	return *ret0, err
}

// IsInG2 is a free data retrieval call binding the contract method 0xd8522ba5.
//
// Solidity: function isInG2(uint256[4] p1) constant returns(bool)
func (_EcOps *EcOpsSession) IsInG2(p1 [4]*big.Int) (bool, error) {
	return _EcOps.Contract.IsInG2(&_EcOps.CallOpts, p1)
}

// IsInG2 is a free data retrieval call binding the contract method 0xd8522ba5.
//
// Solidity: function isInG2(uint256[4] p1) constant returns(bool)
func (_EcOps *EcOpsCallerSession) IsInG2(p1 [4]*big.Int) (bool, error) {
	return _EcOps.Contract.IsInG2(&_EcOps.CallOpts, p1)
}

// ModExp is a free data retrieval call binding the contract method 0x3148f14f.
//
// Solidity: function modExp(uint256 base, uint256 e, uint256 m) constant returns(uint256 r)
func (_EcOps *EcOpsCaller) ModExp(opts *bind.CallOpts, base *big.Int, e *big.Int, m *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _EcOps.contract.Call(opts, out, "modExp", base, e, m)
	return *ret0, err
}

// ModExp is a free data retrieval call binding the contract method 0x3148f14f.
//
// Solidity: function modExp(uint256 base, uint256 e, uint256 m) constant returns(uint256 r)
func (_EcOps *EcOpsSession) ModExp(base *big.Int, e *big.Int, m *big.Int) (*big.Int, error) {
	return _EcOps.Contract.ModExp(&_EcOps.CallOpts, base, e, m)
}

// ModExp is a free data retrieval call binding the contract method 0x3148f14f.
//
// Solidity: function modExp(uint256 base, uint256 e, uint256 m) constant returns(uint256 r)
func (_EcOps *EcOpsCallerSession) ModExp(base *big.Int, e *big.Int, m *big.Int) (*big.Int, error) {
	return _EcOps.Contract.ModExp(&_EcOps.CallOpts, base, e, m)
}

// P is a free data retrieval call binding the contract method 0x9ae8886a.
//
// Solidity: function p() constant returns(uint256)
func (_EcOps *EcOpsCaller) P(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _EcOps.contract.Call(opts, out, "p")
	return *ret0, err
}

// P is a free data retrieval call binding the contract method 0x9ae8886a.
//
// Solidity: function p() constant returns(uint256)
func (_EcOps *EcOpsSession) P() (*big.Int, error) {
	return _EcOps.Contract.P(&_EcOps.CallOpts)
}

// P is a free data retrieval call binding the contract method 0x9ae8886a.
//
// Solidity: function p() constant returns(uint256)
func (_EcOps *EcOpsCallerSession) P() (*big.Int, error) {
	return _EcOps.Contract.P(&_EcOps.CallOpts)
}

// PairingCheck is a free data retrieval call binding the contract method 0x49ac97d1.
//
// Solidity: function pairingCheck(uint256[2] x, uint256[4] w, uint256[2] y, uint256[4] z) constant returns(bool)
func (_EcOps *EcOpsCaller) PairingCheck(opts *bind.CallOpts, x [2]*big.Int, w [4]*big.Int, y [2]*big.Int, z [4]*big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _EcOps.contract.Call(opts, out, "pairingCheck", x, w, y, z)
	return *ret0, err
}

// PairingCheck is a free data retrieval call binding the contract method 0x49ac97d1.
//
// Solidity: function pairingCheck(uint256[2] x, uint256[4] w, uint256[2] y, uint256[4] z) constant returns(bool)
func (_EcOps *EcOpsSession) PairingCheck(x [2]*big.Int, w [4]*big.Int, y [2]*big.Int, z [4]*big.Int) (bool, error) {
	return _EcOps.Contract.PairingCheck(&_EcOps.CallOpts, x, w, y, z)
}

// PairingCheck is a free data retrieval call binding the contract method 0x49ac97d1.
//
// Solidity: function pairingCheck(uint256[2] x, uint256[4] w, uint256[2] y, uint256[4] z) constant returns(bool)
func (_EcOps *EcOpsCallerSession) PairingCheck(x [2]*big.Int, w [4]*big.Int, y [2]*big.Int, z [4]*big.Int) (bool, error) {
	return _EcOps.Contract.PairingCheck(&_EcOps.CallOpts, x, w, y, z)
}

// Q is a free data retrieval call binding the contract method 0xfd3ab282.
//
// Solidity: function q() constant returns(uint256)
func (_EcOps *EcOpsCaller) Q(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _EcOps.contract.Call(opts, out, "q")
	return *ret0, err
}

// Q is a free data retrieval call binding the contract method 0xfd3ab282.
//
// Solidity: function q() constant returns(uint256)
func (_EcOps *EcOpsSession) Q() (*big.Int, error) {
	return _EcOps.Contract.Q(&_EcOps.CallOpts)
}

// Q is a free data retrieval call binding the contract method 0xfd3ab282.
//
// Solidity: function q() constant returns(uint256)
func (_EcOps *EcOpsCallerSession) Q() (*big.Int, error) {
	return _EcOps.Contract.Q(&_EcOps.CallOpts)
}
