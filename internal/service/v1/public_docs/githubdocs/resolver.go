package githubdocs

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/ruf-dev/artel/internal/clients/couchdb"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

// cacheTTL bounds how long a walked tree is served from cache before the next call re-walks
// GitHub. Long enough to comfortably stay under GitHub's ~60 req/hour unauthenticated rate
// limit even under repeated /docs page loads; short enough that content edits show up promptly.
const cacheTTL = 5 * time.Minute

// Resolver serves the docs/public tree of this repo's own GitHub remote in the same shape as
// the CouchDB-backed public docs path. It lazy-fetches on first use (not in NewResolver) and
// caches the whole walked tree — files, folders, and content — as one unit for cacheTTL.
type Resolver struct {
	client *client

	mu       sync.Mutex
	cachedAt time.Time
	notes    []couchdb.NoteEntry
	folders  []string
	content  map[string]string // note path -> raw file content, populated during the same walk
}

// NewResolver constructs a Resolver ready to serve requests. It does not perform any network
// calls itself — the first ListFolders/ListNotes/GetNote/ListTags call triggers the initial walk.
func NewResolver() *Resolver {
	r := &Resolver{
		client: newClient(),
	}

	return r
}

// refresh re-walks GitHub if the cache is empty or older than cacheTTL. It holds mu for the
// duration of any HTTP calls — acceptable given the low traffic this path sees and the short
// TTL, and it avoids a thundering-herd duplicate walk race without a second synchronization
// primitive.
func (r *Resolver) refresh(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.cachedAt.IsZero() && time.Since(r.cachedAt) < cacheTTL {
		return nil
	}

	files, folders, err := walkTree(ctx, r.client, root)
	if err != nil {
		return rerrors.Wrap(err, "error walking github docs tree")
	}

	notes := make([]couchdb.NoteEntry, 0, len(files))
	content := make(map[string]string, len(files))

	// Mtime uses the walk's own fetch time for every note — GitHub's contents API doesn't
	// cheaply expose a per-file last-modified time (that needs one extra rate-limited
	// commits-API call per file), and nothing here sorts/relies on a real per-file mtime.
	now := time.Now()

	for _, f := range files {
		notePath := relPath(f.Path)

		var raw string

		raw, err = r.client.fetchRaw(ctx, f.Path)
		if err != nil {
			return rerrors.Wrap(err, "error fetching github raw content for "+f.Path)
		}

		note := couchdb.NoteEntry{
			Path:  notePath,
			Mtime: now.UnixMilli(),
		}

		notes = append(notes, note)
		content[notePath] = raw
	}

	relFolders := make([]string, 0, len(folders))
	for _, f := range folders {
		relFolders = append(relFolders, relPath(f))
	}

	r.cachedAt = now
	r.notes = notes
	r.folders = relFolders
	r.content = content

	return nil
}

// relPath strips the hardcoded root prefix from a repo-root-relative GitHub path, matching the
// leading/trailing-slash-free convention couchdb.LiveSyncClient.ListNotes/ListFolders use for
// NoteEntry.Path (the raw CouchDB document _id) and its derived folder strings.
func relPath(path string) string {
	rel := strings.TrimPrefix(path, root)
	rel = strings.TrimPrefix(rel, "/")

	return rel
}

func (r *Resolver) ListFolders(ctx context.Context) ([]string, error) {
	err := r.refresh(ctx)
	if err != nil {
		return nil, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	folders := make([]string, len(r.folders))
	copy(folders, r.folders)

	return folders, nil
}

func (r *Resolver) ListNotes(ctx context.Context) ([]couchdb.NoteEntry, error) {
	err := r.refresh(ctx)
	if err != nil {
		return nil, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	notes := make([]couchdb.NoteEntry, len(r.notes))
	copy(notes, r.notes)

	return notes, nil
}

func (r *Resolver) GetNote(ctx context.Context, path string) (couchdb.NoteDoc, error) {
	err := r.refresh(ctx)
	if err != nil {
		return couchdb.NoteDoc{}, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	raw, ok := r.content[path]
	if !ok {
		return couchdb.NoteDoc{}, rerrors.Wrap(user_errors.NotFound)
	}

	note := couchdb.NoteDoc{
		Content: raw,
		Mtime:   r.cachedAt.UnixMilli(),
		Size:    int64(len(raw)),
	}

	return note, nil
}

// ListTags always returns an empty slice: the plain markdown files served here have no
// Obsidian-frontmatter tag convention to parse, and the frontend doesn't call this RPC for the
// docs page today either way.
func (r *Resolver) ListTags(_ context.Context) ([]string, error) {
	return []string{}, nil
}
