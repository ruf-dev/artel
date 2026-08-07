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

	anthropicClient "github.com/ruf-dev/artel/internal/clients/anthropic"
	"github.com/ruf-dev/artel/internal/clients/couchdb"
	"github.com/ruf-dev/artel/internal/clients/googleapi"
	"github.com/ruf-dev/artel/internal/clients/imap"
	openaiClient "github.com/ruf-dev/artel/internal/clients/openai"
	s3client "github.com/ruf-dev/artel/internal/clients/s3"
	"github.com/ruf-dev/artel/internal/clients/smtp"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/repository"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"github.com/ruf-dev/artel/internal/service"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"github.com/ruf-dev/artel/internal/utils"
)

const googleUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"

type googleUserInfo struct {
	Email string `json:"email"`
}

const gitlabDefaultInstanceURL = "https://gitlab.com"

// anthropicDefaultBaseUrl is used when the caller doesn't override the API host (e.g. for a
// proxy/regional endpoint).
const anthropicDefaultBaseUrl = "https://api.anthropic.com"

// openaiDefaultBaseUrl is used when the caller doesn't override the API host (e.g. for a
// proxy/regional endpoint or an OpenAI-compatible third-party provider).
const openaiDefaultBaseUrl = "https://api.openai.com/v1"

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

// telegramConnectionMeta is stored plaintext in metadata (non-sensitive display data).
type telegramConnectionMeta struct {
	BotUsername string `json:"bot_username"`
}

// anthropicConnectionMeta is stored plaintext in metadata (non-sensitive display data).
type anthropicConnectionMeta struct {
	KeyPreview      string `json:"key_preview"`
	DefaultModel    string `json:"default_model"`
	AvailableModels string `json:"available_models"` // comma-joined model IDs; genericDetails only keeps string values
	LastVerifiedAt  string `json:"last_verified_at"` // RFC3339
}

// openaiConnectionMeta is stored plaintext in metadata (non-sensitive display data).
type openaiConnectionMeta struct {
	KeyPreview      string `json:"key_preview"`
	DefaultModel    string `json:"default_model"`
	AvailableModels string `json:"available_models"` // comma-joined model IDs; genericDetails only keeps string values
	LastVerifiedAt  string `json:"last_verified_at"` // RFC3339
}

type Service struct {
	connections           repository.ExternalConnectionRepo
	pendingCodes          repository.PendingAuthCodes
	mcpSpreadsheets       repository.McpSpreadsheetsRepo
	mailServerSuggestions repository.MailServerSuggestions
	oauthCfg              *oauth2.Config
	mom                   service.MomService
	couchInstances        repository.CouchInstances
	s3Instances           repository.S3Instances
}

