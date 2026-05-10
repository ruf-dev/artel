package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/ruf-dev/artel/internal/domain"
)

type Repo interface {
	Users() Users
	Vaults() Vaults
	Subscriptions() Subscriptions
	CouchCredentials() CouchCredentials
}

type Users interface {
	Create(ctx context.Context, email string) (domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
}

type Vaults interface {
	Create(ctx context.Context, userID uuid.UUID, name, couchDBName string) (domain.Vault, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.Vault, error)
}

type Subscriptions interface {
	Upsert(ctx context.Context, userID uuid.UUID, active bool) (domain.Subscription, error)
	GetByUser(ctx context.Context, userID uuid.UUID) (domain.Subscription, error)
}

type CouchCredentials interface {
	Store(ctx context.Context, vaultID uuid.UUID, host, username string, passwordPlain []byte) error
	Load(ctx context.Context, vaultID uuid.UUID) (domain.CouchCred, error)
}
