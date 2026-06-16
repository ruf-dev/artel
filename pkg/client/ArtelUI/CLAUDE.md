# ArtelUI Frontend Rules

This file applies to everything under `pkg/client/ArtelUI`. It is referenced from the
repo root [CLAUDE.md](../../../CLAUDE.md) — read it whenever you touch a file in this
directory.

## Feature-Sliced Design layers

- **`pages/<name>/<Name>Page.tsx`** — page shell only. Composes widgets/components and
  owns page-level effects (auth redirects, top-level data fetches via `app/hooks`).
  Small page-specific glue may stay inline in the page file *only if* it isn't reused
  elsewhere and stays small once the heavy pieces below are extracted:
  - a `HeroSegment`/`ContentSegment`-style layout function
  - a one-off "create X" dialog (see `CreateVaultDialog` in `HomePage.tsx`)
- **`widgets/<Name>/<Name>.tsx`** — the unit rendered per-item in a list/grid that owns
  its own related data fetching (e.g. `widgets/VaultCard`, `widgets/McpKeyCard`).
  Composes from `components/`.
- **`components/<Name>/<Name>.tsx`** — either a reusable leaf (atom/molecule: a chip,
  a generic option row, a form field) used in 2+ places, or a self-contained complex
  dialog (e.g. `ManageVaultDialog`, `ManageKeyDialog`). A dialog's own wizard
  steps/sub-rows stay as **private functions in the same file** (see
  `MemberRow`/`InviteRow` in `ManageVaultDialog.tsx`) — don't fragment one dialog
  across multiple files just because it has steps.
- **`app/hooks/`** — Zustand stores / data hooks, shared across pages, widgets, and
  components.

## When to extract vs. keep inline

- Extract to `components/` once a piece of UI or logic (a chip renderer, a
  `connectionLabel`/color-mapping helper, a generic select-option row) is used in
  **2+ places** — don't duplicate the JSX/logic itself. (Small CSS-module *rules* are
  fine to duplicate — see below.)
- Extract to `widgets/` when the thing is rendered per-item in a list and fetches its
  own related data.
- Keep things inline in the page when they're page-specific and used exactly once.

## CSS Modules

- Every `widgets/`/`components/` file gets its own `*.module.css`, scoped to what
  that file renders.
- **Duplicate small shared rules** (button variants, chip styles, modal chrome)
  across modules rather than inventing a shared global stylesheet — this matches the
  existing `HomePage.module.css` / `ManageVaultDialog.module.css` precedent and keeps
  each file self-contained and movable.

## Reference example

`pages/home/HomePage.tsx` + `widgets/VaultCard/VaultCard.tsx` +
`components/ManageVaultDialog/ManageVaultDialog.tsx` is the canonical shape. When a
page file grows past ~300 lines or starts mixing dialog logic, list-item rendering,
and shared utilities in one file, split it the same way — see the `mcp-keys`
refactor (`widgets/McpKeyCard`, `components/ManageKeyDialog`,
`components/ConnectorChip`, `components/SelectOption`) for a worked example of
splitting an existing fat page.

## Component Structure

- **Never create components with more than 3 levels of HTML nesting** — split into
  smaller components instead
- Top-level container element's style class must be named `***Container` (e.g.,
  `HeaderContainer`)
- When wrapping another component with a styled div, use `***Wrapper` for that div's
  style (e.g., `ButtonWrapper`)

## Error and Confirmation Handling

- **Never use `window.alert` or `window.confirm`** — use project-level primitives:
  - **Errors**: `useBakeError()` from `@/app/hooks/useErrorToast` → call
    `bakeError(title, err)` inside `catch` blocks
  - **Confirmations**: `OpenDialog(<ConfirmDialog ... />)` from
    `@/components/ConfirmDialog/ConfirmDialog`
- `ConfirmDialog` props: `title`, `message`, `confirmLabel`, `cancelLabel`, `danger`
  (boolean), `onConfirm` (async callback)
- The `onConfirm` callback is responsible for `try/catch/finally`; `ConfirmDialog`
  closes itself in `finally` after `onConfirm` resolves
