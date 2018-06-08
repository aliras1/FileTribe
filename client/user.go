package client

import (
	"fmt"

	"github.com/golang/glog"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/crypto/sha3"

	"ipfs-share/crypto"
	nw "ipfs-share/networketh"
)

type User struct {
	Name   string
	ID     [32]byte
	Signer crypto.SigningKeyPair
	Boxer  crypto.AnonymBoxer
}

func NewUser(username, password string) *User {
	passwordDigest := sha3.Sum256([]byte(password))
	keySeeds, err := scrypt.Key(
		passwordDigest[:],
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

	var secretBoxerBytes [32]byte
	var publicBoxerBytes [32]byte
	var signingKeyBytes [32]byte
	copy(secretBoxerBytes[:], keySeeds[:32])
	copy(signingKeyBytes[:], keySeeds[32:])

	curve25519.ScalarBaseMult(&publicBoxerBytes, &secretBoxerBytes)
	verifyKey, signingKey := crypto.Ed25519KeyPair(&signingKeyBytes)

	return &User{
		Name: username,
		ID:   sha3.Sum256(append(publicBoxerBytes[:], verifyKey...)),
		Signer: crypto.SigningKeyPair{
			SigningKey: signingKey,
			VerifyKey:  verifyKey,
		},
		Boxer: crypto.AnonymBoxer{
			PublicKey: crypto.AnonymPublicKey{Value: &publicBoxerBytes},
			SecretKey: crypto.AnonymSecretKey{Value: &secretBoxerBytes},
		},
	}
}

func SignUp(username, password, ipfsAddr string, network *nw.Network) (*User, error) {
	user := NewUser(username, password)
	if user == nil {
		return nil, fmt.Errorf("could not generate user: SignUp")
	}

	exists, err := network.IsUserRegistered(user.ID)
	if err != nil {
		return nil, fmt.Errorf("could not check if username '%s', is registered: SignUp: %s", username, err)
	}
	if exists {
		return nil, fmt.Errorf("username '%s' already exists: SignUp", username)
	}

	if err = network.RegisterUser(user.ID, username, *user.Boxer.PublicKey.Value, *user.Signer.VerifyKey.Bytes(), ipfsAddr); err != nil {
		return nil, fmt.Errorf("could not register username '%s': SignUp: %s", username, err)
	}

	return user, nil
}

func SignIn(username, password string, network *nw.Network) (*User, error) {
	user := NewUser(username, password)
	if user == nil {
		return nil, fmt.Errorf("could not generate user: SignIn")
	}

	exists, err := network.IsUserRegistered(user.ID)
	if err != nil {
		return nil, fmt.Errorf("could not check if username '%s' is registered: SignIn: %s", username, err)
	}
	if !exists {
		return nil, fmt.Errorf("username '%s' does not exists: SignIn", username)
	}

	return user, nil
}

func (u *User) SignTransaction(transaction *Transaction) []byte {
	// return u.Signer.VerifyKey.Sign(transaction.Bytes())[:64]
	return nil
}
