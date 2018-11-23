package clients

import (
	"encoding/json"
	"sync"

	"github.com/golang/glog"
	"github.com/pkg/errors"

	comcommon "ipfs-share/client/communication/common"
	"ipfs-share/client/communication/sessions/common"
	"ipfs-share/client/interfaces"
	"ipfs-share/collections"
	"ipfs-share/network"
)

type ConsensusSessionClient struct {
	sessionId         collections.IIdentifier
	msgType			  comcommon.MessageType
	args   			  []interface{}
	approvals         []*network.Approval
	digest            []byte // the original message digest which is signed by the group members
	state             uint8
	user              interfaces.IUser
	group             interfaces.IGroup
	broadcastFunction common.Broadcast
	onSessionClosed   common.SessionClosedCallback
	lock              sync.RWMutex
	onSuccessCallback common.OnClientSuccessCallback
	error             error
}


func (session *ConsensusSessionClient) Error() error {
	return session.error
}

func (session *ConsensusSessionClient) close() {
	session.state = common.EndOfSession
	session.onSessionClosed(session)
}

func (session *ConsensusSessionClient) Abort() {
	session.lock.Lock()
	defer session.lock.Unlock()

	if !session.IsAlive() {
		return
	}

	session.close()
}

func (session *ConsensusSessionClient) State() uint8 {
	session.lock.RLock()
	defer session.lock.RUnlock()

	return session.state
}

func (session *ConsensusSessionClient) Id() collections.IIdentifier {
	return session.sessionId
}

func (session *ConsensusSessionClient) IsAlive() bool {
	session.lock.RLock()
	defer session.lock.RUnlock()

	return session.state == common.EndOfSession
}

func (session *ConsensusSessionClient) Run() {
	glog.Infof("client args raw: %v", session.args)

	payload, err := json.Marshal(session.args)
	if err != nil {
		session.error = errors.Wrap(err, "could not marshal args")
		session.close()
		return
	}

	glog.Infof("client args json: %v", payload)

	msg, err := comcommon.NewMessage(
		session.user.Address(),
		session.msgType,
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

	if err := session.broadcastFunction(encMsg); err != nil {
		session.error = errors.Wrap(err, "could not broadcast group message")
		session.close()
		return
	}
}

func (session *ConsensusSessionClient) NextState(contact *comcommon.Contact, data []byte) {
	session.lock.Lock()
	defer session.lock.Unlock()

	if session.state != 0 {
		return
	}

	if !session.group.IsMember(contact.Address) {
		glog.Warning("user is not a group member")
		return
	}

	sig := data
	if !contact.VerifySignature(session.digest, sig) {
		glog.Warning("invalid approval")
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

	session.onSuccessCallback(session.args, session.approvals)
	session.close()
}