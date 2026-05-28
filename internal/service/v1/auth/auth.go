package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.redsock.ru/rerrors"
	"golang.org/x/crypto/bcrypt"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/repository/pg/tx_manager"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

type Service struct {
	usersRepo        repository.Users
	sessionsRepo     repository.Sessions
	permissionsRepo  repository.UserPermissionsRepo
	subsRepo         repository.Subscriptions
	txManager        tx_manager.TxManager
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
		permissionsRepo:  repo.UserPermissions(),
		subsRepo:         repo.Subscriptions(),
		txManager:        repo.TxManager(),
		jwksClient:       jwks,
		telegramClientId: telegramClientId,
	}, nil
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

type telegramClaims struct {
	jwt.RegisteredClaims
	PhotoUrl string `json:"photo_url"`
}

type telegramUserInfo struct {
	PhotoUrl string `json:"photo_url"`
}

func fetchTelegramUserInfo(ctx context.Context, idToken string) (telegramUserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://oauth.telegram.org/userinfo", nil)
	if err != nil {
		return telegramUserInfo{}, err
	}
	req.Header.Set("Authorization", "Bearer "+idToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return telegramUserInfo{}, err
	}
	defer resp.Body.Close()
	var info telegramUserInfo
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return telegramUserInfo{}, err
	}
	return info, nil
}

func (s *Service) LoginViaTelegram(ctx context.Context, idToken string) (domain.Session, error) {
	claims := &telegramClaims{}
	token, err := jwt.ParseWithClaims(idToken, claims, s.jwksClient.Keyfunc,
		jwt.WithValidMethods([]string{"RS256", "ES256"}))
	if err != nil {
		return domain.Session{}, rerrors.Wrap(err, "parse telegram id_token")
	}

	if !token.Valid {
		return domain.Session{}, user_errors.InvalidTelegramToken
	}

	telegramId := claims.Subject

	photoUrl := claims.PhotoUrl
	if photoUrl == "" {
		if info, fetchErr := fetchTelegramUserInfo(ctx, idToken); fetchErr == nil {
			photoUrl = info.PhotoUrl
		}
	}

	var user domain.User

	err = s.txManager.Execute(
		func(tx *sql.Tx) error {
			usersRepo := s.usersRepo.WithTx(tx)
			permissionsRepo := s.permissionsRepo.WithTx(tx)
			subsRepo := s.subsRepo.WithTx(tx)

			user, err = usersRepo.UpsertByTelegramId(ctx, telegramId, telegramId, photoUrl)
			if err != nil {
				return rerrors.Wrap(err, "upsert telegram user")
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

	photoUrl, err := s.usersRepo.GetTelegramPhotoUrl(ctx, userUuid)
	if err != nil {
		return domain.User{}, domain.UserPermissions{}, rerrors.Wrap(err, "get telegram photo url")
	}
	user.PhotoUrl = photoUrl

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
