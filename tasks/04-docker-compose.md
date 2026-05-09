# Task: Local Development Docker Compose

## Goal
Provide a `docker-compose.yaml` for local development with CouchDB.

## Scope
- CouchDB container with admin credentials matching `config/dev.yaml`
- Expose CouchDB on a non-default port (e.g., 15984) to avoid conflicts
- Persistent volume for CouchDB data at `data/couch/`
- `make up` / `make down` convenience targets in `Makefile`

## Notes
- CouchDB default port is 5984; use 15984 locally
- Single-node mode (`COUCHDB_NODENAME=nonode@nohost`) for simplicity
- Admin credentials: set via env vars in compose, mirror in `config/dev.yaml`
