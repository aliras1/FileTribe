package common

import (
	"encoding/binary"
	"encoding/json"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethcrypto"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

// MessageType ...
type MessageType byte

const (
	GetGroupData 		MessageType = 0
)

type Message struct {
	From    ethcommon.Address `json:"from"`
	Type    MessageType `json:"type"`
	SessionId uint32	`json:"session_id"`
	Payload []byte `json:"payload"`
	Sig     []byte	`json:"sig"`
}

type GroupData byte

const (
	GetGroupKey 		GroupData = 0
	GetProposedGroupKey GroupData = 1
)

type GroupDataMessage struct {
	Group ethcommon.Address `json:"group"`
	Data GroupData `json:"data"`
	Payload []byte
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


func NewMessage(from ethcommon.Address, msgType MessageType, sessionId uint32, payload []byte, signer Signer) (*Message, error) {
	msg := &Message{
		From: from,
		Type: msgType,
		SessionId: sessionId,
		Payload: payload,
	}

	sig, err := signer(msg.Digest())
	if err != nil {
		return nil, errors.Wrap(err, "could not sign message")
	}

	msg.Sig = sig

	return msg, nil
}

func (m *Message) Encode() ([]byte, error) {
	enc, err := json.Marshal(m)
	if err != nil {
		return nil, errors.Wrap(err, "could not encode Message")
	}

	return enc, nil
}

func DecodeMessage(data []byte) (*Message, error) {
	var m Message
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, errors.Wrap(err, "could not decode Message")
	}

	return &m, nil
}

func (m *Message) Validate(contact *Contact) error {
	if !contact.VerifySignature(m.Digest(), m.Sig) {
		return errors.New("invalid message")
	}

	return nil
}

func (m *Message) Digest() []byte {
	sessionIdBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(sessionIdBytes, m.SessionId)

	digest := ethcrypto.Keccak256(
		m.From.Bytes(),
		[]byte{byte(m.Type)},
		sessionIdBytes,
		m.Payload,
	)
	return digest[:]
}

func (m *GroupDataMessage) Encode() ([]byte, error) {
	enc, err := json.Marshal(m)
	if err != nil {
		return nil, errors.Wrap(err, "could not encode GroupDataMessage")
	}

	return enc, nil
}

func DecodeGroupDataMessage(data []byte) (*GroupDataMessage, error) {
	var m GroupDataMessage
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, errors.Wrap(err, "could not decode GroupDataMessage")
	}

	return &m, nil
}

// -----------------------------------------
// just to see how to use cbor...

