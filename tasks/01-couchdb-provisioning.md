# Task: CouchDB Provisioning Layer

## Goal
Implement CouchDB client and provisioning logic for Obsidian vault creation.

## Scope

### Storage layer (`internal/storage/couchdb/`)
- CouchDB HTTP client wrapping the CouchDB REST API
- Create/delete database (maps to an Obsidian vault)
- Create CouchDB user with scoped access to the vault database
- Set database security (members/admins)

### Config
- Add CouchDB connection config to `internal/config/` (host, port, admin credentials)
- Add to `config/config.yaml` and `config/dev.yaml`

### Service layer (`internal/service/`)
- Define `VaultService` interface
- Implement vault provisioning: create DB + user + set permissions atomically

## Notes
- CouchDB admin API lives at `/_users` (user docs) and `/{db}/_security`
- Use `go.redsock.ru/rerrors` for all error wrapping
- Wire in `internal/app/custom.go`
