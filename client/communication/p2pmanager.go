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
		p2pListener: p2pListener,
		stop:        stop,
		ipfs:        ipfs,
	}

	go p2p.connectionListener(port)

	return p2p, nil
}

// AddSession adds a session to the managers session list
func (p2p *P2PManager) AddSession(session sesscommon.ISession) {
	p2p.sessions.Put(session.ID(), session)
}

// Stop gracefully kills all threads and processes
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
				go p2p.handleConnection(p2p.addressBook, (*common.P2PConn)(conn), p2p.stop)

				glog.Infof("%s is serving %s on %s", p2p.account.Name(), conn.RemoteAddr().String(), conn.LocalAddr().String())
			}
		}
	}
}

func (p2p *P2PManager) handleConnection(addressBook *common.AddressBook, conn *common.P2PConn, stop chan struct{}) {
	for {
		select {
		case <-stop:
			{
				close(stop)
				conn.Close()
				return
			}
		default:
			{
				msg, err := conn.ReadMessage(addressBook)
				if err != nil {
					glog.Errorf("%s: could not read from connection: %s", p2p.account.Name(), err)
					return
				}

				contact, err := addressBook.Get(msg.From)
				if err != nil {
					glog.Errorf("could not get contact from address book: %s", err)
					return
				}

				glog.Infof("%s: msg from: %s, sessid: %d", p2p.account.Name(), msg.From.String(), msg.SessionID)

				var session sesscommon.ISession
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
				session = sessionInterface.(sesscommon.ISession)
				go session.NextState(contact, msg.Payload)
			}
		}
	}
}

func (p2p *P2PManager) onSessionClosed(session sesscommon.ISession) {
	glog.Infof("sid %v closed with error: %v", session.ID(), session.Error())
}

// StartGetGroupKeySession start a new session for retrieving a group's current key
func (p2p *P2PManager) StartGetGroupKeySession(
	group ethcommon.Address,
	receiver *common.Contact,
	sender ethcommon.Address,
	onSuccess sesscommon.OnGetGroupKeySuccessCallback,
) error {
	session := clients.NewGetGroupDataSessionClient(
		common.GetGroupKey,
		group,
		nil, // no additional information needed
		receiver,
		sender,
		p2p.signer,
		p2p.onSessionClosed,
		onSuccess)

	p2p.AddSession(session)

	go session.Run()

	return nil
}

// StartGetProposedGroupKeySession starts a new session to get
// a specific proposed key of the group
func (p2p *P2PManager) StartGetProposedGroupKeySession(
	group ethcommon.Address,
	proposer ethcommon.Address,
	receiver *common.Contact,
	sender ethcommon.Address,
	onSuccess sesscommon.OnGetGroupKeySuccessCallback,
) error {
	glog.Info("StartGetProposedGroupKeySession...")

	session := clients.NewGetGroupDataSessionClient(
		common.GetProposedGroupKey,
		group,
		proposer.Bytes(),
		receiver,
		sender,
		p2p.signer,
		p2p.onSessionClosed,
		onSuccess)

	p2p.AddSession(session)

	go session.Run()

	return nil
}
