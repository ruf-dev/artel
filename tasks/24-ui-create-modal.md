---
status: done
---

# Task 24 — UI: Create Modal Redesign

## Goal

Upgrade the create-vault dialog in `HomePage` to match the design: proper modal animations
(`modalIn`/`overlayIn`), close button in header, subtitle text, and an optional passphrase
field (stored locally but not sent to the API yet — the API only accepts `name`).

## Prerequisite

Tasks 20 and 21 must be done first.

## What to change

### `pkg/client/ArtelUI/src/pages/home/HomePage.tsx`

Replace the entire `{dialogOpen && (...)}` block with:

```tsx
{dialogOpen && (
    <div className={cls.Overlay} onClick={closeDialog}>
        <div className={cls.Modal} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true" aria-labelledby="createModalTitle">
            <div className={cls.ModalHead}>
                <h2 className={cls.ModalTitle} id="createModalTitle">New vault</h2>
                <button className={cls.ModalClose} type="button" onClick={closeDialog} aria-label="Close">
                    <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round">
                        <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
                    </svg>
                </button>
            </div>
            <p className={cls.ModalSub}>Give it a name — that's all we need to spin one up.</p>

            <label className={cls.Field}>
                <span className={cls.FieldLabel}>Vault name</span>
                <input
                    ref={inputRef}
                    className={cls.Input}
                    placeholder="e.g. Marketplace inventory"
                    value={vaultName}
                    onChange={e => setVaultName(e.target.value)}
                    onKeyDown={handleKeyDown}
                    disabled={creating}
                    maxLength={48}
                    autoComplete="off"
                />
            </label>

            <div className={cls.ModalActions}>
                <button className={cls.BtnGhost} type="button" onClick={closeDialog} disabled={creating}>
                    Cancel
                </button>
                <button
                    className={cls.BtnPrimary}
                    type="button"
                    onClick={handleCreate}
                    disabled={creating || !vaultName.trim()}
                >
                    {creating ? "Creating…" : "Create vault"}
                </button>
            </div>
        </div>
    </div>
)}
```

### `pkg/client/ArtelUI/src/pages/home/HomePage.module.css`

Replace the existing `.Overlay`, `.Dialog`, `.DialogTitle`, `.DialogInput`, `.DialogActions`,
`.CancelBtn`, `.CreateBtn` blocks with:

```css
@keyframes overlayIn {
    from { opacity: 0; }
    to   { opacity: 1; }
}

@keyframes modalIn {
    from { opacity: 0; transform: translateY(12px) scale(0.985); }
    to   { opacity: 1; transform: translateY(0) scale(1); }
}

.Overlay {
    position: fixed;
    inset: 0;
    z-index: 100;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 24px;
    background: rgba(0,0,0,0.6);
    backdrop-filter: blur(8px);
    -webkit-backdrop-filter: blur(8px);
    animation: overlayIn 200ms ease both;
}

.Modal {
    width: 100%;
    max-width: 460px;
    background: var(--color-bg-card);
    border: 1px solid var(--color-border-light);
    border-radius: 20px;
    padding: 28px 28px 24px;
    box-shadow: 0 32px 80px -20px rgba(0,0,0,0.8);
    animation: modalIn 280ms cubic-bezier(.2,.7,.2,1) both;
}

.ModalHead {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 16px;
    margin-bottom: 6px;
}

.ModalTitle {
    font-size: 22px;
    font-weight: 600;
    letter-spacing: -0.02em;
    color: var(--color-fg-primary);
}

.ModalClose {
    flex-shrink: 0;
    width: 32px;
    height: 32px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    background: transparent;
    border: 0;
    border-radius: 8px;
    color: var(--color-fg-tertiary);
    cursor: pointer;
    transition: all 140ms ease;
}

.ModalClose:hover {
    color: var(--color-fg-primary);
    background: var(--color-bg-hover);
}

.ModalSub {
    font-size: 13px;
    color: var(--color-fg-secondary);
    margin-bottom: 24px;
    line-height: 1.55;
}

.Field {
    display: block;
    margin-bottom: 18px;
}

.FieldLabel {
    display: block;
    font-size: 12px;
    font-weight: 500;
    color: var(--color-fg-secondary);
    margin-bottom: 8px;
    letter-spacing: -0.005em;
}

.Input {
    width: 100%;
    padding: 12px 14px;
    background: var(--color-bg-input);
    border: 1px solid var(--color-border-light);
    border-radius: 10px;
    color: var(--color-fg-primary);
    font-size: 14px;
    font-family: inherit;
    letter-spacing: -0.005em;
    transition: border-color 160ms ease, background 160ms ease;
    outline: none;
    box-sizing: border-box;
}

.Input::placeholder {
    color: var(--color-fg-disabled);
}

.Input:focus {
    border-color: var(--color-border-bright);
    background: #131313;
}

.ModalActions {
    display: flex;
    align-items: center;
    justify-content: flex-end;
    gap: 10px;
    margin-top: 24px;
}

.BtnGhost {
    display: inline-flex;
    align-items: center;
    padding: 10px 16px;
    border-radius: 10px;
    border: 1px solid var(--color-border-light);
    background: transparent;
    color: var(--color-fg-secondary);
    font-size: 13px;
    font-weight: 500;
    font-family: inherit;
    cursor: pointer;
    transition: all 160ms ease;
}

.BtnGhost:hover {
    color: var(--color-fg-primary);
    border-color: var(--color-border-bright);
}

.BtnPrimary {
    display: inline-flex;
    align-items: center;
    padding: 10px 18px;
    border-radius: 10px;
    border: 1px solid var(--coral);
    background: var(--coral);
    color: white;
    font-size: 13px;
    font-weight: 600;
    font-family: inherit;
    cursor: pointer;
    box-shadow: 0 10px 28px -12px var(--coral-glow);
    transition: filter 160ms ease, opacity 160ms ease;
}

.BtnPrimary:hover:not(:disabled) {
    filter: brightness(1.06);
}

.BtnPrimary:disabled {
    opacity: 0.4;
    cursor: not-allowed;
}
```

## Build check

```bash
cd pkg/client/ArtelUI && bun run build
```

Expected: zero errors. The create dialog should now animate in with `modalIn`.
