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
	"strings"

	"github.com/golang/glog"
	"github.com/pkg/errors"
)

// P2PConn is tcp connection to an IPFS p2p dial/stream endpoint
type P2PConn net.TCPConn

// ReadMessage reads a message from the connection socket
func (conn *P2PConn) ReadMessage(addressBook *AddressBook) ([]*Message, error) {
	data := make([]byte, 4096)

	length, err := conn.Read(data)
	if err != nil {
		return nil, errors.Wrapf(err, "could not read from net.Conn")
	}

	data = data[:length]
	messagesBytes, err := splitBulkMessages(data)
	if err != nil {
		return nil, errors.Wrap(err, "could not split up incoming data stream to messages")
	}

	var messages []*Message
	for _, messageBytes := range messagesBytes {
		msg, err := DecodeMessage(messageBytes)
		if err != nil {
			glog.Infof("Error in msg: %s", string(data))
			glog.Errorf("could not unmarshal Message: %s", err)
			continue
		}

		contact, err := addressBook.GetFromAccountAddress(msg.From)
		if err != nil {
			glog.Errorf("could not get contact from address book: %s", err)
			continue
		}

		if err := msg.Verify(contact); err != nil {
			glog.Errorf("invalid message: %s", err)
			continue
		}

		messages = append(messages, msg)
	}

	return messages, nil
}

func splitBulkMessages(data []byte) ([][]byte, error) {
	var messages [][]byte

	counter := 0
	var currMsg []byte
	for _, bit := range data {
		currMsg = append(currMsg, bit)

		if strings.EqualFold(string(bit), "{") {
			counter++
		} else if strings.EqualFold(string(bit), "}") {
			counter--
		}

		if counter == 0 {
			messages = append(messages, currMsg)
			currMsg = nil
		} else if counter < 0 {
			return nil, errors.New("invalid json message")
		}
	}

	return messages, nil
}