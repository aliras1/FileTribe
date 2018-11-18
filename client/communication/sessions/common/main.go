package common

import (
	"math"

	"ipfs-share/client/communication/common"
	"ipfs-share/client/interfaces"
	. "ipfs-share/collections"
)


const EndOfSession = math.MaxUint8

type ISession interface {
	Id() IIdentifier
	IsAlive() bool
	Abort()
	NextState(contact *common.Contact, data []byte)
	State() uint8
	Run()
	Error() error
}

type SessionClosedCallback func(session ISession)

type GetGroupCallback func(id [32]byte) (interfaces.IGroup)

