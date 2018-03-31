package client

import (
	"testing"

	nw "ipfs-share/network"
)

func TestGroup_RegisterAndGetGroup(t *testing.T) {
	n := nw.Network{"http://0.0.0.0:6000"}
	g := NewGroup("test_group")
	n.RegisterGroup(g.GroupName, g.Signer.PublicKey)
	registered, err := n.IsGroupRegistered(g.GroupName)
	if err != nil {
		t.Fatal(err)
	}
	if !registered {
		t.Fatal("not registered")
	}
	key, err := n.GetGroupSigningKey(g.GroupName)
	if err != nil {
		t.Fatal(err)
	}
	if !key.Equals(&g.Signer.PublicKey) {
		t.Fatal("pk's are not the same")
	}
}
