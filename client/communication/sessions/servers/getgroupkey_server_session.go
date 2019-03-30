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

package servers

import (
	"crypto/rand"
	"sync"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/glog"
	"github.com/pkg/errors"

	comcommon "github.com/aliras1/FileTribe/client/communication/common"
	"github.com/aliras1/FileTribe/client/communication/sessions/common"
)

// GetGroupKeySessionServer is a server session that will send the requested
// group data to the sender if he can authenticate itself and has access to
// the given data
type GetGroupKeySessionServer struct {
	sessionID       uint32
	state           uint8
	contact         *comcommon.Contact
	sender          ethcommon.Address
	groupDataMsg    comcommon.GroupDataMessage
	callback        common.CtxCallback
	signer          comcommon.Signer
	challenge       [32]byte
	onSessionClosed common.SessionClosedCallback
	lock            sync.RWMutex
	stop            chan bool
	error           error
	keyType         comcommon.MessageType
}

// Error returns any errors that may have occurred during the session
func (session *GetGroupKeySessionServer) Error() error {
	return session.error
}

func (session *GetGroupKeySessionServer) close() {
	session.state = common.EndOfSession
	session.onSessionClosed(session)
}

// State returns the session's current state
func (session *GetGroupKeySessionServer) State() uint8 {
	session.lock.RLock()
	defer session.lock.RUnlock()

	return session.state
}

// ID returns the session id
func (session *GetGroupKeySessionServer) ID() uint32 {
	return session.sessionID
}

// Abort aborts the session
func (session *GetGroupKeySessionServer) Abort() {
	if !session.isAlive() {
		return
	}

	session.close()
}

// IsAlive returns if the session is active or not
func (session *GetGroupKeySessionServer) IsAlive() bool {
	session.lock.RLock()
	defer session.lock.RUnlock()

	return session.isAlive()
}

func (session *GetGroupKeySessionServer) isAlive() bool {
	return session.state != common.EndOfSession
}

// Run starts the session
func (session *GetGroupKeySessionServer) Run() {
	session.NextState(nil, nil)
}

// NextState moves the FSM's state. For more information see ISession
func (session *GetGroupKeySessionServer) NextState(contact *comcommon.Contact, data []byte) {
	session.lock.Lock()
	defer session.lock.Unlock()

	switch session.state {
	case 0:
		{
			glog.Infof("server [%d] {%s} [0] --> %s", session.sessionID, session.sender.String(), session.contact.AccAddr.String())
			if err := session.callback.IsMember(session.groupDataMsg.Group, session.contact.AccAddr); err != nil {
				session.error = errors.Wrap(err, "could not verify group membership")
				session.close()
				return
			}
			glog.Infof("server [%d] [0][0]", session.sessionID)

			msg, err := comcommon.NewMessage(
				session.sender,
				comcommon.GetGroupData,
				session.sessionID,
				session.challenge[:],
				session.signer,
			)
			if err != nil {
				session.error = errors.New("could not create message")
				session.close()
				return
			}
			glog.Infof("server [%d] [0][1]", session.sessionID)

			encMsg, err := msg.Encode()
			if err != nil {
				session.error = errors.Wrap(err, "could not encode message")
				session.close()
				return
			}
			glog.Infof("server [%d] [0][2]", session.sessionID)

			if err := session.contact.Send(encMsg); err != nil {
				session.error = errors.Wrap(err, "could not send message")
				session.close()
				return
			}
			glog.Infof("server [%d] [0][3]", session.sessionID)

			session.state = 1

			return
		}
	case 1:
		{
			glog.Infof("server [%d] {%s} [1] --> %s", session.sessionID, session.sender.String(), session.contact.AccAddr.String())
			if !session.contact.VerifySignature(session.challenge[:], data) {
				session.error = errors.New("invalid signature")
				session.close()
				return
			}

			var key []byte

			switch session.groupDataMsg.Data {
			case comcommon.GetGroupKey:
				boxer, err := session.callback.GetBoxerOfGroup(session.groupDataMsg.Group)
				if err != nil {
					session.error = errors.Wrap(err, "could not get group boxer")
					session.close()
					return
				}

				data, err := boxer.Encode()
				if err != nil {
					session.error = errors.Wrap(err, "could not marshal group key")
					session.close()
					return
				}
				key = data
			case comcommon.GetProposedGroupKey:
				boxer, err := session.callback.GetProposedBoxerOfGroup(session.groupDataMsg.Group, ethcommon.BytesToAddress(session.groupDataMsg.Payload))
				if err != nil {
					session.error = errors.Wrapf(err, "%s could not get proposed group boxer", session.sender.String())
					session.close()
					return
				}
				glog.Infof("Sending back key: %v", boxer.Key)
				data, err := boxer.Encode()
				if err != nil {
					session.error = errors.Wrap(err, "could not marshal group key")
					session.close()
					return
				}
				key = data
			}

			msg, err := comcommon.NewMessage(
				session.sender,
				comcommon.GetGroupData,
				session.sessionID,
				key,
				session.signer,
			)
			if err != nil {
				session.error = errors.Wrap(err, "could not create message")
				session.close()
				return
			}

			encMsg, err := msg.Encode()
			if err != nil {
				session.error = errors.Wrap(err, "could not encode message")
				session.close()
				return
			}

			if err := session.contact.Send(encMsg); err != nil {
				session.error = errors.Wrap(err, "could not send message")
				session.close()
				return
			}

			session.close()
		}

	default:
		{
			glog.Errorf("session error: called next state in invalid state")
		}
	}
}

// NewGetGroupDataSessionServer creates a new session server that will
// send the requested group data to the sender
func NewGetGroupDataSessionServer(
	msg *comcommon.Message,
	contact *comcommon.Contact,
	sender ethcommon.Address,
	signer comcommon.Signer,
	callback common.CtxCallback,
	onSessionClosed common.SessionClosedCallback,
) (*GetGroupKeySessionServer, error) {

	var challenge [32]byte
	if _, err := rand.Read(challenge[:]); err != nil {
		return nil, errors.Wrap(err, "could not read rand")
	}

	groupDataMsg, err := comcommon.DecodeGroupDataMessage(msg.Payload)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode message payload")
	}

	return &GetGroupKeySessionServer{
		sessionID:       msg.SessionID,
		contact:         contact,
		callback:        callback,
		sender:          sender,
		signer:          signer,
		groupDataMsg:    *groupDataMsg,
		onSessionClosed: onSessionClosed,
		state:           0,
		challenge:       challenge,
	}, nil
}
