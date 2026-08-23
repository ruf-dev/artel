package hub

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"workbenchbridge/internal/chatprotocol"
)

// writeSessionFile writes events to <dir>/<id>.jsonl in the same one-JSON-object-per-line format
// Hub itself writes, so a test can hand-construct a session file that simulates one left behind
// by a previous bridge process.
func writeSessionFile(t *testing.T, dir, id string, events []chatprotocol.Event) {
	t.Helper()

	var lines []byte

	for _, event := range events {
		payload, err := json.Marshal(event)
		if err != nil {
			t.Fatalf("error marshalling event for test session file: %v", err)
		}

		lines = append(lines, payload...)
		lines = append(lines, '\n')
	}

	path := filepath.Join(dir, id+historySessionExt)

	err := os.WriteFile(path, lines, 0o600)
	if err != nil {
		t.Fatalf("error writing test session file %s: %v", path, err)
	}
}

// TestHub_SeqIncrementsAcrossBroadcasts proves Seq is stamped starting at 1 and increases by
// exactly 1 per published event, in call order, for a fresh Hub backed by an empty historyDir.
func TestHub_SeqIncrementsAcrossBroadcasts(t *testing.T) {
	h := New(t.TempDir())

	server := httptest.NewServer(h)
	defer server.Close()

	live := dialChat(t, server.URL)
	defer live.Close()

	h.Broadcast(chatprotocol.Event{Type: chatprotocol.EventSystemInit})
	h.Broadcast(chatprotocol.Event{Type: chatprotocol.EventUserMessage, Text: "hi"})
	h.Broadcast(chatprotocol.Event{Type: chatprotocol.EventTurnDone})

	want := []uint64{1, 2, 3}

	for i, wantSeq := range want {
		event := readOneEvent(t, live)
		if event.Seq != wantSeq {
			t.Fatalf("event %d Seq = %d, want %d", i, event.Seq, wantSeq)
		}
	}
}

// TestHub_RehydratesFromExistingSessionFile simulates a bridge process restart: historyDir
// already holds a session file (written by hand here, matching the real on-disk format) when a
// new Hub is constructed against it. The new Hub must load its events into the backlog in order
// and continue nextSeq from one past the highest Seq found, not restart it at 1.
func TestHub_RehydratesFromExistingSessionFile(t *testing.T) {
	dir := t.TempDir()

	events := []chatprotocol.Event{
		{Type: chatprotocol.EventSystemInit, Seq: 1},
		{Type: chatprotocol.EventUserMessage, Text: "hello", Seq: 2},
		{Type: chatprotocol.EventTurnDone, SessionID: "abc", Seq: 5},
	}

	writeSessionFile(t, dir, "session-before-restart", events)

	h := New(dir)

	if len(h.backlog) != len(events) {
		t.Fatalf("rehydrated backlog length = %d, want %d: %+v", len(h.backlog), len(events), h.backlog)
	}

	for i, want := range events {
		if h.backlog[i].Type != want.Type || h.backlog[i].Seq != want.Seq {
			t.Fatalf("backlog[%d] = %+v, want %+v", i, h.backlog[i], want)
		}
	}

	if h.nextSeq != 6 {
		t.Fatalf("nextSeq after rehydrate = %d, want 6 (max Seq 5 in the file, plus 1)", h.nextSeq)
	}

	// A freshly attaching consumer should replay the rehydrated backlog...
	server := httptest.NewServer(h)
	defer server.Close()

	fresh := dialChat(t, server.URL)
	defer fresh.Close()

	for i, want := range events {
		replayed := readOneEvent(t, fresh)
		if replayed.Type != want.Type || replayed.Seq != want.Seq {
			t.Fatalf("replayed[%d] = %+v, want %+v", i, replayed, want)
		}
	}

	// ...and a new event published after rehydrate must continue the sequence, not restart it.
	h.Broadcast(chatprotocol.Event{Type: chatprotocol.EventAssistantTextDone, Text: "new"})

	next := readOneEvent(t, fresh)
	if next.Seq != 6 {
		t.Fatalf("first Seq published after rehydrate = %d, want 6", next.Seq)
	}
}

