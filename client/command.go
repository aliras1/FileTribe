package client

import (
	"fmt"
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
	fmt.Println(ctx.User.Username + " executing invite cmd...")
	err := ctx.Storage.CreateGroupAccessCAPForUser(i.NewMember, ctx.Group.GroupName, ctx.Group.Boxer, &ctx.User.Boxer, ctx.Network)
	if err != nil {
		return err
	}
	fmt.Println(ctx.User.Username + " 1...")
	ctx.Members = ctx.Members.Append(i.NewMember, ctx.Network)
	fmt.Println(ctx.User.Username + " 1.5...")
	if err = ctx.Save(); err != nil {
		fmt.Println("==> ctx Save")
		fmt.Println(err)
		return err
	}
	fmt.Println(ctx.User.Username + " 2...")
	fmt.Println(ctx.User.Username + " - " + i.From)
	if strings.Compare(i.From, ctx.User.Username) == 0 {
		fmt.Println(ctx.User.Username + " sending invite msg...")
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
