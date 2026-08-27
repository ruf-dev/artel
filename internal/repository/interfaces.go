package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/clients/postgres"
	"github.com/ruf-dev/artel/internal/domain"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"github.com/ruf-dev/artel/internal/repository/pg/tx_manager"
)

type Repo interface {
	Users() Users
	Vaults() Vaults
	Workbenches() Workbenches
	VaultMembers() VaultMembers
	VaultInvites() VaultInvites
	Sessions() Sessions
	Subscriptions() Subscriptions
	SubscriptionPlans() SubscriptionPlansRepo
	CouchAccounts() CouchAccounts
	CouchInstances() CouchInstances
	S3Instances() S3Instances
	PostgresInstances() PostgresInstances
	VaultPostgresDatabases() VaultPostgresDatabases
	DockerHosts() DockerHosts
	UserPermissions() UserPermissionsRepo
	McpKeyRepository() McpKeyRepository
	PendingAuthCodes() PendingAuthCodes
	MailServerSuggestions() MailServerSuggestions
	Prompts() Prompts
	ExternalConnections() ExternalConnectionRepo
	McpSpreadsheets() McpSpreadsheetsRepo
	McpDefinitions() McpDefinitionsRepo
	McpConnectors() McpConnectorsRepo
	Tracts() TractsRepo
	TractTemplates() TractTemplatesRepo
	Triggers() TriggersRepo
	TriggerPresets() TriggerPresetsRepo
	SystemSettings() SystemSettingsRepo
	UserSettings() UserSettingsRepo

	TxManager() tx_manager.TxManager
}

type ListPromptsParams struct {
	Ids    []string
	Limit  uint32
	Offset uint32
}

type Prompts interface {
	List(ctx context.Context, params ListPromptsParams) ([]domain.Prompt, int64, error)
}

