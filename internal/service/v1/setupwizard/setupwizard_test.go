package setupwizard

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

// fakeSettingsRepo is a minimal in-memory repository.SystemSettingsRepo used to test Service
// without a real database — following the fakeDockerHostsRepo convention in
// internal/service/v1/dockerhosts/dockerhosts_test.go.
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

// fakeAuthService is a minimal service.AuthService test double. Only the methods
// Service.CompleteSetup calls (RegisterAdmin, Login, LoginOrRegisterAdminViaTelegram) are
// configurable; every other method is an unused stub.
type fakeAuthService struct {
	registerAdminUser domain.User
	registerAdminErr  error

	loginSession domain.Session
	loginErr     error

	telegramSession domain.Session
	telegramErr     error
}

func (f *fakeAuthService) Register(_ context.Context, _, _ string) (domain.User, error) {
	return domain.User{}, nil
}

func (f *fakeAuthService) Login(_ context.Context, _, _ string) (domain.Session, error) {
	return f.loginSession, f.loginErr
}

func (f *fakeAuthService) Logout(_ context.Context, _ string) error {
	return nil
}

func (f *fakeAuthService) ValidateToken(_ context.Context, _ string) (domain.User, error) {
	return domain.User{}, nil
}

func (f *fakeAuthService) LoginViaTelegram(_ context.Context, _ string) (domain.Session, error) {
	return domain.Session{}, nil
}

func (f *fakeAuthService) GetMe(_ context.Context, _ uuid.UUID) (domain.UserDetails, error) {
	return domain.UserDetails{}, nil
}

func (f *fakeAuthService) CheckIsAdmin(_ context.Context, _ uuid.UUID) error {
	return nil
}

func (f *fakeAuthService) Refresh(_ context.Context, _ string) (domain.Session, error) {
	return domain.Session{}, nil
}

func (f *fakeAuthService) EnsureNoAuthUser(_ context.Context) (domain.User, error) {
	return domain.User{}, nil
}

func (f *fakeAuthService) RegisterAdmin(_ context.Context, email, _ string) (domain.User, error) {
	if f.registerAdminErr != nil {
		return domain.User{}, f.registerAdminErr
	}

	user := f.registerAdminUser
	user.Email = email

	return user, nil
}

func (f *fakeAuthService) CreateUserUnchecked(_ context.Context, _, _ string) (domain.User, error) {
	return domain.User{}, nil
}

func (f *fakeAuthService) LoginOrRegisterAdminViaTelegram(_ context.Context, _ string) (domain.Session, error) {
	return f.telegramSession, f.telegramErr
}

func (f *fakeAuthService) ChangePassword(_ context.Context, _ uuid.UUID, _ string) error {
	return nil
}

