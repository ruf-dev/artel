package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	anthropicClient "github.com/ruf-dev/artel/internal/clients/anthropic"
	"github.com/ruf-dev/artel/internal/clients/couchdb"
	"github.com/ruf-dev/artel/internal/clients/googleapi"
	openaiClient "github.com/ruf-dev/artel/internal/clients/openai"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
)

type Service interface {
	AuthService() AuthService
	VaultService() VaultService
	CouchInstanceService() CouchInstanceService
	S3InstanceService() S3InstanceService
	AdminCouchService() AdminCouchService
	McpService() McpService
	SubscriptionService() SubscriptionService
	PromptService() PromptService
	TaskTrackerService() TaskTrackerService
	NotesService() NotesService
	AdminUsersService() AdminUsersService
	AdminSubscriptionsService() AdminSubscriptionsService
	ExternalConnectionService() ExternalConnectionService
	MomService() MomService
	TractService() TractService
	PublicDocsService() PublicDocsService
	SetupWizardService() SetupWizardService
	AdminSystemSettingsService() AdminSystemSettingsService
	PostgresInstanceService() PostgresInstanceService
	SkillsService() SkillsService
}

type AdminUsersService interface {
	ListUsers(ctx context.Context, req domain.ListUsersReq) ([]domain.User, int64, error)
	GetUser(ctx context.Context, userUuid uuid.UUID) (domain.UserDetails, error)
	GetUserSessions(ctx context.Context, userUuid uuid.UUID) ([]domain.Session, error)
	// CreateUser creates a new non-admin account from the admin users page — the caller is
	// already authenticated and authorized as an administrator at the transport level, so this
	// bypasses AuthService.Register's setup/registration-mode/auth-method checks.
	CreateUser(ctx context.Context, email, password string) (domain.User, error)
	// ChangePassword resets userUuid's password — used by an administrator on the admin users
	// page.
	ChangePassword(ctx context.Context, userUuid uuid.UUID, newPassword string) error
}

// AdminSystemSettingsService lets an administrator inspect and edit the single-row global
// instance configuration (auth methods, registration mode) outside the first-run setup wizard —
// see SetupWizardService for the wizard's own (pre-setup, unauthenticated) version of the same
// underlying settings.
type AdminSystemSettingsService interface {
	// GetSettings also resolves the current default docs vault (domain.Vault{}, zero value, when
	// SystemSettings.DefaultDocsVaultID is unset) so the transport layer can surface its
	// name/slug alongside the raw id without a second round trip.
	GetSettings(ctx context.Context) (domain.SystemSettings, domain.Vault, error)
	// UpdateAuthMethods returns user_errors.AtLeastOneAuthMethodRequired if both are false.
	UpdateAuthMethods(ctx context.Context, passwordEnabled, telegramEnabled bool) error
	UpdateRegistrationMode(ctx context.Context, mode domain.RegistrationMode) error
	// UpdateDefaultDocsVault sets the instance-wide default `/docs` source. When source is
	// DocsSourceGithub, vaultUuid is ignored and the vault-existence check is skipped entirely —
	// the default routes to the built-in GitHub quick-start guide instead. When source is
	// DocsSourceVault (the zero-value/unset case too), a nil vaultUuid clears the default vault
	// and a non-nil vaultUuid must resolve to an existing vault (user_errors.NotFound otherwise) —
	// it is not required to already be published; see PublicDocsService.GetDefaultVault for the
	// read-time IsPublic enforcement.
	UpdateDefaultDocsVault(ctx context.Context, vaultUuid *uuid.UUID, source domain.DocsSource) error
}

// SetupWizardService drives the first-run setup wizard — a short-lived, unauthenticated flow
// that runs once (SystemSettings.SetupCompleted == false) to pick auth methods/registration
// mode and create the first administrator account. Every method except CurrentStatus and
// RegenerateSetupToken requires a wizardSessionToken minted by SubmitToken, which proves
// possession of the one-time setup token printed to server logs/stdout on first boot.
// Wizard sessions are in-memory only (not persisted) and expire after 30 minutes.
type SetupWizardService interface {
	// CurrentStatus returns the instance's current setup/auth-policy state — used by the
	// frontend to decide whether to show the wizard at all.
	CurrentStatus(ctx context.Context) (domain.SystemSettings, error)
	// RegenerateSetupToken mints a new one-time setup token, invalidating any previous one, and
	// returns it in plaintext (only its hash is persisted).
	RegenerateSetupToken(ctx context.Context) (plaintextToken string, err error)
	// SubmitToken exchanges the one-time setup token for a short-lived wizard session token.
	// Returns user_errors.SetupAlreadyCompleted if setup has already finished, or
	// user_errors.WizardSessionInvalid if token does not match.
	SubmitToken(ctx context.Context, token string) (wizardSessionToken string, err error)
	// SelectAuthMethods enables/disables password and telegram login for the instance. Returns
	// user_errors.AtLeastOneAuthMethodRequired if both are false.
	SelectAuthMethods(ctx context.Context, wizardSessionToken string, passwordEnabled, telegramEnabled bool) error
	SelectRegistrationMode(ctx context.Context, wizardSessionToken string, mode domain.RegistrationMode) error
	// CompleteSetup creates the first administrator account (via password or telegram, whichever
	// was enabled by SelectAuthMethods), marks setup complete, and returns a logged-in session
	// for the new admin. Exactly one of password/telegramIdToken must be set.
	CompleteSetup(
		ctx context.Context, wizardSessionToken string, email, password, telegramIdToken string,
	) (domain.Session, error)
}

