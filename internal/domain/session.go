package domain

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Uuid             uuid.UUID
	UserUuid         uuid.UUID
	Token            string
	ExpiresAt        time.Time
	RefreshToken     string
	RefreshExpiresAt time.Time
	CreatedAt        time.Time
}
