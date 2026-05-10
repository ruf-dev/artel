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
}

type AuthService interface {
	Register(ctx context.Context, email, password string) (domain.User, error)
	Login(ctx context.Context, email, password string) (domain.Session, error)
	Logout(ctx context.Context, token string) error
	ValidateToken(ctx context.Context, token string) (uuid.UUID, error)
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
	DeleteCouchInstance(ctx context.Context, id string) error
}
