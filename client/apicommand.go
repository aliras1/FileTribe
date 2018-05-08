package client

import (
	nw "ipfs-share/network"
	"strings"
	"fmt"
	"ipfs-share/ipfs"
)

type ICommand interface {
	Execute(ctx *UserContext, network *nw.Network, ipfs *ipfs.IPFS) (*UserContext, error)
}

type CMDSignUp struct {
	Username string
	Password string
}

func (cmd *CMDSignUp) Execute(ctx *UserContext, network *nw.Network, ipfs *ipfs.IPFS) (*UserContext, error) {
	if ctx != nil {
		return nil, fmt.Errorf("an active user context already running: CMDSignup.Execute")
	}
	uc, err := NewUserContextFromSignUp(cmd.Username, cmd.Password, cmd.Username, network, ipfs)
	if err != nil {
		return nil, fmt.Errorf("could not sign up: CMDSignUp.Execute: %s", err)
	}
	return uc, nil
}

type CMDSignIn struct {
	Username string
	Password string
}

func (cmd *CMDSignIn) Execute(ctx *UserContext, network *nw.Network, ipfs *ipfs.IPFS) (*UserContext, error) {
	if ctx != nil {
		return nil, fmt.Errorf("an active user context already running: CMDSignIn.Execute")
	}
	uc, err := NewUserContextFromSignIn(cmd.Username, cmd.Password, cmd.Username, network, ipfs)
	if err != nil {
		return nil, fmt.Errorf("could not sign up: CMDSignIn.Execute: %s", err)
	}
	return uc, nil
}

type CMDCreateGroup struct {
	GroupName string
}

func (cmd *CMDCreateGroup) Execute(ctx *UserContext, network *nw.Network, ipfs *ipfs.IPFS) (*UserContext, error) {
	if ctx == nil {
		return nil, fmt.Errorf("no active user context CMDCreateGroup.Execute")
	}
	if err := ctx.CreateGroup(cmd.GroupName); err != nil {
		return ctx, fmt.Errorf("could not create group: CmdCreateGroup.Execute: %s", err)
	}
	return ctx, nil
}

type CMDGroupInvite struct {
	GroupName string
	NewMember string
}

func (cmd *CMDGroupInvite) Execute(ctx *UserContext, network *nw.Network, ipfs *ipfs.IPFS) (*UserContext, error) {
	if ctx == nil {
		return nil, fmt.Errorf("no active user context CMDGroupInvite.Execute")
	}
	for _, group := range ctx.Groups {
		if strings.Compare(group.Group.Name, cmd.GroupName) == 0 {
			if err := group.Invite(cmd.NewMember); err != nil {
				return ctx, fmt.Errorf("could not invite user '%s' into group '%s': CMDGroupInvite.Execute: %s", cmd.NewMember, cmd.GroupName, err)
			}
			return ctx, nil
		}
	}
	return ctx, fmt.Errorf("no group '%s' found in group repo: CMDGroupIOnvite.Execute", cmd.GroupName)
}

func NewCommand(raw string) (ICommand, error) {
	args := strings.Split(raw, " ")
	if len(args) < 1 {
		return nil, fmt.Errorf("empty command")
	}
	switch args[0] {
	case "signup":
		if len(args) < 3 {
			return nil, fmt.Errorf("invalid # args in 'signup' command")
		}
		cmd := CMDSignUp{
			Username: args[1],
			Password: args[2],
		}
		return &cmd, nil
	case "signin":
		if len(args) < 3 {
			return nil, fmt.Errorf("invalid # args in 'signin' command")
		}
		cmd := CMDSignIn{
			Username: args[1],
			Password: args[2],
		}
		return &cmd, nil
	case "cgroup":
		if len(args) < 2 {
			return nil, fmt.Errorf("invalid # args in 'cgroup' command")
		}
		cmd := CMDCreateGroup{
			GroupName: args[1],
		}
		return &cmd, nil
	case "igroup":
		if len(args) < 3 {
			return nil, fmt.Errorf("invalid # args in 'igroup' command")
		}
		cmd := CMDGroupInvite{
			GroupName: args[1],
			NewMember: args[2],
		}
		return &cmd, nil
	default:
		return nil, fmt.Errorf("invalid command")
	}
}