# Integrations — What This Is

## The short version

Your Artel vault doesn't have to live in a bubble. "Integrations" (also
called **Connections**) is the feature that links outside accounts — your
email inbox, a GitLab project — to your vault, so an AI assistant working
inside that vault can also read from and act on those services, in the same
conversation.

Instead of you copying an email into the chat so Claude can read it, or
switching to another tab to check a GitLab issue, you connect the account
once and the assistant can reach it directly whenever it's relevant.

This is a companion feature to [AI Connections](../ai-connections/general-info.md)
— you'll need a vault and an MCP key set up first (see that doc) before an AI
assistant can use anything described here.

## Why you'd want this

- **One conversation, several tools.** Ask about unread emails and a stalled
  GitLab issue in the same chat instead of tab-switching between your inbox,
  GitLab, and Claude.
- **Less copy-pasting.** The assistant reads the email or the issue directly
  instead of you pasting the text in yourself.
- **You stay in control.** Every connected account is listed in one place,
  its credentials are encrypted and used only for that connection, and you
  can disconnect it with a single click at any time.

## What's connected today

- **Email** — connect an inbox and Claude can list your folders, read
  messages, and send replies, all from inside your vault's conversation.
- **GitLab** — connect a GitLab account and Claude can look at issues, merge
  requests, and project activity, and react to events GitLab sends over.

Step-by-step setup for both is in [quickstart.md](quickstart.md).

More integrations (task trackers, spreadsheets, and others) are being built
on the same underlying system and will appear in the Connections page as
they're finished — there's nothing to set up for those yet.

## Bring Your Own Key (BYOK) — a newer, early-stage piece

Artel also lets you connect your own Claude (Anthropic) API key. This is a
different kind of connection — it's not about reading your inbox or a
project board, it's about letting Artel's own automation features run an AI
step using an API key that's yours, rather than a shared one. It's newer
than Email/GitLab and still filling in — see
[detailed-configuration.md](detailed-configuration.md) for exactly what
works today.

## Who this is for

Anyone who already has a vault and wants their AI assistant to do more than
read and write notes — in particular, if you regularly juggle an inbox and a
GitLab project alongside your notes, connecting both means one conversation
can touch all three.

## A word on security

Connecting an account means Artel stores the credentials needed to talk to
that service (a password, a token, an API key) in encrypted form. They're
never shown back to you in full, never shared between different
integrations, and never used for anything other than the connection you set
up. Removing a connection deletes the stored credentials — you can always
reconnect later.
