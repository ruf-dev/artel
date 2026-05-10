package v1

import (
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service"
	"github.com/ruf-dev/artel/internal/service/v1/users"
	"github.com/ruf-dev/artel/internal/service/v1/vault"
)

type Services struct {
	Vault service.VaultService
	User  service.UserService
}

func New(db repository.Repo) *Services {
	return &Services{
		Vault: vault.New(db),
		User:  users.New(db),
	}
}
