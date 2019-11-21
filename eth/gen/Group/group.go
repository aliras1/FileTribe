// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Group

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

// GroupABI is the input ABI used to generate the binding from.
const GroupABI = "[{\"inputs\":[{\"internalType\":\"contractIFileTribeDApp\",\"name\":\"fileTribe\",\"type\":\"address\"},{\"internalType\":\"contractIAccount\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"ipfsHash\",\"type\":\"bytes\"},{\"internalType\":\"uint256[4]\",\"name\":\"verifyKey\",\"type\":\"uint256[4]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"msg\",\"type\":\"uint256\"}],\"name\":\"Debug\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"groupId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"ipfsHash\",\"type\":\"bytes\"}],\"name\":\"GroupUpdateIpfsHash\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractIGroup\",\"name\":\"groupAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractIAccount\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"InvitationAccepted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractIGroup\",\"name\":\"group\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractIAccount\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"InvitationDeclined\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractIGroup\",\"name\":\"group\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractIAccount\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"InvitationSent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractIGroup\",\"name\":\"group\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"ipfsHash\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"contractIAccount\",\"name\":\"proposer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"IpfsHashChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractIGroup\",\"name\":\"group\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractIAccount\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"MemberLeft\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractIGroup\",\"name\":\"group\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractIConsensus\",\"name\":\"consensus\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"NewConsensus\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"_g1\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"_g2\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_ipfsHash\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"_verifyKey\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"isMember\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"newIpfsHash\",\"type\":\"bytes\"},{\"internalType\":\"uint256[2]\",\"name\":\"sig\",\"type\":\"uint256[2]\"}],\"name\":\"commitWithGroupSig\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"newIpfsHash\",\"type\":\"bytes\"}],\"name\":\"commit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"onChangeIpfsHashConsensus\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"leave\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"memberOwner\",\"type\":\"address\"}],\"name\":\"kick\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"contractIAccount\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"invite\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"join\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"decline\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"getConsensus\",\"outputs\":[{\"internalType\":\"contractIConsensus\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"memberOwners\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"threshold\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// GroupBin is the compiled bytecode used for deploying new contracts.
var GroupBin = "0x60c060405260016080908152600260a08190526200002091600c91620004ec565b5060405180608001604052807f198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c281526020017f1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed81526020017f090689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b81526020017f12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa815250600e906004620000d492919062000534565b50348015620000e257600080fd5b50604051620021df380380620021df83398181016040526101008110156200010957600080fd5b815160208301516040808501805191519395929483019291846401000000008211156200013557600080fd5b9083019060208201858111156200014b57600080fd5b82516401000000008111828201881017156200016657600080fd5b82525081516020918201929091019080838360005b83811015620001955781810151838201526020016200017b565b50505050905090810190601f168015620001c35780820380516001836020036101000a031916815260200191505b5060405260200180516040519392919084640100000000821115620001e757600080fd5b908301906020820185811115620001fd57600080fd5b82516401000000008111828201881017156200021857600080fd5b82525081516020918201929091019080838360005b83811015620002475781810151838201526020016200022d565b50505050905090810190601f168015620002755780820380516001836020036101000a031916815260200191505b506040819052600080546001600160a01b0319166001600160a01b038981169190911780835560209490940195508894509290921691907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908290a3506000846001600160a01b0316638da5cb5b6040518163ffffffff1660e01b815260040160206040518083038186803b1580156200030e57600080fd5b505afa15801562000323573d6000803e3d6000fd5b505050506040513d60208110156200033a57600080fd5b505190506200034d600383600462000534565b50600180546001600160a01b0319166001600160a01b03881617905583516200037e90600290602087019062000565565b5082516200039490600990602086019062000565565b506000600a8190556007805460018181019092557fa66cc928b5edb82af9bd49922954155ab7b0942694bea4ce44661d9a8736c6880180546001600160a01b038086166001600160a01b0319928316811790935591845260086020908152604080862080548c8616941684179055935484517f15db20690000000000000000000000000000000000000000000000000000000081526004810193909352935193909216936315db206993602480840194939192918390030190829087803b1580156200045f57600080fd5b505af115801562000474573d6000803e3d6000fd5b505050506040513d60208110156200048b57600080fd5b50516001600160a01b039182166000908152600860205260409020600101805475010000000000000000000000000000000000000000006001600160a01b0319909116939092169290921760ff60a81b191617905550620005f79350505050565b826002810192821562000522579160200282015b8281111562000522578251829060ff1690559160200191906001019062000500565b5062000530929150620005d7565b5090565b826004810192821562000522579160200282015b828111156200052257825182559160200191906001019062000548565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10620005a857805160ff191683800117855562000522565b828001600101855582156200052257918201828111156200052257825182559160200191906001019062000548565b620005f491905b80821115620005305760008155600101620005de565b90565b611bd880620006076000396000f3fe608060405234801561001057600080fd5b50600436106101375760003560e01c8063a230c524116100b8578063d53d64041161007c578063d53d6404146103f7578063d66d9e1914610414578063e7acfda11461041c578063e8f738e114610474578063ebe1c2301461049a578063f2fde38b1461056957610137565b8063a230c5241461033c578063ab04010714610362578063b688a3631461036a578063cc2867d714610372578063d28d8852146103ef57610137565b8063731f6fd8116100ff578063731f6fd8146102495780638379088e146102665780638da5cb5b146102d65780638f32d59b146102fa57806396c551751461031657610137565b806312f8ad051461013c57806342cde4e81461016b5780634b77c4681461017357806366fd3cd81461019b578063715018a614610241575b600080fd5b6101596004803603602081101561015257600080fd5b503561058f565b60408051918252519081900360200190f35b6101596105a3565b6101996004803603602081101561018957600080fd5b50356001600160a01b03166105ae565b005b610199600480360360208110156101b157600080fd5b8101906020810181356401000000008111156101cc57600080fd5b8201836020820111156101de57600080fd5b8035906020019184600183028401116401000000008311171561020057600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295506107bf945050505050565b610199610a6b565b6101596004803603602081101561025f57600080fd5b5035610ac6565b6101996004803603602081101561027c57600080fd5b81019060208101813564010000000081111561029757600080fd5b8201836020820111156102a957600080fd5b803590602001918460018302840111640100000000831117156102cb57600080fd5b509092509050610ad3565b6102de610d78565b604080516001600160a01b039092168252519081900360200190f35b610302610d87565b604080519115158252519081900360200190f35b6101996004803603602081101561032c57600080fd5b50356001600160a01b0316610d98565b6103026004803603602081101561035257600080fd5b50356001600160a01b0316610f45565b610199610f6d565b61019961108c565b61037a6112b5565b6040805160208082528351818301528351919283929083019185019080838360005b838110156103b457818101518382015260200161039c565b50505050905090810190601f1680156103e15780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b61037a611343565b6101596004803603602081101561040d57600080fd5b503561139b565b6101996113a8565b61042461154b565b60408051602080825283518183015283519192839290830191858101910280838360005b83811015610460578181015183820152602001610448565b505050509050019250505060405180910390f35b6102de6004803603602081101561048a57600080fd5b50356001600160a01b03166115ad565b610199600480360360608110156104b057600080fd5b8101906020810181356401000000008111156104cb57600080fd5b8201836020820111156104dd57600080fd5b803590602001918460018302840111640100000000831117156104ff57600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152505060408051808201825293969594818101949350915060029083908390808284376000920191909152509194506115ce9350505050565b6101996004803603602081101561057f57600080fd5b50356001600160a01b0316611982565b6003816004811061059c57fe5b0154905081565b600754600290045b90565b6105b733610f45565b6105c057600080fd5b6000816001600160a01b0316638da5cb5b6040518163ffffffff1660e01b815260040160206040518083038186803b1580156105fb57600080fd5b505afa15801561060f573d6000803e3d6000fd5b505050506040513d602081101561062557600080fd5b50516001600160a01b038116600090815260086020526040902060010154909150600160a81b900460ff161561068c5760405162461bcd60e51b815260040180806020018281038252602a815260200180611b7a602a913960400191505060405180910390fd5b6001600160a01b038181166000908152600b602052604090205416156106f9576040805162461bcd60e51b815260206004820181905260248201527f6163636f756e742068617320616c7265616479206265656e20696e7669746564604482015290519081900360640190fd5b816001600160a01b031663eec30bfd6040518163ffffffff1660e01b8152600401600060405180830381600087803b15801561073457600080fd5b505af1158015610748573d6000803e3d6000fd5b5050604080513081526001600160a01b038616602082015281517f4509853340a1799d2c49c6cb19f9988fc2fa7f8c37aa55a7d865e19a260cfdc99450908190039091019150a16001600160a01b039081166000908152600b602052604090208054919092166001600160a01b0319909116179055565b6107c833610f45565b6107d157600080fd5b6007546001141561092e57600154604080516395184d3b60e01b815233600482015290517ff7084e4afb12f9dc96cee0f3c56fee76afe9598f22f0e3ce04f21b4c18cdc26792309285926001600160a01b03909216916395184d3b91602480820192602092909190829003018186803b15801561084d57600080fd5b505afa158015610861573d6000803e3d6000fd5b505050506040513d602081101561087757600080fd5b5051600a54604080516001600160a01b0380871682528416918101919091526060810182905260806020828101828152865192840192909252855160a084019187019080838360005b838110156108d85781810151838201526020016108c0565b50505050905090810190601f1680156109055780820380516001836020036101000a031916815260200191505b509550505050505060405180910390a18051610928906009906020840190611a05565b50610a68565b336000908152600860209081526040808320600190810154600a54835163bade603360e01b8152920160248301819052600483019384528651604484015286516001600160a01b0390921695869563bade6033958995939490938493606490910192870191908190849084905b838110156109b357818101518382015260200161099b565b50505050905090810190601f1680156109e05780820380516001836020036101000a031916815260200191505b509350505050600060405180830381600087803b158015610a0057600080fd5b505af1158015610a14573d6000803e3d6000fd5b5050600a54604080513081526001600160a01b0386166020820152600190920182820152517fe407c1bec7341bd0a3dcb85da18f4e0b5126fcc79738f8f1dde296455e9c35aa9350908190036060019150a1505b50565b610a73610d87565b610a7c57600080fd5b600080546040516001600160a01b03909116907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a3600080546001600160a01b0319169055565b600c816002811061059c57fe5b6000336001600160a01b0316635d1ca6316040518163ffffffff1660e01b815260040160206040518083038186803b158015610b0e57600080fd5b505afa158015610b22573d6000803e3d6000fd5b505050506040513d6020811015610b3857600080fd5b5051600a549091508111610b88576040805162461bcd60e51b8152602060048201526012602482015271436f6e73656e73757320657870697265642160701b604482015290519081900360640190fd5b6000336001600160a01b031663e9790d026040518163ffffffff1660e01b815260040160206040518083038186803b158015610bc357600080fd5b505afa158015610bd7573d6000803e3d6000fd5b505050506040513d6020811015610bed57600080fd5b505160408051638da5cb5b60e01b815290519192506000916001600160a01b03841691638da5cb5b916004808301926020929190829003018186803b158015610c3557600080fd5b505afa158015610c49573d6000803e3d6000fd5b505050506040513d6020811015610c5f57600080fd5b50516001600160a01b03808216600090815260086020526040902060010154919250163314610cbf5760405162461bcd60e51b8152600401808060200182810382526027815260200180611b536027913960400191505060405180910390fd5b610ccb60098686611a83565b5082600a819055507ff7084e4afb12f9dc96cee0f3c56fee76afe9598f22f0e3ce04f21b4c18cdc26730868685600a5460405180866001600160a01b03166001600160a01b0316815260200180602001846001600160a01b03166001600160a01b031681526020018381526020018281038252868682818152602001925080828437600083820152604051601f909101601f19169092018290039850909650505050505050a15050505050565b6000546001600160a01b031690565b6000546001600160a01b0316331490565b610da133610f45565b610daa57600080fd5b6001600160a01b038082166000908152600860205260408082205481516315fe21e560e31b81523060048201529151931692839263aff10f28926024808201939182900301818387803b158015610e0057600080fd5b505af1158015610e14573d6000803e3d6000fd5b505050506001600160a01b0382166000908152600860205260408120600101805460ff60a81b191690555b600754811015610efd57826001600160a01b031660078281548110610e6057fe5b6000918252602090912001546001600160a01b03161415610ef557600780546000198101908110610e8d57fe5b600091825260209091200154600780546001600160a01b039092169183908110610eb357fe5b600091825260209091200180546001600160a01b0319166001600160a01b03929092169190911790556007805490610eef906000198301611af1565b50610efd565b600101610e3f565b604080513081526001600160a01b038416602082015281517fa0ae7458a482bd2780bed4edd090c3f9e9e3e600171cdc4709da09e5bd233e18929181900390910190a1505050565b6001600160a01b0316600090815260086020526040902060010154600160a81b900460ff1690565b336000908152600b60205260409020546001600160a01b031680610fd2576040805162461bcd60e51b81526020600482015260176024820152761858d8dbdd5b9d081dd85cc81b9bdd081a5b9d9a5d1959604a1b604482015290519081900360640190fd5b806001600160a01b031663b061d9a96040518163ffffffff1660e01b8152600401600060405180830381600087803b15801561100d57600080fd5b505af1158015611021573d6000803e3d6000fd5b5050336000908152600b602090815260409182902080546001600160a01b031916905581513081526001600160a01b0386169181019190915281517f7ea317cacaf634632d88c1fd20864b1f54f609bb063c678280aac6312715edef9450908190039091019150a150565b336000908152600b60205260409020546001600160a01b0316806110f1576040805162461bcd60e51b81526020600482015260176024820152761858d8dbdd5b9d081dd85cc81b9bdd081a5b9d9a5d1959604a1b604482015290519081900360640190fd5b6007805460018082019092557fa66cc928b5edb82af9bd49922954155ab7b0942694bea4ce44661d9a8736c688018054336001600160a01b03199182168117909255600091825260086020908152604080842080549093166001600160a01b03878116918217909455945481516315db206960e01b8152600481019690965290519216936315db206993602480830194928390030190829087803b15801561119857600080fd5b505af11580156111ac573d6000803e3d6000fd5b505050506040513d60208110156111c257600080fd5b5051336000908152600860205260408082206001018054600160a81b6001600160a01b03199091166001600160a01b039586161760ff60a81b1916179055805163020fa19160e61b81529051928416926383e864409260048084019391929182900301818387803b15801561123657600080fd5b505af115801561124a573d6000803e3d6000fd5b5050336000908152600b602090815260409182902080546001600160a01b031916905581513081526001600160a01b0386169181019190915281517f4d7c243e154e530692e62f8539db65779f5cb85d58831956361697addede5adb9450908190039091019150a150565b6009805460408051602060026001851615610100026000190190941693909304601f8101849004840282018401909252818152929183018282801561133b5780601f106113105761010080835404028352916020019161133b565b820191906000526020600020905b81548152906001019060200180831161131e57829003601f168201915b505050505081565b6002805460408051602060018416156101000260001901909316849004601f8101849004840282018401909252818152929183018282801561133b5780601f106113105761010080835404028352916020019161133b565b600e816004811061059c57fe5b6113b133610f45565b6113ba57600080fd5b336000908152600860205260408082205481516315fe21e560e31b815230600482015291516001600160a01b0390911692839263aff10f28926024808301939282900301818387803b15801561140f57600080fd5b505af1158015611423573d6000803e3d6000fd5b5050336000908152600860205260408120600101805460ff60a81b191690559150505b60075481101561150457336001600160a01b03166007828154811061146757fe5b6000918252602090912001546001600160a01b031614156114fc5760078054600019810190811061149457fe5b600091825260209091200154600780546001600160a01b0390921691839081106114ba57fe5b600091825260209091200180546001600160a01b0319166001600160a01b039290921691909117905560078054906114f6906000198301611af1565b50611504565b600101611446565b604080513081526001600160a01b038416602082015281517fa0ae7458a482bd2780bed4edd090c3f9e9e3e600171cdc4709da09e5bd233e18929181900390910190a15050565b606060078054806020026020016040519081016040528092919081815260200182805480156115a357602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311611585575b5050505050905090565b6001600160a01b039081166000908152600860205260409020600101541690565b6115d733610f45565b6115e057600080fd5b6115e8611b1a565b6040516378b8c33160e11b815260206004820181815285516024840152855173__ecOps_________________________________9363f17186629388939283926044019185019080838360005b8381101561164d578181015183820152602001611635565b50505050905090810190601f16801561167a5780820380516001836020036101000a031916815260200191505b5092505050604080518083038186803b15801561169657600080fd5b505af41580156116aa573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525060408110156116cf57600080fd5b50604080516349ac97d160e01b815291925073__ecOps_________________________________916349ac97d19184916003918791600e916004909101908190869080838360005b8381101561172f578181015183820152602001611717565b5050509201915085905060046020028201915b8154815260200190600101908083116117425750849050604080838360005b83811015611779578181015183820152602001611761565b5050509201915083905060046020028201915b81548152602001906001019080831161178c57505094505050505060206040518083038186803b1580156117bf57600080fd5b505af41580156117d3573d6000803e3d6000fd5b505050506040513d60208110156117e957600080fd5b5051611830576040805162461bcd60e51b8152602060048201526011602482015270696e76616c6964207369676e617475726560781b604482015290519081900360640190fd5b600154604080516395184d3b60e01b815233600482015290517ff7084e4afb12f9dc96cee0f3c56fee76afe9598f22f0e3ce04f21b4c18cdc26792309287926001600160a01b03909216916395184d3b91602480820192602092909190829003018186803b1580156118a157600080fd5b505afa1580156118b5573d6000803e3d6000fd5b505050506040513d60208110156118cb57600080fd5b5051600a54604080516001600160a01b0380871682528416918101919091526060810182905260806020828101828152865192840192909252855160a084019187019080838360005b8381101561192c578181015183820152602001611914565b50505050905090810190601f1680156119595780820380516001836020036101000a031916815260200191505b509550505050505060405180910390a1825161197c906009906020860190611a05565b50505050565b61198a610d87565b61199357600080fd5b610a68816001600160a01b0381166119aa57600080fd5b600080546040516001600160a01b03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a3600080546001600160a01b0319166001600160a01b0392909216919091179055565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10611a4657805160ff1916838001178555611a73565b82800160010185558215611a73579182015b82811115611a73578251825591602001919060010190611a58565b50611a7f929150611b38565b5090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10611ac45782800160ff19823516178555611a73565b82800160010185558215611a73579182015b82811115611a73578235825591602001919060010190611ad6565b815481835581811115611b1557600083815260209020611b15918101908301611b38565b505050565b60405180604001604052806002906020820280388339509192915050565b6105ab91905b80821115611a7f5760008155600101611b3e56fe436f6e73656e73757320646f6573206e6f742062656c6f6e6720746f207468652067726f757021546865207573657220746f20626520696e766974656420697320616c72656164792061206d656d626572a265627a7a72315820a3c134ac558e2b397182543cdbf0273cc02e524b506ac0ebb9ad9c337114384b64736f6c634300050c0032"

