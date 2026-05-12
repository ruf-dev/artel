# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**artel** is a Go service that automates Obsidian in-cloud vault creation. It provisions CouchDB instances, creates
users and resources within them, and exposes a React UI at `pkg/client/ArtelUI` (Bun + React). Generated
by [RedSock CLI (rscli)](https://github.com/Red-Sock/rscli).

## Common Commands

```bash
go build -v ./...              # Build all packages
go run ./cmd/service/main.go   # Run locally (prod config)
go run ./cmd/service/main.go -dev   # Run with dev config (./config/dev.yaml)
golangci-lint run ./...        # Lint (uses .golangci.yaml, exit code 42 = issues found)
go test ./...                  # Run all tests
go mod tidy                    # Tidy dependencies
make build-local-container     # Build Docker image for linux/arm64
```

### Frontend (pkg/client/ArtelUI)

```bash
cd pkg/client/ArtelUI
bun install
bun run dev       # Dev server
bun run build     # Production build
```

### Code Generation

```bash
make warmup    # Install proto dependencies (protopack mod download)
moti g         # Generate protos (preferred — runs protoc via moti.yaml config)
sqlc generate  # Generate Go from SQL queries (internal/repository/pg/queries → internal/repository/pg/generated)
make lint      # golangci-lint run ./...
```

## Architecture

### rscli Application Structure

rscli scaffolds a fixed two-part app lifecycle:

- **`internal/app/app.go`** — Generated, DO NOT EDIT. Wires config → `Custom.Init` → `Custom.Start`. Handles graceful
  shutdown via `go.redsock.ru/toolbox/closer`.
- **`internal/app/custom.go`** — User-editable entry point. `Custom.Init` wires repositories, services, and transports.
  `Custom.Start` launches them. `Custom.Stop` tears them down.
- **`internal/app/config.go`** — Generated. Initializes context and loads config via matreshka.

All new application wiring belongs in `custom.go` — never modify `app.go` or `config.go`.

### Configuration (matreshka)

Config is YAML-based using the [matreshka](https://go.vervstack.ru/matreshka) framework:

- `config/config.yaml` — production config (committed)
- `config/dev.yaml` — local dev overrides (run with `-dev` flag)
- `config/config_template.yaml` — template for generated config

Env vars parsed into `internal/config/environment.go` (`EnvironmentConfig`). Data source connection strings and server
addresses go in separate config structs added to `internal/config/`.

The `-config <path>` flag overrides the config file at runtime.

**Adding a new env variable**: edit `config/config.yaml` (add the entry under `environment:`), then run
`rscli-dev project tidy`. This regenerates `internal/config/environment.go` and the new field appears in
`EnvironmentConfig`. Never edit `environment.go` by hand.

### Planned Layers

```
cmd/service/main.go
  → internal/app/custom.go (Init/Start/Stop)
    → internal/transport/   (HTTP/gRPC handlers)
    → internal/service/     (business logic)
    → internal/clients/     (external service clients, e.g. CouchDB)
pkg/client/ArtelUI/         (React UI, built with Bun)
```

### Service Layer Layout

Interfaces live in `internal/service/interfaces.go` (package `service`). Implementations live under
`internal/service/v1/`:

```
internal/service/
  interfaces.go          ← all service interfaces, package service
  v1/
    impl.go              ← Services struct (container for all v1 implementations)
    vault/
      vault.go           ← VaultService implementation
```

`v1.Services` is constructed in `custom.go` and its fields (e.g. `services.Vault`) are passed to transports. Transport
handlers accept interfaces from the `service` package, never concrete v1 types.

### Key Dependencies

| Package                        | Purpose                                        |
|--------------------------------|------------------------------------------------|
| `go.redsock.ru/rerrors`        | Error wrapping with `rerrors.Wrap(err, "msg")` |
| `go.redsock.ru/toolbox/closer` | Graceful shutdown registry                     |
| `github.com/rs/zerolog`        | Structured logging via `log.Info().Msg(...)`   |
| `go.vervstack.ru/matreshka`    | Config loading                                 |
| `golang.org/x/sync/errgroup`   | Concurrent startup                             |

## Go Coding Rules

- **Never check function errors inline.** Always split: `err := func()` on one line, then `if err != nil` on the next.
- **Never create struct/value literals inline in a function call.** Assign to a named variable first, then pass it.
- `internal/app/app.go` and `internal/app/config.go` are generated — do not edit them. `custom.go` is the user-editable
  counterpart.
- Use `rerrors.Wrap(err, "context message")` for all error wrapping.
- **No all-caps field names.** Use mixed-case: `Id` not `ID`. For `uuid.UUID` fields use the type name as the suffix:
  `Uuid` (primary key), `UserUuid`, `VaultUuid`, etc.