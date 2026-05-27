package couchdb

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/service/user_errors"
)

type LiveSyncClient struct {
	baseURL  string
	dbName   string
	username string
	password string
	http     *http.Client
}

func NewLiveSyncClient(baseURL, dbName, username, password string) *LiveSyncClient {
	return &LiveSyncClient{
		baseURL:  baseURL,
		dbName:   dbName,
		username: username,
		password: password,
		http:     &http.Client{},
	}
}

func (c *LiveSyncClient) ListNotes(ctx context.Context) ([]NoteEntry, error) {
	url := fmt.Sprintf("%s/%s/_all_docs?include_docs=true", c.baseURL, c.dbName)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, rerrors.Wrap(err, "failed to create request")
	}

	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, rerrors.Wrap(err, "failed to execute request")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, rerrors.Wrap(err, "failed to read response body")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, rerrors.New(fmt.Sprintf("unexpected status %d: %s", resp.StatusCode, string(body)))
	}

	var allDocsResp struct {
		Rows []struct {
			Id  string `json:"id"`
			Doc struct {
				Type    string `json:"type"`
				Deleted bool   `json:"deleted"`
				Mtime   int64  `json:"mtime"`
			} `json:"doc"`
		} `json:"rows"`
	}

	err = json.Unmarshal(body, &allDocsResp)
	if err != nil {
		return nil, rerrors.Wrap(err, "failed to unmarshal response")
	}

	var notes []NoteEntry

	for _, row := range allDocsResp.Rows {
		if strings.HasPrefix(row.Id, "_design/") {
			continue
		}

		if row.Doc.Type != "plain" {
			continue
		}

		if row.Doc.Deleted {
			continue
		}

		entry := NoteEntry{
			Path:  row.Id,
			Mtime: row.Doc.Mtime,
		}

		notes = append(notes, entry)
	}

	return notes, nil
}

func (c *LiveSyncClient) ReadNote(ctx context.Context, path string) (NoteDoc, error) {
	encodedId := url.PathEscape(path)
	url := fmt.Sprintf("%s/%s/%s", c.baseURL, c.dbName, encodedId)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return NoteDoc{}, rerrors.Wrap(err, "failed to create request")
	}

	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return NoteDoc{}, rerrors.Wrap(err, "failed to execute request")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return NoteDoc{}, rerrors.Wrap(err, "failed to read response body")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return NoteDoc{}, rerrors.New(fmt.Sprintf("unexpected status %d: %s", resp.StatusCode, string(body)))
	}

	var couchDoc struct {
		Id      string `json:"_id"`
		Rev     string `json:"_rev"`
		Data    string `json:"data"`
		Mtime   int64  `json:"mtime"`
		Ctime   int64  `json:"ctime"`
		Size    int64  `json:"size"`
		Deleted bool   `json:"deleted"`
	}

	err = json.Unmarshal(body, &couchDoc)
	if err != nil {
		return NoteDoc{}, rerrors.Wrap(err, "failed to unmarshal response")
	}

	decoded, err := base64.StdEncoding.DecodeString(couchDoc.Data)
	if err != nil {
		return NoteDoc{}, rerrors.Wrap(err, "failed to decode base64 data")
	}

	result := NoteDoc{
		Id:      couchDoc.Id,
		Rev:     couchDoc.Rev,
		Content: string(decoded),
		Mtime:   couchDoc.Mtime,
		Ctime:   couchDoc.Ctime,
		Size:    couchDoc.Size,
		Deleted: couchDoc.Deleted,
	}

	return result, nil
}

