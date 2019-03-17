package communication

import (
	"bytes"
	"encoding/base64"
	"github.com/golang/glog"
	"github.com/pkg/errors"

	"github.com/aliras1/FileTribe/client/communication/common"
	"github.com/aliras1/FileTribe/client/communication/sessions"
	sesscommon "github.com/aliras1/FileTribe/client/communication/sessions/common"
	"github.com/aliras1/FileTribe/client/fs"
	"github.com/aliras1/FileTribe/client/interfaces"
	ipfsapi "github.com/aliras1/FileTribe/ipfs"
)


type GroupConnection struct {
	group         interfaces.IGroup
	repo          *fs.GroupRepo
	account       interfaces.IAccount
	addressBook   *common.AddressBook
	sessionClosed sesscommon.SessionClosedCallback
	p2p           *P2PManager

	ipfs ipfsapi.IIpfs

	channelState  chan []byte
	channelStop   chan bool

	groupSubscription ipfsapi.IPubSubSubscription
}

func NewGroupConnection(
	group interfaces.IGroup,
	repo *fs.GroupRepo,
	user interfaces.IAccount,
	addressBook *common.AddressBook,
	sessionClosed sesscommon.SessionClosedCallback,
	p2p *P2PManager,
	ipfs ipfsapi.IIpfs,
	) *GroupConnection {

	glog.Infof("Creating group connection...")

	conn := GroupConnection{
		group:         group,
		repo:          repo,
		account:       user,
		addressBook:   addressBook,
		sessionClosed: sessionClosed,
		p2p:           p2p,
		channelStop:   make(chan bool),
		ipfs:          ipfs,
	}

	glog.Infof("subscribing to ipfs pubsub topic %s", group.Address().String())

	sub, err := ipfs.PubSubSubscribe(group.Address().String())
	if err != nil {
		glog.Errorf("could not ipfs subscribe to topic %s", group.Address().String())
		return nil
	}

	conn.groupSubscription = sub

	go conn.connectionListener()

	return &conn
}


func (conn *GroupConnection) Broadcast(msg []byte) error {
	topic := conn.group.Address().String()

	boxer := conn.group.Boxer()
	encMsg := boxer.BoxSeal(msg)
	msgString := base64.URLEncoding.EncodeToString(encMsg)

	if err := conn.ipfs.PubSubPublish(topic, msgString); err != nil {
		return errors.Wrap(err, "could not send group message to all members")
	}

	return nil
}

func (conn *GroupConnection) connectionListener() {
	glog.Infof("%s: GroupConnection for group '%s' is running...", conn.account.Name(), conn.group.Address().String())
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
					glog.Warningf("could not get next pubsub record: %s", err)
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
					glog.Warningf("could not decode pubsub record message: %s", err)
					continue
				}

				if bytes.Equal(msg.From.Bytes(), conn.account.ContractAddress().Bytes()) {
					continue
				}

				if !conn.group.IsMember(msg.From) {
					glog.Warningf("non group member %v has written to the group channel", msg.From.Bytes())
					continue
				}

				contact, err := conn.addressBook.Get(msg.From)
				if err != nil {
					glog.Errorf("could not get contact from address book: %s", err)
					continue
				}

				if err := msg.Validate(contact); err != nil {
					glog.Warningf("invalid pubsub message to group %v from account %v", conn.group.Address().String(), msg.From.Bytes())
					continue
				}

				session, err := sessions.NewGroupServerSession(
					msg,
					contact,
					conn.account,
					conn.group,
					conn.repo,
					conn.sessionClosed)
				if err != nil {
					glog.Errorf("could not create new group session server: %s", err)
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
