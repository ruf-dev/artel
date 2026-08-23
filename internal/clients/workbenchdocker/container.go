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
// deploy/workbench/entrypoint.sh) creates/attaches to for the interactive terminal view (ttyd +
// per-tab tmux windows, see tmux_tabs.go). Shared as a constant rather than duplicated so the two
// stay in sync.
const tmuxSessionName = "workbench"

// envDropDir is the in-container directory StartContainer's environment injection writes each
// injected variable into, one file per variable (file name = variable name, file content = raw
// value, no trailing newline). The in-container chat bridge (deploy/workbench/bridge) reads this
// directory back before every `claude` invocation.
//
// A directory of files rather than a real process environment because the Docker Engine API has
// no way to attach env to an already-created container (see StartContainer), and rather than the
// tmux session environment it used to be, because env injection no longer targets tmux at all —
// tmux itself is still present in the image for the separate interactive-terminal feature (see
// tmuxSessionName), it's just no longer how secrets reach the `claude` invocations the chat
// bridge drives. It is mounted as a tmpfs by CreateContainer (see envDropTmpfsOptions), so
// injected secrets live only in the container's page cache — never in the image, never on the
// mounted workspace volume.
const envDropDir = "/run/workbench/env"

// envDropTmpfsOptions is the tmpfs mount option string CreateContainer declares envDropDir with:
// mode 0700 so only root (the only user the workbench image runs anything as) can traverse it,
// plus noexec/nosuid/nodev since it only ever holds secret material, never anything runnable.
const envDropTmpfsOptions = "rw,noexec,nosuid,nodev,mode=0700"

// A workbench container runs two independently-addressable in-container servers: the chat bridge
// (bridgePort) and, restored alongside the tmux-tab terminal view, ttyd (ttydPort) — each gets its
// own fixed host-port range and its own address-resolution method (ContainerAddress / TtydAddress)
// below.

// bridgePort is the port the in-container chat bridge's WebSocket/HTTP server binds to — must
// match the port the bridge listens on (deploy/workbench/bridge). Unchanged from the port ttyd
// used before the bridge replaced it, so the reverse proxy in internal/transport/vaults_api and
// everything that resolves its address needs no adjustment.
//
// Published to a host port in [bridgeHostPortRangeStart, bridgeHostPortRangeEnd] (see CreateContainer)
// rather than left unpublished: the configured daemon (docker_hosts.url) is frequently itself a
// docker:dind *container* (nested Docker), whose inner networks (including each workbench
// container's own dedicated network — see containerNetworkName) live in a network namespace
// private to that container — a workbench container's IP on any of them is only ever reachable
// from processes inside the dind container's own netns, never from Artel's process, regardless of
// which host or execution mode (bare binary, containerized) Artel runs as. Publishing the port
// routes through the dind container's own (reachable) address instead — see ContainerAddress.
const bridgePort = 7681

// bridgeNatPort is bridgePort in github.com/docker/go-connections/nat's "<port>/<proto>" key
// format, used to declare/read back the port binding in CreateContainer/ContainerAddress.
var bridgeNatPort = nat.Port(strconv.Itoa(bridgePort) + "/tcp")

// bridgeHostPortRangeStart/End bound the host ports CreateContainer publishes bridgePort to — a
// *fixed, pre-known* range rather than a Docker-assigned random one (HostPort ""), because a
// random port is only reachable from wherever already has direct network-level access to
// whichever container/process actually bound it (a docker:dind daemon's own netns). A fixed range
// can additionally be published, once, on the outer dind container itself (e.g. `-p
// 20000-20099:20000-20099` — see tests/docker-compose.yaml's test-dockerd), which is what makes
// it reachable uniformly: from a
// bare Artel process on the same host (Docker's own host-port-forwarding, the same mechanism that
// already makes docker_hosts.url's own daemon port reachable), from Artel running as a sibling
// container on the daemon's network, and from a real (non-dind) remote daemon alike. 100 slots is
// a hardcoded prototype limit, same spirit as workbenchCpuLimitNanoCpus/workbenchMemLimitBytes —
// revisit if concurrent workbench count ever approaches it.
const (
	bridgeHostPortRangeStart = 20000
	bridgeHostPortRangeEnd   = 20099
)

// ttydPort is the port the in-container ttyd server (the interactive tmux-tab terminal, restored
// alongside the chat bridge) binds to. Distinct from bridgePort so both servers can run in the
// same container at once.
const ttydPort = 7682

// ttydNatPort is ttydPort in nat's "<port>/<proto>" key format — same role as bridgeNatPort, for
// ttyd's own binding.
var ttydNatPort = nat.Port(strconv.Itoa(ttydPort) + "/tcp")

