package common

import (
	"math"

	"ipfs-share/client/communication/common"
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
