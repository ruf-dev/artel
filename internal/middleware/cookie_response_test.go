package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/metadata"

	"github.com/ruf-dev/artel/internal/middleware"
)

func serverMetadataContext(headerPairs ...string) (*httptest.ResponseRecorder, func() error) {
	rec := httptest.NewRecorder()

	sm := runtime.ServerMetadata{HeaderMD: metadata.Pairs(headerPairs...)}
	ctx := runtime.NewServerMetadataContext(context.Background(), sm)

	opt := middleware.CookieForwardResponseOption(true)

	run := func() error {
		return opt(ctx, rec, nil)
	}

	return rec, run
}

func TestCookieForwardResponseOption_SetsAccessAndRefreshCookies(t *testing.T) {
	expiry := time.Now().Add(time.Hour).UTC().Format(time.RFC3339)

	rec, run := serverMetadataContext(
		middleware.SetCookieAccessTokenKey, "access-value",
		middleware.SetCookieAccessTokenExpiryKey, expiry,
		middleware.SetCookieRefreshTokenKey, "refresh-value",
		middleware.SetCookieRefreshTokenExpiryKey, expiry,
	)

	err := run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cookies := rec.Result().Cookies()

	byName := make(map[string]*http.Cookie, len(cookies))
	for _, c := range cookies {
		byName[c.Name] = c
	}

	access, ok := byName[middleware.AccessTokenCookieName]
	if !ok {
		t.Fatal("expected access_token cookie to be set")
	}
	if access.Value != "access-value" {
		t.Fatalf("expected access_token value 'access-value', got %q", access.Value)
	}
	if !access.HttpOnly {
		t.Fatal("expected access_token cookie to be HttpOnly")
	}
	if !access.Secure {
		t.Fatal("expected access_token cookie to be Secure (secure=true was passed)")
	}
	if access.SameSite != http.SameSiteLaxMode {
		t.Fatalf("expected SameSite=Lax, got %v", access.SameSite)
	}
	if access.Path != middleware.CookiePath {
		t.Fatalf("expected path %q, got %q", middleware.CookiePath, access.Path)
	}

	refresh, ok := byName[middleware.RefreshTokenCookieName]
	if !ok {
		t.Fatal("expected refresh_token cookie to be set")
	}
	if refresh.Value != "refresh-value" {
		t.Fatalf("expected refresh_token value 'refresh-value', got %q", refresh.Value)
	}
	if !refresh.HttpOnly {
		t.Fatal("expected refresh_token cookie to be HttpOnly")
	}

	csrf, ok := byName[middleware.CSRFCookieName]
	if !ok {
		t.Fatal("expected csrf_token cookie to be set")
	}
	if csrf.HttpOnly {
		t.Fatal("expected csrf_token cookie to NOT be HttpOnly")
	}
	if csrf.Value == "" {
		t.Fatal("expected csrf_token cookie to have a non-empty value")
	}
}

func TestCookieForwardResponseOption_ClearsCookiesOnLogout(t *testing.T) {
	rec, run := serverMetadataContext(middleware.ClearAuthCookiesKey, "1")

	err := run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cookies := rec.Result().Cookies()
	if len(cookies) != 3 {
		t.Fatalf("expected 3 cookies cleared, got %d: %v", len(cookies), cookies)
	}

	for _, c := range cookies {
		if c.MaxAge != -1 {
			t.Fatalf("expected cookie %q to have MaxAge=-1, got %d", c.Name, c.MaxAge)
		}
	}
}

func TestCookieForwardResponseOption_NoOpWhenNoRelevantMetadata(t *testing.T) {
	rec, run := serverMetadataContext("some-other-key", "value")

	err := run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cookies := rec.Result().Cookies()
	if len(cookies) != 0 {
		t.Fatalf("expected no cookies set, got %v", cookies)
	}
}
