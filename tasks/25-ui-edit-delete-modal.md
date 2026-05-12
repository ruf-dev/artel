---
status: done
---

# Task 25 — UI: Edit / Delete Vault Modal

## Goal

Add an edit modal to `HomePage` that opens when the edit button on a `VaultCard` is clicked.
The modal shows the vault name (read-only display) and a danger zone with a "Delete" button
that calls `VaultsAPI.DeleteVault` and refreshes the list.

(No rename/description update — `UpdateVault` is not in the API yet.)

## Prerequisites

Tasks 20, 21, 22, 24 must be done first. Task 22 added `editVaultId` state and `onEdit` prop to
VaultCard. Task 24 added `.Modal`, `.ModalHead`, `.ModalTitle`, `.ModalClose`, `.BtnGhost`,
`.BtnPrimary`, `.Overlay`, `.ModalActions` CSS classes to `HomePage.module.css`.

## Modify: `pkg/client/ArtelUI/src/pages/home/HomePage.tsx`

### 1. Existing imports — ensure these are present:

```tsx
import {VaultsAPI} from "@/app/api/artel/vaults.pb.ts"
```

(It's already used in `handleCreate`, so it should be there.)

### 2. Add state for deleting:

After the `const [editVaultId, setEditVaultId] = useState<string | null>(null)` line, add:

```tsx
const [deleting, setDeleting] = useState(false)
```

### 3. Add a computed variable for the vault being edited:

Right before the `return (` statement, add:

```tsx
const editVault = editVaultId ? vaults.find(v => v.id === editVaultId) ?? null : null
```

### 4. Add delete handler function (alongside `handleCreate` and `handleLogout`):

```tsx
async function handleDelete() {
    if (!editVaultId) return
    setDeleting(true)
    try {
        await VaultsAPI.DeleteVault({id: editVaultId}, auth.getInitReq())
        setEditVaultId(null)
        void fetchVaults()
    } finally {
        setDeleting(false)
    }
}
```

### 5. Add edit modal JSX — place it right after the create modal `{dialogOpen && (...)}` block:

```tsx
{editVault && (
    <div className={cls.Overlay} onClick={() => setEditVaultId(null)}>
        <div className={cls.Modal} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true" aria-labelledby="editModalTitle">
            <div className={cls.ModalHead}>
                <h2 className={cls.ModalTitle} id="editModalTitle">Edit vault</h2>
                <button className={cls.ModalClose} type="button" onClick={() => setEditVaultId(null)} aria-label="Close">
                    <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round">
                        <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
                    </svg>
                </button>
            </div>
            <p className={cls.ModalSub}>Rename or delete this vault.</p>

            <div className={cls.FieldLabel} style={{marginBottom: 4}}>Vault name</div>
            <div className={cls.VaultNameDisplay}>{editVault.name}</div>

            <div className={cls.DangerZone}>
                <div className={cls.DangerZoneRow}>
                    <div>
                        <div className={cls.DangerZoneTitle}>Delete this vault</div>
                        <div className={cls.DangerZoneSub}>Permanent. Connection string stops working immediately.</div>
                    </div>
                    <button
                        className={cls.BtnDanger}
                        type="button"
                        onClick={handleDelete}
                        disabled={deleting}
                    >
                        {deleting ? "Deleting…" : "Delete"}
                    </button>
                </div>
            </div>

            <div className={cls.ModalActions}>
                <button className={cls.BtnGhost} type="button" onClick={() => setEditVaultId(null)}>
                    Close
                </button>
            </div>
        </div>
    </div>
)}
```

Also add `Escape` key handling — in the existing `handleKeyDown` function (which handles the create modal input), ensure the component also handles global Escape. Add a `useEffect`:

```tsx
useEffect(() => {
    function onKey(e: KeyboardEvent) {
        if (e.key === "Escape") {
            closeDialog()
            setEditVaultId(null)
        }
    }
    document.addEventListener("keydown", onKey)
    return () => document.removeEventListener("keydown", onKey)
}, [])
```

Place this `useEffect` alongside the other effects.

### 6. Import `KeyboardEvent` from react if not already (for the input handler type):

The existing file already has `import type {KeyboardEvent} from "react"` — leave it.

## Modify: `pkg/client/ArtelUI/src/pages/home/HomePage.module.css`

Add these blocks (append to the end of the file):

```css
.VaultNameDisplay {
    padding: 12px 14px;
    background: var(--color-bg-input);
    border: 1px solid var(--color-border-light);
    border-radius: 10px;
    color: var(--color-fg-primary);
    font-size: 14px;
    margin-bottom: 18px;
    word-break: break-word;
}

.DangerZone {
    margin-top: 8px;
    padding: 16px;
    background: rgba(193,37,61,0.06);
    border: 1px solid rgba(193,37,61,0.25);
    border-radius: 12px;
}

.DangerZoneRow {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 14px;
}

.DangerZoneTitle {
    font-size: 13px;
    font-weight: 600;
    color: var(--color-fg-primary);
    margin-bottom: 2px;
}

.DangerZoneSub {
    font-size: 12px;
    color: var(--color-fg-secondary);
    line-height: 1.45;
}

.BtnDanger {
    flex-shrink: 0;
    display: inline-flex;
    align-items: center;
    padding: 8px 14px;
    border-radius: 10px;
    border: 1px solid rgba(193,37,61,0.4);
    background: transparent;
    color: var(--color-error);
    font-size: 13px;
    font-weight: 500;
    font-family: inherit;
    cursor: pointer;
    transition: background 160ms ease, border-color 160ms ease;
}

.BtnDanger:hover:not(:disabled) {
    background: rgba(193,37,61,0.10);
    border-color: var(--color-error);
}

.BtnDanger:disabled {
    opacity: 0.4;
    cursor: not-allowed;
}
```

## Build check

```bash
cd pkg/client/ArtelUI && bun run build
```

Expected: zero errors. The edit button on each vault card should open a modal with vault name
display + delete danger zone.
