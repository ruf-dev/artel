# Night-shift

An autonomous "night shift" that pulls ideas from the Artel Trello board and
drives each through product analysis → architecture analysis → isolated-worktree
implementation, testing server-side and opening small **draft** PRs for review.
The user is the project manager and the sole architectural authority.

The orchestration is a machine-wide Claude Code skill, not repo code. This repo
holds only the wiring:

| Piece | Location |
|---|---|
| Overseer skill | `~/.claude/skills/night-shift/` (canonical: vault `.profiles/night-shift/`) |
| Role subagents | `~/.claude/agents/{product-analyst,architect,developer}.md` |
| Analyst profiles | `~/.claude/profiles/{product-analyst,architect}/` |
| Developer profile | `~/.claude/rules/` (the existing coder profile) |
| **Project wiring** | **`.claude/night-shift.yaml`** (this repo) |

## Model

- A card's Trello **column is its state**. Every step is idempotent and decides
  what to do from the current column. All per-card state lives in Trello
  (description + comments + attachments) — no external store.
- Each overseer wake is a **bounded sweep**: snapshot the board, spawn one role
  subagent per actionable card up to `parallelism_default`, collect status,
  sleep. The overseer holds only a `{card, column, status}` table.
- Trello access is via the `trello` MoM through the Artel MCP (read tools plus
  `update_card` / `create_list` / `update_list` / `create_label` /
  `add_attachment`, added in migration 081).

## Boundaries

Autonomous, retryable actions only: branch, push, **draft** PR, card
comment/move/attach, worktree, run tests, docker test infra, commit on a feature
branch. Never: merge, mark a PR ready, force-push, delete a branch, commit to
`master`, migrate a shared DB, deploy.

## Running it

```
/night-shift setup                 # once — fills night-shift.yaml, fixes the board's lists + blocked label
/night-shift limit=1               # first real run — stops after one card reaches IN REVIEW
/night-shift parallelism=2         # unbounded run once the first is validated
```

`make test-e2e` is serialised by the overseer (one developer runs it at a time —
the shared docker stack on ports 15434/15985/19000 cannot take concurrent runs).
`graphify update .` is **not** run overnight; the morning review session runs it
once.
