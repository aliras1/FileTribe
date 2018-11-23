package servers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"

	"github.com/golang/glog"
	"github.com/pkg/errors"

	comcommon "ipfs-share/client/communication/common"
	"ipfs-share/client/communication/sessions/common"
	"ipfs-share/client/fs"
	"ipfs-share/client/interfaces"
	"ipfs-share/collections"
	"ipfs-share/crypto"
)

func NewChangeGroupKeySessionServer(
	msg *comcommon.Message,
	contact *comcommon.Contact,
	user interfaces.IUser,
	getGroupData common.GetGroupDataCallback,
	onSuccess common.OnServerSuccessCallback,
	sessionClosed common.SessionClosedCallback,
) (common.ISession, error) {

	encNewIpfsHash, newBoxer, groupId, err := extractArgs(msg.Payload)
	if err != nil {
		return nil, errors.Wrap(err, "could not extract args")
	}

	group, repo := getGroupData(groupId)
	if group == nil || repo == nil {
		return nil, errors.New("no group found")
	}

	leader := group.Leader()
	if !bytes.Equal(leader.Bytes(), contact.Address.Bytes()) {
		return nil, errors.New("user is not the group leader")
	}

	glog.Infof("server raw key :%v", newBoxer.Key)

	session := &ConsensusSessionServer{
		sessionId:         collections.NewUint32Id(msg.SessionId),
		msgType:           comcommon.ChangeKey,
		user:              user,
		group:             group,
		repo:              repo,
		onSessionClosed:   sessionClosed,
		onSuccess: 		   onSuccess,
		args:              []interface{} {encNewIpfsHash, newBoxer},
		validateOperation: ValidateChangeGroupKey,
		contact:           contact,
		state:             0,
	}

	return session, nil
}


func ValidateChangeGroupKey(
	args []interface{},
	user interfaces.IUser,
	contact *comcommon.Contact,
	group interfaces.IGroup,
	repo *fs.GroupRepo,
) (sig []byte, err error) {

	if len(args) < 2 {
		return nil, errors.New("args should be of length 2")
	}

	boxer := args[1].(crypto.SymmetricKey)
	encNewIpfsHash := args[0].([]byte)

	newIpfsHash, ok := boxer.BoxOpen(encNewIpfsHash)
	if !ok {
		err = errors.New("could not decrypt new ipfs hash")
		return
	}
	if err = repo.IsValidChangeKey(string(newIpfsHash), &contact.Address, boxer); err != nil {
		err = errors.Wrap(err, "invalid change key")
		return
	}

	hasher := crypto.NewKeccak256Hasher()
	digest := hasher.Sum(group.EncryptedIpfsHash(), encNewIpfsHash, contact.Address.Bytes())
	glog.Errorf("signer digest: %v", digest)
	signer := user.Signer()
	sig, err = signer.Sign(digest)
	if err != nil {
		err = errors.Wrap(err, "could not sign approval")
		return
	}

	return
}

func extractArgs(raw []byte) (encNewIpfsHash []byte, newBoxer crypto.SymmetricKey, groupId [32]byte, err error) {
	var argsJson []interface{}
	if err = json.Unmarshal(raw, &argsJson); err != nil {
		err = errors.Wrap(err, "could not unmarshal args")
		return
	}

	if len(argsJson) < 3 {
		err = errors.New("args should be of length 3")
		return
	}

	encNewIpfsHash, err = base64.StdEncoding.DecodeString(argsJson[0].(string))
	if err != nil {
		err = errors.Wrap(err, "could not decode new encrypted ipfs hash")
		return
	}

	glog.Infof("server json key :%v", argsJson[1])

	newBoxerEncoded, err := base64.StdEncoding.DecodeString(argsJson[1].(string))
	if err != nil {
		err = errors.Wrap(err, "could not base64 decode new key")
		return
	}

	newBoxerPtr, err := crypto.DecodeSymmetricKey(newBoxerEncoded)
	if err != nil {
		err = errors.Wrap(err, "could not decode new group key")
		return
	}
	newBoxer = *newBoxerPtr

	glog.Infof("groupId: %v", argsJson[2])
	groupIdRaw, err := base64.StdEncoding.DecodeString(argsJson[2].(string))
	if err != nil {
		err = errors.Wrap(err, "could not decode group id")
		return
	}

	copy(groupId[:], groupIdRaw)

	return
}