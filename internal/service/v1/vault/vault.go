package vault

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/clients/couchdb"
	s3client "github.com/ruf-dev/artel/internal/clients/s3"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/repository"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"github.com/ruf-dev/artel/internal/repository/pg/tx_manager"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

type Service struct {
	vaultsRepo         repository.Vaults
	vaultMembersRepo   repository.VaultMembers
	vaultInvitesRepo   repository.VaultInvites
	couchAccountsRepo  repository.CouchAccounts
	couchInstancesRepo repository.CouchInstances
	s3InstancesRepo    repository.S3Instances

	txManager tx_manager.TxManager
}

func New(
	repo repository.Repo,
) *Service {
	return &Service{
		vaultsRepo:         repo.Vaults(),
		vaultMembersRepo:   repo.VaultMembers(),
		vaultInvitesRepo:   repo.VaultInvites(),
		couchAccountsRepo:  repo.CouchAccounts(),
		couchInstancesRepo: repo.CouchInstances(),
		s3InstancesRepo:    repo.S3Instances(),

		txManager: repo.TxManager(),
	}
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

			instanceWithAccount, err := couchInstanceRepo.RandomPick(ctx)
			if err != nil {
				if errors.Is(err, user_errors.NotFound) {
					return rerrors.Wrap(user_errors.NoCouchDbInstance)
				}

				return rerrors.Wrap(err, "pick random couch instance")
			}

			couchClient, err := newCouchClient(instanceWithAccount.Instance)
			if err != nil {
				return rerrors.Wrap(err, "init couch client")
			}

			liveSyncPassphrase, err := generatePassword()
			if err != nil {
				return rerrors.Wrap(err, "generate livesync passphrase")
			}

			vault, err = s.ensureVaultExists(
				ctx,
				couchClient,
				uc,
				instanceWithAccount,
				vaultName,
				liveSyncPassphrase,
				vaultsRepo,
			)
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

			err = couchClient.SetDatabaseSecurity(ctx, vault.CouchDBName, instanceWithAccount.Account.CouchUsername)
			if err != nil {
				return rerrors.Wrap(err, "error setting database security")
			}

			err = vaultMembersRepo.Add(ctx, vault.Uuid, uc.UserUuid, artel_q.VaultRoleOwner)
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

func (s *Service) AddMember(ctx context.Context, vaultID, targetUserUuid uuid.UUID, role artel_q.VaultRole) error {
	err := s.vaultMembersRepo.Add(ctx, vaultID, targetUserUuid, role)
	if err != nil {
		return rerrors.Wrap(err, "add vault member")
	}

	return nil
}

func (s *Service) ListMembers(ctx context.Context, vaultID uuid.UUID) ([]domain.VaultMemberInfo, error) {
	members, err := s.vaultMembersRepo.ListByVaultWithUsers(ctx, vaultID)
	if err != nil {
		return nil, rerrors.Wrap(err, "list vault members")
	}

	return members, nil
}

func (s *Service) CreateInviteLink(
	ctx context.Context,
	vaultID uuid.UUID,
	role artel_q.VaultRole,
) (domain.VaultInvite, error) {
	if role == artel_q.VaultRoleOwner {
		return domain.VaultInvite{}, rerrors.New("cannot create invite with owner role")
	}

	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return domain.VaultInvite{}, rerrors.Wrap(user_errors.Unauthenticated)
	}

	membership, err := s.vaultMembersRepo.Get(ctx, vaultID, uc.UserUuid)
	if err != nil {
		return domain.VaultInvite{}, rerrors.Wrap(err, "get vault membership")
	}

	if membership.Role != artel_q.VaultRoleOwner {
		return domain.VaultInvite{}, rerrors.Wrap(user_errors.NotVaultOwner)
	}

	token, err := generatePassword()
	if err != nil {
		return domain.VaultInvite{}, rerrors.Wrap(err, "generate invite token")
	}

	invite, err := s.vaultInvitesRepo.Create(ctx, vaultID, uc.UserUuid, role, token)
	if err != nil {
		return domain.VaultInvite{}, rerrors.Wrap(err, "create invite link")
	}

	return invite, nil
}

