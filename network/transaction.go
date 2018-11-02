package network

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	. "ipfs-share/collections"
)

type Transaction struct {
	tx *types.Transaction
}

func (t *Transaction) Id() IIdentifier {
	return NewStringId(t.tx.Hash().String())
}

func (t *Transaction) Receipt(network INetwork) (*types.Receipt, error) {
	receipt, err := network.TransactionReceipt(t.tx)
	if err != nil{

	}

	if err != nil {
		return nil, errors.Wrapf(err, "could not get tx receipt")
	}

	return receipt, nil
}