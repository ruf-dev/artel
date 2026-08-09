package external_connections_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (e *ExternalConnectionsImpl) CheckPostgresConnection(
	ctx context.Context,
	req *pb.CheckPostgresConnection_Request,
) (*pb.CheckPostgresConnection_Response, error) {
	err := e.svc.CheckPostgresConnection(
		ctx, req.Host, int(req.Port), req.Database, req.Username, req.Password, req.SslMode,
	)
	if err != nil {
		return nil, rerrors.Wrap(err, "check postgres connection")
	}

	return &pb.CheckPostgresConnection_Response{}, nil
}
