package email_accounts_api

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
)

func (e *EmailAccountsImpl) DeleteEmailAccount(ctx context.Context, req *pb.DeleteEmailAccount_Request) (*pb.DeleteEmailAccount_Response, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse account id")
	}

	if err := e.emailSvc.DeleteAccount(ctx, id); err != nil {
		return nil, rerrors.Wrap(err, "delete email account")
	}

	return &pb.DeleteEmailAccount_Response{}, nil
}
