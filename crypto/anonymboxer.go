package crypto

import (
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/nacl/box"
	"golang.org/x/crypto/sha3"
)

type AnonymPublicKey struct {
	Value *[32]byte
}

type AnonymSecretKey struct {
	Value *[32]byte
}

type AnonymBoxer struct {
	PublicKey AnonymPublicKey
	SecretKey AnonymSecretKey
}

func getNonce(pk1, pk2 *[32]byte) *[24]byte {
	var nonce [24]byte
	digest := sha3.Sum512(append(pk1[:], pk2[:]...))
	copy(nonce[:], digest[:24])
	return &nonce
}

func random() ([]byte, error) {
	var r [32]byte
	n, err := rand.Read(r[:])
	if err != nil {
		return nil, err
	}
	if n != 32 {
		return nil, fmt.Errorf("could not read enough bytes from random: random()")
	}
	return r[:], nil
}

func (pk AnonymPublicKey) Seal(m []byte) ([]byte, error) {
	ephemeral_pk, ephemeral_sk, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("could not generate ephemeral key: AnonymPublicKey.Seal(): %s", err)
	}

	nonce := getNonce(ephemeral_pk, pk.Value)
	r, err := random()
	if err != nil {
		return nil, err
	}
	m = append(r, m...)

	ct := append(ephemeral_pk[:], box.Seal(nil, m, nonce, pk.Value, ephemeral_sk)...)
	return ct, nil
}

func (boxer AnonymBoxer) Seal(m []byte) ([]byte, error) {
	return boxer.PublicKey.Seal(m)
}

func (boxer AnonymBoxer) Open(ct []byte) ([]byte, error) {
	if len(ct) <= 32 {
		return nil, fmt.Errorf("invalid cipher text: not long enough")
	}
	var ephemeral_pk [32]byte
	copy(ephemeral_pk[:], ct[:32])
	nonce := getNonce(&ephemeral_pk, boxer.PublicKey.Value)
	m, ok := box.Open(nil, ct[32:], nonce, &ephemeral_pk, boxer.SecretKey.Value)

	if !ok {
		return nil, fmt.Errorf("could not decrypt")
	}
	return m[32:], nil
}
