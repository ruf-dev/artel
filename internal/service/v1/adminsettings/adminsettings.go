// Package adminsettings implements AdminSystemSettingsService — an administrator's view over
// the same single-row global instance configuration the first-run setup wizard
// (internal/service/v1/setupwizard) edits pre-setup, exposed post-setup for an admin to revisit.
package adminsettings

import (
	"context"

	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

type Service struct {
	settingsRepo repository.SystemSettingsRepo
}

func New(settingsRepo repository.SystemSettingsRepo) *Service {
	return &Service{settingsRepo: settingsRepo}
}

func (s *Service) GetSettings(ctx context.Context) (domain.SystemSettings, error) {
	settings, err := s.settingsRepo.Get(ctx)
	if err != nil {
		return domain.SystemSettings{}, rerrors.Wrap(err, "get system settings")
	}

	return settings, nil
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
