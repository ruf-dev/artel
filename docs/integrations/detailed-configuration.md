# Managing Connections, Permissions, and BYOK

## Where everything lives

All of this is on the **Connections** page, which currently has a couple of
tabs: **External Connections** (Email, GitLab, and similar service
accounts — see [quickstart.md](quickstart.md)) and **BYOK** (your own AI
provider keys, covered below).

## Editing a connection

- **Email**: open the Email card, and next to an existing account you'll
  find an edit option. You can update the email address or server settings
  there; if you leave the password field blank, Artel keeps the password
  you already saved and only updates the other fields.
- **GitLab**: open the GitLab card to see the connected account's details.

## Removing a connection

Every connection type has a **Disconnect** (or **Remove**) button, guarded
by a confirmation prompt so you don't remove one by accident. Removing a
connection deletes the stored credentials. You can reconnect the same
account at any time by going through the connect steps again — nothing
about your vault or notes is affected either way.

## Plans and access

Some integrations are part of your Artel plan. If a "Connect" option looks
disabled or a card is missing entirely, that feature may not be included in
your current plan or may need to be turned on for your account — check with
your workspace administrator.

## A current limitation: connecting vs. granting

Connecting a service (Email, GitLab, or a BYOK key) stores the credentials
and makes the connection *exist*, but a given AI assistant key only gets to
use a connection once that key has specifically been granted access to it.
This lets a vault owner, for example, issue one AI key for themselves and
another for a teammate, and decide separately what each one can touch.

Right now, Artel doesn't yet have an in-app control for doing that granting
— it's a known gap the team is aware of and working on. Practically, this
means a freshly connected integration may not show up for your AI assistant
to use immediately. There's nothing you need to configure differently in
the meantime; it isn't something you did wrong.

## Bring Your Own Key (BYOK)

### What it's for

BYOK lets you hand Artel your own Claude (Anthropic) API key. Rather than
Artel's automation features drawing on one shared, metered AI account,
they draw on your key — usage happens under your own Anthropic account and
billing, not Artel's.

To be clear about what this does *not* do: it's unrelated to the AI
assistant (Claude Desktop, Claude Code, etc.) that reads and writes your
vault through an MCP key — that connection already brings its own model
access. BYOK is specifically for Artel's own automation features that need
to call an AI model as one step in a larger flow.

### What works today

You can already:

- Go to the Connections page's **BYOK** tab and open the Claude (Anthropic)
  card.
- Enter your Anthropic API key (and optionally a custom base URL and default
  model, for advanced setups).
- Have Artel verify the key actually works before saving it.
- See the connection's status (a preview of the key, the model it's using)
  and disconnect it at any time, the same as any other connection.

### What's still coming

- Support for OpenAI-format keys — right now only Claude (Anthropic) keys
  are supported; an OpenAI option is planned but not available yet.
- "Bring your own storage" connections (your own CouchDB, S3-compatible
  bucket, or WebDAV server) — these appear on the BYOK tab today only as
  "coming soon" placeholder cards with no functionality behind them yet.
- Cost estimates in plain dollar terms — usage is tracked, but a friendly
  cost breakdown isn't built yet.

Because this piece is newer and smaller in scope than Email/GitLab, treat it
as an early capability: useful if you specifically want AI automations
running on your own key today, but not yet a broad, general-purpose feature
the way Email and GitLab connections are.
