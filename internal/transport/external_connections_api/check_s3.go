package external_connections_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (e *ExternalConnectionsImpl) CheckS3Connection(
	ctx context.Context,
	req *pb.CheckS3Connection_Request,
) (*pb.CheckS3Connection_Response, error) {
	err := e.svc.CheckS3Connection(
		ctx,
		req.Endpoint,
		req.Region,
		req.AccessKey,
		req.SecretKey,
		req.UseSsl,
		req.PathStyle,
	)
	if err != nil {
		return nil, rerrors.Wrap(err, "check s3 connection")
	}

	return &pb.CheckS3Connection_Response{}, nil
}
