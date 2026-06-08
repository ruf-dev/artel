# Step 2 — Components

Build these as isolated, reusable components before assembling the layouts. Listed smallest → largest.

---

## 1. ArtelMark (Brand Icon)

The coral disc with a geometric "А" cut-out. Used as the app logo and as the icon inside WikiChips.

```
Shape: SVG circle (cx=50 cy=50 r=50), filled #FF4B3E
Cut-out path: M50 22L28 78h10l5.5-14h13L62 78h10ZM46.5 56L50 38l3.5 18Z
Technique: SVG <mask> — white rect minus black path reveals the А
Animation (logo only): artelPulse — see tokens
```

**Sizes used:**
| Context | Size |
|---------|------|
| Top bar (all layouts) | 22×22px |
| Icon sidebar top | 22×22px |
| WikiChip inline | 11×11px |

**Note:** every SVG mask needs a globally unique `id`. Generate with `useId()` or a stable counter seeded on mount.

---

## 2. WikiChip

Inline reference to another note. Appears inside rendered note body text and task lists.

```
Display:        inline-flex
Align:          center
Gap:            5px
Background:     rgba(255,75,62,0.08)   [--coral-dim]
Border:         1px solid rgba(255,75,62,0.22)   [--coral-border]
Border-radius:  5px
Padding:        1px 8px 1px 5px
Font:           Comfortaa 12px
Color:          rgba(255,255,255,0.92)  [--text-primary]
Cursor:         pointer
White-space:    nowrap
Vertical-align: middle
Line-height:    1.6
User-select:    none
```

**Left icon:** `<ArtelMark size={11} />` — no pulse animation when used inline.

**Hover state:**
```
Background: rgba(255,75,62,0.14)
Transition: background 150ms ease
```

**Props:** `name: string` — the linked note's display title.

**On click:** navigate to that note (or open in split pane if supported).

---

## 3. RawLink

How wiki-links appear in **Raw Edit** mode — shows the raw `[[brackets]]` with a coral tint, no ArtelMark icon.

```
Display:        inline
Background:     rgba(255,75,62,0.11)
Border:         1px solid rgba(255,75,62,0.25)
Border-radius:  3px
Padding:        0 4px
Color:          rgba(255,160,130,0.9)
Font:           JetBrains Mono 11px
```

Rendered text: `[[{name}]]` — include the brackets literally.

---

## 4. ModeBar

The three-mode switcher: Edit · Preview · Read. Used in TitleBar (Layout A) and TopBar (Layouts B & C).

**Container:**
```
Display:        inline-flex
Padding:        3px
Background:     #141414   [--bg-elevated]
Border:         1px solid rgba(255,255,255,0.06)
Border-radius:  999px
Font:           JetBrains Mono 10px
Letter-spacing: 0.06em
Gap:            0
```

**Tab (inactive):**
```
Padding:        5px 13px
Border-radius:  999px
Background:     transparent
Color:          rgba(255,255,255,0.32)
Cursor:         pointer
```

**Tab (active):**
```
Background:     #0d0d0d   [--bg-card]
Color:          rgba(255,255,255,0.92)
Box-shadow:     inset 0 0 0 1px rgba(255,255,255,0.11)
```

**Props:** `active: 'edit' | 'preview' | 'read'`

**Behavior:** clicking a tab switches the editor mode (see `04-interactions.md`).

---

## 5. SearchBar

Text input that filters the note tree. Sits at the top of the sidebar.

```
Margin:         8px 10px
Display:        flex
Align:          center
Gap:            8px
Background:     rgba(255,255,255,0.035)
Border:         1px solid rgba(255,255,255,0.06)
Border-radius:  6px
Padding:        7px 10px
```

**Icon:** magnifier SVG, 12×12px, `stroke: rgba(255,255,255,0.18)`, `stroke-width: 1.6`.

**Placeholder text:**
```
Font:           JetBrains Mono 11px
Letter-spacing: 0.04em
Color:          rgba(255,255,255,0.18)
Content:        "Search notes…"
```

**Focus state:**
```
Border-color: rgba(255,255,255,0.22)
Background:   rgba(255,255,255,0.055)
```

---

## 6. SidebarLabel (Section Header)

