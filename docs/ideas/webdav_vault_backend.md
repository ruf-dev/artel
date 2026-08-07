# WebDAV Vault Sync Backend — Discussion Notes

Idea doc, 2026-08-07. Captures a design discussion, not an implementation plan — no code was
written. Companion to the BYOK-storage work (S3/CouchDB bring-your-own-connection, implemented
separately) — this doc covers WebDAV specifically, which was scoped out of that implementation
pass and left here for whoever picks it up next.

## Problem

Today Artel provisions exactly one vault sync backend: CouchDB, via Obsidian's community
[Self-hosted LiveSync](https://github.com/vrtmrz/obsidian-livesync) plugin. `internal/service/v1/vault/vault.go`'s
`CreateVault` creates a CouchDB database + per-user account on a pooled (or, after the BYOK work,
user-owned) CouchDB instance, and `internal/livesync/builder.go` generates an
`obsidian://setuplivesync?settings=...` URI that one-tap-configures the LiveSync plugin with those
connection details. There is no equivalent for users who'd rather sync over WebDAV — a protocol
several other Obsidian sync plugins speak (notably
[remotely-save](https://github.com/remotely-save/remotely-save)), and one some users may already
have infrastructure for (a NAS, Nextcloud, etc.) without wanting to stand up CouchDB.

WebDAV has zero references anywhere in the codebase today — this would be a new sync backend
architecture, not a config tweak.

## What "provisioning a WebDAV vault" would need to decide

1. **Which plugin to target.** Different WebDAV-capable Obsidian sync plugins have different
   config shapes and no shared "setup URI" standard the way LiveSync has one official generator.
   `remotely-save` is the closest analog but its settings-import format (if any exists) hasn't
   been checked against this codebase's needs — this is real research, not engineering, and should
   happen before any client work starts.
2. **How much Artel does vs. the user.** Two ends of a spectrum, not fully explored:
   - **Credential-only (cheap):** Artel just verifies the endpoint is reachable (a PROPFIND/GET
     check with the supplied credentials) and stores them via BYOK. The user manually points their
     WebDAV-sync plugin at the endpoint using the shown credentials — no per-vault directory
     creation, no generated setup URI. Smallest correct slice; ships without pinning down a plugin
     config format.
   - **Full auto-setup (mirrors LiveSync):** Artel creates a per-vault subdirectory on the WebDAV
     endpoint (`MKCOL`) and generates a plugin-specific one-tap setup URI/config, the same shape as
     `internal/livesync/builder.go` does for CouchDB. Blocked on deciding (1) first.
3. **Where WebDAV fits the existing sync-vs-binary-storage split.** CouchDB and S3 already play
   two different roles that shouldn't be conflated for WebDAV: CouchDB is the *sync* backend
   (LiveSync replicates against it); S3 is one of two *binary/attachment* stores, abstracted behind
   `storage.BinaryStore` (`internal/storage/binary_store.go`), the other implementation being
   `couchdb.BinaryStoreAdapter`, selected per-vault via `internal/clients/vaultbucket/vaultbucket.go`
   based on `domain.Vault.UseCouchDBForBinaries`/`S3InstanceUuid`. WebDAV as a *sync* backend (item
   2 above) is a separate concern from WebDAV as a third `BinaryStore` implementation for
   attachments — the latter is a much smaller, self-contained addition (implement the interface,
   no new domain/vault-schema changes) that could ship independently of the sync-backend question.

## Why this isn't a MoM `http` tool

MoM's `http` action (`internal/domain/mom.go`'s `HttpAction`, executed by
`internal/service/v1/mcp/executors/http.go`) was considered as a way to avoid a bespoke Go client —
per `pkg/mom_examples/README.md`, it supports arbitrary methods/headers with `__secrets.*`
credential interpolation, so a basic-auth `PUT`/`GET` is technically expressible. Rejected because:

- MoM is JSON-oriented (request/response shaping assumes JSON), with no story for WebDAV's
  `PROPFIND`/XML multistatus responses (directory listing, existence checks).
- `ExecuteToolForKey`/`ExecuteToolWithSecrets` (`internal/service/v1/mom/mom.go`) are callable from
  arbitrary Go code, not just tract/agent paths — so reachability isn't the blocker — but per this
  project's CLAUDE.md, MoM is documented as the pattern for *outbound declarative REST calls for
  LLM tool-calling*, not general-purpose infra provisioning. Vault provisioning is transactional
  Postgres+backend setup (see `CreateVault`'s `txManager.Execute`), which doesn't fit MoM's
  single-tool-call execution model either.
- `ToolAction` is a closed three-way union (`imap`/`smtp`/`http`) and CLAUDE.md explicitly frames
  adding a fourth action type as "rare — discuss before adding." A `webdav` action type was
  considered and rejected for the same JSON/XML mismatch reason as the `http` reuse above.

**Conclusion:** a bespoke `internal/clients/webdav` package (per CLAUDE.md's own carve-out for
protocols that don't fit `http`/`imap`/`smtp`) is the right shape, following `docs/go-style.md`
conventions (verb-first error wrapping, `New(deps) *Client` constructor, context-first I/O
methods) — same pattern as `internal/clients/s3`/`internal/clients/couchdb`.

## Open questions / not decided

- Which WebDAV-sync Obsidian plugin (if any) to target for auto-setup, and whether its config
  format is even importable/URI-shareable the way LiveSync's is — unresearched.
- Credential-only MVP vs. full auto-provisioning (item 2 above) — no decision; credential-only is
  the cheaper validation step if this gets picked up.
- Whether WebDAV-as-`BinaryStore` (attachments) and WebDAV-as-sync-backend should ship as one
  feature or two — they're independent enough to split, per item 3 above.
- Whether a WebDAV BYOK connection (once it exists) should sync into an owned pool row the way the
  S3/CouchDB BYOK work does for `couch_instances`/`s3_instances` (see `migrations/066_*` and
  `vault.go`'s `PickForUser` resolver) — depends on the sync-backend decision above; a
  credential-only MVP has no pool table to sync into.
