package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/ruf-dev/artel/internal/domain"
)

type Repo interface {
	Users() Users
	Vaults() Vaults
	VaultMembers() VaultMembers
	Sessions() Sessions
	Subscriptions() Subscriptions
	CouchAccounts() CouchAccounts
	CouchInstances() CouchInstances
}

type Users interface {
	Create(ctx context.Context, email, passwordHash string) (domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type Vaults interface {
	Create(ctx context.Context, userID, couchInstanceID uuid.UUID, name, couchDBName string) (domain.Vault, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.Vault, error)
	ListByMembership(ctx context.Context, userID uuid.UUID) ([]domain.Vault, error)
	Delete(ctx context.Context, vaultID uuid.UUID) error
}

type VaultMembers interface {
	Add(ctx context.Context, vaultID, userID uuid.UUID, role string) error
	Remove(ctx context.Context, vaultID, userID uuid.UUID) error
	Get(ctx context.Context, vaultID, userID uuid.UUID) (domain.VaultMember, error)
	ListByVault(ctx context.Context, vaultID uuid.UUID) ([]domain.VaultMember, error)
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
	Create(ctx context.Context, userID, instanceID uuid.UUID, username string, passwordPlain []byte) (domain.CouchAccount, error)
	GetByUserAndInstance(ctx context.Context, userID, instanceID uuid.UUID) (domain.CouchAccount, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.CouchAccount, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CouchInstances interface {
	Register(ctx context.Context, url, username string, passwordPlain []byte) (uuid.UUID, error)
	Get(ctx context.Context, id uuid.UUID) (domain.CouchInstance, error)
	RandomPick(ctx context.Context) (domain.CouchInstance, error)
	List(ctx context.Context) ([]domain.CouchInstance, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
