# Artel Quick Start

Welcome to Artel — a managed cloud platform for your Obsidian notes. Artel provisions a private,
CouchDB-backed vault for you in seconds, so you can sync notes across every device using the
[Obsidian LiveSync](https://github.com/vrtmrz/obsidian-livesync) plugin, without ever running or
maintaining a server yourself. You can also browse and edit your notes directly in the browser,
no Obsidian install required.

## Getting Started

1. **Sign up or log in.** Create an account with your email and password, or sign in with
   Telegram — whichever your instance has enabled.
2. **Create a vault.** From your dashboard, click **Create Vault** and give it a name. Artel
   provisions a dedicated, isolated database for it right away.
3. **Connect Obsidian.** Open your new vault's card and copy the CouchDB connection URL, then
   paste it into the Obsidian LiveSync plugin's settings. Prefer not to type it in by hand? Use
   **Fast Setup** to copy a ready-made prompt for Claude that configures LiveSync for you
   automatically.
4. **Start writing.** Once connected, every note you create syncs automatically across all of
   your devices — desktop, mobile, wherever Obsidian runs.

## Sharing and Collaboration

Vaults aren't limited to a single person. Invite teammates as members with their own accounts —
no need to share passwords or credentials. Each vault also supports **publishing**: pick a public
slug and anyone with the link can browse a read-only, web-rendered view of your notes, exactly
what's serving this page.

## Beyond Sync: AI and Automation

Artel goes further than plain file sync:

- **MCP integration** — generate a scoped API key for a vault and let Claude (or any MCP-capable
  AI tool) read, search, and write your notes directly, with fine-grained, revocable access.
- **Connections** — link email, Trello, GitLab, and other external services so your AI tools can
  act across your whole workflow, not just your notes.
- **Tracts** — build automated workflows that react to triggers (a webhook, a schedule, an
  incoming email) and chain steps together, from simple notifications to multi-step pipelines.

## Where to Go Next

Once you're signed in, your dashboard is the starting point for creating vaults, managing
members, generating MCP keys, and wiring up connections and tracts. If you're an administrator,
the admin panel lets you manage the underlying infrastructure (CouchDB instances, registration
policy, and this very default docs page) for your whole Artel instance.
