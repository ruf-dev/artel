package notes

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"strings"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/clients/couchdb"
	"github.com/ruf-dev/artel/internal/clients/vaultbucket"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"github.com/ruf-dev/artel/internal/storage"
	"github.com/ruf-dev/artel/internal/utils"
	"go.redsock.ru/rerrors"
)

type Service struct {
	vaults         repository.Vaults
	vaultMembers   repository.VaultMembers
	couchInstances repository.CouchInstances
	couchAccounts  repository.CouchAccounts
	s3Instances    repository.S3Instances
	subscriptions  service.SubscriptionService
}

func New(repo repository.Repo, subscriptions service.SubscriptionService) *Service {
	notesSvc := &Service{
		vaults:         repo.Vaults(),
		vaultMembers:   repo.VaultMembers(),
		couchInstances: repo.CouchInstances(),
		couchAccounts:  repo.CouchAccounts(),
		s3Instances:    repo.S3Instances(),
		subscriptions:  subscriptions,
	}

	return notesSvc
}

// resolveVault runs the auth/permission/membership checks shared by liveSyncClient (markdown-
// only callers) and clientAndBucket (folder export/import/delete, which also needs the vault's
// binary-storage settings), and returns the vault row.
func (s *Service) resolveVault(ctx context.Context, vaultID uuid.UUID) (domain.Vault, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return domain.Vault{}, rerrors.Wrap(user_errors.Unauthenticated)
	}

	err := s.subscriptions.CheckFeature(ctx, uc.UserUuid, domain.FeatureNotes)
	if err != nil {
		return domain.Vault{}, err
	}

	_, err = s.vaultMembers.Get(ctx, vaultID, uc.UserUuid)
	if err != nil {
		return domain.Vault{}, rerrors.Wrap(user_errors.Unauthenticated, "not a vault member")
	}

	vault, err := s.vaults.GetByID(ctx, vaultID)
	if err != nil {
		return domain.Vault{}, rerrors.Wrap(err, "error getting vault")
	}

	return vault, nil
}

// resolveVaultForRead is resolveVault's read-only counterpart: it grants access to vault
// members exactly like resolveVault, but also allows a non-member authenticated caller through
// when the vault has been published (domain.Vault.IsPublic). Used only by the read-only RPCs
// (ListFolders/ListNotes/GetNote/ListTags) — write paths keep using resolveVault via
// liveSyncClient/clientAndBucket unchanged.
func (s *Service) resolveVaultForRead(ctx context.Context, vaultID uuid.UUID) (domain.Vault, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return domain.Vault{}, rerrors.Wrap(user_errors.Unauthenticated)
	}

	err := s.subscriptions.CheckFeature(ctx, uc.UserUuid, domain.FeatureNotes)
	if err != nil {
		return domain.Vault{}, err
	}

	_, err = s.vaultMembers.Get(ctx, vaultID, uc.UserUuid)
	if err == nil {
		return s.vaults.GetByID(ctx, vaultID)
	}

	vault, err := s.vaults.GetByID(ctx, vaultID)
	if err != nil {
		return domain.Vault{}, rerrors.Wrap(user_errors.Unauthenticated, "not a vault member")
	}

	if !vault.IsPublic {
		return domain.Vault{}, rerrors.Wrap(user_errors.Unauthenticated, "not a vault member")
	}

	return vault, nil
}

func (s *Service) liveSyncClientForVault(ctx context.Context, vault domain.Vault) (*couchdb.LiveSyncClient, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return nil, rerrors.Wrap(user_errors.Unauthenticated)
	}

	instance, err := s.couchInstances.Get(ctx, vault.CouchInstanceUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "error getting couch instance")
	}

	account, err := s.couchAccounts.GetByUserAndInstance(ctx, uc.UserUuid, vault.CouchInstanceUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "error getting couch account")
	}

	client := couchdb.NewLiveSyncClient(
		instance.Url,
		vault.CouchDBName,
		account.CouchUsername,
		account.CouchPassword,
	)

	return client, nil
}

