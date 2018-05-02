package network

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
	"testing"
)

func TestNetwork_SendMessage(t *testing.T) {
	n := Network{"http://0.0.0.0:6000"}
	n.SendMessage("from", "to", "type", "hello friend!")
	n.SendMessage("from", "to", "type", "hello friend, again!")
	messages, err := n.GetMessages("to")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(*messages[0])
	fmt.Println(*messages[1])
}

func TestNetwork_GetGroupState(t *testing.T) {
	owner := "testuser"
	groupName := "supergroup"
	network := Network{"http://0.0.0.0:6000"}
	if err := network.RegisterGroup(groupName, owner); err != nil {
		t.Fatal(err)
	}
	originalState := sha256.Sum256([]byte(owner))
	originalStateBase64 := base64.StdEncoding.EncodeToString(originalState[:])
	state, err := network.GetGroupState(groupName)
	if err != nil {
		t.Fatal(err)
	}
	stateBase64 := base64.StdEncoding.EncodeToString(state)
	if strings.Compare(originalStateBase64, stateBase64) != 0 {
		t.Fatal("state hashes do not match")
	}
}

func TestNetwork_GetGroupPrevState(t *testing.T) {
	owner := "testuser"
	groupName := "supergroup"
	network := Network{"http://0.0.0.0:6000"}
	if err := network.RegisterGroup(groupName, owner); err != nil {
		t.Fatal(err)
	}
	originalState := sha256.Sum256([]byte(owner))
	originalStateBase64 := base64.StdEncoding.EncodeToString(originalState[:])
	fmt.Println(originalStateBase64)

	newState := []byte{2}
	network.GroupInvite(groupName, "testuser2", newState)

	newStateFromNetwork, err := network.GetGroupPrevState(groupName, newState)
	newStateBase64FromNetwork := base64.StdEncoding.EncodeToString(newStateFromNetwork)
	fmt.Println(string(newStateBase64FromNetwork))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Compare(originalStateBase64, string(newStateBase64FromNetwork)) != 0 {
		t.Fatal("state hashes do not match")
	}
}
