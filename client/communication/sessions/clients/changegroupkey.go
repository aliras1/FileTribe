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



func NewChangeGroupKeySessionClient(
	newBoxer crypto.SymmetricKey,
	encNewIpfsHash []byte,
	user interfaces.IUser,
	group interfaces.IGroup,
	broadcastFunction common.Broadcast,
	onSessionClosed common.SessionClosedCallback,
	onSuccess common.OnClientSuccessCallback,
) (common.ISession, error) {

	sessionId := rand.Uint32()
	hasher := crypto.NewKeccak256Hasher()

	digest := hasher.Sum(group.EncryptedIpfsHash(), encNewIpfsHash, user.Address().Bytes())

	signer := user.Signer()
	sig, err := signer.Sign(digest)
	if err != nil {
		return nil, errors.Wrap(err, "could not sign own Commit digest")
	}

	keyEncoded, err := newBoxer.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode new group key")
	}

	groupId := group.Id().Data().([32]byte)

	session := &ConsensusSessionClient{
		sessionId:         collections.NewUint32Id(sessionId),
		msgType:           comcommon.ChangeKey,
		user:              user,
		group:             group,
		broadcastFunction: broadcastFunction,
		onSessionClosed:   onSessionClosed,
		args:              []interface{} {encNewIpfsHash, keyEncoded, groupId[:]},
		onSuccessCallback: onSuccess,
		approvals:         []*network.Approval{ {From: user.Address(), Signature: sig} },
		digest:            digest,
		state:             0,
	}

	return session, nil
}
