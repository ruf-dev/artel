package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/ruf-dev/artel/internal/domain"
)

type Service interface {
	AuthService() AuthService
	VaultService() VaultService
	CouchInstanceService() CouchInstanceService
	McpService() McpService
	EmailService() EmailService
	SubscriptionService() SubscriptionService
	PromptService() PromptService
}

type AuthService interface {
	Register(ctx context.Context, email, password string) (domain.User, error)
	Login(ctx context.Context, email, password string) (domain.Session, error)
	Logout(ctx context.Context, token string) error
	ValidateToken(ctx context.Context, token string) (uuid.UUID, error)
	LoginViaTelegram(ctx context.Context, idToken string) (domain.Session, error)
	GetMe(ctx context.Context, userUuid uuid.UUID) (domain.User, domain.UserPermissions, error)
	CheckIsAdmin(ctx context.Context, userUuid uuid.UUID) error
}

type VaultService interface {
	CreateVault(ctx context.Context, name string) (domain.Vault, error)
	GetVault(ctx context.Context, vaultID uuid.UUID) (domain.Vault, error)
	ListVaults(ctx context.Context) ([]domain.Vault, error)
	DeleteVault(ctx context.Context, vaultID uuid.UUID) error
	AddMember(ctx context.Context, vaultID, targetUserUuid uuid.UUID) error
	RemoveMember(ctx context.Context, vaultID, targetUserUuid uuid.UUID) error
}

type CouchInstanceService interface {
	RegisterCouchInstance(ctx context.Context, url, username, password string) (string, error)
	GetCouchInstance(ctx context.Context, id string) (domain.CouchInstance, error)
	ListCouchInstances(ctx context.Context) ([]domain.CouchInstance, error)
	UpdateCouchInstance(ctx context.Context, id, url, username, password string) error
	DeleteCouchInstance(ctx context.Context, id string) error
}

type McpService interface {
	// CreateKey generates a new bearer token, stores it hashed, returns the raw token once.
	CreateKey(ctx context.Context, vaultID uuid.UUID, name string) (rawToken string, key domain.McpKey, err error)
	ListKeys(ctx context.Context, vaultID uuid.UUID) ([]domain.McpKey, error)
	RevokeKey(ctx context.Context, keyID uuid.UUID) error
	// ResolveKey validates the raw bearer token and returns vault+couch context.
	ResolveKey(ctx context.Context, rawToken string) (domain.McpKeyContext, error)
}

type SubscriptionService interface {
	CheckActive(ctx context.Context, userUuid uuid.UUID) error
}

type ListPromptsParams struct {
	Ids      []string
	Page     uint32
	PageSize uint32
}

type PromptService interface {
	ListPrompts(ctx context.Context, params ListPromptsParams) ([]domain.Prompt, int64, error)
}

type EmailService interface {
	// Account management — called from gRPC transport only.
	AddAccount(ctx context.Context, account domain.EmailAccount) (domain.EmailAccount, error)
	ListAccounts(ctx context.Context, userUuid uuid.UUID) ([]domain.EmailAccount, error)
	DeleteAccount(ctx context.Context, accountUuid uuid.UUID) error

	// Email operations — called from MCP tools.
	ListFolders(ctx context.Context, accountUuid uuid.UUID) ([]string, error)
	ListEmails(ctx context.Context, accountUuid uuid.UUID, limit int) ([]domain.EmailMeta, error)
	ReadEmail(ctx context.Context, accountUuid uuid.UUID, id string) (domain.EmailMessage, error)
	SendEmail(ctx context.Context, accountUuid uuid.UUID, to, subject, body string) error
}