func (s *Service) liveSyncClient(ctx context.Context, vaultID uuid.UUID) (*couchdb.LiveSyncClient, error) {
	vault, err := s.resolveVault(ctx, vaultID)
	if err != nil {
		return nil, err
	}

	return s.liveSyncClientForVault(ctx, vault)
}

// liveSyncClientForRead mirrors liveSyncClient but resolves the vault via resolveVaultForRead,
// so a non-member authenticated caller can still read a published vault's notes.
func (s *Service) liveSyncClientForRead(ctx context.Context, vaultID uuid.UUID) (*couchdb.LiveSyncClient, error) {
	vault, err := s.resolveVaultForRead(ctx, vaultID)
	if err != nil {
		return nil, err
	}

	return s.liveSyncClientForVault(ctx, vault)
}

// clientAndBucket resolves both halves of a vault's content storage: the CouchDB livesync
// client for markdown notes, and the storage.BinaryStore (S3 or CouchDB-adapter, may be nil)
// for everything else.
func (s *Service) clientAndBucket(ctx context.Context, vaultID uuid.UUID) (*couchdb.LiveSyncClient, storage.BinaryStore, error) {
	vault, err := s.resolveVault(ctx, vaultID)
	if err != nil {
		return nil, nil, err
	}

	client, err := s.liveSyncClientForVault(ctx, vault)
	if err != nil {
		return nil, nil, err
	}

	bucket, err := vaultbucket.Resolve(ctx, vault, client, s.s3Instances)
	if err != nil {
		return nil, nil, rerrors.Wrap(err, "error resolving vault bucket")
	}

	return client, bucket, nil
}

func (s *Service) ListFolders(ctx context.Context, vaultID uuid.UUID) ([]string, error) {
	client, err := s.liveSyncClientForRead(ctx, vaultID)
	if err != nil {
		return nil, err
	}

	return client.ListFolders(ctx)
}

func (s *Service) ListNotes(ctx context.Context, vaultID uuid.UUID) ([]couchdb.NoteEntry, error) {
	client, err := s.liveSyncClientForRead(ctx, vaultID)
	if err != nil {
		return nil, err
	}

	return client.ListNotes(ctx)
}

func (s *Service) GetNote(ctx context.Context, vaultID uuid.UUID, path string) (couchdb.NoteDoc, error) {
	client, err := s.liveSyncClientForRead(ctx, vaultID)
	if err != nil {
		return couchdb.NoteDoc{}, err
	}

	return client.ReadNote(ctx, path)
}

func (s *Service) ListTags(ctx context.Context, vaultID uuid.UUID) ([]string, error) {
	client, err := s.liveSyncClientForRead(ctx, vaultID)
	if err != nil {
		return nil, rerrors.Wrap(err)
	}

	tags, err := client.ListTags(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "list tags")
	}

	return tags, nil
}

func (s *Service) SaveNote(ctx context.Context, vaultID uuid.UUID, path, content string) error {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return rerrors.Wrap(user_errors.Unauthenticated)
	}

	if len(content) > domain.MaxSingleFileBytes {
		return user_errors.FileTooLarge
	}

	client, err := s.liveSyncClient(ctx, vaultID)
	if err != nil {
		return err
	}

	err = s.subscriptions.CheckStorageQuota(ctx, uc.UserUuid)
	if err != nil {
		return err
	}

	err = client.WriteNote(ctx, path, content)
	if err != nil {
		return rerrors.Wrap(err, "error writing note")
	}

	return nil
}

// MoveNote renames/moves a single note. A move onto an already-occupied destination is
// rejected with user_errors.AlreadyExists rather than silently overwriting (client.MoveNote
// upserts on write, so without this check a move onto an existing path would clobber it).
func (s *Service) MoveNote(ctx context.Context, vaultID uuid.UUID, oldPath, newPath string) error {
	client, err := s.liveSyncClient(ctx, vaultID)
	if err != nil {
		return err
	}

	if oldPath == newPath {
		return nil
	}

	notes, err := client.ListNotes(ctx)
	if err != nil {
		return rerrors.Wrap(err, "error listing notes")
	}

	for _, n := range notes {
		if n.Path == newPath {
			return user_errors.AlreadyExists
		}
	}

	err = client.MoveNote(ctx, oldPath, newPath)
	if err != nil {
		return rerrors.Wrap(err, "error moving note")
	}

	return nil
}

