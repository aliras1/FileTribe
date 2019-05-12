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
	"github.com/aliras1/FileTribe/tribecrypto"
)

// OnGetGroupKeySuccessCallback is called when a group key is retrieved successfully
type OnGetGroupKeySuccessCallback func(proposalKey []byte, boxer tribecrypto.SymmetricKey)

// SessionClosedCallback is called when a session is closed
type SessionClosedCallback func(session ISession)

// GetGroupDataCallback is used for retrieving group data
type GetGroupDataCallback func(addr ethcommon.Address) (interfaces.Group, *fs.GroupRepo)

// Broadcast ...
type Broadcast func(msg []byte) error

// CtxCallback is used by sessions to get specific group data
// without access to all GroupContext functions and data
type CtxCallback interface {
	// IsMember returns if an account is a member of a group
	IsMember(group ethcommon.Address, account ethcommon.Address) error

	// GetBoxerOfGroup returns the current key of a group
	GetBoxerOfGroup(group ethcommon.Address) (tribecrypto.SymmetricKey, error)

	// GetProposedBoxerOfGroup returns the proposed group key of a group
	// that was suggested the given member
	GetProposedBoxerOfGroup(group ethcommon.Address, proposalKey []byte) (tribecrypto.SymmetricKey, error)
}
