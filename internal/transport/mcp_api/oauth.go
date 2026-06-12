package mcp_api

import (
	"crypto/rand"
	"crypto/sha256"
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

	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service"
)

type OAuthHandler struct {
	authSvc  service.AuthService
	vaultSvc service.VaultService
	mcpSvc   service.McpService

	pendingCodes repository.PendingAuthCodes
}

func NewOAuthHandler(authSvc service.AuthService, vaultSvc service.VaultService, mcpSvc service.McpService, pendingCodes repository.PendingAuthCodes) *OAuthHandler {
	return &OAuthHandler{
		authSvc:      authSvc,
		vaultSvc:     vaultSvc,
		mcpSvc:       mcpSvc,
		pendingCodes: pendingCodes,
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

// ServeProtectedResourceMeta serves GET /.well-known/oauth-protected-resource (RFC 9728)
func (h *OAuthHandler) ServeProtectedResourceMeta(w http.ResponseWriter, r *http.Request) {
	base := requestBaseURL(r)
	meta := map[string]any{
		"resource":              base + "/mcp",
		"authorization_servers": []string{base},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(meta)
}

// WellKnown serves GET /.well-known/oauth-authorization-server
func (h *OAuthHandler) WellKnown(w http.ResponseWriter, r *http.Request) {
	base := requestBaseURL(r)
	meta := map[string]any{
		"issuer":                                base,
		"authorization_endpoint":                base + "/authorize",
		"token_endpoint":                        base + "/token",
		"registration_endpoint":                 base + "/register",
		"response_types_supported":              []string{"code"},
		"code_challenge_methods_supported":      []string{"S256"},
		"grant_types_supported":                 []string{"authorization_code"},
		"token_endpoint_auth_methods_supported": []string{"none"},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(meta)
}

// ServeRegistration handles POST /register (RFC 7591 dynamic client registration)
func (h *OAuthHandler) ServeRegistration(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	var req map[string]any
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	resp := map[string]any{
		"client_id":                  randomHex(16),
		"client_id_issued_at":        time.Now().Unix(),
		"grant_types":                []string{"authorization_code"},
		"response_types":             []string{"code"},
		"token_endpoint_auth_method": "none",
	}
	if v, ok := req["redirect_uris"]; ok {
		resp["redirect_uris"] = v
	}
	if v, ok := req["client_name"]; ok {
		resp["client_name"] = v
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// ServeOAuthLogin handles POST /oauth/login — validates Telegram id_token, returns session + vaults
func (h *OAuthHandler) ServeOAuthLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IdToken string `json:"id_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.IdToken == "" {
		jsonErr(w, "invalid request", http.StatusBadRequest)
		return
	}

	session, err := h.authSvc.LoginViaTelegram(r.Context(), req.IdToken)
	if err != nil {
		jsonErr(w, "authentication failed", http.StatusUnauthorized)
		return
	}

	user, err := h.authSvc.ValidateToken(r.Context(), session.Token)
	if err != nil {
		jsonErr(w, "authentication failed", http.StatusUnauthorized)
		return
	}

	uc := user_context.UserContext{UserUuid: user.Uuid}
	ctx := user_context.WithUserContext(r.Context(), uc)
	vaults, err := h.vaultSvc.ListVaults(ctx)
	if err != nil {
		jsonErr(w, "failed to load vaults", http.StatusInternalServerError)
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"session_token": session.Token,
		"vaults":        items,
	})
}

// ServeOAuthVaults handles POST /oauth/vaults — validates existing session and returns vault list
func (h *OAuthHandler) ServeOAuthVaults(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SessionToken string `json:"session_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.SessionToken == "" {
		jsonErr(w, "invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.authSvc.ValidateToken(r.Context(), req.SessionToken)
	if err != nil {
		jsonErr(w, "session expired", http.StatusUnauthorized)
		return
	}

	uc := user_context.UserContext{UserUuid: user.Uuid}
	ctx := user_context.WithUserContext(r.Context(), uc)
	vaults, err := h.vaultSvc.ListVaults(ctx)
	if err != nil {
		jsonErr(w, "failed to load vaults", http.StatusInternalServerError)
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"vaults": items})
}

// ServeOAuthVault handles POST /oauth/vault — creates MCP key and returns redirect URL
func (h *OAuthHandler) ServeOAuthVault(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SessionToken  string `json:"session_token"`
		VaultId       string `json:"vault_id"`
		ClientId      string `json:"client_id"`
		RedirectUri   string `json:"redirect_uri"`
		CodeChallenge string `json:"code_challenge"`
		State         string `json:"state"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonErr(w, "invalid request", http.StatusBadRequest)
		return
	}

	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		jsonErr(w, "invalid vault_id", http.StatusBadRequest)
		return
	}

	user, err := h.authSvc.ValidateToken(r.Context(), req.SessionToken)
	if err != nil {
		jsonErr(w, "session expired", http.StatusUnauthorized)
		return
	}

	uc := user_context.UserContext{UserUuid: user.Uuid}
	ctx := user_context.WithUserContext(r.Context(), uc)
	rawToken, _, err := h.mcpSvc.CreateKey(ctx, vaultID, "Claude MCP")
	if err != nil {
		jsonErr(w, "failed to create access key", http.StatusInternalServerError)
		return
	}

	authCode := randomHex(16)
	expiresAt := time.Now().Add(5 * time.Minute)
	if err := h.pendingCodes.Create(r.Context(), authCode, rawToken, req.CodeChallenge, req.RedirectUri, req.ClientId, expiresAt); err != nil {
		log.Error().Err(err).Msg("failed to store pending auth code")
		jsonErr(w, "internal error", http.StatusInternalServerError)
		return
	}

	redirectURL := fmt.Sprintf("%s?code=%s&state=%s", req.RedirectUri, authCode, req.State)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"redirect_url": redirectURL})
}

// ServeToken handles POST /token — PKCE verification and returns access token
func (h *OAuthHandler) ServeToken(w http.ResponseWriter, r *http.Request) {
	vals, err := parseTokenRequest(r)
	if err != nil {
		oauthTokenError(w, "invalid_request", "cannot parse request")
		return
	}

	if vals["grant_type"] != "authorization_code" {
		oauthTokenError(w, "unsupported_grant_type", "only authorization_code is supported")
		return
	}

	code := vals["code"]
	pending, err := h.pendingCodes.Get(r.Context(), code)
	if err != nil {
		oauthTokenError(w, "invalid_grant", "authorization code expired or not found")
		return
	}

	//err = h.pendingCodes.Delete(r.Context(), code);
	//if  err != nil {
	//	log.Error().Err(err).Msg("failed to delete pending auth code")
	//}

	if time.Now().After(pending.ExpiresAt) {
		oauthTokenError(w, "invalid_grant", "authorization code expired")
		return
	}
	if pending.ClientId != "" && pending.ClientId != vals["client_id"] {
		oauthTokenError(w, "invalid_client", "client_id mismatch")
		return
	}
	if pending.RedirectUri != "" && pending.RedirectUri != vals["redirect_uri"] {
		oauthTokenError(w, "invalid_grant", "redirect_uri mismatch")
		return
	}
	if pending.CodeChallenge != "" && !pkceVerify(vals["code_verifier"], pending.CodeChallenge) {
		oauthTokenError(w, "invalid_grant", "PKCE verification failed")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"access_token": pending.RawToken,
		"token_type":   "Bearer",
		"expires_in":   86400,
	})
}

func parseTokenRequest(r *http.Request) (map[string]string, error) {
	ct := r.Header.Get("Content-Type")
	if strings.HasPrefix(ct, "application/json") {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		var m map[string]string
		if err := json.Unmarshal(body, &m); err != nil {
			return nil, err
		}
		return m, nil
	}
	if err := r.ParseForm(); err != nil {
		return nil, err
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
	json.NewEncoder(w).Encode(map[string]string{"error": code, "error_description": desc})
}

func jsonErr(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
