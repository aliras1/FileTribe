package client

import (
	"crypto/rand"
	"errors"

	fs "ipfs-share/client/filestorage"
	"ipfs-share/crypto"
	nw "ipfs-share/network"
)

type Group struct {
	GroupName string
	Boxer     crypto.SymmetricKey
}

func NewGroup(groupName string) *Group {
	var secretKeyBytes [32]byte
	rand.Read(secretKeyBytes[:])
	return &Group{groupName, crypto.SymmetricKey{secretKeyBytes, rand.Reader}}
}

func (g *Group) Save(storage *fs.Storage) error {
	return storage.StoreGroupAccessCAP(g.GroupName, g.Boxer)
}

func (g *Group) Register(network *nw.Network) error {
	registered, err := network.IsGroupRegistered(g.GroupName)
	if err != nil {
		return err
	}
	if registered {
		return errors.New("group name already exists")
	}
	err = network.RegisterGroup(g.GroupName)
	if err != nil {
		return err
	}
	return nil
}
