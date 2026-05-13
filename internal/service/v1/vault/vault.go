package vault

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/clients/couchdb"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/repository/pg/tx_manager"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

type Service struct {
	vaultsRepo         repository.Vaults
	vaultMembersRepo   repository.VaultMembers
	couchAccountsRepo  repository.CouchAccounts
	couchInstancesRepo repository.CouchInstances

	txManager tx_manager.TxManager
}

func New(
	repo repository.Repo,
) *Service {
	return &Service{
		vaultsRepo:         repo.Vaults(),
		vaultMembersRepo:   repo.VaultMembers(),
		couchAccountsRepo:  repo.CouchAccounts(),
		couchInstancesRepo: repo.CouchInstances(),

		txManager: repo.TxManager(),
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

func (s *Service) CreateVault(ctx context.Context, vaultName string) (domain.Vault, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return domain.Vault{}, rerrors.Wrap(user_errors.Unauthenticated)
	}

	var vault domain.Vault

	err := s.txManager.Execute(
		func(tx *sql.Tx) error {
			couchInstanceRepo := s.couchInstancesRepo.WithTx(tx)
			couchAccountsRepo := s.couchAccountsRepo.WithTx(tx)
			vaultsRepo := s.vaultsRepo.WithTx(tx)
			vaultMembersRepo := s.vaultMembersRepo.WithTx(tx)

			instanceWithAccount, err := pickCouchInstance(ctx, couchInstanceRepo, uc.UserUuid)
			if err != nil {
				return rerrors.Wrap(err, "pick couch admin client")
			}

			couchClient := newCouchClient(instanceWithAccount.Instance)

			vault, err = s.ensureVaultExists(ctx, couchClient, uc, instanceWithAccount, vaultName, vaultsRepo)
			if err != nil {
				return rerrors.Wrap(err, "create couch database")
			}

			err = s.ensureCouchUserExists(ctx,
				couchClient,
				uc, &instanceWithAccount,
				couchAccountsRepo,
			)
			if err != nil {
				return rerrors.Wrap(err, "error ensuring couch user exists")
			}

			err = vaultMembersRepo.Add(ctx, vault.Uuid, uc.UserUuid, "owner")
			if err != nil {
				return rerrors.Wrap(err, "add vault owner member")
			}

			return nil
		})
	if err != nil {
		return domain.Vault{}, rerrors.Wrap(err, "error during tx")
	}

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

func (s *Service) AddMember(ctx context.Context, vaultID, targetUserUuid uuid.UUID) error {
	err := s.vaultMembersRepo.Add(ctx, vaultID, targetUserUuid, "member")
	if err != nil {
		return rerrors.Wrap(err, "add vault member")
	}

	return nil
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

func (s *Service) RemoveMember(ctx context.Context, vaultID, targetUserUuid uuid.UUID) error {
	err := s.vaultMembersRepo.Remove(ctx, vaultID, targetUserUuid)
	if err != nil {
		return rerrors.Wrap(err, "remove vault member")
	}

	return nil
}

func (s *Service) ensureCouchUserExists(ctx context.Context,
	adminClient *couchdb.Client,
	uc user_context.UserContext,
	instanceWithAccountPtr *domain.CouchInstanceWithAccount,
	couchAccountsRepo repository.CouchAccounts,
) (err error) {

	var couchPassword string
	if instanceWithAccountPtr.Account != nil {
		couchPassword = instanceWithAccountPtr.Account.CouchPassword
	} else {
		couchPassword, err = generatePassword()
		if err != nil {
			return rerrors.Wrap(err, "generate couch password")
		}
	}

	err = adminClient.CreateUser(ctx, uc.UserName, couchPassword, []string{})
	if err != nil && !errors.Is(err, user_errors.UserAlreadyExistInCouchDb) {
		return rerrors.Wrap(err, "create user in couch")
	}

	account, err := couchAccountsRepo.Upsert(ctx,
		uc.UserUuid,
		instanceWithAccountPtr.Instance.Uuid,
		uc.UserName,
		couchPassword,
	)
	if err != nil {
		return rerrors.Wrap(err, "error saving couch account")
	}

	instanceWithAccountPtr.Account = &account

	return nil
}

func (s *Service) ensureVaultExists(ctx context.Context,
	adminClient *couchdb.Client,
	uc user_context.UserContext,
	instanceWithAccount domain.CouchInstanceWithAccount,
	vaultName string,
	vaultsRepo repository.Vaults,
) (domain.Vault, error) {

	databaseName := sanitizeCouchDBName(uc.UserName + "-" + vaultName + "-vault")

	err := adminClient.CreateDatabase(ctx, databaseName)
	if err != nil {
		if !errors.Is(err, user_errors.CouchDbDatabaseAlreadyExists) {
			return domain.Vault{}, rerrors.Wrap(err, "create database in couch")
		}
	}

	vault, err := vaultsRepo.Upsert(ctx,
		uc.UserUuid,
		instanceWithAccount.Instance.Uuid,
		vaultName, databaseName,
		//TODO move to sql
		"ready")
	if err != nil {
		return domain.Vault{}, rerrors.Wrap(err, "create vault")
	}

	return vault, nil
}

func pickCouchInstance(
	ctx context.Context,
	couchInstance repository.CouchInstances,
	userUuid uuid.UUID,
) (domain.CouchInstanceWithAccount, error) {
	instance, err := couchInstance.Pick(ctx, userUuid)
	if err != nil {
		return domain.CouchInstanceWithAccount{}, rerrors.Wrap(err, "pick random couch instance")
	}

	return instance, nil
}

func newCouchClient(instance domain.CouchInstance) *couchdb.Client {
	cfg := couchdb.Config{
		BaseURL:  instance.Url,
		User:     instance.Username,
		Password: instance.Password,
	}

	adminClient := couchdb.New(cfg)

	return adminClient
}

func sanitizeCouchDBName(name string) string {
	name = strings.ToLower(name)

	name = strings.Replace(name, " ", "_", -1)

	return name
}
