package filestorage

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"ipfs-share/ipfs"
	nw "ipfs-share/network"
)

type UserStorage struct {
	RootDir  []*File
	Username string
	DataPath string
	IPFS     *ipfs.IPFS
	Network  *nw.Network
}

func (u *UserStorage) Init(path string, i *ipfs.IPFS) error {
	u.IPFS = i
	u.DataPath = path + "/data/"
	var errorArray []error

	err := os.Mkdir(u.DataPath, 0770)
	if err != nil {
		errorArray = append(errorArray, err)
	}

	err = os.Mkdir(u.DataPath+"/public", 0770)
	if err != nil {
		errorArray = append(errorArray, err)
	}

	err = os.MkdirAll(u.DataPath+"/userdata/root", 0770)
	if err != nil {
		errorArray = append(errorArray, err)
	}

	f, err := os.Create(u.DataPath + "/userdata/caps.json")
	f.Close()
	if err != nil {
		errorArray = append(errorArray, err)
	}

	f, err = os.Create(u.DataPath + "/userdata/shared.json")
	f.Close()
	if err != nil {
		errorArray = append(errorArray, err)
	}

	u.build()
	return fmt.Errorf("%s", errorArray)
}

func (u *UserStorage) build() error {
	return errors.New("not implemented")
}

func (u *UserStorage) IsFileInRootDir(filePath string) bool {
	for _, i := range u.RootDir {
		if strings.Compare(i.Path, filePath) == 0 {
			return true
		}
	}
	return false
}

func (u *UserStorage) AddAndShareFile(filePath string, shareWith []string) error {
	if u.IsFileInRootDir(filePath) {
		return errors.New("file already in root dir")
	}
	merkleNode, err := u.IPFS.AddFile(filePath)
	if err != nil {
		return err
	}
	f := File{filePath, merkleNode.Hash, u.Username, []string{}, []string{}}
	err = f.Share(shareWith, u.DataPath+"/public/for/", u.Network)
	if err != nil {
		return err
	}
	u.RootDir = append(u.RootDir, &f)
	return nil
}
