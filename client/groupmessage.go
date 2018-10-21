package client

import (
	"ipfs-share/crypto"
	"reflect"
	"encoding/json"
	"ipfs-share/networketh"
	"encoding/binary"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	. "ipfs-share/collections"
	"ipfs-share/ipfs"
)

type GroupMessage struct {
	From ethcommon.Address `json:"from"`
	SessionId uint32 `json:"session_id"`
	Type MessageType `json:"type"`
	Payload []byte `json:"payload"`
	Sig []byte `json:"sig"`
}

func NewGroupMessage(from ethcommon.Address, msgType MessageType, sessionId uint32, payload []byte, signer *crypto.Signer) (*GroupMessage, error) {
	msg := &GroupMessage{
		From: from,
		Type: msgType,
		SessionId: sessionId,
		Payload: payload,
	}

	sig, err := signer.Sign(msg.Digest())
	if err != nil {
		return nil, errors.Wrap(err, "could not sign message")
	}

	msg.Sig = sig

	return msg, nil
}

func (m *GroupMessage) Encode() ([]byte, error) {
	enc, err := json.Marshal(m)
	if err != nil {
		errors.Wrap(err, "could not encode Message")
	}

	return enc, nil
}

func DecodeGroupMessage(data []byte) (*GroupMessage, error) {
	var m GroupMessage
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, errors.Wrap(err, "could not decode Message")
	}

	return &m, nil
}

func (m *GroupMessage) Validate(network networketh.INetwork, ipfs ipfs.IIpfs) (*Contact, error) {
	netContact, err := network.GetUser(m.From)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get user %s", m.From.String())
	}

	contact := NewContact(netContact, ipfs)

	if !contact.VerifySignature(m.Digest(), m.Sig) {
		return nil, errors.New("invalid message")
	}

	return contact, nil
}

func (m *GroupMessage) Digest() []byte {
	sessionIdBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(sessionIdBytes, m.SessionId)

	hasher := crypto.NewKeccak256Hasher()
	digest := hasher.Sum(
		m.From.Bytes(),
		[]byte{byte(m.Type)},
		sessionIdBytes,
		m.Payload,
	)
	return digest[:]
}

func InArray(array interface{}, item interface{}) (int, bool) {
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(item, s.Index(i).Interface()) == true {
				return i, true
			}
		}
	}
	return -1, false
}

type AddFileGroupMessage struct {
	NewFileCapIpfsHash string
	OldGroupIpfsPath string
	NewGroupIpfsPath string
}

func (msg *AddFileGroupMessage) Encode() ([]byte, error) {
	enc, err := json.Marshal(msg)
	if err != nil {
		return nil, errors.Wrapf(err, "could not marshal AddFileGroupMessage")
	}

	return enc, nil
}

func DecodeAddFileGroupMessage(data []byte) (*AddFileGroupMessage, error) {
	var msg AddFileGroupMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, errors.Wrapf(err, "could not unmarshal AddFileGroupMessage")
	}

	return &msg, nil
}

type HeartBeat struct {
	From string `json:"from"`
	Rand []byte `json:"rand"`
}

type Proposal struct {
	From     string   `json:"from"`
	CMD      string   `json:"cmd"`
	Args     []string `json:"args"`
	PrevHash [32]byte `json:"prev_hash"`
	NewHash  [32]byte `json:"new_hash"`
}

type Approval networketh.Approval

func (approval *Approval) Id() IIdentifier {
	return NewAddressId(&approval.From)
}

func (approval *Approval) Validate(digest []byte, contact *Contact) error {
	// signed := append(approval.Signature, rawTransaction...)
	// verifyKey, err := network.GetUserVerifyKey(approval.From)
	// if err != nil {
	// 	return fmt.Errorf("could not get user verify key: ValidateApproval: %s", err)
	// }
	// _, ok := verifyKey.Open(nil, signed)
	// if !ok {
	// 	return fmt.Errorf("invalid approval: ValidateApproval")
	// }
	return nil
}
