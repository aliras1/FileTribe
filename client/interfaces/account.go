package interfaces

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/chequebook"

	ethacc "github.com/aliras1/FileTribe/eth/gen/Account"
	"github.com/aliras1/FileTribe/tribecrypto"
)

type IAccount interface {
	ContractAddress() ethcommon.Address
	Contract() *ethacc.Account
	Name() string
	Boxer() tribecrypto.AnonymBoxer
	SetContract(addr ethcommon.Address, backend chequebook.Backend) error
	Save() error
}
