package filestorage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	nw "ipfs-share/network"
)

type EntryProvider interface {
}

type File struct {
	Path       string   `json:"path"`
	Hash       string   `json:"ipfs_hash"`
	Owner      string   `json:"owner"`
	SharedWith []string `json:"shared_with"`
	WAccess    []string `json:"w_access"`
}

func (f *File) Share(shareWith []string, baseDirPath string, network *nw.Network) error {
	// TODO check if there exists a more efficient way to merge 2 lists
	for _, i := range shareWith {
		skip := false
		for j := 0; j < len(f.SharedWith) && !skip; j++ {
			if strings.Compare(i, f.SharedWith[j]) == 0 {
				skip = true
			}
		}
		if !skip {
			// add to share list
			f.SharedWith = append(f.SharedWith, i)

			// make new capability into for_X directory
			err := os.Mkdir(baseDirPath+i, 0770)
			if err != nil {
				fmt.Println(err) /* TODO check for permission errors */
			}
			jsonMap := make(map[string]string)
			jsonMap["name"] = path.Base(f.Path)
			jsonMap["hash"] = path.Base(f.Hash)
			byteJson, err := json.Marshal(jsonMap)
			err = ioutil.WriteFile(baseDirPath+i+"/"+path.Base(f.Path)+".json", byteJson, 0644)
			if err != nil {
				return err
			}
			// send share message
			err = network.SendMessage(f.Owner, i, path.Base(f.Path))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
