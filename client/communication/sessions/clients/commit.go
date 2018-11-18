package clients

import (
	"encoding/hex"
	"math/rand"
	"sync"

	"github.com/golang/glog"
	"github.com/pkg/errors"

	comcommon "ipfs-share/client/communication/common"
	"ipfs-share/client/communication/sessions/common"
	"ipfs-share/client/interfaces"
	"ipfs-share/collections"
	"ipfs-share/crypto"
	"ipfs-share/network"
)

type Broadcast func(msg []byte) error

type CommitOnSuccessCallback func(encIpfsHash []byte, approvals []*network.Approval)

type CommitGroupSessionClient struct {
	sessionId         collections.IIdentifier
	encNewIpfsHash    []byte
	approvals         []*network.Approval
	digest            []byte // the original message digest which is signed by the group members
	state             uint8
	user              interfaces.IUser
	group             interfaces.IGroup
	broadcastFunction Broadcast
	closedChan        chan common.ISession
	lock              sync.RWMutex
	onSuccessCallback CommitOnSuccessCallback
	error             error
}

func NewCommitChangesGroupSessionClient(
	newIpfsHash string,
	user interfaces.IUser,
	group interfaces.IGroup,
	broadcastFunction Broadcast,
	closedChan chan common.ISession,
	onSuccess CommitOnSuccessCallback,
) common.ISession {

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

	session := &CommitGroupSessionClient{
		sessionId:         collections.NewUint32Id(sessionId),
		user:              user,
		group:             group,
		broadcastFunction: broadcastFunction,
		closedChan:        closedChan,
		encNewIpfsHash:    encNewIpfsHash,
		approvals:         []*network.Approval{ {From: user.Address(), Signature: sig} },
		digest:            digest,
		state:             0,
		onSuccessCallback: onSuccess,
	}

	glog.Infof("----> Digest: %s, old: %v, new: %v", hex.EncodeToString(session.digest), group.IpfsHash, encNewIpfsHash)

	return session
}

func (session *CommitGroupSessionClient) Error() error {
	return session.error
}

func (session *CommitGroupSessionClient) close() {
	session.state = common.EndOfSession
	session.closedChan <- session
}

func (session *CommitGroupSessionClient) Abort() {
	session.lock.Lock()
	defer session.lock.Unlock()

	if !session.IsAlive() {
		return
	}

	session.close()
}

func (session *CommitGroupSessionClient) State() uint8 {
	session.lock.RLock()
	defer session.lock.RUnlock()

	return session.state
}

func (session *CommitGroupSessionClient) Id() collections.IIdentifier {
	return session.sessionId
}

func (session *CommitGroupSessionClient) IsAlive() bool {
	session.lock.RLock()
	defer session.lock.RUnlock()

	return session.state == common.EndOfSession
}

func (session *CommitGroupSessionClient) Run() {
	msg, err := comcommon.NewMessage(
		session.user.Address(),
		comcommon.Commit,
		session.sessionId.Data().(uint32),
		session.encNewIpfsHash,
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

func (session *CommitGroupSessionClient) NextState(contact *comcommon.Contact, data []byte) {
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