// TestHub_ResetRotatesSessionFile proves Reset's on-disk rotation contract: after Reset,
// historyDir holds two session files — the pre-Reset one, untouched, and a new one containing
// only the Reset event.
func TestHub_ResetRotatesSessionFile(t *testing.T) {
	dir := t.TempDir()

	h := New(dir)

	h.Broadcast(chatprotocol.Event{Type: chatprotocol.EventSystemInit})
	h.Broadcast(chatprotocol.Event{Type: chatprotocol.EventUserMessage, Text: "before reset"})

	before, err := filepath.Glob(filepath.Join(dir, "*"+historySessionExt))
	if err != nil {
		t.Fatalf("error globbing history dir: %v", err)
	}

	if len(before) != 1 {
		t.Fatalf("history dir before Reset has %d session files, want 1: %v", len(before), before)
	}

	oldPath := before[0]

	oldEventsBefore, err := readSessionFile(oldPath)
	if err != nil {
		t.Fatalf("error reading old session file before Reset: %v", err)
	}

	if len(oldEventsBefore) != 2 {
		t.Fatalf("old session file before Reset has %d events, want 2", len(oldEventsBefore))
	}

	h.Reset(chatprotocol.Event{Type: chatprotocol.EventNewChat})

	after, err := filepath.Glob(filepath.Join(dir, "*"+historySessionExt))
	if err != nil {
		t.Fatalf("error globbing history dir after Reset: %v", err)
	}

	if len(after) != 2 {
		t.Fatalf("history dir after Reset has %d session files, want 2: %v", len(after), after)
	}

	oldEventsAfter, err := readSessionFile(oldPath)
	if err != nil {
		t.Fatalf("error reading old session file after Reset: %v", err)
	}

	if len(oldEventsAfter) != len(oldEventsBefore) {
		t.Fatalf("old session file's event count changed after Reset: got %d, want %d",
			len(oldEventsAfter), len(oldEventsBefore))
	}

	for i := range oldEventsBefore {
		if oldEventsAfter[i].Type != oldEventsBefore[i].Type || oldEventsAfter[i].Seq != oldEventsBefore[i].Seq {
			t.Fatalf("old session file event %d changed after Reset: got %+v, want %+v",
				i, oldEventsAfter[i], oldEventsBefore[i])
		}
	}

	var newPath string

	for _, path := range after {
		if path != oldPath {
			newPath = path
		}
	}

	if newPath == "" {
		t.Fatalf("could not find a new session file among %v (old one was %s)", after, oldPath)
	}

	newEvents, err := readSessionFile(newPath)
	if err != nil {
		t.Fatalf("error reading new session file: %v", err)
	}

	if len(newEvents) != 1 {
		t.Fatalf("new session file has %d events, want 1 (only the Reset event): %+v", len(newEvents), newEvents)
	}

	if newEvents[0].Type != chatprotocol.EventNewChat {
		t.Fatalf("new session file's event type = %q, want %q", newEvents[0].Type, chatprotocol.EventNewChat)
	}
}

// TestHub_ServeHistoryList proves the /history listing shape: one summary per session file, with
// a non-empty id/lastActivityAt and a first-user-message preview truncated to
// firstUserMessagePreviewRunes runes.
func TestHub_ServeHistoryList(t *testing.T) {
	h := New(t.TempDir())

	longText := "first question about the widget that broke in production yesterday afternoon during the big demo"

	h.Broadcast(chatprotocol.Event{Type: chatprotocol.EventUserMessage, Text: longText})
	h.Broadcast(chatprotocol.Event{Type: chatprotocol.EventAssistantTextDone, Text: "answer"})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/history", nil)

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("/history status = %d, want %d", rec.Code, http.StatusOK)
	}

	var summaries []sessionSummary

	err := json.NewDecoder(rec.Body).Decode(&summaries)
	if err != nil {
		t.Fatalf("error decoding /history response: %v", err)
	}

	if len(summaries) != 1 {
		t.Fatalf("/history returned %d summaries, want 1: %+v", len(summaries), summaries)
	}

	summary := summaries[0]

	if summary.ID == "" {
		t.Fatalf("summary.ID is empty: %+v", summary)
	}

	if summary.LastActivityAt == "" {
		t.Fatalf("summary.LastActivityAt is empty: %+v", summary)
	}

	wantPreview := truncateRunes(longText, firstUserMessagePreviewRunes)
	if summary.FirstUserMessage != wantPreview {
		t.Fatalf("summary.FirstUserMessage = %q, want %q", summary.FirstUserMessage, wantPreview)
	}
}

// TestHub_ServeHistoryDetail proves /history/<id> returns the full ordered event array for that
// session.
func TestHub_ServeHistoryDetail(t *testing.T) {
	h := New(t.TempDir())

	h.Broadcast(chatprotocol.Event{Type: chatprotocol.EventUserMessage, Text: "hello"})
	h.Broadcast(chatprotocol.Event{Type: chatprotocol.EventAssistantTextDone, Text: "hi"})

	h.mu.Lock()
	id := h.sessionID
	h.mu.Unlock()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/history/"+id, nil)

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("/history/%s status = %d, want %d", id, rec.Code, http.StatusOK)
	}

	var events []chatprotocol.Event

	err := json.NewDecoder(rec.Body).Decode(&events)
	if err != nil {
		t.Fatalf("error decoding /history/%s response: %v", id, err)
	}

	if len(events) != 2 {
		t.Fatalf("/history/%s returned %d events, want 2: %+v", id, len(events), events)
	}

	if events[0].Type != chatprotocol.EventUserMessage || events[1].Type != chatprotocol.EventAssistantTextDone {
		t.Fatalf("/history/%s events = %+v, want [user_message, assistant_text_done]", id, events)
	}
}

// TestHub_ServeHistoryDetail_NotFound proves a session id with no matching file responds 404.
func TestHub_ServeHistoryDetail_NotFound(t *testing.T) {
	h := New(t.TempDir())

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/history/does-not-exist", nil)

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("/history/does-not-exist status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

// TestHub_ServeHistoryDetail_RejectsUnsafeID proves a session id containing a ".." traversal
// segment is rejected with 400 before ever being joined into a filesystem path.
func TestHub_ServeHistoryDetail_RejectsUnsafeID(t *testing.T) {
	h := New(t.TempDir())

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/history/placeholder", nil)
	req.URL.Path = "/history/../secret"

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
