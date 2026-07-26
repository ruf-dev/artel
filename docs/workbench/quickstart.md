# Workbench: current status (there's no quickstart yet)

This page is deliberately not a normal "get started in 3 steps" guide, because those steps
don't exist in the Artel app yet. Here's the honest picture: what's actually built, and what
using Workbench will look and feel like once the rest of it ships.

## What's actually available right now

Workbench has been built and works — but only at the technical/API level, behind the scenes.
Specifically, as of today:

- Artel can create a private environment for your vault, start it, and stop it.
- Artel can authenticate that environment either with an API key or by walking through a
  Claude subscription login (see below).
- **None of this has a page, button, or screen in the Artel app yet.** There is nothing to
  click. If you're not comfortable calling APIs directly, there is currently no way for you to
  use Workbench yourself.

In short: the engine has been built and tested, but the dashboard hasn't been built. This is a
normal, intentional stage of building a feature — it's documented here so you know what to
expect (and don't go looking for a "Workbench" tab that isn't there).

## What using it will look like once the UI ships

This section describes the intended experience once the missing UI pieces are built. Nothing
below is available today — treat it as a preview, not instructions.

1. **Turning it on.** From your vault, you'll be able to start your Workbench environment with
   a single action (something as simple as a button). This spins up your private assistant
   environment if it isn't already running.

2. **Choosing how it's powered.** The first time you start it, you'll choose one of two ways to
   authenticate:
   - **Use your own Claude API key** — the same kind of API key you'd connect under Artel's
     Integrations settings. If you already have one connected, this option should just work.
   - **Log in with your Claude subscription** — instead of a key, you sign in with your
     existing Claude.ai account (Pro or Max). Artel will show you a link to open, you sign in
     through Anthropic's own sign-in page in your browser, and then copy a short code back into
     Artel to confirm it's you. Nothing about your account credentials passes through or is
     stored by Artel — Artel only relays what you paste, the same way you might copy a
     verification code from one app to another.

   See [detailed-configuration.md](detailed-configuration.md) for a closer look at how these
   two options differ and when you'd pick one over the other.

3. **Using it.** Once running, your assistant stays on continuously — you'd be able to give it
   work, step away, and come back later to find it's kept going. Exactly how you'll *reach* the
   assistant to type to it (a chat box, a terminal-style screen, etc.) is one of the pieces
   still being designed — see the note below.

4. **Pausing and resuming.** You'll be able to stop it when you don't need it and start it again
   later without losing its state — similar to putting a computer to sleep rather than
   resetting it.

## What's explicitly missing (don't expect these yet)

- **No app screen for any of the above.** Starting, stopping, and logging in all currently
  require calling Artel's backend directly — there is no button in the product.
- **No way to actually reach/type into the running assistant.** Even for someone using the
  backend directly today, there isn't yet a real "terminal in your browser" experience to
  interact with the assistant session live. That's planned as a separate, dedicated piece of
  future work.
- **No syncing with your vault's notes.** The assistant's own workspace is currently separate
  from your vault's notes — they don't automatically share files. That connection is planned
  for later and isn't built yet.

## Bottom line

If you're curious about Workbench, the best next step today is understanding what it's for
(see [general-info.md](general-info.md)) and how its authentication options will work
(see [detailed-configuration.md](detailed-configuration.md)) — not looking for it in the app,
since it isn't there yet. This page will be rewritten into a real quickstart once the
corresponding UI ships.