// AdminSubscriptionsService lets an admin inspect and edit a user's plan assignment and
// overrides — orthogonal to SubscriptionService's Free/Paid enforcement-mode split, since
// managing a subscription isn't itself an enforcement decision.
type AdminSubscriptionsService interface {
	ListPlans(ctx context.Context) ([]domain.SubscriptionPlan, error)
	// GetUserSubscription returns both the raw per-user row (for editing) and the merged
	// plan+override view (for display).
	GetUserSubscription(ctx context.Context, userUuid uuid.UUID) (domain.Subscription, domain.EffectiveSubscription, error)
	UpdateUserSubscription(ctx context.Context, sub domain.Subscription) (domain.EffectiveSubscription, error)
}

type AuthService interface {
	// Register creates a new self-registered account. Returns user_errors.SetupNotCompleted if
	// instance setup hasn't finished yet, user_errors.SelfRegistrationDisabled if
	// RegistrationMode is admin_only, or user_errors.AuthMethodDisabled if password auth is off.
	Register(ctx context.Context, email, password string) (domain.User, error)
	// Login returns user_errors.SetupNotCompleted if instance setup hasn't finished yet, or
	// user_errors.AuthMethodDisabled if password auth is off.
	Login(ctx context.Context, email, password string) (domain.Session, error)
	Logout(ctx context.Context, token string) error
	ValidateToken(ctx context.Context, token string) (domain.User, error)
	// LoginViaTelegram returns user_errors.SetupNotCompleted if instance setup hasn't finished
	// yet, or user_errors.AuthMethodDisabled if telegram auth is off.
	LoginViaTelegram(ctx context.Context, idToken string) (domain.Session, error)
	GetMe(ctx context.Context, userUuid uuid.UUID) (domain.UserDetails, error)
	CheckIsAdmin(ctx context.Context, userUuid uuid.UUID) error
	// Refresh exchanges a still-valid refresh token for a new access/refresh token pair,
	// rotating the refresh token in the process. Returns user_errors.InvalidRefreshToken
	// if the token is unknown, already rotated, or expired.
	Refresh(ctx context.Context, refreshToken string) (domain.Session, error)
	// EnsureNoAuthUser creates (on first run) and caches the fixed local-dev user used to
	// bypass authentication entirely when the server is configured with NoAuthEnabled.
	// Every subsequent ValidateToken call returns this user regardless of the token presented.
	EnsureNoAuthUser(ctx context.Context) (domain.User, error)

	// RegisterAdmin creates a new administrator account, bypassing all setup/registration-mode/
	// auth-method checks — used only by SetupWizardService.CompleteSetup's password path.
	RegisterAdmin(ctx context.Context, email, password string) (domain.User, error)
	// CreateUserUnchecked creates a new non-admin account, bypassing all setup/registration-mode/
	// auth-method checks — used only by AdminUsersService.CreateUser, where the caller is already
	// authenticated and authorized as an administrator at the transport level.
	CreateUserUnchecked(ctx context.Context, email, password string) (domain.User, error)
	// LoginOrRegisterAdminViaTelegram is the wizard-only variant of LoginViaTelegram: no
	// setup/auth-method checks, and the resulting user (new or pre-existing) is promoted to
	// administrator — used only by SetupWizardService.CompleteSetup's telegram path.
	LoginOrRegisterAdminViaTelegram(ctx context.Context, idToken string) (domain.Session, error)
	// ChangePassword resets userUuid's password. Callers are responsible for authorizing the
	// caller before invoking this — used by AdminUsersService.ChangePassword.
	ChangePassword(ctx context.Context, userUuid uuid.UUID, newPassword string) error
}

