#!/bin/sh
# Starts (or re-attaches to) a single named tmux session running `claude`, then serves that
# session as a real interactive terminal over HTTP/WebSocket via ttyd, and keeps this process
# (container PID 1) alive so the container itself keeps running independently of whether
# anything is currently attached.
#
# Session name must match tmuxSessionName in
# internal/clients/workbenchdocker/container.go — that's how a later `docker exec` (used to
# inject env like ANTHROPIC_API_KEY without ever writing it to this container's image or
# volume) finds this session.
set -e

SESSION_NAME="workbench"
# Must match ttydPort in internal/clients/workbenchdocker/container.go, which publishes it to a
# Docker-assigned random host port — the backend's terminal reverse-proxy dials that published
# port on the configured Docker daemon's own address, not this container's workbench-net IP
# directly (see ContainerAddress's doc comment for why).
TTYD_PORT=7681

if ! tmux has-session -t "$SESSION_NAME" 2>/dev/null; then
  tmux new-session -d -s "$SESSION_NAME" "claude"
fi

# Exit cleanly on `docker stop` instead of waiting out the SIGKILL timeout — the loop below is
# PID 1, and a shell with no trap installed ignores signals with default dispositions as PID 1.
trap 'exit 0' TERM INT

# PID 1: supervise ttyd rather than exec'ing it, so a ttyd crash or restart never takes the
# container down with it — the tmux server (and `claude`, plus the session environment injected
# into it at start time) lives on as a separate background process, and the next ttyd just
# re-attaches to the same session. -W makes the terminal writable rather than read-only.
while true; do
  # Explicit fontFamily so the terminal matches macOS Terminal.app/iTerm2's default
  # monospace stack rather than relying on xterm.js's own built-in default. (The actual
  # cause of Claude Code's box-drawing/status glyphs rendering as "_"/"*" was a missing
  # UTF-8 locale, fixed via LANG/LC_ALL in the Dockerfile — see the comment there.)
  ttyd -W -p "$TTYD_PORT" \
    -t 'fontFamily="Menlo, Monaco, Courier New, monospace"' \
    tmux attach -t "$SESSION_NAME" || true
  sleep 1
done
