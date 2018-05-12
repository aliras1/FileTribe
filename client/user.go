package client

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/scrypt"
	"github.com/golang/glog"

	"ipfs-share/crypto"
	nw "ipfs-share/network"
)

type User struct {
	Name   string
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
		glog.Error("error while scrypt: NewUser: %s", err)
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
		Name:username,
		PublicKeyHash:  hash256.Sum(append(publicBoxBytes[:], publicSignKey...)),
		Signer:  crypto.SigningKeyPair{
			 PublicKey: publicSignKey,
			 SecretKey:  secretSignKey,
		},
		Boxer:  crypto.BoxingKeyPair{
			PublicKey:   publicBoxBytes,
			SecretKey:  secretBoxBytes,
			RNG:  rand.Reader,
		},
	}
}

func SignUp(username, password, ipfsAddr string, network *nw.Network) (*User, error) {
	exists, err := network.IsUsernameRegistered(username)
	if err != nil {
		return nil, fmt.Errorf("could not check if username '%s', is registered: SignUp: %s", username, err)
	}
	if exists {
		return nil, fmt.Errorf("username '%s' already exists: SignUp", username)
	}
	user := NewUser(username, password)
	if user == nil {
		return nil, fmt.Errorf("could not generate user: SignUp")
	}
	if err = network.RegisterUsername(username, user.PublicKeyHash); err != nil {
		return nil, fmt.Errorf("could not register username '%s': SignUp: %s", username, err)
	}
	if err := network.PutVerifyKey(user.PublicKeyHash, user.Signer.PublicKey); err != nil {
		return nil, fmt.Errorf("could not put verify key: SignUp: %s", err)
	}
	if err := network.PutBoxingKey(user.PublicKeyHash, user.Boxer.PublicKey); err != nil {
		return nil, fmt.Errorf("could not put boxing key: SignUp: %s", err)
	}
	if err := network.PutIPFSAddr(user.PublicKeyHash, ipfsAddr); err != nil {
		return nil, fmt.Errorf("could not put ipfs address: SignUp: %s", err)
	}

	return user, nil
}

func SignIn(username, password string, network *nw.Network) (*User, error) {
	exists, err := network.IsUsernameRegistered(username)
	if err != nil {
		return nil, fmt.Errorf("could not check if username '%s' is registered: SignIn: %s", username, err)
	}
	if !exists {
		return nil, fmt.Errorf("username '%s' does not exists: SignIn", username)
	}
	user := NewUser(username, password)
	if user == nil {
		return nil, fmt.Errorf("could not generate user: SignIn")
	}
	publicKeyHash, err := network.GetUserPublicKeyHash(username)
	if err != nil {
		return nil, fmt.Errorf("could not get user public key hash: SignIn: %s", err)
	}
	if !publicKeyHash.Equals(&user.PublicKeyHash) {
		return nil, fmt.Errorf("incorrect password: SignIn")
	}

	return user, nil
}

func (u *User) SignTransaction(transaction *Transaction) []byte {
	return u.Signer.SecretKey.Sign(nil, transaction.Bytes())[:64]
}
