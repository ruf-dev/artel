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

// TODO func should be idempotent
func (s *Service) CreateVault(ctx context.Context, name string) (domain.Vault, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return domain.Vault{}, rerrors.Wrap(user_errors.Unauthenticated)
	}

	//TODO pick by user or random by default.
	//TODO move to separate func and return client
	instance, err := s.couchInstancesRepo.RandomPick(ctx)
	if err != nil {
		return domain.Vault{}, rerrors.Wrap(err, "pick random couch instance")
	}

	couchCfg := couchdb.Config{
		BaseURL:  instance.Url,
		User:     instance.Username,
		Password: instance.Password,
	}
	adminClient := couchdb.New(couchCfg)

	//TODO make idempotent
	err = s.ensureCouchAccount(ctx, uc.UserUuid, instance, adminClient)
	if err != nil {
		return domain.Vault{}, rerrors.Wrap(err, "ensure couch account")
	}

	//TODO make idempotent
	err = s.createCouchDatabase(ctx, adminClient, name)
	if err != nil {
		return domain.Vault{}, rerrors.Wrap(err, "create couch database")
	}

	vault, err := s.vaultsRepo.Create(ctx, uc.UserUuid, instance.Uuid, name, name)
	if err != nil {
		return domain.Vault{}, rerrors.Wrap(err, "create vault")
	}
	//TODO move to top. Create metadata first, add status 'provisioning' and change status by the end of this func.
	err = s.vaultMembersRepo.Add(ctx, vault.Uuid, uc.UserUuid, "owner")
	if err != nil {
		return domain.Vault{}, rerrors.Wrap(err, "add vault owner member")
	}

	vault.CouchDBURL = adminClient.DatabaseURL(name)

	return vault, nil
}

func (s *Service) ensureCouchAccount(ctx context.Context, userUuid uuid.UUID, instance domain.CouchInstance, adminClient *couchdb.Client) error {
	account, err := s.couchAccountsRepo.GetByUserAndInstance(ctx, userUuid, instance.Uuid)
	if err == nil {
		_ = account
		return nil
	}

	err = s.createCouchUser(ctx, userUuid, instance, adminClient)
	if err != nil {
		return rerrors.Wrap(err, "create couch user")
	}

	return nil
}

func (s *Service) createCouchUser(ctx context.Context, userUuid uuid.UUID, instance domain.CouchInstance, adminClient *couchdb.Client) error {
	couchUsername := "artel-" + userUuid.String()

	couchPassword, err := generatePassword()
	if err != nil {
		return rerrors.Wrap(err, "generate couch password")
	}

	err = adminClient.CreateUser(ctx, couchUsername, couchPassword, []string{})
	if err != nil {
		return rerrors.Wrap(err, "create user in couch")
	}

	_, err = s.couchAccountsRepo.Create(ctx, userUuid, instance.Uuid, couchUsername, []byte(couchPassword))
	if err != nil {
		return rerrors.Wrap(err, "store couch account")
	}

	return nil
}

func (s *Service) createCouchDatabase(ctx context.Context, adminClient *couchdb.Client, databaseName string) error {
	err := adminClient.CreateDatabase(ctx, databaseName)
	if err != nil {
		return rerrors.Wrap(err, "create database in couch")
	}
	return nil
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