type VaultService interface {
	CreateVault(ctx context.Context, name string) (domain.Vault, error)
	GetVault(ctx context.Context, vaultID uuid.UUID) (domain.Vault, error)
	ListVaults(ctx context.Context) ([]domain.Vault, error)
	DeleteVault(ctx context.Context, vaultID uuid.UUID) error
	AddMember(ctx context.Context, vaultID, targetUserUuid uuid.UUID, role artel_q.VaultRole) error
	RemoveMember(ctx context.Context, vaultID, targetUserUuid uuid.UUID) error
	ListMembers(ctx context.Context, vaultID uuid.UUID) ([]domain.VaultMemberInfo, error)
	CreateInviteLink(ctx context.Context, vaultID uuid.UUID, role artel_q.VaultRole) (domain.VaultInvite, error)
	ListInviteLinks(ctx context.Context, vaultID uuid.UUID) ([]domain.VaultInvite, error)
	RevokeInviteLink(ctx context.Context, inviteID uuid.UUID) error
	AcceptInvite(ctx context.Context, token string) error
	// LinkS3Bucket links vaultID to an S3 bucket. When s3InstanceID is nil, the instance is
	// auto-resolved for the caller (see vault.Service.resolveS3Instance) instead of requiring an
	// explicit id.
	LinkS3Bucket(ctx context.Context, vaultID uuid.UUID, s3InstanceID *uuid.UUID, bucketName string) error
	UnlinkS3Bucket(ctx context.Context, vaultID uuid.UUID) error
	SetUseCouchDBForBinaries(ctx context.Context, vaultID uuid.UUID, useCouchDB bool) error
	PublishVault(ctx context.Context, vaultID uuid.UUID, slug string) (domain.Vault, error)
	UnpublishVault(ctx context.Context, vaultID uuid.UUID) error

	// EnablePostgresDatabase provisions a Postgres database+role for vaultID on a pool instance
	// resolved via PostgresInstances.PickForUser (BYOK-owned preferred, admin pool fallback).
	// Returns user_errors.PostgresDatabaseAlreadyEnabled if vaultID already has one.
	EnablePostgresDatabase(ctx context.Context, vaultID uuid.UUID) (domain.VaultPostgresDatabase, error)
	// GetPostgresDatabase returns vaultID's Postgres database row, Valid: false if none enabled.
	GetPostgresDatabase(ctx context.Context, vaultID uuid.UUID) (sql.Null[domain.VaultPostgresDatabase], error)
	// DisablePostgresDatabase best-effort drops vaultID's provisioned database+role, then deletes
	// the row. No-op if none enabled.
	DisablePostgresDatabase(ctx context.Context, vaultID uuid.UUID) error
}

// WorkbenchService manages the per-(vault, user) Docker workbench container — every vault
// member gets their own workbench, tracked through domain.WorkbenchStatus's state machine; see
// domain.WorkbenchAuthMode for the two auth modes StartWorkbench implements.
//
// CreateWorkbench/GetWorkbench/StartWorkbench/StopWorkbench/DeleteWorkbench all resolve the
// calling user from ctx internally (via user_context.GetUserContext, same idiom as
// vault.Service.requireVaultMember) rather than taking an explicit userID parameter — they're
// only ever called from gRPC handlers behind the auth interceptor, which injects it. Never trust
// a client-supplied user id for these. ResolveTerminalTarget/ResolveTerminalShellTarget are the
// exceptions: their only callers (the chat-bridge and ttyd reverse-proxy handlers, respectively)
// sit outside that interceptor chain and authenticate the request themselves, so they pass userID
// explicitly. The four terminal-tab methods below (ListTerminalTabs/CreateTerminalTab/
// SelectTerminalTab/CloseTerminalTab) resolve the caller from ctx too, same as
// Create/Get/Start/Stop/Delete — they are not part of that exception.
type WorkbenchService interface {
	CreateWorkbench(ctx context.Context, vaultID uuid.UUID) (domain.Workbench, error)
	GetWorkbench(ctx context.Context, vaultID uuid.UUID) (domain.Workbench, error)
	StartWorkbench(ctx context.Context, vaultID uuid.UUID, authMode domain.WorkbenchAuthMode) (domain.Workbench, error)
	// StopWorkbench returns the paths where a genuine edit conflict was detected while syncing
	// the workbench's edited files back into the vault (see notes.Service.SyncFromWorkbench) —
	// the workbench itself is still marked stopped regardless.
	StopWorkbench(ctx context.Context, vaultID uuid.UUID) (conflicts []string, err error)
	DeleteWorkbench(ctx context.Context, vaultID uuid.UUID) error
	// DeleteWorkbenchesForVault tears down every member's workbench for vaultID — used by full
	// vault deletion, which must clean up every user's Docker resources, not just the caller's.
	DeleteWorkbenchesForVault(ctx context.Context, vaultID uuid.UUID) error
	// ResolveTerminalTarget returns the "http://<host>:<port>" base URL of userID's running
	// workbench's in-container chat bridge, for the reverse-proxy handler
	// (internal/transport/vaults_api/workbench_terminal.go) to forward to. Returns
	// user_errors.WorkbenchNotRunning when the workbench isn't running.
	ResolveTerminalTarget(ctx context.Context, vaultID, userID uuid.UUID) (string, error)
	// ResolveTerminalShellTarget returns the "http://<host>:<port>" base URL of userID's running
	// workbench's in-container ttyd server (the interactive tmux-tab terminal), for the
	// terminal-shell reverse-proxy handler
	// (internal/transport/vaults_api/workbench_terminal_shell.go) to forward to. Returns
	// user_errors.WorkbenchNotRunning when the workbench isn't running.
	ResolveTerminalShellTarget(ctx context.Context, vaultID, userID uuid.UUID) (string, error)
	// ListTerminalTabs lists the calling user's own running workbench's terminal tabs (tmux
	// windows) for vaultID, in tmux's own window order.
	ListTerminalTabs(ctx context.Context, vaultID uuid.UUID) ([]domain.TerminalTab, error)
	// CreateTerminalTab opens a new terminal tab (tmux window, running `claude`) in the calling
	// user's own running workbench for vaultID. There is no name parameter — tabs are always
	// auto-named by tmux, never user-supplied.
	CreateTerminalTab(ctx context.Context, vaultID uuid.UUID) (domain.TerminalTab, error)
	// SelectTerminalTab makes tabID the calling user's own running workbench's current tmux
	// window for vaultID.
	SelectTerminalTab(ctx context.Context, vaultID uuid.UUID, tabID string) error
	// CloseTerminalTab closes tabID in the calling user's own running workbench for vaultID.
	// Returns user_errors.WorkbenchCannotCloseLastTab if tabID is the workbench's only remaining
	// tab.
	CloseTerminalTab(ctx context.Context, vaultID uuid.UUID, tabID string) error
}

