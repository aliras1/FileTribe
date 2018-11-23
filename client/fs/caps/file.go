package caps

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"strings"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"ipfs-share/crypto"
)

type FileCap struct {
	Id              [32]byte
	FileName        string
	IpfsHash        string
	DataKey         crypto.FileBoxer
	WriteAccessList []ethcommon.Address // if empty --> everyone has write access to it
}

func (cap *FileCap) Equal(other *FileCap) bool {
	if !bytes.Equal(cap.Id[:], other.Id[:]) {
		return false
	}

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
	var id [32]byte
	if _, err := rand.Read(id[:]); err != nil {
		return nil, errors.Wrap(err, "could not read from crypto/rand")
	}

	var key [32]byte
	if _, err := rand.Read(key[:]); err != nil {
		return nil, errors.Wrap(err, "could not read from crypto/rand")
	}

	boxer := crypto.FileBoxer{Key: key}

	return &FileCap{
		Id:              id,
		DataKey:         boxer,
		IpfsHash:        "",
		FileName:        fileName,
		WriteAccessList: hasWriteAccess,
	}, nil
}