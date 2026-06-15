package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/ruf-dev/artel/internal/clients/couchdb"
	"github.com/ruf-dev/artel/internal/clients/googleapi"
	"github.com/ruf-dev/artel/internal/clients/trello"
	"github.com/ruf-dev/artel/internal/domain"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
)

type Service interface {
	AuthService() AuthService
	VaultService() VaultService
	CouchInstanceService() CouchInstanceService
	AdminCouchService() AdminCouchService
	McpService() McpService
	EmailService() EmailService
	SubscriptionService() SubscriptionService
	PromptService() PromptService
	TaskTrackerService() TaskTrackerService
	NotesService() NotesService
	AdminUsersService() AdminUsersService
	ExternalConnectionService() ExternalConnectionService
	MomService() MomService
}

type AdminUsersService interface {
	ListUsers(ctx context.Context, req domain.ListUsersReq) ([]domain.User, int64, error)
	GetUser(ctx context.Context, userUuid uuid.UUID) (domain.UserDetails, error)
	GetUserSessions(ctx context.Context, userUuid uuid.UUID) ([]domain.Session, error)
}

type AuthService interface {
	Register(ctx context.Context, email, password string) (domain.User, error)
	Login(ctx context.Context, email, password string) (domain.Session, error)
	Logout(ctx context.Context, token string) error
	ValidateToken(ctx context.Context, token string) (domain.User, error)
	LoginViaTelegram(ctx context.Context, idToken string) (domain.Session, error)
	GetMe(ctx context.Context, userUuid uuid.UUID) (domain.User, domain.UserPermissions, error)
	CheckIsAdmin(ctx context.Context, userUuid uuid.UUID) error
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
	SetKeyAccess(ctx context.Context, keyID uuid.UUID, vaultID uuid.UUID, emailAccountID *uuid.UUID) error
	// ResolveKey validates the raw bearer token and returns vault+couch context.
	ResolveKey(ctx context.Context, rawToken string) (domain.McpKeyContext, error)
}

type SubscriptionService interface {
	CheckActive(ctx context.Context, userUuid uuid.UUID) error
}

type ListPromptsParams struct {
	Ids      []string
	Page     uint32
	PageSize uint32
}

type PromptService interface {
	ListPrompts(ctx context.Context, params ListPromptsParams) ([]domain.Prompt, int64, error)
}

type EmailService interface {
	// Account management — called from gRPC transport only.
	AddAccount(ctx context.Context, account domain.EmailAccount) (domain.EmailAccount, error)
	ListAccounts(ctx context.Context, userUuid uuid.UUID) ([]domain.EmailAccount, error)
	DeleteAccount(ctx context.Context, accountUuid uuid.UUID) error
	ListMailServerSuggestions(ctx context.Context, domainPrefix string) ([]domain.MailServerSuggestion, error)

	// Email operations — called from MCP tools.
	ListFolders(ctx context.Context, accountUuid uuid.UUID) ([]string, error)
	ListEmails(ctx context.Context, accountUuid uuid.UUID, limit int) ([]domain.EmailMeta, error)
	ReadEmail(ctx context.Context, accountUuid uuid.UUID, id string) (domain.EmailMessage, error)
	SendEmail(ctx context.Context, accountUuid uuid.UUID, to, subject, body string) error
}

type TaskTrackerService interface {
	AddTracker(ctx context.Context, tracker domain.TaskTracker, creds trello.TaskTrackerCredentials) (domain.TaskTracker, []domain.TrelloBoard, error)
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
}

type MomService interface {
	ListToolsForKey(ctx context.Context, keyId uuid.UUID) ([]domain.McpToolDef, error)
	ExecuteToolForKey(ctx context.Context, keyId uuid.UUID, toolName string, params map[string]interface{}) (string, error)
}

type ExternalConnectionService interface {
	InitiateGoogleOAuth(ctx context.Context, origin string) (authURL string, err error)
	HandleGoogleOAuthCallback(ctx context.Context, code string, state string) (domain.ExternalConnectionMeta, error)
	DisconnectProvider(ctx context.Context, provider string) error
	ListConnections(ctx context.Context) ([]domain.ExternalConnectionMeta, error)
	GetGoogleClient(ctx context.Context) (*googleapi.Client, error)
	GetPickerToken(ctx context.Context) (string, error)
	AddSpreadsheet(ctx context.Context, spreadsheetId string, name string) (domain.McpSpreadsheet, error)
	ListSpreadsheets(ctx context.Context) ([]domain.McpSpreadsheet, error)
	RemoveSpreadsheet(ctx context.Context, spreadsheetId string) error
}
