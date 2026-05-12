---
status: todo
---

# Task 20 — UI: Design Tokens — Coral Palette + JetBrains Mono

## Goal

Migrate the color palette from teal/cyan to coral/red as specified in the design, and add JetBrains Mono
as the monospace font. No component changes — only CSS variable values and font imports.

## Files to modify

### 1. `pkg/client/ArtelUI/index.html`

Add a Google Fonts preconnect + stylesheet link in `<head>` before `</head>`:

```html
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=Comfortaa:wght@300;400;500;600;700&family=JetBrains+Mono:wght@300;400;500&display=swap" rel="stylesheet">
```

### 2. `pkg/client/ArtelUI/src/colors_and_type.css`

Replace the `:root` block entirely with the following. Keep ALL existing variable names
so existing components don't break — just change values and add new tokens:

```css
:root {
  /* ── Coral accent (new primary accent) ─────────────────── */
  --coral:              #FF4B3E;
  --coral-glow:         rgba(255, 75, 62, 0.35);
  --coral-dim:          rgba(255, 75, 62, 0.12);
  --coral-border:       rgba(255, 75, 62, 0.30);

  /* Backgrounds */
  --color-bg-base:        #000000;
  --color-bg-secondary:   #0d0d0d;
  --color-bg-accent:      #141414;
  --color-bg-auth:        #0d0d0d;

  /* Foregrounds */
  --color-fg-primary:     rgba(255,255,255,0.92);
  --color-fg-secondary:   rgba(255,255,255,0.55);
  --color-fg-tertiary:    rgba(255,255,255,0.32);
  --color-fg-accent:      #FF4B3E;
  --color-fg-disabled:    rgba(255,255,255,0.18);

  /* Borders */
  --color-border:         rgba(255,255,255,0.06);
  --color-border-light:   rgba(255,255,255,0.11);
  --color-border-bright:  rgba(255,255,255,0.22);

  /* Status */
  --color-error:          #c1253d;
  --color-warning:        #ffa500;
  --color-info:           rgba(255,255,255,0.55);

  /* ── Semantic Aliases (keep names, update values) ───────── */
  --main-bg-color:              var(--color-bg-base);
  --secondary-bg-color:         var(--color-bg-secondary);
  --accent-bg-color:            var(--color-bg-accent);
  --main-fg-color:              var(--color-fg-primary);
  --secondary-fg-color:         var(--color-fg-secondary);
  --thirdy-fg-color:            var(--color-fg-tertiary);
  --accent-fg-color:            var(--color-fg-accent);
  --disabled-fg-color:          var(--color-fg-disabled);
  --element-border-color:       var(--color-border);
  --element-border-color-lighter: var(--color-border-light);
  --error-color:                var(--color-error);
  --warning-color:              var(--color-warning);
  --info-color:                 var(--color-info);

  /* ── Layout Tokens ──────────────────────────────────────── */
  --header-height:        56px;
  --border-radius:        1em;
  --border-radius-pill:   999px;
  --border-radius-circle: 999px;
  --action-button-size:   2em;

  /* ── Spacing Scale ──────────────────────────────────────── */
  --space-1:  0.25em;
  --space-2:  0.5em;
  --space-3:  0.75em;
  --space-4:  1em;
  --space-5:  1.5em;
  --space-6:  2em;
  --space-8:  3em;

  /* ── Elevation / Shadow ─────────────────────────────────── */
  --shadow-glow-coral: 0 10px 28px -12px var(--coral-glow);
  --shadow-glow-ring:  0 0 0 0 var(--coral-glow);

  /* ── Accent / Glow Variants (coral) ─────────────────────── */
  --color-accent-glow:      var(--coral-glow);
  --color-accent-dim:       var(--coral-dim);
  --color-accent-dim-hover: rgba(255, 75, 62, 0.18);
  --color-accent-surface:   rgba(255, 75, 62, 0.07);
  --color-accent-border:    var(--coral-border);
  --color-accent-shadow:    rgba(255, 75, 62, 0.20);

  /* ── Surface Variants ───────────────────────────────────── */
  --color-bg-card:       #0d0d0d;
  --color-bg-elevated:   #141414;
  --color-bg-input:      #111111;
  --color-bg-hover:      rgba(255,255,255,0.04);
  --color-bg-hover-mid:  rgba(255,255,255,0.07);
  --color-bg-submitted:  rgba(255,255,255,0.08);
  --color-shadow-deep:   rgba(0,0,0,0.80);
  --color-shadow-drop:   rgba(0,0,0,0.70);

  /* ── Status colors ──────────────────────────────────────── */
  --status-online:       #4ade80;
  --status-online-glow:  rgba(74,222,128,0.35);
  --status-locked:       rgba(255,255,255,0.30);

  /* ── Misc ───────────────────────────────────────────────── */
  --color-fg-white: #ffffff;
  --font-mono: 'JetBrains Mono', ui-monospace, SFMono-Regular, Menlo, monospace;

  /* ── Animation Timing ───────────────────────────────────── */
  --ease-quick:   0.125s ease-in-out;
  --ease-default: 0.2s ease;
  --ease-slow:    0.3s ease-in-out;
}
```

Also update the `* { font-family: ... }` rule in the Typography section to keep Comfortaa:

```css
* {
  font-family: 'Comfortaa', sans-serif;
}
```

No other changes.

## Build check

```bash
cd pkg/client/ArtelUI && bun run build
```

Expected: zero errors. The app should still render identically in structure, just with coral accent color instead of teal.
