// Package workbenchdocker is a thin typed wrapper around the Docker Engine SDK
// (github.com/docker/docker/client), used to create/start/stop/remove the containers and
// volumes that back a single workbench (a per-user sandbox running the `claude` CLI inside
// a persistent tmux session).
//
// This package talks to whichever daemon `WorkbenchDockerHost` points at — see
// docs/workbench/02_docker_topology.md for why that daemon is expected to be a dedicated
// second dockerd process, not the daemon Artel's own containers run on. It intentionally does
// not create the `workbench-net` network itself (assumed to pre-exist on the configured
// daemon) and does not expose any inbound port on the containers it creates.
package workbenchdocker

import (
	"github.com/docker/docker/client"

	"go.redsock.ru/rerrors"
)

const (
	// workbenchImage is the image every workbench container is created from. Hardcoded for
	// the prototype — see docs/workbench/02_docker_topology.md.
	workbenchImage = "artel-workbench:latest"

	// workbenchNetworkName is the dedicated, isolated Docker network workbench containers are
	// attached to. Not created by this client — assumed to pre-exist on the configured daemon.
	workbenchNetworkName = "workbench-net"

	// workspaceMountPath is the fixed in-container path the workbench's named volume is
	// mounted at.
	workspaceMountPath = "/workspace"

	// workbenchLabelKey/workbenchLabelValue tag every container/volume this client creates,
	// for operational visibility even on a dedicated daemon.
	workbenchLabelKey   = "artel.workbench"
	workbenchLabelValue = "true"

	// Resource limits, hardcoded at create time for the prototype (see
	// docs/workbench/02_docker_topology.md, "Resource limits").
	workbenchCpuLimitNanoCpus = 1_000_000_000          // 1 CPU
	workbenchMemLimitBytes    = 2 * 1024 * 1024 * 1024 // 2GB
)

// Client wraps the Docker SDK client with the narrow surface a workbench's container/volume
// lifecycle needs.
type Client struct {
	cli  *client.Client
	host string
}

// New constructs a Client talking to the Docker daemon at host (e.g.
// "unix:///var/run/docker-workbenches.sock" or "tcp://host:2376"). API version negotiation is
// enabled so the client stays compatible with the daemon regardless of exactly which API
// version it speaks.
func New(host string) (*Client, error) {
	cli, err := client.NewClientWithOpts(
		client.WithHost(host),
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, rerrors.Wrap(err, "creating docker client")
	}

	c := &Client{
		cli:  cli,
		host: host,
	}

	return c, nil
}
