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

package fs

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/pkg/errors"
	"github.com/sergi/go-diff/diffmatchpatch"

	"github.com/aliras1/FileTribe/tribecrypto"
)

type DiffNode struct {
	Hash      []byte
	Diff      []diffmatchpatch.Diff
	Next      string
	NextBoxer tribecrypto.FileBoxer
}

func (diff *DiffNode) Encode() ([]byte, error) {
	data, err := json.Marshal(diff)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func DecodeDiffNode(data []byte) (*DiffNode, error) {
	var diff DiffNode
	if err := json.Unmarshal(data, &diff); err != nil {
		return nil, err
	}
	return &diff, nil
}

func (diff *DiffNode) Encrypt(boxer tribecrypto.FileBoxer) (io.Reader, error) {
	data, err := diff.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode diff node")
	}
	encData, err := boxer.Seal(bytes.NewReader(data))
	if err != nil {
		return nil, errors.Wrap(err, "could not encrypt data")
	}

	return encData, nil
}