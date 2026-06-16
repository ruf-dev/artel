package mcp

import (
	"github.com/ruf-dev/artel/internal/repository"
)

const tokenPrefix = "artel_vtk_"

const bcryptCost = 12

type McpServiceImpl struct {
	mcpKeys        repository.McpKeyRepository
	vaults         repository.Vaults
	vaultMembers   repository.VaultMembers
	couchInstances repository.CouchInstances
}

func New(
	mcpKeys repository.McpKeyRepository,
	vaults repository.Vaults,
	vaultMembers repository.VaultMembers,
	couchInstances repository.CouchInstances,
) *McpServiceImpl {
	return &McpServiceImpl{
		mcpKeys:        mcpKeys,
		vaults:         vaults,
		vaultMembers:   vaultMembers,
		couchInstances: couchInstances,
	}
}
