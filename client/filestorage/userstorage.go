package filestorage

import (
	"fmt"
	"ipfs-share/ipfs"
	"os"
	"strings"
)

type UserStorage struct {
	username string
	dataPath string
	RootDir  *Directory
	ipfs     *ipfs.IPFS
}

func (u *UserStorage) Init(path string, i *ipfs.IPFS) error {
	u.ipfs = i
	u.dataPath = path + "/data"
	var errorArray []error

	err := os.Mkdir(u.dataPath, 0700)
	if err != nil {
		errorArray = append(errorArray, err)
	}

	err = os.Mkdir(u.dataPath+"/public", 0700)
	if err != nil {
		errorArray = append(errorArray, err)
	}

	err = os.MkdirAll(u.dataPath+"/userdata/root", 0700)
	if err != nil {
		errorArray = append(errorArray, err)
	}

	f, err := os.Create(u.dataPath + "/userdata/caps.json")
	f.Close()
	if err != nil {
		errorArray = append(errorArray, err)
	}

	f, err = os.Create(u.dataPath + "/userdata/shared.json")
	f.Close()
	if err != nil {
		errorArray = append(errorArray, err)
	}

	u.build()
	return fmt.Errorf("%s", errorArray)
}

func (u *UserStorage) build() error {
	u.RootDir = NewDir(u.dataPath+"/userdata/root", "", u.username, []string{}, []string{})
	// TODO read .json files and build up the user storage
	return nil
}

func searchDir(directory *Directory, name string) *Entry {
	if strings.HasSuffix(directory.Path, name) {
		return &directory.Entry
	}
	for _, f := range directory.Files {
		if strings.HasSuffix(f.Path, name) {
			fmt.Println(f.Path + "    " + name)
			return f
		}
	}
	for _, f := range directory.SubDirectories {
		return searchDir(f, name)
	}
	return nil
}

func (u *UserStorage) AddAndShareDirectory(path string, shareWith, wAccess []string) error {
	dp := &Directory{
		Entry{path, "", u.username, shareWith, wAccess},
		[]*Directory{},
		[]*Entry{},
	}

	err := u.RootDir.AddLocalDirRecursively(dp)
	merkleNodes, err := u.ipfs.AddDir(path)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if merkleNodes == nil {
		return nil
	}

	for _, m := range merkleNodes {
		e := searchDir(dp, m.Name)
		if e == nil {
			fmt.Println("NIL by: " + m.Name)
		} else if strings.Compare(e.IPFSAddr, "") != 0 {
			fmt.Println("already visited: " + m.Name)
		} else {
			e.IPFSAddr = m.Hash
		}
	}

	fmt.Println("--------- list ----------")
	dp.List()
	fmt.Println("--------- end ----------")

	return err
}
