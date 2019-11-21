// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package GroupFactory

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

// GroupFactoryABI is the input ABI used to generate the binding from.
const GroupFactoryABI = "[{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"contractIAccount\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"create\",\"outputs\":[{\"internalType\":\"contractIGroup\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"contractIFileTribeDApp\",\"name\":\"parent\",\"type\":\"address\"}],\"name\":\"setParent\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// GroupFactoryBin is the compiled bytecode used for deploying new contracts.
var GroupFactoryBin = "0x608060405234801561001057600080fd5b50600080546001600160a01b0319163390811780835560405191926001600160a01b0391909116917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908290a3506126558061006d6000396000f3fe60806040523480156200001157600080fd5b50600436106200006a5760003560e01c80631499c592146200006f578063715018a6146200009a5780638da5cb5b14620000a45780638f32d59b14620000ca578063a15ab08d14620000e8578063f2fde38b146200016e575b600080fd5b62000098600480360360208110156200008757600080fd5b50356001600160a01b031662000197565b005b62000098620001cd565b620000ae6200022b565b604080516001600160a01b039092168252519081900360200190f35b620000d46200023a565b604080519115158252519081900360200190f35b620000ae600480360360408110156200010057600080fd5b6001600160a01b0382351691908101906040810160208201356401000000008111156200012c57600080fd5b8201836020820111156200013f57600080fd5b803590602001918460018302840111640100000000831117156200016257600080fd5b5090925090506200024b565b62000098600480360360208110156200018657600080fd5b50356001600160a01b0316620003a2565b620001a16200023a565b620001ab57600080fd5b600180546001600160a01b0319166001600160a01b0392909216919091179055565b620001d76200023a565b620001e157600080fd5b600080546040516001600160a01b03909116907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a3600080546001600160a01b0319169055565b6000546001600160a01b031690565b6000546001600160a01b0316331490565b60015460408051600080825260a082018352602082018181528284018290526060830182905260808301829052925190936001600160a01b03169287928792879290620002989062000433565b6001600160a01b038088168252861660208201526040810160608201608080840190859080838360005b83811015620002dc578181015183820152602001620002c2565b5050505090500183810383528787828181526020019250808284376000838201819052601f909101601f191690920185810384528751815287516020918201939189019250908190849084905b838110156200034357818101518382015260200162000329565b50505050905090810190601f168015620003715780820380516001836020036101000a031916815260200191505b5098505050505050505050604051809103906000f08015801562000399573d6000803e3d6000fd5b50949350505050565b620003ac6200023a565b620003b657600080fd5b620003c181620003c4565b50565b6001600160a01b038116620003d857600080fd5b600080546040516001600160a01b03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a3600080546001600160a01b0319166001600160a01b0392909216919091179055565b6121df80620004428339019056fe60c060405260016080908152600260a08190526200002091600c91620004ec565b5060405180608001604052807f198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c281526020017f1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed81526020017f090689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b81526020017f12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa815250600e906004620000d492919062000534565b50348015620000e257600080fd5b50604051620021df380380620021df83398181016040526101008110156200010957600080fd5b815160208301516040808501805191519395929483019291846401000000008211156200013557600080fd5b9083019060208201858111156200014b57600080fd5b82516401000000008111828201881017156200016657600080fd5b82525081516020918201929091019080838360005b83811015620001955781810151838201526020016200017b565b50505050905090810190601f168015620001c35780820380516001836020036101000a031916815260200191505b5060405260200180516040519392919084640100000000821115620001e757600080fd5b908301906020820185811115620001fd57600080fd5b82516401000000008111828201881017156200021857600080fd5b82525081516020918201929091019080838360005b83811015620002475781810151838201526020016200022d565b50505050905090810190601f168015620002755780820380516001836020036101000a031916815260200191505b506040819052600080546001600160a01b0319166001600160a01b038981169190911780835560209490940195508894509290921691907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908290a3506000846001600160a01b0316638da5cb5b6040518163ffffffff1660e01b815260040160206040518083038186803b1580156200030e57600080fd5b505afa15801562000323573d6000803e3d6000fd5b505050506040513d60208110156200033a57600080fd5b505190506200034d600383600462000534565b50600180546001600160a01b0319166001600160a01b03881617905583516200037e90600290602087019062000565565b5082516200039490600990602086019062000565565b506000600a8190556007805460018181019092557fa66cc928b5edb82af9bd49922954155ab7b0942694bea4ce44661d9a8736c6880180546001600160a01b038086166001600160a01b0319928316811790935591845260086020908152604080862080548c8616941684179055935484517f15db20690000000000000000000000000000000000000000000000000000000081526004810193909352935193909216936315db206993602480840194939192918390030190829087803b1580156200045f57600080fd5b505af115801562000474573d6000803e3d6000fd5b505050506040513d60208110156200048b57600080fd5b50516001600160a01b039182166000908152600860205260409020600101805475010000000000000000000000000000000000000000006001600160a01b0319909116939092169290921760ff60a81b191617905550620005f79350505050565b826002810192821562000522579160200282015b8281111562000522578251829060ff1690559160200191906001019062000500565b5062000530929150620005d7565b5090565b826004810192821562000522579160200282015b828111156200052257825182559160200191906001019062000548565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10620005a857805160ff191683800117855562000522565b828001600101855582156200052257918201828111156200052257825182559160200191906001019062000548565b620005f491905b80821115620005305760008155600101620005de565b90565b611bd880620006076000396000f3fe608060405234801561001057600080fd5b50600436106101375760003560e01c8063a230c524116100b8578063d53d64041161007c578063d53d6404146103f7578063d66d9e1914610414578063e7acfda11461041c578063e8f738e114610474578063ebe1c2301461049a578063f2fde38b1461056957610137565b8063a230c5241461033c578063ab04010714610362578063b688a3631461036a578063cc2867d714610372578063d28d8852146103ef57610137565b8063731f6fd8116100ff578063731f6fd8146102495780638379088e146102665780638da5cb5b146102d65780638f32d59b146102fa57806396c551751461031657610137565b806312f8ad051461013c57806342cde4e81461016b5780634b77c4681461017357806366fd3cd81461019b578063715018a614610241575b600080fd5b6101596004803603602081101561015257600080fd5b503561058f565b60408051918252519081900360200190f35b6101596105a3565b6101996004803603602081101561018957600080fd5b50356001600160a01b03166105ae565b005b610199600480360360208110156101b157600080fd5b8101906020810181356401000000008111156101cc57600080fd5b8201836020820111156101de57600080fd5b8035906020019184600183028401116401000000008311171561020057600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295506107bf945050505050565b610199610a6b565b6101596004803603602081101561025f57600080fd5b5035610ac6565b6101996004803603602081101561027c57600080fd5b81019060208101813564010000000081111561029757600080fd5b8201836020820111156102a957600080fd5b803590602001918460018302840111640100000000831117156102cb57600080fd5b509092509050610ad3565b6102de610d78565b604080516001600160a01b039092168252519081900360200190f35b610302610d87565b604080519115158252519081900360200190f35b6101996004803603602081101561032c57600080fd5b50356001600160a01b0316610d98565b6103026004803603602081101561035257600080fd5b50356001600160a01b0316610f45565b610199610f6d565b61019961108c565b61037a6112b5565b6040805160208082528351818301528351919283929083019185019080838360005b838110156103b457818101518382015260200161039c565b50505050905090810190601f1680156103e15780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b61037a611343565b6101596004803603602081101561040d57600080fd5b503561139b565b6101996113a8565b61042461154b565b60408051602080825283518183015283519192839290830191858101910280838360005b83811015610460578181015183820152602001610448565b505050509050019250505060405180910390f35b6102de6004803603602081101561048a57600080fd5b50356001600160a01b03166115ad565b610199600480360360608110156104b057600080fd5b8101906020810181356401000000008111156104cb57600080fd5b8201836020820111156104dd57600080fd5b803590602001918460018302840111640100000000831117156104ff57600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152505060408051808201825293969594818101949350915060029083908390808284376000920191909152509194506115ce9350505050565b6101996004803603602081101561057f57600080fd5b50356001600160a01b0316611982565b6003816004811061059c57fe5b0154905081565b600754600290045b90565b6105b733610f45565b6105c057600080fd5b6000816001600160a01b0316638da5cb5b6040518163ffffffff1660e01b815260040160206040518083038186803b1580156105fb57600080fd5b505afa15801561060f573d6000803e3d6000fd5b505050506040513d602081101561062557600080fd5b50516001600160a01b038116600090815260086020526040902060010154909150600160a81b900460ff161561068c5760405162461bcd60e51b815260040180806020018281038252602a815260200180611b7a602a913960400191505060405180910390fd5b6001600160a01b038181166000908152600b602052604090205416156106f9576040805162461bcd60e51b815260206004820181905260248201527f6163636f756e742068617320616c7265616479206265656e20696e7669746564604482015290519081900360640190fd5b816001600160a01b031663eec30bfd6040518163ffffffff1660e01b8152600401600060405180830381600087803b15801561073457600080fd5b505af1158015610748573d6000803e3d6000fd5b5050604080513081526001600160a01b038616602082015281517f4509853340a1799d2c49c6cb19f9988fc2fa7f8c37aa55a7d865e19a260cfdc99450908190039091019150a16001600160a01b039081166000908152600b602052604090208054919092166001600160a01b0319909116179055565b6107c833610f45565b6107d157600080fd5b6007546001141561092e57600154604080516395184d3b60e01b815233600482015290517ff7084e4afb12f9dc96cee0f3c56fee76afe9598f22f0e3ce04f21b4c18cdc26792309285926001600160a01b03909216916395184d3b91602480820192602092909190829003018186803b15801561084d57600080fd5b505afa158015610861573d6000803e3d6000fd5b505050506040513d602081101561087757600080fd5b5051600a54604080516001600160a01b0380871682528416918101919091526060810182905260806020828101828152865192840192909252855160a084019187019080838360005b838110156108d85781810151838201526020016108c0565b50505050905090810190601f1680156109055780820380516001836020036101000a031916815260200191505b509550505050505060405180910390a18051610928906009906020840190611a05565b50610a68565b336000908152600860209081526040808320600190810154600a54835163bade603360e01b8152920160248301819052600483019384528651604484015286516001600160a01b0390921695869563bade6033958995939490938493606490910192870191908190849084905b838110156109b357818101518382015260200161099b565b50505050905090810190601f1680156109e05780820380516001836020036101000a031916815260200191505b509350505050600060405180830381600087803b158015610a0057600080fd5b505af1158015610a14573d6000803e3d6000fd5b5050600a54604080513081526001600160a01b0386166020820152600190920182820152517fe407c1bec7341bd0a3dcb85da18f4e0b5126fcc79738f8f1dde296455e9c35aa9350908190036060019150a1505b50565b610a73610d87565b610a7c57600080fd5b600080546040516001600160a01b03909116907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a3600080546001600160a01b0319169055565b600c816002811061059c57fe5b6000336001600160a01b0316635d1ca6316040518163ffffffff1660e01b815260040160206040518083038186803b158015610b0e57600080fd5b505afa158015610b22573d6000803e3d6000fd5b505050506040513d6020811015610b3857600080fd5b5051600a549091508111610b88576040805162461bcd60e51b8152602060048201526012602482015271436f6e73656e73757320657870697265642160701b604482015290519081900360640190fd5b6000336001600160a01b031663e9790d026040518163ffffffff1660e01b815260040160206040518083038186803b158015610bc357600080fd5b505afa158015610bd7573d6000803e3d6000fd5b505050506040513d6020811015610bed57600080fd5b505160408051638da5cb5b60e01b815290519192506000916001600160a01b03841691638da5cb5b916004808301926020929190829003018186803b158015610c3557600080fd5b505afa158015610c49573d6000803e3d6000fd5b505050506040513d6020811015610c5f57600080fd5b50516001600160a01b03808216600090815260086020526040902060010154919250163314610cbf5760405162461bcd60e51b8152600401808060200182810382526027815260200180611b536027913960400191505060405180910390fd5b610ccb60098686611a83565b5082600a819055507ff7084e4afb12f9dc96cee0f3c56fee76afe9598f22f0e3ce04f21b4c18cdc26730868685600a5460405180866001600160a01b03166001600160a01b0316815260200180602001846001600160a01b03166001600160a01b031681526020018381526020018281038252868682818152602001925080828437600083820152604051601f909101601f19169092018290039850909650505050505050a15050505050565b6000546001600160a01b031690565b6000546001600160a01b0316331490565b610da133610f45565b610daa57600080fd5b6001600160a01b038082166000908152600860205260408082205481516315fe21e560e31b81523060048201529151931692839263aff10f28926024808201939182900301818387803b158015610e0057600080fd5b505af1158015610e14573d6000803e3d6000fd5b505050506001600160a01b0382166000908152600860205260408120600101805460ff60a81b191690555b600754811015610efd57826001600160a01b031660078281548110610e6057fe5b6000918252602090912001546001600160a01b03161415610ef557600780546000198101908110610e8d57fe5b600091825260209091200154600780546001600160a01b039092169183908110610eb357fe5b600091825260209091200180546001600160a01b0319166001600160a01b03929092169190911790556007805490610eef906000198301611af1565b50610efd565b600101610e3f565b604080513081526001600160a01b038416602082015281517fa0ae7458a482bd2780bed4edd090c3f9e9e3e600171cdc4709da09e5bd233e18929181900390910190a1505050565b6001600160a01b0316600090815260086020526040902060010154600160a81b900460ff1690565b336000908152600b60205260409020546001600160a01b031680610fd2576040805162461bcd60e51b81526020600482015260176024820152761858d8dbdd5b9d081dd85cc81b9bdd081a5b9d9a5d1959604a1b604482015290519081900360640190fd5b806001600160a01b031663b061d9a96040518163ffffffff1660e01b8152600401600060405180830381600087803b15801561100d57600080fd5b505af1158015611021573d6000803e3d6000fd5b5050336000908152600b602090815260409182902080546001600160a01b031916905581513081526001600160a01b0386169181019190915281517f7ea317cacaf634632d88c1fd20864b1f54f609bb063c678280aac6312715edef9450908190039091019150a150565b336000908152600b60205260409020546001600160a01b0316806110f1576040805162461bcd60e51b81526020600482015260176024820152761858d8dbdd5b9d081dd85cc81b9bdd081a5b9d9a5d1959604a1b604482015290519081900360640190fd5b6007805460018082019092557fa66cc928b5edb82af9bd49922954155ab7b0942694bea4ce44661d9a8736c688018054336001600160a01b03199182168117909255600091825260086020908152604080842080549093166001600160a01b03878116918217909455945481516315db206960e01b8152600481019690965290519216936315db206993602480830194928390030190829087803b15801561119857600080fd5b505af11580156111ac573d6000803e3d6000fd5b505050506040513d60208110156111c257600080fd5b5051336000908152600860205260408082206001018054600160a81b6001600160a01b03199091166001600160a01b039586161760ff60a81b1916179055805163020fa19160e61b81529051928416926383e864409260048084019391929182900301818387803b15801561123657600080fd5b505af115801561124a573d6000803e3d6000fd5b5050336000908152600b602090815260409182902080546001600160a01b031916905581513081526001600160a01b0386169181019190915281517f4d7c243e154e530692e62f8539db65779f5cb85d58831956361697addede5adb9450908190039091019150a150565b6009805460408051602060026001851615610100026000190190941693909304601f8101849004840282018401909252818152929183018282801561133b5780601f106113105761010080835404028352916020019161133b565b820191906000526020600020905b81548152906001019060200180831161131e57829003601f168201915b505050505081565b6002805460408051602060018416156101000260001901909316849004601f8101849004840282018401909252818152929183018282801561133b5780601f106113105761010080835404028352916020019161133b565b600e816004811061059c57fe5b6113b133610f45565b6113ba57600080fd5b336000908152600860205260408082205481516315fe21e560e31b815230600482015291516001600160a01b0390911692839263aff10f28926024808301939282900301818387803b15801561140f57600080fd5b505af1158015611423573d6000803e3d6000fd5b5050336000908152600860205260408120600101805460ff60a81b191690559150505b60075481101561150457336001600160a01b03166007828154811061146757fe5b6000918252602090912001546001600160a01b031614156114fc5760078054600019810190811061149457fe5b600091825260209091200154600780546001600160a01b0390921691839081106114ba57fe5b600091825260209091200180546001600160a01b0319166001600160a01b039290921691909117905560078054906114f6906000198301611af1565b50611504565b600101611446565b604080513081526001600160a01b038416602082015281517fa0ae7458a482bd2780bed4edd090c3f9e9e3e600171cdc4709da09e5bd233e18929181900390910190a15050565b606060078054806020026020016040519081016040528092919081815260200182805480156115a357602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311611585575b5050505050905090565b6001600160a01b039081166000908152600860205260409020600101541690565b6115d733610f45565b6115e057600080fd5b6115e8611b1a565b6040516378b8c33160e11b815260206004820181815285516024840152855173__ecOps_________________________________9363f17186629388939283926044019185019080838360005b8381101561164d578181015183820152602001611635565b50505050905090810190601f16801561167a5780820380516001836020036101000a031916815260200191505b5092505050604080518083038186803b15801561169657600080fd5b505af41580156116aa573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525060408110156116cf57600080fd5b50604080516349ac97d160e01b815291925073__ecOps_________________________________916349ac97d19184916003918791600e916004909101908190869080838360005b8381101561172f578181015183820152602001611717565b5050509201915085905060046020028201915b8154815260200190600101908083116117425750849050604080838360005b83811015611779578181015183820152602001611761565b5050509201915083905060046020028201915b81548152602001906001019080831161178c57505094505050505060206040518083038186803b1580156117bf57600080fd5b505af41580156117d3573d6000803e3d6000fd5b505050506040513d60208110156117e957600080fd5b5051611830576040805162461bcd60e51b8152602060048201526011602482015270696e76616c6964207369676e617475726560781b604482015290519081900360640190fd5b600154604080516395184d3b60e01b815233600482015290517ff7084e4afb12f9dc96cee0f3c56fee76afe9598f22f0e3ce04f21b4c18cdc26792309287926001600160a01b03909216916395184d3b91602480820192602092909190829003018186803b1580156118a157600080fd5b505afa1580156118b5573d6000803e3d6000fd5b505050506040513d60208110156118cb57600080fd5b5051600a54604080516001600160a01b0380871682528416918101919091526060810182905260806020828101828152865192840192909252855160a084019187019080838360005b8381101561192c578181015183820152602001611914565b50505050905090810190601f1680156119595780820380516001836020036101000a031916815260200191505b509550505050505060405180910390a1825161197c906009906020860190611a05565b50505050565b61198a610d87565b61199357600080fd5b610a68816001600160a01b0381166119aa57600080fd5b600080546040516001600160a01b03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a3600080546001600160a01b0319166001600160a01b0392909216919091179055565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10611a4657805160ff1916838001178555611a73565b82800160010185558215611a73579182015b82811115611a73578251825591602001919060010190611a58565b50611a7f929150611b38565b5090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10611ac45782800160ff19823516178555611a73565b82800160010185558215611a73579182015b82811115611a73578235825591602001919060010190611ad6565b815481835581811115611b1557600083815260209020611b15918101908301611b38565b505050565b60405180604001604052806002906020820280388339509192915050565b6105ab91905b80821115611a7f5760008155600101611b3e56fe436f6e73656e73757320646f6573206e6f742062656c6f6e6720746f207468652067726f757021546865207573657220746f20626520696e766974656420697320616c72656164792061206d656d626572a265627a7a72315820a3c134ac558e2b397182543cdbf0273cc02e524b506ac0ebb9ad9c337114384b64736f6c634300050c0032a265627a7a723158208e620204dc8615163340c5898e8ac0e4be25e774ff6e4a5d97036197524297eb64736f6c634300050c0032"

