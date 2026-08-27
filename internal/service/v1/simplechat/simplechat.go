// Package simplechat implements service.SimpleChatService: the in-process, container-free
// "Simple Chat" agent. A vault member picks an OpenRouter BYOK model and chats with an agent
// that can call Artel's existing MCP tools, gated by per-tool permission prompts. Unlike the
// Docker workbench there is no container involved — the whole turn runs in the artel server
// process (see run_turn.go).
//
// A chat thread's transcript lives as a JSONL doc in its vault's own CouchDB database, under the
// reserved .chat_history/<user_id>/<chat_id>.jsonl path (see domain.ChatHistoryPath) — the same
// per-vault storage notes and skills already use, hidden from the regular Notes/Files listings.
// A chat thread is personal to its creator: every read/write path checks the caller owns it, not
// merely that they are a member of its vault — see ownedChat.
package simplechat

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/ruf-dev/artel/internal/clients/couchdb"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

// openrouterDefaultBaseUrl mirrors externalconnections.openrouterDefaultBaseUrl (unexported
// there, so duplicated here rather than introducing a cross-package const just for this). It
// re-defaults credentials persisted before that package always resolved a blank BaseUrl at
// connection time (see resolveOpenRouterCredentials).
const openrouterDefaultBaseUrl = "https://openrouter.ai/api/v1"

// mcpService is the narrow subset of service.McpService the agent loop depends on — kept local
// (rather than importing the whole service package) to match the dependency style used by
// internal/service/v1/workbench.
type mcpService interface {
	BuildKeyContext(ctx context.Context, vaultUuid, userUuid uuid.UUID) (domain.McpKeyContext, error)
	ListTools(ctx context.Context) ([]domain.McpToolDef, error)
	ExecuteTool(
		ctx context.Context, keyCtx domain.McpKeyContext, toolName string, params map[string]interface{},
	) (domain.ToolExecResult, error)
}

// Service's vault/CouchDB dependencies mirror internal/service/v1/skills.Service exactly — both
// resolve a vault's LiveSyncClient the same way (vaults + couchInstances + couchAccounts).
type Service struct {
	vaults         repository.Vaults
	vaultMembers   repository.VaultMembers
	couchInstances repository.CouchInstances
	couchAccounts  repository.CouchAccounts
	connections    repository.ExternalConnectionRepo
	systemSettings repository.SystemSettingsRepo
	userSettings   repository.UserSettingsRepo

	mcp mcpService
}

func New(
	vaults repository.Vaults,
	vaultMembers repository.VaultMembers,
	couchInstances repository.CouchInstances,
	couchAccounts repository.CouchAccounts,
	connections repository.ExternalConnectionRepo,
	systemSettings repository.SystemSettingsRepo,
	userSettings repository.UserSettingsRepo,
	mcp mcpService,
) *Service {
	return &Service{
		vaults:         vaults,
		vaultMembers:   vaultMembers,
		couchInstances: couchInstances,
		couchAccounts:  couchAccounts,
		connections:    connections,
		systemSettings: systemSettings,
		userSettings:   userSettings,
		mcp:            mcp,
	}
}

// CreateChat starts a new chat thread for the calling vault member. It refuses up front when the
// caller has no OpenRouter BYOK connection, so the failure surfaces at creation time rather than
// on the first turn once the websocket is already open. The thread's transcript file is written
// immediately with a header line only. ListChats returns it so the workbench can restore the
// selected mode and thread even before the first message is sent.
func (s *Service) CreateChat(
	ctx context.Context, vaultId uuid.UUID, model string, vaultAccess bool,
) (domain.SimpleChat, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return domain.SimpleChat{}, rerrors.Wrap(user_errors.Unauthenticated)
	}
	userId := uc.UserUuid

	vault, err := s.resolveMemberVault(ctx, vaultId, userId)
	if err != nil {
		return domain.SimpleChat{}, err
	}

	_, err = s.resolveOpenRouterCredentials(ctx, userId)
	if err != nil {
		return domain.SimpleChat{}, rerrors.Wrap(err, "error resolving openrouter credentials")
	}

	client, err := s.liveSyncClientForVault(ctx, vault, userId)
	if err != nil {
		return domain.SimpleChat{}, err
	}

	now := time.Now().UTC()
	chatId := uuid.New()

	header := domain.SimpleChatHeader{
		Uuid:           chatId,
		VaultUuid:      vaultId,
		UserUuid:       userId,
		Model:          model,
		VaultAccess:    vaultAccess,
		CreatedAt:      now,
		UpdatedAt:      now,
		LastActivityAt: now,
	}

	file := domain.SimpleChatFile{Header: header}

	err = s.writeChatFile(ctx, client, userId, chatId, file)
	if err != nil {
		return domain.SimpleChat{}, rerrors.Wrap(err, "error creating simple chat")
	}

	return header.ToSimpleChat(), nil
}

