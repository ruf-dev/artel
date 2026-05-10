package pg

import (
	"database/sql"

	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/repository/pg/couchcreds"
	"github.com/ruf-dev/artel/internal/repository/pg/subscriptions"
	"github.com/ruf-dev/artel/internal/repository/pg/users"
	"github.com/ruf-dev/artel/internal/repository/pg/vaults"
)

type Repos struct {
	Users            repository.Users
	Vaults           repository.Vaults
	Subscriptions    repository.Subscriptions
	CouchCredentials repository.CouchCredentials
}

func New(db *sql.DB, encryptionKey []byte) *Repos {
	return &Repos{
		Users:            users.New(db),
		Vaults:           vaults.New(db),
		Subscriptions:    subscriptions.New(db),
		CouchCredentials: couchcreds.New(db, encryptionKey),
	}
}
