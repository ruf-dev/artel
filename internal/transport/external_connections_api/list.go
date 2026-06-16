package external_connections_api

import (
	"context"

	"go.redsock.ru/rerrors"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
)

func (e *ExternalConnectionsImpl) ListConnections(ctx context.Context, _ *pb.ListConnections_Request) (*pb.ListConnections_Response, error) {
	metas, err := e.svc.ListConnections(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing connections")
	}

	connections := make([]*pb.ExternalConnectionInfo, len(metas))
	for i, m := range metas {
		connections[i] = ConnectionToProto(m)
	}

	resp := &pb.ListConnections_Response{
		Connections: connections,
	}
	return resp, nil
}
