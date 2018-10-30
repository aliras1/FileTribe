package client

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/golang/glog"
	ipfsapi "ipfs-share/ipfs"
	"ipfs-share/collections"
	"github.com/pkg/errors"
	"ipfs-share/utils"
	"io"
	"bytes"
)

// IFile is an interface for the files
// which can be shared
type IFile interface {
	Share()
}

type PendingFile File

func (pending *PendingFile) Id() collections.IIdentifier {
	return collections.NewStringId(pending.DataPath)
}

func (pending *PendingFile) SaveMetadata() error {
	jsonBytes, err := json.Marshal((*File)(pending))
	if err != nil {
		return errors.Wrapf(err, "could not marshal file '%s': File.save", pending.Cap.Id)
	}
	if err := utils.WriteFile(pending.PendingPath, jsonBytes); err != nil {
		return errors.Wrapf(err, "could not write file '%s': File.save", pending.Cap.Id)
	}

	return nil
}

// File represents a file that
// is shared in a peer to peer mode
type File struct {
	Cap      *FileCap
	DataPath string
	CapPath string
	PendingPath string
}

func (f *File) Id() collections.IIdentifier {
	return collections.NewBytesId(f.Cap.Id)
}

func NewGroupFile(filePath string, groupCtx *GroupContext, hasWriteAccess []ethcommon.Address) (*File, error) {
	cap, err := NewGroupFileCap(path.Base(filePath), filePath, hasWriteAccess, groupCtx.Ipfs, groupCtx.Storage)
	if err != nil {
		return nil, errors.Wrap(err, "could not create cap for NewFile")
	}

	idString := base64.URLEncoding.EncodeToString(cap.Id[:])
	capPath := groupCtx.Storage.GetGroupFileCapDir(groupCtx.Group.Id.ToString()) + idString
	pendingPath := groupCtx.Storage.GetGroupFilePendingDir(groupCtx.Group.Id.ToString()) + idString

	file := &File{
		Cap:      cap,
		DataPath: filePath,
		CapPath: capPath,
		PendingPath: pendingPath,
	}

	if err := file.SaveMetadata(); err != nil {
		return nil, errors.Wrap(err, "could not save file meta data")
	}

	return file, nil
}

func NewGroupFileFromCap(cap *FileCap, groupCtx *GroupContext) *File {
	idString := base64.URLEncoding.EncodeToString(cap.Id[:])
	capPath := groupCtx.Storage.GetGroupFileCapDir(groupCtx.Group.Id.ToString()) + idString
	dataPath := groupCtx.Storage.GetGroupFileDataDir(groupCtx.Group.Id.ToString()) + cap.FileName
	pendingPath := groupCtx.Storage.GetGroupFilePendingDir(groupCtx.Group.Id.ToString()) + idString

	file := &File{
		Cap:      cap,
		DataPath: dataPath,
		CapPath:   capPath,
		PendingPath: pendingPath,
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

	if err := file.SaveMetadata(); err != nil {
		return nil, errors.Wrapf(err, "could not save file '%s': NewFileFromCap", cap.FileName)
	}

	go file.Download(storage, ipfs)

	return file, nil
}

// Downloads, decrypts and verifies the content of file from Ipfs
func (f *File) Download(storage *Storage, ipfs ipfsapi.IIpfs) {
	// TODO: check if current file data's hash does not match with cap's ipfs hash

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

// SaveMetadata saves FileMetaData to disk
func (f *File) SaveMetadata() error {
	jsonBytes, err := json.Marshal(f)
	if err != nil {
		return errors.Wrapf(err, "could not marshal file '%s': File.save", f.Cap.Id)
	}
	if err := utils.WriteFile(f.CapPath, jsonBytes); err != nil {
		return errors.Wrapf(err, "could not write file '%s': File.save", f.Cap.Id)
	}

	return nil
}


func (f *File) Encrypt() (io.Reader, error) {
	data, err := ioutil.ReadFile(f.DataPath)
	if err != nil {
		return nil, errors.Wrap(err, "could not read file")
	}
	encData, err := f.Cap.DataKey.Seal(bytes.NewReader(data))
	if err != nil {
		return nil, errors.Wrap(err, "could not encrypt data")
	}

	return encData, nil
}

func GetCapListFromFileList(files []*File) []*FileCap {
	var l []*FileCap
	for _, file := range files {
		l = append(l, file.Cap)
	}
	return l
}

// Share file with a set of users, described by shareWith. Encrypted
// capabilities are made and copied in the corresponding 'public/for/'
// directories. The 'public' directory is re-published into IPNS. After
// that, notification messages are sent out.
func (f *File) Share(ctx *UserContext) error {

	//f.SharedWith = append(f.SharedWith, friend.Contact.Address)
	//
	//if err := f.SaveMetadata(groupCtx.Storage); err != nil {
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
