# Workbench — Design Plan

Status: **planning only, nothing in this folder is implemented yet.**

## Goal

Give a user a standalone, long-running Claude Code instance that Artel provisions and hosts —
one Docker container per user, created alongside their vault — so they get 24/7 agent access
without owning a always-on machine themselves. This is the natural next step past
[BYOK](../byok/README.md): BYOK let a user's own Anthropic key be *used* by Artel; Workbench lets
Artel *run Claude Code itself* on the user's behalf, authenticated either via that same BYOK key
or via the user's own Claude subscription login.

The feature is entirely optional at the deployment level: it only exists when Artel is configured
with a Docker daemon to talk to. No `DOCKER_HOST`-equivalent config set → the workbench service
isn't wired in, no schema surface is exposed, nothing changes for self-hosters without Docker.

Explicitly **out of scope for this pass**: syncing the workbench's local notes with the vault's
CouchDB (LiveSync or otherwise). The workbench prototype only needs to create a container, start
it, and get a `claude` session logged in inside it — what that session's filesystem contains is a
separate design pass once the lifecycle/login mechanics are proven out.

## Documents in this folder

| File | Covers |
|---|---|
| [01_data_model_and_lifecycle.md](01_data_model_and_lifecycle.md) | `workbenches` table, status state machine, hook points into vault create/delete |
| [02_docker_topology.md](02_docker_topology.md) | How Artel talks to Docker, daemon placement options, network isolation, config gating |
| [03_auth_and_login_flow.md](03_auth_and_login_flow.md) | API-key injection vs. subscription login, the login-flow unknown that needs a spike first |
| [04_task_breakdown.md](04_task_breakdown.md) | Ordered, file-level implementation tasks referencing the above |

## Key architectural decisions (summary — see linked docs for the "why")

- **Naming**: the feature and its service/table are called "workbench" throughout (not
  "workspace"/"sandbox") — user's call, kept consistent everywhere.
- **One container + one volume per user**, created (not started) at the same time as
  `VaultService.CreateVault`, via a sibling call at the transport-handler level — not a dependency
  baked into `VaultService` itself.
- **Not Docker-in-Docker.** DinD (nested privileged daemon) solves a problem the user doesn't
  have yet (untrusted nested `docker build`/`docker run` inside a workbench). The actual ask —
  keeping workbench containers from being mixed in with Artel's own infra containers — is a
  grouping/labeling problem, solved with a dedicated daemon or network, not nesting.
- **Docker access is a config-gated optional dependency**, same shape as `SubscriptionsEnabled`
  gating `PaidService` vs `FreeService` — absence of config means absence of the whole subsystem,
  not a degraded fallback.
- **Two auth modes, both reusing existing plumbing**: BYOK API key (already-built
  `external_connections` Anthropic row, injected as an env var) or Claude subscription login
  (headless `claude` login flow, mechanics unconfirmed — flagged as the first thing to validate,
  not assumed).
- **Notes/file sync is deliberately undesigned here** — the workbench prototype's container has a
  volume and a shell; what populates that volume is the next design pass.
