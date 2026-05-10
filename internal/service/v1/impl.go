package v1

import (
	"github.com/ruf-dev/artel/internal/clients/couchdb"
	"github.com/ruf-dev/artel/internal/service"
	"github.com/ruf-dev/artel/internal/service/v1/users"
	"github.com/ruf-dev/artel/internal/service/v1/vault"
)

type Services struct {
	Vault service.VaultService
	User  service.UserService
}

func New(db *couchdb.Client) *Services {
	return &Services{
		Vault: vault.New(db),
		User:  users.New(db),
	}
}
