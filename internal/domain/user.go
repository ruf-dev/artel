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
}

type ArtelIdentity struct {
	Email        string
	PasswordHash string
}

type UserIdentities struct {
	Telegram *TelegramIdentity
	Artel    *ArtelIdentity
}

// UserPermissions now holds only the internal admin-access flag — per-feature gating (emails,
// notes, task trackers, spreadsheets, connectors, tract) moved to the subscription layer, see
// SubscriptionFeature/FeatureSet in subscription.go.
type UserPermissions struct {
	UserUuid        uuid.UUID
	IsAdministrator bool
}

type UserDetails struct {
	User
	Permissions           UserPermissions
	EffectiveSubscription EffectiveSubscription
}
