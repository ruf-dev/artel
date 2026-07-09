# Handoff: lint/tooling parity gaps vs. ZpotifyUI

Generated while porting ArtelUI's rules (layer enforcement, banned raw DOM elements,
`z-index` ban, structural limits) over to `~/ruf/zpotify/pkg/client/ZpotifyUI` — the
sibling frontend sharing `@vervstack/chures`, `classnames`, `framer-motion`. That pass
went one direction (Artel → Zpotify). This is the reverse: things ZpotifyUI's
`eslint.config.js` / tooling already enforces that ArtelUI only has as unenforced
`CLAUDE.md` prose, or doesn't have at all. Nothing here is implemented in Artel yet —
this is a punch list for whoever picks it up next.

## No CSS linter at all

Zpotify has a `stylelint.config.js` (extends `stylelint-config-standard` +
`-standard-scss` + `-css-modules`) wired into `bun lint`. Artel has **no CSS linting
whatsoever** — every CSS-only rule in Artel's own `CLAUDE.md` (no `z-index`, no
`!important`, rem-only font sizes) is currently "catch it in review" by explicit
admission in the doc. Adding stylelint would let Artel actually enforce:
- `declaration-no-important: true`
- `declaration-property-value-disallowed-list: { 'z-index': [/.*/] }` — the CSS half
  of the existing `z-index` ban (Artel's `no-restricted-syntax` rule only catches
  `zIndex` in JS/TSX object literals, not `.module.css` files)
- `declaration-property-unit-disallowed-list: { 'font-size': ['px', 'em'] }` — Artel's
  CLAUDE.md says "use rem units for font sizes" but nothing enforces it
- `selector-class-pattern` for CSS Modules class naming

Zpotify's `stylelint.config.js` is a ready-made template — copy it and adjust the
per-file overrides (Zpotify has `src/colors_and_type.css`/`src/sizes.css` token-file
carve-outs that won't map 1:1 to Artel's file names).

## No Prettier

Zpotify runs `eslint-plugin-prettier` + `prettierRecommended` as part of the same
flat config, with a `.prettierrc` (`semi: true`, `singleQuote: true`, `tabWidth: 4`,
`trailingComma: all`, `printWidth: 120`, `jsxSingleQuote: false`). Artel has no
Prettier dependency, config, or lint integration — formatting is presumably
unenforced/manual. Worth a deliberate decision (not just a silent copy) since
introducing Prettier to an existing codebase reformats every file it touches on
`--fix`.

## Structural/style ESLint rules Zpotify enforces that Artel only documents

Artel's `CLAUDE.md` states several of these as prose conventions but has no matching
ESLint rule — meaning they're review-time-only today, same category as the `z-index`
CSS gap above:

- **`func-style: ['warn', 'declaration', { allowArrowFunctions: false }]`** — Artel's
  "Coding rules" equivalent (components as named functions, no `const fn = () => {}`)
  isn't in ArtelUI's `CLAUDE.md` at all currently and has zero lint backing.
- **`react/no-multi-comp: ['warn', { ignoreStateless: false }]`** — Artel's "One React
  component per file" rule (CLAUDE.md "Local colocation" section) is unenforced.
- **`max-lines: ['warn', { max: 300, ... }]`** and **`max-lines-per-function: ['warn',
  { max: 100, ... }]`** — Artel's CLAUDE.md gestures at this informally ("When a page
  file grows past ~300 lines... split it") but has no lint teeth.
- **`react/forbid-component-props: ['warn', { forbid: ['style'] }]`** — bans the
  inline `style` prop outright (CSS Modules only). Artel doesn't ban it explicitly
  anywhere.
- **`no-console: ['warn', { allow: ['warn', 'error'] }]`** — no equivalent in Artel.
- **`no-restricted-syntax` / `ObjectPattern[properties.length>6]`** — interesting
  inversion: this rule's *message* is already copy-pasted verbatim into Artel's own
  `CLAUDE.md` ("State ownership in components" section references destructuring, and
  a `>6 destructured bindings` rule text exists in Zpotify's CLAUDE.md that reads like
  it originated from a shared convention) — but Artel's `eslint.config.js` never
  actually wired up the rule. Worth checking whether this was meant to be there from
  the start.
- **`import/order`** — Zpotify enforces import grouping (react/react-router first,
  `@/**` next, relative/CSS last, blank line between groups). Artel has no import
  ordering rule; imports are whatever order was typed.
- **`local/no-relative-imports`** (custom rule banning `'./'`/`'../'` imports, forcing
  `@/` alias) — Artel's `tsconfig.json` already defines the same `@/*` path alias
  (`tsconfig.json:26-27`), so the alias itself works, but nothing stops relative
  imports from creeping in. Zpotify's custom rule (`eslint.config.js:11-30`, a
  ~15-line inline plugin) is directly portable.
- **`@typescript-eslint/no-restricted-imports` scoping generated API clients to one
  layer** — Zpotify bans importing `@/app/api/**` (gRPC clients) from anywhere except
  `src/processes/` and `src/shared/api/`, forcing all backend calls through a single
  process-function layer. Artel doesn't have an equivalent boundary; not a drop-in
  copy since Artel's layer names differ, but the *pattern* (lint-enforce "generated
  client only callable from one designated layer") is worth considering for whatever
  Artel's gRPC-client entry point is.

## Suggested order of attack

1. `local/no-relative-imports` — cheapest, self-contained, immediate value given the
   alias already exists.
2. `func-style`, `no-multi-comp`, `no-console`, `forbid-component-props` — small
   config additions, low fallout risk, run `eslint .` after each to check for
   pre-existing violations and decide `warn` vs `error` per Zpotify's own precedent
   (start at `warn` if the codebase has a real baseline, document it in `CLAUDE.md`
   like the existing "105 pre-existing warnings" note).
3. `max-lines`/`max-lines-per-function`, `import/order` — same process, expect more
   fallout since these are codebase-wide shape rules.
4. Stylelint — bigger lift (new dependency, new config file, `bun lint` script
   change), but directly reusable from Zpotify's `stylelint.config.js`.
5. Prettier — flag as a separate decision, not a mechanical port; reformatting the
   whole repo is a one-way door worth discussing first.
