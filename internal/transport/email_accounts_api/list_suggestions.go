package email_accounts_api

import (
	"context"

	"go.redsock.ru/rerrors"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/domain"
)

func (e *EmailAccountsImpl) ListMailServerSuggestions(ctx context.Context, req *pb.ListMailServerSuggestions_Request) (*pb.ListMailServerSuggestions_Response, error) {
	suggestions, err := e.emailSvc.ListMailServerSuggestions(ctx, req.Domain)
	if err != nil {
		return nil, rerrors.Wrap(err, "list mail server suggestions")
	}

	return &pb.ListMailServerSuggestions_Response{
		Suggestions: toProtoSuggestions(suggestions),
	}, nil
}

func toProtoSuggestions(suggestions []domain.MailServerSuggestion) []*pb.MailServerSuggestionInfo {
	result := make([]*pb.MailServerSuggestionInfo, len(suggestions))
	for i, s := range suggestions {
		result[i] = &pb.MailServerSuggestionInfo{
			Domain:   s.Domain,
			Smtp:     s.Smtp,
			SmtpPort: int32(s.SmtpPort),
			Imap:     s.Imap,
			ImapPort: int32(s.ImapPort),
		}
	}
	return result
}
