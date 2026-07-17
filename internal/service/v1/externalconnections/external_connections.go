package externalconnections

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"
	"golang.org/x/oauth2"

	"github.com/ruf-dev/artel/internal/clients/googleapi"
	"github.com/ruf-dev/artel/internal/clients/imap"
	"github.com/ruf-dev/artel/internal/clients/smtp"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/repository"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"github.com/ruf-dev/artel/internal/utils"
)

const googleUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"

type googleUserInfo struct {
	Email string `json:"email"`
}

const gitlabDefaultInstanceURL = "https://gitlab.com"

var gitlabValidationClient = &http.Client{Timeout: 10 * time.Second}

const trelloValidationURL = "https://api.trello.com/1/members/me"

var trelloValidationClient = &http.Client{Timeout: 10 * time.Second}

// emailConnectionCheckTimeout bounds each of the IMAP/SMTP dial+auth round trips in
// CheckEmailConnection, so a host that never responds fails fast instead of hanging the request.
const emailConnectionCheckTimeout = 10 * time.Second

type gitlabUserInfo struct {
	Username string `json:"username"`
}

// gitlabConnectionMeta is stored plaintext in metadata (non-sensitive display data).
type gitlabConnectionMeta struct {
	Username         string `json:"username"`
	InstanceUrl      string `json:"instance_url"`
	WebhookSecretSet string `json:"webhook_secret_set"`
}

// trelloConnectionMeta is stored plaintext in metadata (non-sensitive display data).
type trelloConnectionMeta struct {
	FullName string `json:"full_name"`
}

type Service struct {
	connections           repository.ExternalConnectionRepo
	pendingCodes          repository.PendingAuthCodes
	mcpSpreadsheets       repository.McpSpreadsheetsRepo
	mailServerSuggestions repository.MailServerSuggestions
	oauthCfg              *oauth2.Config
}

func New(
	connections repository.ExternalConnectionRepo,
	pendingCodes repository.PendingAuthCodes,
	mcpSpreadsheets repository.McpSpreadsheetsRepo,
	mailServerSuggestions repository.MailServerSuggestions,
	oauthCfg *oauth2.Config,
) *Service {
	return &Service{
		connections:           connections,
		pendingCodes:          pendingCodes,
		mcpSpreadsheets:       mcpSpreadsheets,
		mailServerSuggestions: mailServerSuggestions,
		oauthCfg:              oauthCfg,
	}
}

func (s *Service) InitiateGoogleOAuth(ctx context.Context, origin string) (string, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return "", user_errors.Unauthenticated
	}

	redirectURL := origin + "/connections/google/callback"
	state := randomHex(16)
	expiresAt := time.Now().Add(15 * time.Minute)

	err := s.pendingCodes.Create(ctx, state, uc.UserUuid.String(), "", redirectURL, "", expiresAt)
	if err != nil {
		return "", rerrors.Wrap(err, "error storing oauth state")
	}

	cfgCopy := *s.oauthCfg
	cfgCopy.RedirectURL = redirectURL
	authURL := cfgCopy.AuthCodeURL(state, oauth2.AccessTypeOffline)

	return authURL, nil
}

