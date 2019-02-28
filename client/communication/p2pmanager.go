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
		return nil, errors.New("could not create P2P listener")
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
	defer close(p2p.stop)

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

				glog.Infof("%s is serving %s on %s", p2p.account.Name, conn.RemoteAddr().String(), conn.LocalAddr().String())
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
					glog.Errorf("could not read from connection: %s", err)
					continue
				}

				contact, err := addressBook.Get(msg.From)
				if err != nil {
					glog.Errorf("could not get contact from address book: %s", err)
					continue
				}

				address := p2p.account.ContractAddress()
				glog.Infof("%s (%s): msg from: %s, sessid: %d", p2p.account.Name(), address.String(), msg.From.String(), msg.SessionId)

				var session sesscommon.ISession
				sessionInterface := p2p.sessions.Get(msg.SessionId)
				if sessionInterface == nil {
					session, err = servers.NewGetGroupKeySessionServer(msg, contact, p2p.account.ContractAddress(), p2p.signer, p2p.ctxCallback, p2p.onSessionClosed)
					if err != nil {
						glog.Error("could not create new session: %s", err)
						continue
					}

					p2p.sessions.Put(session.Id(), session)

					go session.Run()
					continue
				}

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
	msgType common.MessageType,
	group ethcommon.Address,
	receiver *common.Contact,
	sender ethcommon.Address,
	onSuccess sesscommon.OnGetGroupKeySuccessCallback,
) error {
	session := clients.NewGetGroupKeySessionClient(
		msgType,
		group,
		receiver,
		sender,
		p2p.signer,
		p2p.onSessionClosed,
		onSuccess)

	p2p.AddSession(session)

	go session.Run()

	return nil
}