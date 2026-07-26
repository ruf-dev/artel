# Docker Topology

## Gating: DB-backed pool, not a single config value

Docker hosts are no longer a single startup-time config field. They're admin-managed rows in the
`docker_hosts` table (migration `061_docker_hosts.sql`) — one row per Docker daemon endpoint
(`url`, e.g. `unix:///var/run/docker-workbenches.sock` or `tcp://host:2376`, plus three optional
TLS/mTLS fields — migration `062_docker_hosts_tls.sql`'s `ca_cert_enc`/`client_cert_enc`/
`client_key_enc`, encrypted at rest the same way as `couch_instances`/`s3_instances` credentials),
CRUD'd via the `/api/docker_hosts/*` routes (`internal/transport/docker_hosts_api`) and the Docker
Api admin tab, the same shape as CouchDB/S3 instance management. The cert fields are write-only —
`GetDockerHost`/`ListDockerHosts` never return them, only `RegisterDockerHost`/`UpdateDockerHost`
accept them, and only `repository.DockerHosts.GetWithCreds` (used internally when resolving a
Docker client, never by the admin API) decrypts and reads them back.

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

For local dev/testing, `tests/docker-compose.yaml`'s `test-dockerd` service (`docker:dind`,
`DOCKER_TLS_CERTDIR=""`, exposed at `tcp://localhost:12375`) automates a throwaway stand-in for
this option — plain, unauthenticated TCP, since it's never reachable outside localhost.
`make setup-dev-env` (or `make setup-dev-docker-host` on its own) brings it up, creates
`workbench-net` inside it, and registers it as a `docker_hosts` row via
`scripts/setup_dev_docker_host.sh`. Note dind runs its own separate inner `dockerd` with its own
network namespace, so `workbench-net` has to be created *inside* that daemon specifically — a
compose-level `networks:` block alone does not reach it; the script handles this via
`docker -H tcp://localhost:12375 network create workbench-net`.

### Option C — dedicated remote host (implemented)
Register `tcp://workbench-host.internal:2376` as a `docker_hosts` row with `ca_cert`/
`client_cert`/`client_key` set (mutual TLS required — an unauthenticated Docker API over plain TCP
is unauthenticated root on that host — never do this even internally). Best isolation, and lets
workbench capacity scale independently of Artel's app servers, but is real ops work (a box, cert
issuance/rotation, firewall rules). The registration/storage/client-wiring side is implemented
(migration `062_docker_hosts_tls.sql`, `internal/clients/workbenchdocker.TLSConfig`); actually
provisioning and pointing at a remote host is still deferred past the prototype — revisit once
workbench container count or resource needs actually justify separate hardware.

**Prototype decision: Option B.** Revisit only if/when there's a concrete reason (capacity,
compliance, blast-radius requirement) to actually run Option C in production — the code path
exists now, but standing up a real remote host is still out of scope for the prototype.

## Client wiring

New `internal/clients/workbenchdocker` package (naming: avoid colliding with a hypothetical future
generic `internal/clients/docker` if Artel ever needs Docker access for something unrelated),
wrapping `github.com/docker/docker/client` — follows the existing
`internal/clients/{couchdb,anthropic,imap,smtp}` convention of one small typed wrapper per
external system. Unlike those, a single cached instance isn't handed to `WorkbenchService` at
construction time — with more than one registered host, `WorkbenchService` builds a fresh
`workbenchdocker.New(host.Url, tlsCfg)` per call, pointed at whichever host the workbench in
question is pinned to (mirrors `couchinstances.Service`'s per-call `couchdb.New(cfg)`, not a
cached client). `tlsCfg` comes from `DockerHosts.GetWithCreds`; an empty `TLSConfig{}` (the local
unix-socket/dind case) reproduces exactly the original no-TLS `client.WithHost` behavior — `New`
only builds a custom `*http.Client` (`tls.X509KeyPair` + an `x509.CertPool`, since the Docker
SDK's `WithTLSClientConfig` wants file paths and this repo passes decrypted secrets as in-memory
values) when at least one TLS field is set.

Minimal surface needed for the prototype:

```go
type Client struct { /* docker/docker/client.Client + configured host */ }

type TLSConfig struct { CaCert, ClientCert, ClientKey string } // empty = no TLS

func New(host string, tlsCfg TLSConfig) (*Client, error)
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