// ttydHostPortRangeStart/End bound the host ports CreateContainer publishes ttydPort to — a second
// fixed range, disjoint from [bridgeHostPortRangeStart, bridgeHostPortRangeEnd], for the same
// reasons documented on that range. Must also be published on the outer dind container itself
// (e.g. `-p 20100-20199:20100-20199`) alongside the bridge range.
const (
	ttydHostPortRangeStart = 20100
	ttydHostPortRangeEnd   = 20199
)

// CreateOpts configures a new workbench container. Deliberately excludes any secret/env value: a
// workbench container is created with no auth env vars at all; those are only decided and
// supplied later, at StartContainer time.
type CreateOpts struct {
	// Name is the container name (e.g. "workbench-<vault_id>-<user_id>"). Also the seed for the
	// dedicated per-container Docker network CreateContainer creates — see containerNetworkName.
	Name string
	// VolumeName is the pre-created named volume to mount at homeMountPath.
	VolumeName string
}

// containerNetworkName derives the dedicated, per-container Docker network name CreateContainer
// creates and attaches solely to this one container, from name (CreateOpts.Name — already unique
// per (vault, user), the same identifier the container and volume names share). One network per
// container, rather than every workbench container sharing workbenchNetworkName, so no two
// workbench containers can ever reach each other's ports.
func containerNetworkName(name string) string {
	return name + "-net"
}

// newWorkbenchHostConfig builds the HostConfig for a new workbench container: mounts, resource
// limits, port bindings, the secret-injection tmpfs, and the isolation hardening applied to every
// workbench container regardless of caller — dropping every Linux capability (CapDrop), blocking
// setuid/setgid/file-capability privilege escalation (SecurityOpt's no-new-privileges), and
// capping the number of PIDs the container's cgroup may create (PidsLimit) as a fork-bomb guard.
// Split out from CreateContainer so these fields can be asserted on directly in tests without a
// live Docker daemon.
func newWorkbenchHostConfig(volumeName string, portBindings nat.PortMap) *container.HostConfig {
	homeMount := mount.Mount{
		Type:   mount.TypeVolume,
		Source: volumeName,
		Target: homeMountPath,
	}

	resources := container.Resources{
		NanoCPUs:  workbenchCpuLimitNanoCpus,
		Memory:    workbenchMemLimitBytes,
		PidsLimit: &pidsLimit,
	}

	// envDropDir is a tmpfs, not part of the container's writable layer and not part of the
	// mounted volume, so the secrets injectEnv drops into it never outlive the container's own
	// runtime and never reach durable storage.
	tmpfs := map[string]string{
		envDropDir: envDropTmpfsOptions,
	}

	capDrop := []string{"ALL"}
	securityOpt := []string{"no-new-privileges"}

	hostConfig := &container.HostConfig{
		Mounts:       []mount.Mount{homeMount},
		Resources:    resources,
		PortBindings: portBindings,
		Tmpfs:        tmpfs,
		CapDrop:      capDrop,
		SecurityOpt:  securityOpt,
	}

	return hostConfig
}

