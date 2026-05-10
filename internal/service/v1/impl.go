package v1

import (
	"github.com/ruf-dev/artel/internal/clients/couchdb"
	"github.com/ruf-dev/artel/internal/repository/pg"
	"github.com/ruf-dev/artel/internal/service"
	"github.com/ruf-dev/artel/internal/service/v1/auth"
	"github.com/ruf-dev/artel/internal/service/v1/couchinstances"
	"github.com/ruf-dev/artel/internal/service/v1/users"
	"github.com/ruf-dev/artel/internal/service/v1/vault"
)

type Services struct {
	Auth          service.AuthService
	Vault         service.VaultService
	User          service.UserService
	CouchInstance service.CouchInstanceService
}

func New(repo *pg.Repos, signingKey []byte) *Services {

	pool := couchdb.NewPool(repo.CouchCredentials)

	return &Services{
		Auth:          auth.New(signingKey),
		Vault:         vault.New(repo.Vaults, repo.CouchCredentials, pool),
		User:          users.New(repo.Users),
		CouchInstance: couchinstances.New(repo.CouchInstances),
	}
}

func (s *Services) AuthService() service.AuthService {
	return s.Auth
}

func (s *Services) UserService() service.UserService {
	return s.User
}

func (s *Services) VaultService() service.VaultService {
	return s.Vault
}

func (s *Services) CouchInstanceService() service.CouchInstanceService {
	return s.CouchInstance
}
