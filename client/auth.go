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
	"github.com/aliras1/FileTribe/client/interfaces"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/pkg/errors"
)

// Auth stores all the information of an Ethereum account with which one can sign transactions
type auth struct {
	wallet  *hdwallet.Wallet
	account accounts.Account
	txOpts  *bind.TransactOpts
}

// NewAuth creates an Auth object from a mnemonic
func NewAuth(mnemonic string) (interfaces.Auth, error) {
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return nil, errors.Wrap(err, "could not get wallet from mnemonic")
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0") // path string: Metamask compatible BIP44 derivation
	account, err := wallet.Derive(path, true)
	if err != nil {
		return nil, errors.Wrap(err, "could not derive account from wallet")
	}

	txOpts := &bind.TransactOpts{
		From:     account.Address,
		GasLimit: 8000000,
		Signer: func(signer types.Signer, address ethcommon.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != account.Address {
				return nil, errors.New("not authorized to sign this account")
			}

			return wallet.SignTx(account, tx, nil)
		},
	}

	return &auth{
		wallet:  wallet,
		account: account,
		txOpts:  txOpts,
	}, nil
}

// Sign signs a hash with its Ethereum account's private key
func (auth *auth) Sign(hash []byte) ([]byte, error) {
	return auth.wallet.SignHash(auth.account, hash)
}

// Address returns the address of the account
func (auth *auth) Address() ethcommon.Address {
	return auth.account.Address
}

// TxOpts returns the transaction options
func (auth *auth) TxOpts() *bind.TransactOpts {
	return auth.txOpts
}
