package external_connections_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (e *ExternalConnectionsImpl) ListMailServerSuggestions(
	ctx context.Context,
	req *pb.ListMailServerSuggestions_Request,
) (*pb.ListMailServerSuggestions_Response, error) {
	suggestions, err := e.svc.ListMailServerSuggestions(ctx, req.Domain)
	if err != nil {
		return nil, rerrors.Wrap(err, "list mail server suggestions")
	}

	pbSuggestions := make([]*pb.MailServerSuggestion, len(suggestions))
	for i, s := range suggestions {
		pbSuggestions[i] = &pb.MailServerSuggestion{
			Domain:   s.Domain,
			Smtp:     s.Smtp,
			SmtpPort: int32(s.SmtpPort),
			Imap:     s.Imap,
			ImapPort: int32(s.ImapPort),
		}
	}

	return &pb.ListMailServerSuggestions_Response{
		Suggestions: pbSuggestions,
	}, nil
}
