# Frontend: Connections Page Restructure

## Layout

```
HeroSegment ("Connected services", N services)
  Tabs: [ External Connections | BYOK ]
  ContentSegment (tab-aware)
```

Tabs sit directly under the hero, above the card grid — per the request ("new tab will be
located under the hero of the connections page"). `pages/connections/ConnectionsPage.tsx` starts
owning tab state (`useState<"external" | "byok">("external")`, or a `searchParams` entry if deep
linking to the BYOK tab is wanted, e.g. from a Tract step editor's "add a key" shortcut) and
passes it down to `ContentSegment`.

## No `Tabs` primitive exists yet

Grepped: chures has no `Tabs` export today. This needs a small new component before the page can
be built. Two placements are defensible under the ArtelUI tiering rules
(`pkg/client/ArtelUI/CLAUDE.md`):

- **Tier 2 atom** (`components/atoms/Tabs/Tabs.tsx`) if it's treated as "chures doesn't have this
  primitive yet" — same justification as `components/atoms/Input`. Needs the required one-line
  gap comment (`// TODO: chures has no tab primitive yet, drop this wrapper once it does`).
- **Tier 3** (`components/Tabs/Tabs.tsx`) if it's judged more composite/product-specific than a
  raw-element wrapper.

Recommendation: **tier 2 atom**, since it's a genuine "chures gap," and it'll likely be reused
elsewhere (this codebase has precedent for multi-section pages — e.g. dialogs with
screens/steps). Props: `{tabs: {key: string; label: string}[]; active: string; onChange: (key:
string) => void}`. No `z-index` — the active-tab indicator must come from DOM order / a
`::after` border on the active button, not stacking (per the project's "never use z-index"
rule).

## `ContentSegment.tsx` becomes tab-aware

Current `ContentSegment` unconditionally renders the External Connections grid. Split it:

- `pages/connections/components/ContentSegment/ContentSegment.tsx` keeps rendering exactly what
  it renders today (Google Sheets/Trello/Miro/Email/GitLab `ProviderCard` grid) — this becomes
  the "External Connections" tab body, unchanged.
- New sibling `pages/connections/components/BYOKSection/BYOKSection.tsx` (colocated under
  `pages/connections/components/`, since it's page-local — see the ArtelUI colocation rules) owns
  the BYOK grid.
- `ConnectionsPage.tsx` renders one or the other based on the active tab. Neither needs to be
  mounted-but-hidden; a simple conditional render is fine (no shared scroll-position state to
  preserve across tabs today).

## `BYOKSection.tsx`

Two groups of cards, in this order:

1. **LLM keys** — one `ProviderCard`-equivalent per configured connection
   (`ExternalProvider.EXTERNAL_PROVIDER_ANTHROPIC` today, `_OPENAI` once stage 2 ships), reusing
   the existing `widgets/ProviderCard/ProviderCard.tsx` shape (icon, name, connected-state,
   click-to-manage) — same component the External Connections tab already uses, not a fork. A
   new `LLM_BYOK_PROVIDERS` array (mirrors the existing `PROVIDERS` const in `ContentSegment.tsx`)
   lists `{provider: EXTERNAL_PROVIDER_ANTHROPIC, name: "Claude (Anthropic)"}` and — greyed
   out/disabled — `{provider: EXTERNAL_PROVIDER_OPENAI, name: "OpenAI"}` until stage 2.
2. **Mock/"coming soon" cards** — CouchDB instance, S3/MinIO bucket, WebDAV. These are *not*
   backed by any `ExternalProvider` enum value yet (don't add placeholder enum entries for
   things with no service behind them) — just static, disabled cards with a tooltip. A new
   `ComingSoonCard` component (project-wide, `components/ComingSoonCard/ComingSoonCard.tsx`,
   tier 3 — used for 3 distinct placeholders) takes `{icon, name, hint}` and renders a
   non-interactive card (no `onClick`, `aria-disabled`) with the hint delivered via tooltip.

### Tooltip — correction: already exists, no new atom needed

Earlier drafts of this doc assumed chures had no tooltip primitive. That's wrong — the app already
has a global `react-tooltip` instance mounted in `src/app/routing/Router.tsx` (`<Tooltip
id="root-tooltip" .../>`), and the existing convention (see `components/TemplateInput/TemplateInput.tsx`,
`components/ProviderChip/ProviderChip.tsx`, etc.) is simply to add `data-tooltip-id="root-tooltip"`
+ `data-tooltip-content="..."` to any element that needs a hint. `ComingSoonCard` and the disabled
OpenAI tab should use this existing pattern — **no new Tooltip atom**.

## `ManageLlmKeyDialog`

New dialog, `dialogs/ManageLlmKeyDialog/ManageLlmKeyDialog.tsx`, structured like
`ManageGitlabDialog.tsx`/`ManageTrelloDialog.tsx`:

- Not connected → `ConnectForm` (colocated `dialogs/ManageLlmKeyDialog/components/ConnectForm/`):
  - Vendor tabs: "Claude" (active) / "OpenAI" (visibly present, disabled, with a tooltip
    explaining "coming soon" — this is the one place stage 2's existence is hinted at in the UI
    now, per "keep in mind there is also another type of connection").
  - API key input (`type="password"`-style masked field — check whether `components/atoms/Input`
    already supports a masked variant; if not, that's a small atom extension, not a new
    component).
  - Optional "Base URL" field, collapsed behind an "Advanced" disclosure — most users never touch
    it; it exists for Claude Platform on AWS / regional endpoints.
  - "Test connection" button → `checkLlmKeyConnection`. On success, shows the detected vendor
    label + a preview of the model list (confirms to the user *what* they just connected to,
    directly answering "how can we check it's working").
  - "Save" button → `addLlmKeyConnection`. Disabled until a successful Test (same spirit as
    requiring a working GitLab token before persisting).
- Connected → `ConnectedContent` (colocated, same folder pattern as GitLab's
  `ConnectedContent.tsx`):
  - Key preview (`sk-ant-...ab12`), vendor label, last-verified timestamp.
  - Default model picker (persists via `AddLlmKeyConnection`'s `default_model` field — used to
    pre-fill new Call LLM steps, not a hard constraint on them).
  - Usage summary for this connection — see [05_metrics_and_usage.md](05_metrics_and_usage.md).
  - "Rotate key" (opens `ConnectForm` pre-filled with the base URL/model, key field blank) and
    "Disconnect" (`ConfirmDialog`, same as GitLab's `handleDisconnect`).

## What does NOT change

- `pages/connections/GoogleOAuthCallbackPage.tsx` — untouched.
- Every existing External Connections dialog (`ManageEmailDialog`, `ManageGitlabDialog`,
  `ManageTrelloDialog`, `ConnectionDetailDialog`, `GoogleSheetsConnectionContent`) — untouched,
  they just now render under the "External Connections" tab instead of the page root.
- `useExternalConnections`/`ExternalConnectionsService` — extended (see
  [02_connection_lifecycle.md](02_connection_lifecycle.md)), not restructured.
