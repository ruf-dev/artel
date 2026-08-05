package auth_api

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/clients/couchdb"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service"
	"github.com/stretchr/testify/require"
)

// fakeAuthService is a minimal in-memory service.AuthService used to test authHandler.GetConfig
// without a real database — following the fakeCouchInstancesRepo convention in
// internal/service/v1/couchinstances/couchinstances_test.go (a plain struct with configurable
// result/err fields, not an embedded-interface trick). GetConfig never calls any of these
// methods, so they all panic if reached.
type fakeAuthService struct{}

func (f *fakeAuthService) Register(_ context.Context, _, _ string) (domain.User, error) {
	panic("not implemented")
}

func (f *fakeAuthService) Login(_ context.Context, _, _ string) (domain.Session, error) {
	panic("not implemented")
}

func (f *fakeAuthService) Logout(_ context.Context, _ string) error {
	panic("not implemented")
}

func (f *fakeAuthService) ValidateToken(_ context.Context, _ string) (domain.User, error) {
	panic("not implemented")
}

func (f *fakeAuthService) LoginViaTelegram(_ context.Context, _ string) (domain.Session, error) {
	panic("not implemented")
}

func (f *fakeAuthService) GetMe(_ context.Context, _ uuid.UUID) (domain.UserDetails, error) {
	panic("not implemented")
}

func (f *fakeAuthService) CheckIsAdmin(_ context.Context, _ uuid.UUID) error {
	panic("not implemented")
}

func (f *fakeAuthService) Refresh(_ context.Context, _ string) (domain.Session, error) {
	panic("not implemented")
}

func (f *fakeAuthService) EnsureNoAuthUser(_ context.Context) (domain.User, error) {
	panic("not implemented")
}

func (f *fakeAuthService) RegisterAdmin(_ context.Context, _, _ string) (domain.User, error) {
	panic("not implemented")
}

func (f *fakeAuthService) CreateUserUnchecked(_ context.Context, _, _ string) (domain.User, error) {
	panic("not implemented")
}

func (f *fakeAuthService) LoginOrRegisterAdminViaTelegram(_ context.Context, _ string) (domain.Session, error) {
	panic("not implemented")
}

func (f *fakeAuthService) ChangePassword(_ context.Context, _ uuid.UUID, _ string) error {
	panic("not implemented")
}

// fakeS3InstanceService is a minimal in-memory service.S3InstanceService used to test
// authHandler.GetConfig. Only HasS3Instances is configurable — GetConfig never calls the rest.
type fakeS3InstanceService struct {
	hasResult bool
	hasErr    error
}

func (f *fakeS3InstanceService) RegisterS3Instance(
	_ context.Context, _, _, _, _ string, _, _ bool,
) (string, error) {
	panic("not implemented")
}

func (f *fakeS3InstanceService) GetS3Instance(_ context.Context, _ string) (domain.S3Instance, error) {
	panic("not implemented")
}

func (f *fakeS3InstanceService) ListS3Instances(_ context.Context) ([]domain.S3Instance, error) {
	panic("not implemented")
}

func (f *fakeS3InstanceService) UpdateS3Instance(
	_ context.Context, _, _, _, _, _ string, _, _ bool,
) error {
	panic("not implemented")
}

func (f *fakeS3InstanceService) DeleteS3Instance(_ context.Context, _ string) error {
	panic("not implemented")
}

func (f *fakeS3InstanceService) TestS3Instance(_ context.Context, _ string) error {
	panic("not implemented")
}

func (f *fakeS3InstanceService) HasS3Instances(_ context.Context) (bool, error) {
	if f.hasErr != nil {
		return false, f.hasErr
	}

	return f.hasResult, nil
}

// fakeCouchInstanceService is a minimal in-memory service.CouchInstanceService used to test
// authHandler.GetConfig. Only HasCouchInstances is configurable — GetConfig never calls the rest.
type fakeCouchInstanceService struct {
	hasResult bool
	hasErr    error
}

func (f *fakeCouchInstanceService) RegisterCouchInstance(_ context.Context, _, _, _ string) (string, error) {
	panic("not implemented")
}

func (f *fakeCouchInstanceService) GetCouchInstance(_ context.Context, _ string) (domain.CouchInstance, error) {
	panic("not implemented")
}

func (f *fakeCouchInstanceService) ListCouchInstances(_ context.Context) ([]domain.CouchInstance, error) {
	panic("not implemented")
}

func (f *fakeCouchInstanceService) UpdateCouchInstance(_ context.Context, _, _, _, _ string) error {
	panic("not implemented")
}

func (f *fakeCouchInstanceService) DeleteCouchInstance(_ context.Context, _ string) error {
	panic("not implemented")
}

func (f *fakeCouchInstanceService) SetupCouchInstance(_ context.Context, _ string) error {
	panic("not implemented")
}

