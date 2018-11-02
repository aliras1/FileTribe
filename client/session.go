package client

import (
	"math/rand"
	"math"
	"ipfs-share/crypto"
	"sync"
	"github.com/golang/glog"

	. "ipfs-share/collections"
)


const EndOfSession = math.MaxUint8

type ISession interface {
	Id() IIdentifier
	IsAlive() bool
	Abort()
	NextState(contact *Contact, data []byte)
	GetState() uint8
	Run()
}

func NewServerSession(msg *Message, contact *Contact, ctx *UserContext) ISession {
	switch msg.Type {
	case GetGroupKey:
		{
			return NewGetGroupKeySessionServer(msg, contact, ctx)
		}

	default:
		return nil
	}
}

type GetGroupKeySessionServer struct {
	sessionId IIdentifier
	state     uint8
	contact   *Contact
	groupId [32]byte
	challenge [32]byte
	ctx  *UserContext
	boxer crypto.SymmetricKey
	lock sync.RWMutex
	stop chan bool
}

func (session *GetGroupKeySessionServer) close() {
	session.state = EndOfSession
	session.ctx.P2P.SessionClosedChan <- session.Id()
}

func (session *GetGroupKeySessionServer) GetState() uint8 {
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
			groupId := NewBytesId(session.groupId)

			groupCtxInterface := session.ctx.Groups.Get(groupId)
			if groupCtxInterface == nil {
				session.close()
				glog.Errorf("%s: no groupContext found", session.ctx.User.Name)
				return
			}

			groupCtx := groupCtxInterface.(*GroupContext)

			if err := groupCtx.Update(); err != nil {
				session.close()
				glog.Errorf("could not update the state of group %s", err)
				return
			}

			if !groupCtx.Group.IsMember(session.contact.Address) {
				session.close()
				glog.Errorf("non group member %s asked for key", session.contact.Address.String())
				return
			}

			msg, err := NewMessage(
				session.ctx.User.Address,
				GetGroupKey,
				session.sessionId.Data().(uint32),
				session.challenge[:],
				session.ctx.User.Signer,
			)
			if err != nil {
				glog.Errorf("could not create Message: %s", err)
				return
			}

			encMsg, err := msg.Encode()
			if err != nil {
				glog.Errorf("could not encode Message: %s", err)
				return
			}

			if err := session.contact.Send(encMsg); err != nil {
				session.close()
				glog.Errorf("could not send message to %s: %s", session.contact.Address.String(), err)
				return
			}

			session.boxer = groupCtx.Group.Boxer()
			session.state = 1

			return
		}
	case 1:
		{
			if !session.contact.VerifySignature(session.challenge[:], data) {
				session.close()
				glog.Errorf("invalid signature")
				return
			}

			key, err := session.boxer.Encode()
			if err != nil {
				session.close()
				glog.Errorf("could not marshal group boxer: %s", err)
				return
			}

			msg, err := NewMessage(
				session.ctx.User.Address,
				GetGroupKey,
				session.sessionId.Data().(uint32),
				key,
				session.ctx.User.Signer,
			)
			if err != nil {
				glog.Errorf("could not create Message: %s", err)
				return
			}

			encMsg, err := msg.Encode()
			if err != nil {
				glog.Errorf("could not encode Message: %s", err)
				return
			}

			if err := session.contact.Send(encMsg); err != nil {
				session.close()
				glog.Errorf("could not send message: %s", err)
				return
			}

			session.close()
			glog.Errorf("session ended")
		}

	default:
		{
			glog.Errorf("session error: called next state in invalid state")
		}
	}
}

func NewGetGroupKeySessionServer(msg *Message, contact *Contact, ctx *UserContext) *GetGroupKeySessionServer {
	var challenge [32]byte
	if _, err := rand.Read(challenge[:]); err != nil {
		return nil
	}

	var groupId [32]byte
	copy(groupId[:], msg.Payload)

	return &GetGroupKeySessionServer{
		sessionId: NewUint32Id(msg.SessionId),
		contact:   contact,
		groupId: groupId,
		state:     0,
		challenge: challenge,
		ctx:  ctx,
	}
}

