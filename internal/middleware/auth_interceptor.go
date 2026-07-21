package middleware

import (
	"context"

	"go.redsock.ru/rerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/service"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

const authHeader = "authorization"

type authMiddleware struct {
	ignoredPaths        map[string]struct{}
	isDebugEnabled      bool
	noAuthEnabled       bool
	authService         service.AuthService
	subscriptionService service.SubscriptionService
}

func (am *authMiddleware) isIgnored(path string) bool {
	_, ok := am.ignoredPaths[path]

	return ok
}

type authOption func(*authMiddleware)

func WithIgnoredPathAuthOption(p ...string) authOption {
	return func(am *authMiddleware) {
		for _, path := range p {
			am.ignoredPaths[path] = struct{}{}
		}
	}
}

func WithDebug(b bool) authOption {
	return func(am *authMiddleware) {
		am.isDebugEnabled = b
	}
}

// WithNoAuth makes the interceptor a pass-through for every RPC, injecting the fixed
// local-dev user (see auth.Service.EnsureNoAuthUser) instead of validating any token.
// Intended for local development only.
func WithNoAuth(b bool) authOption {
	return func(am *authMiddleware) {
		am.noAuthEnabled = b
	}
}

func GrpcAuthInterceptor(srv service.Service, opts ...authOption) grpc.ServerOption {
	ac := &authMiddleware{
		ignoredPaths:        make(map[string]struct{}),
		authService:         srv.AuthService(),
		subscriptionService: srv.SubscriptionService(),
	}

	for _, opt := range opts {
		opt(ac)
	}

	return grpc.ChainUnaryInterceptor(
		func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
			if ac.noAuthEnabled {
				ctxWithUser, err := ac.authWithNoAuth(ctx)
				if err != nil {
					return nil, err
				}

				return handler(ctxWithUser, req)
			}

			if ac.isIgnored(info.FullMethod) {
				return handler(ctx, req)
			}

			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return nil, status.Error(codes.FailedPrecondition, user_errors.NoMetadataInContext.Error())
			}

			ctxWithUser, err := ac.authWithSession(ctx, md)
			if err == nil {
				return handler(ctxWithUser, req)
			}

			if ac.isDebugEnabled {
				debugErr := ac.authWithDebugHeaders(ctx, md)
				if debugErr == nil {
					return handler(ctx, req)
				}
			}

			return nil, err
		})
}

func (am *authMiddleware) authWithSession(ctx context.Context, md metadata.MD) (context.Context, error) {
	auth := md.Get(authHeader)
	if len(auth) == 0 {
		return nil, status.Error(codes.Unauthenticated, user_errors.NoAuthHeader.Error())
	}

	user, err := am.authService.ValidateToken(ctx, auth[0])
	if err != nil {
		return nil, rerrors.Wrap(err)
	}

	err = am.subscriptionService.CheckActive(ctx, user.Uuid)
	if err != nil {
		return nil, rerrors.Wrap(err)
	}

	uc := user_context.UserContext{
		UserUuid: user.Uuid,
		UserName: user.Username,
		Roles:    nil,
	}

	ctxWithUser := user_context.WithUserContext(ctx, uc)

	return ctxWithUser, nil
}

// authWithNoAuth injects the fixed local-dev user without validating any token or
// checking subscription status — the dev user has no real subscription to be active.
func (am *authMiddleware) authWithNoAuth(ctx context.Context) (context.Context, error) {
	user, err := am.authService.ValidateToken(ctx, "")
	if err != nil {
		return nil, rerrors.Wrap(err)
	}

	uc := user_context.UserContext{
		UserUuid: user.Uuid,
		UserName: user.Username,
		Roles:    nil,
	}

	return user_context.WithUserContext(ctx, uc), nil
}

func (am *authMiddleware) authWithDebugHeaders(_ context.Context, _ metadata.MD) (err error) {
	return status.Error(codes.Unimplemented, user_errors.DebugNotSupported.Error())
}