type CouchInstanceService interface {
	RegisterCouchInstance(ctx context.Context, url, username, password string) (string, error)
	GetCouchInstance(ctx context.Context, id string) (domain.CouchInstance, error)
	ListCouchInstances(ctx context.Context) ([]domain.CouchInstance, error)
	UpdateCouchInstance(ctx context.Context, id, url, username, password string) error
	DeleteCouchInstance(ctx context.Context, id string) error
	SetupCouchInstance(ctx context.Context, id string) error
	GetCouchInstanceStatus(ctx context.Context, id string) (couchdb.SetupStatus, error)
	HasCouchInstances(ctx context.Context) (bool, error)
}

// DockerHostService is admin CRUD over the pool of Docker daemons that back per-vault workbench
// containers. Unlike CouchInstanceService/S3InstanceService there's no setup/status concept and
// no single credential blob; instead there
// are three optional TLS/mTLS fields for the remote-daemon case (migrations/
// 062_docker_hosts_tls.sql).
type DockerHostService interface {
	RegisterDockerHost(ctx context.Context, url, caCert, clientCert, clientKey string) (string, error)
	GetDockerHost(ctx context.Context, id string) (domain.DockerHost, error)
	ListDockerHosts(ctx context.Context) ([]domain.DockerHost, error)
	// UpdateDockerHost patches url unconditionally; caCert/clientCert/clientKey are three-way
	// patch pointers: nil leaves the stored cert untouched, a pointer to "" clears it, a pointer
	// to a non-empty PEM re-encrypts and stores it.
	UpdateDockerHost(ctx context.Context, id, url string, caCert, clientCert, clientKey *string) error
	DeleteDockerHost(ctx context.Context, id string) error
	HasDockerHosts(ctx context.Context) (bool, error)
}

type S3InstanceService interface {
	RegisterS3Instance(
		ctx context.Context, endpoint, region, accessKey, secretKey string, useSSL, pathStyle bool,
	) (string, error)
	GetS3Instance(ctx context.Context, id string) (domain.S3Instance, error)
	ListS3Instances(ctx context.Context) ([]domain.S3Instance, error)
	UpdateS3Instance(
		ctx context.Context, id, endpoint, region, accessKey, secretKey string, useSSL, pathStyle bool,
	) error
	DeleteS3Instance(ctx context.Context, id string) error
	TestS3Instance(ctx context.Context, id string) error
	HasS3Instances(ctx context.Context) (bool, error)
}

type AdminCouchService interface {
	ListUsers(ctx context.Context, instanceId string) ([]couchdb.UserListEntry, error)
	DeleteUser(ctx context.Context, instanceId, username string) error
	ChangeUserPassword(ctx context.Context, instanceId, username, newPassword string) error
	GrantDatabaseAccess(ctx context.Context, instanceId, dbName, username string) error
	RevokeDatabaseAccess(ctx context.Context, instanceId, dbName, username string) error
	ListDatabases(ctx context.Context, instanceId string) ([]string, error)
	GetUserDatabaseAccess(ctx context.Context, instanceId, username string) ([]string, error)
}

