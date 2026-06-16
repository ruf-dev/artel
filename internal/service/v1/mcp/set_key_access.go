package mcp

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

func (s *McpServiceImpl) SetKeyAccess(ctx context.Context, keyID uuid.UUID, vaultID uuid.UUID) error {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return user_errors.Unauthenticated
	}

	_, err := s.vaultMembers.Get(ctx, vaultID, uc.UserUuid)
	if err != nil {
		return rerrors.Wrap(user_errors.Unauthenticated, "user is not a member of target vault")
	}

	err = s.mcpKeys.SetMcpKeyAccess(ctx, keyID, uc.UserUuid, vaultID)
	if err != nil {
		return rerrors.Wrap(err, "set mcp key access")
	}
	return nil
}
