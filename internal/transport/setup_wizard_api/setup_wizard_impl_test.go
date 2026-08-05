package setup_wizard_api

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service"
	"github.com/stretchr/testify/require"
)

// fakeSetupWizardService is a minimal in-memory service.SetupWizardService used to test
// SetupWizardImpl's handlers without a real database — following the fakeAuthService convention
// in internal/transport/auth_api/auth_impl_test.go (a plain struct with configurable
// result/err fields per method, not an embedded-interface trick). Every method panics if a test
// doesn't configure the field it reads, so an accidental extra call surfaces immediately.
type fakeSetupWizardService struct {
	currentStatusResult domain.SystemSettings
	currentStatusErr    error

	submitTokenResult string
	submitTokenErr    error

	selectAuthMethodsErr error

	selectRegistrationModeErr error

	completeSetupResult domain.Session
	completeSetupErr    error

	completeSetupCalledWith struct {
		wizardSessionToken string
		email              string
		password           string
		telegramIdToken    string
	}
}

func (f *fakeSetupWizardService) CurrentStatus(_ context.Context) (domain.SystemSettings, error) {
	if f.currentStatusErr != nil {
		return domain.SystemSettings{}, f.currentStatusErr
	}

	return f.currentStatusResult, nil
}

func (f *fakeSetupWizardService) RegenerateSetupToken(_ context.Context) (string, error) {
	panic("not implemented")
}

func (f *fakeSetupWizardService) SubmitToken(_ context.Context, _ string) (string, error) {
	if f.submitTokenErr != nil {
		return "", f.submitTokenErr
	}

	return f.submitTokenResult, nil
}

func (f *fakeSetupWizardService) SelectAuthMethods(_ context.Context, _ string, _, _ bool) error {
	return f.selectAuthMethodsErr
}

func (f *fakeSetupWizardService) SelectRegistrationMode(_ context.Context, _ string, _ domain.RegistrationMode) error {
	return f.selectRegistrationModeErr
}

func (f *fakeSetupWizardService) CompleteSetup(
	_ context.Context, wizardSessionToken, email, password, telegramIdToken string,
) (domain.Session, error) {
	f.completeSetupCalledWith.wizardSessionToken = wizardSessionToken
	f.completeSetupCalledWith.email = email
	f.completeSetupCalledWith.password = password
	f.completeSetupCalledWith.telegramIdToken = telegramIdToken

	if f.completeSetupErr != nil {
		return domain.Session{}, f.completeSetupErr
	}

	return f.completeSetupResult, nil
}

var _ service.SetupWizardService = (*fakeSetupWizardService)(nil)

func TestSetupWizardImpl_GetStatus(t *testing.T) {
	testCases := []struct {
		name           string
		setupCompleted bool
		tokenHash      string
		wantPending    bool
	}{
		{name: "setup not completed, no token issued", setupCompleted: false, tokenHash: "", wantPending: false},
		{name: "setup not completed, token pending", setupCompleted: false, tokenHash: "hash", wantPending: true},
		{name: "setup completed", setupCompleted: true, tokenHash: "", wantPending: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			settings := domain.SystemSettings{SetupCompleted: tc.setupCompleted, SetupTokenHash: tc.tokenHash}
			svc := &fakeSetupWizardService{currentStatusResult: settings}
			impl := New(svc, false)

			resp, err := impl.GetStatus(context.Background(), &artel_api.GetStatus_Request{})
			require.NoError(t, err)
			require.Equal(t, tc.setupCompleted, resp.SetupCompleted)
			require.Equal(t, tc.wantPending, resp.TokenPending)
		})
	}
}

