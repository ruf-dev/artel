package mcp

import (
	"context"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

func (s *ServiceImpl) RemoveConnector(ctx context.Context, keyID uuid.UUID, mcpName string) error {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return user_errors.Unauthenticated
	}

	key, err := s.mcpKeys.GetMcpKeyByID(ctx, keyID)
	if err != nil {
		return rerrors.Wrap(err, "get mcp key")
	}

	if key.UserUuid != uc.UserUuid {
		return rerrors.Wrap(user_errors.NotFound, "mcp key not found")
	}

	err = s.mcpConnectors.Delete(ctx, keyID, mcpName)
	if err != nil {
		return rerrors.Wrap(err, "delete mcp connector")
	}

	return nil
}
