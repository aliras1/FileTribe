package common

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	"ipfs-share/client/fs"
	"ipfs-share/client/interfaces"
	"ipfs-share/collections"
	"ipfs-share/crypto"
)

type OnGetGroupKeySuccessCallback func(address ethcommon.Address, boxer tribecrypto.SymmetricKey)

type OnClientSuccessCallback func(args []interface{})

type OnServerSuccessCallback func(args []interface{}, groupId collections.IIdentifier)

type SessionClosedCallback func(session ISession)

type GetGroupDataCallback func(addr ethcommon.Address) (interfaces.IGroup, *fs.GroupRepo)

type Broadcast func(msg []byte) error

type CtxCallback interface {
	IsMember(group ethcommon.Address, account ethcommon.Address) error

	Boxer(group ethcommon.Address) (tribecrypto.SymmetricKey, error)

	ProposedBoxer(group ethcommon.Address, proposer ethcommon.Address) (tribecrypto.SymmetricKey, error)
}
