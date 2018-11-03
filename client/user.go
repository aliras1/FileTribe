package client

import (
	"io/ioutil"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/golang/glog"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/crypto/sha3"

	"ipfs-share/crypto"
)

type IUser interface {
	Address() ethcommon.Address
	Name() string
	Signer() crypto.Signer
	Boxer() crypto.AnonymBoxer
}

type User struct {
	address ethcommon.Address
	name    string
	signer  *crypto.Signer
	boxer   crypto.AnonymBoxer
}

func (user *User) Address() ethcommon.Address {
	return user.address
}

func (user *User) Name() string {
	return user.name
}

func (user *User) Signer() crypto.Signer {
	return *user.signer
}

func (user *User) Boxer() crypto.AnonymBoxer {
	return user.boxer
}

func NewUser(username, password, ethKeyPath string) (IUser, error) {
	passwordDigest := sha3.Sum256([]byte(password))
	keySeeds, err := scrypt.Key(
		passwordDigest[:],
		[]byte(username),
		32768,
		8,
		1,
		32,
	)
	if err != nil {
		glog.Error("error while scrypt: NewUser: %s", err)
		return nil, err
	}

	var secretBoxerBytes [32]byte
	var publicBoxerBytes [32]byte
	copy(secretBoxerBytes[:], keySeeds)
	curve25519.ScalarBaseMult(&publicBoxerBytes, &secretBoxerBytes)

	ethKeyData, err := ioutil.ReadFile(ethKeyPath)
	if err != nil {
		return nil, err
	}
	key, err := keystore.DecryptKey(ethKeyData, password)
	if err != nil {
		return nil, err
	}

	return &User{
		address: key.Address,
		name:    username,
		signer:  &crypto.Signer{PrivateKey: key.PrivateKey},
		boxer: crypto.AnonymBoxer{
			PublicKey:  crypto.AnonymPublicKey{Value: publicBoxerBytes},
			PrivateKey: crypto.AnonymPrivateKey{Value: secretBoxerBytes},
		},
	}, nil
}
