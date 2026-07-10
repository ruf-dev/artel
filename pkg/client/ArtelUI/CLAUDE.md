# ArtelUI Frontend Rules

This file applies to everything under `pkg/client/ArtelUI`. It is referenced from the
repo root [CLAUDE.md](../../../CLAUDE.md) — read it whenever you touch a file in this
directory.

## Component hierarchy

Five tiers, lowest to highest. **A file may only import from its own tier or a tier
strictly below it** — never sideways across an unrelated subtree, never upward. This
is enforced for tiers 2-5 by `import-x/no-restricted-paths` in `eslint.config.js`
(plus `import-x/no-cycle` catching any cycle regardless of tier); the one thing the
linter can't check is colocation scope (below), which is a review-time rule.

1. **`@vervstack/chures`** — the UI atom library: `Button`, `Input`, `Dropdown`,
   `Toggle`, `ConfirmDialog`, `InfoDialog`, `ModalActions`, `ModalClose`, `Loader`,
   `LoadingWrapper`, `Toaster`/`useToaster`, and the `icons/` set. Always reach for a
   chures component instead of a raw `<input>`, `<button>`, or a bespoke
   confirm/alert dialog. Enforced by the `no-restricted-syntax` rule banning raw
   `<button>`/`<input>` JSX (an exception carves out `components/atoms/**`, since the
   wrapper there has to render the primitive once).
2. **`components/atoms/<Name>/<Name>.tsx` + `.module.css`** — thin wrappers that
   exist *only* because chures doesn't cover the need. Every file here must open with
   a one-line comment explaining the gap, e.g. `// TODO: chures has no multiline
   variant yet, drop this wrapper once it does`. May only import chures and other
   atoms (not tier-3 `components/`). `components/shared/Input` predates this
   convention and is the folder it replaces — don't add new atoms under `shared/`,
   and give it the TODO comment (and ideally move it under `atoms/`) next time you
   touch it.
3. **`components/<Name>/<Name>.tsx`** — project-wide components too specific to the
   product to be a chures-style atom: a reusable leaf used in 2+ places (a chip, a
   generic option row, a form field). May import tiers 1-2 only.
4. **`widgets/<Name>/<Name>.tsx`** and **`segments/<Name>/<Name>.tsx`** — same tier,
   two different roles:
   - A **widget** is a unit reused across pages that owns its own related data
     fetching (`widgets/VaultCard`, `widgets/McpKeyCard`).
   - A **segment** at this top level (`src/segments/Topbar`) is global app-shell
     chrome. Keep this set small and deliberate — it's for the handful of things
     that are genuinely app-wide (the top bar, the global dialog mount), not a
     general-purpose category. A segment specific to one page/dialog is a *local*
     segment (see colocation below), not a new entry under `src/segments/`.

   May import tiers 1-3 only.
5. **`pages/<Name>/<Name>Page.tsx`** and **`dialogs/<Name>/<Name>.tsx`** — top-level
   citizens. A page is a route; a dialog is opened via `OpenDialog(<X/>)` from
   `@/app/hooks/Dialog`. May import tiers 1-4, plus a page may import a dialog to
   open it (`OpenDialog(<CreateVaultDialog/>)` in `HomePage.tsx`).

   **Any tier may import from `dialogs/` to call `OpenDialog(<X/>)`** — a dialog is
   opened from wherever its trigger lives, not just from pages. A `components/` leaf,
   a `widgets/`/`segments/` unit, or a page can all reach into `dialogs/` for this one
   purpose (e.g. `OpenDialog(<FastSetupDialog/>)` from `widgets/VaultCard/VaultCard.tsx`).
   This is the one sanctioned upward exception to the tier rule, carved out of
   `import-x/no-restricted-paths` for every zone below `dialogs/`. **A dialog must
   never import a page** (enforced by `no-restricted-paths`) — dialogs don't know who
   opened them.

## Local colocation

