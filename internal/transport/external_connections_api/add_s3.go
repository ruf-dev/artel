package external_connections_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (e *ExternalConnectionsImpl) AddS3Connection(
	ctx context.Context,
	req *pb.AddS3Connection_Request,
) (*pb.AddS3Connection_Response, error) {
	conn, err := e.svc.AddS3Connection(
		ctx,
		req.Endpoint,
		req.Region,
		req.AccessKey,
		req.SecretKey,
		req.UseSsl,
		req.PathStyle,
	)
	if err != nil {
		return nil, rerrors.Wrap(err, "add s3 connection")
	}

	resp := &pb.AddS3Connection_Response{
		Connection: ConnectionToProto(conn),
	}

	return resp, nil
}
