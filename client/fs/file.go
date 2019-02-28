package fs

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"sync"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/getlantern/deepcopy"
	"github.com/golang-collections/collections/stack"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/sergi/go-diff/diffmatchpatch"

	"ipfs-share/client/fs/caps"
	"ipfs-share/collections"
	"ipfs-share/crypto"
	ipfsapi "ipfs-share/ipfs"
	"ipfs-share/utils"
)

// IFile is an interface for the files
// which can be shared
type IFile interface {
	Share()
}

// File represents a file that
// is shared in a peer to peer mode
type File struct {
	Cap            *caps.FileCap
	PendingChanges *caps.FileCap
	DataPath       string
	CapPath        string
	OrigPath       string
	lock           sync.RWMutex
}

func (f *File) Id() collections.IIdentifier {
	return collections.NewStringId(f.Cap.FileName)
}

func NewGroupFile(filePath string, writeAccessList []ethcommon.Address, groupId string, storage *Storage) (*File, error) {
	if writeAccessList == nil {
		return nil, errors.New("writeAccessList can not be nil")
	}

	fileName := path.Base(filePath)

	cap, err := caps.NewGroupFileCap(fileName, writeAccessList)
	if err != nil {
		return nil, errors.Wrap(err, "could not create cap for NewFile")
	}
	var pendingChanges *caps.FileCap
	if err := deepcopy.Copy(&pendingChanges, cap); err != nil {
		return nil, errors.Wrap(err, "could not deep copy cap")
	}

	idString := base64.URLEncoding.EncodeToString(cap.Id[:])
	capPath := storage.GroupFileCapDir(groupId) + idString
	origPath := storage.GroupFileOrigDir(groupId) + fileName

	file := &File{
		Cap:            cap,
		PendingChanges: pendingChanges,
		DataPath:       filePath,
		CapPath:        capPath,
		OrigPath:       origPath,
	}

	if err := file.SaveMetadata(); err != nil {
		return nil, errors.Wrap(err, "could not save file meta data")
	}

	return file, nil
}

func NewGroupFileFromCap(cap *caps.FileCap, groupId string, storage *Storage) (*File, error) {
	idString := base64.URLEncoding.EncodeToString(cap.Id[:])
	capPath := storage.GroupFileCapDir(groupId) + idString
	dataPath := storage.GroupFileDataDir(groupId) + cap.FileName
	pendingPath := storage.GroupFileOrigDir(groupId) + idString

	var pendingChanges *caps.FileCap
	if err := deepcopy.Copy(&pendingChanges, cap); err != nil {
		return nil, errors.Wrap(err, "could not deep copy cap")
	}

	file := &File{
		Cap:            cap,
		PendingChanges: pendingChanges,
		DataPath:       dataPath,
		CapPath:        capPath,
		OrigPath:       pendingPath,
	}

	return file, nil
}


// LoadPTPFile loads a File from the disk
func LoadPTPFile(filePath string) (*File, error) {
	bytesFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not read file '%s': LoadPTPFile: %s", filePath, err)
	}
	var file File
	if err := json.Unmarshal(bytesFile, &file); err != nil {
		return nil, fmt.Errorf("could not unmarshal file '%s': LoadPTPFile: %s", filePath, err)
	}
	return &file, nil
}

// NewFileFromCap creates a new File instance from a shared
// capability
func NewFileFromCap(dataDir, capDir string, cap *caps.FileCap, ipfs ipfsapi.IIpfs, storage *Storage) (*File, error) {
	baseName := base64.URLEncoding.EncodeToString(cap.Id[:])
	dataPath := dataDir + baseName
	capPath := capDir + baseName

	file := &File{
		Cap:      cap,
		DataPath: dataPath,
		CapPath:   capPath,
	}

	if err := file.SaveMetadata(); err != nil {
		return nil, errors.Wrapf(err, "could not save file '%s': NewFileFromCap", cap.FileName)
	}

	go file.Download(storage, ipfs)

	return file, nil
}

func (f *File) Update(cap *caps.FileCap, storage *Storage, ipfs ipfsapi.IIpfs) error {
	f.lock.Lock()
	defer f.lock.Unlock()

	oldIpfsHash := f.Cap.IpfsHash
	f.Cap = cap
	if strings.Compare(oldIpfsHash, cap.IpfsHash) != 0 {
		if err := f.SaveMetadata(); err != nil {
			return errors.Wrap(err, "could not save file meta data")
		}

		if err := deepcopy.Copy(&f.PendingChanges, f.Cap); err != nil {
			return errors.Wrap(err, "could not deep copy cap top pending changes")
		}

		go f.Download(storage, ipfs)
	}
	return nil
}

// Downloads, decrypts and verifies the content of file from Ipfs
func (f *File) Download(storage *Storage, ipfs ipfsapi.IIpfs) {
	dmp := diffmatchpatch.New()
	patchStack := stack.New()

	currentDiffIpfsHash := f.Cap.IpfsHash
	currentDiffBoxer := f.Cap.DataKey
	currentStr := ""
	var origHash []byte
	if utils.FileExists(f.OrigPath) {
		origData, err := ioutil.ReadFile(f.OrigPath)
		if err != nil {
			glog.Errorf("could not read original file: %s", err)
			return
		}

		hasher := tribecrypto.NewKeccak256Hasher()
		origHash = hasher.Sum(origData)

		currentStr = string(origData)
	}

	for {
		data, err := storage.DownloadAndDecryptWithFileBoxer(currentDiffBoxer, currentDiffIpfsHash, ipfs)
		if err != nil {
			glog.Error("could not download and decrypt diff node")
			return
		}

		diff, err := DecodeDiffNode(data)
		if err != nil {
			glog.Errorf("download err: could not decode diff node: %s", err)
			return
		}

		patch := dmp.PatchMake(diff.Diff)
		patchStack.Push(patch)

		// there is no next element
		if strings.Compare(diff.Next, "") == 0 {
			break
		}
		// we found our state
		if bytes.Equal(diff.Hash, origHash) {
			break
		}

		currentDiffIpfsHash = diff.Next
		currentDiffBoxer = diff.NextBoxer
	}

	for {
		patchInt := patchStack.Pop()
		if patchInt == nil {
			break
		}

		currentStr, _ = dmp.PatchApply(patchInt.([]diffmatchpatch.Patch), currentStr)
	}

	if err := utils.WriteFile(f.OrigPath, []byte(currentStr)); err != nil {
		glog.Errorf("download err: could not write orig file: %s", err)
	}
	if err := utils.WriteFile(f.DataPath, []byte(currentStr)); err != nil {
		glog.Errorf("download err: could not data orig file: %s", err)
	}
}

