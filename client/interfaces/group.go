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
	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/aliras1/FileTribe/tribecrypto"
)

type IGroup interface {
	Address() ethcommon.Address
	Name() string
	IpfsHash() string
	SetIpfsHash(encIpfsHash []byte) error
	EncryptedIpfsHash() []byte
	AddMember(user ethcommon.Address)
	RemoveMember(user ethcommon.Address)
	IsMember(user ethcommon.Address) bool
	CountMembers() int
	Members() []ethcommon.Address
	Boxer() tribecrypto.SymmetricKey
	SetBoxer(boxer tribecrypto.SymmetricKey)
	Update(name string, members []ethcommon.Address, encIpfsHash []byte) error
	Encode() ([]byte, error)
	Save() error
}