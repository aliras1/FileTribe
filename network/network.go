package network

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"ipfs-share/crypto"
)

type Message struct {
	From    string `json:"from"`
	To      string `json:"to,omitempty"`
	Message string `json:"message"`
}

type Network struct {
	Address string
}

func (n *Network) Get(path string, id string) ([]byte, error) {
	resp, err := http.Get(fmt.Sprintf(n.Address + path + id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, nil
}

func (n *Network) GetUserPublicKeyHash(username string) (crypto.PublicKeyHash, error) {
	base64PublicKeyHash, err := n.Get("/get/user/publickeyhash/", username)
	if err != nil {
		return nil, err
	}
	return crypto.Base64ToPublicKeyHash(string(base64PublicKeyHash))
}

func (n *Network) GetUserSigningKey(username string) (crypto.PublicSigningKey, error) {
	base64PublicSigningKey, err := n.Get("/get/user/signkey/", username)
	if err != nil {
		return nil, err
	}
	return crypto.Base64ToPublicSigningKey(string(base64PublicSigningKey))
}

func (n *Network) GetUserBoxingKey(username string) (crypto.PublicBoxingKey, error) {
	base64PublicBoxingKey, err := n.Get("/get/user/boxkey/", username)
	if err != nil {
		return [32]byte{}, err
	}
	return crypto.Base64ToPublicBoxingKey(string(base64PublicBoxingKey))
}

func (n *Network) GetUserIPFSAddr(username string) (string, error) {
	bytesIPFSAddr, err := n.Get("/get/user/ipfsaddr/", username)
	if err != nil {
		return "", err
	}
	return string(bytesIPFSAddr), nil
}

func (n *Network) IsUsernameRegistered(username string) (bool, error) {
	boolString, err := n.Get("/is/username/registered/", username)
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(string(boolString))
}

func (n *Network) GetMessages(username string) ([]*Message, error) {
	resp, err := n.Get("/get/messages/", username)
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

func (n *Network) RegisterUsername(username string, hash crypto.PublicKeyHash) error {
	return n.Put(
		fmt.Sprintf("/register/username/%s", username),
		"application/octet-stream",
		[]byte(hash.ToBase64()),
	)
}

func (n *Network) SendMessage(from, to, msg string) error {
	m := Message{from, to, msg}
	fmt.Println(m)
	byteJson, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return n.Put("/send/message", "application/json", byteJson)
}
