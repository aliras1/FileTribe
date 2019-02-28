package tribecrypto

import (
	"hash"
	"github.com/golang/glog"
	"github.com/ethereum/go-ethereum/crypto/sha3"
)

type IHasher interface {
	Sum(...[]byte) []byte
}

type Keccak256Hasher struct {
	hash hash.Hash
}

func (hasher *Keccak256Hasher) Sum(args ...[]byte) []byte {
	hasher.hash.Reset()

	for _, arg := range args {
		_, err := hasher.hash.Write(arg)
		if err != nil {
			glog.Error("error while keccak256 write: %s", err)
			return nil
		}
	}

	return hasher.hash.Sum(nil)
}

func NewKeccak256Hasher() *Keccak256Hasher {
	return &Keccak256Hasher{
		hash: sha3.NewKeccak256(),
	}
}

