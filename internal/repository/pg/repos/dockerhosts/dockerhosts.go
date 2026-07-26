// Package dockerhosts is the pure-DB layer for the admin-managed pool of Docker daemons backing
// per-vault workbench containers — see docs/workbench/02_docker_topology.md. Unlike
// couchinstances/s3instances, docker hosts carry no credentials, so this repo has no encryptor.
package dockerhosts

import (
	"context"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/clients/sqldb"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"github.com/ruf-dev/artel/internal/repository/pg/pg_err"
	"go.redsock.ru/rerrors"
)

type Repo struct {
	q *artel_q.Queries
}

func New(db sqldb.DB) *Repo {
	return &Repo{
		q: artel_q.New(db),
	}
}

func (r *Repo) Register(ctx context.Context, url string) (uuid.UUID, error) {
	id, err := r.q.RegisterDockerHost(ctx, url)
	if err != nil {
		return uuid.UUID{}, pg_err.UnwrapPgErr(err)
	}

	return id, nil
}

func (r *Repo) Get(ctx context.Context, id uuid.UUID) (domain.DockerHost, error) {
	row, err := r.q.GetDockerHost(ctx, id)
	if err != nil {
		return domain.DockerHost{}, pg_err.UnwrapPgErr(err)
	}

	host := domain.DockerHost{
		Uuid:      row.ID,
		Url:       row.Url,
		CreatedAt: row.CreatedAt,
	}

	return host, nil
}

func (r *Repo) List(ctx context.Context) ([]domain.DockerHost, error) {
	rows, err := r.q.ListDockerHosts(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "list docker hosts")
	}

	hosts := make([]domain.DockerHost, len(rows))
	for i, row := range rows {
		hosts[i] = domain.DockerHost{
			Uuid:      row.ID,
			Url:       row.Url,
			CreatedAt: row.CreatedAt,
		}
	}

	return hosts, nil
}

func (r *Repo) Update(ctx context.Context, id uuid.UUID, url string) error {
	params := artel_q.UpdateDockerHostParams{
		ID:  id,
		Url: url,
	}

	err := r.q.UpdateDockerHost(ctx, params)
	if err != nil {
		return rerrors.Wrap(err, "update docker host")
	}

	return nil
}

func (r *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeleteDockerHost(ctx, id)
	if err != nil {
		return rerrors.Wrap(err, "delete docker host")
	}

	return nil
}

func (r *Repo) Exists(ctx context.Context) (bool, error) {
	exists, err := r.q.DockerHostExists(ctx)
	if err != nil {
		return false, rerrors.Wrap(err, "check docker host existence")
	}

	return exists, nil
}

// PickLeastLoaded returns the docker_hosts row currently backing the fewest live ('created' or
// 'running') workbenches, ties broken by oldest first. Returns user_errors.NotFound (via
// pg_err.UnwrapPgErr) when no docker hosts are registered at all.
func (r *Repo) PickLeastLoaded(ctx context.Context) (domain.DockerHost, error) {
	row, err := r.q.PickLeastLoadedDockerHost(ctx)
	if err != nil {
		return domain.DockerHost{}, pg_err.UnwrapPgErr(err)
	}

	host := domain.DockerHost{
		Uuid:      row.ID,
		Url:       row.Url,
		CreatedAt: row.CreatedAt,
	}

	return host, nil
}

func (r *Repo) WithTx(tx sqldb.DB) repository.DockerHosts {
	return New(tx)
}
