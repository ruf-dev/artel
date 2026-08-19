package telegram_webhook

// telegramUpdate is the minimal subset of Telegram's Update object this relay understands — just
// enough to route a plain text message or an inline-keyboard button press. Deliberately not a
// full Bot API SDK type: per the root CLAUDE.md's MoM carve-out, inbound webhook ingestion is
// exactly the kind of thing that gets hand-rolled structs rather than a bespoke Go SDK client.
type telegramUpdate struct {
	Message       *telegramMessage       `json:"message"`
	CallbackQuery *telegramCallbackQuery `json:"callback_query"`
}

type telegramMessage struct {
	MessageId int64        `json:"message_id"`
	Chat      telegramChat `json:"chat"`
	From      telegramUser `json:"from"`
	Text      string       `json:"text"`
}

type telegramChat struct {
	Id int64 `json:"id"`
}

type telegramUser struct {
	Id int64 `json:"id"`
}

type telegramCallbackQuery struct {
	Id      string           `json:"id"`
	Data    string           `json:"data"`
	Message *telegramMessage `json:"message"`
}

// inlineKeyboardButton mirrors Telegram's InlineKeyboardButton: exactly one of CallbackData/Url
// is set per button in this relay's usage (permission_request buttons use CallbackData,
// auth_link's button uses Url).
type inlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data,omitempty"`
	Url          string `json:"url,omitempty"`
}

// inlineKeyboardMarkup mirrors Telegram's InlineKeyboardMarkup — marshaled to a JSON string and
// passed as the reply_markup param to the telegram.send_message / edit_message_text MoM tools
// (see migrations/073_telegram_webhook_mom.sql's comment on why a string, not a nested object).
type inlineKeyboardMarkup struct {
	InlineKeyboard [][]inlineKeyboardButton `json:"inline_keyboard"`
}

// emptyInlineKeyboard clears a message's keyboard when passed as reply_markup — an explicit empty
// array, since omitting reply_markup entirely leaves the existing keyboard untouched instead.
const emptyInlineKeyboard = `{"inline_keyboard":[]}`

// telegramMessageResponse is the shape of a successful sendMessage/editMessageText response, just
// enough to recover the message_id assistant-text-delta coalescing edits later.
type telegramMessageResponse struct {
	Ok     bool `json:"ok"`
	Result struct {
		MessageId int64 `json:"message_id"`
	} `json:"result"`
}
