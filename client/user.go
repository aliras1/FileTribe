package client

import (
	"fmt"
	"io/ioutil"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/glog"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/crypto/sha3"

	"ipfs-share/crypto"
	nw "ipfs-share/networketh"
)

type User struct {
	Address ethcommon.Address
	Name    string
	Signer  *crypto.Signer
	Boxer   crypto.AnonymBoxer
}

func NewUser(username, password, ethKeyPath string) (*User, error) {
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

	json, err := ioutil.ReadFile(ethKeyPath)
	if err != nil {
		return nil, err
	}
	key, err := keystore.DecryptKey(json, password)
	if err != nil {
		return nil, err
	}

	return &User{
		Address: key.Address,
		Name:    username,
		Signer:  &crypto.Signer{PrivateKey: key.PrivateKey},
		Boxer: crypto.AnonymBoxer{
			PublicKey:  crypto.AnonymPublicKey{Value: &publicBoxerBytes},
			PrivateKey: crypto.AnonymPrivateKey{Value: &secretBoxerBytes},
		},
	}, nil
}

func SignUp(username, password, ipfsPeerId, ethKeyPath string, network nw.INetwork) (*User, error) {
	user, err := NewUser(username, password, ethKeyPath)
	if err != nil {
		return nil, fmt.Errorf("could not generate user: SignUp: %s", err)
	}

	exists, err := network.IsUserRegistered(user.Address)
	if err != nil {
		return nil, fmt.Errorf("could not check if username '%s', is registered: SignUp: %s", username, err)
	}
	if exists {
		return nil, fmt.Errorf("username '%s' already exists: SignUp", username)
	}

	pk := ethcrypto.CompressPubkey(&user.Signer.PrivateKey.PublicKey)
	if err = network.RegisterUser(username, ipfsPeerId, *user.Boxer.PublicKey.Value, pk); err != nil {
		return nil, fmt.Errorf("could not register username '%s': SignUp: %s", username, err)
	}

	return user, nil
}

func SignIn(username, password, keyStore string, network nw.INetwork) (*User, error) {
	user, err := NewUser(username, password, keyStore)
	if err != nil {
		return nil, fmt.Errorf("could not generate user: SignIn")
	}

	exists, err := network.IsUserRegistered(user.Address)
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
