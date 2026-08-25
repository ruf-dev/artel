// Package simplechat implements service.SimpleChatService: the in-process, container-free
// "Simple Chat" agent. A vault member picks an OpenRouter BYOK model and chats with an agent
// that can call Artel's existing MCP tools, gated by per-tool permission prompts. Unlike the
// Docker workbench there is no container involved — the whole turn runs in the artel server
// process (see run_turn.go).
//
// A chat thread is personal to its creator: every read/write path checks the caller owns the
// row, not merely that they are a member of its vault.
package simplechat

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

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

type Service struct {
	chatsRepo      repository.SimpleChats
	messagesRepo   repository.SimpleChatMessages
	allowancesRepo repository.SimpleChatToolAllowances
	vaultMembers   repository.VaultMembers
	connections    repository.ExternalConnectionRepo

	mcp mcpService
}

func New(
	chats repository.SimpleChats,
	messages repository.SimpleChatMessages,
	allowances repository.SimpleChatToolAllowances,
	vaultMembers repository.VaultMembers,
	connections repository.ExternalConnectionRepo,
	mcp mcpService,
) *Service {
	return &Service{
		chatsRepo:      chats,
		messagesRepo:   messages,
		allowancesRepo: allowances,
		vaultMembers:   vaultMembers,
		connections:    connections,
		mcp:            mcp,
	}
}

// CreateChat starts a new chat thread for the calling vault member. It refuses up front when the
// caller has no OpenRouter BYOK connection, so the failure surfaces at creation time rather than
// on the first turn once the websocket is already open.
func (s *Service) CreateChat(
	ctx context.Context, vaultId uuid.UUID, model string, vaultAccess bool,
) (domain.SimpleChat, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return domain.SimpleChat{}, rerrors.Wrap(user_errors.Unauthenticated)
	}
	userId := uc.UserUuid

	_, err := s.vaultMembers.Get(ctx, vaultId, userId)
	if err != nil {
		return domain.SimpleChat{}, rerrors.Wrap(user_errors.SimpleChatRequiresVaultMembership)
	}

	_, err = s.resolveOpenRouterCredentials(ctx, userId)
	if err != nil {
		return domain.SimpleChat{}, rerrors.Wrap(err, "error resolving openrouter credentials")
	}

	chat, err := s.chatsRepo.Create(ctx, vaultId, userId, model, vaultAccess)
	if err != nil {
		return domain.SimpleChat{}, rerrors.Wrap(err, "error creating simple chat")
	}

	return chat, nil
}

// ListChats returns the caller's own threads in the vault, most recently active first.
func (s *Service) ListChats(ctx context.Context, vaultId uuid.UUID) ([]domain.SimpleChat, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return nil, rerrors.Wrap(user_errors.Unauthenticated)
	}

	chats, err := s.chatsRepo.ListByVaultAndUser(ctx, vaultId, uc.UserUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing simple chats")
	}

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

	chat, err := s.ownedChat(ctx, chatId, uc.UserUuid)
	if err != nil {
		return domain.SimpleChat{}, nil, rerrors.Wrap(err, "error getting owned simple chat")
	}

	messages, err := s.messagesRepo.ListByChatID(ctx, chatId)
	if err != nil {
		return domain.SimpleChat{}, nil, rerrors.Wrap(err, "error listing simple chat messages")
	}

	return chat, messages, nil
}

// DeleteChat removes an owned thread. Its messages and remembered tool allowances go with it via
// the ON DELETE CASCADE foreign keys in migrations 076/077.
func (s *Service) DeleteChat(ctx context.Context, chatId uuid.UUID) error {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return rerrors.Wrap(user_errors.Unauthenticated)
	}

	_, err := s.ownedChat(ctx, chatId, uc.UserUuid)
	if err != nil {
		return rerrors.Wrap(err, "error getting owned simple chat")
	}

	err = s.chatsRepo.Delete(ctx, chatId)
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
	return s.ownedChat(ctx, chatId, userUuid)
}

func (s *Service) ownedChat(ctx context.Context, chatId, userUuid uuid.UUID) (domain.SimpleChat, error) {
	chat, err := s.chatsRepo.GetByID(ctx, chatId)
	if err != nil {
		return domain.SimpleChat{}, rerrors.Wrap(err, "error getting simple chat by id")
	}

	if chat.UserUuid != userUuid {
		return domain.SimpleChat{}, rerrors.Wrap(user_errors.SimpleChatNotOwned, chatId.String())
	}

	return chat, nil
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

	return creds, nil
}
