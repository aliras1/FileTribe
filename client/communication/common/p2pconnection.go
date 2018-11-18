package common

import (
	"net"

	"github.com/golang/glog"
	"github.com/pkg/errors"

	"ipfs-share/collections"
	ipfsapi "ipfs-share/ipfs"
	"ipfs-share/network"
)

type P2PConn net.TCPConn

func (conn *P2PConn) ReadMessage(addressBook *collections.ConcurrentCollection, network network.INetwork, ipfs ipfsapi.IIpfs) (*Message, *Contact, error) {
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
