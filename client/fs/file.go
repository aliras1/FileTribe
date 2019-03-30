// Copyright (c) 2019 Laszlo Sari
//
// FileTribe is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// FileTribe is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//

package fs

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"sync"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/golang-collections/collections/stack"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/sergi/go-diff/diffmatchpatch"

	"github.com/aliras1/FileTribe/client/fs/meta"
	ipfsapi "github.com/aliras1/FileTribe/ipfs"
	"github.com/aliras1/FileTribe/tribecrypto"
	"github.com/aliras1/FileTribe/utils"
)

// IFile is an interface for the files
// which can be shared
type IFile interface {
	Share()
}

// File represents a file that
// is shared in a peer to peer mode
type File struct {
	Cap            *meta.FileMeta
	PendingChanges *meta.FileMeta
	DataPath       string
	CapPath        string
	OrigPath       string
	lock           sync.RWMutex
}


func NewGroupFile(filePath string, writeAccessList []ethcommon.Address, groupId string, storage *Storage) (*File, error) {
	if writeAccessList == nil {
		return nil, errors.New("writeAccessList can not be nil")
	}

	fileName := path.Base(filePath)
	fileMeta, err := meta.NewGroupFileMeta(fileName, writeAccessList)
	if err != nil {
		return nil, errors.Wrap(err, "could not create fileMeta for NewFile")
	}

	var pendingChanges meta.FileMeta
	if err := deepcopy(&pendingChanges, fileMeta); err != nil {
		return nil, errors.Wrap(err, "could not deep copy fileMeta")
	}

	capPath := storage.GroupFileCapDir(groupId) + fileMeta.FileName
	origPath := storage.GroupFileOrigDir(groupId) + fileMeta.FileName

	file := &File{
		Cap:            fileMeta,
		PendingChanges: &pendingChanges,
		DataPath:       filePath,
		CapPath:        capPath,
		OrigPath:       origPath,
	}

	if err := file.SaveMetadata(); err != nil {
		return nil, errors.Wrap(err, "could not save file meta data")
	}

	return file, nil
}

func NewGroupFileFromCap(cap *meta.FileMeta, groupId string, storage *Storage) (*File, error) {
	capPath := storage.GroupFileCapDir(groupId) + cap.FileName
	dataPath := storage.GroupFileDataDir(groupId) + cap.FileName
	pendingPath := storage.GroupFileOrigDir(groupId) + cap.FileName

	var pendingChanges *meta.FileMeta
	if err := deepcopy(&pendingChanges, cap); err != nil {
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
func NewFileFromCap(dataDir, capDir string, cap *meta.FileMeta, ipfs ipfsapi.IIpfs, storage *Storage) (*File, error) {
	dataPath := dataDir + cap.FileName
	capPath := capDir + cap.FileName

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

func (f *File) Update(cap *meta.FileMeta, storage *Storage, ipfs ipfsapi.IIpfs) error {
	f.lock.Lock()
	defer f.lock.Unlock()

	oldIpfsHash := f.Cap.IpfsHash
	f.Cap = cap
	if strings.Compare(oldIpfsHash, cap.IpfsHash) != 0 {
		if err := f.SaveMetadata(); err != nil {
			return errors.Wrap(err, "could not save file meta data")
		}

		if err := deepcopy(&f.PendingChanges, f.Cap); err != nil {
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

		origHash = ethcrypto.Keccak256(origData)
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

	if err := utils.CreateAndWriteFile(f.OrigPath, []byte(currentStr)); err != nil {
		glog.Errorf("download err: could not write orig file: %s", err)
	}
	if err := utils.CreateAndWriteFile(f.DataPath, []byte(currentStr)); err != nil {
		glog.Errorf("download err: could not data orig file: %s", err)
	}
}

// SaveMetadata saves FileMetaData to disk
func (f *File) SaveMetadata() error {
	jsonBytes, err := json.Marshal(f)
	if err != nil {
		return errors.Wrapf(err, "could not marshal file '%s'", f.Cap.FileName)
	}
	glog.Infof("%v", f)
	if err := utils.CreateAndWriteFile(f.CapPath, jsonBytes); err != nil {
		return errors.Wrapf(err, "could not write file '%s'", f.Cap.FileName)
	}

	return nil
}

func GetCapListFromFileList(files []*File) []*meta.FileMeta {
	var l []*meta.FileMeta
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

		diff.Hash = ethcrypto.Keccak256(originalData)

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

	var newKey [32]byte
	if _, err := rand.Read(newKey[:]); err != nil {
		return "", errors.Wrap(err, "could not read from crypto.rand")
	}

	f.PendingChanges.DataKey = tribecrypto.FileBoxer{Key: newKey}

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

func deepcopy(dst, src interface{}) error {
	data, _ := json.Marshal(src)
	return json.Unmarshal(data, dst)
}