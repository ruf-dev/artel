package workbenchdocker

import (
	"context"
	"net"
	"strconv"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"

	"go.redsock.ru/rerrors"
)

// tmuxSessionName is the name of the tmux session the workbench image's entrypoint (see
// deploy/workbench/entrypoint.sh) creates/attaches to, and that StartContainer's environment
// injection targets. Shared as a constant rather than duplicated so the two stay in sync.
const tmuxSessionName = "workbench"

// ttydPort is the port ttyd binds to inside every workbench container — must match the port
// passed to `ttyd -p` in deploy/workbench/entrypoint.sh.
//
// Published to a host port in [ttydHostPortRangeStart, ttydHostPortRangeEnd] (see CreateContainer)
// rather than left unpublished: the configured daemon (docker_hosts.url) is frequently itself a
// docker:dind *container* (nested Docker), whose inner bridge network (workbenchNetworkName
// included) lives in a network namespace private to that container — a workbench container's IP
// on it is only ever reachable from processes inside the dind container's own netns, never from
// Artel's process, regardless of which host or execution mode (bare binary, containerized) Artel
// runs as. Publishing the port routes through the dind container's own (reachable) address
// instead — see ContainerAddress.
const ttydPort = 7681

// ttydNatPort is ttydPort in github.com/docker/go-connections/nat's "<port>/<proto>" key format,
// used to declare/read back the port binding in CreateContainer/ContainerAddress.
var ttydNatPort = nat.Port(strconv.Itoa(ttydPort) + "/tcp")

// ttydHostPortRangeStart/End bound the host ports CreateContainer publishes ttydPort to — a
// *fixed, pre-known* range rather than a Docker-assigned random one (HostPort ""), because a
// random port is only reachable from wherever already has direct network-level access to
// whichever container/process actually bound it (a docker:dind daemon's own netns). A fixed range
// can additionally be published, once, on the outer dind container itself (e.g. `-p
// 20000-20099:20000-20099` — see tests/docker-compose.yaml's test-dockerd and
// docs/tasks/workbench/02_docker_topology.md), which is what makes it reachable uniformly: from a
// bare Artel process on the same host (Docker's own host-port-forwarding, the same mechanism that
// already makes docker_hosts.url's own daemon port reachable), from Artel running as a sibling
// container on the daemon's network, and from a real (non-dind) remote daemon alike. 100 slots is
// a hardcoded prototype limit, same spirit as workbenchCpuLimitNanoCpus/workbenchMemLimitBytes —
// revisit if concurrent workbench count ever approaches it.
const (
	ttydHostPortRangeStart = 20000
	ttydHostPortRangeEnd   = 20099
)

// CreateOpts configures a new workbench container. Deliberately excludes any secret/env
// value — per docs/workbench/01_data_model_and_lifecycle.md, a workbench container is created
// with no auth env vars at all; those are only decided and supplied later, at StartContainer
// time.
type CreateOpts struct {
	// Name is the container name (e.g. "workbench-<vault_id>-<user_id>").
	Name string
	// VolumeName is the pre-created named volume to mount at workspaceMountPath.
	VolumeName string
}

