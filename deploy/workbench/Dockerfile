# Minimal image for a single workbench container: a long-running `claude` (Claude Code CLI)
# session inside tmux, so a later `docker exec` can attach/detach without losing the session.
#
# See docs/workbench/02_docker_topology.md and docs/workbench/04_task_breakdown.md for the
# surrounding design. This is deliberately a plain single-session prototype — no multi-session
# management, no supervisor process beyond entrypoint.sh itself.
#
# node:20-slim (not debian:bookworm-slim) because @anthropic-ai/claude-code is an npm package
# and this avoids a separate Node.js install step.
FROM node:20-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends tmux ca-certificates curl \
    && rm -rf /var/lib/apt/lists/*

# ttyd exposes the tmux session above as a real interactive terminal over HTTP/WebSocket (see
# entrypoint.sh), which the backend reverse-proxies to. A pinned upstream static release binary,
# so no C build toolchain has to land in this image; only the two architectures artel is actually
# built for are mapped, and anything else fails the build loudly rather than silently shipping
# an image with no terminal in it.
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

WORKDIR /workspace

COPY entrypoint.sh /entrypoint.sh
COPY tmux.conf /etc/tmux.conf
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
