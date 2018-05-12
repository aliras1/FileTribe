package ipfs

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"path"
	"strconv"
	"strings"
	"github.com/golang/glog"
	"github.com/whyrusleeping/tar-utils"

	"ipfs-share/crypto"
	"sync"
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
	var msgLVL2 []byte
	msgLVL2, err = base64.URLEncoding.DecodeString(string(msgLVL1))
	if err != nil {
		return nil, err
	}
	return msgLVL2, nil
}

func (psm *PubsubMessage) Decrypt(key crypto.SymmetricKey) ([]byte, bool) {
	messageBytes, err := psm.Decode()
	if err != nil {
		glog.Errorf("could not decode pubsub message: PubsubMessage.Decrypt: %s", err)
		return nil, false
	}
	return key.BoxOpen(messageBytes)
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

type PubSubSubscription struct {
	conn net.Conn
	channel chan PubsubMessage
}

type IPFS struct {
	host    string
	port    string
	version string

	pubSubSubscriptions map[string]PubSubSubscription
	mux sync.Mutex
}

func NewIPFS(host string, port int) (*IPFS, error) {
	p := strconv.FormatInt(int64(port), 10)
	pubSubSubscriptions := make(map[string]PubSubSubscription)
	ipfs := IPFS{
		host:  host,
		port:  p,
		version:  "/api/v0/",
		pubSubSubscriptions: pubSubSubscriptions,
	}
	_, err := ipfs.Version()
	if err != nil {
		return nil, fmt.Errorf("could not connect to ipfs daemon: " + err.Error())
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
		err = fmt.Errorf("could not unmarshal response: " + err.Error())
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
		glog.Errorf("could not read dir '%s': IPFS.AddDir: %s", dirPath, err)
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
			err = fmt.Errorf("could not unmarshal response: " + err.Error())
			return nil, err
		}
	}
	return merkleNodes, err
}

func (i *IPFS) Get(filePath, hash string) error {
	b, err := i.getRequest("get?arg=" + hash + "&archive=true")
	if err != nil {
		return fmt.Errorf("error while ipfs api request: IPFS.Get: %s", err)
	}
	extractor := &tar.Extractor{
		Path:     filePath,
		Progress: nil,
	}
	if err := extractor.Extract(bytes.NewReader(b)); err != nil {
		return fmt.Errorf("error while extracting: IPFS.Get: %s", err)
	}
	return nil
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

func (i *IPFS) KillAll() {
	glog.Info("Killing all ipfs subs....")
	i.mux.Lock()
	for key, subscription := range i.pubSubSubscriptions {
		glog.Info("killing channel '%s': IPFS.KillAll", key)
		subscription.conn.Close()
		close(subscription.channel)
	}
	i.mux.Unlock()
}

func (i *IPFS) Kill(channelName string) {
	glog.Infof("Killing channel '%s': IPFS.Kill", channelName)
	i.mux.Lock()
	subscription, ok := i.pubSubSubscriptions[channelName]
	if !ok {
		glog.Warningf("subscription '%s' not found: IPFS.Kill", channelName)
		i.mux.Unlock()
		return
	}
	subscription.conn.Close()
	close(subscription.channel)
	delete(i.pubSubSubscriptions, channelName)
	i.mux.Unlock()
	fmt.Println(i.pubSubSubscriptions)
}

func (i *IPFS) PubsubPublish(channel string, message []byte) error {
	glog.Infof("sending pubsub message on channel '%s'...", channel)
	if _, err := i.getRequest("pubsub/pub?arg=" + channel + "&arg=" + base64.URLEncoding.EncodeToString(message)); err != nil {
		return fmt.Errorf("could not make get request: IPFS.PubsubPublish: %s", err)
	}
	return nil
}

func (i *IPFS) PubsubSubscribe(username, channel string, dst chan PubsubMessage) {
	conn, err := net.Dial("tcp", "127.0.0.1:5001")
	if err != nil {
		glog.Errorf("could not reach ipfs daemon: IPFS.PubsubSubscribe: %s", err)
		return
	}
	if _, ok := i.pubSubSubscriptions[channel]; ok {
		glog.Warningf("user '%s' already subscribed on topic '%s': IPFS.PubsubSubscribe", username, channel)
		return
	}
	i.pubSubSubscriptions[channel] = PubSubSubscription{
		channel: dst,
		conn: conn,
	}
	fmt.Println(i.pubSubSubscriptions)
	req := "GET /api/v0/pubsub/sub?arg=" + channel + " HTTP/1.1\nHost: localhost:5001\n\n"
	conn.Write([]byte(req))
	if _, err := bufio.NewReader(conn).ReadString('\n'); err != nil { // HTTP 200 response
		glog.Errorf("could not read 'HTTP 200' response from ipfs daemon: IPFS.PubsubSubscribe: %s", err)
		return
	}
	glog.Infof("user '%s' on topic '%s' runs...", username, channel)
	// pubsub messages
	for {
		rawStr, err := bufio.NewReader(conn).ReadString('}')
		if err != nil {
			glog.Errorf("could not read response from ipfs daemon: IPFS.PubsubSubscribe: %s", err)
			return
		}
		split := strings.Split(rawStr, "\n")
		if len(split) < 2 {
			glog.Warningf("invalid pubsub message: IPFS.PubsubSubscribe")
			continue
		}
		var msg PubsubMessage
		if err := json.Unmarshal([]byte((split)[1]), &msg); err != nil {
			glog.Warningf("could not unmarshal pubsub message: IPFS.PubsubSubscribe: %s", err)
			continue
		}

		i.mux.Lock()
		if _, ok := i.pubSubSubscriptions[channel]; ok {
			dst <- msg
		}
		i.mux.Unlock()
	}
}

func (i *IPFS) Resolve(anyPath string) (string, error) {
	glog.Info("IPFS resolving...: " + anyPath)
	resp, err := i.getRequest("resolve?arg=" + anyPath + "&recursive=true")
	if err != nil {
		return "", err
	}
	var hash IPFSNameResolvedHash
	if err := json.Unmarshal(resp, &hash); err != nil {
		return "", fmt.Errorf("could not unmarshal IPFSNameResolvedHash '%s': IPFS.Resolve: %s", anyPath, err)
	}
	glog.Info("Resolved")
	if strings.Compare(hash.Path, "") == 0 {
		return "", fmt.Errorf("could not resolve path '%s': IPFS.Resolve", anyPath)
	}
	return hash.Path, nil
}

func (i *IPFS) Version() (string, error) {
	version, err := i.getRequest("version")
	return string(version), err
}

func (i *IPFS) getRequest(path string) ([]byte, error) {
	resp, err := http.Get(i.host + ":" + i.port + "/" + i.version + path)
	if err != nil {
		return nil, fmt.Errorf("error while http.get: IPFS.getRequest: %s", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading http response body: IPFS.getRequest: %s", err)
	}
	return body, nil
}