Any page or dialog can define its own **local** components/widgets/segments in a
colocated subfolder: `pages/<Page>/components/`, `pages/<Page>/widgets/`,
`pages/<Page>/segments/`, `dialogs/<Dialog>/components/`, `dialogs/<Dialog>/widgets/`.
See `pages/tract-canvas/*`, `dialogs/AddTriggerDialog/*`, and
`dialogs/ManageVaultDialog/*` for the reference shape.

- **One React component per file.** A `.tsx` file renders exactly one component — no
  private/local helper components defined inline in the same file, however small (a
  row renderer, a badge, a sub-dialog). Split each into its own file under the
  nearest colocation folder above. See `dialogs/ManageVaultDialog/components/*`
  (`MemberRow`, `InviteRow`, `RoleBadge`, `DangerZoneText`, `RoleOption`,
  `CreateInviteLinkDialog`) for the reference shape — including a dialog-local
  sub-dialog (`CreateInviteLinkDialog`, opened via `OpenDialog` from within
  `ManageVaultDialog`) living in the owning dialog's local `components/`, since
  there's no separate "local dialogs" folder — a local sub-dialog is still a
  `components/`-tier concept relative to its parent.
- These are **local by default**: only the owning page/dialog (and its own
  descendants) may import them. Don't reach into another page's or dialog's local
  folder from outside it — if the same piece is needed in two places, promote it to
  the matching global tier (`src/components`, `src/widgets`) instead of importing
  across subtrees. This isn't linter-enforced (no generic way to express "same
  subtree only" as a static zone) — same as the `z-index` rule, catch it in review.
- **Screens** (`dialogs/<Dialog>/screens/`) are a dialog-only concept, one level
  below the dialog itself, for a dialog with multiple sequential steps (see
  `AddTriggerDialog/screens/`). Same local-only rule as above.

## Messy dialog/page logic

If a page or dialog's `.tsx` file is accumulating non-rendering logic (data shaping,
layout math, orchestration), extract it to a colocated `processes/<name>.ts` file
(see `pages/tract-canvas/processes/tractCanvasLayout.ts`) instead of letting the
component file grow — same tool as the project-wide `src/processes/` used by
`app/hooks`, just colocated when the logic is specific to one page/dialog.

## Known debt (documented, not migrated)

- `RunTractDialog`, `StepPickerDialog`, `S3InstanceFormDialog` are dialogs living
  under `src/components/` instead of `src/dialogs/` — legacy from before this doc
  existed. Don't repeat the pattern for new dialogs; move these opportunistically if
  you're already touching one, but it's not worth a dedicated migration.
- 105 pre-existing `no-restricted-syntax` warnings (raw `<button>`/`<input>`,
  template-literal `className`) are a known baseline as of this rule's introduction
  — fix the ones you touch, don't feel obligated to sweep the whole codebase.
- ~560 pre-existing stylelint `unit-disallowed-list` warnings (`px`/`em`/`vh`/`vw`
  used where a `rem` size token from `sizes.css` should be) are a known baseline as
  of this rule's introduction — same deal, fix a file's sizes when you're already
  touching it, don't sweep the whole codebase in one pass. New code should use an
  existing `--*` token from `sizes.css` or add a new one there rather than writing a
  raw non-rem unit.
- The global dialog mount lives at `pages/segments/Dialog.tsx`, not `src/segments/`
  where `Topbar` lives — two different locations for the same "app-level segment"
  concept. New app-level segments go in `src/segments/`; don't add a third location.

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
- **Combining multiple/conditional classes**: never build `className` with a template
  literal (`` `${cls.Foo} ${cond ? cls.Bar : ""}` ``) — use `cn()` from
  `@/app/utils/cn.ts` (a `classnames` wrapper) instead, e.g.
  `cn(cls.Foo, cond && cls.Bar)`. Enforced by the `no-restricted-syntax` ESLint rule
  banning template literals in `className`.

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
  smaller components instead. This is enforced by the `react/jsx-max-depth` ESLint rule
  (`max: 3`).
