package githubdocs

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func newTestClient(contentsBase, rawBase string) *client {
	httpClient := &http.Client{}

	c := &client{
		httpClient:   httpClient,
		contentsBase: contentsBase,
		rawBase:      rawBase,
	}

	return c
}

func TestClient_ListDir_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := fmt.Sprintf("/repos/%s/%s/contents/%s", owner, repo, root)
		require.Equal(t, expectedPath, r.URL.Path)
		require.Equal(t, ref, r.URL.Query().Get("ref"))

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`[
			{"name":"quick-start.md","path":"docs/public/quick-start.md","type":"file","size":42},
			{"name":"guides","path":"docs/public/guides","type":"dir","size":0}
		]`))
		require.NoError(t, err)
	}))
	defer server.Close()

	c := newTestClient(server.URL, server.URL)

	entries, err := c.listDir(context.Background(), root)
	require.NoError(t, err)
	require.Len(t, entries, 2)
	require.Equal(t, "docs/public/quick-start.md", entries[0].Path)
	require.Equal(t, "file", entries[0].Type)
	require.Equal(t, "docs/public/guides", entries[1].Path)
	require.Equal(t, "dir", entries[1].Type)
}

func TestClient_ListDir_NotFoundIsEmptyNotError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	c := newTestClient(server.URL, server.URL)

	entries, err := c.listDir(context.Background(), root)
	require.NoError(t, err)
	require.Empty(t, entries)
}

func TestClient_ListDir_ServerErrorIsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	c := newTestClient(server.URL, server.URL)

	_, err := c.listDir(context.Background(), root)
	require.Error(t, err)
}

func TestClient_FetchRaw_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := fmt.Sprintf("/%s/%s/%s/docs/public/quick-start.md", owner, repo, ref)
		require.Equal(t, expectedPath, r.URL.Path)

		_, err := w.Write([]byte("# Quick Start\n\nHello."))
		require.NoError(t, err)
	}))
	defer server.Close()

	c := newTestClient(server.URL, server.URL)

	content, err := c.fetchRaw(context.Background(), "docs/public/quick-start.md")
	require.NoError(t, err)
	require.Equal(t, "# Quick Start\n\nHello.", content)
}

func TestClient_FetchRaw_NotFoundIsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	c := newTestClient(server.URL, server.URL)

	_, err := c.fetchRaw(context.Background(), "docs/public/missing.md")
	require.Error(t, err)
}
