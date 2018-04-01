package crypto

import (
	"golang.org/x/crypto/nacl/secretbox"
	"io"
)

type SymmetricKey struct {
	Key [32]byte  `json:"key"`
	RNG io.Reader `json:"-"`
}

func (k *SymmetricKey) GetNonce() [24]byte {
	var nonce [24]byte
	k.RNG.Read(nonce[:])
	return nonce
}

func (k *SymmetricKey) BoxSeal(message []byte) []byte {
	nonce := k.GetNonce()
	enc := secretbox.Seal(nonce[:], message, &nonce, &k.Key)
	return enc
}

func (k *SymmetricKey) BoxOpen(bytesBox []byte) ([]byte, bool) {
	var nonce [24]byte
	copy(nonce[:], bytesBox[:24])
	return secretbox.Open(nil, bytesBox[24:], &nonce, &k.Key)
}
