package client

import (
	"context"
	"ipfs-share/network"
	"ipfs-share/ipfs"
	"strings"
	"net"
	"github.com/pkg/errors"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/golang/glog"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	. "ipfs-share/collections"
	"bytes"
)

type Contact struct {
	*network.Contact
	conn *P2PConn
	ipfs ipfs.IIpfs
}

func NewContact(contact *network.Contact, ipfs ipfs.IIpfs) *Contact {
	return &Contact{
		Contact:contact,
		ipfs:ipfs,
	}
}

func (contact *Contact) Id() IIdentifier {
	return NewAddressId(&contact.Address)
}

func (contact *Contact) Send(data []byte) error {
	if contact.conn == nil {
		conn, err := contact.dialP2PConn(contact.ipfs)
		if err != nil {
			return errors.Wrap(err, "could not dial P2P connection")
		}
		contact.conn = conn
	}

	glog.Infof("sending P2P msg from %s to %s", contact.conn.RemoteAddr().String(), contact.conn.LocalAddr().String())

	if _, err := contact.conn.Write(data); err != nil {
		return errors.Wrap(err, "could not send data")
	}

	return nil
}

func (contact *Contact) VerifySignature(digest, signature []byte) bool {
	pk, err := ethcrypto.SigToPub(digest, signature)
	if err != nil {
		glog.Warningf("could not get pk from sig: %s", err)
	}

	otherAddress := ethcrypto.PubkeyToAddress(*pk)
	return bytes.Equal(contact.Address.Bytes(), otherAddress.Bytes())
}

func (contact *Contact) dialP2PConn(ipfs ipfs.IIpfs) (*P2PConn, error) {
	id, _ := ipfs.ID()
	glog.Infof("user with ipfs %s is P2P dialing to %s", id.ID, contact.IpfsPeerId)

	if contact.conn != nil {
		return contact.conn, nil
	}

	glog.Error(contact.IpfsPeerId)
	stream, err := ipfs.P2PStreamDial(context.Background(), contact.IpfsPeerId, p2pProtocolName, "")
	if err != nil {
		return nil, errors.Wrapf(err, "could not dial to stream %s", contact.IpfsPeerId)
	}

	multiAddress := ma.StringCast(stream.Address)
	var protocolValues []string
	for _, protocol := range multiAddress.Protocols() {
		value, _ := multiAddress.ValueForProtocol(protocol.Code)
		protocolValues = append(protocolValues, value)
	}

	hostAddress := strings.Join(protocolValues, ":")
	tcpAddr, err := net.ResolveTCPAddr("tcp", hostAddress)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		println("Dial failed:", err.Error())
		return nil, errors.Wrap(err, "could not dial to tcp server")
	}
	if err := conn.SetKeepAlive(true); err != nil {
		return nil, errors.Wrapf(err, "could not set keep alive to true")
	}

	contact.conn = (*P2PConn)(conn)

	return contact.conn, nil
}