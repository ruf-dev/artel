// Package githubdocs serves the built-in "Artel Quick Start" docs tree straight from this
// repo's own public GitHub remote (docs/public/ on the master branch), so a fresh Artel instance
// has something real to point an anonymous /docs visitor at before any CouchDB vault has been
// published. See internal/domain.ReservedGithubDocsSlug and internal/service/v1/public_docs for
// how this plugs into the existing public-docs surface — this package is self-contained and
// knows nothing about that wiring.
package githubdocs

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ruf-dev/artel/internal/utils"
	"go.redsock.ru/rerrors"
)

const (
	owner = "ruf-dev"
	repo  = "artel"
	ref   = "master"
	root  = "docs/public"

	contentsAPIBase = "https://api.github.com"
	rawContentBase  = "https://raw.githubusercontent.com"

	httpTimeout = 10 * time.Second
)

// contentsEntry mirrors the subset of GitHub's "get repository content" API response this
// package reads: https://docs.github.com/en/rest/repos/contents#get-repository-content
type contentsEntry struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"` // "file" or "dir"
	Size int64  `json:"size"`
}

// client is an unauthenticated GitHub REST client. No token is used — this is a low-traffic,
// public-repo read path, and unauthenticated requests are capped at ~60 req/hour per GitHub's
// rate limit, which Resolver's 5-minute cache TTL comfortably stays under (one walk of the small
// fixed docs/public tree per cache refresh, not per request).
//
// contentsBase/rawBase default to the real GitHub hosts but are overridable so tests can point
// this at an httptest.Server instead of hitting the network.
type client struct {
	httpClient   *http.Client
	contentsBase string
	rawBase      string
}

func newClient() *client {
	httpClient := &http.Client{
		Timeout: httpTimeout,
	}

	c := &client{
		httpClient:   httpClient,
		contentsBase: contentsAPIBase,
		rawBase:      rawContentBase,
	}

	return c
}

// listDir lists the entries of a single directory in the repo tree. A 404 (the path doesn't
// exist yet — e.g. before docs/public/ has been authored/merged) is treated as a normal empty
// state, not an error: it degrades to "no notes" rather than failing the whole /docs page.
func (c *client) listDir(ctx context.Context, path string) ([]contentsEntry, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/contents/%s?ref=%s", c.contentsBase, owner, repo, path, ref)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, rerrors.Wrap(err, "error building github contents request")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, rerrors.Wrap(err, "error calling github contents api")
	}
	defer utils.CloseWithLog(resp.Body, "error closing github contents response body")

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, rerrors.New(fmt.Sprintf("unexpected github contents api status %d for path %q", resp.StatusCode, path))
	}

	var entries []contentsEntry

	err = json.NewDecoder(resp.Body).Decode(&entries)
	if err != nil {
		return nil, rerrors.Wrap(err, "error decoding github contents response")
	}

	return entries, nil
}

// fetchRaw fetches the raw content of a single file at path via raw.githubusercontent.com.
func (c *client) fetchRaw(ctx context.Context, path string) (string, error) {
	url := fmt.Sprintf("%s/%s/%s/%s/%s", c.rawBase, owner, repo, ref, path)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", rerrors.Wrap(err, "error building github raw content request")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", rerrors.Wrap(err, "error calling github raw content endpoint")
	}
	defer utils.CloseWithLog(resp.Body, "error closing github raw content response body")

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", rerrors.New(fmt.Sprintf("unexpected github raw content status %d for path %q", resp.StatusCode, path))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", rerrors.Wrap(err, "error reading github raw content response body")
	}

	return string(body), nil
}
