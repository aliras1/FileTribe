package client

import (
	"sync"
	"github.com/golang/glog"
	. "ipfs-share/collections"
	"ipfs-share/crypto"
	"math/rand"
	"ipfs-share/network"
	"encoding/hex"
	"github.com/pkg/errors"
)

func NewGroupSessionServer(msg *Message, contact *Contact,  user IUser, group IGroup, repo *GroupRepo, closedChan chan ISession) (ISession, error) {
	switch msg.Type {
	case AddFile:
		{
			return NewCommitChangesGroupSessionServer(msg, contact, user, group, repo, closedChan)
		}
	default:
		{
			return nil, errors.New("invalid message type")
		}
	}
}

type CommitChangesOnSuccessCallback func(encIpfsHash []byte, approvals []*network.Approval)

type CommitChangesGroupSessionClient struct {
	sessionId          IIdentifier
	encNewIpfsHash     []byte
	approvals          []*network.Approval
	digest             []byte // the original message digest which is signed by the group members
	state              uint8
	user IUser
	group IGroup
	groupConnection *GroupConnection
	closedChan chan ISession
	lock               sync.RWMutex
	approvalsCountChan chan int
	onSuccessCallback CommitChangesOnSuccessCallback
	error error
}

func NewCommitChangesGroupSessionClient(
	newIpfsHash string,
	user IUser,
	group IGroup,
	groupConnection *GroupConnection,
	closedChan chan ISession,
	onSuccess CommitChangesOnSuccessCallback,
	) ISession {

	sessionId := rand.Uint32()
	hasher := crypto.NewKeccak256Hasher()

	boxer := group.Boxer()
	encNewIpfsHash := boxer.BoxSeal([]byte(newIpfsHash))
	glog.Errorf("master enc base64: %v", encNewIpfsHash)

	digest := hasher.Sum(group.EncryptedIpfsHash(), encNewIpfsHash)
	glog.Errorf("verif digest: %v", digest)

	signer := user.Signer()
	sig, err := signer.Sign(digest)
	if err != nil {
		glog.Error("could not sign own CommitChanges digest")
	}

	session := &CommitChangesGroupSessionClient{
		sessionId:          NewUint32Id(sessionId),
		user: user,
		group: group,
		groupConnection: groupConnection,
		closedChan:closedChan,
		encNewIpfsHash:     encNewIpfsHash,
		approvals:          []*network.Approval{ {From: user.Address(), Signature: sig} },
		digest:             digest,
		state:              0,
		approvalsCountChan: make(chan int),
		onSuccessCallback: onSuccess,
	}

	glog.Infof("----> Digest: %s, old: %v, new: %v", hex.EncodeToString(session.digest), group.IpfsHash, encNewIpfsHash)

	return session
}

func (session *CommitChangesGroupSessionClient) Error() error {
	return session.error
}

func (session *CommitChangesGroupSessionClient) close() {
	session.state = EndOfSession
	session.closedChan <- session
}

func (session *CommitChangesGroupSessionClient) Abort() {
	session.lock.Lock()
	defer session.lock.Unlock()

	if !session.IsAlive() {
		return
	}

	session.close()
}

func (session *CommitChangesGroupSessionClient) State() uint8 {
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
	commitChangesGroupMsg := CommitGroupMessage{
		NewGroupIpfsHash: session.encNewIpfsHash,
	}
	payload, err := commitChangesGroupMsg.Encode()
	if err != nil {
		session.error = errors.Wrap(err, "could not encode commit changes group message")
		session.close()
		return
	}

	msg, err := NewMessage(
		session.user.Address(),
		AddFile,
		session.sessionId.Data().(uint32),
		payload,
		session.user.Signer(),
	)
	if err != nil {
		session.error = errors.Wrap(err, "could not create commit changes group message")
		session.close()
		return
	}

	encMsg, err := msg.Encode()
	if err != nil {
		session.error = errors.Wrap(err, "could not encode group message")
		session.close()
		return
	}

	if err := session.groupConnection.Broadcast(encMsg); err != nil {
		session.error = errors.Wrap(err, "could not broadcast group message")
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
	if len(session.approvals) <= session.group.CountMembers() / 2 {
		return
	}

	session.onSuccessCallback(session.encNewIpfsHash, session.approvals)
	session.close()
}

type CommitChangesGroupSessionServer struct {
	sessionId          IIdentifier
	newFileCapIpfsHash string
	encNewIpfsHash     []byte
	user IUser
	group IGroup
	repo *GroupRepo
	closedChan chan ISession
	state              uint8
	lock               sync.RWMutex
	contact            *Contact
	error error
}

func (session *CommitChangesGroupSessionServer) Error() error {
	return session.error
}

func (session *CommitChangesGroupSessionServer) close() {
	session.state = EndOfSession
	session.closedChan <- session
}

func (session *CommitChangesGroupSessionServer) Abort() {
	session.lock.Lock()
	defer session.lock.Unlock()

	if !session.IsAlive() {
		return
	}

	session.close()
}

func (session *CommitChangesGroupSessionServer) State() uint8 {
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

	boxer := session.group.Boxer()
	newIpfsHash, ok := boxer.BoxOpen(session.encNewIpfsHash)
	if !ok {
		session.error = errors.New("could not decrypt new ipfs hash")
		return
	}
	if err := session.repo.isValidChangeSet(string(newIpfsHash), &session.contact.Address); err != nil {
		session.error = errors.Wrap(err, "invalid change set")
		return
	}

	hasher := crypto.NewKeccak256Hasher()
	digest := hasher.Sum(session.group.EncryptedIpfsHash(), session.encNewIpfsHash)
	glog.Errorf("signer digest: %v", digest)
	signer := session.user.Signer()
	sig, err := signer.Sign(digest)
	if err != nil {
		session.error = errors.Wrap(err, "could not sign approval")
		return
	}

	msg, err := NewMessage(
		session.user.Address(),
		AddFile,
		session.sessionId.Data().(uint32),
		sig,
		session.user.Signer(),
	)
	if err != nil {
		session.error = errors.Wrap(err, "could not create group message")
		return
	}

	encMsg, err := msg.Encode()
	if err != nil {
		session.error = errors.Wrap(err, "could not encode group message")
		return
	}

	if err := session.contact.Send(encMsg); err != nil {
		session.error = errors.Wrap(err, "could not send message")
	}
}

func (session *CommitChangesGroupSessionServer) NextState(contact *Contact, data []byte) { }

func NewCommitChangesGroupSessionServer(msg *Message, contact *Contact, user IUser, group IGroup, repo *GroupRepo, closedChan chan ISession) (*CommitChangesGroupSessionServer, error) {
	addFileMsg, err := DecodeAddFileGroupMessage(msg.Payload)
	if err != nil {
		return nil, errors.Wrap(err, "could not read rand")
	}

	session := &CommitChangesGroupSessionServer{
		user:       user,
		group:group,
		repo:repo,
		closedChan:closedChan,
		sessionId:      NewUint32Id(msg.SessionId),
		encNewIpfsHash: addFileMsg.NewGroupIpfsHash,
		contact:        contact,
		state:          0,
	}

	return session, nil
}