package clients

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"ipfs-share/client/fs/caps"
	"ipfs-share/crypto"
	"math/rand"
	"sync"

	comcommon "ipfs-share/client/communication/common"
	"ipfs-share/client/communication/sessions/common"
)

type GetGroupKeySessionClient struct {
	sessionId uint32
	msgType   comcommon.MessageType
	state     uint8
	receiver  *comcommon.Contact
	groupAddr ethcommon.Address

	sender          ethcommon.Address
	onSessionClosed common.SessionClosedCallback
	signer          *tribecrypto.Signer

	lock 				sync.RWMutex
	stop			    chan bool
	error				error
	onSuccessCallback   common.OnGetGroupKeySuccessCallback
}

func (session *GetGroupKeySessionClient) Error() error {
	return session.error
}

func (session *GetGroupKeySessionClient) close() {
	session.state = common.EndOfSession
	session.onSessionClosed(session)
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

func (session *GetGroupKeySessionClient) Id() uint32 {
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
				session.sender,
				comcommon.GetGroupKey,
				session.sessionId,
				session.groupAddr.Bytes(),
				session.signer,
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

			if err := session.receiver.Send(encMsg); err != nil {
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
			sig, err := session.signer.Sign(data)
			if err != nil {
				session.error = errors.Wrap(err, "could not sign challange")
				session.close()
			}

			msg, err := comcommon.NewMessage(
				session.sender,
				comcommon.GetGroupKey,
				session.sessionId,
				sig,
				session.signer,
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

			if err := session.receiver.Send(encMsg); err != nil {
				session.error = errors.Wrap(err, "could not send message")
				session.close()
				return
			}

			session.state = 2

			return
		}
	case 2:
		{
			key, err := tribecrypto.DecodeSymmetricKey(data)
			if err != nil {
				session.error = errors.Wrap(err, "could not decode group key")
				session.close()
				return
			}

			groupCap := &caps.GroupAccessCap{
				Address: session.groupAddr,
				Boxer:   *key,
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
	msgType comcommon.MessageType,
	groupAddr ethcommon.Address,
	contact *comcommon.Contact,
	sender ethcommon.Address,
	signer *tribecrypto.Signer,
	onSessionClosed common.SessionClosedCallback,
	onSuccess common.OnGetGroupKeySuccessCallback,
) *GetGroupKeySessionClient {

	return &GetGroupKeySessionClient{
		sessionId:         rand.Uint32(),
		msgType:		   msgType,
		groupAddr:         groupAddr,
		receiver:          contact,
		state:             0,
		sender:            sender,
		signer:			   signer,
		onSessionClosed:   onSessionClosed,
		stop:              make(chan bool),
		onSuccessCallback: onSuccess,
	}
}