func (s *Service) HandleGoogleOAuthCallback(
	ctx context.Context,
	code string,
	state string,
) (domain.ExternalConnectionMeta, error) {
	pendingCode, err := s.pendingCodes.Get(ctx, state)
	if err != nil {
		return domain.ExternalConnectionMeta{}, user_errors.InvalidOAuthState
	}

	userUuid, err := uuid.Parse(pendingCode.RawToken)
	if err != nil {
		return domain.ExternalConnectionMeta{}, user_errors.InvalidOAuthState
	}

	cfgCopy := *s.oauthCfg
	cfgCopy.RedirectURL = pendingCode.RedirectUri

	token, err := cfgCopy.Exchange(ctx, code)
	if err != nil {
		return domain.ExternalConnectionMeta{}, rerrors.Wrap(err, "error exchanging oauth code")
	}

	httpClient := cfgCopy.Client(ctx, token)

	resp, err := httpClient.Get(googleUserInfoURL)
	if err != nil {
		return domain.ExternalConnectionMeta{}, rerrors.Wrap(err, "error fetching google userinfo")
	}

	defer utils.CloseWithLog(resp.Body, "google userinfo response")

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.ExternalConnectionMeta{}, rerrors.Wrap(err, "error reading google userinfo response")
	}

	var userInfo googleUserInfo

	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return domain.ExternalConnectionMeta{}, rerrors.Wrap(err, "error parsing google userinfo")
	}

	scope, _ := token.Extra("scope").(string)
	creds := domain.GoogleOAuthCredentials{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
		Scope:        scope,
		TokenType:    token.TokenType,
	}

	credJSON, err := json.Marshal(creds)
	if err != nil {
		return domain.ExternalConnectionMeta{}, rerrors.Wrap(err, "error marshaling google credentials")
	}

	meta := domain.GoogleConnectionMeta{Email: userInfo.Email, Scopes: scope}

	metaJSON, err := json.Marshal(meta)
	if err != nil {
		return domain.ExternalConnectionMeta{}, rerrors.Wrap(err, "error marshaling google connection meta")
	}

	conn := domain.ExternalConnection{
		UserUuid:        userUuid,
		Provider:        domain.ProviderGoogleSheets,
		ProviderType:    artel_q.ExternalProviderTypeGoogleOauth,
		CredentialsJSON: json.RawMessage(credJSON),
		Metadata:        json.RawMessage(metaJSON),
	}

	saved, err := s.connections.Upsert(ctx, conn)
	if err != nil {
		return domain.ExternalConnectionMeta{}, rerrors.Wrap(err, "error saving google connection")
	}

	err = s.pendingCodes.Delete(ctx, state)
	if err != nil {
		return domain.ExternalConnectionMeta{}, rerrors.Wrap(err, "error deleting oauth state")
	}

	return toMeta(saved, userInfo.Email), nil
}

func (s *Service) DisconnectProvider(ctx context.Context, provider string) error {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return user_errors.Unauthenticated
	}

	err := s.connections.Delete(ctx, uc.UserUuid, provider)
	if err != nil {
		return rerrors.Wrap(err, "error disconnecting provider")
	}

	return nil
}

// DisconnectConnection removes a single connection by id, unlike DisconnectProvider which removes
// whichever single row exists for a provider — needed for providers like Trello that can hold more
// than one connection per user.
func (s *Service) DisconnectConnection(ctx context.Context, id string) error {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return user_errors.Unauthenticated
	}

	connUuid, err := uuid.Parse(id)
	if err != nil {
		return rerrors.Wrap(err, "parse uuid")
	}

	err = s.connections.DeleteByID(ctx, uc.UserUuid, connUuid)
	if err != nil {
		return rerrors.Wrap(err, "error disconnecting connection")
	}

	return nil
}

func (s *Service) ListConnections(ctx context.Context) ([]domain.ExternalConnectionMeta, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return nil, user_errors.Unauthenticated
	}

	conns, err := s.connections.ListByUser(ctx, uc.UserUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing external connections")
	}

	metas := make([]domain.ExternalConnectionMeta, len(conns))

	for i, conn := range conns {
		displayName := extractDisplayName(conn)
		metas[i] = toMeta(conn, displayName)
	}

	return metas, nil
}

func (s *Service) GetGoogleClient(ctx context.Context) (*googleapi.Client, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return nil, user_errors.Unauthenticated
	}

	creds, err := s.freshGoogleCreds(ctx, uc.UserUuid)
	if err != nil {
		return nil, err
	}

	client := googleapi.New(ctx, creds, s.oauthCfg)

	return client, nil
}

func (s *Service) GetPickerToken(ctx context.Context) (string, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return "", user_errors.Unauthenticated
	}

	creds, err := s.freshGoogleCreds(ctx, uc.UserUuid)
	if err != nil {
		return "", err
	}

	return creds.AccessToken, nil
}

