package client

import (
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/pkg/errors"
)

type Auth struct {
	wallet  *hdwallet.Wallet
	account accounts.Account
	Address ethcommon.Address
	TxOpts  *bind.TransactOpts
}

func NewAuth(mnemonic string) (*Auth, error) {
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
		From: account.Address,
		GasLimit: 8000000,
		Signer: func(signer types.Signer, address ethcommon.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != account.Address {
				return nil, errors.New("not authorized to sign this account")
			}

			return wallet.SignTx(account, tx, nil)
		},
	}

	return &Auth{
		wallet:wallet,
		account:account,
		Address: account.Address,
		TxOpts:  txOpts,
	}, nil
}

func (auth *Auth) Sign(hash []byte) ([]byte, error) {
	return auth.wallet.SignHash(auth.account, hash)
}