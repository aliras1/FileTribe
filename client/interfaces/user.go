package interfaces

import (
	ethcommon "github.com/ethereum/go-ethereum/common"

	"ipfs-share/crypto"
)

type IUser interface {
	Address() ethcommon.Address
	Name() string
	Signer() crypto.Signer
	Boxer() crypto.AnonymBoxer
}
