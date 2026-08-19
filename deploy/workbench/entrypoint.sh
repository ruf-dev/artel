#!/bin/sh
# Runs two independent things that both drive `claude` inside this container, as this
# container's PID 1: the headless chat bridge (see ./bridge/README.md), which drives `claude` one
# turn at a time over its own WebSocket port, and a tmux session exposed as a real interactive
# terminal by ttyd on a separate port — restored alongside the bridge so a browser can still get a
# raw, per-tab `claude` TUI (one tmux window per tab), with zero shared history between the two.
#
# Session name must match tmuxSessionName in internal/clients/workbenchdocker/container.go — the
# terminal-tab RPCs (ListTerminalTabs et al., see tmux_tabs.go) target it directly by name.
set -e

SESSION_NAME="workbench"
# Must match ttydPort in internal/clients/workbenchdocker/container.go — 7682, distinct from the
# bridge's own port (bridgePort there, 7681, which workbench-bridge binds to internally and needs
# no mention here), so both servers can be published to the host side by side.
TTYD_PORT=7682

if ! tmux has-session -t "$SESSION_NAME" 2>/dev/null; then
  tmux new-session -d -s "$SESSION_NAME" "claude"
fi

# Neither the bridge nor ttyd is exec'd straight into PID 1: this shell has to stay alive as PID 1
# so a `docker stop` (which signals PID 1 only) can be forwarded to both children explicitly via
# the trap below, and so either one exiting on its own doesn't silently leave the other running
# with no supervisor at all.
/usr/local/bin/workbench-bridge &
BRIDGE_PID=$!

# Supervise ttyd in a loop rather than running it once: ttyd is a crash-prone C binary being
# supervised from the outside because a crash in it must not take `claude`'s tmux session down
# with it — the tmux server (and `claude`, plus whatever env is injected into its session at start
# time) lives on as a separate process, and the next loop iteration just re-attaches to the same
# session. -W makes the terminal writable rather than read-only.
(
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
) &
TTYD_LOOP_PID=$!

# Exit cleanly on `docker stop` instead of waiting out the SIGKILL timeout: forward the signal to
# both children so the bridge gets its own SIGTERM handling (main.go) a chance to run, and the
# ttyd supervision loop above actually stops instead of respawning ttyd forever. The tmux
# server/`claude` process are deliberately left running here, same as the pre-bridge version of
# this script — they're torn down by the container exiting right after, not by this trap, since a
# workbench's tmux session is meant to outlive any one ttyd attach/detach cycle, not this script's
# own lifetime.
trap 'kill -TERM "$BRIDGE_PID" "$TTYD_LOOP_PID" 2>/dev/null; wait "$BRIDGE_PID" 2>/dev/null; exit 0' TERM INT

# PID 1 waits on the bridge specifically, not a bare `wait` (which would also block forever on the
# ttyd loop's own `while true`) — the bridge exiting, planned or crashed, is what ends this
# container's lifecycle, same as it being exec'd directly used to mean before ttyd/tmux came back.
wait "$BRIDGE_PID"
