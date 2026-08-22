// Package authlogin drives `claude setup-token`, the CLI's flow for minting a long-lived (1-year)
// OAuth token from a Claude subscription — the credential the bridge uses for
// WORKBENCH_AUTH_MODE=subscription_login (see internal/service/v1/workbench/workbench.go in the
// main module).
//
// # Why a pty
//
// `claude setup-token` was tested both ways against real CLI 2.1.234 before writing this package:
//
//   - Under a plain os/exec with piped stdin/stdout, it produces no usable output at all — a
//     single newline byte and then nothing, indefinitely. It never even prints the sign-in URL.
//   - Under a pty (this package uses github.com/creack/pty), it renders exactly as it does for an
//     interactive user: full ANSI cursor positioning, truecolor text, and OSC 8 hyperlinks
//     carrying the sign-in URL.
//
// So a pty is not an enhancement here, it is required for the flow to produce any output to
// parse. No other pty library was evaluated.
//
// # What was verified vs. assumed
//
// The sign-in link (an OSC 8 hyperlink) and the invalid-code error message were captured from a
// real run and are parsed by internal patterns confirmed against that capture — see parse.go and
// its testdata fixture. The final success screen (the one-year token itself) was *not* observed:
// completing that requires an actual Claude subscription and a real browser login, out of scope
// for this change. Token extraction (see tokenPattern in parse.go) is a best-effort regex over
// the well-known "sk-ant-" credential prefix and should be re-verified against real output before
// being relied on in production.
//
// # Credentials-file fallback
//
// Because tokenPattern was never confirmed against a real success screen, Run does not rely on it
// alone: `claude setup-token` (like `claude login`) persists its result as a side effect to
// $CLAUDE_CONFIG_DIR/.credentials.json (default $HOME/.claude/.credentials.json) — a
// {"claudeAiOauth":{"accessToken":"sk-ant-...",...}} document, read by every later `claude`
// invocation (headless or interactive) regardless of env vars. This is the same file `claude
// login` writes, confirmed from the CLI's own bundled source. Run polls for that file appearing
// or changing after the process starts (see readFreshToken) as a second, independent path to a
// token — one that doesn't depend on recognizing the literal printed success text at all, only on
// the process having completed successfully. Whichever detection fires first (regex match on the
// pty transcript, or the credentials file showing up) wins; either way the token is written to
// envStore the same way, so downstream code (main.go's ensureAuthenticated gate, claudecli.Runner)
// doesn't need to know which path produced it.
package authlogin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/creack/pty"

	"workbenchbridge/internal/chatprotocol"
	"workbenchbridge/internal/envdrop"
)

// credentialsPollInterval bounds how quickly the credentials-file fallback (see package doc)
// notices `claude setup-token` finishing when the pty transcript itself never produces a
// recognized SignalToken. Short enough not to be perceptible against a human pasting a code.
const credentialsPollInterval = 500 * time.Millisecond

// readChunkSize is the buffer size for each pty Read. setup-token's screen redraws are small
// (a spinner frame, a line of text); there is no benefit to a larger buffer and every Read is
// copied into the Parser's growing transcript anyway.
const readChunkSize = 4096

// ptyWinsize is the window size given to the pty setup-token runs under. pty.Start (no size
// specified) leaves the kernel's default 0x0 winsize in place; setup-token's Ink-based TUI treats
// that as "can't detect width" and falls back to wrapping its own output at a defaulted 80
// columns — counting the invisible OSC 8 escape bytes toward that budget the same as visible
// characters, since it wraps the fully-composed line rather than the rendered text. That
// inserts a real line break partway through the sign-in URL's escape sequence, at whatever byte
// offset 80 columns happens to land on for that particular run's client_id/state values — not
// even at a consistent point, since those are per-run random values of varying rendered width.
// A wide, fixed window size removes the incentive to wrap this line at all.
var ptyWinsize = &pty.Winsize{Rows: 50, Cols: 2000}

// Runner drives one claude setup-token attempt.
type Runner struct {
	// binary is the `claude` executable to run; a field so tests could point it at a stub script
	// (parse.go's logic is what's actually unit-tested, since a stub can't reproduce a real pty
	// session's byte-for-byte output — see parse_test.go).
	binary string
}

