package interfaces

import (
	ethcommon "github.com/ethereum/go-ethereum/common"

	ethacc "github.com/aliras1/FileTribe/eth/gen/Account"
	"github.com/aliras1/FileTribe/tribecrypto"
)

type IAccount interface {
	ContractAddress() ethcommon.Address
	Contract() *ethacc.Account
	Name() string
	Boxer() tribecrypto.AnonymBoxer
	SetContract(contract *ethacc.Account)
	SetContractAddress(addr ethcommon.Address)
	Save() error
}