func New(
	connections repository.ExternalConnectionRepo,
	pendingCodes repository.PendingAuthCodes,
	mcpSpreadsheets repository.McpSpreadsheetsRepo,
	mailServerSuggestions repository.MailServerSuggestions,
	oauthCfg *oauth2.Config,
	mom service.MomService,
	couchInstances repository.CouchInstances,
	s3Instances repository.S3Instances,
) *Service {
	return &Service{
		connections:           connections,
		pendingCodes:          pendingCodes,
		mcpSpreadsheets:       mcpSpreadsheets,
		mailServerSuggestions: mailServerSuggestions,
		oauthCfg:              oauthCfg,
		mom:                   mom,
		couchInstances:        couchInstances,
		s3Instances:           s3Instances,
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

	err = s.cleanupOwnedInstance(ctx, uc.UserUuid, provider)
	if err != nil {
		return rerrors.Wrap(err, "error cleaning up owned storage instance")
	}

	return nil
}

// cleanupOwnedInstance removes the caller's owned couch/s3 pool row (migrations/066_instance_owner_scoping.sql)
// after disconnecting a BYOK couchdb/s3 provider, but only when no vault currently references it
// — see repository.CouchInstances/S3Instances.DeleteOwnedIfUnreferenced. No-op for every other
// provider.
//
// Known limitation, intentionally not handled further: if a vault already references the owned
// instance, the row is left in place so that vault keeps working. It remains resolvable via
// PickForUser for the same user until a new BYOK connection overwrites it (syncOwnedS3Instance /
// syncOwnedCouchInstance update the existing row in place rather than creating a second one).
func (s *Service) cleanupOwnedInstance(ctx context.Context, userUuid uuid.UUID, provider string) error {
	switch provider {
	case domain.ProviderCouchDB:
		err := s.couchInstances.DeleteOwnedIfUnreferenced(ctx, userUuid)
		if err != nil {
			return rerrors.Wrap(err, "delete owned couch instance if unreferenced")
		}
	case domain.ProviderS3:
		err := s.s3Instances.DeleteOwnedIfUnreferenced(ctx, userUuid)
		if err != nil {
			return rerrors.Wrap(err, "delete owned s3 instance if unreferenced")
		}
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

// AddTelegramConnection validates the bot token against Telegram, then persists it.
func (s *Service) AddTelegramConnection(
	ctx context.Context,
	botToken string,
) (domain.ExternalConnectionMeta, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return domain.ExternalConnectionMeta{}, user_errors.Unauthenticated
	}

	botUsername, err := s.validateTelegramToken(ctx, botToken)
	if err != nil {
		return domain.ExternalConnectionMeta{}, err
	}

	creds := domain.TelegramCredentials{
		BotToken: botToken,
	}

	credJSON, err := json.Marshal(creds)
	if err != nil {
		return domain.ExternalConnectionMeta{}, rerrors.Wrap(err, "error marshaling telegram credentials")
	}

	meta := telegramConnectionMeta{
		BotUsername: botUsername,
	}

	metaJSON, err := json.Marshal(meta)
	if err != nil {
		return domain.ExternalConnectionMeta{}, rerrors.Wrap(err, "error marshaling telegram meta")
	}

	conn := domain.ExternalConnection{
		UserUuid:        uc.UserUuid,
		Provider:        domain.ProviderTelegram,
		ProviderType:    artel_q.ExternalProviderTypeApiKey,
		CredentialsJSON: json.RawMessage(credJSON),
		Metadata:        json.RawMessage(metaJSON),
	}

	saved, err := s.connections.Upsert(ctx, conn)
	if err != nil {
		return domain.ExternalConnectionMeta{}, rerrors.Wrap(err, "error saving telegram connection")
	}

	return toMeta(saved, botUsername), nil
}

// CheckTelegramConnection pings Telegram with the given bot token without persisting anything,
// letting the caller confirm the token works before committing to AddTelegramConnection.
func (s *Service) CheckTelegramConnection(
	ctx context.Context,
	botToken string,
) (string, error) {
	_, ok := user_context.GetUserContext(ctx)
	if !ok {
		return "", user_errors.Unauthenticated
	}

	botUsername, err := s.validateTelegramToken(ctx, botToken)
	if err != nil {
		return "", err
	}

	return botUsername, nil
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
// before anything is persisted, by running the same trello.get_current_user MoM tool every later
// trello tool call goes through — just with a caller-supplied secrets map instead of one resolved
// from a persisted external_connections row, since that row doesn't exist yet at this point.
func (s *Service) validateTrelloCredentials(ctx context.Context, apiKey, apiToken string) (string, error) {
	secrets := map[string]interface{}{
		"api_key":   apiKey,
		"api_token": apiToken,
	}

	result, err := s.mom.ExecuteToolWithSecrets(ctx, "trello", "get_current_user", secrets, nil)
	if err != nil {
		return "", rerrors.Wrap(user_errors.TrelloValidationFailed, "error validating trello credentials")
	}

	var member trelloMemberResponse

	err = json.Unmarshal([]byte(result), &member)
	if err != nil {
		return "", rerrors.Wrap(err, "error parsing trello validation response")
	}

	return member.FullName, nil
}

// AddAnthropicConnection validates the API key against Anthropic, then persists it. A blank
// apiKey means "keep the current key" when editing an existing connection, mirroring
// AddEmailConnection's storedEmailPassword fallback.
func (s *Service) AddAnthropicConnection(
	ctx context.Context,
	apiKey, baseUrl, defaultModel string,
) (domain.ExternalConnectionMeta, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return domain.ExternalConnectionMeta{}, user_errors.Unauthenticated
	}

	if apiKey == "" {
		var err error

		apiKey, err = s.storedAnthropicApiKey(ctx, uc.UserUuid)
		if err != nil {
			return domain.ExternalConnectionMeta{}, err
		}
	}

	models, err := s.validateAnthropicKey(ctx, apiKey, baseUrl, defaultModel)
	if err != nil {
		return domain.ExternalConnectionMeta{}, err
	}

	resolvedDefaultModel := recommendedDefaultAnthropicModel(defaultModel, models)

	creds := domain.AnthropicKeyCredentials{
		ApiKey:  apiKey,
		BaseUrl: baseUrl,
	}

	credJSON, err := json.Marshal(creds)
	if err != nil {
		return domain.ExternalConnectionMeta{}, rerrors.Wrap(err, "marshal anthropic credentials")
	}

	meta := anthropicConnectionMeta{
		KeyPreview:      anthropicKeyPreview(apiKey),
		DefaultModel:    resolvedDefaultModel,
		AvailableModels: strings.Join(anthropicModelIds(models), ","),
		LastVerifiedAt:  time.Now().UTC().Format(time.RFC3339),
	}

	metaJSON, err := json.Marshal(meta)
	if err != nil {
		return domain.ExternalConnectionMeta{}, rerrors.Wrap(err, "marshal anthropic meta")
	}

	conn := domain.ExternalConnection{
		UserUuid:        uc.UserUuid,
		Provider:        domain.ProviderAnthropic,
		ProviderType:    artel_q.ExternalProviderTypeApiKey,
		CredentialsJSON: json.RawMessage(credJSON),
		Metadata:        json.RawMessage(metaJSON),
	}

	saved, err := s.connections.Upsert(ctx, conn)
	if err != nil {
		return domain.ExternalConnectionMeta{}, rerrors.Wrap(err, "save anthropic connection")
	}

	displayName := "Claude (Anthropic)"
	if resolvedDefaultModel != "" {
		displayName = resolvedDefaultModel
	}

	return toMeta(saved, displayName), nil
}

// CheckAnthropicConnection pings the provider with the given key without persisting anything,
// letting the caller confirm it works — and preview the available model catalog — before
// committing to AddAnthropicConnection.
func (s *Service) CheckAnthropicConnection(
	ctx context.Context,
	apiKey, baseUrl, defaultModel string,
) ([]anthropicClient.ModelInfo, string, error) {
	_, ok := user_context.GetUserContext(ctx)
	if !ok {
		return nil, "", user_errors.Unauthenticated
	}

	models, err := s.validateAnthropicKey(ctx, apiKey, baseUrl, defaultModel)
	if err != nil {
		return nil, "", err
	}

	recommendedDefault := recommendedDefaultAnthropicModel(defaultModel, models)

	return models, recommendedDefault, nil
}

// storedAnthropicApiKey resolves the API key of the user's existing anthropic connection, used to
// fall back when the caller submits a blank key (edit forms leave it blank to mean "keep the
// current key").
func (s *Service) storedAnthropicApiKey(ctx context.Context, userUuid uuid.UUID) (string, error) {
	existing, err := s.connections.GetByUserAndProvider(ctx, userUuid, domain.ProviderAnthropic)
	if err != nil {
		return "", rerrors.Wrap(err, "error loading existing anthropic connection")
	}

	if !existing.Valid {
		return "", user_errors.LlmKeyRequired
	}

	var creds domain.AnthropicKeyCredentials

	err = json.Unmarshal(existing.V.CredentialsJSON, &creds)
	if err != nil {
		return "", rerrors.Wrap(err, "error parsing existing anthropic credentials")
	}

	return creds.ApiKey, nil
}

// GetAnthropicApiKey returns userUuid's decrypted BYOK Anthropic api key. Returns
// user_errors.LlmKeyRequired if userUuid has no anthropic external_connections row — callers
// that want workbench-specific error semantics should translate that case themselves (see
// internal/service/v1/workbench).
func (s *Service) GetAnthropicApiKey(ctx context.Context, userUuid uuid.UUID) (string, error) {
	apiKey, err := s.storedAnthropicApiKey(ctx, userUuid)
	if err != nil {
		return "", err
	}

	return apiKey, nil
}

// AddGenericConnection stores arbitrary key/value credentials for a community-authored MCP
// connector, since provider is free TEXT and doesn't require a bespoke RPC per integration.
// Unlike the provider-specific Add*Connection methods, it does not validate the credentials
// against any external API — the caller (a MoM http-action tool) is responsible for that.
func (s *Service) AddGenericConnection(
	ctx context.Context,
	provider string,
	credentials map[string]string,
) (domain.ExternalConnectionMeta, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return domain.ExternalConnectionMeta{}, user_errors.Unauthenticated
	}

	if provider == "" {
		return domain.ExternalConnectionMeta{}, user_errors.GenericProviderRequired
	}

	credJSON, err := json.Marshal(credentials)
	if err != nil {
		return domain.ExternalConnectionMeta{}, rerrors.Wrap(err, "marshal generic credentials")
	}

	conn := domain.ExternalConnection{
		UserUuid:        uc.UserUuid,
		Provider:        provider,
		ProviderType:    artel_q.ExternalProviderTypeApiKey,
		CredentialsJSON: json.RawMessage(credJSON),
	}

	saved, err := s.connections.Upsert(ctx, conn)
	if err != nil {
		return domain.ExternalConnectionMeta{}, rerrors.Wrap(err, "save generic connection")
	}

	return toMeta(saved, provider), nil
}

// AddS3Connection validates connectivity against the S3-compatible endpoint, then persists it and
// syncs the caller's owned s3_instances pool row (migrations/066_instance_owner_scoping.sql) so
// new vaults resolve to it automatically instead of the shared admin pool — see
// vault.Service.LinkS3Bucket's PickForUser.
func (s *Service) AddS3Connection(
	ctx context.Context,
	endpoint, region, accessKey, secretKey string,
	useSSL, pathStyle bool,
) (domain.ExternalConnectionMeta, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return domain.ExternalConnectionMeta{}, user_errors.Unauthenticated
	}

	err := s.testS3Connection(ctx, endpoint, region, accessKey, secretKey, useSSL, pathStyle)
	if err != nil {
		return domain.ExternalConnectionMeta{}, err
	}

	creds := domain.S3KeyCredentials{
		Endpoint:  endpoint,
		Region:    region,
		AccessKey: accessKey,
		SecretKey: secretKey,
		UseSSL:    useSSL,
		PathStyle: pathStyle,
	}

	credJSON, err := json.Marshal(creds)
	if err != nil {
		return domain.ExternalConnectionMeta{}, rerrors.Wrap(err, "marshal s3 credentials")
	}

	conn := domain.ExternalConnection{
		UserUuid:        uc.UserUuid,
		Provider:        domain.ProviderS3,
		ProviderType:    artel_q.ExternalProviderTypeApiKey,
		CredentialsJSON: json.RawMessage(credJSON),
	}

	saved, err := s.connections.Upsert(ctx, conn)
	if err != nil {
		return domain.ExternalConnectionMeta{}, rerrors.Wrap(err, "save s3 connection")
	}

	err = s.syncOwnedS3Instance(ctx, uc.UserUuid, endpoint, region, accessKey, secretKey, useSSL, pathStyle)
	if err != nil {
		return domain.ExternalConnectionMeta{}, rerrors.Wrap(err, "sync owned s3 instance")
	}

	return toMeta(saved, endpoint), nil
}

