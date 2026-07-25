package external_connections_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (e *ExternalConnectionsImpl) AddGenericConnection(
	ctx context.Context,
	req *pb.AddGenericConnection_Request,
) (*pb.AddGenericConnection_Response, error) {
	conn, err := e.svc.AddGenericConnection(ctx, req.Provider, req.Credentials)
	if err != nil {
		return nil, rerrors.Wrap(err, "add generic connection")
	}

	resp := &pb.AddGenericConnection_Response{
		Connection: ConnectionToProto(conn),
	}

	return resp, nil
}
