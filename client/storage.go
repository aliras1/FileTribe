package client

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/glog"
	ipfsapi "ipfs-share/ipfs"
	"ipfs-share/utils"

	"github.com/pkg/errors"
)

const (
	CAP_EXT string = ".cap"
)

// ContextData is a struct of data from UserContext
// that are stored on disk.
type ContextData struct {
	Groups         []*GroupContext
}

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
	pendingPath     string
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
	storage.pendingPath = storage.dataPath + "userdata/pending_changes/"
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
	os.MkdirAll(storage.pendingPath, 0770)
	os.MkdirAll(storage.capsGAPath, 0770)
	os.MkdirAll(storage.fileRootPath, 0770)
	os.MkdirAll(storage.myFilesPath, 0770)
	os.MkdirAll(storage.sharedPath, 0770)
	os.MkdirAll(storage.tmpPath, 0770)
	os.MkdirAll(storage.ipfsFilesPath, 0770)
	os.MkdirAll(storage.contextDataPath, 0770)

	return &storage
}

func (storage *Storage) GetUserFilesPath() string {
	return storage.fileRootPath
}

func loadFilesOfAddress(address ethcommon.Address, baseDir string) (map[[32]byte]*File, error) {
	fileMap := make(map[[32]byte]*File)

	currentDir := baseDir + "/" + address.String()
	files, err := ioutil.ReadDir(currentDir)
	if err != nil {
		return fileMap, fmt.Errorf("could not read dir: '%s': Storage.BuildRepo: %s", currentDir, err)
	}

	for _, filePTPFile := range files {
		filePTPPath := currentDir + "/" + filePTPFile.Name()

		filePTP, err := LoadPTPFile(filePTPPath)
		if err != nil {
			glog.Warningf("could not load file ptp: '%s': Storage.BuildRepo: %s", filePTPPath, err)
			continue
		}

		fileMap[filePTP.Cap.Id] = filePTP
	}

	return fileMap, nil
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


// +------------------------------+
// |   Group specific functions   |
// +------------------------------+

// Gets all the locally stored group access capabilities from
// directory data/userdata/caps/GA/
func (storage *Storage) GetGroupCaps() ([]GroupAccessCap, error) {
	var caps []GroupAccessCap
	// read capabilities from caps and try to refresh them
	groupCapFiles, err := ioutil.ReadDir(storage.capsGAPath)
	if err != nil {
		return caps, err
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
		var cap GroupAccessCap
		if err := json.Unmarshal(capBytes, &cap); err != nil {
			glog.Warning("could not unmarshal group cap: Storage.GetGroupCaps: %s", err)
			continue
		}
		cap.Boxer.RNG = rand.Reader
		caps = append(caps, cap)
	}
	return caps, nil
}

func (storage *Storage) SaveGroupCap(groupId string, data []byte) error {
	capPath := storage.capsGAPath + "/" + groupId + CAP_EXT
	if err := utils.WriteFile(capPath, data); err != nil {
		return errors.Wrap(err, "could not write group cap file")
	}
	return nil
}

func (storage *Storage) GetGroupFileCapDir(id string) string {
	return storage.capsPath + id + "/"
}

func (storage *Storage) GetGroupFilePendingDir(id string) string {
	return storage.pendingPath + id + "/"
}

func (storage *Storage) GetGroupFileDataDir(id string) string {
	return storage.fileRootPath + id + "/"
}

// Creates the directory structure needed by a group
func (storage *Storage) MakeGroupDir(id string) {
	os.MkdirAll(storage.capsPath + id, 0770)
	os.MkdirAll(storage.fileRootPath + id, 0770)
	os.MkdirAll(storage.pendingPath + id, 0770)
}


// +------------------------------+
// |       Helper functions       |
// +------------------------------+

func (storage *Storage) PublishPublicDir(ipfs ipfsapi.IIpfs) error {
	glog.Info("Publishing...")
	t := time.Now()
	publicDir := storage.dataPath + "/public"
	hash, err := ipfs.AddDir(publicDir)
	if err != nil {
		return fmt.Errorf("could not ipfs add dir: Storage.PublishPublicDir: %s", err)
	}
	glog.Info("ipfs add: ", time.Since(t))

	if err := ipfs.Publish("", hash); err != nil {
		return fmt.Errorf("could not ipfs name publish: Storage.PublishPublicDir: %s", err)
	}

	glog.Info("ipfs add n pub: ", time.Since(t))
	glog.Info("Publishing ended")
	return nil
}

func (storage *Storage) MakeForDirectory(dirName string, ipfs ipfsapi.IIpfs) (string, error) {
	dirPath := storage.publicForPath + "/" + dirName
	os.Mkdir(dirPath, 0770)
	hash, err := ipfs.AddDir(dirPath)
	if err != nil {
		return "", fmt.Errorf("could not publish: Storage.MakeForDirectory: %s", err)
	}
	return hash, nil
}

func (storage *Storage) InitMyCapsByFriend(friendAddress string) {
	dirPath := storage.capsPath + "/" + friendAddress
	os.Mkdir(dirPath, 0770)
}

func (storage *Storage) DownloadTmpFile(ipfsHash string, ipfs ipfsapi.IIpfs) (string, error) {
	filePath := storage.tmpPath + "/" + ipfsHash
	if err := ipfs.Get(ipfsHash, filePath); err != nil {
		return "", fmt.Errorf("could not ipfs get into tmp path '%s': Storage.downloadToTmp: %s", filePath, err)
	}
	return filePath, nil
}

