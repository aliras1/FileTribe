package client

import (
	"sync"
	"github.com/golang/glog"
	. "ipfs-share/collections"
	"ipfs-share/crypto"
	"math/rand"
	"ipfs-share/network"
	"encoding/hex"
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
	sessionId          IIdentifier
	encNewIpfsHash     []byte
	approvals          []*network.Approval
	digest             []byte // the original message digest which is signed by the group members
	state              uint8
	groupCtx           *GroupContext
	lock               sync.RWMutex
	approvalsCountChan chan int
}

func NewCommitChangesGroupSessionClient(newIpfsHash string, groupCtx *GroupContext) ISession {
	sessionId := rand.Uint32()
	hasher := crypto.NewKeccak256Hasher()

	boxer := groupCtx.Group.Boxer()
	encNewIpfsHash := boxer.BoxSeal([]byte(newIpfsHash))
	glog.Errorf("master enc base64: %v", encNewIpfsHash)

	digest := hasher.Sum(groupCtx.Group.EncryptedIpfsHash(), encNewIpfsHash)
	glog.Errorf("verif digest: %v", digest)

	session := &CommitChangesGroupSessionClient{
		sessionId:          NewUint32Id(sessionId),
		groupCtx:           groupCtx,
		encNewIpfsHash:     encNewIpfsHash,
		approvals:          []*network.Approval{},
		digest:             digest,
		state:              0,
		approvalsCountChan: make(chan int),
	}

	glog.Infof("----> Digest: %s, old: %v, new: %v", hex.EncodeToString(session.digest), groupCtx.Group.IpfsHash, encNewIpfsHash)

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
	addFileGroupMsg := CommitGroupMessage{
		NewGroupIpfsHash: session.encNewIpfsHash,
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
	if len(session.approvals) <= session.groupCtx.Group.CountMembers() / 2 {
		return
	}

	groupId := session.groupCtx.Group.Id().Data().([32]byte)
	if err := session.groupCtx.Network.UpdateGroupIpfsHash(groupId, session.encNewIpfsHash, session.approvals); err != nil {
		session.close()
		glog.Errorf("could not send update group ipfs hash transaction: %s", err)
		return
	}
}

type CommitChangesGroupSessionServer struct {
	sessionId          IIdentifier
	newFileCapIpfsHash string
	encNewIpfsHash     []byte
	groupCtx           *GroupContext
	state              uint8
	lock               sync.RWMutex
	contact            *Contact
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

	glog.Errorf("verif enc: %v", session.encNewIpfsHash)

	boxer := session.groupCtx.Group.Boxer()
	newIpfsHash, ok := boxer.BoxOpen(session.encNewIpfsHash)
	if !ok {
		glog.Errorf("could not decrypt new ipfs hash")
		return
	}
	if err := session.groupCtx.Repo.IsValidChangeSet(string(newIpfsHash), &session.contact.Address); err != nil {
		glog.Errorf("invalid change set: %s", err)
		return
	}

	hasher := crypto.NewKeccak256Hasher()
	digest := hasher.Sum(session.groupCtx.Group.EncryptedIpfsHash(), session.encNewIpfsHash)
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
		groupCtx:       ctx,
		sessionId:      NewUint32Id(msg.SessionId),
		encNewIpfsHash: addFileMsg.NewGroupIpfsHash,
		contact:        contact,
		state:          0,
	}

	return session
}