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
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/nacl/box"
	"golang.org/x/crypto/sha3"
)

// The whole tribecrypto is a mess and should be removed
// TODO: check if ipfs p2p services are encrypted. If so
// there is really no reason for this to exist anymore

// AnonymPublicKey is the public key part of an anonym boxer
type AnonymPublicKey struct {
	Value [32]byte
}

// AnonymPrivateKey is the private key part of an anonym boxer
type AnonymPrivateKey struct {
	Value [32]byte
}

// AnonymBoxer is a helper type for handling nacl boxes
type AnonymBoxer struct {
	PublicKey  AnonymPublicKey
	PrivateKey AnonymPrivateKey
}

func getNonce(pk1, pk2 *[32]byte) *[24]byte {
	var nonce [24]byte
	digest := sha3.Sum512(append(pk1[:], pk2[:]...))
	copy(nonce[:], digest[:24])
	return &nonce
}

// Seal encrypts the provided bytes
func (pk AnonymPublicKey) Seal(m []byte) ([]byte, error) {
	ephemeralPk, ephemeralSk, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("could not generate ephemeral key: AnonymPublicKey.Seal(): %s", err)
	}

	nonce := getNonce(ephemeralPk, &pk.Value)
	var r [32]byte
	_, err = rand.Read(r[:])
	if err != nil {
		return nil, err
	}
	m = append(r[:], m...)

	ct := append(ephemeralPk[:], box.Seal(nil, m, nonce, &pk.Value, ephemeralSk)...)
	return ct, nil
}

// Seal encrypts te provided bytes
func (boxer AnonymBoxer) Seal(m []byte) ([]byte, error) {
	return boxer.PublicKey.Seal(m)
}

// Open decrypts the provided cipher text
func (boxer AnonymBoxer) Open(ct []byte) ([]byte, error) {
	if len(ct) <= 32 {
		return nil, fmt.Errorf("invalid cipher text: not long enough")
	}
	var ephemeralPk [32]byte
	copy(ephemeralPk[:], ct[:32])
	nonce := getNonce(&ephemeralPk, &boxer.PublicKey.Value)
	m, ok := box.Open(nil, ct[32:], nonce, &ephemeralPk, &boxer.PrivateKey.Value)

	if !ok {
		return nil, fmt.Errorf("could not decrypt")
	}
	return m[32:], nil
}

// AuthSeal is the authenticated version of AnonymBoxer.Seal
func AuthSeal(message []byte, otherPK *AnonymPublicKey, mySK *AnonymPrivateKey) ([]byte, error) {
	var nonce [24]byte
	_, err := rand.Read(nonce[:])
	if err != nil {
		return nil, err
	}
	ct := box.Seal(nonce[:], message, &nonce, &otherPK.Value, &mySK.Value)
	return ct, nil
}

// AuthOpen is the authenticated version of AnonymBoxer.Open
func AuthOpen(bytesBox []byte, otherPK *AnonymPublicKey, mySK *AnonymPrivateKey) ([]byte, bool) {
	var nonce [24]byte
	copy(nonce[:], bytesBox[:24])
	return box.Open(nil, bytesBox[24:], &nonce, &otherPK.Value, &mySK.Value)
}
