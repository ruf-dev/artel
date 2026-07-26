package docker_hosts_api

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service"
	"go.redsock.ru/rerrors"
	"google.golang.org/grpc"
)

type DockerHostsImpl struct {
	artel_api.UnimplementedDockerHostsAPIServer
	dockerHostSvc service.DockerHostService
}

func NewDockerHostsImpl(dockerHostSvc service.DockerHostService) *DockerHostsImpl {
	return &DockerHostsImpl{dockerHostSvc: dockerHostSvc}
}

func (d *DockerHostsImpl) Register(srv grpc.ServiceRegistrar) {
	artel_api.RegisterDockerHostsAPIServer(srv, d)
}

func (d *DockerHostsImpl) Gateway(
	ctx context.Context,
	endpoint string,
	opts ...grpc.DialOption,
) (string, http.Handler) {
	gwMux := runtime.NewServeMux()

	err := artel_api.RegisterDockerHostsAPIHandlerFromEndpoint(ctx, gwMux, endpoint, opts)
	if err != nil {
		log.Error().Err(err).Msg("error registering docker hosts grpc-gateway handler")
	}

	return "/api/docker_hosts/", gwMux
}

func (d *DockerHostsImpl) RegisterDockerHost(
	ctx context.Context,
	req *artel_api.RegisterDockerHost_Request,
) (*artel_api.RegisterDockerHost_Response, error) {
	id, err := d.dockerHostSvc.RegisterDockerHost(ctx, req.Url)
	if err != nil {
		return nil, rerrors.Wrap(err, "register docker host")
	}

	resp := &artel_api.RegisterDockerHost_Response{
		Id: id,
	}

	return resp, nil
}

func (d *DockerHostsImpl) GetDockerHost(
	ctx context.Context,
	req *artel_api.GetDockerHost_Request,
) (*artel_api.GetDockerHost_Response, error) {
	host, err := d.dockerHostSvc.GetDockerHost(ctx, req.Id)
	if err != nil {
		return nil, rerrors.Wrap(err, "get docker host")
	}

	resp := &artel_api.GetDockerHost_Response{
		Id:        host.Uuid.String(),
		Url:       host.Url,
		CreatedAt: host.CreatedAt.String(),
	}

	return resp, nil
}

func (d *DockerHostsImpl) ListDockerHosts(
	ctx context.Context,
	req *artel_api.ListDockerHosts_Request,
) (*artel_api.ListDockerHosts_Response, error) {
	hosts, err := d.dockerHostSvc.ListDockerHosts(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "list docker hosts")
	}

	respHosts := make([]*artel_api.GetDockerHost_Response, len(hosts))
	for i, host := range hosts {
		respHosts[i] = &artel_api.GetDockerHost_Response{
			Id:        host.Uuid.String(),
			Url:       host.Url,
			CreatedAt: host.CreatedAt.String(),
		}
	}

	resp := &artel_api.ListDockerHosts_Response{
		Hosts: respHosts,
	}

	return resp, nil
}

func (d *DockerHostsImpl) UpdateDockerHost(
	ctx context.Context,
	req *artel_api.UpdateDockerHost_Request,
) (*artel_api.UpdateDockerHost_Response, error) {
	err := d.dockerHostSvc.UpdateDockerHost(ctx, req.Id, req.Url)
	if err != nil {
		return nil, rerrors.Wrap(err, "update docker host")
	}

	return &artel_api.UpdateDockerHost_Response{}, nil
}

func (d *DockerHostsImpl) DeleteDockerHost(
	ctx context.Context,
	req *artel_api.DeleteDockerHost_Request,
) (*artel_api.DeleteDockerHost_Response, error) {
	err := d.dockerHostSvc.DeleteDockerHost(ctx, req.Id)
	if err != nil {
		return nil, rerrors.Wrap(err, "delete docker host")
	}

	resp := &artel_api.DeleteDockerHost_Response{}

	return resp, nil
}
