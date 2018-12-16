package fs

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"io/ioutil"
	"ipfs-share/client/fs/caps"
	ipfsapi "ipfs-share/ipfs"
	"ipfs-share/utils"
	"os"
	"path"

	"bytes"
	"github.com/pkg/errors"
	"ipfs-share/crypto"
)

const (
	CAP_EXT string = ".cap"
)

// Storage is a struct of the directory paths and has
// functions that are responsible for the file level
// functionalities
type Storage struct {
	dataPath        string
	publicPath      string
	publicFilesPath string
	publicForPath   string
	userDataPath    string
	capsPath        string
	origPath        string
	capsGAPath      string // group access caps
	fileRootPath    string
	sharedPath      string
	tmpPath         string
	myFilesPath     string
	ipfsFilesPath   string
	contextDataPath string
}

// NewStorage creates the directory structure
func NewStorage(dataPath string) *Storage {
	var storage Storage
	storage.dataPath = "./" + path.Clean(dataPath + "/data") + "/"
	storage.publicPath = storage.dataPath + "public/"
	storage.publicFilesPath = storage.dataPath + "public/files/"
	storage.publicForPath = storage.dataPath + "public/for/"
	storage.userDataPath = storage.dataPath + "userdata/"
	storage.capsPath = storage.dataPath + "userdata/caps/"
	storage.origPath = storage.dataPath + "userdata/orig/"
	storage.capsGAPath = storage.dataPath + "userdata/caps/GA/"
	storage.fileRootPath = storage.dataPath + "userdata/root/"
	storage.myFilesPath = storage.dataPath + "userdata/root/MyFiles/"
	storage.sharedPath = storage.dataPath + "userdata/shared/"
	storage.tmpPath = storage.dataPath + "userdata/tmp/"
	storage.contextDataPath = storage.dataPath + "userdata/context/"

	os.Mkdir(storage.dataPath, 0770)
	os.MkdirAll(storage.publicFilesPath, 0770)
	os.MkdirAll(storage.publicForPath, 0770)
	os.MkdirAll(storage.capsPath, 0770)
	os.MkdirAll(storage.origPath, 0770)
	os.MkdirAll(storage.capsGAPath, 0770)
	os.MkdirAll(storage.fileRootPath, 0770)
	os.MkdirAll(storage.myFilesPath, 0770)
	os.MkdirAll(storage.sharedPath, 0770)
	os.MkdirAll(storage.tmpPath, 0770)
	os.MkdirAll(storage.ipfsFilesPath, 0770)
	os.MkdirAll(storage.contextDataPath, 0770)

	return &storage
}

func (storage *Storage) UserFilesPath() string {
	return storage.fileRootPath
}

func (storage *Storage) CopyFileIntoPublicDir(filePath string) error {
	publicFilePath := storage.publicFilesPath + "/" + path.Base(filePath)
	return utils.CopyFile(filePath, publicFilePath)
}

func (storage *Storage) CopyFileIntoMyFiles(filePath string) (string, error) {
	newFilePath := storage.myFilesPath + "/" + path.Base(filePath)
	return newFilePath, utils.CopyFile(filePath, newFilePath)
}

func (storage *Storage) CopyFileIntoGroupFiles(filePath, groupName string) error {
	groupFilesPath := storage.fileRootPath + "/" + groupName
	os.Mkdir(groupFilesPath, 0770)
	newFilePath := groupFilesPath + "/" + path.Base(filePath)
	return utils.CopyFile(filePath, newFilePath)
}

func (storage *Storage) SaveAccountData(data []byte) error {
	path := storage.userDataPath + "account.dat"

	if err := utils.WriteFile(path, data); err != nil {
		return errors.Wrapf(err, "could not write to file: %s", path)
	}

	return nil
}

func (storage *Storage) LoadAccountData() ([]byte, error) {
	path := storage.userDataPath + "account.dat"

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "could not read file: %s", path)
	}

	return data, nil
}

// +------------------------------+
// |   Group specific functions   |
// +------------------------------+

