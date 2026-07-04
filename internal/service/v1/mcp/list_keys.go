package mcp

import (
	"context"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/domain"
	"go.redsock.ru/rerrors"
)

func (s *McpServiceImpl) ListKeys(ctx context.Context, vaultID uuid.UUID) ([]domain.McpKey, error) {
	keys, err := s.mcpKeys.ListMcpKeysByVault(ctx, vaultID)
	if err != nil {
		return nil, rerrors.Wrap(err, "list mcp keys")
	}

	return keys, nil
}
