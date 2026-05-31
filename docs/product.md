# Artel — Product Overview

## What Is Artel?

Artel is a managed cloud platform for provisioning and sharing personal knowledge vaults built on [Obsidian LiveSync](https://github.com/vrtmrz/obsidian-livesync). Instead of running your own CouchDB server and configuring it by hand, Artel gives you a fully isolated vault in seconds — accessible from any Obsidian client and from AI tools like Claude.

---

## Pain Points Solved

| Pain | How Artel Fixes It |
|------|--------------------|
| Setting up CouchDB for Obsidian LiveSync requires server access and manual configuration | Artel provisions a ready-to-use vault with one button click |
| Multi-device Obsidian sync is fragile without a self-hosted backend | Each vault gets a stable CouchDB URL that all your devices connect to |
| Sharing a knowledge base with teammates means giving them your credentials | Vault membership with per-user accounts — no credential sharing |
| Letting Claude or other AI tools read your notes is not straightforward | MCP keys give fine-grained, revocable API access per vault |
| Email is separate from your knowledge base | Connect email accounts and let Claude read, search, and send emails in your vault context |
| Obsidian LiveSync setup is a multi-step tutorial nobody wants to follow | Fast Setup copies the exact Claude prompt that configures LiveSync for you automatically |

---

## User Roles

**Vault Owner / Member** — a regular user who creates vaults, invites teammates, generates MCP keys, and connects email accounts.

**Administrator** — registers and manages the underlying CouchDB server instances that power all vaults. Admins control the infrastructure pool that vault creation draws from.

---

## Feature Areas

### Vault Management
- Create a named vault → instantly provisioned on a CouchDB instance
- View the CouchDB connection URL for Obsidian LiveSync setup
- Rename or delete a vault
- Share a vault with other Artel users (add / remove members)
- Vault status indicator (active / provisioning)

### MCP Keys — AI Tool Access
- Generate a bearer token (`artel_vtk_…`) scoped to a single vault
- Token is shown once and never stored in plaintext
- List all active keys for a vault with previews
- Revoke any key at any time

### Email Accounts
- Connect an email account via IMAP/SMTP credentials
- List email folders and messages from Claude
- Read full email content
- Send emails — all from within your vault context
- Permission-gated: enabled per user by admin

### Fast Setup — Claude Configuration Helper
- Each vault card has a "Fast Setup" option
- Shows the pre-built Claude prompt that configures LiveSync automatically
- One-click copy — paste it into Claude and your vault is wired up

### Authentication
- Email + password registration and login
- Telegram OAuth (sign in with Telegram account, no password needed)
- Sessions are token-based and expire automatically

### Admin Panel
- List all registered CouchDB instances (URL, status)
- Add a new CouchDB server (URL, admin credentials)
- Edit or remove instances
- Restricted to administrator accounts only

---

## User Journeys

### New User: First Vault
1. Sign up with email/password or Telegram
2. Click **Create Vault**, enter a name
3. Open the vault card → copy the CouchDB URL
4. Configure Obsidian LiveSync with the URL — or use **Fast Setup** to let Claude do it

### Power User: Connect Claude
1. Open vault card → MCP tab → **Generate Key**
2. Copy the `artel_vtk_…` token
3. Add it to Claude's MCP server config (`artel` server)
4. Claude can now read notes, write notes, list files, and browse folders in the vault

### Team Lead: Shared Knowledge Base
1. Create a vault for the team
2. Add teammates by their Artel username
3. Each member gets their own CouchDB credentials — no shared passwords
4. All members sync via Obsidian to the same vault

### Email-Integrated Workflow
1. Admin enables email access for the user
2. User adds an email account (IMAP host, SMTP host, credentials)
3. Claude lists folders and reads emails in the vault context
4. Claude can send email replies without leaving the conversation