// NewRunner returns a Runner. An empty binary falls back to "claude".
func NewRunner(binary string) *Runner {
	if binary == "" {
		binary = "claude"
	}

	return &Runner{binary: binary}
}

// Run executes one claude setup-token attempt to completion, emitting chatprotocol events as the
// flow progresses and consuming codes for any auth_code_submit the caller receives while it's
// running (the caller — main.go — owns reading the hub's inbound channel; Run only owns the pty).
//
// Run returns nil once the flow produces a token (already persisted into env via
// envStore.Write(envdrop.OauthTokenVar, ...) before returning), or an error if the process exits
// without one or ctx is cancelled first.
func (r *Runner) Run(ctx context.Context, envStore *envdrop.Store, emit func(chatprotocol.Event), codes <-chan string) error {
	startTime := time.Now()

	cmd := exec.CommandContext(ctx, r.binary, "setup-token")

	ptmx, err := pty.StartWithSize(cmd, ptyWinsize)
	if err != nil {
		return fmt.Errorf("error starting claude setup-token under a pty: %w", err)
	}
	defer ptmx.Close()

	output, readErrCh := readLoop(ptmx)

	parser := NewParser()

	ticker := time.NewTicker(credentialsPollInterval)
	defer ticker.Stop()

	for {
		select {
		case chunk, ok := <-output:
			if !ok {
				return r.finish(cmd, <-readErrCh, envStore, startTime)
			}

			for _, signal := range parser.Feed(chunk) {
				done, err := r.handleSignal(signal, ptmx, envStore, emit)
				if done {
					return err
				}
			}
		case <-ticker.C:
			done, err := r.checkCredentialsFile(cmd, envStore, startTime)
			if done {
				return err
			}
		case code := <-codes:
			_, err := ptmx.Write([]byte(code + "\r"))
			if err != nil {
				emit(errorEvent("error writing auth code to claude setup-token: " + err.Error()))
			}
		case <-ctx.Done():
			_ = cmd.Process.Kill()

			return ctx.Err()
		}
	}
}

// handleSignal reacts to one Parser signal, returning done=true once the flow has a final
// outcome (a token) that Run should return immediately for.
func (r *Runner) handleSignal(signal Signal, ptmx io.Writer, envStore *envdrop.Store, emit func(chatprotocol.Event)) (bool, error) {
	switch signal.Kind {
	case SignalAuthLink:
		event := chatprotocol.Event{Type: chatprotocol.EventAuthLink, URL: signal.URL}
		emit(event)
	case SignalCodeNeeded:
		event := chatprotocol.Event{Type: chatprotocol.EventAuthCodeNeeded}
		emit(event)
	case SignalError:
		emit(errorEvent(signal.Error))

		// The CLI's own error message says "Press Enter to retry" — advance past that prompt so
		// it re-shows the paste-code screen, then ask the caller for another code.
		_, err := ptmx.Write([]byte("\r"))
		if err != nil {
			emit(errorEvent("error advancing past claude setup-token's retry prompt: " + err.Error()))
		}

		event := chatprotocol.Event{Type: chatprotocol.EventAuthCodeNeeded}
		emit(event)
	case SignalToken:
		err := envStore.Write(envdrop.OauthTokenVar, signal.Token)
		if err != nil {
			emit(errorEvent("error persisting the OAuth token: " + err.Error()))

			return true, fmt.Errorf("error persisting the OAuth token: %w", err)
		}

		return true, nil
	}

	return false, nil
}

// finish is reached once the pty's read side hits EOF, i.e. claude setup-token has exited without
// its transcript ever producing a recognized SignalToken (a token found that way returns from the
// main loop directly via handleSignal). Before declaring failure, it makes one last check of the
// credentials file (see package doc's "Credentials-file fallback") in case the process wrote it
// and exited between the last tick and EOF.
func (r *Runner) finish(cmd *exec.Cmd, readErr error, envStore *envdrop.Store, startTime time.Time) error {
	waitErr := cmd.Wait()
	if waitErr != nil {
		return fmt.Errorf("claude setup-token exited without producing a token: %w", waitErr)
	}

	if readErr != nil && !errors.Is(readErr, io.EOF) {
		return fmt.Errorf("error reading claude setup-token output: %w", readErr)
	}

	token, found, err := readFreshToken(startTime)
	if err == nil && found {
		writeErr := envStore.Write(envdrop.OauthTokenVar, token)
		if writeErr != nil {
			return fmt.Errorf("error persisting the OAuth token: %w", writeErr)
		}

		return nil
	}

	return errors.New("claude setup-token exited without producing a token")
}

