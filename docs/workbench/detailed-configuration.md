# Workbench: how it works (concepts, not a how-to)

Like the rest of this documentation set, this page describes how Workbench is designed to
work. As explained in [quickstart.md](quickstart.md), there's no app screen for any of this
yet — this page exists so you understand the concepts and options in advance.

## The two ways to power your Workbench

When you start a Workbench, it needs a way to authenticate as you with Claude. There are two
options, and you pick one each time you start it.

### Option 1: Your own Claude API key

This reuses the same API key connection Artel already supports elsewhere (under
Integrations — bring-your-own-key, or "BYOK"). If you've already connected an Anthropic API
key to your Artel account, Workbench can use that same key.

- **How it works, in plain terms:** Artel securely hands your key to your private environment
  when it starts, and your assistant authenticates with it immediately — no extra steps, no
  login screen.
- **Billing:** Usage is billed per the API key's own metered pricing (the standard "pay for
  what you use" API pricing), same as any other API-key usage.
- **Best for:** People who already have an API key connected, or who prefer usage-based
  billing over a subscription.
- **Note:** Setting up the API key connection itself is covered under Artel's Integrations
  documentation — this page is only about Workbench using a key you've already connected.

### Option 2: Log in with your Claude subscription

Instead of a key, you can sign in with your existing Claude.ai account — the same login you'd
use on claude.ai, including a Pro or Max subscription.

- **How it works, in plain terms:** When you choose this option, Artel starts your environment
  with no key attached. Your assistant then generates a normal Anthropic sign-in link — the
  same kind of "sign in with your account" link you'd see anywhere else. You open that link,
  which takes you to Anthropic's own official sign-in page (not an Artel page) in your browser,
  sign in there as usual, and Anthropic shows you a short one-time code. You copy that code and
  paste it back into Artel to confirm the sign-in. Artel never sees or stores your password —
  it only relays the short code you choose to paste back, much like copying a two-factor code
  from an authenticator app into a website.
- **Billing:** Usage draws from your existing Claude subscription's included usage, instead of
  separate metered API charges. This is the main reason this option exists — for people who
  already pay for Pro/Max and would rather use that included usage for their always-on
  assistant instead of paying for API usage on top of it.
- **Best for:** People who already have a Claude Pro or Max subscription and want their
  Workbench usage to come out of that plan.

Either option can be used independently — you don't need both, and you choose per session
which one you're using.

## The lifecycle: what state is your Workbench in?

A Workbench environment moves through a small number of plain-language states over its life.
You don't need to manage these directly (once there's a UI, this will just be reflected as
a status), but it helps to know what they mean:

| State | What it means in plain terms |
|---|---|
| **Created** | Your private environment has been reserved and set up, but hasn't been turned on yet. Nothing is running, and nothing is being billed. |
| **Running** | Your assistant is live and active — this is the "on" state. |
| **Stopped** | Paused. Your environment isn't running right now, but nothing has been thrown away — all your assistant's state and files are kept, ready to resume. Think of this like a laptop that's asleep rather than one that's been wiped. |
| **Removed** | Fully torn down — this happens if you delete the associated vault, or if the Workbench is explicitly removed. Unlike "stopped," this does not preserve anything; you'd be starting fresh next time. |

A Workbench typically moves: **Created → Running → Stopped → Running** (you can pause and
resume repeatedly), and only reaches **Removed** when you (or the deletion of your vault)
explicitly tear it down.

Each vault gets exactly one Workbench environment, created automatically alongside the vault
itself — you won't need to separately request or provision one.

## What's explicitly not built yet

To avoid any confusion about the scope of what exists, three things are worth calling out
specifically as **not yet available**, even though the rest of the lifecycle above is working:

1. **A real terminal you can use.** Right now there isn't a way to actually open a screen in
   your browser or phone and type to your running assistant interactively. Getting the
   environment running and logged in is done; giving you a live window into it to actually work
   with it is a separate, dedicated piece of work planned for later.

2. **Syncing with your vault's notes.** Your Workbench's own files currently live separately
   from your vault's notes — they are not automatically kept in sync. If you're picturing your
   assistant reading and writing directly into your Obsidian notes, that connection doesn't
   exist yet; it's planned as a future addition once the rest of Workbench is proven out.

3. **Any user-facing screen at all.** As covered in [quickstart.md](quickstart.md), starting,
   stopping, and logging in currently all require direct backend access — there's no app screen
   for any part of this yet.

None of these are needed for Workbench's core idea (an always-on assistant environment) to
work — they're additional pieces being built on top of a foundation that's already solid.
