package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"
	"golang.org/x/crypto/bcrypt"

	"github.com/ruf-dev/artel/internal/domain"
	repository "github.com/ruf-dev/artel/internal/repository"
)

type Service struct {
	usersRepo    repository.Users
	sessionsRepo repository.Sessions
}

func New(usersRepo repository.Users, sessionsRepo repository.Sessions) *Service {
	return &Service{
		usersRepo:    usersRepo,
		sessionsRepo: sessionsRepo,
	}
}

func (s *Service) Register(ctx context.Context, email, password string) (domain.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, rerrors.Wrap(err, "generate password hash")
	}

	user, err := s.usersRepo.Create(ctx, email, string(hash))
	if err != nil {
		return domain.User{}, rerrors.Wrap(err, "create user")
	}

	return user, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (domain.Session, error) {
	user, err := s.usersRepo.GetByEmail(ctx, email)
	if err != nil {
		return domain.Session{}, rerrors.Wrap(err, "get user by email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return domain.Session{}, rerrors.Wrap(err, "password mismatch")
	}

	token := generateToken()

	expiresAt := time.Now().Add(24 * time.Hour)

	session, err := s.sessionsRepo.Create(ctx, user.Uuid, token, expiresAt)
	if err != nil {
		return domain.Session{}, rerrors.Wrap(err, "create session")
	}

	return session, nil
}

func (s *Service) Logout(ctx context.Context, token string) error {
	err := s.sessionsRepo.Delete(ctx, token)
	if err != nil {
		return rerrors.Wrap(err, "delete session")
	}

	return nil
}

func (s *Service) ValidateToken(ctx context.Context, token string) (uuid.UUID, error) {
	session, err := s.sessionsRepo.GetByToken(ctx, token)
	if err != nil {
		return uuid.UUID{}, rerrors.Wrap(err, "get session by token")
	}

	if time.Now().After(session.ExpiresAt) {
		return uuid.UUID{}, rerrors.New("session expired")
	}

	return session.UserUuid, nil
}

func generateToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(b)
}
