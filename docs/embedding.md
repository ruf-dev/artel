# Embedding artel

How to run artel as a library inside another Go application, instead of as
its own `cmd/service` binary.

## Why `pkg/app`

`internal/app` — what `cmd/service/main.go` calls directly — cannot be
imported from outside the `github.com/ruf-dev/artel` module; that's Go's
`internal/` visibility rule. `pkg/app` exists as the sanctioned public
re-export for exactly this case: `app.App` is a type alias for
`internal/app.App`, so its exported fields and methods (`Start`,
`Custom.Transport`, ...) are usable from an external module without ever
importing `internal/app`. An embedding app depends on
`github.com/ruf-dev/artel` and writes:

```go
import "github.com/ruf-dev/artel/pkg/app"

a, err := app.New()
// ...
err = a.Start()
```

## Requirements

- **Postgres reachable.** `app.New()` opens the Postgres connection and runs
  the goose migrations embedded in `github.com/ruf-dev/artel/migrations`
  (`migrations.ApplyMigration`, called from `internal/clients/postgres.Migrate`)
  synchronously as part of startup. The `.sql` files are compiled into the
  binary via `go:embed`, so no `migrations/` folder needs to exist on disk
  for the embedding app. No CouchDB, S3, or MinIO connection is needed at
  startup — those are dialed lazily on first real use.
- **A config skeleton file.** Artel's config loader
  (`internal/config/load.go`) reads a YAML file first, then lets env vars
  override values on nodes that already exist in it — env vars alone can't
  materialize a new server or datasource entry with zero file present. So
  the embedding app needs its own `./config/config.yaml` (or a path passed
  via `-config`) declaring the skeleton shape: one `servers` entry named
  `MASTER`, one `data_sources` entry with `resource_name: postgres`, and the
  `environment` list. Use `config/config_template.yaml` in this repo as the
  starting point. Runtime values (host, port, credentials, migrations
  folder, ...) are then supplied by env vars named per `config/.env.example`
  (e.g. `DATA_SOURCES_POSTGRES_HOST`, `SERVERS_MASTER_PORT`,
  `ENVIRONMENT_LOG_LEVEL`). This file-skeleton-plus-env-vars split is the
  intended, supported configuration path — not a workaround.

## Extending the embedded instance

`app.New()` populates `a.Custom.Transport` (an
`*internal/transport.ServersManager`, reachable without importing
`internal/transport` since its type is never spelled out by the caller)
before returning. Register additional routes/services on it *after*
`app.New()` returns and *before* calling `a.Start()` — the transport only
starts accepting connections once `Start()` runs:

```go
a, err := app.New()
// handle err

a.Custom.Transport.AddHttpHandler("/my-path", myHandler)
a.Custom.Transport.AddImplementation(myGrpcImpl)

err = a.Start()
// handle err
```

Mirror `internal/app/custom.go`'s own `AddHttpHandler`/`AddImplementation`
calls near the end of `Custom.Init` for style, and avoid paths artel already
registers there (`/`, `/mcp`, `/webhooks/gitlab/`, `/tract/hook/`,
`/.well-known/...`, `/register`, `/oauth/...`, `/token`).

## Full example

`docs/examples/embed-service/` is a complete, tested, buildable example: a
genuinely separate Go module (its own `go.mod`, importing artel via a
`replace` back to a local checkout) that registers one extra HTTP handler
and has a test exercising the whole startup → request → graceful-shutdown
lifecycle against a real Postgres via `go.redsock.ru/toolbox/closer`. Its
`required.env` is a ready-to-copy env file for the full set of variables
described above. See its `README.md` for exact run steps.
