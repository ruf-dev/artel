package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Vault struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        string
	CouchDBName string
	CreatedAt   time.Time
}

type Subscription struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CouchCred struct {
	ID          uuid.UUID
	VaultID     uuid.UUID
	Host        string
	Username    string
	PasswordEnc []byte
	CreatedAt   time.Time
}

type Users interface {
	Create(ctx context.Context, email string) (User, error)
	GetByID(ctx context.Context, id uuid.UUID) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
}

type Vaults interface {
	Create(ctx context.Context, userID uuid.UUID, name, couchDBName string) (Vault, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]Vault, error)
}

type Subscriptions interface {
	Upsert(ctx context.Context, userID uuid.UUID, active bool) (Subscription, error)
	GetByUser(ctx context.Context, userID uuid.UUID) (Subscription, error)
}

type CouchCredentials interface {
	Store(ctx context.Context, vaultID uuid.UUID, host, username string, passwordPlain []byte) error
	Load(ctx context.Context, vaultID uuid.UUID) (CouchCred, error)
}
