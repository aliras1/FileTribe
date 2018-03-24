package crypto

import (
	"bytes"
	"crypto/sha512"
	"encoding/base64"

	"golang.org/x/crypto/ed25519"

	"ipfs-share/crypto/ed25519/util/edwards25519"
)

type PublicSigningKey ed25519.PublicKey
type SecretSigningKey ed25519.PrivateKey

//https://github.com/golang/crypto/blob/master/nacl/sign/sign.go

const Overhead = 64

func (ssk *SecretSigningKey) Sign(out, message []byte) []byte {
	sig := ed25519.Sign(ed25519.PrivateKey(*ssk), message)
	ret, out := sliceForAppend(out, Overhead+len(message))
	copy(out, sig)
	copy(out[Overhead:], message)
	return ret
}

// Open verifies a signed message produced by Sign and appends the message to
// out, which must not overlap the signed message. The output will be Overhead
// bytes smaller than the signed message.
func (psk *PublicSigningKey) Open(out, signedMessage []byte) ([]byte, bool) {
	if len(signedMessage) < Overhead {
		return nil, false
	}
	if !ed25519.Verify(ed25519.PublicKey(*psk), signedMessage[Overhead:], signedMessage[:Overhead]) {
		return nil, false
	}
	ret, out := sliceForAppend(out, len(signedMessage)-Overhead)
	copy(out, signedMessage[Overhead:])
	return ret, true
}

// sliceForAppend takes a slice and a requested number of bytes. It returns a
// slice with the contents of the given slice followed by that many bytes and a
// second slice that aliases into it and contains only the extra bytes. If the
// original slice has sufficient capacity then no allocation is performed.
func sliceForAppend(in []byte, n int) (head, tail []byte) {
	if total := len(in) + n; cap(in) >= total {
		head = in[:total]
	} else {
		head = make([]byte, total)
		copy(head, in)
	}
	tail = head[len(in):]
	return
}

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