// CheckS3Connection pings the S3-compatible endpoint with the given credentials without
// persisting anything, letting the caller confirm they work before committing to
// AddS3Connection.
func (s *Service) CheckS3Connection(
	ctx context.Context,
	endpoint, region, accessKey, secretKey string,
	useSSL, pathStyle bool,
) error {
	_, ok := user_context.GetUserContext(ctx)
	if !ok {
		return user_errors.Unauthenticated
	}

	err := s.testS3Connection(ctx, endpoint, region, accessKey, secretKey, useSSL, pathStyle)
	if err != nil {
		return err
	}

	return nil
}

// testS3Connection confirms endpoint/credentials work by listing buckets — a bucket-less
// reachability check (mirrors s3instances.Service.TestS3Instance, which does the same for
// admin-pool instances).
func (s *Service) testS3Connection(
	ctx context.Context,
	endpoint, region, accessKey, secretKey string,
	useSSL, pathStyle bool,
) error {
	cfg := s3client.Config{
		Endpoint:  endpoint,
		Region:    region,
		AccessKey: accessKey,
		SecretKey: secretKey,
		UseSSL:    useSSL,
		PathStyle: pathStyle,
	}

	err := s3client.TestConnection(ctx, cfg)
	if err != nil {
		return rerrors.Wrap(user_errors.S3ConnectionValidationFailed, "error connecting to s3-compatible endpoint")
	}

	return nil
}

