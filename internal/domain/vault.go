package domain

import (
	"time"

	"github.com/google/uuid"
)

type Vault struct {
	Uuid        uuid.UUID
	UserUuid    uuid.UUID
	Name        string
	CouchDBName string
	CouchDBURL  string
	CreatedAt   time.Time
}

type CouchCred struct {
	Uuid      uuid.UUID
	VaultUuid uuid.UUID
	Host      string
	Username  string
	Password  []byte
	CreatedAt time.Time
}
