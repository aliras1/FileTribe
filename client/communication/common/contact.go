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
	"bytes"
	"context"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"net"
	"strings"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/glog"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pkg/errors"

	. "github.com/aliras1/FileTribe/collections"
	"github.com/aliras1/FileTribe/ipfs"
	"github.com/aliras1/FileTribe/tribecrypto"
)

const (
	P2PProtocolName = "ipfsShareP2P"
)


type Contact struct {
	AccAddr    ethcommon.Address
	Address    ethcommon.Address
	Name       string
	IpfsPeerId string
	Boxer      tribecrypto.AnonymPublicKey
	conn       *P2PConn
	ipfs       ipfs.IIpfs
}

func NewContact(
	address ethcommon.Address,
	accAddr ethcommon.Address,
	name string,
	ipfsPeerId string,
	boxer tribecrypto.AnonymPublicKey,
	ipfs ipfs.IIpfs) *Contact {

	return &Contact{
		AccAddr:accAddr,
		Address: 	address,
		Name: 		name,
		IpfsPeerId: ipfsPeerId,
		Boxer: 		boxer,
		ipfs:		ipfs,
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
		return false
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

	stream, err := ipfs.P2PStreamDial(context.Background(), contact.IpfsPeerId, P2PProtocolName, "")
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
