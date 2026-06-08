# Handoff: Artel Notes Manager

## Overview
A full-screen, dark-mode notes management interface for the **artel** product — a connected-knowledge workspace. Inspired by Obsidian, it features a document tree sidebar, three distinct editing modes, and "wiki-link" chips that beautifully reference other notes.

This is a **greenfield feature** — no production notes UI exists yet. The design references here should be implemented as a new section/route within the existing artel codebase.

## About the Design Files
The `.html` files in this bundle are **high-fidelity design prototypes** built in React/Babel for presentation purposes. They are **not production code to ship directly**. The task is to recreate these designs in the target codebase using its established framework, router, and component patterns.

All design values (colors, spacing, typography) come from the **artel design system** (`artel-tokens.css`, `artel-shared.jsx`), which is already live in the project.

## Fidelity
**High-fidelity.** Colors, typography, spacing, border radii, shadows, and interactions are all specified to pixel precision. Implement as close to spec as possible; the design tokens are the single source of truth.

---

## Delivery Steps

The handoff is split into four layers — implement in this order:

| Step | File | Contents |
|------|------|----------|
| 1 | `01-tokens.md` | Design tokens: colors, type scale, spacing |
| 2 | `02-components.md` | Every atomic UI component with exact specs |
| 3 | `03-layouts.md` | Three page layouts and their structure |
| 4 | `04-interactions.md` | Modes, state, hover/click behaviors |

---

## Tech Constraints

- **Font stack:** `Comfortaa` (all UI + note body text), `JetBrains Mono` (code, metadata, labels, source-edit mode). Both already loaded via Google Fonts in the project.
- **Color mode:** Dark only. No light-mode variant required at this stage.
- **Scrollbars:** Hidden globally (`scrollbar-width: none` / `::-webkit-scrollbar { display: none }`).
- **Animations:** Keep the `artelPulse` glow on the brand mark. Mode transitions should be instant (no fade needed for v1).
- **Coral rule:** `#FF4B3E` appears **only** on the brand mark, active tree item border, filled checkboxes, wiki-link chip accents, and blockquote left borders. Nowhere else.

---

## Files in This Bundle

| File | Purpose |
|------|---------|
| `Artel Notes.html` | Full design canvas — pan/zoom to compare all three layouts |
| `notes-layouts.jsx` | Source for all UI components used in the mockups |
| `artel-tokens.css` | Design token source of truth |
| `artel-shared.jsx` | Shared brand components (ArtelMark, icons, etc.) |
| `01-tokens.md` | Token reference for implementation |
| `02-components.md` | Component specs |
| `03-layouts.md` | Layout specs |
| `04-interactions.md` | Interaction specs |

---

## Questions for the Developer
Before starting, confirm:
1. What routing library is in use? (The three layouts map to a single route `/notes/:noteId` with a `?mode=` param or equivalent.)
2. Is there an existing state management layer (Zustand, Redux, context) the notes store should plug into?
3. Is there an existing markdown parser dependency? If not, recommend **`marked`** + **`DOMPurify`** for the Preview/Read modes.
