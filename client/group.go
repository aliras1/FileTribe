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
	"encoding/json"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/glog"
	"github.com/pkg/errors"

	"github.com/aliras1/FileTribe/client/fs"
	"github.com/aliras1/FileTribe/client/fs/meta"
	"github.com/aliras1/FileTribe/client/interfaces"
	ethgroup "github.com/aliras1/FileTribe/eth/gen/Group"
	"github.com/aliras1/FileTribe/tribecrypto"
	"github.com/aliras1/FileTribe/tribecrypto/curves"
)

type group struct {
	data 			  *interfaces.GroupData
	lock              sync.RWMutex
	storage           *fs.Storage
	contract          *ethgroup.Group
}

// NewGroupFromMeta creates a new Group from an existing GroupMeta file stored on disk
func NewGroupFromMeta(meta *meta.GroupMeta, groupContract *ethgroup.Group, storage *fs.Storage) interfaces.Group {
	group := &group{
		storage: storage,
		contract: groupContract,
		data: &interfaces.GroupData{
			Address:meta.Address,
			Boxer:meta.Boxer,
			VerifyKey:meta.VerifyKey,
		},
	}

	if err := group.Update(); err != nil {
		glog.Warningf("could not update group: %s", err)
	}

	return group
}

// NewGroupFromGroupData creates a new group from GroupData and tries to refresh its
// data with up to date information from the blockchain
func NewGroupFromGroupData(data *interfaces.GroupData, groupContract *ethgroup.Group, storage *fs.Storage) interfaces.Group {
	group := &group{
		storage:storage,
		data:data,
		contract:groupContract,
	}

	if err := group.Update(); err != nil {
		glog.Warningf("could not update group: %s", err)
	}

	return group
}

