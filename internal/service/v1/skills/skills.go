package skills

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/clients/couchdb"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

// Service manages skills stored as CouchDB notes under a vault's reserved .skills/ folder,
// plus the always-present, non-vault-backed system "skill-creator" skill (see SystemSkill).
// Vault resolution/auth mirrors notes.Service exactly (resolveVault -> liveSyncClientForVault).
type Service struct {
	vaults         repository.Vaults
	vaultMembers   repository.VaultMembers
	couchInstances repository.CouchInstances
	couchAccounts  repository.CouchAccounts
	subscriptions  service.SubscriptionService
}

func New(repo repository.Repo, subscriptions service.SubscriptionService) *Service {
	skillsSvc := &Service{
		vaults:         repo.Vaults(),
		vaultMembers:   repo.VaultMembers(),
		couchInstances: repo.CouchInstances(),
		couchAccounts:  repo.CouchAccounts(),
		subscriptions:  subscriptions,
	}

	return skillsSvc
}

// resolveVault runs the same auth/permission/membership checks notes.Service.resolveVault
// does, gating on domain.FeatureNotes since skills live in the same vault-scoped CouchDB
// storage as notes.
func (s *Service) resolveVault(ctx context.Context, vaultUuid uuid.UUID) (domain.Vault, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return domain.Vault{}, rerrors.Wrap(user_errors.Unauthenticated)
	}

	err := s.subscriptions.CheckFeature(ctx, uc.UserUuid, domain.FeatureNotes)
	if err != nil {
		return domain.Vault{}, err
	}

	_, err = s.vaultMembers.Get(ctx, vaultUuid, uc.UserUuid)
	if err != nil {
		return domain.Vault{}, rerrors.Wrap(user_errors.Unauthenticated, "not a vault member")
	}

	vault, err := s.vaults.GetByID(ctx, vaultUuid)
	if err != nil {
		return domain.Vault{}, rerrors.Wrap(err, "error getting vault")
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

func (s *Service) liveSyncClient(ctx context.Context, vaultUuid uuid.UUID) (*couchdb.LiveSyncClient, error) {
	vault, err := s.resolveVault(ctx, vaultUuid)
	if err != nil {
		return nil, err
	}

	return s.liveSyncClientForVault(ctx, vault)
}

// skillPath returns the vault-relative CouchDB path a skill note lives at for the given slug
// and hot-plug state.
func skillPath(slug string, hotPlug bool) string {
	if hotPlug {
		return domain.SkillsActiveFolder + slug + ".md"
	}

	return domain.SkillsLibraryFolder + slug + ".md"
}

// skillFromNote reconstructs a domain.Skill from a note's path and raw content, parsing its
// frontmatter for name/description/storage_mode and deriving slug + hot-plug state from the
// path itself.
func skillFromNote(path, content string) domain.Skill {
	fields, body := couchdb.ParseFrontmatter(content)

	slug := strings.TrimSuffix(strings.TrimPrefix(path, domain.SkillsActiveFolder), ".md")
	hotPlug := strings.HasPrefix(path, domain.SkillsActiveFolder)

	if !hotPlug {
		slug = strings.TrimSuffix(strings.TrimPrefix(path, domain.SkillsLibraryFolder), ".md")
	}

	skill := domain.Skill{
		Slug:        slug,
		Name:        fields["name"],
		Description: fields["description"],
		StorageMode: domain.SkillStorageMode(fields["storage_mode"]),
		Body:        body,
		IsHotPlug:   hotPlug,
		IsSystem:    false,
	}

	return skill
}

// renderNote serializes a skill's frontmatter + body into the note content written to CouchDB.
// Name/description are flattened to a single line each — the frontmatter format is flat
// scalars only (see couchdb.ParseFrontmatter).
func renderNote(name, description string, storageMode domain.SkillStorageMode, body string) string {
	flatName := strings.ReplaceAll(name, "\n", " ")
	flatDescription := strings.ReplaceAll(description, "\n", " ")

	var sb strings.Builder

	sb.WriteString("---\n")
	fmt.Fprintf(&sb, "name: %s\n", flatName)
	fmt.Fprintf(&sb, "description: %s\n", flatDescription)
	fmt.Fprintf(&sb, "storage_mode: %s\n", storageMode)
	sb.WriteString("---\n")
	sb.WriteString(body)

	return sb.String()
}

// ListSkills lists every skill visible in vaultUuid: the system skill-creator skill (always
// first, not vault-backed), followed by every note under .skills/active/ and .skills/library/.
func (s *Service) ListSkills(ctx context.Context, vaultUuid uuid.UUID) ([]domain.Skill, error) {
	client, err := s.liveSyncClient(ctx, vaultUuid)
	if err != nil {
		return nil, err
	}

	notes, err := client.ListSkillNotes(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing skill notes")
	}

	result := make([]domain.Skill, 0, len(notes)+1)
	result = append(result, SystemSkill())

	for _, n := range notes {
		note, readErr := client.ReadNote(ctx, n.Path)
		if readErr != nil {
			return nil, rerrors.Wrap(readErr, "error reading skill note")
		}

		result = append(result, skillFromNote(n.Path, note.Content))
	}

	return result, nil
}

// locateSkill finds which of the two skill folders currently holds slug, returning its full
// path and hot-plug state. Listing (rather than a blind ReadNote-and-catch-404) avoids relying
// on error-type detection for a routine "does this exist" check.
func (s *Service) locateSkill(
	ctx context.Context, client *couchdb.LiveSyncClient, slug string,
) (path string, hotPlug bool, err error) {
	notes, err := client.ListSkillNotes(ctx)
	if err != nil {
		return "", false, rerrors.Wrap(err, "error listing skill notes")
	}

	activePath := skillPath(slug, true)
	libraryPath := skillPath(slug, false)

	for _, n := range notes {
		switch n.Path {
		case activePath:
			return activePath, true, nil
		case libraryPath:
			return libraryPath, false, nil
		}
	}

	return "", false, user_errors.NotFound
}

// GetSkillBody returns the full skill (including body) for slug. The system skill-creator
// skill is served directly from SystemSkill without touching CouchDB.
func (s *Service) GetSkillBody(ctx context.Context, vaultUuid uuid.UUID, slug string) (domain.Skill, error) {
	if slug == SkillCreatorSlug {
		return SystemSkill(), nil
	}

	client, err := s.liveSyncClient(ctx, vaultUuid)
	if err != nil {
		return domain.Skill{}, err
	}

	path, _, err := s.locateSkill(ctx, client, slug)
	if err != nil {
		return domain.Skill{}, err
	}

	note, err := client.ReadNote(ctx, path)
	if err != nil {
		return domain.Skill{}, rerrors.Wrap(err, "error reading skill note")
	}

	return skillFromNote(path, note.Content), nil
}

// existingSlugs collects every slug currently in use in either skill folder, plus the reserved
// system skill slug, for CreateSkill's dedup pass.
func (s *Service) existingSlugs(ctx context.Context, client *couchdb.LiveSyncClient) (map[string]bool, error) {
	notes, err := client.ListSkillNotes(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing skill notes")
	}

	taken := make(map[string]bool, len(notes)+1)
	taken[SkillCreatorSlug] = true

	for _, n := range notes {
		slug := strings.TrimSuffix(strings.TrimPrefix(n.Path, domain.SkillsActiveFolder), ".md")
		if strings.HasPrefix(n.Path, domain.SkillsLibraryFolder) {
			slug = strings.TrimSuffix(strings.TrimPrefix(n.Path, domain.SkillsLibraryFolder), ".md")
		}

		taken[slug] = true
	}

	return taken, nil
}

// slugify turns name into a filename-safe slug: lowercased, non-alphanumeric runs collapsed to
// a single dash, leading/trailing dashes trimmed. Falls back to "skill" if that leaves nothing.
func slugify(name string) string {
	lower := strings.ToLower(name)

	var sb strings.Builder

	prevDash := false

	for _, r := range lower {
		isAlnum := (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')

		switch {
		case isAlnum:
			sb.WriteRune(r)

			prevDash = false
		case !prevDash && sb.Len() > 0:
			sb.WriteRune('-')

			prevDash = true
		}
	}

	slug := strings.TrimRight(sb.String(), "-")
	if slug == "" {
		slug = "skill"
	}

	return slug
}

// dedupeSlug appends -2, -3, ... to base until it no longer collides with taken.
func dedupeSlug(base string, taken map[string]bool) string {
	if !taken[base] {
		return base
	}

	for i := 2; ; i++ {
		candidate := fmt.Sprintf("%s-%d", base, i)
		if !taken[candidate] {
			return candidate
		}
	}
}

// CreateSkill slugifies name (deduping against every slug already in use in either folder, and
// against the reserved system skill slug), quota-checks against the vault's effective plan
// caps, then writes the new skill note to active/ or library/ depending on hotPlug.
func (s *Service) CreateSkill(
	ctx context.Context,
	vaultUuid uuid.UUID,
	name string,
	description string,
	storageMode domain.SkillStorageMode,
	body string,
	hotPlug bool,
) (domain.Skill, error) {
	client, err := s.liveSyncClient(ctx, vaultUuid)
	if err != nil {
		return domain.Skill{}, err
	}

	err = s.subscriptions.CheckSkillLimit(ctx, vaultUuid, hotPlug)
	if err != nil {
		return domain.Skill{}, err
	}

	taken, err := s.existingSlugs(ctx, client)
	if err != nil {
		return domain.Skill{}, err
	}

	slug := dedupeSlug(slugify(name), taken)
	path := skillPath(slug, hotPlug)
	content := renderNote(name, description, storageMode, body)

	err = client.WriteNote(ctx, path, content)
	if err != nil {
		return domain.Skill{}, rerrors.Wrap(err, "error writing skill note")
	}

	skill := domain.Skill{
		Slug:        slug,
		Name:        name,
		Description: description,
		StorageMode: storageMode,
		Body:        body,
		IsHotPlug:   hotPlug,
		IsSystem:    false,
	}

	return skill, nil
}

// UpdateSkill rewrites an existing skill's frontmatter/body in place — the slug (and therefore
// its path and hot-plug state) never changes on update; use SetSkillHotPlug for that.
func (s *Service) UpdateSkill(
	ctx context.Context,
	vaultUuid uuid.UUID,
	slug string,
	name string,
	description string,
	storageMode domain.SkillStorageMode,
	body string,
) (domain.Skill, error) {
	if slug == SkillCreatorSlug {
		return domain.Skill{}, user_errors.SystemSkillNotEditable
	}

	client, err := s.liveSyncClient(ctx, vaultUuid)
	if err != nil {
		return domain.Skill{}, err
	}

	path, hotPlug, err := s.locateSkill(ctx, client, slug)
	if err != nil {
		return domain.Skill{}, err
	}

	content := renderNote(name, description, storageMode, body)

	err = client.WriteNote(ctx, path, content)
	if err != nil {
		return domain.Skill{}, rerrors.Wrap(err, "error writing skill note")
	}

	skill := domain.Skill{
		Slug:        slug,
		Name:        name,
		Description: description,
		StorageMode: storageMode,
		Body:        body,
		IsHotPlug:   hotPlug,
		IsSystem:    false,
	}

	return skill, nil
}

// SetSkillHotPlug moves slug between the active/ and library/ folders. Enabling hot-plug
// quota-checks against the hot-plug cap first; disabling never fails on quota (it only frees up
// capacity).
func (s *Service) SetSkillHotPlug(ctx context.Context, vaultUuid uuid.UUID, slug string, hotPlug bool) (domain.Skill, error) {
	if slug == SkillCreatorSlug {
		return domain.Skill{}, user_errors.SystemSkillNotEditable
	}

	client, err := s.liveSyncClient(ctx, vaultUuid)
	if err != nil {
		return domain.Skill{}, err
	}

	path, currentHotPlug, err := s.locateSkill(ctx, client, slug)
	if err != nil {
		return domain.Skill{}, err
	}

	if currentHotPlug == hotPlug {
		note, readErr := client.ReadNote(ctx, path)
		if readErr != nil {
			return domain.Skill{}, rerrors.Wrap(readErr, "error reading skill note")
		}

		return skillFromNote(path, note.Content), nil
	}

	if hotPlug {
		err = s.subscriptions.CheckSkillLimit(ctx, vaultUuid, true)
		if err != nil {
			return domain.Skill{}, err
		}
	}

	newPath := skillPath(slug, hotPlug)

	err = client.MoveNote(ctx, path, newPath)
	if err != nil {
		return domain.Skill{}, rerrors.Wrap(err, "error moving skill note")
	}

	note, err := client.ReadNote(ctx, newPath)
	if err != nil {
		return domain.Skill{}, rerrors.Wrap(err, "error reading moved skill note")
	}

	return skillFromNote(newPath, note.Content), nil
}

// DeleteSkill removes slug from whichever folder it currently lives in.
func (s *Service) DeleteSkill(ctx context.Context, vaultUuid uuid.UUID, slug string) error {
	if slug == SkillCreatorSlug {
		return user_errors.SystemSkillNotEditable
	}

	client, err := s.liveSyncClient(ctx, vaultUuid)
	if err != nil {
		return err
	}

	path, _, err := s.locateSkill(ctx, client, slug)
	if err != nil {
		return err
	}

	err = client.DeleteNote(ctx, path)
	if err != nil {
		return rerrors.Wrap(err, "error deleting skill note")
	}

	return nil
}
