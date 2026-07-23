# Docker Topology

## Gating: presence of config, not a feature flag

Add a config field the same way every other data-source connection string is added per
`docs/architecture.md` ("Adding a new env variable"): an entry under `environment:` in
`config/config.yaml`, then `rscli-dev project tidy` regenerates `internal/config/environment.go`.

```
WorkbenchDockerHost string   // e.g. "unix:///var/run/docker-workbenches.sock" or "tcp://host:2376"
```

When empty, `custom.go` simply does not construct `WorkbenchService` and does not register its
transport routes — the same "absent, not degraded" shape as everything gated behind
`SubscriptionsEnabled` today, except here it's a service that either exists or doesn't, not a
`FreeService`/`PaidService` swap (there's no meaningful no-op workbench).

## Where the daemon actually runs

`DOCKER_HOST` (or the SDK-equivalent connection option) is just a client-side pointer to *a*
daemon — it carries no implication about where that daemon lives. Three real options:

### Option A — same daemon as everything else (`unix:///var/run/docker.sock`)
Zero extra infra: mount the existing socket, tag workbench containers with a label
(`artel.workbench=true`) and put them on their own bridge network for egress control. Simplest to
stand up, but workbench containers share a security boundary (the daemon itself) with whatever
else runs on that host.

### Option B — second `dockerd` process, same VM (recommended for the prototype)
Run a second daemon bound to its own socket and data-root:

```
dockerd --data-root /var/lib/docker-workbenches \
        -H unix:///var/run/docker-workbenches.sock
```

`WorkbenchDockerHost=unix:///var/run/docker-workbenches.sock`. This gives a genuinely separate
daemon — own storage driver, own container/image namespace, a crash or resource exhaustion in one
doesn't touch the other — without provisioning new hardware. This is the sweet spot for the
prototype: real separation, one `systemd` unit to add, no new machine, no TLS cert management.

### Option C — dedicated remote host
`WorkbenchDockerHost=tcp://workbench-host.internal:2376`, mutual TLS required (an unauthenticated
Docker API over plain TCP is unauthenticated root on that host — never do this even internally).
Best isolation, and lets workbench capacity scale independently of Artel's app servers, but is
real ops work (a box, cert issuance/rotation, firewall rules). Defer past the prototype; revisit
once workbench container count or resource needs actually justify separate hardware.

**Prototype decision: Option B.** Revisit only if/when there's a concrete reason (capacity,
compliance, blast-radius requirement) to move to Option C — don't build the mTLS/remote-host path
speculatively.

## Client wiring

New `internal/clients/workbenchdocker` package (naming: avoid colliding with a hypothetical future
generic `internal/clients/docker` if Artel ever needs Docker access for something unrelated),
wrapping `github.com/docker/docker/client` — follows the existing
`internal/clients/{couchdb,anthropic,imap,smtp}` convention of one small typed wrapper per
external system, constructed once in `custom.go` with the configured host and handed to
`WorkbenchService`.

Minimal surface needed for the prototype:

```go
type Client struct { /* docker/docker/client.Client + configured host */ }

func New(host string) (*Client, error)
func (c *Client) CreateContainer(ctx context.Context, opts CreateOpts) (containerID string, err error)
func (c *Client) StartContainer(ctx context.Context, containerID string, env map[string]string) error
func (c *Client) StopContainer(ctx context.Context, containerID string) error
func (c *Client) RemoveContainer(ctx context.Context, containerID string) error
func (c *Client) CreateVolume(ctx context.Context, name string) error
func (c *Client) RemoveVolume(ctx context.Context, name string) error
```

## Network isolation

Workbench containers go on a dedicated Docker network (`workbench-net`) with no route to Artel's
internal service network (Postgres, CouchDB, MinIO). Required outbound: HTTPS to
`api.anthropic.com` (and, for subscription login, whatever Anthropic's OAuth/console endpoints
are) plus Artel's own gRPC gateway if/when the container ever needs to call back in — not needed
for the prototype's create/start/login scope. No inbound ports exposed per container; all access
goes through Artel-side proxying (deferred design — see the README's "out of scope" note; the
prototype's access mechanism is `docker exec`/`attach` driven from the backend, not a
user-reachable port).

## Resource limits

Hardcode CPU/mem limits at `docker create` time for the prototype (e.g. 1 CPU, 2GB RAM, a fixed
volume size if the storage driver supports quotas) — a config/tier-driven limit is a billing
concern to design once the metering story (see `docs/workbench/04_task_breakdown.md`'s deferred
stages) is actually being built, not before.
