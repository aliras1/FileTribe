package client

import (
	"ipfs-share/crypto"
	nw "ipfs-share/network"
)

type ICommand interface {
	Execute(network *nw.Network) error
}

type InviteCMD struct {
	User      string
	Boxer     *crypto.BoxingKeyPair
	NewMember string
	GroupCtx  *GroupContext
}

func (i *InviteCMD) Execute(network *nw.Network) error {
	i.GroupCtx.Members = append(i.GroupCtx.Members, i.NewMember)
	err := i.GroupCtx.Storage.CreateGroupAccessCAPForUser(i.NewMember, i.GroupCtx.Group.GroupName, i.GroupCtx.Group.Boxer, i.Boxer, network)
	if err != nil {
		return err
	}
	err = i.GroupCtx.Save()
	if err != nil {
		return err
	}
	return network.SendMessage(i.User, i.NewMember, "GROUP INVITE", i.GroupCtx.Group.GroupName+".json")
}
