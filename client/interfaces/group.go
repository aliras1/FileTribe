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

package interfaces

import (
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/aliras1/FileTribe/tribecrypto"
	"github.com/aliras1/FileTribe/tribecrypto/curves"
)

// Group is a mirror to the group data on the blockchain
type Group interface {
	Address() ethcommon.Address
	Name() string
	IpfsHash() string
	SetIpfsHash(encIpfsHash []byte) error
	EncryptedIpfsHash() []byte
	AddMember(user ethcommon.Address)
	RemoveMember(user ethcommon.Address)
	IsMember(user ethcommon.Address) bool
	CountMembers() int
	MemberOwners() []ethcommon.Address
	Boxer() tribecrypto.SymmetricKey
	CheckBoxer(newBoxer tribecrypto.SymmetricKey) error
	SetBoxer(boxer tribecrypto.SymmetricKey) error
	VerifyKey() curves.Point
	SetVerifyKey(vk curves.Point)
	SignKey() *big.Int
	SetSignKey(sk *big.Int)
	Update() error
	Save() error
}

// GroupData is the collection of group information
// that can be exported to disk
type GroupData struct {
	Address           ethcommon.Address
	Name              string
	IpfsHash          string
	EncryptedIpfsHash []byte

	MemberOwners []ethcommon.Address
	Boxer        tribecrypto.SymmetricKey

	VerifyKey curves.Point
	SignKey           *big.Int

}