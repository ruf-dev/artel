package user_errors

import (
	"net/http"

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
	CouchDbAccountLocked         = rerrors.New("couch admin account is locked after authentication failures; update instance credentials", codes.FailedPrecondition)

	McpKeyRevoked   = rerrors.New("mcp key revoked", codes.FailedPrecondition, rerrors.WithHttpStatus(http.StatusBadRequest))
	McpInvalidToken = rerrors.New("invalid mcp token", codes.FailedPrecondition, rerrors.WithHttpStatus(http.StatusBadRequest))

	// auth
	SessionExpired         = rerrors.New("session expired", codes.Unauthenticated)
	InvalidTelegramToken   = rerrors.New("invalid telegram token", codes.Unauthenticated)
	UnsupportedLoginMethod = rerrors.New("unsupported login method", codes.InvalidArgument)

	// subscription
	NoActiveSubscription = rerrors.New("no active subscription", codes.PermissionDenied)

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

	// livesync file-type constraints
	UseReadNoteForTextFiles       = rerrors.New("use read_note for text files", codes.FailedPrecondition)
	UseDeleteNoteForTextFiles     = rerrors.New("use delete_note for text files", codes.FailedPrecondition)
	UseMoveNoteForTextFiles       = rerrors.New("use move_note for text files", codes.FailedPrecondition)
	ChunkedBinaryMoveNotSupported = rerrors.New("move of chunked binary files is not supported; use Obsidian to move large files", codes.FailedPrecondition)

	// imap
	InvalidEmailId       = rerrors.New("invalid email id", codes.InvalidArgument)
	EmailMessageNotFound = rerrors.New("message not found", codes.NotFound)
)