// DeployGroupFactory deploys a new Ethereum contract, binding an instance of GroupFactory to it.
func DeployGroupFactory(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *GroupFactory, error) {
	parsed, err := abi.JSON(strings.NewReader(GroupFactoryABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(GroupFactoryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &GroupFactory{GroupFactoryCaller: GroupFactoryCaller{contract: contract}, GroupFactoryTransactor: GroupFactoryTransactor{contract: contract}, GroupFactoryFilterer: GroupFactoryFilterer{contract: contract}}, nil
}

// GroupFactory is an auto generated Go binding around an Ethereum contract.
type GroupFactory struct {
	GroupFactoryCaller     // Read-only binding to the contract
	GroupFactoryTransactor // Write-only binding to the contract
	GroupFactoryFilterer   // Log filterer for contract events
}

// GroupFactoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type GroupFactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GroupFactoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GroupFactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GroupFactoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GroupFactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GroupFactorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GroupFactorySession struct {
	Contract     *GroupFactory     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GroupFactoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GroupFactoryCallerSession struct {
	Contract *GroupFactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// GroupFactoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GroupFactoryTransactorSession struct {
	Contract     *GroupFactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// GroupFactoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type GroupFactoryRaw struct {
	Contract *GroupFactory // Generic contract binding to access the raw methods on
}

// GroupFactoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GroupFactoryCallerRaw struct {
	Contract *GroupFactoryCaller // Generic read-only contract binding to access the raw methods on
}

// GroupFactoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GroupFactoryTransactorRaw struct {
	Contract *GroupFactoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGroupFactory creates a new instance of GroupFactory, bound to a specific deployed contract.
func NewGroupFactory(address common.Address, backend bind.ContractBackend) (*GroupFactory, error) {
	contract, err := bindGroupFactory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GroupFactory{GroupFactoryCaller: GroupFactoryCaller{contract: contract}, GroupFactoryTransactor: GroupFactoryTransactor{contract: contract}, GroupFactoryFilterer: GroupFactoryFilterer{contract: contract}}, nil
}

// NewGroupFactoryCaller creates a new read-only instance of GroupFactory, bound to a specific deployed contract.
func NewGroupFactoryCaller(address common.Address, caller bind.ContractCaller) (*GroupFactoryCaller, error) {
	contract, err := bindGroupFactory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GroupFactoryCaller{contract: contract}, nil
}

// NewGroupFactoryTransactor creates a new write-only instance of GroupFactory, bound to a specific deployed contract.
func NewGroupFactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*GroupFactoryTransactor, error) {
	contract, err := bindGroupFactory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GroupFactoryTransactor{contract: contract}, nil
}

// NewGroupFactoryFilterer creates a new log filterer instance of GroupFactory, bound to a specific deployed contract.
func NewGroupFactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*GroupFactoryFilterer, error) {
	contract, err := bindGroupFactory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GroupFactoryFilterer{contract: contract}, nil
}

