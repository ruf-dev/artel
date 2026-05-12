---
status: todo
---

# Task 21 — UI: Topbar + Hero Section

## Goal

Replace the bare header in `HomePage` with a polished sticky topbar (animated brand mark, user
dropdown with logout) and add a hero section below it (eyebrow, title, vault count, "New vault"
inline button). Remove the floating action button (FAB) — the "New vault" button lives in the hero now.

## Prerequisite

Task 20 must be done first (coral tokens + fonts).

## New file: `pkg/client/ArtelUI/src/components/Topbar/Topbar.tsx`

```tsx
import {useState} from "react"
import cls from "./Topbar.module.css"

interface Props {
    onLogout: () => void
}

export default function Topbar({onLogout}: Props) {
    const [menuOpen, setMenuOpen] = useState(false)

    return (
        <header className={cls.Topbar}>
            <a className={cls.Brand} href="/" aria-label="Artel home">
                <svg className={cls.BrandMark} viewBox="0 0 100 100" aria-hidden="true">
                    <defs>
                        <mask id="brand-a-cut">
                            <rect width="100" height="100" fill="white"/>
                            <path d="M 50 22 L 28 78 L 38 78 L 43.5 64 L 56.5 64 L 62 78 L 72 78 Z M 46.5 56 L 50 38 L 53.5 56 Z" fill="black"/>
                        </mask>
                    </defs>
                    <circle cx="50" cy="50" r="50" fill="#FF4B3E" mask="url(#brand-a-cut)"/>
                </svg>
                <span className={cls.BrandWord}>artel</span>
            </a>

            <div className={cls.Right}>
                <div className={cls.UserWrap}>
                    <button
                        className={cls.UserPill}
                        type="button"
                        aria-expanded={menuOpen}
                        aria-haspopup="menu"
                        onClick={e => { e.stopPropagation(); setMenuOpen(v => !v) }}
                    >
                        <span className={cls.Avatar}>
                            <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="white" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                                <circle cx="12" cy="7" r="4"/>
                            </svg>
                        </span>
                        <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"
                            style={{transform: menuOpen ? 'rotate(180deg)' : 'none', transition: 'transform 200ms ease', color: 'rgba(255,255,255,0.4)'}}>
                            <polyline points="6 9 12 15 18 9"/>
                        </svg>
                    </button>

                    {menuOpen && (
                        <>
                            <div className={cls.Backdrop} onClick={() => setMenuOpen(false)}/>
                            <div className={cls.Menu} role="menu">
                                <button
                                    className={`${cls.MenuItem} ${cls.MenuItemDanger}`}
                                    role="menuitem"
                                    type="button"
                                    onClick={() => { setMenuOpen(false); onLogout() }}
                                >
                                    <svg viewBox="0 0 24 24" width="15" height="15" fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
                                        <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
                                        <polyline points="16 17 21 12 16 7"/>
                                        <line x1="21" y1="12" x2="9" y2="12"/>
                                    </svg>
                                    <span>Log out</span>
                                </button>
                            </div>
                        </>
                    )}
                </div>
            </div>
        </header>
    )
}
```

## New file: `pkg/client/ArtelUI/src/components/Topbar/Topbar.module.css`

```css
.Topbar {
    position: sticky;
    top: 0;
    z-index: 50;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 16px;
    padding: 18px 32px;
    background: rgba(0,0,0,0.7);
    backdrop-filter: saturate(140%) blur(16px);
    -webkit-backdrop-filter: saturate(140%) blur(16px);
    border-bottom: 1px solid var(--color-border);
    flex-shrink: 0;
}

@media (max-width: 720px) {
    .Topbar { padding: 14px 18px; }
}

.Brand {
    display: flex;
    align-items: center;
    gap: 10px;
    text-decoration: none;
    color: inherit;
}

@keyframes artelPulse {
    0%   { filter: drop-shadow(0 0 0    rgba(255,75,62,0.7)); }
    60%  { filter: drop-shadow(0 0 1.2rem rgba(255,75,62,0));  }
    100% { filter: drop-shadow(0 0 0    rgba(255,75,62,0.7)); }
}

.BrandMark {
    width: 24px;
    height: 24px;
    animation: artelPulse 2.6s ease-in-out infinite;
}

.BrandWord {
    font-weight: 700;
    font-size: 16px;
    letter-spacing: -0.02em;
    color: var(--color-fg-primary);
}

.Right {
    display: flex;
    align-items: center;
    gap: 12px;
}

.UserWrap {
    position: relative;
}

.UserPill {
    display: inline-flex;
    align-items: center;
    gap: 10px;
    padding: 5px 12px 5px 5px;
    background: var(--color-bg-card);
    border: 1px solid var(--color-border);
    border-radius: 999px;
    cursor: pointer;
    color: var(--color-fg-primary);
    transition: border-color var(--ease-default), background var(--ease-default);
}

.UserPill:hover {
    border-color: var(--color-border-light);
    background: var(--color-bg-elevated);
}

.Avatar {
    width: 28px;
    height: 28px;
    border-radius: 50%;
    background: linear-gradient(135deg, #FF4B3E 0%, #ff8266 100%);
    display: inline-flex;
    align-items: center;
    justify-content: center;
}

.Backdrop {
    position: fixed;
    inset: 0;
    z-index: 10;
}

.Menu {
    position: absolute;
    top: calc(100% + 8px);
    right: 0;
    min-width: 180px;
    padding: 6px;
    background: var(--color-bg-elevated);
    border: 1px solid var(--color-border-light);
    border-radius: 12px;
    box-shadow: 0 24px 64px -16px rgba(0,0,0,0.7);
    z-index: 20;
    animation: menuIn 160ms ease both;
}

@keyframes menuIn {
    from { opacity: 0; transform: translateY(-4px); }
    to   { opacity: 1; transform: translateY(0); }
}

.MenuItem {
    display: flex;
    align-items: center;
    gap: 10px;
    width: 100%;
    padding: 9px 12px;
    border: 0;
    background: transparent;
    color: var(--color-fg-primary);
    border-radius: 8px;
    font-size: 13px;
    font-family: inherit;
    text-align: left;
    cursor: pointer;
    transition: background 140ms ease, color 140ms ease;
}

.MenuItem:hover {
    background: var(--color-bg-hover);
}

.MenuItemDanger {
    color: var(--color-fg-secondary);
}

.MenuItemDanger:hover {
    background: rgba(193,37,61,0.10);
    color: var(--color-error);
}
```

