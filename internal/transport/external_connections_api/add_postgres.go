package external_connections_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (e *ExternalConnectionsImpl) AddPostgresConnection(
	ctx context.Context,
	req *pb.AddPostgresConnection_Request,
) (*pb.AddPostgresConnection_Response, error) {
	conn, err := e.svc.AddPostgresConnection(
		ctx, req.Host, int(req.Port), req.Database, req.Username, req.Password, req.SslMode,
	)
	if err != nil {
		return nil, rerrors.Wrap(err, "add postgres connection")
	}

	resp := &pb.AddPostgresConnection_Response{
		Connection: ConnectionToProto(conn),
	}

	return resp, nil
}
