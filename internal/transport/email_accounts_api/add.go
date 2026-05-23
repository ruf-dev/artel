package email_accounts_api

import (
	"context"

	"go.redsock.ru/rerrors"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

func (e *EmailAccountsImpl) AddEmailAccount(ctx context.Context, req *pb.AddEmailAccount_Request) (*pb.AddEmailAccount_Response, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return nil, user_errors.Unauthenticated
	}

	account := domain.EmailAccount{
		UserUuid: uc.UserUuid,
		Email:    req.Email,
		ImapHost: req.ImapHost,
		ImapPort: int(req.ImapPort),
		SmtpHost: req.SmtpHost,
		SmtpPort: int(req.SmtpPort),
		Password: req.Password,
	}

	created, err := e.emailSvc.AddAccount(ctx, account)
	if err != nil {
		return nil, rerrors.Wrap(err, "add email account")
	}

	return &pb.AddEmailAccount_Response{
		Account: toProto(created),
	}, nil
}
