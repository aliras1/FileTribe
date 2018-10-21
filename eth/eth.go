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
const EthABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"UserRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"GroupRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"groupId\",\"type\":\"bytes32\"}],\"name\":\"GroupInvitation\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"groupId\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"ipfsPath\",\"type\":\"string\"}],\"name\":\"GroupUpdateIpfsPath\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"MessageSent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"msg\",\"type\":\"bytes\"}],\"name\":\"Debug\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"ipfsPeerId\",\"type\":\"string\"},{\"name\":\"boxingKey\",\"type\":\"bytes32\"}],\"name\":\"registerUser\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"address\"}],\"name\":\"isUserRegistered\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"address\"}],\"name\":\"getUser\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"},{\"name\":\"\",\"type\":\"string\"},{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\"},{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"ipfsPath\",\"type\":\"string\"}],\"name\":\"createGroup\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"groupId\",\"type\":\"bytes32\"}],\"name\":\"getGroup\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"},{\"name\":\"\",\"type\":\"address[]\"},{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"groupId\",\"type\":\"bytes32\"},{\"name\":\"newMember\",\"type\":\"address\"}],\"name\":\"inviteUser\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"groupId\",\"type\":\"bytes32\"},{\"name\":\"newIpfsPath\",\"type\":\"string\"},{\"name\":\"members\",\"type\":\"address[]\"},{\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"name\":\"vs\",\"type\":\"uint8[]\"}],\"name\":\"updateGroupIpfsPath\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"sendMessage\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"self\",\"type\":\"address[]\"}],\"name\":\"heapSort\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// EthBin is the compiled bytecode used for deploying new contracts.
const EthBin = `0x608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550612548806100606000396000f3006080604052600436106100a4576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063163f7522146100a957806368f6d82a146101045780636f77926b146101c15780637b799d87146102f85780637e21f4791461034957806382646a58146104cc5780638da5cb5b14610535578063abab02b51461058c578063b567d4ba146105f2578063e89d7c9d14610750575b600080fd5b3480156100b557600080fd5b506100ea600480360381019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061080d565b604051808215151515815260200191505060405180910390f35b34801561011057600080fd5b506101bf600480360381019080803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290803590602001908201803590602001908080601f01602080910402602001604051908101604052809392919081815260200183838082843782019150505050505091929192908035600019169060200190929190505050610866565b005b3480156101cd57600080fd5b50610202600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610ae6565b6040518080602001806020018460001916600019168152602001838103835286818151815260200191508051906020019080838360005b83811015610254578082015181840152602081019050610239565b50505050905090810190601f1680156102815780820380516001836020036101000a031916815260200191505b50838103825285818151815260200191508051906020019080838360005b838110156102ba57808201518184015260208101905061029f565b50505050905090810190601f1680156102e75780820380516001836020036101000a031916815260200191505b509550505050505060405180910390f35b34801561030457600080fd5b506103476004803603810190808035600019169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610df4565b005b34801561035557600080fd5b506104ca6004803603810190808035600019169060200190929190803590602001908201803590602001908080601f01602080910402602001604051908101604052809392919081815260200183838082843782019150505050505091929192908035906020019082018035906020019080806020026020016040519081016040528093929190818152602001838360200280828437820191505050505050919291929080359060200190820180359060200190808060200260200160405190810160405280939291908181526020018383602002808284378201915050505050509192919290803590602001908201803590602001908080602002602001604051908101604052809392919081815260200183836020028082843782019150505050505091929192908035906020019082018035906020019080806020026020016040519081016040528093929190818152602001838360200280828437820191505050505050919291929050505061113b565b005b3480156104d857600080fd5b50610533600480360381019080803590602001908201803590602001908080601f01602080910402602001604051908101604052809392919081815260200183838082843782019150505050505091929192905050506116f0565b005b34801561054157600080fd5b5061054a61178f565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561059857600080fd5b506105f0600480360381019080803590602001908201803590602001908080602002602001604051908101604052809392919081815260200183836020028082843782019150505050505091929192905050506117b4565b005b3480156105fe57600080fd5b506106216004803603810190808035600019169060200190929190505050611cb0565b60405180806020018060200180602001848103845287818151815260200191508051906020019080838360005b8381101561066957808201518184015260208101905061064e565b50505050905090810190601f1680156106965780820380516001836020036101000a031916815260200191505b50848103835286818151815260200191508051906020019060200280838360005b838110156106d25780820151818401526020810190506106b7565b50505050905001848103825285818151815260200191508051906020019080838360005b838110156107115780820151818401526020810190506106f6565b50505050905090810190601f16801561073e5780820380516001836020036101000a031916815260200191505b50965050505050505060405180910390f35b34801561075c57600080fd5b5061080b6004803603810190808035600019169060200190929190803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290505050611f7c565b005b6000600160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060040160009054906101000a900460ff169050919050565b600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060040160009054906101000a900460ff1615151561092b576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260178152602001807f557365726e616d6520616c72656164792065786973747300000000000000000081525060200191505060405180910390fd5b82600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000019080519060200190610981929190612442565b5081600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010190805190602001906109d8929190612442565b5080600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600201816000191690555060018060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060040160006101000a81548160ff0219169083151502179055507f54db7a5cb4735e1aac1f53db512d3390390bb6637bd30ad4bf9fc98667d9b9b933604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390a1505050565b6060806000610af36124c2565b600160008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060040160009054906101000a900460ff161515610bb7576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260138152602001807f5573657220646f6573206e6f742065786973740000000000000000000000000081525060200191505060405180910390fd5b600160008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060a06040519081016040529081600082018054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610c9b5780601f10610c7057610100808354040283529160200191610c9b565b820191906000526020600020905b815481529060010190602001808311610c7e57829003601f168201915b50505050508152602001600182018054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610d3d5780601f10610d1257610100808354040283529160200191610d3d565b820191906000526020600020905b815481529060010190602001808311610d2057829003601f168201915b5050505050815260200160028201546000191660001916815260200160038201805480602002602001604051908101604052809291908181526020018280548015610dab57602002820191906000526020600020905b81546000191681526020019060010190808311610d93575b505050505081526020016004820160009054906101000a900460ff1615151515815250509050806000015181602001518260400151829250819150935093509350509193909250565b6001151560026000846000191660001916815260200190815260200160002060040160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff161515141515610ed8576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260138152602001807f557365722063616e206e6f7420696e766974650000000000000000000000000081525060200191505060405180910390fd5b600160008273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060040160009054906101000a900460ff161515610f9c576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f43616e206e6f7420696e76697465206e6f6e206578697374656e74207573657281525060200191505060405180910390fd5b6002600083600019166000191681526020019081526020016000206002018190806001815401808255809150509060018203906000526020600020016000909192909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050600160008273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206003018290806001815401808255809150509060018203906000526020600020016000909192909190915090600019169055507f9478e2f0a42543d96af3b3661efc5aaa23dd42c9f8c970c1e4f4bd01ab42374a338284604051808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018260001916600019168152602001935050505060405180910390a15050565b60008060026000896000191660001916815260200190815260200160002060050160009054906101000a900460ff1615156111de576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f67726f757020646f6573206e6f7420657869737400000000000000000000000081525060200191505060405180910390fd5b855185511415156111ee57600080fd5b855184511415156111fe57600080fd5b8551835114151561120e57600080fd5b60028060008a600019166000191681526020019081526020016000206002018054905081151561123a57fe5b04865111151561124957600080fd5b6002600089600019166000191681526020019081526020016000206003018760405180838054600181600116156101000203166002900480156112c35780601f106112a15761010080835404028352918201916112c3565b820191906000526020600020905b8154815290600101906020018083116112af575b505082805190602001908083835b6020831015156112f657805182526020820191506020810190506020830392506112d1565b6001836020036101000a0380198251168184511680821785525050505050509050019250505060405180910390209150600090505b85518110156114ff5761135588878381518110151561134657fe5b90602001906020020151612277565b15156113ef576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602c8152602001807f696e76616c696420617070726f76616c3a2075736572206973206e6f7420612081526020017f67726f7570206d656d626572000000000000000000000000000000000000000081525060400191505060405180910390fd5b611458868281518110151561140057fe5b9060200190602002015183858481518110151561141957fe5b90602001906020020151888581518110151561143157fe5b90602001906020020151888681518110151561144957fe5b90602001906020020151612353565b15156114f2576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260238152602001807f696e76616c696420617070726f76616c3a20696e76616c6964207369676e617481526020017f757265000000000000000000000000000000000000000000000000000000000081525060400191505060405180910390fd5b808060010191505061132b565b611508866117b4565b600090505b8551811015611607576000811415611524576115fa565b856001820381518110151561153557fe5b9060200190602002015173ffffffffffffffffffffffffffffffffffffffff16868281518110151561156357fe5b9060200190602002015173ffffffffffffffffffffffffffffffffffffffff16141515156115f9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601c8152602001807f6475706c696361746520617070726f76616c732064657465637465640000000081525060200191505060405180910390fd5b5b808060010191505061150d565b86600260008a600019166000191681526020019081526020016000206003019080519060200190611639929190612442565b507f9f994d8b60858ce79f51a3eafd335696b0dd099a89b74d21735673e0926c0145888860405180836000191660001916815260200180602001828103825283818151815260200191508051906020019080838360005b838110156116ab578082015181840152602081019050611690565b50505050905090810190601f1680156116d85780820380516001836020036101000a031916815260200191505b50935050505060405180910390a15050505050505050565b7f8c5261668696ce22758910d05bab8f186d6eb247ceac2af2e82c7dc17669b036816040518080602001828103825283818151815260200191508051906020019080838360005b83811015611752578082015181840152602081019050611737565b50505050905090810190601f16801561177f5780820380516001836020036101000a031916815260200191505b509250505060405180910390a150565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000806000806000806000600188510396506117cf8761240f565b95508594505b6000861015156119df578594506117eb8661242d565b93505b86841115156119c657600184019250849150878481518110151561180e57fe5b9060200190602002015173ffffffffffffffffffffffffffffffffffffffff16888381518110151561183c57fe5b9060200190602002015173ffffffffffffffffffffffffffffffffffffffff161015611866578391505b8683111580156118cf5750878381518110151561187f57fe5b9060200190602002015173ffffffffffffffffffffffffffffffffffffffff1688838151811015156118ad57fe5b9060200190602002015173ffffffffffffffffffffffffffffffffffffffff16105b156118d8578291505b848214156118eb576001870193506119c1565b87828151811015156118f957fe5b906020019060200201519050878581518110151561191357fe5b90602001906020020151888381518110151561192b57fe5b9060200190602002019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff168152505080888681518110151561197657fe5b9060200190602002019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff16815250508194506119be8561242d565b93505b6117ee565b60008614156119d4576119df565b6001860395506117d5565b5b6000871115611ca65787878151811015156119f757fe5b906020019060200201519050876000815181101515611a1257fe5b906020019060200201518888815181101515611a2a57fe5b9060200190602002019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff168152505080886000815181101515611a7657fe5b9060200190602002019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff168152505060018703965060009450611ac6600061242d565b93505b8684111515611ca1576001840192508491508784815181101515611ae957fe5b9060200190602002015173ffffffffffffffffffffffffffffffffffffffff168883815181101515611b1757fe5b9060200190602002015173ffffffffffffffffffffffffffffffffffffffff161015611b41578391505b868311158015611baa57508783815181101515611b5a57fe5b9060200190602002015173ffffffffffffffffffffffffffffffffffffffff168883815181101515611b8857fe5b9060200190602002015173ffffffffffffffffffffffffffffffffffffffff16105b15611bb3578291505b84821415611bc657600187019350611c9c565b8782815181101515611bd457fe5b9060200190602002015190508785815181101515611bee57fe5b906020019060200201518883815181101515611c0657fe5b9060200190602002019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff1681525050808886815181101515611c5157fe5b9060200190602002019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff1681525050819450611c998561242d565b93505b611ac9565b6119e0565b5050505050505050565b606080606060026000856000191660001916815260200190815260200160002060050160009054906101000a900460ff161515611d55576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260158152602001807f47726f757020646f6573206e6f7420657869737473000000000000000000000081525060200191505060405180910390fd5b600260008560001916600019168152602001908152602001600020600101600260008660001916600019168152602001908152602001600020600201600260008760001916600019168152602001908152602001600020600301828054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015611e445780601f10611e1957610100808354040283529160200191611e44565b820191906000526020600020905b815481529060010190602001808311611e2757829003601f168201915b5050505050925081805480602002602001604051908101604052809291908181526020018280548015611ecc57602002820191906000526020600020905b8160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019060010190808311611e82575b50505050509150808054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015611f685780601f10611f3d57610100808354040283529160200191611f68565b820191906000526020600020905b815481529060010190602001808311611f4b57829003601f168201915b505050505090509250925092509193909250565b60026000846000191660001916815260200190815260200160002060050160009054906101000a900460ff16151515612043576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260288152602001807f412067726f757020776974682074686520676976656e20696420616c7265616481526020017f792065786973747300000000000000000000000000000000000000000000000081525060400191505060405180910390fd5b3360026000856000191660001916815260200190815260200160002060000160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508160026000856000191660001916815260200190815260200160002060010190805190602001906120d2929190612442565b506002600084600019166000191681526020019081526020016000206002013390806001815401808255809150509060018203906000526020600020016000909192909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050806002600085600019166000191681526020019081526020016000206003019080519060200190612187929190612442565b50600160026000856000191660001916815260200190815260200160002060050160006101000a81548160ff021916908315150217905550600160026000856000191660001916815260200190815260200160002060040160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055507fb78b34f3219f25d6305837697a0e5d110975d6be50317c9a2e815823306c35768360405180826000191660001916815260200191505060405180910390a1505050565b600080600090505b60026000856000191660001916815260200190815260200160002060020180549050811015612347578273ffffffffffffffffffffffffffffffffffffffff16600260008660001916600019168152602001908152602001600020600201828154811015156122ea57fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16141561233a576001915061234c565b808060010191505061227f565b600091505b5092915050565b60008573ffffffffffffffffffffffffffffffffffffffff16600186868686604051600081526020016040526040518085600019166000191681526020018460ff1660ff1681526020018360001916600019168152602001826000191660001916815260200194505050505060206040516020810390808403906000865af11580156123e3573d6000803e3d6000fd5b5050506020604051035173ffffffffffffffffffffffffffffffffffffffff1614905095945050505050565b60008060018303905060028181151561242457fe5b04915050919050565b60008060028302905060018101915050919050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061248357805160ff19168380011785556124b1565b828001600101855582156124b1579182015b828111156124b0578251825591602001919060010190612495565b5b5090506124be91906124f7565b5090565b60a060405190810160405280606081526020016060815260200160008019168152602001606081526020016000151581525090565b61251991905b808211156125155760008160009055506001016124fd565b5090565b905600a165627a7a72305820530511e5cde54d85e38398a8ee3bd83f7ef674ca8ce91275954345846cc2628c0029`

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
