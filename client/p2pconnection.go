package client

import (
	"ipfs-share/ipfs"
	"context"
	"net"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"time"
	"ipfs-share/network"
	. "ipfs-share/collections"
)

const (
	p2pProtocolName = "ipfsShareP2P"
)

type P2PConn net.TCPConn

type P2PServer struct {
	sessions *ConcurrentCollection
	SessionClosedChan chan IIdentifier
	p2pListener *ipfs.P2PListener
	ctx *UserContext
	stop chan bool
	stopConnectionChannels []chan bool
}

func NewP2PConnection(port string, ctx *UserContext) (*P2PServer, error) {
	stop := make(chan bool)
	closedSession := make(chan IIdentifier)

	p2pListener, err := ctx.Ipfs.P2PListen(context.Background(), p2pProtocolName, "/ip4/127.0.0.1/tcp/" + port)
	if err != nil {
		return nil, errors.New("could not create P2P listener")
	}

	sessions := NewConcurrentCollection()

	p2p := &P2PServer{
		SessionClosedChan: closedSession,
		sessions: sessions,
		p2pListener: p2pListener,
		ctx: ctx,
		stop: stop,
	}

	go p2p.closedSessionListener()
	go p2p.connectionListener(port)

	return p2p, nil
}

func (p2p *P2PServer) AddSession(session ISession) {
	if err := p2p.sessions.Append(session); err != nil {
		glog.Warningf("could not append elem: %s", err)
	}
}

func (p2p *P2PServer) Stop() {
	p2p.stop <- true
}

func (p2p *P2PServer) closedSessionListener() {
	for sessionId := range p2p.SessionClosedChan {
		//p2p.sessions.DeleteWithId(sessionId)
		glog.Infof("session %d closed", sessionId.Data().(uint32))
	}
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
				for _, stopConnectionChannel := range p2p.stopConnectionChannels {
					stopConnectionChannel <- true
				}
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
				stopConnectionChannel := make(chan bool)
				p2p.stopConnectionChannels = append(p2p.stopConnectionChannels, stopConnectionChannel)
				conn.SetKeepAlive(true)
				go p2p.handleConnection(p2p.ctx.AddressBook, (*P2PConn)(conn), stopConnectionChannel)
				glog.Infof("%s is serving %s on %s", p2p.ctx.User.Name, conn.RemoteAddr().String(), conn.LocalAddr().String())
			}
		}
	}
}

func (p2p *P2PServer) handleConnection(addressBook *ConcurrentCollection, conn *P2PConn, stop chan bool) {
	defer close(stop)
	defer conn.Close()
	defer glog.Infof("exiting")

	for {
		select {
		case <- stop:
			{
				return
			}
		default:
			{
				msg, contact, err := conn.ReadMessage(addressBook, p2p.ctx.Network, p2p.ctx.Ipfs)
				if err != nil {
					glog.Errorf("could not read from connection: %s", err)
					continue
				}

				glog.Infof("%s (%s): msg from: %s, sessid: %d", p2p.ctx.User.Name, p2p.ctx.User.Address.String(), msg.From.String(), msg.SessionId)

				sessionId := NewUint32Id(msg.SessionId)
				var session ISession
				sessionInterface := p2p.sessions.Get(sessionId)
				if sessionInterface == nil {
					session = NewServerSession(msg, contact, p2p.ctx)
					if err := p2p.sessions.Append(session); err != nil {
						glog.Warningf("could not append elem: %s", err)
					}
					go session.Run()
					continue
				}

				session = sessionInterface.(ISession)
				go session.NextState(contact, msg.Payload)
			}
		}
	}
}

func (conn *P2PConn) ReadMessage(addressBook *ConcurrentCollection, network network.INetwork, ipfs ipfs.IIpfs) (*Message, *Contact, error) {
	data := make([]byte, 4096)

	length, err := conn.Read(data)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "could not read from net.Conn")
	}

	data = data[:length]

	msg, err := DecodeMessage(data)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not unmarshal Message")
	}

	contact, err := msg.Validate(network, ipfs)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "invalid message")
	}

	if err := addressBook.Append(contact); err != nil {
		glog.Warningf("could not append elem: %s", err)
	}
	contact = addressBook.Get(contact.Id()).(*Contact)

	return msg, contact, nil
}
