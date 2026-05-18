package v1

import (
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/config"
	"github.com/ruf-dev/artel/internal/repository/pg"
	"github.com/ruf-dev/artel/internal/service"
	"github.com/ruf-dev/artel/internal/service/v1/auth"
	"github.com/ruf-dev/artel/internal/service/v1/couchinstances"
	"github.com/ruf-dev/artel/internal/service/v1/email"
	"github.com/ruf-dev/artel/internal/service/v1/mcp"
	"github.com/ruf-dev/artel/internal/service/v1/vault"
)

type Services struct {
	Auth          service.AuthService
	Vault         service.VaultService
	CouchInstance service.CouchInstanceService
	Mcp           service.McpService
	Email         service.EmailService
}

func New(repo *pg.Repos, cfg config.EnvironmentConfig) (*Services, error) {
	authSvc, err := auth.New(repo, cfg.TelegramClientID)
	if err != nil {
		return nil, rerrors.Wrap(err, "init auth service")
	}

	return &Services{
		Auth:          authSvc,
		Vault:         vault.New(repo),
		CouchInstance: couchinstances.New(repo),
		Mcp:           mcp.New(repo.McpKeyRepository(), repo.Vaults(), repo.CouchInstances()),
		Email:         email.New(repo.EmailAccounts()),
	}, nil
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

func (s *Services) McpService() service.McpService {
	return s.Mcp
}

func (s *Services) EmailService() service.EmailService {
	return s.Email
}
