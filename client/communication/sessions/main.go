package sessions

import (
	"github.com/pkg/errors"
	"ipfs-share/client/fs"
	"ipfs-share/collections"

	comcommon "ipfs-share/client/communication/common"
	"ipfs-share/client/communication/sessions/common"
	"ipfs-share/client/communication/sessions/servers"
	"ipfs-share/client/interfaces"
)



func NewServerSession(
	msg *comcommon.Message,
	contact *comcommon.Contact,
	user interfaces.IUser,
	callback common.CtxCallback,
	sessionClosed common.SessionClosedCallback,
) (common.ISession, error) {
	switch msg.Type {
	case comcommon.GetGroupKey:
		{
			return servers.NewGetGroupKeySessionServer(
				msg,
				contact,
				user,
				callback.GetGroupData,
				sessionClosed)
		}
	case comcommon.ChangeKey:
		{
			return servers.NewChangeGroupKeySessionServer(
				msg,
				contact,
				user,
				callback.GetGroupData,
				callback.OnChangeGroupKeyServerSessionSuccess,
				sessionClosed)
		}
	default:
		{
			return nil, errors.New("invalid message type")
		}
	}
}

func NewGroupServerSession(
	msg *comcommon.Message,
	contact *comcommon.Contact,
	user interfaces.IUser,
	group interfaces.IGroup,
	repo *fs.GroupRepo,
	sessionClosed common.SessionClosedCallback,
) (common.ISession, error) {
	switch msg.Type {

	case comcommon.Commit:
		{
			return servers.NewCommitChangesGroupSessionServer(
				msg,
				contact,
				user,
				group,
				repo,
				func(args []interface{}, groupId collections.IIdentifier) {},
				sessionClosed)
		}
	default:
		{
			return nil, errors.New("invalid message type")
		}
	}
}
