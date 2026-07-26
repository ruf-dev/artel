# Self-Hosting Artel

## What Is Artel, Briefly

Artel gives you a private, cloud-based notebook — a **vault** — built on
[Obsidian](https://obsidian.md/), the popular note-taking app. Normally, keeping the same
Obsidian notes in sync across your phone, laptop, and tablet means running and maintaining your
own server, which is a technical chore most people never attempt. Artel sets that server up for
you, instantly, and (in its normal hosted form) it also connects your vault to AI assistants and
other tools like email and GitLab.

**This document is about a different way to get that same server: running it yourself, on your
own machine or your own hardware, instead of using Artel's hosted service.**

## Hosted vs. Self-Hosted

Most Artel users don't need this page at all — signing up on Artel's hosted service takes
seconds and someone else handles the server. Self-hosting is for people who would rather run
their own private copy of the whole stack:

- **Full data ownership.** Your notes, sync data, and attachments live entirely on hardware you
  control — nothing is stored on Artel's servers.
- **No dependency on a third-party service staying online.** Your vault works as long as your own
  machine or server does.
- **Fine for personal or small-group use.** Running it for yourself, your household, or a small
  team is straightforward. It is not intended as a production, multi-node deployment — see the
  note below.

If neither of those points matter to you, the hosted service is simpler and requires nothing
beyond a web browser. Keep reading only if you specifically want to run your own instance.

## What You'll Need

Self-hosting Artel is packaged as a single **Docker "omnibus" image** — one container that
bundles everything the service needs to run (its database, its file storage, and the Artel app
itself). You don't need to understand or configure any of those pieces individually; the image
sets them up on its own the first time it starts.

To follow along, you should have:

- **Docker installed** on the machine you'll run it on (Docker Desktop on Mac/Windows, or Docker
  Engine on Linux). If `docker --version` works in your terminal, you're set.
- **Basic terminal comfort.** You'll be copy-pasting a handful of commands into a terminal window.
  You don't need to understand what each one does internally, but you should be comfortable
  running commands and reading their output.
- A few minutes for the first startup, since the container sets up its internal storage the very
  first time it runs.

There's no need to install or configure a database, or anything else separately — the single
container is the entire installation.

*(There is also a more advanced, multi-container setup for larger, production-style deployments,
which splits the pieces the omnibus image bundles together into separate containers. This
document doesn't cover that — it's a different, more involved path meant for people already
comfortable operating multi-service infrastructure. Everything below is about the single-container
omnibus image, which is what almost every self-hoster wants.)*

## Next Step

Ready to get it running? Head to [quickstart.md](quickstart.md) for the copy-paste steps to have
your own private Artel instance up in a few minutes.
