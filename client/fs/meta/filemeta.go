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

type FileMeta struct {
	FileName        string
	IpfsHash        string
	DataKey         tribecrypto.FileBoxer
	WriteAccessList []ethcommon.Address // if empty --> everyone has write access to it
}

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


func (meta *FileMeta) Encode() ([]byte, error) {
	data, err := json.Marshal(meta)
	if err != nil {
		return nil, errors.Wrap(err, "could not json marshal: FileMeta.Encode")
	}

	return data, nil
}

func EncodeFileMetaList(lst []*FileMeta) ([]byte, error) {
	data, err := json.Marshal(lst)
	if err != nil {
		return nil, errors.Wrap(err, "could not json marshal: FileCapList")
	}

	return data, nil
}

func DecodeFileCap(data []byte) (*FileMeta, error) {
	var cap FileMeta
	if err := json.Unmarshal(data, &cap); err != nil {
		return nil, errors.Wrap(err, "could not json unmarshal: FileMeta:")
	}

	return &cap, nil
}

func DecodeFileCapList(data []byte) ([]*FileMeta, error) {
	var cap []*FileMeta
	if err := json.Unmarshal(data, &cap); err != nil {
		return nil, errors.Wrap(err, "could not json unmarshal: FileMeta:")
	}

	return cap, nil
}


func NewGroupFileCap(fileName string, hasWriteAccess []ethcommon.Address) (*FileMeta, error) {
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