// syncOwnedS3Instance mirrors a BYOK s3 connection into userUuid's owned s3_instances pool row,
// updating it in place if one already exists (so a second BYOK save doesn't create a duplicate
// pool row) or registering a new one otherwise.
func (s *Service) syncOwnedS3Instance(
	ctx context.Context,
	userUuid uuid.UUID,
	endpoint, region, accessKey, secretKey string,
	useSSL, pathStyle bool,
) error {
	owned, err := s.s3Instances.GetOwned(ctx, userUuid)
	if err != nil {
		return rerrors.Wrap(err, "get owned s3 instance")
	}

	if owned.Valid {
		err = s.s3Instances.Update(ctx, owned.V.Uuid, endpoint, region, useSSL, pathStyle, accessKey, []byte(secretKey))
		if err != nil {
			return rerrors.Wrap(err, "update owned s3 instance")
		}

		return nil
	}

	_, err = s.s3Instances.RegisterOwned(ctx, userUuid, endpoint, region, useSSL, pathStyle, accessKey, []byte(secretKey))
	if err != nil {
		return rerrors.Wrap(err, "register owned s3 instance")
	}

	return nil
}

// AddCouchDBConnection validates connectivity against the CouchDB server, then persists it and
// syncs the caller's owned couch_instances pool row (migrations/066_instance_owner_scoping.sql)
// so new vaults resolve to it automatically instead of the shared admin pool — see
// vault.Service.CreateVault's PickForUser.
func (s *Service) AddCouchDBConnection(
	ctx context.Context,
	url, username, password string,
) (domain.ExternalConnectionMeta, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return domain.ExternalConnectionMeta{}, user_errors.Unauthenticated
	}

	err := s.testCouchDBConnection(ctx, url, username, password)
	if err != nil {
		return domain.ExternalConnectionMeta{}, err
	}

	creds := domain.CouchDBKeyCredentials{
		URL:      url,
		Username: username,
		Password: password,
	}

	credJSON, err := json.Marshal(creds)
	if err != nil {
		return domain.ExternalConnectionMeta{}, rerrors.Wrap(err, "marshal couchdb credentials")
	}

	conn := domain.ExternalConnection{
		UserUuid:        uc.UserUuid,
		Provider:        domain.ProviderCouchDB,
		ProviderType:    artel_q.ExternalProviderTypeApiKey,
		CredentialsJSON: json.RawMessage(credJSON),
	}

	saved, err := s.connections.Upsert(ctx, conn)
	if err != nil {
		return domain.ExternalConnectionMeta{}, rerrors.Wrap(err, "save couchdb connection")
	}

	err = s.syncOwnedCouchInstance(ctx, uc.UserUuid, url, username, password)
	if err != nil {
		return domain.ExternalConnectionMeta{}, rerrors.Wrap(err, "sync owned couch instance")
	}

	return toMeta(saved, url), nil
}