func (c *LiveSyncClient) WriteNote(ctx context.Context, path, content string) error {
	encodedId := url.PathEscape(path)
	url := fmt.Sprintf("%s/%s/%s", c.baseURL, c.dbName, encodedId)

	var rev string

	getReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return rerrors.Wrap(err, "failed to create get request")
	}

	getReq.SetBasicAuth(c.username, c.password)
	getReq.Header.Set("Accept", "application/json")

	getResp, err := c.http.Do(getReq)
	if err != nil {
		return rerrors.Wrap(err, "failed to execute get request")
	}

	defer getResp.Body.Close()

	if getResp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(getResp.Body)
		if err != nil {
			return rerrors.Wrap(err, "failed to read get response body")
		}

		var existingDoc struct {
			Rev string `json:"_rev"`
		}

		err = json.Unmarshal(body, &existingDoc)
		if err != nil {
			return rerrors.Wrap(err, "failed to unmarshal get response")
		}

		rev = existingDoc.Rev
	} else if getResp.StatusCode != http.StatusNotFound {
		body, err := io.ReadAll(getResp.Body)
		if err != nil {
			return rerrors.Wrap(err, "failed to read get response body")
		}

		return rerrors.New(fmt.Sprintf("unexpected status %d: %s", getResp.StatusCode, string(body)))
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(content))

	doc := map[string]interface{}{
		"_id":   path,
		"data":  encoded,
		"mtime": 0,
		"ctime": 0,
		"size":  len(content),
		"type":  "plain",
	}

	if rev != "" {
		doc["_rev"] = rev
	}

	docBody, err := json.Marshal(doc)
	if err != nil {
		return rerrors.Wrap(err, "failed to marshal document")
	}

	putReq, err := http.NewRequestWithContext(ctx, http.MethodPut, url, strings.NewReader(string(docBody)))
	if err != nil {
		return rerrors.Wrap(err, "failed to create put request")
	}

	putReq.SetBasicAuth(c.username, c.password)
	putReq.Header.Set("Content-Type", "application/json")
	putReq.Header.Set("Accept", "application/json")

	putResp, err := c.http.Do(putReq)
	if err != nil {
		return rerrors.Wrap(err, "failed to execute put request")
	}

	defer putResp.Body.Close()

	body, err := io.ReadAll(putResp.Body)
	if err != nil {
		return rerrors.Wrap(err, "failed to read put response body")
	}

	if putResp.StatusCode < 200 || putResp.StatusCode >= 300 {
		return rerrors.New(fmt.Sprintf("unexpected status %d: %s", putResp.StatusCode, string(body)))
	}

	return nil
}

func (c *LiveSyncClient) DeleteNote(ctx context.Context, path string) error {
	note, err := c.ReadNote(ctx, path)
	if err != nil {
		return rerrors.Wrap(err, "failed to read note")
	}

	encodedId := url.PathEscape(path)
	url := fmt.Sprintf("%s/%s/%s", c.baseURL, c.dbName, encodedId)

	doc := map[string]interface{}{
		"_id":     path,
		"_rev":    note.Rev,
		"deleted": true,
	}

	docBody, err := json.Marshal(doc)
	if err != nil {
		return rerrors.Wrap(err, "failed to marshal document")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, strings.NewReader(string(docBody)))
	if err != nil {
		return rerrors.Wrap(err, "failed to create request")
	}

	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return rerrors.Wrap(err, "failed to execute request")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return rerrors.Wrap(err, "failed to read response body")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return rerrors.New(fmt.Sprintf("unexpected status %d: %s", resp.StatusCode, string(body)))
	}

	return nil
}

func (c *LiveSyncClient) MoveNote(ctx context.Context, oldPath, newPath string) error {
	note, err := c.ReadNote(ctx, oldPath)
	if err != nil {
		return rerrors.Wrap(err, "failed to read note from old path")
	}

	err = c.WriteNote(ctx, newPath, note.Content)
	if err != nil {
		return rerrors.Wrap(err, "failed to write note to new path")
	}

	err = c.DeleteNote(ctx, oldPath)
	if err != nil {
		return rerrors.Wrap(err, "failed to delete note from old path")
	}

	return nil
}

func (c *LiveSyncClient) ListFolders(ctx context.Context) ([]string, error) {
	notes, err := c.ListNotes(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "failed to list notes")
	}

	folderMap := make(map[string]bool)

	for _, note := range notes {
		dir := filepath.Dir(note.Path)

		if dir != "." {
			folderMap[dir] = true
		}
	}

	var folders []string

	for folder := range folderMap {
		folders = append(folders, folder)
	}

	return folders, nil
}

func (c *LiveSyncClient) ListTags(ctx context.Context) ([]string, error) {
	notes, err := c.ListNotes(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "failed to list notes")
	}

	tagMap := make(map[string]bool)

	for _, note := range notes {
		noteDoc, err := c.ReadNote(ctx, note.Path)
		if err != nil {
			continue
		}

		tags := extractTags(noteDoc.Content)

		for _, tag := range tags {
			tagMap[tag] = true
		}
	}

	var tags []string

	for tag := range tagMap {
		tags = append(tags, tag)
	}

	return tags, nil
}

