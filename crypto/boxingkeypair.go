package crypto

import (
	"bytes"
	"encoding/base64"
	"io"
)

type PublicBoxingKey [32]byte
type SecretBoxingKey [32]byte

func (p *PublicBoxingKey) Bytes() *[32]byte {
	pBytes := [32]byte(*p)
	return &pBytes
}

func (s *SecretBoxingKey) Bytes() *[32]byte {
	sBytes := [32]byte(*s)
	return &sBytes
}

type BoxingKeyPair struct {
	PublicKey PublicBoxingKey
	SecretKey SecretBoxingKey
	RNG       io.Reader
}

func (p *PublicBoxingKey) ToBase64() string {
	return base64.StdEncoding.EncodeToString((*p)[:])
}

func (b *BoxingKeyPair) GetNonce() *[24]byte {
	var nonce [24]byte
	b.RNG.Read(nonce[:])
	return &nonce
}

func (p *PublicBoxingKey) Equals(other *PublicBoxingKey) bool {
	if bytes.Compare((*p)[:], (*other)[:]) != 0 {
		return false
	}
	return true
}

func Base64ToPublicBoxingKey(src string) (PublicBoxingKey, error) {
	pBytes, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return [32]byte{}, err
	}
	publicBoxingKey := PublicBoxingKey{}
	copy(publicBoxingKey[:], pBytes)
	return publicBoxingKey, nil
}
