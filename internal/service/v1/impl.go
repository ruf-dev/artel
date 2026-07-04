package v1

import (
	"github.com/ruf-dev/artel/internal/config"
	"github.com/ruf-dev/artel/internal/repository/pg"
	"github.com/ruf-dev/artel/internal/service"
	"github.com/ruf-dev/artel/internal/service/v1/admincouchsvc"
	"github.com/ruf-dev/artel/internal/service/v1/adminusers"
	"github.com/ruf-dev/artel/internal/service/v1/auth"
	"github.com/ruf-dev/artel/internal/service/v1/couchinstances"
	externalconns "github.com/ruf-dev/artel/internal/service/v1/externalconnections"
	"github.com/ruf-dev/artel/internal/service/v1/mcp"
	"github.com/ruf-dev/artel/internal/service/v1/mom"
	"github.com/ruf-dev/artel/internal/service/v1/notes"
	"github.com/ruf-dev/artel/internal/service/v1/prompt"
	"github.com/ruf-dev/artel/internal/service/v1/s3instances"
	"github.com/ruf-dev/artel/internal/service/v1/subscription"
	"github.com/ruf-dev/artel/internal/service/v1/tasktracker"
	"github.com/ruf-dev/artel/internal/service/v1/vault"
	"golang.org/x/oauth2"
	googleoauth "golang.org/x/oauth2/google"
)

type Services struct {
	Auth                service.AuthService
	Vault               service.VaultService
	CouchInstance       service.CouchInstanceService
	S3Instance          service.S3InstanceService
	AdminCouch          service.AdminCouchService
	Mcp                 service.McpService
	Subscription        service.SubscriptionService
	Prompt              service.PromptService
	TaskTracker         service.TaskTrackerService
	Notes               service.NotesService
	AdminUsers          service.AdminUsersService
	ExternalConnections service.ExternalConnectionService
	Mom                 service.MomService
	// Tract is constructed in internal/app/custom.go (not here) — it depends on Mcp and Mom,
	// which must already exist as service.McpService/service.MomService values to compose the
	// tract.ToolExecutor without the tract package importing internal/service/v1/mcp.
	Tract service.TractService
}

func New(repo *pg.Repos, cfg config.EnvironmentConfig) (*Services, error) {
	oauthCfg := &oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/drive.file",
		},
		Endpoint: googleoauth.Endpoint,
	}

	return &Services{
		Auth:                auth.New(repo, cfg.TelegramClientID),
		Vault:               vault.New(repo),
		CouchInstance:       couchinstances.New(repo),
		S3Instance:          s3instances.New(repo),
		AdminCouch:          admincouchsvc.New(repo),
		Mcp:                 mcp.New(repo.McpKeyRepository(), repo.Vaults(), repo.VaultMembers(), repo.CouchInstances(), repo.S3Instances(), repo.McpConnectors(), repo.McpDefinitions(), repo.ExternalConnections()),
		Subscription:        subscription.New(repo.Subscriptions()),
		Prompt:              prompt.New(repo.Prompts()),
		TaskTracker:         tasktracker.New(repo.TaskTrackers()),
		Notes:               notes.New(repo),
		AdminUsers:          adminusers.New(repo.Users(), repo.Sessions()),
		ExternalConnections: externalconns.New(repo.ExternalConnections(), repo.PendingAuthCodes(), repo.McpSpreadsheets(), repo.MailServerSuggestions(), oauthCfg),
		Mom:                 mom.New(repo.McpDefinitions(), repo.McpConnectors(), repo.ExternalConnections()),
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

func (s *Services) S3InstanceService() service.S3InstanceService {
	return s.S3Instance
}

func (s *Services) AdminCouchService() service.AdminCouchService {
	return s.AdminCouch
}

func (s *Services) McpService() service.McpService {
	return s.Mcp
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

func (s *Services) MomService() service.MomService {
	return s.Mom
}

func (s *Services) TractService() service.TractService {
	return s.Tract
}