Uppercase monospace label that groups sidebar sections.

```
Padding:        12px 16px 5px
Font:           JetBrains Mono 9px
Letter-spacing: 0.12em
Text-transform: uppercase
Color:          rgba(255,255,255,0.18)
```

Text content examples: `★ FAVORITES`, `⏱ RECENT`, `ALL NOTES`, `# TAGS`, `◉ GRAPH`

---

## 7. TreeItem

Single row in the file/folder tree. Supports files, folders, and active state.

```
Display:        flex
Align:          center
Gap:            5px
Padding:        4px 16px 4px [16 + depth * 12]px
Font:           Comfortaa 12px
Cursor:         pointer
```

**Inactive state:**
```
Color:          rgba(255,255,255,0.55)
Background:     transparent
Border-left:    2px solid transparent
```

**Active state:**
```
Color:          rgba(255,255,255,0.92)
Background:     rgba(255,75,62,0.08)
Border-left:    2px solid #FF4B3E
```

**Hover state (inactive only):**
```
Background:     rgba(255,255,255,0.035)
```

**Folder expand arrow:**
```
Size:           8×8px
Fill:           currentColor
Opacity:        0.35
Transform:      rotate(0deg) collapsed / rotate(90deg) expanded
Transition:     transform 150ms ease
```

**File icon:** 10×10px SVG (document outline), `opacity: 0.35`, `stroke: currentColor`.  
**Folder icon:** 10×10px SVG (folder outline), `opacity: 0.45`, `stroke: currentColor`.

**Props:**
```ts
interface TreeItemProps {
  name: string
  active?: boolean
  depth?: number        // 0 = root, 1 = nested, etc.
  isFolder?: boolean
  isOpen?: boolean      // only for folders
  onClick?: () => void
  onToggle?: () => void // folders: expand/collapse
}
```

---

## 8. TagPill

Compact tag chip shown in the Tags section of the sidebar.

```
Display:        inline-block
Padding:        3px 8px
Background:     rgba(255,255,255,0.04)
Border:         1px solid rgba(255,255,255,0.06)
Border-radius:  4px
Font:           JetBrains Mono 10px
Letter-spacing: 0.06em
Color:          rgba(255,255,255,0.32)
Cursor:         pointer
```

**Hover:**
```
Background:     rgba(255,255,255,0.07)
Color:          rgba(255,255,255,0.55)
```

**Props:** `tag: string` — include the `#` prefix in the content.

---

## 9. MiniGraph

Decorative knowledge-graph preview at the bottom of the full sidebar. SVG only — not interactive in v1.

```
Container:  padding 2px 10px 10px, width 100%
SVG:        viewBox="0 0 236 120", height 76px, width 100%
```

**Nodes:** circles
- Active node (current note): `r=5`, `fill=#FF4B3E`, `filter: drop-shadow(0 0 4px #FF4B3E)`
- Other nodes: `r=3`, `fill=rgba(255,255,255,0.22)`

**Edges:** `<line>`, `stroke=rgba(255,255,255,0.08)`, `strokeWidth=1.5`

In v1, render a hardcoded graph of 6 nodes representing the current note's connections. In v2, generate dynamically from actual link data.

---

## 10. TitleBar

Note-level header shown in **Layout A** (Classic). Contains title, metadata, and ModeBar.

```
Height:         46px
Background:     #000000
Border-bottom:  1px solid rgba(255,255,255,0.06)
Padding:        0 32px
Display:        flex
Align:          center
Gap:            16px
```

**Left — Note title:**
```
Font:     Comfortaa 17px, weight 600, letter-spacing -0.01em
Color:    rgba(255,255,255,0.92)
```

**Left — Metadata (below title, same block):**
```
Font:     JetBrains Mono 10px
Color:    rgba(255,255,255,0.18)
Margin-top: 1px
Content:  "edited {timeAgo} · {wordCount} words · {linkCount} links"
```

**Right:** `<ModeBar active={mode} />`

---

## 11. TopBar

Full-width app chrome bar. Present in all three layouts (different content per layout).

```
Height:         44px
Background:     #0a0a0a
Border-bottom:  1px solid rgba(255,255,255,0.06)
Padding:        0 16px
Display:        flex
Align:          center
Gap:            10px
```

