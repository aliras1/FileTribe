package filestorage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"ipfs-share/ipfs"
	nw "ipfs-share/network"
)

type ReadCAP struct {
	Name     string `json:"name"`
	IPNSPath string `json:"ipns_path"`
	IPFSHash string `json:"ipfs_hash"`
	Owner    string `json:"owner"`
}

func (rc *ReadCAP) Store(capPath string) error {
	bytesJSON, err := json.Marshal(rc)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(capPath, bytesJSON, 0644)
}

func NewReadCAPFromFile(capPath string) (*ReadCAP, error) {
	bytesFile, err := ioutil.ReadFile(capPath)
	if err != nil {
		return nil, err
	}
	var cap ReadCAP
	err = json.Unmarshal(bytesFile, &cap)
	return &cap, err
}

type UserStorage struct {
	RootDir  []*File
	DataPath string
	Username string
	IPFS     *ipfs.IPFS
	Network  *nw.Network
}

func NewUserStorage(dataPath, username string, network *nw.Network, ipfs *ipfs.IPFS) *UserStorage {
	var us UserStorage
	us.Username = username
	us.IPFS = ipfs
	us.Network = network
	us.DataPath = "./" + path.Clean(dataPath+"/data/")
	us.RootDir = []*File{}

	os.Mkdir(us.DataPath, 0770)
	os.MkdirAll(us.DataPath+"/public/files", 0770)
	os.MkdirAll(us.DataPath+"/userdata/root", 0770)
	os.MkdirAll(us.DataPath+"/userdata/caps", 0770)
	os.MkdirAll(us.DataPath+"/userdata/shared", 0770)
	err := us.build()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &us
}

// It builds up the file structure based on saved data. One part of files
// comes from capabilities which can be found in data/userdata/caps.
// These files contain information about files that are shared with the
// current user. The function appends the representation of those shared
// files into the file structure and checks if they have been updated since
// last time or not. The other half of files comes from data/userdata/shared.
// These files are JSON representation of a File that were shared by the
// user.
func (us *UserStorage) build() error {
	// read capabilities from caps and try to refresh them
	capsPath := us.DataPath + "/userdata/caps"
	entries, err := ioutil.ReadDir(capsPath)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue // do not care about directories
		}
		cap, err := NewReadCAPFromFile(capsPath + "/" + entry.Name())
		if err != nil {
			continue // do not care about trash files
		}
		changed, err := us.RefreshCAP(cap)
		if err != nil {
			return err
		}
		file := us.addFileToRootFromCap(cap)
		if changed {
			fmt.Println("changed")
			us.downloadFileFromCap(file, cap)
		}
	}

	// read share information
	sharedPath := us.DataPath + "/userdata/shared"
	entries, err = ioutil.ReadDir(sharedPath)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		file, err := NewFileFromShared(sharedPath + "/" + entry.Name())
		if err != nil {
			fmt.Println(err)
			continue
		}
		us.addFileToRootDir(file)
	}
	return nil
}

// Checks if the by ReadCap represented file has changed since last time
// or not. It is done via checking the IPFS hash of the file. If it has
// the function returns true. If the file is not present the function
// return true as well. Otherwise it returns false.
func (us *UserStorage) RefreshCAP(cap *ReadCAP) (bool, error) {
	fileExists := false
	if _, err := os.Stat(us.DataPath + "/userdata/root/" + cap.Name); err == nil {
		fileExists = true
	}
	resolvedHash, err := us.IPFS.Resolve(cap.IPNSPath)
	fileChanged := strings.Compare(resolvedHash, cap.IPFSHash) != 0
	if err != nil {
		return false, err
	}
	if fileChanged {
		cap.IPFSHash = resolvedHash
		cap.Store(us.DataPath + "/userdata/caps/" + cap.Name + ".json")
	}
	if !fileExists || fileChanged {
		return true, nil
	}
	return false, nil
}

func (us *UserStorage) addFileToRootDir(file *File) {
	us.RootDir = append(us.RootDir, file)
}

func (us *UserStorage) getFileFromRootDir(name string) *File {
	for _, file := range us.RootDir {
		if strings.Compare(path.Base(file.Path), name) == 0 {
			return file
		}
	}
	return nil
}

func (us *UserStorage) addFileToRootFromCap(cap *ReadCAP) *File {
	if us.IsFileInRootDir(cap.Name) {
		return us.getFileFromRootDir(cap.Name)
	}
	filePath := us.DataPath + "/userdata/root/" + cap.Name
	file := &File{path.Clean(filePath), cap.IPNSPath, cap.Owner, []string{}, []string{}}
	us.addFileToRootDir(file)
	return file
}

func (us *UserStorage) downloadFileFromCap(file *File, cap *ReadCAP) error {
	err := us.IPFS.Get(file.Path, file.IPNSPath)
	if err != nil {
		return err
	}
	return nil
}

func (us *UserStorage) IsFileInRootDir(filePath string) bool {
	for _, i := range us.RootDir {
		if strings.Compare(path.Base(i.Path), path.Base(filePath)) == 0 {
			return true
		}
	}
	return false
}

func (us *UserStorage) CreateFileIntoPublicDir(filePath string) error {
	publicFilePath := us.DataPath + "/public/files/" + path.Base(filePath)
	return CopyFile(filePath, publicFilePath)
}

func (us *UserStorage) AddAndShareFile(filePath, owner string, shareWith []string) error {
	if us.IsFileInRootDir(filePath) {
		return errors.New("file already in root dir")
	}
	merkleNode, err := us.IPFS.AddFile(filePath)
	ipfsID, err := us.IPFS.ID()
	ipnsPath := "/ipns/" + ipfsID.ID + "/files/" + path.Base(filePath)
	if err != nil {
		return err
	}
	file := File{filePath, ipnsPath, owner, []string{}, []string{}}
	err = file.Share(shareWith, us.DataPath+"/public/for/", merkleNode.Hash, us.Network, us.IPFS, us)
	if err != nil {
		return err
	}
	us.addFileToRootDir(&file)
	return nil
}

func (us *UserStorage) StoreFileMetaData(f *File) error {
	byteJson, err := json.Marshal(f)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(us.DataPath+"/userdata/shared/"+path.Base(f.Path)+".json", byteJson, 0644)
}

func (us *UserStorage) CreateCapabilityFile(f *File, forPath, ipfsHash string) error {
	err := os.MkdirAll(forPath, 0770)
	if err != nil {
		fmt.Println(err) /* TODO check for permission errors */
	}
	readCAP := ReadCAP{path.Base(f.Path), f.IPNSPath, ipfsHash, us.Username}
	fmt.Print("create cap file: ")
	fmt.Println(f.IPNSPath)
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

func (us *UserStorage) AddFileFromIPNS(capName, capIPFSHash, ipfsAddr string) error {
	// download cap file
	tmpFilePath := us.DataPath + "/userdata/caps/" + capName
	err := us.IPFS.Get(tmpFilePath, capIPFSHash)
	if err != nil {
		return err
	}

	readCAP, err := NewReadCAPFromFile(tmpFilePath)
	if err != nil {
		return err
	}
	file := us.addFileToRootFromCap(readCAP)
	if file == nil {
		return errors.New("could not add file to root dir from cap")
	}
	return us.downloadFileFromCap(file, readCAP)
}

func (us *UserStorage) List() {
	fmt.Println(us.Username)
	for _, f := range us.RootDir {
		fmt.Print("\t--> ")
		fmt.Println(*f)
	}
}

func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dst) // creates if file doesn't exist
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile) // check first var for number of bytes copied
	if err != nil {
		return err
	}
	return destFile.Sync()
}
