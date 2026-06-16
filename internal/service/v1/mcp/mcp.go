package mcp

import (
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service/v1/mcp/executors"
)

const tokenPrefix = "artel_vtk_"

const bcryptCost = 12

type McpServiceImpl struct {
	mcpKeys             repository.McpKeyRepository
	vaults              repository.Vaults
	vaultMembers        repository.VaultMembers
	couchInstances      repository.CouchInstances
	mcpConnectors       repository.McpConnectorsRepo
	mcpDefinitions      repository.McpDefinitionsRepo
	externalConnections repository.ExternalConnectionRepo
	vaultExecutor       *executors.VaultExecutor
}

func New(
	mcpKeys repository.McpKeyRepository,
	vaults repository.Vaults,
	vaultMembers repository.VaultMembers,
	couchInstances repository.CouchInstances,
	mcpConnectors repository.McpConnectorsRepo,
	mcpDefinitions repository.McpDefinitionsRepo,
	externalConnections repository.ExternalConnectionRepo,
) *McpServiceImpl {
	return &McpServiceImpl{
		mcpKeys:             mcpKeys,
		vaults:              vaults,
		vaultMembers:        vaultMembers,
		couchInstances:      couchInstances,
		mcpConnectors:       mcpConnectors,
		mcpDefinitions:      mcpDefinitions,
		externalConnections: externalConnections,
		vaultExecutor:       executors.NewVaultExecutor(),
	}
}
