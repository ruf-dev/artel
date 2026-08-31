package mcp_api

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/ruf-dev/artel/internal/middleware"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service"
	"go.redsock.ru/rerrors"
)

// csrfHeaderName is the request header the browser echoes the csrf_token cookie back in for the
// cookie-authenticated OAuth consent endpoints. These handlers run outside the grpc-gateway, so
// the value arrives verbatim rather than under grpc-gateway's "Grpc-Metadata-" prefix.
const csrfHeaderName = "X-Csrf-Token"

const (
	randomHexLen               = 16
	authCodeTTL                = 5 * time.Minute
	oauthAccessTokenTTLSeconds = 86400 // 24h
	grantTypeAuthorizationCode = "authorization_code"
)

type OAuthHandler struct {
	authSvc  service.AuthService
	vaultSvc service.VaultService
	mcpSvc   service.McpService

	pendingCodes repository.PendingAuthCodes
}

func NewOAuthHandler(
	authSvc service.AuthService,
	vaultSvc service.VaultService,
	mcpSvc service.McpService,
	pendingCodes repository.PendingAuthCodes,
) *OAuthHandler {
	return &OAuthHandler{
		authSvc:      authSvc,
		vaultSvc:     vaultSvc,
		mcpSvc:       mcpSvc,
		pendingCodes: pendingCodes,
	}
}

func writeJSON(w http.ResponseWriter, v any) {
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		log.Error().Err(err).Msg("failed to write json response")
	}
}

func requestBaseURL(r *http.Request) string {
	scheme := "http"
	if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
		scheme = proto
	} else if r.TLS != nil {
		scheme = "https"
	}

	host := r.Header.Get("X-Forwarded-Host")
	if host == "" {
		host = r.Host
	}

	return scheme + "://" + host
}

// ServeProtectedResourceMeta serves GET /.well-known/oauth-protected-resource (RFC 9728).
func (h *OAuthHandler) ServeProtectedResourceMeta(w http.ResponseWriter, r *http.Request) {
	base := requestBaseURL(r)
	meta := map[string]any{
		"resource":              base + "/mcp",
		"authorization_servers": []string{base},
	}

	w.Header().Set("Content-Type", "application/json")
	writeJSON(w, meta)
}

// WellKnown serves GET /.well-known/oauth-authorization-server.
func (h *OAuthHandler) WellKnown(w http.ResponseWriter, r *http.Request) {
	base := requestBaseURL(r)
	meta := map[string]any{
		"issuer":                                base,
		"authorization_endpoint":                base + "/authorize",
		"token_endpoint":                        base + "/token",
		"registration_endpoint":                 base + "/register",
		"response_types_supported":              []string{"code"},
		"code_challenge_methods_supported":      []string{"S256"},
		"grant_types_supported":                 []string{grantTypeAuthorizationCode},
		"token_endpoint_auth_methods_supported": []string{"none"},
	}

	w.Header().Set("Content-Type", "application/json")
	writeJSON(w, meta)
}

// ServeRegistration handles POST /register (RFC 7591 dynamic client registration).
func (h *OAuthHandler) ServeRegistration(writer http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(writer, "bad request", http.StatusBadRequest)

		return
	}

	var req map[string]any

	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(writer, "bad request", http.StatusBadRequest)

		return
	}

	resp := map[string]any{
		"client_id":                  randomHex(randomHexLen),
		"client_id_issued_at":        time.Now().Unix(),
		"grant_types":                []string{grantTypeAuthorizationCode},
		"response_types":             []string{"code"},
		"token_endpoint_auth_method": "none",
	}
	if v, ok := req["redirect_uris"]; ok {
		resp["redirect_uris"] = v
	}

	if v, ok := req["client_name"]; ok {
		resp["client_name"] = v
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	writeJSON(writer, resp)
}

