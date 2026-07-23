package workbenchdocker

import (
	"bytes"
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/pkg/stdcopy"

	"go.redsock.ru/rerrors"
)

// tmuxSessionName is the name of the tmux session the workbench image's entrypoint (see
// deploy/workbench/entrypoint.sh) creates/attaches to, and that StartContainer's environment
// injection targets. Shared as a constant rather than duplicated so the two stay in sync.
const tmuxSessionName = "workbench"

// CreateOpts configures a new workbench container. Deliberately excludes any secret/env
// value — per docs/workbench/01_data_model_and_lifecycle.md, a workbench container is created
// with no auth env vars at all; those are only decided and supplied later, at StartContainer
// time.
type CreateOpts struct {
	// Name is the container name (e.g. "workbench-<vault_id>").
	Name string
	// VolumeName is the pre-created named volume to mount at workspaceMountPath.
	VolumeName string
}

// CreateContainer creates (but does not start) a workbench container: the hardcoded workbench
// image, attached to the dedicated workbench-net network (assumed to pre-exist on the
// configured daemon — see docs/workbench/02_docker_topology.md), with opts.VolumeName mounted
// at workspaceMountPath, hardcoded CPU/memory limits, no exposed ports, and labeled for
// operational visibility.
func (c *Client) CreateContainer(ctx context.Context, opts CreateOpts) (string, error) {
	labels := map[string]string{
		workbenchLabelKey: workbenchLabelValue,
	}

	containerConfig := &container.Config{
		Image:  workbenchImage,
		Labels: labels,
	}

	workspaceMount := mount.Mount{
		Type:   mount.TypeVolume,
		Source: opts.VolumeName,
		Target: workspaceMountPath,
	}

	resources := container.Resources{
		NanoCPUs: workbenchCpuLimitNanoCpus,
		Memory:   workbenchMemLimitBytes,
	}

	hostConfig := &container.HostConfig{
		Mounts:    []mount.Mount{workspaceMount},
		Resources: resources,
	}

	workbenchEndpoint := &network.EndpointSettings{}
	endpointsConfig := map[string]*network.EndpointSettings{
		workbenchNetworkName: workbenchEndpoint,
	}
	networkingConfig := &network.NetworkingConfig{
		EndpointsConfig: endpointsConfig,
	}

	resp, err := c.cli.ContainerCreate(ctx, containerConfig, hostConfig, networkingConfig, nil, opts.Name)
	if err != nil {
		return "", rerrors.Wrap(err, "creating workbench container")
	}

	return resp.ID, nil
}

// StartContainer starts an already-created workbench container and injects env into it.
//
// The Docker Engine API has no way to attach environment variables to a container at `docker
// start` time — env is otherwise only settable at `docker create` time, which is deliberately
// too early here (see CreateOpts). Once started, env is instead propagated into the running
// container's tmux session (tmuxSessionName) via a docker exec, so the value only ever lives in
// the tmux server's in-memory session state — never baked into the image, never written to the
// mounted volume. See docs/workbench/03_auth_and_login_flow.md for the full login-flow design
// this is feeding into (unconfirmed/spike-flagged there as of this writing); this is the
// mechanism this client offers to keep that property intact.
func (c *Client) StartContainer(ctx context.Context, containerID string, env map[string]string) error {
	startOptions := container.StartOptions{}

	err := c.cli.ContainerStart(ctx, containerID, startOptions)
	if err != nil {
		return rerrors.Wrap(err, "starting workbench container")
	}

	err = c.injectEnv(ctx, containerID, env)
	if err != nil {
		return rerrors.Wrap(err, "injecting env into workbench container")
	}

	return nil
}

