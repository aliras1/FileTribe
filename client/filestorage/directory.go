package filestorage

import (
	"fmt"
	"io/ioutil"
)

type Directory struct {
	Entry
	SubDirectories []*Directory `json:"sub_directories"`
	Files          []*Entry     `json:"files"`
}

func NewDir(path, ipfsAddr, owner string, sharedWith, wAccess []string) *Directory {
	e := Entry{path, ipfsAddr, owner, sharedWith, wAccess}
	return &Directory{e, []*Directory{}, []*Entry{}}
}

func (d *Directory) AddLocalDirRecursively(other *Directory) error {
	files, err := ioutil.ReadDir(other.Path)
	if err != nil {
		fmt.Println(err)
		return err
	}

	for _, f := range files {
		if !f.IsDir() {
			e := Entry{other.Path + "/" + f.Name(), "", other.Owner, []string{}, []string{}}
			other.Files = append(other.Files, &e)
		} else {
			e := Entry{other.Path + "/" + f.Name(), "", other.Owner, []string{}, []string{}}
			newDir := Directory{e, []*Directory{}, []*Entry{}}
			other.AddLocalDirRecursively(&newDir)
		}
	}
	d.SubDirectories = append(d.SubDirectories, other)
	return nil
}

func (d *Directory) List() {
	fmt.Println(d.Path + ": " + d.IPFSAddr)
	for _, f := range d.Files {
		fmt.Println(f.Path + ": " + f.IPFSAddr)
	}
	for _, f := range d.SubDirectories {
		f.List()
	}
}

func (d *Directory) mkdir(name string) *Entry {
	return nil
}
