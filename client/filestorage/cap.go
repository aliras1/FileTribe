package filestorage

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"errors"
	"fmt"
	"ipfs-share/crypto"
	"ipfs-share/ipfs"
	nw "ipfs-share/network"
	"os"
)

type GroupAccessCAP struct {
	GroupName string
	Boxer     crypto.SymmetricKey
}

func (cap *GroupAccessCAP) Store(storage *Storage) error {
	bytesJSON, err := json.Marshal(cap)
	if err != nil {
		return err
	}
	capPath := storage.capsGAPath + "/" + cap.GroupName + ".json"
	return WriteFile(capPath, bytesJSON)
}

type ReadCAP struct {
	FileName  string                  `json:"name"`
	IPFSHash  string                  `json:"ipfs_hash"`
	Owner     string                  `json:"owner"`
	VerifyKey crypto.PublicSigningKey `json:"verify_key"`
}

// Store ReadCAP in json format locally
func (cap *ReadCAP) Store(storage *Storage) error {
	bytesJSON, err := json.Marshal(cap)
	if err != nil {
		return err
	}
	capPath := storage.capsPath + "/" + cap.FileName + ".json"
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
func (cap *ReadCAP) Refresh(username string, boxer *crypto.BoxingKeyPair, storage *Storage, network *nw.Network, ipfs *ipfs.IPFS) (bool, error) {
	newCap, err := DownloadReadCAP(cap.Owner, username, cap.FileName+".json", boxer, storage, network, ipfs)
	if err != nil {
		return false, err
	}
	fileChanged := strings.Compare(newCap.IPFSHash, cap.IPFSHash) != 0
	if fileChanged {
		cap = newCap
		cap.Store(storage)
		return true, nil
	}
	return false, nil
}

func DownloadReadCAP(fromUser, username, capName string, boxer *crypto.BoxingKeyPair, storage *Storage, network *nw.Network, ipfs *ipfs.IPFS) (*ReadCAP, error) {
	capBytes, err := downloadCAP(fromUser, username, capName, boxer, storage, network, ipfs)
	if err != nil {
		return nil, err
	}
	var cap ReadCAP
	if err := json.Unmarshal(capBytes, &cap); err != nil {
		return nil, err
	}
	return &cap, nil
}

// Downloads the capability identified by capName from
// /ipns/from/for/username/capName
func downloadCAP(fromUser, username, capName string, boxer *crypto.BoxingKeyPair, storage *Storage, network *nw.Network, ipfs *ipfs.IPFS) ([]byte, error) {
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
	tmpFilePath := storage.tmpPath + "/" + capName
	err = ipfs.Get(tmpFilePath, ipnsPath)
	if err != nil {
		return nil, err
	}
	bytesEnc, err := ioutil.ReadFile(tmpFilePath)
	if err != nil {
		return nil, err
	}
	fmt.Print("group cap bytes: ")
	fmt.Println(bytesEnc)
	bytesDecr, success := boxer.BoxOpen(bytesEnc, &otherPK)
	if !success {
		fmt.Println("trying decrypt cap")
		fmt.Println(fromUser)
		fmt.Println(ipnsPath)
		return nil, errors.New("could not decrypt capability")
	}
	os.Remove(tmpFilePath)
	// TODO  ATTENTION: it does not save anything
	//if err := WriteFile(capFilePath, bytesDecr); err != nil {
	//	return nil, err
	//}
	return bytesDecr, nil
}
