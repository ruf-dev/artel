package service

import "context"

type VaultService interface {
	CreateVault(ctx context.Context, name string) error
	GetVault(ctx context.Context, name string) (Vault, error)
	ListVaults(ctx context.Context) ([]Vault, error)
	DeleteVault(ctx context.Context, name string) error
}

type UserService interface {
	CreateUser(ctx context.Context, username, password string, roles []string) error
	GetUser(ctx context.Context, username string) (User, error)
	UpdateUser(ctx context.Context, username, password string, roles []string) error
	DeleteUser(ctx context.Context, username string) error
}

type Vault struct {
	Name  string
	DBURL string
}

type User struct {
	Username string
	Roles    []string
}