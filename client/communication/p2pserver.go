package communication

import (
	"ipfs-share/client/communication/sessions"
	sesscommon "ipfs-share/client/communication/sessions/common"
	"ipfs-share/client/interfaces"
	ipfsapi "ipfs-share/ipfs"
	"context"
	"net"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"time"
	"ipfs-share/network"
	. "ipfs-share/collections"
	"ipfs-share/client/communication/common"
)


type P2PHandleConnection func(addressBook *ConcurrentCollection, conn *common.P2PConn, stop chan struct{})

type P2PServer struct {
	user                 interfaces.IUser
	sessions             *ConcurrentCollection
	addressBook          *ConcurrentCollection
	SessionClosedChan    chan sesscommon.ISession
	p2pListener          *ipfsapi.P2PListener
	getGroupDataCallback sesscommon.GetGroupCallback
	stop                 chan struct{}
	stopConnection       chan struct{}
	ipfs                 ipfsapi.IIpfs
	network              network.INetwork
}

func NewP2PConnection(
	port string,
	user interfaces.IUser,
	addressBook *ConcurrentCollection,
	getGroupCallback sesscommon.GetGroupCallback,
	ipfs ipfsapi.IIpfs,
	network network.INetwork,
	) (*P2PServer, error) {

	stop := make(chan struct{})
	closedSession := make(chan sesscommon.ISession)

	p2pListener, err := ipfs.P2PListen(context.Background(), common.P2PProtocolName, "/ip4/127.0.0.1/tcp/" + port)
	if err != nil {
		return nil, errors.New("could not create P2P listener")
	}

	p2p := &P2PServer{
		user:                 user,
		addressBook:          addressBook,
		getGroupDataCallback: getGroupCallback,
		SessionClosedChan:    closedSession,
		sessions:             NewConcurrentCollection(),
		p2pListener:          p2pListener,
		stop:                 stop,
		ipfs:                 ipfs,
		network:              network,
	}

	go p2p.connectionListener(port)

	return p2p, nil
}

func (p2p *P2PServer) AddSession(session sesscommon.ISession) {
	if err := p2p.sessions.Append(session); err != nil {
		glog.Warningf("could not append elem: %s", err)
	}
}

func (p2p *P2PServer) Stop() {
	close(p2p.stop)
}


func (p2p *P2PServer) connectionListener(port string) {
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

				glog.Infof("%s is serving %s on %s", p2p.user.Name, conn.RemoteAddr().String(), conn.LocalAddr().String())
			}
		}
	}
}

func (p2p *P2PServer) handleConnection(addressBook *ConcurrentCollection, conn *common.P2PConn, stop chan struct{}) {
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
				msg, contact, err := conn.ReadMessage(addressBook, p2p.network, p2p.ipfs)
				if err != nil {
					glog.Errorf("could not read from connection: %s", err)
					continue
				}

				address := p2p.user.Address()
				glog.Infof("%s (%s): msg from: %s, sessid: %d", p2p.user.Name, address.String(), msg.From.String(), msg.SessionId)

				sessionId := NewUint32Id(msg.SessionId)
				var session sesscommon.ISession
				sessionInterface := p2p.sessions.Get(sessionId)
				if sessionInterface == nil {
					session, err = sessions.NewServerSession(msg, contact, p2p.user, p2p.getGroupDataCallback, p2p.SessionClosedChan)
					if err != nil {
						glog.Error("could not create new session: %s", err)
						continue
					}

					if err := p2p.sessions.Append(session); err != nil {
						glog.Warningf("could not append elem: %s", err)
					}
					go session.Run()
					continue
				}

				session = sessionInterface.(sesscommon.ISession)
				go session.NextState(contact, msg.Payload)
			}
		}
	}
}

