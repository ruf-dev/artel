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

RUN npm install -g @anthropic-ai/claude-code

# Without a UTF-8 locale, Claude Code's own terminal-capability detection assumes the
# terminal can't render Unicode and falls back to ASCII-safe glyphs (e.g. "_"/"*"/"."
# instead of box-drawing lines, "❯", "✔"). node:20-slim has no locale configured by
# default (LANG is unset, glibc falls back to POSIX/C) even though the C.utf8 locale
# ships built in and needs no locale-gen step. Applies to every `claude` invocation, whether
# started by the bridge or inside the tmux session, and to any later `docker exec` calls, which
# inherit the image's ENV.
ENV LANG=C.utf8
ENV LC_ALL=C.utf8

WORKDIR /workspace

COPY --from=bridge-builder /workbench-bridge /usr/local/bin/workbench-bridge
COPY entrypoint.sh /entrypoint.sh
COPY tmux.conf /etc/tmux.conf
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
