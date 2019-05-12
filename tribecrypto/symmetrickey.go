// Copyright (c) 2019 Laszlo Sari
//
// FileTribe is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// FileTribe is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//

package tribecrypto

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"io"

	"github.com/pkg/errors"
	"golang.org/x/crypto/nacl/secretbox"
)

// SymmetricKey is a helper struct for handling nacl.secretbox
type SymmetricKey struct {
	Key [32]byte  `json:"key"`
	RNG io.Reader `json:"-"`
}

func (k *SymmetricKey) getNonce() [24]byte {
	var nonce [24]byte
	k.RNG.Read(nonce[:])
	return nonce
}

// BoxSeal encrypts the provided message
func (k *SymmetricKey) BoxSeal(message []byte) []byte {
	nonce := k.getNonce()
	enc := secretbox.Seal(nonce[:], message, &nonce, &k.Key)
	return enc
}

// BoxOpen decrypts the provided cipher text
func (k *SymmetricKey) BoxOpen(bytesBox []byte) ([]byte, bool) {
	if len(bytesBox) < 24 {
		return []byte{}, false
	}

	var nonce [24]byte
	copy(nonce[:], bytesBox[:24])
	return secretbox.Open(nil, bytesBox[24:], &nonce, &k.Key)
}

// Encode encodes the secret key
func (k *SymmetricKey) Encode() ([]byte, error) {
	enc, err := json.Marshal(k)
	if err != nil {
		return []byte{}, errors.Wrap(err, "could not encode SymmetricKey")
	}

	return enc, nil
}

// DecodeSymmetricKey decodes and returns an encoded SymmetricKey
func DecodeSymmetricKey(data []byte) (*SymmetricKey, error) {
	var k SymmetricKey
	if err := json.Unmarshal(data, &k); err != nil {
		return nil, errors.Wrap(err, "could not decode SymmetricKey")
	}
	k.RNG = rand.Reader

	return &k, nil
}

func IsBoxerNotNull(boxer SymmetricKey) bool {
	nullKey := [32]byte{0}
	return bytes.Equal(boxer.Key[:], nullKey[:])
}