// ServeOAuthLogin handles POST /api/oauth/login — validates Telegram id_token, returns session + vaults.
func (h *OAuthHandler) ServeOAuthLogin(writer http.ResponseWriter, request *http.Request) {
	var req struct {
		IdToken string `json:"id_token"`
	}

	err := json.NewDecoder(request.Body).Decode(&req)
	if err != nil || req.IdToken == "" {
		jsonErr(writer, "invalid request", http.StatusBadRequest)

		return
	}

	session, err := h.authSvc.LoginViaTelegram(request.Context(), req.IdToken)
	if err != nil {
		jsonErr(writer, "authentication failed", http.StatusUnauthorized)

		return
	}

	user, err := h.authSvc.ValidateToken(request.Context(), session.Token)
	if err != nil {
		jsonErr(writer, "authentication failed", http.StatusUnauthorized)

		return
	}

	uc := user_context.UserContext{UserUuid: user.Uuid}
	ctx := user_context.WithUserContext(request.Context(), uc)

	vaults, err := h.vaultSvc.ListVaults(ctx)
	if err != nil {
		jsonErr(writer, "failed to load vaults", http.StatusInternalServerError)

		return
	}

	type vaultItem struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}

	items := make([]vaultItem, len(vaults))
	for i, v := range vaults {
		items[i] = vaultItem{Id: v.Uuid.String(), Name: v.Name}
	}

	writer.Header().Set("Content-Type", "application/json")
	writeJSON(writer, map[string]any{
		"session_token": session.Token,
		"vaults":        items,
	})
}

// resolveOAuthSessionToken picks the session token for the cookie-or-body-authenticated OAuth
// consent endpoints (ServeOAuthVaults, ServeOAuthVault). A non-empty bodyToken is the Telegram /
// localStorage path (McpLogin stores mcpSessionToken and passes it explicitly) — taken as-is,
// with no cookie and no CSRF check. An empty bodyToken means the httpOnly browser-session path:
// the token comes from the access_token cookie, and because that path is driven purely by
// ambient cookies it must clear a double-submit CSRF check (csrf_token cookie vs X-Csrf-Token
// header). ok == false means a response has already been written and the caller must return.
func resolveOAuthSessionToken(w http.ResponseWriter, r *http.Request, bodyToken string) (token string, ok bool) {
	if bodyToken != "" {
		return bodyToken, true
	}

	accessCookie, err := r.Cookie(middleware.AccessTokenCookieName)
	if err != nil || accessCookie.Value == "" {
		jsonErr(w, "session expired", http.StatusUnauthorized)

		return "", false
	}

	csrfCookie, err := r.Cookie(middleware.CSRFCookieName)
	if err != nil || csrfCookie.Value == "" {
		jsonErr(w, "csrf token missing or invalid", http.StatusForbidden)

		return "", false
	}

	headerToken := r.Header.Get(csrfHeaderName)
	if subtle.ConstantTimeCompare([]byte(csrfCookie.Value), []byte(headerToken)) != 1 {
		jsonErr(w, "csrf token missing or invalid", http.StatusForbidden)

		return "", false
	}

	return accessCookie.Value, true
}

