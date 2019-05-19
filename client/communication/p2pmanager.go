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

package communication

import (
	"context"
	"github.com/aliras1/FileTribe/tribecrypto"
	"net"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/glog"
	"github.com/pkg/errors"

	"github.com/aliras1/FileTribe/client/communication/common"
	"github.com/aliras1/FileTribe/client/communication/sessions/clients"
	sesscommon "github.com/aliras1/FileTribe/client/communication/sessions/common"
	"github.com/aliras1/FileTribe/client/communication/sessions/servers"
	"github.com/aliras1/FileTribe/client/interfaces"
	. "github.com/aliras1/FileTribe/collections"
	ipfsapi "github.com/aliras1/FileTribe/ipfs"
)

// P2PManager is responsible for managing all the incoming libp2p connections
type P2PManager struct {
	account        interfaces.Account
	signer         common.Signer
	sessions       *Map
	messagesCh     chan *common.Message
	addressBook    *common.AddressBook
	p2pListener    *ipfsapi.P2PListener
	ctxCallback    sesscommon.CtxCallback
	stop           chan struct{}
	stopConnection chan struct{}
	ipfs           ipfsapi.IIpfs
}

// NewP2PManager creates a new P2PManager
func NewP2PManager(
	port string,
	account interfaces.Account,
	signer common.Signer,
	addressBook *common.AddressBook,
	ctxCallback sesscommon.CtxCallback,
	ipfs ipfsapi.IIpfs,
) (*P2PManager, error) {

	stop := make(chan struct{})

	p2pListener, err := ipfs.P2PListen(context.Background(), common.P2PProtocolName, "/ip4/127.0.0.1/tcp/"+port)
	if err != nil {
		return nil, errors.Wrap(err, "could not create P2P listener")
	}

	p2p := &P2PManager{
		account:     account,
		signer:      signer,
		addressBook: addressBook,
		ctxCallback: ctxCallback,
		sessions:    NewConcurrentMap(),
		messagesCh:  make(chan *common.Message),
		p2pListener: p2pListener,
		stop:        stop,
		ipfs:        ipfs,
	}

	go p2p.connectionListener(port)

	return p2p, nil
}

// AddSession adds a session to the managers session list
func (p2p *P2PManager) AddSession(session sesscommon.Session) {
	p2p.sessions.Put(session.ID(), session)
}

// Cancel gracefully kills all threads and processes
func (p2p *P2PManager) Stop() {
	close(p2p.stop)
}

func (p2p *P2PManager) connectionListener(port string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:"+port)
	if err != nil {
		glog.Errorf("could not resolve tcp address: %s", err)
	}

	l, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		glog.Error("could not listen to port %s", port)
		return
	}
	defer l.Close()

	glog.Infof("listening on %s", tcpAddr.String())

	go p2p.handleMassages()

	for {
		select {
		case <-p2p.stop:
			{
				glog.Infof("stopping P2P connection")

				return
			}
		default:
			{
				l.SetDeadline(time.Now().Add(1e9))
				conn, err := l.AcceptTCP()
				if err != nil {
					//glog.Errorf("could not accept connection")
					continue
				}

				conn.SetKeepAlive(true)
				go p2p.handleConnection((*common.P2PConn)(conn))

				glog.Infof("%s is serving %s on %s", p2p.account.Name(), conn.RemoteAddr().String(), conn.LocalAddr().String())
			}
		}
	}
}

func (p2p *P2PManager) handleConnection(conn *common.P2PConn) {
	for {
		select {
		case <-p2p.stop:
			{
				close(p2p.stop)
				conn.Close()
				return
			}
		default:
			{
				messages, err := conn.ReadMessage(p2p.addressBook)
				if err != nil {
					glog.Errorf("%s: could not read from connection: %s", p2p.account.Name(), err)
					return
				}

				for _, msg := range messages {
					p2p.messagesCh <- msg
				}
			}
		}
	}
}

func (p2p *P2PManager) handleMassages() {
	for msg := range p2p.messagesCh {
		contact, err := p2p.addressBook.GetFromAccountAddress(msg.From)
		if err != nil {
			glog.Errorf("could not get contact from address book: %s", err)
			continue
		}

		glog.Infof("%s: msg from: %s, sessid: %d", p2p.account.Name(), msg.From.String(), msg.SessionID)

		var session sesscommon.Session
		sessionInterface := p2p.sessions.Get(msg.SessionID)

		if sessionInterface == nil {
			session, err = servers.NewGetGroupDataSessionServer(msg, contact, p2p.account.ContractAddress(), p2p.signer, p2p.ctxCallback, p2p.onSessionClosed)
			if err != nil {
				glog.Error("could not create new session: %s", err)
				continue
			}

			p2p.sessions.Put(session.ID(), session)
			go session.Run()
			continue
		}

		// TODO: fix bug: store original msg.from in the session and
		// check if the current sender is equal to that
		session = sessionInterface.(sesscommon.Session)
		go session.NextState(contact, msg.Payload)
	}
}

func (p2p *P2PManager) onSessionClosed(session sesscommon.Session) {
	glog.Infof("sid %v closed with error: %v", session.ID(), session.Error())
}

// StartGetGroupKeySession start a new session for retrieving a group's current key
func (p2p *P2PManager) StartGetGroupKeySession(
	group ethcommon.Address,
	receiverOwner ethcommon.Address,
	sender ethcommon.Address,
	resultCh chan tribecrypto.SymmetricKey,
) (sesscommon.Session, error) {
	contact, err := p2p.addressBook.GetFromOwnerAddress(receiverOwner)
	if err != nil {
		return nil, errors.Wrap(err, "could not get contact from owner")
	}

	session := clients.NewGetGroupDataSessionClient(
		common.GetGroupKey,
		group,
		nil, // no additional information needed
		contact,
		sender,
		p2p.signer,
		p2p.onSessionClosed,
		resultCh)

	p2p.AddSession(session)

	go session.Run()

	return session, nil
}

// StartGetProposedGroupKeySession starts a new session to get
// a specific proposed key of the group
func (p2p *P2PManager) StartGetProposedGroupKeySession(
	group ethcommon.Address,
	proposalKey []byte,
	receiverOwner ethcommon.Address,
	sender ethcommon.Address,
	resultCh chan tribecrypto.SymmetricKey,
) (sesscommon.Session, error) {
	glog.Info("StartGetProposedGroupKeySession...")

	contact, err := p2p.addressBook.GetFromOwnerAddress(receiverOwner)
	if err != nil {
		return nil, errors.Wrap(err, "could not get contact from owner")
	}

	session := clients.NewGetGroupDataSessionClient(
		common.GetProposedGroupKey,
		group,
		[]byte(proposalKey),
		contact,
		sender,
		p2p.signer,
		p2p.onSessionClosed,
		resultCh)

	p2p.AddSession(session)

	go session.Run()

	return session, nil
}
