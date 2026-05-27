package couch_instances_api

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"go.redsock.ru/rerrors"
	"google.golang.org/grpc"

	"github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service"
)

type CouchInstancesImpl struct {
	artel_api.UnimplementedCouchInstancesAPIServer
	couchInstanceSvc service.CouchInstanceService
}

func NewCouchInstancesImpl(couchInstanceSvc service.CouchInstanceService) *CouchInstancesImpl {
	return &CouchInstancesImpl{couchInstanceSvc: couchInstanceSvc}
}

func (c *CouchInstancesImpl) Register(srv grpc.ServiceRegistrar) {
	artel_api.RegisterCouchInstancesAPIServer(srv, c)
}

func (c *CouchInstancesImpl) Gateway(ctx context.Context, endpoint string, opts ...grpc.DialOption) (string, http.Handler) {
	gwMux := runtime.NewServeMux()

	err := artel_api.RegisterCouchInstancesAPIHandlerFromEndpoint(ctx, gwMux, endpoint, opts)
	if err != nil {
		log.Error().Err(err).Msg("error registering couch instances grpc-gateway handler")
	}

	return "/api/couch/", gwMux
}

func (c *CouchInstancesImpl) RegisterCouchInstance(ctx context.Context, req *artel_api.RegisterCouchInstance_Request) (*artel_api.RegisterCouchInstance_Response, error) {
	id, err := c.couchInstanceSvc.RegisterCouchInstance(ctx, req.Url, req.Username, req.Password)
	if err != nil {
		return nil, rerrors.Wrap(err, "register couch instance")
	}

	resp := &artel_api.RegisterCouchInstance_Response{
		Id: id,
	}
	return resp, nil
}

func (c *CouchInstancesImpl) GetCouchInstance(ctx context.Context, req *artel_api.GetCouchInstance_Request) (*artel_api.GetCouchInstance_Response, error) {
	instance, err := c.couchInstanceSvc.GetCouchInstance(ctx, req.Id)
	if err != nil {
		return nil, rerrors.Wrap(err, "get couch instance")
	}

	resp := &artel_api.GetCouchInstance_Response{
		Id:        instance.Uuid.String(),
		Url:       instance.Url,
		Username:  instance.Username,
		CreatedAt: instance.CreatedAt.String(),
	}
	return resp, nil
}

func (c *CouchInstancesImpl) ListCouchInstances(ctx context.Context, req *artel_api.ListCouchInstances_Request) (*artel_api.ListCouchInstances_Response, error) {
	instances, err := c.couchInstanceSvc.ListCouchInstances(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "list couch instances")
	}

	respInstances := make([]*artel_api.GetCouchInstance_Response, len(instances))
	for i, instance := range instances {
		respInstances[i] = &artel_api.GetCouchInstance_Response{
			Id:        instance.Uuid.String(),
			Url:       instance.Url,
			Username:  instance.Username,
			CreatedAt: instance.CreatedAt.String(),
		}
	}

	resp := &artel_api.ListCouchInstances_Response{
		Instances: respInstances,
	}
	return resp, nil
}

func (c *CouchInstancesImpl) UpdateCouchInstance(ctx context.Context, req *artel_api.UpdateCouchInstance_Request) (*artel_api.UpdateCouchInstance_Response, error) {
	err := c.couchInstanceSvc.UpdateCouchInstance(ctx, req.Id, req.Url, req.Username, req.Password)
	if err != nil {
		return nil, rerrors.Wrap(err, "update couch instance")
	}

	return &artel_api.UpdateCouchInstance_Response{}, nil
}

func (c *CouchInstancesImpl) DeleteCouchInstance(ctx context.Context, req *artel_api.DeleteCouchInstance_Request) (*artel_api.DeleteCouchInstance_Response, error) {
	err := c.couchInstanceSvc.DeleteCouchInstance(ctx, req.Id)
	if err != nil {
		return nil, rerrors.Wrap(err, "delete couch instance")
	}

	resp := &artel_api.DeleteCouchInstance_Response{}
	return resp, nil
}
