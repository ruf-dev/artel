package external_connections_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (e *ExternalConnectionsImpl) AddEmailConnection(
	ctx context.Context,
	req *pb.AddEmailConnection_Request,
) (*pb.AddEmailConnection_Response, error) {
	conn, err := e.svc.AddEmailConnection(
		ctx,
		req.Email,
		req.ImapHost,
		int(req.ImapPort),
		req.SmtpHost,
		int(req.SmtpPort),
		req.Password,
	)
	if err != nil {
		return nil, rerrors.Wrap(err, "add email connection")
	}

	return &pb.AddEmailConnection_Response{
		Connection: ConnectionToProto(conn),
	}, nil
}
