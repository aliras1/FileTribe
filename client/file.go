package client

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/golang/glog"
	ipfsapi "ipfs-share/ipfs"
	"ipfs-share/collections"
	"github.com/pkg/errors"
	"ipfs-share/utils"
	"io"
	"bytes"
	"github.com/getlantern/deepcopy"
	"sync"
	"strings"
)

// IFile is an interface for the files
// which can be shared
type IFile interface {
	Share()
}

// File represents a file that
// is shared in a peer to peer mode
type File struct {
	Cap      *FileCap
	PendingChanges *FileCap
	DataPath string
	CapPath string
	PendingPath string
	lock sync.RWMutex
}

func (f *File) Id() collections.IIdentifier {
	return collections.NewStringId(f.Cap.FileName)
}

func NewGroupFile(filePath string, writeAccessList []ethcommon.Address, group IGroup, storage *Storage, ipfs ipfsapi.IIpfs) (*File, error) {
	if writeAccessList == nil {
		return nil, errors.New("writeAccessList can not be nil")
	}

	cap, err := NewGroupFileCap(path.Base(filePath), filePath, writeAccessList, ipfs, storage)
	if err != nil {
		return nil, errors.Wrap(err, "could not create cap for NewFile")
	}
	var pendingChanges *FileCap
	if err := deepcopy.Copy(&pendingChanges, cap); err != nil {
		return nil, errors.Wrap(err, "could not deep copy cap")
	}

	idString := base64.URLEncoding.EncodeToString(cap.Id[:])
	capPath := storage.GroupFileCapDir(group.Id().ToString()) + idString
	pendingPath := storage.GroupFilePendingDir(group.Id().ToString()) + idString

	file := &File{
		Cap:      cap,
		PendingChanges: pendingChanges,
		DataPath: filePath,
		CapPath: capPath,
		PendingPath: pendingPath,
	}

	if err := file.SaveMetadata(); err != nil {
		return nil, errors.Wrap(err, "could not save file meta data")
	}

	return file, nil
}

func NewGroupFileFromCap(cap *FileCap, groupId string, storage *Storage) (*File, error) {
	idString := base64.URLEncoding.EncodeToString(cap.Id[:])
	capPath := storage.GroupFileCapDir(groupId) + idString
	dataPath := storage.GroupFileDataDir(groupId) + cap.FileName
	pendingPath := storage.GroupFilePendingDir(groupId) + idString

	var pendingChanges *FileCap
	if err := deepcopy.Copy(&pendingChanges, cap); err != nil {
		return nil, errors.Wrap(err, "could not deep copy cap")
	}

	file := &File{
		Cap:      cap,
		PendingChanges: pendingChanges,
		DataPath: dataPath,
		CapPath:   capPath,
		PendingPath: pendingPath,
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
func NewFileFromCap(dataDir, capDir string, cap *FileCap, ipfs ipfsapi.IIpfs, storage *Storage) (*File, error) {
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

func (f *File) Update(cap *FileCap, storage *Storage, ipfs ipfsapi.IIpfs) error {
	oldIpfsHash := f.Cap.IpfsHash
	f.Cap = cap
	if strings.Compare(oldIpfsHash, cap.IpfsHash) != 0 {
		if err := f.SaveMetadata(); err != nil {
			return errors.Wrap(err, "could not save file meta data")
		}

		go f.Download(storage, ipfs)
	}
	return nil
}

// Downloads, decrypts and verifies the content of file from Ipfs
func (f *File) Download(storage *Storage, ipfs ipfsapi.IIpfs) {
	tmpFilePath, err := storage.DownloadTmpFile(f.Cap.IpfsHash, ipfs)
	if err != nil {
		glog.Errorf("could not ipfs get '%s': File.download: %s", f.Cap.IpfsHash, err)
	}

	encReader, err := os.Open(tmpFilePath)
	if err != nil {
		glog.Errorf("could not read file '%s': File.download: %s", tmpFilePath, err)
	}
	defer func() {
		if err := encReader.Close(); err != nil {
			glog.Warningf("could not close tmp file '%s': %s", tmpFilePath, err)
		}
		if err := os.Remove(tmpFilePath); err != nil {
			glog.Warningf("could not delete tmp file '%s': %s", tmpFilePath, err)
		}
	}()

	fout, err := os.Create(f.DataPath)
	if err != nil {
		glog.Errorf("could not create file '%s': File.download: %s", f.DataPath, err)
	}

	defer func() {
		if err := fout.Close(); err != nil {
			glog.Warningf("could not close file '%s': %s", f.DataPath, err)
		}
	}()

	if err := f.Cap.DataKey.Open(encReader, fout); err != nil {
		glog.Errorf("could not dycrypt file '%s': File.download: %s", f.DataPath, err)
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


func (f *File) Encrypt() (io.Reader, error) {
	data, err := ioutil.ReadFile(f.DataPath)
	if err != nil {
		return nil, errors.Wrap(err, "could not read file")
	}
	encData, err := f.Cap.DataKey.Seal(bytes.NewReader(data))
	if err != nil {
		return nil, errors.Wrap(err, "could not encrypt data")
	}

	return encData, nil
}

func GetCapListFromFileList(files []*File) []*FileCap {
	var l []*FileCap
	for _, file := range files {
		l = append(l, file.Cap)
	}
	return l
}

func (f *File) ResetPendingChanges() error {
	var pendingChanges *FileCap
	if err := deepcopy.Copy(&pendingChanges, f.Cap); err != nil {
		return errors.Wrap(err, "could not deep copy file cap")
	}
	f.PendingChanges = pendingChanges

	return nil
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