// injectEnv sets each entry of env into the container's tmux session environment
// (tmuxSessionName), one docker exec per variable. Each exec's own process environment
// (ExecOptions.Env) carries the value in; only the variable's name ever appears in argv, never
// the value, keeping it out of `docker exec`/process-list logs.
func (c *Client) injectEnv(ctx context.Context, containerID string, env map[string]string) error {
	for key, value := range env {
		execOptions := container.ExecOptions{
			Cmd: []string{"/bin/sh", "-c", `exec tmux set-environment -t "$1" "$2" "$VALUE"`, "_", tmuxSessionName, key},
			Env: []string{"VALUE=" + value},
		}

		created, err := c.cli.ContainerExecCreate(ctx, containerID, execOptions)
		if err != nil {
			return rerrors.Wrap(err, "creating exec for env var "+key)
		}

		execStartOptions := container.ExecStartOptions{}

		err = c.cli.ContainerExecStart(ctx, created.ID, execStartOptions)
		if err != nil {
			return rerrors.Wrap(err, "starting exec for env var "+key)
		}
	}

	return nil
}

// CapturePane returns the current visible contents of the workbench's tmux pane
// (tmuxSessionName), via `tmux capture-pane -p`. Used to observe the subscription_login TUI
// flow (login URL, OAuth errors, ...) without disturbing it — see
// docs/workbench/03_auth_and_login_flow.md, "Mechanism (confirmed)".
func (c *Client) CapturePane(ctx context.Context, containerID string) (string, error) {
	execOptions := container.ExecOptions{
		Cmd:          []string{"tmux", "capture-pane", "-t", tmuxSessionName, "-p"},
		AttachStdout: true,
		AttachStderr: true,
	}

	created, err := c.cli.ContainerExecCreate(ctx, containerID, execOptions)
	if err != nil {
		return "", rerrors.Wrap(err, "creating exec for capture-pane")
	}

	attachOptions := container.ExecAttachOptions{}

	hijacked, err := c.cli.ContainerExecAttach(ctx, created.ID, attachOptions)
	if err != nil {
		return "", rerrors.Wrap(err, "attaching exec for capture-pane")
	}
	defer hijacked.Close()

	var stdout, stderr bytes.Buffer

	_, err = stdcopy.StdCopy(&stdout, &stderr, hijacked.Reader)
	if err != nil {
		return "", rerrors.Wrap(err, "reading capture-pane output")
	}

	return stdout.String(), nil
}

// SendKeys relays keys into the workbench's tmux pane (tmuxSessionName) via
// `tmux send-keys ... <keys> Enter`, followed by Enter as a separate key. keys is arbitrary,
// user-controlled input (an OAuth code or first-run keystrokes the user types/pastes) — it is
// passed as a single, literal argv element of ExecOptions.Cmd, never through a shell
// (`/bin/sh -c ...`), so there is no shell-metacharacter injection surface: the exec'd process is
// `tmux` itself, invoked directly, and keys is one opaque argument to it, not text that gets
// re-parsed as shell syntax. See docs/workbench/03_auth_and_login_flow.md, "Mechanism
// (confirmed)", step 4.
func (c *Client) SendKeys(ctx context.Context, containerID string, keys string) error {
	execOptions := container.ExecOptions{
		Cmd: []string{"tmux", "send-keys", "-t", tmuxSessionName, keys, "Enter"},
	}

	created, err := c.cli.ContainerExecCreate(ctx, containerID, execOptions)
	if err != nil {
		return rerrors.Wrap(err, "creating exec for send-keys")
	}

	execStartOptions := container.ExecStartOptions{}

	err = c.cli.ContainerExecStart(ctx, created.ID, execStartOptions)
	if err != nil {
		return rerrors.Wrap(err, "starting exec for send-keys")
	}

	return nil
}

// StopContainer gracefully stops a running workbench container, leaving it (and its volume)
// intact for a later restart.
func (c *Client) StopContainer(ctx context.Context, containerID string) error {
	stopOptions := container.StopOptions{}

	err := c.cli.ContainerStop(ctx, containerID, stopOptions)
	if err != nil {
		return rerrors.Wrap(err, "stopping workbench container")
	}

	return nil
}

// RemoveContainer force-removes a workbench container. It does not remove the container's
// volume — call RemoveVolume separately.
func (c *Client) RemoveContainer(ctx context.Context, containerID string) error {
	removeOptions := container.RemoveOptions{
		Force: true,
	}

	err := c.cli.ContainerRemove(ctx, containerID, removeOptions)
	if err != nil {
		return rerrors.Wrap(err, "removing workbench container")
	}

	return nil
}
