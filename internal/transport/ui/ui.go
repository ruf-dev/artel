package ui

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed all:dist
var distFS embed.FS

func NewHandler() http.Handler {
	sub, err := fs.Sub(distFS, "dist")
	if err != nil {
		panic(err)
	}
	return &spaHandler{fs: sub}
}

type spaHandler struct {
	fs fs.FS
}

func (h *spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")

	if path != "" {
		f, err := h.fs.Open(path)
		if err == nil {
			f.Close()
			http.FileServer(http.FS(h.fs)).ServeHTTP(w, r)
			return
		}
	}

	// index.html is served as a fallback for every unmatched path (e.g. a
	// reload on /mcp_keys), so every such path becomes its own browser cache
	// key for this response. Without an explicit no-store, a stale cached
	// copy keeps referencing asset filenames from a previous build (Vite
	// content-hashes them) that no longer exist after a redeploy, breaking
	// the page on reload while "/" - which gets revalidated far more often -
	// keeps working.
	w.Header().Set("Cache-Control", "no-store")
	http.ServeFileFS(w, r, h.fs, "index.html")
}
