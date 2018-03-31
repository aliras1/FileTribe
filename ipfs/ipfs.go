package ipfs

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/whyrusleeping/tar-utils"
	"io/ioutil"
	"ipfs-share/crypto"
	"log"
	"net"
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

type PubsubMessage struct {
	From             string   `json:"from"`
	Data             string   `json:"data"`
	Seqno            string   `json:"seqno"`
	TopicIDs         []string `json:"topicIDs"`
	XXX_unrecognized []uint8  `json:"XXX_unrecognized,omitempty"`
}

func (psm *PubsubMessage) Decode() ([]byte, error) {
	var msgLVL1 []byte
	msgLVL1, err := base64.StdEncoding.DecodeString(psm.Data)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(msgLVL1))
	var msgLVL2 []byte
	msgLVL2, err = base64.URLEncoding.DecodeString(string(msgLVL1))
	if err != nil {
		return nil, err
	}
	return msgLVL2, nil
}

func (psm *PubsubMessage) Verify(key crypto.PublicSigningKey) ([]byte, bool) {
	signedGroupMsg, err := psm.Decode()
	if err != nil {
		log.Println(err)
		return nil, false
	}
	return key.Open(nil, signedGroupMsg)
}

type IPFSNameResolvedHash struct {
	Path string `json:"Path"`
}

type ListObjects struct {
	Objects []struct {
		Hash  string `json:"Hash"`
		Links []struct {
			Name string `json:"FileName"`
			Hash string `json:"Hash"`
			Size int    `json:"Size"`
			Type int    `json:"Type"`
		} `json:"Links"`
	} `json:"Objects"`
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

func (i *IPFS) PubsubPublish(channel string, message []byte) error {
	_, err := i.getRequest("pubsub/pub?arg=" + channel + "&arg=" + base64.URLEncoding.EncodeToString(message))
	return err
}

func (i *IPFS) PubsubSubscribe(channel string, dst chan PubsubMessage) error {
	conn, err := net.Dial("tcp", "127.0.0.1:5001")
	if err != nil {
		return err
	}
	req := "GET /api/v0/pubsub/sub?arg=" + channel + " HTTP/1.1\nHost: localhost:5001\n\n"
	conn.Write([]byte(req))
	_, err = bufio.NewReader(conn).ReadString('\n') // HTTP 200 response
	if err != nil {
		return err
	}
	// pubsub messages
	for {
		rawStr, err := bufio.NewReader(conn).ReadString('}')
		if err != nil {
			return err
		}
		split := strings.Split(rawStr, "\n")
		if len(split) < 2 {
			log.Println("invalid pubsub message")
			continue
		}
		var msg PubsubMessage
		err = json.Unmarshal([]byte((split)[1]), &msg)
		if err != nil {
			log.Println(err)
			continue
		}
		dst <- msg
	}
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