// CheckCouchDBConnection pings the CouchDB server with the given credentials without persisting
// anything, letting the caller confirm they work before committing to AddCouchDBConnection.
func (s *Service) CheckCouchDBConnection(ctx context.Context, url, username, password string) error {
	_, ok := user_context.GetUserContext(ctx)
	if !ok {
		return user_errors.Unauthenticated
	}

	err := s.testCouchDBConnection(ctx, url, username, password)
	if err != nil {
		return err
	}

	return nil
}

// testCouchDBConnection confirms url/credentials work via GetSetupStatus — a read-only check
// (DBExists + a GET against /_cluster_setup) that never mutates the target server, unlike
// Client.Setup which provisions the _users/_replicator system databases.
func (s *Service) testCouchDBConnection(ctx context.Context, url, username, password string) error {
	cfg := couchdb.Config{
		BaseURL:  url,
		User:     username,
		Password: password,
	}

	client, err := couchdb.New(cfg)
	if err != nil {
		return rerrors.Wrap(user_errors.CouchDBConnectionValidationFailed, "error creating couchdb client")
	}

	_, err = client.GetSetupStatus(ctx)
	if err != nil {
		return rerrors.Wrap(user_errors.CouchDBConnectionValidationFailed, "error connecting to couchdb server")
	}

	return nil
}

// syncOwnedCouchInstance mirrors a BYOK couchdb connection into userUuid's owned couch_instances
// pool row, updating it in place if one already exists (so a second BYOK save doesn't create a
// duplicate pool row) or registering a new one otherwise.
func (s *Service) syncOwnedCouchInstance(ctx context.Context, userUuid uuid.UUID, url, username, password string) error {
	owned, err := s.couchInstances.GetOwned(ctx, userUuid)
	if err != nil {
		return rerrors.Wrap(err, "get owned couch instance")
	}

	if owned.Valid {
		err = s.couchInstances.Update(ctx, owned.V.Uuid, url, username, []byte(password))
		if err != nil {
			return rerrors.Wrap(err, "update owned couch instance")
		}

		return nil
	}

	_, err = s.couchInstances.RegisterOwned(ctx, userUuid, url, username, []byte(password))
	if err != nil {
		return rerrors.Wrap(err, "register owned couch instance")
	}

	return nil
}