// GetGroupKeyFromAddress tries to get the group key of a group with the given address
func GetGroupKeyFromAddress(address ethcommon.Address, ctx *UserContext) error {
	//_, members, _, _, err := ctx.network.GetGroup(groupId)
	//if err != nil {
	//	errors.Wrap(err, "could not get group from network")
	//}
	//
	//for _, member := range members {
	//	if bytes.Equal(member.Bytes(), ctx.account.EthAccountAddress().Bytes()) {
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
	//	contact := ctx.addressBook.Get(NewAddressId(&c.EthAccountAddress)).(*comcommon.Contact)
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

// Save saves the Group to disk
func (g *group) Save() error {
	dataEnc, err := json.Marshal(g.data)
	if err != nil {
		return errors.Wrap(err, "could not encode group data")
	}

	if err := g.storage.SaveGroupData(dataEnc, g.data.Address.String()); err != nil {
		return errors.Wrap(err, "could not save group groupMeta")
	}

	return nil
}

// SetIpfsHash decrypts an encrypted IPFS hash and stores both
// the decrypted one and the encrypted one.
// Note that it is not enough to only receive ipfsHash and
// encrypt it and use that as encryptedIpfsHash because
// the encryption uses a random element when producing the
// cipher text, therefore on each instance of the ipfs-share
// daemon, the ecnryptedIpfsHashes will be different
func (g *group) SetIpfsHash(encIpfsHash []byte) error {
	g.lock.Lock()
	defer g.lock.Unlock()

	ipfsHash, ok := g.data.Boxer.BoxOpen(encIpfsHash)
	if !ok {
		return errors.New("could not decrypt encrypted ipfs hash")
	}

	g.data.IpfsHash = string(ipfsHash)
	g.data.EncryptedIpfsHash = encIpfsHash

	return nil
}

// Update updates the Group data with the provided parameters
func (g *group) Update() error {
	g.lock.Lock()
	defer g.lock.Unlock()

	name, err := g.contract.Name(&bind.CallOpts{Pending: true})
	if err != nil {
		return errors.Wrap(err, "could not get group name")
	}

	memberOwners, err := g.contract.MemberOwners(&bind.CallOpts{Pending: true})
	if err != nil {
		return errors.Wrap(err, "could not get group member owners")
	}

	encIpfsHash, err := g.contract.IpfsHash(&bind.CallOpts{Pending: true})
	if err != nil {
		return errors.Wrap(err, "could not get group ipfs hash")
	}

	if len(encIpfsHash) > 0 {
		ipfsHash, ok := g.data.Boxer.BoxOpen(encIpfsHash)
		if !ok {
			return errors.New("could not decrypt ipfs hash")
		}
		g.data.IpfsHash = string(ipfsHash)
	}

	g.data.Name = name
	g.data.MemberOwners = memberOwners
	g.data.EncryptedIpfsHash = encIpfsHash

	return nil
}

// IsMember checks if an account is a group member or not
func (g *group) IsMember(memberOwner ethcommon.Address) bool {
	g.lock.RLock()
	defer g.lock.RUnlock()

	return g.isMember(memberOwner)
}

func (g *group) isMember(memberOwner ethcommon.Address) bool {
	for _, m := range g.data.MemberOwners {
		if bytes.Equal(m.Bytes(), memberOwner.Bytes()) {
			return true
		}
	}
	return false
}

// AddMember adds an account to the member list
func (g *group) AddMember(accountOwner ethcommon.Address) {
	g.lock.Lock()
	defer g.lock.Unlock()

	if !g.isMember(accountOwner) {
		g.data.MemberOwners = append(g.data.MemberOwners, accountOwner)
	}
}

// RemoveMember removes an account from the member list
func (g *group) RemoveMember(account ethcommon.Address) {
	g.lock.Lock()
	defer g.lock.Unlock()

	for i, m := range g.data.MemberOwners {
		if bytes.Equal(m.Bytes(), account.Bytes()) {
			g.data.MemberOwners = append(g.data.MemberOwners[:i], g.data.MemberOwners[i+1:]...)
			return
		}
	}
}

// MemberOwners returns the list of members
func (g *group) MemberOwners() []ethcommon.Address {
	g.lock.RLock()
	defer g.lock.RUnlock()

	var list []ethcommon.Address
	for _, member := range g.data.MemberOwners {
		list = append(list, member)
	}

	return list
}

// CountMembers returns the number of members
func (g *group) CountMembers() int {
	g.lock.RLock()
	defer g.lock.RUnlock()

	return len(g.data.MemberOwners)
}

// Address returns the group's smart contract address
func (g *group) Address() ethcommon.Address {
	g.lock.RLock()
	defer g.lock.RUnlock()

	return g.data.Address
}

// Name returns the group name
func (g *group) Name() string {
	g.lock.RLock()
	defer g.lock.RUnlock()

	return g.data.Name
}

// IpfsHash returns the group's current IPFS hash
func (g *group) IpfsHash() string {
	g.lock.RLock()
	defer g.lock.RUnlock()

	return g.data.IpfsHash
}

// EncryptedIpfsHash returns the group's current encrypted IPFS hash
func (g *group) EncryptedIpfsHash() []byte {
	g.lock.RLock()
	defer g.lock.RUnlock()

	return g.data.EncryptedIpfsHash
}

// Boxer returns the group key
func (g *group) Boxer() tribecrypto.SymmetricKey {
	g.lock.RLock()
	defer g.lock.RUnlock()

	return g.data.Boxer
}

// SetBoxer is a setter for the group key
func (g *group) SetBoxer(boxer tribecrypto.SymmetricKey) error {
	g.lock.Lock()
	defer g.lock.Unlock()

	encIpfsHash, err := g.contract.IpfsHash(&bind.CallOpts{Pending:true})
	if err != nil {
		return errors.Wrap(err, "could not get the most recent ipfs hash of the group")
	}

	_, ok := boxer.BoxOpen(encIpfsHash)
	if !ok {
		return errors.New("could not decrypt encIpfsHash")
	}

	g.data.EncryptedIpfsHash = encIpfsHash
	g.data.Boxer = boxer

	return nil
}

func (g* group) CheckBoxer(newBoxer tribecrypto.SymmetricKey) error {
	g.lock.RLock()
	defer g.lock.RUnlock()

	encIpfsHash, err := g.contract.IpfsHash(&bind.CallOpts{Pending:true})
	if err != nil {
		return errors.Wrap(err, "could not get the most recent ipfs hash of the group")
	}

	_, ok := newBoxer.BoxOpen(encIpfsHash)
	if !ok {
		return errors.New("could not decrypt encIpfsHash")
	}

	return nil
}

func (g *group) VerifyKey() curves.Point {
	return g.data.VerifyKey
}


func (g *group) SetVerifyKey(vk curves.Point) {
	g.lock.RLock()
	defer g.lock.RUnlock()

	g.data.VerifyKey = vk
}

func (g *group) SignKey() *big.Int {
	g.lock.RLock()
	defer g.lock.RUnlock()

	return g.data.SignKey
}

func (g *group) SetSignKey(sk *big.Int) {
	g.lock.Lock()
	defer g.lock.Unlock()

	g.data.SignKey = sk
}