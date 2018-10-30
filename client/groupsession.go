package client

import (
	"sync"
	"github.com/golang/glog"
	. "ipfs-share/collections"
	"ipfs-share/crypto"
	"math/rand"
	"ipfs-share/network"
	"encoding/hex"
	"encoding/base64"
)

func NewGroupSessionServer(msg *Message, contact *Contact, ctx *GroupContext) ISession {
	switch msg.Type {
	case AddFile:
		{
			return NewCommitChangesGroupSessionServer(msg, contact, ctx)
		}
	default:
		{
			return nil
		}
	}
}

type CommitChangesGroupSessionClient struct {
	sessionId            IIdentifier
	encNewIpfsPathBase64 string
	approvals            []*network.Approval
	digest               []byte // the original message digest which is signed by the group members
	state                uint8
	groupCtx             *GroupContext
	lock                 sync.RWMutex
	approvalsCountChan   chan int
}

func NewCommitChangesGroupSessionClient(newIpfsPath string, groupCtx *GroupContext) ISession {
	sessionId := rand.Uint32()
	hasher := crypto.NewKeccak256Hasher()

	encNewIpfsPathBase64 := base64.URLEncoding.EncodeToString(groupCtx.Group.Boxer.BoxSeal([]byte(newIpfsPath)))
	glog.Errorf("master enc base64: %v", encNewIpfsPathBase64)

	digest := hasher.Sum([]byte(groupCtx.Group.EncryptedIpfsHash), []byte(encNewIpfsPathBase64))
	glog.Errorf("verif digest: %v", digest)

	session := &CommitChangesGroupSessionClient{
		sessionId:            NewUint32Id(sessionId),
		groupCtx:             groupCtx,
		encNewIpfsPathBase64: string(encNewIpfsPathBase64),
		approvals:            []*network.Approval{},
		digest:               digest,
		state:                0,
		approvalsCountChan:   make(chan int),
	}

	glog.Infof("----> Digest: %s, old: %s, new: %s", hex.EncodeToString(session.digest), groupCtx.Group.IpfsHash, encNewIpfsPathBase64)

	return session
}

func (session *CommitChangesGroupSessionClient) close() {
	session.state = EndOfSession
	session.groupCtx.P2P.SessionClosedChan <- session.Id()
}

func (session *CommitChangesGroupSessionClient) Abort() {
	session.lock.Lock()
	defer session.lock.Unlock()

	if !session.IsAlive() {
		return
	}

	session.close()
}

func (session *CommitChangesGroupSessionClient) GetState() uint8 {
	session.lock.RLock()
	defer session.lock.RUnlock()

	return session.state
}

func (session *CommitChangesGroupSessionClient) Id() IIdentifier {
	return session.sessionId
}

func (session *CommitChangesGroupSessionClient) IsAlive() bool {
	session.lock.RLock()
	defer session.lock.RUnlock()

	return session.state == EndOfSession
}

func (session *CommitChangesGroupSessionClient) Run() {
	addFileGroupMsg := AddFileGroupMessage{
		NewGroupIpfsPath: session.encNewIpfsPathBase64,
		NewFileCapIpfsHash: "",
	}
	payload, err := addFileGroupMsg.Encode()
	if err != nil {
		glog.Errorf("could not encode add file group message %s", err)
		session.close()
		return
	}

	msg, err := NewMessage(
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

func (session *CommitChangesGroupSessionClient) NextState(contact *Contact, data []byte) {
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

	approval := &network.Approval{
		From: contact.Address,
		Signature: sig,
	}

	session.approvals = append (session.approvals, approval)
	if len(session.approvals) <= len(session.groupCtx.Group.Members) / 2 {
		return
	}

	groupId := session.groupCtx.Group.Id.Data().([32]byte)
	if err := session.groupCtx.Network.UpdateGroupIpfsPath(groupId, session.encNewIpfsPathBase64, session.approvals); err != nil {
		session.close()
		glog.Errorf("could not send update group ipfs path transaction: %s", err)
		return
	}
}

type CommitChangesGroupSessionServer struct {
	sessionId            IIdentifier
	newFileCapIpfsHash   string
	encNewIpfsPathBase64 string
	groupCtx             *GroupContext
	state                uint8
	lock                 sync.RWMutex
	contact              *Contact
}

func (session *CommitChangesGroupSessionServer) close() {
	session.state = EndOfSession
	session.groupCtx.P2P.SessionClosedChan <- session.Id()
}

func (session *CommitChangesGroupSessionServer) Abort() {
	session.lock.Lock()
	defer session.lock.Unlock()

	if !session.IsAlive() {
		return
	}

	session.close()
}

func (session *CommitChangesGroupSessionServer) GetState() uint8 {
	session.lock.RLock()
	defer session.lock.RUnlock()

	return session.state
}

func (session *CommitChangesGroupSessionServer) Id() IIdentifier {
	return session.sessionId
}

func (session *CommitChangesGroupSessionServer) IsAlive() bool {
	session.lock.RLock()
	defer session.lock.RUnlock()

	return session.state == EndOfSession
}

func (session *CommitChangesGroupSessionServer) Run() {
	defer session.close()

	glog.Errorf("verif enc: %v", session.encNewIpfsPathBase64)
	encIpfsPath, err := base64.URLEncoding.DecodeString(session.encNewIpfsPathBase64)
	if err != nil {
		glog.Errorf("could not base64 decode encrypted new ipfs path")
	}

	newIpfsPath, ok := session.groupCtx.Group.Boxer.BoxOpen([]byte(encIpfsPath))
	if !ok {
		glog.Errorf("could not decrypt new ipfs path")
		return
	}
	if err := session.groupCtx.Repo.IsValidChangeSet(string(newIpfsPath), &session.contact.Address); err != nil {
		glog.Errorf("add file operation is invalid: %s", err)
		return
	}

	hasher := crypto.NewKeccak256Hasher()
	digest := hasher.Sum([]byte(session.groupCtx.Group.EncryptedIpfsHash), []byte(session.encNewIpfsPathBase64))
	glog.Errorf("signer digest: %v", digest)
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

func (session *CommitChangesGroupSessionServer) NextState(contact *Contact, data []byte) { }

func NewCommitChangesGroupSessionServer(msg *Message, contact *Contact, ctx *GroupContext) *CommitChangesGroupSessionServer {
	addFileMsg, err := DecodeAddFileGroupMessage(msg.Payload)
	if err != nil {
		glog.Errorf("could not create CommitChangesGroupSessionServer: %s", err)
		return nil
	}

	session := &CommitChangesGroupSessionServer{
		groupCtx:             ctx,
		sessionId:            NewUint32Id(msg.SessionId),
		encNewIpfsPathBase64: addFileMsg.NewGroupIpfsPath,
		contact:              contact,
		state:                0,
	}

	return session
}