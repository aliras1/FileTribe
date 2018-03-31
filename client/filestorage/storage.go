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

	"ipfs-share/crypto"
	"ipfs-share/ipfs"
	nw "ipfs-share/network"
)

type Storage struct {
	dataPath        string
	publicPath      string
	publicFilesPath string
	publicForPath   string
	userDataPath    string
	capsPath        string
	storagePath     string
	sharedPath      string
	tmpPath         string
}

func NewStorage(dataPath string) *Storage {
	var us Storage
	us.dataPath = "./" + path.Clean(dataPath+"/data/")
	us.publicPath = us.dataPath + "/public"
	us.publicFilesPath = us.dataPath + "/public/files"
	us.publicForPath = us.dataPath + "/public/for"
	us.userDataPath = us.dataPath + "/userdata"
	us.capsPath = us.dataPath + "/userdata/caps"
	us.storagePath = us.dataPath + "/userdata/root"
	us.sharedPath = us.dataPath + "/userdata/shared"
	us.tmpPath = us.dataPath + "/userdata/tmp"

	os.Mkdir(us.dataPath, 0770)
	os.MkdirAll(us.publicFilesPath, 0770)
	os.MkdirAll(us.publicForPath, 0770)
	os.MkdirAll(us.capsPath, 0770)
	os.MkdirAll(us.storagePath, 0770)
	os.MkdirAll(us.sharedPath, 0770)
	os.MkdirAll(us.tmpPath, 0770)

	return &us
}

func (us *Storage) ExistsFile(filePath string) bool {
	if _, err := os.Stat(filePath); err == nil {
		return true
	}
	return false
}

// It builds up the file repo based on saved data. One part of files
// comes from capabilities which can be found in data/userdata/caps.
// These files contain information about files that are shared with the
// current user. The function appends the representation of those shared
// files into the file structure and checks if they have been updated since
// last time or not. The other half of files comes from data/userdata/shared.
// These files are JSON representation of a File that were shared by the
// user.
func (us *Storage) BuildRepo(ipfs *ipfs.IPFS) ([]*File, error) {
	var repo []*File
	// read capabilities from caps and try to refresh them
	entries, err := ioutil.ReadDir(us.capsPath)
	if err != nil {
		return []*File{}, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue // do not care about directories
		}
		cap, err := NewReadCAPFromFile(us.capsPath + "/" + entry.Name())
		if err != nil {
			continue // do not care about trash files
		}
		changed, err := cap.RefreshCAP(us, ipfs)
		if err != nil {
			return nil, err
		}
		file, err := NewFileFromCAP(cap, us, ipfs)
		if changed {
			fmt.Println("changed")
			us.DownloadFileFromCap(file, cap, ipfs)
		}
		repo = append(repo, file)
	}

	// read share information
	entries, err = ioutil.ReadDir(us.sharedPath)
	if err != nil {
		return []*File{}, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		file, err := NewFileFromShared(us.sharedPath + "/" + entry.Name())
		if err != nil {
			fmt.Println(err)
			continue
		}
		repo = append(repo, file)
	}
	return repo, nil
}

// Downloads and verifies the file, described in the capability.
func (us *Storage) DownloadFileFromCap(file *File, cap *ReadCAP, ipfs *ipfs.IPFS) error {
	tmpFilePath := us.tmpPath + "/" + path.Base(file.Path)
	err := ipfs.Get(tmpFilePath, file.IPNSPath)
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

func (us *Storage) CopyFileIntoPublicDir(filePath string) error {
	publicFilePath := us.publicFilesPath + "/" + path.Base(filePath)
	return CopyFile(filePath, publicFilePath)
}

// Signs the files with the Write key and then the function adds
// it to IPFS. The function returns with the future IPNS path of
// the file and with the IPFS hash of the file
func (us *Storage) SignAndAddFileToIPFS(filePath string, writeKey crypto.SecretSigningKey, ipfs *ipfs.IPFS) (string, string, error) {
	publicPath := us.publicFilesPath + "/" + path.Base(filePath)
	bytesFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", "", err
	}
	signedFile := writeKey.Sign(nil, bytesFile)
	if err := WriteFile(publicPath, signedFile); err != nil {
		return "", "", err
	}
	merkleNode, err := ipfs.AddFile(publicPath)
	ipfsID, err := ipfs.ID()
	ipnsPath := "/ipns/" + ipfsID.ID + "/files/" + path.Base(publicPath)
	if err != nil {
		return "", "", err
	}
	return ipnsPath, "/ipfs/" + merkleNode.Hash, nil
}

// Saves a File object (containing meta-data of a file) in json format
// locally
func (us *Storage) StoreFileMetaData(f *File) error {
	byteJson, err := json.Marshal(f)
	if err != nil {
		return err
	}
	return WriteFile(us.sharedPath+"/"+path.Base(f.Path)+".json", byteJson)
}

func (us *Storage) CreateCAPForUser(f *File, forUserPath, ipfsHash string, boxer *crypto.BoxingKeyPair, network *nw.Network) error {
	err := os.MkdirAll(forUserPath, 0770)
	if err != nil {
		fmt.Println(err) /* TODO check for permission errors */
	}
	readCAP := ReadCAP{path.Base(f.Path), f.IPNSPath, ipfsHash, f.Owner, f.VerifyKey}
	byteJSON, err := json.Marshal(readCAP)
	otherPK, err := network.GetUserBoxingKey(path.Base(forUserPath))
	if err != nil {
		return err
	}

	encJSON := boxer.BoxSeal(byteJSON, &otherPK)
	return ioutil.WriteFile(forUserPath+"/"+path.Base(f.Path)+".json", encJSON, 0644)
}

func (us *Storage) PublishPublicDir(ipfs *ipfs.IPFS) error {
	publicDir := us.dataPath + "/public"
	merkleNodes, err := ipfs.AddDir(publicDir)
	if err != nil {
		return err
	}
	for _, mn := range merkleNodes {
		if strings.Compare(mn.Name, "public") == 0 {
			err = ipfs.NamePublish(mn.Hash)
			if err != nil {
				return err
			}
			break
		}
	}
	return nil
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
