package pg

import (
	"database/sql"

	"github.com/ruf-dev/artel/internal/repository"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"github.com/ruf-dev/artel/internal/repository/pg/repos/couchaccounts"
	"github.com/ruf-dev/artel/internal/repository/pg/repos/couchinstances"
	"github.com/ruf-dev/artel/internal/repository/pg/repos/externalconnections"
	"github.com/ruf-dev/artel/internal/repository/pg/repos/mailserversuggestions"
	"github.com/ruf-dev/artel/internal/repository/pg/repos/mcpconnectors"
	"github.com/ruf-dev/artel/internal/repository/pg/repos/mcpdefinitions"
	"github.com/ruf-dev/artel/internal/repository/pg/repos/mcpkeys"
	"github.com/ruf-dev/artel/internal/repository/pg/repos/mcpspreadsheets"
	"github.com/ruf-dev/artel/internal/repository/pg/repos/pendingauthcodes"
	"github.com/ruf-dev/artel/internal/repository/pg/repos/prompts"
	"github.com/ruf-dev/artel/internal/repository/pg/repos/s3instances"
	"github.com/ruf-dev/artel/internal/repository/pg/repos/sessions"
	"github.com/ruf-dev/artel/internal/repository/pg/repos/subscriptions"
	"github.com/ruf-dev/artel/internal/repository/pg/repos/tasktrackers"
	"github.com/ruf-dev/artel/internal/repository/pg/repos/tracts"
	"github.com/ruf-dev/artel/internal/repository/pg/repos/triggerpresets"
	"github.com/ruf-dev/artel/internal/repository/pg/repos/triggers"
	"github.com/ruf-dev/artel/internal/repository/pg/repos/userpermissions"
	"github.com/ruf-dev/artel/internal/repository/pg/repos/users"
	"github.com/ruf-dev/artel/internal/repository/pg/repos/vaultinvites"
	"github.com/ruf-dev/artel/internal/repository/pg/repos/vaultmembers"
	"github.com/ruf-dev/artel/internal/repository/pg/repos/vaults"
	"github.com/ruf-dev/artel/internal/repository/pg/tx_manager"
)

type Repos struct {
	users                 repository.Users
	vaults                repository.Vaults
	vaultMembers          repository.VaultMembers
	vaultInvitesRepo      repository.VaultInvites
	sessions              repository.Sessions
	subscriptions         repository.Subscriptions
	couchAccounts         repository.CouchAccounts
	couchInstances        repository.CouchInstances
	s3Instances           repository.S3Instances
	userPermissions       repository.UserPermissionsRepo
	mcpKey                repository.McpKeyRepository
	pendingAuthCodes      repository.PendingAuthCodes
	mailServerSuggestions repository.MailServerSuggestions
	promptsRepo           repository.Prompts
	taskTrackersRepo      repository.TaskTrackerRepo
	externalConnections   repository.ExternalConnectionRepo
	mcpSpreadsheets       repository.McpSpreadsheetsRepo
	mcpDefinitions        repository.McpDefinitionsRepo
	mcpConnectors         repository.McpConnectorsRepo
	tracts                repository.TractsRepo
	triggers              repository.TriggersRepo
	triggerPresets        repository.TriggerPresetsRepo

	txManager tx_manager.TxManager
}

func (r Repos) Vaults() repository.Vaults {
	return r.vaults
}

func (r Repos) VaultMembers() repository.VaultMembers {
	return r.vaultMembers
}

func (r Repos) VaultInvites() repository.VaultInvites {
	return r.vaultInvitesRepo
}

func (r Repos) Sessions() repository.Sessions {
	return r.sessions
}

func (r Repos) Subscriptions() repository.Subscriptions {
	return r.subscriptions
}

func (r Repos) CouchAccounts() repository.CouchAccounts {
	return r.couchAccounts
}

func (r Repos) CouchInstances() repository.CouchInstances {
	return r.couchInstances
}

func (r Repos) S3Instances() repository.S3Instances {
	return r.s3Instances
}

func (r Repos) Users() repository.Users {
	return r.users
}

func (r Repos) TxManager() tx_manager.TxManager {
	return r.txManager
}

func (r Repos) McpKeyRepository() repository.McpKeyRepository {
	return r.mcpKey
}

func (r Repos) PendingAuthCodes() repository.PendingAuthCodes {
	return r.pendingAuthCodes
}

func (r Repos) MailServerSuggestions() repository.MailServerSuggestions {
	return r.mailServerSuggestions
}

func (r Repos) UserPermissions() repository.UserPermissionsRepo {
	return r.userPermissions
}

func (r Repos) Prompts() repository.Prompts {
	return r.promptsRepo
}

func (r Repos) TaskTrackers() repository.TaskTrackerRepo {
	return r.taskTrackersRepo
}

func (r Repos) ExternalConnections() repository.ExternalConnectionRepo {
	return r.externalConnections
}

func (r Repos) McpSpreadsheets() repository.McpSpreadsheetsRepo {
	return r.mcpSpreadsheets
}

func (r Repos) McpDefinitions() repository.McpDefinitionsRepo {
	return r.mcpDefinitions
}

func (r Repos) McpConnectors() repository.McpConnectorsRepo {
	return r.mcpConnectors
}

func (r Repos) Tracts() repository.TractsRepo {
	return r.tracts
}

func (r Repos) Triggers() repository.TriggersRepo {
	return r.triggers
}

func (r Repos) TriggerPresets() repository.TriggerPresetsRepo {
	return r.triggerPresets
}

func New(db *sql.DB, encryptionKey []byte) *Repos {
	q := artel_q.New(newLoggingDB(db))
	txManager := tx_manager.New(db)

	return &Repos{
		vaults:           vaults.New(db, encryptionKey),
		vaultMembers:     vaultmembers.New(db),
		vaultInvitesRepo: vaultinvites.New(db),
		couchInstances:   couchinstances.New(db, encryptionKey),
		s3Instances:      s3instances.New(db, encryptionKey),

		users:                 users.New(q, db),
		sessions:              sessions.New(q),
		subscriptions:         subscriptions.New(q),
		couchAccounts:         couchaccounts.New(db, encryptionKey),
		userPermissions:       userpermissions.New(q),
		mcpKey:                mcpkeys.New(q),
		pendingAuthCodes:      pendingauthcodes.New(q),
		mailServerSuggestions: mailserversuggestions.New(q),
		promptsRepo:           prompts.New(db),
		taskTrackersRepo:      tasktrackers.New(q, encryptionKey),
		externalConnections:   externalconnections.New(q, encryptionKey),
		mcpSpreadsheets:       mcpspreadsheets.New(q),
		mcpDefinitions:        mcpdefinitions.New(q),
		mcpConnectors:         mcpconnectors.New(q),
		tracts:                tracts.New(q, txManager),
		triggers:              triggers.New(q),
		triggerPresets:        triggerpresets.New(q),

		txManager: txManager,
	}
}
