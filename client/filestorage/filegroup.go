package filestorage

import (
	"fmt"
	"io/ioutil"
	"path"

	"ipfs-share/crypto"
	"ipfs-share/ipfs"
)

type FileGroup struct {
	Name string
	//Path string
	IPFSHash string
}

func NewSharedFileGroup(filePath, groupName string, dataKey crypto.SymmetricKey, storage *Storage, ipfs *ipfs.IPFS) (*FileGroup, error) {
	fileName := path.Base(filePath)

	// First just encrypt and add to ipfs, then call file.Share()
	// which makes the share transaction, after which everyone,
	// including the sharer, ipfs.get the new file, decrypts it
	// and saves it in the storage.GroupFiles...
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not read file '%s': NewSharedFileGroup: %s", filePath, err)
	}
	encFileBytes := dataKey.BoxSeal(fileBytes)
	ipfsFilePath := storage.ipfsFilesPath + "/" + fileName
	if err := WriteFile(ipfsFilePath, encFileBytes); err != nil {
		return nil, fmt.Errorf("could not copy file into ipfs file: NewSharedFileGroup: %s", err)
	}
	merkleNode, err := ipfs.AddFile(ipfsFilePath)
	if err != nil {
		return nil, fmt.Errorf("could not ipfs add file '%s': NewSharedFileGroup: %s", ipfsFilePath, err)
	}
	ipfsHash := "/ipfs/" + merkleNode.Hash

	file := &FileGroup{
		Name: fileName,
		//Path: newFilePath,
		IPFSHash: ipfsHash,
	}
	// TODO:  it might be a good idea to store this FileGroup in
	// TODO:  public/for/group/files/ to be enable others to pull
	// TODO:  all group data in one step
	return file, nil
}

