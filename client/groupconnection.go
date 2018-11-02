package client

import (
	"github.com/golang/glog"
	ipfsapi "ipfs-share/ipfs"
	"bytes"
	"encoding/base64"
	"github.com/pkg/errors"
)


type GroupConnection struct {
	groupCtx *GroupContext

	channelState  chan []byte
	channelStop   chan bool

	groupSubscription ipfsapi.IPubSubSubscription
}

func NewGroupConnection(groupCtx *GroupContext) *GroupConnection {
	glog.Infof("Creating group connection...")

	conn := GroupConnection{
		groupCtx: groupCtx,
		channelStop: make(chan bool),
	}

	glog.Infof("%s: subscribing to ipfs pubsub topic %s", groupCtx.User.Name, groupCtx.Group.Id().ToString())

	sub, err := groupCtx.Ipfs.PubSubSubscribe(groupCtx.Group.Id().ToString())
	if err != nil {
		glog.Errorf("%s: could not ipfs subscribe to topic %s", groupCtx.User.Name, groupCtx.Group.Id().ToString())
		return nil
	}

	conn.groupSubscription = sub

	go conn.connectionListener()

	return &conn
}


func (conn *GroupConnection) SendAll(msg []byte) error {
	id := conn.groupCtx.Group.Id().Data().([32]byte)
	topic := base64.URLEncoding.EncodeToString(id[:])

	boxer := conn.groupCtx.Group.Boxer()
	encMsg := boxer.BoxSeal(msg)
	msgString := base64.URLEncoding.EncodeToString(encMsg)

	if err := conn.groupCtx.Ipfs.PubSubPublish(topic, msgString); err != nil {
		return errors.Wrap(err, "could not send group message to all members")
	}

	return nil
}

func (conn *GroupConnection) connectionListener() {
	glog.Infof("GroupConnection for user '%s' group '%s' is running...", conn.groupCtx.User.Name, conn.groupCtx.Group.Id().ToString())
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

				glog.Infof("%s got a group message", conn.groupCtx.User.Name)

				encMsg, err := base64.URLEncoding.DecodeString((string)(pubsubRecord.Data))
				if err != nil {
					glog.Warningf("could not url decode group message: %s", err)
					continue
				}

				boxer := conn.groupCtx.Group.Boxer()
				msgData, ok := boxer.BoxOpen(encMsg)
				if !ok {
					glog.Warningf("could not decrypt pubsub message")
					continue
				}

				msg, err := DecodeMessage(msgData)
				if err != nil {
					glog.Warning("could not decode pubsub record message")
					continue
				}

				if !conn.groupCtx.Group.IsMember(msg.From) {
					glog.Warningf("non group member %v has written to the group channel", msg.From.Bytes())
					continue
				}

				contact, err := msg.Validate(conn.groupCtx.Network, conn.groupCtx.Ipfs)
				if err != nil {
					glog.Warningf("invalid pubsub message to group %v from user %v", conn.groupCtx.Group.Id().Data(), msg.From.Bytes())
					continue
				}

				// TODO: check this with Ipfs address at the beginning
				if bytes.Equal(contact.Address.Bytes(), conn.groupCtx.User.Address.Bytes()) {
					continue
				}

				// append new contact to address book. if one already exists, therefore
				// it's P2P connection is not null, we will not try to create a new one
				// later
				if err := conn.groupCtx.AddressBook.Append(contact); err != nil {
					glog.Warningf("could not append elem: %s", err)
				}
				contact = conn.groupCtx.AddressBook.Get(contact.Id()).(*Contact)

				session := NewGroupSessionServer(msg, contact, conn.groupCtx)
				conn.groupCtx.P2P.AddSession(session)
				go session.Run()
			}
		}
	}
}

func (conn *GroupConnection) Kill() {
	conn.channelStop <- true
}
