package domain

import (
	"time"

	"github.com/google/uuid"
)

type TaskTracker struct {
	Uuid      uuid.UUID
	UserUuid  uuid.UUID
	Type      string
	Name      string
	ApiKey    string
	ApiToken  string
	CreatedAt time.Time
}
