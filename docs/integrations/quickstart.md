# Quickstart — Connecting Email and GitLab

## Before you start

You'll need:

- An Artel account and at least one vault.
- Ideally, an MCP key already set up so an AI assistant is connected to that
  vault — see [AI Connections](../ai-connections/general-info.md) if you
  haven't done that yet. You can connect Email/GitLab first and set that up
  after, but the AI can't use the connection until both exist.

Both connections live in the same place: open the **Connections** page and
stay on the **External Connections** tab (the default tab).

## Connect an email account

1. On the Connections page, find the **Email** card and click it to open the
   "Email" dialog, then choose **Add account**.
2. Enter the email address you want to connect.
3. Artel tries to fill in the mail server settings for you (the IMAP and
   SMTP host/port fields) based on the address you typed — for well-known
   providers this happens automatically. If it doesn't recognize the
   provider, enter the IMAP host/port (for reading mail) and SMTP host/port
   (for sending mail) yourself; your provider's help pages will list these.
4. Enter the account's password.
   - **Many providers require an "app-specific password" instead of your
     normal login password** for this kind of connection (Gmail and Yandex
     are common examples). If Artel recognizes your provider, it shows a
     link to the page where you can generate one.
5. Click **Check** (or the verify button next to it) to confirm Artel can
   actually log in with what you entered.
6. Once verification succeeds, click **Add account**.

Your inbox is now connected. You can add more than one email account if you
need to.

## Connect a GitLab account

1. On the Connections page, find the **GitLab** card and click it to open
   the "GitLab" dialog.
2. **Instance URL** — leave this blank (or click the button to fill in the
   default) if you use `gitlab.com`. If your organization runs its own
   GitLab server, enter that server's URL instead.
3. **Personal access token** — GitLab uses a token instead of your password
   for this kind of connection. If you don't have one yet, click the
   generated link next to the field to jump straight to GitLab's "create a
   personal access token" page, create one there, and paste it back in.
4. Click the verify/check button to confirm the token works against that
   GitLab instance.
5. Once it's verified, click **Connect**.

Your GitLab account is now connected — Claude can create and list merge
requests and issues, and comment on them, within your vault's conversations.

## What your AI assistant can do once connected

With Email connected: list folders, read messages, and send replies.

With GitLab connected: look at issues and merge requests, comment on them,
and react to project activity.

## One thing to know: connecting isn't quite the same as "the AI can use it"

Connecting an account makes it available in principle. For a *specific* AI
assistant (a specific MCP key) to actually use it, that key also needs to be
granted access to the connection. Right now, Artel doesn't yet have a
control in the app for doing that granting yourself — see
[detailed-configuration.md](detailed-configuration.md) for the current
state of this. It doesn't change anything about the steps above; it's worth
knowing if a newly connected account doesn't immediately show up for your
assistant to use.
