# Codex repository guidance

This repository is primarily maintained with Claude Code. Treat the existing
Claude instructions as the canonical project context and load them at the
start of every session:

- Read [CLAUDE.md](CLAUDE.md) completely.
- Read [CLAUDE.local.md](CLAUDE.local.md) if it exists. It is intentionally
  local-only and may contain user-specific preferences.
- When working under `pkg/client/ArtelUI`, also read
  [pkg/client/ArtelUI/CLAUDE.md](pkg/client/ArtelUI/CLAUDE.md).

Follow those instructions unless a direct user request overrides them. In
particular, preserve existing work in the working tree, follow the Go and
frontend rules, use the documented test commands, and honor the commit
message convention.

## Graphify

Use the repository's graphify knowledge graph as part of normal codebase
work. The detailed workflow is in
[.claude/skills/graphify/SKILL.md](.claude/skills/graphify/SKILL.md); read it
when using graphify.

- If `graphify-out/graph.json` exists, run a focused `graphify query` before
  investigating a codebase question. Expand query terms only from
  `graphify-out/.vocab.txt` as described by the skill.
- Use `graphify path` for relationships and `graphify explain` for a focused
  concept.
- After code changes, run `graphify update .` when the graph is available and
  the graphify interpreter guard succeeds.
- Do not commit `graphify-out/`; it is local and ignored.

