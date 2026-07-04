package userpermissions

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"go.redsock.ru/rerrors"
)

type Repo struct {
	q *artel_q.Queries
}

func New(q *artel_q.Queries) *Repo {
	return &Repo{q: q}
}

func (r *Repo) WithTx(tx *sql.Tx) repository.UserPermissionsRepo {
	return &Repo{q: r.q.WithTx(tx)}
}

func (r *Repo) Get(ctx context.Context, userUuid uuid.UUID) (domain.UserPermissions, error) {
	row, err := r.q.GetUserPermissions(ctx, userUuid)
	if err != nil {
		return domain.UserPermissions{}, rerrors.Wrap(err, "get user permissions")
	}

	return domain.UserPermissions{
		UserUuid:        row.UserID,
		IsAdministrator: row.IsAdministrator,
		HasEmails:       row.HasEmails,
		HasTaskTrackers: row.HasTaskTrackers,
		HasNotes:        row.HasNotes,
		HasSpreadsheets: row.HasSpreadsheets,
	}, nil
}

func (r *Repo) Upsert(
	ctx context.Context,
	userUuid uuid.UUID,
	isAdmin bool,
	hasEmails bool,
	hasTaskTrackers bool,
	hasNotes bool,
	hasSpreadsheets bool,
) (domain.UserPermissions, error) {
	params := artel_q.UpsertUserPermissionsParams{
		UserID:          userUuid,
		IsAdministrator: isAdmin,
		HasEmails:       hasEmails,
		HasTaskTrackers: hasTaskTrackers,
		HasNotes:        hasNotes,
		HasSpreadsheets: hasSpreadsheets,
	}

	row, err := r.q.UpsertUserPermissions(ctx, params)
	if err != nil {
		return domain.UserPermissions{}, rerrors.Wrap(err, "upsert user permissions")
	}

	return domain.UserPermissions{
		UserUuid:        row.UserID,
		IsAdministrator: row.IsAdministrator,
		HasEmails:       row.HasEmails,
		HasTaskTrackers: row.HasTaskTrackers,
		HasNotes:        row.HasNotes,
		HasSpreadsheets: row.HasSpreadsheets,
	}, nil
}

func (r *Repo) CreateDefault(ctx context.Context, userUuid uuid.UUID) error {
	err := r.q.CreateDefaultUserPermissions(ctx, userUuid)
	if err != nil {
		return rerrors.Wrap(err, "create default user permissions")
	}

	return nil
}
