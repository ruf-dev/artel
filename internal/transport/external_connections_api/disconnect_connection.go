package external_connections_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (e *ExternalConnectionsImpl) DisconnectConnection(
	ctx context.Context,
	req *pb.DisconnectConnection_Request,
) (*pb.DisconnectConnection_Response, error) {
	err := e.svc.DisconnectConnection(ctx, req.Id)
	if err != nil {
		return nil, rerrors.Wrap(err, "error disconnecting connection")
	}

	return &pb.DisconnectConnection_Response{}, nil
}
