package workbenchdocker

import (
	"bytes"
	"context"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"

	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/domain"
)

// tmuxWindowListFormat asks tmux to print one line per window as
// "#{window_id}\t#{window_active}\t#{window_name}". A tab is used as the field separator rather
// than the more obvious ":" or " " because a tmux window name (freeform, user/`claude`-title-set)
// can legally contain either of those but can never contain a literal tab — so this is the only
// separator that can't collide with real window-name content and be misparsed.
const tmuxWindowListFormat = "#{window_id}\t#{window_active}\t#{window_name}"

// ListTmuxWindows lists every tmux window (surfaced to the browser as a "terminal tab") in
// containerID's workbench session, in tmux's own window order.
func (c *Client) ListTmuxWindows(ctx context.Context, containerID string) ([]domain.TerminalTab, error) {
	script := `exec tmux list-windows -t "$1" -F "` + tmuxWindowListFormat + `"`

	stdout, err := c.execTmuxCommand(ctx, containerID, script, []string{})
	if err != nil {
		return nil, rerrors.Wrap(err, "listing tmux windows")
	}

	return parseTmuxWindowList(stdout), nil
}

// parseTmuxWindowList parses ListTmuxWindows' captured stdout (one "#{window_id}\t#{window_active}\t
// #{window_name}" line per tmux window) into []domain.TerminalTab. Kept as a separate, pure,
// I/O-free function — deliberately exported for testing (lowercase but unit-testable from this
// package's own _test.go file) — so the parsing logic can be exercised without a Docker daemon.
func parseTmuxWindowList(output string) []domain.TerminalTab {
	lines := strings.Split(output, "\n")

	tabs := make([]domain.TerminalTab, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}

		fields := strings.Split(line, "\t")
		if len(fields) != 3 {
			continue
		}

		tab := domain.TerminalTab{
			ID:     fields[0],
			Active: fields[1] == "1",
			Name:   fields[2],
		}

		tabs = append(tabs, tab)
	}

	return tabs
}

// NewTmuxWindow creates a new tmux window in containerID's workbench session, running `claude` as
// its command — the same command the entrypoint's initial `tmux new-session` starts (see
// deploy/workbench/entrypoint.sh) — and returns it as an active domain.TerminalTab.
//
// `-P -F '#{window_id}'` makes `tmux new-window` print just the new window's id to stdout, so the
// new tab's id is known without a second ListTmuxWindows round trip. Name is deliberately left
// empty: tmux's automatic-rename (on by default — see deploy/workbench/tmux.conf) only settles
// the window's name once `claude` starts printing its title escape sequence, which hasn't
// happened yet by the time this call returns. A caller that needs the settled name should do a
// fresh ListTmuxWindows shortly after.
func (c *Client) NewTmuxWindow(ctx context.Context, containerID string) (domain.TerminalTab, error) {
	script := `exec tmux new-window -t "$1" -P -F "#{window_id}" claude`

	stdout, err := c.execTmuxCommand(ctx, containerID, script, []string{})
	if err != nil {
		return domain.TerminalTab{}, rerrors.Wrap(err, "creating tmux window")
	}

	windowID := strings.TrimSpace(stdout)

	tab := domain.TerminalTab{
		ID:     windowID,
		Active: true,
	}

	return tab, nil
}

// SelectTmuxWindow makes windowID the "current" window of containerID's workbench session — since
// ttyd just attaches to the session and mirrors whatever window is currently current, this alone
// is what makes a tab switch reach the browser. tmux window ids ("@N") are unique across the whole
// tmux server, not just within one session, so no session-qualified target is needed.
//
// windowID is passed via the exec's Env ($WINDOW_ID), never string-concatenated into Cmd or the
// script text, even though every caller today only ever passes a server-issued "@N" id — same
// defense-in-depth this package's injectEnv already applies to caller-influenced values, keeping
// them out of argv/process-list logs.
func (c *Client) SelectTmuxWindow(ctx context.Context, containerID, windowID string) error {
	script := `exec tmux select-window -t "$WINDOW_ID"`
	env := []string{"WINDOW_ID=" + windowID}

	_, err := c.execTmuxCommand(ctx, containerID, script, env)
	if err != nil {
		return rerrors.Wrap(err, "selecting tmux window")
	}

	return nil
}

