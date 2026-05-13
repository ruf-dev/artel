package vaults_api

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service"
)

type VaultsImpl struct {
	pb.UnimplementedVaultsAPIServer
	vaultSvc service.VaultService
}

func NewVaultsImpl(vaultSvc service.VaultService) *VaultsImpl {
	return &VaultsImpl{vaultSvc: vaultSvc}
}

func (v *VaultsImpl) Register(srv grpc.ServiceRegistrar) {
	pb.RegisterVaultsAPIServer(srv, v)
}

func (v *VaultsImpl) Gateway(ctx context.Context, endpoint string, opts ...grpc.DialOption) (string, http.Handler) {
	gwMux := runtime.NewServeMux()

	err := pb.RegisterVaultsAPIHandlerFromEndpoint(ctx, gwMux, endpoint, opts)
	if err != nil {
		log.Error().Err(err).Msg("error registering vaults grpc-gateway handler")
	}

	return "/api/vaults/", gwMux
}
