package client

import (
	"sync"
	"github.com/golang/glog"
	. "ipfs-share/collections"
	"ipfs-share/crypto"
	"math/rand"
	"ipfs-share/networketh"
	"encoding/hex"
)

func NewServerGroupSession(msg *GroupMessage, contact *Contact, ctx *GroupContext) ISession {
	switch msg.Type {
	case AddFile:
		{
			return NewAddFileGroupSessionServer(msg, contact, ctx)
		}
	default:
		{
			return nil
		}
	}
}

type AddFileGroupSessionClient struct {
	sessionId            IIdentifier
	newIpfsPath          string
	approvals            []*networketh.Approval
	digest               []byte // the original message digest which is signed by the group members
	state                uint8
	groupCtx *GroupContext
	lock                 sync.RWMutex
	approvalsCountChan   chan int
}

func NewAddFileGroupSessionClient(newIpfsPath string, groupCtx *GroupContext) ISession {
	sessionId := rand.Uint32()
	hasher := crypto.NewKeccak256Hasher()

	session := &AddFileGroupSessionClient{
		sessionId: NewUint32Id(sessionId),
		groupCtx: groupCtx,
		newIpfsPath: newIpfsPath,
		approvals: []*networketh.Approval{},
		digest: hasher.Sum([]byte(groupCtx.Group.IPFSPath), []byte(newIpfsPath)),
		state: 0,
		approvalsCountChan: make(chan int),
	}

	glog.Infof("----> Digest: %s, old: %s, new: %s", hex.EncodeToString(session.digest), groupCtx.Group.IPFSPath, newIpfsPath)

	return session
}

func (session *AddFileGroupSessionClient) close() {
	session.state = EndOfSession
	session.groupCtx.P2P.SessionClosedChan <- session.Id()
}

func (session *AddFileGroupSessionClient) Abort() {
	session.lock.Lock()
	defer session.lock.Unlock()

	if !session.IsAlive() {
		return
	}

	session.close()
}

func (session *AddFileGroupSessionClient) GetState() uint8 {
	session.lock.RLock()
	defer session.lock.RUnlock()

	return session.state
}

func (session *AddFileGroupSessionClient) Id() IIdentifier {
	return session.sessionId
}

func (session *AddFileGroupSessionClient) IsAlive() bool {
	session.lock.RLock()
	defer session.lock.RUnlock()

	return session.state == EndOfSession
}

func (session *AddFileGroupSessionClient) Run() {
	addFileGroupMsg := AddFileGroupMessage{
		NewGroupIpfsPath: session.newIpfsPath,
		NewFileCapIpfsHash: "",
	}
	payload, err := addFileGroupMsg.Encode()
	if err != nil {
		glog.Errorf("could not encode add file group message %s", err)
		session.close()
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
		session.close()
		return
	}

	encMsg, err := msg.Encode()
	if err != nil {
		glog.Errorf("could not encode group message: %s", err)
		session.close()
		return
	}

	if err := session.groupCtx.GroupConnection.SendAll(encMsg); err != nil {
		glog.Errorf("could not send to all group message: %s", err)
		session.close()
		return
	}
}

func (session *AddFileGroupSessionClient) NextState(contact *Contact, data []byte) {
	session.lock.Lock()
	defer session.lock.Unlock()

	if session.state != 0 {
		return
	}

	sig := data
	if !contact.VerifySignature(session.digest, sig) {
		glog.Errorf("invalid approval")
		return
	}

	approval := &networketh.Approval{
		From: contact.Address,
		Signature: sig,
	}

	session.approvals = append (session.approvals, approval)
	if len(session.approvals) <= len(session.groupCtx.Group.Members) / 2 {
		return
	}

	groupId := session.groupCtx.Group.Id.Data().([32]byte)
	if err := session.groupCtx.Network.UpdateGroupIpfsPath(groupId, session.newIpfsPath, session.approvals); err != nil {
		session.close()
		glog.Errorf("could not send update group ipfs path transaction: %s", err)
		return
	}
}

type AddFileGroupSessionServer struct {
	sessionId          IIdentifier
	newFileCapIpfsHash string
	newIpfsPath        string
	groupCtx           *GroupContext
	state              uint8
	lock               sync.RWMutex
	contact            *Contact
}

func (session *AddFileGroupSessionServer) close() {
	session.state = EndOfSession
	session.groupCtx.P2P.SessionClosedChan <- session.Id()
}

func (session *AddFileGroupSessionServer) Abort() {
	session.lock.Lock()
	defer session.lock.Unlock()

	if !session.IsAlive() {
		return
	}

	session.close()
}

func (session *AddFileGroupSessionServer) GetState() uint8 {
	session.lock.RLock()
	defer session.lock.RUnlock()

	return session.state
}

func (session *AddFileGroupSessionServer) Id() IIdentifier {
	return session.sessionId
}

func (session *AddFileGroupSessionServer) IsAlive() bool {
	session.lock.RLock()
	defer session.lock.RUnlock()

	return session.state == EndOfSession
}

func (session *AddFileGroupSessionServer) Run() {
	defer session.close()

	if err := session.groupCtx.Repo.IsValidAddFile(session.newIpfsPath); err != nil {
		glog.Errorf("add file operation is invalid: %s", err)
		return
	}

	hasher := crypto.NewKeccak256Hasher()
	digest := hasher.Sum([]byte(session.groupCtx.Group.IPFSPath), []byte(session.newIpfsPath))
	sig, err := session.groupCtx.User.Signer.Sign(digest)
	if err != nil {
		glog.Errorf("could not sign approval: %s", err)
		return
	}

	msg, err := NewMessage(
		session.groupCtx.User.Address,
		AddFile,
		session.sessionId.Data().(uint32),
		sig,
		session.groupCtx.User.Signer,
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

func (session *AddFileGroupSessionServer) NextState(contact *Contact, data []byte) { }

func NewAddFileGroupSessionServer(msg *GroupMessage, contact *Contact, ctx *GroupContext) *AddFileGroupSessionServer {
	addFileMsg, err := DecodeAddFileGroupMessage(msg.Payload)
	if err != nil {
		glog.Errorf("could not create AddFileGroupSessionServer: %s", err)
		return nil
	}

	session := &AddFileGroupSessionServer{
		groupCtx:    ctx,
		sessionId:   NewUint32Id(msg.SessionId),
		newIpfsPath: addFileMsg.NewGroupIpfsPath,
		contact:     contact,
		state:       0,
	}

	return session
}