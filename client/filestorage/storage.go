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

	"crypto/rand"
	"ipfs-share/crypto"
	"ipfs-share/ipfs"
	nw "ipfs-share/network"
	"log"
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
	ipfsPath        string
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
	storage.ipfsPath = storage.dataPath + "/userdata/ipfs" // signed and encrypted files, that are added to ipfs are stored here

	os.Mkdir(storage.dataPath, 0770)
	os.MkdirAll(storage.publicFilesPath, 0770)
	os.MkdirAll(storage.publicForPath, 0770)
	os.MkdirAll(storage.capsPath, 0770)
	os.MkdirAll(storage.capsGAPath, 0770)
	os.MkdirAll(storage.fileRootPath, 0770)
	os.MkdirAll(storage.myFilesPath, 0770)
	os.MkdirAll(storage.sharedPath, 0770)
	os.MkdirAll(storage.tmpPath, 0770)
	os.MkdirAll(storage.ipfsPath, 0770)

	return &storage
}

func (storage *Storage) ExistsFile(filePath string) bool {
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
// These files are JSON representation of a FilePTP that were shared by the
// user.
func (storage *Storage) BuildRepo(username string, boxer *crypto.BoxingKeyPair, network *nw.Network, ipfs *ipfs.IPFS) ([]*FilePTP, error) {
	var repo []*FilePTP
	// read capabilities from caps and try to refresh them
	entries, err := ioutil.ReadDir(storage.capsPath)
	if err != nil {
		return []*FilePTP{}, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue // do not care about directories
		}
		cap, err := NewReadCAPFromFile(storage.capsPath + "/" + entry.Name())
		if err != nil {
			continue // do not care about trash files
		}
		changed, err := cap.Refresh(username, boxer, storage, network, ipfs)
		if err != nil {
			return nil, err
		}
		file, err := NewFileFromCAP(cap, storage, ipfs)
		if changed {
			fmt.Println("changed")
			storage.DownloadFileFromCap(cap, ipfs)
		}
		repo = append(repo, file)
	}

	// read share information
	entries, err = ioutil.ReadDir(storage.sharedPath)
	if err != nil {
		return []*FilePTP{}, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		file, err := NewFileFromShared(storage.sharedPath + "/" + entry.Name())
		if err != nil {
			fmt.Println(err)
			continue
		}
		repo = append(repo, file)
	}
	return repo, nil
}

// Downloads and verifies the file, described in the capability.
func (storage *Storage) DownloadFileFromCap(cap *ReadCAP, ipfs *ipfs.IPFS) error {
	tmpFilePath := storage.tmpPath + "/" + path.Base(cap.FileName)
	err := ipfs.Get(tmpFilePath, cap.IPFSHash)
	if err != nil {
		return err
	}
	bytesSignedFile, err := ioutil.ReadFile(tmpFilePath)
	if err != nil {
		return err
	}
	os.Remove(tmpFilePath)
	// make a directory to the owner
	dirPath := storage.fileRootPath + "/" + cap.Owner
	os.MkdirAll(dirPath, 0770)

	bytesRawFile, ok := cap.VerifyKey.Open(nil, bytesSignedFile)
	if !ok {
		return errors.New("by downloadFileFromCap(): could not verify file")
	}
	filePath := dirPath + "/" + cap.FileName
	if err := WriteFile(filePath, bytesRawFile); err != nil {
		return err
	}
	return nil
}

func (storage *Storage) CopyFileIntoPublicDir(filePath string) error {
	publicFilePath := storage.publicFilesPath + "/" + path.Base(filePath)
	return CopyFile(filePath, publicFilePath)
}

func (storage *Storage) CopyFileIntoMyFiles(filePath string) (string, error) {
	newFilePath := storage.myFilesPath + "/" + path.Base(filePath)
	return newFilePath, CopyFile(filePath, newFilePath)
}

// Signs the files with the Write key and then the function adds
// it to IPFS. The function returns with the with the IPFS hash
// of the file
func (storage *Storage) SignAndAddFileToIPFS(filePath string, writeKey crypto.SecretSigningKey, ipfs *ipfs.IPFS) (string, error) {
	ipfsPath := storage.ipfsPath + "/" + path.Base(filePath)
	bytesFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	signedFile := writeKey.Sign(nil, bytesFile)
	if err := WriteFile(ipfsPath, signedFile); err != nil {
		return "", err
	}
	merkleNode, err := ipfs.AddFile(ipfsPath)
	if err != nil {
		return "", err
	}
	return "/ipfs/" + merkleNode.Hash, nil
}

// Saves a FilePTP object (containing meta-data of a file) in json format
// locally
func (storage *Storage) StoreFileMetaData(f *FilePTP) error {
	byteJson, err := json.Marshal(f)
	if err != nil {
		return err
	}
	return WriteFile(storage.sharedPath+"/"+path.Base(f.Name)+".json", byteJson)
}

func (storage *Storage) CreateFileReadCAPForUser(f *FilePTP, username, ipfsHash string, boxer *crypto.BoxingKeyPair, network *nw.Network) error {
	cap := ReadCAP{path.Base(f.Name), ipfsHash, f.Owner, f.VerifyKey}
	capBytes, err := json.Marshal(cap)
	if err != nil {
		return err
	}
	return storage.createFileForUser(username, path.Base(f.Name), capBytes, boxer, network)
}

