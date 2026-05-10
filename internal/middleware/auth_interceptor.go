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
)

const authHeader = "authorization"

type authMiddleware struct {
	ignoredPaths   map[string]struct{}
	isDebugEnabled bool
	authService    service.AuthService
	userService    service.UserService
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

func GrpcAuthInterceptor(srv service.Service, opts ...authOption) grpc.ServerOption {
	ac := &authMiddleware{
		ignoredPaths: make(map[string]struct{}),
		authService:  srv.AuthService(),
		userService:  srv.UserService(),
	}

	for _, opt := range opts {
		opt(ac)
	}

	return grpc.ChainUnaryInterceptor(
		func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
			if ac.isIgnored(info.FullMethod) {
				return handler(ctx, req)
			}

			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				err := rerrors.New("error getting metadata from context")
				return nil, status.Error(codes.FailedPrecondition, err.Error())
			}

			ctxWithUser, err := ac.authWithSession(ctx, md)
			if err == nil {
				return handler(ctxWithUser, req)
			}

			if ac.isDebugEnabled {
				ctxWithDebug, debugErr := ac.authWithDebugHeaders(ctx, md)
				if debugErr == nil {
					return handler(ctxWithDebug, req)
				}
			}

			return nil, err
		})
}

func (am *authMiddleware) authWithSession(ctx context.Context, md metadata.MD) (context.Context, error) {
	auth := md.Get(authHeader)
	if len(auth) == 0 {
		err := rerrors.New("error getting auth header")
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	userUuid, err := am.authService.AuthWithToken(ctx, auth[0])
	if err != nil {
		wrapped := rerrors.Wrap(err, "auth with token")
		return nil, status.Error(codes.Unauthenticated, wrapped.Error())
	}

	uc := user_context.UserContext{
		UserUuid: userUuid,
	}

	user, err := am.userService.GetUser(ctx, userUuid.String())
	if err != nil {
		wrapped := rerrors.Wrap(err, "get user")
		return nil, status.Error(codes.Internal, wrapped.Error())
	}

	uc.Roles = user.Roles

	ctxWithUser := user_context.WithUserContext(ctx, uc)
	return ctxWithUser, nil
}

func (am *authMiddleware) authWithDebugHeaders(ctx context.Context, md metadata.MD) (context.Context, error) {
	err := rerrors.New("debug not supported")
	return nil, status.Error(codes.Unimplemented, err.Error())
}
