package servers

import (
	"sync"

	"github.com/pkg/errors"

	comcommon "ipfs-share/client/communication/common"
	"ipfs-share/client/communication/sessions/common"
	"ipfs-share/client/fs"
	"ipfs-share/client/interfaces"
	. "ipfs-share/collections"
)

type ValidateOperationCallback func(
	args []interface{},
	user interfaces.IUser,
	contact *comcommon.Contact,
	group interfaces.IGroup,
	repo *fs.GroupRepo,
	) (sig []byte, err error)

type ConsensusSessionServer struct {
	sessionId         IIdentifier
	msgType           comcommon.MessageType
	args              []interface{}
	validateOperation ValidateOperationCallback
	onSuccess         common.OnServerSuccessCallback
	user              interfaces.IUser
	group             interfaces.IGroup
	repo              *fs.GroupRepo
	onSessionClosed   common.SessionClosedCallback
	state             uint8
	lock              sync.RWMutex
	contact           *comcommon.Contact
	error             error
}

func (session *ConsensusSessionServer) Error() error {
	return session.error
}

func (session *ConsensusSessionServer) close() {
	session.state = common.EndOfSession
	session.onSessionClosed(session)
}

func (session *ConsensusSessionServer) Abort() {
	session.lock.Lock()
	defer session.lock.Unlock()

	if !session.IsAlive() {
		return
	}

	session.close()
}

func (session *ConsensusSessionServer) State() uint8 {
	session.lock.RLock()
	defer session.lock.RUnlock()

	return session.state
}

func (session *ConsensusSessionServer) Id() IIdentifier {
	return session.sessionId
}

func (session *ConsensusSessionServer) IsAlive() bool {
	session.lock.RLock()
	defer session.lock.RUnlock()

	return session.state == common.EndOfSession
}

func (session *ConsensusSessionServer) Run() {
	defer session.close()

	sig, err := session.validateOperation(session.args, session.user, session.contact, session.group, session.repo)
	if err != nil {
		session.error = errors.Wrap(err, "invalid operation")
		return
	}

	msg, err := comcommon.NewMessage(
		session.user.Address(),
		session.msgType,
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

	session.onSuccess(session.args, session.group.Id())
}

func (session *ConsensusSessionServer) NextState(contact *comcommon.Contact, data []byte) { }