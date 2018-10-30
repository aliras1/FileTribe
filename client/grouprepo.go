package client

import (
	. "ipfs-share/collections"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"github.com/golang/glog"
	"bytes"
	"strings"
	"github.com/getlantern/deepcopy"
	"io"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"ipfs-share/utils"
)

type IpfsAddOperation func(reader io.Reader) (string, error)

type GroupRepo struct {
	files *ConcurrentCollection
	pendingChanges *ConcurrentCollection
	groupCtx *GroupContext
	ipfsHash string
}

func NewGroupRepo(groupCtx *GroupContext) (*GroupRepo, error) {
	groupCtx.Storage.MakeGroupDir(groupCtx.Group.Id.ToString())

	var caps []*FileCap

	data, err := EncodeFileCapList(caps)
	if err != nil {
		return nil, errors.Wrap(err, "could not encode empty cap list")
	}
	encData := groupCtx.Group.Boxer.BoxSeal(data)
	ipfsHash, err := groupCtx.Ipfs.Add(bytes.NewReader(encData))
	if err != nil {
		return nil, errors.Wrap(err, "could not add empty list to ipfs")
	}

	return &GroupRepo{
		groupCtx: groupCtx,
		files: NewConcurrentCollection(),
		ipfsHash: ipfsHash,
		pendingChanges: NewConcurrentCollection(),
	}, nil
}

func NewGroupRepoFromIpfs(ipfsHash string, groupCtx *GroupContext) (*GroupRepo, error) {
	repo := &GroupRepo{
		groupCtx: groupCtx,
		ipfsHash: ipfsHash,
	}

	caps, err := repo.getGroupFileCapsFromIpfs(ipfsHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not get group file caps")
	}

	var files *ConcurrentCollection
	for _, cap := range caps {
		files.Append(NewGroupFileFromCap(cap, groupCtx))
	}
	repo.files = files

	return repo, nil
}

func (repo *GroupRepo) GetFileFromPath(filePath string) *File {
	for f := range repo.files.Iterator() {
		if strings.Compare(filePath, f.(*File).DataPath) == 0 {
			return f.(*File)
		}
	}
	return nil
}

func (repo *GroupRepo) GetPendingChanges() error {
	dir := repo.groupCtx.Storage.GetGroupFileDataDir(repo.groupCtx.Group.Id.ToString())
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return errors.Wrap(err, "could not open group file data dir")
	}

	for _, f := range files {
		filePath := dir + f.Name()
		var pendingFile *File

		pendingInt := repo.pendingChanges.Get(NewStringId(filePath))
		// file is not in pending changes --> create and add
		if pendingInt == nil {
			pendingFile, err = NewGroupFile(filePath, repo.groupCtx, []ethcommon.Address{})
			if err != nil {
				return errors.Wrap(err, "could not create new group file")
			}
			repo.pendingChanges.Append(pendingFile)
		}

		encData, err := pendingFile.Encrypt()
		if err != nil {
			return errors.Wrap(err, "could not encrypt file data")
		}
		newIpfsHash, err := repo.groupCtx.Ipfs.Add(encData)
		if err != nil {
			return errors.Wrap(err, "could not ipfs add file")
		}
		pendingFile.Cap.IpfsHash = newIpfsHash

		if err := (*PendingFile)(pendingFile).SaveMetadata(); err != nil {
			return errors.Wrap(err, "could not save pending meta data")
		}
	}

	return nil
}

func (repo *GroupRepo) CommitChanges() (string, error) {
	err := repo.GetPendingChanges()
	if err != nil {
		return "", errors.Wrap(err, "could not get pending changes")
	}

	var caps []*FileCap
	for pendingFileInterface := range repo.pendingChanges.Iterator() {
		caps = append(caps, pendingFileInterface.(*File).Cap)
	}

	data, err := EncodeFileCapList(caps)
	if err != nil {
		return "", errors.Wrap(err, "could not encode file cap list")
	}

	encData := repo.groupCtx.Group.Boxer.BoxSeal(data)
	newIpfsHash, err := repo.groupCtx.Ipfs.Add(bytes.NewReader(encData))
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
		fileInterface := repo.files.Get(NewBytesId(newCap.Id))
		if fileInterface == nil {
			// TODO: new file --> check if user can add new files to repo
			continue
		}
		file := fileInterface.(*File)
		if len(file.Cap.WriteAccess) != 0 {
			if _, inArray := utils.InArray(file.Cap.WriteAccess, address); !inArray {
				return errors.New("invalid file modification request: member has no write access")
			}
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
		fileInterface := repo.files.Get(NewBytesId(cap.Id))
		if fileInterface == nil {
			file = NewGroupFileFromCap(cap, repo.groupCtx)
			repo.files.Append(file)
		} else {
			file = fileInterface.(*File)
		}

		// if hashes do not match or this is a new file
		if strings.Compare(file.Cap.IpfsHash, cap.IpfsHash) != 0 || fileInterface == nil {
			if err := file.SaveMetadata(); err != nil {
				return errors.Wrap(err, "could not save file meta data")
			}

			go file.Download(repo.groupCtx.Storage, repo.groupCtx.Ipfs)
		}
	}

	repo.ipfsHash = newIpfsHash
	repo.pendingChanges.Reset()

	return nil
}

func (repo *GroupRepo) getGroupFileCapsFromIpfs(ipfsHash string) ([]*FileCap, error) {
	path := repo.groupCtx.Storage.tmpPath + ipfsHash
	if err := repo.groupCtx.Ipfs.Get(ipfsHash, path); err != nil {
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

	data, ok := repo.groupCtx.Group.Boxer.BoxOpen(encData)
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

func (repo *GroupRepo) GrantWriteAccessToFile(path string, user ethcommon.Address) error {
	var pending *PendingFile

	pendingInt := repo.pendingChanges.Get(NewStringId(path))
	if pendingInt == nil {
		file, err := NewGroupFile(path, repo.groupCtx, []ethcommon.Address{})
		if err != nil {
			return errors.Wrap(err, "could not get pending change")
		}
		pending = (*PendingFile)(file)
	} else {
		pending = pendingInt.(*PendingFile)
	}

	for _, hasW := range pending.Cap.WriteAccess {
		if bytes.Equal(hasW.Bytes(), user.Bytes()) {
			return errors.New("user already has Write access")
		}
	}

	for _, hasW := range pending.Cap.WriteAccess {
		if bytes.Equal(hasW.Bytes(), repo.groupCtx.User.Address.Bytes()) {
			return errors.New("can not grant write access: user has no Write access")
		}
	}

	pending.Cap.WriteAccess = append(pending.Cap.WriteAccess, user)
	if err := pending.SaveMetadata(); err != nil {
		return errors.Wrap(err, "could not save pending meta data")
	}

	return nil
}