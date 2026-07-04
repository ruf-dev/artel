package storage

import (
	"path/filepath"
	"strings"
)

// IsMarkdown is the strict, extension-only routing decision between the two backends:
// true -> CouchDB note path, false -> S3. This is deliberately narrower than
// couchdb.MimeTypeForPath's "text/*" bucket (which also covers .txt/.csv/.yaml/.json/etc) —
// only .md/.markdown stay in CouchDB; everything else routes to S3.
func IsMarkdown(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".md", ".markdown":
		return true
	default:
		return false
	}
}
