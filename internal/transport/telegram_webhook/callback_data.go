package telegram_webhook

import (
	"strings"

	"github.com/ruf-dev/artel/internal/chatprotocol"
)

// permissionDecision mirrors chatprotocol.PermissionDecision — kept as its own type rather than a
// direct alias so this package's callback_data encoding doesn't silently accept a bridge protocol
// value that isn't actually one of the three buttons this relay renders.
type permissionDecision = chatprotocol.PermissionDecision

const (
	decisionAllowOnce   = chatprotocol.DecisionAllowOnce
	decisionAllowAlways = chatprotocol.DecisionAllowAlways
	decisionDeny        = chatprotocol.DecisionDeny
)

// callbackDataSep separates the correlation id from the decision in an inline-keyboard button's
// callback_data. Telegram caps callback_data at 64 bytes, so this intentionally carries nothing
// but the chatprotocol.Event.ID (a bridge-assigned correlation id, not the full tool input) plus
// a short decision token — see permission_request handling in session.go.
const callbackDataSep = ":"

// encodeCallbackData packs eventId + decision into one callback_data string. eventId is assumed
// not to contain callbackDataSep (bridge-assigned correlation ids are plain UUIDs in practice);
// decodeCallbackData splits from the right specifically to tolerate it anyway.
func encodeCallbackData(eventId string, decision permissionDecision) string {
	return eventId + callbackDataSep + string(decision)
}

// decodeCallbackData reverses encodeCallbackData, splitting on the last separator so an eventId
// that happened to contain the separator itself still round-trips. ok is false for anything that
// doesn't parse as "<id><sep><one of the three known decisions>" — e.g. stale/foreign
// callback_data from a message sent before a code change, or a crafted request.
func decodeCallbackData(data string) (eventId string, decision permissionDecision, ok bool) {
	idx := strings.LastIndex(data, callbackDataSep)
	if idx < 0 {
		return "", "", false
	}

	eventId = data[:idx]

	decisionStr := permissionDecision(data[idx+1:])
	if eventId == "" {
		return "", "", false
	}

	switch decisionStr {
	case decisionAllowOnce, decisionAllowAlways, decisionDeny:
		return eventId, decisionStr, true
	default:
		return "", "", false
	}
}
