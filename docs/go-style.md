# Go Coding Rules

- **Never check function errors inline.** Always split: `err := func()` on one line, then `if err != nil` on the next.
- **Never create struct/value literals inline in a function call.** Assign to a named variable first, then pass it.
- `internal/app/app.go` and `internal/app/config.go` are generated — do not edit them. `custom.go` is the user-editable
  counterpart.
- Use `rerrors.Wrap(err, "context message")` for all error wrapping.
- **Never create errors in-place.** Define all user-facing errors in `internal/service/user_errors`, then reference or wrap them at the call site.
- **No all-caps field names.** Use mixed-case: `Id` not `ID`. For `uuid.UUID` fields use the type name as the suffix:
  `Uuid` (primary key), `UserUuid`, `VaultUuid`, etc.
- **Repository functions must be pure DB operations.** No branching logic (`if user exists → update, else create`).
  Use `INSERT ... ON CONFLICT DO UPDATE` for upserts. Multi-step conditional flows belong in the service layer.
- **Repo functions that may return no row must return `sql.Null[T]`, not an error.**
  Absorb `sql.ErrNoRows` inside the repo and return `sql.Null[T]{Valid: false}, nil`.
  The service checks `if !result.Valid` — never `errors.Is(err, sql.ErrNoRows)` at the service layer.
- **Add `FOR UPDATE` to SELECT queries that will be followed by a write on the same row(s) inside the same transaction.**
  Safe to call outside a transaction — the lock is released immediately.
- **Repository functions that operate on more than one field of a domain entity should accept the domain struct, not individual fields.**
  Use `UpsertTelegramIdentity(ctx, identity domain.TelegramIdentity)`, not `(ctx, userUuid, telegramId, photoUrl)`.