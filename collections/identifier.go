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

package collections

import (
	"bytes"
	"encoding/base64"
	"strings"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

type IIdentifier interface {
	Equal(other IIdentifier) bool
	Data() interface{}
	ToString() string
}

type StringId struct {
	data string
}

func (id *StringId) Equal(other IIdentifier) bool {
	castedOther := other.(*StringId)
	return strings.Compare(id.data, castedOther.data) == 0
}

func (id *StringId) Data() interface{} {
	return id.data
}

func (id *StringId) ToString() string {
	return id.data
}

func NewStringId(id string) IIdentifier {
	return &StringId{data: id}
}

type Uint32Id struct {
	data uint32
}

func NewUint32Id(id uint32) *Uint32Id {
	return &Uint32Id{data: id}
}

func (id *Uint32Id) ToString() string {
	return string(id.data)
}

func (id *Uint32Id) Equal(other IIdentifier) bool {
	castedOther := other.(*Uint32Id)
	return id.data == castedOther.data
}

func (id *Uint32Id) Data() interface{} {
	return id.data
}

type BytesId struct {
	data [32]byte
}

func (id *BytesId) ToString() string {
	return base64.URLEncoding.EncodeToString(id.data[:])
}

func NewBytesId(id [32]byte) *BytesId {
	return &BytesId{data: id}
}

func (id *BytesId) Equal(other IIdentifier) bool {
	castedOther := other.(*BytesId)
	return bytes.Equal(id.data[:], castedOther.data[:])
}

func (id *BytesId) Data() interface{} {
	return id.data
}

type AddressId struct {
	data *ethcommon.Address
}

func (id *AddressId) ToString() string {
	return id.data.String()
}

func NewAddressId(address *ethcommon.Address) *AddressId {
	return &AddressId{data: address}
}

func (id *AddressId) Equal(other IIdentifier) bool {
	castedOther := other.(*AddressId)
	return bytes.Equal(id.data.Bytes(), castedOther.data.Bytes())
}

func (id *AddressId) Data() interface{} {
	return id.data
}