func (s *Service) ListInviteLinks(ctx context.Context, vaultID uuid.UUID) ([]domain.VaultInvite, error) {
	invites, err := s.vaultInvitesRepo.ListByVault(ctx, vaultID)
	if err != nil {
		return nil, rerrors.Wrap(err, "list invite links")
	}

	return invites, nil
}

func (s *Service) RevokeInviteLink(ctx context.Context, inviteID uuid.UUID) error {
	err := s.vaultInvitesRepo.Revoke(ctx, inviteID)
	if err != nil {
		return rerrors.Wrap(err, "revoke invite link")
	}

	return nil
}

func (s *Service) AcceptInvite(ctx context.Context, token string) error {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return rerrors.Wrap(user_errors.Unauthenticated)
	}

	invite, err := s.vaultInvitesRepo.GetByToken(ctx, token)
	if err != nil {
		return rerrors.Wrap(err, "get invite by token")
	}

	if invite.RevokedAt != nil {
		return rerrors.Wrap(user_errors.InviteLinkRevoked)
	}

	err = s.vaultMembersRepo.Add(ctx, invite.VaultUuid, uc.UserUuid, invite.Role)
	if err != nil {
		return rerrors.Wrap(err, "add vault member via invite")
	}

	return nil
}

func (s *Service) DeleteVault(ctx context.Context, vaultID uuid.UUID) error {
	vault, err := s.vaultsRepo.GetByID(ctx, vaultID)
	if err != nil {
		return rerrors.Wrap(err, "get vault by id")
	}

	instance, err := s.couchInstancesRepo.Get(ctx, vault.CouchInstanceUuid)
	if err != nil {
		return rerrors.Wrap(err, "get couch instance")
	}

	couchClient, err := newCouchClient(instance)
	if err != nil {
		return rerrors.Wrap(err, "init couch client")
	}

	err = couchClient.DeleteDatabase(ctx, vault.CouchDBName)
	if err != nil {
		return rerrors.Wrap(err, "delete couch database")
	}

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

func (s *Service) LinkS3Bucket(ctx context.Context, vaultID, s3InstanceID uuid.UUID, bucketName string) error {
	instance, err := s.s3InstancesRepo.Get(ctx, s3InstanceID)
	if err != nil {
		return rerrors.Wrap(err, "get s3 instance")
	}

	cfg := s3client.Config{
		Endpoint:  instance.Endpoint,
		Region:    instance.Region,
		AccessKey: instance.AccessKey,
		SecretKey: instance.SecretKey,
		UseSSL:    instance.UseSSL,
		PathStyle: instance.PathStyle,
	}

	client, err := s3client.New(cfg, bucketName)
	if err != nil {
		return rerrors.Wrap(err, "init s3 client")
	}

	err = client.EnsureBucket(ctx)
	if err != nil {
		return rerrors.Wrap(err, "ensure bucket exists")
	}

	err = s.vaultsRepo.LinkS3Bucket(ctx, vaultID, s3InstanceID, bucketName)
	if err != nil {
		return rerrors.Wrap(err, "link vault s3 bucket")
	}

	return nil
}

func (s *Service) UnlinkS3Bucket(ctx context.Context, vaultID uuid.UUID) error {
	err := s.vaultsRepo.UnlinkS3Bucket(ctx, vaultID)
	if err != nil {
		return rerrors.Wrap(err, "unlink vault s3 bucket")
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
	liveSyncPassphrase string,
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
		// TODO move to sql
		"ready",
		liveSyncPassphrase)
	if err != nil {
		return domain.Vault{}, rerrors.Wrap(err, "create vault")
	}

	return vault, nil
}

func newCouchClient(instance domain.CouchInstance) (*couchdb.Client, error) {
	cfg := couchdb.Config{
		BaseURL:  instance.Url,
		User:     instance.Username,
		Password: instance.Password,
	}

	adminClient, err := couchdb.New(cfg)
	if err != nil {
		return nil, rerrors.Wrap(err, "creating couch client")
	}

	return adminClient, nil
}

func sanitizeCouchDBName(name string) string {
	name = strings.ToLower(name)

	name = strings.ReplaceAll(name, " ", "_")

	return name
}

func generatePassword() (string, error) {
	b := make([]byte, 16)

	_, err := rand.Read(b)
	if err != nil {
		return "", rerrors.Wrap(err, "generate random bytes")
	}

	return hex.EncodeToString(b), nil
}
