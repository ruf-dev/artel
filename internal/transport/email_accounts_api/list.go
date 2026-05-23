package email_accounts_api

import (
	"context"

	"go.redsock.ru/rerrors"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

func (e *EmailAccountsImpl) ListEmailAccounts(ctx context.Context, _ *pb.ListEmailAccounts_Request) (*pb.ListEmailAccounts_Response, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return nil, user_errors.Unauthenticated
	}

	accounts, err := e.emailSvc.ListAccounts(ctx, uc.UserUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "list email accounts")
	}

	items := make([]*pb.EmailAccountInfo, len(accounts))
	for i, a := range accounts {
		items[i] = toProto(a)
	}

	return &pb.ListEmailAccounts_Response{Accounts: items}, nil
}
