package external_connections_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (e *ExternalConnectionsImpl) AddTrelloConnection(
	ctx context.Context,
	req *pb.AddTrelloConnection_Request,
) (*pb.AddTrelloConnection_Response, error) {
	conn, err := e.svc.AddTrelloConnection(ctx, req.ApiKey, req.ApiToken)
	if err != nil {
		return nil, rerrors.Wrap(err, "add trello connection")
	}

	resp := &pb.AddTrelloConnection_Response{
		Connection: ConnectionToProto(conn),
	}

	return resp, nil
}
