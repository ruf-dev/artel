# Detailed Configuration

This covers managing multiple keys, revoking access, exactly what tools an AI
assistant gets and why, and (for the technically curious) the raw API calls
behind the scenes. If you just want to get connected, see
[quickstart.md](quickstart.md) instead.

## Managing multiple keys

You can generate more than one key for the same vault — for example, one for
Claude Desktop on your laptop, one for Claude Code, one for a teammate. Each
key is independent: it has its own name, its own on/off switch, and (in the
future) its own set of granted integrations.

Open a vault's **MCP** tab in Artel to see every key you've generated for
that vault, along with:

- The name you gave it
- A short preview of the code (enough to recognize which key is which — the
  full code is never shown again after creation)
- When it was created

There's no limit baked into this guide on how many keys you can have per
vault — create as many as you need to keep each device or tool separately
revocable.

## Revoking a key

If a device is lost, a teammate leaves, or you just want to rotate a key,
revoke it from the same **MCP** tab. Revoking is immediate and permanent —
the AI tool using that key loses access right away, and there's no way to
"undo" a revoke. If you still need that tool connected, generate a fresh key
and update the tool's configuration with the new code (see
[quickstart.md](quickstart.md) Step 2).

Revoking one key never affects your other keys or other vaults.

## What tools does the AI actually get?

Every valid, non-revoked key gives an AI assistant the same core set of
vault abilities — no extra setup needed beyond having the key:

| Ability | What it does |
|---|---|
| List files | See every file and folder in the vault |
| Read a file | Open any note or file — text notes come back as text; images and PDFs come back as viewable content |
| Write a file | Create a new note/file or update an existing one |
| Delete a file | Remove a file |
| Move a file | Rename or relocate a file |
| List folders | See the folder structure |
| List tags | See every tag used in the vault |
| Get file details | Size, last-modified time, and similar metadata for one file |
| See connections | Ask "what other services are linked to this key" — e.g. it'll report `email` if that's connected and granted |

A couple of practical notes:

- Very large binary files (anything that isn't a Markdown note) are stored in
  a separate file bucket linked to your vault. If your vault doesn't have one
  linked, the AI can still read and write your Markdown notes fine, but
  non-Markdown file operations will fail with a clear "no file storage linked"
  message.
- "Move" can't jump a file between a note and a non-note file type in a single
  step — if the AI needs to do that, it'll read the old file, write a new one
  at the new location, and delete the old one. You don't need to do anything
  differently; this is just how it's handled internally.

## Extra abilities from other integrations

Artel lets you connect other services to your account — today that's email,
with more planned. Once a service is connected **and** that specific key has
been granted access to it, the AI gains extra tools for it. For email, that
currently means things like listing folders, reading messages, and sending
replies — all without leaving your conversation with the AI.

### Current limitation — please read this before relying on it

As of today, connecting a service (like email) to your account is not
automatically enough to make its tools available to a key, and — importantly
— **there is no button in the Artel UI yet to grant a connected service to a
specific key.** The underlying support for "grant this integration to this
key" exists in Artel, but nothing currently switches it on for you. In
practice, this means:

- Calling `connections` from the AI will report a service as connected once
  you've set it up.
- But the AI still won't be able to actually *use* that service's tools (e.g.
  `send_email`) yet, because no key has been granted to it — you'll see a
  "tool not found" style error if the AI tries.

This is a known gap, not a bug on your end — nothing you do in the current
Artel UI can grant email/GitLab access to a key today. If having a specific
key use a connected integration is important to you, keep an eye on Artel's
release notes; this is expected to get a UI in a future update.

## Advanced: talking to the vault directly (no Claude required)

Everything above happens automatically once you've configured Claude
Desktop or Claude Code. If you're comfortable with the command line or a
little scripting, you can also talk to the vault's connection endpoint
directly — useful for testing a key, or wiring up your own script.

The vault's assistant endpoint speaks a standard JSON-based protocol over a
single web address (`/mcp`) with your key sent as an `Authorization: Bearer`
header on every request.

**Check the connection and list tools with curl:**

```bash
curl -X POST https://your-artel-host/mcp \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer artel_vtk_YOUR_RAW_TOKEN" \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{}}}'

curl -X POST https://your-artel-host/mcp \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer artel_vtk_YOUR_RAW_TOKEN" \
  -d '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}'
```

**Read a note with curl:**

```bash
curl -X POST https://your-artel-host/mcp \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer artel_vtk_YOUR_RAW_TOKEN" \
  -d '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"read_file","arguments":{"path":"Daily/2026-05-12.md"}}}'
```

**Python example:**

```python
import httpx

MCP_URL = "https://your-artel-host/mcp"
TOKEN = "artel_vtk_YOUR_RAW_TOKEN"

headers = {
    "Authorization": f"Bearer {TOKEN}",
    "Content-Type": "application/json",
}

def rpc(method: str, params: dict, req_id: int = 1):
    payload = {"jsonrpc": "2.0", "id": req_id, "method": method, "params": params}
    resp = httpx.post(MCP_URL, json=payload, headers=headers)
    resp.raise_for_status()
    return resp.json()

# Initialize
rpc("initialize", {"protocolVersion": "2025-06-18", "capabilities": {}})

# List files
result = rpc("tools/call", {"name": "list_files", "arguments": {}})
print(result["result"]["content"][0]["text"])

# Write a note (Markdown path -> plain text content)
rpc("tools/call", {
    "name": "write_file",
    "arguments": {"path": "AI/hello.md", "content": "# Hello from Python"}
})
```

### Creating and managing keys via the API

If you're scripting key management instead of using the MCP tab, the same
actions are available as plain web API calls (using your normal Artel login
session, not the MCP key):

**Create a key:**

```bash
curl -X POST https://your-artel-host/api/mcp/VAULT_UUID/mcp-keys \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_SESSION_TOKEN" \
  -d '{"name": "my-llm-key"}'
```

The response includes a `rawToken` field — that's the full `artel_vtk_...`
code, shown only in this response and never retrievable again.

**List keys for a vault:**

```bash
curl https://your-artel-host/api/mcp/VAULT_UUID/mcp-keys \
  -H "Authorization: Bearer YOUR_SESSION_TOKEN"
```

**Revoke a key:**

```bash
curl -X DELETE https://your-artel-host/api/mcp/VAULT_UUID/mcp-keys/KEY_ID \
  -H "Authorization: Bearer YOUR_SESSION_TOKEN"
```

This is exactly what the MCP tab in the Artel UI does behind the scenes — the
UI is the recommended way to manage keys day to day; the API is here for
anyone automating setup.
