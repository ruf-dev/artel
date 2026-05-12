---
status: done
---

# Task 15 — CouchDB LiveSync Client Extension

## Goal

Add CouchDB operations for reading and writing Obsidian notes in the LiveSync document format.

## Context

Obsidian Self-Hosted LiveSync stores notes as CouchDB documents. The document structure:

```json
{
  "_id": "path/to/note.md",
  "_rev": "1-abc...",
  "data": "<base64-encoded note content>",
  "mtime": 1700000000000,
  "ctime": 1700000000000,
  "size": 1234,
  "deleted": false,
  "type": "plain"
}
```

For large files, content is split into chunks (`type: "newnote"` with `children` array), but for MVP only handle `type: "plain"` (files under ~100KB).

## File to Create

`internal/clients/couchdb/livesync.go`

Implement a `LiveSyncClient` struct with the following methods. It connects to CouchDB using a base URL, database name, username, and password (Basic Auth on every request).

```go
type LiveSyncClient struct {
    baseURL  string // e.g. "http://localhost:5984"
    dbName   string
    username string
    password string
    http     *http.Client
}

func NewLiveSyncClient(baseURL, dbName, username, password string) *LiveSyncClient

// ListNotes returns all non-deleted plain-type document IDs and their mtime.
func (c *LiveSyncClient) ListNotes(ctx context.Context) ([]NoteEntry, error)

// ReadNote fetches a note by its CouchDB _id (which is the file path).
func (c *LiveSyncClient) ReadNote(ctx context.Context, path string) (NoteDoc, error)

// WriteNote creates or updates a note. Fetches current _rev first if updating.
func (c *LiveSyncClient) WriteNote(ctx context.Context, path, content string) error

// DeleteNote marks a note as deleted (sets deleted:true, updates rev).
func (c *LiveSyncClient) DeleteNote(ctx context.Context, path string) error

// MoveNote reads from oldPath, writes to newPath, then deletes oldPath.
func (c *LiveSyncClient) MoveNote(ctx context.Context, oldPath, newPath string) error

// ListFolders returns unique folder prefixes derived from all note paths.
func (c *LiveSyncClient) ListFolders(ctx context.Context) ([]string, error)

// ListTags scans all note content for YAML frontmatter `tags:` fields and returns unique tags.
func (c *LiveSyncClient) ListTags(ctx context.Context) ([]string, error)

// GetNoteMetadata returns NoteDoc with content omitted (just metadata fields).
func (c *LiveSyncClient) GetNoteMetadata(ctx context.Context, path string) (NoteDoc, error)
```

Supporting types:

```go
type NoteEntry struct {
    Path  string
    Mtime int64
}

type NoteDoc struct {
    Id      string  // CouchDB _id = file path
    Rev     string  // CouchDB _rev
    Content string  // decoded from base64 "data" field
    Mtime   int64
    Ctime   int64
    Size    int64
    Deleted bool
}
```

## Implementation Details

- Use `encoding/json` for request/response bodies.
- Use `encoding/base64` (StdEncoding) to encode/decode the `data` field.
- For `ListNotes`: call `GET /{db}/_all_docs?include_docs=true`, filter rows where `doc.type == "plain"` and `doc.deleted != true`, and exclude rows whose `_id` starts with `_design/`.
- For `WriteNote`: first try `GET /{db}/{encoded_id}` to get `_rev`. If 404, create new. Then `PUT /{db}/{encoded_id}` with the document body.
- URL-encode the `_id` path using `url.PathEscape`.
- Set `Content-Type: application/json` and `Accept: application/json` on all requests.
- Use Basic Auth: `req.SetBasicAuth(username, password)`.
- Return `rerrors.Wrap(err, "context")` from `go.redsock.ru/rerrors` for all errors.

## Coding Rules

- Never check errors inline; always `err := f()` then `if err != nil`.
- Never create struct literals inline in a function call.
- No all-caps field names.

## Verification

- `go build ./...` must pass.
