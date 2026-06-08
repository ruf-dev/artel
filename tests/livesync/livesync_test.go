package livesync_test

import (
	"context"
	"fmt"
	"math/rand/v2"
	"os"
	"strings"
	"testing"

	kivik "github.com/go-kivik/kivik/v4"
	kivikcouch "github.com/go-kivik/kivik/v4/couchdb"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/ruf-dev/artel/internal/clients/couchdb"
)

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// noteID returns a unique, suite-scoped document ID derived from the test name and a suffix.
// Using the test name as a prefix keeps docs isolated within the shared DB.
func noteID(t *testing.T, suffix string) string {
	name := strings.ReplaceAll(t.Name(), "/", "_")
	return fmt.Sprintf("%s/%s.md", name, suffix)
}

type VaultSuite struct {
	suite.Suite
	admin  *couchdb.Client
	vault  *couchdb.LiveSyncClient
	rawDB  *kivik.DB // direct kivik handle for verifying raw document structure
	dbName string
}

func TestVault(t *testing.T) {
	suite.Run(t, new(VaultSuite))
}

func (s *VaultSuite) SetupSuite() {
	couchURL := envOrDefault("COUCH_URL", "http://localhost:15985")
	couchUser := envOrDefault("COUCH_USER", "admin")
	couchPass := envOrDefault("COUCH_PASS", "admin")

	s.dbName = fmt.Sprintf("test_vault_%08x", rand.Uint32())

	s.admin = couchdb.New(couchdb.Config{
		BaseURL:  couchURL,
		User:     couchUser,
		Password: couchPass,
	})

	ctx := context.Background()

	err := s.admin.Setup(ctx)
	s.Require().NoError(err, "CouchDB system setup failed — is the container running?")

	err = s.admin.CreateDatabase(ctx, s.dbName)
	s.Require().NoError(err)

	s.vault = couchdb.NewLiveSyncClient(couchURL, s.dbName, couchUser, couchPass)

	kivikClient, err := kivik.New("couch", couchURL, kivikcouch.BasicAuth(couchUser, couchPass))
	s.Require().NoError(err)
	s.rawDB = kivikClient.DB(s.dbName)
}

func (s *VaultSuite) TeardownSuite() {
	_ = s.admin.DeleteDatabase(context.Background(), s.dbName)
}

// ── Write / Read ─────────────────────────────────────────────────────────────

func (s *VaultSuite) TestWriteRead_Simple() {
	t := s.T()
	ctx := context.Background()
	path := noteID(t, "simple")
	content := "# Hello\nThis is a simple note."

	err := s.vault.WriteNote(ctx, path, content)
	require.NoError(t, err)

	doc, err := s.vault.ReadNote(ctx, path)
	require.NoError(t, err)
	require.Equal(t, content, doc.Content)
	require.False(t, doc.Deleted)
}

func (s *VaultSuite) TestWriteRead_MultipleWrites() {
	t := s.T()
	ctx := context.Background()
	path := noteID(t, "overwrite")

	err := s.vault.WriteNote(ctx, path, "first version")
	require.NoError(t, err)

	err = s.vault.WriteNote(ctx, path, "second version")
	require.NoError(t, err)

	err = s.vault.WriteNote(ctx, path, "third version")
	require.NoError(t, err)

	doc, err := s.vault.ReadNote(ctx, path)
	require.NoError(t, err)
	require.Equal(t, "third version", doc.Content, "expected latest overwrite, not earlier revision")
}

func (s *VaultSuite) TestWriteRead_EmptyContent() {
	t := s.T()
	ctx := context.Background()
	path := noteID(t, "empty")

	err := s.vault.WriteNote(ctx, path, "")
	require.NoError(t, err)

	doc, err := s.vault.ReadNote(ctx, path)
	require.NoError(t, err)
	require.Equal(t, "", doc.Content)
}

func (s *VaultSuite) TestWriteRead_Unicode() {
	t := s.T()
	ctx := context.Background()
	path := noteID(t, "unicode")
	content := "こんにちは世界 — em-dash: —, emoji: 🔥, accents: éàü"

	err := s.vault.WriteNote(ctx, path, content)
	require.NoError(t, err)

	doc, err := s.vault.ReadNote(ctx, path)
	require.NoError(t, err)
	require.Equal(t, content, doc.Content)
}

