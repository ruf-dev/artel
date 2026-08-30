# ArtelUI Frontend Rules — artel specifics

General frontend rules — component shape, the 5-tier import hierarchy, colocation, state
ownership, async style, memoization, CSS structure, the `z-index` and `!important` bans, the
`--is-mobile` technique, test shape — live in the machine-wide coder profile and are injected
automatically when a `.ts`/`.tsx`/`.css` file is edited:

- `~/.claude/rules/40-frontend-react.md`
- `~/.claude/rules/41-frontend-style.md`
- `~/.claude/rules/42-css.md`
- `~/.claude/rules/43-frontend-testing.md`

**Don't duplicate them here.** This file holds only what's specific to this codebase.

Where artel doesn't currently satisfy a profile rule, it's tracked in
[docs/profile-drift.md](../../../docs/profile-drift.md) — read that before "fixing" a baseline.

## Tier 1 — the chures atom library

`@vervstack/chures` provides: `Button`, `Input`, `Dropdown`, `Toggle`, `ConfirmDialog`,
`InfoDialog`, `ModalActions`, `ModalClose`, `Loader`, `LoadingWrapper`, `Toaster`/`useToaster`,
and the `icons/` set. Always reach for one of these before writing a wrapper.

If chures can't do what you need, add a `components/atoms/<Name>/` wrapper with a one-line comment
naming the gap, e.g. `// TODO: chures has no multiline variant yet, drop this wrapper once it does`.

`components/shared/Input` predates the `atoms/` convention and is the folder it replaces. Don't
add new atoms under `shared/`; give it the TODO comment (and ideally move it under `atoms/`) next
time you touch it.

## Error and confirmation primitives

- **Errors**: `useBakeError()` from `@/app/hooks/useErrorToast` → `bakeError(title, err)` inside
  `.catch()`. Backed by chures' `useToaster`.
- **Confirmations**: `OpenDialog(<ConfirmDialog ... />)`, imported straight from
  `@vervstack/chures` — no local wrapper. Props: `title`, `message`, `confirmLabel`,
  `cancelLabel`, `danger`, `onConfirm` (async). The `onConfirm` callback owns its own error
  handling; `ConfirmDialog` closes itself in `finally` after it resolves.

## Lint enforcement in this repo

| Profile rule | Enforced by |
| --- | --- |
| Tier import direction | `import-x/no-restricted-paths` + `import-x/no-cycle` in `eslint.config.js` |
| No raw `<button>`/`<input>` | `no-restricted-syntax` (carve-out for `components/atoms/**`) |
| No template-literal `className` | `no-restricted-syntax` |
| Max JSX depth 3 | `react/jsx-max-depth` |
| No `z-index` | `no-restricted-syntax` (JS) + `declaration-property-value-disallowed-list` (stylelint) |
| No `!important` | `declaration-no-important` |
| rem size tokens | `unit-disallowed-list` |
| Dialogs cap height + scroll | custom `artel/dialog-scrollable` rule (`stylelint-rules/dialog-scrollable.js`), scoped to `src/**/*Dialog/*Dialog.module.css` |

`bun run lint` runs both; `bun run lint:js` / `bun run lint:css` individually.

The one thing the linter can't check is **colocation scope** — "only the owning page/dialog may
import its local `components/`". Catch that in review.

## The `OpenDialog` upward exception

Any tier may import from `dialogs/` solely to call `OpenDialog(<X/>)` — carved out of
`no-restricted-paths` for every zone below `dialogs/`. A dialog must never import a page (also
enforced). See `widgets/VaultCard/VaultCard.tsx` opening `FastSetupDialog`.

## artel tokens worth knowing

- `--standard-button-height` (`sizes.css`) — shared height for compact/icon action buttons. Also
  backs `--chures-input-height`'s mobile value, so the chat input, model-selector dropdown, and
  composer buttons all shrink together. If a control's height doesn't match its neighbours,
  reuse this rather than inventing a value.
- `--dialog-max-height` (`80vh`) / `--dialog-max-height-lg` (`90vh`) — reach for `-lg` on dense
  multi-field forms or dialogs containing a list (`ManageVaultDialog`, `S3InstanceFormDialog`);
  the plain token is the default (`WebhookDetailsDialog`, `TokenRevealDialog`).
- `--content-padding-mobile` — shared mobile content padding.

## Reference shapes

- Canonical page: `pages/home/HomePage.tsx` + `widgets/VaultCard/VaultCard.tsx` +
  `components/ManageVaultDialog/ManageVaultDialog.tsx`.
- Splitting a fat page: the `mcp-keys` refactor — `widgets/McpKeyCard`,
  `components/ManageKeyDialog`, `components/ConnectorChip`, `components/SelectOption`.
- Colocation: `pages/tract-canvas/*`, `dialogs/AddTriggerDialog/*` (including `screens/`),
  `dialogs/ManageVaultDialog/components/*` (`MemberRow`, `InviteRow`, `RoleBadge`,
  `DangerZoneText`, `RoleOption`, and the dialog-local sub-dialog `CreateInviteLinkDialog`).
- Colocated logic: `pages/tract-canvas/processes/tractCanvasLayout.ts`.
- Portaled floating UI: `components/TemplateInput/TemplateInput.tsx`.
- The DOM-reorder + `column-reverse` stacking fix:
  `pages/notes/components/NotesSidebar/components/SidebarTopBar/*` and the admin
  `TabBar` / `AdminPage.tsx` `StickyTabsWrapper`.
- Specificity fix without `!important`: `DrawerCloseButton.tsx`/`.module.css`,
  `MobileTopBar.module.css`, `DesktopNotesShell.module.css`.

## Known debt — documented, not migrated

- `RunTractDialog`, `StepPickerDialog`, `S3InstanceFormDialog` live under `src/components/`
  instead of `src/dialogs/` — legacy. Don't repeat it; move them opportunistically if you're
  already in one, but it isn't worth a dedicated migration.
- The global dialog mount is at `pages/segments/Dialog.tsx`, not `src/segments/` where `Topbar`
  lives. New app-level segments go in `src/segments/`; don't add a third location.
- Lint baselines: 105 `no-restricted-syntax` warnings, ~560 stylelint `unit-disallowed-list`
  warnings. Fix what you touch; don't sweep.
