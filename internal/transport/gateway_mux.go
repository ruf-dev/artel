package transport

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/ruf-dev/artel/internal/middleware"
)

// NewGatewayMux builds the grpc-gateway mux shared by every *Impl.Gateway() method.
// CookieToMetadataAnnotator lets browser callers authenticate via the access_token cookie
// instead of an explicit Authorization header; RequestSchemeAnnotator records whether the
// originating request was secure (X-Forwarded-Proto or r.TLS), auto-detected per request rather
// than from a config value; CookieForwardResponseOption turns x-set-cookie-*/x-clear-auth-cookies
// gRPC response metadata into real Set-Cookie headers, using that per-request scheme to decide
// each cookie's Secure attribute. All three are no-ops until a handler explicitly sets those
// keys, so this is safe to use uniformly regardless of whether this project ends up using
// cookie-based auth.
func NewGatewayMux() *runtime.ServeMux {
	return runtime.NewServeMux(
		runtime.WithMetadata(middleware.CookieToMetadataAnnotator),
		runtime.WithMetadata(middleware.RequestSchemeAnnotator),
		runtime.WithForwardResponseOption(middleware.CookieForwardResponseOption()),
	)
}
