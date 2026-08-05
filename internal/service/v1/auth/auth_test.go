package auth

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"
	"golang.org/x/crypto/bcrypt"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
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

// fakeUsersRepo is a minimal in-memory repository.Users used to test Service without a real
// database — following the fakeDockerHostsRepo convention in
// internal/service/v1/dockerhosts/dockerhosts_test.go.
type fakeUsersRepo struct {
	byID map[uuid.UUID]domain.User
}

func newFakeUsersRepo() *fakeUsersRepo {
	return &fakeUsersRepo{byID: map[uuid.UUID]domain.User{}}
}

func (f *fakeUsersRepo) Create(_ context.Context, email, passwordHash string) (domain.User, error) {
	user := domain.User{Uuid: uuid.New(), Email: email, PasswordHash: passwordHash}
	f.byID[user.Uuid] = user

	return user, nil
}

func (f *fakeUsersRepo) GetByID(_ context.Context, id uuid.UUID) (domain.User, error) {
	user, ok := f.byID[id]
	if !ok {
		return domain.User{}, errors.New("not found")
	}

	return user, nil
}

func (f *fakeUsersRepo) GetByEmail(_ context.Context, email string) (domain.User, error) {
	for _, user := range f.byID {
		if user.Email == email {
			return user, nil
		}
	}

	return domain.User{}, errors.New("not found")
}

func (f *fakeUsersRepo) FindByEmail(_ context.Context, email string) (sql.Null[domain.User], error) {
	for _, user := range f.byID {
		if user.Email == email {
			return sql.Null[domain.User]{V: user, Valid: true}, nil
		}
	}

	return sql.Null[domain.User]{}, nil
}

func (f *fakeUsersRepo) GetByTelegramId(_ context.Context, _ string) (sql.Null[domain.User], error) {
	return sql.Null[domain.User]{}, nil
}

func (f *fakeUsersRepo) CreateByUsername(_ context.Context, username, photoUrl string) (domain.User, error) {
	user := domain.User{Uuid: uuid.New(), Username: username, PhotoUrl: photoUrl}
	f.byID[user.Uuid] = user

	return user, nil
}

func (f *fakeUsersRepo) UpsertTelegramIdentity(_ context.Context, _ domain.TelegramIdentity) error {
	return nil
}

func (f *fakeUsersRepo) GetTelegramPhotoUrl(_ context.Context, _ uuid.UUID) (string, error) {
	return "", nil
}

func (f *fakeUsersRepo) UpdatePhotoUrl(_ context.Context, _ uuid.UUID, _ string) error {
	return nil
}

func (f *fakeUsersRepo) UpdatePasswordHash(_ context.Context, userUuid uuid.UUID, passwordHash string) error {
	user, ok := f.byID[userUuid]
	if !ok {
		return errors.New("not found")
	}

	user.PasswordHash = passwordHash
	f.byID[userUuid] = user

	return nil
}

func (f *fakeUsersRepo) ListAll(_ context.Context, _ domain.ListUsersReq) ([]domain.User, int64, error) {
	return nil, 0, nil
}

func (f *fakeUsersRepo) GetDetailsById(_ context.Context, _ uuid.UUID) (domain.UserDetails, error) {
	return domain.UserDetails{}, nil
}

func (f *fakeUsersRepo) Delete(_ context.Context, _ uuid.UUID) error {
	return nil
}

func (f *fakeUsersRepo) WithTx(_ *sql.Tx) repository.Users {
	return f
}

// fakePermissionsRepo is a minimal in-memory repository.UserPermissionsRepo.
type fakePermissionsRepo struct {
	byUser map[uuid.UUID]bool
}

func newFakePermissionsRepo() *fakePermissionsRepo {
	return &fakePermissionsRepo{byUser: map[uuid.UUID]bool{}}
}

func (f *fakePermissionsRepo) Get(_ context.Context, userUuid uuid.UUID) (domain.UserPermissions, error) {
	perms := domain.UserPermissions{UserUuid: userUuid, IsAdministrator: f.byUser[userUuid]}

	return perms, nil
}

func (f *fakePermissionsRepo) Upsert(
	_ context.Context, userUuid uuid.UUID, isAdmin bool,
) (domain.UserPermissions, error) {
	f.byUser[userUuid] = isAdmin

	perms := domain.UserPermissions{UserUuid: userUuid, IsAdministrator: isAdmin}

	return perms, nil
}

