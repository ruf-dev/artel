package vault

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/clients/couchdb"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

type Service struct {
	vaultsRepo         repository.Vaults
	vaultMembersRepo   repository.VaultMembers
	couchAccountsRepo  repository.CouchAccounts
	couchInstancesRepo repository.CouchInstances
}

func New(
	vaultsRepo repository.Vaults,
	vaultMembersRepo repository.VaultMembers,
	couchAccountsRepo repository.CouchAccounts,
	couchInstancesRepo repository.CouchInstances,
) *Service {
	return &Service{
		vaultsRepo:         vaultsRepo,
		vaultMembersRepo:   vaultMembersRepo,
		couchAccountsRepo:  couchAccountsRepo,
		couchInstancesRepo: couchInstancesRepo,
	}
}

func generatePassword() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", rerrors.Wrap(err, "generate random bytes")
	}
	return hex.EncodeToString(b), nil
}

func (s *Service) CreateVault(ctx context.Context, name string) (domain.Vault, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return domain.Vault{}, rerrors.Wrap(user_errors.Unauthenticated)
	}

	instance, err := s.couchInstancesRepo.RandomPick(ctx)
	if err != nil {
		return domain.Vault{}, rerrors.Wrap(err, "pick random couch instance")
	}

	account, err := s.couchAccountsRepo.GetByUserAndInstance(ctx, uc.UserUuid, instance.Uuid)
	if err != nil {
		couchUsername := "artel-" + uc.UserUuid.String()

		couchPassword, err := generatePassword()
		if err != nil {
			return domain.Vault{}, rerrors.Wrap(err, "generate couch password")
		}

		adminCfg := couchdb.Config{
			BaseURL:  instance.Url,
			User:     instance.Username,
			Password: instance.Password,
		}
		adminClient := couchdb.New(adminCfg)

		err = adminClient.CreateUser(ctx, couchUsername, couchPassword, []string{})
		if err != nil {
			return domain.Vault{}, rerrors.Wrap(err, "create couch user")
		}

		account, err = s.couchAccountsRepo.Create(ctx, uc.UserUuid, instance.Uuid, couchUsername, []byte(couchPassword))
		if err != nil {
			return domain.Vault{}, rerrors.Wrap(err, "store couch account")
		}
	}

	_ = account

	couchDBName := name

	adminCfg := couchdb.Config{
		BaseURL:  instance.Url,
		User:     instance.Username,
		Password: instance.Password,
	}
	adminClient := couchdb.New(adminCfg)

	err = adminClient.CreateDatabase(ctx, couchDBName)
	if err != nil {
		return domain.Vault{}, rerrors.Wrap(err, "create couch database")
	}

	vault, err := s.vaultsRepo.Create(ctx, uc.UserUuid, instance.Uuid, name, couchDBName)
	if err != nil {
		return domain.Vault{}, rerrors.Wrap(err, "create vault")
	}

	err = s.vaultMembersRepo.Add(ctx, vault.Uuid, uc.UserUuid, "owner")
	if err != nil {
		return domain.Vault{}, rerrors.Wrap(err, "add vault owner member")
	}

	vault.CouchDBURL = adminClient.DatabaseURL(couchDBName)

	return vault, nil
}

func (s *Service) GetVault(ctx context.Context, vaultID uuid.UUID) (domain.Vault, error) {
	vault, err := s.vaultsRepo.GetByID(ctx, vaultID)
	if err != nil {
		return domain.Vault{}, rerrors.Wrap(err, "get vault by id")
	}

	instance, err := s.couchInstancesRepo.Get(ctx, vault.CouchInstanceUuid)
	if err != nil {
		return domain.Vault{}, rerrors.Wrap(err, "get couch instance")
	}

	vault.CouchDBURL = instance.Url + "/" + vault.CouchDBName

	return vault, nil
}

func (s *Service) ListVaults(ctx context.Context) ([]domain.Vault, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return nil, rerrors.Wrap(user_errors.Unauthenticated)
	}

	vaults, err := s.vaultsRepo.ListByMembership(ctx, uc.UserUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "list vaults by membership")
	}

	for i, vault := range vaults {
		instance, err := s.couchInstancesRepo.Get(ctx, vault.CouchInstanceUuid)
		if err != nil {
			return nil, rerrors.Wrap(err, "get couch instance")
		}
		vaults[i].CouchDBURL = instance.Url + "/" + vault.CouchDBName
	}

	return vaults, nil
}

func (s *Service) DeleteVault(ctx context.Context, vaultID uuid.UUID) error {
	vault, err := s.vaultsRepo.GetByID(ctx, vaultID)
	if err != nil {
		return rerrors.Wrap(err, "get vault by id")
	}

	_ = vault

	err = s.vaultsRepo.Delete(ctx, vaultID)
	if err != nil {
		return rerrors.Wrap(err, "delete vault")
	}

	return nil
}

func (s *Service) AddMember(ctx context.Context, vaultID, targetUserUuid uuid.UUID) error {
	err := s.vaultMembersRepo.Add(ctx, vaultID, targetUserUuid, "member")
	if err != nil {
		return rerrors.Wrap(err, "add vault member")
	}

	return nil
}

func (s *Service) RemoveMember(ctx context.Context, vaultID, targetUserUuid uuid.UUID) error {
	err := s.vaultMembersRepo.Remove(ctx, vaultID, targetUserUuid)
	if err != nil {
		return rerrors.Wrap(err, "remove vault member")
	}

	return nil
}
