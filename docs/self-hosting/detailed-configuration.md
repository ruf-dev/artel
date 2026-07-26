# Detailed Configuration

Everything you might need after your first startup: the available settings, where your data
lives, how to back it up or wipe it, and what to check if something isn't working.

## Configuration Options

If you're using the `docker-compose.omnibus.yaml` file, these environment variables can be set
on the `artel` service:

| Variable | Default | Notes |
|---|---|---|
| `COUCHDB_USER` | `admin` | The sync database's admin username. Only applied on first boot — set it before your very first `up`/`run` if you want something other than `admin`. Changing it afterward has no effect on an already-initialized instance. |
| `COUCHDB_PASSWORD` | `admin` | Same as above, for the password. |
| `SERVERS_MASTER_PORT` | `1551` | The port the Artel app listens on inside the container. If you change this, also update the port mapping (the `-p` flag or the compose file's `ports:` section) to match. |

If you're using plain `docker run`, set the same variables with `-e NAME=value` flags before the
image name, as shown in [quickstart.md](quickstart.md).

Everything else — internal database credentials, storage secrets, and the key used to encrypt
saved credentials — is generated automatically the first time the container starts, and stored
inside the container's own data volume. There is nothing else you need to configure to get a
working instance; the defaults above are the only knobs most people will ever touch.

## Data Persistence

All of your data — notes, sync history, attachments, accounts — lives under `/data` inside the
`artel-data` named volume. Concretely, that volume contains:

- Your sync database's data and indexes
- Your file/attachment storage
- Generated secrets used internally by the instance

As long as that volume exists, your data survives container restarts and even container
recreation. Stopping the instance with `docker compose down` (or `docker stop`/`docker rm` for a
plain `docker run` container) does **not** touch the volume — starting it back up again picks up
exactly where you left off.

### Backing Up

Since everything lives in the single `artel-data` Docker volume, backing up your instance means
backing up that volume. The simplest approach is Docker's own volume backup pattern — for
example:

```bash
docker run --rm -v artel-data:/data -v "$(pwd)":/backup alpine tar czf /backup/artel-backup.tar.gz -C /data .
```

This copies the entire contents of the volume into a single archive file in your current
directory, which you can store wherever you keep other backups. To restore, create a fresh
`artel-data` volume and extract the archive back into it before starting the container.

### Wiping Everything and Starting Fresh

If you want to reset your instance completely — deleting all accounts, vaults, and notes stored
in it — remove the volume along with the container:

```bash
docker compose -f docker-compose.omnibus.yaml down -v
```

(The `-v` is what actually deletes the volume — a plain `down` leaves your data intact.) The next
`up` starts completely from scratch, as if it were the very first time.

## Troubleshooting

If you started the container with plain `docker run` instead of `docker compose`, swap any
`docker compose -f docker-compose.omnibus.yaml logs/exec ...` command below for the plain-Docker
equivalent: `docker logs artel` / `docker exec artel ...`.

- **Something looks stuck or broken on startup.** Check the container's logs for output from each
  internal service:

  ```bash
  docker compose -f docker-compose.omnibus.yaml logs -f
  ```

- **The instance won't come up and you suspect a database problem.** Database setup failures are
  logged to a file inside the container:

  ```bash
  docker compose -f docker-compose.omnibus.yaml exec artel cat /data/postgres-init.log
  ```

- **One part seems broken but the rest of the app still loads.** The container runs each internal
  service independently, so a single one crashing or restarting in a loop doesn't take the whole
  instance down. Check the same `logs -f` output above to see which specific service is
  unhappy — the log lines are labeled by service, so you can tell which one is failing.

If none of the above turns up an obvious cause, a clean restart (`down` followed by `up -d`,
*without* `-v` so your data is preserved) is a reasonable first thing to try before digging
further.