## Modify: `pkg/client/ArtelUI/src/pages/home/HomePage.tsx`

1. Add import at top:
   ```tsx
   import Topbar from "@/components/Topbar/Topbar.tsx"
   ```

2. Replace the entire `<header className={cls.Header}>...</header>` block with:
   ```tsx
   <Topbar onLogout={handleLogout}/>
   ```

3. Add a hero section inside `<div className={cls.Root}>` right after `<Topbar>`:
   ```tsx
   <div className={cls.Hero}>
       <div className={cls.HeroTitles}>
           <div className={cls.Eyebrow}>Workspace</div>
           <h1 className={cls.HeroTitle}>Your vaults</h1>
           <p className={cls.HeroSub}>
               <b>{loading ? "…" : `${vaults.length} ${vaults.length === 1 ? "vault" : "vaults"}`}</b>
               {" · "}<span>all systems operational</span>
           </p>
       </div>
       <button className={cls.NewVaultBtn} onClick={openDialog}>
           <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
               <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
           </svg>
           New vault
       </button>
   </div>
   ```

4. Remove the `<button className={cls.Fab} ...>+</button>` line entirely.

## Modify: `pkg/client/ArtelUI/src/pages/home/HomePage.module.css`

Remove the `.Header`, `.Logo`, `.LogoutBtn`, `.Fab` blocks entirely.

Add these new blocks:

```css
.Hero {
    display: flex;
    align-items: flex-end;
    justify-content: space-between;
    gap: 24px;
    max-width: 1200px;
    margin: 0 auto;
    padding: 56px 32px 0;
    animation: fadeUp 500ms cubic-bezier(.2,.7,.2,1) both;
}

@keyframes fadeUp {
    from { opacity: 0; transform: translateY(8px); }
    to   { opacity: 1; transform: translateY(0); }
}

.HeroTitles {
    min-width: 0;
}

.Eyebrow {
    font-family: var(--font-mono);
    font-size: 11px;
    letter-spacing: 0.12em;
    color: var(--color-fg-tertiary);
    text-transform: uppercase;
    margin-bottom: 12px;
}

.HeroTitle {
    font-weight: 600;
    font-size: 42px;
    letter-spacing: -0.03em;
    line-height: 1.05;
    color: var(--color-fg-primary);
}

.HeroSub {
    font-size: 14px;
    color: var(--color-fg-secondary);
    margin-top: 10px;
    letter-spacing: -0.005em;
}

.HeroSub b {
    color: var(--color-fg-primary);
    font-weight: 600;
}

.NewVaultBtn {
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
    letter-spacing: -0.005em;
    cursor: pointer;
    white-space: nowrap;
    box-shadow: 0 10px 28px -12px var(--coral-glow);
    transition: filter 160ms ease, box-shadow 160ms ease;
    flex-shrink: 0;
}

.NewVaultBtn:hover {
    filter: brightness(1.06);
    box-shadow: 0 14px 32px -10px var(--coral-glow);
}

.NewVaultBtn:active {
    transform: translateY(1px);
}

@media (max-width: 720px) {
    .Hero {
        flex-direction: column;
        align-items: flex-start;
        gap: 18px;
        padding: 36px 18px 0;
    }
    .HeroTitle { font-size: 32px; }
}
```

Also update `.Content` padding to account for the hero:

```css
.Content {
    flex: 1;
    padding: 32px 32px 80px;
    max-width: 1200px;
    margin: 0 auto;
    width: 100%;
    overflow-y: auto;
}

@media (max-width: 720px) {
    .Content { padding: 24px 18px 60px; }
}
```

## Build check

```bash
cd pkg/client/ArtelUI && bun run build
```

Expected: zero errors.