**Left — Brand group:**
```
<ArtelMark size={22} />
<span> artel </span>   → Comfortaa 14px, weight 700, letter-spacing -0.02em
<span> notes </span>   → JetBrains Mono 10px, letter-spacing 0.06em, color --text-disabled
                         padding-left 8px, border-left 1px solid --border-dim
```

**Center (Layout C only) — Breadcrumb:**
```
Font:     JetBrains Mono 11px
Color:    --text-disabled  (ancestors) / --text-primary (current)
Separator: "  ›  " in --text-disabled
```

**Right — New button:**
```
Display:        inline-flex
Align:          center
Gap:            5px
Padding:        5px 10px
Background:     rgba(255,255,255,0.04)
Border:         1px solid rgba(255,255,255,0.06)
Border-radius:  6px
Font:           JetBrains Mono 11px
Color:          rgba(255,255,255,0.32)
Icon:           plus SVG 10×10px, stroke currentColor, strokeWidth 1.8
```

**Right (optional) — ModeBar:** present in Layout B (in TopBar) and Layout C (in TopBar). In Layout A the ModeBar is inside TitleBar instead.

---

## 12. IconSidebar

Compact vertical icon rail used in **Layout B**. Width: 44px.

```
Width:          44px
Background:     #0d0d0d
Border-right:   1px solid rgba(255,255,255,0.06)
Display:        flex
Flex-direction: column
Align:          center
Padding-top:    10px
Gap:            4px
```

**Top:** `<ArtelMark size={22} />` with pulse animation, `margin-bottom: 10px`.

**Icon button (each):**
```
Width:          32px
Height:         32px
Border-radius:  6px
Display:        flex
Align:          center
Justify:        center
Cursor:         pointer
```

**Active:**
```
Background:     rgba(255,255,255,0.07)
Color:          rgba(255,255,255,0.92)
```

**Inactive:**
```
Background:     transparent
Color:          rgba(255,255,255,0.32)
```

**Hover (inactive):**
```
Background:     rgba(255,255,255,0.04)
```

**Icons (SVG outlines, 15–16px, strokeWidth 1.5–1.6):**
1. Search (magnifier)
2. Files (tree / file list)  ← default active
3. Favorites (star)
4. Tags (tag/label)
5. Graph (5-node graph)

**Tooltip:** show label on hover (`title` attribute or custom tooltip) — e.g. "Files", "Favorites".

---

## 13. LineNumbers

Narrow column shown to the left of the raw editor in **Layout B**.

```
Width:          38px
Background:     #080808
Border-right:   1px solid rgba(255,255,255,0.06)
Padding:        28px 0   (matches editor content padding-top)
Flex-shrink:    0
```

**Each line:**
```
Height:         21px   (matches JetBrains Mono 12px line-height 1.75)
Display:        flex
Align:          center
Justify:        flex-end
Padding-right:  10px
Font:           JetBrains Mono 10px
Color:          rgba(255,255,255,0.18)
```

Count: render enough lines to fill the visible area (lazy or virtualized in v2).

---

## 14. RightPanel

Collapsible panel in **Layout C** (Read mode). Shows Outline, Backlinks, and Properties.

```
Width:          200px
Background:     #0d0d0d
Border-left:    1px solid rgba(255,255,255,0.06)
Flex-shrink:    0
Display:        flex
Flex-direction: column
Overflow:       hidden
```

**Outline section:**
Heading items rendered as:
```
Padding:        4px 16px
Font:           Comfortaa 11px
Color:          rgba(255,255,255,0.32)
Cursor:         pointer
Content prefix: "— "
```
Scroll-to on click (smooth scroll to heading anchor in v2).

**Backlinks section:**
Each backlink:
```
Padding:        6px 16px
```
Note name: `Comfortaa 11px`, `--text-secondary`  
Context hint: `JetBrains Mono 9px`, `--text-disabled`, `margin-top: 1px`

**Properties section:**
Key–value pairs:
```
Display:        flex
Justify:        space-between
Padding:        3px 16px
Font:           JetBrains Mono 10px
Key color:      --text-disabled
Value color:    --text-muted
```

Fields: `created`, `modified`, `words`, `links`.
