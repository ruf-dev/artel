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
	return domain.UserSettings{UserUuid: row.UserID, UserPrompt: row.UserPrompt}, nil
}

var _ repository.UserSettingsRepo = (*Repo)(nil)
