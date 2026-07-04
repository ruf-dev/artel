package external_connections_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (e *ExternalConnectionsImpl) AddGitlabConnection(ctx context.Context, req *pb.AddGitlabConnection_Request) (*pb.AddGitlabConnection_Response, error) {
	conn, err := e.svc.AddGitlabConnection(ctx, req.PersonalAccessToken, req.InstanceUrl)
	if err != nil {
		return nil, rerrors.Wrap(err, "add gitlab connection")
	}

	resp := &pb.AddGitlabConnection_Response{
		Connection: ConnectionToProto(conn),
	}

	return resp, nil
}
