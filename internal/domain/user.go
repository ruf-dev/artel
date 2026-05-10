package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Uuid      uuid.UUID
	Email     string
	Roles     []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Subscription struct {
	Uuid      uuid.UUID
	UserUuid  uuid.UUID
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
