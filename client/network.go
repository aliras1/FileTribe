package client

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"encoding/json"
	"ipfs-share/crypto"
)

type Network struct {
	Address string
}

func (n *Network) Get(path string, id string) ([]byte, error) {
	resp, err := http.Get(fmt.Sprintf("%s%s%s", n.Address, path, id))
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

func (n *Network) isUsernameRegistered(username string) (bool, error) {
	boolString, err := n.Get("/is/username/registered/", username)
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(string(boolString))
}

func (n *Network) Put(path string, contentType string, data []byte) error {
	resp, err := http.Post(
		fmt.Sprintf("%s%s", n.Address, path),
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

func (n *Network) RegisterUsername(username string, hash crypto.PublicKeyHash) error {
	return n.Put(
		fmt.Sprintf("/register/username/%s", username),
		"application/octet-stream",
		[]byte(hash.ToBase64()),
	)
}

func (n *Network) SendMessage(from, to, msg string) error {
	jsonMap := make(map[string]string)
	jsonMap["from"] = from
	jsonMap["to"] = to
	jsonMap["msg"] = msg
	byteJson, err := json.Marshal(jsonMap)
	if err != nil {
		return err
	}
	return n.Put("/send/message", "application/json", byteJson)
}