func (s *VaultSuite) TestWriteRead_LargeNote() {
	t := s.T()
	ctx := context.Background()
	path := noteID(t, "large")
	// 64 KB of content
	content := strings.Repeat("Lorem ipsum dolor sit amet. ", 64*1024/28+1)[:64*1024]

	err := s.vault.WriteNote(ctx, path, content)
	require.NoError(t, err)

	doc, err := s.vault.ReadNote(ctx, path)
	require.NoError(t, err)
	require.Equal(t, content, doc.Content)
}

// ── List Notes ────────────────────────────────────────────────────────────────

func (s *VaultSuite) TestListNotes_AfterWrites() {
	t := s.T()
	ctx := context.Background()

	paths := make([]string, 5)
	for i := range paths {
		paths[i] = noteID(t, fmt.Sprintf("note%d", i))
		err := s.vault.WriteNote(ctx, paths[i], fmt.Sprintf("content %d", i))
		require.NoError(t, err)
	}

	notes, err := s.vault.ListNotes(ctx)
	require.NoError(t, err)

	listed := make(map[string]bool, len(notes))
	for _, n := range notes {
		listed[n.Path] = true
	}

	for _, p := range paths {
		require.True(t, listed[p], "expected %q in ListNotes result", p)
	}
}

func (s *VaultSuite) TestListNotes_DeletedNotIncluded() {
	t := s.T()
	ctx := context.Background()
	path := noteID(t, "delete_list")

	err := s.vault.WriteNote(ctx, path, "will be deleted")
	require.NoError(t, err)

	err = s.vault.DeleteNote(ctx, path)
	require.NoError(t, err)

	notes, err := s.vault.ListNotes(ctx)
	require.NoError(t, err)

	for _, n := range notes {
		require.NotEqual(t, path, n.Path, "deleted note should not appear in ListNotes")
	}
}

func (s *VaultSuite) TestWriteMany_ListAll() {
	t := s.T()
	ctx := context.Background()

	const count = 10
	paths := make([]string, count)
	for i := range paths {
		paths[i] = noteID(t, fmt.Sprintf("many%d", i))
		err := s.vault.WriteNote(ctx, paths[i], fmt.Sprintf("note %d", i))
		require.NoError(t, err)
	}

	notes, err := s.vault.ListNotes(ctx)
	require.NoError(t, err)

	prefix := strings.ReplaceAll(t.Name(), "/", "_")
	var matching int
	for _, n := range notes {
		if strings.HasPrefix(n.Path, prefix) {
			matching++
		}
	}
	require.Equal(t, count, matching, "expected exactly %d notes with test prefix", count)
}

// ── List Folders ──────────────────────────────────────────────────────────────

func (s *VaultSuite) TestListFolders() {
	t := s.T()
	ctx := context.Background()
	prefix := strings.ReplaceAll(t.Name(), "/", "_")

	nested := fmt.Sprintf("%s/subdir/deep/note.md", prefix)
	shallow := fmt.Sprintf("%s/top/note.md", prefix)

	err := s.vault.WriteNote(ctx, nested, "deep note")
	require.NoError(t, err)
	err = s.vault.WriteNote(ctx, shallow, "shallow note")
	require.NoError(t, err)

	folders, err := s.vault.ListFolders(ctx)
	require.NoError(t, err)

	folderSet := make(map[string]bool, len(folders))
	for _, f := range folders {
		folderSet[f] = true
	}

	require.True(t, folderSet[prefix+"/subdir/deep"], "expected nested folder in ListFolders")
	require.True(t, folderSet[prefix+"/top"], "expected top-level folder in ListFolders")
}

// ── List Tags ─────────────────────────────────────────────────────────────────

func (s *VaultSuite) TestListTags() {
	t := s.T()
	ctx := context.Background()
	path := noteID(t, "tagged")
	content := "---\ntags: [alpha, beta, gamma]\n---\n\nNote body."

	err := s.vault.WriteNote(ctx, path, content)
	require.NoError(t, err)

	tags, err := s.vault.ListTags(ctx)
	require.NoError(t, err)

	tagSet := make(map[string]bool, len(tags))
	for _, tag := range tags {
		tagSet[tag] = true
	}

	require.True(t, tagSet["alpha"], "expected tag 'alpha'")
	require.True(t, tagSet["beta"], "expected tag 'beta'")
	require.True(t, tagSet["gamma"], "expected tag 'gamma'")
}

// ── Delete ────────────────────────────────────────────────────────────────────

