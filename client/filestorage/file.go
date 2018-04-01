package filestorage

import (
	"crypto/rand"
	"encoding/json"
	"io/ioutil"
	"path"

	"golang.org/x/crypto/ed25519"

	"ipfs-share/crypto"
	"ipfs-share/ipfs"
	nw "ipfs-share/network"
)

type File struct {
	// TODO  don't publish file on IPNS, refresh every time
	// TODO  the capability instead
	Path       string                  `json:"path"`
	IPNSPath   string                  `json:"ipnsPath"`
	IPFSAddr   string                  `json:"ipfs_addr"`
	Owner      string                  `json:"owner"`
	SharedWith []string                `json:"shared_with"`
	WAccess    []string                `json:"w_access"`
	VerifyKey  crypto.PublicSigningKey `json:"verify_key"`
	WriteKey   crypto.SecretSigningKey `json:"write_key"`
}

// New File object from local data, found in /userdata/shared
func NewFileFromShared(filePath string) (*File, error) {
	bytesFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var file File
	err = json.Unmarshal(bytesFile, &file)
	return &file, err
}

// New File object from local data, found in /userdata/caps
func NewFileFromCAP(cap *ReadCAP, storage *Storage, ipfs *ipfs.IPFS) (*File, error) {
	filePath := storage.storagePath + "/" + cap.FileName
	file := File{path.Clean(filePath), cap.IPNSPath, cap.IPFSHash, cap.Owner, []string{}, []string{}, cap.VerifyKey, crypto.SecretSigningKey{}}
	if !storage.ExistsFile(storage.storagePath + "/" + cap.FileName) {
		err := storage.DownloadFileFromCap(&file, cap, ipfs)
		if err != nil {
			return nil, err
		}
	}
	return &file, nil
}

// Create a new shared file object from a local file
func NewSharedFile(filePath, owner string, storage *Storage, ipfs *ipfs.IPFS) (*File, error) {
	vk, wk, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	verifyKey := crypto.PublicSigningKey(vk)
	writeKey := crypto.SecretSigningKey(wk)

	ipnsPath, ipfsAddr, err := storage.SignAndAddFileToIPFS(filePath, writeKey, ipfs)
	if err != nil {
		return nil, err
	}
	return &File{filePath, ipnsPath, ipfsAddr, owner, []string{}, []string{}, verifyKey, writeKey}, nil
}

// Share file with a set of users, described by shareWith. Encrypted
// capabilities are made and copied in the corresponding 'public/for/'
// directories. The 'public' directory is re-published into IPNS. After
// that, notification messages are sent out.
func (f *File) Share(shareWith []string, boxer *crypto.BoxingKeyPair, us *Storage, network *nw.Network, ipfs *ipfs.IPFS) error {
	var newUsers []string
	for _, user := range shareWith {
		// add to share list
		f.SharedWith = append(f.SharedWith, user)
		// make new capability into for_X directory
		err := us.CreateFileReadCAPForUser(f, us.publicForPath+"/"+user, f.IPFSAddr, boxer, network)
		newUsers = append(newUsers, user)
		if err != nil {
			return err
		}
	}
	us.StoreFileMetaData(f)
	err := us.PublishPublicDir(ipfs)
	if err != nil {
		return err
	}
	// send share messages
	for _, user := range newUsers {
		err = network.SendMessage(f.Owner, user, "GROUP INVITE", path.Base(f.Path)+".json")
		if err != nil {
			return err
		}
	}
	return nil
}
