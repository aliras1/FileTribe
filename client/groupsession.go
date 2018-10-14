package client

import (
	"sync"
	"github.com/golang/glog"
	. "ipfs-share/collections"
	"github.com/pkg/errors"
	"ipfs-share/crypto"
	"math/rand"
)

func NewServerGroupSession(groupId [32]byte, msg *GroupMessage, contact *Contact, ctx *GroupContext) ISession {
	switch msg.Type {
	case AddFile:
		{
			return NewAddFileServerGroupSession(groupId, msg, contact, ctx)
		}
	default:
		{
			return nil
		}
	}
}

type AddFileClientGroupSession struct {
	sessionId            IIdentifier
	newIpfsPath          string
	oldIpfsPath string
	approvals            *ConcurrentCollection
	digest               []byte // the original message digest which is signed by the group members
	state                uint8
	groupCtx *GroupContext
	lock                 sync.Mutex
	stop                 chan bool
	approvalsCountChan   chan int
}

func NewAddFileClientGroupSession(newIpfsPath, oldIpfsPath string, groupCtx *GroupContext) ISession {
	sessionId := rand.Uint32()
	hasher := crypto.NewKeccak256Hasher()

	session := &AddFileClientGroupSession{
		sessionId: NewUint32Id(sessionId),
		groupCtx: groupCtx,
		newIpfsPath: newIpfsPath,
		oldIpfsPath: oldIpfsPath,
		approvals: NewConcurrentCollection(),
		digest: hasher.Sum([]byte(newIpfsPath), []byte(oldIpfsPath)),
		state: 0,
		stop: make(chan bool),
		approvalsCountChan: make(chan int),
	}

	return session
}

func (session *AddFileClientGroupSession) close() {
	session.state = EndOfSession
	close(session.stop)
	close(session.approvalsCountChan)
}

func (session *AddFileClientGroupSession) Abort() {
	session.state = EndOfSession
	session.stop <- true
}

func (session *AddFileClientGroupSession) GetState() uint8 {
	return session.state
}

func (session *AddFileClientGroupSession) Id() IIdentifier {
	return session.sessionId
}

func (session *AddFileClientGroupSession) IsAlive() bool {
	return session.state == EndOfSession
}

func (session *AddFileClientGroupSession) Run() {
	addFileGroupMsg := AddFileGroupMessage{
		NewGroupIpfsPath: session.newIpfsPath,
		OldGroupIpfsPath: session.oldIpfsPath,
		NewFileCapIpfsHash: "",
	}
	payload, err := addFileGroupMsg.Encode()
	if err != nil {
		glog.Errorf("could not encoder add file group message %s", err)
		return
	}

	msg, err := NewGroupMessage(
		session.groupCtx.User.Address,
		AddFile,
		session.sessionId.Data().(uint32),
		payload,
		session.groupCtx.User.Signer,
	)
	if err != nil {
		glog.Errorf("could not create add file group message %s", err)
		return
	}

	encMsg, err := msg.Encode()
	if err != nil {
		glog.Errorf("could not encode group message: %s", err)
		return
	}

	if err := session.groupCtx.GroupConnection.SendAll(encMsg); err != nil {
		glog.Errorf("could not send to all group message: %s", err)
		return
	}
}

func (session *AddFileClientGroupSession) NextState(contact *Contact, data []byte) error {
	session.lock.Lock()
	defer session.lock.Unlock()

	if session.state != 0 {
		return errors.New("end of session")
	}

	sig := data
	if !contact.VerifyKey.Verify(session.digest, sig) {
		return errors.New("invalid approval")
	}

	approval := &Approval{
		From: contact.Address,
		Signature: sig,
	}
	session.approvals.Append(approval)

	if session.approvals.Count() < len(session.groupCtx.Group.Members) / 2 {
		return nil
	}

	groupId := session.groupCtx.Group.Id.Data().([32]byte)
	if err := session.groupCtx.Network.UpdateGroupIpfsPath(groupId, session.newIpfsPath, session.approvals); err != nil {
		return errors.Wrap(err, "could not send update group ipfs path transaction")
	}

	session.Abort()

	return nil
}

type AddFileServerGroupSession struct {
	sessionId            IIdentifier
	groupId              [32]byte
	newFileCapIpfsHash string
	oldIpfsPath string
	newIpfsPath          string
	ctx *GroupContext
	state                uint8
	lock                 sync.Mutex
	stop                 chan bool
	contact *Contact
}

func (session *AddFileServerGroupSession) close() {
	session.state = EndOfSession
	//close(session.stop)
	//session.conn.Close()
}

func (session *AddFileServerGroupSession) Abort() {
	session.state = EndOfSession
	session.stop <- true
}

func (session *AddFileServerGroupSession) GetState() uint8 {
	return session.state
}

func (session *AddFileServerGroupSession) Id() IIdentifier {
	return session.sessionId
}

func (session *AddFileServerGroupSession) IsAlive() bool {
	return session.state == EndOfSession
}

func (session *AddFileServerGroupSession) Run() {
	defer session.close()
	defer glog.Info("add file server group session ended")

	// TODO: check if valid

	hasher := crypto.NewKeccak256Hasher()
	digest := hasher.Sum([]byte(session.oldIpfsPath), []byte(session.newIpfsPath))
	sig, err := session.ctx.User.Signer.Sign(digest)
	if err != nil {
		glog.Errorf("could not sign approval: %s", err)
		return
	}

	msg, err := NewMessage(
		session.ctx.User.Address,
		AddFile,
		session.sessionId.Data().(uint32),
		sig,
		session.ctx.User.Signer,
	)
	if err != nil {
		glog.Errorf("could not create group message: %s", err)
		return
	}

	encMsg, err := msg.Encode()
	if err != nil {
		glog.Errorf("could not encode group message: %s", err)
		return
	}

	if err := session.contact.Send(encMsg); err != nil {
		glog.Errorf("could not send message: %s", err)
	}
}

func (session *AddFileServerGroupSession) NextState(contact *Contact, data []byte) error {
	return nil
}

func NewAddFileServerGroupSession(groupId [32]byte, msg *GroupMessage, contact *Contact, ctx *GroupContext) *AddFileServerGroupSession {
	addFileMsg, err := DecodeAddFileGroupMessage(msg.Payload)
	if err != nil {
		glog.Errorf("could not create AddFileServerGroupSession: %s", err)
		return nil
	}

	session := &AddFileServerGroupSession{
		ctx: ctx,
		sessionId: NewUint32Id(msg.SessionId),
		oldIpfsPath: addFileMsg.OldGroupIpfsPath,
		newIpfsPath: addFileMsg.NewGroupIpfsPath,
		stop: make(chan bool),
		groupId: groupId,
		contact: contact,

		state: 0,
	}

	return session
}