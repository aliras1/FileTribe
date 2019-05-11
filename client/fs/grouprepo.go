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
	"io/ioutil"
	"strings"
	"sync"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/glog"
	"github.com/pkg/errors"

	"github.com/aliras1/FileTribe/client/fs/meta"
	"github.com/aliras1/FileTribe/client/interfaces"
	. "github.com/aliras1/FileTribe/collections"
	"github.com/aliras1/FileTribe/ipfs"
	"github.com/aliras1/FileTribe/tribecrypto"
)

// GroupRepo is responsible for managing and maintaining a group's file repository
type GroupRepo struct {
	files   *Map
	group   interfaces.Group
	ipfs    ipfs.IIpfs
	storage *Storage
	user    ethcommon.Address

	ipfsHash string

	lock sync.RWMutex
}

// NewGroupRepo creates a new GroupRepo
func NewGroupRepo(group interfaces.Group, user ethcommon.Address, storage *Storage, ipfs ipfs.IIpfs) (*GroupRepo, error) {
	storage.MakeGroupDir(group.Name(), group.Address().String())

	metas, err := storage.GetGroupFileMetas(group.Address().String())
	if err != nil {
		glog.Warningf("could not load group file metas: %s", err)
	}

	data, err := meta.EncodeFileMetaList(metas)
	if err != nil {
		return nil, errors.Wrap(err, "could not encode empty file meta list")
	}

	boxer := group.Boxer()
	encData := boxer.BoxSeal(data)
	ipfsHash, err := ipfs.Add(bytes.NewReader(encData))
	if err != nil {
		return nil, errors.Wrap(err, "could not add empty list to ipfs")
	}

	return &GroupRepo{
		group:    group,
		files:    NewConcurrentMap(),
		ipfsHash: ipfsHash,
		storage:  storage,
		ipfs:     ipfs,
		user:     user,
	}, nil
}

// IpfsHash returns the current IPFS hash of the group repository
func (repo *GroupRepo) IpfsHash() string {
	repo.lock.RLock()
	defer repo.lock.RUnlock()

	return repo.ipfsHash
}

// Get retrieves a file from the repo
func (repo *GroupRepo) Get(fileName string) *File {
	repo.lock.RLock()
	defer repo.lock.RUnlock()

	fileInt := repo.files.Get(fileName)
	if fileInt == nil {
		return nil
	}

	return fileInt.(*File)
}

// Files returns a list of the repo's files
func (repo *GroupRepo) Files() []*File {
	repo.lock.RLock()
	defer repo.lock.RUnlock()

	var files []*File

	for fileInt := range repo.files.VIterator() {
		files = append(files, fileInt.(*File))
	}

	return files
}

func (repo *GroupRepo) getPendingChanges() ([]*meta.FileMeta, error) {
	dir := repo.storage.GroupFileDataDir(repo.group.Name())
	filesInLocalDir, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrap(err, "could not open group file data dir")
	}
	var listPendingChanges []*meta.FileMeta

	for _, f := range filesInLocalDir {
		filePath := dir + f.Name()
		glog.Infof("file path: %s", filePath)
		var file *File

		fileInt := repo.files.Get(f.Name())

		// if current file is not in repo --> create new
		if fileInt == nil {
			file, err = NewGroupFile(filePath, []ethcommon.Address{repo.user}, repo.group.Address().String(), repo.storage)
			if err != nil {
				return nil, errors.Wrap(err, "could not create new group file")
			}
			repo.files.Put(file.Meta.FileName, file)
		} else {
			file = fileInt.(*File)
		}

		newIpfsHash, err := file.UploadDiff(repo.ipfs)
		if err != nil {
			return nil, errors.Wrap(err, "could not upload file diff")
		}

		file.PendingChanges.IpfsHash = newIpfsHash

		if err := file.SaveMetadata(); err != nil {
			return nil, errors.Wrap(err, "could not save pending meta data")
		}

		listPendingChanges = append(listPendingChanges, file.PendingChanges)
	}

	return listPendingChanges, nil
}

// CommitChanges encrypts and adds the repo's changes to IPFS
func (repo *GroupRepo) CommitChanges(boxer tribecrypto.SymmetricKey) (string, error) {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	pendingChanges, err := repo.getPendingChanges()
	if err != nil {
		return "", errors.Wrap(err, "could not get pending changes")
	}

	data, err := meta.EncodeFileMetaList(pendingChanges)
	if err != nil {
		return "", errors.Wrap(err, "could not encode file meta list")
	}

	encData := boxer.BoxSeal(data)
	newIpfsHash, err := repo.ipfs.Add(bytes.NewReader(encData))
	if err != nil {
		return "", errors.Wrap(err, "could not add new file metas to ipfs")
	}

	return newIpfsHash, nil
}

