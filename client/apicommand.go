package client

import (
	"strings"
	"fmt"
	"io/ioutil"
	"ipfs-share/ipfs"
	nw "ipfs-share/networketh"
)

type ICommand interface {
	Execute(ctx *UserContext, network *nw.Network, ipfs *ipfs.IPFS) (*UserContext, string, error)
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
	case "ptpshare":
		if len(args) < 3 {
			return nil, fmt.Errorf("invalid # args in 'ptpshare' command")
		}
		cmd := CMDPTPAddAndShareFile{
			ShareWith: args[1],
			FilePath: args[2],
		}
		return &cmd, nil
	case "gshare":
		if len(args) < 2 {
			return nil, fmt.Errorf("invalid # args in 'gshare' command")
		}
		cmd := CMDGroupAddAndShareFile{
			GroupName: args[1],
			Path: args[2],
		}
		return &cmd, nil
	case "ls":
		cmd := CMDList{}
		return &cmd, nil
	case "cat":
		if len(args) < 2 {
			return nil, fmt.Errorf("invalid # args in 'cat' command")
		}
		cmd := CMDCat{
			Path: args[1],
		}
		return &cmd, nil
	case "signout":
		cmd := CMDSignOut{}
		return &cmd, nil
	default:
		return nil, fmt.Errorf("invalid command")
	}
}

type CMDSignUp struct {
	Username string
	Password string
}

func (cmd *CMDSignUp) Execute(ctx *UserContext, network *nw.Network, ipfs *ipfs.IPFS) (*UserContext, string, error) {
	if ctx != nil {
		return nil, "", fmt.Errorf("an active user context already running: CMDSignup.Execute")
	}
	uc, err := NewUserContextFromSignUp(cmd.Username, cmd.Password, cmd.Username, network, ipfs)
	if err != nil {
		return nil, "", fmt.Errorf("could not sign up: CMDSignUp.Execute: %s", err)
	}
	return uc, "", nil
}

type CMDSignIn struct {
	Username string
	Password string
}

func (cmd *CMDSignIn) Execute(ctx *UserContext, network *nw.Network, ipfs *ipfs.IPFS) (*UserContext, string, error) {
	if ctx != nil {
		return nil, "", fmt.Errorf("an active user context already running: CMDSignIn.Execute")
	}
	uc, err := NewUserContextFromSignIn(cmd.Username, cmd.Password, cmd.Username, network, ipfs)
	if err != nil {
		return nil, "", fmt.Errorf("could not sign in: CMDSignIn.Execute: %s", err)
	}
	return uc, "", nil
}

type CMDSignOut struct {

}

func (cmd *CMDSignOut) Execute(ctx *UserContext, network *nw.Network, ipfs *ipfs.IPFS) (*UserContext, string, error) {
	if ctx == nil {
		return nil, "", fmt.Errorf("no active user context running: CMDSignOut.Execute")
	}
	ctx.SignOut()
	return nil, "", nil
}

type CMDCreateGroup struct {
	GroupName string
}

func (cmd *CMDCreateGroup) Execute(ctx *UserContext, network *nw.Network, ipfs *ipfs.IPFS) (*UserContext, string, error) {
	if ctx == nil {
		return nil, "", fmt.Errorf("no active user context CMDCreateGroup.Execute")
	}
	if err := ctx.CreateGroup(cmd.GroupName); err != nil {
		return ctx, "", fmt.Errorf("could not create group: CmdCreateGroup.Execute: %s", err)
	}
	return ctx, "", nil
}

type CMDGroupInvite struct {
	GroupName string
	NewMember string
}

func (cmd *CMDGroupInvite) Execute(ctx *UserContext, network *nw.Network, ipfs *ipfs.IPFS) (*UserContext, string, error) {
	if ctx == nil {
		return nil, "", fmt.Errorf("no active user context CMDGroupInvite.Execute")
	}
	for _, group := range ctx.Groups {
		if strings.Compare(group.Group.Name, cmd.GroupName) == 0 {
			if err := group.Invite(cmd.NewMember); err != nil {
				return ctx, "", fmt.Errorf("could not invite user '%s' into group '%s': CMDGroupInvite.Execute: %s", cmd.NewMember, cmd.GroupName, err)
			}
			return ctx, "", nil
		}
	}
	return ctx, "", fmt.Errorf("no group '%s' found in group repo: CMDGroupIOnvite.Execute", cmd.GroupName)
}

type CMDPTPAddAndShareFile struct {
	FilePath string
	ShareWith string
}

func (cmd *CMDPTPAddAndShareFile) Execute(ctx *UserContext, network *nw.Network, ipfs *ipfs.IPFS) (*UserContext, string, error) {
	if ctx == nil {
		return nil, "", fmt.Errorf("no active user context CMDPTPAddAndShareFile.Execute")
	}
	if err := ctx.AddAndShareFile(cmd.FilePath, []string{cmd.ShareWith}); err != nil {
		return ctx, "", fmt.Errorf("could not add and share file '%s' with user '%s': CMDPTPAddAndShareFile.Execute: %s", cmd.FilePath, cmd.ShareWith, err)
	}
	return ctx, "", nil
}

type CMDGroupAddAndShareFile struct {
	GroupName string
	Path string
}

func (cmd *CMDGroupAddAndShareFile) Execute(ctx *UserContext, network *nw.Network, ipfs *ipfs.IPFS) (*UserContext, string, error) {
	if ctx == nil {
		return nil, "", fmt.Errorf("no active user context CMDPTPAddAndShareFile.Execute")
	}
	for _, groupCtx := range ctx.Groups {
		if strings.Compare(cmd.GroupName, groupCtx.Group.Name) == 0 {
			if err := groupCtx.AddAndShareFile(cmd.Path); err != nil {
				return ctx, "", fmt.Errorf("could not group add and share file: CMDGroupAddAndShareFile: %s", err)
			}
			return ctx, "", nil
		}
	}
	return ctx, "", fmt.Errorf("no group named '%s', exists: CMDGroupAddAndShareFile", cmd.GroupName)
}

type CMDList struct {

}

func (cmd *CMDList) Execute(ctx *UserContext, network *nw.Network, ipfs *ipfs.IPFS) (*UserContext, string, error) {
	if ctx == nil {
		return nil, "", fmt.Errorf("no active user context CMDList.Execute")
	}
	return ctx, ctx.List(), nil
}

type CMDCat struct {
	Path string
}

func (cmd *CMDCat) Execute(ctx *UserContext, network *nw.Network, ipfs *ipfs.IPFS) (*UserContext, string, error) {
	if ctx == nil {
		return nil, "", fmt.Errorf("no active user context CMDList.Execute")
	}
	filePath := ctx.Storage.GetUserFilesPath() + "/" + cmd.Path
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return ctx, "", fmt.Errorf("could not read file '%s': CMDCat.Execute", filePath)
	}
	fmt.Printf(string(fileBytes) + "\n")
	return ctx, "", nil
}
