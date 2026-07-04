package external_connections_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (e *ExternalConnectionsImpl) GetGooglePickerToken(
	ctx context.Context,
	_ *pb.GooglePickerToken_Request,
) (*pb.GooglePickerToken_Response, error) {
	accessToken, err := e.svc.GetPickerToken(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "error getting picker token")
	}

	resp := &pb.GooglePickerToken_Response{
		AccessToken: accessToken,
	}

	return resp, nil
}
