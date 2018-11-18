package fs

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"
	"sync"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/getlantern/deepcopy"
	"github.com/pkg/errors"

	"ipfs-share/client/fs/caps"
	"ipfs-share/client/interfaces"
	. "ipfs-share/collections"
	"ipfs-share/crypto"
	"ipfs-share/ipfs"
)

type IpfsAddOperation func(reader io.Reader) (string, error)

type GroupRepo struct {
	files *ConcurrentCollection
	group interfaces.IGroup
	ipfs ipfs.IIpfs
	storage *Storage
	user ethcommon.Address

	ipfsHash string

	lock sync.RWMutex
}

func NewGroupRepo(group interfaces.IGroup, user ethcommon.Address, storage *Storage, ipfs ipfs.IIpfs) (*GroupRepo, error) {
	storage.MakeGroupDir(group.Id().ToString())

	var capabilities []*caps.FileCap

	data, err := caps.EncodeFileCapList(capabilities)
	if err != nil {
		return nil, errors.Wrap(err, "could not encode empty cap list")
	}

	boxer := group.Boxer()
	encData := boxer.BoxSeal(data)
	ipfsHash, err := ipfs.Add(bytes.NewReader(encData))
	if err != nil {
		return nil, errors.Wrap(err, "could not add empty list to ipfs")
	}

	return &GroupRepo{
		group: group,
		files: NewConcurrentCollection(),
		ipfsHash: ipfsHash,
		storage: storage,
		ipfs: ipfs,
		user: user,
	}, nil
}

func NewGroupRepoFromIpfs(ipfsHash string, group interfaces.IGroup, user ethcommon.Address, storage *Storage, ipfs ipfs.IIpfs) (*GroupRepo, error) {
	repo := &GroupRepo{
		ipfsHash: ipfsHash,
		group: group,
		storage: storage,
		ipfs: ipfs,
		user: user,
	}

	capabilities, err := repo.getGroupFileCapsFromIpfs(ipfsHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not get group file caps")
	}

	var files *ConcurrentCollection
	for _, cap := range capabilities {
		file, err := NewGroupFileFromCap(cap, group.Id().ToString(), storage)
		if err != nil {
			return nil, errors.Wrap(err, "could not create new file from cap")
		}
		files.Append(file)
	}
	repo.files = files

	return repo, nil
}

func (repo *GroupRepo) IpfsHash() string {
	repo.lock.RLock()
	defer repo.lock.RUnlock()

	return repo.ipfsHash
}

func (repo *GroupRepo) Get(id IIdentifier) *File {
	repo.lock.RLock()
	defer repo.lock.RUnlock()

	fileInt := repo.files.Get(id)
	if fileInt == nil {
		return nil
	}

	return fileInt.(*File)
}

func (repo *GroupRepo) Files() []*File {
	repo.lock.RLock()
	defer repo.lock.RUnlock()

	var files []*File

	for fileInt := range repo.files.Iterator() {
		files = append(files, fileInt.(*File))
	}

	return files
}

