package s3instances

import (
	"context"

	"github.com/google/uuid"
	s3client "github.com/ruf-dev/artel/internal/clients/s3"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	"go.redsock.ru/rerrors"
)

type Service struct {
	s3InstancesRepo repository.S3Instances
}

func New(repo repository.Repo) *Service {
	return &Service{
		s3InstancesRepo: repo.S3Instances(),
	}
}

func (s *Service) RegisterS3Instance(ctx context.Context, endpoint, region, accessKey, secretKey string, useSSL, pathStyle bool) (string, error) {
	id, err := s.s3InstancesRepo.Register(ctx, endpoint, region, useSSL, pathStyle, accessKey, []byte(secretKey))
	if err != nil {
		return "", rerrors.Wrap(err, "register s3 instance")
	}

	return id.String(), nil
}

func (s *Service) GetS3Instance(ctx context.Context, id string) (domain.S3Instance, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return domain.S3Instance{}, rerrors.Wrap(err, "parse uuid")
	}

	instance, err := s.s3InstancesRepo.Get(ctx, uid)
	if err != nil {
		return domain.S3Instance{}, rerrors.Wrap(err, "get s3 instance")
	}

	return instance, nil
}

func (s *Service) ListS3Instances(ctx context.Context) ([]domain.S3Instance, error) {
	instances, err := s.s3InstancesRepo.List(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "list s3 instances")
	}

	return instances, nil
}

func (s *Service) UpdateS3Instance(ctx context.Context, id, endpoint, region, accessKey, secretKey string, useSSL, pathStyle bool) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return rerrors.Wrap(err, "parse uuid")
	}

	err = s.s3InstancesRepo.Update(ctx, uid, endpoint, region, useSSL, pathStyle, accessKey, []byte(secretKey))
	if err != nil {
		return rerrors.Wrap(err, "update s3 instance")
	}

	return nil
}

func (s *Service) DeleteS3Instance(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return rerrors.Wrap(err, "parse uuid")
	}

	err = s.s3InstancesRepo.Delete(ctx, uid)
	if err != nil {
		return rerrors.Wrap(err, "delete s3 instance")
	}

	return nil
}

func (s *Service) TestS3Instance(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return rerrors.Wrap(err, "parse uuid")
	}

	instance, err := s.s3InstancesRepo.Get(ctx, uid)
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

	err = s3client.TestConnection(ctx, cfg)
	if err != nil {
		return rerrors.Wrap(err, "test s3 connection")
	}

	return nil
}
