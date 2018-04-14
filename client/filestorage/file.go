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

type File interface {
	Share()
	//NewFileFromCAP(cap CAP, )
}

type FilePTP struct {
	// TODO  don't publish file on IPNS, refresh every time
	// TODO  the capability instead
	Name       string                  `json:"name"`
	Owner      string                  `json:"owner"`
	Path       string                  `json:"path"`
	SharedWith []string                `json:"shared_with"`
	WAccess    []string                `json:"w_access"`
	VerifyKey  crypto.PublicSigningKey `json:"verify_key"`
	WriteKey   crypto.SecretSigningKey `json:"write_key"`
}

// New FilePTP object from local data, found in /userdata/shared
func NewFileFromShared(filePath string) (*FilePTP, error) {
	bytesFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var file FilePTP
	err = json.Unmarshal(bytesFile, &file)
	return &file, err
}

// New FilePTP object from local data, found in /userdata/caps
func NewFileFromCAP(cap *ReadCAP, storage *Storage, ipfs *ipfs.IPFS) (*FilePTP, error) {
	filePath := storage.fileRootPath + "/" + cap.Owner + "/" + cap.FileName
	file := FilePTP{cap.FileName, cap.Owner, filePath, []string{}, []string{}, cap.VerifyKey, crypto.SecretSigningKey{}}
	if !storage.ExistsFile(storage.fileRootPath + "/" + cap.FileName) {
		err := storage.DownloadFileFromCap(cap, ipfs)
		if err != nil {
			return nil, err
		}
	}
	return &file, nil
}

// Create a new shared file object from a local file
func NewSharedFile(filePath, owner string, storage *Storage, ipfs *ipfs.IPFS) (*FilePTP, error) {
	newFilePath, err := storage.CopyFileIntoMyFiles(filePath)
	if err != nil {
		return nil, err
	}

	vk, wk, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	verifyKey := crypto.PublicSigningKey(vk)
	writeKey := crypto.SecretSigningKey(wk)

	return &FilePTP{path.Base(filePath), owner, newFilePath, []string{}, []string{}, verifyKey, writeKey}, nil
}

// Share file with a set of users, described by shareWith. Encrypted
// capabilities are made and copied in the corresponding 'public/for/'
// directories. The 'public' directory is re-published into IPNS. After
// that, notification messages are sent out.
func (f *FilePTP) Share(shareWith []string, boxer *crypto.BoxingKeyPair, storage *Storage, network *nw.Network, ipfs *ipfs.IPFS) error {
	var newUsers []string
	for _, user := range shareWith {
		// add to share list
		f.SharedWith = append(f.SharedWith, user)
		ipfsAddr, err := storage.SignAndAddFileToIPFS(f.Path, f.WriteKey, ipfs)
		if err != nil {
			return err
		}
		// make new capability into for_X directory
		if err := storage.CreateFileReadCAPForUser(f, user, ipfsAddr, boxer, network); err != nil {
			return err
		}
		newUsers = append(newUsers, user)
	}
	storage.StoreFileMetaData(f)
	err := storage.PublishPublicDir(ipfs)
	if err != nil {
		return err
	}
	// send share messages
	for _, user := range newUsers {
		err = network.SendMessage(f.Owner, user, "PTP READ CAP", path.Base(f.Name)+".json")
		if err != nil {
			return err
		}
	}
	return nil
}
