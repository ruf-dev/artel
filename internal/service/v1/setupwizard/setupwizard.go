// Package setupwizard implements the first-run setup wizard — see
// internal/service/interfaces.go's SetupWizardService doc comment for the flow this drives.
package setupwizard

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"sync"
	"time"

	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

// wizardSessionTTL bounds how long a wizard session token minted by SubmitToken stays valid.
const wizardSessionTTL = 30 * time.Minute

// wizardSession is the in-memory record backing a wizard session token — never persisted, since
// the wizard only ever runs once per instance and a server restart mid-wizard is fine to require
// re-submitting the setup token for.
type wizardSession struct {
	createdAt time.Time
}

type Service struct {
	settingsRepo repository.SystemSettingsRepo
	authSvc      service.AuthService

	mu             sync.Mutex
	wizardSessions map[string]wizardSession
}

func New(settingsRepo repository.SystemSettingsRepo, authSvc service.AuthService) *Service {
	return &Service{
		settingsRepo:   settingsRepo,
		authSvc:        authSvc,
		wizardSessions: make(map[string]wizardSession),
	}
}

func (s *Service) CurrentStatus(ctx context.Context) (domain.SystemSettings, error) {
	settings, err := s.settingsRepo.Get(ctx)
	if err != nil {
		return domain.SystemSettings{}, rerrors.Wrap(err, "get system settings")
	}

	return settings, nil
}

func (s *Service) RegenerateSetupToken(ctx context.Context) (string, error) {
	plaintext := generateRandomToken()

	hash := hashToken(plaintext)

	err := s.settingsRepo.SetSetupToken(ctx, hash, time.Now())
	if err != nil {
		return "", rerrors.Wrap(err, "set setup token")
	}

	return plaintext, nil
}

func (s *Service) SubmitToken(ctx context.Context, token string) (string, error) {
	settings, err := s.settingsRepo.Get(ctx)
	if err != nil {
		return "", rerrors.Wrap(err, "get system settings")
	}

	if settings.SetupCompleted {
		return "", user_errors.SetupAlreadyCompleted
	}

	presentedHash := hashToken(token)

	match := subtle.ConstantTimeCompare([]byte(presentedHash), []byte(settings.SetupTokenHash)) == 1
	if settings.SetupTokenHash == "" || !match {
		return "", user_errors.WizardSessionInvalid
	}

	wizardSessionToken := generateRandomToken()

	s.mu.Lock()
	s.wizardSessions[wizardSessionToken] = wizardSession{createdAt: time.Now()}
	s.mu.Unlock()

	return wizardSessionToken, nil
}

// validateWizardSession rejects a missing or expired wizard session token — every wizard step
// after SubmitToken must call this first.
func (s *Service) validateWizardSession(wizardSessionToken string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.wizardSessions[wizardSessionToken]
	if !ok {
		return user_errors.WizardSessionInvalid
	}

	if time.Since(session.createdAt) > wizardSessionTTL {
		delete(s.wizardSessions, wizardSessionToken)

		return user_errors.WizardSessionInvalid
	}

	return nil
}

func (s *Service) clearWizardSession(wizardSessionToken string) {
	s.mu.Lock()
	delete(s.wizardSessions, wizardSessionToken)
	s.mu.Unlock()
}

func (s *Service) SelectAuthMethods(
	ctx context.Context, wizardSessionToken string, passwordEnabled, telegramEnabled bool,
) error {
	err := s.validateWizardSession(wizardSessionToken)
	if err != nil {
		return err
	}

	if !passwordEnabled && !telegramEnabled {
		return user_errors.AtLeastOneAuthMethodRequired
	}

	err = s.settingsRepo.UpdateAuthMethods(ctx, passwordEnabled, telegramEnabled)
	if err != nil {
		return rerrors.Wrap(err, "update auth methods")
	}

	return nil
}

func (s *Service) SelectRegistrationMode(
	ctx context.Context, wizardSessionToken string, mode domain.RegistrationMode,
) error {
	err := s.validateWizardSession(wizardSessionToken)
	if err != nil {
		return err
	}

	err = s.settingsRepo.UpdateRegistrationMode(ctx, mode)
	if err != nil {
		return rerrors.Wrap(err, "update registration mode")
	}

	return nil
}

func (s *Service) CompleteSetup(
	ctx context.Context, wizardSessionToken string, email, password, telegramIdToken string,
) (domain.Session, error) {
	err := s.validateWizardSession(wizardSessionToken)
	if err != nil {
		return domain.Session{}, err
	}

	settings, err := s.settingsRepo.Get(ctx)
	if err != nil {
		return domain.Session{}, rerrors.Wrap(err, "get system settings")
	}

	if settings.SetupCompleted {
		return domain.Session{}, user_errors.SetupAlreadyCompleted
	}

	var session domain.Session

	switch {
	case password != "":
		if !settings.PasswordAuthEnabled {
			return domain.Session{}, user_errors.AuthMethodNotEnabled
		}

		session, err = s.completeSetupViaPassword(ctx, email, password)
	case telegramIdToken != "":
		if !settings.TelegramAuthEnabled {
			return domain.Session{}, user_errors.AuthMethodNotEnabled
		}

		session, err = s.completeSetupViaTelegram(ctx, telegramIdToken)
	default:
		return domain.Session{}, user_errors.AuthMethodNotEnabled
	}

	if err != nil {
		return domain.Session{}, err
	}

	s.clearWizardSession(wizardSessionToken)

	return session, nil
}

func (s *Service) completeSetupViaPassword(ctx context.Context, email, password string) (domain.Session, error) {
	_, err := s.authSvc.RegisterAdmin(ctx, email, password)
	if err != nil {
		return domain.Session{}, rerrors.Wrap(err, "register admin")
	}

	// CompleteSetup must run before Login: Login gates on settings.SetupCompleted.
	err = s.settingsRepo.CompleteSetup(ctx)
	if err != nil {
		return domain.Session{}, rerrors.Wrap(err, "complete setup")
	}

	session, err := s.authSvc.Login(ctx, email, password)
	if err != nil {
		return domain.Session{}, rerrors.Wrap(err, "login admin")
	}

	return session, nil
}

func (s *Service) completeSetupViaTelegram(ctx context.Context, telegramIdToken string) (domain.Session, error) {
	session, err := s.authSvc.LoginOrRegisterAdminViaTelegram(ctx, telegramIdToken)
	if err != nil {
		return domain.Session{}, rerrors.Wrap(err, "login or register admin via telegram")
	}

	err = s.settingsRepo.CompleteSetup(ctx)
	if err != nil {
		return domain.Session{}, rerrors.Wrap(err, "complete setup")
	}

	return session, nil
}

func generateRandomToken() string {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(b)
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))

	return hex.EncodeToString(sum[:])
}
