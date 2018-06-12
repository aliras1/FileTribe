package filestorage

import (
	"github.com/ethereum/go-ethereum/common"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/golang/glog"

	"ipfs-share/crypto"
	"ipfs-share/ipfs"
	nw "ipfs-share/networketh"
)

type Storage struct {
	dataPath        string
	publicPath      string
	publicFilesPath string
	publicForPath   string
	userDataPath    string
	capsPath        string
	capsGAPath      string // group access caps
	fileRootPath    string
	sharedPath      string
	tmpPath         string
	myFilesPath     string
	ipfsFilesPath   string
}

func NewStorage(dataPath string) *Storage {
	var storage Storage
	storage.dataPath = "./" + path.Clean(dataPath+"/data/")
	storage.publicPath = storage.dataPath + "/public"
	storage.publicFilesPath = storage.dataPath + "/public/files"
	storage.publicForPath = storage.dataPath + "/public/for"
	storage.userDataPath = storage.dataPath + "/userdata"
	storage.capsPath = storage.dataPath + "/userdata/caps"
	storage.capsGAPath = storage.dataPath + "/userdata/caps/GA"
	storage.fileRootPath = storage.dataPath + "/userdata/root"
	storage.myFilesPath = storage.dataPath + "/userdata/root/MyFiles"
	storage.sharedPath = storage.dataPath + "/userdata/shared"
	storage.tmpPath = storage.dataPath + "/userdata/tmp"
	storage.ipfsFilesPath = storage.dataPath + "/userdata/ipfs" // signed and encrypted files, that are added to ipfs are stored here

	return &storage
}

func (storage *Storage) Init() {
	os.Mkdir(storage.dataPath, 0770)
	os.MkdirAll(storage.publicFilesPath, 0770)
	os.MkdirAll(storage.publicForPath, 0770)
	os.MkdirAll(storage.capsPath, 0770)
	os.MkdirAll(storage.capsGAPath, 0770)
	os.MkdirAll(storage.fileRootPath, 0770)
	os.MkdirAll(storage.myFilesPath, 0770)
	os.MkdirAll(storage.sharedPath, 0770)
	os.MkdirAll(storage.tmpPath, 0770)
	os.MkdirAll(storage.ipfsFilesPath, 0770)
}

func (storage *Storage) ExistsFile(filePath string) bool {
	if _, err := os.Stat(filePath); err == nil {
		return true
	}
	return false
}

func (storage *Storage) GetUserFilesPath() string {
	return storage.fileRootPath
}

// It builds up the file repo based on saved data. One part of files
// comes from capabilities which can be found in data/userdata/caps.
// These files contain information about files that are shared with the
// current user. The function appends the representation of those shared
// files into the file structure and checks if they have been updated since
// last time or not. The other half of files comes from data/userdata/shared.
// These files are JSON representation of a FilePTP that were shared by the
// user.
func (storage *Storage) BuildRepo(username string, network *nw.Network, ipfs *ipfs.IPFS) ([]*FilePTP, error) {
	glog.Info("Building repo...")
	var repo []*FilePTP
	// read capabilities from caps and try to refresh them
	entries, err := ioutil.ReadDir(storage.capsPath)
	if err != nil {
		return nil, fmt.Errorf("could not read directory '%s': Storage.BuildRepo: %s", storage.capsPath, err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue // do not care about directories
		}
		path := storage.capsPath + "/" + entry.Name()
		file, err := NewFile(path)
		if err != nil {
			glog.Warningf("invalid file '%s': Storage.BuildRepo: %s\n", path, err)
			continue // do not care about trash files
		}
		if err := file.Refresh(storage, ipfs); err != nil {
			return nil, fmt.Errorf("could not refresh file '%s': Storage.BuildRepo: %s", file.Name, err)
		}
		repo = append(repo, file)
	}
	glog.Info("Repo ready")
	return repo, nil
}

func (storage *Storage) CopyFileIntoPublicDir(filePath string) error {
	publicFilePath := storage.publicFilesPath + "/" + path.Base(filePath)
	return CopyFile(filePath, publicFilePath)
}

