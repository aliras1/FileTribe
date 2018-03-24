package filestorage

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"golang.org/x/crypto/ed25519"

	"ipfs-share/crypto"
	"ipfs-share/ipfs"
	nw "ipfs-share/network"
)

type ReadCAP struct {
	Name      string                  `json:"name"`
	IPNSPath  string                  `json:"ipns_path"`
	IPFSHash  string                  `json:"ipfs_hash"`
	Owner     string                  `json:"owner"`
	VerifyKey crypto.PublicSigningKey `json:"verify_key"`
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
	Storage  []*File
	Username string
	IPFS     *ipfs.IPFS
	Network  *nw.Network

	DataPath        string
	PublicPath      string
	PublicFilesPath string
	PublicForPath   string
	UserDataPath    string
	CapsPath        string
	StoragePath     string
	SharedPath      string
	TmpPath         string
}

func NewUserStorage(dataPath, username string, network *nw.Network, ipfs *ipfs.IPFS) *UserStorage {
	var us UserStorage
	us.Username = username
	us.IPFS = ipfs
	us.Network = network
	us.Storage = []*File{}
	us.DataPath = "./" + path.Clean(dataPath+"/data/")
	us.PublicPath = us.DataPath + "/public"
	us.PublicFilesPath = us.DataPath + "/public/files"
	us.PublicForPath = us.DataPath + "/public/for"
	us.UserDataPath = us.DataPath + "/userdata"
	us.CapsPath = us.DataPath + "/userdata/caps"
	us.StoragePath = us.DataPath + "/userdata/root"
	us.SharedPath = us.DataPath + "/userdata/shared"
	us.TmpPath = us.DataPath + "/userdata/tmp"

	os.Mkdir(us.DataPath, 0770)
	os.MkdirAll(us.PublicFilesPath, 0770)
	os.MkdirAll(us.PublicForPath, 0770)
	os.MkdirAll(us.CapsPath, 0770)
	os.MkdirAll(us.StoragePath, 0770)
	os.MkdirAll(us.SharedPath, 0770)
	os.MkdirAll(us.TmpPath, 0770)
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
	entries, err := ioutil.ReadDir(us.CapsPath)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue // do not care about directories
		}
		cap, err := NewReadCAPFromFile(us.CapsPath + "/" + entry.Name())
		if err != nil {
			continue // do not care about trash files
		}
		changed, err := us.RefreshCAP(cap)
		if err != nil {
			return err
		}
		file := us.addFileToStorageFromCap(cap)
		if changed {
			fmt.Println("changed")
			us.downloadFileFromCap(file, cap)
		}
	}

	// read share information
	entries, err = ioutil.ReadDir(us.SharedPath)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		file, err := NewFileFromShared(us.SharedPath + "/" + entry.Name())
		if err != nil {
			fmt.Println(err)
			continue
		}
		us.addFileToStorage(file)
	}
	return nil
}

// Checks if the by ReadCap represented file has changed since last time
// or not. It is done via checking the IPFS hash of the file. If it has
// the function returns true. If the file is not present the function
// return true as well. Otherwise it returns false.
func (us *UserStorage) RefreshCAP(cap *ReadCAP) (bool, error) {
	// TODO this could be a CAP function
	fileExists := false
	if _, err := os.Stat(us.StoragePath + "/" + cap.Name); err == nil {
		fileExists = true
	}
	resolvedHash, err := us.IPFS.Resolve(cap.IPNSPath)
	fileChanged := strings.Compare(resolvedHash, cap.IPFSHash) != 0
	if err != nil {
		return false, err
	}
	if fileChanged {
		cap.IPFSHash = resolvedHash
		cap.Store(us.CapsPath + "/" + cap.Name + ".json")
	}
	if !fileExists || fileChanged {
		return true, nil
	}
	return false, nil
}

func (us *UserStorage) addFileToStorage(file *File) {
	us.Storage = append(us.Storage, file)
}

func (us *UserStorage) getFileFromStorage(name string) *File {
	for _, file := range us.Storage {
		if strings.Compare(path.Base(file.Path), name) == 0 {
			return file
		}
	}
	return nil
}

func (us *UserStorage) addFileToStorageFromCap(cap *ReadCAP) *File {
	if us.IsFileInStorage(cap.Name) {
		return us.getFileFromStorage(cap.Name)
	}
	filePath := us.StoragePath + "/" + cap.Name
	file := &File{path.Clean(filePath), cap.IPNSPath, cap.Owner, []string{}, []string{}, cap.VerifyKey, crypto.SecretSigningKey{}}
	us.addFileToStorage(file)
	return file
}

