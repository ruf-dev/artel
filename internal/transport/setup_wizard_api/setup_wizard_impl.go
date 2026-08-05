package setup_wizard_api

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware"
	"github.com/ruf-dev/artel/internal/service"
	"github.com/ruf-dev/artel/internal/transport"
	"go.redsock.ru/rerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// SetupWizardImpl implements the SetupWizardAPI proto interface. Every RPC here is
// unauthenticated at the interceptor level (see WithIgnoredPathAuthOption in
// internal/app/custom.go) — CompleteSetup/SelectAuthMethods/SelectRegistrationMode instead gate
// on the wizardSessionToken minted by SubmitToken, checked inside SetupWizardService itself.
type SetupWizardImpl struct {
	artel_api.UnimplementedSetupWizardAPIServer
	svc          service.SetupWizardService
	cookieSecure bool
}

func New(svc service.SetupWizardService, cookieSecure bool) *SetupWizardImpl {
	return &SetupWizardImpl{svc: svc, cookieSecure: cookieSecure}
}

func (s *SetupWizardImpl) Register(srv grpc.ServiceRegistrar) {
	artel_api.RegisterSetupWizardAPIServer(srv, s)
}

func (s *SetupWizardImpl) Gateway(ctx context.Context, endpoint string, opts ...grpc.DialOption) (string, http.Handler) {
	gwMux := transport.NewGatewayMux(s.cookieSecure)

	err := artel_api.RegisterSetupWizardAPIHandlerFromEndpoint(ctx, gwMux, endpoint, opts)
	if err != nil {
		log.Error().Err(err).Msg("error registering setup wizard grpc-gateway handler")
	}

	return "/api/setup_wizard/", gwMux
}

func (s *SetupWizardImpl) GetStatus(
	ctx context.Context,
	_ *artel_api.GetStatus_Request,
) (*artel_api.GetStatus_Response, error) {
	settings, err := s.svc.CurrentStatus(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "error getting setup wizard status")
	}

	resp := &artel_api.GetStatus_Response{
		SetupCompleted: settings.SetupCompleted,
		TokenPending:   settings.SetupTokenHash != "",
	}

	return resp, nil
}

func (s *SetupWizardImpl) SubmitToken(
	ctx context.Context,
	req *artel_api.SubmitToken_Request,
) (*artel_api.SubmitToken_Response, error) {
	wizardSessionToken, err := s.svc.SubmitToken(ctx, req.GetToken())
	if err != nil {
		return nil, rerrors.Wrap(err, "error submitting setup token")
	}

	resp := &artel_api.SubmitToken_Response{WizardSessionToken: wizardSessionToken}

	return resp, nil
}

func (s *SetupWizardImpl) SelectAuthMethods(
	ctx context.Context,
	req *artel_api.SelectAuthMethods_Request,
) (*artel_api.SelectAuthMethods_Response, error) {
	err := s.svc.SelectAuthMethods(ctx, req.GetWizardSessionToken(), req.GetPasswordEnabled(), req.GetTelegramEnabled())
	if err != nil {
		return nil, rerrors.Wrap(err, "error selecting auth methods")
	}

	return &artel_api.SelectAuthMethods_Response{}, nil
}

func (s *SetupWizardImpl) SelectRegistrationMode(
	ctx context.Context,
	req *artel_api.SelectRegistrationMode_Request,
) (*artel_api.SelectRegistrationMode_Response, error) {
	mode := registrationModeFromProto(req.GetMode())

	err := s.svc.SelectRegistrationMode(ctx, req.GetWizardSessionToken(), mode)
	if err != nil {
		return nil, rerrors.Wrap(err, "error selecting registration mode")
	}

	return &artel_api.SelectRegistrationMode_Response{}, nil
}

func registrationModeFromProto(mode artel_api.RegistrationMode) domain.RegistrationMode {
	if mode == artel_api.RegistrationMode_SELF_REGISTER {
		return domain.RegistrationModeSelfRegister
	}

	return domain.RegistrationModeAdminOnly
}

func (s *SetupWizardImpl) CompleteSetup(
	ctx context.Context,
	req *artel_api.CompleteSetup_Request,
) (*artel_api.CompleteSetup_Response, error) {
	email := ""
	password := ""
	telegramIdToken := ""

	switch method := req.GetMethod().(type) {
	case *artel_api.CompleteSetup_Request_Password:
		email = method.Password.GetEmail()
		password = method.Password.GetPassword()
	case *artel_api.CompleteSetup_Request_Telegram:
		telegramIdToken = method.Telegram.GetIdToken()
	}

	session, err := s.svc.CompleteSetup(ctx, req.GetWizardSessionToken(), email, password, telegramIdToken)
	if err != nil {
		return nil, rerrors.Wrap(err, "error completing setup")
	}

	setAuthCookieHeaders(ctx, session)

	resp := &artel_api.CompleteSetup_Response{
		Token:            session.Token,
		ExpiresAt:        timestamppb.New(session.ExpiresAt),
		RefreshToken:     session.RefreshToken,
		RefreshExpiresAt: timestamppb.New(session.RefreshExpiresAt),
	}

	return resp, nil
}

// setAuthCookieHeaders signals CookieForwardResponseOption (registered on the shared gateway
// mux, see transport.NewGatewayMux) to set the browser's httpOnly auth cookies from this
// session — the admin created by CompleteSetup should end up logged in exactly like a normal
// login. Mirrors auth_api.setAuthHandler.setAuthCookieHeaders; best-effort, a failure here must
// not fail CompleteSetup for native gRPC/CLI callers, who never look at cookies.
func setAuthCookieHeaders(ctx context.Context, session domain.Session) {
	md := metadata.Pairs(
		middleware.SetCookieAccessTokenKey, session.Token,
		middleware.SetCookieAccessTokenExpiryKey, session.ExpiresAt.UTC().Format(time.RFC3339),
		middleware.SetCookieRefreshTokenKey, session.RefreshToken,
		middleware.SetCookieRefreshTokenExpiryKey, session.RefreshExpiresAt.UTC().Format(time.RFC3339),
	)

	err := grpc.SetHeader(ctx, md)
	if err != nil {
		log.Error().Err(err).Msg("error setting auth cookie headers")
	}
}