func (c *LiveSyncClient) GetNoteMetadata(ctx context.Context, path string) (NoteDoc, error) {
	note, err := c.ReadNote(ctx, path)
	if err != nil {
		return NoteDoc{}, rerrors.Wrap(err, "failed to read note")
	}

	note.Content = ""

	return note, nil
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

func (c *LiveSyncClient) getDocRev(ctx context.Context, path string) (string, error) {
	encodedId := url.PathEscape(path)
	docURL := fmt.Sprintf("%s/%s/%s", c.baseURL, c.dbName, encodedId)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, docURL, nil)
	if err != nil {
		return "", rerrors.Wrap(err, "failed to create request")
	}
	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", rerrors.Wrap(err, "failed to execute request")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", rerrors.Wrap(err, "failed to read response body")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", rerrors.New(fmt.Sprintf("unexpected status %d: %s", resp.StatusCode, string(body)))
	}

	var doc struct {
		Rev string `json:"_rev"`
	}
	if err := json.Unmarshal(body, &doc); err != nil {
		return "", rerrors.Wrap(err, "failed to unmarshal response")
	}
	return doc.Rev, nil
}

func (c *LiveSyncClient) fetchLeaf(ctx context.Context, hash string) ([]byte, error) {
	encodedId := url.PathEscape(hash)
	leafURL := fmt.Sprintf("%s/%s/%s", c.baseURL, c.dbName, encodedId)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, leafURL, nil)
	if err != nil {
		return nil, rerrors.Wrap(err, "failed to create leaf request")
	}
	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, rerrors.Wrap(err, "failed to execute leaf request")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, rerrors.Wrap(err, "failed to read leaf body")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, rerrors.New(fmt.Sprintf("leaf fetch unexpected status %d: %s", resp.StatusCode, string(body)))
	}

	var leafDoc struct {
		Data string `json:"data"`
	}
	if err := json.Unmarshal(body, &leafDoc); err != nil {
		return nil, rerrors.Wrap(err, "failed to unmarshal leaf")
	}

	decoded, err := base64.StdEncoding.DecodeString(leafDoc.Data)
	if err != nil {
		return nil, rerrors.Wrap(err, "failed to decode leaf base64")
	}
	return decoded, nil
}

func (c *LiveSyncClient) ListFiles(ctx context.Context) ([]FileEntry, error) {
	listURL := fmt.Sprintf("%s/%s/_all_docs?include_docs=true", c.baseURL, c.dbName)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, listURL, nil)
	if err != nil {
		return nil, rerrors.Wrap(err, "failed to create request")
	}
	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, rerrors.Wrap(err, "failed to execute request")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, rerrors.Wrap(err, "failed to read response body")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, rerrors.New(fmt.Sprintf("unexpected status %d: %s", resp.StatusCode, string(body)))
	}

	var allDocsResp struct {
		Rows []struct {
			Id  string `json:"id"`
			Doc struct {
				Type    string `json:"type"`
				Deleted bool   `json:"deleted"`
				Mtime   int64  `json:"mtime"`
			} `json:"doc"`
		} `json:"rows"`
	}
	if err := json.Unmarshal(body, &allDocsResp); err != nil {
		return nil, rerrors.Wrap(err, "failed to unmarshal response")
	}

	var files []FileEntry
	for _, row := range allDocsResp.Rows {
		if strings.HasPrefix(row.Id, "_design/") {
			continue
		}
		if row.Doc.Type != "plain" && row.Doc.Type != "newnote" {
			continue
		}
		if row.Doc.Deleted {
			continue
		}
		mime := MimeTypeForPath(row.Id)
		if strings.HasPrefix(mime, "text/") {
			continue
		}
		files = append(files, FileEntry{
			Path:     row.Id,
			Mtime:    row.Doc.Mtime,
			MimeType: mime,
		})
	}
	return files, nil
}

