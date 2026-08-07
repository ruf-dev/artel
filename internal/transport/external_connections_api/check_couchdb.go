package external_connections_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (e *ExternalConnectionsImpl) CheckCouchDBConnection(
	ctx context.Context,
	req *pb.CheckCouchDBConnection_Request,
) (*pb.CheckCouchDBConnection_Response, error) {
	err := e.svc.CheckCouchDBConnection(ctx, req.Url, req.Username, req.Password)
	if err != nil {
		return nil, rerrors.Wrap(err, "check couchdb connection")
	}

	return &pb.CheckCouchDBConnection_Response{}, nil
}