// DialP2PConn creates 2 blockchain messages. One which is encrypted by
// the public key of the current user and another that is encrypted
// by the public key of the recipent. Both messages are put on the
// blockchain
//func (m *Message) DialP2PConn(otherBoxer, myBoxer *crypto.AnonymPublicKey, address ethcommon.Address, network *nw.Network) error {
//	var r [32]byte
//	if _, err := rand.Read(r[:]); err != nil {
//		return err
//	}
//
//	m.Random = r
//
//	var (
//		handler codec.CborHandle
//		out []byte
//	)
//
//	cbor := codec.NewEncoderBytes(&out, &handler)
//	if err := cbor.Encode(m); err != nil {
//		glog.Errorf("could not cbor encode other message: %s", err)
//		return err
//	}
//
//	var otherCbor = make([]byte, len(out))
//	copy(otherCbor, out)
//
//	glog.Info("otherCbor in Msg.DialP2PConn: ", otherCbor)
//
//	otherMessage, err := otherBoxer.Seal(otherCbor)
//	if err != nil {
//		return err
//	}
//
//	if _, err := rand.Read(r[:]); err != nil {
//		return err
//	}
//
//	m.Random = r
//
//	out = []byte{}
//	if err := cbor.Encode(m); err != nil {
//		glog.Errorf("could not cbor encode other message: %s", err)
//	}
//
//	var myCbor = make([]byte, len(out))
//	copy(myCbor, out)
//
//	myMessage, err := myBoxer.Seal(myCbor)
//	if err != nil {
//		return err
//	}
//
//	if err := network.DialP2PConn(otherMessage, myMessage, address); err != nil {
//		return fmt.Errorf("coudl not send message: Message.DialP2PConn: %s", err)
//	}
//
//	return nil
//}
//
//// Verify checks if the signiture on Payload is ok
//func (m *Message) Verify(verifyKey *crypto.VerifyKey) bool {
//	digest := sha3.Sum256(m.Payload)
//	return verifyKey.Verify(digest[:], m.Sig)
//}
//
//// NewFriendRequest creates a message with a FriendRequest as
//// a payload
//func NewFriendRequest(from, to ethcommon.Address, ipfsHash string, signer *crypto.Signer) (*Message, error) {
//	friendRequest := FriendRequest{
//		To:       to,
//		IpfsHash: ipfsHash,
//	}
//
//	// payload, err := json.Marshal(friendRequest)
//	// if err != nil {
//	// 	return nil, fmt.Errorf("could not marshal payload: NewFriendRequest: %s", err)
//	// }
//
//	var payload []byte
//	var handler codec.CborHandle
//
//	cbor := codec.NewEncoderBytes(&payload, &handler)
//
//	if err := cbor.Encode(friendRequest); err != nil {
//		return nil, fmt.Errorf("could not cbor encode payload: NewFriendRequest: %s", err)
//	}
//
//	digest := sha3.Sum256(payload)
//
//	sig, err := signer.Sign(digest[:])
//	if err != nil {
//		return nil, fmt.Errorf("could not sign message: NewFriendRequest: %s", err)
//	}
//
//	message := &Message{
//		From:    from,
//		Type:    FRIEND_REQUEST,
//		Payload: payload,
//		Sig:     sig[:64],
//	}
//
//	return message, nil
//}
//
//// NewFriendConfirmation creates a message with a FriendConfirmation as
//// a payload
//func NewFriendConfirmation(from, to ethcommon.Address, ipfsHash string, signer *crypto.Signer) (*Message, error) {
//	friendConfirmation := FriendConfirmation{
//		To:       to,
//		IpfsHash: ipfsHash,
//	}
//
//	// payload, err := json.Marshal(friendConfirmation)
//	// if err != nil {
//	// 	return nil, fmt.Errorf("could not marshal payload: NewFriendConfirm: %s", err)
//	// }
//
//	var payload []byte
//	var handler codec.CborHandle
//
//	cbor := codec.NewEncoderBytes(&payload, &handler)
//
//	if err := cbor.Encode(friendConfirmation); err != nil {
//		return nil, fmt.Errorf("could not cbor encode payload: NewFriendConfirmation: %s", err)
//	}
//
//	digest := sha3.Sum256(payload)
//
//	sig, err := signer.Sign(digest[:])
//	if err != nil {
//		return nil, fmt.Errorf("could not sign message: NewFriendConfirm: %s", err)
//	}
//
//	message := &Message{
//		From:    from,
//		Type:    FRIEND_CONFIRM,
//		Payload: payload,
//		Sig:     sig[:64],
//	}
//
//	return message, nil
//}
//
//// NewUpdateDirectory creates a message with a UpdateDirectory as
//// a payload
//func NewUpdateDirectory(from, to ethcommon.Address, ipfsHash string, signer *crypto.Signer) (*Message, error) {
//	updateDirectory := UpdateDirectory{
//		To:       to,
//		IpfsHash: ipfsHash,
//	}
//
//	// payload, err := json.Marshal(updateDirectory)
//	// if err != nil {
//	// 	return nil, fmt.Errorf("could not marshal payload: NewUpdateDirectory: %s", err)
//	// }
//
//	var payload []byte
//	var handler codec.CborHandle
//
//	cbor := codec.NewEncoderBytes(&payload, &handler)
//
//	if err := cbor.Encode(updateDirectory); err != nil {
//		return nil, fmt.Errorf("could not cbor encode payload: NewUpdateDirectory: %s", err)
//	}
//
//	digest := sha3.Sum256(payload)
//
//	sig, err := signer.Sign(digest[:])
//	if err != nil {
//		return nil, fmt.Errorf("could not sign message: NewUpdateDirectory: %s", err)
//	}
//
//	message := &Message{
//		From:    from,
//		Type:    UPDATE_DIRECTORY,
//		Payload: payload,
//		Sig:     sig[:64],
//	}
//
//	return message, nil
//}
