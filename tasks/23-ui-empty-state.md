---
status: done
---

# Task 23 — UI: Empty State Component

## Goal

Replace the bare `<p>No vaults yet…</p>` text in `HomePage` with a polished illustrated empty
state: SVG art of a vault icon, title, subtitle, and a "Create your first vault" CTA button.

## Prerequisite

Tasks 20 and 21 must be done first.

## New file: `pkg/client/ArtelUI/src/components/EmptyState/EmptyState.tsx`

```tsx
import cls from "./EmptyState.module.css"

interface Props {
    onCreateClick: () => void
}

export default function EmptyState({onCreateClick}: Props) {
    return (
        <div className={cls.Root}>
            <div className={cls.Art} aria-hidden="true">
                <svg viewBox="0 0 100 100" fill="none" stroke="rgba(255,255,255,0.35)" strokeWidth="1.4" strokeLinecap="round" strokeLinejoin="round">
                    <rect x="18" y="28" width="64" height="52" rx="8"/>
                    <path d="M30 28v-6a8 8 0 0 1 8-8h24a8 8 0 0 1 8 8v6"/>
                    <circle cx="50" cy="54" r="6" fill="rgba(255,75,62,0.18)" stroke="#FF4B3E"/>
                    <line x1="50" y1="60" x2="50" y2="68"/>
                </svg>
            </div>
            <h2 className={cls.Title}>No vaults yet</h2>
            <p className={cls.Sub}>
                Vaults are encrypted CouchDB instances for your data. Create your first one
                and you'll get a connection string in seconds.
            </p>
            <button className={cls.Cta} type="button" onClick={onCreateClick}>
                <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                    <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
                </svg>
                Create your first vault
            </button>
        </div>
    )
}
```

## New file: `pkg/client/ArtelUI/src/components/EmptyState/EmptyState.module.css`

```css
.Root {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    text-align: center;
    padding: 80px 24px 60px;
    animation: fadeUp 500ms cubic-bezier(.2,.7,.2,1) both;
}

@keyframes fadeUp {
    from { opacity: 0; transform: translateY(8px); }
    to   { opacity: 1; transform: translateY(0); }
}

.Art {
    width: 96px;
    height: 96px;
    margin-bottom: 32px;
}

.Art svg {
    width: 100%;
    height: 100%;
    opacity: 0.85;
}

.Title {
    font-size: 24px;
    font-weight: 600;
    letter-spacing: -0.02em;
    color: var(--color-fg-primary);
    margin-bottom: 12px;
}

.Sub {
    font-size: 14px;
    color: var(--color-fg-secondary);
    line-height: 1.6;
    max-width: 420px;
    margin: 0 auto 28px;
}

.Cta {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    padding: 12px 18px;
    border-radius: 12px;
    border: 1px solid var(--coral);
    background: var(--coral);
    color: white;
    font-size: 14px;
    font-weight: 500;
    font-family: inherit;
    cursor: pointer;
    box-shadow: 0 10px 28px -12px var(--coral-glow);
    transition: filter 160ms ease, box-shadow 160ms ease;
}

.Cta:hover {
    filter: brightness(1.06);
    box-shadow: 0 14px 32px -10px var(--coral-glow);
}

.Cta:active {
    transform: translateY(1px);
}
```

## Modify: `pkg/client/ArtelUI/src/pages/home/HomePage.tsx`

1. Add import:
   ```tsx
   import EmptyState from "@/components/EmptyState/EmptyState.tsx"
   ```

2. Replace the empty branch in the render:

   Find this block:
   ```tsx
   ) : vaults.length === 0 ? (
       <p className={cls.Empty}>No vaults yet. Create one with +</p>
   ) : (
   ```

   Replace with:
   ```tsx
   ) : vaults.length === 0 ? (
       <EmptyState onCreateClick={openDialog}/>
   ) : (
   ```

3. Remove the `.Empty` class usage — it's no longer needed (the CSS class can stay or be removed from the CSS file).

## Build check

```bash
cd pkg/client/ArtelUI && bun run build
```

Expected: zero errors.
