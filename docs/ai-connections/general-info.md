# AI Connections — What This Is

## The short version

Your Artel vault is a private, cloud-hosted notebook. "AI Connections" is the
feature that lets an AI assistant — Claude Desktop, Claude Code, or any other
tool that speaks the same connection language — actually open that notebook
and work in it: read your notes, write new ones, organize folders, and (if
you've connected other services like email) act on those too, all from inside
your conversation with the AI.

Without this feature, you'd have to keep copy-pasting your notes into the chat
every time you wanted the AI to know about them. With it, the AI reaches into
your vault directly.

## How it works, in plain terms

1. You generate a special access code for one vault — called an **MCP key**.
   ("MCP" is just the name of the connection standard; you don't need to know
   more than that it's the thing that lets AI tools plug into apps like
   Artel.) Think of it like an app password: a long random code that proves
   "this AI assistant is allowed to open this specific vault."
2. You paste that code into your AI tool's settings.
3. From then on, that AI tool can list, read, write, move, and delete files in
   that vault, list your tags, and see what other services (like email) are
   connected — but only what you've allowed, and only for that one vault.
4. You can look at every key you've issued, and switch one off at any time.
   Once switched off, that AI tool can no longer reach the vault.

Nothing here touches your Artel password, and a key only ever gives access to
one vault at a time — connecting a key for "Notes" doesn't expose your
"Journal" vault.

## Why you'd want this

- **Stop copy-pasting context.** Instead of explaining your projects to Claude
  every session, Claude can just read your notes directly.
- **Let the AI take notes for you.** Ask Claude to summarize a conversation or
  a document, and have it write the summary straight into your vault.
- **Keep your notes organized without doing it yourself.** The AI can create
  folders, move files around, and tag things based on what you ask it to do.
- **One AI, multiple tools.** If you've also connected other services (like
  email) to your Artel account, granting them to a key lets the same AI
  conversation touch your inbox too, without switching apps. See
  [detailed-configuration.md](detailed-configuration.md) for exactly how that
  works and its current limits.
- **You stay in control.** Every key is scoped to one vault, listed in one
  place, and revocable in one click.

## Who this is for

This is for anyone who already has an Artel vault and wants an AI assistant to
actually *use* it, rather than just being told about it in conversation. In
Artel's own language, this is the "AI Power User" journey: you generate a key,
drop it into Claude's settings, and from then on Claude can search, read, and
write your notes directly — turning your vault into the assistant's long-term
memory for your projects.

You don't need to know anything about servers, APIs, or how the connection
works under the hood to use this feature — the steps below are copy, paste,
done. If you're comfortable with a little more detail (or you're setting this
up for a script rather than a chat app), the advanced section in
[detailed-configuration.md](detailed-configuration.md) covers raw API calls
too.

## Before you start

You need a vault already created in Artel. If you don't have one yet, create
it first from the Artel web app — vault creation itself is covered elsewhere,
not in this guide. Once you have a vault, come back here and continue with
[quickstart.md](quickstart.md).
