package client

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/golang/glog"

	ipfsapi "ipfs-share/ipfs"
	"ipfs-share/collections"
	"github.com/pkg/errors"
	"ipfs-share/utils"
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
	DataPath string
	CapPath string
}

func (f *File) Id() collections.IIdentifier {
	return collections.NewBytesId(f.Cap.Id)
}

func NewGroupFile(srcFilePath string, groupCtx *GroupContext) (*File, error) {
	dataPath := groupCtx.Storage.GetGroupFileDataDir(groupCtx.Group.Id.ToString()) + path.Base(srcFilePath)
	err := groupCtx.Storage.CopyFileIntoGroupFiles(srcFilePath, groupCtx.Group.Id.ToString())
	if err != nil {
		return nil, errors.Wrap(err, "could not copy file: NewFile")
	}

	cap, err := NewGroupFileCap(path.Base(dataPath), dataPath, groupCtx.Ipfs, groupCtx.Storage)
	if err != nil {
		return nil, errors.Wrap(err, "could not create cap for NewFile")
	}

	idString := base64.URLEncoding.EncodeToString(cap.Id[:])
	capPath := groupCtx.Storage.GetGroupFileCapDir(groupCtx.Group.Id.ToString()) + idString

	file := &File{
		Cap:      cap,
		DataPath: dataPath,
		CapPath: capPath,
	}

	if err := file.Save(); err != nil {
		return nil, errors.Wrap(err, "could not save file")
	}

	return file, nil
}

func NewGroupFileFromCap(cap *FileCap, groupCtx *GroupContext) *File {
	idString := base64.URLEncoding.EncodeToString(cap.Id[:])
	capPath := groupCtx.Storage.GetGroupFileCapDir(groupCtx.Group.Id.ToString()) + idString
	dataPath := groupCtx.Storage.GetGroupFileDataDir(groupCtx.Group.Id.ToString()) + cap.FileName

	file := &File{
		Cap:      cap,
		DataPath: dataPath,
		CapPath:   capPath,
	}

	return file
}

func NewFile(srcFilePath, dstDir string, ctx *UserContext) (*File, error) {
	//newPath, err := ctx.Storage.CopyFileIntoMyFiles(srcFilePath)
	//if err != nil {
	//	return nil, fmt.Errorf("could not copy file: NewFilePTP: %s", err)
	//}
	//
	//cap, err := NewFileCap(newPath, ctx.Ipfs, ctx.Storage)
	//if err != nil {
	//	return nil, fmt.Errorf("could not create cap: NewFilePTP: %s", err)
	//}
	//
	//file := &File{
	//	Cap:      cap,
	//	DataPath: newPath,
	//}
	//
	//return file, nil
	return nil, nil
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

	if err := file.Save(); err != nil {
		return nil, errors.Wrapf(err, "could not save file '%s': NewFileFromCap", cap.FileName)
	}

	go file.Download(storage, ipfs)

	return file, nil
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

// Save saves File to disk
func (f *File) Save() error {
	jsonBytes, err := json.Marshal(f)
	if err != nil {
		return errors.Wrapf(err, "could not marshal file '%s': File.save", f.Cap.Id)
	}
	if err := utils.WriteFile(f.CapPath, jsonBytes); err != nil {
		return errors.Wrapf(err, "could not write file '%s': File.save", f.Cap.Id)
	}

	return nil
}

// Share file with a set of users, described by shareWith. Encrypted
// capabilities are made and copied in the corresponding 'public/for/'
// directories. The 'public' directory is re-published into IPNS. After
// that, notification messages are sent out.
func (f *File) Share(ctx *UserContext) error {

	//f.SharedWith = append(f.SharedWith, friend.Contact.Address)
	//
	//if err := f.Save(groupCtx.Storage); err != nil {
	//	return fmt.Errorf("could not save file: FilePTP.Share: %s", err)
	//}
	//
	//// make new capability into for_X directory
	//dirHash, err := groupCtx.Storage.GiveCAPToUser(f.Cap, friend.Contact, groupCtx.Ipfs, groupCtx.Network)
	//if err != nil {
	//	return fmt.Errorf("could not give cap to user: FilePTP.Share: %s", err)
	//}
	//
	//message, err := NewUpdateDirectory(groupCtx.User.Address, friend.Contact.Address, dirHash, groupCtx.User.Signer)
	//if err != nil {
	//	return fmt.Errorf("could not create new update dir message: FilePTP.Share: %s", err)
	//}
	//
	//if err := message.DialP2PConn(&friend.Contact.Boxer, &groupCtx.User.Boxer.PublicKey, groupCtx.User.Address, groupCtx.Network); err != nil {
	//	return fmt.Errorf("could not send message: FilePTP.Share: %s", err)
	//}

	return nil
}
