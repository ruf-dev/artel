package transport

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/ruf-dev/artel/internal/middleware"
)

// NewGatewayMux builds the grpc-gateway mux shared by every *Impl.Gateway() method.
// CookieToMetadataAnnotator lets browser callers authenticate via the access_token cookie
// instead of an explicit Authorization header; CookieForwardResponseOption turns
// x-set-cookie-*/x-clear-auth-cookies gRPC response metadata into real Set-Cookie headers.
// Both are no-ops until a handler explicitly sets those keys, so this is safe to use
// uniformly regardless of whether this project ends up using cookie-based auth.
func NewGatewayMux(cookieSecure bool) *runtime.ServeMux {
	return runtime.NewServeMux(
		runtime.WithMetadata(middleware.CookieToMetadataAnnotator),
		runtime.WithForwardResponseOption(middleware.CookieForwardResponseOption(cookieSecure)),
	)
}
