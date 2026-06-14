package googleapi

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/oauth2"

	"github.com/ruf-dev/artel/internal/domain"
)

var baseHTTPClient = &http.Client{Transport: newLoggingTransport()}

type loggingTransport struct {
	next http.RoundTripper
}

func newLoggingTransport() *loggingTransport {
	return &loggingTransport{next: otelhttp.NewTransport(http.DefaultTransport)}
}

func (t *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()
	resp, err := t.next.RoundTrip(req)
	elapsed := time.Since(start)
	if err != nil {
		log.Ctx(req.Context()).Error().Err(err).Str("method", req.Method).Str("url", req.URL.String()).Dur("dur", elapsed).Msg("google api returned error")
		return resp, err
	}
	log.Ctx(req.Context()).Debug().Str("method", req.Method).Str("url", req.URL.String()).Int("status", resp.StatusCode).Dur("dur", elapsed).Msg("google api request")
	return resp, nil
}

// Client is an authenticated Google API client backed by stored OAuth2 credentials.
type Client struct {
	httpClient *http.Client
}

// New builds a Client from stored credentials. The caller is responsible for ensuring
// the token is fresh (call GetGoogleClient on the service, which handles refresh).
func New(ctx context.Context, creds domain.GoogleOAuthCredentials, cfg *oauth2.Config) *Client {
	token := &oauth2.Token{
		AccessToken:  creds.AccessToken,
		RefreshToken: creds.RefreshToken,
		Expiry:       creds.Expiry,
		TokenType:    creds.TokenType,
	}
	httpClient := cfg.Client(ctx, token)
	return &Client{httpClient: httpClient}
}

// HTTPClient returns the underlying authenticated HTTP client for making Google API calls.
func (c *Client) HTTPClient() *http.Client {
	return c.httpClient
}
