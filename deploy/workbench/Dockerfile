# Minimal image for a single workbench container: two independent ways to drive `claude` (Claude
# Code CLI) side by side. A headless Go relay ("the bridge", see ./bridge/README.md) drives it one
# turn at a time and speaks a normalized JSON event protocol over a WebSocket. Restored alongside
# it: a long-running `claude` session inside tmux, exposed as a real interactive terminal over
# HTTP/WebSocket by ttyd — one tmux window per terminal tab, each its own live `claude` TUI, with
# zero shared history with the chat bridge. No multi-session management beyond tmux's own windows,
# no supervisor process beyond entrypoint.sh itself.
#
# Stage 1: build the bridge binary. Separate Go module (bridge/go.mod) — see bridge/README.md for
# why — so it gets its own builder stage rather than being folded into the final image's own build.
# golang:1.26.3 matches the repo-root Dockerfile's builder stage.
FROM golang:1.26.3 AS bridge-builder

WORKDIR /bridge

COPY bridge/go.mod bridge/go.sum ./
RUN go mod download

COPY bridge/ .
RUN CGO_ENABLED=0 go build -o /workbench-bridge .

# Stage 2: the runtime image.
#
# node:20-slim (not debian:bookworm-slim) because @anthropic-ai/claude-code is an npm package
# and this avoids a separate Node.js install step.
FROM node:20-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends tmux ca-certificates curl \
    && rm -rf /var/lib/apt/lists/*

# ttyd exposes the tmux session above as a real interactive terminal over HTTP/WebSocket (see
# entrypoint.sh), which the backend reverse-proxies to on its own port, independent of the chat
# bridge's. A pinned upstream static release binary, so no C build toolchain has to land in this
# image; only the two architectures artel is actually built for are mapped, and anything else
# fails the build loudly rather than silently shipping an image with no terminal in it.
ARG TTYD_VERSION=1.7.7
RUN set -eu; \
    case "$(uname -m)" in \
      x86_64) ttyd_asset="ttyd.x86_64" ;; \
      aarch64 | arm64) ttyd_asset="ttyd.aarch64" ;; \
      *) echo "unsupported architecture for ttyd: $(uname -m)" >&2; exit 1 ;; \
    esac; \
    curl -fsSL -o /usr/local/bin/ttyd \
      "https://github.com/tsl0922/ttyd/releases/download/${TTYD_VERSION}/${ttyd_asset}"; \
    chmod +x /usr/local/bin/ttyd

# DISABLE_AUTOUPDATER=1 stops `claude` from self-updating in place at runtime. Without it, the
# CLI checks for a newer npm-published version on first real invocation and, if found, rewrites
# its own install (including the /usr/local/bin/claude symlink target) in the background — while
# that's in flight, ANY other process trying to exec `claude` races it and can hit ENOENT ("claude:
# executable file not found in $PATH"), since the two processes sharing this image (the bridge's
# own invocations and the tmux-interactive session) both exec the same bare `claude` concurrently.
# This image is rebuilt regularly via `make workbench-image`, which already reinstalls
# @anthropic-ai/claude-code fresh each time — an ephemeral per-session container self-updating at
# runtime only adds this race (and an unpinned, unaudited version bump) for no benefit.
ENV DISABLE_AUTOUPDATER=1

RUN npm install -g @anthropic-ai/claude-code

# Claude Code's first-run onboarding wizard (theme picker: "Choose the text style that looks
# best with your terminal") is gated purely on the global config's hasCompletedOnboarding field
# (~/.claude.json, resolved as $CLAUDE_CONFIG_DIR/.claude.json when that env var is set, else
# $HOME/.claude.json — this image sets neither, so it's /root/.claude.json), checked once at
# `claude` startup regardless of subcommand. Every workbench container starts with a brand-new
# $HOME, so without this the wizard is the very first thing either `claude` process here produces
# — both entrypoint.sh's tmux REPL and workbench-bridge's `claude setup-token` run share this same
# $HOME. For the bridge specifically this is fatal, not just cosmetic: its authlogin parser
# (bridge/internal/authlogin/parse.go) only recognizes setup-token's own auth-link/code-prompt/
# error/token text, so an onboarding screen in front of that is invisible to it — the
# subscription-login flow hangs forever with no output and no error, since nothing ever matches.
# Pre-seeding hasCompletedOnboarding (and theme, so it isn't reprompted separately) skips the
# wizard entirely, the same way the UTF-8 locale fix below works around another `claude`-first-run
# assumption.
#
# Same file also pre-accepts the per-project "Do you trust the files in this folder?" dialog —
# a second, independent first-run gate keyed by cwd (projects["<path>"].hasTrustDialogAccepted in
# this same config), which the CLI's own error text documents as settable exactly this way ("...or
# set projects[<path>].hasTrustDialogAccepted: true in ~/.claude.json"). Every `claude` invocation
# here runs at /root/vault (this image's WORKDIR, and homeMountPath in
# internal/clients/workbenchdocker/client.go — the volume every workbench container mounts its
# one working directory at), so that's the only project key that ever needs accepting.
RUN mkdir -p /root && printf '%s' '{"hasCompletedOnboarding":true,"theme":"dark","projects":{"/root/vault":{"hasTrustDialogAccepted":true}}}' > /root/.claude.json

# Without a UTF-8 locale, Claude Code's own terminal-capability detection assumes the
# terminal can't render Unicode and falls back to ASCII-safe glyphs (e.g. "_"/"*"/"."
# instead of box-drawing lines, "❯", "✔"). node:20-slim has no locale configured by
# default (LANG is unset, glibc falls back to POSIX/C) even though the C.utf8 locale
# ships built in and needs no locale-gen step. Applies to every `claude` invocation, whether
# started by the bridge or inside the tmux session, and to any later `docker exec` calls, which
# inherit the image's ENV.
ENV LANG=C.utf8
ENV LC_ALL=C.utf8

WORKDIR /root/vault

COPY --from=bridge-builder /workbench-bridge /usr/local/bin/workbench-bridge
COPY entrypoint.sh /entrypoint.sh
COPY tmux.conf /etc/tmux.conf
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
