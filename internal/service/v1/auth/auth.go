package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.redsock.ru/rerrors"
	"golang.org/x/crypto/bcrypt"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

type Service struct {
	usersRepo        repository.Users
	sessionsRepo     repository.Sessions
	jwksClient       keyfunc.Keyfunc
	telegramClientId string
}

func New(repo repository.Repo, telegramClientId string) (*Service, error) {
	jwks, err := keyfunc.NewDefaultCtx(context.Background(), []string{"https://oauth.telegram.org/.well-known/jwks.json"})
	if err != nil {
		return nil, rerrors.Wrap(err, "init telegram jwks client")
	}

	return &Service{
		usersRepo:        repo.Users(),
		sessionsRepo:     repo.Sessions(),
		jwksClient:       jwks,
		telegramClientId: telegramClientId,
	}, nil
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
		return uuid.UUID{}, user_errors.SessionExpired
	}

	return session.UserUuid, nil
}

func (s *Service) LoginViaTelegram(ctx context.Context, idToken string) (domain.Session, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(idToken, claims, s.jwksClient.Keyfunc,
		jwt.WithValidMethods([]string{"RS256", "ES256"}))
	if err != nil {
		return domain.Session{}, rerrors.Wrap(err, "parse telegram id_token")
	}

	if !token.Valid {
		return domain.Session{}, user_errors.InvalidTelegramToken
	}

	telegramId := claims.Subject

	user, err := s.usersRepo.UpsertByTelegramId(ctx, telegramId, telegramId)
	if err != nil {
		return domain.Session{}, rerrors.Wrap(err, "upsert telegram user")
	}

	sessionToken := generateToken()
	expiresAt := time.Now().Add(24 * time.Hour)

	session, err := s.sessionsRepo.Create(ctx, user.Uuid, sessionToken, expiresAt)
	if err != nil {
		return domain.Session{}, rerrors.Wrap(err, "create session")
	}

	return session, nil
}

func generateToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(b)
}
