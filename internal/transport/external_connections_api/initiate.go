package external_connections_api

import (
	"context"

	"go.redsock.ru/rerrors"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
)

func (e *ExternalConnectionsImpl) InitiateGoogleOAuth(ctx context.Context, _ *pb.InitiateGoogleOAuth_Request) (*pb.InitiateGoogleOAuth_Response, error) {
	authURL, err := e.svc.InitiateGoogleOAuth(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "error initiating google oauth")
	}

	resp := &pb.InitiateGoogleOAuth_Response{
		AuthUrl: authURL,
	}
	return resp, nil
}