// checkCredentialsFile is the ticker-driven half of the credentials-file fallback (see package
// doc): it lets Run notice claude setup-token has already succeeded even while the process is
// still alive and its pty transcript never produced a recognized SignalToken — the case that
// otherwise hangs the flow forever, since neither a new signal nor EOF would ever arrive. Returns
// done=true once a fresh token is found, at which point the still-running process is no longer
// needed and is killed so its pty fully closes (unblocking readLoop's goroutine via Run's deferred
// ptmx.Close, which happens regardless of how Run returns).
func (r *Runner) checkCredentialsFile(cmd *exec.Cmd, envStore *envdrop.Store, startTime time.Time) (bool, error) {
	token, found, err := readFreshToken(startTime)
	if err != nil || !found {
		return false, nil
	}

	writeErr := envStore.Write(envdrop.OauthTokenVar, token)

	_ = cmd.Process.Kill()
	_, _ = cmd.Process.Wait()

	if writeErr != nil {
		return true, fmt.Errorf("error persisting the OAuth token: %w", writeErr)
	}

	return true, nil
}

// credentialsPath resolves where `claude` stores its own persistent credentials, matching the
// resolution order confirmed in the CLI's own bundled source: $CLAUDE_CONFIG_DIR/.credentials.json
// when that env var is set, else $HOME/.claude/.credentials.json.
func credentialsPath() string {
	if dir := os.Getenv("CLAUDE_CONFIG_DIR"); dir != "" {
		return filepath.Join(dir, ".credentials.json")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		home = "/root"
	}

	return filepath.Join(home, ".claude", ".credentials.json")
}

// credentialsFile is the subset of claude's own .credentials.json this package reads — see
// package doc's "Credentials-file fallback".
type credentialsFile struct {
	ClaudeAiOauth struct {
		AccessToken string `json:"accessToken"`
	} `json:"claudeAiOauth"`
}

// readFreshToken reports the access token in claude's credentials file, but only if that file's
// mtime is after since (i.e. it was written during this Run, not left over from some earlier
// login in this same $HOME) — a missing file, or one that predates since, is not an error, just a
// "not yet" (found=false).
func readFreshToken(since time.Time) (string, bool, error) {
	path := credentialsPath()

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", false, nil
		}

		return "", false, err
	}

	if !info.ModTime().After(since) {
		return "", false, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", false, err
	}

	var creds credentialsFile

	err = json.Unmarshal(data, &creds)
	if err != nil {
		return "", false, err
	}

	if creds.ClaudeAiOauth.AccessToken == "" {
		return "", false, nil
	}

	return creds.ClaudeAiOauth.AccessToken, true, nil
}

// readLoop copies ptmx into a channel of chunks so Run's select loop can multiplex pty output
// against inbound codes and ctx cancellation. The output channel is closed once Read returns any
// error (a pty's Read returns EOF, wrapped as a PathError, once the child's slave side closes —
// there is no cleaner signal available), with that error delivered on readErrCh for finish to
// inspect.
func readLoop(ptmx io.Reader) (<-chan []byte, <-chan error) {
	output := make(chan []byte)
	readErrCh := make(chan error, 1)

	go func() {
		defer close(output)

		buf := make([]byte, readChunkSize)

		for {
			n, err := ptmx.Read(buf)
			if n > 0 {
				chunk := make([]byte, n)
				copy(chunk, buf[:n])
				output <- chunk
			}

			if err != nil {
				readErrCh <- err

				return
			}
		}
	}()

	return output, readErrCh
}

func errorEvent(text string) chatprotocol.Event {
	return chatprotocol.Event{
		Type: chatprotocol.EventError,
		Text: text,
	}
}
