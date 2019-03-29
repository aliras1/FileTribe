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
