package couchdb

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	kivik "github.com/go-kivik/kivik/v4"
	kivikcouch "github.com/go-kivik/kivik/v4/couchdb"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/service/user_errors"
)

type LiveSyncClient struct {
	db *kivik.DB
}

func NewLiveSyncClient(baseURL, dbName, username, password string) *LiveSyncClient {
	httpClient := &http.Client{Transport: newLoggingTransport()}

	client, err := kivik.New("couch", baseURL,
		kivikcouch.BasicAuth(username, password),
		kivikcouch.OptionHTTPClient(httpClient),
	)
	if err != nil {
		panic("couchdb: failed to create kivik client: " + err.Error())
	}

	return &LiveSyncClient{db: client.DB(dbName)}
}

func (c *LiveSyncClient) ListNotes(ctx context.Context) ([]NoteEntry, error) {
	rows := c.db.AllDocs(ctx, kivik.Params(map[string]any{"include_docs": true}))
	defer rows.Close()

	var notes []NoteEntry
	for rows.Next() {
		id, err := rows.ID()
		if err != nil {
			continue
		}
		if strings.HasPrefix(id, "_design/") {
			continue
		}
		var doc struct {
			Type    string `json:"type"`
			Deleted bool   `json:"deleted"`
			Mtime   int64  `json:"mtime"`
		}
		err = rows.ScanDoc(&doc)
		if err != nil {
			continue
		}
		if doc.Type != "plain" || doc.Deleted {
			continue
		}
		notes = append(notes, NoteEntry{Path: id, Mtime: doc.Mtime})
	}
	err := rows.Err()
	if err != nil {
		return nil, rerrors.Wrap(err, "rows iteration")
	}
	return notes, nil
}

func (c *LiveSyncClient) ReadNote(ctx context.Context, path string) (NoteDoc, error) {
	d := c.db.Get(ctx, path)
	defer d.Close()

	var couchDoc struct {
		Id       string   `json:"_id"`
		Rev      string   `json:"_rev"`
		Data     string   `json:"data"`
		Children []string `json:"children"`
		Mtime    int64    `json:"mtime"`
		Ctime    int64    `json:"ctime"`
		Size     int64    `json:"size"`
		Deleted  bool     `json:"deleted"`
	}
	err := d.ScanDoc(&couchDoc)
	if err != nil {
		return NoteDoc{}, rerrors.Wrap(err, "scan doc")
	}

	var content string
	if len(couchDoc.Children) > 0 {
		var sb strings.Builder
		for _, hash := range couchDoc.Children {
			chunk, err := c.fetchLeaf(ctx, hash)
			if err != nil {
				return NoteDoc{}, rerrors.Wrap(err, fmt.Sprintf("fetch leaf %s", hash))
			}
			sb.WriteString(chunk)
		}
		content = sb.String()
	} else {
		decoded, err := base64.StdEncoding.DecodeString(couchDoc.Data)
		if err != nil {
			return NoteDoc{}, rerrors.Wrap(err, "decode base64")
		}
		content = string(decoded)
	}

	return NoteDoc{
		Id:      couchDoc.Id,
		Rev:     couchDoc.Rev,
		Content: content,
		Mtime:   couchDoc.Mtime,
		Ctime:   couchDoc.Ctime,
		Size:    couchDoc.Size,
		Deleted: couchDoc.Deleted,
	}, nil
}

