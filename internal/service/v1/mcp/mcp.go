package mcp

import (
	"context"

	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service"
	"github.com/ruf-dev/artel/internal/service/v1/mcp/executors"
)

const tokenPrefix = "artel_vtk_" //nolint:gosec // this is a public token prefix, not a credential

const bcryptCost = 12

type ServiceImpl struct {
	mcpKeys             repository.McpKeyRepository
	vaults              repository.Vaults
	vaultMembers        repository.VaultMembers
	couchInstances      repository.CouchInstances
	s3Instances         repository.S3Instances
	mcpConnectors       repository.McpConnectorsRepo
	mcpDefinitions      repository.McpDefinitionsRepo
	externalConnections repository.ExternalConnectionRepo
	subscriptions       service.SubscriptionService
	authSvc             service.AuthService
	vaultExecutor       *executors.VaultExecutor
	tractExecutor       *executors.TractExecutor

	// tractSvc/tractBaseCtx are unset until SetTractService is called from
	// internal/app/custom.go, after TractService is constructed (Tract composes Mcp's
	// ToolExecutor, so Mcp must exist first). The tract-authoring builtin tools degrade with
	// user_errors.TractServiceNotConfigured while tractSvc is nil.
	tractSvc     service.TractService
	tractBaseCtx context.Context
}

func New(
	mcpKeys repository.McpKeyRepository,
	vaults repository.Vaults,
	vaultMembers repository.VaultMembers,
	couchInstances repository.CouchInstances,
	s3Instances repository.S3Instances,
	mcpConnectors repository.McpConnectorsRepo,
	mcpDefinitions repository.McpDefinitionsRepo,
	externalConnections repository.ExternalConnectionRepo,
	subscriptions service.SubscriptionService,
	authSvc service.AuthService,
) *ServiceImpl {
	impl := &ServiceImpl{
		mcpKeys:             mcpKeys,
		vaults:              vaults,
		vaultMembers:        vaultMembers,
		couchInstances:      couchInstances,
		s3Instances:         s3Instances,
		mcpConnectors:       mcpConnectors,
		mcpDefinitions:      mcpDefinitions,
		externalConnections: externalConnections,
		subscriptions:       subscriptions,
		authSvc:             authSvc,
		vaultExecutor:       executors.NewVaultExecutor(),
		tractExecutor:       executors.NewTractExecutor(),
	}

	return impl
}

// SetTractService wires the tract service dependency and the server-lifecycle context used to
// spawn run_tract's async StartRun call. See the doc comment on service.McpService.
func (s *ServiceImpl) SetTractService(baseCtx context.Context, ts service.TractService) {
	s.tractSvc = ts
	s.tractBaseCtx = baseCtx
}
