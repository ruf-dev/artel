// Package adminsettings implements AdminSystemSettingsService — an administrator's view over
// the same single-row global instance configuration the first-run setup wizard
// (internal/service/v1/setupwizard) edits pre-setup, exposed post-setup for an admin to revisit.
package adminsettings

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

type Service struct {
	settingsRepo repository.SystemSettingsRepo
	vaultsRepo   repository.Vaults
}

func New(settingsRepo repository.SystemSettingsRepo, vaultsRepo repository.Vaults) *Service {
	return &Service{settingsRepo: settingsRepo, vaultsRepo: vaultsRepo}
}

func (s *Service) GetSettings(ctx context.Context) (domain.SystemSettings, domain.Vault, error) {
	settings, err := s.settingsRepo.Get(ctx)
	if err != nil {
		return domain.SystemSettings{}, domain.Vault{}, rerrors.Wrap(err, "get system settings")
	}

	if settings.DefaultDocsVaultID == nil {
		return settings, domain.Vault{}, nil
	}

	vault, err := s.vaultsRepo.GetByID(ctx, *settings.DefaultDocsVaultID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// The referenced vault was deleted out from under the FK (ON DELETE SET NULL should
			// normally prevent this from being reachable, but tolerate it rather than failing the
			// whole settings page).
			return settings, domain.Vault{}, nil
		}

		return domain.SystemSettings{}, domain.Vault{}, rerrors.Wrap(err, "get default docs vault")
	}

	return settings, vault, nil
}

func (s *Service) UpdateAuthMethods(ctx context.Context, passwordEnabled, telegramEnabled bool) error {
	if !passwordEnabled && !telegramEnabled {
		return user_errors.AtLeastOneAuthMethodRequired
	}

	err := s.settingsRepo.UpdateAuthMethods(ctx, passwordEnabled, telegramEnabled)
	if err != nil {
		return rerrors.Wrap(err, "update auth methods")
	}

	return nil
}

func (s *Service) UpdateRegistrationMode(ctx context.Context, mode domain.RegistrationMode) error {
	err := s.settingsRepo.UpdateRegistrationMode(ctx, mode)
	if err != nil {
		return rerrors.Wrap(err, "update registration mode")
	}

	return nil
}

// UpdateDefaultDocsVault sets the instance-wide default `/docs` source. When source is
// DocsSourceGithub, the vault-existence check is skipped entirely and the default vault is
// cleared (same clearing path as a nil vaultUuid) since the default now routes to the built-in
// GitHub quick-start guide instead. Otherwise (DocsSourceVault, including the unset/zero-value
// case), a non-nil vaultUuid must resolve to an existing vault — it is not required to already
// be published; PublicDocsService.GetDefaultVault re-checks IsPublic at read time.
func (s *Service) UpdateDefaultDocsVault(ctx context.Context, vaultUuid *uuid.UUID, source domain.DocsSource) error {
	if source == domain.DocsSourceGithub {
		err := s.settingsRepo.UpdateDefaultDocsVault(ctx, nil)
		if err != nil {
			return rerrors.Wrap(err, "update default docs vault")
		}

		err = s.settingsRepo.UpdateDefaultDocsSource(ctx, domain.DocsSourceGithub)
		if err != nil {
			return rerrors.Wrap(err, "update default docs source")
		}

		return nil
	}

	if vaultUuid != nil {
		_, err := s.vaultsRepo.GetByID(ctx, *vaultUuid)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return rerrors.Wrap(user_errors.NotFound, "vault not found")
			}

			return rerrors.Wrap(err, "get vault by id")
		}
	}

	err := s.settingsRepo.UpdateDefaultDocsVault(ctx, vaultUuid)
	if err != nil {
		return rerrors.Wrap(err, "update default docs vault")
	}

	err = s.settingsRepo.UpdateDefaultDocsSource(ctx, domain.DocsSourceVault)
	if err != nil {
		return rerrors.Wrap(err, "update default docs source")
	}

	return nil
}

func (s *Service) UpdateSystemPrompt(ctx context.Context, prompt string) error {
	repo, ok := s.settingsRepo.(interface {
		UpdateSystemPrompt(context.Context, string) error
	})
	if !ok {
		return rerrors.New("system prompt storage is unavailable")
	}
	err := repo.UpdateSystemPrompt(ctx, prompt)
	if err != nil {
		return rerrors.Wrap(err, "update system prompt")
	}
	return nil
}
