package external_connections_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (e *ExternalConnectionsImpl) CheckOpenAIConnection(
	ctx context.Context,
	req *pb.CheckOpenAIConnection_Request,
) (*pb.CheckOpenAIConnection_Response, error) {
	models, recommendedDefaultModel, err := e.svc.CheckOpenAIConnection(ctx, req.ApiKey, req.BaseUrl, req.DefaultModel)
	if err != nil {
		return nil, rerrors.Wrap(err, "check openai connection")
	}

	modelIds := make([]string, len(models))
	for i, m := range models {
		modelIds[i] = m.Id
	}

	resp := &pb.CheckOpenAIConnection_Response{
		AvailableModels:         modelIds,
		RecommendedDefaultModel: recommendedDefaultModel,
	}

	return resp, nil
}
