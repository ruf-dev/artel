package mcp

import (
	"context"

	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

// DeleteCommunityConnector deletes a community connector by name. Returns
// user_errors.NotFound both when name doesn't exist and when the caller isn't its owner —
// deliberately not a distinct "forbidden" error, so a non-owner can't learn who owns a name.
func (s *ServiceImpl) DeleteCommunityConnector(ctx context.Context, name string) error {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return user_errors.Unauthenticated
	}

	err := s.authSvc.CheckIsAdmin(ctx, uc.UserUuid)
	if err != nil {
		return rerrors.Wrap(err, "check caller is admin")
	}

	def, err := s.mcpDefinitions.Get(ctx, name)
	if err != nil {
		return rerrors.Wrap(err, "get mcp definition")
	}

	if !def.Valid {
		return user_errors.NotFound
	}

	notOwner := def.V.OwnerUserUuid == nil || *def.V.OwnerUserUuid != uc.UserUuid
	if notOwner {
		return user_errors.NotFound
	}

	err = s.mcpDefinitions.Delete(ctx, name)
	if err != nil {
		return rerrors.Wrap(err, "delete mcp definition")
	}

	return nil
}
