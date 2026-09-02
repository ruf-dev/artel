package usersettings

import (
	"context"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"go.redsock.ru/rerrors"
)

type Repo struct{ q *artel_q.Queries }

func New(q *artel_q.Queries) *Repo { return &Repo{q: q} }

func (r *Repo) Get(ctx context.Context, userID uuid.UUID) (domain.UserSettings, error) {
	row, err := r.q.GetUserSettings(ctx, userID)
	if err != nil {
		return domain.UserSettings{}, rerrors.Wrap(err, "error getting user settings")
	}

	settings := domain.UserSettings{
		UserUuid:              row.UserID,
		UserPrompt:            row.UserPrompt,
		LikedOpenrouterModels: row.LikedOpenrouterModels,
		LastUsedModel:         row.LastUsedModel,
	}

	return settings, nil
}

func (r *Repo) UpsertLikedModels(ctx context.Context, userID uuid.UUID, modelIds []string) error {
	params := artel_q.UpsertLikedModelsParams{
		UserID:                userID,
		LikedOpenrouterModels: modelIds,
	}

	err := r.q.UpsertLikedModels(ctx, params)
	if err != nil {
		return rerrors.Wrap(err, "error upserting liked models")
	}

	return nil
}

func (r *Repo) UpsertLastUsedModel(ctx context.Context, userID uuid.UUID, model string) error {
	params := artel_q.UpsertLastUsedModelParams{
		UserID:        userID,
		LastUsedModel: model,
	}

	err := r.q.UpsertLastUsedModel(ctx, params)
	if err != nil {
		return rerrors.Wrap(err, "error upserting last used model")
	}

	return nil
}

var _ repository.UserSettingsRepo = (*Repo)(nil)
