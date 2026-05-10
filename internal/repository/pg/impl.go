package pg

import (
	"database/sql"

	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/repository/pg/couchcreds"
	"github.com/ruf-dev/artel/internal/repository/pg/couchinstances"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"github.com/ruf-dev/artel/internal/repository/pg/subscriptions"
	"github.com/ruf-dev/artel/internal/repository/pg/users"
	"github.com/ruf-dev/artel/internal/repository/pg/vaults"
)

type Repos struct {
	Users            repository.Users
	Vaults           repository.Vaults
	Subscriptions    repository.Subscriptions
	CouchCredentials repository.CouchCredentials
	CouchInstances   repository.CouchInstances
}

func New(db *sql.DB, encryptionKey []byte) *Repos {
	q := artel_q.New(db)

	return &Repos{
		Users:            users.New(q),
		Vaults:           vaults.New(q),
		Subscriptions:    subscriptions.New(q),
		CouchCredentials: couchcreds.New(q, encryptionKey),
		CouchInstances:   couchinstances.New(q, encryptionKey),
	}
}
