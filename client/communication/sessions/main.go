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

package sessions

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	comcommon "github.com/aliras1/FileTribe/client/communication/common"
	"github.com/aliras1/FileTribe/client/communication/sessions/common"
	"github.com/aliras1/FileTribe/client/communication/sessions/servers"
	"github.com/aliras1/FileTribe/client/fs"
	"github.com/aliras1/FileTribe/client/interfaces"
)

// NewServerSession ...
func NewServerSession(
	msg *comcommon.Message,
	contact *comcommon.Contact,
	account ethcommon.Address,
	signer comcommon.Signer,
	callback common.CtxCallback,
	sessionClosed common.SessionClosedCallback,
) (common.Session, error) {

	switch msg.Type {
	case comcommon.GetGroupData:
		return servers.NewGetGroupDataSessionServer(
			msg,
			contact,
			account,
			signer,
			callback,
			sessionClosed)
	default:
		return nil, errors.New("invalid message type")
	}
}

// NewGroupServerSession ...
func NewGroupServerSession(
	msg *comcommon.Message,
	contact *comcommon.Contact,
	user interfaces.Account,
	group interfaces.Group,
	repo *fs.GroupRepo,
	sessionClosed common.SessionClosedCallback,
) (common.Session, error) {
	//switch msg.Type {
	//
	//case comcommon.Commit:
	//	{
	//		return servers.NewCommitChangesGroupSessionServer(
	//			msg,
	//			contact,
	//			user,
	//			group,
	//			repo,
	//			func(args []interface{}, groupId collections.IIdentifier) {},
	//			sessionClosed)
	//	}
	//default:
	//	{
	//		return nil, errors.New("invalid message type")
	//	}
	//}
	return nil, errors.New("not implemented")
}
