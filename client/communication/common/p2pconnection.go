package common

import (
	"net"

	"github.com/pkg/errors"
)

type P2PConn net.TCPConn

func (conn *P2PConn) ReadMessage(addressBook *AddressBook) (*Message, error) {
	data := make([]byte, 4096)

	length, err := conn.Read(data)
	if err != nil {
		return nil, errors.Wrapf(err, "could not read from net.Conn")
	}

	data = data[:length]

	msg, err := DecodeMessage(data)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshal Message")
	}

	contact, err := addressBook.Get(msg.From)
	if err != nil {
		return nil, errors.Wrap(err, "could not get contact from address book")
	}

	if err := msg.Validate(contact); err != nil {
		return nil, errors.Wrapf(err, "invalid message")
	}

	return msg, nil
}