func (f *fakePermissionsRepo) CreateDefault(_ context.Context, userUuid uuid.UUID) error {
	if _, ok := f.byUser[userUuid]; !ok {
		f.byUser[userUuid] = false
	}

	return nil
}

func (f *fakePermissionsRepo) WithTx(_ *sql.Tx) repository.UserPermissionsRepo {
	return f
}

// fakeSubsRepo is a minimal repository.Subscriptions — only CreateDefault/WithTx are exercised
// by Service, the rest are unused stubs.
type fakeSubsRepo struct{}

func (f *fakeSubsRepo) Upsert(_ context.Context, _ uuid.UUID, _ bool) (domain.Subscription, error) {
	return domain.Subscription{}, nil
}

func (f *fakeSubsRepo) GetByUser(_ context.Context, _ uuid.UUID) (domain.Subscription, error) {
	return domain.Subscription{}, nil
}

func (f *fakeSubsRepo) CreateDefault(_ context.Context, _ uuid.UUID) error {
	return nil
}

func (f *fakeSubsRepo) GetWithPlan(_ context.Context, _ uuid.UUID) (domain.EffectiveSubscription, error) {
	return domain.EffectiveSubscription{}, nil
}

func (f *fakeSubsRepo) UpsertOverrides(_ context.Context, sub domain.Subscription) (domain.Subscription, error) {
	return sub, nil
}

func (f *fakeSubsRepo) WithTx(_ *sql.Tx) repository.Subscriptions {
	return f
}

// fakeSettingsRepo is a minimal in-memory repository.SystemSettingsRepo.
type fakeSettingsRepo struct {
	settings domain.SystemSettings
}

func (f *fakeSettingsRepo) Get(_ context.Context) (domain.SystemSettings, error) {
	return f.settings, nil
}

func (f *fakeSettingsRepo) GetForUpdate(_ context.Context) (domain.SystemSettings, error) {
	return f.settings, nil
}

func (f *fakeSettingsRepo) UpdateAuthMethods(_ context.Context, passwordEnabled, telegramEnabled bool) error {
	f.settings.PasswordAuthEnabled = passwordEnabled
	f.settings.TelegramAuthEnabled = telegramEnabled

	return nil
}

func (f *fakeSettingsRepo) UpdateRegistrationMode(_ context.Context, mode domain.RegistrationMode) error {
	f.settings.RegistrationMode = mode

	return nil
}

func (f *fakeSettingsRepo) SetSetupToken(_ context.Context, tokenHash string, issuedAt time.Time) error {
	f.settings.SetupTokenHash = tokenHash
	f.settings.SetupTokenIssuedAt = issuedAt

	return nil
}

func (f *fakeSettingsRepo) CompleteSetup(_ context.Context) error {
	f.settings.SetupCompleted = true

	return nil
}

func (f *fakeSettingsRepo) WithTx(_ *sql.Tx) repository.SystemSettingsRepo {
	return f
}

// fakeTxManager runs the transaction body directly against a nil *sql.Tx — every fake repo
// above ignores the tx argument passed to WithTx, so this is enough to exercise
// Service.createUserTx's control flow without a real database.
type fakeTxManager struct{}

func (f *fakeTxManager) Execute(do func(tx *sql.Tx) error) error {
	return do(nil)
}

// newTestAuthService builds a Service wired entirely to in-memory fakes, for tests that don't
// need TestServiceRefresh's narrower fakeSessionsRepo-only setup. settingsRepo isn't returned —
// no test needs to inspect it after construction, only seed it via the settings param.
func newTestAuthService(
	settings domain.SystemSettings,
) (svc *Service, usersRepo *fakeUsersRepo, permissionsRepo *fakePermissionsRepo) {
	usersRepo = newFakeUsersRepo()
	permissionsRepo = newFakePermissionsRepo()
	settingsRepo := &fakeSettingsRepo{settings: settings}

	svc = &Service{
		usersRepo:       usersRepo,
		sessionsRepo:    &fakeSessionsRepo{},
		permissionsRepo: permissionsRepo,
		subsRepo:        &fakeSubsRepo{},
		settingsRepo:    settingsRepo,
		txManager:       &fakeTxManager{},
	}

	return svc, usersRepo, permissionsRepo
}

