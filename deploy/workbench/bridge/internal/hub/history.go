package hub

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"workbenchbridge/internal/chatprotocol"
)

// historySessionExt is the file extension a session's JSONL history file is stored under, both
// when picking the most recent one to rehydrate from (loadHistory) and when minting a new one
// (startNewSession).
const historySessionExt = ".jsonl"

// firstUserMessagePreviewRunes bounds how much of a session's first user_message event is kept
// for the /history listing's row title — see firstUserMessagePreview.
const firstUserMessagePreviewRunes = 80

// loadHistory prepares h's on-disk persistence: it creates h.historyDir if missing, then either
// rehydrates from the most recently modified existing session file in it — a bridge process
// restart mid-session, since h.historyDir survives the same container being stopped and started
// again (see New's doc comment) — or starts a brand-new session file when there's no history to
// continue at all. Called once, from New, before h is shared with any goroutine, so it needs no
// locking of its own.
func (h *Hub) loadHistory() {
	err := os.MkdirAll(h.historyDir, 0o700)
	if err != nil {
		log.Printf("hub: error creating history dir %s: %v", h.historyDir, err)

		return
	}

	path, found, err := latestSessionFile(h.historyDir)
	if err != nil {
		log.Printf("hub: error scanning history dir %s: %v", h.historyDir, err)
	}

	if found {
		h.rehydrateFrom(path)

		return
	}

	h.startNewSession()
	h.nextSeq = 1
}

// latestSessionFile returns the most-recently-modified *.jsonl file directly inside dir, if any.
func latestSessionFile(dir string) (path string, found bool, err error) {
	matches, err := filepath.Glob(filepath.Join(dir, "*"+historySessionExt))
	if err != nil {
		return "", false, err
	}

	var (
		latestPath string
		latestMod  time.Time
	)

	for _, match := range matches {
		info, statErr := os.Stat(match)
		if statErr != nil {
			log.Printf("hub: error stating candidate history file %s: %v", match, statErr)

			continue
		}

		if latestPath == "" || info.ModTime().After(latestMod) {
			latestPath = match
			latestMod = info.ModTime()
		}
	}

	return latestPath, latestPath != "", nil
}

// rehydrateFrom loads path's events into h.backlog in order, advances h.nextSeq past the highest
// Seq seen among them (or to 1, if the file held none), and reopens path for appending as h's
// active session file. This is the bridge-process-restarted-mid-session path.
func (h *Hub) rehydrateFrom(path string) {
	events, err := readSessionFile(path)
	if err != nil {
		log.Printf("hub: error reading session file %s: %v", path, err)
	}

	var maxSeq uint64

	for _, event := range events {
		if event.Seq > maxSeq {
			maxSeq = event.Seq
		}
	}

	h.backlog = events
	h.nextSeq = maxSeq + 1
	h.sessionID = strings.TrimSuffix(filepath.Base(path), historySessionExt)

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		log.Printf("hub: error reopening session file %s for append: %v", path, err)

		return
	}

	h.sessionFile = f
}

// startNewSession mints a fresh session id and creates <h.historyDir>/<id>.jsonl as h's active
// session file. It does not touch h.nextSeq: Seq is monotonic for this Hub's entire lifetime,
// spanning any number of sessions, so only loadHistory (via its two call paths) ever sets it —
// once, at startup.
func (h *Hub) startNewSession() {
	id := newSessionID()
	path := filepath.Join(h.historyDir, id+historySessionExt)

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		log.Printf("hub: error creating session file %s: %v", path, err)

		return
	}

	h.sessionFile = f
	h.sessionID = id
}

// rotateSession closes h's current active session file, if any, and starts a brand-new one via
// startNewSession — the old file is left on disk untouched, becoming a complete, browsable past
// session. Caller must hold h.mu; used by Reset, before the new-chat event is written.
func (h *Hub) rotateSession() {
	if h.sessionFile != nil {
		err := h.sessionFile.Close()
		if err != nil {
			log.Printf("hub: error closing session file for session %s: %v", h.sessionID, err)
		}
	}

	h.startNewSession()
}

// appendToSession writes event as one JSON line to h's active session file, if any. There is no
// application-level buffering to flush: chat traffic is low-volume enough that a Write call per
// event is fine, and a raw *os.File keeps that trivially true rather than needing an explicit
// Flush() after every write. Caller must hold h.mu, and must call this only from within publish,
// after mutateBacklog, so the file and h.backlog can never disagree about which events exist.
func (h *Hub) appendToSession(event chatprotocol.Event) {
	if h.sessionFile == nil {
		return
	}

	payload, err := json.Marshal(event)
	if err != nil {
		log.Printf("hub: error marshalling event for history: %v", err)

		return
	}

	payload = append(payload, '\n')

	_, err = h.sessionFile.Write(payload)
	if err != nil {
		log.Printf("hub: error writing to session file for session %s: %v", h.sessionID, err)
	}
}

