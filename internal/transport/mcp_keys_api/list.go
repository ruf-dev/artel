package mcp_keys_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (m *McpKeysImpl) ListMcpKeys(ctx context.Context, req *pb.ListMcpKeys_Request) (*pb.ListMcpKeys_Response, error) {
	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse vault id")
	}

	keys, err := m.mcpSvc.ListKeys(ctx, vaultID)
	if err != nil {
		return nil, rerrors.Wrap(err, "list mcp keys")
	}

	keyInfos := make([]*pb.McpKeyInfo, 0, len(keys))

	for _, key := range keys {
		keyInfo := &pb.McpKeyInfo{
			Id:         key.Uuid.String(),
			VaultId:    key.VaultUuid.String(),
			Name:       key.Name,
			KeyPreview: key.KeyPreview,
			CreatedAt:  key.CreatedAt.String(),
		}
		keyInfos = append(keyInfos, keyInfo)
	}

	resp := &pb.ListMcpKeys_Response{
		Keys: keyInfos,
	}

	return resp, nil
}
