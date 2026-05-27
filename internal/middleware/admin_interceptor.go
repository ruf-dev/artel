package middleware

import (
	"context"

	"go.redsock.ru/rerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/service"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

func GrpcAdminInterceptor(authSvc service.AuthService, adminPaths ...string) grpc.ServerOption {
	pathSet := make(map[string]struct{}, len(adminPaths))
	for _, p := range adminPaths {
		pathSet[p] = struct{}{}
	}

	return grpc.ChainUnaryInterceptor(
		func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
			if _, ok := pathSet[info.FullMethod]; !ok {
				return handler(ctx, req)
			}

			uc, ok := user_context.GetUserContext(ctx)
			if !ok {
				return nil, status.Error(codes.PermissionDenied, user_errors.NoMetadataInContext.Error())
			}

			err := authSvc.CheckIsAdmin(ctx, uc.UserUuid)
			if err != nil {
				return nil, status.Error(codes.PermissionDenied, rerrors.Wrap(err, "admin check").Error())
			}

			return handler(ctx, req)
		})
}
