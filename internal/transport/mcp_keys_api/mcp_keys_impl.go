package mcp_keys_api

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"go.redsock.ru/rerrors"
	"google.golang.org/grpc"

	"github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service"
)

type McpKeysImpl struct {
	artel_api.UnimplementedMcpKeysAPIServer
	mcpSvc service.McpService
}

func NewMcpKeysImpl(mcpSvc service.McpService) *McpKeysImpl {
	return &McpKeysImpl{mcpSvc: mcpSvc}
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

	return "/api/vaults/", gwMux
}

func (m *McpKeysImpl) CreateMcpKey(ctx context.Context, req *artel_api.CreateMcpKey_Request) (*artel_api.CreateMcpKey_Response, error) {
	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse vault id")
	}

	rawToken, key, err := m.mcpSvc.CreateKey(ctx, vaultID, req.Name)
	if err != nil {
		return nil, rerrors.Wrap(err, "create mcp key")
	}

	keyInfo := &artel_api.McpKeyInfo{
		Id:         key.Uuid.String(),
		VaultId:    key.VaultUuid.String(),
		Name:       key.Name,
		KeyPreview: key.KeyPreview,
		CreatedAt:  key.CreatedAt.String(),
	}

	resp := &artel_api.CreateMcpKey_Response{
		Key:      keyInfo,
		RawToken: rawToken,
	}
	return resp, nil
}

func (m *McpKeysImpl) ListMcpKeys(ctx context.Context, req *artel_api.ListMcpKeys_Request) (*artel_api.ListMcpKeys_Response, error) {
	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse vault id")
	}

	keys, err := m.mcpSvc.ListKeys(ctx, vaultID)
	if err != nil {
		return nil, rerrors.Wrap(err, "list mcp keys")
	}

	keyInfos := make([]*artel_api.McpKeyInfo, 0, len(keys))
	for _, key := range keys {
		keyInfo := &artel_api.McpKeyInfo{
			Id:         key.Uuid.String(),
			VaultId:    key.VaultUuid.String(),
			Name:       key.Name,
			KeyPreview: key.KeyPreview,
			CreatedAt:  key.CreatedAt.String(),
		}
		keyInfos = append(keyInfos, keyInfo)
	}

	resp := &artel_api.ListMcpKeys_Response{
		Keys: keyInfos,
	}
	return resp, nil
}

func (m *McpKeysImpl) RevokeMcpKey(ctx context.Context, req *artel_api.RevokeMcpKey_Request) (*artel_api.RevokeMcpKey_Response, error) {
	keyID, err := uuid.Parse(req.KeyId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse key id")
	}

	err = m.mcpSvc.RevokeKey(ctx, keyID)
	if err != nil {
		return nil, rerrors.Wrap(err, "revoke mcp key")
	}

	resp := &artel_api.RevokeMcpKey_Response{}
	return resp, nil
}
