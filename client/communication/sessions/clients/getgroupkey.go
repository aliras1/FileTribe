package clients

import (
	"ipfs-share/client/fs/caps"
	"math/rand"
	"sync"

	"github.com/golang/glog"
	"github.com/pkg/errors"

	comcommon "ipfs-share/client/communication/common"
	"ipfs-share/client/communication/sessions/common"
	"ipfs-share/client/fs"
	"ipfs-share/client/interfaces"
	"ipfs-share/collections"
	"ipfs-share/crypto"
)

type GetGroupKeyOnSuccessCallback func(cap *caps.GroupAccessCap)

type GetGroupKeySessionClient struct {
	sessionId collections.IIdentifier
	state     uint8
	contact   *comcommon.Contact
	groupId [32]byte

	storage *fs.Storage
	user interfaces.IUser
	closedChan chan common.ISession

	lock sync.RWMutex
	stop chan bool
	error error
	onSuccessCallback GetGroupKeyOnSuccessCallback
}

func (session *GetGroupKeySessionClient) Error() error {
	return session.error
}

func (session *GetGroupKeySessionClient) close() {
	session.state = common.EndOfSession
	session.closedChan <- session
}

func (session *GetGroupKeySessionClient) Abort() {
	session.lock.Lock()
	defer session.lock.Unlock()

	if !session.IsAlive() {
		return
	}

	session.close()
}

func (session *GetGroupKeySessionClient) State() uint8 {
	session.lock.RLock()
	defer session.lock.RUnlock()

	return session.state
}

func (session *GetGroupKeySessionClient) Id() collections.IIdentifier {
	return session.sessionId
}

func (session *GetGroupKeySessionClient) IsAlive() bool {
	session.lock.RLock()
	defer session.lock.RUnlock()

	return session.state == common.EndOfSession
}

func (session *GetGroupKeySessionClient) Run() {
	session.NextState(nil, nil)
}

func (session *GetGroupKeySessionClient) NextState(contact *comcommon.Contact, data []byte) {
	session.lock.Lock()
	defer session.lock.Unlock()

	switch session.state {
	case 0:
		{
			msg, err := comcommon.NewMessage(
				session.user.Address(),
				comcommon.GetGroupKey,
				session.sessionId.Data().(uint32),
				session.groupId[:],
				session.user.Signer(),
			)
			if err != nil {
				session.error = errors.Wrap(err, "could not create message")
				session.close()
				return
			}

			encMsg, err := msg.Encode()
			if err != nil {
				session.error = errors.Wrap(err, "could not encode message")
				session.close()
				return
			}

			if err := session.contact.Send(encMsg); err != nil {
				session.error = errors.Wrap(err, "could not send message")
				session.close()
				return
			}

			session.state = 1

			return
		}
		// Got the challenge
	case 1:
		{
			signer := session.user.Signer()
			sig, err := signer.Sign(data)
			if err != nil {
				session.error = errors.Wrap(err, "could not sign challange")
				session.close()
			}

			msg, err := comcommon.NewMessage(
				session.user.Address(),
				comcommon.GetGroupKey,
				session.sessionId.Data().(uint32),
				sig,
				session.user.Signer(),
			)
			if err != nil {
				session.error = errors.Wrap(err, "could not create message")
				session.close()
				return
			}

			encMsg, err := msg.Encode()
			if err != nil {
				session.error = errors.Wrap(err, "could not encode Message")
				session.close()
				return
			}

			if err := session.contact.Send(encMsg); err != nil {
				session.error = errors.Wrap(err, "could not send message")
				session.close()
				return
			}

			session.state = 2

			return
		}
	case 2:
		{
			key, err := crypto.DecodeSymmetricKey(data)
			if err != nil {
				session.error = errors.Wrap(err, "could not decode group key")
				session.close()
				return
			}

			groupCap := &caps.GroupAccessCap{
				GroupId: session.groupId,
				Boxer:   *key,
			}

			if err := session.storage.SaveGroupAccessCap(groupCap); err != nil {
				session.error = errors.Wrap(err, "could not save group access cap")
				session.close()
				return
			}

			session.onSuccessCallback(groupCap)
			session.close()
			return
		}

	default:
		{
			glog.Errorf("session ended")
		}
	}
}

func NewGetGroupKeySessionClient(
	groupId [32]byte,
	contact *comcommon.Contact,
	user interfaces.IUser,
	storage *fs.Storage,
	closedChan chan common.ISession,
	onSuccess GetGroupKeyOnSuccessCallback,
) *GetGroupKeySessionClient {

	sessionId := rand.Uint32()

	return &GetGroupKeySessionClient{
		sessionId: collections.NewUint32Id(sessionId),
		groupId:   groupId,
		contact:   contact,
		state:     0,
		user:       user,
		storage: storage,
		closedChan: closedChan,
		stop: make(chan bool),
		onSuccessCallback: onSuccess,
	}
}