type McpService interface {
	// CreateKey generates a new bearer token, stores it hashed, returns the raw token once.
	CreateKey(ctx context.Context, vaultID uuid.UUID, name string) (rawToken string, key domain.McpKey, err error)
	ListKeys(ctx context.Context, vaultID uuid.UUID) ([]domain.McpKey, error)
	ListUserKeys(ctx context.Context) ([]domain.McpKey, error)
	RevokeKey(ctx context.Context, keyID uuid.UUID) error
	SetKeyAccess(ctx context.Context, keyID uuid.UUID, vaultID uuid.UUID) error
	// ListConnectors returns the MCP connectors linked to keyID.
	ListConnectors(ctx context.Context, keyID uuid.UUID) ([]domain.McpConnector, error)
	// AddConnector links an MCP definition to keyID via an existing external connection.
	AddConnector(
		ctx context.Context, keyID uuid.UUID, mcpName string, externalConnectionID uuid.UUID,
	) (domain.McpConnector, error)
	// RemoveConnector unlinks mcpName from keyID.
	RemoveConnector(ctx context.Context, keyID uuid.UUID, mcpName string) error
	// ListMomCandidates returns all available MCP definitions paired with the
	// caller's external connections that can satisfy each one.
	ListMomCandidates(ctx context.Context) ([]domain.MomCandidate, error)
	// ListCommunityConnectors returns admin-authored community MoMs (IsCommunity == true),
	// paired with the caller's external connections the same way ListMomCandidates does, plus
	// ViewerIsOwner so callers can tell their own connectors apart from other admins'.
	ListCommunityConnectors(ctx context.Context) ([]domain.MomCandidate, error)
	// DeleteCommunityConnector deletes a community connector by name — the caller must be its
	// owning admin. Returns user_errors.NotFound both when the name doesn't exist and when the
	// caller doesn't own it, so a non-owner can't learn who owns a given name.
	DeleteCommunityConnector(ctx context.Context, name string) error
	// ResolveKey validates the raw bearer token and returns vault+couch context.
	ResolveKey(ctx context.Context, rawToken string) (domain.McpKeyContext, error)
	// ListTools returns the built-in tool definitions (vault tools + connections).
	ListTools(ctx context.Context) ([]domain.McpToolDef, error)
	// IsBuiltinTool reports whether name is a built-in tool (vault tools + connections).
	IsBuiltinTool(name string) bool
	// ExecuteTool executes a built-in tool by name.
	ExecuteTool(
		ctx context.Context, keyCtx domain.McpKeyContext, toolName string, params map[string]interface{},
	) (domain.ToolExecResult, error)
	// ExecuteBuiltinToolForUser executes a built-in tool as userUuid rather than through an
	// MCP key context — used by the tract engine, which has no key.
	ExecuteBuiltinToolForUser(
		ctx context.Context, userUuid uuid.UUID, toolName string, params map[string]interface{},
	) (string, error)
	// SetTractService wires the tract service dependency for the tract-authoring builtin tools
	// (list_tract_actions, create_tract, ...). Called once from internal/app/custom.go after
	// TractService is constructed — Tract composes Mcp's ToolExecutor, so Mcp must exist first;
	// this setter breaks that construction-order cycle instead of Mcp importing Tract's package.
	// baseCtx is the server-lifecycle context (App.Ctx) — run_tract spawns TractService.StartRun
	// against it rather than the per-request ctx, which net/http cancels once the MCP handler's
	// response is written (mirrors tracts_api.TractsImpl.baseCtx).
	SetTractService(baseCtx context.Context, ts TractService)
	// ListHotPlugSkillTools returns one dynamic skill_<slug> tool per hot-plug skill visible in
	// keyCtx's vault (the always-hot-plug system skill-creator skill included). Deliberately kept
	// out of ListTools/IsBuiltinTool — that pair is also relied on by Tract's toolExecutorAdapter
	// for step-name validation via a vault-agnostic, cacheable catalog, which a per-vault dynamic
	// tool set must not leak into. internal/transport/mcp_api's handleToolsList merges this in
	// separately.
	ListHotPlugSkillTools(ctx context.Context, keyCtx domain.McpKeyContext) ([]domain.McpToolDef, error)
	// ExecuteSkillTool runs a dynamic skill_<slug> tool (slug has already had the "skill_" prefix
	// stripped by the caller), returning the skill's body text verbatim.
	ExecuteSkillTool(ctx context.Context, keyCtx domain.McpKeyContext, slug string) (string, error)
}

// SkillsService manages skills stored as CouchDB notes inside a vault's reserved .skills/
// folder, plus the always-present, non-vault-backed system "skill-creator" skill. See
// internal/service/v1/skills.Service for the implementation.
type SkillsService interface {
	ListSkills(ctx context.Context, vaultUuid uuid.UUID) ([]domain.Skill, error)
	GetSkillBody(ctx context.Context, vaultUuid uuid.UUID, slug string) (domain.Skill, error)
	CreateSkill(
		ctx context.Context, vaultUuid uuid.UUID, name string, description string,
		storageMode domain.SkillStorageMode, body string, hotPlug bool,
	) (domain.Skill, error)
	UpdateSkill(
		ctx context.Context, vaultUuid uuid.UUID, slug string, name string, description string,
		storageMode domain.SkillStorageMode, body string,
	) (domain.Skill, error)
	SetSkillHotPlug(ctx context.Context, vaultUuid uuid.UUID, slug string, hotPlug bool) (domain.Skill, error)
	DeleteSkill(ctx context.Context, vaultUuid uuid.UUID, slug string) error
}

// SubscriptionService regulates access to gated functionality (feature flags previously handled
// by UserPermissions) and storage quotas. Two implementations exist behind this interface,
// selected by config.EnvironmentConfig.SubscriptionsEnabled: a no-op "free" implementation
// (always allow, unlimited — used when subscriptions are disabled) and a "paid" implementation
// enforcing each user's effective plan+overrides.
type SubscriptionService interface {
	CheckActive(ctx context.Context, userUuid uuid.UUID) error
	// HasFeature reports whether userUuid's effective subscription grants feature.
	HasFeature(ctx context.Context, userUuid uuid.UUID, feature domain.SubscriptionFeature) (bool, error)
	// CheckFeature is HasFeature wrapped in user_errors.FeatureNotEnabled for callers that just
	// want to gate a call.
	CheckFeature(ctx context.Context, userUuid uuid.UUID, feature domain.SubscriptionFeature) error
	// GetEffective returns userUuid's merged plan+override view.
	GetEffective(ctx context.Context, userUuid uuid.UUID) (domain.EffectiveSubscription, error)
	// GetUsage measures userUuid's current storage footprint on demand (CouchDB + S3, summed
	// across every vault they belong to).
	GetUsage(ctx context.Context, userUuid uuid.UUID) (domain.StorageUsage, error)
	// CheckStorageQuota compares GetUsage against GetEffective's quotas, returning
	// user_errors.CouchStorageQuotaExceeded / S3StorageQuotaExceeded if either is already over.
	CheckStorageQuota(ctx context.Context, userUuid uuid.UUID) error
	// CheckSkillLimit compares the given vault's current skill counts (measured live from its
	// CouchDB .skills/ folders, same on-demand-measurement approach as CheckStorageQuota)
	// against the vault owner's effective plan caps, returning
	// user_errors.SkillLimitExceeded / HotPlugSkillLimitExceeded if already at or over. When
	// wantHotPlug is true both the total and hot-plug caps are checked; otherwise only the
	// total cap is checked.
	CheckSkillLimit(ctx context.Context, vaultUuid uuid.UUID, wantHotPlug bool) error
}