// CreateContainer creates (but does not start) a workbench container: the hardcoded workbench
// image, attached to a dedicated network created just for it (see containerNetworkName), with
// opts.VolumeName mounted at homeMountPath, hardcoded CPU/memory/PID limits, capability/privilege
// hardening (see newWorkbenchHostConfig), and labeled for operational visibility.
//
// The container runs two independently-addressable in-container servers: the chat bridge
// (bridgePort) and ttyd (ttydPort, the interactive tmux-tab terminal). Each is published to a
// host port allocated from its own fixed range ([bridgeHostPortRangeStart, bridgeHostPortRangeEnd]
// and [ttydHostPortRangeStart, ttydHostPortRangeEnd] respectively — see allocateHostPort) rather
// than left unpublished. Neither is exposed to the wider internet by this alone: the configured
// daemon is expected to sit on a network only Artel's own process can reach, and both servers are
// only ever dialed through internal/transport/vaults_api's authenticated reverse proxy, never
// linked directly to a client.
func (c *Client) CreateContainer(ctx context.Context, opts CreateOpts) (string, error) {
	tag, err := c.EnsureImage(ctx)
	if err != nil {
		return "", rerrors.Wrap(err, "error ensuring workbench image")
	}

	bridgeHostPort, err := c.allocateHostPort(ctx, bridgeNatPort, bridgeHostPortRangeStart, bridgeHostPortRangeEnd)
	if err != nil {
		return "", rerrors.Wrap(err, "allocating bridge host port")
	}

	ttydHostPort, err := c.allocateHostPort(ctx, ttydNatPort, ttydHostPortRangeStart, ttydHostPortRangeEnd)
	if err != nil {
		return "", rerrors.Wrap(err, "allocating ttyd host port")
	}

	labels := map[string]string{
		workbenchLabelKey: workbenchLabelValue,
	}

	networkName := containerNetworkName(opts.Name)

	networkCreateOpts := network.CreateOptions{
		Labels: labels,
	}

	_, err = c.cli.NetworkCreate(ctx, networkName, networkCreateOpts)
	if err != nil {
		return "", rerrors.Wrap(err, "creating per-container workbench network")
	}

	exposedPorts := nat.PortSet{
		bridgeNatPort: struct{}{},
		ttydNatPort:   struct{}{},
	}

	containerConfig := &container.Config{
		Image:        tag,
		Labels:       labels,
		ExposedPorts: exposedPorts,
	}

	bridgePortBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: bridgeHostPort,
	}
	ttydPortBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: ttydHostPort,
	}
	portBindings := nat.PortMap{
		bridgeNatPort: []nat.PortBinding{bridgePortBinding},
		ttydNatPort:   []nat.PortBinding{ttydPortBinding},
	}

	hostConfig := newWorkbenchHostConfig(opts.VolumeName, portBindings)

	workbenchEndpoint := &network.EndpointSettings{}
	endpointsConfig := map[string]*network.EndpointSettings{
		networkName: workbenchEndpoint,
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

// allocateHostPort picks a free host port in [rangeStart, rangeEnd] for a new workbench
// container's natPort binding (bridgeNatPort or ttydNatPort), by inspecting every existing
// artel.workbench container (running or not — HostConfig.PortBindings reflects a container's
// declared binding regardless of its current state) and returning the first port in range none
// of them already claim for that same natPort.
//
// This is a best-effort reservation, not a lock: two concurrent CreateContainer calls could both
// observe the same free port and both pick it, in which case the daemon rejects whichever one
// starts second with a "port is already allocated" error — an acceptable failure mode for the
// prototype's expected (low, per-admin) concurrency, same tradeoff as the hardcoded resource
// limits elsewhere in this package.
func (c *Client) allocateHostPort(ctx context.Context, natPort nat.Port, rangeStart, rangeEnd int) (string, error) {
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

		for _, binding := range inspect.HostConfig.PortBindings[natPort] {
			usedPorts[binding.HostPort] = struct{}{}
		}
	}

	for port := rangeStart; port <= rangeEnd; port++ {
		portStr := strconv.Itoa(port)

		_, taken := usedPorts[portStr]
		if !taken {
			return portStr, nil
		}
	}

	return "", rerrors.New("no free host port in range for " + string(natPort))
}

// StartContainer starts an already-created workbench container and injects env into it.
//
// The Docker Engine API has no way to attach environment variables to a container at `docker
// start` time — env is otherwise only settable at `docker create` time, which is deliberately
// too early here (see CreateOpts). Once started, env is instead written into the container's
// envDropDir tmpfs via a docker exec (see injectEnv), so the value only ever lives in that
// tmpfs's pages — never baked into the image, never written to the mounted volume. This is the
// mechanism the api_key/subscription_login flows (internal/service/v1/workbench/workbench.go)
// rely on to keep that property intact.
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

