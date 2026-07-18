# Notes Search & Recall — Architecture Overview

Idea doc, 2026-07-18. Not an implementation plan — captures the agreed direction so a future
plan/migration can start from it instead of re-deriving it.

## Problem

The MCP vault tools (`internal/service/v1/mcp/executors/vault.go`) have no search primitive.
`list_files` returns every path in the vault; `list_tags` (`internal/clients/couchdb/livesync.go:225`)
calls `ListNotes` then `ReadNote` **on every single note** just to regex out `#tags`. There is no
caching or index — every tool call re-scans the whole vault. Today's "search" is: agent lists
everything, then reads candidates into context and searches client-side. That's the bottleneck
this doc addresses.

Also true today: there is no chat/conversation storage in Artel (only an auth `sessions` table —
JWT/refresh tokens, no message content). Conversation/chat history is planned separately, later.
This doc is written so that feature can plug in without reworking what's proposed here.

## Design principle

Memory = notes. Rather than a second "agent memory" store living beside notes the user can see
and edit, an explicit "remember X" defaults to *becoming a note* (existing note edit, or a
fallback note like `Inbox.md` if nothing obvious fits). One source of truth, one search/recall
pipeline serves both real notes and agent-recorded facts.

## Components

| Component | Owns | Status |
|---|---|---|
| CouchDB (per-vault) | Notes — source of truth | exists |
| `notes_index` / `note_tags` / `note_links` (Postgres) | Search index, kept in sync off CouchDB's `_changes` feed | proposed |
| `index_checkpoints` (Postgres) | Last processed `_changes` `seq` per vault DB, for resumable sync | proposed |
| `conversations` / `messages` (Postgres) | Full chat transcripts | planned separately, later |
| `session_summaries` (Postgres) | Short distilled recap per session/topic | proposed |
| MCP tools: `search_notes`, `get_recent_context` | The agent's interface to all of the above | proposed |