// Gets all the locally stored group access capabilities from
// directory data/userdata/caps/GA/
func (storage *Storage) GetGroupCaps() ([]caps.GroupAccessCap, error) {
	var capabilities []caps.GroupAccessCap
	// read capabilities from caps and try to refresh them
	groupCapFiles, err := ioutil.ReadDir(storage.capsGAPath)
	if err != nil {
		return capabilities, err
	}
	for _, groupCapFile := range groupCapFiles {
		if groupCapFile.IsDir() {
			continue // do not care about directories
		}
		filePath := storage.capsGAPath + "/" + groupCapFile.Name()
		capBytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			glog.Warning("could not read file '%s': Storage.GetGroupCaps: %s", filePath, err)
			continue
		}
		var cap caps.GroupAccessCap
		if err := json.Unmarshal(capBytes, &cap); err != nil {
			glog.Warning("could not unmarshal group cap: Storage.GetGroupCaps: %s", err)
			continue
		}
		cap.Boxer.RNG = rand.Reader
		capabilities = append(capabilities, cap)
	}
	return capabilities, nil
}

func (storage *Storage) SaveGroupAccessCap(cap *caps.GroupAccessCap) error {
	capJson, err := json.Marshal(cap)
	if err != nil {
		return errors.Wrap(err, "could not marshal group access capability")
	}

	path := storage.GroupAccessCapDir() + cap.Address.String() + CAP_EXT
	if err := utils.WriteFile(path, capJson); err != nil {
		return errors.Wrap(err, "could not write group cap file")
	}

	return nil
}

func (storage *Storage) GroupAccessCapDir() string {
	return storage.capsGAPath
}

func (storage *Storage) GroupFileCapDir(id string) string {
	return storage.capsPath + id + "/"
}

func (storage *Storage) GroupFileOrigDir(id string) string {
	return storage.origPath + id + "/"
}

func (storage *Storage) GroupFileDataDir(id string) string {
	return storage.fileRootPath + id + "/"
}

// Creates the directory structure needed by a group
func (storage *Storage) MakeGroupDir(id string) {
	os.MkdirAll(storage.capsPath + id, 0770)
	os.MkdirAll(storage.fileRootPath + id, 0770)
	os.MkdirAll(storage.origPath+ id, 0770)
}

func (storage *Storage) DownloadTmpFile(ipfsHash string, ipfs ipfsapi.IIpfs) (string, error) {
	filePath := storage.tmpPath + "/" + ipfsHash
	if err := ipfs.Get(ipfsHash, filePath); err != nil {
		return "", fmt.Errorf("could not ipfs get into tmp path '%s': Storage.downloadToTmp: %s", filePath, err)
	}
	return filePath, nil
}

func (storage *Storage) DownloadAndDecryptWithSymmetricKey(boxer crypto.SymmetricKey, ipfsHash string, ipfs ipfsapi.IIpfs) ([]byte, error) {
	path := storage.tmpPath + ipfsHash
	if err := ipfs.Get(ipfsHash, path); err != nil {
		return nil, errors.Wrapf(err, "could not ipfs get ipfs hash %s", ipfsHash)
	}

	encData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "could not read tmp file: %s", path)
	}
	defer func() {
		if err := os.Remove(path); err != nil {
			glog.Warningf("could not remove tmp file %s", path)
		}
	}()

	data, ok := boxer.BoxOpen(encData)
	if !ok {
		return nil, errors.New("could not decrypt shared group dir")
	}

	return data, nil
}

func (storage *Storage) DownloadAndDecryptWithFileBoxer(boxer crypto.FileBoxer, ipfsHash string, ipfs ipfsapi.IIpfs) ([]byte, error) {
	tmpFilePath, err := storage.DownloadTmpFile(ipfsHash, ipfs)
	if err != nil {
		return nil, errors.Wrapf(err, "could not ipfs get '%s'", ipfsHash)
	}

	encReader, err := os.Open(tmpFilePath)
	if err != nil {
		return nil, errors.Wrapf(err, "download err: could not read file '%s'", tmpFilePath)
	}

	diffBuf := new(bytes.Buffer)
	err = boxer.Open(encReader, diffBuf)
	defer func() {
		if err := encReader.Close(); err != nil {
			glog.Warningf("download err: could not close tmp file '%s': %s", tmpFilePath, err)
		}
		if err := os.Remove(tmpFilePath); err != nil {
			glog.Warningf("download err: could not delete tmp file '%s': %s", tmpFilePath, err)
		}
	}()

	if err != nil {
		return nil, errors.Wrap(err, "download err: could not decrypt file dif")
	}

	return diffBuf.Bytes(), nil
}
