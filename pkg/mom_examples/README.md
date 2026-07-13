MCP of MCP (MoM)

## Format

Each MoM record is a JSON object with `name`, `author`, and `tools` (array of tool definitions).

**Tool structure:**
- `api_description` — what the LLM sees: name, description, input schema (standard MCP tool format)
- `action` — how to execute it: one key = the protocol/client type (`imap`, `smtp`, `http`), value = executor definition

**Dispatch rules:**
- The `action` key drives executor selection — unmarshal the value into the matching executor struct
- For `imap` / `smtp`: `api_description.properties` map directly to Go client function parameters; no remapping needed inside the action
- For `http`: inputs are agile — use `${{params.name}}` to interpolate tool inputs into URL path, query, or headers; use `__secrets.field` to inject values from the linked `external_connections.credentials_enc`

---

## Email MoM

```json
{
  "name": "email",
  "author": "Artel",
  "description": "Email handler. Can send, read, list emails in connected account.",
  "tools": [
    {
      "api_description": {
        "name": "list_email_folders",
        "description": "List IMAP folders in the connected email account",
        "properties": {},
        "required": []
      },
      "action": {
        "imap": {
          "operation": "list_folders"
        }
      }
    },
    {
      "api_description": {
        "name": "list_emails",
        "description": "List recent emails from the inbox",
        "properties": {
          "limit": {
            "type": "integer",
            "description": "Max number of emails to return (default 20)"
          },
          "after_uid": {
            "type": "string",
            "description": "Only return emails with UID greater than this value, oldest-first. Pass \"0\" to start from the oldest email in the inbox. To page forward, pass the id of the last (newest) email from the previous page. Mutually exclusive with before_uid."
          },
          "before_uid": {
            "type": "string",
            "description": "Only return emails with UID less than this value, also oldest-first. To page further back, pass the id of the first (oldest) email from the previous page. Mutually exclusive with after_uid."
          }
        },
        "required": []
      },
      "action": {
        "imap": {
          "operation": "list_messages"
        }
      }
    },
    {
      "api_description": {
        "name": "read_email",
        "description": "Read the full content of an email by its UID",
        "properties": {
          "id": {
            "type": "string",
            "description": "Email UID from list_emails"
          }
        },
        "required": ["id"]
      },
      "action": {
        "imap": {
          "operation": "fetch_message"
        }
      }
    },
    {
      "api_description": {
        "name": "send_email",
        "description": "Send an email from the connected account",
        "properties": {
          "to":      { "type": "string", "description": "Recipient email address" },
          "subject": { "type": "string", "description": "Email subject" },
          "body":    { "type": "string", "description": "Email body (plain text)" }
        },
        "required": ["to", "subject", "body"]
      },
      "action": {
        "smtp": {
          "operation": "send"
        }
      }
    }
  ]
}
```

---

## HTTP Action Example (Trello)

Demonstrates `http` executor features: path param interpolation, headers, query params with secrets.

```json
{
  "name": "trello",
  "author": "community",
  "description": "Trello project management integration",
  "tools": [
    {
      "api_description": {
        "name": "get_board_members",
        "description": "Get all members of a Trello board",
        "properties": {
          "board_id": {
            "type": "string",
            "description": "Trello board ID"
          }
        },
        "required": ["board_id"]
      },
      "action": {
        "http": {
          "method": "GET",
          "url": "https://api.trello.com/1/boards/${{params.board_id}}/members",
          "headers": {
            "Accept": "application/json"
          },
          "query": {
            "key":   "__secrets.api_key",
            "token": "__secrets.api_token"
          },
          "credentials": "trello"
        }
      }
    }
  ]
}
```

**`url`** — `${{params.name}}` interpolates tool call inputs into URL path segments.

**`headers`** — static strings or `__secrets.*` values (e.g. `"Authorization": "Bearer __secrets.token"`).

**`query`** — appended as URL query parameters; `__secrets.*` resolved from linked credentials.

**`credentials`** — value of `external_connections.provider` to look up. The service decrypts that connection's `credentials_enc` JSON and exposes its fields as `__secrets.*`. Example Trello credentials JSON: `{ "api_key": "xxx", "api_token": "yyy" }`.

---

## Implementation Plan

### DB

**Migration 027 — `mcps` table**
```sql
CREATE TABLE mcps (
    name        TEXT        PRIMARY KEY,
    author      TEXT        NOT NULL,
    description TEXT        NOT NULL DEFAULT '',
    tools       JSONB       NOT NULL DEFAULT '[]',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- seed with email MoM above
```

**Migration 028 — `mcp_connectors` (link table)**
```sql
CREATE TABLE mcp_connectors (
    id                     UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    mcp_key_id             UUID        NOT NULL REFERENCES mcp_keys(id) ON DELETE CASCADE,
    mcp_name               TEXT        NOT NULL REFERENCES mcps(name),
    external_connection_id UUID        NOT NULL REFERENCES external_connections(id) ON DELETE CASCADE,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(mcp_key_id, mcp_name)
);
```

No credentials stored here — all secrets stay in `external_connections.credentials_enc`.

### Go layers

1. **Domain** — `McpDefinition`, `McpToolDef`, `ToolAction` (tagged by action key), `EmailConnectorCreds`, `HttpConnectorCreds`
2. **SQL queries + `sqlc generate`** — `mcps.sql` (get by name, list), `mcp_connectors.sql` (insert, list by key, get by key+name, delete)
3. **Repo** — `McpDefinitionsRepo`, `McpConnectorsRepo`
4. **Executors** (`internal/service/v1/mcp/executors/`)
   - `EmailExecutor` — wraps existing `imap.Client` + `smtp.Client`; dispatches by `operation`; reads params from `api_description.properties` values in the tool call
   - `HttpExecutor` — builds `*http.Request`; resolves `${{params.*}}` from tool call inputs and `__secrets.*` from decrypted credentials
5. **MCP service** — `ListToolsForKey(keyId)`, `ExecuteToolForKey(keyId, toolName, params)`, connector CRUD
6. **MCP transport** (`tools.go`) — on `tools/list`: merge DB-driven tools with existing vault tools; on `tools/call`: try generic dispatch first, fall back to hardcoded handlers during migration; remove hardcoded email tools once `mcp_connectors` replaces `email_accounts`