// newSessionID mints a collision-safe session id without pulling in a UUID dependency this module
// doesn't otherwise need: a base36 nanosecond timestamp plus a few bytes of crypto/rand, hex
// encoded.
func newSessionID() string {
	randomSuffix := make([]byte, 4)

	_, err := rand.Read(randomSuffix)
	if err != nil {
		log.Printf("hub: error reading random bytes for session id: %v", err)
	}

	return strconv.FormatInt(time.Now().UnixNano(), 36) + hex.EncodeToString(randomSuffix)
}

// readSessionFile parses path line by line as chatprotocol.Event, in order. A line that fails to
// parse is logged and skipped rather than aborting the whole read — one corrupted line (say, a
// write torn by a crash) shouldn't lose the rest of a session's history.
func readSessionFile(path string) ([]chatprotocol.Event, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		closeErr := f.Close()
		if closeErr != nil {
			log.Printf("hub: error closing session file %s: %v", path, closeErr)
		}
	}()

	events := make([]chatprotocol.Event, 0)

	reader := bufio.NewReader(f)
	lineNum := 0

	for {
		line, readErr := reader.ReadString('\n')

		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			lineNum++

			var event chatprotocol.Event

			unmarshalErr := json.Unmarshal([]byte(trimmed), &event)
			if unmarshalErr != nil {
				log.Printf("hub: skipping malformed history line %d in %s: %v", lineNum, path, unmarshalErr)
			} else {
				events = append(events, event)
			}
		}

		if readErr != nil {
			if readErr != io.EOF {
				log.Printf("hub: error reading session file %s: %v", path, readErr)
			}

			break
		}
	}

	return events, nil
}

// sessionSummary is one row of the /history listing — see Hub.serveHistoryList.
type sessionSummary struct {
	ID               string `json:"id"`
	LastActivityAt   string `json:"lastActivityAt"`
	FirstUserMessage string `json:"firstUserMessage"`
}

// listSessions returns one summary per *.jsonl file directly inside dir. A file that can't be
// stat'd, read, or parsed still gets a summary (with an empty FirstUserMessage) rather than being
// dropped or failing the whole list — see summarizeSession.
func listSessions(dir string) ([]sessionSummary, error) {
	matches, err := filepath.Glob(filepath.Join(dir, "*"+historySessionExt))
	if err != nil {
		return nil, err
	}

	summaries := make([]sessionSummary, 0, len(matches))

	for _, match := range matches {
		summaries = append(summaries, summarizeSession(match))
	}

	return summaries, nil
}

// summarizeSession builds one sessionSummary for the file at path.
func summarizeSession(path string) sessionSummary {
	id := strings.TrimSuffix(filepath.Base(path), historySessionExt)

	info, err := os.Stat(path)
	if err != nil {
		log.Printf("hub: error stating session file %s: %v", path, err)

		return sessionSummary{ID: id}
	}

	events, err := readSessionFile(path)
	if err != nil {
		log.Printf("hub: error reading session file %s: %v", path, err)
	}

	summary := sessionSummary{
		ID:               id,
		LastActivityAt:   info.ModTime().Format(time.RFC3339),
		FirstUserMessage: firstUserMessagePreview(events),
	}

	return summary
}

// firstUserMessagePreview returns the first user_message-typed event's Text, truncated to
// firstUserMessagePreviewRunes runes for a list-row title — empty if events holds none.
func firstUserMessagePreview(events []chatprotocol.Event) string {
	for _, event := range events {
		if event.Type == chatprotocol.EventUserMessage {
			return truncateRunes(event.Text, firstUserMessagePreviewRunes)
		}
	}

	return ""
}

// truncateRunes returns s truncated to at most n runes (not bytes, so a multi-byte character
// is never split).
func truncateRunes(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}

	return string(runes[:n])
}

// isSafeSessionID reports whether id is safe to join into a path under a Hub's historyDir: no
// path separators and no ".." traversal segments.
func isSafeSessionID(id string) bool {
	if id == "" {
		return false
	}

	if strings.ContainsAny(id, `/\`) {
		return false
	}

	return !strings.Contains(id, "..")
}

// serveHistoryList answers GET /history with one summary per session file in h.historyDir — see
// listSessions.
func (h *Hub) serveHistoryList(w http.ResponseWriter) {
	summaries, err := listSessions(h.historyDir)
	if err != nil {
		log.Printf("hub: error listing history sessions in %s: %v", h.historyDir, err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	writeJSON(w, summaries)
}

// serveHistoryDetail answers GET /history/<id> with the full ordered event array parsed from
// <h.historyDir>/<id>.jsonl. Responds 400 if id isn't a safe bare filename (see isSafeSessionID)
// and 404 if no such session file exists.
func (h *Hub) serveHistoryDetail(w http.ResponseWriter, id string) {
	if !isSafeSessionID(id) {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	path := filepath.Join(h.historyDir, id+historySessionExt)

	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)

			return
		}

		log.Printf("hub: error stating history session %s: %v", path, err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	events, err := readSessionFile(path)
	if err != nil {
		log.Printf("hub: error reading history session %s: %v", path, err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	writeJSON(w, events)
}

// writeJSON writes v as an application/json response body.
func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		log.Printf("hub: error writing JSON response: %v", err)
	}
}
