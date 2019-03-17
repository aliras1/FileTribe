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