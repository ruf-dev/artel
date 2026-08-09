package vault

import (
	"context"
	"database/sql"
	"strings"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"github.com/ruf-dev/artel/internal/utils"
)

// EnablePostgresDatabase provisions a Postgres database+role for vaultID on a pool instance
// resolved via postgresInstancesRepo.PickForUser (prefers the caller's BYOK postgres connection,
// falling back to the shared admin pool). The vaultPostgresRepo row is created up front with
// status=provisioning (see repository.VaultPostgresDatabases.Create) before the live provisioning
// calls run, so a mid-flight failure leaves a retryable/inspectable 'error' row rather than
// silently returning an error with no persisted trace — mirrors the Workbench
// provisioning->ready/error lifecycle.
func (s *Service) EnablePostgresDatabase(ctx context.Context, vaultID uuid.UUID) (domain.VaultPostgresDatabase, error) {
	uc, err := s.requireVaultMember(ctx, vaultID)
	if err != nil {
		return domain.VaultPostgresDatabase{}, err
	}

	existing, err := s.vaultPostgresRepo.GetByVaultID(ctx, vaultID)
	if err != nil {
		return domain.VaultPostgresDatabase{}, rerrors.Wrap(err, "error getting vault postgres database")
	}

	if existing.Valid {
		return domain.VaultPostgresDatabase{}, rerrors.Wrap(user_errors.PostgresDatabaseAlreadyEnabled)
	}

	err = s.subscriptions.CheckFeature(ctx, uc.UserUuid, domain.FeaturePostgres)
	if err != nil {
		return domain.VaultPostgresDatabase{}, err
	}

	// NOTE(track-b): there is no dedicated user_errors entry for "no postgres pool instance
	// available" (the couch/s3 equivalent is user_errors.NoCouchDbInstance) — user_errors.go is
	// owned by a parallel task, so this wraps the raw not-found error instead of introducing one
	// in-place. Track D/user_errors follow-up: add a NoPostgresInstance-style error mirroring
	// NoCouchDbInstance.
	instance, err := s.postgresInstancesRepo.PickForUser(ctx, uc.UserUuid)
	if err != nil {
		return domain.VaultPostgresDatabase{}, rerrors.Wrap(err, "error picking postgres instance for user")
	}

	dbName := vaultPostgresDatabaseName(vaultID)
	roleUsername := vaultPostgresRoleUsername(vaultID)

	rolePassword, err := generatePassword()
	if err != nil {
		return domain.VaultPostgresDatabase{}, rerrors.Wrap(err, "error generating postgres role password")
	}

	vaultPostgres, err := s.vaultPostgresRepo.Create(ctx, vaultID, instance.Uuid, dbName, roleUsername, []byte(rolePassword))
	if err != nil {
		return domain.VaultPostgresDatabase{}, rerrors.Wrap(err, "error creating vault postgres database record")
	}

	provisionErr := s.provisionPostgresDatabase(ctx, instance, dbName, roleUsername, rolePassword)
	if provisionErr != nil {
		markErr := s.vaultPostgresRepo.MarkError(ctx, vaultID, provisionErr.Error())
		if markErr != nil {
			return domain.VaultPostgresDatabase{}, rerrors.Wrap(markErr, "error marking vault postgres database as errored")
		}

		return domain.VaultPostgresDatabase{}, rerrors.Wrap(provisionErr, "error provisioning vault postgres database")
	}

	err = s.vaultPostgresRepo.MarkReady(ctx, vaultID)
	if err != nil {
		return domain.VaultPostgresDatabase{}, rerrors.Wrap(err, "error marking vault postgres database as ready")
	}

	vaultPostgres.Status = domain.VaultPostgresStatusReady

	return vaultPostgres, nil
}

// provisionPostgresDatabase runs the admin-side DDL (CREATE ROLE, CREATE DATABASE) against
// instance for a newly-created vault_postgres_databases row.
func (s *Service) provisionPostgresDatabase(
	ctx context.Context,
	instance domain.PostgresInstance,
	dbName, roleUsername, rolePassword string,
) error {
	client, err := s.newAdminClient(instance)
	if err != nil {
		return rerrors.Wrap(err, "error creating postgres admin client")
	}

	defer utils.CloseWithLog(client, "vaultpg admin client")

	err = client.EnsureRole(ctx, roleUsername, rolePassword)
	if err != nil {
		return rerrors.Wrap(err, "error ensuring postgres role")
	}

	err = client.EnsureDatabase(ctx, dbName, roleUsername)
	if err != nil {
		return rerrors.Wrap(err, "error ensuring postgres database")
	}

	return nil
}

