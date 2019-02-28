package interfaces

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethacc "ipfs-share/eth/gen/Account"

	"ipfs-share/crypto"
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
