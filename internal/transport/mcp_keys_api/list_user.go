package mcp_keys_api

import (
	"context"

	"go.redsock.ru/rerrors"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
)

func (m *McpKeysImpl) ListUserMcpKeys(ctx context.Context, _ *pb.ListUserMcpKeys_Request) (*pb.ListUserMcpKeys_Response, error) {
	keys, err := m.mcpSvc.ListUserKeys(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "list user mcp keys")
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
		if key.LastAccessedAt != nil {
			keyInfo.LastAccessedAt = key.LastAccessedAt.String()
		}
		keyInfos = append(keyInfos, keyInfo)
	}

	return &pb.ListUserMcpKeys_Response{Keys: keyInfos}, nil
}