// normalizeFolderPath strips leading/trailing slashes so "" and "/" both mean "vault root".
func normalizeFolderPath(p string) string {
	return strings.Trim(p, "/")
}

// underFolder reports whether notePath is folderPath itself or lives under it; an empty
// folderPath (vault root) matches everything.
func underFolder(notePath, folderPath string) bool {
	if folderPath == "" {
		return true
	}

	return notePath == folderPath || strings.HasPrefix(notePath, folderPath+"/")
}

// relativeToFolder strips folderPath (and its separator) off notePath, so an exported zip's
// entries are relative to the folder being exported rather than to the vault root.
func relativeToFolder(notePath, folderPath string) string {
	if folderPath == "" {
		return notePath
	}

	return strings.TrimPrefix(strings.TrimPrefix(notePath, folderPath), "/")
}

// joinDestPath builds the destination vault path for a zip entry landing under destPath
// ("" meaning vault root).
func joinDestPath(destPath, entryName string) string {
	entryName = strings.TrimPrefix(entryName, "/")

	if destPath == "" {
		return entryName
	}

	return destPath + "/" + entryName
}

type zipEntry struct {
	Name    string
	Content []byte
}

func readZipFileEntry(f *zip.File) ([]byte, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, rerrors.Wrap(err, "error opening zip entry")
	}
	defer utils.CloseWithLog(rc, "zip entry reader")

	content, err := io.ReadAll(rc)
	if err != nil {
		return nil, rerrors.Wrap(err, "error reading zip entry")
	}

	return content, nil
}

func readZipEntries(zipData []byte) ([]zipEntry, error) {
	reader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return nil, rerrors.Wrap(err, "error reading zip archive")
	}

	entries := make([]zipEntry, 0, len(reader.File))

	for _, f := range reader.File {
		if f.FileInfo().IsDir() {
			continue
		}

		content, err := readZipFileEntry(f)
		if err != nil {
			return nil, err
		}

		entry := zipEntry{Name: f.Name, Content: content}
		entries = append(entries, entry)
	}

	return entries, nil
}

func writeZipEntry(zw *zip.Writer, name string, content []byte) error {
	w, err := zw.Create(name)
	if err != nil {
		return rerrors.Wrap(err, "error creating zip entry")
	}

	_, err = w.Write(content)
	if err != nil {
		return rerrors.Wrap(err, "error writing zip entry")
	}

	return nil
}

// ExportFolder zips every note and binary file under folderPath (empty means the whole
// vault), with entries relative to folderPath, so re-importing the zip elsewhere lands its
// contents directly under whatever destination path the user picks.
func (s *Service) ExportFolder(ctx context.Context, vaultID uuid.UUID, folderPath string) ([]byte, error) {
	client, bucket, err := s.clientAndBucket(ctx, vaultID)
	if err != nil {
		return nil, err
	}

	folderPath = normalizeFolderPath(folderPath)

	notes, err := client.ListNotes(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing notes")
	}

	var files []storage.ObjectEntry
	if bucket != nil {
		files, err = bucket.List(ctx)
		if err != nil {
			return nil, rerrors.Wrap(err, "error listing files")
		}
	}

	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)

	for _, n := range notes {
		if !underFolder(n.Path, folderPath) {
			continue
		}

		var note couchdb.NoteDoc
		note, err = client.ReadNote(ctx, n.Path)
		if err != nil {
			return nil, rerrors.Wrap(err, "error reading note")
		}

		err = writeZipEntry(zw, relativeToFolder(n.Path, folderPath), []byte(note.Content))
		if err != nil {
			return nil, err
		}
	}

	for _, f := range files {
		if !underFolder(f.Path, folderPath) {
			continue
		}

		var obj storage.Object
		obj, err = bucket.Get(ctx, f.Path)
		if err != nil {
			return nil, rerrors.Wrap(err, "error reading file")
		}

		err = writeZipEntry(zw, relativeToFolder(f.Path, folderPath), obj.RawBytes)
		if err != nil {
			return nil, err
		}
	}

	err = zw.Close()
	if err != nil {
		return nil, rerrors.Wrap(err, "error closing zip writer")
	}

	return buf.Bytes(), nil
}

