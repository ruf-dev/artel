package external_connections_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (e *ExternalConnectionsImpl) CheckTelegramConnection(
	ctx context.Context,
	req *pb.CheckTelegramConnection_Request,
) (*pb.CheckTelegramConnection_Response, error) {
	botUsername, err := e.svc.CheckTelegramConnection(ctx, req.BotToken)
	if err != nil {
		return nil, rerrors.Wrap(err, "check telegram connection")
	}

	resp := &pb.CheckTelegramConnection_Response{
		BotUsername: botUsername,
	}

	return resp, nil
}
