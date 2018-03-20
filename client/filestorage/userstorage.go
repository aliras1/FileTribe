package filestorage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"ipfs-share/ipfs"
	nw "ipfs-share/network"
)

type ReadCAP struct {
	Name  string `json:"name"`
	Hash  string `json:"hash"`
	Owner string `json:"owner"`
}

type UserStorage struct {
	RootDir  []*File
	DataPath string
	Username string
	IPFS     *ipfs.IPFS
	Network  *nw.Network
}

func NewUserStorage(dataPath, username string, ipfs *ipfs.IPFS, network *nw.Network) *UserStorage {
	var us UserStorage
	us.Username = username
	us.IPFS = ipfs
	us.Network = network
	us.DataPath = "./" + path.Clean(dataPath+"/data/")
	us.RootDir = []*File{}

	os.Mkdir(us.DataPath, 0770)
	os.MkdirAll(us.DataPath+"/tmp", 0770)
	os.MkdirAll(us.DataPath+"/public", 0770)
	os.MkdirAll(us.DataPath+"/userdata/root", 0770)
	f, _ := os.Create(us.DataPath + "/userdata/caps.json")
	f.Close()
	f, _ = os.Create(us.DataPath + "/userdata/shared.json")
	f.Close()
	us.build()
	return &us
}

func (us *UserStorage) build() error {
	return errors.New("not implemented")
}

func (us *UserStorage) List() {
	for _, f := range us.RootDir {
		fmt.Print("\t--> ")
		fmt.Println(*f)
	}
}

func (us *UserStorage) IsFileInRootDir(filePath string) bool {
	for _, i := range us.RootDir {
		if strings.Compare(i.Path, filePath) == 0 {
			return true
		}
	}
	return false
}

func (us *UserStorage) AddAndShareFile(filePath, owner string, shareWith []string) error {
	if us.IsFileInRootDir(filePath) {
		return errors.New("file already in root dir")
	}
	merkleNode, err := us.IPFS.AddFile(filePath)
	if err != nil {
		return err
	}
	f := File{filePath, merkleNode.Hash, owner, []string{}, []string{}}
	err = f.Share(shareWith, us.DataPath+"/public/for/", us.Network, us.IPFS, us)
	if err != nil {
		return err
	}
	us.RootDir = append(us.RootDir, &f)
	return nil
}

func (us *UserStorage) CreateCapabilityFile(f *File, forPath string) error {
	err := os.MkdirAll(forPath, 0770)
	if err != nil {
		fmt.Println(err) /* TODO check for permission errors */
	}
	readCAP := ReadCAP{path.Base(f.Path), path.Base(f.Hash), us.Username}
	byteJson, err := json.Marshal(readCAP)
	return ioutil.WriteFile(forPath+"/"+path.Base(f.Path)+".json", byteJson, 0644)
}

func (us *UserStorage) PublishPublicDir() error {
	publicDir := us.DataPath + "/public"
	merkleNodes, err := us.IPFS.AddDir(publicDir)
	if err != nil {
		return err
	}
	for _, mn := range merkleNodes {
		if strings.Compare(mn.Name, "public") == 0 {
			err = us.IPFS.NamePublish(mn.Hash)
			if err != nil {
				return err
			}
			break
		}
	}
	return nil
}

func (us *UserStorage) Contains(filePath string) bool {
	for _, f := range us.RootDir {
		if strings.Compare(path.Base(f.Path), filePath) == 0 {
			return true
		}
	}
	return false
}

func (us *UserStorage) AddFileFromIPFS(name, hash string) error {
	// download tmp cap file
	tmpFilePath := us.DataPath + "/tmp/" + name
	err := us.IPFS.Get(tmpFilePath, hash)
	if err != nil {
		return err
	}

	bytesCapability, err := ioutil.ReadFile(tmpFilePath)
	if err != nil {
		return err
	}
	// remove tmp cap file
	err = os.Remove(tmpFilePath)
	if err != nil {
		return err
	}
	var readCAP ReadCAP
	err = json.Unmarshal(bytesCapability, &readCAP)
	if err != nil {
		return err
	}
	// download and add to root directory
	filePath := us.DataPath + "/userdata/root/" + readCAP.Name
	f := File{path.Clean(filePath), readCAP.Owner, readCAP.Owner, []string{}, []string{}}
	err = us.IPFS.Get(filePath, readCAP.Hash)
	if err != nil {
		return err
	}
	us.RootDir = append(us.RootDir, &f)
	return nil
}