func TestSetupWizardImpl_GetStatus_PropagatesError(t *testing.T) {
	svcErr := errors.New("db unavailable")
	svc := &fakeSetupWizardService{currentStatusErr: svcErr}
	impl := New(svc, false)

	resp, err := impl.GetStatus(context.Background(), &artel_api.GetStatus_Request{})
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestSetupWizardImpl_SubmitToken(t *testing.T) {
	svc := &fakeSetupWizardService{submitTokenResult: "wizard-session-token"}
	impl := New(svc, false)

	req := &artel_api.SubmitToken_Request{Token: "raw-setup-token"}

	resp, err := impl.SubmitToken(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, "wizard-session-token", resp.WizardSessionToken)
}

func TestSetupWizardImpl_SubmitToken_PropagatesError(t *testing.T) {
	svcErr := errors.New("invalid token")
	svc := &fakeSetupWizardService{submitTokenErr: svcErr}
	impl := New(svc, false)

	resp, err := impl.SubmitToken(context.Background(), &artel_api.SubmitToken_Request{Token: "bad"})
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestSetupWizardImpl_SelectAuthMethods(t *testing.T) {
	svc := &fakeSetupWizardService{}
	impl := New(svc, false)

	req := &artel_api.SelectAuthMethods_Request{
		WizardSessionToken: "wizard-token",
		PasswordEnabled:    true,
		TelegramEnabled:    false,
	}

	resp, err := impl.SelectAuthMethods(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestSetupWizardImpl_SelectAuthMethods_PropagatesError(t *testing.T) {
	svcErr := errors.New("at least one auth method required")
	svc := &fakeSetupWizardService{selectAuthMethodsErr: svcErr}
	impl := New(svc, false)

	req := &artel_api.SelectAuthMethods_Request{WizardSessionToken: "wizard-token"}

	resp, err := impl.SelectAuthMethods(context.Background(), req)
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestSetupWizardImpl_SelectRegistrationMode(t *testing.T) {
	svc := &fakeSetupWizardService{}
	impl := New(svc, false)

	req := &artel_api.SelectRegistrationMode_Request{
		WizardSessionToken: "wizard-token",
		Mode:               artel_api.RegistrationMode_SELF_REGISTER,
	}

	resp, err := impl.SelectRegistrationMode(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestSetupWizardImpl_SelectRegistrationMode_PropagatesError(t *testing.T) {
	svcErr := errors.New("invalid wizard session")
	svc := &fakeSetupWizardService{selectRegistrationModeErr: svcErr}
	impl := New(svc, false)

	req := &artel_api.SelectRegistrationMode_Request{WizardSessionToken: "wizard-token"}

	resp, err := impl.SelectRegistrationMode(context.Background(), req)
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestSetupWizardImpl_CompleteSetup_Password(t *testing.T) {
	expiresAt := time.Now().Add(time.Hour)
	refreshExpiresAt := time.Now().Add(24 * time.Hour)
	session := domain.Session{
		Token:            "access-token",
		ExpiresAt:        expiresAt,
		RefreshToken:     "refresh-token",
		RefreshExpiresAt: refreshExpiresAt,
	}
	svc := &fakeSetupWizardService{completeSetupResult: session}
	impl := New(svc, false)

	req := &artel_api.CompleteSetup_Request{
		WizardSessionToken: "wizard-token",
		Method: &artel_api.CompleteSetup_Request_Password{
			Password: &artel_api.PasswordCredentials{Email: "admin@example.com", Password: "s3cret-p4ss"},
		},
	}

	resp, err := impl.CompleteSetup(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, "access-token", resp.Token)
	require.Equal(t, "refresh-token", resp.RefreshToken)
	require.True(t, resp.ExpiresAt.AsTime().Equal(expiresAt))
	require.True(t, resp.RefreshExpiresAt.AsTime().Equal(refreshExpiresAt))
	require.Equal(t, "wizard-token", svc.completeSetupCalledWith.wizardSessionToken)
	require.Equal(t, "admin@example.com", svc.completeSetupCalledWith.email)
	require.Equal(t, "s3cret-p4ss", svc.completeSetupCalledWith.password)
	require.Empty(t, svc.completeSetupCalledWith.telegramIdToken)
}

func TestSetupWizardImpl_CompleteSetup_Telegram(t *testing.T) {
	svc := &fakeSetupWizardService{completeSetupResult: domain.Session{Token: "access-token"}}
	impl := New(svc, false)

	req := &artel_api.CompleteSetup_Request{
		WizardSessionToken: "wizard-token",
		Method: &artel_api.CompleteSetup_Request_Telegram{
			Telegram: &artel_api.TelegramCredentials{IdToken: "telegram-id-token"},
		},
	}

	resp, err := impl.CompleteSetup(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, "access-token", resp.Token)
	require.Equal(t, "wizard-token", svc.completeSetupCalledWith.wizardSessionToken)
	require.Equal(t, "telegram-id-token", svc.completeSetupCalledWith.telegramIdToken)
	require.Empty(t, svc.completeSetupCalledWith.email)
	require.Empty(t, svc.completeSetupCalledWith.password)
}

func TestSetupWizardImpl_CompleteSetup_PropagatesError(t *testing.T) {
	svcErr := errors.New("wizard session expired")
	svc := &fakeSetupWizardService{completeSetupErr: svcErr}
	impl := New(svc, false)

	req := &artel_api.CompleteSetup_Request{WizardSessionToken: "wizard-token"}

	resp, err := impl.CompleteSetup(context.Background(), req)
	require.Error(t, err)
	require.Nil(t, resp)
}