// validateAnthropicKey resolves baseUrl to anthropicDefaultBaseUrl when blank, then confirms the
// key actually authenticates against the provider by listing its model catalog. ListModels is a
// zero-token metadata call, so it doubles as key validation before anything is persisted.
//
// Non-official Anthropic-compatible providers frequently mirror only the Messages API and don't
// implement GET /v1/models at all, which surfaces as a 404 here rather than an auth failure. In
// that case, when the caller supplied defaultModel, fall back to a minimal Ping against that
// model instead of failing outright — the model catalog stays empty (there's no listing endpoint
// to populate it from) but the key itself is still confirmed to work.
func (s *Service) validateAnthropicKey(
	ctx context.Context,
	apiKey, baseUrl, defaultModel string,
) ([]anthropicClient.ModelInfo, error) {
	client := anthropicClient.New(apiKey, resolveAnthropicBaseUrl(baseUrl))

	models, err := client.ListModels(ctx)
	if err == nil {
		return models, nil
	}

	statusCode, hasStatusCode := anthropicClient.StatusCode(err)

	switch {
	case hasStatusCode && (statusCode == http.StatusUnauthorized || statusCode == http.StatusForbidden):
		return nil, rerrors.Wrap(user_errors.LlmKeyRejected, "anthropic rejected the api key")
	case hasStatusCode && statusCode == http.StatusNotFound && defaultModel != "":
		pingErr := client.Ping(ctx, defaultModel)
		if pingErr != nil {
			return nil, rerrors.Wrap(user_errors.LlmKeyRejected, "error pinging anthropic-compatible endpoint with fallback model")
		}

		return nil, nil
	case hasStatusCode && statusCode == http.StatusNotFound:
		return nil, rerrors.Wrap(user_errors.LlmKeyModelListingUnsupported, "provider does not implement GET /v1/models")
	case hasStatusCode:
		return nil, rerrors.Wrap(user_errors.LlmKeyValidationFailed, "error listing anthropic models")
	default:
		return nil, rerrors.Wrap(user_errors.LlmKeyProviderUnreachable, "error connecting to anthropic-compatible endpoint")
	}
}

// resolveAnthropicBaseUrl returns baseUrl unchanged when set, or anthropicDefaultBaseUrl when
// blank.
func resolveAnthropicBaseUrl(baseUrl string) string {
	if baseUrl == "" {
		return anthropicDefaultBaseUrl
	}

	return baseUrl
}

// recommendedDefaultAnthropicModel returns the caller's chosen default model if set, otherwise
// falls back to the first entry in models. Anthropic's model-list endpoint returns models in
// reverse-chronological order (newest first), so the first entry is the newest model.
//
// TODO(byok): ModelInfo carries no tier/capability field (only Id, DisplayName,
// MaxInputTokens, MaxTokens), so there's no reliable way to prefer e.g. an Opus-tier model over
// a Haiku-tier one from the API response alone. "Newest first" is the simplest honest heuristic
// for this first pass — revisit if the SDK/API ever exposes a tier signal.
func recommendedDefaultAnthropicModel(defaultModel string, models []anthropicClient.ModelInfo) string {
	if defaultModel != "" {
		return defaultModel
	}

	if len(models) == 0 {
		return ""
	}

	return models[0].Id
}

// anthropicModelIds extracts the Id field from each model, preserving order.
func anthropicModelIds(models []anthropicClient.ModelInfo) []string {
	ids := make([]string, len(models))

	for i, m := range models {
		ids[i] = m.Id
	}

	return ids
}

// anthropicKeyPreview returns a display-safe suffix of apiKey (e.g. "...ab12") for
// anthropicConnectionMeta.KeyPreview — never the full key.
func anthropicKeyPreview(apiKey string) string {
	const previewLen = 4

	if len(apiKey) <= previewLen {
		return "..." + apiKey
	}

	return "..." + apiKey[len(apiKey)-previewLen:]
}

// AddOpenAIConnection validates the API key against OpenAI, then persists it. A blank
// apiKey means "keep the current key" when editing an existing connection, mirroring
// AddAnthropicConnection's storedAnthropicApiKey fallback.
func (s *Service) AddOpenAIConnection(
	ctx context.Context,
	apiKey, baseUrl, defaultModel string,
) (domain.ExternalConnectionMeta, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return domain.ExternalConnectionMeta{}, user_errors.Unauthenticated
	}

	if apiKey == "" {
		var err error

		apiKey, err = s.storedOpenAIApiKey(ctx, uc.UserUuid)
		if err != nil {
			return domain.ExternalConnectionMeta{}, err
		}
	}

	models, err := s.validateOpenAIKey(ctx, apiKey, baseUrl, defaultModel)
	if err != nil {
		return domain.ExternalConnectionMeta{}, err
	}

	resolvedDefaultModel := recommendedDefaultOpenAIModel(defaultModel, models)

	creds := domain.OpenAIKeyCredentials{
		ApiKey:  apiKey,
		BaseUrl: baseUrl,
	}

	credJSON, err := json.Marshal(creds)
	if err != nil {
		return domain.ExternalConnectionMeta{}, rerrors.Wrap(err, "marshal openai credentials")
	}

	meta := openaiConnectionMeta{
		KeyPreview:      openaiKeyPreview(apiKey),
		DefaultModel:    resolvedDefaultModel,
		AvailableModels: strings.Join(openaiModelIds(models), ","),
		LastVerifiedAt:  time.Now().UTC().Format(time.RFC3339),
	}

	metaJSON, err := json.Marshal(meta)
	if err != nil {
		return domain.ExternalConnectionMeta{}, rerrors.Wrap(err, "marshal openai meta")
	}

	conn := domain.ExternalConnection{
		UserUuid:        uc.UserUuid,
		Provider:        domain.ProviderOpenAI,
		ProviderType:    artel_q.ExternalProviderTypeApiKey,
		CredentialsJSON: json.RawMessage(credJSON),
		Metadata:        json.RawMessage(metaJSON),
	}

	saved, err := s.connections.Upsert(ctx, conn)
	if err != nil {
		return domain.ExternalConnectionMeta{}, rerrors.Wrap(err, "save openai connection")
	}

	displayName := "GPT (OpenAI)"
	if resolvedDefaultModel != "" {
		displayName = resolvedDefaultModel
	}

	return toMeta(saved, displayName), nil
}

