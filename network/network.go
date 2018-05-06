package network

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"crypto/sha256"
	"encoding/base64"
	"ipfs-share/crypto"
	"log"
)

type Message struct {
	From    string `json:"from"`
	To      string `json:"to,omitempty"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

type Network struct {
	Address string
}

func (n *Network) Get(path string, args ...string) ([]byte, error) {
	url := n.Address + path
	for _, arg := range args {
		url += "/" + arg
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, nil
}

func (n *Network) GetGroupMembers(groupName string) ([]string, error) {
	membersBytes, err := n.Get("/get/group/members", groupName)
	if err != nil {
		return nil, fmt.Errorf("could not get group members: Network.GetGroupMembers: %s", err)
	}
	members := []string{}
	if err := json.Unmarshal(membersBytes, &members); err != nil {
		log.Printf("memberbytes: '%s'\n", string(membersBytes))
		return nil, fmt.Errorf("could not unmarshal group members: Network.GetGroupMembers: %s", err)
	}
	return members, nil
}

func (n *Network) GetUserPublicKeyHash(username string) (crypto.PublicKeyHash, error) {
	base64PublicKeyHash, err := n.Get("/get/user/publickeyhash", username)
	if err != nil {
		return nil, err
	}
	return crypto.Base64ToPublicKeyHash(string(base64PublicKeyHash))
}

func (n *Network) GetUserSigningKey(username string) (crypto.PublicSigningKey, error) {
	base64PublicSigningKey, err := n.Get("/get/user/signkey", username)
	if err != nil {
		return nil, err
	}
	return crypto.Base64ToPublicSigningKey(string(base64PublicSigningKey))
}

func (n *Network) GetUserBoxingKey(username string) (crypto.PublicBoxingKey, error) {
	base64PublicBoxingKey, err := n.Get("/get/user/boxkey", username)
	if err != nil {
		return [32]byte{}, err
	}
	return crypto.Base64ToPublicBoxingKey(string(base64PublicBoxingKey))
}

func (n *Network) GetUserIPFSAddr(username string) (string, error) {
	bytesIPFSAddr, err := n.Get("/get/user/ipfsaddr", username)
	if err != nil {
		return "", err
	}
	return string(bytesIPFSAddr), nil
}

func (n *Network) GetGroupState(groupName string) ([]byte, error) {
	stateBase64Bytes, err := n.Get("/get/group/state", groupName)
	if err != nil {
		return nil, fmt.Errorf("error while getting state of group %s: Network.GetGroupState: %s", groupName, err)
	}
	state, err := base64.StdEncoding.DecodeString(string(stateBase64Bytes))
	if err != nil {
		return nil, fmt.Errorf("error while decoding state of group %s: Network.GetGroupState: %s", groupName, err)
	}
	return state, nil
}

func (n *Network) GetGroupPrevState(groupName string, state []byte) ([]byte, error) {
	stateBase64 := base64.StdEncoding.EncodeToString(state)
	prevStateBase64Bytes, err := n.Get("/get/group/prev/state", groupName, stateBase64)
	if err != nil {
		return nil, fmt.Errorf("error while getting previous state %s of group %s: Network.GetGroupPrevState: %s", state, groupName, err)
	}
	prevState, err := base64.StdEncoding.DecodeString(string(prevStateBase64Bytes))
	if err != nil {
		return nil, fmt.Errorf("error while decoding state of group %s: Network.GetGroupPrevState: %s", groupName, err)
	}
	return prevState, nil
}

func (n *Network) GroupInvite(groupname string, transaction []byte) error {
	err := n.Put(
		"/group/invite/"+groupname,
		"application/json",
		transaction,
	)
	if err != nil {
		return fmt.Errorf("error while inviting into %s: Network.GroupInvite: %s", groupname, err)
	}
	return nil
}

func (n *Network) IsGroupRegistered(groupName string) (bool, error) {
	boolString, err := n.Get("/is/group/registered", groupName)
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(string(boolString))
}

func (n *Network) IsUsernameRegistered(username string) (bool, error) {
	boolString, err := n.Get("/is/username/registered", username)
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(string(boolString))
}

func (n *Network) GetMessages(username string) ([]*Message, error) {
	resp, err := n.Get("/get/messages", username)
	if err != nil {
		return nil, err
	}
	var messages []*Message
	err = json.Unmarshal(resp, &messages)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (n *Network) Put(path string, contentType string, data []byte) error {
	resp, err := http.Post(
		fmt.Sprintf(n.Address+path),
		contentType,
		bytes.NewReader(data),
	)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if len(body) != 0 {
		return errors.New(string(body))
	}
	return nil
}

func (n *Network) PutSigningKey(hash crypto.PublicKeyHash, key crypto.PublicSigningKey) error {
	jsonStr := fmt.Sprintf(`{"hash":"%s", "signkey":"%s"}`, hash.ToBase64(), key.ToBase64())
	return n.Put(
		"/put/signkey",
		"application/json",
		[]byte(jsonStr),
	)
}

func (n *Network) PutBoxingKey(hash crypto.PublicKeyHash, key crypto.PublicBoxingKey) error {
	jsonStr := fmt.Sprintf(`{"hash":"%s", "boxkey":"%s"}`, hash.ToBase64(), key.ToBase64())
	return n.Put(
		"/put/boxkey",
		"application/json",
		[]byte(jsonStr),
	)
}

func (n *Network) PutIPFSAddr(hash crypto.PublicKeyHash, ipfsAddr string) error {
	jsonStr := fmt.Sprintf(`{"hash":"%s", "ipfsaddr":"%s"}`, hash.ToBase64(), ipfsAddr)
	return n.Put(
		"/put/ipfsaddr",
		"application/json",
		[]byte(jsonStr),
	)
}

func (n *Network) RegisterGroup(groupName, owner string) error {
	stateHash := sha256.Sum256([]byte(owner))
	stateHashBase64 := base64.StdEncoding.EncodeToString(stateHash[:])
	jsonStr := fmt.Sprintf(`{"groupname":"%s", "owner":"%s", "state":"%s"}`, groupName, owner, stateHashBase64)
	fmt.Println(jsonStr)
	return n.Put(
		"/register/group",
		"application/json",
		[]byte(jsonStr))
}

func (n *Network) RegisterUsername(username string, hash crypto.PublicKeyHash) error {
	return n.Put(
		fmt.Sprintf("/register/username/%s", username),
		"application/octet-stream",
		[]byte(hash.ToBase64()),
	)
}

func (n *Network) SendMessage(from, to, msgType, msgData string) error {
	m := Message{from, to, msgType, msgData}
	fmt.Println(m)
	byteJson, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return n.Put("/send/message", "application/json", byteJson)
}
