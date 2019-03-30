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

package ipfs

import (
	"context"
	"io"

	ipfsapi "github.com/ipfs/go-ipfs-api"
	ma "github.com/multiformats/go-multiaddr"

	"strconv"
)

// IPubSubSubscription is an interface to IPFS pubsub subscriptions
type IPubSubSubscription interface {
	Cancel() error
	Next() (*ipfsapi.Message, error)
}

// Own version of PubSubSubscription
type PubSubSubscription ipfsapi.PubSubSubscription

// Cancel cancels a pubsub subscription
func (sub *PubSubSubscription) Cancel() error {
	return (*ipfsapi.PubSubSubscription)(sub).Cancel()
}

// Next reads the next message in the subscription
func (sub *PubSubSubscription) Next() (*ipfsapi.Message, error) {
	return ((*ipfsapi.PubSubSubscription)(sub)).Next()
}

type IIpfs interface {
	ID() (*ipfsapi.IdOutput, error)
	PubSubPublish(topic string, data string) error
	PubSubSubscribe(topic string) (IPubSubSubscription, error)
	Get(hash string, outdir string) error
	AddDir(dir string) (string, error)
	Publish(node string, value string) error
	Add(r io.Reader) (string, error)
	P2PListen(ctx context.Context, protocol, maddr string) (*P2PListener, error)
	P2PCloseListener(ctx context.Context, protocol string, closeAll bool) error
	P2PStreamDial(ctx context.Context, peerID, protocol, listenerMaddr string) (*P2PStream, error)
	P2PCloseStream(ctx context.Context, handlerID string, closeAll bool) error
}

type Ipfs struct {
	shell *ipfsapi.Shell
}

func NewIpfs(url string) *Ipfs {
	return &Ipfs{shell: ipfsapi.NewShell(url)}
}

func (ipfs *Ipfs) ID() (*ipfsapi.IdOutput, error) {
	return ipfs.shell.ID()
}

func (ipfs *Ipfs) PubSubPublish(topic string, data string) error {
	return ipfs.shell.PubSubPublish(topic, data)
}

func (ipfs *Ipfs) Get(hash string, outdir string) error {
	return ipfs.shell.Get(hash, outdir)
}

func (ipfs *Ipfs) AddDir(dir string) (string, error) {
	return ipfs.shell.AddDir(dir)
}

func (ipfs *Ipfs) Publish(node string, value string) error {
	return ipfs.shell.Publish(node, value)
}

func (ipfs *Ipfs) Add(r io.Reader) (string, error) {
	return ipfs.shell.Add(r)
}

func (ipfs *Ipfs) PubSubSubscribe(topic string) (IPubSubSubscription, error) {
	sub, err := ipfs.shell.PubSubSubscribe(topic)
	return (*PubSubSubscription)(sub), err
}

type P2PListener struct {
	Protocol string
	Address  string
}

func (ipfs *Ipfs) P2PListen(ctx context.Context, protocol, maddr string) (*P2PListener, error) {
	// TODO: replace with the official api version
	// Note that this feature is not implemented yet by the official api

	if _, err := ma.NewMultiaddr(maddr); err != nil {
		return nil, err
	}
	var response *P2PListener
	err := ipfs.shell.Request("p2p/listener/open").
					  Arguments(protocol, maddr).Exec(ctx, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (ipfs *Ipfs) P2PCloseListener(ctx context.Context, protocol string, closeAll bool) error {
	// TODO: replace with the official api version
	// Note that this feature is not implemented yet by the official api

	req := ipfs.shell.Request("p2p/listener/close").
				      Option("all", strconv.FormatBool(closeAll))
	if protocol != "" {
		req.Arguments(protocol)
	}
	if err := req.Exec(ctx, nil); err != nil {
		return err
	}
	return nil
}

type P2PStream struct {
	Protocol string
	Address  string
}

func (ipfs *Ipfs) P2PStreamDial(ctx context.Context, peerID, protocol, listenerMaddr string) (*P2PStream, error) {
	// TODO: replace with the official api version
	// Note that this feature is not implemented yet by the official api

	var response *P2PStream
	req := ipfs.shell.Request("p2p/stream/dial").
		Arguments(peerID, protocol)
	if listenerMaddr != "" {
		if _, err := ma.NewMultiaddr(listenerMaddr); err != nil {
			return nil, err
		}
		req.Arguments(listenerMaddr)
	}
	if err := req.Exec(ctx, &response); err != nil {
		return nil, err
	}
	return response, nil
}

func (ipfs *Ipfs) P2PCloseStream(ctx context.Context, handlerID string, closeAll bool) error {
	// TODO: replace with the official api version
	// Note that this feature is not implemented yet by the official api

	req := ipfs.shell.Request("p2p/stream/close").
		Option("all", strconv.FormatBool(closeAll))
	if handlerID != "" {
		req.Arguments(handlerID)
	}
	if err := req.Exec(ctx, nil); err != nil {
		return err
	}
	return nil
}