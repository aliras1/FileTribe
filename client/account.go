package client

import (
	"crypto/rand"
	"encoding/json"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/chequebook"
	"github.com/pkg/errors"
	"golang.org/x/crypto/curve25519"

	"github.com/aliras1/FileTribe/client/fs"
	"github.com/aliras1/FileTribe/client/interfaces"
	ethacc "github.com/aliras1/FileTribe/eth/gen/Account"
	"github.com/aliras1/FileTribe/tribecrypto"
)

type AccountData struct {
	Name            string
	ContractAddress ethcommon.Address
	Boxer           tribecrypto.AnonymBoxer
}

type Account struct {
	data     *AccountData
	contract *ethacc.Account
	storage  *fs.Storage
}


func (acc *Account) ContractAddress() ethcommon.Address {
	return acc.data.ContractAddress
}

func (acc *Account) Name() string {
	return acc.data.Name
}


func (acc *Account) Boxer() tribecrypto.AnonymBoxer {
	return acc.data.Boxer
}

func (acc *Account) Contract() *ethacc.Account {
	return acc.contract
}

func (acc *Account) SetContract(address ethcommon.Address, backend chequebook.Backend) error {
	contract, err := ethacc.NewAccount(address, backend)
	if err != nil {
		return err
	}

	acc.data.ContractAddress = address
	acc.contract = contract

	return nil
}

func (acc *Account) Save() error {
	dataEncoded, err := json.Marshal(acc.data)
	if err != nil {
		return errors.Wrap(err, "could not json encode account data")
	}

	if err := acc.storage.SaveAccountData(dataEncoded); err != nil {
		return errors.Wrap(err, "storage could not save account data")
	}

	return nil
}

func NewAccountFromStorage(storage *fs.Storage, backend chequebook.Backend) (interfaces.IAccount, error) {
	dataEncoded, err := storage.LoadAccountData()
	if err != nil {
		return nil, errors.Wrap(err, "could not load account data")
	}

	var accData AccountData
	if err := json.Unmarshal(dataEncoded, &accData); err != nil {
		return nil, errors.Wrap(err, "could not json decode account data")
	}

	contract, err := ethacc.NewAccount(accData.ContractAddress, backend)
	if err != nil {
		return nil, errors.Wrap(err, "could not create account contract instance")
	}

	acc := &Account{
		storage:	storage,
		data:		&accData,
		contract:	contract,
	}

	return acc, nil
}

func NewAccount(username string, storage *fs.Storage) (interfaces.IAccount, error) {
	var secretBoxerBytes [32]byte
	var publicBoxerBytes [32]byte

	if _, err := rand.Read(secretBoxerBytes[:]); err != nil {
		return nil, errors.Wrap(err, "could not read from crypto.rand")
	}

	curve25519.ScalarBaseMult(&publicBoxerBytes, &secretBoxerBytes)

	return &Account{
		data: &AccountData{
			Name:username,
			Boxer: tribecrypto.AnonymBoxer{
				PublicKey:  tribecrypto.AnonymPublicKey{Value: publicBoxerBytes},
				PrivateKey: tribecrypto.AnonymPrivateKey{Value: secretBoxerBytes},
			},
		},
		storage:storage,
	}, nil
}
