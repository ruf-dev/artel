# Auth & Login Flow

`StartWorkbench(ctx, vaultID, authMode)` takes an explicit `authMode` — the user (or the frontend,
based on what's available) picks per-start, it's not inferred silently. Two modes:

## `api_key` mode

Requires an existing BYOK Anthropic connection for the workbench's owning user
(`external_connections`, `provider='anthropic'` — see `docs/byok/01_data_model_and_security.md`).
`StartWorkbench`:

1. Look up the user's Anthropic `external_connections` row; if none exists, fail fast with a
   clear error (`user_errors.WorkbenchMissingAnthropicConnection` or similar) — don't silently
   fall back to `subscription_login`, the caller asked for a specific mode.
2. Decrypt via the same `cryptoutil`/repo path BYOK already uses — no new decryption code.
3. `docker start` the container with `ANTHROPIC_API_KEY` set as an env var passed at start time
   (not baked into the image, not written to the volume).
4. Entrypoint launches `claude` inside a `tmux` session; with the env var present, `claude`
   authenticates immediately, no login prompt.

This is the straightforward path — reuses 100% of BYOK's existing storage/decryption, no new
credential surface. **Build and validate this mode first**; it has no unresolved unknowns.

## `subscription_login` mode

No key injected. Entrypoint launches `claude` with no `ANTHROPIC_API_KEY`, which (per Claude
Code's existing support for headless/CI/SSH environments) should trigger its own account-login
flow rather than erroring — the premise being that a subscription (Pro/Max) login draws from
included usage instead of metered API billing, which is the actual economic reason this mode
exists (see the discussion that motivated this doc: BYOK API pricing vs. subscription-included
usage for a 24/7 agent).

### Task 0 findings — confirmed by direct spike (2026-07-22)

Ran the actual `deploy/workbench/Dockerfile` image with no `ANTHROPIC_API_KEY`, attached via
`docker exec ... tmux capture-pane`, and drove it via `tmux send-keys`. Neither of the two guessed
shapes is exactly right — the real mechanism is a third, simpler one:

- On first run, `claude` shows a one-time interactive **theme-selection prompt**, then a
  **login-method menu** (`Claude account with subscription` / `Anthropic Console account` /
  `3rd-party platform`). Both are ordinary keypress-driven TUI screens — advanced the same way
  (`tmux send-keys ... Enter`).
- After picking "Claude account with subscription", `claude` prints "Opening browser to sign
  in…", then — since there's no browser in the container — falls back on its own to: **"Browser
  didn't open? Use the url below to sign in (c to copy)"** followed by a full OAuth authorize URL,
  then a **`Paste code here if prompted >`** input line.
- Critically, the printed URL's `redirect_uri` is `https://platform.claude.com/oauth/code/callback`
  — an **Anthropic-hosted page**, not `localhost`. So there is no loopback listener started inside
  the container at any point; the browser-side redirect never needs to reach the container at all.
  This rules out the "local-redirect" shape entirely — no port-forward exception is needed.
- It is also not a pure device-code/polling flow — `claude` does not silently poll Anthropic's
  servers to detect completion. Instead, Anthropic's hosted callback page is expected to display a
  short code to the user, which the user copies and **pastes back** into the CLI's waiting input
  line. Confirmed the injection path works mechanically: sent a throwaway invalid string via
  `tmux send-keys -t workbench "<string>" Enter` and `claude` read it as real stdin input,
  responding `OAuth error: Invalid code. Please make sure the full code was copied` — proving
  `docker exec` + `tmux send-keys` is the correct, working primitive for relaying the user's pasted
  code back into the session (not `docker attach`, not an SDK-level Attach stream).
- Session persistence: the lifecycle design already reuses the same container across
  `running ⇄ stopped` transitions (`docker stop`/`docker start`, never `rm`, per
  `01_data_model_and_lifecycle.md`) — so whatever `claude` writes to disk after a successful login
  survives restarts regardless of whether it happens to live on the mounted volume. Only a full
  `removed` teardown (`docker rm` + volume removal) loses it, which is expected and fine — the user
  re-authenticates the next time a workbench is (re-)created for that vault.

**Conclusion: Artel is a pure URL-out / code-in relay, nothing more.** No polling of Anthropic's
API is needed on Artel's backend at all — the actual OAuth exchange happens entirely between the
user's own browser and Anthropic's servers, mediated only by what the user manually copies and
pastes. This is simpler than either guessed shape.

### Mechanism (confirmed)

1. `docker start` with no auth env vars.
2. Backend reads the container's tmux pane (`docker exec ... tmux capture-pane -t workbench -p`,
   polled — no long-lived stream needed) looking for a line matching the printed authorize URL.
3. A new RPC (e.g. `GetWorkbenchLoginPrompt(vaultID)`) surfaces "no prompt yet" / "here's the URL"
   to the frontend, which displays it for the user to open in any browser on any device.
4. A second new RPC (e.g. `SubmitWorkbenchLoginCode(vaultID, code)`) takes whatever the user pastes
   back from Anthropic's callback page and relays it verbatim via
   `docker exec ... tmux send-keys -t workbench "<code>" Enter`.
5. Backend keeps polling `tmux capture-pane` after submission: an `OAuth error: ...` line means
   retry (resurface the input prompt), the login-method-menu/URL disappearing without an error
   means success — mark the workbench `running` with `auth_mode='subscription_login'`.
6. The same generic "read pane, relay keystrokes" primitive also naturally handles the one-time
   theme-selection prompt on a brand-new container — no special-casing needed, Task 7 can just
   forward the user's first keystroke(s) through the same channel before the login URL even
   appears, rather than hardcoding an "auto-press Enter for theme" step.

No token/session is stored by Artel in this mode — the CLI manages its own session state on the
container's volume, same as it would on a user's own machine. Nothing here needs to touch
`external_connections` at all; this mode is deliberately separate from BYOK, not a variant of it.
