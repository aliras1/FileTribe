package crypto

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
)

//https://github.com/golang/crypto/blob/master/nacl/sign/sign.go

type Signer struct {
	PrivateKey *ecdsa.PrivateKey
}

func(s *Signer) Sign(digest []byte) ([]byte, error) {
	return crypto.Sign(digest, s.PrivateKey)
}

type VerifyKey []byte

func (vk *VerifyKey) Verify(digest, signature []byte) bool {
	return crypto.VerifySignature(*vk, digest, signature)
}