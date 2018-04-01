package filestorage

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"ipfs-share/crypto"
	"ipfs-share/ipfs"
	nw "ipfs-share/network"
)

type GroupAccessCAP struct {
	GroupName string
	Boxer     crypto.SymmetricKey
}

type ReadCAP struct {
	FileName  string                  `json:"name"`
	IPNSPath  string                  `json:"ipns_path"`
	IPFSHash  string                  `json:"ipfs_hash"`
	Owner     string                  `json:"owner"`
	VerifyKey crypto.PublicSigningKey `json:"verify_key"`
}

// Store ReadCAP in json format locally
func (rc *ReadCAP) Store(capPath string) error {
	bytesJSON, err := json.Marshal(rc)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(capPath, bytesJSON, 0644)
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

func (s *Storage) DownloadReadCAP(fromUser, username, capName string, boxer *crypto.BoxingKeyPair, network *nw.Network, ipfs *ipfs.IPFS) (*ReadCAP, error) {
	capBytes, err := s.downloadCAP(s.capsPath, fromUser, username, capName, boxer, network, ipfs)
	if err != nil {
		return nil, err
	}
	var cap *ReadCAP
	if err := json.Unmarshal(capBytes, cap); err != nil {
		return nil, err
	}
	return cap, nil
}

func (s *Storage) DownloadGroupAccessCAP(fromUser, username, capName string, boxer *crypto.BoxingKeyPair, network *nw.Network, ipfs *ipfs.IPFS) (*GroupAccessCAP, error) {
	capBytes, err := s.downloadCAP(s.capsGAPath, fromUser, username, capName, boxer, network, ipfs)
	if err != nil {
		return nil, err
	}
	var cap GroupAccessCAP
	if err := json.Unmarshal(capBytes, &cap); err != nil {
		return nil, err
	}
	return &cap, nil
}

// Downloads the capability identified by capName from
// /ipns/from/for/username/capName
func (s *Storage) downloadCAP(basePath, fromUser, username, capName string, boxer *crypto.BoxingKeyPair, network *nw.Network, ipfs *ipfs.IPFS) ([]byte, error) {
	// get address and key
	ipfsAddr, err := network.GetUserIPFSAddr(fromUser)
	if err != nil {
		return nil, err
	}
	otherPK, err := network.GetUserBoxingKey(fromUser)
	if err != nil {
		return nil, err
	}
	ipnsPath := "/ipns/" + ipfsAddr + "/for/" + username + "/" + capName
	// download cap file
	tmpFilePath := s.tmpPath + "/" + capName
	err = ipfs.Get(tmpFilePath, ipnsPath)
	if err != nil {
		return nil, err
	}
	capFilePath := basePath + "/" + capName
	bytesEnc, err := ioutil.ReadFile(tmpFilePath)
	if err != nil {
		return nil, err
	}

	bytesDecr, success := boxer.BoxOpen(bytesEnc, &otherPK)
	if !success {
		return nil, errors.New("could not decrypt capability")
	}
	os.Remove(tmpFilePath)
	if err := WriteFile(capFilePath, bytesDecr); err != nil {
		return nil, err
	}
	return bytesDecr, nil
}

// Checks if the by ReadCap represented file has changed since last time
// or not. It is done via checking the IPFS hash of the file. If it has
// changed, the function returns true. Otherwise it returns false.
func (rc *ReadCAP) RefreshCAP(storage *Storage, ipfs *ipfs.IPFS) (bool, error) {
	resolvedHash, err := ipfs.Resolve(rc.IPNSPath)
	fileChanged := strings.Compare(resolvedHash, rc.IPFSHash) != 0
	if err != nil {
		return false, err
	}
	if fileChanged {
		rc.IPFSHash = resolvedHash
		rc.Store(storage.capsPath + "/" + rc.FileName + ".json")
		return true, nil
	}
	return false, nil
}
