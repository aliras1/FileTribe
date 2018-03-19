package client

import "testing"

func TestNetwork_SendMessage(t *testing.T) {
	n := Network{"http://0.0.0.0:6000"}
	n.SendMessage("from", "to", "hello friend!")
	n.SendMessage("from", "to", "hello friend, again!")
}
