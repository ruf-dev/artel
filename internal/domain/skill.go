package domain

import "strings"

// SkillStorageMode is a hint the skill's frontmatter carries for how an agent following the
// skill should persist data it produces: nothing structured, freeform vault notes, or a
// dedicated Postgres schema the agent designs itself (see skills.SystemSkill's instructional
// body — Artel never auto-creates tables/procedures for structured mode).
type SkillStorageMode string

const (
	SkillStorageNone       SkillStorageMode = "none"
	SkillStorageFreeform   SkillStorageMode = "freeform_notes"
	SkillStorageStructured SkillStorageMode = "structured"
)

// Skill is a plain in-memory value, never a DB row: it is assembled fresh on every read from a
// CouchDB note's frontmatter + path (slug = filename), or — for the system skill — from a Go
// constant. There is no persistent Skill table/repository.
type Skill struct {
	Slug        string
	Name        string
	Description string
	StorageMode SkillStorageMode
	Body        string
	IsHotPlug   bool
	IsSystem    bool
}

const (
	// SkillsFolderPrefix is the vault-path prefix under which every skill note lives, hiding
	// the folder from the regular notes/files listings (mirrors the existing "_design/" skip).
	SkillsFolderPrefix = ".skills/"
	// SkillsActiveFolder holds the hot-plug skill set — each one is synthesized into its own
	// MCP tool. SkillsLibraryFolder holds the rest of the library (available but not injected
	// as a standalone tool).
	SkillsActiveFolder  = ".skills/active/"
	SkillsLibraryFolder = ".skills/library/"
)

// IsSkillPath reports whether path falls under the reserved skills folder.
func IsSkillPath(path string) bool {
	return strings.HasPrefix(path, SkillsFolderPrefix)
}

// SkillToolPrefix is prepended to a hot-plug skill's slug to form its dynamic MCP tool name
// (e.g. skill_my-skill). Synthesized in internal/service/v1/mcp (ListHotPlugSkillTools) and
// matched against in internal/transport/mcp_api (handleToolsCall) to route a tools/call
// request to ExecuteSkillTool instead of the static builtin/MoM dispatch paths.
const SkillToolPrefix = "skill_"
