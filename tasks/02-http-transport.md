# Task: HTTP Transport Layer

## Goal
Expose REST endpoints for vault management operations.

## Scope

### Transport (`internal/transport/`)
- HTTP server setup (net/http or chi router)
- `POST /vault` — create new vault (provisions CouchDB DB + user)
- `DELETE /vault/{id}` — tear down vault
- `GET /vault/{id}` — get vault status/details
- Request/response models in `internal/transport/model/`

### Auth
- Decide and implement auth strategy (JWT, API key, or session)
- Middleware in `internal/transport/middleware/`

### Wire up
- Register HTTP server in `internal/app/custom.go` `Init` and `Start`/`Stop`
- Add server port to config

## Notes
- Keep handler thin: validate input, call service, return JSON
- Use `zerolog` for request logging