// DeployGroup deploys a new Ethereum contract, binding an instance of Group to it.
func DeployGroup(auth *bind.TransactOpts, backend bind.ContractBackend, fileTribe common.Address, account common.Address, name string, ipfsHash []byte, verifyKey [4]*big.Int) (common.Address, *types.Transaction, *Group, error) {
	parsed, err := abi.JSON(strings.NewReader(GroupABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(GroupBin), backend, fileTribe, account, name, ipfsHash, verifyKey)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Group{GroupCaller: GroupCaller{contract: contract}, GroupTransactor: GroupTransactor{contract: contract}, GroupFilterer: GroupFilterer{contract: contract}}, nil
}

// Group is an auto generated Go binding around an Ethereum contract.
type Group struct {
	GroupCaller     // Read-only binding to the contract
	GroupTransactor // Write-only binding to the contract
	GroupFilterer   // Log filterer for contract events
}

// GroupCaller is an auto generated read-only Go binding around an Ethereum contract.
type GroupCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GroupTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GroupTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GroupFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GroupFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GroupSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GroupSession struct {
	Contract     *Group            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GroupCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GroupCallerSession struct {
	Contract *GroupCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// GroupTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GroupTransactorSession struct {
	Contract     *GroupTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GroupRaw is an auto generated low-level Go binding around an Ethereum contract.
type GroupRaw struct {
	Contract *Group // Generic contract binding to access the raw methods on
}

// GroupCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GroupCallerRaw struct {
	Contract *GroupCaller // Generic read-only contract binding to access the raw methods on
}

// GroupTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GroupTransactorRaw struct {
	Contract *GroupTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGroup creates a new instance of Group, bound to a specific deployed contract.
func NewGroup(address common.Address, backend bind.ContractBackend) (*Group, error) {
	contract, err := bindGroup(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Group{GroupCaller: GroupCaller{contract: contract}, GroupTransactor: GroupTransactor{contract: contract}, GroupFilterer: GroupFilterer{contract: contract}}, nil
}

// NewGroupCaller creates a new read-only instance of Group, bound to a specific deployed contract.
func NewGroupCaller(address common.Address, caller bind.ContractCaller) (*GroupCaller, error) {
	contract, err := bindGroup(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GroupCaller{contract: contract}, nil
}

// NewGroupTransactor creates a new write-only instance of Group, bound to a specific deployed contract.
func NewGroupTransactor(address common.Address, transactor bind.ContractTransactor) (*GroupTransactor, error) {
	contract, err := bindGroup(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GroupTransactor{contract: contract}, nil
}

// NewGroupFilterer creates a new log filterer instance of Group, bound to a specific deployed contract.
func NewGroupFilterer(address common.Address, filterer bind.ContractFilterer) (*GroupFilterer, error) {
	contract, err := bindGroup(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GroupFilterer{contract: contract}, nil
}

// bindGroup binds a generic wrapper to an already deployed contract.
func bindGroup(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(GroupABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Group *GroupRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Group.Contract.GroupCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Group *GroupRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Group.Contract.GroupTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Group *GroupRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Group.Contract.GroupTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Group *GroupCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Group.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Group *GroupTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Group.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Group *GroupTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Group.Contract.contract.Transact(opts, method, params...)
}

// G1 is a free data retrieval call binding the contract method 0x731f6fd8.
//
// Solidity: function _g1(uint256 ) constant returns(uint256)
func (_Group *GroupCaller) G1(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Group.contract.Call(opts, out, "_g1", arg0)
	return *ret0, err
}

// G1 is a free data retrieval call binding the contract method 0x731f6fd8.
//
// Solidity: function _g1(uint256 ) constant returns(uint256)
func (_Group *GroupSession) G1(arg0 *big.Int) (*big.Int, error) {
	return _Group.Contract.G1(&_Group.CallOpts, arg0)
}

// G1 is a free data retrieval call binding the contract method 0x731f6fd8.
//
// Solidity: function _g1(uint256 ) constant returns(uint256)
func (_Group *GroupCallerSession) G1(arg0 *big.Int) (*big.Int, error) {
	return _Group.Contract.G1(&_Group.CallOpts, arg0)
}

// G2 is a free data retrieval call binding the contract method 0xd53d6404.
//
// Solidity: function _g2(uint256 ) constant returns(uint256)
func (_Group *GroupCaller) G2(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Group.contract.Call(opts, out, "_g2", arg0)
	return *ret0, err
}

// G2 is a free data retrieval call binding the contract method 0xd53d6404.
//
// Solidity: function _g2(uint256 ) constant returns(uint256)
func (_Group *GroupSession) G2(arg0 *big.Int) (*big.Int, error) {
	return _Group.Contract.G2(&_Group.CallOpts, arg0)
}

// G2 is a free data retrieval call binding the contract method 0xd53d6404.
//
// Solidity: function _g2(uint256 ) constant returns(uint256)
func (_Group *GroupCallerSession) G2(arg0 *big.Int) (*big.Int, error) {
	return _Group.Contract.G2(&_Group.CallOpts, arg0)
}

// IpfsHash is a free data retrieval call binding the contract method 0xcc2867d7.
//
// Solidity: function _ipfsHash() constant returns(bytes)
func (_Group *GroupCaller) IpfsHash(opts *bind.CallOpts) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _Group.contract.Call(opts, out, "_ipfsHash")
	return *ret0, err
}

// IpfsHash is a free data retrieval call binding the contract method 0xcc2867d7.
//
// Solidity: function _ipfsHash() constant returns(bytes)
func (_Group *GroupSession) IpfsHash() ([]byte, error) {
	return _Group.Contract.IpfsHash(&_Group.CallOpts)
}

// IpfsHash is a free data retrieval call binding the contract method 0xcc2867d7.
//
// Solidity: function _ipfsHash() constant returns(bytes)
func (_Group *GroupCallerSession) IpfsHash() ([]byte, error) {
	return _Group.Contract.IpfsHash(&_Group.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0xd28d8852.
//
// Solidity: function _name() constant returns(string)
func (_Group *GroupCaller) Name(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Group.contract.Call(opts, out, "_name")
	return *ret0, err
}

// Name is a free data retrieval call binding the contract method 0xd28d8852.
//
// Solidity: function _name() constant returns(string)
func (_Group *GroupSession) Name() (string, error) {
	return _Group.Contract.Name(&_Group.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0xd28d8852.
//
// Solidity: function _name() constant returns(string)
func (_Group *GroupCallerSession) Name() (string, error) {
	return _Group.Contract.Name(&_Group.CallOpts)
}

// VerifyKey is a free data retrieval call binding the contract method 0x12f8ad05.
//
// Solidity: function _verifyKey(uint256 ) constant returns(uint256)
func (_Group *GroupCaller) VerifyKey(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Group.contract.Call(opts, out, "_verifyKey", arg0)
	return *ret0, err
}

// VerifyKey is a free data retrieval call binding the contract method 0x12f8ad05.
//
// Solidity: function _verifyKey(uint256 ) constant returns(uint256)
func (_Group *GroupSession) VerifyKey(arg0 *big.Int) (*big.Int, error) {
	return _Group.Contract.VerifyKey(&_Group.CallOpts, arg0)
}

// VerifyKey is a free data retrieval call binding the contract method 0x12f8ad05.
//
// Solidity: function _verifyKey(uint256 ) constant returns(uint256)
func (_Group *GroupCallerSession) VerifyKey(arg0 *big.Int) (*big.Int, error) {
	return _Group.Contract.VerifyKey(&_Group.CallOpts, arg0)
}

// GetConsensus is a free data retrieval call binding the contract method 0xe8f738e1.
//
// Solidity: function getConsensus(address owner) constant returns(address)
func (_Group *GroupCaller) GetConsensus(opts *bind.CallOpts, owner common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Group.contract.Call(opts, out, "getConsensus", owner)
	return *ret0, err
}

// GetConsensus is a free data retrieval call binding the contract method 0xe8f738e1.
//
// Solidity: function getConsensus(address owner) constant returns(address)
func (_Group *GroupSession) GetConsensus(owner common.Address) (common.Address, error) {
	return _Group.Contract.GetConsensus(&_Group.CallOpts, owner)
}

// GetConsensus is a free data retrieval call binding the contract method 0xe8f738e1.
//
// Solidity: function getConsensus(address owner) constant returns(address)
func (_Group *GroupCallerSession) GetConsensus(owner common.Address) (common.Address, error) {
	return _Group.Contract.GetConsensus(&_Group.CallOpts, owner)
}

// IsMember is a free data retrieval call binding the contract method 0xa230c524.
//
// Solidity: function isMember(address owner) constant returns(bool)
func (_Group *GroupCaller) IsMember(opts *bind.CallOpts, owner common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Group.contract.Call(opts, out, "isMember", owner)
	return *ret0, err
}

// IsMember is a free data retrieval call binding the contract method 0xa230c524.
//
// Solidity: function isMember(address owner) constant returns(bool)
func (_Group *GroupSession) IsMember(owner common.Address) (bool, error) {
	return _Group.Contract.IsMember(&_Group.CallOpts, owner)
}

// IsMember is a free data retrieval call binding the contract method 0xa230c524.
//
// Solidity: function isMember(address owner) constant returns(bool)
func (_Group *GroupCallerSession) IsMember(owner common.Address) (bool, error) {
	return _Group.Contract.IsMember(&_Group.CallOpts, owner)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Group *GroupCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Group.contract.Call(opts, out, "isOwner")
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Group *GroupSession) IsOwner() (bool, error) {
	return _Group.Contract.IsOwner(&_Group.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Group *GroupCallerSession) IsOwner() (bool, error) {
	return _Group.Contract.IsOwner(&_Group.CallOpts)
}

// MemberOwners is a free data retrieval call binding the contract method 0xe7acfda1.
//
// Solidity: function memberOwners() constant returns(address[])
func (_Group *GroupCaller) MemberOwners(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _Group.contract.Call(opts, out, "memberOwners")
	return *ret0, err
}

// MemberOwners is a free data retrieval call binding the contract method 0xe7acfda1.
//
// Solidity: function memberOwners() constant returns(address[])
func (_Group *GroupSession) MemberOwners() ([]common.Address, error) {
	return _Group.Contract.MemberOwners(&_Group.CallOpts)
}

// MemberOwners is a free data retrieval call binding the contract method 0xe7acfda1.
//
// Solidity: function memberOwners() constant returns(address[])
func (_Group *GroupCallerSession) MemberOwners() ([]common.Address, error) {
	return _Group.Contract.MemberOwners(&_Group.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Group *GroupCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Group.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Group *GroupSession) Owner() (common.Address, error) {
	return _Group.Contract.Owner(&_Group.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Group *GroupCallerSession) Owner() (common.Address, error) {
	return _Group.Contract.Owner(&_Group.CallOpts)
}

// Threshold is a free data retrieval call binding the contract method 0x42cde4e8.
//
// Solidity: function threshold() constant returns(uint256)
func (_Group *GroupCaller) Threshold(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Group.contract.Call(opts, out, "threshold")
	return *ret0, err
}

// Threshold is a free data retrieval call binding the contract method 0x42cde4e8.
//
// Solidity: function threshold() constant returns(uint256)
func (_Group *GroupSession) Threshold() (*big.Int, error) {
	return _Group.Contract.Threshold(&_Group.CallOpts)
}

// Threshold is a free data retrieval call binding the contract method 0x42cde4e8.
//
// Solidity: function threshold() constant returns(uint256)
func (_Group *GroupCallerSession) Threshold() (*big.Int, error) {
	return _Group.Contract.Threshold(&_Group.CallOpts)
}

// Commit is a paid mutator transaction binding the contract method 0x66fd3cd8.
//
// Solidity: function commit(bytes newIpfsHash) returns()
func (_Group *GroupTransactor) Commit(opts *bind.TransactOpts, newIpfsHash []byte) (*types.Transaction, error) {
	return _Group.contract.Transact(opts, "commit", newIpfsHash)
}

// Commit is a paid mutator transaction binding the contract method 0x66fd3cd8.
//
// Solidity: function commit(bytes newIpfsHash) returns()
func (_Group *GroupSession) Commit(newIpfsHash []byte) (*types.Transaction, error) {
	return _Group.Contract.Commit(&_Group.TransactOpts, newIpfsHash)
}

// Commit is a paid mutator transaction binding the contract method 0x66fd3cd8.
//
// Solidity: function commit(bytes newIpfsHash) returns()
func (_Group *GroupTransactorSession) Commit(newIpfsHash []byte) (*types.Transaction, error) {
	return _Group.Contract.Commit(&_Group.TransactOpts, newIpfsHash)
}

// CommitWithGroupSig is a paid mutator transaction binding the contract method 0xebe1c230.
//
// Solidity: function commitWithGroupSig(bytes newIpfsHash, uint256[2] sig) returns()
func (_Group *GroupTransactor) CommitWithGroupSig(opts *bind.TransactOpts, newIpfsHash []byte, sig [2]*big.Int) (*types.Transaction, error) {
	return _Group.contract.Transact(opts, "commitWithGroupSig", newIpfsHash, sig)
}

// CommitWithGroupSig is a paid mutator transaction binding the contract method 0xebe1c230.
//
// Solidity: function commitWithGroupSig(bytes newIpfsHash, uint256[2] sig) returns()
func (_Group *GroupSession) CommitWithGroupSig(newIpfsHash []byte, sig [2]*big.Int) (*types.Transaction, error) {
	return _Group.Contract.CommitWithGroupSig(&_Group.TransactOpts, newIpfsHash, sig)
}

// CommitWithGroupSig is a paid mutator transaction binding the contract method 0xebe1c230.
//
// Solidity: function commitWithGroupSig(bytes newIpfsHash, uint256[2] sig) returns()
func (_Group *GroupTransactorSession) CommitWithGroupSig(newIpfsHash []byte, sig [2]*big.Int) (*types.Transaction, error) {
	return _Group.Contract.CommitWithGroupSig(&_Group.TransactOpts, newIpfsHash, sig)
}

// Decline is a paid mutator transaction binding the contract method 0xab040107.
//
// Solidity: function decline() returns()
func (_Group *GroupTransactor) Decline(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Group.contract.Transact(opts, "decline")
}

// Decline is a paid mutator transaction binding the contract method 0xab040107.
//
// Solidity: function decline() returns()
func (_Group *GroupSession) Decline() (*types.Transaction, error) {
	return _Group.Contract.Decline(&_Group.TransactOpts)
}

// Decline is a paid mutator transaction binding the contract method 0xab040107.
//
// Solidity: function decline() returns()
func (_Group *GroupTransactorSession) Decline() (*types.Transaction, error) {
	return _Group.Contract.Decline(&_Group.TransactOpts)
}

// Invite is a paid mutator transaction binding the contract method 0x4b77c468.
//
// Solidity: function invite(address account) returns()
func (_Group *GroupTransactor) Invite(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _Group.contract.Transact(opts, "invite", account)
}

// Invite is a paid mutator transaction binding the contract method 0x4b77c468.
//
// Solidity: function invite(address account) returns()
func (_Group *GroupSession) Invite(account common.Address) (*types.Transaction, error) {
	return _Group.Contract.Invite(&_Group.TransactOpts, account)
}

// Invite is a paid mutator transaction binding the contract method 0x4b77c468.
//
// Solidity: function invite(address account) returns()
func (_Group *GroupTransactorSession) Invite(account common.Address) (*types.Transaction, error) {
	return _Group.Contract.Invite(&_Group.TransactOpts, account)
}

// Join is a paid mutator transaction binding the contract method 0xb688a363.
//
// Solidity: function join() returns()
func (_Group *GroupTransactor) Join(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Group.contract.Transact(opts, "join")
}

// Join is a paid mutator transaction binding the contract method 0xb688a363.
//
// Solidity: function join() returns()
func (_Group *GroupSession) Join() (*types.Transaction, error) {
	return _Group.Contract.Join(&_Group.TransactOpts)
}

// Join is a paid mutator transaction binding the contract method 0xb688a363.
//
// Solidity: function join() returns()
func (_Group *GroupTransactorSession) Join() (*types.Transaction, error) {
	return _Group.Contract.Join(&_Group.TransactOpts)
}

// Kick is a paid mutator transaction binding the contract method 0x96c55175.
//
// Solidity: function kick(address memberOwner) returns()
func (_Group *GroupTransactor) Kick(opts *bind.TransactOpts, memberOwner common.Address) (*types.Transaction, error) {
	return _Group.contract.Transact(opts, "kick", memberOwner)
}

// Kick is a paid mutator transaction binding the contract method 0x96c55175.
//
// Solidity: function kick(address memberOwner) returns()
func (_Group *GroupSession) Kick(memberOwner common.Address) (*types.Transaction, error) {
	return _Group.Contract.Kick(&_Group.TransactOpts, memberOwner)
}

// Kick is a paid mutator transaction binding the contract method 0x96c55175.
//
// Solidity: function kick(address memberOwner) returns()
func (_Group *GroupTransactorSession) Kick(memberOwner common.Address) (*types.Transaction, error) {
	return _Group.Contract.Kick(&_Group.TransactOpts, memberOwner)
}

// Leave is a paid mutator transaction binding the contract method 0xd66d9e19.
//
// Solidity: function leave() returns()
func (_Group *GroupTransactor) Leave(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Group.contract.Transact(opts, "leave")
}

// Leave is a paid mutator transaction binding the contract method 0xd66d9e19.
//
// Solidity: function leave() returns()
func (_Group *GroupSession) Leave() (*types.Transaction, error) {
	return _Group.Contract.Leave(&_Group.TransactOpts)
}

// Leave is a paid mutator transaction binding the contract method 0xd66d9e19.
//
// Solidity: function leave() returns()
func (_Group *GroupTransactorSession) Leave() (*types.Transaction, error) {
	return _Group.Contract.Leave(&_Group.TransactOpts)
}

// OnChangeIpfsHashConsensus is a paid mutator transaction binding the contract method 0x8379088e.
//
// Solidity: function onChangeIpfsHashConsensus(bytes payload) returns()
func (_Group *GroupTransactor) OnChangeIpfsHashConsensus(opts *bind.TransactOpts, payload []byte) (*types.Transaction, error) {
	return _Group.contract.Transact(opts, "onChangeIpfsHashConsensus", payload)
}

// OnChangeIpfsHashConsensus is a paid mutator transaction binding the contract method 0x8379088e.
//
// Solidity: function onChangeIpfsHashConsensus(bytes payload) returns()
func (_Group *GroupSession) OnChangeIpfsHashConsensus(payload []byte) (*types.Transaction, error) {
	return _Group.Contract.OnChangeIpfsHashConsensus(&_Group.TransactOpts, payload)
}

// OnChangeIpfsHashConsensus is a paid mutator transaction binding the contract method 0x8379088e.
//
// Solidity: function onChangeIpfsHashConsensus(bytes payload) returns()
func (_Group *GroupTransactorSession) OnChangeIpfsHashConsensus(payload []byte) (*types.Transaction, error) {
	return _Group.Contract.OnChangeIpfsHashConsensus(&_Group.TransactOpts, payload)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Group *GroupTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Group.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Group *GroupSession) RenounceOwnership() (*types.Transaction, error) {
	return _Group.Contract.RenounceOwnership(&_Group.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Group *GroupTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Group.Contract.RenounceOwnership(&_Group.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Group *GroupTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Group.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Group *GroupSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Group.Contract.TransferOwnership(&_Group.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Group *GroupTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Group.Contract.TransferOwnership(&_Group.TransactOpts, newOwner)
}

// GroupDebugIterator is returned from FilterDebug and is used to iterate over the raw logs and unpacked data for Debug events raised by the Group contract.
type GroupDebugIterator struct {
	Event *GroupDebug // Event containing the contract specifics and raw log

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
func (it *GroupDebugIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupDebug)
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
		it.Event = new(GroupDebug)
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
func (it *GroupDebugIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupDebugIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupDebug represents a Debug event raised by the Group contract.
type GroupDebug struct {
	Msg *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterDebug is a free log retrieval operation binding the contract event 0x8a36f5a234186d446e36a7df36ace663a05a580d9bea2dd899c6dd76a075d5fa.
//
// Solidity: event Debug(uint256 msg)
func (_Group *GroupFilterer) FilterDebug(opts *bind.FilterOpts) (*GroupDebugIterator, error) {

	logs, sub, err := _Group.contract.FilterLogs(opts, "Debug")
	if err != nil {
		return nil, err
	}
	return &GroupDebugIterator{contract: _Group.contract, event: "Debug", logs: logs, sub: sub}, nil
}

// WatchDebug is a free log subscription operation binding the contract event 0x8a36f5a234186d446e36a7df36ace663a05a580d9bea2dd899c6dd76a075d5fa.
//
// Solidity: event Debug(uint256 msg)
func (_Group *GroupFilterer) WatchDebug(opts *bind.WatchOpts, sink chan<- *GroupDebug) (event.Subscription, error) {

	logs, sub, err := _Group.contract.WatchLogs(opts, "Debug")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupDebug)
				if err := _Group.contract.UnpackLog(event, "Debug", log); err != nil {
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
// Solidity: event Debug(uint256 msg)
func (_Group *GroupFilterer) ParseDebug(log types.Log) (*GroupDebug, error) {
	event := new(GroupDebug)
	if err := _Group.contract.UnpackLog(event, "Debug", log); err != nil {
		return nil, err
	}
	return event, nil
}

// GroupGroupUpdateIpfsHashIterator is returned from FilterGroupUpdateIpfsHash and is used to iterate over the raw logs and unpacked data for GroupUpdateIpfsHash events raised by the Group contract.
type GroupGroupUpdateIpfsHashIterator struct {
	Event *GroupGroupUpdateIpfsHash // Event containing the contract specifics and raw log

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
func (it *GroupGroupUpdateIpfsHashIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupGroupUpdateIpfsHash)
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
		it.Event = new(GroupGroupUpdateIpfsHash)
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
func (it *GroupGroupUpdateIpfsHashIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupGroupUpdateIpfsHashIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupGroupUpdateIpfsHash represents a GroupUpdateIpfsHash event raised by the Group contract.
type GroupGroupUpdateIpfsHash struct {
	GroupId  [32]byte
	IpfsHash []byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterGroupUpdateIpfsHash is a free log retrieval operation binding the contract event 0x54f7f672bf7abc702d283d4028997d3c86d75e43fdbd6c425b50d1d5cdb2229f.
//
// Solidity: event GroupUpdateIpfsHash(bytes32 groupId, bytes ipfsHash)
func (_Group *GroupFilterer) FilterGroupUpdateIpfsHash(opts *bind.FilterOpts) (*GroupGroupUpdateIpfsHashIterator, error) {

	logs, sub, err := _Group.contract.FilterLogs(opts, "GroupUpdateIpfsHash")
	if err != nil {
		return nil, err
	}
	return &GroupGroupUpdateIpfsHashIterator{contract: _Group.contract, event: "GroupUpdateIpfsHash", logs: logs, sub: sub}, nil
}

// WatchGroupUpdateIpfsHash is a free log subscription operation binding the contract event 0x54f7f672bf7abc702d283d4028997d3c86d75e43fdbd6c425b50d1d5cdb2229f.
//
// Solidity: event GroupUpdateIpfsHash(bytes32 groupId, bytes ipfsHash)
func (_Group *GroupFilterer) WatchGroupUpdateIpfsHash(opts *bind.WatchOpts, sink chan<- *GroupGroupUpdateIpfsHash) (event.Subscription, error) {

	logs, sub, err := _Group.contract.WatchLogs(opts, "GroupUpdateIpfsHash")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupGroupUpdateIpfsHash)
				if err := _Group.contract.UnpackLog(event, "GroupUpdateIpfsHash", log); err != nil {
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

// ParseGroupUpdateIpfsHash is a log parse operation binding the contract event 0x54f7f672bf7abc702d283d4028997d3c86d75e43fdbd6c425b50d1d5cdb2229f.
//
// Solidity: event GroupUpdateIpfsHash(bytes32 groupId, bytes ipfsHash)
func (_Group *GroupFilterer) ParseGroupUpdateIpfsHash(log types.Log) (*GroupGroupUpdateIpfsHash, error) {
	event := new(GroupGroupUpdateIpfsHash)
	if err := _Group.contract.UnpackLog(event, "GroupUpdateIpfsHash", log); err != nil {
		return nil, err
	}
	return event, nil
}

// GroupInvitationAcceptedIterator is returned from FilterInvitationAccepted and is used to iterate over the raw logs and unpacked data for InvitationAccepted events raised by the Group contract.
type GroupInvitationAcceptedIterator struct {
	Event *GroupInvitationAccepted // Event containing the contract specifics and raw log

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
func (it *GroupInvitationAcceptedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupInvitationAccepted)
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
		it.Event = new(GroupInvitationAccepted)
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
func (it *GroupInvitationAcceptedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupInvitationAcceptedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupInvitationAccepted represents a InvitationAccepted event raised by the Group contract.
type GroupInvitationAccepted struct {
	GroupAddress common.Address
	Account      common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterInvitationAccepted is a free log retrieval operation binding the contract event 0x4d7c243e154e530692e62f8539db65779f5cb85d58831956361697addede5adb.
//
// Solidity: event InvitationAccepted(address groupAddress, address account)
func (_Group *GroupFilterer) FilterInvitationAccepted(opts *bind.FilterOpts) (*GroupInvitationAcceptedIterator, error) {

	logs, sub, err := _Group.contract.FilterLogs(opts, "InvitationAccepted")
	if err != nil {
		return nil, err
	}
	return &GroupInvitationAcceptedIterator{contract: _Group.contract, event: "InvitationAccepted", logs: logs, sub: sub}, nil
}

// WatchInvitationAccepted is a free log subscription operation binding the contract event 0x4d7c243e154e530692e62f8539db65779f5cb85d58831956361697addede5adb.
//
// Solidity: event InvitationAccepted(address groupAddress, address account)
func (_Group *GroupFilterer) WatchInvitationAccepted(opts *bind.WatchOpts, sink chan<- *GroupInvitationAccepted) (event.Subscription, error) {

	logs, sub, err := _Group.contract.WatchLogs(opts, "InvitationAccepted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupInvitationAccepted)
				if err := _Group.contract.UnpackLog(event, "InvitationAccepted", log); err != nil {
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
// Solidity: event InvitationAccepted(address groupAddress, address account)
func (_Group *GroupFilterer) ParseInvitationAccepted(log types.Log) (*GroupInvitationAccepted, error) {
	event := new(GroupInvitationAccepted)
	if err := _Group.contract.UnpackLog(event, "InvitationAccepted", log); err != nil {
		return nil, err
	}
	return event, nil
}

// GroupInvitationDeclinedIterator is returned from FilterInvitationDeclined and is used to iterate over the raw logs and unpacked data for InvitationDeclined events raised by the Group contract.
type GroupInvitationDeclinedIterator struct {
	Event *GroupInvitationDeclined // Event containing the contract specifics and raw log

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
func (it *GroupInvitationDeclinedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupInvitationDeclined)
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
		it.Event = new(GroupInvitationDeclined)
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
func (it *GroupInvitationDeclinedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupInvitationDeclinedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupInvitationDeclined represents a InvitationDeclined event raised by the Group contract.
type GroupInvitationDeclined struct {
	Group   common.Address
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInvitationDeclined is a free log retrieval operation binding the contract event 0x7ea317cacaf634632d88c1fd20864b1f54f609bb063c678280aac6312715edef.
//
// Solidity: event InvitationDeclined(address group, address account)
func (_Group *GroupFilterer) FilterInvitationDeclined(opts *bind.FilterOpts) (*GroupInvitationDeclinedIterator, error) {

	logs, sub, err := _Group.contract.FilterLogs(opts, "InvitationDeclined")
	if err != nil {
		return nil, err
	}
	return &GroupInvitationDeclinedIterator{contract: _Group.contract, event: "InvitationDeclined", logs: logs, sub: sub}, nil
}

// WatchInvitationDeclined is a free log subscription operation binding the contract event 0x7ea317cacaf634632d88c1fd20864b1f54f609bb063c678280aac6312715edef.
//
// Solidity: event InvitationDeclined(address group, address account)
func (_Group *GroupFilterer) WatchInvitationDeclined(opts *bind.WatchOpts, sink chan<- *GroupInvitationDeclined) (event.Subscription, error) {

	logs, sub, err := _Group.contract.WatchLogs(opts, "InvitationDeclined")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupInvitationDeclined)
				if err := _Group.contract.UnpackLog(event, "InvitationDeclined", log); err != nil {
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
// Solidity: event InvitationDeclined(address group, address account)
func (_Group *GroupFilterer) ParseInvitationDeclined(log types.Log) (*GroupInvitationDeclined, error) {
	event := new(GroupInvitationDeclined)
	if err := _Group.contract.UnpackLog(event, "InvitationDeclined", log); err != nil {
		return nil, err
	}
	return event, nil
}

// GroupInvitationSentIterator is returned from FilterInvitationSent and is used to iterate over the raw logs and unpacked data for InvitationSent events raised by the Group contract.
type GroupInvitationSentIterator struct {
	Event *GroupInvitationSent // Event containing the contract specifics and raw log

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
func (it *GroupInvitationSentIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupInvitationSent)
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
		it.Event = new(GroupInvitationSent)
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
func (it *GroupInvitationSentIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupInvitationSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupInvitationSent represents a InvitationSent event raised by the Group contract.
type GroupInvitationSent struct {
	Group   common.Address
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInvitationSent is a free log retrieval operation binding the contract event 0x4509853340a1799d2c49c6cb19f9988fc2fa7f8c37aa55a7d865e19a260cfdc9.
//
// Solidity: event InvitationSent(address group, address account)
func (_Group *GroupFilterer) FilterInvitationSent(opts *bind.FilterOpts) (*GroupInvitationSentIterator, error) {

	logs, sub, err := _Group.contract.FilterLogs(opts, "InvitationSent")
	if err != nil {
		return nil, err
	}
	return &GroupInvitationSentIterator{contract: _Group.contract, event: "InvitationSent", logs: logs, sub: sub}, nil
}

// WatchInvitationSent is a free log subscription operation binding the contract event 0x4509853340a1799d2c49c6cb19f9988fc2fa7f8c37aa55a7d865e19a260cfdc9.
//
// Solidity: event InvitationSent(address group, address account)
func (_Group *GroupFilterer) WatchInvitationSent(opts *bind.WatchOpts, sink chan<- *GroupInvitationSent) (event.Subscription, error) {

	logs, sub, err := _Group.contract.WatchLogs(opts, "InvitationSent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupInvitationSent)
				if err := _Group.contract.UnpackLog(event, "InvitationSent", log); err != nil {
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

// ParseInvitationSent is a log parse operation binding the contract event 0x4509853340a1799d2c49c6cb19f9988fc2fa7f8c37aa55a7d865e19a260cfdc9.
//
// Solidity: event InvitationSent(address group, address account)
func (_Group *GroupFilterer) ParseInvitationSent(log types.Log) (*GroupInvitationSent, error) {
	event := new(GroupInvitationSent)
	if err := _Group.contract.UnpackLog(event, "InvitationSent", log); err != nil {
		return nil, err
	}
	return event, nil
}

// GroupIpfsHashChangedIterator is returned from FilterIpfsHashChanged and is used to iterate over the raw logs and unpacked data for IpfsHashChanged events raised by the Group contract.
type GroupIpfsHashChangedIterator struct {
	Event *GroupIpfsHashChanged // Event containing the contract specifics and raw log

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
func (it *GroupIpfsHashChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupIpfsHashChanged)
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
		it.Event = new(GroupIpfsHashChanged)
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
func (it *GroupIpfsHashChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupIpfsHashChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupIpfsHashChanged represents a IpfsHashChanged event raised by the Group contract.
type GroupIpfsHashChanged struct {
	Group    common.Address
	IpfsHash []byte
	Proposer common.Address
	Id       *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterIpfsHashChanged is a free log retrieval operation binding the contract event 0xf7084e4afb12f9dc96cee0f3c56fee76afe9598f22f0e3ce04f21b4c18cdc267.
//
// Solidity: event IpfsHashChanged(address group, bytes ipfsHash, address proposer, uint256 id)
func (_Group *GroupFilterer) FilterIpfsHashChanged(opts *bind.FilterOpts) (*GroupIpfsHashChangedIterator, error) {

	logs, sub, err := _Group.contract.FilterLogs(opts, "IpfsHashChanged")
	if err != nil {
		return nil, err
	}
	return &GroupIpfsHashChangedIterator{contract: _Group.contract, event: "IpfsHashChanged", logs: logs, sub: sub}, nil
}

// WatchIpfsHashChanged is a free log subscription operation binding the contract event 0xf7084e4afb12f9dc96cee0f3c56fee76afe9598f22f0e3ce04f21b4c18cdc267.
//
// Solidity: event IpfsHashChanged(address group, bytes ipfsHash, address proposer, uint256 id)
func (_Group *GroupFilterer) WatchIpfsHashChanged(opts *bind.WatchOpts, sink chan<- *GroupIpfsHashChanged) (event.Subscription, error) {

	logs, sub, err := _Group.contract.WatchLogs(opts, "IpfsHashChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupIpfsHashChanged)
				if err := _Group.contract.UnpackLog(event, "IpfsHashChanged", log); err != nil {
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

// ParseIpfsHashChanged is a log parse operation binding the contract event 0xf7084e4afb12f9dc96cee0f3c56fee76afe9598f22f0e3ce04f21b4c18cdc267.
//
// Solidity: event IpfsHashChanged(address group, bytes ipfsHash, address proposer, uint256 id)
func (_Group *GroupFilterer) ParseIpfsHashChanged(log types.Log) (*GroupIpfsHashChanged, error) {
	event := new(GroupIpfsHashChanged)
	if err := _Group.contract.UnpackLog(event, "IpfsHashChanged", log); err != nil {
		return nil, err
	}
	return event, nil
}

// GroupMemberLeftIterator is returned from FilterMemberLeft and is used to iterate over the raw logs and unpacked data for MemberLeft events raised by the Group contract.
type GroupMemberLeftIterator struct {
	Event *GroupMemberLeft // Event containing the contract specifics and raw log

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
func (it *GroupMemberLeftIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupMemberLeft)
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
		it.Event = new(GroupMemberLeft)
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
func (it *GroupMemberLeftIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupMemberLeftIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupMemberLeft represents a MemberLeft event raised by the Group contract.
type GroupMemberLeft struct {
	Group   common.Address
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMemberLeft is a free log retrieval operation binding the contract event 0xa0ae7458a482bd2780bed4edd090c3f9e9e3e600171cdc4709da09e5bd233e18.
//
// Solidity: event MemberLeft(address group, address account)
func (_Group *GroupFilterer) FilterMemberLeft(opts *bind.FilterOpts) (*GroupMemberLeftIterator, error) {

	logs, sub, err := _Group.contract.FilterLogs(opts, "MemberLeft")
	if err != nil {
		return nil, err
	}
	return &GroupMemberLeftIterator{contract: _Group.contract, event: "MemberLeft", logs: logs, sub: sub}, nil
}

// WatchMemberLeft is a free log subscription operation binding the contract event 0xa0ae7458a482bd2780bed4edd090c3f9e9e3e600171cdc4709da09e5bd233e18.
//
// Solidity: event MemberLeft(address group, address account)
func (_Group *GroupFilterer) WatchMemberLeft(opts *bind.WatchOpts, sink chan<- *GroupMemberLeft) (event.Subscription, error) {

	logs, sub, err := _Group.contract.WatchLogs(opts, "MemberLeft")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupMemberLeft)
				if err := _Group.contract.UnpackLog(event, "MemberLeft", log); err != nil {
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

// ParseMemberLeft is a log parse operation binding the contract event 0xa0ae7458a482bd2780bed4edd090c3f9e9e3e600171cdc4709da09e5bd233e18.
//
// Solidity: event MemberLeft(address group, address account)
func (_Group *GroupFilterer) ParseMemberLeft(log types.Log) (*GroupMemberLeft, error) {
	event := new(GroupMemberLeft)
	if err := _Group.contract.UnpackLog(event, "MemberLeft", log); err != nil {
		return nil, err
	}
	return event, nil
}

// GroupNewConsensusIterator is returned from FilterNewConsensus and is used to iterate over the raw logs and unpacked data for NewConsensus events raised by the Group contract.
type GroupNewConsensusIterator struct {
	Event *GroupNewConsensus // Event containing the contract specifics and raw log

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
func (it *GroupNewConsensusIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupNewConsensus)
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
		it.Event = new(GroupNewConsensus)
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
func (it *GroupNewConsensusIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupNewConsensusIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupNewConsensus represents a NewConsensus event raised by the Group contract.
type GroupNewConsensus struct {
	Group     common.Address
	Consensus common.Address
	Id        *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterNewConsensus is a free log retrieval operation binding the contract event 0xe407c1bec7341bd0a3dcb85da18f4e0b5126fcc79738f8f1dde296455e9c35aa.
//
// Solidity: event NewConsensus(address group, address consensus, uint256 id)
func (_Group *GroupFilterer) FilterNewConsensus(opts *bind.FilterOpts) (*GroupNewConsensusIterator, error) {

	logs, sub, err := _Group.contract.FilterLogs(opts, "NewConsensus")
	if err != nil {
		return nil, err
	}
	return &GroupNewConsensusIterator{contract: _Group.contract, event: "NewConsensus", logs: logs, sub: sub}, nil
}

// WatchNewConsensus is a free log subscription operation binding the contract event 0xe407c1bec7341bd0a3dcb85da18f4e0b5126fcc79738f8f1dde296455e9c35aa.
//
// Solidity: event NewConsensus(address group, address consensus, uint256 id)
func (_Group *GroupFilterer) WatchNewConsensus(opts *bind.WatchOpts, sink chan<- *GroupNewConsensus) (event.Subscription, error) {

	logs, sub, err := _Group.contract.WatchLogs(opts, "NewConsensus")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupNewConsensus)
				if err := _Group.contract.UnpackLog(event, "NewConsensus", log); err != nil {
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

// ParseNewConsensus is a log parse operation binding the contract event 0xe407c1bec7341bd0a3dcb85da18f4e0b5126fcc79738f8f1dde296455e9c35aa.
//
// Solidity: event NewConsensus(address group, address consensus, uint256 id)
func (_Group *GroupFilterer) ParseNewConsensus(log types.Log) (*GroupNewConsensus, error) {
	event := new(GroupNewConsensus)
	if err := _Group.contract.UnpackLog(event, "NewConsensus", log); err != nil {
		return nil, err
	}
	return event, nil
}

// GroupOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Group contract.
type GroupOwnershipTransferredIterator struct {
	Event *GroupOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *GroupOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupOwnershipTransferred)
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
		it.Event = new(GroupOwnershipTransferred)
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
func (it *GroupOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupOwnershipTransferred represents a OwnershipTransferred event raised by the Group contract.
type GroupOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Group *GroupFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*GroupOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Group.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &GroupOwnershipTransferredIterator{contract: _Group.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Group *GroupFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *GroupOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Group.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupOwnershipTransferred)
				if err := _Group.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Group *GroupFilterer) ParseOwnershipTransferred(log types.Log) (*GroupOwnershipTransferred, error) {
	event := new(GroupOwnershipTransferred)
	if err := _Group.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}
