package client

import (
	"crypto/rand"
	"errors"

	"ipfs-share/crypto"
	nw "ipfs-share/network"
)

type Group struct {
	GroupName string
	Signer    crypto.SigningKeyPair
}

func NewGroup(groupName string) *Group {
	var secretKeyBytes [32]byte
	rand.Read(secretKeyBytes[:])
	return NewGroupFromKey(groupName, &secretKeyBytes)
}

func NewGroupFromKey(groupName string, secretKeyBytes *[32]byte) *Group {
	pk, sk := crypto.Ed25519KeyPair(secretKeyBytes)
	return &Group{groupName, crypto.SigningKeyPair{pk, sk}}
}

func (g *Group) SignIn(network *nw.Network) error {
	publicKey, err := network.GetGroupSigningKey(g.GroupName)
	if err != nil {
		return err
	}
	if !g.Signer.PublicKey.Equals(&publicKey) {
		return errors.New("invalid group credentials")
	}
	return nil
}

func (g *Group) Register(network *nw.Network) error {
	registered, err := network.IsGroupRegistered(g.GroupName)
	if err != nil {
		return err
	}
	if registered {
		return errors.New("group name already exists")
	}
	err = network.RegisterGroup(g.GroupName, g.Signer.PublicKey)
	if err != nil {
		return err
	}
	return nil
}
