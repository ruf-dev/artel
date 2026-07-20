# Omnibus image

The omnibus image bundles Postgres, CouchDB, Garage (S3) and the artel app together in a single
container, supervised by [s6-overlay](https://github.com/just-containers/s6-overlay). It's meant
for personal/self-hosted use — "run one container and the whole stack is there" — not for
production multi-node deployments. For the standard multi-container setup (separate Postgres/
CouchDB/Garage containers), see `docker-compose.yaml` instead.

Built from [`Dockerfile.omnibus`](../Dockerfile.omnibus); CI publishes it with a `-omnibus` tag
suffix alongside the regular image (`.github/workflows/release.yaml`).

## Quick start

```bash
docker run -d --name artel -p 1551:1551 -p 5984:5984 -p 3900:3900 -v artel-data:/data redsockruf/artel:latest-omnibus
```

Then open http://localhost:1551. Add `-e COUCHDB_USER=... -e COUCHDB_PASSWORD=...` before first
run to override the default `admin`/`admin` CouchDB credentials. See below for the `docker compose`
equivalent (recommended for regular use), or building the image from source instead of pulling it.

## Building from source

```bash
docker compose -f docker-compose.omnibus.yaml up -d
```

This builds the image locally (`context: .`, `dockerfile: Dockerfile.omnibus`) and starts a
single `artel` container. First boot takes longer than subsequent restarts — Postgres, CouchDB
and Garage all run their init steps (create DB/user, write `local.ini`, generate an RPC secret and
assign a single-node storage layout) before the `artel` service itself starts.

Once it's up:

- Web UI + gRPC/HTTP API: http://localhost:1551
- CouchDB (Obsidian LiveSync connects here directly): http://localhost:5984
- Garage S3 API (attachments / external S3 clients): http://localhost:3900

Tear down (data is preserved in the `artel-data` volume):

```bash
docker compose -f docker-compose.omnibus.yaml down
```

Wipe everything and start fresh:

```bash
docker compose -f docker-compose.omnibus.yaml down -v
```

## Configuration

The compose file exposes a few environment variables on the `artel` service:

| Variable | Default | Notes |
|---|---|---|
| `COUCHDB_USER` | `admin` | Only applied on first boot — set before the first `up` if you want something other than `admin`/`admin`. |
| `COUCHDB_PASSWORD` | `admin` | Same as above. |
| `SERVERS_MASTER_PORT` | `1551` | Port the artel app listens on inside the container; also update the `ports:` mapping if you change it. |

Everything else (Postgres credentials, Garage RPC/admin secrets, the app's credential-encryption
key) is generated on first boot and stored under `/data/secrets` — there's nothing else to
configure to get a working instance.

## Data persistence

All state lives under `/data` in the `artel-data` named volume:

- `/data/postgres` — Postgres cluster (`artel_db`, owned by role `artel`)
- `/data/couchdb` — CouchDB data + view indexes + `local.ini`
- `/data/garage` — Garage metadata/data dirs + `garage.toml`
- `/data/secrets` — generated secrets (Garage RPC secret, Garage admin token, creds encryption key)

`docker compose down` followed by `up` again preserves all of it; only `down -v` (or deleting the
volume) resets the instance.

## Using it

1. Open http://localhost:1551 and register/sign in through the web UI to create your artel account
   and provision a vault.
2. Point [Obsidian LiveSync](https://github.com/vrtmrz/obsidian-livesync) at
   `http://localhost:5984` with the CouchDB credentials above (or whatever you overrode them to)
   to sync a vault.
3. If you need direct S3 access to attachments, use the Garage S3 API at `http://localhost:3900`
   (Garage is configured with `s3_region = "garage"` by default, single-node, replication factor 1).

## Troubleshooting

(swap `docker compose -f docker-compose.omnibus.yaml logs/exec ...` for plain `docker logs`/`docker exec artel ...` if you started the container with `docker run` instead of compose)

- Check container logs for per-service startup output: `docker compose -f docker-compose.omnibus.yaml logs -f`.
- Postgres init failures are logged to `/data/postgres-init.log` inside the container
  (`docker compose -f docker-compose.omnibus.yaml exec artel cat /data/postgres-init.log`).
- s6 supervises each service independently (`/etc/services.d/{postgres,couchdb,garage,garage-layout-init,artel}`),
  so an individual service crash-looping won't take down the others — check
  `docker compose -f docker-compose.omnibus.yaml logs -f` for which one is failing.
