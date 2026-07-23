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
    && apt-get install -y --no-install-recommends tmux \
    && rm -rf /var/lib/apt/lists/*

RUN npm install -g @anthropic-ai/claude-code

WORKDIR /workspace

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