func TestServiceRegister(t *testing.T) {
	t.Run("rejects when setup not completed", func(t *testing.T) {
		settings := domain.SystemSettings{SetupCompleted: false}
		svc, _, _ := newTestAuthService(settings)

		_, err := svc.Register(context.Background(), "a@b.com", "pw")
		if !rerrors.Is(err, user_errors.SetupNotCompleted) {
			t.Fatalf("expected SetupNotCompleted, got %v", err)
		}
	})

	t.Run("rejects self-registration when registration mode is admin_only", func(t *testing.T) {
		settings := domain.SystemSettings{
			SetupCompleted:      true,
			PasswordAuthEnabled: true,
			RegistrationMode:    domain.RegistrationModeAdminOnly,
		}
		svc, _, _ := newTestAuthService(settings)

		_, err := svc.Register(context.Background(), "a@b.com", "pw")
		if !rerrors.Is(err, user_errors.SelfRegistrationDisabled) {
			t.Fatalf("expected SelfRegistrationDisabled, got %v", err)
		}
	})

	t.Run("rejects when password auth is disabled", func(t *testing.T) {
		settings := domain.SystemSettings{
			SetupCompleted:      true,
			PasswordAuthEnabled: false,
			RegistrationMode:    domain.RegistrationModeSelfRegister,
		}
		svc, _, _ := newTestAuthService(settings)

		_, err := svc.Register(context.Background(), "a@b.com", "pw")
		if !rerrors.Is(err, user_errors.AuthMethodDisabled) {
			t.Fatalf("expected AuthMethodDisabled, got %v", err)
		}
	})

	t.Run("creates a non-admin user when allowed", func(t *testing.T) {
		settings := domain.SystemSettings{
			SetupCompleted:      true,
			PasswordAuthEnabled: true,
			RegistrationMode:    domain.RegistrationModeSelfRegister,
		}
		svc, _, permissionsRepo := newTestAuthService(settings)

		user, err := svc.Register(context.Background(), "a@b.com", "pw")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if user.Email != "a@b.com" {
			t.Errorf("expected email a@b.com, got %s", user.Email)
		}

		if permissionsRepo.byUser[user.Uuid] {
			t.Errorf("expected self-registered user to not be an administrator")
		}
	})
}

func TestServiceLogin(t *testing.T) {
	t.Run("rejects when setup not completed", func(t *testing.T) {
		settings := domain.SystemSettings{SetupCompleted: false}
		svc, _, _ := newTestAuthService(settings)

		_, err := svc.Login(context.Background(), "a@b.com", "pw")
		if !rerrors.Is(err, user_errors.SetupNotCompleted) {
			t.Fatalf("expected SetupNotCompleted, got %v", err)
		}
	})

	t.Run("rejects when password auth is disabled", func(t *testing.T) {
		settings := domain.SystemSettings{SetupCompleted: true, PasswordAuthEnabled: false}
		svc, _, _ := newTestAuthService(settings)

		_, err := svc.Login(context.Background(), "a@b.com", "pw")
		if !rerrors.Is(err, user_errors.AuthMethodDisabled) {
			t.Fatalf("expected AuthMethodDisabled, got %v", err)
		}
	})

	t.Run("succeeds with correct credentials", func(t *testing.T) {
		settings := domain.SystemSettings{SetupCompleted: true, PasswordAuthEnabled: true}
		svc, usersRepo, _ := newTestAuthService(settings)

		hash, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
		if err != nil {
			t.Fatalf("unexpected error hashing password: %v", err)
		}

		seeded := domain.User{Uuid: uuid.New(), Email: "a@b.com", PasswordHash: string(hash)}
		usersRepo.byID[seeded.Uuid] = seeded

		session, err := svc.Login(context.Background(), "a@b.com", "correct-password")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if session.UserUuid != seeded.Uuid {
			t.Errorf("expected session for user %v, got %v", seeded.Uuid, session.UserUuid)
		}
	})
}

