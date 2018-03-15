package crypto

import (
	"bytes"
	"encoding/base64"
)

type PublicKeyHash []byte

func (p *PublicKeyHash) ToBase64() string {
	return base64.StdEncoding.EncodeToString(*p)
}

func (p *PublicKeyHash) Equals(other *PublicKeyHash) bool {
	if bytes.Compare(*p, *other) != 0 {
		return false
	}
	return true
}

func Base64ToPublicKeyHash(src string) (PublicKeyHash, error) {
	dBytes, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return nil, err
	}
	return PublicKeyHash(dBytes), nil
}
