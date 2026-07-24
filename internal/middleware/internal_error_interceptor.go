package middleware

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const internalErrorMessage = "internal error"

// InternalErrorInterceptor prevents raw/unclassified server errors (e.g. underlying driver or
// crypto errors wrapped via rerrors.Wrap without an explicit codes.X) from reaching grpc clients.
// Any error whose resolved status code is codes.Internal or codes.Unknown is logged in full
// server-side and replaced with a fixed, generic message. Errors that already carry an explicit,
// intentional code (user_errors sentinels, status.Error calls in other interceptors) pass through
// unchanged.
func InternalErrorInterceptor() grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(RedactInternalErrors)
}

// RedactInternalErrors is exported (rather than only reachable via InternalErrorInterceptor) so it
// can be exercised directly in unit tests without spinning up a real grpc.Server.
func RedactInternalErrors(
	ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (any, error) {
	resp, err := handler(ctx, req)
	if err == nil {
		return resp, nil
	}

	st, _ := status.FromError(err)

	if isClientSafeCode(st.Code()) {
		return resp, err
	}

	log.Ctx(ctx).Error().
		Err(err).
		Str("method", info.FullMethod).
		Msg("redacting unclassified internal error before returning it to grpc client")

	return nil, status.Error(codes.Internal, internalErrorMessage)
}

func isClientSafeCode(code codes.Code) bool {
	return code != codes.Internal && code != codes.Unknown
}
