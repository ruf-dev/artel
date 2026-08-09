// Package postgresinstances is the admin CRUD service for the shared pool of Postgres servers
// (admin pool or BYOK) that per-vault databases are provisioned on — mirrors
// internal/service/v1/s3instances.
package postgresinstances

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/clients/vaultpg"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/utils"
)

type Service struct {
	postgresInstancesRepo repository.PostgresInstances
}

func New(repo repository.Repo) *Service {
	return &Service{
		postgresInstancesRepo: repo.PostgresInstances(),
	}
}

func (s *Service) RegisterPostgresInstance(
	ctx context.Context,
	host string, port int, adminDatabase, username, password, sslMode string,
) (string, error) {
	id, err := s.postgresInstancesRepo.Register(ctx, host, port, adminDatabase, username, []byte(password), sslMode)
	if err != nil {
		return "", rerrors.Wrap(err, "register postgres instance")
	}

	return id.String(), nil
}

func (s *Service) GetPostgresInstance(ctx context.Context, id string) (domain.PostgresInstance, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return domain.PostgresInstance{}, rerrors.Wrap(err, "parse uuid")
	}

	instance, err := s.postgresInstancesRepo.Get(ctx, uid)
	if err != nil {
		return domain.PostgresInstance{}, rerrors.Wrap(err, "get postgres instance")
	}

	return instance, nil
}

func (s *Service) ListPostgresInstances(ctx context.Context) ([]domain.PostgresInstance, error) {
	instances, err := s.postgresInstancesRepo.List(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "list postgres instances")
	}

	return instances, nil
}

func (s *Service) UpdatePostgresInstance(
	ctx context.Context,
	id, host string, port int, adminDatabase, username, password, sslMode string,
) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return rerrors.Wrap(err, "parse uuid")
	}

	err = s.postgresInstancesRepo.Update(ctx, uid, host, port, adminDatabase, username, []byte(password), sslMode)
	if err != nil {
		return rerrors.Wrap(err, "update postgres instance")
	}

	return nil
}

func (s *Service) DeletePostgresInstance(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return rerrors.Wrap(err, "parse uuid")
	}

	err = s.postgresInstancesRepo.Delete(ctx, uid)
	if err != nil {
		return rerrors.Wrap(err, "delete postgres instance")
	}

	return nil
}

func (s *Service) HasPostgresInstances(ctx context.Context) (bool, error) {
	exists, err := s.postgresInstancesRepo.Exists(ctx)
	if err != nil {
		return false, rerrors.Wrap(err, "check postgres instances exist")
	}

	return exists, nil
}

// TestPostgresInstance opens a connection to id's admin database and pings it, without mutating
// anything — mirrors s3instances.Service.TestS3Instance. vaultpg.AdminClient does not expose a
// raw ping/query method (only the DDL helpers EnsureRole/EnsureDatabase/DropDatabase/DropRole), so
// this opens its own *sql.DB from vaultpg.Config.DSN() rather than going through AdminClient.
func (s *Service) TestPostgresInstance(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return rerrors.Wrap(err, "parse uuid")
	}

	instance, err := s.postgresInstancesRepo.Get(ctx, uid)
	if err != nil {
		return rerrors.Wrap(err, "get postgres instance")
	}

	cfg := vaultpg.Config{
		Host:     instance.Host,
		Port:     instance.Port,
		Database: instance.AdminDatabase,
		User:     instance.Username,
		Password: instance.Password,
		SSLMode:  instance.SSLMode,
	}

	err = pingPostgres(ctx, cfg)
	if err != nil {
		return rerrors.Wrap(err, "test postgres connection")
	}

	return nil
}

// pingPostgres opens a throwaway connection against cfg and pings it — the minimal connectivity
// check shared by TestPostgresInstance here and externalconnections.Service.checkPostgresConnection,
// since vaultpg.AdminClient keeps its *sql.DB unexported and exposes no raw ping/query method.
func pingPostgres(ctx context.Context, cfg vaultpg.Config) error {
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return rerrors.Wrap(err, "error opening postgres connection")
	}

	defer utils.CloseWithLog(db, "postgres connection test")

	err = db.PingContext(ctx)
	if err != nil {
		return rerrors.Wrap(err, "error pinging postgres server")
	}

	return nil
}
