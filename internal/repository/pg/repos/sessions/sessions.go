package sessions

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/domain"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

type SessionsRepo struct {
	q *artel_q.Queries
}

func New(q *artel_q.Queries) *SessionsRepo {
	return &SessionsRepo{q: q}
}

func (r *SessionsRepo) Create(
	ctx context.Context,
	session domain.Session,
) (domain.Session, error) {
	params := artel_q.CreateSessionParams{
		UserID:           session.UserUuid,
		Token:            session.Token,
		ExpiresAt:        session.ExpiresAt,
		RefreshToken:     sql.NullString{String: session.RefreshToken, Valid: session.RefreshToken != ""},
		RefreshExpiresAt: sql.NullTime{Time: session.RefreshExpiresAt, Valid: !session.RefreshExpiresAt.IsZero()},
	}

	created, err := r.q.CreateSession(ctx, params)
	if err != nil {
		return domain.Session{}, rerrors.Wrap(err, "error creating session")
	}

	return domain.Session{
		Uuid:             created.ID,
		UserUuid:         created.UserID,
		Token:            created.Token,
		ExpiresAt:        created.ExpiresAt,
		RefreshToken:     created.RefreshToken.String,
		RefreshExpiresAt: created.RefreshExpiresAt.Time,
		CreatedAt:        created.CreatedAt,
	}, nil
}

func (r *SessionsRepo) GetByToken(ctx context.Context, token string) (domain.Session, error) {
	session, err := r.q.GetSessionByToken(ctx, token)
	if err != nil {
		return domain.Session{}, rerrors.Wrap(user_errors.Unauthenticated)
	}

	return domain.Session{
		Uuid:      session.ID,
		UserUuid:  session.UserID,
		Token:     session.Token,
		ExpiresAt: session.ExpiresAt,
		CreatedAt: session.CreatedAt,
	}, nil
}

func (r *SessionsRepo) GetByTokenWithUser(ctx context.Context, token string) (domain.Session, domain.User, error) {
	row, err := r.q.GetSessionWithUser(ctx, token)
	if err != nil {
		return domain.Session{}, domain.User{}, rerrors.Wrap(user_errors.Unauthenticated)
	}

	session := domain.Session{
		Uuid:      row.SessionID,
		UserUuid:  row.UserID,
		Token:     row.Token,
		ExpiresAt: row.SessionExpiresAt,
		CreatedAt: row.SessionCreatedAt,
	}

	user := domain.User{
		Uuid:         row.UserID,
		Email:        row.Email.String,
		Username:     row.Username,
		PhotoUrl:     row.PhotoUrl,
		PasswordHash: row.PasswordHash,
		CreatedAt:    row.UserCreatedAt,
		UpdatedAt:    row.UserUpdatedAt,
	}

	return session, user, nil
}

func (r *SessionsRepo) Delete(ctx context.Context, token string) error {
	err := r.q.DeleteSession(ctx, token)
	if err != nil {
		return rerrors.Wrap(err, "failed to delete session")
	}

	return nil
}

func (r *SessionsRepo) GetByUserID(ctx context.Context, userUuid uuid.UUID) ([]domain.Session, error) {
	rows, err := r.q.GetSessionsByUserID(ctx, userUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "error getting sessions by user id")
	}

	result := make([]domain.Session, len(rows))
	for i, row := range rows {
		result[i] = domain.Session{
			Uuid:      row.ID,
			UserUuid:  row.UserID,
			Token:     row.Token,
			ExpiresAt: row.ExpiresAt,
			CreatedAt: row.CreatedAt,
		}
	}

	return result, nil
}

func (r *SessionsRepo) RotateByRefreshToken(
	ctx context.Context,
	oldRefreshToken string,
	newSession domain.Session,
) (sql.Null[domain.Session], error) {
	params := artel_q.RotateSessionParams{
		RefreshToken:     sql.NullString{String: oldRefreshToken, Valid: oldRefreshToken != ""},
		Token:            newSession.Token,
		ExpiresAt:        newSession.ExpiresAt,
		RefreshToken_2:   sql.NullString{String: newSession.RefreshToken, Valid: newSession.RefreshToken != ""},
		RefreshExpiresAt: sql.NullTime{Time: newSession.RefreshExpiresAt, Valid: !newSession.RefreshExpiresAt.IsZero()},
	}

	rotated, err := r.q.RotateSession(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.Null[domain.Session]{}, nil
		}

		return sql.Null[domain.Session]{}, rerrors.Wrap(err, "error rotating session")
	}

	session := domain.Session{
		Uuid:             rotated.ID,
		UserUuid:         rotated.UserID,
		Token:            rotated.Token,
		ExpiresAt:        rotated.ExpiresAt,
		RefreshToken:     rotated.RefreshToken.String,
		RefreshExpiresAt: rotated.RefreshExpiresAt.Time,
		CreatedAt:        rotated.CreatedAt,
	}

	return sql.Null[domain.Session]{V: session, Valid: true}, nil
}
