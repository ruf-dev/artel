package claudecli

import (
	"testing"

	"workbenchbridge/internal/chatprotocol"
)

// findEvent returns the first event of the given type in events, failing the test if none is
// found.
func findEvent(t *testing.T, events []chatprotocol.Event, eventType chatprotocol.EventType) chatprotocol.Event {
	t.Helper()

	for _, event := range events {
		if event.Type == eventType {
			return event
		}
	}

	t.Fatalf("no event of type %q found in %+v", eventType, events)

	return chatprotocol.Event{}
}

// findEvents returns every event of the given type in events, in order.
func findEvents(events []chatprotocol.Event, eventType chatprotocol.EventType) []chatprotocol.Event {
	var matched []chatprotocol.Event

	for _, event := range events {
		if event.Type == eventType {
			matched = append(matched, event)
		}
	}

	return matched
}

// TestTranslator_TextDeltaDoneCorrelation covers the assistant_text_delta / assistant_text_done id
// correlation that translate.go's textBlockSeq machinery exists for — see the regression subtest
// for the concrete bug being fixed.
func TestTranslator_TextDeltaDoneCorrelation(t *testing.T) {
	// RawIndexDivergesFromFinalIndex is the regression test for the bug this file fixes: a block
	// type present in the raw stream (simulated here by starting the text_delta at raw index 1, as
	// if index 0 were a "thinking" block) but absent from the final assembled message (simulated by
	// the "assistant" line's single text block sitting at array position 0) used to produce
	// mismatched delta/done ids, rendering as two separate chat bubbles instead of one being
	// replaced by the other.
	t.Run("RawIndexDivergesFromFinalIndex", func(t *testing.T) {
		translator := NewTranslator()

		messageStart := []byte(`{"type":"stream_event","event":{"type":"message_start","message":{"id":"msg_1"}}}`)
		translator.Translate(messageStart)

		// Raw index 1: index 0 is assumed consumed by a block type (e.g. "thinking") that never
		// surfaces as a translate.go event.
		textDelta := []byte(`{"type":"stream_event","event":{"type":"content_block_delta","index":1,"delta":{"type":"text_delta","text":"Hello"}}}`)
		deltaEvents := translator.Translate(textDelta)

		deltaEvent := findEvent(t, deltaEvents, chatprotocol.EventAssistantTextDelta)

		// Final assembled message has only the one text block, at array position 0 — the absent
		// block type never appears here.
		assistantLine := []byte(`{"type":"assistant","message":{"id":"msg_1","content":[{"type":"text","text":"Hello"}]}}`)
		doneEvents := translator.Translate(assistantLine)

		doneEvent := findEvent(t, doneEvents, chatprotocol.EventAssistantTextDone)

		if deltaEvent.ID != doneEvent.ID {
			t.Fatalf("delta id %q does not match done id %q", deltaEvent.ID, doneEvent.ID)
		}
	})

	// InterleavedToolUse confirms two text blocks separated by a tool_use block still get distinct,
	// stable ids that agree between their own delta and done events, and that the tool_use block's
	// id (which comes from the block itself, not blockId) is unaffected by the text-only counter.
	t.Run("InterleavedToolUse", func(t *testing.T) {
		translator := NewTranslator()

		messageStart := []byte(`{"type":"stream_event","event":{"type":"message_start","message":{"id":"msg_2"}}}`)
		translator.Translate(messageStart)

		firstDeltaLine := []byte(`{"type":"stream_event","event":{"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"first"}}}`)
		firstDeltaEvents := translator.Translate(firstDeltaLine)
		firstDelta := findEvent(t, firstDeltaEvents, chatprotocol.EventAssistantTextDelta)

		secondDeltaLine := []byte(`{"type":"stream_event","event":{"type":"content_block_delta","index":2,"delta":{"type":"text_delta","text":"second"}}}`)
		secondDeltaEvents := translator.Translate(secondDeltaLine)
		secondDelta := findEvent(t, secondDeltaEvents, chatprotocol.EventAssistantTextDelta)

		if firstDelta.ID == secondDelta.ID {
			t.Fatalf("expected distinct ids for the two text deltas, got %q for both", firstDelta.ID)
		}

		assistantLine := []byte(`{"type":"assistant","message":{"id":"msg_2","content":[` +
			`{"type":"text","text":"first"},` +
			`{"type":"tool_use","id":"toolu_1","name":"Bash","input":{}},` +
			`{"type":"text","text":"second"}` +
			`]}}`)
		events := translator.Translate(assistantLine)

		doneEvents := findEvents(events, chatprotocol.EventAssistantTextDone)
		if len(doneEvents) != 2 {
			t.Fatalf("expected 2 assistant_text_done events, got %d: %+v", len(doneEvents), doneEvents)
		}

		firstDone := doneEvents[0]
		secondDone := doneEvents[1]

		if firstDelta.ID != firstDone.ID {
			t.Fatalf("first text block: delta id %q does not match done id %q", firstDelta.ID, firstDone.ID)
		}

		if secondDelta.ID != secondDone.ID {
			t.Fatalf("second text block: delta id %q does not match done id %q", secondDelta.ID, secondDone.ID)
		}

		if firstDone.ID == secondDone.ID {
			t.Fatalf("expected distinct ids for the two text_done events, got %q for both", firstDone.ID)
		}

		toolStarted := findEvent(t, events, chatprotocol.EventToolCallStarted)
		if toolStarted.ID != "toolu_1" {
			t.Fatalf("expected tool_call_started id %q (from part.Id), got %q", "toolu_1", toolStarted.ID)
		}
	})

	// SingleTextBlockNoDivergence covers the common case with no divergence between raw stream
	// index and final assembled index — both are 0 — to confirm the fix didn't break the
	// unremarkable path.
	t.Run("SingleTextBlockNoDivergence", func(t *testing.T) {
		translator := NewTranslator()

		messageStart := []byte(`{"type":"stream_event","event":{"type":"message_start","message":{"id":"msg_3"}}}`)
		translator.Translate(messageStart)

		textDelta := []byte(`{"type":"stream_event","event":{"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"hi"}}}`)
		deltaEvents := translator.Translate(textDelta)
		deltaEvent := findEvent(t, deltaEvents, chatprotocol.EventAssistantTextDelta)

		assistantLine := []byte(`{"type":"assistant","message":{"id":"msg_3","content":[{"type":"text","text":"hi"}]}}`)
		doneEvents := translator.Translate(assistantLine)
		doneEvent := findEvent(t, doneEvents, chatprotocol.EventAssistantTextDone)

		if deltaEvent.ID != doneEvent.ID {
			t.Fatalf("delta id %q does not match done id %q", deltaEvent.ID, doneEvent.ID)
		}

		const expectedId = "msg_3:0"
		if deltaEvent.ID != expectedId {
			t.Fatalf("expected id %q, got %q", expectedId, deltaEvent.ID)
		}
	})
}
