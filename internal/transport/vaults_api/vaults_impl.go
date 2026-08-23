package vaults_api

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service"
	"github.com/ruf-dev/artel/internal/transport"
	"google.golang.org/grpc"
)

// terminalAuthLinkCache is the narrow subset of *WorkbenchTerminalShellHandler GetVault depends
// on: reading back whatever Claude CLI OAuth sign-in link its WS relay most recently detected on
// a vault's terminal output. It's the WS relay handler, not VaultService, that owns this
// in-memory cache — see WorkbenchTerminalShellHandler's doc comment on its authLinks field — so
// VaultsImpl depends on it directly, through this narrow interface, rather than through
// service.WorkbenchService.
type terminalAuthLinkCache interface {
	PendingTerminalAuthLink(vaultID uuid.UUID) string
}

type VaultsImpl struct {
	pb.UnimplementedVaultsAPIServer
	vaultSvc          service.VaultService
	workbenchSvc      service.WorkbenchService
	terminalAuthLinks terminalAuthLinkCache
}

func NewVaultsImpl(
	vaultSvc service.VaultService, workbenchSvc service.WorkbenchService, terminalAuthLinks terminalAuthLinkCache,
) *VaultsImpl {
	return &VaultsImpl{vaultSvc: vaultSvc, workbenchSvc: workbenchSvc, terminalAuthLinks: terminalAuthLinks}
}

func (v *VaultsImpl) Register(srv grpc.ServiceRegistrar) {
	pb.RegisterVaultsAPIServer(srv, v)
}

func (v *VaultsImpl) Gateway(ctx context.Context, endpoint string, opts ...grpc.DialOption) (string, http.Handler) {
	gwMux := transport.NewGatewayMux()

	err := pb.RegisterVaultsAPIHandlerFromEndpoint(ctx, gwMux, endpoint, opts)
	if err != nil {
		log.Error().Err(err).Msg("error registering vaults grpc-gateway handler")
	}

	return "/api/vaults/", gwMux
}