func (s *Service) existingPaths(ctx context.Context, client *couchdb.LiveSyncClient, bucket storage.BinaryStore) (map[string]bool, error) {
	notes, err := client.ListNotes(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing notes")
	}

	existing := make(map[string]bool, len(notes))
	for _, n := range notes {
		existing[n.Path] = true
	}

	if bucket == nil {
		return existing, nil
	}

	files, err := bucket.List(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing files")
	}

	for _, f := range files {
		existing[f.Path] = true
	}

	return existing, nil
}

// CheckImportConflicts is a dry run over an uploaded zip: it reports which of its entries
// would collide with an existing note/file at destPath, without writing anything. The caller
// re-submits the same zip to CommitImport along with a resolution per conflicting path.
func (s *Service) CheckImportConflicts(ctx context.Context, vaultID uuid.UUID, destPath string, zipData []byte) ([]string, error) {
	client, bucket, err := s.clientAndBucket(ctx, vaultID)
	if err != nil {
		return nil, err
	}

	destPath = normalizeFolderPath(destPath)

	entries, err := readZipEntries(zipData)
	if err != nil {
		return nil, err
	}

	existing, err := s.existingPaths(ctx, client, bucket)
	if err != nil {
		return nil, err
	}

	var conflicts []string
	for _, e := range entries {
		dest := joinDestPath(destPath, e.Name)
		if existing[dest] {
			conflicts = append(conflicts, dest)
		}
	}

	return conflicts, nil
}

func (s *Service) writeImportedEntry(
	ctx context.Context,
	client *couchdb.LiveSyncClient,
	bucket storage.BinaryStore,
	path string,
	content []byte,
) error {
	if storage.IsMarkdown(path) {
		err := client.WriteNote(ctx, path, string(content))
		if err != nil {
			return rerrors.Wrap(err, "error writing note")
		}

		return nil
	}

	err := bucket.Put(ctx, path, content, couchdb.MimeTypeForPath(path))
	if err != nil {
		return rerrors.Wrap(err, "error writing file")
	}

	return nil
}

// CommitImport writes an uploaded zip's entries under destPath, applying resolutions for any
// path CheckImportConflicts previously reported as conflicting. Entries with no resolution
// and no conflict are written directly. Binary entries are skipped (not failed) when the
// vault has no bucket, matching the best-effort semantics of DeleteFolder.
func (s *Service) CommitImport(
	ctx context.Context,
	vaultID uuid.UUID,
	destPath string,
	zipData []byte,
	resolutions []domain.ImportResolution,
) (imported int, skipped int, err error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return 0, 0, rerrors.Wrap(user_errors.Unauthenticated)
	}

	client, bucket, err := s.clientAndBucket(ctx, vaultID)
	if err != nil {
		return 0, 0, err
	}

	err = s.subscriptions.CheckStorageQuota(ctx, uc.UserUuid)
	if err != nil {
		return 0, 0, err
	}

	destPath = normalizeFolderPath(destPath)

	entries, err := readZipEntries(zipData)
	if err != nil {
		return 0, 0, err
	}

	resolutionByPath := make(map[string]domain.ImportResolution, len(resolutions))
	for _, r := range resolutions {
		resolutionByPath[r.Path] = r
	}

	for _, e := range entries {
		dest := joinDestPath(destPath, e.Name)
		writePath := dest

		if resolution, ok := resolutionByPath[dest]; ok {
			switch resolution.Action {
			case domain.ImportActionSkip:
				skipped++
				continue
			case domain.ImportActionRename:
				writePath = joinDestPath(destPath, resolution.RenameTo)
			case domain.ImportActionOverwrite:
				// writePath already equals dest
			}
		}

		if !storage.IsMarkdown(writePath) && bucket == nil {
			skipped++
			continue
		}

		if len(e.Content) > domain.MaxSingleFileBytes {
			return imported, skipped, user_errors.FileTooLarge
		}

		err = s.writeImportedEntry(ctx, client, bucket, writePath, e.Content)
		if err != nil {
			return imported, skipped, err
		}

		imported++
	}

	return imported, skipped, nil
}

