// Copyright (c) 2019 Laszlo Sari
//
// FileTribe is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// FileTribe is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//

package common

import (
	"encoding/binary"
	"encoding/json"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

// MessageType is an enumeration of message types
type MessageType byte

const (
	// GetGroupData enum value
	GetGroupData MessageType = 0
)

// Message is a message struct
type Message struct {
	From      ethcommon.Address `json:"from"`
	Type      MessageType       `json:"type"`
	SessionID uint32            `json:"session_id"`
	Payload   []byte            `json:"payload"`
	Sig       []byte            `json:"sig"`
}

// GroupData is an enumeration of which group data wants to be retrieved by peers
type GroupData byte

const (
	// GetGroupKey ...
	GetGroupKey GroupData = 0
	// GetProposedGroupKey ...
	GetProposedGroupKey GroupData = 1
)

// GroupDataMessage is a wrapper for transferring group data like current key, a proposed key
type GroupDataMessage struct {
	Group   ethcommon.Address `json:"group"`
	Data    GroupData         `json:"data"`
	Payload []byte
}

// HeartBeat is a heartbeat message, currently not used
type HeartBeat struct {
	From string `json:"from"`
	Rand []byte `json:"rand"`
}

// NewMessage creates a new message
func NewMessage(from ethcommon.Address, msgType MessageType, sessionID uint32, payload []byte, signer Signer) (*Message, error) {
	msg := &Message{
		From:      from,
		Type:      msgType,
		SessionID: sessionID,
		Payload:   payload,
	}

	sig, err := signer(msg.Digest())
	if err != nil {
		return nil, errors.Wrap(err, "could not sign message")
	}

	msg.Sig = sig

	return msg, nil
}

// Encode encodes the message
func (m *Message) Encode() ([]byte, error) {
	enc, err := json.Marshal(m)
	if err != nil {
		return nil, errors.Wrap(err, "could not encode Message")
	}

	return enc, nil
}

// DecodeMessage decodes a message byte stream
func DecodeMessage(data []byte) (*Message, error) {
	var m Message
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, errors.Wrap(err, "could not decode Message")
	}

	return &m, nil
}

// Verify verifies a message if it really originates from the sender
func (m *Message) Verify(contact *Contact) error {
	if !contact.VerifySignature(m.Digest(), m.Sig) {
		return errors.New("invalid message")
	}

	return nil
}

// Digest returns the message digest
func (m *Message) Digest() []byte {
	sessionIDBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(sessionIDBytes, m.SessionID)

	digest := ethcrypto.Keccak256(
		m.From.Bytes(),
		[]byte{byte(m.Type)},
		sessionIDBytes,
		m.Payload,
	)
	return digest[:]
}

// Encode encodes the group data message
func (m *GroupDataMessage) Encode() ([]byte, error) {
	enc, err := json.Marshal(m)
	if err != nil {
		return nil, errors.Wrap(err, "could not encode GroupDataMessage")
	}

	return enc, nil
}

// DecodeGroupDataMessage decodes a group message data
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
