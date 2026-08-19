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
package authlogin

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"

	"github.com/creack/pty"

	"workbenchbridge/internal/chatprotocol"
	"workbenchbridge/internal/envdrop"
)

// readChunkSize is the buffer size for each pty Read. setup-token's screen redraws are small
// (a spinner frame, a line of text); there is no benefit to a larger buffer and every Read is
// copied into the Parser's growing transcript anyway.
const readChunkSize = 4096

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
	cmd := exec.CommandContext(ctx, r.binary, "setup-token")

	ptmx, err := pty.Start(cmd)
	if err != nil {
		return fmt.Errorf("error starting claude setup-token under a pty: %w", err)
	}
	defer ptmx.Close()

	output, readErrCh := readLoop(ptmx)

	parser := NewParser()

	for {
		select {
		case chunk, ok := <-output:
			if !ok {
				return r.finish(cmd, <-readErrCh)
			}

			for _, signal := range parser.Feed(chunk) {
				done, err := r.handleSignal(signal, ptmx, envStore, emit)
				if done {
					return err
				}
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
// ever producing a token (a token exit returns from the main loop directly via handleSignal). It
// reports whatever went wrong, preferring the process's own exit error over a plain EOF.
func (r *Runner) finish(cmd *exec.Cmd, readErr error) error {
	waitErr := cmd.Wait()
	if waitErr != nil {
		return fmt.Errorf("claude setup-token exited without producing a token: %w", waitErr)
	}

	if readErr != nil && !errors.Is(readErr, io.EOF) {
		return fmt.Errorf("error reading claude setup-token output: %w", readErr)
	}

	return errors.New("claude setup-token exited without producing a token")
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