// DeleteFolder best-effort deletes every note and file under folderPath. CouchDB has no
// cross-document transactions here, so this cannot be all-or-nothing: it keeps going past
// individual failures and reports them in failedPaths rather than aborting.
func (s *Service) DeleteFolder(ctx context.Context, vaultID uuid.UUID, folderPath string) (deletedCount int, failedPaths []string, err error) {
	client, bucket, err := s.clientAndBucket(ctx, vaultID)
	if err != nil {
		return 0, nil, err
	}

	folderPath = normalizeFolderPath(folderPath)

	notes, err := client.ListNotes(ctx)
	if err != nil {
		return 0, nil, rerrors.Wrap(err, "error listing notes")
	}

	var objects []storage.ObjectEntry
	if bucket != nil {
		objects, err = bucket.List(ctx)
		if err != nil {
			return 0, nil, rerrors.Wrap(err, "error listing files")
		}
	}

	for _, n := range notes {
		if !underFolder(n.Path, folderPath) {
			continue
		}

		delErr := client.DeleteNote(ctx, n.Path)
		if delErr != nil {
			failedPaths = append(failedPaths, n.Path)
			continue
		}

		deletedCount++
	}

	for _, o := range objects {
		if !underFolder(o.Path, folderPath) {
			continue
		}

		delErr := bucket.Delete(ctx, o.Path)
		if delErr != nil {
			failedPaths = append(failedPaths, o.Path)
			continue
		}

		deletedCount++
	}

	return deletedCount, failedPaths, nil
}

// SyncFromWorkbench pushes markdown files edited inside a stopped workbench container's
// workspace back into vaultID's CouchDB notes. This is a snapshot sync, not continuous/live
// sync: snapshot is the path→mtime baseline captured when the vault was last materialized into
// the workbench (internal/service/v1/workbench/workbench.go's materializeVault), used to tell an
// untouched vault-side note apart from one edited elsewhere (e.g. the Obsidian app) while the
// workbench was running. Binary/bucket files are explicitly out of scope for this stage (see
// ExportFolder's bucket half) — non-markdown paths in files are silently skipped, never treated
// as conflicts.
//
// Only a failure resolving the vault or listing its current notes aborts the whole call early;
// a single path's write/read failure doesn't abort the rest (best-effort, mirrors DeleteFolder)
// — since the only return values are the conflicts slice and an overall error, such a path is
// folded into conflicts too, so the caller still knows it wasn't safely synced.
func (s *Service) SyncFromWorkbench(
	ctx context.Context, vaultID uuid.UUID, files map[string][]byte, snapshot map[string]int64,
) (conflicts []string, err error) {
	client, err := s.liveSyncClient(ctx, vaultID)
	if err != nil {
		return nil, err
	}

	return syncNotesFromWorkbench(ctx, client, files, snapshot)
}

// syncNotesFromWorkbench holds SyncFromWorkbench's actual conflict-detection logic, factored out
// as a free function purely to keep SyncFromWorkbench itself short.
func syncNotesFromWorkbench(
	ctx context.Context, client *couchdb.LiveSyncClient, files map[string][]byte, snapshot map[string]int64,
) (conflicts []string, err error) {
	currentNotes, err := client.ListNotes(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing notes")
	}

	currentMtime := make(map[string]int64, len(currentNotes))
	for _, n := range currentNotes {
		currentMtime[n.Path] = n.Mtime
	}

	for path, content := range files {
		if !storage.IsMarkdown(path) {
			continue
		}

		mtime, exists := currentMtime[path]
		if !exists || mtime <= snapshot[path] {
			// New file, or unchanged in the vault since materialization — safe to overwrite.
			writeErr := client.WriteNote(ctx, path, string(content))
			if writeErr != nil {
				conflicts = append(conflicts, path)
			}

			continue
		}

		// The vault-side copy changed since materialization (e.g. edited via the Obsidian app
		// while the workbench was running) — only a genuine conflict if its content actually
		// diverged from what the workbench started with; otherwise the "change" is just the
		// workbench's own prior write being observed back.
		current, readErr := client.ReadNote(ctx, path)
		if readErr != nil {
			conflicts = append(conflicts, path)
			continue
		}

		if current.Content == string(content) {
			continue
		}

		conflicts = append(conflicts, path)
	}

	return conflicts, nil
}

