package pg

import (
	"database/sql"

	"github.com/ruf-dev/artel/internal/repository"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"github.com/ruf-dev/artel/internal/repository/pg/repos/couchaccounts"
	"github.com/ruf-dev/artel/internal/repository/pg/repos/couchinstances"
	"github.com/ruf-dev/artel/internal/repository/pg/repos/emailaccounts"
	"github.com/ruf-dev/artel/internal/repository/pg/repos/mailserversuggestions"
	"github.com/ruf-dev/artel/internal/repository/pg/repos/mcpkeys"
	"github.com/ruf-dev/artel/internal/repository/pg/repos/pendingauthcodes"
	"github.com/ruf-dev/artel/internal/repository/pg/repos/prompts"
	"github.com/ruf-dev/artel/internal/repository/pg/repos/sessions"
	"github.com/ruf-dev/artel/internal/repository/pg/repos/subscriptions"
	"github.com/ruf-dev/artel/internal/repository/pg/repos/userpermissions"
	"github.com/ruf-dev/artel/internal/repository/pg/repos/users"
	"github.com/ruf-dev/artel/internal/repository/pg/repos/vaultmembers"
	"github.com/ruf-dev/artel/internal/repository/pg/repos/vaults"
	"github.com/ruf-dev/artel/internal/repository/pg/tx_manager"
)

type Repos struct {
	users                 repository.Users
	vaults                repository.Vaults
	vaultMembers          repository.VaultMembers
	sessions              repository.Sessions
	subscriptions         repository.Subscriptions
	couchAccounts         repository.CouchAccounts
	couchInstances        repository.CouchInstances
	userPermissions       repository.UserPermissionsRepo
	mcpKey                repository.McpKeyRepository
	pendingAuthCodes      repository.PendingAuthCodes
	emailAccounts         repository.EmailAccounts
	mailServerSuggestions repository.MailServerSuggestions
	promptsRepo           repository.Prompts

	txManager tx_manager.TxManager
}

func (r Repos) Vaults() repository.Vaults {
	return r.vaults
}

func (r Repos) VaultMembers() repository.VaultMembers {
	return r.vaultMembers
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

func (r Repos) EmailAccounts() repository.EmailAccounts {
	return r.emailAccounts
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

func New(db *sql.DB, encryptionKey []byte) *Repos {
	q := artel_q.New(db)

	return &Repos{
		vaults:         vaults.New(db),
		vaultMembers:   vaultmembers.New(db),
		couchInstances: couchinstances.New(db, encryptionKey),

		users:                 users.New(q),
		sessions:              sessions.New(q),
		subscriptions:         subscriptions.New(q),
		couchAccounts:         couchaccounts.New(db, encryptionKey),
		userPermissions:       userpermissions.New(q),
		mcpKey:                mcpkeys.New(q),
		pendingAuthCodes:      pendingauthcodes.New(q),
		emailAccounts:         emailaccounts.New(q, encryptionKey),
		mailServerSuggestions: mailserversuggestions.New(q),
		promptsRepo:           prompts.New(db),

		txManager: tx_manager.New(db),
	}
}
