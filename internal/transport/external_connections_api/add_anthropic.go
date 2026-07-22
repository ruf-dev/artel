package external_connections_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (e *ExternalConnectionsImpl) AddAnthropicConnection(
	ctx context.Context,
	req *pb.AddAnthropicConnection_Request,
) (*pb.AddAnthropicConnection_Response, error) {
	conn, err := e.svc.AddAnthropicConnection(ctx, req.ApiKey, req.BaseUrl, req.DefaultModel)
	if err != nil {
		return nil, rerrors.Wrap(err, "add anthropic connection")
	}

	resp := &pb.AddAnthropicConnection_Response{
		Connection: ConnectionToProto(conn),
	}

	return resp, nil
}
