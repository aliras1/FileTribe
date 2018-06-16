package filestorage

import (
	"github.com/ethereum/go-ethereum/common"
	// "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	// "os"
	// "path"
	// "strings"

	// "golang.org/x/crypto/ed25519"
	ipfsapi "github.com/ipfs/go-ipfs-api"

	"ipfs-share/crypto"
	nw "ipfs-share/networketh"
)

type File interface {
	Share()
	//NewFileFromCAP(cap CAP, )
}

type FilePTP struct {
	Name       string           `json:"name"`
	Owner      common.Address   `json:"owner"`
	IPFSHash   string           `json:"ipfs_hash"`
	IPNSPath   string           `json:"ipns_path"`
	Path       string           `json:"path"`
	SharedWith []common.Address `json:"shared_with"`
	WAccess    []string         `json:"w_access"`
	VerifyKey  crypto.VerifyKey `json:"verify_key"`
	WriteKey   crypto.Signer    `json:"write_key"`
	Own        bool             `json:"own"` // current user owns the file?
	// it could be a good idea to hardwire Owner into the file data
	// as well and validate it...
}

// New FilePTP object from local data
func NewFile(filePath string) (*FilePTP, error) {
	bytesFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not read file '%s': NewFile: %s", filePath, err)
	}
	var file FilePTP
	if err := json.Unmarshal(bytesFile, &file); err != nil {
		fmt.Errorf("could not unmarshal file '%s': NewFile: %s", filePath, err)
	}
	return &file, nil
}

func NewFileFromCAP(cap *ReadCAP, storage *Storage, ipfs *ipfsapi.Shell) (*FilePTP, error) {
	filePath := storage.fileRootPath + "/" + base64.StdEncoding.EncodeToString(cap.Owner[:]) + "/" + cap.FileName
	// we could directly ipfs.get with this /ipns/address but we need it's
	// /ipfs/hash to be able to check if the file has changed
	ipfsHash, err := ipfs.Resolve(cap.IPNSPath)
	if err != nil {
		return nil, fmt.Errorf("could not resolve ipns address: NewFileFromCAP: %s", err)
	}

	file := FilePTP{
		Name:       cap.FileName,
		Owner:      cap.Owner,
		IPFSHash:   ipfsHash,
		IPNSPath:   cap.IPNSPath,
		Path:       filePath,
		SharedWith: []common.Address{},
		WAccess:    []string{},
		VerifyKey:  cap.VerifyKey,
		WriteKey:   crypto.Signer{},
		Own:        false,
	}

	if err := file.download(storage, ipfs); err != nil {
		return nil, fmt.Errorf("could not download file '%s': NewFileFromCAP: %s", cap.FileName, err)
	}
	if err := file.save(storage); err != nil {
		return nil, fmt.Errorf("could not save file '%s': NewFileFromCAP: %s", file.Name, err)
	}
	return &file, nil
}

// Downloads and verifies the file from IPFS
func (f *FilePTP) download(storage *Storage, ipfs *ipfsapi.Shell) error {
	// tmpFilePath := storage.tmpPath + "/" + path.Base(f.Name)
	// err := ipfs.Get(tmpFilePath, f.IPFSHash)
	// if err != nil {
	// 	return fmt.Errorf("could not ipfs get '%s': FilePTP.download: %s", f.IPFSHash, err)
	// }
	// bytesSignedFile, err := ioutil.ReadFile(tmpFilePath)
	// if err != nil {
	// 	return fmt.Errorf("could not read file '%s': FilePTP.download: %s", tmpFilePath, err)
	// }
	// os.Remove(tmpFilePath)
	// // make a directory to the owner
	// dirPath := storage.fileRootPath + "/" + base64.StdEncoding.EncodeToString(f.Owner[:])
	// os.MkdirAll(dirPath, 0770)

	// bytesRawFile, ok := f.VerifyKey.Verify(bytesSignedFile)
	// if !ok {
	// 	return fmt.Errorf("could not verify file '%s': FilePTP.download: %s", f.Name, err)
	// }
	// filePath := dirPath + "/" + f.Name
	// if err := WriteFile(filePath, bytesRawFile); err != nil {
	// 	return fmt.Errorf("could not write file '%s': FilePTP.download: %s", f.Name, err)
	// }
	return nil
}

func (f *FilePTP) save(storage *Storage) error {
	// path := storage.capsPath + "/" + f.Name
	// jsonBytes, err := json.Marshal(f)
	// if err != nil {
	// 	return fmt.Errorf("could not marshal file '%s': FilePTP.save: %s", f.Name, err)
	// }
	// if err := WriteFile(path, jsonBytes); err != nil {
	// 	return fmt.Errorf("could not write file '%s': FilePTP.save: %s", path, err)
	// }
	return nil
}

