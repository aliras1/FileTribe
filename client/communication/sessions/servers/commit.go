package servers

import (
	"encoding/base64"
	"encoding/json"
	"github.com/golang/glog"
	"github.com/pkg/errors"

	comcommon "ipfs-share/client/communication/common"
	"ipfs-share/client/communication/sessions/common"
	"ipfs-share/client/fs"
	"ipfs-share/client/interfaces"
	. "ipfs-share/collections"
	"ipfs-share/crypto"
)

func NewCommitChangesGroupSessionServer(
	msg *comcommon.Message,
	contact *comcommon.Contact,
	user interfaces.IUser,
	group interfaces.IGroup,
	repo *fs.GroupRepo,
	onSuccess common.OnServerSuccessCallback,
	sessionClosed common.SessionClosedCallback,
) (common.ISession, error) {

	var argsJson []interface{}
	if err := json.Unmarshal(msg.Payload, &argsJson); err != nil {
		return nil, errors.Wrap(err, "could not unmarshal args")
	}

	if len(argsJson) < 1 {
		return nil, errors.New("args should be of length 1")
	}

	encNewIpfsHash, err := base64.StdEncoding.DecodeString(argsJson[0].(string))
	if err != nil {
		return nil, errors.Wrap(err, "could not decode new encrypted ipfs hash")
	}

	session := &ConsensusSessionServer{
		sessionId:         NewUint32Id(msg.SessionId),
		msgType:           comcommon.Commit,
		user:              user,
		group:             group,
		repo:              repo,
		onSessionClosed:   sessionClosed,
		onSuccess:         onSuccess,
		args:              []interface{} {encNewIpfsHash},
		validateOperation: ValidateCommit,
		contact:           contact,
		state:             0,
	}

	return session, nil
}


func ValidateCommit(
	args []interface{},
	user interfaces.IUser,
	contact *comcommon.Contact,
	group interfaces.IGroup,
	repo *fs.GroupRepo,
) (sig []byte, err error) {

	if len(args) < 1 {
		return nil, errors.New("args should be of length 1")
	}

	glog.Infof("args: %v", args)

	encNewIpfsHash := args[0].([]byte)

	boxer := group.Boxer()
	newIpfsHash, ok := boxer.BoxOpen(encNewIpfsHash)
	if !ok {
		err = errors.New("could not decrypt new ipfs hash")
		return
	}
	if err = repo.IsValidChangeSet(string(newIpfsHash), &contact.Address); err != nil {
		err = errors.Wrap(err, "invalid change set")
		return
	}

	hasher := crypto.NewKeccak256Hasher()
	digest := hasher.Sum(group.EncryptedIpfsHash(), encNewIpfsHash)
	glog.Errorf("signer digest: %v", digest)
	signer := user.Signer()
	sig, err = signer.Sign(digest)
	if err != nil {
		err = errors.Wrap(err, "could not sign approval")
		return
	}

	return
}