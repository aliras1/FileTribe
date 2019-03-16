package client

import (
	"io/ioutil"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/aliras1/FileTribe/tribecrypto"
)

type Auth struct {
	Address ethcommon.Address
	Signer  *tribecrypto.Signer
	TxOpts  *bind.TransactOpts
}

func NewAuth(ethKeyPath string, password string) (*Auth, error) {

	ethKeyData, err := ioutil.ReadFile(ethKeyPath)
	if err != nil {
		return nil, err
	}

	key, err := keystore.DecryptKey(ethKeyData, password)
	if err != nil {
		return nil, err
	}

	txOpts := bind.NewKeyedTransactor(key.PrivateKey)

	return &Auth{
		Address: key.Address,
		Signer:  &tribecrypto.Signer{PrivateKey: key.PrivateKey},
		TxOpts:  txOpts,
	}, nil
}