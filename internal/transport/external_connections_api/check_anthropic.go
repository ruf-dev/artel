package external_connections_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (e *ExternalConnectionsImpl) CheckAnthropicConnection(
	ctx context.Context,
	req *pb.CheckAnthropicConnection_Request,
) (*pb.CheckAnthropicConnection_Response, error) {
	models, recommendedDefaultModel, err := e.svc.CheckAnthropicConnection(ctx, req.ApiKey, req.BaseUrl, req.DefaultModel)
	if err != nil {
		return nil, rerrors.Wrap(err, "check anthropic connection")
	}

	modelIds := make([]string, len(models))
	for i, m := range models {
		modelIds[i] = m.Id
	}

	resp := &pb.CheckAnthropicConnection_Response{
		AvailableModels:         modelIds,
		RecommendedDefaultModel: recommendedDefaultModel,
	}

	return resp, nil
}
