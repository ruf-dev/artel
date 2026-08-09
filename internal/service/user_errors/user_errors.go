package user_errors

import (
	"net/http"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
	"google.golang.org/grpc/codes"
)

var (
	Unauthenticated = rerrors.New("unauthenticated", codes.Unauthenticated)
	NotFound        = rerrors.New("not found", codes.NotFound)
	AlreadyExists   = rerrors.New("already exists", codes.AlreadyExists)

	InvalidCouchDbDatabaseName   = rerrors.New("invalid database name", codes.InvalidArgument)
	CouchDbDatabaseAlreadyExists = rerrors.New("database already exists", codes.FailedPrecondition)
	UserAlreadyExistInCouchDb    = rerrors.New("user already exists in couch db", codes.AlreadyExists)
	CouchDbAccountLocked         = rerrors.New(
		"couch admin account is locked after authentication failures; update instance credentials",
		codes.FailedPrecondition,
	)

	McpKeyRevoked = rerrors.New(
		"mcp key revoked",
		codes.FailedPrecondition,
		rerrors.WithHttpStatus(http.StatusBadRequest),
	)
	McpInvalidToken = rerrors.New(
		"invalid mcp token",
		codes.FailedPrecondition,
		rerrors.WithHttpStatus(http.StatusBadRequest),
	)

	// auth.
	SessionExpired         = rerrors.New("session expired", codes.Unauthenticated)
	InvalidTelegramToken   = rerrors.New("invalid telegram token", codes.Unauthenticated)
	UnsupportedLoginMethod = rerrors.New("unsupported login method", codes.InvalidArgument)
	InvalidRefreshToken    = rerrors.New("invalid refresh token", codes.Unauthenticated)

	// auth: first-run setup wizard / instance auth policy.
	SetupNotCompleted = rerrors.New(
		"instance setup has not been completed yet",
		codes.FailedPrecondition,
		rerrors.WithHttpStatus(http.StatusPreconditionFailed),
	)
	SelfRegistrationDisabled = rerrors.New(
		"self-registration is disabled for this instance; ask an administrator to create your account",
		codes.PermissionDenied,
		rerrors.WithHttpStatus(http.StatusForbidden),
	)
	AuthMethodDisabled = rerrors.New(
		"this authentication method is disabled for this instance",
		codes.FailedPrecondition,
		rerrors.WithHttpStatus(http.StatusPreconditionFailed),
	)
	WizardSessionInvalid = rerrors.New(
		"setup wizard session is invalid or has expired; submit the setup token again",
		codes.Unauthenticated,
		rerrors.WithHttpStatus(http.StatusUnauthorized),
	)
	SetupAlreadyCompleted = rerrors.New(
		"instance setup has already been completed",
		codes.FailedPrecondition,
		rerrors.WithHttpStatus(http.StatusConflict),
	)
	AtLeastOneAuthMethodRequired = rerrors.New(
		"at least one authentication method must be enabled",
		codes.InvalidArgument,
		rerrors.WithHttpStatus(http.StatusBadRequest),
	)
	AuthMethodNotEnabled = rerrors.New(
		"the credential provided does not match any enabled authentication method",
		codes.InvalidArgument,
		rerrors.WithHttpStatus(http.StatusBadRequest),
	)

	// subscription.

	NoActiveSubscription = rerrors.New("no active subscription",
		codes.FailedPrecondition,
		rerrors.WithHttpStatus(http.StatusForbidden),
		rerrors.WithPreconditionFailure("SUB", "subscription", pb.UserErrors_NoSubscription.String()),
	)

	// admin.
	NotAdmin = rerrors.New("not an administrator", codes.PermissionDenied)

	// feature flags.
	NotesNotEnabled   = rerrors.New("notes are not enabled for your account", codes.PermissionDenied)
	FeatureNotEnabled = rerrors.New(
		"this feature is not enabled for your account",
		codes.PermissionDenied,
		rerrors.WithHttpStatus(http.StatusForbidden),
	)

	// subscription: storage quotas.
	CouchStorageQuotaExceeded = rerrors.New(
		"couchdb storage quota exceeded for your subscription",
		codes.ResourceExhausted,
		rerrors.WithHttpStatus(http.StatusForbidden),
	)
	S3StorageQuotaExceeded = rerrors.New(
		"s3 storage quota exceeded for your subscription",
		codes.ResourceExhausted,
		rerrors.WithHttpStatus(http.StatusForbidden),
	)
	FileTooLarge = rerrors.New(
		"file exceeds the maximum allowed size",
		codes.InvalidArgument,
		rerrors.WithHttpStatus(http.StatusRequestEntityTooLarge),
	)

	// vault.
	NotVaultOwner      = rerrors.New("only vault owner can perform this action", codes.PermissionDenied)
	InviteLinkRevoked  = rerrors.New("invite link has been revoked", codes.FailedPrecondition)
	InvalidInviteToken = rerrors.New("invalid invite token", codes.NotFound)
	VaultSlugInvalid   = rerrors.New(
		"slug must be lowercase kebab-case, 3-64 characters",
		codes.InvalidArgument,
	)

	// middleware.
	NoMetadataInContext = rerrors.New("error getting metadata from context", codes.FailedPrecondition)
	NoAuthHeader        = rerrors.New("error getting auth header", codes.Unauthenticated)
	DebugNotSupported   = rerrors.New("debug not supported", codes.Unimplemented)
	CSRFTokenInvalid    = rerrors.New("csrf token missing or invalid", codes.PermissionDenied)

	// mcp tool argument validation.
	McpPathRequired      = rerrors.New("path is required and must be a string", codes.InvalidArgument)
	McpContentRequired   = rerrors.New("content is required and must be a string", codes.InvalidArgument)
	McpOldPathRequired   = rerrors.New("old_path is required and must be a string", codes.InvalidArgument)
	McpNewPathRequired   = rerrors.New("new_path is required and must be a string", codes.InvalidArgument)
	McpAccountIdRequired = rerrors.New("account_id is required and must be a string", codes.InvalidArgument)
	McpIdRequired        = rerrors.New("id is required and must be a string", codes.InvalidArgument)
	McpToRequired        = rerrors.New("to is required and must be a string", codes.InvalidArgument)
	McpSubjectRequired   = rerrors.New("subject is required and must be a string", codes.InvalidArgument)
	McpBodyRequired      = rerrors.New("body is required and must be a string", codes.InvalidArgument)

	// mcp executor.
	McpEmailActionMissing = rerrors.New(
		"email executor: action must have imap or smtp discriminator",
		codes.InvalidArgument,
	)
	McpUnknownImapOperation = rerrors.New("email executor: unknown imap operation", codes.InvalidArgument)
	McpUnknownSmtpOperation = rerrors.New("email executor: unknown smtp operation", codes.InvalidArgument)
	McpToolNotFound         = rerrors.New("tool not found in any connected mcp", codes.NotFound)
	McpConnectionNotOwned   = rerrors.New("external connection not found", codes.NotFound)
	McpConnectorNotFound    = rerrors.New("no mcp connector configured for this key", codes.FailedPrecondition)
	McpActionMissing        = rerrors.New(
		"tool action has no imap, smtp, or http discriminator set",
		codes.InvalidArgument,
	)
	McpSecretFieldMissing = rerrors.New(
		"http executor: referenced __secrets field not found in connected credentials",
		codes.FailedPrecondition,
	)
	McpCredentialsProviderMismatch = rerrors.New(
		"http executor: linked external connection provider does not match tool's declared credentials provider",
		codes.FailedPrecondition,
	)
	McpHttpRequestFailed = rerrors.New("http executor: upstream request failed", codes.Unavailable)

	// mcp: community connectors.
	McpCommunityConnectorNameRequired = rerrors.New(
		"name is required and must be a non-empty string",
		codes.InvalidArgument,
	)
	McpCommunityConnectorDescriptionRequired = rerrors.New(
		"description is required and must be a non-empty string",
		codes.InvalidArgument,
	)
	McpCommunityConnectorToolsRequired = rerrors.New(
		"tools is required and must be a non-empty array",
		codes.InvalidArgument,
	)
	McpCommunityConnectorNameTaken = rerrors.New(
		"this name belongs to a system connector or another admin's community connector",
		codes.AlreadyExists,
	)

	// livesync file-type constraints.
	UseReadNoteForTextFiles       = rerrors.New("use read_note for text files", codes.FailedPrecondition)
	UseDeleteNoteForTextFiles     = rerrors.New("use delete_note for text files", codes.FailedPrecondition)
	UseMoveNoteForTextFiles       = rerrors.New("use move_note for text files", codes.FailedPrecondition)
	ChunkedBinaryMoveNotSupported = rerrors.New(
		"move of chunked binary files is not supported; use Obsidian to move large files",
		codes.FailedPrecondition,
	)

	// notes: folder move.
	InvalidFolderMove = rerrors.New(
		"cannot move a folder into itself or one of its own subfolders",
		codes.InvalidArgument,
	)

	// imap.
	InvalidEmailId       = rerrors.New("invalid email id", codes.InvalidArgument)
	EmailMessageNotFound = rerrors.New("message not found", codes.NotFound)
	EmailCursorConflict  = rerrors.New("after_uid and before_uid cannot both be set", codes.InvalidArgument)

	// external connections.
	GoogleNotConnected = rerrors.New(
		"google account not connected",
		codes.FailedPrecondition,
		rerrors.WithHttpStatus(http.StatusPreconditionFailed),
	)
	InvalidOAuthState = rerrors.New(
		"invalid or expired oauth state",
		codes.InvalidArgument,
		rerrors.WithHttpStatus(http.StatusBadRequest),
	)
	GenericProviderRequired = rerrors.New(
		"provider is required",
		codes.InvalidArgument,
		rerrors.WithHttpStatus(http.StatusBadRequest),
	)

	// email connection.
	EmailImapValidationFailed = rerrors.New(
		"could not connect to the imap server with the provided settings",
		codes.FailedPrecondition,
		rerrors.WithHttpStatus(http.StatusPreconditionFailed),
	)
	EmailSmtpValidationFailed = rerrors.New(
		"could not connect to the smtp server with the provided settings",
		codes.FailedPrecondition,
		rerrors.WithHttpStatus(http.StatusPreconditionFailed),
	)
	EmailPasswordRequired = rerrors.New(
		"password is required to connect this email account",
		codes.InvalidArgument,
		rerrors.WithHttpStatus(http.StatusBadRequest),
	)

	NoCouchDbInstance = rerrors.New("no storage instance",
		codes.FailedPrecondition,
		rerrors.WithHttpStatus(http.StatusFailedDependency),
		rerrors.WithPreconditionFailure("COUCHDB", "instance", pb.UserErrors_NoCouchDbInstance.String()))

	// s3 / binary storage.
	NoS3BucketLinked = rerrors.New(
		"vault has no linked S3 bucket; this operation is unavailable for non-markdown files",
		codes.FailedPrecondition,
		rerrors.WithPreconditionFailure("S3", "bucket", pb.UserErrors_NoS3BucketLinked.String()),
	)

	McpCrossBackendMoveNotSupported = rerrors.New(
		"cannot move a file across the markdown/binary storage boundary in one call; "+
			"write to the new path then delete the old path instead",
		codes.InvalidArgument,
	)

	// gitlab webhook.
	GitlabWebhookSecretMismatch = rerrors.New(
		"gitlab webhook: token does not match configured secret",
		codes.PermissionDenied,
		rerrors.WithHttpStatus(http.StatusUnauthorized),
	)
	GitlabWebhookConnectionNotFound = rerrors.New(
		"gitlab webhook: no external connection found for this webhook id",
		codes.NotFound,
		rerrors.WithHttpStatus(http.StatusNotFound),
	)

	// gitlab connection.
	InvalidInstanceURL = rerrors.New(
		"invalid gitlab instance url",
		codes.InvalidArgument,
		rerrors.WithHttpStatus(http.StatusBadRequest),
	)
	GitlabValidationFailed = rerrors.New(
		"could not verify gitlab token against instance",
		codes.FailedPrecondition,
		rerrors.WithHttpStatus(http.StatusPreconditionFailed),
	)
	GitlabConnectionNotFound = rerrors.New(
		"no gitlab connection found",
		codes.NotFound,
		rerrors.WithHttpStatus(http.StatusNotFound),
	)

	// trello connection.
	TrelloValidationFailed = rerrors.New(
		"could not verify trello credentials",
		codes.FailedPrecondition,
		rerrors.WithHttpStatus(http.StatusPreconditionFailed),
	)

	// telegram connection.
	TelegramValidationFailed = rerrors.New(
		"could not verify telegram bot token",
		codes.FailedPrecondition,
		rerrors.WithHttpStatus(http.StatusPreconditionFailed),
	)

	// anthropic / llm key connection.
	LlmKeyValidationFailed = rerrors.New(
		"could not verify the api key against the provider",
		codes.FailedPrecondition,
		rerrors.WithHttpStatus(http.StatusPreconditionFailed),
	)
	LlmKeyRejected = rerrors.New(
		"the provider rejected the api key",
		codes.FailedPrecondition,
		rerrors.WithHttpStatus(http.StatusPreconditionFailed),
	)
	LlmKeyProviderUnreachable = rerrors.New(
		"could not reach the provider to verify the api key",
		codes.FailedPrecondition,
		rerrors.WithHttpStatus(http.StatusPreconditionFailed),
	)
	LlmKeyModelListingUnsupported = rerrors.New(
		"the provider doesn't support listing models; enter a model id manually to verify the key against it",
		codes.FailedPrecondition,
		rerrors.WithHttpStatus(http.StatusPreconditionFailed),
	)
	LlmKeyRequired = rerrors.New(
		"api key is required to connect this provider",
		codes.InvalidArgument,
		rerrors.WithHttpStatus(http.StatusBadRequest),
	)

	// s3 (byok) connection.
	S3ConnectionValidationFailed = rerrors.New(
		"could not connect to the s3-compatible endpoint with the provided settings",
		codes.FailedPrecondition,
		rerrors.WithHttpStatus(http.StatusPreconditionFailed),
	)

	// couchdb (byok) connection.
	CouchDBConnectionValidationFailed = rerrors.New(
		"could not connect to the couchdb server with the provided settings",
		codes.FailedPrecondition,
		rerrors.WithHttpStatus(http.StatusPreconditionFailed),
	)

	// postgres (byok) connection / per-vault database.
	PostgresConnectionValidationFailed = rerrors.New(
		"could not connect to the postgres server with the provided settings",
		codes.FailedPrecondition,
		rerrors.WithHttpStatus(http.StatusPreconditionFailed),
	)
	PostgresDatabaseAlreadyEnabled = rerrors.New(
		"vault already has a postgres database enabled",
		codes.FailedPrecondition,
		rerrors.WithHttpStatus(http.StatusPreconditionFailed),
		rerrors.WithPreconditionFailure("POSTGRES", "database", pb.UserErrors_PostgresDatabaseAlreadyEnabled.String()),
	)
	NoPostgresDatabaseLinked = rerrors.New(
		"vault has no enabled postgres database; enable one before using postgres tools",
		codes.FailedPrecondition,
		rerrors.WithHttpStatus(http.StatusFailedDependency),
		rerrors.WithPreconditionFailure("POSTGRES", "database", pb.UserErrors_NoPostgresDatabaseLinked.String()),
	)

	// workbench.
	WorkbenchMissingAnthropicConnection = rerrors.New(
		"no anthropic api key connected for this account; connect one before starting this workbench in api_key mode",
		codes.FailedPrecondition,
		rerrors.WithHttpStatus(http.StatusPreconditionFailed),
	)
	WorkbenchAuthModeNotImplemented = rerrors.New(
		"this workbench auth mode is not implemented yet",
		codes.Unimplemented,
	)
	WorkbenchNotConfigured = rerrors.New(
		"workbenches are not enabled for this deployment",
		codes.Unimplemented,
	)
	NoDockerHostsAvailable = rerrors.New(
		"no docker hosts are registered — ask an administrator to add one",
		codes.FailedPrecondition,
	)
	// WorkbenchMissingDockerHost surfaces from resolveClient for a workbenches row that predates
	// the docker_hosts pool (nullable docker_host_id, no backfill — see
	// migrations/061_docker_hosts.sql) — there's no host to resolve a docker client from, so the
	// only recovery is deleting and recreating the workbench.
	WorkbenchMissingDockerHost = rerrors.New(
		"this workbench predates docker host assignment; delete and recreate it",
		codes.FailedPrecondition,
	)

	// tract: builtin execution
	NoVaultForBuiltinTool = rerrors.New(
		"no vault available to run a builtin tool for this tract",
		codes.FailedPrecondition,
	)

	// tract: validation
	TractStepIdInvalid = rerrors.New(
		"step id must match ^[a-z][a-z0-9_]*$ and not be 'trigger'",
		codes.InvalidArgument,
	)
	TractStepIdDuplicate     = rerrors.New("duplicate step id in tract definition", codes.InvalidArgument)
	TractToolNotFound        = rerrors.New("action references a tool that does not exist", codes.InvalidArgument)
	TractConnectionRequired  = rerrors.New("connection_id is required for MoM tool actions", codes.InvalidArgument)
	TractConnectionForbidden = rerrors.New(
		"connection_id is not allowed for builtin tool actions",
		codes.InvalidArgument,
	)
	TractConnectionNotOwned    = rerrors.New("connection_id does not belong to the tract owner", codes.PermissionDenied)
	TractConnectionUuidInvalid = rerrors.New("connection uuid is not a valid uuid", codes.InvalidArgument)
	TractInvalidTemplateRef    = rerrors.New("template reference is not visible from this step", codes.InvalidArgument)
	TractMalformedTemplate     = rerrors.New("malformed template expression", codes.InvalidArgument)
	TractUnknownTemplateVar    = rerrors.New("unknown $-variable in template expression", codes.InvalidArgument)
	TractConditionsOnNonBranch = rerrors.New(
		"only condition steps may declare conditions/then/else",
		codes.InvalidArgument,
	)
	TractParamsOnNonAction = rerrors.New("only action steps may declare params", codes.InvalidArgument)
	TractUnknownStepType   = rerrors.New("unknown step type", codes.InvalidArgument)
	TractStepNotFound      = rerrors.New("referenced step not found in this run", codes.FailedPrecondition)

	// tract: script steps
	TractScriptEngineNotFound = rerrors.New("no script engine registered for this language", codes.InvalidArgument)
	TractScriptParamInvalid   = rerrors.New(
		"script input/output param name must be a valid identifier and unique within its list",
		codes.InvalidArgument,
	)
	TractScriptBindingMismatch = rerrors.New(
		"script step params must have exactly one binding per declared input param",
		codes.InvalidArgument,
	)
	TractScriptRuntimeError       = rerrors.New("script execution failed", codes.FailedPrecondition)
	TractScriptOutputTypeMismatch = rerrors.New(
		"script output does not match its declared output param types",
		codes.FailedPrecondition,
	)

	// tract: llm_call steps
	TractLlmStepInvalid = rerrors.New(
		"llm_call step requires llm_connection_uuid, llm_model, and prompt",
		codes.InvalidArgument,
	)
	TractLlmConnectionProviderMismatch = rerrors.New(
		"llm_connection_uuid does not reference a supported LLM provider connection",
		codes.FailedPrecondition,
	)
	TractLlmCallFailed = rerrors.New("llm call failed", codes.FailedPrecondition)

	// tract: template path navigation (distinct from TractStepNotFound — the base step id was
	// found, but a later segment of the reference path failed to resolve against its output).
	TractTemplateIndexOnNonArray  = rerrors.New("template index used on a value that is not a list", codes.FailedPrecondition)
	TractTemplateIndexOutOfRange  = rerrors.New("template index is out of range", codes.FailedPrecondition)
	TractTemplateFieldOnNonObject = rerrors.New(
		"template field access on a value that is not an object",
		codes.FailedPrecondition,
	)
	TractTemplateFieldNotFound = rerrors.New("referenced field not found in step output", codes.FailedPrecondition)

	// tract: ownership
	TriggerNotFound       = rerrors.New("trigger not found", codes.NotFound)
	TriggerNotOwned       = rerrors.New("trigger does not belong to the caller", codes.PermissionDenied)
	TractNotOwned         = rerrors.New("tract does not belong to the caller", codes.PermissionDenied)
	TractTemplateNotOwned = rerrors.New("template does not belong to the caller", codes.PermissionDenied)

	// tract: provider-linked trigger presets (e.g. gitlab_push)
	TriggerProviderConnectionRequired = rerrors.New(
		"this trigger preset requires an existing connection for its provider - connect it first",
		codes.FailedPrecondition,
		rerrors.WithHttpStatus(http.StatusPreconditionFailed),
	)

	// tract: transport (tracts_api request fields that carry JSON-encoded strings)
	TractRequestFieldInvalidJSON = rerrors.New(
		"request field is not valid JSON",
		codes.InvalidArgument,
		rerrors.WithHttpStatus(http.StatusBadRequest),
	)

	// tract: mcp builtin tools
	TractServiceNotConfigured  = rerrors.New("tract service is not available", codes.Unimplemented)
	McpTractNameRequired       = rerrors.New("name is required and must be a string", codes.InvalidArgument)
	McpTractDefinitionRequired = rerrors.New("definition is required and must be an object", codes.InvalidArgument)
	McpTractUuidRequired       = rerrors.New("tract_uuid is required and must be a string", codes.InvalidArgument)
	McpTriggerUuidRequired     = rerrors.New("trigger_uuid is required and must be a string", codes.InvalidArgument)
	McpTriggerSourceRequired   = rerrors.New("source is required and must be a string", codes.InvalidArgument)
	McpTractUuidInvalid        = rerrors.New("tract_uuid is not a valid uuid", codes.InvalidArgument)
	McpTriggerUuidInvalid      = rerrors.New("trigger_uuid is not a valid uuid", codes.InvalidArgument)
	McpTractUnknownToolName    = rerrors.New("unknown tract builtin tool", codes.InvalidArgument)
)