// ServeOAuthVaults handles POST /api/oauth/vaults — validates existing session and returns vault list.
func (h *OAuthHandler) ServeOAuthVaults(writer http.ResponseWriter, r *http.Request) {
	var req struct {
		SessionToken string `json:"session_token"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		jsonErr(writer, "invalid request", http.StatusBadRequest)

		return
	}

	token, ok := resolveOAuthSessionToken(writer, r, req.SessionToken)
	if !ok {
		return
	}

	user, err := h.authSvc.ValidateToken(r.Context(), token)
	if err != nil {
		jsonErr(writer, "session expired", http.StatusUnauthorized)

		return
	}

	uc := user_context.UserContext{UserUuid: user.Uuid}
	ctx := user_context.WithUserContext(r.Context(), uc)

	vaults, err := h.vaultSvc.ListVaults(ctx)
	if err != nil {
		jsonErr(writer, "failed to load vaults", http.StatusInternalServerError)

		return
	}

	type vaultItem struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}

	items := make([]vaultItem, len(vaults))
	for i, v := range vaults {
		items[i] = vaultItem{Id: v.Uuid.String(), Name: v.Name}
	}

	writer.Header().Set("Content-Type", "application/json")
	writeJSON(writer, map[string]any{"vaults": items})
}

// ServeOAuthVault handles POST /api/oauth/vault — creates MCP key and returns redirect URL.
func (h *OAuthHandler) ServeOAuthVault(writer http.ResponseWriter, request *http.Request) {
	var req struct {
		SessionToken  string `json:"session_token"`
		VaultId       string `json:"vault_id"`
		ClientId      string `json:"client_id"`
		RedirectUri   string `json:"redirect_uri"`
		CodeChallenge string `json:"code_challenge"`
		State         string `json:"state"`
	}

	err := json.NewDecoder(request.Body).Decode(&req)
	if err != nil {
		jsonErr(writer, "invalid request", http.StatusBadRequest)

		return
	}

	token, ok := resolveOAuthSessionToken(writer, request, req.SessionToken)
	if !ok {
		return
	}

	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		jsonErr(writer, "invalid vault_id", http.StatusBadRequest)

		return
	}

	user, err := h.authSvc.ValidateToken(request.Context(), token)
	if err != nil {
		jsonErr(writer, "session expired", http.StatusUnauthorized)

		return
	}

	uc := user_context.UserContext{UserUuid: user.Uuid}
	ctx := user_context.WithUserContext(request.Context(), uc)

	rawToken, _, err := h.mcpSvc.CreateKey(ctx, vaultID, "Claude MCP")
	if err != nil {
		jsonErr(writer, "failed to create access key", http.StatusInternalServerError)

		return
	}

	authCode := randomHex(randomHexLen)
	expiresAt := time.Now().Add(authCodeTTL)

	err = h.pendingCodes.Create(
		request.Context(), authCode, rawToken, req.CodeChallenge, req.RedirectUri, req.ClientId, expiresAt,
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to store pending auth code")
		jsonErr(writer, "internal error", http.StatusInternalServerError)

		return
	}

	redirectURL := fmt.Sprintf("%s?code=%s&state=%s", req.RedirectUri, authCode, req.State)

	writer.Header().Set("Content-Type", "application/json")
	writeJSON(writer, map[string]string{"redirect_url": redirectURL})
}

// ServeToken handles POST /token — PKCE verification and returns access token.
func (h *OAuthHandler) ServeToken(writer http.ResponseWriter, r *http.Request) {
	vals, err := parseTokenRequest(r)
	if err != nil {
		oauthTokenError(writer, "invalid_request", "cannot parse request")

		return
	}

	if vals["grant_type"] != grantTypeAuthorizationCode {
		oauthTokenError(writer, "unsupported_grant_type", "only authorization_code is supported")

		return
	}

	code := vals["code"]

	pending, err := h.pendingCodes.Get(r.Context(), code)
	if err != nil {
		oauthTokenError(writer, "invalid_grant", "authorization code expired or not found")

		return
	}

	if time.Now().After(pending.ExpiresAt) {
		oauthTokenError(writer, "invalid_grant", "authorization code expired")

		return
	}

	if pending.ClientId != "" && pending.ClientId != vals["client_id"] {
		oauthTokenError(writer, "invalid_client", "client_id mismatch")

		return
	}

	if pending.RedirectUri != "" && pending.RedirectUri != vals["redirect_uri"] {
		oauthTokenError(writer, "invalid_grant", "redirect_uri mismatch")

		return
	}

	if pending.CodeChallenge != "" && !pkceVerify(vals["code_verifier"], pending.CodeChallenge) {
		oauthTokenError(writer, "invalid_grant", "PKCE verification failed")

		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writeJSON(writer, map[string]any{
		"access_token": pending.RawToken,
		"token_type":   "Bearer",
		"expires_in":   oauthAccessTokenTTLSeconds,
	})
}

func parseTokenRequest(r *http.Request) (map[string]string, error) {
	ct := r.Header.Get("Content-Type")
	if strings.HasPrefix(ct, "application/json") {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, rerrors.Wrap(err, "failed to read request body")
		}

		var m map[string]string

		err = json.Unmarshal(body, &m)
		if err != nil {
			return nil, rerrors.Wrap(err, "error unmarshaling request body")
		}

		return m, nil
	}

	err := r.ParseForm()
	if err != nil {
		return nil, rerrors.Wrap(err, "error parsing form")
	}

	m := make(map[string]string)

	for k, v := range r.Form {
		if len(v) > 0 {
			m[k] = v[0]
		}
	}

	return m, nil
}

func pkceVerify(verifier, challenge string) bool {
	hash := sha256.Sum256([]byte(verifier))

	return base64.RawURLEncoding.EncodeToString(hash[:]) == challenge
}

func randomHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)

	return hex.EncodeToString(b)
}

func oauthTokenError(w http.ResponseWriter, code, desc string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	writeJSON(w, map[string]string{"error": code, "error_description": desc})
}

func jsonErr(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	writeJSON(w, map[string]string{"error": msg})
}
