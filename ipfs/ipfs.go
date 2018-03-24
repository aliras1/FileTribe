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

type IPFSID struct {
	ID              string   `json:"ID"`
	PublicKey       string   `json:"PublicKey"`
	Addresses       []string `json:"Addresses"`
	AgentVersion    string   `json:"AgentVersion"`
	ProtocolVersion string   `json:"ProtocolVersion"`
}

type IPFSNameResolvedHash struct {
	Path string `json:"Path"`
}

type ListLink struct {
	Name string `json:"Name"`
	Hash string `json:"Hash"`
	Size int    `json:"Size"`
	Type int    `json:"Type"`
}

type ListObject struct {
	Hash  string     `json:"Hash"`
	Links []ListLink `json:"Links"`
}

type ListObjects struct {
	Objects []ListObject `json:"Objects"`
}

type IPFS struct {
	host    string
	port    string
	version string
}

func NewIPFS(host string, port int) (*IPFS, error) {
	p := strconv.FormatInt(int64(port), 10)
	ipfs := IPFS{host, p, "/api/v0/"}
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
			m.AddSubDir(dirPath+"/"+f.Name(), dirName+"/"+f.Name())
		}
	}
	resp, err := m.Send()
	if err != nil {
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
	return merkleNodes, err
}

func (i *IPFS) Get(filePath, hash string) error {
	b, err := i.getRequest("get?arg=" + hash)
	if err != nil {
		return err
	}
	extractor := &tar.Extractor{filePath}
	return extractor.Extract(bytes.NewReader(b))
}

func (i *IPFS) ID() (*IPFSID, error) {
	bytesID, err := i.getRequest("id")
	if err != nil {
		return nil, err
	}
	var id IPFSID
	err = json.Unmarshal(bytesID, &id)
	return &id, err
}

func (i *IPFS) List(pathIPFS string) (*ListObjects, error) {
	bytesListObjectsJSON, err := i.getRequest("ls?arg=" + pathIPFS)
	if err != nil {
		return nil, err
	}
	var listObjects ListObjects
	err = json.Unmarshal(bytesListObjectsJSON, &listObjects)
	if err != nil {
		return nil, err
	}
	return &listObjects, nil
}

func (i *IPFS) NamePublish(hash string) error {
	_, err := i.getRequest("name/publish?arg=" + hash)
	if err != nil {
		return err
	}
	return nil
}

func (i *IPFS) NameResolve(ipnsPath string) (string, error) {
	resp, err := i.getRequest("name/resolve?arg=" + ipnsPath)
	if err != nil {
		return "", err
	}
	var hash IPFSNameResolvedHash
	err = json.Unmarshal(resp, &hash)
	return hash.Path, err
}

func (i *IPFS) Resolve(anyPath string) (string, error) {
	resp, err := i.getRequest("resolve?arg=" + anyPath + "&recursive=true")
	if err != nil {
		return "", err
	}
	var hash IPFSNameResolvedHash
	err = json.Unmarshal(resp, &hash)
	return hash.Path, err
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
