package caps

import (
	"encoding/json"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"ipfs-share/crypto"
)

type GroupAccessCap struct {
	Address ethCommon.Address
	Boxer   crypto.SymmetricKey
}

func (cap *GroupAccessCap) Encode() ([]byte, error) {
	data, err := json.Marshal(cap)
	if err != nil {
		return nil, errors.Wrap(err, "could not json marshal group access capability")
	}

	return data, nil
}