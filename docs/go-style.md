# Go Coding Rules — artel specifics

The general Go rules (statements, naming, errors, imports, logging, layering, transactions,
repository contracts, test style) live in the machine-wide coder profile and are injected
automatically whenever a `.go` file is edited:

- `~/.claude/rules/30-go.md` — language, naming, errors, imports, logging
- `~/.claude/rules/31-go-layering.md` — transport/service/repository, transactions, domain purity
- `~/.claude/rules/32-go-testing.md` — assertions, table tests, mocks

**Don't duplicate them here.** This file holds only what's specific to this repo.

## artel paths

| Concept | Path |
| --- | --- |
| User-facing errors | `internal/service/user_errors` |
| Repository interfaces | `internal/repository/interfaces.go` |
| sqlc queries → generated | `internal/repository/pg/queries/` → `internal/repository/pg/generated/` |
| Service implementations | `internal/service/v1/{domain}/` |
| Transport handlers | `internal/transport/{name}_api/` |

## Generated — never edit

`internal/app/app.go` and `internal/app/config.go` are verv-generated (`verv project tidy`).
`internal/app/custom.go` is the editable counterpart.

## Known drift from the profile

artel does not currently satisfy every profile rule. See
[docs/profile-drift.md](profile-drift.md) for the list and its status — read it before
"fixing" something that's a known, deliberate baseline.
