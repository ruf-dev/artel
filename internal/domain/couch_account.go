package domain

import (
	"time"

	"github.com/google/uuid"
)

type CouchAccount struct {
	Uuid              uuid.UUID
	UserUuid          uuid.UUID
	CouchInstanceUuid uuid.UUID
	CouchUsername     string
	CouchPassword     string
	CreatedAt         time.Time
}
