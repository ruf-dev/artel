# Docker Topology

## Gating: DB-backed pool, not a single config value

Docker hosts are no longer a single startup-time config field. They're admin-managed rows in the
`docker_hosts` table (migration `061_docker_hosts.sql`) — one row per Docker daemon endpoint
(`url` only, e.g. `unix:///var/run/docker-workbenches.sock` or `tcp://host:2376`; no credentials,
unlike `couch_instances`/`s3_instances`), CRUD'd via the `/api/docker_hosts/*` routes
(`internal/transport/docker_hosts_api`) and the Docker Api admin tab, the same shape as CouchDB/S3
instance management.

`WorkbenchService` is now **always constructed** in `custom.go` — there's no "absent config" case
that leaves the whole subsystem unconstructed the way there used to be. Each `CreateWorkbench`
call picks the least-loaded registered host (`DockerHosts.PickLeastLoaded`, by live workbench
count, ties broken oldest-first) and pins the new workbench to it (`workbenches.docker_host_id`);
every subsequent operation on that workbench (start/stop/delete/login flow) resolves its Docker
client from that same pinned host rather than a single shared client instance. "No hosts
registered yet" is therefore a *runtime* error surfaced from `CreateWorkbench`
(`user_errors.NoDockerHostsAvailable`), not a startup-time absence — an admin can register a host
at any time and workbench creation starts working immediately, no restart needed.

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

Register `unix:///var/run/docker-workbenches.sock` as a `docker_hosts` row. This gives a genuinely
separate daemon — own storage driver, own container/image namespace, a crash or resource
exhaustion in one doesn't touch the other — without provisioning new hardware. This is the sweet
spot for the prototype: real separation, one `systemd` unit to add, no new machine, no TLS cert
management.

### Option C — dedicated remote host
Register `tcp://workbench-host.internal:2376` as a `docker_hosts` row, mutual TLS required (an unauthenticated
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
external system. Unlike those, a single cached instance isn't handed to `WorkbenchService` at
construction time — with more than one registered host, `WorkbenchService` builds a fresh
`workbenchdocker.New(host.Url)` per call, pointed at whichever host the workbench in question is
pinned to (mirrors `couchinstances.Service`'s per-call `couchdb.New(cfg)`, not a cached client).

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
