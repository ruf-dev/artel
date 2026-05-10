package v1

import (
	"github.com/ruf-dev/artel/internal/clients/couchdb"
	"github.com/ruf-dev/artel/internal/repository/pg"
	"github.com/ruf-dev/artel/internal/service"
	"github.com/ruf-dev/artel/internal/service/v1/users"
	"github.com/ruf-dev/artel/internal/service/v1/vault"
)

type Services struct {
	Vault service.VaultService
	User  service.UserService
}

func New(repo *pg.Repos) *Services {
	pool := couchdb.NewPool(repo.CouchCredentials)

	return &Services{
		Vault: vault.New(repo.Vaults, repo.CouchCredentials, pool),
		User:  users.New(repo.Users),
	}
}
