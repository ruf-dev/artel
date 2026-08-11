package docker_hosts_api

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
	gwMux := transport.NewGatewayMux()

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
	id, err := d.dockerHostSvc.RegisterDockerHost(
		ctx, req.Url, stringOrEmpty(req.CaCert), stringOrEmpty(req.ClientCert), stringOrEmpty(req.ClientKey),
	)
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
	// req.CaCert/ClientCert/ClientKey are proto3 `optional` *string, and their nil-vs-set
	// presence is passed straight through as the service/repo layer's three-way patch pointer:
	// nil means "the caller didn't touch this cert", not "clear it".
	err := d.dockerHostSvc.UpdateDockerHost(ctx, req.Id, req.Url, req.CaCert, req.ClientCert, req.ClientKey)
	if err != nil {
		return nil, rerrors.Wrap(err, "update docker host")
	}

	return &artel_api.UpdateDockerHost_Response{}, nil
}

// stringOrEmpty dereferences a proto3 `optional string` field, treating "not provided" (nil)
// the same as an explicit empty string — used for Register, where there's no existing value to
// preserve, so there's no third "don't touch" state to distinguish (unlike Update).
func stringOrEmpty(v *string) string {
	if v == nil {
		return ""
	}

	return *v
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