func (c *LiveSyncClient) WriteNote(ctx context.Context, path, content string) error {
	d := c.db.Get(ctx, path)
	defer d.Close()

	var existing struct {
		Rev string `json:"_rev"`
	}
	rev := ""
	err := d.ScanDoc(&existing)
	if err != nil {
		if kivik.HTTPStatus(err) != http.StatusNotFound {
			return rerrors.Wrap(err, "get existing doc")
		}
	} else {
		rev = existing.Rev
	}

	now := time.Now().UnixMilli()
	doc := struct {
		Id    string `json:"_id"`
		Rev   string `json:"_rev,omitempty"`
		Data  string `json:"data"`
		Mtime int64  `json:"mtime"`
		Ctime int64  `json:"ctime"`
		Size  int64  `json:"size"`
		Type  string `json:"type"`
	}{
		Id:    path,
		Rev:   rev,
		Data:  base64.StdEncoding.EncodeToString([]byte(content)),
		Mtime: now,
		Ctime: now,
		Size:  int64(len(content)),
		Type:  "plain",
	}

	_, err = c.db.Put(ctx, path, doc)
	if err != nil {
		return rerrors.Wrap(err, "put doc")
	}
	return nil
}

func (c *LiveSyncClient) DeleteNote(ctx context.Context, path string) error {
	note, err := c.ReadNote(ctx, path)
	if err != nil {
		return rerrors.Wrap(err, "read note")
	}

	doc := struct {
		Id      string `json:"_id"`
		Rev     string `json:"_rev"`
		Type    string `json:"type"`
		Deleted bool   `json:"deleted"`
	}{
		Id:      path,
		Rev:     note.Rev,
		Type:    "plain",
		Deleted: true,
	}

	_, err = c.db.Put(ctx, path, doc)
	if err != nil {
		return rerrors.Wrap(err, "put delete doc")
	}
	return nil
}

func (c *LiveSyncClient) MoveNote(ctx context.Context, oldPath, newPath string) error {
	note, err := c.ReadNote(ctx, oldPath)
	if err != nil {
		return rerrors.Wrap(err, "read note from old path")
	}

	err = c.WriteNote(ctx, newPath, note.Content)
	if err != nil {
		return rerrors.Wrap(err, "write note to new path")
	}

	err = c.DeleteNote(ctx, oldPath)
	if err != nil {
		return rerrors.Wrap(err, "delete note from old path")
	}
	return nil
}

func (c *LiveSyncClient) ListFolders(ctx context.Context) ([]string, error) {
	notes, err := c.ListNotes(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "list notes for folders")
	}

	folderMap := make(map[string]bool)
	for _, note := range notes {
		dir := filepath.Dir(note.Path)
		if dir != "." {
			folderMap[dir] = true
		}
	}

	folders := make([]string, 0, len(folderMap))
	for folder := range folderMap {
		folders = append(folders, folder)
	}
	return folders, nil
}

func (c *LiveSyncClient) ListTags(ctx context.Context) ([]string, error) {
	notes, err := c.ListNotes(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "list notes for tags")
	}

	tagMap := make(map[string]bool)
	for _, note := range notes {
		noteDoc, err := c.ReadNote(ctx, note.Path)
		if err != nil {
			continue
		}
		for _, tag := range extractTags(noteDoc.Content) {
			tagMap[tag] = true
		}
	}

	tags := make([]string, 0, len(tagMap))
	for tag := range tagMap {
		tags = append(tags, tag)
	}
	return tags, nil
}

func (c *LiveSyncClient) GetNoteMetadata(ctx context.Context, path string) (NoteDoc, error) {
	note, err := c.ReadNote(ctx, path)
	if err != nil {
		return NoteDoc{}, rerrors.Wrap(err, "read note")
	}
	note.Content = ""
	return note, nil
}

func (c *LiveSyncClient) fetchLeaf(ctx context.Context, hash string) (string, error) {
	d := c.db.Get(ctx, hash)
	defer d.Close()

	var leafDoc struct {
		Data string `json:"data"`
	}
	err := d.ScanDoc(&leafDoc)
	if err != nil {
		return "", rerrors.Wrap(err, "scan leaf")
	}
	return leafDoc.Data, nil
}

func (c *LiveSyncClient) getDocRev(ctx context.Context, path string) (string, error) {
	d := c.db.Get(ctx, path)
	defer d.Close()

	var doc struct {
		Rev string `json:"_rev"`
	}
	err := d.ScanDoc(&doc)
	if err != nil {
		return "", rerrors.Wrap(err, "get doc rev")
	}
	return doc.Rev, nil
}

