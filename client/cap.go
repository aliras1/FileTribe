package client

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	ipfsapi "ipfs-share/ipfs"

	"ipfs-share/crypto"
	"github.com/pkg/errors"
	"bytes"
	"strings"
	"crypto/rand"
	"os"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"ipfs-share/utils"
)

type GroupAccessCap struct {
	GroupId [32]byte
	Boxer   crypto.SymmetricKey
}

func (cap *GroupAccessCap) Save(storage *Storage) error {
	capJson, err := json.Marshal(cap)
	if err != nil {
		return fmt.Errorf("could not marshal group access capability: GroupAccessCap.SaveMetadata: %s", err)
	}

	groupIdBase64 := base64.URLEncoding.EncodeToString(cap.GroupId[:])
	path := storage.GroupAccessCapDir() + groupIdBase64 + CAP_EXT
	if err := utils.WriteFile(path, capJson); err != nil {
		return errors.Wrap(err, "could not write group cap file")
	}

	return nil
}


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


func NewGroupFileCap(fileName string, filePath string, hasWriteAccess []ethcommon.Address, ipfs ipfsapi.IIpfs, storage *Storage) (*FileCap, error) {
	var id [32]byte
	if _, err := rand.Read(id[:]); err != nil {
		return nil, errors.Wrap(err, "could not read from crypto/rand")
	}

	var key [32]byte
	if _, err := rand.Read(key[:]); err != nil {
		return nil, errors.Wrap(err, "could not read from crypto/rand")
	}

	boxer := crypto.FileBoxer{Key: key}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "could not open file")
	}

 	encBuffer, err := boxer.Seal(file)
 	if err != nil {
 		return nil, errors.Wrap(err, "could not encrypt file")
	}

	hash, err := ipfs.Add(encBuffer)
	if err != nil {
		return nil, errors.Wrap(err, "could not ipfs add file")
	}

	return &FileCap{
		Id:              id,
		DataKey:         boxer,
		IpfsHash:        hash,
		FileName:        fileName,
		WriteAccessList: hasWriteAccess,
	}, nil
}