func (storage *Storage) createFileForUser(user, capName string, data []byte, boxer *crypto.BoxingKeyPair, network *nw.Network) error {
	forUserPath := storage.publicForPath + "/" + user
	err := os.MkdirAll(forUserPath, 0770)
	if err != nil {
		fmt.Println(err) /* TODO check for permission errors */
	}
	otherPK, err := network.GetUserBoxingKey(path.Base(forUserPath))
	if err != nil {
		return err
	}
	encData := boxer.BoxSeal(data, &otherPK)
	return ioutil.WriteFile(forUserPath+"/"+capName+".json", encData, 0644)
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
		capBytes, err := ioutil.ReadFile(storage.capsGAPath + "/" + entry.Name())
		if err != nil {
			log.Println(err)
			continue
		}
		var cap GroupAccessCAP
		if err := json.Unmarshal(capBytes, &cap); err != nil {
			log.Println(err)
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

// Gets the ipfs hash of the given group meta data stored locally
// like members, ACL, etc... If it is not present, it returns with
// a null string
func (storage *Storage) GetLocalGroupDataHash(groupname, data string, ipfs *ipfs.IPFS) (string, error) {
	filePath := storage.publicForPath + "/" + groupname + "/" + data
	if !FileExists(filePath) {
		fmt.Println("do not exist")
		return "", nil
	}
	mn, err := ipfs.AddFile(filePath)
	if err != nil {
		return "", err
	}
	return mn.Hash, nil
}

// Returns the path of the wanted group meta data
func (storage *Storage) GetGroupDataPath(groupname, data string) string {
	return storage.publicForPath + "/" + groupname + "/" + data
}

// Checks if the given meta data is present and if so, whether it
// changed (accorting to other members) since last time or not. If
// changes are detected, it returns with the new ipfs hash of the
// meta data
func (storage *Storage) GroupDataChanged(groupname, data string, activeMembers []string, ipfs *ipfs.IPFS, network *nw.Network) (string, error) {
	memberHash, err := storage.GetLocalGroupDataHash(groupname, "members.json", ipfs)
	if err != nil {
		return "", err
	}
	// only active member is current user
	if len(activeMembers) == 1 {
		fmt.Println("only active member")
		return "", nil
	}
	fmt.Print("active members: ")
	fmt.Println(activeMembers)
	commonHashes := make(map[string]int)
	for _, member := range activeMembers {
		ipns, err := network.GetUserIPFSAddr(member)
		if err != nil {
			log.Println(err)
			continue
		}
		ipnsPath := "/ipns/" + ipns + "/for/" + groupname + "/" + data
		ipfsHash, err := ipfs.Resolve(ipnsPath)
		if err != nil {
			log.Println(err)
			continue
		}
		_, in := commonHashes[ipfsHash]
		if !in {
			commonHashes[ipfsHash] = 1
		} else {
			commonHashes[ipfsHash] += 1
		}
	}
	mostCommonHash := ""
	membersAgreeOnData := 0
	for k, v := range commonHashes {
		if v >= membersAgreeOnData {
			membersAgreeOnData = v
			mostCommonHash = k
		}
	}
	if membersAgreeOnData < len(activeMembers)/2 {
		return "", errors.New("members do not agree on data: " + data)
	}
	if strings.Compare(memberHash, mostCommonHash) == 0 {
		return "", nil
	}
	return mostCommonHash, nil
}

// Downloads the given group meta data
func (storage *Storage) DownloadGroupData(groupName, file, ipfsHash string, ipfs *ipfs.IPFS, network *nw.Network) error {
	filePath := storage.publicForPath + "/" + groupName + "/" + file
	return ipfs.Get(filePath, ipfsHash)
}

func (storage *Storage) CreateGroupAccessCAPForUser(user, group string, key crypto.SymmetricKey, boxer *crypto.BoxingKeyPair, network *nw.Network) error {
	cap := GroupAccessCAP{group, key}
	capBytes, err := json.Marshal(cap)
	if err != nil {
		return err
	}
	return storage.createFileForUser(user, group, capBytes, boxer, network)
}

// Stores the given group meta data in data/public/for/group/
func (storage *Storage) SaveGroupData(groupName, fileName string, boxer crypto.SymmetricKey, data []byte) error {
	// group data goes always into the /public/for/group/ directory
	filePath := storage.publicForPath + "/" + groupName + "/" + fileName
	encData := boxer.BoxSeal(data)
	return WriteFile(filePath, encData)
}

// +------------------------------+
// |     Capability functions     |
// +------------------------------+

func (storage *Storage) StoreGroupAccessCAP(group string, key crypto.SymmetricKey) error {
	cap := GroupAccessCAP{group, key}
	return cap.Store(storage)
}

func (storage *Storage) DownloadGroupAccessCAP(fromUser, username, capName string, boxer *crypto.BoxingKeyPair, network *nw.Network, ipfs *ipfs.IPFS) (*GroupAccessCAP, error) {
	capBytes, err := downloadCAP(fromUser, username, capName, boxer, storage, network, ipfs)
	if err != nil {
		return nil, err
	}
	var cap GroupAccessCAP
	if err := json.Unmarshal(capBytes, &cap); err != nil {
		return nil, err
	}
	return &cap, nil
}

// +------------------------------+
// |       Helper functions       |
// +------------------------------+

func (storage *Storage) PublishPublicDir(ipfs *ipfs.IPFS) error {
	fmt.Println("[*] Publishing...")
	publicDir := storage.dataPath + "/public"
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
	fmt.Println("[*] Publishing ended")
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
		// TODO
		//return err
		fmt.Println(err)
	}
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		return err
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
