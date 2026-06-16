package mcp

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

func (s *McpServiceImpl) ListConnectors(ctx context.Context, keyID uuid.UUID) ([]domain.McpConnector, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return nil, user_errors.Unauthenticated
	}

	key, err := s.mcpKeys.GetMcpKeyByID(ctx, keyID)
	if err != nil {
		return nil, rerrors.Wrap(err, "get mcp key")
	}
	if key.UserUuid != uc.UserUuid {
		return nil, rerrors.Wrap(user_errors.NotFound, "mcp key not found")
	}

	connectors, err := s.mcpConnectors.ListByKey(ctx, keyID)
	if err != nil {
		return nil, rerrors.Wrap(err, "list mcp connectors")
	}
	return connectors, nil
}
