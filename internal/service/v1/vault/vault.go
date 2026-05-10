package vault

import (
	"context"

	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/clients/couchdb"
	"github.com/ruf-dev/artel/internal/service"
)

type Service struct {
	db *couchdb.Client
}

func New(db *couchdb.Client) *Service {
	return &Service{db: db}
}

func (s *Service) CreateVault(ctx context.Context, name string) error {
	err := s.db.CreateDatabase(ctx, name)
	if err != nil {
		return rerrors.Wrap(err, "create vault")
	}
	return nil
}

func (s *Service) GetVault(ctx context.Context, name string) (service.Vault, error) {
	v := service.Vault{
		Name:  name,
		DBURL: s.db.DatabaseURL(name),
	}
	return v, nil
}

func (s *Service) ListVaults(ctx context.Context) ([]service.Vault, error) {
	names, err := s.db.ListDatabases(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "list vaults")
	}

	vaults := make([]service.Vault, 0, len(names))
	for _, name := range names {
		v := service.Vault{
			Name:  name,
			DBURL: s.db.DatabaseURL(name),
		}
		vaults = append(vaults, v)
	}
	return vaults, nil
}

func (s *Service) DeleteVault(ctx context.Context, name string) error {
	err := s.db.DeleteDatabase(ctx, name)
	if err != nil {
		return rerrors.Wrap(err, "delete vault")
	}
	return nil
}