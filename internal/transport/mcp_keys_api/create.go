package mcp_keys_api

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
)

func (m *McpKeysImpl) CreateMcpKey(ctx context.Context, req *pb.CreateMcpKey_Request) (*pb.CreateMcpKey_Response, error) {
	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse vault id")
	}

	rawToken, key, err := m.mcpSvc.CreateKey(ctx, vaultID, req.Name)
	if err != nil {
		return nil, rerrors.Wrap(err, "create mcp key")
	}

	keyInfo := &pb.McpKeyInfo{
		Id:         key.Uuid.String(),
		VaultId:    key.VaultUuid.String(),
		Name:       key.Name,
		KeyPreview: key.KeyPreview,
		CreatedAt:  key.CreatedAt.String(),
	}

	resp := &pb.CreateMcpKey_Response{
		Key:      keyInfo,
		RawToken: rawToken,
	}
	return resp, nil
}
