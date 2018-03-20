package client

import (
	"crypto/rand"
	"crypto/sha256"

	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/scrypt"

	"ipfs-share/crypto"
)

type User struct {
	Username string
	crypto.PublicKeyHash
	Signer crypto.SigningKeyPair
	Boxer  crypto.BoxingKeyPair
}

func NewUser(username, password string) *User {
	hash256 := sha256.New()
	keySeeds, err := scrypt.Key(hash256.Sum([]byte(password)),
		[]byte(username),
		32768,
		8,
		1,
		64,
	)
	if err != nil {
		return nil
	}

	var secretBoxBytes [32]byte
	var publicBoxBytes [32]byte
	var secretSignBytes [32]byte
	copy(secretBoxBytes[:], keySeeds[:32])
	copy(secretSignBytes[:], keySeeds[32:])

	curve25519.ScalarBaseMult(&publicBoxBytes, &secretBoxBytes)
	publicSignKey, secretSignKey := crypto.Ed25519KeyPair(&secretSignBytes)

	return &User{
		username,
		hash256.Sum(append(publicBoxBytes[:], publicSignKey...)),
		crypto.SigningKeyPair{publicSignKey, secretSignKey},
		crypto.BoxingKeyPair{publicBoxBytes, secretBoxBytes, rand.Reader},
	}
}