// ListChats returns the caller's own threads in the vault, most recently active first. Empty
// threads are included because their existence is the durable record that the caller selected
// Simple Chat for this vault.
func (s *Service) ListChats(ctx context.Context, vaultId uuid.UUID) ([]domain.SimpleChat, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return nil, rerrors.Wrap(user_errors.Unauthenticated)
	}
	userId := uc.UserUuid

	vault, err := s.resolveMemberVault(ctx, vaultId, userId)
	if err != nil {
		return nil, err
	}

	client, err := s.liveSyncClientForVault(ctx, vault, userId)
	if err != nil {
		return nil, err
	}

	notes, err := client.ListChatNotes(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing simple chat notes")
	}

	userPrefix := domain.ChatHistoryFolderPrefix + userId.String() + "/"

	chats := make([]domain.SimpleChat, 0, len(notes))

	for _, note := range notes {
		if !strings.HasPrefix(note.Path, userPrefix) {
			continue
		}

		noteDoc, readErr := client.ReadNote(ctx, note.Path)
		if readErr != nil {
			return nil, rerrors.Wrap(readErr, "error reading simple chat note")
		}

		file, decodeErr := domain.DecodeSimpleChatFile([]byte(noteDoc.Content))
		if decodeErr != nil {
			return nil, rerrors.Wrap(decodeErr, "error decoding simple chat file")
		}

		chats = append(chats, file.Header.ToSimpleChat())
	}

	domain.SortSimpleChatsByLastActivityDesc(chats)

	return chats, nil
}

// GetChat returns one owned thread plus its full transcript, ordered by seq.
func (s *Service) GetChat(
	ctx context.Context, chatId uuid.UUID,
) (domain.SimpleChat, []domain.SimpleChatMessage, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return domain.SimpleChat{}, nil, rerrors.Wrap(user_errors.Unauthenticated)
	}

	chat, _, file, err := s.ownedChat(ctx, chatId, uc.UserUuid)
	if err != nil {
		return domain.SimpleChat{}, nil, rerrors.Wrap(err, "error getting owned simple chat")
	}

	return chat, file.Messages, nil
}

// DeleteChat removes an owned thread — its transcript file (messages and remembered tool
// allowances included, both part of the same JSONL doc) is deleted in one call.
func (s *Service) DeleteChat(ctx context.Context, chatId uuid.UUID) error {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return rerrors.Wrap(user_errors.Unauthenticated)
	}

	_, client, _, err := s.ownedChat(ctx, chatId, uc.UserUuid)
	if err != nil {
		return rerrors.Wrap(err, "error getting owned simple chat")
	}

	path := domain.ChatHistoryPath(uc.UserUuid, chatId)

	err = client.DeleteNote(ctx, path)
	if err != nil {
		return rerrors.Wrap(err, "error deleting simple chat")
	}

	return nil
}

// OwnedChat loads a chat and asserts userUuid created it. Exported for the websocket transport,
// which authenticates from a cookie rather than a user_context and so cannot go through the
// context-based methods above.
func (s *Service) OwnedChat(
	ctx context.Context, chatId, userUuid uuid.UUID,
) (domain.SimpleChat, error) {
	chat, _, _, err := s.ownedChat(ctx, chatId, userUuid)

	return chat, err
}

// ownedChat resolves chatId by scanning userUuid's vault memberships for the one whose CouchDB
// holds domain.ChatHistoryPath(userUuid, chatId) — that path is already scoped to userUuid, so
// finding it there is itself the ownership proof; a wrong or someone-else's chat id simply isn't
// found in any of the caller's vaults. Bounded by how many vaults userUuid belongs to, not a
// scan of chat data. Returns the resolved client and decoded file alongside the chat so callers
// (GetChat, DeleteChat) don't have to re-resolve or re-read.
func (s *Service) ownedChat(
	ctx context.Context, chatId, userUuid uuid.UUID,
) (domain.SimpleChat, *couchdb.LiveSyncClient, domain.SimpleChatFile, error) {
	vaults, err := s.vaults.ListByMembership(ctx, userUuid)
	if err != nil {
		return domain.SimpleChat{}, nil, domain.SimpleChatFile{}, rerrors.Wrap(err, "error listing vaults by membership")
	}

	for _, vault := range vaults {
		client, clientErr := s.liveSyncClientForVault(ctx, vault, userUuid)
		if clientErr != nil {
			return domain.SimpleChat{}, nil, domain.SimpleChatFile{}, clientErr
		}

		file, readErr := s.readChatFile(ctx, client, userUuid, chatId)
		if readErr != nil {
			continue
		}

		return file.Header.ToSimpleChat(), client, file, nil
	}

	err = rerrors.Wrap(user_errors.SimpleChatNotOwned, chatId.String())

	return domain.SimpleChat{}, nil, domain.SimpleChatFile{}, err
}

