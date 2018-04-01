package client

import (
	"testing"

	nw "ipfs-share/network"
)

func TestGroup_RegisterAndGetGroup(t *testing.T) {
	n := nw.Network{"http://0.0.0.0:6000"}
	g := NewGroup("test_group")
	n.RegisterGroup(g.GroupName)
	registered, err := n.IsGroupRegistered(g.GroupName)
	if err != nil {
		t.Fatal(err)
	}
	if !registered {
		t.Fatal("not registered")
	}
}
