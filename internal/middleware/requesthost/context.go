// Package requesthost threads the public-facing host of an inbound HTTP request from the
// transport layer into the service layer, mirroring the user_context package's pattern. Used
// by the MCP handler (a plain net/http handler, not grpc-gateway) so builtin tract tools can
// build absolute webhook_url values without the service layer depending on *http.Request.
package requesthost

import "context"

// contextKey is a type for context keys to avoid collisions.
type contextKey string

const hostContextKey contextKey = "request_host"

// WithHost stores the request's public host (scheme-less, e.g. "artel.example.com") in ctx.
func WithHost(ctx context.Context, host string) context.Context {
	return context.WithValue(ctx, hostContextKey, host)
}

// FromContext retrieves the public host stored by WithHost.
func FromContext(ctx context.Context) (string, bool) {
	host, ok := ctx.Value(hostContextKey).(string)

	return host, ok
}
