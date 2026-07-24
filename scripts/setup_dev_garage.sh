#!/usr/bin/env bash
# Bootstraps the local dev Garage container (docker-compose.yaml: artel-garage-s3) after
# `docker compose up -d`: a fresh Garage node has no layout and can't store anything until a
# capacity/zone is assigned and applied, and it has no S3 access key until one is created.
# Safe to re-run — every step here is a no-op if already done.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

CONTAINER="${GARAGE_CONTAINER:-artel-garage-s3}"
ZONE="${GARAGE_ZONE:-z1}"
CAPACITY="${GARAGE_CAPACITY:-1G}"
KEY_NAME="${GARAGE_KEY_NAME:-artel-dev}"
S3_PORT="${GARAGE_S3_PORT:-13900}"
GARAGE_TOML="$REPO_ROOT/.verv/template/garage.toml"

garage() {
  docker exec "$CONTAINER" /garage "$@"
}

echo "--- Waiting for $CONTAINER to be up ---"
if ! docker inspect -f '{{.State.Running}}' "$CONTAINER" >/dev/null 2>&1; then
  echo "$CONTAINER is not running, starting it via docker compose ---"
  (cd "$REPO_ROOT" && docker compose up -d artel-garage-s3)
fi

echo "--- Waiting for Garage admin API ---"
i=0
until garage status >/dev/null 2>&1; do
  i=$((i + 1))
  if [ "$i" -ge 60 ]; then
    echo "Garage did not become ready in time" >&2
    exit 1
  fi
  sleep 1
done

echo "--- Checking cluster layout ---"
layout_version=$(garage layout show 2>/dev/null | awk -F': ' '/Current cluster layout version/{print $2}')

if [ "$layout_version" = "0" ]; then
  node_id=$(garage status 2>/dev/null | awk '/^[0-9a-f]{16}/{print $1; exit}')
  if [ -z "$node_id" ]; then
    echo "Could not determine Garage node id from 'garage status'" >&2
    exit 1
  fi

  echo "Assigning layout to node $node_id (zone=$ZONE, capacity=$CAPACITY)"
  garage layout assign "$node_id" -z "$ZONE" -c "$CAPACITY"
  garage layout apply --version 1
else
  echo "Layout already applied (version $layout_version), skipping"
fi

echo "--- Ensuring API key '$KEY_NAME' exists ---"
key_id=$(garage key list 2>/dev/null | awk -v name="$KEY_NAME" '$3==name{print $1}')

if [ -z "$key_id" ]; then
  create_out=$(garage key create "$KEY_NAME" 2>/dev/null)
  key_id=$(echo "$create_out" | awk -F': +' '/^Key ID/{print $2}')
else
  echo "Key '$KEY_NAME' already exists ($key_id)"
fi

garage key allow --create-bucket "$key_id" >/dev/null

key_info=$(garage key info "$key_id" --show-secret 2>/dev/null)
secret_key=$(echo "$key_info" | awk -F': +' '/^Secret key/{print $2}')

region="garage"
if [ -f "$GARAGE_TOML" ]; then
  region=$(awk -F'"' '/^s3_region/{print $2; exit}' "$GARAGE_TOML")
fi

echo
echo "--- Garage dev setup complete ---"
echo "Endpoint:    localhost:$S3_PORT"
echo "Region:      $region"
echo "Path style:  true"
echo "Access key:  $key_id"
echo "Secret key:  $secret_key"
echo
echo "Register these as an S3 instance in Artel to use Garage for vault storage."
