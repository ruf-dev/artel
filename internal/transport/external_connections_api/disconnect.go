package external_connections_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (e *ExternalConnectionsImpl) DisconnectProvider(
	ctx context.Context,
	req *pb.DisconnectProvider_Request,
) (*pb.DisconnectProvider_Response, error) {
	err := e.svc.DisconnectProvider(ctx, req.Provider)
	if err != nil {
		return nil, rerrors.Wrap(err, "error disconnecting provider")
	}

	return &pb.DisconnectProvider_Response{}, nil
}
