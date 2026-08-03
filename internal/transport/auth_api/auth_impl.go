package auth_api

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware"
	"github.com/ruf-dev/artel/internal/service"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"github.com/ruf-dev/artel/internal/transport"
	"go.redsock.ru/rerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// authHandler implements the AuthAPIServer proto interface.
// Kept separate from AuthImpl because the proto defines a method named "Register"
// which collides with the transport.GrpcImpl interface method Register(grpc.ServiceRegistrar).
type authHandler struct {
	artel_api.UnimplementedAuthAPIServer
	authSvc          service.AuthService
	s3InstanceSvc    service.S3InstanceService
	couchInstanceSvc service.CouchInstanceService
	dockerHostSvc    service.DockerHostService
	telegramClientID string
	noAuthEnabled    bool
	credsEncrypted   bool
}

func (h *authHandler) Register(
	ctx context.Context,
	req *artel_api.Register_Request,
) (*artel_api.Register_Response, error) {
	user, err := h.authSvc.Register(ctx, req.Email, req.Password)
	if err != nil {
		return nil, rerrors.Wrap(err, "register")
	}

	resp := &artel_api.Register_Response{
		Id:    user.Uuid.String(),
		Email: user.Email,
	}

	return resp, nil
}

func (h *authHandler) Login(ctx context.Context, req *artel_api.Login_Request) (*artel_api.Login_Response, error) {
	switch {
	case req.GetPassword() != nil:
		passwordCreds := req.GetPassword()

		session, err := h.authSvc.Login(ctx, passwordCreds.GetEmail(), passwordCreds.GetPassword())
		if err != nil {
			return nil, rerrors.Wrap(err, "login")
		}

		setAuthCookieHeaders(ctx, session)

		resp := &artel_api.Login_Response{
			Token:            session.Token,
			ExpiresAt:        timestamppb.New(session.ExpiresAt),
			RefreshToken:     session.RefreshToken,
			RefreshExpiresAt: timestamppb.New(session.RefreshExpiresAt),
		}

		return resp, nil
	case req.GetTelegram() != nil:
		telegramCreds := req.GetTelegram()

		session, err := h.authSvc.LoginViaTelegram(ctx, telegramCreds.GetIdToken())
		if err != nil {
			return nil, rerrors.Wrap(err, "login via telegram")
		}

		setAuthCookieHeaders(ctx, session)

		resp := &artel_api.Login_Response{
			Token:            session.Token,
			ExpiresAt:        timestamppb.New(session.ExpiresAt),
			RefreshToken:     session.RefreshToken,
			RefreshExpiresAt: timestamppb.New(session.RefreshExpiresAt),
		}

		return resp, nil
	default:
		return nil, user_errors.UnsupportedLoginMethod
	}
}

func (h *authHandler) Refresh(
	ctx context.Context,
	req *artel_api.Refresh_Request,
) (*artel_api.Refresh_Response, error) {
	refreshToken := refreshTokenFromRequest(ctx, req)

	session, err := h.authSvc.Refresh(ctx, refreshToken)
	if err != nil {
		return nil, rerrors.Wrap(err, "error refreshing session")
	}

	setAuthCookieHeaders(ctx, session)

	resp := &artel_api.Refresh_Response{
		Token:            session.Token,
		ExpiresAt:        timestamppb.New(session.ExpiresAt),
		RefreshToken:     session.RefreshToken,
		RefreshExpiresAt: timestamppb.New(session.RefreshExpiresAt),
	}

	return resp, nil
}

// refreshTokenFromRequest returns the request body's refresh_token field, falling back to the
// refresh_token cookie (browsers never send it in the body — see CookieForwardResponseOption)
// when the body field is empty.
func refreshTokenFromRequest(ctx context.Context, req *artel_api.Refresh_Request) string {
	token := req.GetRefreshToken()
	if token != "" {
		return token
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	rawCookieHeader := ""
	cookieHeaders := md.Get(middleware.GatewayCookieMetadataKey)
	if len(cookieHeaders) > 0 {
		rawCookieHeader = cookieHeaders[0]
	}
	if rawCookieHeader == "" {
		return ""
	}

	cookieToken, err := middleware.CookieValueFromRawHeader(rawCookieHeader, middleware.RefreshTokenCookieName)
	if err != nil {
		return ""
	}

	return cookieToken
}

// setAuthCookieHeaders signals CookieForwardResponseOption (registered on the shared gateway
// mux, see transport.NewGatewayMux) to set the browser's httpOnly auth cookies from this
// session. It's a best-effort side channel on top of the unchanged proto response — a failure
// here must not fail Login/Refresh for native gRPC/CLI callers, who never look at cookies.
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

func (h *authHandler) GetConfig(
	ctx context.Context,
	_ *artel_api.GetConfig_Request,
) (*artel_api.GetConfig_Response, error) {
	hasS3, err := h.s3InstanceSvc.HasS3Instances(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "check s3 instances availability")
	}

	hasCouch, err := h.couchInstanceSvc.HasCouchInstances(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "check couch instances availability")
	}

	hasDockerHosts, err := h.dockerHostSvc.HasDockerHosts(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "check docker hosts availability")
	}

	return &artel_api.GetConfig_Response{
		TelegramClientId:     h.telegramClientID,
		IsS3Available:        hasS3,
		NoAuthEnabled:        h.noAuthEnabled,
		CredsEncrypted:       h.credsEncrypted,
		IsCouchAvailable:     hasCouch,
		IsWorkbenchAvailable: hasDockerHosts,
	}, nil
}

