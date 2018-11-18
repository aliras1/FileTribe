package servers

import (
	"sync"

	"github.com/golang/glog"
	"github.com/pkg/errors"

	comcommon "ipfs-share/client/communication/common"
	"ipfs-share/client/communication/sessions/common"
	"ipfs-share/client/fs"
	"ipfs-share/client/interfaces"
	. "ipfs-share/collections"
	"ipfs-share/crypto"
)

type CommitChangesGroupSessionServer struct {
	sessionId          IIdentifier
	newFileCapIpfsHash string
	encNewIpfsHash     []byte
	user interfaces.IUser
	group interfaces.IGroup
	repo *fs.GroupRepo
	sessionClosed common.SessionClosedCallback
	state              uint8
	lock               sync.RWMutex
	contact            *comcommon.Contact
	error error
}

func (session *CommitChangesGroupSessionServer) Error() error {
	return session.error
}

func (session *CommitChangesGroupSessionServer) close() {
	session.state = common.EndOfSession
	session.sessionClosed(session)
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

	return session.state == common.EndOfSession
}

func (session *CommitChangesGroupSessionServer) Run() {
	defer session.close()

	boxer := session.group.Boxer()
	newIpfsHash, ok := boxer.BoxOpen(session.encNewIpfsHash)
	if !ok {
		session.error = errors.New("could not decrypt new ipfs hash")
		return
	}
	if err := session.repo.IsValidChangeSet(string(newIpfsHash), &session.contact.Address); err != nil {
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

	msg, err := comcommon.NewMessage(
		session.user.Address(),
		comcommon.Commit,
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

func (session *CommitChangesGroupSessionServer) NextState(contact *comcommon.Contact, data []byte) { }

func NewCommitChangesGroupSessionServer(
	msg *comcommon.Message,
	contact *comcommon.Contact,
	user interfaces.IUser,
	group interfaces.IGroup,
	repo *fs.GroupRepo,
	sessionClosed common.SessionClosedCallback,
	) (*CommitChangesGroupSessionServer, error) {

	session := &CommitChangesGroupSessionServer{
		user:       user,
		group:group,
		repo:repo,
		sessionClosed: sessionClosed,
		sessionId:      NewUint32Id(msg.SessionId),
		encNewIpfsHash: msg.Payload,
		contact:        contact,
		state:          0,
	}

	return session, nil
}