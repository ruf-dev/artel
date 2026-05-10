package vault

import (
	"context"
	"fmt"

	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/clients/couchdb"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

type Service struct {
	vaultsRepo     repository.Vaults
	couchCredsRepo repository.CouchCredentials

	pool *couchdb.Pool
}

func New(
	vaultsRepo repository.Vaults,
	couchCredsRepo repository.CouchCredentials,
	pool *couchdb.Pool) *Service {
	return &Service{
		vaultsRepo:     vaultsRepo,
		couchCredsRepo: couchCredsRepo,
		pool:           pool,
	}
}

func (s *Service) CreateVault(ctx context.Context, name string) error {
	couchDBName := name
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return rerrors.Wrap(user_errors.Unauthenticated)
	}

	v, err := s.vaultsRepo.Create(ctx, uc.UserUuid, name, couchDBName)
	if err != nil {
		return rerrors.Wrap(err, "create vault metadata")
	}

	couchClient, err := s.pool.GetDefault(ctx)
	if err != nil {
		return rerrors.Wrap(err, "get default couchdb client")
	}

	err = couchClient.CreateDatabase(ctx, couchDBName)
	if err != nil {
		return rerrors.Wrap(err, "create couchdb database")
	}

	cfg := couchClient.Config()
	err = s.couchCredsRepo.Store(ctx, v.Uuid, cfg.BaseURL, cfg.User, []byte(cfg.Password))
	if err != nil {
		return rerrors.Wrap(err, "store couch credentials")
	}

	return nil
}

func (s *Service) GetVault(ctx context.Context, name string) (domain.Vault, error) {
	vault, err := s.vaultsRepo.GetByName(ctx, name)
	if err != nil {
		return domain.Vault{}, rerrors.Wrap(err, "get vault by name")
	}

	cred, err := s.couchCredsRepo.Load(ctx, vault.Uuid)
	if err != nil {
		return domain.Vault{}, rerrors.Wrap(err, "load couch credentials")
	}

	vault.CouchDBURL = fmt.Sprintf("%s/%s", cred.Host, vault.CouchDBName)
	return vault, nil
}

func (s *Service) ListVaults(ctx context.Context) ([]domain.Vault, error) {
	vaults, err := s.vaultsRepo.ListAll(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "list vaults")
	}

	for i, vault := range vaults {
		cred, err := s.couchCredsRepo.Load(ctx, vault.Uuid)
		if err != nil {
			return nil, rerrors.Wrap(err, "load couch credentials")
		}
		vaults[i].CouchDBURL = fmt.Sprintf("%s/%s", cred.Host, vault.CouchDBName)
	}

	return vaults, nil
}

func (s *Service) DeleteVault(ctx context.Context, name string) error {
	vault, err := s.vaultsRepo.GetByName(ctx, name)
	if err != nil {
		return rerrors.Wrap(err, "get vault by name")
	}

	instanceClient, err := s.pool.Get(ctx, vault.Uuid)
	if err != nil {
		return rerrors.Wrap(err, "get couchdb client from pool")
	}

	err = instanceClient.DeleteDatabase(ctx, vault.CouchDBName)
	if err != nil {
		return rerrors.Wrap(err, "delete couchdb database")
	}

	err = s.couchCredsRepo.Delete(ctx, vault.Uuid)
	if err != nil {
		return rerrors.Wrap(err, "delete couch credentials")
	}

	err = s.vaultsRepo.Delete(ctx, vault.Uuid)
	if err != nil {
		return rerrors.Wrap(err, "delete vault metadata")
	}

	return nil
}
