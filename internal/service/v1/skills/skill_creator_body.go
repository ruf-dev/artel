package skills

import "github.com/ruf-dev/artel/internal/domain"

// SkillCreatorSlug identifies the always-present system skill synthesized by SystemSkill.
const SkillCreatorSlug = "skill-creator"

// skillCreatorBody is skill-creator's instructional body — returned verbatim as the tool
// result when an agent calls the skill_skill-creator hot-plug tool. It walks the agent through
// gathering the fields create_skill needs and, for structured mode, through designing the
// Postgres schema itself before writing the skill.
const skillCreatorBody = `# Skill Creator

You are helping the user define a new skill for this vault. A skill is a reusable, named set of
instructions that gets injected as its own MCP tool (when hot-plugged) or kept in the library for
later promotion. Work through the steps below interactively — ask the user one question at a
time rather than demanding all the fields at once, and propose sensible defaults where you can.

## 1. Gather the basics

- **name** — a short, human-readable title (e.g. "Weekly Report Drafting"). Keep it under ~60
  characters; it becomes the skill's slug (lowercased, non-alphanumeric characters collapsed to
  dashes), so avoid punctuation-heavy titles.
- **description** — one or two sentences describing when this skill should fire. This is the
  *keyword-trigger surface*: once the skill is hot-plugged, its description is the only thing
  another agent sees before deciding whether to invoke it. Write it the way you'd write a
  function docstring meant to be matched against a user's request — mention the situations,
  nouns, and verbs that should trigger it.
- **hot_plug** — ask whether this skill should be hot-plugged immediately (injected as its own
  callable tool, subject to the plan's hot-plug cap) or kept in the library (available but not
  synthesized into a standalone tool until promoted later). Default to hot_plug=false for
  experimental or rarely-used skills to conserve the hot-plug cap.

## 2. Choose a storage mode

Every skill declares one of three storage_mode values, describing how an agent following the
skill should persist any data it produces:

- **none** — the skill is pure instructions/process; it reads and reasons but does not persist
  anything of its own. Most skills should start here.
- **freeform_notes** — the skill's output belongs in the vault as ordinary markdown notes (e.g.
  under a dedicated folder). Use this when the output is prose, logs, or anything better kept
  human-readable and lightly structured.
- **structured** — the skill's output belongs in Postgres, in a purpose-built schema (tables,
  columns, maybe stored procedures) that the skill can query and update reliably across runs.

**Important — structured mode is entirely your responsibility, not Artel's.** Choosing
storage_mode=structured here does NOT cause Artel to create any tables, columns, or procedures
automatically. If the user wants structured storage:

1. Design the schema yourself, in this same conversation, using the vault's postgres MCP tools
   (the ones already available to you for querying/executing SQL against the vault's linked
   Postgres database — enable one first via the postgres connection tools if none is linked yet).
2. Create the tables/columns/procedures you designed, using those same postgres tools.
3. Bake the *exact* resulting table and column names into the new skill's body text (step 3
   below), spelled out precisely — the skill's body is the only place this schema knowledge is
   recorded, so any agent invoking the skill later must be able to read the body and know exactly
   which table/column names to use without re-deriving them.

Do not skip schema creation and only describe it in prose "for later" — a structured-mode skill
whose body doesn't name real, already-created tables is broken by construction.

## 3. Write the body

The body is the actual instructional markdown that gets returned verbatim when the skill is
invoked (as skill_<slug> once hot-plugged, or via get_skill_body from the library). Write it as
you would a runbook for another agent with no other context:

- State the goal and when to stop.
- List concrete steps in order, including which tools to call.
- For storage_mode=freeform_notes, name the vault folder/path convention to use.
- For storage_mode=structured, name the exact table(s)/column(s) from step 2, and give example
  queries if the schema is non-trivial.
- Call out edge cases or things to confirm with the user before acting.

Keep it focused — a skill that tries to cover too much ground is harder to trigger correctly and
harder to follow. Prefer splitting unrelated concerns into separate skills.

## 4. Create it

Once you have name, description, storage_mode, body, and the hot_plug decision, call
create_skill with those fields. Read back the created skill's slug to the user so they know how
to refer to it later (e.g. for update_skill, set_skill_hot_plug, or delete_skill).
`

// SystemSkill returns the always-present, non-vault-backed "skill-creator" skill: it is
// synthesized directly from this Go constant on every read, never stored in or read from
// CouchDB, always counted as hot-plug, and never counted against either quota cap.
func SystemSkill() domain.Skill {
	skill := domain.Skill{
		Slug:        SkillCreatorSlug,
		Name:        "Skill Creator",
		Description: "Interactively define a new skill for this vault — gathers name, description, storage mode, and body, and for structured-storage skills walks through designing the Postgres schema before writing it. Use this whenever the user wants to add, define, or set up a new skill.",
		StorageMode: domain.SkillStorageNone,
		Body:        skillCreatorBody,
		IsHotPlug:   true,
		IsSystem:    true,
	}

	return skill
}
