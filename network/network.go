package network

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"ipfs-share/crypto"
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
		return nil, fmt.Errorf("error while http get request '%s': Network.Get: %s", url, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: Network.Get: %s", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http get request '%s' returned with status code '%d'", url, resp.StatusCode)
	}
	return body, nil
}

func (n *Network) GetGroupMembers(groupName string) ([]string, error) {
	membersBytes, err := n.Get("/get/group/members", groupName)
	if err != nil {
		return nil, fmt.Errorf("could not get group members: Network.GetGroupMembers: %s", err)
	}
	var members []string
	if err := json.Unmarshal(membersBytes, &members); err != nil {
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

func (n *Network) GetUserVerifyKey(username string) (crypto.Signer, error) {
	// base64PublicSigningKey, err := n.Get("/get/user/signkey", username)
	// if err != nil {
	// 	return crypto.Signer{}, err
	// }
	// return crypto.Base64ToPublicSigningKey(string(base64PublicSigningKey))
	return crypto.Signer{}, nil
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

func (n *Network) GetGroupOperation(groupName string, state []byte) ([]byte, error) {
	stateBase64 := base64.StdEncoding.EncodeToString(state)
	operationBytes, err := n.Post(
		"/get/group/operation/"+groupName,
		"application/octet-stream",
		[]byte(stateBase64),
	)
	if err != nil {
		return nil, fmt.Errorf("error while getting operation of state %s of group %s: Network.GetGroupOperation: %s", state, groupName, err)
	}
	return operationBytes, nil
}

func (n *Network) GroupInvite(groupname string, transaction []byte) error {
	if _, err := n.Post(
		"/group/invite/"+groupname,
		"application/json",
		transaction,
	); err != nil {
		return fmt.Errorf("error while inviting into %s: Network.GroupInvite: %s", groupname, err)
	}
	return nil
}

func (n *Network) GroupShare(groupname string, transaction []byte) error {
	if _, err := n.Post(
		"/group/share/"+groupname,
		"application/json",
		transaction,
	); err != nil {
		return fmt.Errorf("error while sharing file with group '%s', Network.GroupShare: %s", groupname, err)
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

func (n *Network) Post(path string, contentType string, data []byte) ([]byte, error) {
	url := n.Address + path
	resp, err := http.Post(
		url,
		contentType,
		bytes.NewReader(data),
	)
	if err != nil {
		return nil, fmt.Errorf("error while http post request '%s': Network.Post: %s", url, err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: Network.Post: %s", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("request '%s' returned with status code '%d': Network.Post", url, resp.StatusCode)
	}
	return body, nil
}

func (n *Network) PutVerifyKey(hash crypto.PublicKeyHash, key crypto.Signer) error {
	// jsonStr := fmt.Sprintf(`{"hash":"%s", "signkey":"%s"}`, hash.ToBase64(), key.ToBase64())
	// _, err := n.Post(
	// 	"/put/signkey",
	// 	"application/json",
	// 	[]byte(jsonStr),
	// )
	// if err != nil {
	// 	return fmt.Errorf("error while post: Network.PutVerifyKey: %s", err)
	// }
	return nil
}

func (n *Network) PutBoxingKey(hash crypto.PublicKeyHash, key crypto.PublicBoxingKey) error {
	jsonStr := fmt.Sprintf(`{"hash":"%s", "boxkey":"%s"}`, hash.ToBase64(), key.ToBase64())
	_, err := n.Post(
		"/put/boxkey",
		"application/json",
		[]byte(jsonStr),
	)
	if err != nil {
		return fmt.Errorf("error while post: Network.PutBoxingKey: %s", err)
	}
	return nil
}

func (n *Network) PutIPFSAddr(hash crypto.PublicKeyHash, ipfsAddr string) error {
	jsonStr := fmt.Sprintf(`{"hash":"%s", "ipfsaddr":"%s"}`, hash.ToBase64(), ipfsAddr)
	_, err := n.Post(
		"/put/ipfsaddr",
		"application/json",
		[]byte(jsonStr),
	)
	if err != nil {
		return fmt.Errorf("error while post: Network.PutIPFSAddr: %s", err)
	}
	return nil
}

func (n *Network) RegisterGroup(groupName, owner string) error {
	stateHash := sha256.Sum256([]byte(owner))
	stateHashBase64 := base64.StdEncoding.EncodeToString(stateHash[:])
	jsonStr := fmt.Sprintf(`{"groupname":"%s", "owner":"%s", "state":"%s"}`, groupName, owner, stateHashBase64)

	_, err := n.Post(
		"/register/group",
		"application/json",
		[]byte(jsonStr),
	)
	if err != nil {
		return fmt.Errorf("error while post: Network.PutBoxingKey: %s", err)
	}
	return nil
}

func (n *Network) RegisterUsername(username string, hash crypto.PublicKeyHash) error {
	_, err := n.Post(
		fmt.Sprintf("/register/username/%s", username),
		"application/octet-stream",
		[]byte(hash.ToBase64()),
	)
	if err != nil {
		return fmt.Errorf("error while post: Network.PutBoxingKey: %s", err)
	}
	return nil
}

func (n *Network) SendMessage(from, to, msgType, msgData string) error {
	m := Message{from, to, msgType, msgData}

	byteJSON, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("could not marshal message")
	}
	_, err = n.Post("/send/message", "application/json", byteJSON)
	if err != nil {
		return fmt.Errorf("error while post: Network.PutBoxingKey: %s", err)
	}
	return nil
}
