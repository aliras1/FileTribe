package ipfs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/whyrusleeping/tar-utils"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"strings"
)

type IPFS struct {
	host    string
	port    string
	version string

	Files
}

func NewIPFS(host string, port int) (*IPFS, error) {
	p := strconv.FormatInt(int64(port), 10)
	ipfs := IPFS{host, p, "/api/v0/", Files{}}
	_, err := ipfs.Version()
	if err != nil {
		return nil, errors.New("could not connect to ipfs daemon: " + err.Error())
	}
	return &ipfs, nil
}

func (i *IPFS) AddFile(filePath string) (*MerkleNode, error) {
	fileName := path.Base(filePath)
	url := i.host + ":" + i.port + i.version + "add?"
	m := NewMultipart(url)
	m.AddFile(filePath, fileName)
	resp, err := m.Send()
	if err != nil {
		return nil, err
	}
	var returnMerkleNode MerkleNode
	err = json.Unmarshal(resp, &returnMerkleNode)
	if err != nil {
		err = errors.New("could not unmarshal response: " + err.Error())
		return nil, err
	}
	return &returnMerkleNode, nil
}

func (i *IPFS) AddDir(dirPath string) ([]*MerkleNode, error) {
	dirName := path.Base(dirPath)
	url := i.host + ":" + i.port + i.version + "add?wrap-with-directory=true&pin=false"
	m := NewMultipart(url)

	// list dir
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for _, f := range files {
		if !f.IsDir() {
			m.AddFile(dirPath+"/"+f.Name(), dirName+"/"+f.Name())
		} else {
			fmt.Println(dirName + "/" + f.Name())
			m.AddSubDir(dirPath+"/"+f.Name(), dirName+"/"+f.Name())
		}
	}
	resp, err := m.Send()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	jsonStrings := strings.Split(string(resp), "\n")
	var merkleNodes []*MerkleNode
	for _, j := range jsonStrings {
		if strings.Compare(j, "") == 0 {
			break
		}
		var mn MerkleNode
		err = json.Unmarshal([]byte(j), &mn)
		merkleNodes = append(merkleNodes, &mn)
		if err != nil {
			err = errors.New("could not unmarshal response: " + err.Error())
			return nil, err
		}
	}
	fmt.Println(string(resp))
	return merkleNodes, err
}

func (i *IPFS) Get(hash string) error {
	b, err := i.getRequest("get?arg=" + hash)
	if err != nil {
		return err
	}

	extractor := &tar.Extractor{"/home/aliras/tmp/tmp"}
	return extractor.Extract(bytes.NewReader(b))
}

// Files commands
type Files struct {
}

func List(path string) {

}

func (i *IPFS) Version() (string, error) {
	version, err := i.getRequest("version")
	return string(version), err
}

func (i *IPFS) getRequest(path string) ([]byte, error) {
	resp, err := http.Get(i.host + ":" + i.port + "/" + i.version + path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, nil
}