type ListPromptsParams struct {
	Ids      []string
	Page     uint32
	PageSize uint32
}

type PromptService interface {
	ListPrompts(ctx context.Context, params ListPromptsParams) ([]domain.Prompt, int64, error)
}

type TaskTrackerService interface {
	AddTracker(ctx context.Context, apiKey, apiToken string) (domain.TaskTracker, []domain.TrelloBoard, error)
	ListTrackers(ctx context.Context) ([]domain.TaskTracker, error)
	DeleteTracker(ctx context.Context, trackerUuid uuid.UUID) error
	ListTrelloBoards(ctx context.Context, trackerUuid uuid.UUID) ([]domain.TrelloBoard, error)
}

type NotesService interface {
	ListFolders(ctx context.Context, vaultID uuid.UUID) ([]string, error)
	ListNotes(ctx context.Context, vaultID uuid.UUID) ([]couchdb.NoteEntry, error)
	GetNote(ctx context.Context, vaultID uuid.UUID, path string) (couchdb.NoteDoc, error)
	ListTags(ctx context.Context, vaultID uuid.UUID) ([]string, error)
	SaveNote(ctx context.Context, vaultID uuid.UUID, path, content string) error
	MoveNote(ctx context.Context, vaultID uuid.UUID, oldPath, newPath string) error
	ExportFolder(ctx context.Context, vaultID uuid.UUID, folderPath string) ([]byte, error)
	CheckImportConflicts(ctx context.Context, vaultID uuid.UUID, destPath string, zipData []byte) ([]string, error)
	CommitImport(
		ctx context.Context, vaultID uuid.UUID, destPath string, zipData []byte, resolutions []domain.ImportResolution,
	) (imported int, skipped int, err error)
	DeleteFolder(ctx context.Context, vaultID uuid.UUID, folderPath string) (deletedCount int, failedPaths []string, err error)
	MoveFolder(ctx context.Context, vaultID uuid.UUID, oldPath, newPath string) (movedCount int, err error)
	// SyncFromWorkbench pushes markdown files edited inside a stopped workbench container's
	// workspace back into vaultID's notes — see internal/service/v1/notes.Service.
	// SyncFromWorkbench's doc comment for the full snapshot-sync/conflict-detection semantics.
	SyncFromWorkbench(
		ctx context.Context, vaultID uuid.UUID, files map[string][]byte, snapshot map[string]int64,
	) (conflicts []string, err error)
}

// PublicDocsService is the anonymous, auth-exempt read surface backing PublicDocsAPI (see
// docs/architecture.md). Every method resolves the target vault by slug and only ever returns
// data for a published vault (domain.Vault.IsPublic) — it must never require a user_context.
type PublicDocsService interface {
	GetVaultBySlug(ctx context.Context, slug string) (domain.Vault, error)
	ListFolders(ctx context.Context, slug string) ([]string, error)
	ListNotes(ctx context.Context, slug string) ([]couchdb.NoteEntry, error)
	GetNote(ctx context.Context, slug, path string) (couchdb.NoteDoc, error)
	ListTags(ctx context.Context, slug string) ([]string, error)
	// GetDefaultVault resolves the admin-configured default `/docs` vault (SystemSettings.
	// DefaultDocsVaultID), returning user_errors.NotFound both when unset and when the
	// configured vault is no longer published.
	GetDefaultVault(ctx context.Context) (domain.Vault, error)
}

type MomService interface {
	ListToolsForKey(ctx context.Context, keyId uuid.UUID) ([]domain.McpToolDef, error)
	ExecuteToolForKey(
		ctx context.Context, keyId uuid.UUID, toolName string, params map[string]interface{},
	) (string, error)
	ExecuteToolForUserConnection(
		ctx context.Context, exConnUuid uuid.UUID, mcpName string, toolName string, params map[string]interface{},
	) (string, error)
	// ExecuteToolForConnection executes a tool against a specific external connection without
	// an ownership check — used by the tract engine, which verifies ownership itself since it
	// has no user_context for webhook-triggered runs.
	ExecuteToolForConnection(
		ctx context.Context, exConnUuid uuid.UUID, mcpName string, toolName string, params map[string]interface{},
	) (string, error)
	// ExecuteToolWithSecrets executes a MoM http tool's action directly against a caller-supplied
	// secrets map, without requiring a persisted external_connections row — used to validate
	// externally-supplied credentials against the real provider before anything is saved.
	ExecuteToolWithSecrets(
		ctx context.Context, mcpName, toolName string, secrets, params map[string]interface{},
	) (string, error)
}