// injectEnv writes each entry of env into the container's envDropDir tmpfs, one docker exec and
// one file per variable — the file's name is the variable's name and its content the raw value,
// which is the contract the in-container chat bridge reads back before each `claude` invocation.
//
// Each exec's own process environment (ExecOptions.Env) carries the value in; only the
// variable's name ever appears in argv, never the value, keeping it out of `docker
// exec`/process-list logs. `umask 077` plus the mode=0700 tmpfs means a dropped file is never
// readable by anything but root even transiently, and `printf %s` (not `echo`) writes the value
// byte-exactly, with no trailing newline and no backslash-escape interpretation.
func (c *Client) injectEnv(ctx context.Context, containerID string, env map[string]string) error {
	for key, value := range env {
		execOptions := container.ExecOptions{
			Cmd: []string{
				"/bin/sh", "-c",
				`umask 077; mkdir -p "$1"; printf %s "$VALUE" > "$1/$2"`,
				"_", envDropDir, key,
			},
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

// ContainerAddress returns containerID's "<daemon-host>:<published-port>" address for the chat
// bridge running inside it, as reachable from Artel's own process. Used by the workbench
// reverse-proxy handler in internal/transport/vaults_api.
//
// This deliberately does not resolve the container's IP on its dedicated per-container network
// (see containerNetworkName): the configured daemon (c.host) is commonly itself a docker:dind
// *container*, whose inner networks live in a network namespace private to that container, so a
// workbench container's IP on any of them is never routable from outside the dind container
// regardless of where Artel's own process runs. Reading back the Docker-assigned host port
// CreateContainer published bridgePort to and pairing it with c.host's own hostname instead
// routes through the dind container's (or, for a bare/second-dockerd host, the host's)
// already-reachable address — the same address Artel used to create/start the container in the
// first place.
//
// Fails rather than guessing when the daemon hasn't assigned a host port yet (i.e. the container
// isn't running) — an empty/absent address would otherwise surface much later as an opaque proxy
// dial failure.
func (c *Client) ContainerAddress(ctx context.Context, containerID string) (string, error) {
	return c.publishedAddress(ctx, containerID, bridgeNatPort)
}

// TtydAddress returns containerID's "<daemon-host>:<published-port>" address for the ttyd server
// running inside it (the interactive tmux-tab terminal, restored alongside the chat bridge), as
// reachable from Artel's own process. Mirrors ContainerAddress exactly, reading back
// ttydNatPort's binding instead of bridgeNatPort's — a workbench container now runs two
// independently-addressable in-container servers, so each gets its own address method.
func (c *Client) TtydAddress(ctx context.Context, containerID string) (string, error) {
	return c.publishedAddress(ctx, containerID, ttydNatPort)
}

// publishedAddress resolves containerID's "<daemon-host>:<published-port>" address for natPort
// (bridgeNatPort or ttydNatPort), as reachable from Artel's own process — the shared
// implementation behind ContainerAddress and TtydAddress; see ContainerAddress's doc comment for
// why this reads back the published host port rather than the container's network IP.
func (c *Client) publishedAddress(ctx context.Context, containerID string, natPort nat.Port) (string, error) {
	inspect, err := c.cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return "", rerrors.Wrap(err, "inspecting workbench container")
	}

	if inspect.NetworkSettings == nil {
		return "", rerrors.New("workbench container has no network settings: " + containerID)
	}

	bindings, ok := inspect.NetworkSettings.Ports[natPort]
	if !ok || len(bindings) == 0 {
		return "", rerrors.New("workbench container has no published port " + string(natPort) + ": " + containerID)
	}

	hostPort := bindings[0].HostPort
	if hostPort == "" {
		return "", rerrors.New("workbench container's port " + string(natPort) + " has no assigned host port: " + containerID)
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

// RemoveContainer force-removes a workbench container and the dedicated per-container network
// CreateContainer created for it (see containerNetworkName). It does not remove the container's
// volume — call RemoveVolume separately.
//
// The network(s) containerID is attached to are read back via ContainerInspect before removal —
// once the container is gone there is nothing left to inspect them from — and each is deleted via
// NetworkRemove only after ContainerRemove succeeds, since a network can't be removed while a
// container endpoint is still attached to it. workbenchNetworkName is deliberately never removed
// here even if found attached: it's the old shared network every workbench container used to
// share (see its own doc comment), so removing it out from under a container created before this
// change would break every other container still using it.
func (c *Client) RemoveContainer(ctx context.Context, containerID string) error {
	networkNames, err := c.attachedNetworks(ctx, containerID)
	if err != nil {
		return rerrors.Wrap(err, "inspecting workbench container networks")
	}

	removeOptions := container.RemoveOptions{
		Force: true,
	}

	err = c.cli.ContainerRemove(ctx, containerID, removeOptions)
	if err != nil {
		return rerrors.Wrap(err, "removing workbench container")
	}

	for _, networkName := range networkNames {
		err = c.cli.NetworkRemove(ctx, networkName)
		if err != nil {
			return rerrors.Wrap(err, "removing per-container workbench network "+networkName)
		}
	}

	return nil
}

// attachedNetworks returns the names of every Docker network containerID is currently attached
// to that RemoveContainer should also remove (see filterRemovableNetworks). Returns an empty,
// non-nil slice if the container has no recorded network settings at all rather than erroring,
// since that's a shape ContainerInspect can legitimately return.
func (c *Client) attachedNetworks(ctx context.Context, containerID string) ([]string, error) {
	inspect, err := c.cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return nil, rerrors.Wrap(err, "inspecting workbench container")
	}

	if inspect.NetworkSettings == nil {
		return []string{}, nil
	}

	names := make([]string, 0, len(inspect.NetworkSettings.Networks))

	for name := range inspect.NetworkSettings.Networks {
		names = append(names, name)
	}

	return filterRemovableNetworks(names), nil
}

// filterRemovableNetworks returns the subset of networkNames RemoveContainer should delete: every
// one except workbenchNetworkName. Split out from attachedNetworks as a pure function so the
// exclusion — the one thing that must never regress, since removing the old shared network out
// from under a container created before this change would break every other container still
// using it — can be unit tested without a Docker daemon.
func filterRemovableNetworks(networkNames []string) []string {
	removable := make([]string, 0, len(networkNames))

	for _, name := range networkNames {
		if name == workbenchNetworkName {
			continue
		}

		removable = append(removable, name)
	}

	return removable
}