func (s *Service) AddSpreadsheet(
	ctx context.Context,
	spreadsheetId string,
	name string,
) (domain.McpSpreadsheet, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return domain.McpSpreadsheet{}, user_errors.Unauthenticated
	}

	result, err := s.connections.GetByUserAndProvider(ctx, uc.UserUuid, domain.ProviderGoogleSheets)
	if err != nil {
		return domain.McpSpreadsheet{}, rerrors.Wrap(err, "error getting google connection")
	}

	if !result.Valid {
		return domain.McpSpreadsheet{}, user_errors.GoogleNotConnected
	}

	spreadsheet := domain.McpSpreadsheet{
		UserUuid:             uc.UserUuid,
		ExternalConnectionId: result.V.Uuid,
		SpreadsheetId:        spreadsheetId,
		Name:                 name,
	}

	saved, err := s.mcpSpreadsheets.Insert(ctx, spreadsheet)
	if err != nil {
		return domain.McpSpreadsheet{}, rerrors.Wrap(err, "error adding spreadsheet")
	}

	return saved, nil
}

func (s *Service) ListSpreadsheets(ctx context.Context) ([]domain.McpSpreadsheet, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return nil, user_errors.Unauthenticated
	}

	spreadsheets, err := s.mcpSpreadsheets.ListByUser(ctx, uc.UserUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing spreadsheets")
	}

	return spreadsheets, nil
}

func (s *Service) RemoveSpreadsheet(ctx context.Context, spreadsheetId string) error {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return user_errors.Unauthenticated
	}

	err := s.mcpSpreadsheets.Delete(ctx, uc.UserUuid, spreadsheetId)
	if err != nil {
		return rerrors.Wrap(err, "error removing spreadsheet")
	}

	return nil
}

func (s *Service) AddEmailConnection(
	ctx context.Context,
	email, imapHost string,
	imapPort int,
	smtpHost string,
	smtpPort int,
	password string,
) (domain.ExternalConnectionMeta, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return domain.ExternalConnectionMeta{}, user_errors.Unauthenticated
	}

	if password == "" {
		var err error

		password, err = s.storedEmailPassword(ctx, uc.UserUuid)
		if err != nil {
			return domain.ExternalConnectionMeta{}, err
		}
	}

	creds := domain.EmailCredentials{
		ImapHost: imapHost,
		ImapPort: imapPort,
		SmtpHost: smtpHost,
		SmtpPort: smtpPort,
		Username: email,
		Password: password,
	}

	credJSON, err := json.Marshal(creds)
	if err != nil {
		return domain.ExternalConnectionMeta{}, rerrors.Wrap(err, "marshal email credentials")
	}

	type emailMeta struct {
		Username string `json:"username"`
	}

	metaJSON, err := json.Marshal(emailMeta{Username: email})
	if err != nil {
		return domain.ExternalConnectionMeta{}, rerrors.Wrap(err, "marshal email meta")
	}

	conn := domain.ExternalConnection{
		UserUuid:        uc.UserUuid,
		Provider:        domain.ProviderEmail,
		ProviderType:    artel_q.ExternalProviderTypePassword,
		CredentialsJSON: json.RawMessage(credJSON),
		Metadata:        json.RawMessage(metaJSON),
	}

	saved, err := s.connections.Upsert(ctx, conn)
	if err != nil {
		return domain.ExternalConnectionMeta{}, rerrors.Wrap(err, "save email connection")
	}

	return toMeta(saved, email), nil
}

