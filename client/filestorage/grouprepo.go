package filestorage

import (
	"log"
)

type GroupRepo struct {
	Files []*FileGroup
}

func (repo *GroupRepo) Bytes() []byte {
	var data []byte
	for _, file := range repo.Files {
		data = append(data, []byte(file.IPFSHash)...)
	}
	return data
}

func (repo *GroupRepo) Append(file *FileGroup) {
	repo.Files = append(repo.Files, file)
}

func (repo *GroupRepo) List() {
	for _, file := range repo.Files {
		log.Println("\t--> " + file.Name)
	}
}