// resolveMemberVault checks userId is a member of vaultId and returns the vault, for building
// its LiveSyncClient. CreateChat and ListChats both start from a known vaultId and need this;
// GetChat/DeleteChat/OwnedChat only have a bare chatId to start from and use ownedChat's
// membership-scan instead.
func (s *Service) resolveMemberVault(ctx context.Context, vaultId, userId uuid.UUID) (domain.Vault, error) {
	_, err := s.vaultMembers.Get(ctx, vaultId, userId)
	if err != nil {
		return domain.Vault{}, rerrors.Wrap(user_errors.SimpleChatRequiresVaultMembership)
	}

	vault, err := s.vaults.GetByID(ctx, vaultId)
	if err != nil {
		return domain.Vault{}, rerrors.Wrap(err, "error getting vault")
	}

	return vault, nil
}

// liveSyncClientForVault mirrors internal/service/v1/skills.Service.liveSyncClientForVault
// exactly, except userUuid is passed explicitly rather than pulled from ctx via user_context —
// the websocket transport (OwnedChat/RunTurn) authenticates from a cookie, not a user_context.
func (s *Service) liveSyncClientForVault(
	ctx context.Context, vault domain.Vault, userUuid uuid.UUID,
) (*couchdb.LiveSyncClient, error) {
	instance, err := s.couchInstances.Get(ctx, vault.CouchInstanceUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "error getting couch instance")
	}

	account, err := s.couchAccounts.GetByUserAndInstance(ctx, userUuid, vault.CouchInstanceUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "error getting couch account")
	}

	client := couchdb.NewLiveSyncClient(
		instance.Url,
		vault.CouchDBName,
		account.CouchUsername,
		account.CouchPassword,
	)

	return client, nil
}

// readChatFile reads and decodes chatUuid's transcript from client.
func (s *Service) readChatFile(
	ctx context.Context, client *couchdb.LiveSyncClient, userUuid, chatUuid uuid.UUID,
) (domain.SimpleChatFile, error) {
	path := domain.ChatHistoryPath(userUuid, chatUuid)

	note, err := client.ReadNote(ctx, path)
	if err != nil {
		return domain.SimpleChatFile{}, rerrors.Wrap(err, "error reading simple chat note")
	}

	file, err := domain.DecodeSimpleChatFile([]byte(note.Content))
	if err != nil {
		return domain.SimpleChatFile{}, rerrors.Wrap(err, "error decoding simple chat file")
	}

	return file, nil
}

// writeChatFile encodes file and writes it back to client as a whole — every mutation (a new
// message, a tool allowance, a last-activity bump) is a full read-modify-write of the same doc.
// There is no locking beyond WriteNote's own _rev optimistic-concurrency fetch: a genuine
// concurrent write (e.g. two tabs open on the same chat) can lose an update, which is accepted
// here since a turn runs single-goroutine and sequential.
func (s *Service) writeChatFile(
	ctx context.Context, client *couchdb.LiveSyncClient, userUuid, chatUuid uuid.UUID, file domain.SimpleChatFile,
) error {
	path := domain.ChatHistoryPath(userUuid, chatUuid)

	content, err := domain.EncodeSimpleChatFile(file)
	if err != nil {
		return rerrors.Wrap(err, "error encoding simple chat file")
	}

	err = client.WriteNote(ctx, path, string(content))
	if err != nil {
		return rerrors.Wrap(err, "error writing simple chat note")
	}

	return nil
}

// chatClient resolves an already-known chat's vault and builds its LiveSyncClient — the shared
// first step behind mutateChatFile and currentChatFile.
func (s *Service) chatClient(ctx context.Context, chat domain.SimpleChat) (*couchdb.LiveSyncClient, error) {
	vault, err := s.vaults.GetByID(ctx, chat.VaultUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "error getting vault")
	}

	client, err := s.liveSyncClientForVault(ctx, vault, chat.UserUuid)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// currentChatFile is the read-only counterpart to mutateChatFile — for a caller that only needs
