package service

import (
	"context"

	"github.com/ruf-dev/artel/internal/storage/couchdb"
)

type VaultService interface {
	CreateVault(ctx context.Context, name string) error
	DeleteVault(ctx context.Context, name string) error
}

type vaultService struct {
	db *couchdb.Client
}

func NewVaultService(db *couchdb.Client) VaultService {
	return &vaultService{db: db}
}

func (vs *vaultService) CreateVault(ctx context.Context, name string) error {
	return vs.db.CreateDatabase(ctx, name)
}

func (vs *vaultService) DeleteVault(ctx context.Context, name string) error {
	return vs.db.DeleteDatabase(ctx, name)
}
