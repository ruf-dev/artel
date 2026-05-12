---
status: done
---

# Task 22 — UI: VaultCard Redesign

## Goal

Rebuild `VaultCard` to match the design: edit icon button, status chip (always Online since the API
doesn't expose status yet), connection string bar with styled copy button. Description is shown as
empty/greyed placeholder when absent from the API.

## Prerequisite

Tasks 20 and 21 must be done first (tokens and hero layout).

## Context

- `VaultItem` type (from generated proto): `{ id?: string; name?: string; dbUrl?: string }`
- There is no status, itemCount, description, or passProtected in the API — handle gracefully
- The `onEdit` prop is a callback so `HomePage` can open the edit modal (task 25)
- Copy uses `navigator.clipboard`

## Replace: `pkg/client/ArtelUI/src/pages/home/VaultCard.tsx`

```tsx
import {useState} from "react"
import cls from "@/pages/home/VaultCard.module.css"
import {VaultItem} from "@/app/api/artel/vaults.pb.ts"

interface Props {
    vault: VaultItem
    onEdit?: (id: string) => void
}

export default function VaultCard({vault, onEdit}: Props) {
    const [copied, setCopied] = useState(false)

    async function handleCopy() {
        const text = vault.dbUrl ?? ""
        if (navigator.clipboard?.writeText) {
            await navigator.clipboard.writeText(text).catch(() => {})
        } else {
            const ta = document.createElement("textarea")
            ta.value = text
            document.body.appendChild(ta)
            ta.select()
            try { document.execCommand("copy") } catch {}
            document.body.removeChild(ta)
        }
        setCopied(true)
        setTimeout(() => setCopied(false), 1500)
    }

    return (
        <article className={cls.Card}>
            <div className={cls.Head}>
                <h3 className={cls.Name}>{vault.name}</h3>
                {onEdit && (
                    <button
                        className={cls.IconBtn}
                        type="button"
                        onClick={() => onEdit(vault.id ?? "")}
                        title="Edit vault"
                        aria-label="Edit vault"
                    >
                        <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
                            <path d="M12 20h9"/><path d="M16.5 3.5a2.12 2.12 0 0 1 3 3L7 19l-4 1 1-4z"/>
                        </svg>
                    </button>
                )}
            </div>

            <div className={cls.Meta}>
                <span className={cls.Chip}>
                    <span className={cls.StatusDot}/>
                    <span>Online</span>
                </span>
            </div>

            <div className={cls.ConnBar}>
                <code className={cls.ConnString} title={vault.dbUrl}>{vault.dbUrl}</code>
                <button
                    className={`${cls.CopyBtn} ${copied ? cls.CopyBtnCopied : ""}`}
                    type="button"
                    onClick={handleCopy}
                >
                    <svg viewBox="0 0 24 24" width="12" height="12" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
                        <rect x="9" y="9" width="13" height="13" rx="2"/>
                        <path d="M5 15V5a2 2 0 0 1 2-2h10"/>
                    </svg>
                    <span>{copied ? "Copied" : "Copy"}</span>
                </button>
            </div>
        </article>
    )
}
```

## Replace: `pkg/client/ArtelUI/src/pages/home/VaultCard.module.css`

```css
.Card {
    position: relative;
    display: flex;
    flex-direction: column;
    padding: 22px 22px 18px;
    background: var(--color-bg-card);
    border: 1px solid var(--color-border);
    border-radius: 18px;
    transition: border-color 200ms ease, background 200ms ease;
    animation: fadeUp 500ms cubic-bezier(.2,.7,.2,1) both;
    min-width: 0;
}

@keyframes fadeUp {
    from { opacity: 0; transform: translateY(8px); }
    to   { opacity: 1; transform: translateY(0); }
}

.Card:hover {
    border-color: var(--color-border-light);
    background: #101010;
}

.Head {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 12px;
    margin-bottom: 14px;
}

.Name {
    font-weight: 600;
    font-size: 18px;
    letter-spacing: -0.015em;
    line-height: 1.25;
    color: var(--color-fg-primary);
    word-break: break-word;
}

.IconBtn {
    flex-shrink: 0;
    width: 32px;
    height: 32px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    background: transparent;
    border: 1px solid transparent;
    border-radius: 8px;
    color: var(--color-fg-tertiary);
    cursor: pointer;
    transition: all 140ms ease;
}

.IconBtn:hover {
    background: var(--color-bg-hover);
    border-color: var(--color-border-light);
    color: var(--color-fg-primary);
}

.Meta {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    margin-bottom: 16px;
}

.Chip {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 5px 10px;
    background: var(--color-bg-elevated);
    border: 1px solid var(--color-border);
    border-radius: 999px;
    font-size: 11px;
    font-family: var(--font-mono);
    letter-spacing: 0.02em;
    color: var(--color-fg-secondary);
}

.StatusDot {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background: var(--status-online);
    box-shadow: 0 0 8px var(--status-online-glow);
}

.ConnBar {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 10px 10px 14px;
    background: var(--color-bg-input);
    border: 1px solid var(--color-border);
    border-radius: 10px;
    margin-top: auto;
}

.ConnString {
    flex: 1;
    min-width: 0;
    font-family: var(--font-mono);
    font-size: 11.5px;
    color: var(--color-fg-secondary);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.CopyBtn {
    flex-shrink: 0;
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 6px 10px;
    background: var(--color-bg-elevated);
    border: 1px solid var(--color-border-light);
    border-radius: 7px;
    color: var(--color-fg-secondary);
    font-size: 11px;
    font-family: var(--font-mono);
    letter-spacing: 0.04em;
    cursor: pointer;
    text-transform: uppercase;
    transition: all 140ms ease;
}

.CopyBtn:hover {
    color: var(--color-fg-primary);
    border-color: var(--color-border-bright);
    background: var(--color-bg-hover-mid);
}

.CopyBtnCopied {
    color: var(--status-online) !important;
    border-color: var(--status-online-glow) !important;
}
```

## Modify: `pkg/client/ArtelUI/src/pages/home/HomePage.tsx`

The `VaultCard` now accepts an `onEdit` prop. Pass a handler from `HomePage`:

1. Add state for edit modal (just the vault id for now):
   ```tsx
   const [editVaultId, setEditVaultId] = useState<string | null>(null)
   ```

2. In the grid, pass `onEdit`:
   ```tsx
   <VaultCard key={v.id} vault={v} onEdit={id => setEditVaultId(id)}/>
   ```

(The full edit modal UI will be wired in task 25.)

## Build check

```bash
cd pkg/client/ArtelUI && bun run build
```

Expected: zero errors.
