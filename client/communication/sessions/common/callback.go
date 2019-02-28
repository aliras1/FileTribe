package common

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	"ipfs-share/client/fs"
	"ipfs-share/client/fs/caps"
	"ipfs-share/client/interfaces"
	"ipfs-share/collections"
	"ipfs-share/crypto"
)

type OnGetGroupKeySuccessCallback func(cap *caps.GroupAccessCap)

type OnClientSuccessCallback func(args []interface{})

type OnServerSuccessCallback func(args []interface{}, groupId collections.IIdentifier)

type SessionClosedCallback func(session ISession)

type GetGroupDataCallback func(addr ethcommon.Address) (interfaces.IGroup, *fs.GroupRepo)

type Broadcast func(msg []byte) error

type CtxCallback interface {
	IsMember(group ethcommon.Address, account ethcommon.Address) error

	Boxer(group ethcommon.Address) (tribecrypto.SymmetricKey, error)

	ProposedBoxer(group ethcommon.Address) (tribecrypto.SymmetricKey, error)
}
