// Copyright (c) 2019 Laszlo Sari
//
// FileTribe is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// FileTribe is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//

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

func (acc *Account) SetContract(contract *ethacc.Account) {
	acc.contract = contract
}

func (acc *Account) SetContractAddress(addr ethcommon.Address) {
	acc.data.ContractAddress = addr
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