func (storage *Storage) CopyFileIntoMyFiles(filePath string) (string, error) {
	newFilePath := storage.myFilesPath + "/" + path.Base(filePath)
	return newFilePath, CopyFile(filePath, newFilePath)
}

func (storage *Storage) CopyFileIntoGroupFiles(filePath, groupName string) (string, error) {
	groupFilesPath := storage.fileRootPath + "/" + groupName
	os.Mkdir(groupFilesPath, 0770)
	newFilePath := groupFilesPath + "/" + path.Base(filePath)
	return newFilePath, CopyFile(filePath, newFilePath)
}

func (storage *Storage) createFileForUser(userID common.Address, capName string, data []byte, network *nw.Network) error {
	userIDBase64 := base64.StdEncoding.EncodeToString(userID[:])
	forUserPath := storage.publicForPath + "/" + userIDBase64
	err := os.MkdirAll(forUserPath, 0770)
	if err != nil {
		glog.Warningf("error while creating dir: Storage.createFileForUser: %s", err) /* TODO check for permission errors */
	}
	_, publicKey, _, _, err := network.GetUser(userID)
	if err != nil {
		return fmt.Errorf("could not get public boxing key: Storage.createFileForUser: %s", err)
	}
	boxer := crypto.AnonymPublicKey{Value: &publicKey}
	encData, err := boxer.Seal(data)
	if err != nil {
		return fmt.Errorf("error while encryption: Storage.createFileForUser: %s", err)
	}
	if err := ioutil.WriteFile(forUserPath+"/"+capName+".json", encData, 0644); err != nil {
		return fmt.Errorf("could not write file: Storage.createFileForUser: %s", err)
	}
	return nil
}

// +------------------------------+
// |   Group specific functions   |
// +------------------------------+

// Gets all the locally stored group access capabilities from
// directory data/userdata/caps/GA/
func (storage *Storage) GetGroupCAPs() ([]GroupAccessCAP, error) {
	var caps []GroupAccessCAP
	// read capabilities from caps and try to refresh them
	entries, err := ioutil.ReadDir(storage.capsGAPath)
	if err != nil {
		return caps, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue // do not care about directories
		}
		filePath := storage.capsGAPath + "/" + entry.Name()
		capBytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			glog.Warning("could not read file '%s': Storage.GetGroupCAPs: %s", filePath, err)
			continue
		}
		var cap GroupAccessCAP
		if err := json.Unmarshal(capBytes, &cap); err != nil {
			glog.Warning("could not unmarshal group cap: Storage.GetGroupCAPs: %s", err)
			continue
		}
		cap.Boxer.RNG = rand.Reader
		caps = append(caps, cap)
	}
	return caps, nil
}

// Creates the directory structure needed by a group
func (storage *Storage) CreateGroupStorage(groupName string) {
	os.MkdirAll(storage.publicForPath+"/"+groupName, 0770)
	os.MkdirAll(storage.fileRootPath+"/"+groupName, 0770)
}

func (storage *Storage) DownloadGroupFile(fileGroup *FileGroup, groupname string, boxer *crypto.SymmetricKey, ipfs *ipfs.IPFS) error {
	// TODO: choice: download to ipfsFiles and host it
	tmpFilePath, err := storage.downloadToTmp(fileGroup.IPFSHash, ipfs)
	if err != nil {
		return fmt.Errorf("could not donwload to tmp: storage.DownloadGroupFile: %s", err)
	}
	encFileBytes, err := ioutil.ReadFile(tmpFilePath)
	if err != nil {
		return fmt.Errorf("could not read tmp file '%s': storage.DownloadGroupFile: %s", tmpFilePath, err)
	}
	fileBytes, ok := boxer.BoxOpen(encFileBytes)
	if !ok {
		return fmt.Errorf("could not decrypt file: storage.DownloadGroupFile")
	}
	groupPath := storage.fileRootPath + "/" + groupname
	os.Mkdir(groupPath, 0770)
	filePath := groupPath + "/" + fileGroup.Name

	if err := WriteFile(filePath, fileBytes); err != nil {
		return fmt.Errorf("could not file '%s': Storage.DownloadGroupFile: %s", filePath, err)
	}
	return nil
}

