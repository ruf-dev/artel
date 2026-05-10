package v1

import (
	"database/sql"

	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/repository/v1/couchcreds"
	"github.com/ruf-dev/artel/internal/repository/v1/subscriptions"
	"github.com/ruf-dev/artel/internal/repository/v1/users"
	"github.com/ruf-dev/artel/internal/repository/v1/vaults"
)

type Repos struct {
	Users          repository.Users
	Vaults         repository.Vaults
	Subscriptions  repository.Subscriptions
	CouchCredentials repository.CouchCredentials
}

func New(db *sql.DB, encryptionKey []byte) *Repos {
	return &Repos{
		Users:          users.New(db),
		Vaults:         vaults.New(db),
		Subscriptions:  subscriptions.New(db),
		CouchCredentials: couchcreds.New(db, encryptionKey),
	}
}
