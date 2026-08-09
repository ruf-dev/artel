package domain

import (
	"time"

	"github.com/google/uuid"
)

type VaultPostgresStatus string

const (
	VaultPostgresStatusProvisioning VaultPostgresStatus = "provisioning"
	VaultPostgresStatusReady        VaultPostgresStatus = "ready"
	VaultPostgresStatusError        VaultPostgresStatus = "error"
)

// VaultPostgresDatabase is the per-vault Postgres database+role provisioned on a PostgresInstance
// — mirrors the Workbench provisioning→ready/error lifecycle.
type VaultPostgresDatabase struct {
	VaultUuid            uuid.UUID
	PostgresInstanceUuid uuid.UUID
	DatabaseName         string
	RoleUsername         string
	RolePassword         string
	Status               VaultPostgresStatus
	ErrorMessage         string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}
