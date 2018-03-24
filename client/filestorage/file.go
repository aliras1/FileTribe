package filestorage

import (
	"encoding/json"
	"io/ioutil"
	"path"

	"ipfs-share/crypto"
	"ipfs-share/ipfs"
	nw "ipfs-share/network"
)

type File struct {
	Path       string   `json:"path"`
	IPNSPath   string   `json:"ipnsPath"`
	Owner      string   `json:"owner"`
	SharedWith []string `json:"shared_with"`
	WAccess    []string `json:"w_access"`
}

func NewFileFromShared(filePath string) (*File, error) {
	bytesFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var file File
	err = json.Unmarshal(bytesFile, &file)
	return &file, err
}

func (f *File) Share(shareWith []string, baseDirPath, ipfsHash string, boxer *crypto.BoxingKeyPair, network *nw.Network, ipfs *ipfs.IPFS, us *UserStorage) error {
	var newUsers []string
	for _, user := range shareWith {
		// add to share list
		f.SharedWith = append(f.SharedWith, user)
		// make new capability into for_X directory
		err := us.CreateCapabilityFile(f, baseDirPath+user, "/ipfs/"+ipfsHash, boxer)
		newUsers = append(newUsers, user)
		if err != nil {
			return err
		}
	}
	us.CreateFileIntoPublicDir(f.Path)
	us.StoreFileMetaData(f)
	err := us.PublishPublicDir()
	if err != nil {
		return err
	}
	// send share messages
	for _, user := range newUsers {
		err = network.SendMessage(f.Owner, user, path.Base(f.Path))
		if err != nil {
			return err
		}
	}
	return nil
}
