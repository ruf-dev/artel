package external_connections_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (e *ExternalConnectionsImpl) GenerateGitlabWebhookSecret(
	ctx context.Context,
	_ *pb.GenerateGitlabWebhookSecret_Request,
) (*pb.GenerateGitlabWebhookSecret_Response, error) {
	conn, webhookSecret, err := e.svc.GenerateGitlabWebhookSecret(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "generate gitlab webhook secret")
	}

	resp := &pb.GenerateGitlabWebhookSecret_Response{
		Connection:    ConnectionToProto(conn),
		WebhookSecret: webhookSecret,
	}

	return resp, nil
}
