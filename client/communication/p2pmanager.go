package communication

import (
	"context"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"ipfs-share/client/communication/common"
	"ipfs-share/client/communication/sessions/clients"
	sesscommon "ipfs-share/client/communication/sessions/common"
	"ipfs-share/client/communication/sessions/servers"
	"ipfs-share/client/interfaces"
	. "ipfs-share/collections"
	"ipfs-share/crypto"
	ipfsapi "ipfs-share/ipfs"
	"net"
	"time"
)


type P2PHandleConnection func(addressBook *Map, conn *common.P2PConn, stop chan struct{})

type P2PManager struct {
	account        interfaces.IAccount
	signer		   *tribecrypto.Signer
	sessions       *Map
	addressBook    *common.AddressBook
	p2pListener    *ipfsapi.P2PListener
	ctxCallback    sesscommon.CtxCallback
	stop           chan struct{}
	stopConnection chan struct{}
	ipfs           ipfsapi.IIpfs
}

func NewP2PManager(
	port string,
	account interfaces.IAccount,
	signer *tribecrypto.Signer,
	addressBook *common.AddressBook,
	ctxCallback sesscommon.CtxCallback,
	ipfs ipfsapi.IIpfs,
	) (*P2PManager, error) {

	stop := make(chan struct{})

	p2pListener, err := ipfs.P2PListen(context.Background(), common.P2PProtocolName, "/ip4/127.0.0.1/tcp/" + port)
	if err != nil {
		return nil, errors.Wrap(err, "could not create P2P listener")
	}

	p2p := &P2PManager{
		account:     account,
		signer:		 signer,
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

func (p2p *P2PManager) AddSession(session sesscommon.ISession) {
	p2p.sessions.Put(session.Id(), session)
}

func (p2p *P2PManager) Stop() {
	close(p2p.stop)
}


func (p2p *P2PManager) connectionListener(port string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:" + port)
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
		case <- p2p.stop:
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
	defer func() {
		close(stop)
		conn.Close()
		glog.Infof("exiting")
	}()

	for {
		select {
		case <- stop:
			{
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

				glog.Infof("%s: msg from: %s, sessid: %d", p2p.account.Name(), msg.From.String(), msg.SessionId)

				var session sesscommon.ISession
				sessionInterface := p2p.sessions.Get(msg.SessionId)

				glog.Info("d0")
				if sessionInterface == nil {
					glog.Info("d0_0")
					session, err = servers.NewGetGroupDataSessionServer(msg, contact, p2p.account.ContractAddress(), p2p.signer, p2p.ctxCallback, p2p.onSessionClosed)
					if err != nil {
						glog.Error("could not create new session: %s", err)
						continue
					}
					glog.Info("d0_1")
					p2p.sessions.Put(session.Id(), session)
					glog.Info("d0_2")
					go session.Run()
					continue
				}

				glog.Info("d1")
				session = sessionInterface.(sesscommon.ISession)
				go session.NextState(contact, msg.Payload)
			}
		}
	}
}

func (p2p *P2PManager) onSessionClosed(session sesscommon.ISession) {
	glog.Infof("sid %v closed with error: %v", session.Id(), session.Error())
}

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