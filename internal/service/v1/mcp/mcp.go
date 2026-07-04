package mcp

import (
	"context"

	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service"
	"github.com/ruf-dev/artel/internal/service/v1/mcp/executors"
)

const tokenPrefix = "artel_vtk_"

const bcryptCost = 12

type McpServiceImpl struct {
	mcpKeys             repository.McpKeyRepository
	vaults              repository.Vaults
	vaultMembers        repository.VaultMembers
	couchInstances      repository.CouchInstances
	s3Instances         repository.S3Instances
	mcpConnectors       repository.McpConnectorsRepo
	mcpDefinitions      repository.McpDefinitionsRepo
	externalConnections repository.ExternalConnectionRepo
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
) *McpServiceImpl {
	return &McpServiceImpl{
		mcpKeys:             mcpKeys,
		vaults:              vaults,
		vaultMembers:        vaultMembers,
		couchInstances:      couchInstances,
		s3Instances:         s3Instances,
		mcpConnectors:       mcpConnectors,
		mcpDefinitions:      mcpDefinitions,
		externalConnections: externalConnections,
		vaultExecutor:       executors.NewVaultExecutor(),
		tractExecutor:       executors.NewTractExecutor(),
	}
}

// SetTractService wires the tract service dependency and the server-lifecycle context used to
// spawn run_tract's async StartRun call. See the doc comment on service.McpService.
func (s *McpServiceImpl) SetTractService(baseCtx context.Context, ts service.TractService) {
	s.tractSvc = ts
	s.tractBaseCtx = baseCtx
}
