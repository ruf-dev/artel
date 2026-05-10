package v1

import (
	"github.com/ruf-dev/artel/internal/repository/pg"
	"github.com/ruf-dev/artel/internal/service"
	"github.com/ruf-dev/artel/internal/service/v1/auth"
	"github.com/ruf-dev/artel/internal/service/v1/couchinstances"
	"github.com/ruf-dev/artel/internal/service/v1/vault"
)

type Services struct {
	Auth          service.AuthService
	Vault         service.VaultService
	CouchInstance service.CouchInstanceService
}

func New(repo *pg.Repos) *Services {
	return &Services{
		Auth:          auth.New(repo.Users, repo.Sessions),
		Vault:         vault.New(repo.Vaults, repo.VaultMembers, repo.CouchAccounts, repo.CouchInstances),
		CouchInstance: couchinstances.New(repo.CouchInstances),
	}
}

func (s *Services) AuthService() service.AuthService {
	return s.Auth
}

func (s *Services) VaultService() service.VaultService {
	return s.Vault
}

func (s *Services) CouchInstanceService() service.CouchInstanceService {
	return s.CouchInstance
}
