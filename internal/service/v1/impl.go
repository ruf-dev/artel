package v1

import (
	"github.com/ruf-dev/artel/internal/clients/couchdb"
	"github.com/ruf-dev/artel/internal/service"
	"github.com/ruf-dev/artel/internal/service/v1/vault"
)

type Services struct {
	Vault service.VaultService
}

func New(db *couchdb.Client) *Services {
	return &Services{
		Vault: vault.New(db),
	}
}
