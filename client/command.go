package client

import (
	"strings"
)

type ICommand interface {
	Execute(ctx *GroupContext) error
}

type InviteCMD struct {
	From      string
	NewMember string
}

func (i *InviteCMD) Execute(ctx *GroupContext) error {
	err := ctx.Storage.CreateGroupAccessCAPForUser(i.NewMember, ctx.Group.GroupName, ctx.Group.Boxer, &ctx.User.Boxer, ctx.Network)
	if err != nil {
		return err
	}
	ctx.Members = ctx.Members.Append(i.NewMember, ctx.Network)
	err = ctx.Save()
	if err != nil {
		return err
	}
	if strings.Compare(i.From, ctx.User.Username) == 0 {
		return ctx.Network.SendMessage(i.From, i.NewMember, "GROUP INVITE", ctx.Group.GroupName+".json")
	}
	return nil
}

func CMDFromProposal(proposal Proposal) ICommand {
	switch proposal.CMD {
	case "invite":
		if len(proposal.Args) < 1 {
			return nil
		}
		cmd := InviteCMD{proposal.From, proposal.Args[0]}
		return &cmd
	default:
		return nil
	}
}
