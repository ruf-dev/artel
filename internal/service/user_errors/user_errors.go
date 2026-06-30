package user_errors

import (
	"net/http"

	"go.redsock.ru/rerrors"
	"google.golang.org/grpc/codes"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
)

var (
	Unauthenticated = rerrors.New("unauthenticated", codes.Unauthenticated)
	NotFound        = rerrors.New("not found", codes.NotFound)
	AlreadyExists   = rerrors.New("already exists", codes.AlreadyExists)

	InvalidCouchDbDatabaseName   = rerrors.New("invalid database name", codes.InvalidArgument)
	CouchDbDatabaseAlreadyExists = rerrors.New("database already exists", codes.FailedPrecondition)
	UserAlreadyExistInCouchDb    = rerrors.New("user already exists in couch db", codes.AlreadyExists)
	CouchDbAccountLocked         = rerrors.New("couch admin account is locked after authentication failures; update instance credentials", codes.FailedPrecondition)

	McpKeyRevoked   = rerrors.New("mcp key revoked", codes.FailedPrecondition, rerrors.WithHttpStatus(http.StatusBadRequest))
	McpInvalidToken = rerrors.New("invalid mcp token", codes.FailedPrecondition, rerrors.WithHttpStatus(http.StatusBadRequest))

	// auth
	SessionExpired         = rerrors.New("session expired", codes.Unauthenticated)
	InvalidTelegramToken   = rerrors.New("invalid telegram token", codes.Unauthenticated)
	UnsupportedLoginMethod = rerrors.New("unsupported login method", codes.InvalidArgument)

	// subscription

	NoActiveSubscription = rerrors.New("no active subscription",
		codes.FailedPrecondition,
		rerrors.WithHttpStatus(http.StatusForbidden),
		rerrors.WithPreconditionFailure("SUB", "subscription", pb.UserErrors_NoSubscription.String()),
	)

	// admin
	NotAdmin = rerrors.New("not an administrator", codes.PermissionDenied)

	// feature flags
	NotesNotEnabled = rerrors.New("notes are not enabled for your account", codes.PermissionDenied)

	// vault
	NotVaultOwner      = rerrors.New("only vault owner can perform this action", codes.PermissionDenied)
	InviteLinkRevoked  = rerrors.New("invite link has been revoked", codes.FailedPrecondition)
	InvalidInviteToken = rerrors.New("invalid invite token", codes.NotFound)

	// middleware
	NoMetadataInContext = rerrors.New("error getting metadata from context", codes.FailedPrecondition)
	NoAuthHeader        = rerrors.New("error getting auth header", codes.Unauthenticated)
	DebugNotSupported   = rerrors.New("debug not supported", codes.Unimplemented)

	// task tracker
	TrelloInvalidCredentials = rerrors.New("invalid trello api key or token", codes.InvalidArgument, rerrors.WithHttpStatus(http.StatusBadRequest))

	// mcp tool argument validation
	McpPathRequired      = rerrors.New("path is required and must be a string", codes.InvalidArgument)
	McpContentRequired   = rerrors.New("content is required and must be a string", codes.InvalidArgument)
	McpOldPathRequired   = rerrors.New("old_path is required and must be a string", codes.InvalidArgument)
	McpNewPathRequired   = rerrors.New("new_path is required and must be a string", codes.InvalidArgument)
	McpAccountIdRequired = rerrors.New("account_id is required and must be a string", codes.InvalidArgument)
	McpIdRequired        = rerrors.New("id is required and must be a string", codes.InvalidArgument)
	McpToRequired        = rerrors.New("to is required and must be a string", codes.InvalidArgument)
	McpSubjectRequired   = rerrors.New("subject is required and must be a string", codes.InvalidArgument)
	McpBodyRequired      = rerrors.New("body is required and must be a string", codes.InvalidArgument)

	// mcp executor
	McpEmailActionMissing          = rerrors.New("email executor: action must have imap or smtp discriminator", codes.InvalidArgument)
	McpUnknownImapOperation        = rerrors.New("email executor: unknown imap operation", codes.InvalidArgument)
	McpUnknownSmtpOperation        = rerrors.New("email executor: unknown smtp operation", codes.InvalidArgument)
	McpToolNotFound                = rerrors.New("tool not found in any connected mcp", codes.NotFound)
	McpConnectionNotOwned          = rerrors.New("external connection not found", codes.NotFound)
	McpConnectorNotFound           = rerrors.New("no mcp connector configured for this key", codes.FailedPrecondition)
	McpActionMissing               = rerrors.New("tool action has no imap, smtp, or http discriminator set", codes.InvalidArgument)
	McpSecretFieldMissing          = rerrors.New("http executor: referenced __secrets field not found in connected credentials", codes.FailedPrecondition)
	McpCredentialsProviderMismatch = rerrors.New("http executor: linked external connection provider does not match tool's declared credentials provider", codes.FailedPrecondition)
	McpHttpRequestFailed           = rerrors.New("http executor: upstream request failed", codes.Unavailable)

	// livesync file-type constraints
	UseReadNoteForTextFiles       = rerrors.New("use read_note for text files", codes.FailedPrecondition)
	UseDeleteNoteForTextFiles     = rerrors.New("use delete_note for text files", codes.FailedPrecondition)
	UseMoveNoteForTextFiles       = rerrors.New("use move_note for text files", codes.FailedPrecondition)
	ChunkedBinaryMoveNotSupported = rerrors.New("move of chunked binary files is not supported; use Obsidian to move large files", codes.FailedPrecondition)

	// imap
	InvalidEmailId       = rerrors.New("invalid email id", codes.InvalidArgument)
	EmailMessageNotFound = rerrors.New("message not found", codes.NotFound)

	// external connections
	GoogleNotConnected = rerrors.New("google account not connected", codes.FailedPrecondition, rerrors.WithHttpStatus(http.StatusPreconditionFailed))
	InvalidOAuthState  = rerrors.New("invalid or expired oauth state", codes.InvalidArgument, rerrors.WithHttpStatus(http.StatusBadRequest))

	NoCouchDbInstance = rerrors.New("no storage instance",
		codes.FailedPrecondition,
		rerrors.WithHttpStatus(http.StatusFailedDependency),
		rerrors.WithPreconditionFailure("COUCHDB", "instance", pb.UserErrors_NoCouchDbInstance.String()))

	// gitlab webhook
	GitlabWebhookSecretMismatch     = rerrors.New("gitlab webhook: token does not match configured secret", codes.PermissionDenied, rerrors.WithHttpStatus(http.StatusUnauthorized))
	GitlabWebhookConnectionNotFound = rerrors.New("gitlab webhook: no external connection found for this webhook id", codes.NotFound, rerrors.WithHttpStatus(http.StatusNotFound))

	// gitlab connection
	InvalidInstanceURL       = rerrors.New("invalid gitlab instance url", codes.InvalidArgument, rerrors.WithHttpStatus(http.StatusBadRequest))
	GitlabValidationFailed   = rerrors.New("could not verify gitlab token against instance", codes.FailedPrecondition, rerrors.WithHttpStatus(http.StatusPreconditionFailed))
	GitlabConnectionNotFound = rerrors.New("no gitlab connection found", codes.NotFound, rerrors.WithHttpStatus(http.StatusNotFound))
)
