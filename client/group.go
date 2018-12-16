package client

import (
	"bytes"
	"crypto/rand"
	"ipfs-share/client/fs"
	"ipfs-share/client/fs/caps"
	"sync"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"ipfs-share/client/interfaces"
	"ipfs-share/crypto"
)

type Group struct {
	address           ethcommon.Address
	name              string
	ipfsHash          string
	encryptedIpfsHash []byte
	members           []ethcommon.Address
	boxer             crypto.SymmetricKey
	lock 			  sync.RWMutex
	storage           *fs.Storage
}

func NewGroup(address ethcommon.Address, groupName string, storage *fs.Storage) interfaces.IGroup {
	var secretKeyBytes [32]byte
	rand.Read(secretKeyBytes[:])
	boxer := crypto.SymmetricKey{
		Key: secretKeyBytes,
		RNG: rand.Reader,
	}

	return &Group{
		address:	address,
		name: 		groupName,
		boxer: 		boxer,
		storage:    storage,
	}
}

func NewGroupFromCap(cap *caps.GroupAccessCap, storage *fs.Storage) interfaces.IGroup {
	return &Group {
		address: 	cap.Address,
		boxer: 		cap.Boxer,
		storage:	storage,
	}
}

func NewGroupFromId(groupId [32]byte, ctx *UserContext) error {
	//_, members, _, _, err := ctx.network.GetGroup(groupId)
	//if err != nil {
	//	errors.Wrap(err, "could not get group from network")
	//}
	//
	//for _, member := range members {
	//	if bytes.Equal(member.Bytes(), ctx.account.Address().Bytes()) {
	//		continue
	//	}
	//	c, err := ctx.network.GetUser(member)
	//	if err != nil {
	//		glog.Warningf("could not get account in Group.GetKey(): %s", err)
	//		continue
	//	}
	//
	//	if err := ctx.addressBook.Append(comcommon.NewContact(c, ctx.ipfs)); err != nil {
	//		glog.Warningf("could not append elem: %s", err)
	//	}
	//	contact := ctx.addressBook.Get(NewAddressId(&c.Address)).(*comcommon.Contact)
	//
	//	err = ctx.p2p.StartGetGroupKeySession(groupId, contact, ctx.account, ctx.storage, func(cap *caps.GroupAccessCap) {
	//		groupCtx, err := NewGroupContextFromCAP(
	//			cap,
	//			ctx.account,
	//			ctx.p2p,
	//			ctx.addressBook,
	//			ctx.network,
	//			ctx.ipfs,
	//			ctx.storage,
	//			ctx.transactions,
	//		)
	//		if err != nil {
	//			glog.Error(err, "could not create group context")
	//			return
	//		}
	//
	//		if err := groupCtx.Update(); err != nil {
	//			glog.Error(err, "could not update group")
	//			return
	//		}
	//
	//		if err := ctx.groups.Put(groupCtx); err != nil {
	//			glog.Warningf("could not append elem: %s", err)
	//		}
	//	})
	//	if err != nil {
	//		glog.Errorf("could not start get group key session: %s", err)
	//	}
	//}

	return errors.New("not implemented")
}

func (g *Group) Encode() ([]byte, error) {
	g.lock.RLock()
	defer g.lock.RUnlock()

	cap := caps.GroupAccessCap{
		Address: g.address,
		Boxer:   g.boxer,
	}

	data, err := cap.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode group data")
	}

	return data, nil
}

func (g *Group) Save() error {
	//data, err := g.Encode()
	//if err != nil {
	//	return errors.Wrap(err, "could not encode group")
	//}

	cap := caps.GroupAccessCap{
		Address: g.address,
		Boxer:   g.boxer,
	}

	if err := g.storage.SaveGroupAccessCap(&cap); err != nil {
		return errors.Wrap(err, "could not save group cap")
	}

	return nil
}

// Note that it is not enough to only receive ipfsHash and
// encrypt it and use that as encryptedIpfsHash because
// the encryption uses a random element when producing the
// cipher text, therefore on each instance of the ipfs-share
// daemon, the ecnryptedIpfsHash's will be different
func (g *Group) SetIpfsHash(encIpfsHash []byte) error {
	g.lock.Lock()
	defer g.lock.Unlock()

	ipfsHash, ok := g.boxer.BoxOpen(encIpfsHash)
	if !ok {
		return errors.New("could not decrypt encrypted ipfs hash")
	}

	g.ipfsHash = string(ipfsHash)
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

func (g *Group) IsMember(account ethcommon.Address) bool {
	g.lock.RLock()
	defer g.lock.RUnlock()

	return g.isMember(account)
}

func (g *Group) isMember(user ethcommon.Address) bool {
	for _, m := range g.members {
		if bytes.Equal(m.Bytes(), user.Bytes()) {
			return true
		}
	}
	return false
}

func (g *Group) AddMember(account ethcommon.Address) {
	g.lock.Lock()
	defer g.lock.Unlock()

	if !g.isMember(account) {
		g.members = append(g.members, account)
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

func (g *Group) Address() ethcommon.Address {
	g.lock.RLock()
	defer g.lock.RUnlock()

	return g.address
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

func (g *Group) SetBoxer(boxer crypto.SymmetricKey) {
	g.lock.Lock()
	defer g.lock.Unlock()

	g.boxer = boxer
}