// CheckEmailConnection verifies that both the IMAP and SMTP settings actually work, without
// persisting anything — backs the "Check" button in the add/edit email account form.
func (s *Service) CheckEmailConnection(
	ctx context.Context,
	email, imapHost string,
	imapPort int,
	smtpHost string,
	smtpPort int,
	password string,
) error {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return user_errors.Unauthenticated
	}

	if password == "" {
		var err error

		password, err = s.storedEmailPassword(ctx, uc.UserUuid)
		if err != nil {
			return err
		}
	}

	imapClient := imap.New(imapHost, imapPort, email, password)

	imapCtx, cancelImap := context.WithTimeout(ctx, emailConnectionCheckTimeout)
	defer cancelImap()

	err := imapClient.TestConnection(imapCtx)
	if err != nil {
		return rerrors.Wrap(user_errors.EmailImapValidationFailed, "error connecting to imap server")
	}

	smtpClient := smtp.New(smtpHost, smtpPort, email, password)

	smtpCtx, cancelSmtp := context.WithTimeout(ctx, emailConnectionCheckTimeout)
	defer cancelSmtp()

	err = smtpClient.TestConnection(smtpCtx)
	if err != nil {
		return rerrors.Wrap(user_errors.EmailSmtpValidationFailed, "error connecting to smtp server")
	}

	return nil
}

// storedEmailPassword resolves the password of the user's existing email connection, used to
// fall back when the caller submits a blank password (edit forms leave it blank to mean "keep
// the current password" for both saving and pinging).
func (s *Service) storedEmailPassword(ctx context.Context, userUuid uuid.UUID) (string, error) {
	existing, err := s.connections.GetByUserAndProvider(ctx, userUuid, domain.ProviderEmail)
	if err != nil {
		return "", rerrors.Wrap(err, "error loading existing email connection")
	}

	if !existing.Valid {
		return "", user_errors.EmailPasswordRequired
	}

	var creds domain.EmailCredentials

	err = json.Unmarshal(existing.V.CredentialsJSON, &creds)
	if err != nil {
		return "", rerrors.Wrap(err, "error parsing existing email credentials")
	}

	return creds.Password, nil
}

func (s *Service) AddGitlabConnection(
	ctx context.Context,
	personalAccessToken, instanceUrl string,
) (domain.ExternalConnectionMeta, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return domain.ExternalConnectionMeta{}, user_errors.Unauthenticated
	}

	normalizedInstanceUrl, err := normalizeGitlabInstanceURL(instanceUrl)
	if err != nil {
		return domain.ExternalConnectionMeta{}, err
	}

	username, err := s.validateGitlabToken(ctx, normalizedInstanceUrl, personalAccessToken)
	if err != nil {
		return domain.ExternalConnectionMeta{}, err
	}

	creds := domain.GitlabCredentials{
		PersonalAccessToken: personalAccessToken,
		InstanceUrl:         normalizedInstanceUrl,
	}

	credJSON, err := json.Marshal(creds)
	if err != nil {
		return domain.ExternalConnectionMeta{}, rerrors.Wrap(err, "marshal gitlab credentials")
	}

	meta := gitlabConnectionMeta{
		Username:         username,
		InstanceUrl:      normalizedInstanceUrl,
		WebhookSecretSet: "false",
	}

	metaJSON, err := json.Marshal(meta)
	if err != nil {
		return domain.ExternalConnectionMeta{}, rerrors.Wrap(err, "marshal gitlab meta")
	}

	conn := domain.ExternalConnection{
		UserUuid:        uc.UserUuid,
		Provider:        domain.ProviderGitlab,
		ProviderType:    artel_q.ExternalProviderTypeApiKey,
		CredentialsJSON: json.RawMessage(credJSON),
		Metadata:        json.RawMessage(metaJSON),
	}

	saved, err := s.connections.Upsert(ctx, conn)
	if err != nil {
		return domain.ExternalConnectionMeta{}, rerrors.Wrap(err, "save gitlab connection")
	}

	return toMeta(saved, username), nil
}

// CheckGitlabConnection pings the instance with the given token without persisting anything,
// letting the caller confirm the credentials work before committing to AddGitlabConnection.
func (s *Service) CheckGitlabConnection(
	ctx context.Context,
	personalAccessToken, instanceUrl string,
) (string, error) {
	_, ok := user_context.GetUserContext(ctx)
	if !ok {
		return "", user_errors.Unauthenticated
	}

	normalizedInstanceUrl, err := normalizeGitlabInstanceURL(instanceUrl)
	if err != nil {
		return "", err
	}

	username, err := s.validateGitlabToken(ctx, normalizedInstanceUrl, personalAccessToken)
	if err != nil {
		return "", err
	}

	return username, nil
}