// to inspect the current file (RunTurn's history/next-seq, resolvePermission's stored
// allowances) without writing anything back.
func (s *Service) currentChatFile(ctx context.Context, chat domain.SimpleChat) (domain.SimpleChatFile, error) {
	client, err := s.chatClient(ctx, chat)
	if err != nil {
		return domain.SimpleChatFile{}, err
	}

	file, err := s.readChatFile(ctx, client, chat.UserUuid, chat.Uuid)
	if err != nil {
		return domain.SimpleChatFile{}, err
	}

	return file, nil
}

// mutateChatFile is the shared read-modify-write helper run_turn.go uses to append a message,
// upsert a tool allowance, or bump last-activity — one call in, one full file rewrite out.
func (s *Service) mutateChatFile(
	ctx context.Context, chat domain.SimpleChat, mutate func(*domain.SimpleChatFile),
) error {
	client, err := s.chatClient(ctx, chat)
	if err != nil {
		return err
	}

	file, err := s.readChatFile(ctx, client, chat.UserUuid, chat.Uuid)
	if err != nil {
		return err
	}

	mutate(&file)

	err = s.writeChatFile(ctx, client, chat.UserUuid, chat.Uuid, file)
	if err != nil {
		return err
	}

	return nil
}

// appendMessage appends msg to chat's transcript file as one full read-modify-write. Seq/ChatUuid
// on msg are stamped here rather than by the caller — see domain.SimpleChatFile.NextSeq.
func (s *Service) appendMessage(ctx context.Context, chat domain.SimpleChat, msg domain.SimpleChatMessage) error {
	appendFn := func(file *domain.SimpleChatFile) {
		msg.ChatUuid = chat.Uuid
		msg.Seq = file.NextSeq()

		if msg.CreatedAt.IsZero() {
			msg.CreatedAt = time.Now().UTC()
		}

		file.Messages = append(file.Messages, msg)
	}

	err := s.mutateChatFile(ctx, chat, appendFn)
	if err != nil {
		return rerrors.Wrap(err, "error appending simple chat message")
	}

	return nil
}

// getToolAllowance returns chat's remembered decision for toolName, if any — the CouchDB-era
// replacement for repository.SimpleChatToolAllowances.Get.
func (s *Service) getToolAllowance(ctx context.Context, chat domain.SimpleChat, toolName string) (string, bool, error) {
	file, err := s.currentChatFile(ctx, chat)
	if err != nil {
		return "", false, rerrors.Wrap(err, "error reading simple chat file")
	}

	decision, ok := file.Header.ToolAllowances[toolName]

	return decision, ok, nil
}

// upsertToolAllowance remembers decision for (chat, toolName) — the CouchDB-era replacement for
// repository.SimpleChatToolAllowances.Upsert.
func (s *Service) upsertToolAllowance(ctx context.Context, chat domain.SimpleChat, toolName, decision string) error {
	upsertFn := func(file *domain.SimpleChatFile) {
		if file.Header.ToolAllowances == nil {
			file.Header.ToolAllowances = make(map[string]string)
		}

		file.Header.ToolAllowances[toolName] = decision
	}

	err := s.mutateChatFile(ctx, chat, upsertFn)
	if err != nil {
		return rerrors.Wrap(err, "error upserting simple chat tool allowance")
	}

	return nil
}

// resolveOpenRouterCredentials loads the caller's OpenRouter BYOK credentials, mapping an absent
// connection onto SimpleChatMissingOpenRouterConnection.
func (s *Service) resolveOpenRouterCredentials(
	ctx context.Context, userUuid uuid.UUID,
) (domain.OpenAIKeyCredentials, error) {
	result, err := s.connections.GetByUserAndProvider(ctx, userUuid, domain.ProviderOpenRouter)
	if err != nil {
		return domain.OpenAIKeyCredentials{}, rerrors.Wrap(err, "error getting openrouter connection")
	}

	if !result.Valid {
		return domain.OpenAIKeyCredentials{}, user_errors.SimpleChatMissingOpenRouterConnection
	}

	var creds domain.OpenAIKeyCredentials

	err = json.Unmarshal(result.V.CredentialsJSON, &creds)
	if err != nil {
		return domain.OpenAIKeyCredentials{}, rerrors.Wrap(err, "error parsing openrouter credentials")
	}

	// Connections persisted before externalconnections.AddOpenAIConnection always resolved a
	// blank BaseUrl at write time may still have an empty BaseUrl stored. Re-default it here so
	// the OpenAI SDK client doesn't silently fall back to its own hardcoded OpenAI endpoint.
	if creds.BaseUrl == "" {
		creds.BaseUrl = openrouterDefaultBaseUrl
	}

	return creds, nil
}
