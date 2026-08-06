package external_connections_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (e *ExternalConnectionsImpl) AddTelegramConnection(
	ctx context.Context,
	req *pb.AddTelegramConnection_Request,
) (*pb.AddTelegramConnection_Response, error) {
	conn, err := e.svc.AddTelegramConnection(ctx, req.BotToken)
	if err != nil {
		return nil, rerrors.Wrap(err, "add telegram connection")
	}

	resp := &pb.AddTelegramConnection_Response{
		Connection: ConnectionToProto(conn),
	}

	return resp, nil
}
