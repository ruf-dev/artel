package couchinstances

import (
	"context"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/clients/couchdb"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	"go.redsock.ru/rerrors"
)

type Service struct {
	couchInstancesRepo repository.CouchInstances
}

func New(repo repository.Repo) *Service {
	return &Service{
		couchInstancesRepo: repo.CouchInstances(),
	}
}

func (s *Service) RegisterCouchInstance(ctx context.Context, url, username, password string) (string, error) {
	id, err := s.couchInstancesRepo.Register(ctx, url, username, []byte(password))
	if err != nil {
		return "", rerrors.Wrap(err, "register couch instance")
	}

	cfg := couchdb.Config{
		BaseURL:  url,
		User:     username,
		Password: password,
	}

	client, err := couchdb.New(cfg)
	if err != nil {
		return "", rerrors.Wrap(err, "creating couch client")
	}

	err = client.Setup(ctx)
	if err != nil {
		return "", rerrors.Wrap(err, "setup couch instance")
	}

	return id.String(), nil
}

func (s *Service) GetCouchInstance(ctx context.Context, id string) (domain.CouchInstance, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return domain.CouchInstance{}, rerrors.Wrap(err, "parse uuid")
	}

	instance, err := s.couchInstancesRepo.Get(ctx, uid)
	if err != nil {
		return domain.CouchInstance{}, rerrors.Wrap(err, "get couch instance")
	}

	return instance, nil
}

func (s *Service) ListCouchInstances(ctx context.Context) ([]domain.CouchInstance, error) {
	instances, err := s.couchInstancesRepo.List(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "list couch instances")
	}

	return instances, nil
}

func (s *Service) UpdateCouchInstance(ctx context.Context, id, url, username, password string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return rerrors.Wrap(err, "parse uuid")
	}

	err = s.couchInstancesRepo.Update(ctx, uid, url, username, []byte(password))
	if err != nil {
		return rerrors.Wrap(err, "update couch instance")
	}

	return nil
}

func (s *Service) DeleteCouchInstance(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return rerrors.Wrap(err, "parse uuid")
	}

	err = s.couchInstancesRepo.Delete(ctx, uid)
	if err != nil {
		return rerrors.Wrap(err, "delete couch instance")
	}

	return nil
}

func (s *Service) SetupCouchInstance(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return rerrors.Wrap(err, "parse uuid")
	}

	instance, err := s.couchInstancesRepo.Get(ctx, uid)
	if err != nil {
		return rerrors.Wrap(err, "get couch instance")
	}

	cfg := couchdb.Config{
		BaseURL:  instance.Url,
		User:     instance.Username,
		Password: instance.Password,
	}

	client, err := couchdb.New(cfg)
	if err != nil {
		return rerrors.Wrap(err, "creating couch client")
	}

	err = client.Setup(ctx)
	if err != nil {
		return rerrors.Wrap(err, "setup couch instance")
	}

	return nil
}

func (s *Service) GetCouchInstanceStatus(ctx context.Context, id string) (couchdb.SetupStatus, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return couchdb.SetupStatus{}, rerrors.Wrap(err, "parse uuid")
	}

	instance, err := s.couchInstancesRepo.Get(ctx, uid)
	if err != nil {
		return couchdb.SetupStatus{}, rerrors.Wrap(err, "get couch instance")
	}

	cfg := couchdb.Config{
		BaseURL:  instance.Url,
		User:     instance.Username,
		Password: instance.Password,
	}

	client, err := couchdb.New(cfg)
	if err != nil {
		return couchdb.SetupStatus{}, rerrors.Wrap(err, "creating couch client")
	}

	status, err := client.GetSetupStatus(ctx)
	if err != nil {
		return couchdb.SetupStatus{}, rerrors.Wrap(err, "get setup status")
	}

	return status, nil
}