func (repo *GroupRepo) getPendingChanges() ([]*caps.FileCap, error) {
	dir := repo.storage.GroupFileDataDir(repo.group.Id().ToString())
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrap(err, "could not open group file data dir")
	}
	var listPendingChanges []*caps.FileCap

	for _, f := range files {
		filePath := dir + f.Name()
		var file *File

		fileInt := repo.files.Get(NewStringId(f.Name()))

		// if current file is not in repo --> create new
		if fileInt == nil {
			file, err = NewGroupFile(filePath, []ethcommon.Address{repo.user}, repo.group.Id().ToString(), repo.storage)
			if err != nil {
				return nil, errors.Wrap(err, "could not create new group file")
			}
			repo.files.Append(file)
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

func (repo *GroupRepo) CommitChanges(boxer crypto.SymmetricKey) (string, error) {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	pendingChanges, err := repo.getPendingChanges()
	if err != nil {
		return "", errors.Wrap(err, "could not get pending changes")
	}

	data, err := caps.EncodeFileCapList(pendingChanges)
	if err != nil {
		return "", errors.Wrap(err, "could not encode file cap list")
	}

	encData := boxer.BoxSeal(data)
	newIpfsHash, err := repo.ipfs.Add(bytes.NewReader(encData))
	if err != nil {
		return "", errors.Wrap(err, "could not add new file caps to ipfs")
	}

	return newIpfsHash, nil
}


func (repo *GroupRepo) getFileCaps() []*caps.FileCap {
	repo.lock.RLock()
	defer repo.lock.RUnlock()

	var capabilities []*caps.FileCap
	for fileInterface := range repo.files.Iterator() {
		file := fileInterface.(*File)
		var capCopy *caps.FileCap
		deepcopy.Copy(capCopy, file.Cap)
		capabilities = append(capabilities,  capCopy)
	}

	return capabilities
}

func (repo *GroupRepo) IsValidChangeSet(newIpfsHash string, address *ethcommon.Address) error {
	repo.lock.RLock()
	defer repo.lock.RUnlock()

	newCaps, err := repo.getGroupFileCapsFromIpfs(newIpfsHash)
	if err != nil {
		return errors.Wrap(err, "could not get requested group changes")
	}

	// TODO: handle delete
	for _, newCap := range newCaps {
		if newCap.WriteAccessList == nil {
			return errors.New("new write access list can not be nil")
		}

		fileInt := repo.files.Get(NewStringId(newCap.FileName))
		if fileInt == nil {
			// new file, nothing to check
			continue
		}

		file := fileInt.(*File)
		if strings.Compare(file.Cap.IpfsHash, newCap.IpfsHash) == 0 {
			// no changes
			continue
		}

		// check if user has write access to the current file
		hasWriteAccess := false
		for _, hasW := range file.Cap.WriteAccessList {
			if bytes.Equal(hasW.Bytes(), address.Bytes()) {
				hasWriteAccess = true
				break
			}
		}

		if !hasWriteAccess {
			return errors.New("member has no write access")
		}

		// check if new DiffNode is correct
		if err := repo.checkDiffNode(file, newCap.DataKey, newCap.IpfsHash); err != nil {
			return errors.Wrap(err, "invalid new DiffNode")
		}
	}

	return nil
}

func (repo *GroupRepo) checkDiffNode(file *File, newBoxer crypto.FileBoxer, newIpfsHash string) error {
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


	if strings.Compare(newDiff.Next, file.Cap.IpfsHash) != 0 {
		return errors.New("next ipfs hash is not the current ipfs hash")
	}

	fileData, err := ioutil.ReadFile(file.OrigPath)
	if err != nil {
		return errors.Wrap(err, "could not read orig file")
	}

	hasher := crypto.NewKeccak256Hasher()
	hash := hasher.Sum(fileData)

	if !bytes.Equal(newDiff.Hash, hash) {
		return errors.New("new diff prev hash does not match with current hash")
	}

	return nil
}

func (repo *GroupRepo) Update(newIpfsHash string) error {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	caps, err := repo.getGroupFileCapsFromIpfs(newIpfsHash)
	if err != nil {
		return errors.Wrap(err, "could not get group file caps from ipfs")
	}

	for _, cap := range caps {
		var file *File
		fileInterface := repo.files.Get(NewStringId(cap.FileName))
		if fileInterface == nil {
			file, err := NewGroupFileFromCap(cap, repo.group.Id().ToString(), repo.storage)
			if err != nil {
				return errors.Wrap(err, "could not create new group file from cap")
			}
			repo.files.Append(file)
			go file.Download(repo.storage, repo.ipfs)
		} else {
			file = fileInterface.(*File)
			if err := file.Update(cap, repo.storage, repo.ipfs); err != nil {
				return errors.Wrap(err, "could not Update group file")
			}
		}

		return nil
	}

	repo.ipfsHash = newIpfsHash

	return nil
}

func (repo *GroupRepo) getGroupFileCapsFromIpfs(ipfsHash string) ([]*caps.FileCap, error) {
	data, err := repo.storage.DownloadAndDecryptWithSymmetricKey(repo.group.Boxer(), ipfsHash, repo.ipfs)
	if err != nil {
		return nil, errors.Wrap(err, "could not download group data")
	}

	capabilities, err := caps.DecodeFileCapList(data)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode file cap list")
	}

	return capabilities, nil
}

func (repo *GroupRepo) ReEncrypt(boxer crypto.SymmetricKey) (string, error) {
	for fileInt := range repo.files.Iterator() {
		file := fileInt.(*File)
		if err := file.ChangeKey(repo.ipfs); err != nil {
			return "", errors.Wrap(err, "could not change file key")
		}
	}

	ipfsHash, err := repo.CommitChanges(boxer)
	if err != nil {
		return "", errors.Wrap(err, "could not commit changes")
	}

	return ipfsHash, nil
}
