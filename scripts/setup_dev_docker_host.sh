#!/usr/bin/env bash
# Registers the local dind daemon (tests/docker-compose.yaml: test-dockerd) as a docker_hosts
# row so workbench.Service.CreateWorkbench has somewhere to schedule containers in local dev —
# see docs/workbench/02_docker_topology.md. Safe to re-run: RegisterDockerHost is called only if
# no docker_hosts row for DOCKER_HOST_URL already exists (docker_hosts.url is UNIQUE per
# migrations/061_docker_hosts.sql, so a duplicate call is harmless either way, but this avoids
# noise from a failed duplicate-key POST on every re-run).
#
# Assumes both `docker compose -f tests/docker-compose.yaml up -d test-dockerd` and the Artel app
# itself are already up (this script only talks to test-dockerd's Docker API and the Artel HTTP
# gateway — it does not start either).
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

DOCKERD_CONTAINER="${DOCKERD_CONTAINER:-artel-test-dockerd}"
DOCKER_HOST_PORT="${DOCKER_HOST_PORT:-12375}"
DOCKER_HOST_URL="${DOCKER_HOST_URL:-tcp://localhost:$DOCKER_HOST_PORT}"
WORKBENCH_NETWORK="${WORKBENCH_NETWORK:-workbench-net}"
ARTEL_URL="${ARTEL_URL:-http://localhost:1551}"
ARTEL_AUTH_TOKEN="${ARTEL_AUTH_TOKEN:-}"

echo "--- Waiting for $DOCKERD_CONTAINER (test-dockerd) to be up ---"
if ! docker inspect -f '{{.State.Running}}' "$DOCKERD_CONTAINER" >/dev/null 2>&1; then
  echo "$DOCKERD_CONTAINER is not running, starting it via docker compose ---"
  (cd "$REPO_ROOT" && docker compose -f tests/docker-compose.yaml up -d test-dockerd)
fi

echo "--- Waiting for test-dockerd's Docker API on $DOCKER_HOST_URL ---"
i=0
until docker -H "$DOCKER_HOST_URL" version >/dev/null 2>&1; do
  i=$((i + 1))
  if [ "$i" -ge 60 ]; then
    echo "test-dockerd's Docker API did not become ready in time" >&2
    exit 1
  fi
  sleep 1
done

# workbench-net is assumed by internal/clients/workbenchdocker/client.go to pre-exist on the
# configured daemon. A compose-level `networks:` block only creates it in the host running
# compose, not inside dind's own inner dockerd (a separate daemon/network namespace) — so it has
# to be created explicitly against test-dockerd's own Docker API.
echo "--- Ensuring '$WORKBENCH_NETWORK' network exists inside test-dockerd ---"
if ! docker -H "$DOCKER_HOST_URL" network inspect "$WORKBENCH_NETWORK" >/dev/null 2>&1; then
  docker -H "$DOCKER_HOST_URL" network create "$WORKBENCH_NETWORK" >/dev/null
else
  echo "Network '$WORKBENCH_NETWORK' already exists, skipping"
fi

echo "--- Waiting for Artel app at $ARTEL_URL ---"
i=0
until curl -sf -o /dev/null "$ARTEL_URL/api/docker_hosts/list" \
  -X POST -H 'Content-Type: application/json' \
  ${ARTEL_AUTH_TOKEN:+-H "Authorization: Bearer $ARTEL_AUTH_TOKEN"} \
  -d '{}'; do
  i=$((i + 1))
  if [ "$i" -ge 60 ]; then
    echo "Artel app did not become reachable at $ARTEL_URL in time" >&2
    echo "(this endpoint is admin-gated — if ENVIRONMENT_NO_AUTH_ENABLED=false, set ARTEL_AUTH_TOKEN too)" >&2
    exit 1
  fi
  sleep 1
done

echo "--- Checking for an existing docker_hosts row for $DOCKER_HOST_URL ---"
existing=$(curl -sf -X POST -H 'Content-Type: application/json' \
  ${ARTEL_AUTH_TOKEN:+-H "Authorization: Bearer $ARTEL_AUTH_TOKEN"} \
  -d '{}' "$ARTEL_URL/api/docker_hosts/list" \
  | grep -c "\"$DOCKER_HOST_URL\"" || true)

if [ "$existing" -gt 0 ]; then
  echo "A docker_hosts row for $DOCKER_HOST_URL already exists, skipping registration"
else
  echo "--- Registering $DOCKER_HOST_URL as a docker_hosts row ---"
  register_payload=$(printf '{"url":"%s"}' "$DOCKER_HOST_URL")

  curl -sf -X POST -H 'Content-Type: application/json' \
    ${ARTEL_AUTH_TOKEN:+-H "Authorization: Bearer $ARTEL_AUTH_TOKEN"} \
    -d "$register_payload" "$ARTEL_URL/api/docker_hosts/add" >/dev/null
fi

echo
echo "--- Local dind docker host setup complete ---"
echo "Docker host URL: $DOCKER_HOST_URL"
echo "Network:         $WORKBENCH_NETWORK (created inside test-dockerd)"
