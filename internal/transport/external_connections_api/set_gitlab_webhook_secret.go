package external_connections_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (e *ExternalConnectionsImpl) SetGitlabWebhookSecret(ctx context.Context, req *pb.SetGitlabWebhookSecret_Request) (*pb.SetGitlabWebhookSecret_Response, error) {
	conn, err := e.svc.SetGitlabWebhookSecret(ctx, req.WebhookSecret)
	if err != nil {
		return nil, rerrors.Wrap(err, "set gitlab webhook secret")
	}

	resp := &pb.SetGitlabWebhookSecret_Response{
		Connection: ConnectionToProto(conn),
	}

	return resp, nil
}
