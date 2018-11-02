package client

import (
	"math/rand"
	"math"
	"ipfs-share/crypto"
	"sync"
	"github.com/golang/glog"

	. "ipfs-share/collections"
	"github.com/pkg/errors"
)


const EndOfSession = math.MaxUint8

type ISession interface {
	Id() IIdentifier
	IsAlive() bool
	Abort()
	NextState(contact *Contact, data []byte)
	State() uint8
	Run()
	Error() error
}

func NewServerSession(msg *Message, contact *Contact, user IUser, groups *ConcurrentCollection, closedChan chan ISession) (ISession, error) {
	switch msg.Type {
	case GetGroupKey:
		{
			return NewGetGroupKeySessionServer(msg, contact, user, groups, closedChan)
		}

	default:
		return nil, errors.New("invalid message type")
	}
}

type GetGroupKeySessionServer struct {
	sessionId IIdentifier
	state     uint8
	contact   *Contact
	user IUser
	group IGroup
	challenge [32]byte
	closedChan chan ISession
	lock sync.RWMutex
	stop chan bool
	error error
}

func (session *GetGroupKeySessionServer) Error() error {
	return session.error
}

func (session *GetGroupKeySessionServer) close() {
	session.state = EndOfSession
	session.closedChan <- session
}

func (session *GetGroupKeySessionServer) State() uint8 {
	session.lock.RLock()
	defer session.lock.RUnlock()

	return session.state
}

func (session *GetGroupKeySessionServer) Id() IIdentifier {
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

func (session *GetGroupKeySessionServer) NextState(contact *Contact, data []byte) {
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

			msg, err := NewMessage(
				session.user.Address(),
				GetGroupKey,
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

			msg, err := NewMessage(
				session.user.Address(),
				GetGroupKey,
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

func NewGetGroupKeySessionServer(msg *Message, contact *Contact, user IUser, groups *ConcurrentCollection, closedChan chan ISession) (*GetGroupKeySessionServer, error) {
	var challenge [32]byte
	if _, err := rand.Read(challenge[:]); err != nil {
		return nil, errors.Wrap(err, "could not read rand")
	}

	var groupId [32]byte
	copy(groupId[:], msg.Payload)

	groupInt := groups.Get(NewBytesId(groupId))
	if groupInt == nil {
		return nil, errors.New("no group found")
	}

	return &GetGroupKeySessionServer{
		sessionId: NewUint32Id(msg.SessionId),
		contact:   contact,
		user: user,
		group: groupInt.(*GroupContext).Group,
		closedChan: closedChan,
		state:     0,
		challenge: challenge,
	}, nil
}

// ---------------------

type GetGroupKeyOnSuccessCallback func(cap *GroupAccessCap)

type GetGroupKeySessionClient struct {
	sessionId IIdentifier
	state     uint8
	contact   *Contact
	groupId [32]byte

	storage *Storage
	user IUser
	closedChan chan ISession

	lock sync.RWMutex
	stop chan bool
	error error
	onSuccessCallback GetGroupKeyOnSuccessCallback
}

func (session *GetGroupKeySessionClient) Error() error {
	return session.error
}

func (session *GetGroupKeySessionClient) close() {
	session.state = EndOfSession
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

func (session *GetGroupKeySessionClient) Id() IIdentifier {
	return session.sessionId
}

func (session *GetGroupKeySessionClient) IsAlive() bool {
	session.lock.RLock()
	defer session.lock.RUnlock()

	return session.state == EndOfSession
}

func (session *GetGroupKeySessionClient) Run() {
	session.NextState(nil, nil)
}

func (session *GetGroupKeySessionClient) NextState(contact *Contact, data []byte) {
	session.lock.Lock()
	defer session.lock.Unlock()

	switch session.state {
	case 0:
		{
			msg, err := NewMessage(
				session.user.Address(),
				GetGroupKey,
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

			msg, err := NewMessage(
				session.user.Address(),
				GetGroupKey,
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

			groupCap := &GroupAccessCap{
				GroupId: session.groupId,
				Boxer:   *key,
			}

			if err := groupCap.Save(session.storage); err != nil {
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
	contact *Contact,
	user IUser,
	storage *Storage,
	closedChan chan ISession,
	onSuccess GetGroupKeyOnSuccessCallback,
	) *GetGroupKeySessionClient {

	sessionId := rand.Uint32()

	return &GetGroupKeySessionClient{
		sessionId: NewUint32Id(sessionId),
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