// CheckOpenAIConnection pings the provider with the given key without persisting anything,
// letting the caller confirm it works — and preview the available model catalog — before
// committing to AddOpenAIConnection.
func (s *Service) CheckOpenAIConnection(
	ctx context.Context,
	apiKey, baseUrl, defaultModel string,
) ([]openaiClient.ModelInfo, string, error) {
	_, ok := user_context.GetUserContext(ctx)
	if !ok {
		return nil, "", user_errors.Unauthenticated
	}

	models, err := s.validateOpenAIKey(ctx, apiKey, baseUrl, defaultModel)
	if err != nil {
		return nil, "", err
	}

	recommendedDefault := recommendedDefaultOpenAIModel(defaultModel, models)

	return models, recommendedDefault, nil
}

// storedOpenAIApiKey resolves the API key of the user's existing openai connection, used to
// fall back when the caller submits a blank key (edit forms leave it blank to mean "keep the
// current key").
func (s *Service) storedOpenAIApiKey(ctx context.Context, userUuid uuid.UUID) (string, error) {
	existing, err := s.connections.GetByUserAndProvider(ctx, userUuid, domain.ProviderOpenAI)
	if err != nil {
		return "", rerrors.Wrap(err, "error loading existing openai connection")
	}

	if !existing.Valid {
		return "", user_errors.LlmKeyRequired
	}

	var creds domain.OpenAIKeyCredentials

	err = json.Unmarshal(existing.V.CredentialsJSON, &creds)
	if err != nil {
		return "", rerrors.Wrap(err, "error parsing existing openai credentials")
	}

	return creds.ApiKey, nil
}

// validateOpenAIKey resolves baseUrl to openaiDefaultBaseUrl when blank, then confirms the
// key actually authenticates against the provider by listing its model catalog. ListModels is a
// zero-token metadata call, so it doubles as key validation before anything is persisted.
//
// Non-official OpenAI-compatible providers frequently mirror only the Chat Completions API and
// don't implement GET /v1/models at all, which surfaces as a 404 here rather than an auth
// failure. In that case, when the caller supplied defaultModel, fall back to a minimal Ping
// against that model instead of failing outright — the model catalog stays empty (there's no
// listing endpoint to populate it from) but the key itself is still confirmed to work.
func (s *Service) validateOpenAIKey(
	ctx context.Context,
	apiKey, baseUrl, defaultModel string,
) ([]openaiClient.ModelInfo, error) {
	client := openaiClient.New(apiKey, resolveOpenAIBaseUrl(baseUrl))

	models, err := client.ListModels(ctx)
	if err == nil {
		return models, nil
	}

	statusCode, hasStatusCode := openaiClient.StatusCode(err)

	switch {
	case hasStatusCode && (statusCode == http.StatusUnauthorized || statusCode == http.StatusForbidden):
		return nil, rerrors.Wrap(user_errors.LlmKeyRejected, "openai rejected the api key")
	case hasStatusCode && statusCode == http.StatusNotFound && defaultModel != "":
		pingErr := client.Ping(ctx, defaultModel)
		if pingErr != nil {
			return nil, rerrors.Wrap(user_errors.LlmKeyRejected, "error pinging openai-compatible endpoint with fallback model")
		}

		return nil, nil
	case hasStatusCode && statusCode == http.StatusNotFound:
		return nil, rerrors.Wrap(user_errors.LlmKeyModelListingUnsupported, "provider does not implement GET /v1/models")
	case hasStatusCode:
		return nil, rerrors.Wrap(user_errors.LlmKeyValidationFailed, "error listing openai models")
	default:
		return nil, rerrors.Wrap(user_errors.LlmKeyProviderUnreachable, "error connecting to openai-compatible endpoint")
	}
}