func (s *VaultSuite) TestDelete_ThenReadReturnsDeleted() {
	t := s.T()
	ctx := context.Background()
	path := noteID(t, "del_read")

	err := s.vault.WriteNote(ctx, path, "will be deleted")
	require.NoError(t, err)

	err = s.vault.DeleteNote(ctx, path)
	require.NoError(t, err)

	doc, err := s.vault.ReadNote(ctx, path)
	require.NoError(t, err, "ReadNote on deleted doc should not error")
	require.True(t, doc.Deleted, "expected Deleted=true after DeleteNote")
}

func (s *VaultSuite) TestDelete_ThenWriteAgain() {
	t := s.T()
	ctx := context.Background()
	path := noteID(t, "del_rewrite")

	err := s.vault.WriteNote(ctx, path, "original")
	require.NoError(t, err)

	err = s.vault.DeleteNote(ctx, path)
	require.NoError(t, err)

	err = s.vault.WriteNote(ctx, path, "resurrected")
	require.NoError(t, err, "WriteNote after delete should succeed (must reuse deleted doc rev)")

	doc, err := s.vault.ReadNote(ctx, path)
	require.NoError(t, err)
	require.Equal(t, "resurrected", doc.Content)
	require.False(t, doc.Deleted, "note written after delete should not be marked deleted")
}

func (s *VaultSuite) TestDelete_NonExistent() {
	t := s.T()
	ctx := context.Background()
	path := noteID(t, "never_written")

	err := s.vault.DeleteNote(ctx, path)
	// DeleteNote calls ReadNote first; a non-existent path should return an error.
	require.Error(t, err, "deleting a non-existent note should return an error")
}

// ── LiveSync plugin compatibility ─────────────────────────────────────────────
//
// The LiveSync Obsidian plugin iterates meta.children unconditionally:
//   for (const child of meta.children) { ... }
// If children is absent from the CouchDB document, meta.children is undefined
// and this throws TypeError: meta.children is not iterable, breaking sync.
// Every document artel writes must include children: [] for non-chunked content.

func (s *VaultSuite) TestLiveSyncCompat_WriteNote_ChildrenIsEmptyArray() {
	t := s.T()
	ctx := context.Background()
	path := noteID(t, "compat_write")

	err := s.vault.WriteNote(ctx, path, "# Note\nContent.")
	require.NoError(t, err)

	d := s.rawDB.Get(ctx, path)
	defer d.Close()

	var raw struct {
		Children []string `json:"children"`
	}
	err = d.ScanDoc(&raw)
	require.NoError(t, err)
	require.NotNil(t, raw.Children, "children field must be an empty array, not missing — LiveSync iterates it unconditionally")
	require.Empty(t, raw.Children)
}

func (s *VaultSuite) TestLiveSyncCompat_DeleteNote_ChildrenIsEmptyArray() {
	t := s.T()
	ctx := context.Background()
	path := noteID(t, "compat_delete")

	err := s.vault.WriteNote(ctx, path, "will be deleted")
	require.NoError(t, err)

	err = s.vault.DeleteNote(ctx, path)
	require.NoError(t, err)

	d := s.rawDB.Get(ctx, path)
	defer d.Close()

	var raw struct {
		Children []string `json:"children"`
		Deleted  bool     `json:"deleted"`
	}
	err = d.ScanDoc(&raw)
	require.NoError(t, err)
	require.True(t, raw.Deleted)
	require.NotNil(t, raw.Children, "children field must be present even on soft-deleted docs")
	require.Empty(t, raw.Children)
}

func (s *VaultSuite) TestLiveSyncCompat_MoveFile_DestinationChildrenIsEmptyArray() {
	t := s.T()
	ctx := context.Background()
	prefix := strings.ReplaceAll(t.Name(), "/", "_")
	oldPath := prefix + "/source.bin"
	newPath := prefix + "/dest.bin"

	// Insert a plain binary doc directly — simulates a file written by the LiveSync plugin.
	_, err := s.rawDB.Put(ctx, oldPath, map[string]interface{}{
		"_id":   oldPath,
		"data":  "aGVsbG8=", // base64("hello")
		"type":  "plain",
		"mtime": int64(0),
		"ctime": int64(0),
		"size":  int64(5),
	})
	require.NoError(t, err)

	err = s.vault.MoveFile(ctx, oldPath, newPath)
	require.NoError(t, err)

	d := s.rawDB.Get(ctx, newPath)
	defer d.Close()

	var raw struct {
		Children []string `json:"children"`
	}
	err = d.ScanDoc(&raw)
	require.NoError(t, err)
	require.NotNil(t, raw.Children, "children field must be present on move destination")
	require.Empty(t, raw.Children)
}
