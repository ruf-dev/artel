# MCP Connection Examples

The Artel MCP server speaks **JSON-RPC 2.0 over HTTP POST** at `/mcp`.  
Authentication: `Authorization: Bearer artel_vtk_...` on every request.

---

## Claude Desktop

Add to `~/.config/claude/claude_desktop_config.json` (macOS: `~/Library/Application Support/Claude/claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "artel-vault": {
      "url": "https://your-artel-host/mcp",
      "headers": {
        "Authorization": "Bearer artel_vtk_YOUR_RAW_TOKEN"
      }
    }
  }
}
```

Restart Claude Desktop. The vault tools (`list_notes`, `read_note`, `write_note`, etc.) will appear automatically.

---

## Claude Code (CLI)

```bash
claude mcp add artel-vault \
  --transport http \
  --url https://your-artel-host/mcp \
  --header "Authorization: Bearer artel_vtk_YOUR_RAW_TOKEN"
```

Or add directly to `.claude/settings.json` in your project:

```json
{
  "mcpServers": {
    "artel-vault": {
      "type": "http",
      "url": "https://your-artel-host/mcp",
      "headers": {
        "Authorization": "Bearer artel_vtk_YOUR_RAW_TOKEN"
      }
    }
  }
}
```

---

## Raw HTTP (curl)

**Initialize session:**

```bash
curl -X POST https://your-artel-host/mcp \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer artel_vtk_YOUR_RAW_TOKEN" \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{}}}'
```

**List available tools:**

```bash
curl -X POST https://your-artel-host/mcp \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer artel_vtk_YOUR_RAW_TOKEN" \
  -d '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}'
```

**Read a note:**

```bash
curl -X POST https://your-artel-host/mcp \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer artel_vtk_YOUR_RAW_TOKEN" \
  -d '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"read_note","arguments":{"path":"Daily/2026-05-12.md"}}}'
```

---

## Python (httpx)

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

# List notes
result = rpc("tools/call", {"name": "list_notes", "arguments": {}})
print(result["result"]["content"][0]["text"])

# Write a note
rpc("tools/call", {
    "name": "write_note",
    "arguments": {"path": "AI/hello.md", "content": "# Hello from Python"}
})
```

---

## Available Tools

| Tool | Description | Required args |
|------|-------------|---------------|
| `list_files` | List all files in the vault (notes and binary files) | — |
| `read_file` | Read any file by path (text or base64 binary) | `path` |
| `write_note` | Create or update a note | `path`, `content` |
| `delete_file` | Delete any file by path | `path` |
| `move_file` | Move/rename any file | `old_path`, `new_path` |
| `list_folders` | List all folders | — |
| `list_tags` | List all tags | — |
| `get_note_metadata` | Get id, rev, mtime, ctime, size | `path` |
| `connections` | List the MoMs (e.g. `email`) connected to this key | — |

Text tools return `{"content": [{"type": "text", "text": "..."}]}`; `read_file` on an image/PDF
returns an `image`/`document` content block with base64 `data` instead. Beyond these built-in
tools, a key may also expose extra tools from connected MoMs (e.g. `send_email`) — see
[tools.md](tools.md) for exactly when those appear.
