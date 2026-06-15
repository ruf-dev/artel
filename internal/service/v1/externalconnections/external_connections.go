package externalconnections

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"time"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"
	"golang.org/x/oauth2"

	"github.com/ruf-dev/artel/internal/clients/googleapi"
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

type Service struct {
	connections     repository.ExternalConnectionRepo
	pendingCodes    repository.PendingAuthCodes
	mcpSpreadsheets repository.McpSpreadsheetsRepo
	oauthCfg        *oauth2.Config
}

func New(
	connections repository.ExternalConnectionRepo,
	pendingCodes repository.PendingAuthCodes,
	mcpSpreadsheets repository.McpSpreadsheetsRepo,
	oauthCfg *oauth2.Config,
) *Service {
	return &Service{
		connections:     connections,
		pendingCodes:    pendingCodes,
		mcpSpreadsheets: mcpSpreadsheets,
		oauthCfg:        oauthCfg,
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

func (s *Service) HandleGoogleOAuthCallback(ctx context.Context, code string, state string) (domain.ExternalConnectionMeta, error) {
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

	creds, conn, err := s.freshGoogleCreds(ctx, uc.UserUuid)
	if err != nil {
		return nil, err
	}

	_ = conn
	client := googleapi.New(ctx, creds, s.oauthCfg)
	return client, nil
}

func (s *Service) GetPickerToken(ctx context.Context) (string, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return "", user_errors.Unauthenticated
	}

	creds, _, err := s.freshGoogleCreds(ctx, uc.UserUuid)
	if err != nil {
		return "", err
	}

	return creds.AccessToken, nil
}

func (s *Service) AddSpreadsheet(ctx context.Context, spreadsheetId string, name string) (domain.McpSpreadsheet, error) {
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

func (s *Service) AddEmailConnection(ctx context.Context, email, imapHost string, imapPort int, smtpHost string, smtpPort int, password string) (domain.ExternalConnectionMeta, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return domain.ExternalConnectionMeta{}, user_errors.Unauthenticated
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

// freshGoogleCreds loads the stored Google credentials for the given user, refreshes the
// access token if it is within 5 minutes of expiry, persists any refreshed token, and
// returns the up-to-date credentials together with the connection row.
func (s *Service) freshGoogleCreds(ctx context.Context, userUuid uuid.UUID) (domain.GoogleOAuthCredentials, domain.ExternalConnection, error) {
	result, err := s.connections.GetByUserAndProvider(ctx, userUuid, domain.ProviderGoogleSheets)
	if err != nil {
		return domain.GoogleOAuthCredentials{}, domain.ExternalConnection{}, rerrors.Wrap(err, "error getting google connection")
	}
	if !result.Valid {
		return domain.GoogleOAuthCredentials{}, domain.ExternalConnection{}, user_errors.GoogleNotConnected
	}

	conn := result.V

	var creds domain.GoogleOAuthCredentials
	err = json.Unmarshal(conn.CredentialsJSON, &creds)
	if err != nil {
		return domain.GoogleOAuthCredentials{}, domain.ExternalConnection{}, rerrors.Wrap(err, "error parsing google credentials")
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
			return domain.GoogleOAuthCredentials{}, domain.ExternalConnection{}, rerrors.Wrap(err, "error refreshing google token")
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
				return domain.GoogleOAuthCredentials{}, domain.ExternalConnection{}, rerrors.Wrap(err, "error marshaling refreshed credentials")
			}

			conn.CredentialsJSON = json.RawMessage(credJSON)
			conn, err = s.connections.Upsert(ctx, conn)
			if err != nil {
				return domain.GoogleOAuthCredentials{}, domain.ExternalConnection{}, rerrors.Wrap(err, "error storing refreshed credentials")
			}
		}
	}

	return creds, conn, nil
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

func randomHex(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