// AddTrelloConnection validates the API key/token pair against Trello, then persists it.
func (s *Service) AddTrelloConnection(
	ctx context.Context,
	apiKey, apiToken string,
) (domain.ExternalConnectionMeta, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return domain.ExternalConnectionMeta{}, user_errors.Unauthenticated
	}

	fullName, err := s.validateTrelloCredentials(ctx, apiKey, apiToken)
	if err != nil {
		return domain.ExternalConnectionMeta{}, err
	}

	creds := domain.APIKeyCredentials{
		APIKey:   apiKey,
		APIToken: apiToken,
	}

	credJSON, err := json.Marshal(creds)
	if err != nil {
		return domain.ExternalConnectionMeta{}, rerrors.Wrap(err, "error marshaling trello credentials")
	}

	meta := trelloConnectionMeta{
		FullName: fullName,
	}

	metaJSON, err := json.Marshal(meta)
	if err != nil {
		return domain.ExternalConnectionMeta{}, rerrors.Wrap(err, "error marshaling trello meta")
	}

	conn := domain.ExternalConnection{
		UserUuid:        uc.UserUuid,
		Provider:        domain.ProviderTrello,
		ProviderType:    artel_q.ExternalProviderTypeApiKey,
		CredentialsJSON: json.RawMessage(credJSON),
		Metadata:        json.RawMessage(metaJSON),
	}

	saved, err := s.connections.Insert(ctx, conn)
	if err != nil {
		return domain.ExternalConnectionMeta{}, rerrors.Wrap(err, "error saving trello connection")
	}

	return toMeta(saved, fullName), nil
}

// CheckTrelloConnection pings Trello with the given credentials without persisting anything,
// letting the caller confirm they work before committing to AddTrelloConnection.
func (s *Service) CheckTrelloConnection(
	ctx context.Context,
	apiKey, apiToken string,
) (string, error) {
	_, ok := user_context.GetUserContext(ctx)
	if !ok {
		return "", user_errors.Unauthenticated
	}

	fullName, err := s.validateTrelloCredentials(ctx, apiKey, apiToken)
	if err != nil {
		return "", err
	}

	return fullName, nil
}

// trelloMemberResponse is Trello's GET /1/members/me response, trimmed to the one field used
// here.
type trelloMemberResponse struct {
	FullName string `json:"fullName"`
}

// validateTrelloCredentials confirms the key/token pair actually authenticates against Trello
// before anything is persisted. This can't go through the generic MoM http executor the way
// every later trello tool call does, because that executor resolves __secrets.* from an
// existing external_connections row — and at this point the connection doesn't exist yet (that's
// exactly what this call is validating before we create it). So it stays a small direct request
// here, mirroring validateGitlabToken above, rather than pulling in a dedicated client package
// for one unauthenticated-style GET.
func (s *Service) validateTrelloCredentials(ctx context.Context, apiKey, apiToken string) (string, error) {
	reqUrl := trelloValidationURL + "?key=" + url.QueryEscape(apiKey) + "&token=" + url.QueryEscape(apiToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl, nil)
	if err != nil {
		return "", rerrors.Wrap(err, "error building trello validation request")
	}

	resp, err := trelloValidationClient.Do(req)
	if err != nil {
		return "", rerrors.Wrap(user_errors.TrelloValidationFailed, "error connecting to trello")
	}
	defer utils.CloseWithLog(resp.Body, "trello validation response")

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", rerrors.Wrap(user_errors.TrelloValidationFailed, "trello rejected the credentials")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", rerrors.Wrap(err, "error reading trello validation response")
	}

	var member trelloMemberResponse

	err = json.Unmarshal(body, &member)
	if err != nil {
		return "", rerrors.Wrap(err, "error parsing trello validation response")
	}

	return member.FullName, nil
}

