package filestorage

import (
	"path"

	"ipfs-share/ipfs"
	nw "ipfs-share/network"
)

type EntryProvider interface {
}

type File struct {
	Path       string   `json:"path"`
	Hash       string   `json:"ipfs_hash"`
	Owner      string   `json:"owner"`
	SharedWith []string `json:"shared_with"`
	WAccess    []string `json:"w_access"`
}

func (f *File) Share(shareWith []string, baseDirPath string, network *nw.Network, ipfs *ipfs.IPFS, us *UserStorage) error {
	for _, user := range shareWith {
		// add to share list
		f.SharedWith = append(f.SharedWith, user)
		// make new capability into for_X directory
		err := us.CreateCapabilityFile(f, baseDirPath+user)
		if err != nil {
			return err
		}
	}
	// re-publish the public directory
	err := us.PublishPublicDir()
	if err != nil {
		return err
	}
	for _, user := range shareWith {
		// send share message
		err = network.SendMessage(f.Owner, user, path.Base(f.Path))
		if err != nil {
			return err
		}
	}
	return nil
}
