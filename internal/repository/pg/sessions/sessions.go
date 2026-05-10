package sessions

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/domain"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
)

type SessionsRepo struct {
	q *artel_q.Queries
}

func New(q *artel_q.Queries) *SessionsRepo {
	return &SessionsRepo{q: q}
}

func (r *SessionsRepo) Create(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) (domain.Session, error) {
	params := artel_q.CreateSessionParams{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
	}
	session, err := r.q.CreateSession(ctx, params)
	if err != nil {
		return domain.Session{}, rerrors.Wrap(err, "failed to create session")
	}

	return domain.Session{
		Uuid:      session.ID,
		UserUuid:  session.UserID,
		Token:     session.Token,
		ExpiresAt: session.ExpiresAt,
		CreatedAt: session.CreatedAt,
	}, nil
}

func (r *SessionsRepo) GetByToken(ctx context.Context, token string) (domain.Session, error) {
	session, err := r.q.GetSessionByToken(ctx, token)
	if err != nil {
		return domain.Session{}, rerrors.Wrap(err, "failed to get session by token")
	}

	return domain.Session{
		Uuid:      session.ID,
		UserUuid:  session.UserID,
		Token:     session.Token,
		ExpiresAt: session.ExpiresAt,
		CreatedAt: session.CreatedAt,
	}, nil
}

func (r *SessionsRepo) Delete(ctx context.Context, token string) error {
	err := r.q.DeleteSession(ctx, token)
	if err != nil {
		return rerrors.Wrap(err, "failed to delete session")
	}

	return nil
}
