package network

import (
	"fmt"
	"testing"
)

func TestNetwork_SendMessage(t *testing.T) {
	n := Network{"http://0.0.0.0:6000"}
	n.SendMessage("from", "to", "hello friend!")
	n.SendMessage("from", "to", "hello friend, again!")
	messages, err := n.GetMessages("to")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(*messages[0])
	fmt.Println(*messages[1])
}
