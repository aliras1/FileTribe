package client

import (
	"crypto/rand"
	"crypto/sha256"

	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/scrypt"

	"errors"
	"ipfs-share/crypto"
	nw "ipfs-share/network"
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

func SignUp(username, password, ipfsAddr string, network *nw.Network) (*User, error) {
	exists, err := network.IsUsernameRegistered(username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("user already exists")
	}
	user := NewUser(username, password)
	if user == nil {
		return nil, errors.New("could not generate user")
	}
	err = network.RegisterUsername(username, user.PublicKeyHash)
	if err != nil {
		return nil, err
	}
	network.PutSigningKey(user.PublicKeyHash, user.Signer.PublicKey)
	network.PutBoxingKey(user.PublicKeyHash, user.Boxer.PublicKey)
	network.PutIPFSAddr(user.PublicKeyHash, ipfsAddr)
	return user, nil
}

func SignIn(username, password string, network *nw.Network) (*User, error) {
	exists, err := network.IsUsernameRegistered(username)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("username does not exists")
	}
	user := NewUser(username, password)
	if user == nil {
		return nil, errors.New("could not generate user")
	}
	publicKeyHash, err := network.GetUserPublicKeyHash(username)
	if err != nil {
		return nil, err
	}
	if !publicKeyHash.Equals(&user.PublicKeyHash) {
		return nil, errors.New("incorrect password")
	}
	return user, nil
}

func (u *User) SignTransaction(transaction *Transaction) []byte {
	return u.Signer.SecretKey.Sign(nil, transaction.Bytes())[:64]
}
