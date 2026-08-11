package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

// csrfTokenNonceBytes is the size of the fresh nonce minted for the csrf_token cookie on every
// login/refresh — the CSRF token is independent of the session tokens, never derived from them.
const csrfTokenNonceBytes = 32

// insecureCookieWarnOnce logs the plain-HTTP notice (see requestWasSecure) at most once per
// process — every subsequent insecure cookie issuance stays silent.
var insecureCookieWarnOnce sync.Once

// CookieForwardResponseOption returns a runtime.WithForwardResponseOption hook for the shared
// gateway mux (see transport.NewGatewayMux). It is a no-op for almost every RPC — it only acts
// when the handler explicitly set one of the x-set-cookie-*/x-clear-auth-cookies metadata keys
// via grpc.SetHeader, which only a project's own login/refresh/logout handlers would do.
//
// The Secure attribute on every cookie it sets is decided per-request by requestWasSecure, fed by
// RequestSchemeAnnotator — no config knob involved.
func CookieForwardResponseOption() func(context.Context, http.ResponseWriter, proto.Message) error {
	return func(ctx context.Context, w http.ResponseWriter, _ proto.Message) error {
		serverMD, ok := runtime.ServerMetadataFromContext(ctx)
		if !ok {
			return nil
		}

		secure := requestWasSecure(ctx)

		if metadataValue(serverMD.HeaderMD, ClearAuthCookiesKey) == ClearAuthCookiesValue {
			warnIfInsecure(secure)
			clearAuthCookies(w, secure)
			return nil
		}

		accessToken := metadataValue(serverMD.HeaderMD, SetCookieAccessTokenKey)
		if accessToken == "" {
			return nil
		}

		warnIfInsecure(secure)

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

// warnIfInsecure logs the plain-HTTP notice at most once per process, the first time a cookie
// gets issued (set or cleared) without Secure — purely informational now that scheme detection is
// automatic and there's no config value an operator could get wrong.
func warnIfInsecure(secure bool) {
	if secure {
		return
	}

	insecureCookieWarnOnce.Do(func() {
		log.Warn().Msg("issuing auth cookies without Secure — this instance is serving over plain HTTP")
	})
}

// requestWasSecure reports whether RequestSchemeAnnotator determined the originating HTTP
// request was secure — see internal/middleware/request_scheme_annotator.go. Absence of
// RequestSecureKey (including when ctx carries no incoming metadata at all) means false.
func requestWasSecure(ctx context.Context) bool {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return false
	}

	return metadataValue(md, RequestSecureKey) == RequestSecureValue
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
