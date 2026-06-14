package v1

import (
	"golang.org/x/oauth2"
	googleoauth "golang.org/x/oauth2/google"

	"github.com/ruf-dev/artel/internal/config"
	"github.com/ruf-dev/artel/internal/repository/pg"
	"github.com/ruf-dev/artel/internal/service"
	"github.com/ruf-dev/artel/internal/service/v1/admincouchsvc"
	adminusers "github.com/ruf-dev/artel/internal/service/v1/adminusers"
	"github.com/ruf-dev/artel/internal/service/v1/auth"
	"github.com/ruf-dev/artel/internal/service/v1/couchinstances"
	"github.com/ruf-dev/artel/internal/service/v1/email"
	externalconns "github.com/ruf-dev/artel/internal/service/v1/externalconnections"
	"github.com/ruf-dev/artel/internal/service/v1/mcp"
	"github.com/ruf-dev/artel/internal/service/v1/notes"
	"github.com/ruf-dev/artel/internal/service/v1/prompt"
	"github.com/ruf-dev/artel/internal/service/v1/subscription"
	"github.com/ruf-dev/artel/internal/service/v1/tasktracker"
	"github.com/ruf-dev/artel/internal/service/v1/vault"
)

type Services struct {
	Auth                service.AuthService
	Vault               service.VaultService
	CouchInstance       service.CouchInstanceService
	AdminCouch          service.AdminCouchService
	Mcp                 service.McpService
	Email               service.EmailService
	Subscription        service.SubscriptionService
	Prompt              service.PromptService
	TaskTracker         service.TaskTrackerService
	Notes               service.NotesService
	AdminUsers          service.AdminUsersService
	ExternalConnections service.ExternalConnectionService
}

func New(repo *pg.Repos, cfg config.EnvironmentConfig) (*Services, error) {
	oauthCfg := &oauth2.Config{
		ClientID:     cfg.GoogleClientId,
		ClientSecret: cfg.GoogleClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/drive.file",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: googleoauth.Endpoint,
	}

	return &Services{
		Auth:                auth.New(repo, cfg.TelegramClientID),
		Vault:               vault.New(repo),
		CouchInstance:       couchinstances.New(repo),
		AdminCouch:          admincouchsvc.New(repo),
		Mcp:                 mcp.New(repo.McpKeyRepository(), repo.Vaults(), repo.VaultMembers(), repo.EmailAccounts(), repo.CouchInstances(), repo.UserPermissions()),
		Email:               email.New(repo.EmailAccounts(), repo.MailServerSuggestions()),
		Subscription:        subscription.New(repo.Subscriptions()),
		Prompt:              prompt.New(repo.Prompts()),
		TaskTracker:         tasktracker.New(repo.TaskTrackers()),
		Notes:               notes.New(repo),
		AdminUsers:          adminusers.New(repo.Users(), repo.Sessions()),
		ExternalConnections: externalconns.New(repo.ExternalConnections(), repo.PendingAuthCodes(), oauthCfg),
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

func (s *Services) AdminCouchService() service.AdminCouchService {
	return s.AdminCouch
}

func (s *Services) McpService() service.McpService {
	return s.Mcp
}

func (s *Services) EmailService() service.EmailService {
	return s.Email
}

func (s *Services) SubscriptionService() service.SubscriptionService {
	return s.Subscription
}

func (s *Services) PromptService() service.PromptService {
	return s.Prompt
}

func (s *Services) TaskTrackerService() service.TaskTrackerService {
	return s.TaskTracker
}

func (s *Services) NotesService() service.NotesService {
	return s.Notes
}

func (s *Services) AdminUsersService() service.AdminUsersService {
	return s.AdminUsers
}

func (s *Services) ExternalConnectionService() service.ExternalConnectionService {
	return s.ExternalConnections
}