// CreateContainer creates (but does not start) a workbench container: the hardcoded workbench
// image, attached to the dedicated workbench-net network (assumed to pre-exist on the
// configured daemon — see docs/workbench/02_docker_topology.md), with opts.VolumeName mounted
// at workspaceMountPath, hardcoded CPU/memory limits, and labeled for operational visibility.
//
// ttydPort is published to a host port allocated from [ttydHostPortRangeStart,
// ttydHostPortRangeEnd] (see allocateTtydHostPort) rather than left unpublished. It is not
// exposed to the wider internet by this alone: the configured daemon is expected to sit on a
// network only Artel's own process can reach (see docs/workbench/02_docker_topology.md's "Network
// isolation"), and ttyd itself is only ever dialed through internal/transport/vaults_api's
// authenticated reverse proxy, never linked directly to a client.
func (c *Client) CreateContainer(ctx context.Context, opts CreateOpts) (string, error) {
	tag, err := c.EnsureImage(ctx)
	if err != nil {
		return "", rerrors.Wrap(err, "error ensuring workbench image")
	}

	hostPort, err := c.allocateTtydHostPort(ctx)
	if err != nil {
		return "", rerrors.Wrap(err, "allocating ttyd host port")
	}

	labels := map[string]string{
		workbenchLabelKey: workbenchLabelValue,
	}

	exposedPorts := nat.PortSet{
		ttydNatPort: struct{}{},
	}

	containerConfig := &container.Config{
		Image:        tag,
		Labels:       labels,
		ExposedPorts: exposedPorts,
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

	ttydPortBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: hostPort,
	}
	portBindings := nat.PortMap{
		ttydNatPort: []nat.PortBinding{ttydPortBinding},
	}

	hostConfig := &container.HostConfig{
		Mounts:       []mount.Mount{workspaceMount},
		Resources:    resources,
		PortBindings: portBindings,
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

// allocateTtydHostPort picks a free host port in [ttydHostPortRangeStart, ttydHostPortRangeEnd]
// for a new workbench container's ttyd binding, by inspecting every existing artel.workbench
// container (running or not — HostConfig.PortBindings reflects a container's declared binding
// regardless of its current state) and returning the first port in range none of them already
// claim.
//
// This is a best-effort reservation, not a lock: two concurrent CreateContainer calls could both
// observe the same free port and both pick it, in which case the daemon rejects whichever one
// starts second with a "port is already allocated" error — an acceptable failure mode for the
// prototype's expected (low, per-admin) concurrency, same tradeoff as the hardcoded resource
// limits elsewhere in this package.
func (c *Client) allocateTtydHostPort(ctx context.Context) (string, error) {
	labelFilter := filters.NewArgs(filters.Arg("label", workbenchLabelKey+"="+workbenchLabelValue))
	listOptions := container.ListOptions{
		All:     true,
		Filters: labelFilter,
	}

	existing, err := c.cli.ContainerList(ctx, listOptions)
	if err != nil {
		return "", rerrors.Wrap(err, "listing existing workbench containers")
	}

	usedPorts := make(map[string]struct{}, len(existing))

	for _, summary := range existing {
		inspect, err := c.cli.ContainerInspect(ctx, summary.ID)
		if err != nil {
			return "", rerrors.Wrap(err, "inspecting existing workbench container")
		}

		if inspect.HostConfig == nil {
			continue
		}

		for _, binding := range inspect.HostConfig.PortBindings[ttydNatPort] {
			usedPorts[binding.HostPort] = struct{}{}
		}
	}

	for port := ttydHostPortRangeStart; port <= ttydHostPortRangeEnd; port++ {
		portStr := strconv.Itoa(port)

		_, taken := usedPorts[portStr]
		if !taken {
			return portStr, nil
		}
	}

	return "", rerrors.New("no free ttyd host port in range")
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

// ContainerAddress returns containerID's "<daemon-host>:<published-port>" address for the ttyd
// server running inside it, as reachable from Artel's own process. Used by the terminal
// reverse-proxy handler in internal/transport/vaults_api.
//
// This deliberately does not resolve the container's IP on workbenchNetworkName: the configured
// daemon (c.host) is commonly itself a docker:dind *container*, whose inner networks — including
// workbenchNetworkName — live in a network namespace private to that container, so a workbench
// container's IP on it is never routable from outside the dind container regardless of where
// Artel's own process runs. Reading back the Docker-assigned host port CreateContainer published
// ttydPort to and pairing it with c.host's own hostname instead routes through the dind
// container's (or, for a bare/second-dockerd host, the host's) already-reachable address — the
// same address Artel used to create/start the container in the first place.
//
// Fails rather than guessing when the daemon hasn't assigned a host port yet (i.e. the container
// isn't running) — an empty/absent address would otherwise surface much later as an opaque proxy
// dial failure.
func (c *Client) ContainerAddress(ctx context.Context, containerID string) (string, error) {
	inspect, err := c.cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return "", rerrors.Wrap(err, "inspecting workbench container")
	}

	if inspect.NetworkSettings == nil {
		return "", rerrors.New("workbench container has no network settings: " + containerID)
	}

	bindings, ok := inspect.NetworkSettings.Ports[ttydNatPort]
	if !ok || len(bindings) == 0 {
		return "", rerrors.New("workbench container has no published ttyd port: " + containerID)
	}

	hostPort := bindings[0].HostPort
	if hostPort == "" {
		return "", rerrors.New("workbench container's ttyd port has no assigned host port: " + containerID)
	}

	daemonHost, err := c.daemonHost()
	if err != nil {
		return "", rerrors.Wrap(err, "resolving docker daemon host")
	}

	return net.JoinHostPort(daemonHost, hostPort), nil
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
