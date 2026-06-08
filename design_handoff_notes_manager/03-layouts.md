# Step 3 — Layouts (Pages)

Three layouts share the same route and note data — the user switches between them via ModeBar or a global preference. All are full-viewport (min 100vw × 100vh), dark background `#000`.

---

## Shared Structure

```
┌─────────────────────────────────┐
│           TopBar (44px)         │  ← present in A and C; in B it's the breadcrumb bar
├──────────┬──────────────────────┤
│          │                      │
│ Sidebar  │    Editor Area       │  ← varies per layout
│          │                      │
└──────────┴──────────────────────┘
```

All three layouts use `display: flex; flex-direction: column; height: 100vh; overflow: hidden`.

---

## Layout A — Classic (Live Preview / WYSIWYG)

**When to show:** default mode, or when user selects "Preview" in ModeBar.  
**Reference artboard:** "A · Classic — Live Preview" in `Artel Notes.html`

### Column structure
```
[TopBar — full width, 44px]
[NoteSidebar 240px] | [divider 1px] | [Editor flex:1]
```

### TopBar (Layout A)
- Brand group (left): ArtelMark 22px + "artel" + "notes"
- Right: New button only (ModeBar is inside TitleBar below)

### Sidebar (240px)
Full sidebar — see component specs. Sections top to bottom:
1. SearchBar
2. ★ FAVORITES section + items
3. Divider
4. ⏱ RECENT section + items
5. Divider
6. ALL NOTES section + tree
7. Divider
8. \# TAGS section + tag pills (flex-wrap)
9. Divider
10. ◉ GRAPH section + MiniGraph SVG

The sidebar scrolls independently (`overflow-y: auto`). Content below tags is shown only in full sidebar (not slim variant).

### Editor Area (flex: 1)
```
[TitleBar — 46px, note title + metadata + ModeBar]
[NoteContent — flex: 1, overflow-y: auto, padding 30px 40px]
```

### NoteContent — Preview mode rendering
Render parsed markdown as styled HTML. See `01-tokens.md` for type scale.

**Block elements (top → bottom):**
| Markdown | Rendered element | Key styles |
|----------|-----------------|------------|
| `# H1` | `<h1>` | 26px/700/Comfortaa, −0.02em, mb 14px |
| `## H2` | `<h2>` | 16px/600/Comfortaa, −0.01em, mt 26px mb 10px |
| `> blockquote` | `<blockquote>` | border-left 2px solid rgba(255,75,62,0.4), pl 14px, italic, color --text-muted |
| paragraph | `<p>` | 14px/400/Comfortaa, lh 1.72, mb 14px, color --text-secondary |
| `- [ ] task` | div+checkbox+span | 14px checkbox 14×14px, br 3px, open: border --border-mid |
| `- [x] task` | same + strikethrough | checkbox filled #FF4B3E with ✓, text color --text-disabled, text-decoration line-through |
| `` ` `` code | `<code>` | inline: JetBrains Mono 11px, bg #0d0d0d, px 3px, br 3px |
| ` ``` ` block | `<pre><code>` | bg #0d0d0d, border --border-dim, br 6px, p 12px 16px, JetBrains Mono 11px, lh 1.65 |
| `[[Link]]` | `<WikiChip name="Link" />` | replace all `[[…]]` with WikiChip before rendering |

---

## Layout B — Focused (Raw / Source Edit)

**When to show:** when user selects "Edit" (raw source) in ModeBar.  
**Reference artboard:** "B · Focused — Raw Edit" in `Artel Notes.html`

### Column structure
```
[BreadcrumbBar — full width, 44px]
[IconSidebar 44px] | [LineNumbers 38px] | [RawEditor flex:1]
```

No traditional sidebar — the icon sidebar is the only navigation chrome.

### BreadcrumbBar (44px, replaces TopBar)
```
Background:     #0a0a0a
Border-bottom:  1px solid rgba(255,255,255,0.06)
Padding:        0 20px
Display:        flex
Align:          center
Gap:            6px
```

