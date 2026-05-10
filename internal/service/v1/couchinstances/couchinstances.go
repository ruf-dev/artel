package couchinstances

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
)

type Service struct {
	repo repository.CouchInstances
}

func New(repo repository.CouchInstances) *Service {
	return &Service{repo: repo}
}

func (s *Service) RegisterCouchInstance(ctx context.Context, url, username, password string) (string, error) {
	id, err := s.repo.Register(ctx, url, username, []byte(password))
	if err != nil {
		return "", rerrors.Wrap(err, "register couch instance")
	}

	return id.String(), nil
}

func (s *Service) GetCouchInstance(ctx context.Context, id string) (domain.CouchInstance, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return domain.CouchInstance{}, rerrors.Wrap(err, "parse uuid")
	}

	instance, err := s.repo.Get(ctx, uid)
	if err != nil {
		return domain.CouchInstance{}, rerrors.Wrap(err, "get couch instance")
	}

	return instance, nil
}

func (s *Service) ListCouchInstances(ctx context.Context) ([]domain.CouchInstance, error) {
	instances, err := s.repo.List(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "list couch instances")
	}

	return instances, nil
}

func (s *Service) DeleteCouchInstance(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return rerrors.Wrap(err, "parse uuid")
	}

	err = s.repo.Delete(ctx, uid)
	if err != nil {
		return rerrors.Wrap(err, "delete couch instance")
	}

	return nil
}