func (c *LiveSyncClient) ReadFile(ctx context.Context, path string) (FileDoc, error) {
	mime := MimeTypeForPath(path)
	encodedId := url.PathEscape(path)
	docURL := fmt.Sprintf("%s/%s/%s", c.baseURL, c.dbName, encodedId)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, docURL, nil)
	if err != nil {
		return FileDoc{}, rerrors.Wrap(err, "failed to create request")
	}
	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return FileDoc{}, rerrors.Wrap(err, "failed to execute request")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return FileDoc{}, rerrors.Wrap(err, "failed to read response body")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return FileDoc{}, rerrors.New(fmt.Sprintf("unexpected status %d: %s", resp.StatusCode, string(body)))
	}

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
	if err := json.Unmarshal(body, &couchDoc); err != nil {
		return FileDoc{}, rerrors.Wrap(err, "failed to unmarshal response")
	}

	var rawBytes []byte
	switch couchDoc.Type {
	case "newnote":
		for _, hash := range couchDoc.Children {
			chunk, err := c.fetchLeaf(ctx, hash)
			if err != nil {
				return FileDoc{}, rerrors.Wrap(err, fmt.Sprintf("failed to fetch leaf %s", hash))
			}
			rawBytes = append(rawBytes, chunk...)
		}
	default:
		decoded, err := base64.StdEncoding.DecodeString(couchDoc.Data)
		if err != nil {
			return FileDoc{}, rerrors.Wrap(err, "failed to decode base64 data")
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
		return rerrors.Wrap(err, "failed to get doc rev")
	}

	encodedId := url.PathEscape(path)
	docURL := fmt.Sprintf("%s/%s/%s", c.baseURL, c.dbName, encodedId)

	doc := map[string]interface{}{
		"_id":     path,
		"_rev":    rev,
		"deleted": true,
	}
	docBody, err := json.Marshal(doc)
	if err != nil {
		return rerrors.Wrap(err, "failed to marshal document")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, docURL, strings.NewReader(string(docBody)))
	if err != nil {
		return rerrors.Wrap(err, "failed to create request")
	}
	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return rerrors.Wrap(err, "failed to execute request")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return rerrors.Wrap(err, "failed to read response body")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return rerrors.New(fmt.Sprintf("unexpected status %d: %s", resp.StatusCode, string(body)))
	}
	return nil
}

func (c *LiveSyncClient) MoveFile(ctx context.Context, oldPath, newPath string) error {
	encodedId := url.PathEscape(oldPath)
	docURL := fmt.Sprintf("%s/%s/%s", c.baseURL, c.dbName, encodedId)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, docURL, nil)
	if err != nil {
		return rerrors.Wrap(err, "failed to create request")
	}
	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return rerrors.Wrap(err, "failed to execute request")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return rerrors.Wrap(err, "failed to read response body")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return rerrors.New(fmt.Sprintf("unexpected status %d: %s", resp.StatusCode, string(body)))
	}

	var couchDoc struct {
		Rev  string `json:"_rev"`
		Type string `json:"type"`
	}
	if err := json.Unmarshal(body, &couchDoc); err != nil {
		return rerrors.Wrap(err, "failed to unmarshal response")
	}

	if couchDoc.Type == "newnote" {
		return user_errors.ChunkedBinaryMoveNotSupported
	}

	file, err := c.ReadFile(ctx, oldPath)
	if err != nil {
		return rerrors.Wrap(err, "failed to read file")
	}

	encoded := base64.StdEncoding.EncodeToString(file.RawBytes)
	newDoc := map[string]interface{}{
		"_id":   newPath,
		"data":  encoded,
		"mtime": file.Mtime,
		"ctime": file.Ctime,
		"size":  file.Size,
		"type":  "plain",
	}
	newDocBody, err := json.Marshal(newDoc)
	if err != nil {
		return rerrors.Wrap(err, "failed to marshal new document")
	}

	newEncodedId := url.PathEscape(newPath)
	newDocURL := fmt.Sprintf("%s/%s/%s", c.baseURL, c.dbName, newEncodedId)

	putReq, err := http.NewRequestWithContext(ctx, http.MethodPut, newDocURL, strings.NewReader(string(newDocBody)))
	if err != nil {
		return rerrors.Wrap(err, "failed to create put request")
	}
	putReq.SetBasicAuth(c.username, c.password)
	putReq.Header.Set("Content-Type", "application/json")
	putReq.Header.Set("Accept", "application/json")

	putResp, err := c.http.Do(putReq)
	if err != nil {
		return rerrors.Wrap(err, "failed to execute put request")
	}
	defer putResp.Body.Close()

	putBody, err := io.ReadAll(putResp.Body)
	if err != nil {
		return rerrors.Wrap(err, "failed to read put response body")
	}
	if putResp.StatusCode < 200 || putResp.StatusCode >= 300 {
		return rerrors.New(fmt.Sprintf("unexpected status %d: %s", putResp.StatusCode, string(putBody)))
	}

	return c.DeleteFile(ctx, oldPath)
}
