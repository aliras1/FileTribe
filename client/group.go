package client

import (
	"crypto/rand"
	"fmt"

	"ipfs-share/crypto"
	nw "ipfs-share/network"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"bytes"
	. "ipfs-share/collections"
	"encoding/base64"
)

type Group struct {
	Id                IIdentifier
	Name              string
	IpfsHash          string
	EncryptedIpfsHash string
	Members           []ethcommon.Address
	Boxer             crypto.SymmetricKey
}

func NewGroup(groupName string) *Group {
	var id [32]byte
	rand.Read(id[:])

	glog.Infof("group id: %s", base64.URLEncoding.EncodeToString(id[:]))

	var secretKeyBytes [32]byte
	rand.Read(secretKeyBytes[:])

	return &Group{
		Id:   NewBytesId(id),
		Name: groupName,
		Boxer: crypto.SymmetricKey{
			Key: secretKeyBytes,
			RNG: rand.Reader,
		},
		IpfsHash: "init_ipfs",
		EncryptedIpfsHash: "init_enc",
	}
}

func NewGroupFromId(groupId [32]byte, ctx *UserContext) error {
	_, members, _, err := ctx.Network.GetGroup(groupId)
	if err != nil {
		errors.Wrap(err, "could not get group from network")
	}

	for _, member := range members {
		if bytes.Equal(member.Bytes(), ctx.User.Address.Bytes()) {
			continue
		}
		c, err := ctx.Network.GetUser(member)
		if err != nil {
			glog.Warningf("could not get user in Group.GetKey(): %s", err)
			continue
		}

		if err := ctx.AddressBook.Append(NewContact(c, ctx.Ipfs)); err != nil {
			glog.Warningf("could not append elem: %s", err)
		}
		contact := ctx.AddressBook.Get(NewAddressId(&c.Address)).(*Contact)

		session := NewGetGroupKeySessionClient(groupId, contact, ctx)
		if err := ctx.P2P.sessions.Append(session); err != nil {
			glog.Warningf("could not append elem: %s", err)
		}

		go session.Run()
	}

	return nil
}

func (g *Group) Save(storage *Storage) error {
	cap := GroupAccessCap{
		GroupId: g.Id.Data().([32]byte),
		Boxer:   g.Boxer,
	}
	if err := cap.Save(storage); err != nil {
		return fmt.Errorf("could not store ga cap: Group.Save: %s", err)
	}
	return nil
}

func (g *Group) CreateOnNetwork(owner string, network nw.INetwork) error {
	glog.Infof("Registering group '%s'", g.Name)

	encIpfsPath := base64.URLEncoding.EncodeToString(g.Boxer.BoxSeal([]byte(g.IpfsHash)))
	g.EncryptedIpfsHash = encIpfsPath
	if err := network.CreateGroup(g.Id.Data().([32]byte), g.Name, encIpfsPath); err != nil {
		// TODO: check error message
		return fmt.Errorf("could not register group: Group.CreateOnNetwork: %s", err)
	}

	return nil
}
