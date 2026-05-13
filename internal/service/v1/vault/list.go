package vault

import (
	"context"

	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/livesync"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

func (s *Service) ListVaults(ctx context.Context) ([]domain.Vault, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return nil, rerrors.Wrap(user_errors.Unauthenticated)
	}

	vaults, err := s.vaultsRepo.ListByMembership(ctx, uc.UserUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "list vaults by membership")
	}

	for i, vault := range vaults {
		//TODO remake onto one sql call
		instance, err := s.couchInstancesRepo.Get(ctx, vault.CouchInstanceUuid)
		if err != nil {
			return nil, rerrors.Wrap(err, "get couch instance")
		}

		dbUser, err := s.couchAccountsRepo.GetByUserAndInstance(ctx, uc.UserUuid, instance.Uuid)
		if err != nil {
			return nil, rerrors.Wrap(err, "get user and instance")
		}

		livesyncConfig := livesync.DefaultConfig(instance.Url,
			dbUser.CouchUsername, dbUser.CouchPassword, vault.CouchDBName, "1")

		livesyncConfig.Encrypt = false
		livesyncConfig.UsePathObfuscation = false

		vaults[i].CouchDBURL, err = livesync.GenerateSetupURI(livesyncConfig, "1")
		if err != nil {
			return nil, rerrors.Wrap(err, "generate setup uri")
		}
	}

	return vaults, nil
}