**Left — breadcrumb path:**
Each segment: `JetBrains Mono 11px`  
Active (last segment): `--text-primary`  
Ancestors: `--text-disabled`  
Separator: `›` in `--text-disabled`

**Right:** `<ModeBar active="edit" />`

### IconSidebar (44px)
See component spec. Active section: Files (tree icon).

### LineNumbers (38px)
See component spec. Synced with editor scroll position.

### RawEditor (flex: 1)
Monospace text area, dark background `#000`.  
`padding: 28px 32px`  
`font: JetBrains Mono 12px, line-height 1.75`

**Syntax display rules** (highlight without a full lexer — pattern-match lines):
| Pattern | Style |
|---------|-------|
| Lines starting with `# ` | `--text-primary`, weight 700, font-size 15px |
| Lines starting with `## ` | `--text-primary`, weight 500, font-size 12px |
| Lines starting with `> ` | `--text-muted`, font-style italic |
| Lines starting with `- [x]` | `--text-muted`, text-decoration line-through |
| Lines starting with `- [ ]` or `- ` | `--text-secondary` |
| Lines starting with ` ``` ` | `--text-muted` |
| Lines inside code fence | `--text-muted`, bg `rgba(255,255,255,0.02)`, pl 16px |
| Empty lines | `--text-disabled` (space placeholder) |

**`[[wiki-link]]` in raw mode:** replace with `<RawLink name="…" />` inline. The surrounding `[[` / `]]` brackets appear literally.

**Cursor blink:**  
Insert a blinking cursor span at the active cursor position:
```css
width: 2px; height: 1em; background: #FF4B3E;
animation: blinkCaret 1s step-end infinite;
display: inline-block; vertical-align: text-bottom; margin-left: 2px;
```
(`@keyframes blinkCaret` is defined in `artel-tokens.css`)

In v1, the cursor can be static (placed at the end of the last edited line).

---

## Layout C — Panel (Read Mode + Right Panel)

**When to show:** when user selects "Read" in ModeBar.  
**Reference artboard:** "C · Panel — Read + Backlinks" in `Artel Notes.html`

### Column structure
```
[TopBar — full width, 44px, includes breadcrumb + ModeBar]
[SlimSidebar 220px] | [divider] | [ReadContent flex:1] | [divider] | [RightPanel 200px]
```

### TopBar (Layout C)
```
Left:    ArtelMark + "artel" + "notes"
Center:  Breadcrumb — "Projects  ›  {noteTitle}"
Right:   ModeBar active="read" | New button
```

### SlimSidebar (220px)
Same as full sidebar but:
- Width: 220px (20px narrower)
- Graph section **omitted** (not enough room)
- Otherwise identical sections

### ReadContent (flex: 1)
No TitleBar. Note content starts directly at top of content area.  
`padding: 30px 40px`  
Rendered markdown exactly like Layout A's Preview mode.  
No editing affordances — all text `user-select: text` but no `contentEditable`.

**Reading line width:** optionally cap at `max-width: 720px; margin: 0 auto` for comfortable reading. This is a v2 decision — v1 can let it fill the space.

### RightPanel (200px)
See component spec. Sections: Outline → Backlinks → Properties.

---

## Responsive Behavior (v2, not required in v1)

| Viewport | Behavior |
|----------|----------|
| < 900px | Sidebar collapses to icon-only (44px) in all layouts |
| < 640px | Sidebar hidden by default; toggle via hamburger in TopBar |
| Layout C < 1100px | Right panel collapses to a toggle button |

---

## Route / URL Design (recommendation)

```
/notes                       → redirect to last opened note or inbox
/notes/:noteId               → Layout A (Preview) by default
/notes/:noteId?mode=edit     → Layout B (Raw Edit)
/notes/:noteId?mode=read     → Layout C (Read)
```

Persist last mode to `localStorage` per note or globally.
