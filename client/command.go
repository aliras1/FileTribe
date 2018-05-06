package client

import (
	"fmt"
	"strings"
	"log"
)

type ICommand interface {
	Execute(ctx *GroupContext) error
}

type InviteCMD struct {
	From      string
	NewMember string
}

func NewCommand(operation *Operation) (ICommand, error) {
	switch operation.Type {
	case "INVITE":
		args := strings.Split(operation.Data, " ")
		if len(args) < 2 {
			return nil, fmt.Errorf("invalid #args in operation data: NewCommand")
		}
		cmd := InviteCMD{
			From: args[0],
			NewMember: args[1],
		}
		return &cmd, nil
	default:
		return nil, fmt.Errorf("invalid operation type: NewCommand")
	}
}

func (i *InviteCMD) Execute(groupCtx *GroupContext) error {
	log.Printf("[*] %s executes invite cmd...", groupCtx.User.Name)
	groupCtx.Members = groupCtx.Members.Append(i.NewMember, groupCtx.Network)
	if err := groupCtx.Storage.CreateGroupAccessCAPForUser(
		i.NewMember,
		groupCtx.Group.Name,
		groupCtx.Group.Boxer,
		&groupCtx.User.Boxer,
		groupCtx.Network,
	); err != nil {
		return fmt.Errorf("could not create ga cap for user '%s': InviteCMD.Execute: %s", i.NewMember, err)
	}
	if err := groupCtx.Storage.PublishPublicDir(groupCtx.IPFS); err != nil {
		return fmt.Errorf("could not publish public dir: InviteCMD.Execute: %s", err)
	}
	// the proposer invites the new member
	if strings.Compare(i.From, groupCtx.User.Name) == 0 {
		if err := groupCtx.Network.SendMessage(
			i.From,
			i.NewMember,
			"GROUP INVITE",
			groupCtx.Group.Name+".json",
		); err != nil {
			return fmt.Errorf("user '%s'could not send message to user '%s': InviteCMD.Execute: %s", i.From, i.NewMember, err)
		}
	}
	return nil
}