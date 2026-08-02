package domain

import (
	"time"

	"github.com/google/uuid"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
)

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
