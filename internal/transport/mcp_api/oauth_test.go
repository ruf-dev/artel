package mcp_api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ruf-dev/artel/internal/middleware"
)

func TestResolveOAuthSessionToken_BodyTokenBypassesCookieAndCSRF(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/oauth/vaults", nil)

	rec := httptest.NewRecorder()
	token, ok := resolveOAuthSessionToken(rec, req, "explicit-body-token")

	require.True(t, ok)
	require.Equal(t, "explicit-body-token", token)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestResolveOAuthSessionToken_CookieBranch(t *testing.T) {
	const csrfValue = "csrf-abc"

	cases := []struct {
		name         string
		accessCookie string
		csrfCookie   string
		csrfHeader   string
		wantOk       bool
		wantToken    string
		wantStatus   int
	}{
		{
			name:         "valid cookie and matching csrf",
			accessCookie: "cookie-token",
			csrfCookie:   csrfValue,
			csrfHeader:   csrfValue,
			wantOk:       true,
			wantToken:    "cookie-token",
			wantStatus:   http.StatusOK,
		},
		{
			name:         "no access cookie",
			accessCookie: "",
			csrfCookie:   csrfValue,
			csrfHeader:   csrfValue,
			wantOk:       false,
			wantStatus:   http.StatusUnauthorized,
		},
		{
			name:         "csrf header missing",
			accessCookie: "cookie-token",
			csrfCookie:   csrfValue,
			csrfHeader:   "",
			wantOk:       false,
			wantStatus:   http.StatusForbidden,
		},
		{
			name:         "csrf header mismatch",
			accessCookie: "cookie-token",
			csrfCookie:   csrfValue,
			csrfHeader:   "wrong",
			wantOk:       false,
			wantStatus:   http.StatusForbidden,
		},
		{
			name:         "csrf cookie missing",
			accessCookie: "cookie-token",
			csrfCookie:   "",
			csrfHeader:   csrfValue,
			wantOk:       false,
			wantStatus:   http.StatusForbidden,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/oauth/vaults", nil)
			if tc.accessCookie != "" {
				req.AddCookie(&http.Cookie{Name: middleware.AccessTokenCookieName, Value: tc.accessCookie})
			}
			if tc.csrfCookie != "" {
				req.AddCookie(&http.Cookie{Name: middleware.CSRFCookieName, Value: tc.csrfCookie})
			}
			if tc.csrfHeader != "" {
				req.Header.Set(csrfHeaderName, tc.csrfHeader)
			}

			rec := httptest.NewRecorder()
			token, ok := resolveOAuthSessionToken(rec, req, "")

			require.Equal(t, tc.wantOk, ok)
			require.Equal(t, tc.wantToken, token)
			require.Equal(t, tc.wantStatus, rec.Code)
		})
	}
}
