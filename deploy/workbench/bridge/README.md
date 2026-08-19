# workbench-bridge

The process a workbench container runs as PID 1 (see `../entrypoint.sh`). It replaces the previous
tmux + ttyd + interactive `claude` stack with a headless relay: it drives the `claude` CLI in
`--print` mode, one process per chat turn, and speaks the normalized JSON event protocol from
`internal/chatprotocol` over a WebSocket on port 7681 — the same port ttyd used, so the backend's
reverse proxy (`internal/transport/vaults_api/workbench_terminal.go`) needs no changes.

This is a **separate Go module**: it is compiled in its own builder stage inside the workbench
image (`../Dockerfile`) and therefore cannot import `github.com/ruf-dev/artel`.
`internal/chatprotocol/events.go` is a hand-mirrored, byte-identical copy of the main module's
`internal/chatprotocol/events.go` — keep the two in sync on every change.

## Pieces

| file | role |
| --- | --- |
| `main.go` | wiring, signal handling, "bridge stays up" supervision |
| `internal/envdrop` | reads the `/run/workbench/env` tmpfs the backend's `injectEnv` drops secrets into, and writes back the OAuth token the login flow obtains |
| `internal/hub` | WebSocket fan-out: broadcast + per-bridge-lifetime backlog replay for late joiners, inbound events from any consumer |
| `internal/claudecli` | `claude` stream-json translation (`translate.go`, pure) and per-turn subprocess management (`runner.go`) |
| `internal/permissions` | the `PreToolUse` HTTP hook endpoint plus the decision broker that blocks it on a consumer's `permission_decision` |
| `internal/authlogin` | the `claude setup-token` device-login flow, driven over a pty |

## How permission prompts work

`claude` has no `--permission-prompt-tool`. The extension point is a **`PreToolUse` hook of type
`http`**, registered in a settings file the bridge writes at startup and passes to every `claude`
invocation via `--settings`. Claude Code POSTs the tool call to the bridge and blocks
synchronously on the response, which is what lets a human approve a tool call mid-turn:

```
claude  --PreToolUse-->  bridge :8787/hook  --permission_request-->  WS consumer
        <--allow/deny--                     <--permission_decision--
```

A hook response is always an explicit `allow` or `deny`; the bridge never relies on Claude Code's
own fallback, and never passes a permissive `--permission-mode`. `allow_always` is remembered per
tool name for the lifetime of the bridge process.
