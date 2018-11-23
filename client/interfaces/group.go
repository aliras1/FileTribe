package interfaces

import (
	ethcommon "github.com/ethereum/go-ethereum/common"

	. "ipfs-share/collections"
	"ipfs-share/crypto"
)

type IGroup interface {
	Id() IIdentifier
	Name() string
	IpfsHash() string
	SetIpfsHash(ipfsHash string, encIpfsHash []byte) error
	EncryptedIpfsHash() []byte
	AddMember(user ethcommon.Address)
	RemoveMember(user ethcommon.Address)
	IsMember(user ethcommon.Address) bool
	CountMembers() int
	Members() []ethcommon.Address
	Boxer() crypto.SymmetricKey
	SetBoxer(boxer crypto.SymmetricKey)
	Update(name string, members []ethcommon.Address, encIpfsHash []byte, leader ethcommon.Address) error
	Encode() ([]byte, error)
	Leader() ethcommon.Address
}