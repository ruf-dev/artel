package external_connections_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (e *ExternalConnectionsImpl) CheckGitlabConnection(
	ctx context.Context,
	req *pb.CheckGitlabConnection_Request,
) (*pb.CheckGitlabConnection_Response, error) {
	username, err := e.svc.CheckGitlabConnection(ctx, req.PersonalAccessToken, req.InstanceUrl)
	if err != nil {
		return nil, rerrors.Wrap(err, "check gitlab connection")
	}

	resp := &pb.CheckGitlabConnection_Response{
		Username: username,
	}

	return resp, nil
}
