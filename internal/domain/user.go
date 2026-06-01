package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Uuid         uuid.UUID
	Email        string
	Username     string
	PhotoUrl     string
	PasswordHash string
	Roles        []string
	Identities   UserIdentities
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type TelegramIdentity struct {
	UserUuid   uuid.UUID
	TelegramId string
	PhotoUrl   string
}

type ArtelIdentity struct {
	Email        string
	PasswordHash string
}

type UserIdentities struct {
	Telegram *TelegramIdentity
	Artel    *ArtelIdentity
}

type UserPermissions struct {
	UserUuid        uuid.UUID
	IsAdministrator bool
	HasEmails       bool
	HasTaskTrackers bool
}

type Subscription struct {
	UserUuid uuid.UUID
	Active   bool
}
