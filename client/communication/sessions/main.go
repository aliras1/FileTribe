package sessions

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	comcommon "github.com/aliras1/FileTribe/client/communication/common"
	"github.com/aliras1/FileTribe/client/communication/sessions/common"
	"github.com/aliras1/FileTribe/client/communication/sessions/servers"
	"github.com/aliras1/FileTribe/client/fs"
	"github.com/aliras1/FileTribe/client/interfaces"
	"github.com/aliras1/FileTribe/tribecrypto"
)



func NewServerSession(
	msg *comcommon.Message,
	contact *comcommon.Contact,
	account ethcommon.Address,
	signer *tribecrypto.Signer,
	callback common.CtxCallback,
	sessionClosed common.SessionClosedCallback,
) (common.ISession, error) {

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

func NewGroupServerSession(
	msg *comcommon.Message,
	contact *comcommon.Contact,
	user interfaces.IAccount,
	group interfaces.IGroup,
	repo *fs.GroupRepo,
	sessionClosed common.SessionClosedCallback,
) (common.ISession, error) {
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