// TractService owns tract/trigger CRUD, trigger↔tract links, and run lifecycle (StartRun,
// the startup sweep). CreateTract/UpdateTract return non-fatal warnings alongside the
// persisted tract — e.g. a trigger.* ref not found in any linked trigger's payload schema.
type TractService interface {
	CreateTract(
		ctx context.Context, name string, description string, def domain.TractDefinition,
	) (domain.Tract, []string, error)
	GetTract(ctx context.Context, id uuid.UUID) (domain.Tract, error)
	ListTracts(ctx context.Context) ([]domain.Tract, error)
	UpdateTract(
		ctx context.Context, id uuid.UUID, name string, description string, def domain.TractDefinition,
	) (domain.Tract, []string, error)
	SetTractEnabled(ctx context.Context, id uuid.UUID, enabled bool) error
	DeleteTract(ctx context.Context, id uuid.UUID) error

	// PublishTemplate snapshots tractUuid's current definition as a new public template.
	PublishTemplate(ctx context.Context, tractUuid uuid.UUID, category string) (domain.TractTemplate, error)
	// UnpublishTemplate removes a template — publisher-only.
	UnpublishTemplate(ctx context.Context, templateUuid uuid.UUID) error
	// ListTemplates lists published templates; mineOnly restricts to the caller's own.
	ListTemplates(ctx context.Context, category string, mineOnly bool) ([]domain.TractTemplate, error)
	// GetTemplate returns one template — any authenticated user may read any template.
	GetTemplate(ctx context.Context, templateUuid uuid.UUID) (domain.TractTemplate, error)
	// InstantiateTemplate copies templateUuid into a brand-new tract owned by the caller.
	// connections maps MoM name -> connection uuid for every non-builtin action step the
	// template uses.
	InstantiateTemplate(
		ctx context.Context, templateUuid uuid.UUID, name, description string, connections map[string]uuid.UUID,
	) (domain.Tract, []string, error)

	// CreateTrigger returns the raw webhook token once — only its hash is persisted.
	CreateTrigger(
		ctx context.Context, name string, kind string, source string,
		config json.RawMessage, payloadSchema domain.ToolSchema,
	) (domain.Trigger, string, error)
	GetTrigger(ctx context.Context, id uuid.UUID) (domain.Trigger, error)
	ListTriggers(ctx context.Context) ([]domain.Trigger, error)
	SetTriggerEnabled(ctx context.Context, id uuid.UUID, enabled bool) error
	DeleteTrigger(ctx context.Context, id uuid.UUID) error
	// RotateTriggerToken invalidates triggerUuid's current webhook URL and mints a new one,
	// returning the raw token once (only its hash is persisted) — same one-time-reveal contract as
	// CreateTrigger.
	RotateTriggerToken(ctx context.Context, id uuid.UUID) (domain.Trigger, string, error)
	LinkTrigger(ctx context.Context, triggerUuid uuid.UUID, tractUuid uuid.UUID, filters []domain.TractCondition) error
	UnlinkTrigger(ctx context.Context, triggerUuid uuid.UUID, tractUuid uuid.UUID) error
	ListLinksByTract(ctx context.Context, tractUuid uuid.UUID) ([]repository.TractTriggerLink, error)

	// StartRun persists the run then walks the definition; the caller decides whether to run
	// it synchronously or as `go TractService.StartRun(...)` against a server-lifecycle ctx.
	StartRun(
		ctx context.Context, tract domain.Tract, payload json.RawMessage, startedBy string, triggerUuid uuid.UUID,
	) (domain.TractRun, error)
	// SweepStaleRuns marks stale 'running' runs/steps 'failed' — call once at app init.
	SweepStaleRuns(ctx context.Context, threshold time.Time) error
	// ListRuns returns tractUuid's most recent runs (most recent first), capped at limit.
	ListRuns(ctx context.Context, tractUuid uuid.UUID, limit int32) ([]domain.TractRun, error)
	// GetRun returns one run plus its ordered step rows.
	GetRun(ctx context.Context, id uuid.UUID) (domain.TractRun, []domain.TractRunStep, error)

	// ListTractTools returns the action picker's tool catalog: builtins (Mcp == "artel") plus
	// every mcp_tools row across every MoM.
	ListTractTools(ctx context.Context) ([]domain.McpToolRef, error)
	// ListTriggerSources returns the webhook preset catalog (gitlab_push, generic, ...).
	ListTriggerSources(ctx context.Context) ([]domain.TriggerPreset, error)
}

