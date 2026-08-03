package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

// csrfTokenNonceBytes is the size of the fresh nonce minted for the csrf_token cookie on every
// Login/Refresh — the CSRF token is independent of the session tokens, never derived from them.
const csrfTokenNonceBytes = 32

// CookieForwardResponseOption returns a runtime.WithForwardResponseOption hook for the shared
// gateway mux (see transport.NewGatewayMux). It is a no-op for almost every RPC — it only acts
// when the handler explicitly set one of the x-set-cookie-*/x-clear-auth-cookies metadata keys
// via grpc.SetHeader, which only auth_api's Login/Refresh/Logout handlers do (see auth_impl.go).
//
// secure controls the Secure attribute on every cookie it sets (config.EnvironmentConfig.CookieSecure —
// false for plain-HTTP local dev, true otherwise).
func CookieForwardResponseOption(secure bool) func(context.Context, http.ResponseWriter, proto.Message) error {
	return func(ctx context.Context, w http.ResponseWriter, _ proto.Message) error {
		serverMD, ok := runtime.ServerMetadataFromContext(ctx)
		if !ok {
			return nil
		}

		if metadataValue(serverMD.HeaderMD, ClearAuthCookiesKey) == ClearAuthCookiesValue {
			clearAuthCookies(w, secure)
			return nil
		}

		accessToken := metadataValue(serverMD.HeaderMD, SetCookieAccessTokenKey)
		if accessToken == "" {
			return nil
		}

		accessExpiry := parseCookieExpiry(metadataValue(serverMD.HeaderMD, SetCookieAccessTokenExpiryKey))
		setAuthCookie(w, AccessTokenCookieName, accessToken, accessExpiry, true, secure, CookiePath)

		refreshToken := metadataValue(serverMD.HeaderMD, SetCookieRefreshTokenKey)
		if refreshToken != "" {
			refreshExpiry := parseCookieExpiry(metadataValue(serverMD.HeaderMD, SetCookieRefreshTokenExpiryKey))
			setAuthCookie(w, RefreshTokenCookieName, refreshToken, refreshExpiry, true, secure, CookiePath)
		}

		csrfToken, err := generateCSRFToken()
		if err != nil {
			log.Error().Err(err).Msg("error generating csrf token")
			return nil
		}
		setAuthCookie(w, CSRFCookieName, csrfToken, accessExpiry, false, secure, CSRFCookiePath)

		return nil
	}
}

func metadataValue(md metadata.MD, key string) string {
	vals := md.Get(key)
	if len(vals) == 0 {
		return ""
	}

	return vals[0]
}

func parseCookieExpiry(v string) time.Time {
	if v == "" {
		return time.Time{}
	}

	t, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return time.Time{}
	}

	return t
}

func setAuthCookie(w http.ResponseWriter, name, value string, expires time.Time, httpOnly, secure bool, path string) {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     path,
		HttpOnly: httpOnly,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	}
	if !expires.IsZero() {
		cookie.Expires = expires
	}

	http.SetCookie(w, &cookie)
}

// clearAuthCookies expires all three auth cookies (Max-Age=0), mirrored on the client by the
// browser dropping them immediately. Each cookie must be cleared with the same Path it was set
// with — browsers key deletion on (name, Path), not just name.
func clearAuthCookies(w http.ResponseWriter, secure bool) {
	expireCookie(w, AccessTokenCookieName, true, secure, CookiePath)
	expireCookie(w, RefreshTokenCookieName, true, secure, CookiePath)
	expireCookie(w, CSRFCookieName, false, secure, CSRFCookiePath)
}

func expireCookie(w http.ResponseWriter, name string, httpOnly, secure bool, path string) {
	cookie := http.Cookie{
		Name:     name,
		Value:    "",
		Path:     path,
		HttpOnly: httpOnly,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	}

	http.SetCookie(w, &cookie)
}

func generateCSRFToken() (string, error) {
	buf := make([]byte, csrfTokenNonceBytes)

	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(buf), nil
}
