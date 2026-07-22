package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	anthropicClient "github.com/ruf-dev/artel/internal/clients/anthropic"
	"github.com/ruf-dev/artel/internal/clients/couchdb"
	"github.com/ruf-dev/artel/internal/clients/googleapi"
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
}

type AdminUsersService interface {
	ListUsers(ctx context.Context, req domain.ListUsersReq) ([]domain.User, int64, error)
	GetUser(ctx context.Context, userUuid uuid.UUID) (domain.UserDetails, error)
	GetUserSessions(ctx context.Context, userUuid uuid.UUID) ([]domain.Session, error)
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
	Register(ctx context.Context, email, password string) (domain.User, error)
	Login(ctx context.Context, email, password string) (domain.Session, error)
	Logout(ctx context.Context, token string) error
	ValidateToken(ctx context.Context, token string) (domain.User, error)
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
	LinkS3Bucket(ctx context.Context, vaultID, s3InstanceID uuid.UUID, bucketName string) error
	UnlinkS3Bucket(ctx context.Context, vaultID uuid.UUID) error
	SetUseCouchDBForBinaries(ctx context.Context, vaultID uuid.UUID, useCouchDB bool) error
}

type CouchInstanceService interface {
	RegisterCouchInstance(ctx context.Context, url, username, password string) (string, error)
	GetCouchInstance(ctx context.Context, id string) (domain.CouchInstance, error)
	ListCouchInstances(ctx context.Context) ([]domain.CouchInstance, error)
	UpdateCouchInstance(ctx context.Context, id, url, username, password string) error
	DeleteCouchInstance(ctx context.Context, id string) error
	SetupCouchInstance(ctx context.Context, id string) error
	GetCouchInstanceStatus(ctx context.Context, id string) (couchdb.SetupStatus, error)
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
	AddAnthropicConnection(
		ctx context.Context, apiKey, baseUrl, defaultModel string,
	) (domain.ExternalConnectionMeta, error)
	CheckAnthropicConnection(
		ctx context.Context, apiKey, baseUrl, defaultModel string,
	) (models []anthropicClient.ModelInfo, recommendedDefaultModel string, err error)
}