func (repo *GroupRepo) getFileMetas() []*meta.FileMeta {
	repo.lock.RLock()
	defer repo.lock.RUnlock()

	var fileMetas []*meta.FileMeta
	for fileInterface := range repo.files.VIterator() {
		file := fileInterface.(*File)
		var metaCopy *meta.FileMeta
		deepcopy(metaCopy, file.Meta)
		fileMetas = append(fileMetas, metaCopy)
	}

	return fileMetas
}

// IsValidChangeSet verifies if a proposed change set is valid or not
func (repo *GroupRepo) IsValidChangeSet(newIpfsHash string, boxer tribecrypto.SymmetricKey, address ethcommon.Address) error {
	repo.lock.RLock()
	defer repo.lock.RUnlock()

	newMetas, err := repo.getGroupFileMetasFromIpfs(newIpfsHash, boxer)
	if err != nil {
		return errors.Wrap(err, "could not get requested group changes")
	}

	// TODO: handle delete
	for _, newMeta := range newMetas {
		if newMeta.WriteAccessList == nil {
			return errors.New("new write access list can not be nil")
		}

		fileInt := repo.files.Get(newMeta.FileName)
		if fileInt == nil {
			// new file, nothing to check
			continue
		}

		file := fileInt.(*File)
		if strings.Compare(file.Meta.IpfsHash, newMeta.IpfsHash) == 0 {
			// no changes
			continue
		}

		// check if user has write access to the current file
		hasWriteAccess := false
		for _, hasW := range file.Meta.WriteAccessList {
			if bytes.Equal(hasW.Bytes(), address.Bytes()) {
				hasWriteAccess = true
				break
			}
		}

		if !hasWriteAccess {
			return errors.New("member has no write access")
		}

		// check if new DiffNode is correct
		if err := repo.isDiffNodeValid(file, newMeta.DataKey, newMeta.IpfsHash); err != nil {
			return errors.Wrap(err, "invalid new DiffNode")
		}
	}

	return nil
}

func (repo *GroupRepo) isDiffNodeValid(file *File, newBoxer tribecrypto.FileBoxer, newIpfsHash string) error {
	repo.lock.RLock()
	defer repo.lock.RUnlock()

	data, err := repo.storage.DownloadAndDecryptWithFileBoxer(newBoxer, newIpfsHash, repo.ipfs)
	if err != nil {
		return errors.Wrap(err, "could not download and decrypt new diff node")
	}

	newDiff, err := DecodeDiffNode(data)
	if err != nil {
		return errors.Wrap(err, "could not decode new DiffNode")
	}

	if strings.Compare(newDiff.Next, file.Meta.IpfsHash) != 0 {
		return errors.New("next ipfs hash is not the current ipfs hash")
	}

	fileData, err := ioutil.ReadFile(file.OrigPath)
	if err != nil {
		return errors.Wrap(err, "could not read orig file")
	}

	hash := ethcrypto.Keccak256(fileData)

	if !bytes.Equal(newDiff.Hash, hash) {
		return errors.New("new diff prev hash does not match with current hash")
	}

	return nil
}

// Update updates the group repository according to the new IPFS hash
func (repo *GroupRepo) Update(newIpfsHash string) error {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	if strings.EqualFold(newIpfsHash, "") || strings.EqualFold(newIpfsHash, repo.ipfsHash) {
		return nil
	}

	fileMetas, err := repo.getGroupFileMetasFromIpfs(newIpfsHash, repo.group.Boxer())
	if err != nil {
		return errors.Wrap(err, "could not get group file fileMetas from ipfs")
	}

	for _, fileMeta := range fileMetas {
		var file *File
		var err error
		fileInterface := repo.files.Get(fileMeta.FileName)
		if fileInterface == nil {
			file, err = NewGroupFileFromMeta(fileMeta, repo.group.Address().String(), repo.group.Name(), repo.storage)
			if err != nil {
				return errors.Wrap(err, "could not create new group file from fileMeta")
			}

			repo.files.Put(file.Meta.FileName, file)
			go file.Download(repo.storage, repo.ipfs)
		} else {
			file = fileInterface.(*File)
			if err := file.Update(fileMeta, repo.storage, repo.ipfs); err != nil {
				return errors.Wrap(err, "could not Update group file")
			}
		}

		if err := file.SaveMetadata(); err != nil {
			return errors.Wrap(err, "could not save file meta data on disk")
		}
	}

	repo.ipfsHash = newIpfsHash

	return nil
}

func (repo *GroupRepo) getGroupFileMetasFromIpfs(ipfsHash string, boxer tribecrypto.SymmetricKey) ([]*meta.FileMeta, error) {
	data, err := repo.storage.DownloadAndDecryptWithSymmetricKey(boxer, ipfsHash, repo.ipfs)
	if err != nil {
		return nil, errors.Wrap(err, "could not download group data")
	}

	fileMetas, err := meta.DecodeFileMetaList(data)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode file meta list")
	}

	return fileMetas, nil
}
