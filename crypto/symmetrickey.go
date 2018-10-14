package crypto

import (
	"golang.org/x/crypto/nacl/secretbox"
	"io"
	"encoding/json"
	"github.com/pkg/errors"
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

func (k *SymmetricKey) Encode() ([]byte, error) {
	enc, err := json.Marshal(k)
	if err != nil {
		errors.Wrap(err, "could not encode SymmetricKey")
	}

	return enc, nil
}

func DecodeSymmetricKey(data []byte) (*SymmetricKey, error) {
	var k SymmetricKey
	if err := json.Unmarshal(data, &k); err != nil {
		return nil, errors.Wrap(err, "could not decode SymmetricKey")
	}

	return &k, nil
}