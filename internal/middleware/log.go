package middleware

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func LogInterceptor() grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(
		func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
			fields := log.Debug().
				Str("method", info.FullMethod).
				Any("request", req)

			defer func() {
				fields.Msg("incoming GRPC request")
			}()

			resp, err = handler(ctx, req)

			fields = fields.
				Err(err).
				Any("response", resp)

			return resp, err
		})
}
