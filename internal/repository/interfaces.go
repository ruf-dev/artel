package repository

import (
	"context"
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
	McpKeyRepository() McpKeyRepository

	TxManager() tx_manager.TxManager
}

type Users interface {
	Create(ctx context.Context, email, passwordHash string) (domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByTelegramId(ctx context.Context, telegramId string) (domain.User, error)
	UpsertByTelegramId(ctx context.Context, telegramId string, username string) (domain.User, error)
}

type Vaults interface {
	Create(ctx context.Context, userID, couchInstanceID uuid.UUID, name, couchDBName, status string) (domain.Vault, error)
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
}

type CouchAccounts interface {
	Create(ctx context.Context, userID, instanceID uuid.UUID, username string, passwordPlain string) (domain.CouchAccount, error)
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
	Delete(ctx context.Context, id uuid.UUID) error

	WithTx(tx sqldb.DB) CouchInstances
}

type McpKeyRepository interface {
	CreateMcpKey(ctx context.Context, vaultID, userID uuid.UUID, name string, keyHash []byte, keyPreview string) (domain.McpKey, error)
	ListMcpKeysByVault(ctx context.Context, vaultID uuid.UUID) ([]domain.McpKey, error)
	GetMcpKeyByID(ctx context.Context, id uuid.UUID) (domain.McpKey, error)
	ListActiveMcpKeys(ctx context.Context, vaultID uuid.UUID) ([]domain.McpKey, error)
	RevokeMcpKey(ctx context.Context, id uuid.UUID) error
}