func (us *UserStorage) downloadFileFromCap(file *File, cap *ReadCAP) error {
	tmpFilePath := us.TmpPath + "/" + path.Base(file.Path)
	err := us.IPFS.Get(tmpFilePath, file.IPNSPath)
	if err != nil {
		return err
	}
	bytesSignedFile, err := ioutil.ReadFile(tmpFilePath)
	if err != nil {
		return err
	}
	os.Remove(tmpFilePath)
	bytesRawFile, ok := file.VerifyKey.Open(nil, bytesSignedFile)
	if !ok {
		return errors.New("by downloadFileFromCap(): could not verify file")
	}
	if err := WriteFile(file.Path, bytesRawFile); err != nil {
		return err
	}
	return nil
}

func (us *UserStorage) IsFileInStorage(filePath string) bool {
	for _, i := range us.Storage {
		if strings.Compare(path.Base(i.Path), path.Base(filePath)) == 0 {
			return true
		}
	}
	return false
}

func (us *UserStorage) CopyFileIntoPublicDir(filePath string) error {
	publicFilePath := us.PublicFilesPath + "/" + path.Base(filePath)
	return CopyFile(filePath, publicFilePath)
}

func (us *UserStorage) AddAndShareFile(filePath, owner string, shareWith []string, boxer *crypto.BoxingKeyPair) error {
	if us.IsFileInStorage(filePath) {
		return errors.New("file already in root dir")
	}
	vk, wk, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}
	verifyKey := crypto.PublicSigningKey(vk)
	writeKey := crypto.SecretSigningKey(wk)

	publicPath := us.PublicFilesPath + "/" + path.Base(filePath)
	err = us.SignAndCopyFileToPublicDir(filePath, publicPath, writeKey)
	if err != nil {
		return err
	}

	merkleNode, err := us.IPFS.AddFile(publicPath)
	ipfsID, err := us.IPFS.ID()
	ipnsPath := "/ipns/" + ipfsID.ID + "/files/" + path.Base(publicPath)
	if err != nil {
		return err
	}
	file := File{filePath, ipnsPath, owner, []string{}, []string{}, verifyKey, writeKey}
	err = file.Share(shareWith, us.PublicForPath+"/", merkleNode.Hash, boxer, us.Network, us.IPFS, us)
	if err != nil {
		return err
	}
	us.addFileToStorage(&file)
	return nil
}

func (us *UserStorage) SignAndCopyFileToPublicDir(filePath, publicPath string, writeKey crypto.SecretSigningKey) error {
	bytesFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	signedFile := writeKey.Sign(nil, bytesFile)
	if err := WriteFile(publicPath, signedFile); err != nil {
		return err
	}
	return nil
}

func (us *UserStorage) StoreFileMetaData(f *File) error {
	byteJson, err := json.Marshal(f)
	if err != nil {
		return err
	}
	return WriteFile(us.SharedPath+"/"+path.Base(f.Path)+".json", byteJson)
}

func (us *UserStorage) CreateCapabilityFile(f *File, forPath, ipfsHash string, boxer *crypto.BoxingKeyPair) error {
	err := os.MkdirAll(forPath, 0770)
	if err != nil {
		fmt.Println(err) /* TODO check for permission errors */
	}
	readCAP := ReadCAP{path.Base(f.Path), f.IPNSPath, ipfsHash, us.Username, f.VerifyKey}
	byteJSON, err := json.Marshal(readCAP)
	otherPK, err := us.Network.GetUserBoxingKey(path.Base(forPath))
	if err != nil {
		return err
	}

	encJSON := boxer.BoxSeal(byteJSON, &otherPK)
	return ioutil.WriteFile(forPath+"/"+path.Base(f.Path)+".json", encJSON, 0644)
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

func (us *UserStorage) AddFileFromIPNS(capName, capIPFSHash, ipfsAddr string, otherPK *crypto.PublicBoxingKey, boxer *crypto.BoxingKeyPair) error {
	// download cap file
	tmpFilePath := us.TmpPath + "/" + capName
	err := us.IPFS.Get(tmpFilePath, capIPFSHash)
	if err != nil {
		return err
	}
	capFilePath := us.CapsPath + "/" + capName
	bytesEnc, err := ioutil.ReadFile(tmpFilePath)
	if err != nil {
		return err
	}

	bytesDecr, success := boxer.BoxOpen(bytesEnc, otherPK)
	if !success {
		return errors.New("could not decrypt capability")
	}
	os.Remove(tmpFilePath)
	if err := WriteFile(capFilePath, bytesDecr); err != nil {
		return err
	}

	readCAP, err := NewReadCAPFromFile(capFilePath)
	if err != nil {
		return errors.New("error by NewReadCAPFromFile: " + err.Error())
	}
	file := us.addFileToStorageFromCap(readCAP)
	if file == nil {
		return errors.New("could not add file to root dir from cap")
	}
	return us.downloadFileFromCap(file, readCAP)
}

func (us *UserStorage) List() {
	fmt.Println(us.Username)
	for _, f := range us.Storage {
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

func WriteFile(filePath string, data []byte) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		return err
	}
	f.Sync()
	return nil
}
