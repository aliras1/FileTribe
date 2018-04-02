package filestorage

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"ipfs-share/crypto"
	"ipfs-share/ipfs"
)

type GroupAccessCAP struct {
	GroupName string
	Boxer     crypto.SymmetricKey
}

func (cap *GroupAccessCAP) Store(capPath string) error {
	bytesJSON, err := json.Marshal(cap)
	if err != nil {
		return err
	}
	return WriteFile(capPath, bytesJSON)
}

type ReadCAP struct {
	FileName  string                  `json:"name"`
	IPNSPath  string                  `json:"ipns_path"`
	IPFSHash  string                  `json:"ipfs_hash"`
	Owner     string                  `json:"owner"`
	VerifyKey crypto.PublicSigningKey `json:"verify_key"`
}

// Store ReadCAP in json format locally
func (cap *ReadCAP) Store(capPath string) error {
	bytesJSON, err := json.Marshal(cap)
	if err != nil {
		return err
	}
	return WriteFile(capPath, bytesJSON)
}

// New ReadCAP from local json file
func NewReadCAPFromFile(capPath string) (*ReadCAP, error) {
	bytesFile, err := ioutil.ReadFile(capPath)
	if err != nil {
		return nil, err
	}
	var cap ReadCAP
	err = json.Unmarshal(bytesFile, &cap)
	return &cap, err
}

// Checks if the by ReadCap represented file has changed since last time
// or not. It is done via checking the IPFS hash of the file. If it has
// changed, the function returns true. Otherwise it returns false.
func (cap *ReadCAP) Refresh(storage *Storage, ipfs *ipfs.IPFS) (bool, error) {
	resolvedHash, err := ipfs.Resolve(cap.IPNSPath)
	fileChanged := strings.Compare(resolvedHash, cap.IPFSHash) != 0
	if err != nil {
		return false, err
	}
	if fileChanged {
		cap.IPFSHash = resolvedHash
		cap.Store(storage.capsPath + "/" + cap.FileName + ".json")
		return true, nil
	}
	return false, nil
}