// GenerateGitlabWebhookSecret mints a fresh random secret for the caller's GitLab connection and
// returns it once — only the connection's encrypted credentials retain it after this call. The
// caller pastes the returned value into GitLab's own webhook config; gitlab_webhook.Handler
// compares inbound X-Gitlab-Token deliveries against the stored value.
func (s *Service) GenerateGitlabWebhookSecret(
	ctx context.Context,
) (domain.ExternalConnectionMeta, string, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return domain.ExternalConnectionMeta{}, "", user_errors.Unauthenticated
	}

	result, err := s.connections.GetByUserAndProvider(ctx, uc.UserUuid, domain.ProviderGitlab)
	if err != nil {
		return domain.ExternalConnectionMeta{}, "", rerrors.Wrap(err, "error getting gitlab connection")
	}

	if !result.Valid {
		return domain.ExternalConnectionMeta{}, "", user_errors.GitlabConnectionNotFound
	}

	conn := result.V

	var creds domain.GitlabCredentials

	err = json.Unmarshal(conn.CredentialsJSON, &creds)
	if err != nil {
		return domain.ExternalConnectionMeta{}, "", rerrors.Wrap(err, "error parsing gitlab credentials")
	}

	webhookSecret := randomHex(32)
	creds.WebhookSecret = webhookSecret

	credJSON, err := json.Marshal(creds)
	if err != nil {
		return domain.ExternalConnectionMeta{}, "", rerrors.Wrap(err, "marshal gitlab credentials")
	}

	var meta gitlabConnectionMeta

	err = json.Unmarshal(conn.Metadata, &meta)
	if err != nil {
		return domain.ExternalConnectionMeta{}, "", rerrors.Wrap(err, "error parsing gitlab meta")
	}

	meta.WebhookSecretSet = "true"

	metaJSON, err := json.Marshal(meta)
	if err != nil {
		return domain.ExternalConnectionMeta{}, "", rerrors.Wrap(err, "marshal gitlab meta")
	}

	conn.CredentialsJSON = json.RawMessage(credJSON)
	conn.Metadata = json.RawMessage(metaJSON)

	saved, err := s.connections.Upsert(ctx, conn)
	if err != nil {
		return domain.ExternalConnectionMeta{}, "", rerrors.Wrap(err, "save gitlab connection")
	}

	return toMeta(saved, meta.Username), webhookSecret, nil
}

// normalizeGitlabInstanceURL defaults to gitlab.com when blank, strips a trailing slash,
// and rejects schemes/hosts that would let the server be tricked into calling back into
// its own internal network (the server makes outbound requests to this host both for the
// validation ping below and for every later MoM tool call).
func normalizeGitlabInstanceURL(instanceUrl string) (string, error) {
	trimmed := strings.TrimRight(strings.TrimSpace(instanceUrl), "/")
	if trimmed == "" {
		trimmed = gitlabDefaultInstanceURL
	}

	parsed, err := url.Parse(trimmed)
	if err != nil {
		return "", rerrors.Wrap(user_errors.InvalidInstanceURL, "error parsing instance url")
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", user_errors.InvalidInstanceURL
	}

	if isLocalGitlabHost(parsed.Hostname()) {
		return "", user_errors.InvalidInstanceURL
	}

	return trimmed, nil
}

func isLocalGitlabHost(host string) bool {
	if host == "" || strings.EqualFold(host, "localhost") {
		return true
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return false
	}

	return ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsPrivate() ||
		ip.IsUnspecified()
}

