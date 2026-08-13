package public_docs

import (
	"context"

	"github.com/ruf-dev/artel/internal/clients/couchdb"
)

// DocsResolver abstracts "browse a folder/note tree" over either a published CouchDB vault or
// the GitHub-backed quick-start guide (internal/service/v1/public_docs/githubdocs), letting
// Service.resolveDocsSource branch once instead of duplicating CouchDB-vs-GitHub logic across
// every method below.
type DocsResolver interface {
	ListFolders(ctx context.Context) ([]string, error)
	ListNotes(ctx context.Context) ([]couchdb.NoteEntry, error)
	GetNote(ctx context.Context, path string) (couchdb.NoteDoc, error)
	ListTags(ctx context.Context) ([]string, error)
}

// couchResolver adapts an already-authenticated CouchDB LiveSyncClient to DocsResolver.
type couchResolver struct {
	client *couchdb.LiveSyncClient
}

func (c couchResolver) ListFolders(ctx context.Context) ([]string, error) {
	return c.client.ListFolders(ctx)
}

func (c couchResolver) ListNotes(ctx context.Context) ([]couchdb.NoteEntry, error) {
	return c.client.ListNotes(ctx)
}

func (c couchResolver) GetNote(ctx context.Context, path string) (couchdb.NoteDoc, error) {
	return c.client.ReadNote(ctx, path)
}

func (c couchResolver) ListTags(ctx context.Context) ([]string, error) {
	return c.client.ListTags(ctx)
}
