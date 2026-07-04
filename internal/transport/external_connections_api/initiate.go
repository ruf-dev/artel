package external_connections_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (e *ExternalConnectionsImpl) InitiateGoogleOAuth(
	ctx context.Context,
	req *pb.InitiateGoogleOAuth_Request,
) (*pb.InitiateGoogleOAuth_Response, error) {
	authURL, err := e.svc.InitiateGoogleOAuth(ctx, req.Origin)
	if err != nil {
		return nil, rerrors.Wrap(err, "error initiating google oauth")
	}

	resp := &pb.InitiateGoogleOAuth_Response{
		AuthUrl: authURL,
	}

	return resp, nil
}
