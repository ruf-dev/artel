package githubdocs

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ruf-dev/artel/internal/service/user_errors"
	"github.com/stretchr/testify/require"
)

// newFakeGithubServer builds a single httptest.Server that fakes both the contents API and the
// raw content endpoint for a small two-level fake tree:
//
//	docs/public/quick-start.md
//	docs/public/guides/ (dir)
//	docs/public/guides/setup.md
//
// requestCount is incremented once per incoming HTTP request, letting tests assert on cache
// behavior (how many times the fake server was actually hit).
func newFakeGithubServer(t *testing.T, requestCount *int64) *httptest.Server {
	t.Helper()

	contentsPrefix := fmt.Sprintf("/repos/%s/%s/contents/", owner, repo)
	rawPrefix := fmt.Sprintf("/%s/%s/%s/", owner, repo, ref)

	handler := func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(requestCount, 1)

		switch r.URL.Path {
		case contentsPrefix + "docs/public":
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write([]byte(`[
				{"name":"quick-start.md","path":"docs/public/quick-start.md","type":"file","size":10},
				{"name":"guides","path":"docs/public/guides","type":"dir","size":0}
			]`))
			require.NoError(t, err)

		case contentsPrefix + "docs/public/guides":
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write([]byte(`[
				{"name":"setup.md","path":"docs/public/guides/setup.md","type":"file","size":10}
			]`))
			require.NoError(t, err)

		case rawPrefix + "docs/public/quick-start.md":
			_, err := w.Write([]byte("# Quick Start"))
			require.NoError(t, err)

		case rawPrefix + "docs/public/guides/setup.md":
			_, err := w.Write([]byte("# Setup"))
			require.NoError(t, err)

		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	t.Cleanup(server.Close)

	return server
}

func newTestResolver(server *httptest.Server) *Resolver {
	c := newTestClient(server.URL, server.URL)

	r := &Resolver{
		client: c,
	}

	return r
}

func TestResolver_ListNotesAndFolders(t *testing.T) {
	var requestCount int64

	server := newFakeGithubServer(t, &requestCount)
	r := newTestResolver(server)

	notes, err := r.ListNotes(context.Background())
	require.NoError(t, err)
	require.Len(t, notes, 2)

	paths := make([]string, 0, len(notes))
	for _, n := range notes {
		paths = append(paths, n.Path)
	}
	require.ElementsMatch(t, []string{"quick-start.md", "guides/setup.md"}, paths)

	folders, err := r.ListFolders(context.Background())
	require.NoError(t, err)
	require.Equal(t, []string{"guides"}, folders)
}

func TestResolver_GetNote(t *testing.T) {
	var requestCount int64

	server := newFakeGithubServer(t, &requestCount)
	r := newTestResolver(server)

	note, err := r.GetNote(context.Background(), "quick-start.md")
	require.NoError(t, err)
	require.Equal(t, "# Quick Start", note.Content)
	require.Equal(t, int64(len("# Quick Start")), note.Size)

	_, err = r.GetNote(context.Background(), "does-not-exist.md")
	require.Error(t, err)
	require.True(t, errors.Is(err, user_errors.NotFound), "expected user_errors.NotFound, got %v", err)
}

func TestResolver_ListTags(t *testing.T) {
	var requestCount int64

	server := newFakeGithubServer(t, &requestCount)
	r := newTestResolver(server)

	tags, err := r.ListTags(context.Background())
	require.NoError(t, err)
	require.Equal(t, []string{}, tags)
}

func TestResolver_CachesWithinTTL(t *testing.T) {
	var requestCount int64

	server := newFakeGithubServer(t, &requestCount)
	r := newTestResolver(server)

	_, err := r.ListNotes(context.Background())
	require.NoError(t, err)

	firstCount := atomic.LoadInt64(&requestCount)
	require.Positive(t, firstCount)

	_, err = r.ListNotes(context.Background())
	require.NoError(t, err)

	require.Equal(t, firstCount, atomic.LoadInt64(&requestCount), "second call within TTL must not re-hit the server")

	// Rewind cachedAt past cacheTTL to simulate expiry, then confirm a third call re-hits the
	// server.
	r.mu.Lock()
	r.cachedAt = time.Now().Add(-cacheTTL - time.Second)
	r.mu.Unlock()

	_, err = r.ListNotes(context.Background())
	require.NoError(t, err)

	require.Greater(t, atomic.LoadInt64(&requestCount), firstCount, "call after TTL expiry must re-hit the server")
}

func TestRelPath(t *testing.T) {
	require.Equal(t, "quick-start.md", relPath("docs/public/quick-start.md"))
	require.Equal(t, "guides", relPath("docs/public/guides"))
	require.True(t, !strings.HasPrefix(relPath("docs/public/quick-start.md"), "/"))
}
