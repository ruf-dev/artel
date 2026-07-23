#!/bin/sh
# Starts (or re-attaches to) a single named tmux session running `claude`, and keeps this
# process (container PID 1) alive so the container itself keeps running independently of
# whether anything is currently attached to the tmux session.
#
# Session name must match tmuxSessionName in
# internal/clients/workbenchdocker/container.go — that's how a later `docker exec` (used to
# inject env like ANTHROPIC_API_KEY without ever writing it to this container's image or
# volume) finds this session.
set -e

SESSION_NAME="workbench"

if ! tmux has-session -t "$SESSION_NAME" 2>/dev/null; then
  tmux new-session -d -s "$SESSION_NAME" "claude"
fi

# PID 1: block forever. The tmux server keeps running as a background process regardless of
# whether this process is doing anything, so this is only here to keep the container alive.
exec tail -f /dev/null
