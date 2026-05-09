# Task 06 — Vault Note Storage via MCP

## Goal

Extend the MCP server with note-level operations so Claude can read and write notes inside a vault (CouchDB database).

## Scope

- Add `NoteService` interface to `internal/service/interfaces.go`
- Implement in `internal/service/v1/note/note.go`
- Add CouchDB document-level methods to `internal/clients/couchdb/` (get, put, delete doc)
- Expose MCP tools: `write_note`, `read_note`, `delete_note`, `list_notes`

## Data Model

A note is a CouchDB document in the vault database:
```json
{
  "_id": "<note-id>",
  "title": "...",
  "content": "...",
  "updated_at": "..."
}
```

## Acceptance Criteria

- `go build ./...` + `go test ./...` pass
- `write_note` creates or updates a document in the named vault
- `read_note` returns document content
- `list_notes` returns all document IDs + titles in a vault
- `delete_note` removes a document

## Notes

- Vault must exist before writing notes; return a clear error if it does not
- Note IDs are caller-supplied strings (e.g. file path like `folder/note.md`)
- CouchDB `_rev` handling: fetch current rev before update/delete
