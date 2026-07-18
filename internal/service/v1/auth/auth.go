package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"strconv"
	"time"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"
	"golang.org/x/crypto/bcrypt"

	"github.com/ruf-dev/artel/internal/client/telegram"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/repository/pg/tx_manager"
	"github.com/ruf-dev/artel/internal/service"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

const (
	accessTokenTTL  = time.Hour
	refreshTokenTTL = 30 * 24 * time.Hour
)

type Service struct {
	usersRepo       repository.Users
	sessionsRepo    repository.Sessions
	permissionsRepo repository.UserPermissionsRepo
	subsRepo        repository.Subscriptions
	subscriptions   service.SubscriptionService
	txManager       tx_manager.TxManager
	tgParser        telegram.TokenParser
}

func New(repo repository.Repo, telegramClientId string, subscriptions service.SubscriptionService) *Service {
	tgParser := telegram.NewTokenParser(
		"https://oauth.telegram.org/.well-known/jwks.json",
		"https://oauth.telegram.org",
		telegramClientId,
	)

	return &Service{
		usersRepo:       repo.Users(),
		sessionsRepo:    repo.Sessions(),
		permissionsRepo: repo.UserPermissions(),
		subsRepo:        repo.Subscriptions(),
		subscriptions:   subscriptions,
		txManager:       repo.TxManager(),
		tgParser:        tgParser,
	}
}

func (s *Service) Register(ctx context.Context, email, password string) (domain.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, rerrors.Wrap(err, "generate password hash")
	}

	var user domain.User

	err = s.txManager.Execute(func(tx *sql.Tx) error {
		user, err = s.usersRepo.WithTx(tx).Create(ctx, email, string(hash))
		if err != nil {
			return rerrors.Wrap(err, "create user")
		}

		err = s.permissionsRepo.WithTx(tx).CreateDefault(ctx, user.Uuid)
		if err != nil {
			return rerrors.Wrap(err, "create default permissions")
		}

		err = s.subsRepo.WithTx(tx).CreateDefault(ctx, user.Uuid)
		if err != nil {
			return rerrors.Wrap(err, "create default subscription")
		}

		return nil
	})
	if err != nil {
		return domain.User{}, err
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

	token, expiresAt, refreshToken, refreshExpiresAt := newTokenPair()

	newSession := domain.Session{
		UserUuid:         user.Uuid,
		Token:            token,
		ExpiresAt:        expiresAt,
		RefreshToken:     refreshToken,
		RefreshExpiresAt: refreshExpiresAt,
	}

	session, err := s.sessionsRepo.Create(ctx, newSession)
	if err != nil {
		return domain.Session{}, rerrors.Wrap(err, "error creating session")
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

func (s *Service) ValidateToken(ctx context.Context, token string) (domain.User, error) {
	session, user, err := s.sessionsRepo.GetByTokenWithUser(ctx, token)
	if err != nil {
		return domain.User{}, rerrors.Wrap(err, "get session with user")
	}

	if time.Now().After(session.ExpiresAt) {
		return domain.User{}, user_errors.SessionExpired
	}

	return user, nil
}

func (s *Service) LoginViaTelegram(ctx context.Context, idToken string) (domain.Session, error) {
	claims, err := s.tgParser.ParseAndVerifyIdToken(idToken)
	if err != nil {
		return domain.Session{}, rerrors.Wrap(
			user_errors.InvalidTelegramToken,
			"error validating telegram token",
			err.Error(),
		)
	}

	telegramId := strconv.FormatInt(claims.Id, 10)

	var user domain.User

	err = s.txManager.Execute(
		func(tx *sql.Tx) error {
			usersRepo := s.usersRepo.WithTx(tx)
			permissionsRepo := s.permissionsRepo.WithTx(tx)
			subsRepo := s.subsRepo.WithTx(tx)

			var userValue sql.Null[domain.User]

			userValue, err = usersRepo.GetByTelegramId(ctx, telegramId)
			if err != nil {
				return rerrors.Wrap(err, "get user by telegram id")
			}

			if !userValue.Valid {
				userValue.V, err = usersRepo.CreateByUsername(ctx, claims.Login, claims.Picture)
				if err != nil {
					return rerrors.Wrap(err, "create user")
				}
			}

			user = userValue.V

			identity := domain.TelegramIdentity{
				UserUuid:   user.Uuid,
				TelegramId: telegramId,
			}

			err = usersRepo.UpsertTelegramIdentity(ctx, identity)
			if err != nil {
				return rerrors.Wrap(err, "upsert telegram identity")
			}

			err = usersRepo.UpdatePhotoUrl(ctx, user.Uuid, claims.Picture)
			if err != nil {
				return rerrors.Wrap(err, "update user photo url")
			}

			err = permissionsRepo.CreateDefault(ctx, user.Uuid)
			if err != nil {
				return rerrors.Wrap(err, "create default permissions")
			}

			err = subsRepo.CreateDefault(ctx, user.Uuid)
			if err != nil {
				return rerrors.Wrap(err, "create default subscription")
			}

			return nil
		})
	if err != nil {
		return domain.Session{}, err
	}

	token, expiresAt, refreshToken, refreshExpiresAt := newTokenPair()

	newSession := domain.Session{
		UserUuid:         user.Uuid,
		Token:            token,
		ExpiresAt:        expiresAt,
		RefreshToken:     refreshToken,
		RefreshExpiresAt: refreshExpiresAt,
	}

	session, err := s.sessionsRepo.Create(ctx, newSession)
	if err != nil {
		return domain.Session{}, rerrors.Wrap(err, "error creating session")
	}

	return session, nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (domain.Session, error) {
	token, expiresAt, newRefreshToken, newRefreshExpiresAt := newTokenPair()

	newSession := domain.Session{
		Token:            token,
		ExpiresAt:        expiresAt,
		RefreshToken:     newRefreshToken,
		RefreshExpiresAt: newRefreshExpiresAt,
	}

	rotated, err := s.sessionsRepo.RotateByRefreshToken(ctx, refreshToken, newSession)
	if err != nil {
		return domain.Session{}, rerrors.Wrap(err, "error rotating session")
	}

	if !rotated.Valid {
		return domain.Session{}, user_errors.InvalidRefreshToken
	}

	return rotated.V, nil
}

func (s *Service) GetMe(ctx context.Context, userUuid uuid.UUID) (domain.UserDetails, error) {
	details, err := s.usersRepo.GetDetailsById(ctx, userUuid)
	if err != nil {
		return domain.UserDetails{}, rerrors.Wrap(err, "get user details")
	}

	effective, err := s.subscriptions.GetEffective(ctx, userUuid)
	if err != nil {
		return domain.UserDetails{}, rerrors.Wrap(err, "get effective subscription")
	}

	details.EffectiveSubscription = effective

	return details, nil
}

func (s *Service) CheckIsAdmin(ctx context.Context, userUuid uuid.UUID) error {
	perms, err := s.permissionsRepo.Get(ctx, userUuid)
	if err != nil {
		return rerrors.Wrap(err, "get user permissions")
	}

	if !perms.IsAdministrator {
		return user_errors.NotAdmin
	}

	return nil
}

func newTokenPair() (token string, expiresAt time.Time, refreshToken string, refreshExpiresAt time.Time) {
	token = generateToken()
	expiresAt = time.Now().Add(accessTokenTTL)
	refreshToken = generateToken()
	refreshExpiresAt = time.Now().Add(refreshTokenTTL)

	return token, expiresAt, refreshToken, refreshExpiresAt
}

func generateToken() string {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(b)
}
