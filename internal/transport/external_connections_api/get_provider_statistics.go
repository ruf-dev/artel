package external_connections_api

import (
	"context"

	"go.redsock.ru/rerrors"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

// GetProviderStatistics dispatches by req.Provider to the provider-specific statistics fetch —
// currently only OpenRouter implements this; every other provider returns
// user_errors.ProviderStatisticsUnsupported until it gets its own oneof case (see
// GetProviderStatistics.Response's oneof in api/grpc/external_connections.proto).
func (e *ExternalConnectionsImpl) GetProviderStatistics(
	ctx context.Context,
	req *pb.GetProviderStatistics_Request,
) (*pb.GetProviderStatistics_Response, error) {
	switch req.Provider {
	case pb.ExternalProvider_EXTERNAL_PROVIDER_OPENROUTER:
		stats, err := e.svc.GetOpenRouterStatistics(ctx)
		if err != nil {
			return nil, rerrors.Wrap(err, "get openrouter statistics")
		}

		resp := &pb.GetProviderStatistics_Response{
			Statistics: &pb.GetProviderStatistics_Response_Openrouter{
				Openrouter: toOpenRouterStatisticsProto(stats),
			},
		}

		return resp, nil
	default:
		return nil, user_errors.ProviderStatisticsUnsupported
	}
}
