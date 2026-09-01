package external_connections_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (e *ExternalConnectionsImpl) EnsureArtelTelegramConnection(
	ctx context.Context,
	_ *pb.EnsureArtelTelegramConnection_Request,
) (*pb.EnsureArtelTelegramConnection_Response, error) {
	connID, err := e.svc.EnsureArtelTelegramConnection(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "ensure artel telegram connection")
	}

	resp := &pb.EnsureArtelTelegramConnection_Response{
		ExternalConnectionId: connID.String(),
	}

	return resp, nil
}