func (f *fakeCouchInstanceService) GetCouchInstanceStatus(
	_ context.Context, _ string,
) (couchdb.SetupStatus, error) {
	panic("not implemented")
}

func (f *fakeCouchInstanceService) HasCouchInstances(_ context.Context) (bool, error) {
	if f.hasErr != nil {
		return false, f.hasErr
	}

	return f.hasResult, nil
}

// fakeDockerHostService is a minimal in-memory service.DockerHostService used to test
// authHandler.GetConfig. Only HasDockerHosts is configurable — GetConfig never calls the rest.
type fakeDockerHostService struct {
	hasResult bool
	hasErr    error
}

func (f *fakeDockerHostService) RegisterDockerHost(_ context.Context, _, _, _, _ string) (string, error) {
	panic("not implemented")
}

func (f *fakeDockerHostService) GetDockerHost(_ context.Context, _ string) (domain.DockerHost, error) {
	panic("not implemented")
}

func (f *fakeDockerHostService) ListDockerHosts(_ context.Context) ([]domain.DockerHost, error) {
	panic("not implemented")
}

func (f *fakeDockerHostService) UpdateDockerHost(_ context.Context, _, _ string, _, _, _ *string) error {
	panic("not implemented")
}

func (f *fakeDockerHostService) DeleteDockerHost(_ context.Context, _ string) error {
	panic("not implemented")
}

func (f *fakeDockerHostService) HasDockerHosts(_ context.Context) (bool, error) {
	if f.hasErr != nil {
		return false, f.hasErr
	}

	return f.hasResult, nil
}

// fakeSetupWizardService is a minimal in-memory service.SetupWizardService used to test
// authHandler.GetConfig. Only CurrentStatus is configurable — GetConfig never calls the rest.
type fakeSetupWizardService struct {
	status domain.SystemSettings
	err    error
}

func (f *fakeSetupWizardService) CurrentStatus(_ context.Context) (domain.SystemSettings, error) {
	if f.err != nil {
		return domain.SystemSettings{}, f.err
	}

	return f.status, nil
}

func (f *fakeSetupWizardService) RegenerateSetupToken(_ context.Context) (string, error) {
	panic("not implemented")
}

func (f *fakeSetupWizardService) SubmitToken(_ context.Context, _ string) (string, error) {
	panic("not implemented")
}

func (f *fakeSetupWizardService) SelectAuthMethods(_ context.Context, _ string, _, _ bool) error {
	panic("not implemented")
}

func (f *fakeSetupWizardService) SelectRegistrationMode(_ context.Context, _ string, _ domain.RegistrationMode) error {
	panic("not implemented")
}

func (f *fakeSetupWizardService) CompleteSetup(
	_ context.Context, _, _, _, _ string,
) (domain.Session, error) {
	panic("not implemented")
}

var (
	_ service.AuthService          = (*fakeAuthService)(nil)
	_ service.S3InstanceService    = (*fakeS3InstanceService)(nil)
	_ service.CouchInstanceService = (*fakeCouchInstanceService)(nil)
	_ service.DockerHostService    = (*fakeDockerHostService)(nil)
	_ service.SetupWizardService   = (*fakeSetupWizardService)(nil)
)