func (h *authHandler) Logout(ctx context.Context, req *artel_api.Logout_Request) (*artel_api.Logout_Response, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	tokens := md.Get("authorization")

	token := ""
	if len(tokens) > 0 {
		token = tokens[0]
	}

	err := h.authSvc.Logout(ctx, token)
	if err != nil {
		return nil, rerrors.Wrap(err, "logout")
	}

	clearAuthCookieHeaders(ctx)

	return &artel_api.Logout_Response{}, nil
}

// clearAuthCookieHeaders signals CookieForwardResponseOption to expire all three auth cookies
// on the browser. Best-effort, same rationale as setAuthCookieHeaders.
func clearAuthCookieHeaders(ctx context.Context) {
	md := metadata.Pairs(middleware.ClearAuthCookiesKey, middleware.ClearAuthCookiesValue)

	err := grpc.SetHeader(ctx, md)
	if err != nil {
		log.Error().Err(err).Msg("error setting clear-auth-cookies header")
	}
}

func (h *authHandler) GetMe(ctx context.Context, _ *artel_api.GetMe_Request) (*artel_api.GetMe_Response, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	tokens := md.Get("authorization")

	token := ""
	if len(tokens) > 0 {
		token = tokens[0]
	}

	validatedUser, err := h.authSvc.ValidateToken(ctx, token)
	if err != nil {
		return nil, rerrors.Wrap(err, "validate token")
	}

	details, err := h.authSvc.GetMe(ctx, validatedUser.Uuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "get me")
	}

	resp := &artel_api.GetMe_Response{
		Id:       details.Uuid.String(),
		Username: details.Username,
		Email:    details.Email,
		PhotoUrl: details.PhotoUrl,
		Permissions: &artel_api.Permissions{
			IsAdministrator: details.Permissions.IsAdministrator,
			HasEmails:       details.EffectiveSubscription.Features.Emails,
			HasTaskTrackers: details.EffectiveSubscription.Features.TaskTrackers,
			HasNotes:        details.EffectiveSubscription.Features.Notes,
		},
	}

	return resp, nil
}

// AuthImpl satisfies transport.GrpcImpl and transport.GrpcWithGateway.
type AuthImpl struct {
	handler      *authHandler
	cookieSecure bool
}

func NewAuthImpl(
	authSvc service.AuthService, telegramClientID string, s3InstanceSvc service.S3InstanceService,
	couchInstanceSvc service.CouchInstanceService, noAuthEnabled bool,
	credsEncrypted bool, dockerHostSvc service.DockerHostService, cookieSecure bool,
) *AuthImpl {
	return &AuthImpl{
		handler: &authHandler{
			authSvc:          authSvc,
			telegramClientID: telegramClientID,
			s3InstanceSvc:    s3InstanceSvc,
			couchInstanceSvc: couchInstanceSvc,
			dockerHostSvc:    dockerHostSvc,
			noAuthEnabled:    noAuthEnabled,
			credsEncrypted:   credsEncrypted,
		},
		cookieSecure: cookieSecure,
	}
}

func (a *AuthImpl) Register(srv grpc.ServiceRegistrar) {
	artel_api.RegisterAuthAPIServer(srv, a.handler)
}

func (a *AuthImpl) Gateway(ctx context.Context, endpoint string, opts ...grpc.DialOption) (string, http.Handler) {
	gwMux := transport.NewGatewayMux(a.cookieSecure)

	err := artel_api.RegisterAuthAPIHandlerFromEndpoint(ctx, gwMux, endpoint, opts)
	if err != nil {
		log.Error().Err(err).Msg("error registering auth grpc-gateway handler")
	}

	return "/api/auth/", gwMux
}
