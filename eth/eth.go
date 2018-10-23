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
const EthABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"UserRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"GroupRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"groupId\",\"type\":\"bytes32\"}],\"name\":\"GroupInvitation\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"groupId\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"ipfsPath\",\"type\":\"string\"}],\"name\":\"GroupUpdateIpfsPath\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"msg\",\"type\":\"bytes\"}],\"name\":\"Debug\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"ipfsPeerId\",\"type\":\"string\"},{\"name\":\"boxingKey\",\"type\":\"bytes32\"}],\"name\":\"registerUser\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"address\"}],\"name\":\"isUserRegistered\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"address\"}],\"name\":\"getUser\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"},{\"name\":\"\",\"type\":\"string\"},{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\"},{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"ipfsPath\",\"type\":\"string\"}],\"name\":\"createGroup\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"groupId\",\"type\":\"bytes32\"}],\"name\":\"getGroup\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"},{\"name\":\"\",\"type\":\"address[]\"},{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"groupId\",\"type\":\"bytes32\"},{\"name\":\"newMember\",\"type\":\"address\"}],\"name\":\"inviteUser\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"groupId\",\"type\":\"bytes32\"},{\"name\":\"newIpfsPath\",\"type\":\"string\"},{\"name\":\"members\",\"type\":\"address[]\"},{\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"name\":\"vs\",\"type\":\"uint8[]\"}],\"name\":\"updateGroupIpfsPath\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"self\",\"type\":\"address[]\"}],\"name\":\"heapSort\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// EthBin is the compiled bytecode used for deploying new contracts.
const EthBin = `0x608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550612435806100606000396000f300608060405260043610610099576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063163f75221461009e57806368f6d82a146100f95780636f77926b146101b65780637b799d87146102ed5780637e21f4791461033e5780638da5cb5b146104c1578063abab02b514610518578063b567d4ba1461057e578063e89d7c9d146106dc575b600080fd5b3480156100aa57600080fd5b506100df600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610799565b604051808215151515815260200191505060405180910390f35b34801561010557600080fd5b506101b4600480360381019080803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290803590602001908201803590602001908080601f016020809104026020016040519081016040528093929190818152602001838380828437820191505050505050919291929080356000191690602001909291905050506107f2565b005b3480156101c257600080fd5b506101f7600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610a72565b6040518080602001806020018460001916600019168152602001838103835286818151815260200191508051906020019080838360005b8381101561024957808201518184015260208101905061022e565b50505050905090810190601f1680156102765780820380516001836020036101000a031916815260200191505b50838103825285818151815260200191508051906020019080838360005b838110156102af578082015181840152602081019050610294565b50505050905090810190601f1680156102dc5780820380516001836020036101000a031916815260200191505b509550505050505060405180910390f35b3480156102f957600080fd5b5061033c6004803603810190808035600019169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610d80565b005b34801561034a57600080fd5b506104bf6004803603810190808035600019169060200190929190803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290803590602001908201803590602001908080602002602001604051908101604052809392919081815260200183836020028082843782019150505050505091929192908035906020019082018035906020019080806020026020016040519081016040528093929190818152602001838360200280828437820191505050505050919291929080359060200190820180359060200190808060200260200160405190810160405280939291908181526020018383602002808284378201915050505050509192919290803590602001908201803590602001908080602002602001604051908101604052809392919081815260200183836020028082843782019150505050505091929192905050506110c7565b005b3480156104cd57600080fd5b506104d661167c565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561052457600080fd5b5061057c600480360381019080803590602001908201803590602001908080602002602001604051908101604052809392919081815260200183836020028082843782019150505050505091929192905050506116a1565b005b34801561058a57600080fd5b506105ad6004803603810190808035600019169060200190929190505050611b9d565b60405180806020018060200180602001848103845287818151815260200191508051906020019080838360005b838110156105f55780820151818401526020810190506105da565b50505050905090810190601f1680156106225780820380516001836020036101000a031916815260200191505b50848103835286818151815260200191508051906020019060200280838360005b8381101561065e578082015181840152602081019050610643565b50505050905001848103825285818151815260200191508051906020019080838360005b8381101561069d578082015181840152602081019050610682565b50505050905090810190601f1680156106ca5780820380516001836020036101000a031916815260200191505b50965050505050505060405180910390f35b3480156106e857600080fd5b506107976004803603810190808035600019169060200190929190803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290505050611e69565b005b6000600160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060040160009054906101000a900460ff169050919050565b600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060040160009054906101000a900460ff161515156108b7576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260178152602001807f557365726e616d6520616c72656164792065786973747300000000000000000081525060200191505060405180910390fd5b82600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600001908051906020019061090d92919061232f565b5081600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600101908051906020019061096492919061232f565b5080600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600201816000191690555060018060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060040160006101000a81548160ff0219169083151502179055507f54db7a5cb4735e1aac1f53db512d3390390bb6637bd30ad4bf9fc98667d9b9b933604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390a1505050565b6060806000610a7f6123af565b600160008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060040160009054906101000a900460ff161515610b43576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260138152602001807f5573657220646f6573206e6f742065786973740000000000000000000000000081525060200191505060405180910390fd5b600160008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060a06040519081016040529081600082018054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610c275780601f10610bfc57610100808354040283529160200191610c27565b820191906000526020600020905b815481529060010190602001808311610c0a57829003601f168201915b50505050508152602001600182018054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610cc95780601f10610c9e57610100808354040283529160200191610cc9565b820191906000526020600020905b815481529060010190602001808311610cac57829003601f168201915b5050505050815260200160028201546000191660001916815260200160038201805480602002602001604051908101604052809291908181526020018280548015610d3757602002820191906000526020600020905b81546000191681526020019060010190808311610d1f575b505050505081526020016004820160009054906101000a900460ff1615151515815250509050806000015181602001518260400151829250819150935093509350509193909250565b6001151560026000846000191660001916815260200190815260200160002060040160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff161515141515610e64576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260138152602001807f557365722063616e206e6f7420696e766974650000000000000000000000000081525060200191505060405180910390fd5b600160008273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060040160009054906101000a900460ff161515610f28576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f43616e206e6f7420696e76697465206e6f6e206578697374656e74207573657281525060200191505060405180910390fd5b6002600083600019166000191681526020019081526020016000206002018190806001815401808255809150509060018203906000526020600020016000909192909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050600160008273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206003018290806001815401808255809150509060018203906000526020600020016000909192909190915090600019169055507f9478e2f0a42543d96af3b3661efc5aaa23dd42c9f8c970c1e4f4bd01ab42374a338284604051808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018260001916600019168152602001935050505060405180910390a15050565b60008060026000896000191660001916815260200190815260200160002060050160009054906101000a900460ff16151561116a576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f67726f757020646f6573206e6f7420657869737400000000000000000000000081525060200191505060405180910390fd5b8551855114151561117a57600080fd5b8551845114151561118a57600080fd5b8551835114151561119a57600080fd5b60028060008a60001916600019168152602001908152602001600020600201805490508115156111c657fe5b0486511115156111d557600080fd5b60026000896000191660001916815260200190815260200160002060030187604051808380546001816001161561010002031660029004801561124f5780601f1061122d57610100808354040283529182019161124f565b820191906000526020600020905b81548152906001019060200180831161123b575b505082805190602001908083835b602083101515611282578051825260208201915060208101905060208303925061125d565b6001836020036101000a0380198251168184511680821785525050505050509050019250505060405180910390209150600090505b855181101561148b576112e18887838151811015156112d257fe5b90602001906020020151612164565b151561137b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602c8152602001807f696e76616c696420617070726f76616c3a2075736572206973206e6f7420612081526020017f67726f7570206d656d626572000000000000000000000000000000000000000081525060400191505060405180910390fd5b6113e4868281518110151561138c57fe5b906020019060200201518385848151811015156113a557fe5b9060200190602002015188858151811015156113bd57fe5b9060200190602002015188868151811015156113d557fe5b90602001906020020151612240565b151561147e576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260238152602001807f696e76616c696420617070726f76616c3a20696e76616c6964207369676e617481526020017f757265000000000000000000000000000000000000000000000000000000000081525060400191505060405180910390fd5b80806001019150506112b7565b611494866116a1565b600090505b85518110156115935760008114156114b057611586565b85600182038151811015156114c157fe5b9060200190602002015173ffffffffffffffffffffffffffffffffffffffff1686828151811015156114ef57fe5b9060200190602002015173ffffffffffffffffffffffffffffffffffffffff1614151515611585576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601c8152602001807f6475706c696361746520617070726f76616c732064657465637465640000000081525060200191505060405180910390fd5b5b8080600101915050611499565b86600260008a6000191660001916815260200190815260200160002060030190805190602001906115c592919061232f565b507f9f994d8b60858ce79f51a3eafd335696b0dd099a89b74d21735673e0926c0145888860405180836000191660001916815260200180602001828103825283818151815260200191508051906020019080838360005b8381101561163757808201518184015260208101905061161c565b50505050905090810190601f1680156116645780820380516001836020036101000a031916815260200191505b50935050505060405180910390a15050505050505050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000806000806000806000600188510396506116bc876122fc565b95508594505b6000861015156118cc578594506116d88661231a565b93505b86841115156118b35760018401925084915087848151811015156116fb57fe5b9060200190602002015173ffffffffffffffffffffffffffffffffffffffff16888381518110151561172957fe5b9060200190602002015173ffffffffffffffffffffffffffffffffffffffff161015611753578391505b8683111580156117bc5750878381518110151561176c57fe5b9060200190602002015173ffffffffffffffffffffffffffffffffffffffff16888381518110151561179a57fe5b9060200190602002015173ffffffffffffffffffffffffffffffffffffffff16105b156117c5578291505b848214156117d8576001870193506118ae565b87828151811015156117e657fe5b906020019060200201519050878581518110151561180057fe5b90602001906020020151888381518110151561181857fe5b9060200190602002019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff168152505080888681518110151561186357fe5b9060200190602002019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff16815250508194506118ab8561231a565b93505b6116db565b60008614156118c1576118cc565b6001860395506116c2565b5b6000871115611b935787878151811015156118e457fe5b9060200190602002015190508760008151811015156118ff57fe5b90602001906020020151888881518110151561191757fe5b9060200190602002019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff16815250508088600081518110151561196357fe5b9060200190602002019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff1681525050600187039650600094506119b3600061231a565b93505b8684111515611b8e5760018401925084915087848151811015156119d657fe5b9060200190602002015173ffffffffffffffffffffffffffffffffffffffff168883815181101515611a0457fe5b9060200190602002015173ffffffffffffffffffffffffffffffffffffffff161015611a2e578391505b868311158015611a9757508783815181101515611a4757fe5b9060200190602002015173ffffffffffffffffffffffffffffffffffffffff168883815181101515611a7557fe5b9060200190602002015173ffffffffffffffffffffffffffffffffffffffff16105b15611aa0578291505b84821415611ab357600187019350611b89565b8782815181101515611ac157fe5b9060200190602002015190508785815181101515611adb57fe5b906020019060200201518883815181101515611af357fe5b9060200190602002019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff1681525050808886815181101515611b3e57fe5b9060200190602002019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff1681525050819450611b868561231a565b93505b6119b6565b6118cd565b5050505050505050565b606080606060026000856000191660001916815260200190815260200160002060050160009054906101000a900460ff161515611c42576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260158152602001807f47726f757020646f6573206e6f7420657869737473000000000000000000000081525060200191505060405180910390fd5b600260008560001916600019168152602001908152602001600020600101600260008660001916600019168152602001908152602001600020600201600260008760001916600019168152602001908152602001600020600301828054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015611d315780601f10611d0657610100808354040283529160200191611d31565b820191906000526020600020905b815481529060010190602001808311611d1457829003601f168201915b5050505050925081805480602002602001604051908101604052809291908181526020018280548015611db957602002820191906000526020600020905b8160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019060010190808311611d6f575b50505050509150808054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015611e555780601f10611e2a57610100808354040283529160200191611e55565b820191906000526020600020905b815481529060010190602001808311611e3857829003601f168201915b505050505090509250925092509193909250565b60026000846000191660001916815260200190815260200160002060050160009054906101000a900460ff16151515611f30576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260288152602001807f412067726f757020776974682074686520676976656e20696420616c7265616481526020017f792065786973747300000000000000000000000000000000000000000000000081525060400191505060405180910390fd5b3360026000856000191660001916815260200190815260200160002060000160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550816002600085600019166000191681526020019081526020016000206001019080519060200190611fbf92919061232f565b506002600084600019166000191681526020019081526020016000206002013390806001815401808255809150509060018203906000526020600020016000909192909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505080600260008560001916600019168152602001908152602001600020600301908051906020019061207492919061232f565b50600160026000856000191660001916815260200190815260200160002060050160006101000a81548160ff021916908315150217905550600160026000856000191660001916815260200190815260200160002060040160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055507fb78b34f3219f25d6305837697a0e5d110975d6be50317c9a2e815823306c35768360405180826000191660001916815260200191505060405180910390a1505050565b600080600090505b60026000856000191660001916815260200190815260200160002060020180549050811015612234578273ffffffffffffffffffffffffffffffffffffffff16600260008660001916600019168152602001908152602001600020600201828154811015156121d757fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614156122275760019150612239565b808060010191505061216c565b600091505b5092915050565b60008573ffffffffffffffffffffffffffffffffffffffff16600186868686604051600081526020016040526040518085600019166000191681526020018460ff1660ff1681526020018360001916600019168152602001826000191660001916815260200194505050505060206040516020810390808403906000865af11580156122d0573d6000803e3d6000fd5b5050506020604051035173ffffffffffffffffffffffffffffffffffffffff1614905095945050505050565b60008060018303905060028181151561231157fe5b04915050919050565b60008060028302905060018101915050919050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061237057805160ff191683800117855561239e565b8280016001018555821561239e579182015b8281111561239d578251825591602001919060010190612382565b5b5090506123ab91906123e4565b5090565b60a060405190810160405280606081526020016060815260200160008019168152602001606081526020016000151581525090565b61240691905b808211156124025760008160009055506001016123ea565b5090565b905600a165627a7a72305820c3159bbc53a9fa02742b9123d34256e0ea8af0d8f7b032a684d2a0b3ce1d5c360029`

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

// GetGroup is a free data retrieval call binding the contract method 0xb567d4ba.
//
// Solidity: function getGroup(groupId bytes32) constant returns(string, address[], string)
func (_Eth *EthCaller) GetGroup(opts *bind.CallOpts, groupId [32]byte) (string, []common.Address, string, error) {
	var (
		ret0 = new(string)
		ret1 = new([]common.Address)
		ret2 = new(string)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
	}
	err := _Eth.contract.Call(opts, out, "getGroup", groupId)
	return *ret0, *ret1, *ret2, err
}

// GetGroup is a free data retrieval call binding the contract method 0xb567d4ba.
//
// Solidity: function getGroup(groupId bytes32) constant returns(string, address[], string)
func (_Eth *EthSession) GetGroup(groupId [32]byte) (string, []common.Address, string, error) {
	return _Eth.Contract.GetGroup(&_Eth.CallOpts, groupId)
}

// GetGroup is a free data retrieval call binding the contract method 0xb567d4ba.
//
// Solidity: function getGroup(groupId bytes32) constant returns(string, address[], string)
func (_Eth *EthCallerSession) GetGroup(groupId [32]byte) (string, []common.Address, string, error) {
	return _Eth.Contract.GetGroup(&_Eth.CallOpts, groupId)
}

// GetUser is a free data retrieval call binding the contract method 0x6f77926b.
//
// Solidity: function getUser(id address) constant returns(string, string, bytes32)
func (_Eth *EthCaller) GetUser(opts *bind.CallOpts, id common.Address) (string, string, [32]byte, error) {
	var (
		ret0 = new(string)
		ret1 = new(string)
		ret2 = new([32]byte)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
	}
	err := _Eth.contract.Call(opts, out, "getUser", id)
	return *ret0, *ret1, *ret2, err
}

// GetUser is a free data retrieval call binding the contract method 0x6f77926b.
//
// Solidity: function getUser(id address) constant returns(string, string, bytes32)
func (_Eth *EthSession) GetUser(id common.Address) (string, string, [32]byte, error) {
	return _Eth.Contract.GetUser(&_Eth.CallOpts, id)
}

// GetUser is a free data retrieval call binding the contract method 0x6f77926b.
//
// Solidity: function getUser(id address) constant returns(string, string, bytes32)
func (_Eth *EthCallerSession) GetUser(id common.Address) (string, string, [32]byte, error) {
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

// CreateGroup is a paid mutator transaction binding the contract method 0xe89d7c9d.
//
// Solidity: function createGroup(id bytes32, name string, ipfsPath string) returns()
func (_Eth *EthTransactor) CreateGroup(opts *bind.TransactOpts, id [32]byte, name string, ipfsPath string) (*types.Transaction, error) {
	return _Eth.contract.Transact(opts, "createGroup", id, name, ipfsPath)
}

// CreateGroup is a paid mutator transaction binding the contract method 0xe89d7c9d.
//
// Solidity: function createGroup(id bytes32, name string, ipfsPath string) returns()
func (_Eth *EthSession) CreateGroup(id [32]byte, name string, ipfsPath string) (*types.Transaction, error) {
	return _Eth.Contract.CreateGroup(&_Eth.TransactOpts, id, name, ipfsPath)
}

// CreateGroup is a paid mutator transaction binding the contract method 0xe89d7c9d.
//
// Solidity: function createGroup(id bytes32, name string, ipfsPath string) returns()
func (_Eth *EthTransactorSession) CreateGroup(id [32]byte, name string, ipfsPath string) (*types.Transaction, error) {
	return _Eth.Contract.CreateGroup(&_Eth.TransactOpts, id, name, ipfsPath)
}

// HeapSort is a paid mutator transaction binding the contract method 0xabab02b5.
//
// Solidity: function heapSort(self address[]) returns()
func (_Eth *EthTransactor) HeapSort(opts *bind.TransactOpts, self []common.Address) (*types.Transaction, error) {
	return _Eth.contract.Transact(opts, "heapSort", self)
}

// HeapSort is a paid mutator transaction binding the contract method 0xabab02b5.
//
// Solidity: function heapSort(self address[]) returns()
func (_Eth *EthSession) HeapSort(self []common.Address) (*types.Transaction, error) {
	return _Eth.Contract.HeapSort(&_Eth.TransactOpts, self)
}

// HeapSort is a paid mutator transaction binding the contract method 0xabab02b5.
//
// Solidity: function heapSort(self address[]) returns()
func (_Eth *EthTransactorSession) HeapSort(self []common.Address) (*types.Transaction, error) {
	return _Eth.Contract.HeapSort(&_Eth.TransactOpts, self)
}

// InviteUser is a paid mutator transaction binding the contract method 0x7b799d87.
//
// Solidity: function inviteUser(groupId bytes32, newMember address) returns()
func (_Eth *EthTransactor) InviteUser(opts *bind.TransactOpts, groupId [32]byte, newMember common.Address) (*types.Transaction, error) {
	return _Eth.contract.Transact(opts, "inviteUser", groupId, newMember)
}

// InviteUser is a paid mutator transaction binding the contract method 0x7b799d87.
//
// Solidity: function inviteUser(groupId bytes32, newMember address) returns()
func (_Eth *EthSession) InviteUser(groupId [32]byte, newMember common.Address) (*types.Transaction, error) {
	return _Eth.Contract.InviteUser(&_Eth.TransactOpts, groupId, newMember)
}

// InviteUser is a paid mutator transaction binding the contract method 0x7b799d87.
//
// Solidity: function inviteUser(groupId bytes32, newMember address) returns()
func (_Eth *EthTransactorSession) InviteUser(groupId [32]byte, newMember common.Address) (*types.Transaction, error) {
	return _Eth.Contract.InviteUser(&_Eth.TransactOpts, groupId, newMember)
}

// RegisterUser is a paid mutator transaction binding the contract method 0x68f6d82a.
//
// Solidity: function registerUser(name string, ipfsPeerId string, boxingKey bytes32) returns()
func (_Eth *EthTransactor) RegisterUser(opts *bind.TransactOpts, name string, ipfsPeerId string, boxingKey [32]byte) (*types.Transaction, error) {
	return _Eth.contract.Transact(opts, "registerUser", name, ipfsPeerId, boxingKey)
}

// RegisterUser is a paid mutator transaction binding the contract method 0x68f6d82a.
//
// Solidity: function registerUser(name string, ipfsPeerId string, boxingKey bytes32) returns()
func (_Eth *EthSession) RegisterUser(name string, ipfsPeerId string, boxingKey [32]byte) (*types.Transaction, error) {
	return _Eth.Contract.RegisterUser(&_Eth.TransactOpts, name, ipfsPeerId, boxingKey)
}

// RegisterUser is a paid mutator transaction binding the contract method 0x68f6d82a.
//
// Solidity: function registerUser(name string, ipfsPeerId string, boxingKey bytes32) returns()
func (_Eth *EthTransactorSession) RegisterUser(name string, ipfsPeerId string, boxingKey [32]byte) (*types.Transaction, error) {
	return _Eth.Contract.RegisterUser(&_Eth.TransactOpts, name, ipfsPeerId, boxingKey)
}

// UpdateGroupIpfsPath is a paid mutator transaction binding the contract method 0x7e21f479.
//
// Solidity: function updateGroupIpfsPath(groupId bytes32, newIpfsPath string, members address[], rs bytes32[], ss bytes32[], vs uint8[]) returns()
func (_Eth *EthTransactor) UpdateGroupIpfsPath(opts *bind.TransactOpts, groupId [32]byte, newIpfsPath string, members []common.Address, rs [][32]byte, ss [][32]byte, vs []uint8) (*types.Transaction, error) {
	return _Eth.contract.Transact(opts, "updateGroupIpfsPath", groupId, newIpfsPath, members, rs, ss, vs)
}

// UpdateGroupIpfsPath is a paid mutator transaction binding the contract method 0x7e21f479.
//
// Solidity: function updateGroupIpfsPath(groupId bytes32, newIpfsPath string, members address[], rs bytes32[], ss bytes32[], vs uint8[]) returns()
func (_Eth *EthSession) UpdateGroupIpfsPath(groupId [32]byte, newIpfsPath string, members []common.Address, rs [][32]byte, ss [][32]byte, vs []uint8) (*types.Transaction, error) {
	return _Eth.Contract.UpdateGroupIpfsPath(&_Eth.TransactOpts, groupId, newIpfsPath, members, rs, ss, vs)
}

// UpdateGroupIpfsPath is a paid mutator transaction binding the contract method 0x7e21f479.
//
// Solidity: function updateGroupIpfsPath(groupId bytes32, newIpfsPath string, members address[], rs bytes32[], ss bytes32[], vs uint8[]) returns()
func (_Eth *EthTransactorSession) UpdateGroupIpfsPath(groupId [32]byte, newIpfsPath string, members []common.Address, rs [][32]byte, ss [][32]byte, vs []uint8) (*types.Transaction, error) {
	return _Eth.Contract.UpdateGroupIpfsPath(&_Eth.TransactOpts, groupId, newIpfsPath, members, rs, ss, vs)
}

// EthDebugIterator is returned from FilterDebug and is used to iterate over the raw logs and unpacked data for Debug events raised by the Eth contract.
type EthDebugIterator struct {
	Event *EthDebug // Event containing the contract specifics and raw log

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
func (it *EthDebugIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthDebug)
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
		it.Event = new(EthDebug)
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
func (it *EthDebugIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthDebugIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthDebug represents a Debug event raised by the Eth contract.
type EthDebug struct {
	Msg []byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterDebug is a free log retrieval operation binding the contract event 0x31fe5ce544784f2e2e62fd1c70cb1ebe6cd1f8e17a185909f82c430b6c7ec470.
//
// Solidity: e Debug(msg bytes)
func (_Eth *EthFilterer) FilterDebug(opts *bind.FilterOpts) (*EthDebugIterator, error) {

	logs, sub, err := _Eth.contract.FilterLogs(opts, "Debug")
	if err != nil {
		return nil, err
	}
	return &EthDebugIterator{contract: _Eth.contract, event: "Debug", logs: logs, sub: sub}, nil
}

// WatchDebug is a free log subscription operation binding the contract event 0x31fe5ce544784f2e2e62fd1c70cb1ebe6cd1f8e17a185909f82c430b6c7ec470.
//
// Solidity: e Debug(msg bytes)
func (_Eth *EthFilterer) WatchDebug(opts *bind.WatchOpts, sink chan<- *EthDebug) (event.Subscription, error) {

	logs, sub, err := _Eth.contract.WatchLogs(opts, "Debug")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthDebug)
				if err := _Eth.contract.UnpackLog(event, "Debug", log); err != nil {
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

// EthGroupInvitationIterator is returned from FilterGroupInvitation and is used to iterate over the raw logs and unpacked data for GroupInvitation events raised by the Eth contract.
type EthGroupInvitationIterator struct {
	Event *EthGroupInvitation // Event containing the contract specifics and raw log

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
func (it *EthGroupInvitationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthGroupInvitation)
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
		it.Event = new(EthGroupInvitation)
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
func (it *EthGroupInvitationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthGroupInvitationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthGroupInvitation represents a GroupInvitation event raised by the Eth contract.
type EthGroupInvitation struct {
	From    common.Address
	To      common.Address
	GroupId [32]byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterGroupInvitation is a free log retrieval operation binding the contract event 0x9478e2f0a42543d96af3b3661efc5aaa23dd42c9f8c970c1e4f4bd01ab42374a.
//
// Solidity: e GroupInvitation(from address, to address, groupId bytes32)
func (_Eth *EthFilterer) FilterGroupInvitation(opts *bind.FilterOpts) (*EthGroupInvitationIterator, error) {

	logs, sub, err := _Eth.contract.FilterLogs(opts, "GroupInvitation")
	if err != nil {
		return nil, err
	}
	return &EthGroupInvitationIterator{contract: _Eth.contract, event: "GroupInvitation", logs: logs, sub: sub}, nil
}

// WatchGroupInvitation is a free log subscription operation binding the contract event 0x9478e2f0a42543d96af3b3661efc5aaa23dd42c9f8c970c1e4f4bd01ab42374a.
//
// Solidity: e GroupInvitation(from address, to address, groupId bytes32)
func (_Eth *EthFilterer) WatchGroupInvitation(opts *bind.WatchOpts, sink chan<- *EthGroupInvitation) (event.Subscription, error) {

	logs, sub, err := _Eth.contract.WatchLogs(opts, "GroupInvitation")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthGroupInvitation)
				if err := _Eth.contract.UnpackLog(event, "GroupInvitation", log); err != nil {
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

// EthGroupRegisteredIterator is returned from FilterGroupRegistered and is used to iterate over the raw logs and unpacked data for GroupRegistered events raised by the Eth contract.
type EthGroupRegisteredIterator struct {
	Event *EthGroupRegistered // Event containing the contract specifics and raw log

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
func (it *EthGroupRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthGroupRegistered)
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
		it.Event = new(EthGroupRegistered)
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
func (it *EthGroupRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthGroupRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthGroupRegistered represents a GroupRegistered event raised by the Eth contract.
type EthGroupRegistered struct {
	Id  [32]byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterGroupRegistered is a free log retrieval operation binding the contract event 0xb78b34f3219f25d6305837697a0e5d110975d6be50317c9a2e815823306c3576.
//
// Solidity: e GroupRegistered(id bytes32)
func (_Eth *EthFilterer) FilterGroupRegistered(opts *bind.FilterOpts) (*EthGroupRegisteredIterator, error) {

	logs, sub, err := _Eth.contract.FilterLogs(opts, "GroupRegistered")
	if err != nil {
		return nil, err
	}
	return &EthGroupRegisteredIterator{contract: _Eth.contract, event: "GroupRegistered", logs: logs, sub: sub}, nil
}

// WatchGroupRegistered is a free log subscription operation binding the contract event 0xb78b34f3219f25d6305837697a0e5d110975d6be50317c9a2e815823306c3576.
//
// Solidity: e GroupRegistered(id bytes32)
func (_Eth *EthFilterer) WatchGroupRegistered(opts *bind.WatchOpts, sink chan<- *EthGroupRegistered) (event.Subscription, error) {

	logs, sub, err := _Eth.contract.WatchLogs(opts, "GroupRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthGroupRegistered)
				if err := _Eth.contract.UnpackLog(event, "GroupRegistered", log); err != nil {
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

// EthGroupUpdateIpfsPathIterator is returned from FilterGroupUpdateIpfsPath and is used to iterate over the raw logs and unpacked data for GroupUpdateIpfsPath events raised by the Eth contract.
type EthGroupUpdateIpfsPathIterator struct {
	Event *EthGroupUpdateIpfsPath // Event containing the contract specifics and raw log

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
func (it *EthGroupUpdateIpfsPathIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthGroupUpdateIpfsPath)
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
		it.Event = new(EthGroupUpdateIpfsPath)
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
func (it *EthGroupUpdateIpfsPathIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthGroupUpdateIpfsPathIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthGroupUpdateIpfsPath represents a GroupUpdateIpfsPath event raised by the Eth contract.
type EthGroupUpdateIpfsPath struct {
	GroupId  [32]byte
	IpfsPath string
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterGroupUpdateIpfsPath is a free log retrieval operation binding the contract event 0x9f994d8b60858ce79f51a3eafd335696b0dd099a89b74d21735673e0926c0145.
//
// Solidity: e GroupUpdateIpfsPath(groupId bytes32, ipfsPath string)
func (_Eth *EthFilterer) FilterGroupUpdateIpfsPath(opts *bind.FilterOpts) (*EthGroupUpdateIpfsPathIterator, error) {

	logs, sub, err := _Eth.contract.FilterLogs(opts, "GroupUpdateIpfsPath")
	if err != nil {
		return nil, err
	}
	return &EthGroupUpdateIpfsPathIterator{contract: _Eth.contract, event: "GroupUpdateIpfsPath", logs: logs, sub: sub}, nil
}

// WatchGroupUpdateIpfsPath is a free log subscription operation binding the contract event 0x9f994d8b60858ce79f51a3eafd335696b0dd099a89b74d21735673e0926c0145.
//
// Solidity: e GroupUpdateIpfsPath(groupId bytes32, ipfsPath string)
func (_Eth *EthFilterer) WatchGroupUpdateIpfsPath(opts *bind.WatchOpts, sink chan<- *EthGroupUpdateIpfsPath) (event.Subscription, error) {

	logs, sub, err := _Eth.contract.WatchLogs(opts, "GroupUpdateIpfsPath")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthGroupUpdateIpfsPath)
				if err := _Eth.contract.UnpackLog(event, "GroupUpdateIpfsPath", log); err != nil {
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
