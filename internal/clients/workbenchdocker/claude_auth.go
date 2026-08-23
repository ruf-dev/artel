package workbenchdocker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types/container"

	"go.redsock.ru/rerrors"
)

// claudeCredentialsCheckScript resolves claude's own credentials file path the same way
// deploy/workbench/bridge/internal/authlogin's credentialsPath() does —
// $CLAUDE_CONFIG_DIR/.credentials.json when that env var is set inside the container, else
// $HOME/.claude/.credentials.json — and tests for its existence. `test -f` needs no captured
// output: its exit code alone (0 = exists, non-zero = doesn't) is the whole answer.
const claudeCredentialsCheckScript = `test -f "${CLAUDE_CONFIG_DIR:-$HOME/.claude}/.credentials.json"`

// CheckClaudeLoggedIn reports whether containerID's claude CLI currently has a completed login —
// ground truth read straight from the container's filesystem (see claudeCredentialsCheckScript),
// unlike internal/transport/vaults_api/workbench_terminal_shell.go's terminal-output-based
// sign-in-link detection, which only ever observes what scrolled through a terminal WS relay and
// has no way to notice a login completing after the user stops watching.
//
// A "command ran, file doesn't exist" outcome (test's own non-zero exit) is reported as
// false, nil — not an error. Only a failure of the exec machinery itself (Create/Attach/Inspect)
// is returned as an error. Mirrors execTmuxCommand's Create/Attach/Inspect shape; unlike it, there
// is nothing caller-controlled in the command, so no positional-arg substitution is needed.
func (c *Client) CheckClaudeLoggedIn(ctx context.Context, containerID string) (bool, error) {
	cmd := []string{"/bin/sh", "-c", claudeCredentialsCheckScript}
	execOptions := container.ExecOptions{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
	}

	created, err := c.cli.ContainerExecCreate(ctx, containerID, execOptions)
	if err != nil {
		return false, rerrors.Wrap(err, "creating claude login check exec")
	}

	execAttachOptions := container.ExecAttachOptions{}

	resp, err := c.cli.ContainerExecAttach(ctx, created.ID, execAttachOptions)
	if err != nil {
		return false, rerrors.Wrap(err, "attaching to claude login check exec")
	}
	defer resp.Close()

	_, err = io.Copy(io.Discard, resp.Reader)
	if err != nil {
		return false, rerrors.Wrap(err, "reading claude login check exec output")
	}

	inspect, err := c.cli.ContainerExecInspect(ctx, created.ID)
	if err != nil {
		return false, rerrors.Wrap(err, "inspecting claude login check exec")
	}

	return inspect.ExitCode == 0, nil
}