func TestCurrentStatus(t *testing.T) {
	settingsRepo := &fakeSettingsRepo{settings: domain.SystemSettings{SetupCompleted: true}}
	svc := New(settingsRepo, &fakeAuthService{})

	status, err := svc.CurrentStatus(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !status.SetupCompleted {
		t.Errorf("expected CurrentStatus to reflect the repo's SetupCompleted=true")
	}
}

func TestSubmitToken(t *testing.T) {
	t.Run("succeeds with the correct token", func(t *testing.T) {
		settingsRepo := &fakeSettingsRepo{}
		svc := New(settingsRepo, &fakeAuthService{})

		plaintext, err := svc.RegenerateSetupToken(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		wizardToken, err := svc.SubmitToken(context.Background(), plaintext)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if wizardToken == "" {
			t.Errorf("expected a non-empty wizard session token")
		}
	})

	t.Run("rejects an incorrect token", func(t *testing.T) {
		settingsRepo := &fakeSettingsRepo{}
		svc := New(settingsRepo, &fakeAuthService{})

		_, err := svc.RegenerateSetupToken(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		_, err = svc.SubmitToken(context.Background(), "wrong-token")
		if !rerrors.Is(err, user_errors.WizardSessionInvalid) {
			t.Fatalf("expected WizardSessionInvalid, got %v", err)
		}
	})

	t.Run("rejects when no token has ever been issued", func(t *testing.T) {
		settingsRepo := &fakeSettingsRepo{}
		svc := New(settingsRepo, &fakeAuthService{})

		_, err := svc.SubmitToken(context.Background(), "")
		if !rerrors.Is(err, user_errors.WizardSessionInvalid) {
			t.Fatalf("expected WizardSessionInvalid, got %v", err)
		}
	})

	t.Run("rejects when setup is already completed", func(t *testing.T) {
		settings := domain.SystemSettings{SetupCompleted: true}
		settingsRepo := &fakeSettingsRepo{settings: settings}
		svc := New(settingsRepo, &fakeAuthService{})

		_, err := svc.SubmitToken(context.Background(), "whatever")
		if !rerrors.Is(err, user_errors.SetupAlreadyCompleted) {
			t.Fatalf("expected SetupAlreadyCompleted, got %v", err)
		}
	})
}

func TestSelectAuthMethods(t *testing.T) {
	t.Run("rejects a missing wizard session", func(t *testing.T) {
		svc := New(&fakeSettingsRepo{}, &fakeAuthService{})

		err := svc.SelectAuthMethods(context.Background(), "does-not-exist", true, false)
		if !rerrors.Is(err, user_errors.WizardSessionInvalid) {
			t.Fatalf("expected WizardSessionInvalid, got %v", err)
		}
	})

	t.Run("rejects an expired wizard session", func(t *testing.T) {
		svc := New(&fakeSettingsRepo{}, &fakeAuthService{})

		wizardToken := "expired-token"
		svc.wizardSessions[wizardToken] = wizardSession{createdAt: time.Now().Add(-wizardSessionTTL - time.Minute)}

		err := svc.SelectAuthMethods(context.Background(), wizardToken, true, false)
		if !rerrors.Is(err, user_errors.WizardSessionInvalid) {
			t.Fatalf("expected WizardSessionInvalid, got %v", err)
		}
	})

	t.Run("rejects disabling both auth methods", func(t *testing.T) {
		settingsRepo := &fakeSettingsRepo{}
		svc := New(settingsRepo, &fakeAuthService{})

		wizardToken := "valid-token"
		svc.wizardSessions[wizardToken] = wizardSession{createdAt: time.Now()}

		err := svc.SelectAuthMethods(context.Background(), wizardToken, false, false)
		if !rerrors.Is(err, user_errors.AtLeastOneAuthMethodRequired) {
			t.Fatalf("expected AtLeastOneAuthMethodRequired, got %v", err)
		}
	})

	t.Run("persists the selected auth methods", func(t *testing.T) {
		settingsRepo := &fakeSettingsRepo{}
		svc := New(settingsRepo, &fakeAuthService{})

		wizardToken := "valid-token"
		svc.wizardSessions[wizardToken] = wizardSession{createdAt: time.Now()}

		err := svc.SelectAuthMethods(context.Background(), wizardToken, true, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !settingsRepo.settings.PasswordAuthEnabled || settingsRepo.settings.TelegramAuthEnabled {
			t.Errorf("expected password enabled and telegram disabled, got %+v", settingsRepo.settings)
		}
	})
}

func TestSelectRegistrationMode(t *testing.T) {
	t.Run("rejects an invalid wizard session", func(t *testing.T) {
		svc := New(&fakeSettingsRepo{}, &fakeAuthService{})

		err := svc.SelectRegistrationMode(context.Background(), "does-not-exist", domain.RegistrationModeSelfRegister)
		if !rerrors.Is(err, user_errors.WizardSessionInvalid) {
			t.Fatalf("expected WizardSessionInvalid, got %v", err)
		}
	})

	t.Run("persists the selected registration mode", func(t *testing.T) {
		settingsRepo := &fakeSettingsRepo{}
		svc := New(settingsRepo, &fakeAuthService{})

		wizardToken := "valid-token"
		svc.wizardSessions[wizardToken] = wizardSession{createdAt: time.Now()}

		err := svc.SelectRegistrationMode(context.Background(), wizardToken, domain.RegistrationModeSelfRegister)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if settingsRepo.settings.RegistrationMode != domain.RegistrationModeSelfRegister {
			t.Errorf("expected self_register mode, got %v", settingsRepo.settings.RegistrationMode)
		}
	})
}

func TestCompleteSetup(t *testing.T) {
	t.Run("rejects an invalid wizard session", func(t *testing.T) {
		svc := New(&fakeSettingsRepo{}, &fakeAuthService{})

		_, err := svc.CompleteSetup(context.Background(), "does-not-exist", "a@b.com", "pw", "")
		if !rerrors.Is(err, user_errors.WizardSessionInvalid) {
			t.Fatalf("expected WizardSessionInvalid, got %v", err)
		}
	})

	t.Run("rejects when setup is already completed", func(t *testing.T) {
		settings := domain.SystemSettings{SetupCompleted: true, PasswordAuthEnabled: true}
		settingsRepo := &fakeSettingsRepo{settings: settings}
		svc := New(settingsRepo, &fakeAuthService{})

		wizardToken := "valid-token"
		svc.wizardSessions[wizardToken] = wizardSession{createdAt: time.Now()}

		_, err := svc.CompleteSetup(context.Background(), wizardToken, "a@b.com", "pw", "")
		if !rerrors.Is(err, user_errors.SetupAlreadyCompleted) {
			t.Fatalf("expected SetupAlreadyCompleted, got %v", err)
		}
	})

	t.Run("rejects a password completion when password auth was not selected", func(t *testing.T) {
		settings := domain.SystemSettings{PasswordAuthEnabled: false}
		settingsRepo := &fakeSettingsRepo{settings: settings}
		svc := New(settingsRepo, &fakeAuthService{})

		wizardToken := "valid-token"
		svc.wizardSessions[wizardToken] = wizardSession{createdAt: time.Now()}

		_, err := svc.CompleteSetup(context.Background(), wizardToken, "a@b.com", "pw", "")
		if !rerrors.Is(err, user_errors.AuthMethodNotEnabled) {
			t.Fatalf("expected AuthMethodNotEnabled, got %v", err)
		}
	})

	t.Run("rejects when neither password nor telegram id token is provided", func(t *testing.T) {
		settings := domain.SystemSettings{PasswordAuthEnabled: true, TelegramAuthEnabled: true}
		settingsRepo := &fakeSettingsRepo{settings: settings}
		svc := New(settingsRepo, &fakeAuthService{})

		wizardToken := "valid-token"
		svc.wizardSessions[wizardToken] = wizardSession{createdAt: time.Now()}

		_, err := svc.CompleteSetup(context.Background(), wizardToken, "", "", "")
		if !rerrors.Is(err, user_errors.AuthMethodNotEnabled) {
			t.Fatalf("expected AuthMethodNotEnabled, got %v", err)
		}
	})

	t.Run("completes setup via password and returns an admin session", func(t *testing.T) {
		adminUuid := uuid.New()
		expectedSession := domain.Session{UserUuid: adminUuid, Token: "session-token"}

		settings := domain.SystemSettings{PasswordAuthEnabled: true}
		settingsRepo := &fakeSettingsRepo{settings: settings}
		authSvc := &fakeAuthService{
			registerAdminUser: domain.User{Uuid: adminUuid},
			loginSession:      expectedSession,
		}
		svc := New(settingsRepo, authSvc)

		wizardToken := "valid-token"
		svc.wizardSessions[wizardToken] = wizardSession{createdAt: time.Now()}

		session, err := svc.CompleteSetup(context.Background(), wizardToken, "admin@b.com", "pw", "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if session.UserUuid != adminUuid {
			t.Errorf("expected session for admin %v, got %v", adminUuid, session.UserUuid)
		}

		if !settingsRepo.settings.SetupCompleted {
			t.Errorf("expected setup to be marked completed")
		}

		if _, stillExists := svc.wizardSessions[wizardToken]; stillExists {
			t.Errorf("expected the wizard session to be cleared after completion")
		}
	})

	t.Run("completes setup via telegram and returns an admin session", func(t *testing.T) {
		adminUuid := uuid.New()
		expectedSession := domain.Session{UserUuid: adminUuid, Token: "tg-session-token"}

		settings := domain.SystemSettings{TelegramAuthEnabled: true}
		settingsRepo := &fakeSettingsRepo{settings: settings}
		authSvc := &fakeAuthService{telegramSession: expectedSession}
		svc := New(settingsRepo, authSvc)

		wizardToken := "valid-token"
		svc.wizardSessions[wizardToken] = wizardSession{createdAt: time.Now()}

		session, err := svc.CompleteSetup(context.Background(), wizardToken, "", "", "tg-id-token")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if session.UserUuid != adminUuid {
			t.Errorf("expected session for admin %v, got %v", adminUuid, session.UserUuid)
		}

		if !settingsRepo.settings.SetupCompleted {
			t.Errorf("expected setup to be marked completed")
		}

		if _, stillExists := svc.wizardSessions[wizardToken]; stillExists {
			t.Errorf("expected the wizard session to be cleared after completion")
		}
	})
}
