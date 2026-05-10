package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/domain"
)

type VaultService interface {
	CreateVault(ctx context.Context, name string) error
	GetVault(ctx context.Context, name string) (domain.Vault, error)
	ListVaults(ctx context.Context) ([]domain.Vault, error)
	DeleteVault(ctx context.Context, name string) error
}

type UserService interface {
	CreateUser(ctx context.Context, username, password string, roles []string) error
	GetUser(ctx context.Context, username string) (domain.User, error)
	UpdateUser(ctx context.Context, username, password string, roles []string) error
	DeleteUser(ctx context.Context, username string) error
}

type AuthService interface {
	AuthWithToken(ctx context.Context, token string) (uuid.UUID, error)
}

type Service interface {
	AuthService() AuthService
	UserService() UserService
}
