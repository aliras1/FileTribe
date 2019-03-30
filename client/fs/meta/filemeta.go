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
	"bytes"
	"crypto/rand"
	"encoding/json"
	"strings"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/aliras1/FileTribe/tribecrypto"
)

// FileMeta stores all data that is necessary to reach and read a file from IPFS
type FileMeta struct {
	FileName        string
	IpfsHash        string
	DataKey         tribecrypto.FileBoxer
	WriteAccessList []ethcommon.Address // if empty --> everyone has write access to it
}

// Equal decides if two files are identical to each other or not
func (meta *FileMeta) Equal(other *FileMeta) bool {
	if strings.Compare(meta.FileName, other.FileName) != 0 {
		return false
	}

	if strings.Compare(meta.IpfsHash, other.IpfsHash) != 0 {
		return false
	}

	if !bytes.Equal(meta.DataKey.Key[:], other.DataKey.Key[:]) {
		return false
	}

	if len(meta.WriteAccessList) != len(other.WriteAccessList) {
		return false
	}

	for i := 0; i < len(meta.WriteAccessList); i++ {
		if !bytes.Equal(meta.WriteAccessList[i].Bytes(), other.WriteAccessList[i].Bytes()) {
			return false
		}
	}

	return true
}

// Encode encodes a file meta
func (meta *FileMeta) Encode() ([]byte, error) {
	data, err := json.Marshal(meta)
	if err != nil {
		return nil, errors.Wrap(err, "could not json marshal: FileMeta.Encode")
	}

	return data, nil
}

// EncodeFileMetaList encodes a list of file metas
func EncodeFileMetaList(lst []*FileMeta) ([]byte, error) {
	data, err := json.Marshal(lst)
	if err != nil {
		return nil, errors.Wrap(err, "could not json marshal: FileMetaList")
	}

	return data, nil
}

// DecodeFileMeta decodes a single file meta
func DecodeFileMeta(data []byte) (*FileMeta, error) {
	var meta FileMeta
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, errors.Wrap(err, "could not json unmarshal: FileMeta:")
	}

	return &meta, nil
}

// DecodeFileMetaList decodes a list of file metas
func DecodeFileMetaList(data []byte) ([]*FileMeta, error) {
	var meta []*FileMeta
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, errors.Wrap(err, "could not json unmarshal: FileMeta:")
	}

	return meta, nil
}

// NewFileMeta creates a new file meta
func NewFileMeta(fileName string, hasWriteAccess []ethcommon.Address) (*FileMeta, error) {
	var key [32]byte
	if _, err := rand.Read(key[:]); err != nil {
		return nil, errors.Wrap(err, "could not read from crypto/rand")
	}

	boxer := tribecrypto.FileBoxer{Key: key}

	return &FileMeta{
		DataKey:         boxer,
		IpfsHash:        "",
		FileName:        fileName,
		WriteAccessList: hasWriteAccess,
	}, nil
}
