---
id: "11"
title: "Telegram Auth"
status: "done"
model: "qwen2.5-coder:3b"

Plan: Telegram OIDC Auth Feature

  Flow (Login Library mode — what the article describes)

  1. Frontend loads telegram-login.js from oauth.telegram.org
  2. User clicks Telegram button → popup opens; Telegram runs a full PKCE OIDC flow internally
  3. Telegram delivers { id_token, user } via postMessage callback (onTelegramAuth)
  4. Frontend POSTs { id_token } to backend
  5. Backend verifies the JWT (RS256/ES256) against JWKS at https://oauth.telegram.org/.well-known/jwks.json, validates aud=client_id, iss, exp, then upserts user by
  sub (telegram_id) and creates a session

---

Step-by-step implementation plan

Step 1 — DB migration (migrations/007_telegram_auth.sql)

- Users table have new field - username.
- New table identities_telegram created - it contains telegram id, user id and telegram session related data.

Step 2 — sqlc queries (internal/repository/pg/queries/users.sql)

Add two queries:
-- name: GetUserByTelegramId :one
SELECT ... FROM users JOIN identities_telegram ON ... WHERE telegram_id = $1;

-- name: UpsertTelegramUser :one
INSERT INTO identities_telegram (telegram_id) VALUES ($1)
ON CONFLICT (telegram_id) DO UPDATE SET updated_at = NOW()
RETURNING ...;
Then run sqlc generate to regenerate users.sql.go and models.go.

Step 3 — Domain (internal/domain/user.go)

- New domain structs: TelegramIdentity (telegram session), ArtelIdentity (log + pass) 
- New domain struct: UserIdentities containing all identities 
- Add UserIdentities field to domain.User
- User now have Username

Step 4 — Repository interface (internal/repository/interfaces.go)

type Users interface {
// existing...
GetByTelegramId(ctx context.Context, telegramId string) (domain.User, error)
UpsertByTelegramId(ctx context.Context, telegramId string) (domain.User, error)
}

Step 5 — Repository implementation

Implement the two new methods in the pg layer using the sqlc-generated queries.

Step 6 — Config (config/config.yaml)

Add under environment::
telegram_client_id: ""
Then run rscli-dev project tidy to regenerate internal/config/environment.go. This exposes
EnvironmentConfig.TelegramClientId.

Step 7 — Go dependencies

go get github.com/golang-jwt/jwt/v5
go get github.com/MicahParks/keyfunc/v3

- jwt/v5 for parsing/verifying JWT
- keyfunc/v3 for auto-caching JWKS from Telegram

Step 8 — Service interface (internal/service/interfaces.go)

type AuthService interface {
// existing...
LoginViaTelegram(ctx context.Context, idToken string) (domain.Session, error)
}

Step 9 — Auth service implementation (internal/service/v1/auth/auth.go)

- Store jwksClient in Service struct (initialized once at New() from config)
- Implement LoginViaTelegram:
- Get signing key from JWKS client using id_token
- Call jwt.ParseWithClaims with RS256/ES256 algorithms, validate aud and iss
- Extract sub (telegram_id)
- Call usersRepo.UpsertByTelegramId(ctx, telegramId)
- Call sessionsRepo.Create(ctx, user.Uuid, token, expiresAt)
- Return session

Step 10 — Proto (api/grpc/auth.proto)

message TelegramLogin {
message Request { string id_token = 1; }
message Response {
string token = 1;
google.protobuf.Timestamp expires_at = 2;
}
}
Use it inside Login rpc method

Then run moti g to regenerate Go bindings.


Step 12 — Frontend

src/app/api/api.ts — add loginTelegram(idToken: string) API call.

src/processes/Auth.ts — add LoginViaTelegram(idToken: string): Promise<Session> to AuthService.

src/pages/init/InitPage.tsx — replace the disabled "Telegram (coming soon)" button:

- Inject telegram-login.js script tag dynamically
- Expose window.onTelegramAuth that calls svc.LoginViaTelegram(data.id_token)
- Use bot_id from a Vite env variable (VITE_TELEGRAM_BOT_ID)

index.html — no change needed (script loaded dynamically)

vite.config.ts — expose VITE_TELEGRAM_BOT_ID env var (already supported by Vite)

---
Key decisions

┌─────────────────────────────────────┬───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│ Question │ Decision │
├─────────────────────────────────────┼───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│ Email handling for Telegram users │ Allow NULL — already DEFAULT '', just drop NOT NULL │
src/app/api/api.ts — add loginTelegram(idToken: string) API call.

src/processes/Auth.ts — add LoginViaTelegram(idToken: string): Promise<Session> to AuthService.

src/pages/init/InitPage.tsx — replace the disabled "Telegram (coming soon)" button:

- Inject telegram-login.js script tag dynamically
- Expose window.onTelegramAuth that calls svc.LoginViaTelegram(data.id_token)
- Use bot_id from a Vite env variable (VITE_TELEGRAM_BOT_ID)

index.html — no change needed (script loaded dynamically)

vite.config.ts — expose VITE_TELEGRAM_BOT_ID env var (already supported by Vite)

---
Key decisions

┌─────────────────────────────────────┬───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│ Question │ Decision │
├─────────────────────────────────────┼───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│ Email handling for Telegram users │ Allow NULL — already DEFAULT '', just drop NOT NULL │
├─────────────────────────────────────┼───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│ JWKS caching │ MicahParks/keyfunc/v3 — handles auto-refresh with TTL │
├─────────────────────────────────────┼───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│ Where is client_id? │ Backend reads from EnvironmentConfig.TelegramClientId; frontend reads VITE_TELEGRAM_BOT_ID (
these are the same value) │
├─────────────────────────────────────┼───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│ Token format │ Reuse existing opaque hex token + DB session (no new JWT infrastructure)
│
├─────────────────────────────────────┼───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│ Separate telegram_identities table? │ No — simple telegram_id column on users; add later if multiple identity
providers are needed │
└─────────────────────────────────────┴───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘
