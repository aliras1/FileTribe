package meta

import (
	"encoding/json"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/aliras1/FileTribe/tribecrypto"
)

type GroupMeta struct {
	Address ethCommon.Address
	Boxer   tribecrypto.SymmetricKey
}

func (meta *GroupMeta) Encode() ([]byte, error) {
	data, err := json.Marshal(meta)
	if err != nil {
		return nil, errors.Wrap(err, "could not json marshal group access capability")
	}

	return data, nil
}