type Users interface {
	Create(ctx context.Context, email, passwordHash string) (domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (sql.Null[domain.User], error)

	GetByTelegramId(ctx context.Context, telegramId string) (sql.Null[domain.User], error)
	CreateByUsername(ctx context.Context, username, photoUrl string) (domain.User, error)
	UpsertTelegramIdentity(ctx context.Context, identity domain.TelegramIdentity) error
	GetTelegramPhotoUrl(ctx context.Context, userUuid uuid.UUID) (string, error)
	UpdatePhotoUrl(ctx context.Context, userUuid uuid.UUID, photoUrl string) error
	UpdatePasswordHash(ctx context.Context, userUuid uuid.UUID, passwordHash string) error

	ListAll(ctx context.Context, req domain.ListUsersReq) ([]domain.User, int64, error)
	GetDetailsById(ctx context.Context, id uuid.UUID) (domain.UserDetails, error)

	Delete(ctx context.Context, id uuid.UUID) error

	WithTx(tx *sql.Tx) Users
}

type Vaults interface {
	Upsert(
		ctx context.Context, userID, couchInstanceID uuid.UUID, name, couchDBName, status, passphrase string,
	) (domain.Vault, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.Vault, error)
	GetByNameAndUser(ctx context.Context, userID uuid.UUID, name string) (domain.Vault, error)
	UpdateStatus(ctx context.Context, vaultID uuid.UUID, status string) error
	SetLiveSyncPassphrase(ctx context.Context, vaultID uuid.UUID, passphrase string) error
	ListByMembership(ctx context.Context, userID uuid.UUID) ([]domain.Vault, error)
	Delete(ctx context.Context, vaultID uuid.UUID) error

	LinkS3Bucket(ctx context.Context, vaultID, s3InstanceID uuid.UUID, bucketName string) error
	UnlinkS3Bucket(ctx context.Context, vaultID uuid.UUID) error

	SetUseCouchDBForBinaries(ctx context.Context, vaultID uuid.UUID, value bool) error

	PublishVault(ctx context.Context, vaultID uuid.UUID, slug string) (domain.Vault, error)
	UnpublishVault(ctx context.Context, vaultID uuid.UUID) error
	GetBySlug(ctx context.Context, slug string) (domain.Vault, error)

	WithTx(tx postgres.DB) Vaults
}

// Workbenches is the pure-DB layer for the per-(vault, user) Docker workbench container: its
// Mark* methods below correspond to the domain.WorkbenchStatus transitions (created →
// configuring → running → stopped → removed).
type Workbenches interface {
	Create(ctx context.Context, vaultID, userID uuid.UUID, volumeName string, dockerHostID uuid.UUID) (domain.Workbench, error)
	GetByVaultAndUser(ctx context.Context, vaultID, userID uuid.UUID) (domain.Workbench, error)
	MarkContainerCreated(ctx context.Context, vaultID, userID uuid.UUID, containerID string) error
	MarkConfiguring(ctx context.Context, vaultID, userID uuid.UUID) error
	MarkRunning(ctx context.Context, vaultID, userID uuid.UUID, authMode domain.WorkbenchAuthMode) error
	MarkStopped(ctx context.Context, vaultID, userID uuid.UUID) error
	MarkRemoved(ctx context.Context, vaultID, userID uuid.UUID) error
	Delete(ctx context.Context, vaultID, userID uuid.UUID) error

	// ListByVaultID returns every workbench row for vaultID across all members — used only by
	// vault deletion, which must tear down every member's Docker resources, not just one.
	ListByVaultID(ctx context.Context, vaultID uuid.UUID) ([]domain.Workbench, error)

	// GetMostRecentByUser returns the most recently started (falling back to most recently
	// created) of userID's workbenches across every vault they're a member of — used by
	// internal/transport/telegram_webhook to pick a single workbench to relay chat into (v1
	// scope: no vault picker). sql.Null[domain.Workbench]{Valid: false} means the user has no
	// workbench at all yet.
	GetMostRecentByUser(ctx context.Context, userID uuid.UUID) (sql.Null[domain.Workbench], error)

	// GetContentSnapshot returns the path→mtime baseline captured at the (vault, user)
	// workbench's last materialization, or a nil map if it has never been started.
	GetContentSnapshot(ctx context.Context, vaultID, userID uuid.UUID) (map[string]int64, error)
	// SetContentSnapshot overwrites the (vault, user) workbench's path→mtime baseline — called
	// once per materialization (workbench start).
	SetContentSnapshot(ctx context.Context, vaultID, userID uuid.UUID, snapshot map[string]int64) error

	WithTx(tx postgres.DB) Workbenches
}

type VaultMembers interface {
	Add(ctx context.Context, vaultID, userID uuid.UUID, role artel_q.VaultRole) error
	Remove(ctx context.Context, vaultID, userID uuid.UUID) error
	Get(ctx context.Context, vaultID, userID uuid.UUID) (domain.VaultMember, error)
	ListByVault(ctx context.Context, vaultID uuid.UUID) ([]domain.VaultMember, error)
	ListByVaultWithUsers(ctx context.Context, vaultID uuid.UUID) ([]domain.VaultMemberInfo, error)

	WithTx(tx postgres.DB) VaultMembers
}

type VaultInvites interface {
	Create(
		ctx context.Context, vaultID, createdBy uuid.UUID, role artel_q.VaultRole, token string,
	) (domain.VaultInvite, error)
	GetByToken(ctx context.Context, token string) (domain.VaultInvite, error)
	ListByVault(ctx context.Context, vaultID uuid.UUID) ([]domain.VaultInvite, error)
	Revoke(ctx context.Context, id uuid.UUID) error

	WithTx(tx postgres.DB) VaultInvites
}

type Sessions interface {
	Create(ctx context.Context, session domain.Session) (domain.Session, error)
	GetByToken(ctx context.Context, token string) (domain.Session, error)
	GetByTokenWithUser(ctx context.Context, token string) (domain.Session, domain.User, error)
	Delete(ctx context.Context, token string) error
	GetByUserID(ctx context.Context, userUuid uuid.UUID) ([]domain.Session, error)
	// RotateByRefreshToken atomically swaps the session identified by oldRefreshToken for the
	// token/expiry fields in newSession, provided the old refresh token has not already expired
	// or been rotated away. Returns Valid: false if no matching row was updated.
	RotateByRefreshToken(ctx context.Context, oldRefreshToken string, newSession domain.Session) (sql.Null[domain.Session], error)
}

type Subscriptions interface {
	Upsert(ctx context.Context, userID uuid.UUID, active bool) (domain.Subscription, error)
	GetByUser(ctx context.Context, userID uuid.UUID) (domain.Subscription, error)
	CreateDefault(ctx context.Context, userID uuid.UUID) error
	// GetWithPlan returns the merged plan+override view — the plan's defaults with the user's
	// feature/quota overrides applied on top.
	GetWithPlan(ctx context.Context, userID uuid.UUID) (domain.EffectiveSubscription, error)
	// UpsertOverrides overwrites a user's plan assignment and overrides in place — used by the
	// migration backfill and any future admin override tooling, not a live per-request path.
	UpsertOverrides(ctx context.Context, sub domain.Subscription) (domain.Subscription, error)

	WithTx(tx *sql.Tx) Subscriptions
}

type SubscriptionPlansRepo interface {
	Get(ctx context.Context, planKey string) (domain.SubscriptionPlan, error)
	List(ctx context.Context) ([]domain.SubscriptionPlan, error)

	WithTx(tx *sql.Tx) SubscriptionPlansRepo
}

type CouchAccounts interface {
	Upsert(
		ctx context.Context, userID, instanceID uuid.UUID, username string, passwordPlain string,
	) (domain.CouchAccount, error)
	GetByUserAndInstance(ctx context.Context, userID, instanceID uuid.UUID) (domain.CouchAccount, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.CouchAccount, error)
	UpdatePassword(ctx context.Context, username string, instanceID uuid.UUID, passwordPlain string) error
	Delete(ctx context.Context, id uuid.UUID) error

	WithTx(tx postgres.DB) CouchAccounts
}

type CouchInstances interface {
	Register(ctx context.Context, url, username string, passwordPlain []byte) (uuid.UUID, error)
	Get(ctx context.Context, id uuid.UUID) (domain.CouchInstance, error)
	RandomPick(ctx context.Context) (domain.CouchInstanceWithAccount, error)
	// PickForUser resolves storage instance selection for vault creation: prefers an instance
	// userID owns (a BYOK couchdb connection), falling back to RandomPick's shared admin pool
	// when they have none.
	PickForUser(ctx context.Context, userID uuid.UUID) (domain.CouchInstanceWithAccount, error)
	// GetOwned returns the instance owned by userID, if any — used by the service layer to
	// decide between RegisterOwned (insert) and Update (already owns one) when syncing a BYOK
	// couchdb connection.
	GetOwned(ctx context.Context, userID uuid.UUID) (sql.Null[domain.CouchInstance], error)
	List(ctx context.Context) ([]domain.CouchInstance, error)
	Update(ctx context.Context, id uuid.UUID, url, username string, passwordPlain []byte) error
	// RegisterOwned is like Register but stamps owner_user_id, marking the row as ownerUserID's
	// BYOK instance rather than an admin pool entry.
	RegisterOwned(
		ctx context.Context, ownerUserID uuid.UUID, url, username string, passwordPlain []byte,
	) (uuid.UUID, error)
	Delete(ctx context.Context, id uuid.UUID) error
	// DeleteOwnedIfUnreferenced removes ownerUserID's owned instance row, but only if no vault
	// currently references it — called on BYOK couchdb disconnect. No-op if ownerUserID has no
	// owned row.
	DeleteOwnedIfUnreferenced(ctx context.Context, ownerUserID uuid.UUID) error
	Exists(ctx context.Context) (bool, error)

	WithTx(tx postgres.DB) CouchInstances
}

type S3Instances interface {
	Register(
		ctx context.Context, endpoint, region string, useSSL, pathStyle bool, accessKey string, secretKeyPlain []byte,
	) (uuid.UUID, error)
	Get(ctx context.Context, id uuid.UUID) (domain.S3Instance, error)
	// PickForUser resolves storage instance selection for vault S3 linking: prefers an instance
	// userID owns (a BYOK s3 connection), falling back to a random pick from the shared admin
	// pool when they have none.
	PickForUser(ctx context.Context, userID uuid.UUID) (domain.S3Instance, error)
	// GetOwned returns the instance owned by userID, if any — used by the service layer to
	// decide between RegisterOwned (insert) and Update (already owns one) when syncing a BYOK
	// s3 connection.
	GetOwned(ctx context.Context, userID uuid.UUID) (sql.Null[domain.S3Instance], error)
	List(ctx context.Context) ([]domain.S3Instance, error)
	Update(
		ctx context.Context, id uuid.UUID, endpoint, region string,
		useSSL, pathStyle bool, accessKey string, secretKeyPlain []byte,
	) error
	// RegisterOwned is like Register but stamps owner_user_id, marking the row as ownerUserID's
	// BYOK instance rather than an admin pool entry.
	RegisterOwned(
		ctx context.Context, ownerUserID uuid.UUID, endpoint, region string,
		useSSL, pathStyle bool, accessKey string, secretKeyPlain []byte,
	) (uuid.UUID, error)
	Delete(ctx context.Context, id uuid.UUID) error
	// DeleteOwnedIfUnreferenced removes ownerUserID's owned instance row, but only if no vault
	// currently references it — called on BYOK s3 disconnect. No-op if ownerUserID has no owned
	// row.
	DeleteOwnedIfUnreferenced(ctx context.Context, ownerUserID uuid.UUID) error
	Exists(ctx context.Context) (bool, error)

	WithTx(tx postgres.DB) S3Instances
}

// PostgresInstances is the pure-DB layer for the admin-managed pool of Postgres servers (admin
// pool or BYOK) that per-vault databases are provisioned on — mirrors CouchInstances/S3Instances.
// Unlike those two, domain.PostgresInstance surfaces OwnerUserUuid directly (per its doc comment),
// so Get/List/RandomPick/etc. here populate it instead of leaving owner-scoping entirely to
// GetOwned/RegisterOwned.
type PostgresInstances interface {
	Register(
		ctx context.Context, host string, port int, adminDatabase, username string, passwordPlain []byte, sslMode string,
	) (uuid.UUID, error)
	Get(ctx context.Context, id uuid.UUID) (domain.PostgresInstance, error)
	RandomPick(ctx context.Context) (domain.PostgresInstance, error)
	// PickForUser resolves storage instance selection for vault postgres provisioning: prefers an
	// instance userID owns (a BYOK postgres connection), falling back to RandomPick's shared admin
	// pool when they have none.
	PickForUser(ctx context.Context, userID uuid.UUID) (domain.PostgresInstance, error)
	// GetOwned returns the instance owned by userID, if any — used by the service layer to
	// decide between RegisterOwned (insert) and Update (already owns one) when syncing a BYOK
	// postgres connection.
	GetOwned(ctx context.Context, userID uuid.UUID) (sql.Null[domain.PostgresInstance], error)
	List(ctx context.Context) ([]domain.PostgresInstance, error)
	Update(
		ctx context.Context, id uuid.UUID, host string, port int, adminDatabase, username string,
		passwordPlain []byte, sslMode string,
	) error
	// RegisterOwned is like Register but stamps owner_user_id, marking the row as ownerUserID's
	// BYOK instance rather than an admin pool entry.
	RegisterOwned(
		ctx context.Context, ownerUserID uuid.UUID, host string, port int, adminDatabase, username string,
		passwordPlain []byte, sslMode string,
	) (uuid.UUID, error)
	Delete(ctx context.Context, id uuid.UUID) error
	// DeleteOwnedIfUnreferenced removes ownerUserID's owned instance row, but only if no vault
	// postgres database currently references it — called on BYOK postgres disconnect. No-op if
	// ownerUserID has no owned row.
	DeleteOwnedIfUnreferenced(ctx context.Context, ownerUserID uuid.UUID) error
	Exists(ctx context.Context) (bool, error)

	WithTx(tx postgres.DB) PostgresInstances
}

// VaultPostgresDatabases is the pure-DB layer for the per-vault Postgres database+role
// provisioned on a PostgresInstance — mirrors the Workbench provisioning→ready/error lifecycle
// (see domain.VaultPostgresDatabase).
type VaultPostgresDatabases interface {
	Create(
		ctx context.Context, vaultID, postgresInstanceID uuid.UUID, dbName, roleUsername string, rolePasswordPlain []byte,
	) (domain.VaultPostgresDatabase, error)
	// GetByVaultID returns Valid: false rather than an error when vaultID has no Postgres database
	// provisioned — a vault without Postgres enabled is a normal state, not an error.
	GetByVaultID(ctx context.Context, vaultID uuid.UUID) (sql.Null[domain.VaultPostgresDatabase], error)
	MarkReady(ctx context.Context, vaultID uuid.UUID) error
	MarkError(ctx context.Context, vaultID uuid.UUID, message string) error
	Delete(ctx context.Context, vaultID uuid.UUID) error

	WithTx(tx postgres.DB) VaultPostgresDatabases
}

// DockerHosts is the pure-DB layer for the admin-managed pool of Docker daemons that back
// per-vault workbench containers. Unlike CouchInstances/S3Instances there's no single credential
// blob; instead there are three optional TLS/mTLS fields (migrations/062_docker_hosts_tls.sql)
// for the remote-daemon case.
type DockerHosts interface {
	Register(ctx context.Context, url, caCert, clientCert, clientKey string) (uuid.UUID, error)
	Get(ctx context.Context, id uuid.UUID) (domain.DockerHost, error)
	// GetWithCreds is like Get but also decrypts and populates the TLS fields — used only to
	// build a TLS-configured Docker client, never for admin list/view.
	GetWithCreds(ctx context.Context, id uuid.UUID) (domain.DockerHost, error)
	List(ctx context.Context) ([]domain.DockerHost, error)
	// Update patches url unconditionally; caCert/clientCert/clientKey are three-way patch
	// pointers: nil leaves the stored cert untouched, a pointer to "" clears it, a pointer to a
	// non-empty PEM re-encrypts and stores it.
	Update(ctx context.Context, id uuid.UUID, url string, caCert, clientCert, clientKey *string) error
	Delete(ctx context.Context, id uuid.UUID) error
	Exists(ctx context.Context) (bool, error)
	// PickLeastLoaded returns the docker host currently backing the fewest live workbenches —
	// used by workbench.Service.CreateWorkbench to spread new workbenches across the pool.
	PickLeastLoaded(ctx context.Context) (domain.DockerHost, error)

	WithTx(tx postgres.DB) DockerHosts
}

type UserPermissionsRepo interface {
	Get(ctx context.Context, userUuid uuid.UUID) (domain.UserPermissions, error)
	Upsert(ctx context.Context, userUuid uuid.UUID, isAdmin bool) (domain.UserPermissions, error)
	CreateDefault(ctx context.Context, userUuid uuid.UUID) error

	WithTx(tx *sql.Tx) UserPermissionsRepo
}

type ExternalConnectionRepo interface {
	Upsert(ctx context.Context, conn domain.ExternalConnection) (domain.ExternalConnection, error)
	Insert(ctx context.Context, conn domain.ExternalConnection) (domain.ExternalConnection, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.ExternalConnection, error)
	GetByUserAndProvider(
		ctx context.Context, userUuid uuid.UUID, provider string,
	) (sql.Null[domain.ExternalConnection], error)
	ListByUser(ctx context.Context, userUuid uuid.UUID) ([]domain.ExternalConnection, error)
	Delete(ctx context.Context, userUuid uuid.UUID, provider string) error
	DeleteByID(ctx context.Context, userUuid uuid.UUID, id uuid.UUID) error
}

type McpSpreadsheetsRepo interface {
	Insert(ctx context.Context, spreadsheet domain.McpSpreadsheet) (domain.McpSpreadsheet, error)
	ListByUser(ctx context.Context, userUuid uuid.UUID) ([]domain.McpSpreadsheet, error)
	Delete(ctx context.Context, userUuid uuid.UUID, spreadsheetId string) error
}

type McpKeyRepository interface {
	CreateMcpKey(
		ctx context.Context, vaultID, userID, keyID uuid.UUID, name string, keyHash []byte, keyPreview string,
	) (domain.McpKey, error)
	ListMcpKeysByVault(ctx context.Context, vaultID uuid.UUID) ([]domain.McpKey, error)
	ListMcpKeysByUser(ctx context.Context, userUuid uuid.UUID) ([]domain.McpKey, error)
	GetMcpKeyByID(ctx context.Context, id uuid.UUID) (domain.McpKey, error)
	ListActiveMcpKeys(ctx context.Context, vaultID uuid.UUID) ([]domain.McpKey, error)
	RevokeMcpKey(ctx context.Context, id uuid.UUID) error
	SetMcpKeyAccess(ctx context.Context, keyUuid, userUuid, vaultUuid uuid.UUID) error
	TouchLastAccessed(ctx context.Context, keyUuid uuid.UUID) error
}

type PendingAuthCodes interface {
	Create(ctx context.Context, code, rawToken, codeChallenge, redirectUri, clientId string, expiresAt time.Time) error
	Get(ctx context.Context, code string) (domain.PendingAuthCode, error)
	Delete(ctx context.Context, code string) error
	DeleteExpired(ctx context.Context) error
}

type MailServerSuggestions interface {
	ListByDomain(ctx context.Context, domainPrefix string) ([]domain.MailServerSuggestion, error)
}

type McpDefinitionsRepo interface {
	Upsert(ctx context.Context, def domain.McpDefinition) (domain.McpDefinition, error)
	Get(ctx context.Context, name string) (sql.Null[domain.McpDefinition], error)
	List(ctx context.Context) ([]domain.McpDefinition, error)
	Delete(ctx context.Context, name string) error
	GetTool(ctx context.Context, mcpName string, toolName string) (sql.Null[domain.McpToolDef], error)
	ListAllTools(ctx context.Context) ([]domain.McpToolRef, error)
}

type McpConnectorsRepo interface {
	Insert(ctx context.Context, connector domain.McpConnector) (domain.McpConnector, error)
	Get(ctx context.Context, mcpKeyUuid uuid.UUID, mcpName string) (sql.Null[domain.McpConnector], error)
	ListByKey(ctx context.Context, mcpKeyUuid uuid.UUID) ([]domain.McpConnector, error)
	Delete(ctx context.Context, mcpKeyUuid uuid.UUID, mcpName string) error
}

// TractsRepo is the pure-DB layer for tracts + their runs/run-steps. The step tree itself is
// NOT a separate table — Tract.Definition is persisted as one JSONB column.
type TractsRepo interface {
	Create(ctx context.Context, tract domain.Tract) (domain.Tract, error)
	Get(ctx context.Context, id uuid.UUID) (sql.Null[domain.Tract], error)
	ListByUser(ctx context.Context, userUuid uuid.UUID) ([]domain.Tract, error)
	// Update overwrites name/description/definition.
	Update(ctx context.Context, tract domain.Tract) (domain.Tract, error)
	SetEnabled(ctx context.Context, id uuid.UUID, enabled bool) error
	Delete(ctx context.Context, id uuid.UUID) error

	// InsertRun persists a tract_runs row before the engine walks the step tree
	// (persist-before-apply).
	InsertRun(ctx context.Context, run domain.TractRun) (domain.TractRun, error)
	GetRun(ctx context.Context, id uuid.UUID) (sql.Null[domain.TractRun], error)
	ListRunsByTract(ctx context.Context, tractUuid uuid.UUID, limit int32) ([]domain.TractRun, error)
	UpdateRunStatus(ctx context.Context, id uuid.UUID, status domain.TractRunStatus, errMsg string) error

	// InsertRunStep persists a running row for one executed step before it executes
	// (persist-before-apply, append-only history).
	InsertRunStep(ctx context.Context, step domain.TractRunStep) (domain.TractRunStep, error)
	UpdateRunStepFinish(
		ctx context.Context, id uuid.UUID, status domain.TractRunStepStatus, output json.RawMessage, errMsg string,
	) error
	ListRunStepsByRun(ctx context.Context, runUuid uuid.UUID) ([]domain.TractRunStep, error)

	// SweepStaleRuns/SweepStaleRunSteps mark stale 'running' rows 'failed' — call once at app init.
	SweepStaleRuns(ctx context.Context, threshold time.Time) error
	SweepStaleRunSteps(ctx context.Context, threshold time.Time) error
}

// TractTemplatesRepo is the pure-DB layer for published tract templates — immutable snapshots
// of a Tract's definition, browsable/copyable by any user. See domain.TractTemplate.
type TractTemplatesRepo interface {
	Create(ctx context.Context, template domain.TractTemplate) (domain.TractTemplate, error)
	Get(ctx context.Context, id uuid.UUID) (sql.Null[domain.TractTemplate], error)
	// ListAll returns published templates ordered by install_count desc, published_at desc.
	// category == "" means no filter.
	ListAll(ctx context.Context, category string) ([]domain.TractTemplate, error)
	ListByOwner(ctx context.Context, ownerUuid uuid.UUID) ([]domain.TractTemplate, error)
	Delete(ctx context.Context, id uuid.UUID) error
	IncrementInstallCount(ctx context.Context, id uuid.UUID) error
}

// TriggersRepo is the pure-DB layer for standalone triggers and their tract links.
type TriggersRepo interface {
	Create(ctx context.Context, trigger domain.Trigger) (domain.Trigger, error)
	// Get looks up a trigger by its stable primary key (owner-facing CRUD).
	Get(ctx context.Context, id uuid.UUID) (sql.Null[domain.Trigger], error)
	// GetByTriggerUuid looks up a trigger by its rotatable webhook routing id — used by the
	// inbound webhook handler, which only knows the routing id embedded in the fired URL.
	GetByTriggerUuid(ctx context.Context, triggerUuid uuid.UUID) (sql.Null[domain.Trigger], error)
	ListByUser(ctx context.Context, userUuid uuid.UUID) ([]domain.Trigger, error)
	SetEnabled(ctx context.Context, id uuid.UUID, enabled bool) error
	Delete(ctx context.Context, id uuid.UUID) error
	// RotateSecret overwrites trigger_uuid, secret_hash, and token_suffix in place, keyed by the
	// trigger's stable primary key id — invalidates the trigger's current webhook URL/token
	// without touching anything else about the trigger or its tract links.
	RotateSecret(ctx context.Context, id uuid.UUID, newTriggerUuid uuid.UUID, secretHash []byte, tokenSuffix string) (domain.Trigger, error)

	Link(ctx context.Context, link domain.TriggerTractLink) error
	Unlink(ctx context.Context, triggerUuid uuid.UUID, tractUuid uuid.UUID) error
	// ListLinksByTract returns tractUuid's linked triggers with Trigger populated (Tract left
	// zero) — the tract editor's "wired up triggers" view.
	ListLinksByTract(ctx context.Context, tractUuid uuid.UUID) ([]TractTriggerLink, error)
	// ListLinksByTrigger returns triggerUuid's linked tracts with Tract populated (Trigger left
	// zero) — the webhook handler's fan-out: one delivery may start runs on several tracts.
	ListLinksByTrigger(ctx context.Context, triggerUuid uuid.UUID) ([]TractTriggerLink, error)

	// LinkToProvider attaches triggerId to a shared provider connection (see
	// trigger_provider_links) instead of it minting its own trigger_uuid/secret_hash webhook URL.
	LinkToProvider(ctx context.Context, triggerId uuid.UUID, externalConnectionId uuid.UUID) error
	// ListByExternalConnection returns every trigger (Matchers populated) linked to one shared
	// provider connection - the gitlab_webhook handler's fan-out lookup.
	ListByExternalConnection(ctx context.Context, externalConnectionId uuid.UUID) ([]domain.Trigger, error)
}

// TriggerPresetsRepo is the pure-DB read layer for the trigger_presets catalog table - see
// domain.TriggerPreset doc comment.
type TriggerPresetsRepo interface {
	List(ctx context.Context) ([]domain.TriggerPreset, error)
	GetByKey(ctx context.Context, key string) (sql.Null[domain.TriggerPreset], error)
}

// SystemSettingsRepo is the pure-DB layer for the single-row (id=1) global instance
// configuration — see domain.SystemSettings and migrations/064_system_settings.sql.
type SystemSettingsRepo interface {
	Get(ctx context.Context) (domain.SystemSettings, error)
	GetForUpdate(ctx context.Context) (domain.SystemSettings, error)
	UpdateAuthMethods(ctx context.Context, passwordEnabled, telegramEnabled bool) error
	UpdateRegistrationMode(ctx context.Context, mode domain.RegistrationMode) error
	SetSetupToken(ctx context.Context, tokenHash string, issuedAt time.Time) error
	CompleteSetup(ctx context.Context) error
	UpdateDefaultDocsVault(ctx context.Context, vaultUuid *uuid.UUID) error
	UpdateDefaultDocsSource(ctx context.Context, source domain.DocsSource) error

	WithTx(tx *sql.Tx) SystemSettingsRepo
}

type UserSettingsRepo interface {
	Get(ctx context.Context, userID uuid.UUID) (domain.UserSettings, error)
}

// TractTriggerLink is a read-side join projection — see TriggersRepo doc comments above for
// which field each List method populates.
type TractTriggerLink struct {
	TractUuid   uuid.UUID
	TriggerUuid uuid.UUID
	Trigger     domain.Trigger
	Tract       domain.Tract
	Filters     []domain.TractCondition
}
