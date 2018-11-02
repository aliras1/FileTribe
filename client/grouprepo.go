package client

import (
	. "ipfs-share/collections"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"github.com/golang/glog"
	"bytes"
	"github.com/getlantern/deepcopy"
	"io"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"sync"
	"ipfs-share/ipfs"
)

type IpfsAddOperation func(reader io.Reader) (string, error)

type GroupRepo struct {
	files *ConcurrentCollection
	group IGroup
	ipfs ipfs.IIpfs
	storage *Storage
	user ethcommon.Address

	ipfsHash string

	lock sync.RWMutex
}

func NewGroupRepo(group IGroup, user ethcommon.Address, storage *Storage, ipfs ipfs.IIpfs) (*GroupRepo, error) {
	storage.MakeGroupDir(group.Id().ToString())

	var caps []*FileCap

	data, err := EncodeFileCapList(caps)
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

func NewGroupRepoFromIpfs(ipfsHash string, group IGroup, user ethcommon.Address, storage *Storage, ipfs ipfs.IIpfs) (*GroupRepo, error) {
	repo := &GroupRepo{
		ipfsHash: ipfsHash,
		group: group,
		storage: storage,
		ipfs: ipfs,
		user: user,
	}

	caps, err := repo.getGroupFileCapsFromIpfs(ipfsHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not get group file caps")
	}

	var files *ConcurrentCollection
	for _, cap := range caps {
		file, err := NewGroupFileFromCap(cap, group.Id().ToString(), storage)
		if err != nil {
			return nil, errors.Wrap(err, "could not create new file from cap")
		}
		files.Append(file)
	}
	repo.files = files

	return repo, nil
}


func (repo *GroupRepo) GetPendingChanges() ([]*FileCap, error) {
	dir := repo.storage.GroupFileDataDir(repo.group.Id().ToString())
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrap(err, "could not open group file data dir")
	}
	var listPendingChanges []*FileCap

	for _, f := range files {
		filePath := dir + f.Name()
		var file *File

		fileInt := repo.files.Get(NewStringId(f.Name()))

		// if current file is not in repo --> create new
		if fileInt == nil {
			file, err = NewGroupFile(filePath, []ethcommon.Address{repo.user}, repo.group, repo.storage, repo.ipfs)
			if err != nil {
				return nil, errors.Wrap(err, "could not create new group file")
			}
			repo.files.Append(file)
		} else {
			file = fileInt.(*File)
		}

		encData, err := file.Encrypt()
		if err != nil {
			return nil, errors.Wrap(err, "could not encrypt file data")
		}
		newIpfsHash, err := repo.ipfs.Add(encData)
		if err != nil {
			return nil, errors.Wrap(err, "could not ipfs add file")
		}
		file.PendingChanges.IpfsHash = newIpfsHash

		if err := file.SaveMetadata(); err != nil {
			return nil, errors.Wrap(err, "could not save pending meta data")
		}

		listPendingChanges = append(listPendingChanges, file.PendingChanges)
	}

	return listPendingChanges, nil
}

func (repo *GroupRepo) CommitChanges() (string, error) {
	pendingChanges, err := repo.GetPendingChanges()
	if err != nil {
		return "", errors.Wrap(err, "could not get pending changes")
	}

	data, err := EncodeFileCapList(pendingChanges)
	if err != nil {
		return "", errors.Wrap(err, "could not encode file cap list")
	}

	boxer := repo.group.Boxer()
	encData := boxer.BoxSeal(data)
	newIpfsHash, err := repo.ipfs.Add(bytes.NewReader(encData))
	if err != nil {
		return "", errors.Wrap(err, "could not add new file caps to ipfs")
	}

	return newIpfsHash, nil
}


func (repo *GroupRepo) getFileCaps() []*FileCap {
	var caps []*FileCap
	for fileInterface := range repo.files.Iterator() {
		file := fileInterface.(*File)
		var capCopy *FileCap
		deepcopy.Copy(capCopy, file.Cap)
		caps = append(caps,  capCopy)
	}
	return caps
}

func (repo *GroupRepo) IsValidChangeSet(newIpfsHash string, address *ethcommon.Address) error {
	newCaps, err := repo.getGroupFileCapsFromIpfs(newIpfsHash)
	if err != nil {
		return errors.Wrap(err, "could not get requested group changes")
	}

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
	}

	return nil
}

func (repo *GroupRepo) Update(newIpfsHash string) error {
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
				return errors.Wrap(err, "could not create enw group file from cap")
			}
			repo.files.Append(file)
			go file.Download(repo.storage, repo.ipfs)
		} else {
			file = fileInterface.(*File)
			if err := file.Update(cap, repo.storage, repo.ipfs); err != nil {
				return errors.Wrap(err, "could not update group file")
			}
		}

		return nil
	}

	repo.ipfsHash = newIpfsHash
	if err := repo.ResetPendingChanges(); err != nil {
		return errors.Wrap(err, "could not reset repo pending changes")
	}

	return nil
}

func (repo *GroupRepo) ResetPendingChanges() error {
	for fileInt := range repo.files.Iterator() {
		if err := fileInt.(*File).ResetPendingChanges(); err != nil {
			return errors.Wrap(err, "could not reset file pending changes")
		}
	}
	return nil
}

func (repo *GroupRepo) getGroupFileCapsFromIpfs(ipfsHash string) ([]*FileCap, error) {
	path := repo.storage.tmpPath + ipfsHash
	if err := repo.ipfs.Get(ipfsHash, path); err != nil {
		return nil, errors.Wrapf(err, "could not ipfs get shared group dir: %s", ipfsHash)
	}

	encData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "could not read shared group dir: %s", path)
	}
	defer func() {
		if err := os.Remove(path); err != nil {
			glog.Warningf("could not remove tmp file %s", path)
		}
	}()

	boxer := repo.group.Boxer()
	data, ok := boxer.BoxOpen(encData)
	if !ok {
		return nil, errors.New("could not decrypt shared group dir")
	}

	caps, err := DecodeFileCapList(data)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode file cap list")
	}

	return caps, nil
}

func (repo *GroupRepo) List() []string {
	var l []string
	for fileInt := range repo.files.Iterator() {
		l = append(l, fileInt.(*File).Id().ToString())
	}
	return l
}
