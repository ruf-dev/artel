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
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

type Service struct {
	usersRepo       repository.Users
	sessionsRepo    repository.Sessions
	permissionsRepo repository.UserPermissionsRepo
	subsRepo        repository.Subscriptions
	txManager       tx_manager.TxManager
	tgParser        telegram.TokenParser
}

func New(repo repository.Repo, telegramClientId string) *Service {
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

		if err = s.permissionsRepo.WithTx(tx).CreateDefault(ctx, user.Uuid); err != nil {
			return rerrors.Wrap(err, "create default permissions")
		}

		if err = s.subsRepo.WithTx(tx).CreateDefault(ctx, user.Uuid); err != nil {
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

func (s *Service) ValidateToken(ctx context.Context, token string) (domain.User, error) {
	session, err := s.sessionsRepo.GetByToken(ctx, token)
	if err != nil {
		return domain.User{}, rerrors.Wrap(err, "get session by token")
	}

	if time.Now().After(session.ExpiresAt) {
		return domain.User{}, user_errors.SessionExpired
	}

	user, err := s.usersRepo.GetByID(ctx, session.UserUuid)
	if err != nil {
		return domain.User{}, rerrors.Wrap(err, "get user by id")
	}

	return user, nil
}

func (s *Service) LoginViaTelegram(ctx context.Context, idToken string) (domain.Session, error) {
	claims, err := s.tgParser.ParseAndVerifyIdToken(idToken)
	if err != nil {
		return domain.Session{}, rerrors.Wrap(user_errors.InvalidTelegramToken, "error validating telegram token", err.Error())
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

	sessionToken := generateToken()
	expiresAt := time.Now().Add(24 * time.Hour)

	session, err := s.sessionsRepo.Create(ctx, user.Uuid, sessionToken, expiresAt)
	if err != nil {
		return domain.Session{}, rerrors.Wrap(err, "create session")
	}

	return session, nil
}

func (s *Service) GetMe(ctx context.Context, userUuid uuid.UUID) (domain.User, domain.UserPermissions, error) {
	user, err := s.usersRepo.GetByID(ctx, userUuid)
	if err != nil {
		return domain.User{}, domain.UserPermissions{}, rerrors.Wrap(err, "get user by id")
	}

	perms, err := s.permissionsRepo.Get(ctx, userUuid)
	if err != nil {
		return domain.User{}, domain.UserPermissions{}, rerrors.Wrap(err, "get user permissions")
	}

	return user, perms, nil
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

func generateToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(b)
}