// Returns the path of the wanted group meta data
func (storage *Storage) GetGroupDataPath(groupname, data string) string {
	return storage.publicForPath + "/" + groupname + "/" + data
}

// Downloads the given group meta data
func (storage *Storage) DownloadGroupData(groupName, file, ipfsHash string, ipfs *ipfs.IPFS, network *nw.Network) error {
	filePath := storage.publicForPath + "/" + groupName + "/" + file
	return ipfs.Get(filePath, ipfsHash)
}

func (storage *Storage) CreateGroupAccessCAPForUser(userID common.Address, group string, key crypto.SymmetricKey, network *nw.Network) error {
	cap := GroupAccessCAP{group, key}
	capBytes, err := json.Marshal(cap)
	if err != nil {
		return fmt.Errorf("could not marshal group access capability: Storage.CreateGroupAccessCAPForUser: %s", err)
	}
	if err := storage.createFileForUser(userID, group, capBytes, network); err != nil {
		return fmt.Errorf("could not create cap for user: Storage.CreateGroupAccessCAPForUser: %s", err)
	}
	return nil
}

// Stores the given group meta data in data/public/for/group/
func (storage *Storage) SaveGroupData(groupName, fileName string, boxer crypto.SymmetricKey, data []byte) error {
	// group data goes always into the /public/for/group/ directory
	filePath := storage.publicForPath + "/" + groupName + "/" + fileName
	encData := boxer.BoxSeal(data)
	return WriteFile(filePath, encData)
}

func (storage *Storage) StoreGroupAccessCAP(group string, key crypto.SymmetricKey) error {
	cap := GroupAccessCAP{group, key}
	return cap.Store(storage)
}

func (storage *Storage) DownloadGroupAccessCAP(fromUserID, userID common.Address, capName string, boxer *crypto.AnonymBoxer, network *nw.Network, ipfs *ipfs.IPFS) (*GroupAccessCAP, error) {
	capBytes, err := downloadCAP(fromUserID, userID, capName, boxer, storage, network, ipfs)
	if err != nil {
		return nil, fmt.Errorf("could not download group cap: Storage.DownloadGroupAccessCAP: %s", err)
	}
	var cap GroupAccessCAP
	if err := json.Unmarshal(capBytes, &cap); err != nil {
		return nil, fmt.Errorf("could not unmarshal ga cap: Storage.StoreGroupAccessCAP: %s", err)
	}
	return &cap, nil
}

// +------------------------------+
// |       Helper functions       |
// +------------------------------+

func (storage *Storage) PublishPublicDir(ipfs *ipfs.IPFS) error {
	glog.Info("Publishing...")
	publicDir := storage.dataPath + "/public"
	merkleNodes, err := ipfs.AddDir(publicDir)
	if err != nil {
		return fmt.Errorf("could not ipfs add dir: Storage.PublishPublicDir: %s", err)
	}
	for _, mn := range merkleNodes {
		if strings.Compare(mn.Name, "public") == 0 {
			if err := ipfs.NamePublish(mn.Hash); err != nil {
				return fmt.Errorf("could not ipfs name publish: Storage.PublishPublicDir: %s", err)
			}
			break
		}
	}
	glog.Info("Publishing ended")
	return nil
}

func (storage *Storage) downloadToTmp(ipfsHash string, ipfs *ipfs.IPFS) (string, error) {
	filePath := storage.tmpPath + "/" + path.Base(ipfsHash)
	if err := ipfs.Get(filePath, ipfsHash); err != nil {
		return "", fmt.Errorf("could not ipfs get into tmp path '%s': Storage.downloadToTmp: %s", filePath, err)
	}
	return filePath, nil
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
		glog.Warningf("file '%s' already exists: WriteFile", filePath)
	}
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		return fmt.Errorf("could not write file '%s': WriteFile: %s", filePath, err)
	}
	f.Sync()
	return nil
}

func FileExists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}
