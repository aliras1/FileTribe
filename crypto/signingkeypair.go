package crypto

import (
	"bytes"
	"crypto/sha512"
	"encoding/base64"

	"golang.org/x/crypto/ed25519"

	"ipfs-share-go/crypto/ed25519/util/edwards25519"
)

type PublicSigningKey ed25519.PublicKey
type SecretSigningKey ed25519.PrivateKey

type SigningKeyPair struct {
	PublicKey PublicSigningKey
	SecretKey SecretSigningKey
}

func (s *PublicSigningKey) ToBase64() string {
	return base64.StdEncoding.EncodeToString(*s)
}

func (p *PublicSigningKey) Equals(other *PublicSigningKey) bool {
	if bytes.Compare(*p, *other) != 0 {
		return false
	}
	return true
}

func Ed25519KeyPair(sk *[32]byte) (PublicSigningKey, SecretSigningKey) {
	// equivalent to https://github.com/golang/crypto/blob/master/ed25519/ed25519.go
	// the only difference is that it uses sk instead of an RNG

	privateKey := make([]byte, 64)
	publicKey := make([]byte, 32)
	copy(privateKey[:32], (*sk)[:])

	digest := sha512.Sum512(privateKey[:32])
	digest[0] &= 248
	digest[31] &= 127
	digest[31] |= 64

	var A edwards25519.ExtendedGroupElement
	var hBytes [32]byte
	copy(hBytes[:], digest[:])
	edwards25519.GeScalarMultBase(&A, &hBytes)
	var publicKeyBytes [32]byte
	A.ToBytes(&publicKeyBytes)

	copy(privateKey[32:], publicKeyBytes[:])
	copy(publicKey, publicKeyBytes[:])
	return publicKey, privateKey
}

func Base64ToPublicSigningKey(src string) (PublicSigningKey, error) {
	pBytes, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return nil, err
	}
	return PublicSigningKey(pBytes), nil
}
