package filestorage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"ipfs-share/crypto"
	"ipfs-share/ipfs"
	nw "ipfs-share/network"
	"os"
	"path"
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
	IPNSPath  string                  `json:"ipns_path"`
	Owner     string                  `json:"owner"`
	VerifyKey crypto.PublicSigningKey `json:"verify_key"`
}

// Store ReadCAP in json format locally
func (cap *ReadCAP) Store(storage *Storage) error {
	bytesJSON, err := json.Marshal(cap)
	if err != nil {
		return fmt.Errorf("could not marshal cap '%s': ReadCAP.Store: %s", cap.FileName, err)
	}
	capPath := storage.capsPath + "/" + cap.FileName + ".json"
	if err := WriteFile(capPath, bytesJSON); err != nil {
		return fmt.Errorf("could not write file '%s': ReadCAP.Store: %s", capPath, err)
	}
	return nil
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

func (cap *ReadCAP) Refresh(username string, boxer *crypto.BoxingKeyPair, storage *Storage, network *nw.Network, ipfs *ipfs.IPFS) (bool, error) {
	return false, fmt.Errorf("not implemented: ReadCAP.Refresh")
}

func DownloadReadCAP(fromUser, username, capName string, boxer *crypto.BoxingKeyPair, storage *Storage, network *nw.Network, ipfs *ipfs.IPFS) (*ReadCAP, error) {
	capBytes, err := downloadCAP(fromUser, username, capName, boxer, storage, network, ipfs)
	if err != nil {
		return nil, fmt.Errorf("could not download ReadCAP '%s': DownloadReadCAP: %s", capName, err)
	}
	var cap ReadCAP
	if err := json.Unmarshal(capBytes, &cap); err != nil {
		return nil, fmt.Errorf("could not unmarsharl ReadCAP '%s': DownloadReadCAP: %s", capName, err)
	}
	return &cap, nil
}

// Downloads the capability identified by capName from
// /ipns/from/for/username/capName
func downloadCAP(fromUser, username, capName string, boxer *crypto.BoxingKeyPair, storage *Storage, network *nw.Network, ipfs *ipfs.IPFS) ([]byte, error) {
	// get address and key
	ipfsAddr, err := network.GetUserIPFSAddr(fromUser)
	if err != nil {
		return nil, fmt.Errorf("could not get ipfs address of user '%s': downloadCAP: %s", username, err)
	}
	otherPK, err := network.GetUserBoxingKey(fromUser)
	if err != nil {
		return nil, fmt.Errorf("could not get public boxing key of user '%s': downloadCAP: %s", username, err)
	}
	ipnsPath := "/ipns/" + ipfsAddr + "/for/" + username + "/" + capName
	// download cap file
	tmpFilePath := storage.tmpPath + "/" + capName
	err = ipfs.Get(tmpFilePath, ipnsPath)
	if err != nil {
		return nil, fmt.Errorf("could not ipfs get file '%s': downloadCAP: %s", ipnsPath, err)
	}
	bytesEnc, err := ioutil.ReadFile(tmpFilePath)
	if err != nil {
		return nil, fmt.Errorf("could not read file '%s': downloadCAP: %s", tmpFilePath, err)
	}
	fmt.Print("group cap bytes: ")
	fmt.Println(bytesEnc)
	bytesDecr, success := boxer.BoxOpen(bytesEnc, &otherPK)
	if !success {
		fmt.Println("trying decrypt cap")
		fmt.Println(fromUser)
		fmt.Println(ipnsPath)
		return nil, fmt.Errorf("could not decrypt cap '%s': downloadCAP: %s", capName, err)
	}
	os.Remove(tmpFilePath)
	return bytesDecr, nil
}

func CreateFileReadCAPForUser(f *FilePTP, username, ipnsAddr string, boxer *crypto.BoxingKeyPair, storage *Storage, network *nw.Network) error {
	cap := ReadCAP{path.Base(f.Name), ipnsAddr, f.Owner, f.VerifyKey}
	capBytes, err := json.Marshal(cap)
	if err != nil {
		return fmt.Errorf("could not marshal cap '%s': CreateFileReadCAPForUser: '%s'", cap.FileName, err)
	}
	if err := storage.createFileForUser(username, path.Base(f.Name), capBytes, boxer, network); err != nil {
		return fmt.Errorf("could not create file '%s' for user '%s': CreateFileReadCAPForUser: %s", cap.FileName, username, err)
	}
	return nil
}
