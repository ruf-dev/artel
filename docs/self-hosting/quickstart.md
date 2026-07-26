# Quickstart: Running Your Own Artel Instance

This is the shortest path from "I have Docker installed" to "I'm looking at my own private Artel
instance in a browser." If you haven't already, skim [general-info.md](general-info.md) first to
make sure self-hosting is what you want.

## Option 1: The Fastest Way — `docker run`

If you just want to try it out, this single command downloads the pre-built image and starts
everything:

```bash
docker run -d --name artel -p 1551:1551 -p 5984:5984 -p 3900:3900 -v artel-data:/data redsockruf/artel:latest-omnibus
```

A quick note on what that command does: `-d` runs it in the background, `--name artel` gives the
container a friendly name, the three `-p` flags expose the ports you'll use (explained below),
and `-v artel-data:/data` creates a persistent storage volume so your data survives restarts.

Once it's running, open **http://localhost:1551** in your browser.

By default, the internal sync database is created with the username/password `admin`/`admin`.
If you want something else, add these flags *before the image name* and *before the first time
you run it* (they only take effect on first boot):

```bash
docker run -d --name artel -p 1551:1551 -p 5984:5984 -p 3900:3900 \
  -e COUCHDB_USER=yourname -e COUCHDB_PASSWORD=yourpassword \
  -v artel-data:/data redsockruf/artel:latest-omnibus
```

This is the quickest way to try things out. For regular, ongoing use, the `docker compose`
approach below is recommended instead, since it keeps your configuration in a file you can
re-run and version instead of retyping a long command.

## Option 2: The Recommended Way — `docker compose`

This builds the image from source on your own machine and starts it via a compose file already
included in the Artel repository.

```bash
docker compose -f docker-compose.omnibus.yaml up -d
```

The first startup takes a bit longer than later ones — behind the scenes, the container is
setting up its internal database and storage for the first time. Subsequent restarts are much
faster.

Once it's up, you have three things running, all inside the one container:

- **Web UI + API** — http://localhost:1551 — this is the page you'll actually use.
- **Sync endpoint** — http://localhost:5984 — this is what the Obsidian sync plugin connects to
  directly; you won't need to open this in a browser yourself.
- **S3-compatible storage API** — http://localhost:3900 — used for attachments, and available if
  you want to point an external S3 client at your data directly.

## First Login

1. Open **http://localhost:1551**.
2. Register a new account and sign in, exactly as you would on the hosted service.
3. Create your first vault from the dashboard.

From here on, using your self-hosted instance works identically to Artel's hosted service — see
the [Vault documentation](../vault/quickstart.md) for creating and syncing your first notebook.
That part of Artel doesn't know or care whether it's self-hosted.

## Stopping and Restarting

If you used `docker compose`, stop the container without losing any data with:

```bash
docker compose -f docker-compose.omnibus.yaml down
```

Your data stays intact in the `artel-data` volume — running `up -d` again picks up right where
you left off. See [detailed-configuration.md](detailed-configuration.md) for how data persistence
works, how to fully wipe an instance, and how to troubleshoot startup problems.
