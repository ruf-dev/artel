# Workbench: your always-on coding assistant

> **Early access notice:** Workbench exists and works today, but only "under the hood."
> There is no button or page for it in Artel yet. See
> [quickstart.md](quickstart.md) for exactly what that means and what's coming.

## What is Workbench?

Workbench is a private, always-on coding assistant environment that Artel runs for you in the
cloud. Think of it as your own personal copy of Claude Code — Anthropic's AI coding agent —
except instead of living on your laptop (and only working while your laptop is on and open),
it lives on Artel's servers and is available 24 hours a day, 7 days a week.

Normally, running an AI coding agent like Claude Code means it only works while your own
computer is turned on, awake, and connected to the internet. Close your laptop, and the agent
stops. Workbench removes that limitation: Artel provisions and hosts the environment for you,
so your assistant keeps running whether or not your own device is on.

## Why it matters

- **No always-on machine required.** You don't need to leave a laptop running, rent your own
  server, or babysit a process. Artel keeps it available for you.
- **One per vault.** When you have a vault (your Obsidian cloud notebook) set up with Artel,
  a matching Workbench environment is created alongside it — a dedicated space tied to that
  vault, not shared with other users.
- **Persistent.** A Workbench can be paused and resumed. Pausing it doesn't erase your work —
  it's more like putting a computer to sleep than shutting it down and wiping it, so your
  assistant's state picks back up where it left off.
- **Flexible billing.** You'll be able to choose how the assistant is powered: with your own
  Claude API key (usage-based billing), or by logging into your own Claude subscription
  (Pro/Max), so usage draws from the plan you already pay for instead of separate metered
  billing. See [detailed-configuration.md](detailed-configuration.md) for how these two options
  differ.

## Who it's for

Anyone who wants an AI coding/agent assistant that's available continuously — for example, to
keep working on something overnight, to hand off a longer task and check back later, or simply
to avoid needing to keep a personal machine running just to have access to Claude Code. If
you're an Artel user who already has (or plans to have) a vault, Workbench is designed to sit
right alongside it.

## What Workbench is *not* (yet)

To set expectations clearly:

- It is **not** something you can click into and use today — there's no on/off switch or chat
  screen in the Artel app yet. See [quickstart.md](quickstart.md).
- It does **not** yet give you a way to actually type into or view your running assistant from
  a browser or phone — that "open a terminal to it" experience is a future piece of work, not
  part of what exists today.
- It does **not** yet keep its files in sync with your vault's notes. The assistant's own
  workspace and your vault's notes are separate things for now; connecting them is planned as
  a later step, not something you should rely on yet.

None of this is a bug or an oversight — Workbench is being built in stages, and the stage that's
done so far intentionally covers the "engine," not the "dashboard." See
[quickstart.md](quickstart.md) for the honest current status and
[detailed-configuration.md](detailed-configuration.md) for more on how the pieces will fit
together once the rest ships.