// GetPostgresDatabase returns vaultID's Postgres database row, Valid: false when Postgres has
// never been enabled for this vault.
func (s *Service) GetPostgresDatabase(ctx context.Context, vaultID uuid.UUID) (sql.Null[domain.VaultPostgresDatabase], error) {
	_, err := s.requireVaultMember(ctx, vaultID)
	if err != nil {
		return sql.Null[domain.VaultPostgresDatabase]{}, err
	}

	result, err := s.vaultPostgresRepo.GetByVaultID(ctx, vaultID)
	if err != nil {
		return sql.Null[domain.VaultPostgresDatabase]{}, rerrors.Wrap(err, "error getting vault postgres database")
	}

	return result, nil
}

// DisablePostgresDatabase best-effort drops vaultID's provisioned database+role (tolerating a
// server that's gone or unreachable — DropDatabase/DropRole are already IF EXISTS/idempotent, see
// vaultpg.AdminClient), then deletes the row. No-op if Postgres was never enabled for this vault.
// Only the vault owner may disable it, mirroring PublishVault/UnpublishVault's owner check for
// other destructive per-vault settings.
func (s *Service) DisablePostgresDatabase(ctx context.Context, vaultID uuid.UUID) error {
	err := s.requireVaultOwner(ctx, vaultID)
	if err != nil {
		return err
	}

	existing, err := s.vaultPostgresRepo.GetByVaultID(ctx, vaultID)
	if err != nil {
		return rerrors.Wrap(err, "error getting vault postgres database")
	}

	if !existing.Valid {
		return nil
	}

	instance, err := s.postgresInstancesRepo.Get(ctx, existing.V.PostgresInstanceUuid)
	if err != nil {
		return rerrors.Wrap(err, "error getting postgres instance")
	}

	client, err := s.newAdminClient(instance)
	if err != nil {
		return rerrors.Wrap(err, "error creating postgres admin client")
	}

	utils.CallWithLog(func() error {
		return client.DropDatabase(ctx, existing.V.DatabaseName)
	}, "drop vault postgres database")

	utils.CallWithLog(func() error {
		return client.DropRole(ctx, existing.V.RoleUsername)
	}, "drop vault postgres role")

	utils.CloseWithLog(client, "vaultpg admin client")

	err = s.vaultPostgresRepo.Delete(ctx, vaultID)
	if err != nil {
		return rerrors.Wrap(err, "error deleting vault postgres database")
	}

	return nil
}

// requireVaultMember resolves the caller and confirms they belong to vaultID, mirroring
// notes.Service.resolveVault's membership check — enabling/inspecting a vault's Postgres database
// is available to any member, not just the owner.
func (s *Service) requireVaultMember(ctx context.Context, vaultID uuid.UUID) (user_context.UserContext, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return user_context.UserContext{}, rerrors.Wrap(user_errors.Unauthenticated)
	}

	_, err := s.vaultMembersRepo.Get(ctx, vaultID, uc.UserUuid)
	if err != nil {
		return user_context.UserContext{}, rerrors.Wrap(err, "error getting vault membership")
	}

	return uc, nil
}

// requireVaultOwner resolves the caller and confirms they own vaultID, mirroring the owner check
// in PublishVault/UnpublishVault/SetUseCouchDBForBinaries.
func (s *Service) requireVaultOwner(ctx context.Context, vaultID uuid.UUID) error {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return rerrors.Wrap(user_errors.Unauthenticated)
	}

	membership, err := s.vaultMembersRepo.Get(ctx, vaultID, uc.UserUuid)
	if err != nil {
		return rerrors.Wrap(err, "error getting vault membership")
	}

	if membership.Role != artel_q.VaultRoleOwner {
		return rerrors.Wrap(user_errors.NotVaultOwner)
	}

	return nil
}

// vaultPostgresDatabaseName / vaultPostgresRoleUsername derive deterministic Postgres identifiers
// from vaultID, sidestepping the need to sanitize a user-supplied name — a UUID with dashes
// stripped is 32 characters, well under Postgres's 63-byte identifier limit even with the prefix.
func vaultPostgresDatabaseName(vaultID uuid.UUID) string {
	return "vault_" + strings.ReplaceAll(vaultID.String(), "-", "")
}

func vaultPostgresRoleUsername(vaultID uuid.UUID) string {
	return "vaultpg_" + strings.ReplaceAll(vaultID.String(), "-", "")
}