// destPathUnderFolder rewrites a path living under oldFolder so it lives under newFolder
// instead, preserving everything past the oldFolder prefix.
func destPathUnderFolder(path, oldFolder, newFolder string) string {
	rel := relativeToFolder(path, oldFolder)
	return joinDestPath(newFolder, rel)
}

// MoveFolder renames folderPath to newFolderPath by moving every note and binary file whose
// path is folderPath itself or lives under it. Unlike DeleteFolder, this is all-or-nothing:
// every destination path is checked for conflicts against everything NOT being moved before
// anything is touched, and the whole move is rejected with user_errors.AlreadyExists if any
// collide, rather than partially applying and reporting failures.
func (s *Service) MoveFolder(ctx context.Context, vaultID uuid.UUID, oldFolderPath, newFolderPath string) (movedCount int, err error) {
	oldFolderPath = normalizeFolderPath(oldFolderPath)
	newFolderPath = normalizeFolderPath(newFolderPath)

	if oldFolderPath == newFolderPath {
		return 0, nil
	}

	if oldFolderPath != "" && underFolder(newFolderPath, oldFolderPath) {
		return 0, user_errors.InvalidFolderMove
	}

	client, bucket, err := s.clientAndBucket(ctx, vaultID)
	if err != nil {
		return 0, err
	}

	notes, err := client.ListNotes(ctx)
	if err != nil {
		return 0, rerrors.Wrap(err, "error listing notes")
	}

	var files []storage.ObjectEntry
	if bucket != nil {
		files, err = bucket.List(ctx)
		if err != nil {
			return 0, rerrors.Wrap(err, "error listing files")
		}
	}

	var movingNotes []couchdb.NoteEntry
	for _, n := range notes {
		if underFolder(n.Path, oldFolderPath) {
			movingNotes = append(movingNotes, n)
		}
	}

	var movingFiles []storage.ObjectEntry
	for _, f := range files {
		if underFolder(f.Path, oldFolderPath) {
			movingFiles = append(movingFiles, f)
		}
	}

	existing := make(map[string]bool, len(notes)+len(files))
	for _, n := range notes {
		existing[n.Path] = true
	}

	for _, f := range files {
		existing[f.Path] = true
	}

	moving := make(map[string]bool, len(movingNotes)+len(movingFiles))
	for _, n := range movingNotes {
		moving[n.Path] = true
	}

	for _, f := range movingFiles {
		moving[f.Path] = true
	}

	noteDest := make(map[string]string, len(movingNotes))
	for _, n := range movingNotes {
		dest := destPathUnderFolder(n.Path, oldFolderPath, newFolderPath)
		if existing[dest] && !moving[dest] {
			return 0, user_errors.AlreadyExists
		}

		noteDest[n.Path] = dest
	}

	fileDest := make(map[string]string, len(movingFiles))
	for _, f := range movingFiles {
		dest := destPathUnderFolder(f.Path, oldFolderPath, newFolderPath)
		if existing[dest] && !moving[dest] {
			return 0, user_errors.AlreadyExists
		}

		fileDest[f.Path] = dest
	}

	for _, n := range movingNotes {
		moveErr := client.MoveNote(ctx, n.Path, noteDest[n.Path])
		if moveErr != nil {
			return movedCount, rerrors.Wrap(moveErr, "error moving note")
		}

		movedCount++
	}

	for _, f := range movingFiles {
		moveErr := bucket.Move(ctx, f.Path, fileDest[f.Path])
		if moveErr != nil {
			return movedCount, rerrors.Wrap(moveErr, "error moving file")
		}

		movedCount++
	}

	return movedCount, nil
}
