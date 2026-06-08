# Step 1 — Design Tokens

All values are already defined in `artel-tokens.css`. This doc exists as a quick reference when building components — use the CSS variables in production, not raw hex values.

---

## Colors

### Surfaces
| Token | Value | Usage |
|-------|-------|-------|
| `--bg-base` | `#000000` | Editor area background |
| `--bg-card` | `#0d0d0d` | Sidebar, right panel, icon sidebar |
| `--bg-topbar` | `#0a0a0a` | Top bar, line-numbers column *(not in tokens file — add it)* |
| `--bg-elevated` | `#141414` | ModeBar container, dropdowns |
| `--bg-hover` | `rgba(255,255,255,0.04)` | Button hover, New button rest |
| `--bg-mid` | `rgba(255,255,255,0.07)` | Active icon button |

### Borders
| Token | Value | Usage |
|-------|-------|-------|
| `--border-dim` | `rgba(255,255,255,0.06)` | Most dividers and outlines |
| `--border-mid` | `rgba(255,255,255,0.11)` | ModeBar active tab inset, focused inputs |
| `--border-bright` | `rgba(255,255,255,0.22)` | Emphasis borders |

### Text
| Token | Value | Usage |
|-------|-------|-------|
| `--text-primary` | `rgba(255,255,255,0.92)` | Headings, active items, body |
| `--text-secondary` | `rgba(255,255,255,0.55)` | Body text, inactive tree items |
| `--text-muted` | `rgba(255,255,255,0.32)` | Metadata, placeholders, done tasks |
| `--text-disabled` | `rgba(255,255,255,0.18)` | Line numbers, section labels, icon strokes |

### Brand / Coral (≤5% of screen surface)
| Token | Value | Usage |
|-------|-------|-------|
| `--coral` | `#FF4B3E` | Brand mark, active tree border, filled checkbox, blockquote bar |
| `--coral-dim` | `rgba(255,75,62,0.08)` | WikiChip background, active tree item background |
| `--coral-border` | `rgba(255,75,62,0.22)` | WikiChip border |
| `--coral-glow` | `rgba(255,75,62,0.35)` | Brand mark pulse shadow |

> **Coral rule:** do not use coral for hover states, link colors, or decorative fills.  
> Reserve it strictly for the cases listed above.

---

## Typography

### Fonts
```css
--font:      'Comfortaa', sans-serif;
--font-mono: 'JetBrains Mono', ui-monospace, SFMono-Regular, Menlo, monospace;
```

### Note Content Type Scale (Preview / Read modes)
| Element | Font | Size | Weight | Line-height | Letter-spacing | Color |
|---------|------|------|--------|-------------|----------------|-------|
| H1 | Comfortaa | 26px | 700 | — | −0.02em | `--text-primary` |
| H2 | Comfortaa | 16px | 600 | — | −0.01em | `--text-primary` |
| H3 | Comfortaa | 14px | 600 | — | 0 | `--text-primary` |
| Body | Comfortaa | 14px | 400 | 1.72 | 0 | `--text-secondary` |
| Blockquote | Comfortaa | 13.5px | 400 italic | 1.72 | 0 | `--text-muted` |
| Inline code | JetBrains Mono | 11px | 400 | 1.65 | 0 | `--text-secondary` |

### Raw Edit Type Scale
| Element | Font | Size | Weight | Color |
|---------|------|------|--------|-------|
| `# H1` | JetBrains Mono | 15px | 700 | `--text-primary` |
| `## H2` | JetBrains Mono | 12px | 500 | `--text-primary` |
| Body | JetBrains Mono | 12px | 400 | `--text-secondary` |
| Muted (blockquote, done tasks) | JetBrains Mono | 12px | 400 | `--text-muted` |
| `[[RawLink]]` | JetBrains Mono | 11px | 400 | `rgba(255,160,130,0.9)` |

### UI Chrome Type Scale
| Element | Font | Size | Weight | Letter-spacing | Color |
|---------|------|------|--------|----------------|-------|
| App name "artel" | Comfortaa | 14px | 700 | −0.02em | `--text-primary` |
| Section label "FAVORITES" | JetBrains Mono | 9px | 400 | 0.12em | `--text-disabled` |
| Tree item | Comfortaa | 12px | 400 | 0 | `--text-secondary` |
| Tag pill | JetBrains Mono | 10px | 400 | 0.06em | `--text-muted` |
| Breadcrumb | JetBrains Mono | 11px | 400 | 0 | `--text-disabled` |
| Mode tab | JetBrains Mono | 10px | 400 | 0.06em | varies |
| Metadata row | JetBrains Mono | 10px | 400 | 0 | `--text-disabled` |
| Line numbers | JetBrains Mono | 10px | 400 | 0 | `--text-disabled` |
| Note title (TitleBar) | Comfortaa | 17px | 600 | −0.01em | `--text-primary` |

---

## Spacing & Layout

### Key Dimensions
| Name | Value | Where |
|------|-------|-------|
| Top bar height | 44px | All layouts |
| Title bar height | 46px | Layout A only |
| Full sidebar width | 240px | Layout A |
| Slim sidebar width | 220px | Layout C |
| Icon sidebar width | 44px | Layout B |
| Line numbers width | 38px | Layout B |
| Right panel width | 200px | Layout C |
| Editor content padding | 30px 40px | Preview/Read |
| Raw content padding | 28px 32px | Raw Edit |
| Sidebar content padding (horizontal) | 16px | All sidebar items |
| Tree indent per depth level | 12px | Tree items |

### Sidebar Item Heights
| Component | Height |
|-----------|--------|
| Search bar | ~34px (7px padding × 2 + ~20px content) |
| Section label | ~28px |
| Tree item | ~26px (4px padding × 2 + ~18px line) |
| Divider | 1px + 5px margin each side |
| Tag pill row | auto (flex-wrap) |

---

## Shadows & Effects

| Name | Value | Usage |
|------|-------|-------|
| Brand mark pulse | `drop-shadow(0 0 1.1rem rgba(255,75,62,0))` → `drop-shadow(0 0 0 rgba(255,75,62,0.7))` | ArtelMark animation |
| Code block | none (border only) | Code blocks |
| WikiChip hover | background → `rgba(255,75,62,0.14)` | On hover |

---

## Border Radii
| Element | Radius |
|---------|--------|
| WikiChip | 5px |
| RawLink | 3px |
| ModeBar container | 999px (pill) |
| ModeBar tab | 999px |
| Search bar | 6px |
| New button | 6px |
| Tree item active state | 0 (full-width) |
| Tag pill | 4px |
| Code block | 6px |
| Checkbox | 3px |
| Icon button | 6px |
