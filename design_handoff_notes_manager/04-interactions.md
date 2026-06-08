# Step 4 — Interactions & Behavior

---

## Mode Switching

The three modes share one route and one note. Switching mode is a **URL param change** (or equivalent state update) — the note data never reloads.

### Mode → Layout mapping
| Mode tab | Layout rendered | URL param |
|----------|----------------|-----------|
| Preview | Layout A (Classic) | `?mode=preview` or default |
| Edit | Layout B (Focused, raw source) | `?mode=edit` |
| Read | Layout C (Panel + right rail) | `?mode=read` |

### Transition
No fade or slide in v1 — instant swap. The editor scroll position should be preserved across mode switches (store in a ref, restore on mount).

### Persistence
- Store the last used mode per note in `localStorage` keyed by `artel.notes.mode.{noteId}`.
- On first visit to a note, default to `preview`.

---

## Note Tree

### Expand / Collapse Folders
- Click on a folder row (anywhere on the row) → toggle `isOpen` state.
- Arrow icon rotates 0° → 90° on open (`transform: rotate(90deg)`, `transition: 150ms ease`).
- Children animate in/out with a simple max-height collapse (`max-height: 0 → max-height: 500px`, `overflow: hidden`, `transition: 200ms ease`).
- Expand state persists in `localStorage` keyed by `artel.notes.tree.{folderId}`.

### Active Note Highlight
- The currently open note gets `active` state: coral left border + coral-dim background.
- If the active note is inside a collapsed folder, auto-expand that folder on load.

### Hover State
- Any inactive tree item on hover: `background: rgba(255,255,255,0.035)`.
- Transition: `background 120ms ease`.

### Right-click context menu (v2)
Not required for v1. When implemented: Rename, Delete, Move to, New note here.

---

## WikiChip Interactions

### Hover
```
background: rgba(255,75,62,0.14)   (was 0.08)
transition: background 150ms ease
cursor: pointer
```

### Click
Navigate to the linked note: `router.push('/notes/' + noteSlug)`.  
In v2: open in right pane / split view.

### Tooltip (v2)
Show a preview popover on hover-hold (300ms delay): note title + first 2 lines of content.

---

## Search

### Trigger
- Clicking SearchBar focuses the input.
- Keyboard shortcut: `Cmd/Ctrl + K` opens search (or `Cmd/Ctrl + P` — pick one, be consistent with rest of app).
- `Escape` clears and blurs.

### Filter behavior
- Filter tree in real-time as user types.
- Match against note titles (case-insensitive substring).
- Hide non-matching items; keep parent folders visible if any child matches.
- Highlight matching substring in tree item text (v2).
- Empty state: "No notes match" message in `--text-muted`, `JetBrains Mono 11px`.

---

## New Note Button

1. Creates a new note with an auto-generated title ("Untitled" + timestamp or sequential number).
2. Adds it to the tree under the currently active folder (or root if none).
3. Navigates to the new note in `edit` mode so the user can immediately type a title.
4. The title in TitleBar becomes `contentEditable` on creation — auto-focus + select-all.

---

## Tags

### Sidebar tag pills
- Click: filter the tree to show only notes tagged with that tag.
- Active tag: background `rgba(255,255,255,0.08)`, color `--text-secondary`.
- Click again: deselect (show all notes).
- Multiple tags: AND filter (only notes matching all selected tags).

---

## Outline (Right Panel, Layout C)

- Each item is a heading (H1/H2/H3) extracted from the note.
- Click: smooth-scroll the ReadContent area to that heading anchor.
- Highlight the item corresponding to the heading currently in the viewport (intersection observer, v2).

---

## Backlinks (Right Panel, Layout C)

- List notes that contain `[[Current Note Title]]`.
- Click: navigate to that note (preserving current mode).
- Context hint: show the sentence containing the link (truncated to 1 line, v2).

---

## Editor Behavior (Raw Edit, Layout B)

### Typing
In v1, a standard `<textarea>` or `contentEditable` div is acceptable.  
In v2, consider CodeMirror or Monaco for:
- Syntax highlighting
- Auto-closing brackets for `[[`
- `[[` → autocomplete from note titles

### `[[` auto-complete (v2)
When user types `[[`, show a dropdown of note titles filtered by subsequent typing.  
Each option: note title + tag list in muted text.  
Selecting inserts `[[Note Title]]` and closes dropdown.

### Saving
- **Auto-save** on every keystroke with a 500ms debounce.
- Show save state in metadata row: "saving…" → "saved" → "edited Xs ago".
- On network error: "save failed" in `--color-error` (#c1253d).

### Keyboard shortcuts
| Shortcut | Action |
|----------|--------|
| `Cmd/Ctrl + S` | Force-save immediately |
| `Cmd/Ctrl + B` | Wrap selection in `**bold**` |
| `Cmd/Ctrl + I` | Wrap selection in `*italic*` |
| `Cmd/Ctrl + K` | Insert `[[` → triggers link autocomplete |
| `Tab` | Indent (add 2 spaces or increase list depth) |
| `Shift+Tab` | Dedent |

---

## ArtelMark Animation

The brand mark in the **TopBar** and **IconSidebar** uses the `artelPulse` animation:

```css
@keyframes artelPulse {
  0%   { filter: drop-shadow(0 0 0 rgba(255,75,62,0.7)); }
  60%  { filter: drop-shadow(0 0 1.1rem rgba(255,75,62,0)); }
  100% { filter: drop-shadow(0 0 0 rgba(255,75,62,0.7)); }
}
animation: artelPulse 2.6s ease-in-out infinite;
```

WikiChip ArtelMark icons do **not** animate — `animation: none`.

---

## Icon Sidebar Panels (Layout B, v2)

In v1 the icon sidebar is navigation-only (clicking a section icon switches to Layout A with that panel open).

In v2, clicking an icon in Layout B expands a **temporary overlay panel** (240px wide) that slides in over the editor, auto-closes when the editor is clicked:

```
Transition: transform 200ms cubic-bezier(0.2, 0.7, 0.3, 1)
Transform:  translateX(-100%) → translateX(0)
```

---

## Empty States

| Situation | Message | Placement |
|-----------|---------|-----------|
| No notes in vault | "Your vault is empty. Create your first note." | Center of editor area |
| No search results | "No notes match "{query}"" | Below search bar, replacing tree |
| No backlinks | "No notes link here yet." | In backlinks section of RightPanel |
| No tags | *(hide # TAGS section entirely)* | — |