// ---------------------

type GetGroupKeySessionClient struct {
	sessionId IIdentifier
	state     uint8
	contact   *Contact
	ctx  *UserContext
	groupId [32]byte
	resultChannel chan crypto.SymmetricKey
	lock sync.RWMutex
	stop chan bool
}

func (session *GetGroupKeySessionClient) close() {
	session.state = EndOfSession
	session.ctx.P2P.SessionClosedChan <- session.Id()
}

func (session *GetGroupKeySessionClient) Abort() {
	session.lock.Lock()
	defer session.lock.Unlock()

	if !session.IsAlive() {
		return
	}

	session.close()
}

func (session *GetGroupKeySessionClient) GetState() uint8 {
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
				session.ctx.User.Address,
				GetGroupKey,
				session.sessionId.Data().(uint32),
				session.groupId[:],
				session.ctx.User.Signer,
			)
			if err != nil {
				glog.Errorf("could not create Message: %s", err)
				return
			}

			encMsg, err := msg.Encode()
			if err != nil {
				glog.Errorf("could not encode Message: %s", err)
				return
			}

			if err := session.contact.Send(encMsg); err != nil {
				session.close()
				glog.Errorf("could not send message to %s: %s", session.contact.Address.String(), err)
				return
			}

			session.state = 1

			return
		}
	// Got the challenge
	case 1:
		{
			sig, err := session.ctx.User.Signer.Sign(data)
			if err != nil {
				glog.Errorf("could not sign challenge: %s", err)
			}

			msg, err := NewMessage(
				session.ctx.User.Address,
				GetGroupKey,
				session.sessionId.Data().(uint32),
				sig,
				session.ctx.User.Signer,
			)
			if err != nil {
				glog.Errorf("could not create Message: %s", err)
				return
			}

			encMsg, err := msg.Encode()
			if err != nil {
				glog.Errorf("could not encode Message: %s", err)
				return
			}

			if err := session.contact.Send(encMsg); err != nil {
				session.close()
				glog.Errorf("could not send message to %s: %s", session.contact.Address.String(), err)
				return
			}

			session.state = 2

			return
		}
	case 2:
		{
			key, err := crypto.DecodeSymmetricKey(data)
			if err != nil {
				glog.Errorf("could not decode group key: %s", err)
				return
			}

			groupCap := &GroupAccessCap{
				GroupId: session.groupId,
				Boxer:   *key,
			}

			if err := groupCap.Save(session.ctx.Storage); err != nil {
				glog.Errorf("could not save group access groupCap: %s", err)
				return
			}

			groupCtx, err := NewGroupContextFromCAP(
				groupCap,
				session.ctx.User,
				session.ctx.P2P,
				session.ctx.AddressBook,
				session.ctx.Network,
				session.ctx.Ipfs,
				session.ctx.Storage,
			)
			if err != nil {
				glog.Errorf("could not create group context: %s", err)
				return
			}

			if err := groupCtx.Update(); err != nil {
				glog.Errorf("could not update group: %s", err)
				return
			}

			if err := session.ctx.Groups.Append(groupCtx); err != nil {
				glog.Warningf("could not append elem: %s", err)
			}

			session.close()

			return
		}

	default:
		{
			glog.Errorf("session ended")
		}
	}
}

func NewGetGroupKeySessionClient(groupId [32]byte, contact *Contact, ctx *UserContext) *GetGroupKeySessionClient {
	sessionId := rand.Uint32()

	return &GetGroupKeySessionClient{
		sessionId: NewUint32Id(sessionId),
		groupId:   groupId,
		contact:   contact,
		state:     0,
		ctx:       ctx,
		stop: make(chan bool),
	}
}
