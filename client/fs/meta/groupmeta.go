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