func TestAuthHandler_GetConfig(t *testing.T) {
	testCases := []struct {
		name                string
		hasS3               bool
		hasCouch            bool
		noAuthEnabled       bool
		credsEncrypted      bool
		hasDockerHosts      bool
		setupCompleted      bool
		passwordAuthEnabled bool
		telegramAuthEnabled bool
		registrationMode    domain.RegistrationMode
	}{
		{
			name:                "all flags false",
			hasS3:               false,
			hasCouch:            false,
			noAuthEnabled:       false,
			credsEncrypted:      false,
			hasDockerHosts:      false,
			setupCompleted:      false,
			passwordAuthEnabled: false,
			telegramAuthEnabled: false,
			registrationMode:    domain.RegistrationModeAdminOnly,
		},
		{
			name:                "all flags true",
			hasS3:               true,
			hasCouch:            true,
			noAuthEnabled:       true,
			credsEncrypted:      true,
			hasDockerHosts:      true,
			setupCompleted:      true,
			passwordAuthEnabled: true,
			telegramAuthEnabled: true,
			registrationMode:    domain.RegistrationModeSelfRegister,
		},
		{
			name:                "mixed flags",
			hasS3:               true,
			hasCouch:            false,
			noAuthEnabled:       false,
			credsEncrypted:      true,
			hasDockerHosts:      false,
			setupCompleted:      true,
			passwordAuthEnabled: false,
			telegramAuthEnabled: true,
			registrationMode:    domain.RegistrationModeAdminOnly,
		},
		{
			name:                "workbench available only",
			hasS3:               false,
			hasCouch:            false,
			noAuthEnabled:       false,
			credsEncrypted:      false,
			hasDockerHosts:      true,
			setupCompleted:      false,
			passwordAuthEnabled: false,
			telegramAuthEnabled: false,
			registrationMode:    domain.RegistrationModeAdminOnly,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s3Svc := &fakeS3InstanceService{hasResult: tc.hasS3}
			couchSvc := &fakeCouchInstanceService{hasResult: tc.hasCouch}
			dockerHostSvc := &fakeDockerHostService{hasResult: tc.hasDockerHosts}
			settings := domain.SystemSettings{
				SetupCompleted:      tc.setupCompleted,
				PasswordAuthEnabled: tc.passwordAuthEnabled,
				TelegramAuthEnabled: tc.telegramAuthEnabled,
				RegistrationMode:    tc.registrationMode,
			}
			setupWizardSvc := &fakeSetupWizardService{status: settings}

			handler := &authHandler{
				authSvc:          &fakeAuthService{},
				s3InstanceSvc:    s3Svc,
				couchInstanceSvc: couchSvc,
				dockerHostSvc:    dockerHostSvc,
				setupWizardSvc:   setupWizardSvc,
				telegramClientID: "test-client-id",
				noAuthEnabled:    tc.noAuthEnabled,
				credsEncrypted:   tc.credsEncrypted,
			}

			resp, err := handler.GetConfig(context.Background(), nil)
			require.NoError(t, err)
			require.Equal(t, "test-client-id", resp.TelegramClientId)
			require.Equal(t, tc.hasS3, resp.IsS3Available)
			require.Equal(t, tc.hasCouch, resp.IsCouchAvailable)
			require.Equal(t, tc.noAuthEnabled, resp.NoAuthEnabled)
			require.Equal(t, tc.credsEncrypted, resp.CredsEncrypted)
			require.Equal(t, tc.hasDockerHosts, resp.IsWorkbenchAvailable)
			require.Equal(t, tc.setupCompleted, resp.SetupCompleted)
			require.Equal(t, tc.passwordAuthEnabled, resp.PasswordAuthEnabled)
			require.Equal(t, tc.telegramAuthEnabled, resp.TelegramAuthEnabled)
			require.Equal(t, tc.registrationMode == domain.RegistrationModeSelfRegister, resp.SelfRegisterEnabled)
		})
	}
}

func TestAuthHandler_GetConfig_PropagatesS3Error(t *testing.T) {
	s3Err := errors.New("s3 backend unavailable")
	s3Svc := &fakeS3InstanceService{hasErr: s3Err}
	couchSvc := &fakeCouchInstanceService{}

	handler := &authHandler{
		authSvc:          &fakeAuthService{},
		s3InstanceSvc:    s3Svc,
		couchInstanceSvc: couchSvc,
	}

	resp, err := handler.GetConfig(context.Background(), nil)
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestAuthHandler_GetConfig_PropagatesCouchError(t *testing.T) {
	couchErr := errors.New("couch backend unavailable")
	s3Svc := &fakeS3InstanceService{}
	couchSvc := &fakeCouchInstanceService{hasErr: couchErr}

	handler := &authHandler{
		authSvc:          &fakeAuthService{},
		s3InstanceSvc:    s3Svc,
		couchInstanceSvc: couchSvc,
	}

	resp, err := handler.GetConfig(context.Background(), nil)
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestAuthHandler_GetConfig_PropagatesDockerHostsError(t *testing.T) {
	dockerHostsErr := errors.New("docker hosts backend unavailable")
	s3Svc := &fakeS3InstanceService{}
	couchSvc := &fakeCouchInstanceService{}
	dockerHostSvc := &fakeDockerHostService{hasErr: dockerHostsErr}

	handler := &authHandler{
		authSvc:          &fakeAuthService{},
		s3InstanceSvc:    s3Svc,
		couchInstanceSvc: couchSvc,
		dockerHostSvc:    dockerHostSvc,
	}

	resp, err := handler.GetConfig(context.Background(), nil)
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestAuthHandler_GetConfig_PropagatesSetupWizardError(t *testing.T) {
	setupErr := errors.New("setup wizard status unavailable")
	s3Svc := &fakeS3InstanceService{}
	couchSvc := &fakeCouchInstanceService{}
	dockerHostSvc := &fakeDockerHostService{}
	setupWizardSvc := &fakeSetupWizardService{err: setupErr}

	handler := &authHandler{
		authSvc:          &fakeAuthService{},
		s3InstanceSvc:    s3Svc,
		couchInstanceSvc: couchSvc,
		dockerHostSvc:    dockerHostSvc,
		setupWizardSvc:   setupWizardSvc,
	}

	resp, err := handler.GetConfig(context.Background(), nil)
	require.Error(t, err)
	require.Nil(t, resp)
}
