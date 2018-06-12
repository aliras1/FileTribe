package filestorage

import (
	"github.com/ethereum/go-ethereum/common"
	"encoding/base64"
	"encoding/json"
	"fmt"
	// "io/ioutil"
	// "os"
	"path"

	// "github.com/golang/glog"

	"ipfs-share/crypto"
	"ipfs-share/ipfs"
	nw "ipfs-share/networketh"
)

type GroupAccessCAP struct {
	GroupName string
	Boxer     crypto.SymmetricKey
}

func (cap *GroupAccessCAP) Store(storage *Storage) error {
	bytesJSON, err := json.Marshal(cap)
	if err != nil {
		return fmt.Errorf("could not marshal group access capability: GroupAccessCAP.Store: %s", err)
	}
	capPath := storage.capsGAPath + "/" + cap.GroupName + ".json"
	if err := WriteFile(capPath, bytesJSON); err != nil {
		return fmt.Errorf("could not write group cap file: GroupAccessCapability.Store: %s", err)
	}
	return nil
}

type ReadCAP struct {
	FileName  string            `json:"name"`
	IPNSPath  string            `json:"ipns_path"`
	Owner     common.Address            `json:"owner"`
	VerifyKey crypto.VerifyKey `json:"verify_key"`
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

func (cap *ReadCAP) Refresh(username string, boxer *crypto.BoxingKeyPair, storage *Storage, network *nw.Network, ipfs *ipfs.IPFS) (bool, error) {
	return false, fmt.Errorf("not implemented: ReadCAP.Refresh")
}

func DownloadReadCAP(fromUserID, userID common.Address, capName string, boxer *crypto.AnonymBoxer, storage *Storage, network *nw.Network, ipfs *ipfs.IPFS) (*ReadCAP, error) {
	capBytes, err := downloadCAP(fromUserID, userID, capName, boxer, storage, network, ipfs)
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
func downloadCAP(fromUserID, userID common.Address, capName string, boxer *crypto.AnonymBoxer, storage *Storage, network *nw.Network, ipfs *ipfs.IPFS) ([]byte, error) {
	// glog.Info("Downloading CAP...")
	// // get address and key
	// _, _, _, ipfsAddr, err := network.GetUser(fromUserID)
	// if err != nil {
	// 	return nil, fmt.Errorf("could not get data of user '%s': downloadCAP: %s", base64.StdEncoding.EncodeToString(fromUserID[:]), err)
	// }

	// ipnsPath := "/ipns/" + ipfsAddr + "/for/" + base64.StdEncoding.EncodeToString(userID[:]) + "/" + capName
	// // download cap file
	// tmpFilePath := storage.tmpPath + "/" + capName
	// err = ipfs.Get(tmpFilePath, ipnsPath)
	// if err != nil {
	// 	return nil, fmt.Errorf("could not ipfs get file '%s': downloadCAP: %s", ipnsPath, err)
	// }
	// bytesEnc, err := ioutil.ReadFile(tmpFilePath)
	// if err != nil {
	// 	return nil, fmt.Errorf("could not read file '%s': downloadCAP: %s", tmpFilePath, err)
	// }
	// bytesDecr, err := boxer.Open(bytesEnc)
	// if err != nil {
	// 	return nil, fmt.Errorf(
	// 		"could not decrypt cap '%s' from user '%s' with path '%s': downloadCAP",
	// 		capName,
	// 		base64.StdEncoding.EncodeToString(fromUserID[:]),
	// 		ipnsPath,
	// 	)
	// }
	// os.Remove(tmpFilePath)
	// glog.Info("\t<-- CAP Downloaded")
	// return bytesDecr, nil
	return nil, fmt.Errorf("not implemented: downloadCAP")
}

func CreateFileReadCAPForUser(f *FilePTP, userID common.Address, ipnsAddr string, storage *Storage, network *nw.Network) error {
	cap := ReadCAP{path.Base(f.Name), ipnsAddr, f.Owner, f.VerifyKey}
	capBytes, err := json.Marshal(cap)
	if err != nil {
		return fmt.Errorf("could not marshal cap '%s': CreateFileReadCAPForUser: '%s'", cap.FileName, err)
	}
	if err := storage.createFileForUser(userID, path.Base(f.Name), capBytes, network); err != nil {
		return fmt.Errorf(
			"could not create file '%s' for user '%s': CreateFileReadCAPForUser: %s",
			cap.FileName,
			base64.StdEncoding.EncodeToString(userID[:]),
			err,
		)
	}
	return nil
}
