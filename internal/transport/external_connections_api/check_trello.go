package external_connections_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (e *ExternalConnectionsImpl) CheckTrelloConnection(
	ctx context.Context,
	req *pb.CheckTrelloConnection_Request,
) (*pb.CheckTrelloConnection_Response, error) {
	fullName, err := e.svc.CheckTrelloConnection(ctx, req.ApiKey, req.ApiToken)
	if err != nil {
		return nil, rerrors.Wrap(err, "check trello connection")
	}

	resp := &pb.CheckTrelloConnection_Response{
		FullName: fullName,
	}

	return resp, nil
}
