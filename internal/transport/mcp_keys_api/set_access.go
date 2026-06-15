package mcp_keys_api

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
)

func (m *McpKeysImpl) SetMcpKeyAccess(ctx context.Context, req *pb.SetMcpKeyAccess_Request) (*pb.SetMcpKeyAccess_Response, error) {
	keyID, err := uuid.Parse(req.KeyId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse key id")
	}

	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse vault id")
	}

	err = m.mcpSvc.SetKeyAccess(ctx, keyID, vaultID)
	if err != nil {
		return nil, rerrors.Wrap(err, "set mcp key access")
	}

	return &pb.SetMcpKeyAccess_Response{}, nil
}