// bindGroupFactory binds a generic wrapper to an already deployed contract.
func bindGroupFactory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(GroupFactoryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GroupFactory *GroupFactoryRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _GroupFactory.Contract.GroupFactoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GroupFactory *GroupFactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GroupFactory.Contract.GroupFactoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GroupFactory *GroupFactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GroupFactory.Contract.GroupFactoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GroupFactory *GroupFactoryCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _GroupFactory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GroupFactory *GroupFactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GroupFactory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GroupFactory *GroupFactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GroupFactory.Contract.contract.Transact(opts, method, params...)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_GroupFactory *GroupFactoryCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _GroupFactory.contract.Call(opts, out, "isOwner")
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_GroupFactory *GroupFactorySession) IsOwner() (bool, error) {
	return _GroupFactory.Contract.IsOwner(&_GroupFactory.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_GroupFactory *GroupFactoryCallerSession) IsOwner() (bool, error) {
	return _GroupFactory.Contract.IsOwner(&_GroupFactory.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_GroupFactory *GroupFactoryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _GroupFactory.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_GroupFactory *GroupFactorySession) Owner() (common.Address, error) {
	return _GroupFactory.Contract.Owner(&_GroupFactory.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_GroupFactory *GroupFactoryCallerSession) Owner() (common.Address, error) {
	return _GroupFactory.Contract.Owner(&_GroupFactory.CallOpts)
}

// Create is a paid mutator transaction binding the contract method 0xa15ab08d.
//
// Solidity: function create(address account, string name) returns(address)
func (_GroupFactory *GroupFactoryTransactor) Create(opts *bind.TransactOpts, account common.Address, name string) (*types.Transaction, error) {
	return _GroupFactory.contract.Transact(opts, "create", account, name)
}

// Create is a paid mutator transaction binding the contract method 0xa15ab08d.
//
// Solidity: function create(address account, string name) returns(address)
func (_GroupFactory *GroupFactorySession) Create(account common.Address, name string) (*types.Transaction, error) {
	return _GroupFactory.Contract.Create(&_GroupFactory.TransactOpts, account, name)
}

// Create is a paid mutator transaction binding the contract method 0xa15ab08d.
//
// Solidity: function create(address account, string name) returns(address)
func (_GroupFactory *GroupFactoryTransactorSession) Create(account common.Address, name string) (*types.Transaction, error) {
	return _GroupFactory.Contract.Create(&_GroupFactory.TransactOpts, account, name)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_GroupFactory *GroupFactoryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GroupFactory.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_GroupFactory *GroupFactorySession) RenounceOwnership() (*types.Transaction, error) {
	return _GroupFactory.Contract.RenounceOwnership(&_GroupFactory.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_GroupFactory *GroupFactoryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _GroupFactory.Contract.RenounceOwnership(&_GroupFactory.TransactOpts)
}

// SetParent is a paid mutator transaction binding the contract method 0x1499c592.
//
// Solidity: function setParent(address parent) returns()
func (_GroupFactory *GroupFactoryTransactor) SetParent(opts *bind.TransactOpts, parent common.Address) (*types.Transaction, error) {
	return _GroupFactory.contract.Transact(opts, "setParent", parent)
}

// SetParent is a paid mutator transaction binding the contract method 0x1499c592.
//
// Solidity: function setParent(address parent) returns()
func (_GroupFactory *GroupFactorySession) SetParent(parent common.Address) (*types.Transaction, error) {
	return _GroupFactory.Contract.SetParent(&_GroupFactory.TransactOpts, parent)
}

// SetParent is a paid mutator transaction binding the contract method 0x1499c592.
//
// Solidity: function setParent(address parent) returns()
func (_GroupFactory *GroupFactoryTransactorSession) SetParent(parent common.Address) (*types.Transaction, error) {
	return _GroupFactory.Contract.SetParent(&_GroupFactory.TransactOpts, parent)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_GroupFactory *GroupFactoryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _GroupFactory.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_GroupFactory *GroupFactorySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _GroupFactory.Contract.TransferOwnership(&_GroupFactory.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_GroupFactory *GroupFactoryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _GroupFactory.Contract.TransferOwnership(&_GroupFactory.TransactOpts, newOwner)
}

// GroupFactoryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the GroupFactory contract.
type GroupFactoryOwnershipTransferredIterator struct {
	Event *GroupFactoryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *GroupFactoryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupFactoryOwnershipTransferred)
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
		it.Event = new(GroupFactoryOwnershipTransferred)
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
func (it *GroupFactoryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupFactoryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupFactoryOwnershipTransferred represents a OwnershipTransferred event raised by the GroupFactory contract.
type GroupFactoryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_GroupFactory *GroupFactoryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*GroupFactoryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _GroupFactory.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &GroupFactoryOwnershipTransferredIterator{contract: _GroupFactory.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_GroupFactory *GroupFactoryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *GroupFactoryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _GroupFactory.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupFactoryOwnershipTransferred)
				if err := _GroupFactory.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_GroupFactory *GroupFactoryFilterer) ParseOwnershipTransferred(log types.Log) (*GroupFactoryOwnershipTransferred, error) {
	event := new(GroupFactoryOwnershipTransferred)
	if err := _GroupFactory.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}
