# Getting MCP Credentials

Artel exposes an MCP server at `/mcp`. Access requires a **vault API key** (bearer token) that is tied to a specific vault.

## Step 1 — Get a Vault ID

You need the UUID of the vault you want to connect to. The vault UUID is displayed in the Artel UI after creating a vault, or returned in the vault list API.

## Step 2 — Create an MCP Key

Call the REST API:

```http
POST /api/mcp/{vault_id}/mcp-keys
Content-Type: application/json

{"name": "my-llm-key"}
```

**curl example:**

```bash
curl -X POST https://your-artel-host/api/mcp/VAULT_UUID/mcp-keys \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_SESSION_TOKEN" \
  -d '{"name": "my-llm-key"}'
```

**Response:**

```json
{
  "key": {
    "id": "3e4a1b2c-...",
    "vaultId": "VAULT_UUID",
    "name": "my-llm-key",
    "keyPreview": "artel_vtk_3",
    "createdAt": "2026-05-12T10:00:00Z"
  },
  "rawToken": "artel_vtk_3e4a1b2c..._secrethex..."
}
```

**Save `rawToken` immediately — it is shown only once and cannot be recovered.**

The `keyPreview` field (first 12 characters) is stored permanently and lets you identify the key later.

## Step 3 — List or Revoke Keys

List keys for a vault:

```bash
curl https://your-artel-host/api/mcp/VAULT_UUID/mcp-keys \
  -H "Authorization: Bearer YOUR_SESSION_TOKEN"
```

Revoke a key:

```bash
curl -X DELETE https://your-artel-host/api/mcp/VAULT_UUID/mcp-keys/KEY_ID \
  -H "Authorization: Bearer YOUR_SESSION_TOKEN"
```

## Token Format

MCP tokens follow the pattern `artel_vtk_{uuid_hex}_{secret_hex}`. The UUID encodes the key ID so the server can look up the key without a full table scan. The secret is bcrypt-hashed at rest (cost 12).
