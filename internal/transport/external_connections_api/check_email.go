package external_connections_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (e *ExternalConnectionsImpl) CheckEmailConnection(
	ctx context.Context,
	req *pb.CheckEmailConnection_Request,
) (*pb.CheckEmailConnection_Response, error) {
	err := e.svc.CheckEmailConnection(
		ctx,
		req.Email,
		req.ImapHost,
		int(req.ImapPort),
		req.SmtpHost,
		int(req.SmtpPort),
		req.Password,
	)
	if err != nil {
		return nil, rerrors.Wrap(err, "check email connection")
	}

	return &pb.CheckEmailConnection_Response{}, nil
}
