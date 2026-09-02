// Package usersettings implements service.UserSettingsService: per-user preference storage
// synced across devices (liked OpenRouter models, last used model), replacing what used to be
// pure localStorage on the frontend. Every method resolves the caller from ctx and operates only
// on that user's own row — there is no id parameter from the client.
package usersettings

import (
	"context"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

type Service struct {
	userSettings repository.UserSettingsRepo
}

func New(userSettings repository.UserSettingsRepo) *Service {
	return &Service{userSettings: userSettings}
}

func (s *Service) GetUserSettings(ctx context.Context) (domain.UserSettings, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return domain.UserSettings{}, rerrors.Wrap(user_errors.Unauthenticated)
	}

	settings, err := s.userSettings.Get(ctx, uc.UserUuid)
	if err != nil {
		return domain.UserSettings{}, rerrors.Wrap(err, "error getting user settings")
	}

	return settings, nil
}

func (s *Service) SetLikedModels(ctx context.Context, modelIds []string) error {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return rerrors.Wrap(user_errors.Unauthenticated)
	}

	err := s.userSettings.UpsertLikedModels(ctx, uc.UserUuid, modelIds)
	if err != nil {
		return rerrors.Wrap(err, "error setting liked models")
	}

	return nil
}

func (s *Service) SetLastUsedModel(ctx context.Context, model string) error {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return rerrors.Wrap(user_errors.Unauthenticated)
	}

	err := s.userSettings.UpsertLastUsedModel(ctx, uc.UserUuid, model)
	if err != nil {
		return rerrors.Wrap(err, "error setting last used model")
	}

	return nil
}
