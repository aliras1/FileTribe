package filestorage

import (
	"fmt"
)

type Entry interface {
}

type Directory struct {
	Path    string
	Entries []*Entry
}

func (d *Directory) mkdir(name string) *Entry {
	newD := Entry(Directory{fmt.Sprintf("%s/%s", d.Path, name), nil})
	d.Entries = append(d.Entries, &newD)
	return &newD
}
