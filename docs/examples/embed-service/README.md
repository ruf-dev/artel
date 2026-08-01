# embed-service

A minimal, working example of embedding artel as a library inside a
**separate Go module** — proving `github.com/ruf-dev/artel/pkg/app` works
the way an external consumer would actually use it, not just as an in-repo
import. This directory has its own `go.mod`; it is not part of the
`github.com/ruf-dev/artel` module and is invisible to `go build ./...` run
from the repo root.

## Why `pkg/app`

`internal/app` (what `cmd/service/main.go` uses directly) can't be imported
from outside the `github.com/ruf-dev/artel` module — that's Go's
`internal/` visibility rule. `pkg/app` is the sanctioned public re-export
for exactly this "run artel inside my own Go program" case: `app.App` is a
type alias for `internal/app.App`, so its exported fields and methods
(`Start`, `Custom.Transport`, ...) stay callable from outside the module
without ever importing `internal/app` yourself. See the doc comment on
`pkg/app/app.go` for the full rationale.

## What this demonstrates

- `main.go` calls `app.New()`, registers one extra HTTP handler
  (`/example/hello`) on the embedded app's transport via
  `a.Custom.Transport.AddHttpHandler`, then calls `a.Start()` — the same
  three-call shape as `cmd/service/main.go`, just from outside the module.
- `main_test.go` starts the real app (against a real Postgres), hits that
  handler over HTTP, then shuts the app down cleanly via
  `go.redsock.ru/toolbox/closer.Close()` — the same shutdown path
  production takes on `SIGINT`/`SIGTERM`.

## Prerequisites

- Go toolchain (matching the root repo's `go 1.25.5`).
- Docker, for Postgres. No CouchDB, S3, or MinIO is needed — `svcv1.New`
  wires only Postgres-backed repos and pure Go structs at startup; the
  other backends are dialed lazily on first real use.

## Run steps

```bash
# from this directory
docker compose up -d

go test ./...    # starts the embedded app, hits /example/hello, shuts it down

# runs it for real; required.env has everything go run . needs
set -a && source required.env && set +a
go run .         # curl http://localhost:8080/example/hello

docker compose down
```

## The config-skeleton-plus-env-vars pattern

Artel's config loader (`internal/config/load.go`) reads a YAML file, then
lets env vars override values on nodes that already exist in it — it
cannot materialize a brand-new server or datasource purely from env vars
with zero file present. So embedding artel means shipping a small config
**skeleton** file that declares the shape (one `servers` entry named
`MASTER`, one `data_sources` entry with `resource_name: postgres`, the
`environment` list) — see `config/config.yaml` in this directory, adapted
from the root repo's `config/config_template.yaml`. Actual runtime values
(host, port, credentials, ...) are then supplied by env vars named per
`config/.env.example` in the root repo (e.g.
`DATA_SOURCES_POSTGRES_HOST`, `SERVERS_MASTER_PORT`) — see the env vars set
in `main_test.go` for a complete example set. This is the intended,
supported way to configure an embedded instance, not a workaround.
