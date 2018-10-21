package client

import (
	"fmt"
	"strings"
)

type SignedBy struct {
	Username  string `json:"username"`
	Signature []byte `json:"signature"`
}

type RawOperation struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type Transaction struct {
	PrevState []byte       `json:"prev_state"`
	State     []byte       `json:"state"`
	Operation RawOperation `json:"operation"`
	SignedBy  []SignedBy   `json:"signed_by"`
}

func (t *Transaction) Bytes() []byte {
	var transactionBytes []byte
	transactionBytes = append(transactionBytes, t.PrevState...)
	transactionBytes = append(transactionBytes, t.State...)
	transactionBytes = append(transactionBytes, []byte(t.Operation.Type)...)
	transactionBytes = append(transactionBytes, []byte(t.Operation.Data)...)
	return transactionBytes
}

type IOperation interface {
	Execute(ctx *GroupContext) error
	RawOperation() RawOperation
	Validate(state []byte, groupCtx *GroupContext) error
}

func NewOperation(operation *RawOperation) (IOperation, error) {
	switch operation.Type {
	case "INVITE":
		args := strings.Split(operation.Data, " ")
		if len(args) < 2 {
			return nil, fmt.Errorf("invalid #args in operation 'INVITE' data: NewOperation")
		}
		cmd := InviteOperation{
			From:      args[0],
			NewMember: args[1],
		}
		return &cmd, nil
	case "SHARE":
		args := strings.Split(operation.Data, " ")
		if len(args) < 3 {
			return nil, fmt.Errorf("invalid #args in operation 'SHARE' data: NewOperation")
		}
		cmd := ShareFileOperation{
			Owner:    args[0],
			FileName: args[1],
			IPFSHash: args[2],
		}
		return &cmd, nil
	default:
		return nil, fmt.Errorf("invalid operation type: NewOperation")
	}
}

type InviteOperation struct {
	From      string
	NewMember string
}

func NewInviteOperation(from string, newMember string) IOperation {
	inviteOperation := InviteOperation{
		From:      from,
		NewMember: newMember,
	}
	return &inviteOperation
}

func (i *InviteOperation) Validate(state []byte, groupCtx *GroupContext) error {
	// newMembers := groupCtx.Members.Append(i.NewMember, groupCtx.Network)
	// newState := groupCtx.CalculateState(newMembers, groupCtx.Repo)
	// if !bytes.Equal(newState[:], state) {
	// 	return fmt.Errorf("invalid new state in transaction proposal: GroupConnection.validateTransaction")
	// }
	return nil
}

func (i *InviteOperation) RawOperation() RawOperation {
	rawOperation := RawOperation{
		Type: "INVITE",
		Data: i.From + " " + i.NewMember,
	}
	return rawOperation
}

func (i *InviteOperation) Execute(groupCtx *GroupContext) error {
	// glog.Infof("User '%s' executing invite cmd...\n", groupCtx.User.Name)

	// groupCtx.Members = groupCtx.Members.Append(i.NewMember, groupCtx.Network)
	// if err := groupCtx.Storage.CreateGroupAccessCAPForUser(
	// 	i.NewMember,
	// 	groupCtx.Group.Name,
	// 	groupCtx.Group.Boxer,
	// 	&groupCtx.User.Boxer,
	// 	groupCtx.Network,
	// ); err != nil {
	// 	return fmt.Errorf("could not create ga cap for user '%s': InviteOperation.Execute: %s", i.NewMember, err)
	// }
	// if err := groupCtx.Storage.PublishPublicDir(groupCtx.Ipfs); err != nil {
	// 	return fmt.Errorf("could not publish public dir: InviteOperation.Execute: %s", err)
	// }
	// // the proposer invites the new member
	// if strings.Compare(i.From, groupCtx.User.Name) == 0 {
	// 	glog.Info("\t--> Invite proposer sending chain message...")
	// 	if err := groupCtx.Network.DialP2PConn(
	// 		i.From,
	// 		i.NewMember,
	// 		"GROUP INVITE",
	// 		groupCtx.Group.Name+".json",
	// 	); err != nil {
	// 		return fmt.Errorf("user '%s' could not send message to user '%s': InviteOperation.Execute: %s", i.From, i.NewMember, err)
	// 	}
	// }
	return nil
}

type ShareFileOperation struct {
	Owner    string
	FileName string
	IPFSHash string
}

func NewShareFileOperation(owner, fileName, ipfsHash string) IOperation {
	shareFileOperation := &ShareFileOperation{
		Owner:    owner,
		FileName: fileName,
		IPFSHash: ipfsHash,
	}
	return shareFileOperation
}

func (share *ShareFileOperation) RawOperation() RawOperation {
	rawOperation := RawOperation{
		Type: "SHARE",
		Data: share.Owner + " " + share.FileName + " " + share.IPFSHash,
	}
	return rawOperation
}

func (share *ShareFileOperation) Validate(state []byte, groupCtx *GroupContext) error {
	//return fmt.Errorf("not implemented: ShareFileOperation.Validate")
	return nil // nothing to validate
}

func (share *ShareFileOperation) Execute(groupCtx *GroupContext) error {
	//file := &FileGroup{
	//	Name:     share.FileName,
	//	IPFSHash: share.IPFSHash,
	//}
	//if err := groupCtx.Storage.DownloadGroupFile(file,
	//	groupCtx.Group.Name,
	//	&groupCtx.Group.Boxer,
	//	groupCtx.Ipfs,
	//); err != nil {
	//	return fmt.Errorf("could not download group share file: ShareFileOperation: %s", err)
	//}
	//
	//groupCtx.Repo.Append(file)
	return nil
}
