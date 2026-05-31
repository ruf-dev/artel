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

## Naming
- **Interface naming: no `I` prefix.** Use `AudioService`, not `IAudioService`.
- **Receiver names: single letter.** Use the type's first letter: `s *Service`, `a *App`.
- **Package names: lowercase, single word.** e.g. `authimpl`, `pg`.
- **Type names: PascalCase; constants: SCREAMING_SNAKE_CASE.**

## Functions
- **Named declarations only.** No anonymous assignments: `var f = func(){}` is forbidden; use `func name() {}`.
- **Context is always the first parameter** in any function that touches I/O.
- **Constructor pattern: `func New(deps) *Type`.** Constructors accept interfaces, never concrete types.
- **Deferred cleanup: `defer utils.CloseWithLog(resource, "description")`.** Use this helper, not bare `defer resource.Close()`.

## Errors
- **Error wrap messages must start with a verb:** `rerrors.Wrap(err, "error creating user")`, `"error upserting session"`.
- **Attach gRPC codes at definition time:** `rerrors.New("unauthorized", codes.Unauthenticated)`. Never set codes at the call site.
- **Attach HTTP status:** `rerrors.WithHttpStatus(http.StatusUnauthorized)`. Chain onto the error at definition.
- **Never use naked `errors.New` or `fmt.Errorf`.** All errors go through `rerrors`.
- **Call `rerrors.SetSeparator(':')` in `main()`.**

## Proto / gRPC Rules

- All RPCs must use `POST` with `body: "*"` in the HTTP gateway annotation:
  ```proto
  option (google.api.http) = {
      post: "/api/your-service/method"
      body: "*"
  };
  ```
