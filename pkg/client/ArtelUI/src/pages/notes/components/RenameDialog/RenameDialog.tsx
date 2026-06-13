import { useState } from "react"
import cls from "./RenameDialog.module.css"
import { useDialog, useDialogKeyboard } from "@/app/hooks/Dialog"

interface Props {
    currentPath: string
    onConfirm: (newPath: string) => Promise<void>
}

export default function RenameDialog({ currentPath, onConfirm }: Props) {
    const { CloseDialog } = useDialog()
    const [newPath, setNewPath] = useState(currentPath)
    const [loading, setLoading] = useState(false)

    async function handleConfirm() {
        if (!newPath.trim() || newPath === currentPath) return
        setLoading(true)
        try {
            await onConfirm(newPath.trim())
        } finally {
            setLoading(false)
            CloseDialog()
        }
    }

    const handleKeyDown = useDialogKeyboard(handleConfirm)

    return (
        <div className={cls.RenameContainer} role="dialog" aria-modal="true">
            <h2 className={cls.RenameTitle}>Rename / Move</h2>
            <p className={cls.RenameLabel}>New path</p>
            <input
                className={cls.RenameInput}
                type="text"
                value={newPath}
                onChange={e => setNewPath(e.target.value)}
                onKeyDown={handleKeyDown}
                disabled={loading}
                autoFocus
                spellCheck={false}
            />
            <div className={cls.RenameActions}>
                <button
                    className={cls.BtnCancel}
                    type="button"
                    onClick={CloseDialog}
                    disabled={loading}
                >
                    Cancel
                </button>
                <button
                    className={cls.BtnConfirm}
                    type="button"
                    onClick={handleConfirm}
                    disabled={loading || !newPath.trim() || newPath === currentPath}
                >
                    {loading ? "…" : "Move"}
                </button>
            </div>
        </div>
    )
}