// SaveMetadata saves FileMetaData to disk
func (f *File) SaveMetadata() error {
	jsonBytes, err := json.Marshal(f)
	if err != nil {
		return errors.Wrapf(err, "could not marshal file '%s': File.save", f.Cap.Id)
	}
	if err := utils.WriteFile(f.CapPath, jsonBytes); err != nil {
		return errors.Wrapf(err, "could not write file '%s': File.save", f.Cap.Id)
	}

	return nil
}

func GetCapListFromFileList(files []*File) []*caps.FileCap {
	var l []*caps.FileCap
	for _, file := range files {
		l = append(l, file.Cap)
	}
	return l
}

func (f *File) GrantWriteAccess(user, target ethcommon.Address) error {
	f.lock.Lock()
	defer f.lock.Unlock()

	hasWriteAccess := false
	for _, hasW := range f.PendingChanges.WriteAccessList {
		if bytes.Equal(hasW.Bytes(), user.Bytes()) {
			hasWriteAccess = true
		}

		if bytes.Equal(hasW.Bytes(), target.Bytes()) {
			return errors.New("user already has Write access")
		}
	}

	if !hasWriteAccess {
		return errors.New("user has no Write access")
	}

	f.PendingChanges.WriteAccessList = append(f.PendingChanges.WriteAccessList, target)
	if err := f.SaveMetadata(); err != nil {
		return errors.Wrap(err, "could not save file meta data")
	}

	return nil
}

func (f *File) RevokeWriteAccess(user, target ethcommon.Address) error {
	f.lock.Lock()
	defer f.lock.Unlock()

	userHasWriteAccess := false
	targetHasWriteAccess := false
	for _, hasW := range f.PendingChanges.WriteAccessList {
		if bytes.Equal(hasW.Bytes(), user.Bytes()) {
			userHasWriteAccess = true
		}

		if bytes.Equal(hasW.Bytes(), target.Bytes()) {
			targetHasWriteAccess = true
		}

		if userHasWriteAccess && targetHasWriteAccess {
			break
		}
	}

	if !userHasWriteAccess {
		return errors.New("user has no Write access")
	}
	if !targetHasWriteAccess {
		return errors.New("target has no Write access")
	}

	for i, hasW := range f.PendingChanges.WriteAccessList {
		if bytes.Equal(hasW.Bytes(), target.Bytes()) {
			f.PendingChanges.WriteAccessList = append(f.PendingChanges.WriteAccessList[:i], f.PendingChanges.WriteAccessList[i+1:]...)
			break
		}
	}

	return nil
}

func (f *File) diff(boxer tribecrypto.FileBoxer) (*DiffNode, error) {
	f.lock.RLock()
	defer f.lock.RUnlock()

	dmp := diffmatchpatch.New()

	diff := &DiffNode{
		Hash:      nil,
		Next:      "",
		NextBoxer: boxer,
	}

	originalStr := ""

	if utils.FileExists(f.OrigPath) {
		originalData, err := ioutil.ReadFile(f.OrigPath)
		if err != nil {
			return nil, errors.Wrap(err, "could not read original file")
		}

		hasher := tribecrypto.NewKeccak256Hasher()
		hash := hasher.Sum(originalData)

		diff.Hash = hash

		originalStr = string(originalData)
	}

	currentData, err := ioutil.ReadFile(f.DataPath)
	if err != nil {
		return nil, errors.Wrap(err, "could not read current file")
	}

	diffs := dmp.DiffMain(originalStr, string(currentData), true)

	diff.Diff = diffs
	diff.Next = f.Cap.IpfsHash

	return diff, nil
}

// uploads the diff
func (f *File) UploadDiff(ipfs ipfsapi.IIpfs) (string, error) {
	f.lock.RLock()
	defer f.lock.RUnlock()

	diff, err := f.diff(f.PendingChanges.DataKey)
	if err != nil {
		return "", errors.Wrap(err, "could not get file diff")
	}

	encData, err := diff.Encrypt(f.PendingChanges.DataKey)
	if err != nil {
		return "", errors.Wrap(err, "could not encrypt file diff")
	}

	newIpfsHash, err := ipfs.Add(encData)
	if err != nil {
		return "", errors.Wrap(err, "could not ipfs add file diff")
	}

	return newIpfsHash, nil
}

func (f *File) ChangeKey(ipfs ipfsapi.IIpfs) error {
	var newKey [32]byte
	if _, err := rand.Read(newKey[:]); err != nil {
		return errors.Wrap(err, "could not read from crypto.rand")
	}

	f.PendingChanges.DataKey = tribecrypto.FileBoxer{Key: newKey}

	ipfsHash, err := f.UploadDiff(ipfs)
	if err != nil {
		return errors.Wrap(err, "could not upload file diff to ipfs")
	}

	f.PendingChanges.IpfsHash = ipfsHash

	return nil
}