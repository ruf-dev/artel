package mcp_keys_api

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service"
	"google.golang.org/grpc"
)

type McpKeysImpl struct {
	artel_api.UnimplementedMcpKeysAPIServer
	mcpSvc service.McpService
	momSvc service.MomService
}

func NewMcpKeysImpl(mcpSvc service.McpService, momSvc service.MomService) *McpKeysImpl {
	return &McpKeysImpl{mcpSvc: mcpSvc, momSvc: momSvc}
}

func (m *McpKeysImpl) Register(srv grpc.ServiceRegistrar) {
	artel_api.RegisterMcpKeysAPIServer(srv, m)
}

func (m *McpKeysImpl) Gateway(ctx context.Context, endpoint string, opts ...grpc.DialOption) (string, http.Handler) {
	gwMux := runtime.NewServeMux()

	err := artel_api.RegisterMcpKeysAPIHandlerFromEndpoint(ctx, gwMux, endpoint, opts)
	if err != nil {
		log.Error().Err(err).Msg("error registering mcp keys grpc-gateway handler")
	}

	return "/api/mcp/", gwMux
}
