package auth

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

type fakeSessionsRepo struct {
	stored domain.Session
}

func (f *fakeSessionsRepo) Create(_ context.Context, session domain.Session) (domain.Session, error) {
	return session, nil
}

func (f *fakeSessionsRepo) GetByToken(_ context.Context, _ string) (domain.Session, error) {
	return domain.Session{}, nil
}

func (f *fakeSessionsRepo) GetByTokenWithUser(_ context.Context, _ string) (domain.Session, domain.User, error) {
	return domain.Session{}, domain.User{}, nil
}

func (f *fakeSessionsRepo) Delete(_ context.Context, _ string) error {
	return nil
}

func (f *fakeSessionsRepo) GetByUserID(_ context.Context, _ uuid.UUID) ([]domain.Session, error) {
	return nil, nil
}

func (f *fakeSessionsRepo) RotateByRefreshToken(
	_ context.Context,
	oldRefreshToken string,
	newSession domain.Session,
) (sql.Null[domain.Session], error) {
	if oldRefreshToken == "" || oldRefreshToken != f.stored.RefreshToken {
		return sql.Null[domain.Session]{}, nil
	}

	if time.Now().After(f.stored.RefreshExpiresAt) {
		return sql.Null[domain.Session]{}, nil
	}

	newSession.Uuid = f.stored.Uuid
	newSession.UserUuid = f.stored.UserUuid
	f.stored = newSession

	return sql.Null[domain.Session]{V: newSession, Valid: true}, nil
}

func TestServiceRefresh(t *testing.T) {
	validSession := domain.Session{
		Uuid:             uuid.New(),
		UserUuid:         uuid.New(),
		Token:            "old-access-token",
		ExpiresAt:        time.Now().Add(-time.Minute),
		RefreshToken:     "valid-refresh-token",
		RefreshExpiresAt: time.Now().Add(time.Hour),
	}

	t.Run("valid refresh token rotates and returns a new pair", func(t *testing.T) {
		repo := &fakeSessionsRepo{stored: validSession}
		svc := &Service{sessionsRepo: repo}

		got, err := svc.Refresh(context.Background(), validSession.RefreshToken)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got.Token == validSession.Token {
			t.Errorf("expected a newly generated access token, got the same one back")
		}

		if got.RefreshToken == validSession.RefreshToken {
			t.Errorf("expected a newly generated refresh token, got the same one back")
		}

		if !got.ExpiresAt.After(time.Now()) {
			t.Errorf("expected new access token to not be expired, got expiresAt=%v", got.ExpiresAt)
		}

		if got.UserUuid != validSession.UserUuid {
			t.Errorf("expected rotated session to keep the same user, got %v want %v", got.UserUuid, validSession.UserUuid)
		}
	})

	t.Run("expired refresh token is rejected", func(t *testing.T) {
		expiredSession := domain.Session{
			RefreshToken:     "expired-refresh-token",
			RefreshExpiresAt: time.Now().Add(-time.Hour),
		}
		repo := &fakeSessionsRepo{stored: expiredSession}
		svc := &Service{sessionsRepo: repo}

		_, err := svc.Refresh(context.Background(), expiredSession.RefreshToken)
		if !rerrors.Is(err, user_errors.InvalidRefreshToken) {
			t.Fatalf("expected InvalidRefreshToken, got %v", err)
		}
	})

	t.Run("unknown or already-rotated refresh token is rejected", func(t *testing.T) {
		repo := &fakeSessionsRepo{stored: validSession}
		svc := &Service{sessionsRepo: repo}

		_, err := svc.Refresh(context.Background(), "does-not-exist")
		if !rerrors.Is(err, user_errors.InvalidRefreshToken) {
			t.Fatalf("expected InvalidRefreshToken, got %v", err)
		}
	})
}
