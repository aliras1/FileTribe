package sessions

import (
	"github.com/pkg/errors"
	"ipfs-share/client/communication/sessions/servers"

	comcommon "ipfs-share/client/communication/common"
	"ipfs-share/client/communication/sessions/common"
	"ipfs-share/client/fs"
	"ipfs-share/client/interfaces"
)



func NewServerSession(
	msg *comcommon.Message,
	contact *comcommon.Contact,
	user interfaces.IUser,
	getGroup common.GetGroupCallback,
	closedChan chan common.ISession,
	) (common.ISession, error) {

	switch msg.Type {
	case comcommon.GetGroupKey:
		{
			return servers.NewGetGroupKeySessionServer(msg, contact, user, getGroup, closedChan)
		}

	default:
		return nil, errors.New("invalid message type")
	}
}

func NewGroupSessionServer(
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
			return servers.NewCommitChangesGroupSessionServer(msg, contact, user, group, repo, sessionClosed)
		}
	default:
		{
			return nil, errors.New("invalid message type")
		}
	}
}

