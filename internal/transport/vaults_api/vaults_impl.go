package vaults_api

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service"
	"github.com/ruf-dev/artel/internal/transport"
	"google.golang.org/grpc"
)

type VaultsImpl struct {
	pb.UnimplementedVaultsAPIServer
	vaultSvc     service.VaultService
	workbenchSvc service.WorkbenchService
	cookieSecure bool
}

func NewVaultsImpl(
	vaultSvc service.VaultService, workbenchSvc service.WorkbenchService, cookieSecure bool,
) *VaultsImpl {
	return &VaultsImpl{vaultSvc: vaultSvc, workbenchSvc: workbenchSvc, cookieSecure: cookieSecure}
}

func (v *VaultsImpl) Register(srv grpc.ServiceRegistrar) {
	pb.RegisterVaultsAPIServer(srv, v)
}

func (v *VaultsImpl) Gateway(ctx context.Context, endpoint string, opts ...grpc.DialOption) (string, http.Handler) {
	gwMux := transport.NewGatewayMux(v.cookieSecure)

	err := pb.RegisterVaultsAPIHandlerFromEndpoint(ctx, gwMux, endpoint, opts)
	if err != nil {
		log.Error().Err(err).Msg("error registering vaults grpc-gateway handler")
	}

	return "/api/vaults/", gwMux
}
