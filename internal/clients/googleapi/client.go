package googleapi

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/ruf-dev/artel/internal/domain"
)

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
