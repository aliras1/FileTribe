package tribecrypto

import (
	"crypto/rand"
	"encoding/json"
	"github.com/pkg/errors"
	"golang.org/x/crypto/nacl/secretbox"
	"io"
)

type SymmetricKey struct {
	Key [32]byte  `json:"key"`
	RNG io.Reader `json:"-"`
}

func NewSymmetricKey() (*SymmetricKey, error) {
	var k [32]byte
	if _, err := rand.Read(k[:]); err != nil {
		return &SymmetricKey{}, errors.Wrap(err, "could not read from crypto.rand")
	}

	return &SymmetricKey{Key: k, RNG: rand.Reader}, nil
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
	if len(bytesBox) < 24 {
		return []byte{}, false
	}

	var nonce [24]byte
	copy(nonce[:], bytesBox[:24])
	return secretbox.Open(nil, bytesBox[24:], &nonce, &k.Key)
}

func (k *SymmetricKey) Encode() ([]byte, error) {
	enc, err := json.Marshal(k)
	if err != nil {
		return []byte{}, errors.Wrap(err, "could not encode SymmetricKey")
	}

	return enc, nil
}

func DecodeSymmetricKey(data []byte) (*SymmetricKey, error) {
	var k SymmetricKey
	if err := json.Unmarshal(data, &k); err != nil {
		return nil, errors.Wrap(err, "could not decode SymmetricKey")
	}
	k.RNG = rand.Reader

	return &k, nil
}