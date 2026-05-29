package domain

import (
	"time"

	"github.com/google/uuid"
)

type McpKey struct {
	Uuid       uuid.UUID
	VaultUuid  uuid.UUID
	UserUuid   uuid.UUID
	Name       string
	KeyHash    []byte
	KeyPreview string
	CreatedAt  time.Time
	RevokedAt  *time.Time
}

// McpKeyContext is resolved from a raw bearer token; contains everything
// needed to connect to CouchDB for the associated vault.
type McpKeyContext struct {
	KeyUuid   uuid.UUID
	VaultUuid uuid.UUID
	UserUuid  uuid.UUID
	CouchURL  string
	CouchDb   string
	CouchUser string
	CouchPass string
	HasEmails bool
}
