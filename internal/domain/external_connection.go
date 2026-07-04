package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
)

const (
	ProviderGoogleSheets = "google_sheets"
	ProviderTrello       = "trello"
	ProviderMiro         = "miro"
	ProviderEmail        = "email"
	ProviderGitlab       = "gitlab"
)

type ExternalConnection struct {
	Uuid            uuid.UUID
	UserUuid        uuid.UUID
	Provider        string
	ProviderType    artel_q.ExternalProviderType
	CredentialsJSON json.RawMessage
	Metadata        json.RawMessage
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type ExternalConnectionMeta struct {
	Uuid         uuid.UUID
	Provider     string
	ProviderType artel_q.ExternalProviderType
	DisplayName  string
	Metadata     json.RawMessage
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// GoogleOAuthCredentials is stored encrypted in credentials_enc.
type GoogleOAuthCredentials struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Expiry       time.Time `json:"expiry"`
	Scope        string    `json:"scope"`
	TokenType    string    `json:"token_type"`
}

// GoogleConnectionMeta is stored plaintext in metadata (non-sensitive display data).
type GoogleConnectionMeta struct {
	Email  string `json:"email"`
	Scopes string `json:"scopes"`
}

// APIKeyCredentials is stored encrypted in credentials_enc for api_key providers.
type APIKeyCredentials struct {
	APIKey   string `json:"api_key"`
	APIToken string `json:"api_token"`
}

// EmailCredentials is stored encrypted in credentials_enc for email (IMAP/SMTP) providers.
type EmailCredentials struct {
	ImapHost string `json:"imap_host"`
	ImapPort int    `json:"imap_port"`
	SmtpHost string `json:"smtp_host"`
	SmtpPort int    `json:"smtp_port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// GitlabCredentials is stored encrypted in credentials_enc for the gitlab provider.
type GitlabCredentials struct {
	PersonalAccessToken string `json:"personal_access_token"`
	InstanceUrl         string `json:"instance_url"`
	WebhookSecret       string `json:"webhook_secret,omitempty"`
}
