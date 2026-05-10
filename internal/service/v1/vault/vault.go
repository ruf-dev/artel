package vault

import (
	"context"

	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
)

type Service struct {
	vaultRepo repository.Vaults
}

func New(db repository.Repo) *Service {
	return &Service{
		vaultRepo: db.Vaults(),
	}
}

func (s *Service) CreateVault(ctx context.Context, name string) error {
	err := s.vaultRepo.Create(ctx, name)
	if err != nil {
		return rerrors.Wrap(err, "create vault")
	}
	return nil
}

func (s *Service) GetVault(ctx context.Context, name string) (domain.Vault, error) {
	v := domain.Vault{
		Name:        name,
		CouchDBName: name,
		CouchDBURL:  s.db.DatabaseURL(name),
	}
	return v, nil
}

func (s *Service) ListVaults(ctx context.Context) ([]domain.Vault, error) {
	names, err := s.db.ListDatabases(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "list vaults")
	}

	vaults := make([]domain.Vault, 0, len(names))
	for _, name := range names {
		v := domain.Vault{
			Name:        name,
			CouchDBName: name,
			CouchDBURL:  s.db.DatabaseURL(name),
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
