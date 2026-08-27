package domain

import (
	"time"

	"github.com/google/uuid"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
)

// ReservedGithubDocsSlug is a sentinel `/docs/:slug` value that never belongs to a real vault —
// vault.validateSlug rejects it outright. public_docs.Service recognizes it and routes to the
// GitHub-backed quick-start resolver (internal/service/v1/public_docs/githubdocs) instead of
// resolving a vaults row, letting GetDefaultVault point unauthenticated `/docs` visitors at the
// built-in quick-start guide without any frontend route/RPC changes.
const ReservedGithubDocsSlug = "github-quickstart"

type Vault struct {
	Uuid                  uuid.UUID
	UserUuid              uuid.UUID
	CouchInstanceUuid     uuid.UUID
	Name                  string
	CouchDBName           string
	CouchDBURL            string
	LiveSyncPassphrase    string
	Status                string
	S3InstanceUuid        *uuid.UUID // nil when vault has no linked bucket
	S3BucketName          string     // "" when S3InstanceUuid is nil
	UseCouchDBForBinaries bool
	CreatedAt             time.Time
	IsPublic              bool
	Slug                  string // "" when vault has never been published
	Prompt                string
	UseSystemPrompt       bool

	// MyRole is the calling user's membership role on this vault ("owner"/"reader"/"maintainer"),
	// or "" if they have no membership row. Only populated by ListByMembership — every other
	// repo method (GetByID, CreateVault, etc.) has no "calling user" concept and leaves it "".
	MyRole string
}

type VaultMember struct {
	Uuid      uuid.UUID
	VaultUuid uuid.UUID
	UserUuid  uuid.UUID
	Role      artel_q.VaultRole
	CreatedAt time.Time
}

type VaultMemberInfo struct {
	VaultMember
	Email    string
	Username string
}

type VaultInvite struct {
	Uuid      uuid.UUID
	VaultUuid uuid.UUID
	CreatedBy uuid.UUID
	Role      artel_q.VaultRole
	Token     string
	RevokedAt *time.Time
	CreatedAt time.Time
}
