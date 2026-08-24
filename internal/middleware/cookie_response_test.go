package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/metadata"
)

// buildForwardResponseCtx wires up the two contexts CookieForwardResponseOption reads from: the
// incoming gRPC metadata (requestWasSecure, set by RequestSchemeAnnotator in production) and the
// grpc-gateway ServerMetadata (the x-set-cookie-*/x-clear-auth-cookies keys a handler sets via
// grpc.SetHeader).
func buildForwardResponseCtx(secure bool, headerMD metadata.MD) context.Context {
	ctx := context.Background()

	if secure {
		incomingMD := metadata.Pairs(RequestSecureKey, RequestSecureValue)
		ctx = metadata.NewIncomingContext(ctx, incomingMD)
	} else {
		ctx = metadata.NewIncomingContext(ctx, metadata.MD{})
	}

	serverMD := runtime.ServerMetadata{HeaderMD: headerMD}
	ctx = runtime.NewServerMetadataContext(ctx, serverMD)

	return ctx
}

func TestCookieForwardResponseOption_SetCookie(t *testing.T) {
	tests := []struct {
		name       string
		secure     bool
		wantSecure bool
	}{
		{name: "secure request sets Secure cookies", secure: true, wantSecure: true},
		{name: "insecure request sets non-Secure cookies", secure: false, wantSecure: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headerMD := metadata.Pairs(
				SetCookieAccessTokenKey, "access-token-value",
				SetCookieRefreshTokenKey, "refresh-token-value",
			)
			ctx := buildForwardResponseCtx(tt.secure, headerMD)

			w := httptest.NewRecorder()

			opt := CookieForwardResponseOption()

			err := opt(ctx, w, nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			cookies := w.Result().Cookies()
			if len(cookies) == 0 {
				t.Fatal("expected cookies to be set, got none")
			}

			for _, c := range cookies {
				if c.Secure != tt.wantSecure {
					t.Fatalf("cookie %q Secure = %v, want %v", c.Name, c.Secure, tt.wantSecure)
				}
			}
		})
	}
}

func TestCookieForwardResponseOption_Logout(t *testing.T) {
	tests := []struct {
		name       string
		secure     bool
		wantSecure bool
	}{
		{name: "logout over secure request clears Secure cookies", secure: true, wantSecure: true},
		{name: "logout over insecure request clears non-Secure cookies", secure: false, wantSecure: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headerMD := metadata.Pairs(ClearAuthCookiesKey, ClearAuthCookiesValue)
			ctx := buildForwardResponseCtx(tt.secure, headerMD)

			w := httptest.NewRecorder()

			opt := CookieForwardResponseOption()

			err := opt(ctx, w, nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			cookies := w.Result().Cookies()
			if len(cookies) == 0 {
				t.Fatal("expected cleared cookies to be set, got none")
			}

			for _, c := range cookies {
				if c.MaxAge != -1 {
					t.Fatalf("cookie %q MaxAge = %d, want -1 (expired)", c.Name, c.MaxAge)
				}

				if c.Secure != tt.wantSecure {
					t.Fatalf("cookie %q Secure = %v, want %v", c.Name, c.Secure, tt.wantSecure)
				}
			}
		})
	}
}

func TestCookieForwardResponseOption_NoServerMetadata(t *testing.T) {
	w := httptest.NewRecorder()

	opt := CookieForwardResponseOption()

	err := opt(context.Background(), w, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(w.Result().Cookies()) != 0 {
		t.Fatal("expected no cookies without ServerMetadata in context")
	}
}

func TestRequestWasSecure(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
		want bool
	}{
		{
			name: "no incoming metadata",
			ctx:  context.Background(),
			want: false,
		},
		{
			name: "incoming metadata without the secure key",
			ctx:  metadata.NewIncomingContext(context.Background(), metadata.MD{}),
			want: false,
		},
		{
			name: "incoming metadata with the secure key",
			ctx: metadata.NewIncomingContext(
				context.Background(), metadata.Pairs(RequestSecureKey, RequestSecureValue),
			),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := requestWasSecure(tt.ctx)
			if got != tt.want {
				t.Fatalf("requestWasSecure() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestCookieForwardResponseOption_CSRFExpiryTracksRefreshToken guards against the csrf_token
// cookie reverting to the access token's short expiry: isAuthenticated() on the frontend reads
// only csrf_token's presence, synchronously and without an API call, so if it expired alongside
// the (much shorter-lived) access token, users would be bounced to the login page long before
// their still-valid refresh token actually ran out.
func TestCookieForwardResponseOption_CSRFExpiryTracksRefreshToken(t *testing.T) {
	accessExpiry := time.Now().Add(time.Hour).UTC().Format(time.RFC3339)
	refreshExpiry := time.Now().Add(30 * 24 * time.Hour).UTC().Format(time.RFC3339)

	headerMD := metadata.Pairs(
		SetCookieAccessTokenKey, "access-token-value",
		SetCookieAccessTokenExpiryKey, accessExpiry,
		SetCookieRefreshTokenKey, "refresh-token-value",
		SetCookieRefreshTokenExpiryKey, refreshExpiry,
	)
	ctx := buildForwardResponseCtx(true, headerMD)

	w := httptest.NewRecorder()

	opt := CookieForwardResponseOption()

	err := opt(ctx, w, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var csrfCookie *http.Cookie
	for _, c := range w.Result().Cookies() {
		if c.Name == CSRFCookieName {
			csrfCookie = c
		}
	}
	if csrfCookie == nil {
		t.Fatal("expected csrf_token cookie to be set")
	}

	wantExpiry := parseCookieExpiry(refreshExpiry)
	if !csrfCookie.Expires.Equal(wantExpiry) {
		t.Fatalf("csrf_token cookie Expires = %v, want %v (refresh token expiry)", csrfCookie.Expires, wantExpiry)
	}
}

// TestSetCookieHeaderContainsSecureAttribute is a thin sanity check on the raw Set-Cookie header
// text (not just the parsed http.Cookie.Secure field checked above), since Secure is the attribute
// browsers actually key their "silently drop this cookie over plain HTTP" behavior on — the bug
// this whole feature exists to fix.
func TestSetCookieHeaderContainsSecureAttribute(t *testing.T) {
	headerMD := metadata.Pairs(SetCookieAccessTokenKey, "access-token-value")

	secureCtx := buildForwardResponseCtx(true, headerMD)
	insecureCtx := buildForwardResponseCtx(false, headerMD)

	opt := CookieForwardResponseOption()

	secureW := httptest.NewRecorder()

	err := opt(secureCtx, secureW, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	insecureW := httptest.NewRecorder()

	err = opt(insecureCtx, insecureW, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	secureHeader := strings.Join(secureW.Result().Header.Values("Set-Cookie"), "; ")
	if !strings.Contains(secureHeader, "Secure") {
		t.Fatalf("expected Secure attribute in Set-Cookie header, got: %s", secureHeader)
	}

	insecureHeader := strings.Join(insecureW.Result().Header.Values("Set-Cookie"), "; ")
	if strings.Contains(insecureHeader, "Secure") {
		t.Fatalf("expected no Secure attribute in Set-Cookie header, got: %s", insecureHeader)
	}
}
