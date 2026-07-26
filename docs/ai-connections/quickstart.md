# Quickstart: Connect Claude to Your Vault

This walks you through generating an access key for one vault and plugging it
into Claude Desktop or Claude Code. It takes about five minutes.

You'll need: an Artel account with at least one vault already created, and
Claude Desktop or Claude Code installed.

## Step 1 — Generate an MCP key

1. Open the Artel web app and go to the vault you want to connect.
2. Open the vault's **MCP** tab.
3. Click **Generate Key**, and give it a name you'll recognize later (e.g.
   `claude-desktop` or `work-laptop`).
4. Artel shows you a long code that starts with `artel_vtk_...`. This is your
   **key** — it's the password-like code your AI tool will use to open this
   vault.

> **Copy it now.** The full code is shown to you exactly once. After you
> navigate away, Artel only remembers a short preview of it (enough to
> recognize the key in your list later, not enough to reuse it). If you lose
> it, don't worry — just revoke that key and generate a new one.

Keep this code as private as your password. Anyone who has it can read and
change the notes in that vault until you revoke it.

## Step 2 — Add it to your AI tool

Pick whichever tool you're setting up.

### Option A: Claude Desktop

Open (or create) the file:

- **macOS:** `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Other platforms:** `~/.config/claude/claude_desktop_config.json`

Add this, replacing the URL with your actual Artel address and pasting in
your key where shown:

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

If the file already has an `"mcpServers"` section with other entries in it,
just add `"artel-vault"` alongside them rather than replacing the whole file.

Save the file, then **fully restart Claude Desktop** (quit and reopen it, not
just close the window).

### Option B: Claude Code (command line)

Run:

```bash
claude mcp add artel-vault \
  --transport http \
  --url https://your-artel-host/mcp \
  --header "Authorization: Bearer artel_vtk_YOUR_RAW_TOKEN"
```

(Again, swap in your actual Artel address and your key.)

Alternatively, if you'd rather edit a config file directly, add this to your
project's `.claude/settings.json`:

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

## Step 3 — Confirm it worked

Start a new conversation in Claude and ask it something like:

> "List the files in my vault."

or

> "What tools do you have available for my vault?"

If it's connected correctly, Claude will call the vault's tools and show you
real file names from your notebook. If Claude says it has no such tool, or
the connection times out, double check:

- The key was pasted in full, with no line breaks or missing characters.
- The Artel address (`https://your-artel-host/mcp`) is correct — it should
  end in `/mcp`.
- You fully restarted Claude Desktop (or re-ran the `claude mcp add` command)
  after saving the config.
- The key hasn't been revoked (check the vault's MCP tab in Artel — revoked
  keys stop working immediately).

## What Claude can do now

Once connected, Claude can, inside that one vault:

- List every file and folder
- Read any note or file (including images and PDFs)
- Create or update notes
- Delete files
- Move or rename files
- List your tags
- Tell you what other services (like email) are connected to your account —
  though using those services from Claude takes one more setup step; see
  [detailed-configuration.md](detailed-configuration.md).

Claude cannot touch any *other* vault with this key, and it cannot change
your Artel account settings, billing, or membership — the key only opens the
one vault it was generated for.

## Next steps

- To connect a second device, teammate, or tool, generate another key the
  same way — see [detailed-configuration.md](detailed-configuration.md) for
  managing multiple keys and revoking ones you no longer use.
- To let Claude use your email or GitLab from the same conversation, see the
  separate integrations guide for connecting those services, then
  [detailed-configuration.md](detailed-configuration.md) for how (and to what
  extent, today) that access reaches a given key.
