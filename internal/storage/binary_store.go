package storage

import "context"

type ObjectEntry struct {
	Path     string
	Size     int64
	MimeType string
	Mtime    int64 // unix ms
}

type Object struct {
	Path     string
	RawBytes []byte
	MimeType string
	Size     int64
}

// BinaryStore is a CRUD-over-object-keys abstraction backing the non-markdown half of a
// vault's content. There is no CouchDB implementation of this interface — CouchDB's own binary
// methods (couchdb.LiveSyncClient's ReadFile/DeleteFile/MoveFile/ListFiles) stay as direct,
// unwrapped calls and are no longer reachable from MCP tools once the S3 routing lands (a
// confirmed, intentional behavior change owned by a different task, not yours).
type BinaryStore interface {
	Put(ctx context.Context, path string, content []byte, mimeType string) error
	Get(ctx context.Context, path string) (Object, error)
	Stat(ctx context.Context, path string) (ObjectEntry, error)
	Delete(ctx context.Context, path string) error
	Move(ctx context.Context, oldPath, newPath string) error
	List(ctx context.Context) ([]ObjectEntry, error)
}