func TestServiceLoginViaTelegram(t *testing.T) {
	t.Run("rejects when setup not completed", func(t *testing.T) {
		settings := domain.SystemSettings{SetupCompleted: false}
		svc, _, _ := newTestAuthService(settings)

		_, err := svc.LoginViaTelegram(context.Background(), "whatever")
		if !rerrors.Is(err, user_errors.SetupNotCompleted) {
			t.Fatalf("expected SetupNotCompleted, got %v", err)
		}
	})

	t.Run("rejects when telegram auth is disabled", func(t *testing.T) {
		settings := domain.SystemSettings{SetupCompleted: true, TelegramAuthEnabled: false}
		svc, _, _ := newTestAuthService(settings)

		_, err := svc.LoginViaTelegram(context.Background(), "whatever")
		if !rerrors.Is(err, user_errors.AuthMethodDisabled) {
			t.Fatalf("expected AuthMethodDisabled, got %v", err)
		}
	})
}

func TestServiceRegisterAdmin(t *testing.T) {
	t.Run("bypasses setup checks and creates an administrator", func(t *testing.T) {
		// Setup deliberately not completed and no auth method enabled — RegisterAdmin must not
		// consult settings at all, since the wizard runs before setup is marked complete.
		settings := domain.SystemSettings{SetupCompleted: false}
		svc, _, permissionsRepo := newTestAuthService(settings)

		user, err := svc.RegisterAdmin(context.Background(), "admin@b.com", "pw")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !permissionsRepo.byUser[user.Uuid] {
			t.Errorf("expected wizard-created user to be an administrator")
		}
	})
}

func TestServiceCreateUserUnchecked(t *testing.T) {
	t.Run("bypasses registration-mode/auth-method checks and creates a non-admin", func(t *testing.T) {
		settings := domain.SystemSettings{
			SetupCompleted:      true,
			RegistrationMode:    domain.RegistrationModeAdminOnly,
			PasswordAuthEnabled: false,
		}
		svc, _, permissionsRepo := newTestAuthService(settings)

		user, err := svc.CreateUserUnchecked(context.Background(), "u@b.com", "pw")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if permissionsRepo.byUser[user.Uuid] {
			t.Errorf("expected admin-created user to not be an administrator by default")
		}
	})
}

func TestServiceLoginOrRegisterAdminViaTelegram(t *testing.T) {
	t.Run("bypasses setup/auth-method checks", func(t *testing.T) {
		// Setup deliberately not completed and telegram auth disabled — this method must not
		// consult settings at all. The malformed token below fails at JWT parsing before ever
		// reaching the JWKS fetch, so this stays network-free; asserting the error is
		// InvalidTelegramToken (not SetupNotCompleted/AuthMethodDisabled) confirms no settings
		// gate ran first.
		settings := domain.SystemSettings{SetupCompleted: false, TelegramAuthEnabled: false}
		svc, _, _ := newTestAuthService(settings)

		_, err := svc.LoginOrRegisterAdminViaTelegram(context.Background(), "not-a-valid-jwt")
		if rerrors.Is(err, user_errors.SetupNotCompleted) || rerrors.Is(err, user_errors.AuthMethodDisabled) {
			t.Fatalf("expected wizard telegram login to bypass setup gating entirely, got %v", err)
		}

		// rerrors.Wrap(user_errors.InvalidTelegramToken, ...) (in loginOrRegisterViaTelegram)
		// re-wraps with extra context, which rerrors.Is can't see through (it only compares the
		// outermost Error's message) — check the message text instead of the sentinel identity.
		if err == nil || !strings.Contains(err.Error(), "invalid telegram token") {
			t.Fatalf("expected an invalid-telegram-token error for a malformed token, got %v", err)
		}
	})
}

func TestServiceChangePassword(t *testing.T) {
	t.Run("updates the stored password hash", func(t *testing.T) {
		settings := domain.SystemSettings{}
		svc, usersRepo, _ := newTestAuthService(settings)

		seeded := domain.User{Uuid: uuid.New(), Email: "a@b.com", PasswordHash: "old-hash"}
		usersRepo.byID[seeded.Uuid] = seeded

		err := svc.ChangePassword(context.Background(), seeded.Uuid, "new-password")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		updated := usersRepo.byID[seeded.Uuid]

		err = bcrypt.CompareHashAndPassword([]byte(updated.PasswordHash), []byte("new-password"))
		if err != nil {
			t.Errorf("expected stored hash to match the new password: %v", err)
		}
	})
}
