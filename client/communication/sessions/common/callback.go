package common

import (
	"ipfs-share/client/fs"
	"ipfs-share/client/fs/caps"
	"ipfs-share/client/interfaces"
	"ipfs-share/collections"
	"ipfs-share/network"
)

type OnGetGroupKeySuccessCallback func(cap *caps.GroupAccessCap)

type OnClientSuccessCallback func(args []interface{}, approvals []*network.Approval)

type OnServerSuccessCallback func(args []interface{}, groupId collections.IIdentifier)

type SessionClosedCallback func(session ISession)

type GetGroupDataCallback func(id [32]byte) (interfaces.IGroup, *fs.GroupRepo)

type Broadcast func(msg []byte) error

type CtxCallback interface {
	OnChangeGroupKeyServerSessionSuccess(args []interface{}, groupId collections.IIdentifier)

	GetGroupData(id [32]byte) (interfaces.IGroup, *fs.GroupRepo)
}
