// Package telegram_webhook implements the inbound half of Artel's BYO-bot Telegram integration:
// receiving updates from Telegram for a bot the user linked via
// externalconnections.Service.AddTelegramConnection, and relaying them into that user's
// workbench chat session (a headless Claude Code session running in a Docker container, exposed
// as a WebSocket speaking internal/chatprotocol's normalized JSON event protocol).
//
// Structured to mirror internal/transport/gitlab_webhook.Handler closely: a plain http.Handler
// (not gRPC-gateway), narrow injected deps, ServeHTTP does method check -> path-based connection
// id extraction -> external_connections lookup -> constant-time secret comparison -> bounded body
// read -> dispatch, always responding 200 OK so nothing about internal state ever leaks to the
// (unauthenticated-until-the-secret-check) caller.
package telegram_webhook

import (
	"context"
	"crypto/subtle"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

const (
	pathPrefix   = "/webhooks/telegram/"
	maxBodyBytes = 1 << 20 // 1MB
)

// toolExecutor is the narrow subset of service.MomService this handler depends on: executing a
// MoM tool (telegram.send_message / edit_message_text / answer_callback_query, seeded by
// migrations/065_telegram_mom.sql and migrations/073_telegram_webhook_mom.sql) against a specific
// external connection, without an ownership check — this handler has already verified the
// connection via its webhook secret, which stands in for the ownership check
// ExecuteToolForConnection's own doc comment says its other callers must do themselves.
type toolExecutor interface {
	ExecuteToolForConnection(
		ctx context.Context, exConnUuid uuid.UUID, mcpName string, toolName string, params map[string]interface{},
	) (string, error)
}

// recentWorkbenchLookup is the narrow subset of repository.Workbenches this handler depends on —
// just the v1 "which workbench do we relay into" lookup, not the whole workbench CRUD surface.
type recentWorkbenchLookup interface {
	GetMostRecentByUser(ctx context.Context, userID uuid.UUID) (sql.Null[domain.Workbench], error)
}

// bridgeTargetResolver is the narrow subset of service.WorkbenchService this handler depends on —
// mirrors vaults_api.terminalTargetResolver (internal/transport/vaults_api/workbench_terminal.go),
// just the bridge base-URL lookup, not the whole workbench lifecycle surface.
type bridgeTargetResolver interface {
	ResolveTerminalTarget(ctx context.Context, vaultID, userID uuid.UUID) (string, error)
}

// Handler serves POST /webhooks/telegram/{external_connection_id} — the inbound endpoint Telegram
// calls for a user's linked bot. One external_connections row per bot; ChatID is captured from
// the first inbound message and persisted back onto that row.
type Handler struct {
	baseCtx       context.Context
	externalConns repository.ExternalConnectionRepo
	workbenches   recentWorkbenchLookup
	workbenchSvc  bridgeTargetResolver
	mom           toolExecutor

	// sessions holds one *bridgeSession per external connection with a live (or lazily
	// reconnectable) workbench bridge WebSocket, keyed by external_connections.id. A sync.Map
	// rather than a mutex-guarded map because reads (the common case: relaying another message
	// into an already-open session) vastly outnumber writes (dialing a new session).
	sessions sync.Map
}

func New(
	baseCtx context.Context,
	externalConns repository.ExternalConnectionRepo,
	workbenches recentWorkbenchLookup,
	workbenchSvc bridgeTargetResolver,
	mom toolExecutor,
) *Handler {
	return &Handler{
		baseCtx:       baseCtx,
		externalConns: externalConns,
		workbenches:   workbenches,
		workbenchSvc:  workbenchSvc,
		mom:           mom,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	rawId := strings.TrimPrefix(r.URL.Path, pathPrefix)

	exConnId, err := uuid.Parse(rawId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	exConn, err := h.externalConns.GetByID(ctx, exConnId)
	if err != nil {
		wrappedErr := rerrors.Wrap(user_errors.TelegramWebhookConnectionNotFound, "error loading external connection")
		log.Error().
			Err(wrappedErr).
			Str("external_connection_id", exConnId.String()).
			Msg("telegram webhook: connection lookup failed")
		w.WriteHeader(http.StatusNotFound)

		return
	}

	if exConn.Provider != domain.ProviderTelegram {
		log.Error().
			Err(user_errors.TelegramWebhookConnectionNotFound).
			Str("external_connection_id", exConnId.String()).
			Msg("telegram webhook: connection is not a telegram connection")
		w.WriteHeader(http.StatusNotFound)

		return
	}

	var creds domain.TelegramCredentials

	err = json.Unmarshal(exConn.CredentialsJSON, &creds)
	if err != nil {
		wrappedErr := rerrors.Wrap(
			user_errors.TelegramWebhookConnectionNotFound,
			"error unmarshalling telegram credentials",
		)
		log.Error().
			Err(wrappedErr).
			Str("external_connection_id", exConnId.String()).
			Msg("telegram webhook: invalid credentials")
		w.WriteHeader(http.StatusNotFound)

		return
	}

	token := r.Header.Get("X-Telegram-Bot-Api-Secret-Token")

	secretMatches := creds.WebhookSecret != "" &&
		subtle.ConstantTimeCompare([]byte(token), []byte(creds.WebhookSecret)) == 1
	if !secretMatches {
		log.Warn().Str("external_connection_id", exConnId.String()).Msg(user_errors.TelegramWebhookSecretMismatch.Error())
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	limitedBody := io.LimitReader(r.Body, maxBodyBytes)

	body, err := io.ReadAll(limitedBody)
	if err != nil {
		log.Error().Err(err).Str("external_connection_id", exConnId.String()).Msg("telegram webhook: failed to read body")
		w.WriteHeader(http.StatusOK)

		return
	}

	var update telegramUpdate

	err = json.Unmarshal(body, &update)
	if err != nil {
		log.Error().Err(err).Str("external_connection_id", exConnId.String()).Msg("telegram webhook: invalid update payload")
		w.WriteHeader(http.StatusOK)

		return
	}

	// Always 200 from here on — same "never leak internal state" posture as gitlab_webhook.
	h.dispatch(ctx, exConn, creds, update)

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) dispatch(ctx context.Context, exConn domain.ExternalConnection, creds domain.TelegramCredentials, update telegramUpdate) {
	switch {
	case update.CallbackQuery != nil:
		h.handleCallbackQuery(ctx, exConn, *update.CallbackQuery)
	case update.Message != nil && update.Message.Text != "":
		h.handleMessage(ctx, exConn, creds, *update.Message)
	default:
		// Non-text messages (photos, stickers, ...) and any other update kind aren't relayed in
		// v1 — silently ignored, same as gitlab_webhook ignoring events no linked trigger matches.
	}
}

// handleMessage captures ChatID on the connection's first inbound message (if not already set),
// resolves/dials the user's workbench bridge session, and forwards the text either as an
// auth_code_submit (if the bridge is waiting on one) or a plain user_message.
func (h *Handler) handleMessage(
	ctx context.Context, exConn domain.ExternalConnection, creds domain.TelegramCredentials, msg telegramMessage,
) {
	if creds.ChatID == 0 {
		creds.ChatID = msg.Chat.Id

		err := h.persistChatID(ctx, exConn, creds)
		if err != nil {
			log.Error().Err(err).Str("external_connection_id", exConn.Uuid.String()).Msg("telegram webhook: failed to persist chat id")
		}
	}

	sess, err := h.getOrDialSession(ctx, exConn.Uuid, exConn.UserUuid, creds.ChatID)
	if err != nil {
		h.replyToSessionError(ctx, exConn.Uuid, creds.ChatID, err)

		return
	}

	if sess.takeAwaitingAuthCode() {
		sess.sendAuthCodeSubmit(msg.Text)

		return
	}

	sess.sendUserMessage(msg.Text)
}

// handleCallbackQuery translates an inline-keyboard button press into a permission_decision sent
// to the connection's bridge session, acknowledges the callback (required by the Bot API to
// dismiss the button's loading spinner), and edits the original message so the buttons no longer
// look clickable.
func (h *Handler) handleCallbackQuery(ctx context.Context, exConn domain.ExternalConnection, cq telegramCallbackQuery) {
	eventId, decision, ok := decodeCallbackData(cq.Data)
	if !ok {
		log.Warn().Str("external_connection_id", exConn.Uuid.String()).Str("data", cq.Data).
			Msg("telegram webhook: unparsable callback_data")
		h.answerCallbackQuery(ctx, exConn.Uuid, cq.Id, "")

		return
	}

	var creds domain.TelegramCredentials

	err := json.Unmarshal(exConn.CredentialsJSON, &creds)
	if err != nil {
		log.Error().Err(err).Str("external_connection_id", exConn.Uuid.String()).Msg("telegram webhook: invalid credentials")

		return
	}

	sess, err := h.getOrDialSession(ctx, exConn.Uuid, exConn.UserUuid, creds.ChatID)
	if err != nil {
		log.Error().Err(err).Str("external_connection_id", exConn.Uuid.String()).Msg("telegram webhook: failed to resolve session for callback")
		h.answerCallbackQuery(ctx, exConn.Uuid, cq.Id, "could not reach workbench")

		return
	}

	sess.sendPermissionDecision(eventId, decision)

	label := decisionLabel(decision)

	h.answerCallbackQuery(ctx, exConn.Uuid, cq.Id, label)

	if cq.Message == nil {
		return
	}

	sess.finalizePermissionMessage(ctx, eventId, cq.Message.Chat.Id, cq.Message.MessageId, label)
}

// replyToSessionError translates a getOrDialSession failure into a plain-text explanation sent
// back to the Telegram chat, rather than erroring silently — per the v1 scoping rule, a user with
// no workbench at all or a stopped one gets pointed at the web UI instead of a dead relay.
func (h *Handler) replyToSessionError(ctx context.Context, exConnUuid uuid.UUID, chatId int64, err error) {
	log.Error().Err(err).Str("external_connection_id", exConnUuid.String()).Msg("telegram webhook: failed to resolve workbench session")

	if chatId == 0 {
		return
	}

	text := "Something went wrong reaching your workbench. Please try again shortly."

	switch {
	case rerrors.Is(err, user_errors.TelegramNoWorkbench):
		text = "You don't have a workbench set up yet. Head to the Artel web UI to create one first."
	case rerrors.Is(err, user_errors.WorkbenchNotRunning):
		text = "Your workbench isn't running right now. Start it from the Artel web UI, then message me again."
	}

	_, sendErr := h.sendText(ctx, exConnUuid, chatId, text)
	if sendErr != nil {
		log.Error().Err(sendErr).Str("external_connection_id", exConnUuid.String()).Msg("telegram webhook: failed to send session-error explanation")
	}
}

// persistChatID re-marshals creds (now carrying the just-captured ChatID) back onto exConn and
// upserts it — mirrors externalconnections.Service.GenerateGitlabWebhookSecret's
// read-mutate-marshal-Upsert shape for updating one field of a connection's stored credentials.
func (h *Handler) persistChatID(ctx context.Context, exConn domain.ExternalConnection, creds domain.TelegramCredentials) error {
	credJSON, err := json.Marshal(creds)
	if err != nil {
		return rerrors.Wrap(err, "error marshaling telegram credentials")
	}

	exConn.CredentialsJSON = credJSON

	_, err = h.externalConns.Upsert(ctx, exConn)
	if err != nil {
		return rerrors.Wrap(err, "error persisting telegram chat id")
	}

	return nil
}

func decisionLabel(decision permissionDecision) string {
	switch decision {
	case decisionAllowOnce:
		return "Allowed once"
	case decisionAllowAlways:
		return "Allowed always"
	case decisionDeny:
		return "Denied"
	default:
		return string(decision)
	}
}