- Top-level container element's style class must be named `{ComponentName}Container`
  (e.g., `VaultCardHeaderContainer` in `VaultCardHeader.tsx`) — never a generic name
  like `Field` or `Root`.
- All other classes in that component's `.module.css` must be nested inside the
  `{ComponentName}Container` rule using native CSS nesting (`&`), following the
  DOM structure — no top-level sibling classes. See `VaultCardHeader.module.css`
  for the reference pattern.
- When wrapping another component with a styled div, use `***Wrapper` for that div's
  style (e.g., `ButtonWrapper`)

## Buttons

- **Never use a raw `<button>` element** — always use `Button` from
  `@vervstack/chures`. If chures' `Button` can't do what you need, wrap it as a
  `components/atoms/` component (tier 2 above) with a TODO explaining the gap —
  never fork the raw element inline. Enforced by `no-restricted-syntax`.

## Error and Confirmation Handling

- **Never use `window.alert` or `window.confirm`** — enforced by
  `no-restricted-syntax`. Use project-level primitives instead:
  - **Errors**: `useBakeError()` from `@/app/hooks/useErrorToast` → call
    `bakeError(title, err)` inside `catch` blocks (backed by chures' `useToaster`)
  - **Confirmations**: `OpenDialog(<ConfirmDialog ... />)` — `ConfirmDialog` is
    imported straight from `@vervstack/chures`, no local wrapper
- `ConfirmDialog` props: `title`, `message`, `confirmLabel`, `cancelLabel`, `danger`
  (boolean), `onConfirm` (async callback)
- The `onConfirm` callback is responsible for error handling; `ConfirmDialog`
  closes itself in `finally` after `onConfirm` resolves

## State ownership in components

- **Private components** (file-local functions, not exported) that need to change **global state** (Zustand stores, `useDialog`, `useNavigate`, etc.) must call the relevant hook directly — do not thread the action down as a prop.
- **Private components** that need to change **parent local state** (e.g. `useState` in the enclosing component) receive a callback prop for that change — local state belongs to whoever owns it.
- **Public/exported components** that trigger state changes always receive a callback prop — they must not reach into a specific store themselves, because callers control which state is affected.

## Never use z-index

- **`z-index` is forbidden** in every `.module.css`/`.css` file and in inline `style` props —
  no exceptions, no "just this once."
- Stacking order comes from **DOM order** instead: within a stacking context, positioned
  elements paint in document order, so put the element that should be on top **later in the
  DOM** (use flexbox `order` if you need it to *look* earlier visually — `order` doesn't
  affect paint/stacking order, only layout position).
- Floating UI (dropdowns, popovers, tooltips) that must render above unrelated sibling
  content should be **portaled to `document.body`** (`createPortal`) and positioned via
  `getBoundingClientRect()`, not stacked with `z-index` — see
  `components/TemplateInput/TemplateInput.tsx` for the pattern.
- Global overlays (dialogs, toasts) already work without `z-index` because they're mounted
  last, as a sibling after `<Routes>`, in `pages/segments/Dialog.tsx` — follow that precedent
  for any new global overlay instead of reaching for `z-index`.
- Enforced for JS/TSX by the `no-restricted-syntax` rule banning `zIndex` in `eslint.config.js`.
  There is no CSS linter in this project (see root CLAUDE.md — only add one if the user asks),
  so `z-index` in `.module.css` files must be caught in code review.

## Async style

- **Prefer promise chains over `try/catch`** — use `.then().catch().finally()` instead
  of `async/await` with `try/catch` blocks. Exception: best-effort fire-and-forget
  where no error surface is needed (silent `catch {}` is fine there).

## Testing frontend changes

- **Do not start the dev server, spin up a browser, or otherwise self-test frontend
  changes.** Reaching an authenticated screen requires the full Go backend + DB, which
  isn't worth spinning up for a UI change, and a headless smoke test is a poor
  substitute for a human actually looking at it.
- Verify with `tsc -b`/`bun run build` and lint instead, then hand the change back to
  the user with a short note on what to click through to confirm it visually
  (which page/component, what interaction to try).