// resolveOpenAIBaseUrl returns baseUrl unchanged when set, or openaiDefaultBaseUrl when
// blank.
func resolveOpenAIBaseUrl(baseUrl string) string {
	if baseUrl == "" {
		return openaiDefaultBaseUrl
	}

	return baseUrl
}

// openaiModelAllowlist is the preference order recommendedDefaultOpenAIModel falls back to when
// the caller doesn't set a default model. Ordered most- to least-preferred.
var openaiModelAllowlist = []string{"gpt-4o", "gpt-4o-mini", "gpt-4-turbo"}

// recommendedDefaultOpenAIModel returns the caller's chosen default model if set. Otherwise it
// picks the first entry of openaiModelAllowlist present in models.
//
// This differs from recommendedDefaultAnthropicModel's "first entry" heuristic deliberately:
// OpenAI's GET /v1/models response is not reverse-chronological and mixes in non-chat models
// (embeddings, moderation, TTS, etc.), so "first entry" is not a safe signal of "newest/best chat
// model" the way it is for Anthropic. If none of the allowlisted models are present, "" is
// returned rather than guessing at an arbitrary entry.
func recommendedDefaultOpenAIModel(defaultModel string, models []openaiClient.ModelInfo) string {
	if defaultModel != "" {
		return defaultModel
	}

	available := make(map[string]bool, len(models))
	for _, m := range models {
		available[m.Id] = true
	}

	for _, candidate := range openaiModelAllowlist {
		if available[candidate] {
			return candidate
		}
	}

	return ""
}

// openaiModelIds extracts the Id field from each model, preserving order.
func openaiModelIds(models []openaiClient.ModelInfo) []string {
	ids := make([]string, len(models))

	for i, m := range models {
		ids[i] = m.Id
	}

	return ids
}

// openaiKeyPreview returns a display-safe suffix of apiKey (e.g. "...ab12") for
// openaiConnectionMeta.KeyPreview — never the full key.
func openaiKeyPreview(apiKey string) string {
	const previewLen = 4

	if len(apiKey) <= previewLen {
		return "..." + apiKey
	}

	return "..." + apiKey[len(apiKey)-previewLen:]
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
// instance URL actually work together before anything is persisted, by running the same
// gitlab.get_current_user MoM tool every later gitlab tool call goes through — just with a
// caller-supplied secrets map instead of one resolved from a persisted external_connections row,
// since that row doesn't exist yet at this point.
func (s *Service) validateGitlabToken(ctx context.Context, instanceUrl, personalAccessToken string) (string, error) {
	secrets := map[string]interface{}{
		"instance_url":          instanceUrl,
		"personal_access_token": personalAccessToken,
	}

	result, err := s.mom.ExecuteToolWithSecrets(ctx, "gitlab", "get_current_user", secrets, nil)
	if err != nil {
		return "", rerrors.Wrap(user_errors.GitlabValidationFailed, "error validating gitlab token")
	}

	var userInfo gitlabUserInfo

	err = json.Unmarshal([]byte(result), &userInfo)
	if err != nil {
		return "", rerrors.Wrap(err, "error parsing gitlab validation response")
	}

	return userInfo.Username, nil
}

// telegramGetMeResponse mirrors Telegram's real getMe response shape:
// {"ok": true, "result": {"id": 123, "is_bot": true, "first_name": "...", "username": "mybot"}}.
type telegramGetMeResponse struct {
	Ok     bool `json:"ok"`
	Result struct {
		Username string `json:"username"`
	} `json:"result"`
}

// validateTelegramToken pings Telegram's getMe endpoint to confirm the bot token actually works
// before anything is persisted, by running the same telegram.get_me MoM tool every later telegram
// tool call goes through — just with a caller-supplied secrets map instead of one resolved from a
// persisted external_connections row, since that row doesn't exist yet at this point.
func (s *Service) validateTelegramToken(ctx context.Context, botToken string) (string, error) {
	secrets := map[string]interface{}{
		"bot_token": botToken,
	}

	result, err := s.mom.ExecuteToolWithSecrets(ctx, "telegram", "get_me", secrets, nil)
	if err != nil {
		return "", rerrors.Wrap(user_errors.TelegramValidationFailed, "error validating telegram bot token")
	}

	var resp telegramGetMeResponse

	err = json.Unmarshal([]byte(result), &resp)
	if err != nil {
		return "", rerrors.Wrap(err, "error parsing telegram validation response")
	}

	if !resp.Ok {
		return "", rerrors.Wrap(user_errors.TelegramValidationFailed, "error validating telegram bot token")
	}

	return resp.Result.Username, nil
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
