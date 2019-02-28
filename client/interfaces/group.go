package interfaces

import (
	ethcommon "github.com/ethereum/go-ethereum/common"

	"ipfs-share/crypto"
)

type IGroup interface {
	Address() ethcommon.Address
	Name() string
	IpfsHash() string
	SetIpfsHash(encIpfsHash []byte) error
	EncryptedIpfsHash() []byte
	AddMember(user ethcommon.Address)
	RemoveMember(user ethcommon.Address)
	IsMember(user ethcommon.Address) bool
	CountMembers() int
	Members() []ethcommon.Address
	Boxer() tribecrypto.SymmetricKey
	SetBoxer(boxer tribecrypto.SymmetricKey)
	Update(name string, members []ethcommon.Address, encIpfsHash []byte) error
	Encode() ([]byte, error)
	Save() error
}