type ExternalConnectionService interface {
	InitiateGoogleOAuth(ctx context.Context, origin string) (authURL string, err error)
	HandleGoogleOAuthCallback(ctx context.Context, code string, state string) (domain.ExternalConnectionMeta, error)
	DisconnectProvider(ctx context.Context, provider string) error
	DisconnectConnection(ctx context.Context, id string) error
	ListConnections(ctx context.Context) ([]domain.ExternalConnectionMeta, error)
	GetGoogleClient(ctx context.Context) (*googleapi.Client, error)
	GetPickerToken(ctx context.Context) (string, error)
	AddSpreadsheet(ctx context.Context, spreadsheetId string, name string) (domain.McpSpreadsheet, error)
	ListSpreadsheets(ctx context.Context) ([]domain.McpSpreadsheet, error)
	RemoveSpreadsheet(ctx context.Context, spreadsheetId string) error
	AddEmailConnection(
		ctx context.Context, email, imapHost string, imapPort int, smtpHost string, smtpPort int, password string,
	) (domain.ExternalConnectionMeta, error)
	CheckEmailConnection(
		ctx context.Context, email, imapHost string, imapPort int, smtpHost string, smtpPort int, password string,
	) error
	ListMailServerSuggestions(ctx context.Context, domain string) ([]domain.MailServerSuggestion, error)
	AddGitlabConnection(
		ctx context.Context, personalAccessToken, instanceUrl string,
	) (domain.ExternalConnectionMeta, error)
	CheckGitlabConnection(ctx context.Context, personalAccessToken, instanceUrl string) (username string, err error)
	GenerateGitlabWebhookSecret(ctx context.Context) (domain.ExternalConnectionMeta, string, error)
	AddTrelloConnection(ctx context.Context, apiKey, apiToken string) (domain.ExternalConnectionMeta, error)
	CheckTrelloConnection(ctx context.Context, apiKey, apiToken string) (fullName string, err error)
	AddTelegramConnection(ctx context.Context, botToken string) (domain.ExternalConnectionMeta, error)
	CheckTelegramConnection(ctx context.Context, botToken string) (botUsername string, err error)
	AddAnthropicConnection(
		ctx context.Context, apiKey, baseUrl, defaultModel string,
	) (domain.ExternalConnectionMeta, error)
	CheckAnthropicConnection(
		ctx context.Context, apiKey, baseUrl, defaultModel string,
	) (models []anthropicClient.ModelInfo, recommendedDefaultModel string, err error)
	// GetAnthropicApiKey returns userUuid's decrypted BYOK Anthropic api key — used by
	// WorkbenchService to inject ANTHROPIC_API_KEY when starting a workbench in api_key auth
	// mode. Returns user_errors.LlmKeyRequired if userUuid has no anthropic connection.
	GetAnthropicApiKey(ctx context.Context, userUuid uuid.UUID) (string, error)
	// AddOpenAIConnection and CheckOpenAIConnection are shared across every provider that speaks
	// the OpenAI Chat Completions protocol (currently openai and openrouter); provider selects
	// which one, via a domain.Provider* constant — see
	// external_connections_api.OpenAICompatibleProviderFromProto for how the transport layer
	// derives it from the request's proto ExternalProvider field.
	AddOpenAIConnection(
		ctx context.Context, apiKey, baseUrl, defaultModel, provider string,
	) (domain.ExternalConnectionMeta, error)
	CheckOpenAIConnection(
		ctx context.Context, apiKey, baseUrl, defaultModel, provider string,
	) (models []openaiClient.ModelInfo, recommendedDefaultModel string, err error)
	AddGenericConnection(
		ctx context.Context, provider string, credentials map[string]string,
	) (domain.ExternalConnectionMeta, error)
	AddS3Connection(
		ctx context.Context, endpoint, region, accessKey, secretKey string, useSSL, pathStyle bool,
	) (domain.ExternalConnectionMeta, error)
	CheckS3Connection(
		ctx context.Context, endpoint, region, accessKey, secretKey string, useSSL, pathStyle bool,
	) error
	AddCouchDBConnection(
		ctx context.Context, url, username, password string,
	) (domain.ExternalConnectionMeta, error)
	CheckCouchDBConnection(ctx context.Context, url, username, password string) error
	AddPostgresConnection(
		ctx context.Context, host string, port int, database, username, password, sslMode string,
	) (domain.ExternalConnectionMeta, error)
	CheckPostgresConnection(
		ctx context.Context, host string, port int, database, username, password, sslMode string,
	) error
}

// PostgresInstanceService is admin CRUD over the shared pool of Postgres servers (admin pool or
// BYOK) that per-vault databases are provisioned on — mirrors S3InstanceService/CouchInstanceService.
type PostgresInstanceService interface {
	RegisterPostgresInstance(
		ctx context.Context, host string, port int, adminDatabase, username, password, sslMode string,
	) (string, error)
	GetPostgresInstance(ctx context.Context, id string) (domain.PostgresInstance, error)
	ListPostgresInstances(ctx context.Context) ([]domain.PostgresInstance, error)
	UpdatePostgresInstance(
		ctx context.Context, id, host string, port int, adminDatabase, username, password, sslMode string,
	) error
	DeletePostgresInstance(ctx context.Context, id string) error
	// TestPostgresInstance opens a connection to id's admin database and pings it, without
	// mutating anything — mirrors S3InstanceService.TestS3Instance.
	TestPostgresInstance(ctx context.Context, id string) error
	HasPostgresInstances(ctx context.Context) (bool, error)
}
