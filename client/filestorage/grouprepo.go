package filestorage

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

func (repo *GroupRepo) List() string {
	str := ""
	for _, file := range repo.Files {
		str += "\t--> " + file.Name + "\n"
	}
	return str
}