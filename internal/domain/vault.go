package domain

import (
	"time"

	"github.com/google/uuid"
)

type Vault struct {
	Uuid              uuid.UUID
	UserUuid          uuid.UUID
	CouchInstanceUuid uuid.UUID
	Name              string
	CouchDBName       string
	CouchDBURL        string
	CreatedAt         time.Time
}

type VaultMember struct {
	Uuid      uuid.UUID
	VaultUuid uuid.UUID
	UserUuid  uuid.UUID
	Role      string
	CreatedAt time.Time
}
