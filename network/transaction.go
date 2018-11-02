package network

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

type Transaction struct {
	tx *types.Transaction
}

func (t *Transaction) Status(network INetwork) (uint64, error) {
	receipt, err := network.TransactionReceipt(t.tx)
	if err != nil{

	}

	if err != nil {
		return 255, errors.Wrapf(err, "could not get tx receipt")
	}

	return receipt.Status, nil
}