package servers

import (
	"crypto/rand"
	"math"
	"sync"

	"github.com/golang/glog"
	"github.com/pkg/errors"

	comcommon"ipfs-share/client/communication/common"
	"ipfs-share/client/communication/sessions/common"
	"ipfs-share/client/interfaces"
	"ipfs-share/collections"
)

type GetGroupKeySessionServer struct {
	sessionId       collections.IIdentifier
	state           uint8
	contact         *comcommon.Contact
	user            interfaces.IUser
	group           interfaces.IGroup
	challenge       [32]byte
	onSessionClosed common.SessionClosedCallback
	lock            sync.RWMutex
	stop            chan bool
	error           error
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

func (session *GetGroupKeySessionServer) Id() collections.IIdentifier {
	return session.sessionId
}

func (session *GetGroupKeySessionServer) Abort() {
	session.lock.Lock()
	defer session.lock.Unlock()

	if !session.IsAlive() {
		return
	}

	session.close()
}

func (session *GetGroupKeySessionServer) IsAlive() bool {
	session.lock.RLock()
	defer session.lock.RUnlock()

	return session.state == math.MaxUint8
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
			if !session.group.IsMember(session.contact.Address) {
				session.close()
				session.error = errors.New("non group member requested group key")
				return
			}

			msg, err := comcommon.NewMessage(
				session.user.Address(),
				comcommon.GetGroupKey,
				session.sessionId.Data().(uint32),
				session.challenge[:],
				session.user.Signer(),
			)
			if err != nil {
				session.error = errors.New("could not create message")
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
	case 1:
		{
			if !session.contact.VerifySignature(session.challenge[:], data) {
				session.error = errors.New("invalid signature")
				session.close()
				return
			}

			boxer := session.group.Boxer()
			key, err := boxer.Encode()
			if err != nil {
				session.error = errors.Wrap(err, "could not marshal group key")
				session.close()
				return
			}

			msg, err := comcommon.NewMessage(
				session.user.Address(),
				comcommon.GetGroupKey,
				session.sessionId.Data().(uint32),
				key,
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

			session.close()
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
	user interfaces.IUser,
	getGroupData common.GetGroupDataCallback,
	onSessionClosed common.SessionClosedCallback,
) (*GetGroupKeySessionServer, error) {

	var challenge [32]byte
	if _, err := rand.Read(challenge[:]); err != nil {
		return nil, errors.Wrap(err, "could not read rand")
	}

	var groupId [32]byte
	copy(groupId[:], msg.Payload)

	group, _ := getGroupData(groupId)
	if group == nil {
		return nil, errors.New("no group found")
	}

	return &GetGroupKeySessionServer{
		sessionId:       collections.NewUint32Id(msg.SessionId),
		contact:         contact,
		user:            user,
		group:           group,
		onSessionClosed: onSessionClosed,
		state:           0,
		challenge:       challenge,
	}, nil
}