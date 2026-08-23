package claudecli

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"workbenchbridge/internal/chatprotocol"
	"workbenchbridge/internal/envdrop"
)

// maxLineBytes bounds a single stream-json line. Lines are ordinarily small, but a tool_result
// block echoing a large file read can be enormous, and bufio's 64KiB default would split it and
// hand the translator two halves of one JSON object.
const maxLineBytes = 32 * 1024 * 1024

// stderrTailLimit is how much of a failed `claude` process's stderr is reported to consumers —
// enough to carry the actual error, short of pasting a whole stack trace into the chat.
const stderrTailLimit = 4000

// Runner executes one `claude` turn per call and streams the resulting events.
//
// `claude --print` is single-shot: there is no long-lived multi-turn stdin mode. Multi-turn is
// therefore "a fresh process per turn, resumed by session id" — the first turn runs bare, every
// later one adds --resume <session id> read off the previous turn's result line, which is what
// SessionId tracks.
type Runner struct {
	// binary is the `claude` executable to run; a field so tests can point it at a stub script.
	binary string
	// workdir is the cwd every turn runs in — the workbench's mounted volume (see NewRunner's
	// default, /root/vault).
	workdir string
	// settingsPath registers the PreToolUse hook (see ../claudesettings) and is passed to every
	// invocation.
	settingsPath string

	env *envdrop.Store

	// mu serialises turns. Only one `claude` process may run at a time: two concurrent processes
	// resuming the same session would interleave writes into one transcript, and the permission
	// broker's correlation ids would no longer identify which turn a decision belongs to.
	mu sync.Mutex

	sessionMu sync.Mutex
	sessionId string
}

// NewRunner returns a Runner. An empty binary/workdir/settingsPath falls back to the container's
// defaults.
func NewRunner(binary, workdir, settingsPath string, env *envdrop.Store) *Runner {
	if binary == "" {
		binary = "claude"
	}

	if workdir == "" {
		workdir = "/root/vault"
	}

	return &Runner{
		binary:       binary,
		workdir:      workdir,
		settingsPath: settingsPath,
		env:          env,
	}
}

// SessionId returns the session the next turn will resume, empty before the first turn completes.
func (r *Runner) SessionId() string {
	r.sessionMu.Lock()
	defer r.sessionMu.Unlock()

	return r.sessionId
}

// SetSessionId overrides the session the next turn resumes.
func (r *Runner) SetSessionId(id string) {
	r.sessionMu.Lock()
	defer r.sessionMu.Unlock()

	r.sessionId = id
}

// RunTurn runs one chat turn for prompt, calling emit for every event as it is produced — the
// caller broadcasts them, so text reaches consumers token by token rather than at end of turn.
//
// It never returns an error for a failed turn: a crashed or non-zero `claude` is reported to
// consumers as an error event and swallowed here, because the bridge is PID 1 and one bad turn
// must not take the container (or anyone's WebSocket) down with it. The only thing that ends a
// turn early from the outside is ctx being cancelled.
func (r *Runner) RunTurn(ctx context.Context, prompt string, emit func(event chatprotocol.Event)) {
	r.mu.Lock()
	defer r.mu.Unlock()

	values, err := r.env.Read()
	if err != nil {
		emit(errorEvent("error reading injected environment: " + err.Error()))

		return
	}

	cmd := exec.CommandContext(ctx, r.binary, r.args(prompt)...)
	cmd.Dir = r.workdir
	cmd.Env = envdrop.Environ(os.Environ(), values)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		emit(errorEvent("error opening claude stdout: " + err.Error()))

		return
	}

	var stderr strings.Builder

	cmd.Stderr = &stderr

	err = cmd.Start()
	if err != nil {
		emit(errorEvent("error starting claude: " + err.Error()))

		return
	}

	r.stream(stdout, emit)

	err = cmd.Wait()
	if err != nil {
		// A cancelled context is the bridge shutting the turn down on purpose, not a failure
		// worth alarming consumers about.
		if ctx.Err() != nil {
			return
		}

		emit(errorEvent(describeFailure(err, stderr.String())))
	}
}

// args builds the argv for one turn.
//
// --output-format stream-json + --verbose + --include-partial-messages is the combination that
// yields token-level deltas; without --include-partial-messages only whole messages arrive, and
// stream-json rejects being used without --verbose. No --permission-mode is passed on purpose:
// every gated tool call must reach the PreToolUse hook and therefore a human, and a permissive
// mode would bypass it entirely.
func (r *Runner) args(prompt string) []string {
	args := []string{
		"-p", prompt,
		"--output-format", "stream-json",
		"--verbose",
		"--include-partial-messages",
	}

	if r.settingsPath != "" {
		args = append(args, "--settings", r.settingsPath)
	}

	sessionId := r.SessionId()
	if sessionId != "" {
		args = append(args, "--resume", sessionId)
	}

	return args
}

// stream reads stdout line by line, translating and emitting as it goes, and latches the session
// id off the turn_done event so the next turn can resume.
func (r *Runner) stream(stdout io.Reader, emit func(event chatprotocol.Event)) {
	translator := NewTranslator()

	scanner := bufio.NewScanner(stdout)
	scanner.Buffer(make([]byte, 0, 64*1024), maxLineBytes)

	for scanner.Scan() {
		for _, event := range translator.Translate(scanner.Bytes()) {
			if event.Type == chatprotocol.EventTurnDone && event.SessionID != "" {
				r.SetSessionId(event.SessionID)
			}

			emit(event)
		}
	}

	err := scanner.Err()
	if err != nil {
		log.Printf("claudecli: error reading claude output: %v", err)
		emit(errorEvent("error reading claude output: " + err.Error()))
	}
}

// describeFailure renders a non-zero `claude` exit for a consumer, appending whatever the process
// put on stderr since that is usually the only place the real cause appears.
func describeFailure(err error, stderr string) string {
	message := "claude exited with an error: " + err.Error()

	trimmed := strings.TrimSpace(stderr)
	if trimmed == "" {
		return message
	}

	if len(trimmed) > stderrTailLimit {
		trimmed = trimmed[len(trimmed)-stderrTailLimit:]
	}

	return fmt.Sprintf("%s\n%s", message, trimmed)
}

func errorEvent(text string) chatprotocol.Event {
	return chatprotocol.Event{
		Type: chatprotocol.EventError,
		Text: text,
	}
}
