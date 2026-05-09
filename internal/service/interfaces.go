package service

import "context"

type VaultService interface {
	CreateVault(ctx context.Context, name string) error
	DeleteVault(ctx context.Context, name string) error
}
