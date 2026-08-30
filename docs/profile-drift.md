# Drift from the coder profile

The machine-wide coder profile (`~/.claude/rules/`, source of truth in the Artel vault at
`.profiles/developer/`) was rebuilt on 2026-08-29 from a preference quiz run against this repo.
Several answers changed rules artel already had, so artel is now knowingly non-conforming in the
places below.

**Read this before "fixing" any of them** — they're a tracked baseline, not accidents. Fix the
ones you're already touching; don't sweep.

## Go

| # | Rule | artel today |
| --- | --- | --- |
| 1 | Imports in three groups (stdlib / third-party / own module) | Two groups — goimports default, no `-local`/gci config. Needs a formatter setting before it's enforceable. |
| 2 | No `f`-suffixed log calls; values as typed field chains | ~40 `log.Infof` calls |
| 3 | Transactions via a `TxManager.Execute` helper | Hand-rolled `BeginTx` / defer rollback / `Commit`; repos already expose `WithTx` |

## Proto

| # | Rule | artel today |
| --- | --- | --- |
| 4 | GET with path params for reads, POST + `body: "*"` for mutations and lists | **Every** RPC across 10+ `.proto` files is POST + `body: "*"`. This reverses artel's previous rule — the whole API surface is non-conforming. Largest single item here. |

## Frontend

| # | Rule | artel today |
| --- | --- | --- |
| 5 | Everything PascalCase, including page directories | `pages/` uses kebab-case (`tract-canvas`, `closed-alpha`, `mcp-auth`) |
| 6 | Double quotes, no semicolons, 120-char max-len, no `object-curly-spacing`, `any` banned | None of these rules exist in `eslint.config.js`; 53 single-quote imports in test files, 148 lines over 120 chars, 4 files using `any` |
| 7 | Portrait is mobile — no width check | `sizes.css` sets `--is-mobile` from `(orientation: portrait), (width <= 45rem)` |
| 8 | rem lengths and `var()` colors only | ~560 stylelint `unit-disallowed-list` warnings (known baseline) |
| 9 | No raw `<button>`/`<input>`, no template-literal `className` | 105 `no-restricted-syntax` warnings (known baseline) |

## Testing

| # | Rule | artel today |
| --- | --- | --- |
| 10 | E2E covering critical business flows is mandatory; unit tests optional | Replaces the old "all new code covered with tests". E2E suites exist under `tests/`; the gap is coverage breadth, not shape. |

## Migrations

| # | Rule | artel today |
| --- | --- | --- |
| 11 | Every migration ships a Down that exactly reverses the Up | `081_trello_write_tools.sql` adds 5 `trello.*` labels to the `mcp_tool_name` enum via `ALTER TYPE ... ADD VALUE`; its Down removes the tool rows but leaves the labels — Postgres has no `DROP VALUE`. The `mcp_name` / `mcp_tool_name` enums (migration 080) are consumed through their sqlc-generated Go types (`artel_q.McpName` / `artel_q.McpToolName`) and are attached to no column, so a stale label is inert. |

## Not drift — deliberately artel-only

These stay in artel and are **not** part of the profile: the `[Area] Category:` commit convention
and its hook, branch naming, MoM as the default for third-party integrations, and the
`moti`/`sqlc`/`make test-e2e` workflow specifics.
