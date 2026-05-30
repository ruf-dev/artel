package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"

	"github.com/ruf-dev/artel/internal/clients/sqldb"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository/pg/tx_manager"
)

type Repo interface {
	Users() Users
	Vaults() Vaults
	VaultMembers() VaultMembers
	Sessions() Sessions
	Subscriptions() Subscriptions
	CouchAccounts() CouchAccounts
	CouchInstances() CouchInstances
	UserPermissions() UserPermissionsRepo
	McpKeyRepository() McpKeyRepository
	PendingAuthCodes() PendingAuthCodes
	EmailAccounts() EmailAccounts
	Prompts() Prompts

	TxManager() tx_manager.TxManager
}

type ListPromptsParams struct {
	Ids    []string
	Limit  uint32
	Offset uint32
}

type Prompts interface {
	List(ctx context.Context, params ListPromptsParams) ([]domain.Prompt, int64, error)
}

type Users interface {
	Create(ctx context.Context, email, passwordHash string) (domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)

	GetByTelegramId(ctx context.Context, telegramId string) (sql.Null[domain.User], error)
	CreateByUsername(ctx context.Context, username string) (domain.User, error)
	UpsertTelegramIdentity(ctx context.Context, identity domain.TelegramIdentity) error
	GetTelegramPhotoUrl(ctx context.Context, userUuid uuid.UUID) (string, error)

	Delete(ctx context.Context, id uuid.UUID) error

	WithTx(tx *sql.Tx) Users
}

type Vaults interface {
	Upsert(ctx context.Context, userID, couchInstanceID uuid.UUID, name, couchDBName, status string) (domain.Vault, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.Vault, error)
	GetByNameAndUser(ctx context.Context, userID uuid.UUID, name string) (domain.Vault, error)
	UpdateStatus(ctx context.Context, vaultID uuid.UUID, status string) error
	ListByMembership(ctx context.Context, userID uuid.UUID) ([]domain.Vault, error)
	Delete(ctx context.Context, vaultID uuid.UUID) error

	WithTx(tx sqldb.DB) Vaults
}

type VaultMembers interface {
	Add(ctx context.Context, vaultID, userID uuid.UUID, role string) error
	Remove(ctx context.Context, vaultID, userID uuid.UUID) error
	Get(ctx context.Context, vaultID, userID uuid.UUID) (domain.VaultMember, error)
	ListByVault(ctx context.Context, vaultID uuid.UUID) ([]domain.VaultMember, error)

	WithTx(tx sqldb.DB) VaultMembers
}

type Sessions interface {
	Create(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) (domain.Session, error)
	GetByToken(ctx context.Context, token string) (domain.Session, error)
	Delete(ctx context.Context, token string) error
}

type Subscriptions interface {
	Upsert(ctx context.Context, userID uuid.UUID, active bool) (domain.Subscription, error)
	GetByUser(ctx context.Context, userID uuid.UUID) (domain.Subscription, error)
	CreateDefault(ctx context.Context, userID uuid.UUID) error

	WithTx(tx *sql.Tx) Subscriptions
}

type CouchAccounts interface {
	Upsert(ctx context.Context, userID, instanceID uuid.UUID, username string, passwordPlain string) (domain.CouchAccount, error)
	GetByUserAndInstance(ctx context.Context, userID, instanceID uuid.UUID) (domain.CouchAccount, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.CouchAccount, error)
	Delete(ctx context.Context, id uuid.UUID) error

	WithTx(tx sqldb.DB) CouchAccounts
}

type CouchInstances interface {
	Register(ctx context.Context, url, username string, passwordPlain []byte) (uuid.UUID, error)
	Get(ctx context.Context, id uuid.UUID) (domain.CouchInstance, error)
	Pick(ctx context.Context, id uuid.UUID) (domain.CouchInstanceWithAccount, error)
	List(ctx context.Context) ([]domain.CouchInstance, error)
	Update(ctx context.Context, id uuid.UUID, url, username string, passwordPlain []byte) error
	Delete(ctx context.Context, id uuid.UUID) error

	WithTx(tx sqldb.DB) CouchInstances
}

type UserPermissionsRepo interface {
	Get(ctx context.Context, userUuid uuid.UUID) (domain.UserPermissions, error)
	Upsert(ctx context.Context, userUuid uuid.UUID, isAdmin bool, hasEmails bool) (domain.UserPermissions, error)
	CreateDefault(ctx context.Context, userUuid uuid.UUID) error

	WithTx(tx *sql.Tx) UserPermissionsRepo
}

type McpKeyRepository interface {
	CreateMcpKey(ctx context.Context, vaultID, userID, keyID uuid.UUID, name string, keyHash []byte, keyPreview string) (domain.McpKey, error)
	ListMcpKeysByVault(ctx context.Context, vaultID uuid.UUID) ([]domain.McpKey, error)
	GetMcpKeyByID(ctx context.Context, id uuid.UUID) (domain.McpKey, error)
	ListActiveMcpKeys(ctx context.Context, vaultID uuid.UUID) ([]domain.McpKey, error)
	RevokeMcpKey(ctx context.Context, id uuid.UUID) error
}

type PendingAuthCodes interface {
	Create(ctx context.Context, code, rawToken, codeChallenge, redirectUri, clientId string, expiresAt time.Time) error
	Get(ctx context.Context, code string) (domain.PendingAuthCode, error)
	Delete(ctx context.Context, code string) error
	DeleteExpired(ctx context.Context) error
}

type EmailAccounts interface {
	Insert(ctx context.Context, account domain.EmailAccount) (domain.EmailAccount, error)
	GetByUuid(ctx context.Context, id uuid.UUID) (domain.EmailAccount, error)
	ListByUser(ctx context.Context, userUuid uuid.UUID) ([]domain.EmailAccount, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
