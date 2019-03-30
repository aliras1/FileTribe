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

package common

import (
	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/aliras1/FileTribe/client/fs"
	"github.com/aliras1/FileTribe/client/interfaces"
	"github.com/aliras1/FileTribe/collections"
	"github.com/aliras1/FileTribe/tribecrypto"
)

type OnGetGroupKeySuccessCallback func(address ethcommon.Address, boxer tribecrypto.SymmetricKey)

type OnClientSuccessCallback func(args []interface{})

type OnServerSuccessCallback func(args []interface{}, groupId collections.IIdentifier)

type SessionClosedCallback func(session ISession)

type GetGroupDataCallback func(addr ethcommon.Address) (interfaces.IGroup, *fs.GroupRepo)

type Broadcast func(msg []byte) error

type CtxCallback interface {
	IsMember(group ethcommon.Address, account ethcommon.Address) error

	GetBoxerOfGroup(group ethcommon.Address) (tribecrypto.SymmetricKey, error)

	GetProposedBoxerOfGroup(group ethcommon.Address, proposer ethcommon.Address) (tribecrypto.SymmetricKey, error)
}
