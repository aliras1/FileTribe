package client

import (
	. "ipfs-share/collections"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"github.com/golang/glog"
	"bytes"
	"strings"
)

type GroupRepo struct {
	files *ConcurrentCollection
	groupCtx *GroupContext
	ipfsHash string
	pendingChanges map[string][]*File // ipfs hash -> change set
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
		pendingChanges: make(map[string][]*File),
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

// Creates a change set of the new file and returns the ipfs
// hash of the capabilities as if the changes were made for
// real
func (repo *GroupRepo) QueueAddFile(file *File) (string, error) {
	caps := repo.getFileCaps()
	caps = append(caps, file.Cap)
	data, err := EncodeFileCapList(caps)
	if err != nil {
		return "", errors.Wrap(err, "could not encode capability list")
	}
	encData := repo.groupCtx.Group.Boxer.BoxSeal(data)

	hash, err := repo.groupCtx.Ipfs.Add(bytes.NewReader(encData))
	if err != nil {
		return "", errors.Wrap(err, "could not add new group caps list to ipfs")
	}

	repo.pendingChanges[hash] = []*File{file}

	return hash, nil
}

func (repo *GroupRepo) getFileCaps() []*FileCap {
	var caps []*FileCap
	for fileInterface := range repo.files.Iterator() {
		file := fileInterface.(*File)
		caps = append(caps, file.Cap)
	}
	return caps
}

func (repo *GroupRepo) IsValidAddFile(newIpfsHash string) error {
	// TODO: refresh current ipfs hash before checking

	newCaps, err := repo.getGroupFileCapsFromIpfs(newIpfsHash)
	if err != nil {
		return errors.Wrap(err, "could not get group new file caps")
	}

	if len(newCaps) != repo.files.Count() + 1 {
		return errors.New("invalid number of caps")
	}

	for fileInt := range repo.files.Iterator() {
		file := fileInt.(*File)
		var currentCap *FileCap = nil
		for i := 0; i < len(newCaps); i++ {
			if file.Id().Equal(NewBytesId(newCaps[i].Id)) {
				currentCap = newCaps[i]
				newCaps = append(newCaps[:i], newCaps[i+1: ]...)
			}
		}
		if currentCap == nil {
			return errors.New("did not found an existing cap in new caps")
		}
		if !file.Cap.Equal(currentCap) {
			return errors.New("caps have changed")
		}
	}

	// Cache the new file
	repo.pendingChanges[newIpfsHash] = []*File{NewGroupFileFromCap(newCaps[0], repo.groupCtx)}

	return nil
}

func (repo *GroupRepo) Update(newIpfsHash string) error {
	if strings.Compare(newIpfsHash, repo.ipfsHash) == 0 {
		return nil
	}

	var changes []*File
	var ok bool
	changes, ok = repo.pendingChanges[newIpfsHash]
	if !ok {
		caps, err := repo.getGroupFileCapsFromIpfs(newIpfsHash)
		if err != nil {
			return errors.Wrap(err, "could not get group file caps from ipfs")
		}
		for _, cap := range caps {
			changes = append(changes, NewGroupFileFromCap(cap, repo.groupCtx))
		}
	}

	for _, change := range changes {
		change.Save()
		go change.Download(repo.groupCtx.Storage, repo.groupCtx.Ipfs)
		// Update existing file
		err := repo.files.Update(change); if err == nil {
			continue
		}
		// Add new file
		repo.files.Append(change)
	}

	repo.ipfsHash = newIpfsHash

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
