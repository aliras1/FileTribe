package filestorage

import (
	"errors"
	"fmt"
	"ipfs-share/client"
	"ipfs-share/ipfs"
	"os"
	"strings"
)

type UserStorage struct {
	RootDir  []*File
	username string
	dataPath string
	ipfs     *ipfs.IPFS
	network  *client.Network
}

func (u *UserStorage) Init(path string, i *ipfs.IPFS) error {
	u.ipfs = i
	u.dataPath = path + "/data/"
	var errorArray []error

	err := os.Mkdir(u.dataPath, 0770)
	if err != nil {
		errorArray = append(errorArray, err)
	}

	err = os.Mkdir(u.dataPath+"/public", 0770)
	if err != nil {
		errorArray = append(errorArray, err)
	}

	err = os.MkdirAll(u.dataPath+"/userdata/root", 0770)
	if err != nil {
		errorArray = append(errorArray, err)
	}

	f, err := os.Create(u.dataPath + "/userdata/caps.json")
	f.Close()
	if err != nil {
		errorArray = append(errorArray, err)
	}

	f, err = os.Create(u.dataPath + "/userdata/shared.json")
	f.Close()
	if err != nil {
		errorArray = append(errorArray, err)
	}

	u.build()
	return fmt.Errorf("%s", errorArray)
}

func (u *UserStorage) build() error {

	return nil
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
	merkleNode, err := u.ipfs.AddFile(filePath)
	if err != nil {
		return err
	}
	f := File{filePath, merkleNode.Hash, u.username, []string{}, []string{}}
	err = f.Share(shareWith, u.dataPath+"/public/for/", u.network)
	if err != nil {
		return err
	}
	u.RootDir = append(u.RootDir, &f)
	return nil
}
