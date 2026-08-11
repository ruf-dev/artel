package s3_instances_api

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service"
	"github.com/ruf-dev/artel/internal/transport"
	"go.redsock.ru/rerrors"
	"google.golang.org/grpc"
)

type S3InstancesImpl struct {
	artel_api.UnimplementedS3InstancesAPIServer
	s3InstanceSvc service.S3InstanceService
}

func NewS3InstancesImpl(s3InstanceSvc service.S3InstanceService) *S3InstancesImpl {
	return &S3InstancesImpl{s3InstanceSvc: s3InstanceSvc}
}

func (s *S3InstancesImpl) Register(srv grpc.ServiceRegistrar) {
	artel_api.RegisterS3InstancesAPIServer(srv, s)
}

func (s *S3InstancesImpl) Gateway(
	ctx context.Context,
	endpoint string,
	opts ...grpc.DialOption,
) (string, http.Handler) {
	gwMux := transport.NewGatewayMux()

	err := artel_api.RegisterS3InstancesAPIHandlerFromEndpoint(ctx, gwMux, endpoint, opts)
	if err != nil {
		log.Error().Err(err).Msg("error registering s3 instances grpc-gateway handler")
	}

	return "/api/s3/", gwMux
}

func (s *S3InstancesImpl) RegisterS3Instance(
	ctx context.Context,
	req *artel_api.RegisterS3Instance_Request,
) (*artel_api.RegisterS3Instance_Response, error) {
	id, err := s.s3InstanceSvc.RegisterS3Instance(
		ctx,
		req.Endpoint,
		req.Region,
		req.AccessKey,
		req.SecretKey,
		req.UseSsl,
		req.PathStyle,
	)
	if err != nil {
		return nil, rerrors.Wrap(err, "register s3 instance")
	}

	resp := &artel_api.RegisterS3Instance_Response{
		Id: id,
	}

	return resp, nil
}

func (s *S3InstancesImpl) GetS3Instance(
	ctx context.Context,
	req *artel_api.GetS3Instance_Request,
) (*artel_api.GetS3Instance_Response, error) {
	instance, err := s.s3InstanceSvc.GetS3Instance(ctx, req.Id)
	if err != nil {
		return nil, rerrors.Wrap(err, "get s3 instance")
	}

	resp := &artel_api.GetS3Instance_Response{
		Id:        instance.Uuid.String(),
		Endpoint:  instance.Endpoint,
		Region:    instance.Region,
		UseSsl:    instance.UseSSL,
		PathStyle: instance.PathStyle,
		CreatedAt: instance.CreatedAt.String(),
	}

	return resp, nil
}

func (s *S3InstancesImpl) ListS3Instances(
	ctx context.Context,
	req *artel_api.ListS3Instances_Request,
) (*artel_api.ListS3Instances_Response, error) {
	instances, err := s.s3InstanceSvc.ListS3Instances(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "list s3 instances")
	}

	respInstances := make([]*artel_api.GetS3Instance_Response, len(instances))
	for i, instance := range instances {
		respInstances[i] = &artel_api.GetS3Instance_Response{
			Id:        instance.Uuid.String(),
			Endpoint:  instance.Endpoint,
			Region:    instance.Region,
			UseSsl:    instance.UseSSL,
			PathStyle: instance.PathStyle,
			CreatedAt: instance.CreatedAt.String(),
		}
	}

	resp := &artel_api.ListS3Instances_Response{
		Instances: respInstances,
	}

	return resp, nil
}

func (s *S3InstancesImpl) UpdateS3Instance(
	ctx context.Context,
	req *artel_api.UpdateS3Instance_Request,
) (*artel_api.UpdateS3Instance_Response, error) {
	err := s.s3InstanceSvc.UpdateS3Instance(
		ctx,
		req.Id,
		req.Endpoint,
		req.Region,
		req.AccessKey,
		req.SecretKey,
		req.UseSsl,
		req.PathStyle,
	)
	if err != nil {
		return nil, rerrors.Wrap(err, "update s3 instance")
	}

	return &artel_api.UpdateS3Instance_Response{}, nil
}

func (s *S3InstancesImpl) DeleteS3Instance(
	ctx context.Context,
	req *artel_api.DeleteS3Instance_Request,
) (*artel_api.DeleteS3Instance_Response, error) {
	err := s.s3InstanceSvc.DeleteS3Instance(ctx, req.Id)
	if err != nil {
		return nil, rerrors.Wrap(err, "delete s3 instance")
	}

	resp := &artel_api.DeleteS3Instance_Response{}

	return resp, nil
}

func (s *S3InstancesImpl) TestS3Instance(
	ctx context.Context,
	req *artel_api.TestS3Instance_Request,
) (*artel_api.TestS3Instance_Response, error) {
	err := s.s3InstanceSvc.TestS3Instance(ctx, req.Id)
	if err != nil {
		return nil, rerrors.Wrap(err, "test s3 instance")
	}

	resp := &artel_api.TestS3Instance_Response{}

	return resp, nil
}