func (c *LiveSyncClient) ListFiles(ctx context.Context) ([]FileEntry, error) {
	rows := c.db.AllDocs(ctx, kivik.Params(map[string]any{"include_docs": true}))
	defer rows.Close()

	var files []FileEntry
	for rows.Next() {
		id, err := rows.ID()
		if err != nil {
			continue
		}
		if strings.HasPrefix(id, "_design/") {
			continue
		}
		var doc struct {
			Type    string `json:"type"`
			Deleted bool   `json:"deleted"`
			Mtime   int64  `json:"mtime"`
		}
		err = rows.ScanDoc(&doc)
		if err != nil {
			continue
		}
		if (doc.Type != "plain" && doc.Type != "newnote") || doc.Deleted {
			continue
		}
		mime := MimeTypeForPath(id)
		if strings.HasPrefix(mime, "text/") {
			continue
		}
		files = append(files, FileEntry{Path: id, Mtime: doc.Mtime, MimeType: mime})
	}
	err := rows.Err()
	if err != nil {
		return nil, rerrors.Wrap(err, "rows iteration")
	}
	return files, nil
}

func (c *LiveSyncClient) ReadFile(ctx context.Context, path string) (FileDoc, error) {
	mime := MimeTypeForPath(path)
	d := c.db.Get(ctx, path)
	defer d.Close()

	var couchDoc struct {
		Id       string   `json:"_id"`
		Rev      string   `json:"_rev"`
		Type     string   `json:"type"`
		Data     string   `json:"data"`
		Children []string `json:"children"`
		Mtime    int64    `json:"mtime"`
		Ctime    int64    `json:"ctime"`
		Size     int64    `json:"size"`
		Deleted  bool     `json:"deleted"`
	}
	err := d.ScanDoc(&couchDoc)
	if err != nil {
		return FileDoc{}, rerrors.Wrap(err, "scan doc")
	}

	var rawBytes []byte
	if couchDoc.Type == "newnote" {
		for _, hash := range couchDoc.Children {
			raw, err := c.fetchLeaf(ctx, hash)
			if err != nil {
				return FileDoc{}, rerrors.Wrap(err, fmt.Sprintf("fetch leaf %s", hash))
			}
			chunk, err := base64.StdEncoding.DecodeString(raw)
			if err != nil {
				return FileDoc{}, rerrors.Wrap(err, fmt.Sprintf("decode leaf %s", hash))
			}
			rawBytes = append(rawBytes, chunk...)
		}
	} else {
		decoded, err := base64.StdEncoding.DecodeString(couchDoc.Data)
		if err != nil {
			return FileDoc{}, rerrors.Wrap(err, "decode base64")
		}
		rawBytes = decoded
	}

	return FileDoc{
		Id:       couchDoc.Id,
		Rev:      couchDoc.Rev,
		RawBytes: rawBytes,
		MimeType: mime,
		Mtime:    couchDoc.Mtime,
		Ctime:    couchDoc.Ctime,
		Size:     couchDoc.Size,
		Deleted:  couchDoc.Deleted,
	}, nil
}

func (c *LiveSyncClient) DeleteFile(ctx context.Context, path string) error {
	rev, err := c.getDocRev(ctx, path)
	if err != nil {
		return rerrors.Wrap(err, "get doc rev")
	}

	doc := struct {
		Id      string `json:"_id"`
		Rev     string `json:"_rev"`
		Type    string `json:"type"`
		Deleted bool   `json:"deleted"`
	}{
		Id:      path,
		Rev:     rev,
		Type:    "plain",
		Deleted: true,
	}

	_, err = c.db.Put(ctx, path, doc)
	if err != nil {
		return rerrors.Wrap(err, "put delete doc")
	}
	return nil
}

