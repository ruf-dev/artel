// Package dockerhosts implements service.DockerHostService: admin CRUD over the pool of Docker
// daemons that back per-vault workbench containers — see docs/workbench/02_docker_topology.md.
// Unlike couchinstances/s3instances, docker hosts carry no credentials, so there's no
// setup/status concept here.
package dockerhosts

import (
	"context"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	"go.redsock.ru/rerrors"
)

type Service struct {
	dockerHostsRepo repository.DockerHosts
}

func New(repo repository.Repo) *Service {
	return &Service{
		dockerHostsRepo: repo.DockerHosts(),
	}
}

func (s *Service) RegisterDockerHost(ctx context.Context, url, caCert, clientCert, clientKey string) (string, error) {
	id, err := s.dockerHostsRepo.Register(ctx, url, caCert, clientCert, clientKey)
	if err != nil {
		return "", rerrors.Wrap(err, "register docker host")
	}

	return id.String(), nil
}

func (s *Service) GetDockerHost(ctx context.Context, id string) (domain.DockerHost, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return domain.DockerHost{}, rerrors.Wrap(err, "parse uuid")
	}

	host, err := s.dockerHostsRepo.Get(ctx, uid)
	if err != nil {
		return domain.DockerHost{}, rerrors.Wrap(err, "get docker host")
	}

	return host, nil
}

func (s *Service) ListDockerHosts(ctx context.Context) ([]domain.DockerHost, error) {
	hosts, err := s.dockerHostsRepo.List(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "list docker hosts")
	}

	return hosts, nil
}

func (s *Service) UpdateDockerHost(ctx context.Context, id, url string, caCert, clientCert, clientKey *string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return rerrors.Wrap(err, "parse uuid")
	}

	err = s.dockerHostsRepo.Update(ctx, uid, url, caCert, clientCert, clientKey)
	if err != nil {
		return rerrors.Wrap(err, "update docker host")
	}

	return nil
}

func (s *Service) DeleteDockerHost(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return rerrors.Wrap(err, "parse uuid")
	}

	err = s.dockerHostsRepo.Delete(ctx, uid)
	if err != nil {
		return rerrors.Wrap(err, "delete docker host")
	}

	return nil
}

func (s *Service) HasDockerHosts(ctx context.Context) (bool, error) {
	exists, err := s.dockerHostsRepo.Exists(ctx)
	if err != nil {
		return false, rerrors.Wrap(err, "check docker hosts availability")
	}

	return exists, nil
}
