package clients

import (
	"github.com/pkg/errors"
	comcommon "ipfs-share/client/communication/common"
	"ipfs-share/client/communication/sessions/common"
	"ipfs-share/client/interfaces"
	"ipfs-share/collections"
	"ipfs-share/crypto"
	"ipfs-share/network"
	"math/rand"
)



func NewCommitGroupSessionClient(
	newIpfsHash string,
	user interfaces.IUser,
	group interfaces.IGroup,
	broadcastFunction common.Broadcast,
	onSessionClosed common.SessionClosedCallback,
	onSuccess common.OnClientSuccessCallback,
) (common.ISession, error) {

	sessionId := rand.Uint32()
	hasher := crypto.NewKeccak256Hasher()

	boxer := group.Boxer()
	encNewIpfsHash := boxer.BoxSeal([]byte(newIpfsHash))

	digest := hasher.Sum(group.EncryptedIpfsHash(), encNewIpfsHash)

	signer := user.Signer()
	sig, err := signer.Sign(digest)
	if err != nil {
		return nil, errors.Wrap(err, "could not sign own Commit digest")
	}

	session := &ConsensusSessionClient{
		sessionId:         collections.NewUint32Id(sessionId),
		msgType:           comcommon.Commit,
		user:              user,
		group:             group,
		broadcastFunction: broadcastFunction,
		onSessionClosed:   onSessionClosed,
		args:              []interface{} {encNewIpfsHash},
		onSuccessCallback: onSuccess,
		approvals:         []*network.Approval{ {From: user.Address(), Signature: sig} },
		digest:            digest,
		state:             0,
	}

	return session, nil
}
