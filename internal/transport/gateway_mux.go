package transport

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/ruf-dev/artel/internal/middleware"
)

// NewGatewayMux builds the grpc-gateway mux shared by every *Impl.Gateway() method, wiring the
// cookie<->metadata plumbing (see internal/middleware/cookie_annotator.go,
// cookie_response.go): CookieToMetadataAnnotator lets browser callers authenticate via the
// access_token cookie instead of an explicit header, and CookieForwardResponseOption turns the
// x-set-cookie-*/x-clear-auth-cookies metadata auth_api's handlers set into real Set-Cookie
// response headers. It's a no-op for every RPC except auth_api's Login/Refresh/Logout, so it's
// safe to use uniformly across all gateway-registering implementations.
//
// cookieSecure controls the Secure attribute on cookies CookieForwardResponseOption sets (false
// for plain-HTTP local dev, true otherwise — config.EnvironmentConfig.CookieSecure).
func NewGatewayMux(cookieSecure bool) *runtime.ServeMux {
	return runtime.NewServeMux(
		runtime.WithMetadata(middleware.CookieToMetadataAnnotator),
		runtime.WithForwardResponseOption(middleware.CookieForwardResponseOption(cookieSecure)),
	)
}
