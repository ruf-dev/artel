package vault

import (
	"context"

	"github.com/ruf-dev/artel/internal/clients/couchdb"
)

type Service struct {
	db *couchdb.Client
}

func New(db *couchdb.Client) *Service {
	return &Service{db: db}
}

func (s *Service) CreateVault(ctx context.Context, name string) error {
	return s.db.CreateDatabase(ctx, name)
}

func (s *Service) DeleteVault(ctx context.Context, name string) error {
	return s.db.DeleteDatabase(ctx, name)
}
