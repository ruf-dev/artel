package external_connections_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (e *ExternalConnectionsImpl) AddOpenAIConnection(
	ctx context.Context,
	req *pb.AddOpenAIConnection_Request,
) (*pb.AddOpenAIConnection_Response, error) {
	conn, err := e.svc.AddOpenAIConnection(ctx, req.ApiKey, req.BaseUrl, req.DefaultModel)
	if err != nil {
		return nil, rerrors.Wrap(err, "add openai connection")
	}

	resp := &pb.AddOpenAIConnection_Response{
		Connection: ConnectionToProto(conn),
	}

	return resp, nil
}
