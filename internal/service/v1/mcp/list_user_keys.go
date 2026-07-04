package mcp

import (
	"context"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

func (s *ServiceImpl) ListUserKeys(ctx context.Context) ([]domain.McpKey, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return nil, user_errors.Unauthenticated
	}

	keys, err := s.mcpKeys.ListMcpKeysByUser(ctx, uc.UserUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "list user mcp keys")
	}

	return keys, nil
}
