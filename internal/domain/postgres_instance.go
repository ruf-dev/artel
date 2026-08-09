package domain

import (
	"time"

	"github.com/google/uuid"
)

// PostgresInstance is a registered Postgres server (admin pool or BYOK) that vault databases can
// be provisioned on — mirrors CouchInstance/S3Instance.
type PostgresInstance struct {
	Uuid          uuid.UUID
	Host          string
	Port          int
	AdminDatabase string
	Username      string
	Password      string
	SSLMode       string
	// OwnerUserUuid is nil for an admin-pool instance, non-nil for a BYOK instance owned by that
	// user — mirrors CouchInstance/S3Instance owner_user_id scoping (migrations/066).
	OwnerUserUuid *uuid.UUID
	CreatedAt     time.Time
}