// Create a new shared file object from a local file
func NewSharedFilePTP(filePath string, owner common.Address, storage *Storage, ipfs *ipfsapi.Shell) (*FilePTP, error) {
	// newFilePath, err := storage.CopyFileIntoMyFiles(filePath)
	// if err != nil {
	// 	return nil, fmt.Errorf("could not copy file to user storage: NewSharedFilePTP: %s", err)
	// }

	// vk, wk, err := ed25519.GenerateKey(rand.Reader)
	// if err != nil {
	// 	return nil, fmt.Errorf("could not generate signing key pair: NewSharedFilePTP: %s", err)
	// }
	// verifyKey := crypto.VerifyKey(vk)
	// writeKey := crypto.Signer(wk)

	// ipfsID, err := ipfs.ID()
	// if err != nil {
	// 	return nil, fmt.Errorf("could not get ipfs id: NewSharedFilePTP: %s", err)
	// }
	// fileName := path.Base(filePath)
	// file := &FilePTP{
	// 	Name:       fileName,
	// 	Owner:      owner,
	// 	IPFSHash:   "",
	// 	IPNSPath:   "/ipns/" + ipfsID.ID + "/files/" + fileName,
	// 	Path:       newFilePath,
	// 	SharedWith: []common.Address{},
	// 	WAccess:    []string{},
	// 	VerifyKey:  verifyKey,
	// 	WriteKey:   writeKey,
	// 	Own:        true,
	// }
	// if err := file.signAndAddToIPFS(storage, ipfs); err != nil {
	// 	return nil, fmt.Errorf("could not sign and add file '%s' to ipfs: NewSharedFilePTP: %s", fileName, err)
	// }
	// if err := file.save(storage); err != nil {
	// 	return nil, fmt.Errorf("could not save file '%s': NewSharedFilePTP: %s", file.Name, err)
	// }
	return nil, nil
}

// Signs the files with the Write key and then the function adds
// it to IPFS. The function returns with the with the IPFS hash
// of the file
func (f *FilePTP) signAndAddToIPFS(storage *Storage, ipfs *ipfsapi.Shell) error {
	// publicFilePath := storage.publicFilesPath + "/" + f.Name
	// fileBytes, err := ioutil.ReadFile(f.Path)
	// if err != nil {
	// 	return fmt.Errorf("could not read file '%s': FilePTP.signAndAddToIPFS: %s", f.Name, err)
	// }
	// if f.WriteKey == nil {
	// 	return fmt.Errorf("no write key found in file '%s': FilePTP.signAndAddToIPFS", f.Name)
	// }
	// signedFile := f.WriteKey.Sign(fileBytes)
	// if err := WriteFile(publicFilePath, signedFile); err != nil {
	// 	return fmt.Errorf("could not write signed file '%s': FilePTP.signAndAddToIPFS: %s", f.Name, err)
	// }
	// merkleNode, err := ipfs.AddFile(publicFilePath)
	// if err != nil {
	// 	return fmt.Errorf("could not add file '%s' to ipfs: FilePTP.signAndAddToIPFS: %s", f.Name, err)
	// }
	// f.IPFSHash = merkleNode.Hash
	return nil
}

// Share file with a set of users, described by shareWith. Encrypted
// capabilities are made and copied in the corresponding 'public/for/'
// directories. The 'public' directory is re-published into IPNS. After
// that, notification messages are sent out.
func (f *FilePTP) Share(shareWith []common.Address, boxer *crypto.BoxingKeyPair, signer *crypto.Signer, storage *Storage, network *nw.Network, ipfs *ipfsapi.Shell) error {
	// var newUsers []common.Address
	// for _, userID := range shareWith {
	// 	// add to share list
	// 	f.SharedWith = append(f.SharedWith, userID)
	// 	// make new capability into for_X directory
	// 	if err := CreateFileReadCAPForUser(f, userID, f.IPNSPath, storage, network); err != nil {
	// 		return fmt.Errorf(
	// 			"could not create CAP for file '%s' for user '%s': FilePTP.Share: %s",
	// 			f.Name,
	// 			base64.StdEncoding.EncodeToString(userID[:]),
	// 			err,
	// 		)
	// 	}
	// 	// NOTE: we cannot send notification messages here because
	// 	// from efficiency considerations /public directory will be
	// 	// published just once, with all the new CAPs in it
	// 	newUsers = append(newUsers, userID)
	// }
	// if err := f.save(storage); err != nil {
	// 	return fmt.Errorf("could not save file: FilePTP.Share: %s", err)
	// }
	// if err := storage.PublishPublicDir(ipfs); err != nil {
	// 	return fmt.Errorf("could not publish public dir: Share: %s", err)
	// }
	// // send share messages
	// for _, userID := range newUsers {
	// 	_, publicKeyBytes, _, _, err := network.GetUser(userID)
	// 	if err != nil {
	// 		return fmt.Errorf("could not get user data: FilePTP.Share: %s", err)
	// 	}
	// 	publicKey := &crypto.AnonymPublicKey{Value: &publicKeyBytes}

	// 	if err := network.SendMessage(publicKey, signer, f.Owner, "PTP READ CAP", path.Base(f.Name)+".json"); err != nil {
	// 		return fmt.Errorf(
	// 			"could not send 'PTP READ CAP' message to user '%s': Share: %s",
	// 			base64.StdEncoding.EncodeToString(userID[:]),
	// 			err,
	// 		)
	// 	}
	// }
	return nil
}

func (f *FilePTP) Refresh(storage *Storage, ipfs *ipfsapi.Shell) error {
	// if f.Own {
	// 	return nil
	// }
	// newIPFSHash, err := ipfs.Resolve(f.IPNSPath)
	// if err != nil {
	// 	return fmt.Errorf("could not resolve ipns path '%s': FilePTP.Refresh: %s", f.IPNSPath, err)
	// }
	// if strings.Compare(newIPFSHash, f.IPFSHash) != 0 {
	// 	f.IPFSHash = newIPFSHash
	// 	if err := f.download(storage, ipfs); err != nil {
	// 		return fmt.Errorf("could not download file '%s': FilePTP.Refresh: %s", f.Name, err)
	// 	}
	// }
	return nil
}
