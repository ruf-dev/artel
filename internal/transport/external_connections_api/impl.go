package external_connections_api

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service"
)

type ExternalConnectionsImpl struct {
	pb.UnimplementedExternalConnectionsAPIServer
	svc service.ExternalConnectionService
}

func New(svc service.ExternalConnectionService) *ExternalConnectionsImpl {
	return &ExternalConnectionsImpl{svc: svc}
}

func (e *ExternalConnectionsImpl) Register(srv grpc.ServiceRegistrar) {
	pb.RegisterExternalConnectionsAPIServer(srv, e)
}

func (e *ExternalConnectionsImpl) Gateway(ctx context.Context, endpoint string, opts ...grpc.DialOption) (string, http.Handler) {
	gwMux := runtime.NewServeMux()

	err := pb.RegisterExternalConnectionsAPIHandlerFromEndpoint(ctx, gwMux, endpoint, opts)
	if err != nil {
		log.Error().Err(err).Msg("error registering external connections grpc-gateway handler")
	}

	return "/api/external-connections/", gwMux
}
