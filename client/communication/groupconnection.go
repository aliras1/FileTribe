package communication

import (
	"github.com/golang/glog"
	"ipfs-share/client/communication/common"
	"ipfs-share/client/communication/sessions"
	sesscommon "ipfs-share/client/communication/sessions/common"
	"ipfs-share/client/fs"
	"ipfs-share/client/interfaces"
	"ipfs-share/collections"
	ipfsapi "ipfs-share/ipfs"
	"bytes"
	"encoding/base64"
	"github.com/pkg/errors"
	"ipfs-share/network"
)


type GroupConnection struct {
	group interfaces.IGroup
	repo *fs.GroupRepo
	user interfaces.IUser
	addressBook *collections.ConcurrentCollection
	sessionClosed sesscommon.SessionClosedCallback
	p2p *P2PServer

	ipfs ipfsapi.IIpfs
	network network.INetwork

	channelState  chan []byte
	channelStop   chan bool

	groupSubscription ipfsapi.IPubSubSubscription
}

func NewGroupConnection(
	group interfaces.IGroup,
	repo *fs.GroupRepo,
	user interfaces.IUser,
	addressBook *collections.ConcurrentCollection,
	sessionClosed sesscommon.SessionClosedCallback,
	p2p *P2PServer,
	ipfs ipfsapi.IIpfs,
	network network.INetwork,
	) *GroupConnection {

	glog.Infof("Creating group connection...")

	conn := GroupConnection{
		group: group,
		repo: repo,
		user: user,
		addressBook: addressBook,
		sessionClosed: sessionClosed,
		p2p: p2p,
		channelStop: make(chan bool),
		ipfs: ipfs,
		network: network,
	}

	glog.Infof("subscribing to ipfs pubsub topic %s", group.Id().ToString())

	sub, err := ipfs.PubSubSubscribe(group.Id().ToString())
	if err != nil {
		glog.Errorf("could not ipfs subscribe to topic %s", group.Id().ToString())
		return nil
	}

	conn.groupSubscription = sub

	go conn.connectionListener()

	return &conn
}


func (conn *GroupConnection) Broadcast(msg []byte) error {
	id := conn.group.Id().Data().([32]byte)
	topic := base64.URLEncoding.EncodeToString(id[:])

	boxer := conn.group.Boxer()
	encMsg := boxer.BoxSeal(msg)
	msgString := base64.URLEncoding.EncodeToString(encMsg)

	if err := conn.ipfs.PubSubPublish(topic, msgString); err != nil {
		return errors.Wrap(err, "could not send group message to all members")
	}

	return nil
}

func (conn *GroupConnection) connectionListener() {
	glog.Infof("%s: GroupConnection for group '%s' is running...", conn.user.Name(), conn.group.Id().ToString())
	for {
		select {
		case <- conn.channelStop:
			{
				conn.groupSubscription.Cancel()
				close(conn.channelStop)
				return
			}
		default:
			{
				pubsubRecord, err := conn.groupSubscription.Next()
				if err != nil {
					glog.Warning("could not get next pubsub record")
					continue
				}

				glog.Infof("got a group message")

				encMsg, err := base64.URLEncoding.DecodeString((string)(pubsubRecord.Data))
				if err != nil {
					glog.Warningf("could not url decode group message: %s", err)
					continue
				}

				boxer := conn.group.Boxer()
				msgData, ok := boxer.BoxOpen(encMsg)
				if !ok {
					glog.Warningf("could not decrypt pubsub message")
					continue
				}

				msg, err := common.DecodeMessage(msgData)
				if err != nil {
					glog.Warning("could not decode pubsub record message")
					continue
				}

				if !conn.group.IsMember(msg.From) {
					glog.Warningf("non group member %v has written to the group channel", msg.From.Bytes())
					continue
				}

				contact, err := msg.Validate(conn.network, conn.ipfs)
				if err != nil {
					glog.Warningf("invalid pubsub message to group %v from user %v", conn.group.Id().Data(), msg.From.Bytes())
					continue
				}

				// TODO: check this with Ipfs address at the beginning
				address := conn.user.Address()
				if bytes.Equal(contact.Address.Bytes(), address.Bytes()) {
					continue
				}

				// append new contact to address book. if one already exists, therefore
				// it's P2P connection is not null, we will not try to create a new one
				// later
				if err := conn.addressBook.Append(contact); err != nil {
					glog.Warningf("could not append elem: %s", err)
				}
				contact = conn.addressBook.Get(contact.Id()).(*common.Contact)

				session, err := sessions.NewGroupSessionServer(
					msg,
					contact,
					conn.user,
					conn.group,
					conn.repo,
					conn.sessionClosed)
				if err != nil {
					glog.Error("could not create new group session server: %s", err)
					continue
				}

				conn.p2p.AddSession(session)
				go session.Run()
			}
		}
	}
}

func (conn *GroupConnection) Kill() {
	conn.channelStop <- true
}
