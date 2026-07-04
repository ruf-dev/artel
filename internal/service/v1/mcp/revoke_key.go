package mcp

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"
)

func (s *McpServiceImpl) RevokeKey(ctx context.Context, keyID uuid.UUID) error {
	err := s.mcpKeys.RevokeMcpKey(ctx, keyID)
	if err != nil {
		return rerrors.Wrap(err, "revoke mcp key")
	}

	return nil
}