func (c *LiveSyncClient) MoveFile(ctx context.Context, oldPath, newPath string) error {
	d := c.db.Get(ctx, oldPath)
	defer d.Close()

	var couchDoc struct {
		Type string `json:"type"`
	}
	err := d.ScanDoc(&couchDoc)
	if err != nil {
		return rerrors.Wrap(err, "get existing doc type")
	}

	if couchDoc.Type == "newnote" {
		return user_errors.ChunkedBinaryMoveNotSupported
	}

	file, err := c.ReadFile(ctx, oldPath)
	if err != nil {
		return rerrors.Wrap(err, "read file")
	}

	newDoc := struct {
		Id    string `json:"_id"`
		Data  string `json:"data"`
		Mtime int64  `json:"mtime"`
		Ctime int64  `json:"ctime"`
		Size  int64  `json:"size"`
		Type  string `json:"type"`
	}{
		Id:    newPath,
		Data:  base64.StdEncoding.EncodeToString(file.RawBytes),
		Mtime: file.Mtime,
		Ctime: file.Ctime,
		Size:  file.Size,
		Type:  "plain",
	}

	_, err = c.db.Put(ctx, newPath, newDoc)
	if err != nil {
		return rerrors.Wrap(err, "put new doc")
	}

	err = c.DeleteFile(ctx, oldPath)
	if err != nil {
		return rerrors.Wrap(err, "delete old file")
	}
	return nil
}

func extractTags(content string) []string {
	lines := strings.Split(content, "\n")

	if len(lines) == 0 || lines[0] != "---" {
		return nil
	}

	var frontmatter string

	for i := 1; i < len(lines); i++ {
		if lines[i] == "---" {
			frontmatter = strings.Join(lines[1:i], "\n")
			break
		}
	}

	if frontmatter == "" {
		return nil
	}

	tagsRe := regexp.MustCompile(`(?m)^\s*tags:\s*\[([^\]]+)\]|^\s*tags:\s*$`)

	match := tagsRe.FindStringSubmatch(frontmatter)

	if len(match) > 1 && match[1] != "" {
		tagStr := match[1]
		parts := strings.Split(tagStr, ",")

		var tags []string

		for _, part := range parts {
			tag := strings.TrimSpace(part)
			tag = strings.Trim(tag, "\"'")

			if tag != "" {
				tags = append(tags, tag)
			}
		}

		return tags
	}

	lines = strings.Split(frontmatter, "\n")

	var tags []string

	for i, line := range lines {
		if strings.Contains(line, "tags:") {
			for j := i + 1; j < len(lines); j++ {
				if !strings.HasPrefix(lines[j], "  ") && !strings.HasPrefix(lines[j], "-") {
					break
				}

				trimmed := strings.TrimSpace(lines[j])

				if strings.HasPrefix(trimmed, "- ") {
					tag := strings.TrimPrefix(trimmed, "- ")
					tag = strings.Trim(tag, "\"'")

					if tag != "" {
						tags = append(tags, tag)
					}
				}
			}

			break
		}
	}

	return tags
}

type NoteEntry struct {
	Path  string
	Mtime int64
}

type NoteDoc struct {
	Id      string
	Rev     string
	Content string
	Mtime   int64
	Ctime   int64
	Size    int64
	Deleted bool
}

type FileEntry struct {
	Path     string
	Mtime    int64
	MimeType string
}

type FileDoc struct {
	Id       string
	Rev      string
	RawBytes []byte
	MimeType string
	Mtime    int64
	Ctime    int64
	Size     int64
	Deleted  bool
}

func MimeTypeForPath(path string) string {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".md", ".markdown":
		return "text/markdown"
	case ".txt", ".csv", ".yaml", ".yml", ".json", ".xml", ".html", ".htm", ".css", ".js", ".ts":
		return "text/plain"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	case ".pdf":
		return "application/pdf"
	default:
		return "application/octet-stream"
	}
}
