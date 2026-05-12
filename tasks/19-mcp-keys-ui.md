# Task 19 — MCP Keys UI

## Goal

Add an MCP Keys section to the vault card / vault detail view. Users can list, create (with one-time token display), and revoke keys for each vault.

## Context

- Framework: React + TypeScript, Bun bundler.
- Existing vault card: `pkg/client/ArtelUI/src/pages/home/VaultCard.tsx` — a compact card with vault name, CouchDB URL, and copy button.
- Style: CSS modules (`.module.css`), same pattern as `VaultCard.module.css`.
- API calls should follow the same pattern as existing API wrappers in `pkg/client/ArtelUI/src/app/api/`.
- The MCP keys REST API is exposed via grpc-gateway at:
  - `POST   /api/vaults/{vault_id}/mcp-keys`       — body `{ "name": "..." }` → `{ "key": {...}, "raw_token": "artel_vtk_..." }`
  - `GET    /api/vaults/{vault_id}/mcp-keys`         → `{ "keys": [{ "id", "vault_id", "name", "key_preview", "created_at" }] }`
  - `DELETE /api/vaults/{vault_id}/mcp-keys/{key_id}` → `{}`

## Files to Create

### `pkg/client/ArtelUI/src/app/api/artel/mcp_keys.ts`

Plain `fetch`-based API helpers:

```ts
export interface McpKeyInfo {
  id: string
  vaultId: string
  name: string
  keyPreview: string
  createdAt: string
}

export async function listMcpKeys(vaultId: string): Promise<McpKeyInfo[]>
export async function createMcpKey(vaultId: string, name: string): Promise<{ key: McpKeyInfo; rawToken: string }>
export async function revokeMcpKey(vaultId: string, keyId: string): Promise<void>
```

Include `Authorization` header using the session token from wherever it is stored in the app (look at how existing API calls do auth — check `pkg/client/ArtelUI/src/` for the auth token pattern).

### `pkg/client/ArtelUI/src/pages/home/McpKeysSection.tsx`

A React component `McpKeysSection({ vaultId: string })`:

- On mount: fetch and display the list of MCP keys via `listMcpKeys`.
- Each key row: show `name`, `keyPreview`, `createdAt`, and a "Revoke" button.
- "Revoke" calls `revokeMcpKey`, refreshes the list, shows confirmation.
- "Add Key" button: shows an inline form with a name input and "Create" button.
- On create success: show the `rawToken` in a highlighted box with a "Copy" button and a warning "Save this token — it will not be shown again."
- After copying or dismissing: token disappears from UI.

### `pkg/client/ArtelUI/src/pages/home/McpKeysSection.module.css`

CSS module for the component. Keep the same visual style as `VaultCard.module.css`.

## Integration

Extend `VaultCard.tsx` to render `<McpKeysSection vaultId={vault.id} />` below the existing footer, collapsed by default behind a toggle ("MCP Keys ▾").

## Verification

- `bun run build` (from `pkg/client/ArtelUI/`) must succeed with no TypeScript errors.
- Visually: key list loads, create shows raw token once, revoke removes the key.
