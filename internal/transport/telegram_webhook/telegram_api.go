package telegram_webhook

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/service/user_errors"
)

// sendText sends a plain-text message via the telegram.send_message MoM tool and returns the
// message id the reply can later be edited through (assistant-text-delta coalescing, permission
// message finalization). Errors are logged and swallowed by callers throughout this package — a
// failed relay message should never take down the read loop processing the rest of the session.
func (h *Handler) sendText(ctx context.Context, exConnUuid uuid.UUID, chatId int64, text string) (int64, error) {
	params := map[string]interface{}{
		"chat_id": strconv.FormatInt(chatId, 10),
		"text":    text,
	}

	result, err := h.mom.ExecuteToolForConnection(ctx, exConnUuid, "telegram", "send_message", params)
	if err != nil {
		return 0, rerrors.Wrap(err, "error sending telegram message")
	}

	return parseMessageId(result)
}

// sendKeyboard sends a text message with an inline keyboard attached (permission_request's three
// buttons, auth_link's single "Open link" button) via the same telegram.send_message tool,
// reply_markup passed as a JSON-serialized string param — see types.go's inlineKeyboardMarkup doc.
func (h *Handler) sendKeyboard(
	ctx context.Context, exConnUuid uuid.UUID, chatId int64, text string, keyboard inlineKeyboardMarkup,
) (int64, error) {
	markupJSON, err := json.Marshal(keyboard)
	if err != nil {
		return 0, rerrors.Wrap(err, "error marshaling inline keyboard")
	}

	params := map[string]interface{}{
		"chat_id":      strconv.FormatInt(chatId, 10),
		"text":         text,
		"reply_markup": string(markupJSON),
	}

	result, err := h.mom.ExecuteToolForConnection(ctx, exConnUuid, "telegram", "send_message", params)
	if err != nil {
		return 0, rerrors.Wrap(err, "error sending telegram message with keyboard")
	}

	return parseMessageId(result)
}

// editText edits a previously sent message's text via telegram.edit_message_text. When
// clearKeyboard is set, its inline keyboard (if any) is cleared too — used once a permission
// decision or a coalesced delta message reaches its final state.
func (h *Handler) editText(
	ctx context.Context, exConnUuid uuid.UUID, chatId, messageId int64, text string, clearKeyboard bool,
) error {
	params := map[string]interface{}{
		"chat_id":    strconv.FormatInt(chatId, 10),
		"message_id": strconv.FormatInt(messageId, 10),
		"text":       text,
	}

	if clearKeyboard {
		params["reply_markup"] = emptyInlineKeyboard
	}

	_, err := h.mom.ExecuteToolForConnection(ctx, exConnUuid, "telegram", "edit_message_text", params)
	if err != nil {
		return rerrors.Wrap(err, "error editing telegram message")
	}

	return nil
}

// answerCallbackQuery acknowledges an inline-keyboard button press, dismissing its loading
// spinner — required by the Bot API on every callback_query. Errors are logged, not returned:
// callers have already applied the user's decision by this point, so a failed acknowledgement
// shouldn't be treated as the whole callback having failed.
func (h *Handler) answerCallbackQuery(ctx context.Context, exConnUuid uuid.UUID, callbackQueryId, text string) {
	params := map[string]interface{}{
		"callback_query_id": callbackQueryId,
		"text":              text,
	}

	_, err := h.mom.ExecuteToolForConnection(ctx, exConnUuid, "telegram", "answer_callback_query", params)
	if err != nil {
		log.Error().Err(err).Str("external_connection_id", exConnUuid.String()).Msg("telegram webhook: failed to answer callback query")
	}
}

func parseMessageId(result string) (int64, error) {
	var resp telegramMessageResponse

	err := json.Unmarshal([]byte(result), &resp)
	if err != nil {
		return 0, rerrors.Wrap(err, "error parsing telegram message response")
	}

	if !resp.Ok {
		return 0, user_errors.TelegramApiCallFailed
	}

	return resp.Result.MessageId, nil
}
