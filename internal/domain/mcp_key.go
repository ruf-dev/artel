package domain

import (
	"time"

	"github.com/google/uuid"
)

type McpKey struct {
	Uuid           uuid.UUID
	VaultUuid      uuid.UUID
	UserUuid       uuid.UUID
	Name           string
	KeyHash        []byte
	KeyPreview     string
	CreatedAt      time.Time
	RevokedAt      *time.Time
	LastAccessedAt *time.Time
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
	S3        *McpKeyS3Context // nil if the vault has no linked bucket

	UseCouchDBForBinaries bool

	// Postgres is nil unless the vault has an enabled (status=ready) Postgres database.
	Postgres *McpKeyPostgresContext
}

// McpKeyS3Context carries the resolved S3-compatible bucket credentials for the vault behind
// an MCP key, when one is linked.
type McpKeyS3Context struct {
	Endpoint  string
	Region    string
	AccessKey string
	SecretKey string
	UseSSL    bool
	PathStyle bool
	Bucket    string
}

// McpKeyPostgresContext carries the resolved per-vault Postgres role credentials behind an
// MCP key, when the vault has an enabled (status=ready) Postgres database.
type McpKeyPostgresContext struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
	SSLMode  string
}
