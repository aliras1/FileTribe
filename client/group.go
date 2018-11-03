package client

import (
	"crypto/rand"
	"fmt"

	"ipfs-share/crypto"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"bytes"
	. "ipfs-share/collections"
	"encoding/base64"
	"sync"
	"strings"
	"ipfs-share/client/fs"
)

type IGroup interface {
	Id() IIdentifier
	Name() string
	IpfsHash() string
	SetIpfsHash(ipfsHash string, encIpfsHash []byte) error
	EncryptedIpfsHash() []byte
	AddMember(user ethcommon.Address)
	RemoveMember(user ethcommon.Address)
	IsMember(user ethcommon.Address) bool
	CountMembers() int
	Members() []ethcommon.Address
	Boxer() crypto.SymmetricKey
	Update(name string, members []ethcommon.Address, encIpfsHash []byte) error
	Save(storage *fs.Storage) error
}

type Group struct {
	id                IIdentifier
	name              string
	ipfsHash          string
	encryptedIpfsHash []byte
	members           []ethcommon.Address
	boxer             crypto.SymmetricKey
	lock sync.RWMutex
}

func NewGroup(groupName string) IGroup {
	var id [32]byte
	rand.Read(id[:])

	glog.Infof("group id: %s", base64.URLEncoding.EncodeToString(id[:]))

	var secretKeyBytes [32]byte
	rand.Read(secretKeyBytes[:])
	boxer := crypto.SymmetricKey{
		Key: secretKeyBytes,
		RNG: rand.Reader,
	}

	return &Group{
		id:   NewBytesId(id),
		name: groupName,
		boxer: boxer,
	}
}

func NewGroupFromCap(cap *fs.GroupAccessCap) IGroup {
	return &Group {
		id: NewBytesId(cap.GroupId),
		boxer: cap.Boxer,
	}
}

func NewGroupFromId(groupId [32]byte, ctx *UserContext) error {
	_, members, _, err := ctx.network.GetGroup(groupId)
	if err != nil {
		errors.Wrap(err, "could not get group from network")
	}

	for _, member := range members {
		if bytes.Equal(member.Bytes(), ctx.user.Address().Bytes()) {
			continue
		}
		c, err := ctx.network.GetUser(member)
		if err != nil {
			glog.Warningf("could not get user in Group.GetKey(): %s", err)
			continue
		}

		if err := ctx.addressBook.Append(NewContact(c, ctx.ipfs)); err != nil {
			glog.Warningf("could not append elem: %s", err)
		}
		contact := ctx.addressBook.Get(NewAddressId(&c.Address)).(*Contact)

		session := NewGetGroupKeySessionClient(groupId, contact, ctx.user, ctx.storage, ctx.p2p.SessionClosedChan, func(cap *fs.GroupAccessCap) {
			groupCtx, err := NewGroupContextFromCAP(
				cap,
				ctx.user,
				ctx.p2p,
				ctx.addressBook,
				ctx.network,
				ctx.ipfs,
				ctx.storage,
				ctx.transactions,
			)
			if err != nil {
				glog.Error(err, "could not create group context")
				return
			}

			if err := groupCtx.Update(); err != nil {
				glog.Error(err, "could not update group")
				return
			}

			if err := ctx.groups.Append(groupCtx); err != nil {
				glog.Warningf("could not append elem: %s", err)
			}
		})
		if err := ctx.p2p.sessions.Append(session); err != nil {
			glog.Warningf("could not append elem: %s", err)
		}

		go session.Run()
	}

	return nil
}

func (g *Group) Save(storage *fs.Storage) error {
	g.lock.RLock()
	defer g.lock.RUnlock()

	cap := fs.GroupAccessCap{
		GroupId: g.id.Data().([32]byte),
		Boxer:   g.boxer,
	}
	if err := cap.Save(storage); err != nil {
		return fmt.Errorf("could not store ga cap: Group.SaveMetadata: %s", err)
	}
	return nil
}

// Note that it is not enough to only receive ipfsHash and
// encrypt it and use that as encryptedIpfsHash because
// the encryption uses a random element when producing the
// cipher text, therefore on each instance of the ipfs-share
// daemon, the ecnryptedIpfsHash's will be different
func (g *Group) SetIpfsHash(ipfsHash string, encIpfsHash []byte) error {
	g.lock.Lock()
	defer g.lock.Unlock()

	decryptedIpfsHash, ok := g.boxer.BoxOpen(encIpfsHash)
	if !ok {
		return errors.New("could not decrypt encrypted ipfs hash")
	}
	if strings.Compare(string(decryptedIpfsHash), ipfsHash) != 0 {
		return errors.New("decrypt(encIpfsHash, group_key) != ipfsHash")
	}

	g.ipfsHash = ipfsHash
	g.encryptedIpfsHash = encIpfsHash

	return nil
}

func (g *Group) Update(name string, members []ethcommon.Address, encIpfsHash []byte) error {
	g.lock.Lock()
	defer g.lock.Unlock()

	ipfsHash, ok := g.boxer.BoxOpen(encIpfsHash)
	if !ok {
		return errors.New("could not decrypt ipfs hash")
	}

	g.name = name
	g.members = members
	g.ipfsHash = string(ipfsHash)
	g.encryptedIpfsHash = encIpfsHash

	return nil
}

func (g *Group) IsMember(user ethcommon.Address) bool {
	g.lock.RLock()
	defer g.lock.RUnlock()

	for _, m := range g.members {
		if bytes.Equal(m.Bytes(), user.Bytes()) {
			return true
		}
	}
	return false
}

func (g *Group) AddMember(user ethcommon.Address) {
	g.lock.Lock()
	defer g.lock.Unlock()

	if !g.IsMember(user) {
		g.members = append(g.members, user)
	}
}

func (g *Group) RemoveMember(user ethcommon.Address) {
	g.lock.Lock()
	defer g.lock.Unlock()

	for i, m := range g.members {
		if bytes.Equal(m.Bytes(), user.Bytes()) {
			g.members = append(g.members[:i], g.members[i+1:]...)
			return
		}
	}
}

func (g *Group) Members() []ethcommon.Address {
	g.lock.RLock()
	defer g.lock.RUnlock()

	var list []ethcommon.Address
	for _, member := range g.members {
		list = append(list, member)
	}

	return list
}

func (g *Group) CountMembers() int {
	g.lock.RLock()
	defer g.lock.RUnlock()

	return len(g.members)
}

func (g *Group) Id() IIdentifier {
	g.lock.RLock()
	defer g.lock.RUnlock()

	return g.id
}

func (g *Group) Name() string {
	g.lock.RLock()
	defer g.lock.RUnlock()

	return g.name
}

func (g *Group) IpfsHash() string {
	g.lock.RLock()
	defer g.lock.RUnlock()

	return g.ipfsHash
}

func (g *Group) EncryptedIpfsHash() []byte {
	g.lock.RLock()
	defer g.lock.RUnlock()

	return g.encryptedIpfsHash
}

func (g *Group) Boxer() crypto.SymmetricKey {
	g.lock.RLock()
	defer g.lock.RUnlock()

	return g.boxer
}
