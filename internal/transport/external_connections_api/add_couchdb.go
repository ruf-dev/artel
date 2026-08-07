package external_connections_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (e *ExternalConnectionsImpl) AddCouchDBConnection(
	ctx context.Context,
	req *pb.AddCouchDBConnection_Request,
) (*pb.AddCouchDBConnection_Response, error) {
	conn, err := e.svc.AddCouchDBConnection(ctx, req.Url, req.Username, req.Password)
	if err != nil {
		return nil, rerrors.Wrap(err, "add couchdb connection")
	}

	resp := &pb.AddCouchDBConnection_Response{
		Connection: ConnectionToProto(conn),
	}

	return resp, nil
}
