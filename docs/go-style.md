# Go Coding Rules

- **Never check function errors inline.** Always split: `err := func()` on one line, then `if err != nil` on the next.
- **Never create struct/value literals inline in a function call.** Assign to a named variable first, then pass it.
- `internal/app/app.go` and `internal/app/config.go` are generated — do not edit them. `custom.go` is the user-editable
  counterpart.
- Use `rerrors.Wrap(err, "context message")` for all error wrapping.
- **Never create errors in-place.** Define all user-facing errors in `internal/service/user_errors`, then reference or wrap them at the call site.
- **No all-caps field names.** Use mixed-case: `Id` not `ID`. For `uuid.UUID` fields use the type name as the suffix:
  `Uuid` (primary key), `UserUuid`, `VaultUuid`, etc.