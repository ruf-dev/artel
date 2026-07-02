# Artel — About the Project

## What is Artel?

Artel is a service that gives you your own private, cloud-based notebook — and then connects that notebook to the rest of your digital life.

At its heart is the **vault**: a personal knowledge base built on [Obsidian](https://obsidian.md/), the popular note-taking app used by writers, engineers, students, and teams to organize their thoughts as linked notes. Normally, if you want to access the same Obsidian notes from your phone, laptop, and tablet, or share them with a teammate, you need to set up and maintain your own server — a technical, error-prone chore most people never attempt. Artel does that setup for you, instantly, with one click.

But Artel is more than sync. Once you have a vault, it becomes a hub that AI assistants like Claude can plug into — reading your notes, writing new ones, and reaching out to other tools (email, GitLab, and more) on your behalf, all from inside the same conversation.

---

## Why the Vault Is the Core of Everything

The vault isn't just a feature — it's the foundation everything else is built on:

- **It's yours, and only yours.** Each vault gets its own isolated storage. Nothing is mixed with other users' data, and sharing a vault with a teammate never means sharing your password.
- **It's always in sync.** Write a note on your phone, see it instantly on your laptop. No manual exporting, no emailing files to yourself.
- **It's the context AI works from.** When you connect Claude to a vault, every other integration — email, project boards, spreadsheets — plugs into that same shared context. Your vault becomes the memory an AI assistant reasons over.
- **It's simple to set up.** A guided "Fast Setup" hands you a ready-made prompt that configures the technical sync settings for you — no manual server configuration, ever.

In short: the vault is the trustworthy, private foundation. Everything else Artel does is about making that foundation more useful.

---

## Rich Integrations: Your Vault, Connected to Everything Else

Artel is built around a flexible connector system (internally called "MoM"). Think of it as a universal adapter: instead of building a one-off, hard-coded connection for every new service, Artel describes each integration as a simple, structured recipe — what it can do, and how to do it. This means new integrations can be added quickly and safely, without reinventing how AI talks to the outside world each time.

**What's connected today:**
- **Email** — Claude can list your folders, read messages, and send replies, entirely within the context of your vault.
- **GitLab** — connect a GitLab account so Claude can look at issues, boards, and project activity, and react to GitLab events pushed to Artel.
- **Task trackers & spreadsheets** — the groundwork for linking project-tracking tools and spreadsheets into the same assistant workflow is already in place, with more surfacing in the UI soon.

Every integration is permission-gated and credential-isolated: secrets are encrypted and never shared between services, and access can be revoked at any time with a single click.

**Where this is headed:** because integrations all follow the same recipe format, Artel is designed to grow a library of them quickly — Slack, Jira, calendars, and other everyday tools are natural next additions. Longer-term, the roadmap includes a **workflow layer**: chaining tools together so the output of one step (say, "find unread urgent emails") can automatically feed into the next ("create a GitLab issue for each one") — without a person manually relaying information between apps.

---

## Who Uses Artel, and How

### The Individual Note-Taker
Maria journals and plans her freelance projects in Obsidian. She used to worry about losing her notes if her laptop died. With Artel, she creates a vault in seconds, gets a stable sync link, and now her notes are backed up and available on her phone the moment she writes them — no server admin skills required.

### The AI Power User
David wants Claude to actually know his notes instead of him copy-pasting context into every chat. He generates an MCP key for his vault, drops it into Claude's settings, and from then on Claude can search, read, and write notes directly — turning his vault into Claude's long-term memory for his projects.

### The Team Lead
Priya runs a small product team and wants a shared knowledge base without emailing documents around or handing out a shared login. She creates one vault, invites her teammates by username, and each person gets their own private credentials into the same shared space — everyone stays in sync, nobody shares a password.

### The Connected Operator
Jonas manages both his inbox and his team's GitLab backlog and doesn't want to keep switching tabs. He connects his email and GitLab account to his vault, then asks Claude — in one conversation — to summarize unread emails and check on stalled GitLab issues. Both tools answer from the same assistant, grounded in the same notes.

### Looking Ahead: The Automator
In the near future, a user like Jonas won't need to ask Claude to check email and then separately ask it to open a GitLab issue — he'll define a workflow once ("when an urgent client email arrives, open a tracked issue and drop a summary note in my vault"), and Artel will run that chain automatically. This is the direction the platform is heading: from "AI that answers when asked" to "AI that quietly keeps your tools in sync with each other."

---

## In One Sentence

Artel gives you a private, always-synced knowledge vault, then plugs an AI assistant into it and the tools around it — so your notes, your inbox, and your project boards all become one connected, conversational workspace.