Start with Postgres only (`tsvector` + GIN + `pg_trgm`). OpenSearch/pgvector is a phase-2 upgrade
for semantic search once literal keyword+tag+trigram matching stops being "smart enough" — not a
day-one requirement, given this is a multi-tenant system (one shared CouchDB instance can host
many vaults' databases — `internal/domain/couch_instance.go`, `resolve_key.go:102`) and running a
shared search cluster is real ongoing infra weight.

## Why a `_changes`-feed indexer, not a write-path hook

The client type is literally named `LiveSyncClient` because Obsidian's Self-hosted LiveSync plugin
syncs directly from the user's device straight into CouchDB, bypassing Artel's API entirely.
Artel's `write_file`/`delete_file`/`move_file` MCP tools are one writer among several — day to
day, the primary writer is the user's own Obsidian app. Hooking indexing into `write_file` would
silently miss most real edits. The `_changes` feed is the only place that sees every mutation
regardless of origin, so it's the mandatory mechanism, not a convenience.

This also answers "won't a note get written twice?" — no. There is exactly one write path
(CouchDB, whichever client wrote it), and exactly one derived, read-only path (Postgres). Content
lands in CouchDB; only derived metadata (tags, resolved links, `tsvector`) lands in Postgres,
computed by re-reading what CouchDB already has.

**Mechanics:**
- CouchDB's per-database `_changes` endpoint streams `{seq, id, rev, deleted}` in order, resumable
  via `?since=<last_seq>`. A background worker subscribes per vault DB; `index_checkpoints` tracks
  progress so a restart resumes instead of rescanning the whole vault.
- Per changed id: apply the same filter `ListNotes` already uses (skip `_design/` prefixes, skip
  anything not `type=plain`, drop tombstoned docs from the index), fetch the full note, run
  `extractTags` (`livesync.go:239`) plus a wikilink parser, upsert into
  `notes_index`/`note_tags`/`note_links` (delete-then-reinsert per note is fine at this size).
- Binary/S3-backed files never appear in CouchDB's `_changes` (they live in S3) — need a separate,
  coarser signal later (periodic `bucket.List` diff, filename/path only, no content). Flagged as a
  gap, not solved here.

**Trade-off:** indexing is async, so there's a small eventual-consistency window — a
`search_notes` call moments after a `write_file` might not reflect that edit yet. Fine for a
note-taking assistant in the overwhelming majority of cases. If a flow ever needs read-your-own-write
in the same turn, that should be a narrow, deliberate exception (inline upsert right after that one
write) — not a change to the general design.

**Scale lever (not needed yet):** running one `_changes` connection per vault could mean many open
long-polls per CouchDB server as usage grows, since one `CouchInstance` hosts many vaults' DBs.
CouchDB also exposes `_db_updates`, a server-level feed covering every database on that instance —
collapses N-vault connections down to one per physical CouchDB instance. Worth knowing exists;
not needed at current scale.

## Wikilinks are not real paths — resolution has to happen in Artel's code

CouchDB is completely markdown-blind: a note is stored as either a base64 `Data` blob or chunked
`Children` leaf hashes (`docFull` in `livesync.go`). `[[Wikilink]]` syntax is invisible to it —
just bytes inside the blob. So the link graph has to be built entirely by Artel: fetch full text,
run a wikilink parser, resolve.

Obsidian's default link resolution isn't path-based — `[[Some Note]]` resolves by matching a
**basename anywhere in the vault** ("shortest path" mode), not a literal relative path. A raw
wikilink target can't just be treated as a path fragment. Correct resolution needs a basename →
path(s) map built from the vault's own note listing, with each raw link matched against it:

- exact relative path match first, if the link literally is `folder/sub/Note`
- else basename match across the vault
- if the basename match is ambiguous (same basename in multiple folders) — and Obsidian even lets
  the vault owner configure link format (absolute/relative/shortest-path) per vault, a setting
  CouchDB never surfaces — store the raw unresolved link text rather than guessing wrong. A wrong
  edge in `note_links` is worse than a missing one.

## New MCP tools

- **`search_notes(query, tags?, folder?, limit)`** — Postgres FTS over `notes_index`, ranked,
  returns `{path, snippet, score}[]`. Agent calls this first, then `read_file` only on the winner.
  Replaces "fetch everything, search client-side" with "search server-side, fetch what's needed."
- **`get_recent_context`** — session-start recap. Backed by `session_summaries` once conversation
  history exists; until then, degrades to "recently written/touched notes" (derived from
  `notes_index.mtime` or an access log) as a placeholder signal. Same tool contract either way, so
  nothing built now gets thrown away once chat history ships.
- No new "remember" tool — reuse `write_file` (see design principle above).

## Lifecycle A — "remember X"

1. User: *"remember the induction panel is 60cm, it's in the kitchen note."*
2. Agent decides which note this belongs in (existing note, or a fallback like `Kitchen.md` /
   `Inbox.md`) and calls the existing `write_file` tool — no new tool needed.
3. The write lands in CouchDB → the `_changes` worker picks it up → `notes_index`/`note_tags`/
   `note_links` update automatically, same pipeline as any other note edit, no special-casing.
4. (Once conversation storage exists) the turn also gets logged to `messages`, so the fact that a
   kitchen topic was discussed is itself a recall signal, independent of the note content.

## Lifecycle B — recall

1. **Session start** — agent calls `get_recent_context` before answering anything.
2. **In-conversation lookup** — for a concrete question ("what size is my induction panel"), agent
   calls `search_notes(query, tags?)`, gets ranked candidates, then `read_file` on the winner only.

## Phasing

1. **Now:** `notes_index` + `_changes` sync worker + `search_notes` tool (Postgres only — fixes
   the immediate "fetch everything" cost).
2. **Next:** `get_recent_context` backed by the access-log/mtime placeholder.
3. **Once chat/history ships:** wire `session_summaries` generation off real `messages`; upgrade
   `get_recent_context` for free — tool contract doesn't change.
4. **Later, if keyword search isn't smart enough:** semantic search (OpenSearch or pgvector)
   underneath `search_notes` — same signature, better ranking under the hood.

## Open questions / not yet decided

- Wikilink parser: exact grammar to support (`[[Note]]`, `[[Note#heading]]`, `[[Note|alias]]`,
  embeds `![[Note]]`) and what to do with genuinely ambiguous basenames.
- Where `index_checkpoints`/worker lifecycle lives — one goroutine per active vault vs. a shared
  poller — and whether to start with per-vault `_changes` or go straight to `_db_updates`.
- Binary/S3 file indexing (filename/path only) — priority relative to the markdown path.
- Exact shape of the access-log placeholder for `get_recent_context` before conversation history
  exists (which tool calls to log, retention window).