// KillTmuxWindow closes windowID in containerID's workbench session (see SelectTmuxWindow for why
// windowID alone, with no session qualifier, is a valid target, and why it's passed via Env rather
// than Cmd). It deliberately does not guard against killing a session's last window — tmux itself
// would refuse or tear down the whole session depending on version/config, and deciding whether
// that's allowed needs a window count (via ListTmuxWindows) first, which is a service-layer
// concern (see workbench.Service.CloseTerminalTab), not this thin client's.
//
// A vanished windowID (e.g. the user closed it manually from inside the terminal, Ctrl-b w) is a
// real, expected scenario, not a bug — execTmuxCommand surfaces tmux's own "can't find window"
// stderr in the returned error so a caller can tell the two apart from a generic failure.
func (c *Client) KillTmuxWindow(ctx context.Context, containerID, windowID string) error {
	script := `exec tmux kill-window -t "$WINDOW_ID"`
	env := []string{"WINDOW_ID=" + windowID}

	_, err := c.execTmuxCommand(ctx, containerID, script, env)
	if err != nil {
		return rerrors.Wrap(err, "killing tmux window")
	}

	return nil
}

// execTmuxCommand runs script (a static "/bin/sh -c" body — the "$1"/tmuxSessionName positional
// substitution keeps every caller-influenced value out of the shell string itself, same rule as
// injectEnv) inside containerID via docker exec, and returns its captured stdout.
//
// Unlike injectEnv (fire-and-forget ContainerExecStart, output and exit code both ignored), the
// four TerminalTab methods above all need the command's actual output and need to detect tmux
// failures — e.g. select-window/kill-window against a window id that's since vanished, a real
// scenario since a user can also open/close tmux windows manually from inside the terminal itself
// (Ctrl-b c / Ctrl-b w), outside this API's knowledge. So this uses ContainerExecAttach (not
// ExecStart) to get a readable stream, demultiplexes it with stdcopy (Docker frames
// stdout/stderr together with an 8-byte header per chunk in non-tty exec mode — the default
// ExecAttachOptions{} used here), and inspects the exec's exit code afterward, wrapping a non-zero
// exit with the captured stderr so a caller sees e.g. "tmux exited with status 1: can't find
// window: @9" instead of an opaque failure.
func (c *Client) execTmuxCommand(ctx context.Context, containerID string, script string, env []string) (string, error) {
	cmd := []string{"/bin/sh", "-c", script, "_", tmuxSessionName}
	execOptions := container.ExecOptions{
		Cmd:          cmd,
		Env:          env,
		AttachStdout: true,
		AttachStderr: true,
	}

	created, err := c.cli.ContainerExecCreate(ctx, containerID, execOptions)
	if err != nil {
		return "", rerrors.Wrap(err, "creating tmux exec")
	}

	execAttachOptions := container.ExecAttachOptions{}

	resp, err := c.cli.ContainerExecAttach(ctx, created.ID, execAttachOptions)
	if err != nil {
		return "", rerrors.Wrap(err, "attaching to tmux exec")
	}
	defer resp.Close()

	var stdoutBuf, stderrBuf bytes.Buffer

	_, err = stdcopy.StdCopy(&stdoutBuf, &stderrBuf, resp.Reader)
	if err != nil {
		return "", rerrors.Wrap(err, "reading tmux exec output")
	}

	inspect, err := c.cli.ContainerExecInspect(ctx, created.ID)
	if err != nil {
		return "", rerrors.Wrap(err, "inspecting tmux exec")
	}

	if inspect.ExitCode != 0 {
		stderr := strings.TrimSpace(stderrBuf.String())

		return "", rerrors.New("tmux exited with status " + strconv.Itoa(inspect.ExitCode) + ": " + stderr)
	}

	return stdoutBuf.String(), nil
}
