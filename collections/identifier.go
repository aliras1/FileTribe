package collections

import (
	"bytes"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"encoding/base64"
)

type IIdentifier interface {
	Equal(other IIdentifier) bool
	Data() interface{}
	ToString() string
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