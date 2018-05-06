package client

import (
	"crypto/rand"
	"fmt"

	fs "ipfs-share/client/filestorage"
	"ipfs-share/crypto"
	nw "ipfs-share/network"
)

type Group struct {
	Name  string
	Boxer crypto.SymmetricKey
}

func NewGroup(groupName string) *Group {
	var secretKeyBytes [32]byte
	rand.Read(secretKeyBytes[:])
	return &Group{
		Name:groupName,
		Boxer: crypto.SymmetricKey{
			Key: secretKeyBytes,
			RNG: rand.Reader,
		},
	}
}

func (g *Group) Save(storage *fs.Storage) error {
	if err := storage.StoreGroupAccessCAP(g.Name, g.Boxer); err != nil {
		return fmt.Errorf("could not store ga cap: Group.Save: %s", err)
	}
	return nil
}

func (g *Group) Register(owner string, network *nw.Network) error {
	registered, err := network.IsGroupRegistered(g.Name)
	if err != nil {
		return fmt.Errorf("could not check if group is registered: Group.Register: %s", err)
	}
	if registered {
		return fmt.Errorf("group name already exists: Group.Register: %s", err)
	}
	if err := network.RegisterGroup(g.Name, owner); err != nil {
		return fmt.Errorf("could not register group: Group.Register: %s", err)
	}
	return nil
}
