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

package caps

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"strings"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/aliras1/FileTribe/tribecrypto"
)

type FileCap struct {
	FileName        string
	IpfsHash        string
	DataKey         tribecrypto.FileBoxer
	WriteAccessList []ethcommon.Address // if empty --> everyone has write access to it
}

func (cap *FileCap) Equal(other *FileCap) bool {
	if strings.Compare(cap.FileName, other.FileName) != 0 {
		return false
	}

	if strings.Compare(cap.IpfsHash, other.IpfsHash) != 0 {
		return false
	}

	if !bytes.Equal(cap.DataKey.Key[:], other.DataKey.Key[:]) {
		return false
	}

	if len(cap.WriteAccessList) != len(other.WriteAccessList) {
		return false
	}

	for i := 0; i < len(cap.WriteAccessList); i++ {
		if !bytes.Equal(cap.WriteAccessList[i].Bytes(), other.WriteAccessList[i].Bytes()) {
			return false
		}
	}

	return true
}


func (cap *FileCap) Encode() ([]byte, error) {
	data, err := json.Marshal(cap)
	if err != nil {
		return nil, errors.Wrap(err, "could not json marshal: FileCap.Encode")
	}

	return data, nil
}

func EncodeFileCapList(lst []*FileCap) ([]byte, error) {
	data, err := json.Marshal(lst)
	if err != nil {
		return nil, errors.Wrap(err, "could not json marshal: FileCapList")
	}

	return data, nil
}

func DecodeFileCap(data []byte) (*FileCap, error) {
	var cap FileCap
	if err := json.Unmarshal(data, &cap); err != nil {
		return nil, errors.Wrap(err, "could not json unmarshal: FileCap:")
	}

	return &cap, nil
}

func DecodeFileCapList(data []byte) ([]*FileCap, error) {
	var cap []*FileCap
	if err := json.Unmarshal(data, &cap); err != nil {
		return nil, errors.Wrap(err, "could not json unmarshal: FileCap:")
	}

	return cap, nil
}


func NewGroupFileCap(fileName string, hasWriteAccess []ethcommon.Address) (*FileCap, error) {
	var key [32]byte
	if _, err := rand.Read(key[:]); err != nil {
		return nil, errors.Wrap(err, "could not read from crypto/rand")
	}

	boxer := tribecrypto.FileBoxer{Key: key}

	return &FileCap{
		DataKey:         boxer,
		IpfsHash:        "",
		FileName:        fileName,
		WriteAccessList: hasWriteAccess,
	}, nil
}