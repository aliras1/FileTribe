package servers

import (
	"crypto/rand"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/glog"
	"ipfs-share/crypto"
	"sync"

	"github.com/pkg/errors"

	comcommon "ipfs-share/client/communication/common"
	"ipfs-share/client/communication/sessions/common"
)

type GetGroupKeySessionServer struct {
	sessionId       uint32
	state           uint8
	contact         *comcommon.Contact
	sender          ethcommon.Address
	group           ethcommon.Address
	callback        common.CtxCallback
	signer          *crypto.Signer
	challenge       [32]byte
	onSessionClosed common.SessionClosedCallback
	lock            sync.RWMutex
	stop            chan bool
	error           error
	keyType         comcommon.MessageType
}

func (session *GetGroupKeySessionServer) Error() error {
	return session.error
}

func (session *GetGroupKeySessionServer) close() {
	session.state = common.EndOfSession
	session.onSessionClosed(session)
}

func (session *GetGroupKeySessionServer) State() uint8 {
	session.lock.RLock()
	defer session.lock.RUnlock()

	return session.state
}

func (session *GetGroupKeySessionServer) Id() uint32 {
	return session.sessionId
}

func (session *GetGroupKeySessionServer) Abort() {
	if !session.isAlive() {
		return
	}

	session.close()
}

func (session *GetGroupKeySessionServer) IsAlive() bool {
	session.lock.RLock()
	defer session.lock.RUnlock()

	return session.isAlive()
}

func (session *GetGroupKeySessionServer) isAlive() bool {
	return session.state != common.EndOfSession
}

func (session *GetGroupKeySessionServer) Run() {
	session.NextState(nil, nil)
}

func (session *GetGroupKeySessionServer) NextState(contact *comcommon.Contact, data []byte) {
	session.lock.Lock()
	defer session.lock.Unlock()

	switch session.state {
	case 0:
		{
			if err := session.callback.IsMember(session.group, session.contact.AccAddr); err != nil {
				session.error = errors.Wrap(err, "could not verify group membership")
				session.Abort()
				return
			}

			msg, err := comcommon.NewMessage(
				session.sender,
				comcommon.GetGroupKey,
				session.sessionId,
				session.challenge[:],
				session.signer,
			)
			if err != nil {
				session.error = errors.New("could not create message")
				session.Abort()
				return
			}

			encMsg, err := msg.Encode()
			if err != nil {
				session.error = errors.Wrap(err, "could not encode message")
				session.Abort()
				return
			}

			if err := session.contact.Send(encMsg); err != nil {
				session.error = errors.Wrap(err, "could not send message")
				session.Abort()
				return
			}

			session.state = 1

			return
		}
	case 1:
		{
			if !session.contact.VerifySignature(session.challenge[:], data) {
				session.error = errors.New("invalid signature")
				session.Abort()
				return
			}

			var key []byte

			switch session.keyType {
			case comcommon.GetGroupKey:
				boxer, err := session.callback.Boxer(session.group)
				if err != nil {
					session.error = errors.Wrap(err, "could not get group boxer")
					session.Abort()
					return
				}

				data, err := boxer.Encode()
				if err != nil {
					session.error = errors.Wrap(err, "could not marshal group key")
					session.Abort()
					return
				}
				key = data
			case comcommon.GetProposedGroupKey:
				boxer, err := session.callback.ProposedBoxer(session.group)
				if err != nil {
					session.error = errors.Wrap(err, "could not get proposed group boxer")
					session.Abort()
					return
				}
				data, err := boxer.Encode()
				if err != nil {
					session.error = errors.Wrap(err, "could not marshal group key")
					session.Abort()
					return
				}
				key = data
			}


			msg, err := comcommon.NewMessage(
				session.sender,
				comcommon.GetGroupKey,
				session.sessionId,
				key,
				session.signer,
			)
			if err != nil {
				session.error = errors.Wrap(err, "could not create message")
				session.Abort()
				return
			}

			encMsg, err := msg.Encode()
			if err != nil {
				session.error = errors.Wrap(err, "could not encode message")
				session.Abort()
				return
			}

			if err := session.contact.Send(encMsg); err != nil {
				session.error = errors.Wrap(err, "could not send message")
				session.Abort()
				return
			}

			session.Abort()
		}

	default:
		{
			glog.Errorf("session error: called next state in invalid state")
		}
	}
}

func NewGetGroupKeySessionServer(
	msg *comcommon.Message,
	contact *comcommon.Contact,
	sender ethcommon.Address,
	signer *crypto.Signer,
	callback common.CtxCallback,
	onSessionClosed common.SessionClosedCallback,
) (*GetGroupKeySessionServer, error) {

	var challenge [32]byte
	if _, err := rand.Read(challenge[:]); err != nil {
		return nil, errors.Wrap(err, "could not read rand")
	}

	group := ethcommon.BytesToAddress(msg.Payload)

	return &GetGroupKeySessionServer{
		sessionId:       msg.SessionId,
		keyType:		 msg.Type,
		contact:         contact,
		callback:		 callback,
		sender:          sender,
		signer:			 signer,
		group:           group,
		onSessionClosed: onSessionClosed,
		state:           0,
		challenge:       challenge,
	}, nil
}