package crypto

import (
	"crypto/ecdsa"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/glog"
)

//https://github.com/golang/crypto/blob/master/nacl/sign/sign.go

type Signer struct {
	PrivateKey *ecdsa.PrivateKey
}

func(s *Signer) Sign(digest []byte) ([]byte, error) {
	return ethcrypto.Sign(digest, s.PrivateKey)
}

func (s *Signer) Verify(digest, signature []byte) bool {
	pk := ethcrypto.CompressPubkey(&s.PrivateKey.PublicKey)
	return ethcrypto.VerifySignature(pk, digest, signature[:len(signature) - 1])
}


type VerifyKey []byte

func (vk *VerifyKey) Verify(digest, signature []byte) bool {
	addr, err := ethcrypto.SigToPub(digest, signature)
	if err != nil {
		glog.Errorf("error while ecrecover: %s", err)
		return false
	}

	glog.Infof("ecrecover addr: %s", ethcrypto.PubkeyToAddress(*addr).String())

	return ethcrypto.VerifySignature(*vk, digest, signature[:len(signature) - 1])
}