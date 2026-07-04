package middleware

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func PanicInterceptor() grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(
		func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
			defer func() {
				r := recover()
				if r != nil {
					panicErr, ok := r.(error)
					if ok {
						log.Err(panicErr).Msg("panic in grpc handler")
					}
				}
			}()

			return handler(ctx, req)
		},
	)
}
