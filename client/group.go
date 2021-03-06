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

package client

import (
	"bytes"
	"crypto/rand"
	"sync"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/aliras1/FileTribe/client/fs"
	"github.com/aliras1/FileTribe/client/fs/meta"
	"github.com/aliras1/FileTribe/client/interfaces"
	"github.com/aliras1/FileTribe/tribecrypto"
)

// Group is the mirror object of a Group smart contract
type Group struct {
	address           ethcommon.Address
	name              string
	ipfsHash          string
	encryptedIpfsHash []byte
	members           []ethcommon.Address
	boxer             tribecrypto.SymmetricKey
	lock              sync.RWMutex
	storage           *fs.Storage
}

// NewGroup creates a new Group with the given parameters
func NewGroup(address ethcommon.Address, groupName string, storage *fs.Storage) interfaces.IGroup {
	var secretKeyBytes [32]byte
	rand.Read(secretKeyBytes[:])
	boxer := tribecrypto.SymmetricKey{
		Key: secretKeyBytes,
		RNG: rand.Reader,
	}

	return &Group{
		address: address,
		name:    groupName,
		boxer:   boxer,
		storage: storage,
	}
}

// NewGroupFromMeta creates a new Group from an existing GroupMeta file stored on disk
func NewGroupFromMeta(meta *meta.GroupMeta, storage *fs.Storage) interfaces.IGroup {
	return &Group{
		address: meta.Address,
		boxer:   meta.Boxer,
		storage: storage,
	}
}

// GetGroupKeyFromAddress tries to get the group key of a group with the given address
func GetGroupKeyFromAddress(address ethcommon.Address, ctx *UserContext) error {
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
	//	err = ctx.p2p.StartGetGroupKeySession(groupId, contact, ctx.account, ctx.storage, func(cap *caps.GroupMeta) {
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

// Encode encodes the Group
func (g *Group) Encode() ([]byte, error) {
	g.lock.RLock()
	defer g.lock.RUnlock()

	cap := meta.GroupMeta{
		Address: g.address,
		Boxer:   g.boxer,
	}

	data, err := cap.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode group data")
	}

	return data, nil
}

// Save saves the Group to disk
func (g *Group) Save() error {
	//data, err := g.Encode()
	//if err != nil {
	//	return errors.Wrap(err, "could not encode group")
	//}

	cap := meta.GroupMeta{
		Address: g.address,
		Boxer:   g.boxer,
	}

	if err := g.storage.SaveGroupMeta(&cap); err != nil {
		return errors.Wrap(err, "could not save group cap")
	}

	return nil
}

// SetIpfsHash decrypts an encrypted IPFS hash and stores both
// the decrypted one and the encrypted one.
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

// Update updates the Group data with the provided parameters
func (g *Group) Update(name string, members []ethcommon.Address, encIpfsHash []byte) error {
	g.lock.Lock()
	defer g.lock.Unlock()

	if len(encIpfsHash) > 0 {
		ipfsHash, ok := g.boxer.BoxOpen(encIpfsHash)
		if !ok {
			return errors.New("could not decrypt ipfs hash")
		}
		g.ipfsHash = string(ipfsHash)
	}

	g.name = name
	g.members = members
	g.encryptedIpfsHash = encIpfsHash

	return nil
}

// IsMember checks if an account is a group member or not
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

// AddMember adds an account to the member list
func (g *Group) AddMember(account ethcommon.Address) {
	g.lock.Lock()
	defer g.lock.Unlock()

	if !g.isMember(account) {
		g.members = append(g.members, account)
	}
}

// RemoveMember removes an account from the member list
func (g *Group) RemoveMember(account ethcommon.Address) {
	g.lock.Lock()
	defer g.lock.Unlock()

	for i, m := range g.members {
		if bytes.Equal(m.Bytes(), account.Bytes()) {
			g.members = append(g.members[:i], g.members[i+1:]...)
			return
		}
	}
}

// Members returns the list of members
func (g *Group) Members() []ethcommon.Address {
	g.lock.RLock()
	defer g.lock.RUnlock()

	var list []ethcommon.Address
	for _, member := range g.members {
		list = append(list, member)
	}

	return list
}

// CountMembers returns the number of members
func (g *Group) CountMembers() int {
	g.lock.RLock()
	defer g.lock.RUnlock()

	return len(g.members)
}

// Address returns the group's smart contract address
func (g *Group) Address() ethcommon.Address {
	g.lock.RLock()
	defer g.lock.RUnlock()

	return g.address
}

// Name returns the group name
func (g *Group) Name() string {
	g.lock.RLock()
	defer g.lock.RUnlock()

	return g.name
}

// IpfsHash returns the group's current IPFS hash
func (g *Group) IpfsHash() string {
	g.lock.RLock()
	defer g.lock.RUnlock()

	return g.ipfsHash
}

// EncryptedIpfsHash returns the group's current encrypted IPFS hash
func (g *Group) EncryptedIpfsHash() []byte {
	g.lock.RLock()
	defer g.lock.RUnlock()

	return g.encryptedIpfsHash
}

// Boxer returns the group key
func (g *Group) Boxer() tribecrypto.SymmetricKey {
	g.lock.RLock()
	defer g.lock.RUnlock()

	return g.boxer
}

// SetBoxer is a setter for the group key
func (g *Group) SetBoxer(boxer tribecrypto.SymmetricKey) {
	g.lock.Lock()
	defer g.lock.Unlock()

	g.boxer = boxer
}
