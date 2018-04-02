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
	us.capsGAPath = us.dataPath + "/userdata/caps/GA"
	us.storagePath = us.dataPath + "/userdata/root"
	us.sharedPath = us.dataPath + "/userdata/shared"
	us.tmpPath = us.dataPath + "/userdata/tmp"

	os.Mkdir(us.dataPath, 0770)
	os.MkdirAll(us.publicFilesPath, 0770)
	os.MkdirAll(us.publicForPath, 0770)
	os.MkdirAll(us.capsPath, 0770)
	os.MkdirAll(us.capsGAPath, 0770)
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
		changed, err := cap.Refresh(us, ipfs)
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

func (s *Storage) CreateFileReadCAPForUser(f *File, forUserPath, ipfsHash string, boxer *crypto.BoxingKeyPair, network *nw.Network) error {
	user := path.Base(forUserPath)
	cap := ReadCAP{path.Base(f.Path), f.IPNSPath, ipfsHash, f.Owner, f.VerifyKey}
	capBytes, err := json.Marshal(cap)
	if err != nil {
		return err
	}
	return s.createFileForUser(user, path.Base(f.Path), capBytes, boxer, network)
}

func (s *Storage) createFileForUser(user, capName string, data []byte, boxer *crypto.BoxingKeyPair, network *nw.Network) error {
	forUserPath := s.publicForPath + "/" + user
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
func (s *Storage) GetGroupCAPs() ([]GroupAccessCAP, error) {
	var caps []GroupAccessCAP
	// read capabilities from caps and try to refresh them
	entries, err := ioutil.ReadDir(s.capsGAPath)
	if err != nil {
		return caps, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue // do not care about directories
		}
		capBytes, err := ioutil.ReadFile(s.capsGAPath + "/" + entry.Name())
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
func (s *Storage) CreateGroupStorage(groupName string) {
	os.MkdirAll(s.publicForPath+"/"+groupName, 0770)
	os.MkdirAll(s.storagePath+"/"+groupName, 0770)
}

// Gets the ipfs hash of the given group meta data stored locally
// like members, ACL, etc... If it is not present, it returns with
// a null string
func (s *Storage) GetLocalGroupDataHash(groupname, data string, ipfs *ipfs.IPFS) (string, error) {
	filePath := s.publicForPath + "/" + groupname + "/" + data
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
func (s *Storage) GetGroupDataPath(groupname, data string) string {
	return s.publicForPath + "/" + groupname + "/" + data
}

// Checks if the given meta data is present and if so, whether it
// changed (accorting to other members) since last time or not. If
// changes are detected, it returns with the new ipfs hash of the
// meta data
func (s *Storage) GroupDataChanged(groupname, data string, activeMembers []string, ipfs *ipfs.IPFS, network *nw.Network) (string, error) {
	memberHash, err := s.GetLocalGroupDataHash(groupname, "members.json", ipfs)
	if err != nil {
		return "", err
	}
	// only active member is current user
	if len(activeMembers) == 1 {
		fmt.Println("only active")
		return "", nil
	}
	membersAgreeOnData := 0
	commonHash := ""
	for _, member := range activeMembers {
		ipns, err := network.GetUserIPFSAddr(member)
		if err != nil {
			return "", err
		}
		ipnsPath := "/ipns/" + ipns + "/for/" + groupname + "/" + data
		ipfsHash, err := ipfs.Resolve(ipnsPath)
		if err != nil {
			return "", err
		}
		if strings.Compare(commonHash, "") == 0 {
			commonHash = ipfsHash
		} else if strings.Compare(commonHash, ipfsHash) == 0 {
			membersAgreeOnData += 1
		}
	}
	if membersAgreeOnData < len(activeMembers)/2 {
		return "", errors.New("members do not agree on data: " + data)
	}
	if strings.Compare(memberHash, commonHash) == 0 {
		return "", nil
	}
	return commonHash, nil
}

// Downloads the given group meta data
func (s *Storage) DownloadGroupData(groupName, file, ipfsHash string, ipfs *ipfs.IPFS, network *nw.Network) error {
	filePath := s.publicForPath + "/" + groupName + "/" + file
	return ipfs.Get(filePath, ipfsHash)
}

func (s *Storage) CreateGroupAccessCAPForUser(user, group string, key crypto.SymmetricKey, boxer *crypto.BoxingKeyPair, network *nw.Network) error {
	cap := GroupAccessCAP{group, key}
	capBytes, err := json.Marshal(cap)
	if err != nil {
		return err
	}
	return s.createFileForUser(user, group, capBytes, boxer, network)
}

// Stores the given group meta data in data/public/for/group/
func (s *Storage) SaveGroupData(groupName, fileName string, boxer crypto.SymmetricKey, data []byte) error {
	// group data goes always into the /public/for/group/ directory
	filePath := s.publicForPath + "/" + groupName + "/" + fileName
	encData := boxer.BoxSeal(data)
	return WriteFile(filePath, encData)
}

// +------------------------------+
// |     Capability functions     |
// +------------------------------+

func (s *Storage) StoreGroupAccessCAP(group string, key crypto.SymmetricKey) error {
	cap := GroupAccessCAP{group, key}
	filePath := s.capsGAPath + "/" + group + ".json"
	return cap.Store(filePath)
}

func (s *Storage) DownloadReadCAP(fromUser, username, capName string, boxer *crypto.BoxingKeyPair, network *nw.Network, ipfs *ipfs.IPFS) (*ReadCAP, error) {
	capBytes, err := s.downloadCAP(s.capsPath, fromUser, username, capName, boxer, network, ipfs)
	if err != nil {
		return nil, err
	}
	var cap *ReadCAP
	if err := json.Unmarshal(capBytes, cap); err != nil {
		return nil, err
	}
	return cap, nil
}

func (s *Storage) DownloadGroupAccessCAP(fromUser, username, capName string, boxer *crypto.BoxingKeyPair, network *nw.Network, ipfs *ipfs.IPFS) (*GroupAccessCAP, error) {
	capBytes, err := s.downloadCAP(s.capsGAPath, fromUser, username, capName, boxer, network, ipfs)
	if err != nil {
		return nil, err
	}
	var cap GroupAccessCAP
	if err := json.Unmarshal(capBytes, &cap); err != nil {
		return nil, err
	}
	return &cap, nil
}

// Downloads the capability identified by capName from
// /ipns/from/for/username/capName
func (s *Storage) downloadCAP(basePath, fromUser, username, capName string, boxer *crypto.BoxingKeyPair, network *nw.Network, ipfs *ipfs.IPFS) ([]byte, error) {
	// get address and key
	ipfsAddr, err := network.GetUserIPFSAddr(fromUser)
	if err != nil {
		return nil, err
	}
	otherPK, err := network.GetUserBoxingKey(fromUser)
	if err != nil {
		return nil, err
	}
	ipnsPath := "/ipns/" + ipfsAddr + "/for/" + username + "/" + capName
	// download cap file
	tmpFilePath := s.tmpPath + "/" + capName
	err = ipfs.Get(tmpFilePath, ipnsPath)
	if err != nil {
		return nil, err
	}
	capFilePath := basePath + "/" + capName
	bytesEnc, err := ioutil.ReadFile(tmpFilePath)
	if err != nil {
		return nil, err
	}
	bytesDecr, success := boxer.BoxOpen(bytesEnc, &otherPK)
	if !success {
		fmt.Println("trying decrypt cap")
		fmt.Println(fromUser)
		fmt.Println(ipnsPath)
		return nil, errors.New("could not decrypt capability")
	}
	os.Remove(tmpFilePath)
	if err := WriteFile(capFilePath, bytesDecr); err != nil {
		return nil, err
	}
	return bytesDecr, nil
}

// +------------------------------+
// |       Helper functions       |
// +------------------------------+

func (us *Storage) PublishPublicDir(ipfs *ipfs.IPFS) error {
	fmt.Println("[*] Publishing...")
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

func FileExists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}