// validateGitlabToken pings the instance's "who am I" endpoint to confirm the token and
// instance URL actually work together before anything is persisted.
func (s *Service) validateGitlabToken(ctx context.Context, instanceUrl, personalAccessToken string) (string, error) {
	reqUrl := instanceUrl + "/api/v4/user"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl, nil)
	if err != nil {
		return "", rerrors.Wrap(err, "error building gitlab validation request")
	}

	req.Header.Set("PRIVATE-TOKEN", personalAccessToken)

	resp, err := gitlabValidationClient.Do(req)
	if err != nil {
		return "", rerrors.Wrap(user_errors.GitlabValidationFailed, "error connecting to gitlab instance")
	}
	defer utils.CloseWithLog(resp.Body, "gitlab validation response")

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", rerrors.Wrap(user_errors.GitlabValidationFailed, "gitlab instance rejected the token")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", rerrors.Wrap(err, "error reading gitlab validation response")
	}

	var userInfo gitlabUserInfo

	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return "", rerrors.Wrap(err, "error parsing gitlab validation response")
	}

	return userInfo.Username, nil
}

func (s *Service) freshGoogleCreds(
	ctx context.Context,
	userUuid uuid.UUID,
) (domain.GoogleOAuthCredentials, error) {
	result, err := s.connections.GetByUserAndProvider(ctx, userUuid, domain.ProviderGoogleSheets)
	if err != nil {
		return domain.GoogleOAuthCredentials{}, rerrors.Wrap(err, "error getting google connection")
	}

	if !result.Valid {
		return domain.GoogleOAuthCredentials{}, user_errors.GoogleNotConnected
	}

	conn := result.V

	var creds domain.GoogleOAuthCredentials

	err = json.Unmarshal(conn.CredentialsJSON, &creds)
	if err != nil {
		return domain.GoogleOAuthCredentials{}, rerrors.Wrap(err, "error parsing google credentials")
	}

	if time.Until(creds.Expiry) < 5*time.Minute {
		existingToken := &oauth2.Token{
			AccessToken:  creds.AccessToken,
			RefreshToken: creds.RefreshToken,
			Expiry:       creds.Expiry,
			TokenType:    creds.TokenType,
		}

		tokenSource := s.oauthCfg.TokenSource(ctx, existingToken)

		freshToken, err := tokenSource.Token()
		if err != nil {
			return domain.GoogleOAuthCredentials{}, rerrors.Wrap(err, "error refreshing google token")
		}

		if freshToken.AccessToken != creds.AccessToken {
			creds.AccessToken = freshToken.AccessToken
			creds.Expiry = freshToken.Expiry
			creds.TokenType = freshToken.TokenType

			if freshToken.RefreshToken != "" {
				creds.RefreshToken = freshToken.RefreshToken
			}

			credJSON, err := json.Marshal(creds)
			if err != nil {
				return domain.GoogleOAuthCredentials{},
					rerrors.Wrap(err, "error marshaling refreshed credentials")
			}

			conn.CredentialsJSON = credJSON

			_, err = s.connections.Upsert(ctx, conn)
			if err != nil {
				return domain.GoogleOAuthCredentials{}, rerrors.Wrap(
					err,
					"error storing refreshed credentials",
				)
			}
		}
	}

	return creds, nil
}

func toMeta(conn domain.ExternalConnection, displayName string) domain.ExternalConnectionMeta {
	return domain.ExternalConnectionMeta{
		Uuid:         conn.Uuid,
		Provider:     conn.Provider,
		ProviderType: conn.ProviderType,
		DisplayName:  displayName,
		Metadata:     conn.Metadata,
		CreatedAt:    conn.CreatedAt,
		UpdatedAt:    conn.UpdatedAt,
	}
}

func extractDisplayName(conn domain.ExternalConnection) string {
	if conn.Metadata == nil {
		return ""
	}

	var meta domain.GoogleConnectionMeta

	err := json.Unmarshal(conn.Metadata, &meta)
	if err != nil {
		return ""
	}

	return meta.Email
}

func (s *Service) ListMailServerSuggestions(
	ctx context.Context,
	domainPrefix string,
) ([]domain.MailServerSuggestion, error) {
	suggestions, err := s.mailServerSuggestions.ListByDomain(ctx, domainPrefix)
	if err != nil {
		return nil, rerrors.Wrap(err, "list mail server suggestions")
	}

	return suggestions, nil
}

func randomHex(n int) string {
	b := make([]byte, n)

	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(b)
}
